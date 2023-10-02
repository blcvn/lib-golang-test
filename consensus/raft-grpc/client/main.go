package client

import (
	"context"
	"crypto/x509"
	"fmt"

	"github.com/blcvn/lib-golang-test/consensus/raft-grpc/transport"
	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func main() {
	dialer := &transport.SDialer{}
	endpoint := "127.0.0.1:8001"
	conn, err := dialer.Dial(endpoint, remoteVerifierFn(endpoint, []byte(endpoint)))
	if err != nil {
		fmt.Printf("dialer.Dial failed %s \n", err)
		return
	}
	probeConnection := func(conn *grpc.ClientConn) error {
		connState := conn.GetState()
		if connState == connectivity.Connecting {
			return errors.Errorf("connection to is in state %s \n", connState)
		}
		return nil
	}

	clusterClient := orderer.NewClusterClient(conn)

	getStream := func(ctx context.Context) (transport.StepClientStream, error) {
		fmt.Printf("getStream: call clusterClient.Step \n")
		stream, err := clusterClient.Step(ctx)
		if err != nil {
			fmt.Printf("getStream: clusterClient.Step failed %s", err)
			return nil, err
		}
		stepClientStream := &transport.CommClientStream{
			StepClient: stream,
		}
		return stepClientStream, nil
	}
	if err := probeConnection(conn); err != nil {
		fmt.Printf("probeConnection failed %s \n", err)
		return
	}
	ctx, cancel := context.WithCancel(context.TODO())
	stream, err := getStream(ctx)
	if err != nil {
		fmt.Printf("getStream failed %s \n", err)
		cancel()
		return
	}
	err = stream.Auth()
	if err != nil {
		fmt.Printf("ream.Auth failed %s \n", err)
		return
	}
}
func remoteVerifierFn(endpoint string, certificate []byte) transport.RemoteVerifier {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		return nil
	}
}
