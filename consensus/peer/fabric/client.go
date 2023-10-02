package fabric

import (
	"context"

	"github.com/blcvn/lib-golang-test/consensus/peer/config"
	"github.com/blcvn/lib-golang-test/log/flogging"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"google.golang.org/grpc"
)

var (
	otherLogger = flogging.MustGetLogger("nhs.lib.fabric.other.go")
)

func NewGrpcClient(peer config.FabricConfigPeer) (*grpc.ClientConn, error) {
	// TLSRootCA := peer.TLSRootCertFile
	// TLSKeyFile := peer.TLSKeyFile
	// TLSCertFile := peer.TLSCertFile

	// caPEM, err := ioutil.ReadFile(TLSRootCA)
	// if err != nil {
	// 	otherLogger.Errorf("Error while read TLSRoorCA File: %s", err.Error())
	// 	return nil, err
	// }
	// certPool := x509.NewCertPool()
	// if !certPool.AppendCertsFromPEM(caPEM) {
	// 	return nil, fmt.Errorf("failed to add server CA's certificate")
	// }

	// clientCert, err := tls.LoadX509KeyPair(TLSCertFile, TLSKeyFile)
	// if err != nil {
	// 	return nil, err
	// }
	// tlsConfig := &tls.Config{
	// 	ServerName:   peer.PeerName,
	// 	RootCAs:      certPool,
	// 	Certificates: []tls.Certificate{clientCert},
	// }
	// tlsCredentials := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{
		grpc.WithInsecure(), //binhnt: enable insecure connection
		grpc.WithReturnConnectionError(),
		// grpc.WithBlock(),
		// grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(104857600)), // 100MB
		// grpc.WithTimeout(time.Duration(1) * time.Second),
	}
	// opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))

	ctx, cancel := context.WithTimeout(context.Background(), config.PeerGRPCTimeout)
	defer cancel()

	return grpc.DialContext(ctx, peer.PeerAddress, opts...)
}

func NewDeliverClient(peer config.FabricConfigPeer) (pb.DeliverClient, error) {
	conn, err := NewGrpcClient(peer)
	if err != nil {
		return nil, err
	}
	return pb.NewDeliverClient(conn), nil
}
