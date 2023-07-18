package transport

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/orderer"
)

type ServerHandle struct {
	Handler Handler
}

// DispatchSubmit identifies the channel and sender of the submit request and passes it
// to the underlying Handler
func (c *ServerHandle) DispatchSubmit(ctx context.Context, request *orderer.SubmitRequest) error {
	reqCtx, err := c.requestContext(ctx, request)
	if err != nil {
		return err
	}
	return c.Handler.OnSubmit(reqCtx.channel, reqCtx.sender, request)
}

// DispatchConsensus identifies the channel and sender of the step request and passes it
// to the underlying Handler
func (c *ServerHandle) DispatchConsensus(ctx context.Context, request *orderer.ConsensusRequest) error {
	reqCtx, err := c.requestContext(ctx, request)
	if err != nil {
		return err
	}
	return c.Handler.OnConsensus(reqCtx.channel, reqCtx.sender, request)
}

// requestContext identifies the sender and channel of the request and returns
// it wrapped in a requestContext
func (c *ServerHandle) requestContext(ctx context.Context, msg proto.Message) (*requestContext, error) {
	// channel := c.ChanExt.TargetChannel(msg)
	return &requestContext{}, nil
}
