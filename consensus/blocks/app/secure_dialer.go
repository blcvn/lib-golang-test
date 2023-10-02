package app

import (
	"github.com/blcvn/lib-golang-test/consensus/blocks/consensus/common/cluster"
	"google.golang.org/grpc"
)

type AppSecureDialer struct {
}

func (s *AppSecureDialer) Dial(address string, verifyFunc cluster.RemoteVerifier) (*grpc.ClientConn, error) {
	return grpc.Dial(address)
}
