package transport

import (
	"context"
	"fmt"
	"io"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/common/util"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Dispatcher dispatches requests
type ServerDispatcher interface {
	DispatchSubmit(ctx context.Context, request *orderer.SubmitRequest) error
	DispatchConsensus(ctx context.Context, request *orderer.ConsensusRequest) error
}

// Handler handles Step() and Submit() requests and returns a corresponding response
type Handler interface {
	OnConsensus(channel string, sender uint64, req *orderer.ConsensusRequest) error
	OnSubmit(channel string, sender uint64, req *orderer.SubmitRequest) error
}

type requestContext struct {
	channel string
	sender  uint64
}

// StepStream defines the gRPC stream for sending
// transactions, and receiving corresponding responses
type StepStream interface {
	Send(response *orderer.StepResponse) error
	Recv() (*orderer.StepRequest, error)
	grpc.ServerStream
}

// Implement orderer.ClusterServer
type ClusterService struct {
	Dispatcher ServerDispatcher
}

// Step passes an implementation-specific message to another cluster member.
func (s *ClusterService) Step(stream orderer.Cluster_StepServer) error {
	fmt.Printf("ClusterService.Step: new stream \n")
	addr := util.ExtractRemoteAddress(stream.Context())
	for {
		err := s.handleMessage(stream, addr)
		if err == io.EOF {
			fmt.Printf("%s disconnected", addr)
			return nil
		}
		if err != nil {
			return err
		}
		// Else, no error occurred, so we continue to the next iteration
	}
}

func (s *ClusterService) handleMessage(stream orderer.Cluster_StepServer, addr string) error {
	fmt.Printf("ClusterService.handleMessage: message from %s \n", addr)

	request, err := stream.Recv()
	if err == io.EOF {
		fmt.Printf("ClusterService.handleMessage: error %s \n", err)

		return err
	}
	if err != nil {
		fmt.Printf("ClusterService.handleMessage: %s failed: %v \n", addr, err)
		return err
	}

	fmt.Printf("ClusterService.handleMessage: request %+v \n", request)

	if submitReq := request.GetSubmitRequest(); submitReq != nil {
		fmt.Printf("ClusterService.handleMessage: Receive Submit message from %s: %+v", addr, request)
		return s.handleSubmit(submitReq, stream, addr)
	} else if consensusReq := request.GetConsensusRequest(); consensusReq != nil {
		fmt.Printf("ClusterService.handleMessage: Receive Consensus message from %s: %+v", addr, request)

		return s.Dispatcher.DispatchConsensus(stream.Context(), request.GetConsensusRequest())
	}

	return errors.Errorf("message is neither a Submit nor a Consensus request")
}

func (s *ClusterService) handleSubmit(request *orderer.SubmitRequest, stream StepStream, addr string) error {
	err := s.Dispatcher.DispatchSubmit(stream.Context(), request)
	if err != nil {
		fmt.Printf("Handling of Submit() from %s failed: %v", addr, err)
		return err
	}
	return err
}
