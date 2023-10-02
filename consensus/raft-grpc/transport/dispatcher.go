package transport

import (
	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/pkg/errors"
)

// MessageReceiver receives messages
type MessageReceiver interface {
	// Consensus passes the given ConsensusRequest message to the MessageReceiver
	Consensus(req *orderer.ConsensusRequest, sender uint64) error

	// Submit passes the given SubmitRequest message to the MessageReceiver
	Submit(req *orderer.SubmitRequest, sender uint64) error
}

//go:generate mockery -dir . -name ReceiverGetter -case underscore -output mocks

// ReceiverGetter obtains instances of MessageReceiver given a channel ID
type ReceiverGetter interface {
	// ReceiverByChain returns the MessageReceiver if it exists, or nil if it doesn't
	ReceiverByChain(channelID string) MessageReceiver
}

// Dispatcher dispatches Submit and Step requests to the designated per chain instances
type Dispatcher struct {
	ChainSelector ReceiverGetter
}

// OnConsensus notifies the Dispatcher for a reception of a StepRequest from a given sender on a given channel
func (d *Dispatcher) OnConsensus(channel string, sender uint64, request *orderer.ConsensusRequest) error {
	receiver := d.ChainSelector.ReceiverByChain(channel)
	if receiver == nil {
		return errors.Errorf("channel %s doesn't exist", channel)
	}
	return receiver.Consensus(request, sender)
}

// OnSubmit notifies the Dispatcher for a reception of a SubmitRequest from a given sender on a given channel
func (d *Dispatcher) OnSubmit(channel string, sender uint64, request *orderer.SubmitRequest) error {
	receiver := d.ChainSelector.ReceiverByChain(channel)
	if receiver == nil {
		return errors.Errorf("channel %s doesn't exist", channel)
	}
	return receiver.Submit(request, sender)
}
