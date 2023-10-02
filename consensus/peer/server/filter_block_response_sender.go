package server

import (
	"github.com/blcvn/lib-golang-test/consensus/peer/server/deliver"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
)

// filteredBlockResponseSender structure used to send filtered block responses
type filteredBlockResponseSender struct {
	peer.Deliver_DeliverFilteredServer
}

// SendStatusResponse generates status reply proto message
func (fbrs *filteredBlockResponseSender) SendStatusResponse(status common.Status) error {
	response := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_Status{Status: status},
	}
	return fbrs.Send(response)
}

// IsFiltered is a marker method which indicates that this response sender
// sends filtered blocks.
func (fbrs *filteredBlockResponseSender) IsFiltered() bool {
	return true
}

// SendBlockResponse generates deliver response with filtered block message
func (fbrs *filteredBlockResponseSender) SendBlockResponse(
	block *common.Block,
	channelID string,
	chain deliver.Chain,
	signedData *protoutil.SignedData,
) error {
	// Generates filtered block response
	b := blockEvent(*block)
	filteredBlock, err := b.toFilteredBlock()
	if err != nil {
		logger.Warningf("Failed to generate filtered block due to: %s", err)
		return fbrs.SendStatusResponse(common.Status_BAD_REQUEST)
	}
	response := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_FilteredBlock{FilteredBlock: filteredBlock},
	}
	return fbrs.Send(response)
}

func (fbrs *filteredBlockResponseSender) DataType() string {
	return "filtered_block"
}
