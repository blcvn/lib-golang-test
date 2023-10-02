package server

import (
	"github.com/blcvn/lib-golang-test/consensus/peer/server/deliver"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
)

// blockResponseSender structure used to send block responses
type blockResponseSender struct {
	peer.Deliver_DeliverServer
}

// SendStatusResponse generates status reply proto message
func (brs *blockResponseSender) SendStatusResponse(status common.Status) error {
	reply := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_Status{Status: status},
	}
	return brs.Send(reply)
}

// SendBlockResponse generates deliver response with block message.
func (brs *blockResponseSender) SendBlockResponse(
	block *common.Block,
	channelID string,
	chain deliver.Chain,
	signedData *protoutil.SignedData,
) error {
	// Generates filtered block response
	response := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_Block{Block: block},
	}
	return brs.Send(response)
}

func (brs *blockResponseSender) DataType() string {
	return "block"
}
