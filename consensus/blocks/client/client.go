// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/cheggaaa/pb"
	cb "github.com/hyperledger/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/protoutil"
	"google.golang.org/grpc"
)

type Signer struct {
}

func (s *Signer) Sign(message []byte) ([]byte, error) {
	return []byte{}, nil
}
func (s *Signer) Serialize() ([]byte, error) {
	return []byte{}, nil
}

type broadcastClient struct {
	client    ab.AtomicBroadcast_BroadcastClient
	signer    *Signer
	channelID string
}

// newBroadcastClient creates a simple instance of the broadcastClient interface
func newBroadcastClient(client ab.AtomicBroadcast_BroadcastClient, channelID string, signer *Signer) *broadcastClient {
	return &broadcastClient{client: client, channelID: channelID, signer: signer}
}

func (s *broadcastClient) broadcast(transaction []byte) error {
	env, err := protoutil.CreateSignedEnvelope(cb.HeaderType_MESSAGE, s.channelID, s.signer, &cb.ConfigValue{Value: transaction}, 0, 0)
	if err != nil {
		panic(err)
	}
	return s.client.Send(env)
}

func (s *broadcastClient) getAck() error {
	msg, err := s.client.Recv()
	if err != nil {
		return err
	}
	if msg.Status != cb.Status_SUCCESS {
		return fmt.Errorf("got unexpected status: %v - %s", msg.Status, msg.Info)
	}
	return nil
}
func main() {
	serverAddr := "localhost:8080"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	clusterClient := ab.NewClusterClient(conn)
	opts := grpc.EmptyCallOption{}
	stepClient, err := clusterClient.Step(context.TODO(), opts)

	payload := &cb.Envelope{}
	submitReq := &ab.StepRequest_SubmitRequest{
		SubmitRequest: &ab.SubmitRequest{
			Channel:           "test",
			LastValidationSeq: 0,
			Payload:           payload,
		},
	}
	stepReq := &ab.StepRequest{
		Payload: submitReq,
	}
	err = stepClient.Send(stepReq)
	if err != nil {
		fmt.Println("Send failed:", err)
		return
	}

	time.Sleep(10 * time.Second)
}
func client_broadcast() {
	signer := &Signer{}

	var channelID string
	var serverAddr string
	var messages uint64
	var goroutines uint64
	var msgSize uint64
	var bar *pb.ProgressBar

	// flag.StringVar(&serverAddr, "server", fmt.Sprintf("%s:%d", conf.General.ListenAddress, conf.General.ListenPort), "The RPC server to connect to.")
	// flag.StringVar(&channelID, "channelID", "mychannel", "The channel ID to broadcast to.")
	// flag.Uint64Var(&messages, "messages", 1, "The number of messages to broadcast.")
	// flag.Uint64Var(&goroutines, "goroutines", 1, "The number of concurrent go routines to broadcast the messages on")
	// flag.Uint64Var(&msgSize, "size", 1024, "The size in bytes of the data section for the payload")
	// flag.Parse()

	serverAddr = "localhost:8080"
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	messages = 1000
	goroutines = 1
	msgsPerGo := messages / goroutines
	roundMsgs := msgsPerGo * goroutines
	if roundMsgs != messages {
		fmt.Println("Rounding messages to", roundMsgs)
	}
	bar = pb.New64(int64(roundMsgs))
	bar.ShowPercent = true
	bar.ShowSpeed = true
	bar = bar.Start()

	msgData := make([]byte, msgSize)

	var wg sync.WaitGroup
	wg.Add(int(goroutines))
	for i := uint64(0); i < goroutines; i++ {
		go func(i uint64, pb *pb.ProgressBar) {
			client, err := ab.NewAtomicBroadcastClient(conn).Broadcast(context.TODO())
			if err != nil {
				fmt.Println("Error connecting:", err)
				return
			}

			s := newBroadcastClient(client, channelID, signer)
			done := make(chan (struct{}))
			go func() {
				for i := uint64(0); i < msgsPerGo; i++ {
					err = s.getAck()
					if err == nil && bar != nil {
						bar.Increment()
					}
				}
				if err != nil {
					fmt.Printf("\n Error1: %v\n", err)
				}
				close(done)
			}()
			for i := uint64(0); i < msgsPerGo; i++ {
				if err := s.broadcast(msgData); err != nil {
					fmt.Printf("\n broadcast failed: %v\n", err)
					panic(err)
				}
			}
			<-done
			wg.Done()
			client.CloseSend()
		}(i, bar)
	}

	wg.Wait()
	bar.FinishPrint("----------------------broadcast message finish-------------------------------")
}
