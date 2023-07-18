package transport

import (
	"fmt"
	"log"
	"net"

	ab "github.com/hyperledger/fabric-protos-go/orderer"
	"google.golang.org/grpc"
)

func StartGrpcServer(port int, cs *ClusterService) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", port))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	// set up our server options
	serverOpts := []grpc.ServerOption{}

	grpcServer := grpc.NewServer(serverOpts...)

	ab.RegisterClusterServer(grpcServer, cs)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("Start grpc server error : %s", err)
		return err
	}
	return nil
}
