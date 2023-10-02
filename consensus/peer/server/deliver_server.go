package server

import (
	"github.com/blcvn/lib-golang-test/consensus/peer/server/deliver"
	"github.com/blcvn/lib-golang-test/log/flogging"
	"github.com/hyperledger/fabric-protos-go/peer"
)

var logger = flogging.MustGetLogger("common.deliverevents")

type DeliverServer struct {
	DeliverHandler *deliver.Handler
}

// DeliverFiltered sends a stream of blocks to a client after commitment
func (s *DeliverServer) DeliverFiltered(srv peer.Deliver_DeliverFilteredServer) error {
	logger.Debugf("Starting new DeliverFiltered handler")
	defer dumpStacktraceOnPanic()
	// getting policy checker based on resources.Event_FilteredBlock resource name
	deliverServer := &deliver.Server{
		Receiver: srv,
		ResponseSender: &filteredBlockResponseSender{
			Deliver_DeliverFilteredServer: srv,
		},
	}
	return s.DeliverHandler.Handle(srv.Context(), deliverServer)
}

// Deliver sends a stream of blocks to a client after commitment
func (s *DeliverServer) Deliver(srv peer.Deliver_DeliverServer) (err error) {
	logger.Debugf("Starting new Deliver handler")
	defer dumpStacktraceOnPanic()
	// getting policy checker based on resources.Event_Block resource name
	deliverServer := &deliver.Server{
		Receiver: srv,
		ResponseSender: &blockResponseSender{
			Deliver_DeliverServer: srv,
		},
	}
	return s.DeliverHandler.Handle(srv.Context(), deliverServer)
}

// DeliverWithPrivateData sends a stream of blocks and pvtdata to a client after commitment
func (s *DeliverServer) DeliverWithPrivateData(srv peer.Deliver_DeliverWithPrivateDataServer) (err error) {
	logger.Debug("Starting new DeliverWithPrivateData handler")
	defer dumpStacktraceOnPanic()

	// getting policy checker based on resources.Event_Block resource name
	deliverServer := &deliver.Server{
		Receiver: srv,
		ResponseSender: &blockAndPrivateDataResponseSender{
			Deliver_DeliverWithPrivateDataServer: srv,
		},
	}
	err = s.DeliverHandler.Handle(srv.Context(), deliverServer)
	return err
}
