package deliveryservice

import (
	"github.com/blcvn/lib-golang-test/consensus/peer/comm"
	"google.golang.org/grpc"
)

type DialerAdapter struct {
	ClientConfig comm.ClientConfig
}

func (da DialerAdapter) Dial(address string, rootCerts [][]byte) (*grpc.ClientConn, error) {
	cc := da.ClientConfig
	cc.SecOpts.ServerRootCAs = rootCerts
	return cc.Dial(address)
}
