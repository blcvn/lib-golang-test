package deliveryservice

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"google.golang.org/grpc"
)

type DeliverAdapter struct{}

func (DeliverAdapter) Deliver(ctx context.Context, clientConn *grpc.ClientConn) (orderer.AtomicBroadcast_DeliverClient, error) {
	return orderer.NewAtomicBroadcastClient(clientConn).Deliver(ctx)
}
