package app

import (
	"github.com/blcvn/lib-golang-test/blocks/consensus"
	cb "github.com/blcvn/lib-golang-test/blocks/types/common"
)

// Receiver defines a sink for the ordered broadcast messages
type AppReceiver struct {
	consensus.Receiver
}

// Ordered should be invoked sequentially as messages are ordered
// Each batch in `messageBatches` will be wrapped into a block.
// `pending` indicates if there are still messages pending in the receiver.

func (s *AppReceiver) Ordered(msg *cb.Envelope) (messageBatches [][]*cb.Envelope, pending bool) {
	return [][]*cb.Envelope{}, false
}

// Cut returns the current batch and starts a new one
func (s *AppReceiver) Cut() []*cb.Envelope {
	return []*cb.Envelope{}
}
