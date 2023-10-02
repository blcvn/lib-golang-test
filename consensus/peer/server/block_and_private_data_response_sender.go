package server

import (
	"github.com/blcvn/lib-golang-test/consensus/peer/server/deliver"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/protoutil"
)

type Chain interface {
	deliver.Chain
	// Ledger() ledger.PeerLedger
}

// blockResponseSender structure used to send block responses
type blockAndPrivateDataResponseSender struct {
	peer.Deliver_DeliverWithPrivateDataServer
	IdentityDeserializerManager
}

// SendStatusResponse generates status reply proto message
func (bprs *blockAndPrivateDataResponseSender) SendStatusResponse(status common.Status) error {
	reply := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_Status{Status: status},
	}
	return bprs.Send(reply)
}

// SendBlockResponse gets private data and generates deliver response with both block and private data
func (bprs *blockAndPrivateDataResponseSender) SendBlockResponse(
	block *common.Block,
	channelID string,
	chain deliver.Chain,
	signedData *protoutil.SignedData,
) error {
	pvtData, err := bprs.getPrivateData(block, chain, channelID, signedData)
	if err != nil {
		return err
	}

	blockAndPvtData := &peer.BlockAndPrivateData{
		Block:          block,
		PrivateDataMap: pvtData,
	}
	response := &peer.DeliverResponse{
		Type: &peer.DeliverResponse_BlockAndPrivateData{BlockAndPrivateData: blockAndPvtData},
	}
	return bprs.Send(response)
}

func (bprs *blockAndPrivateDataResponseSender) DataType() string {
	return "block_and_pvtdata"
}

// getPrivateData returns private data for the block
func (bprs *blockAndPrivateDataResponseSender) getPrivateData(
	block *common.Block,
	chain deliver.Chain,
	channelID string,
	signedData *protoutil.SignedData,
) (map[uint64]*rwset.TxPvtReadWriteSet, error) {
	// channel, ok := chain.(Chain)
	// if !ok {
	// 	return nil, errors.New("wrong chain type")
	// }

	// pvtData, err := channel.Ledger().GetPvtDataByNum(block.Header.Number, nil)
	// if err != nil {
	// 	logger.Errorf("Error getting private data by block number %d on channel %s", block.Header.Number, channelID)
	// 	return nil, errors.Wrapf(err, "error getting private data by block number %d", block.Header.Number)
	// }

	seqs2Namespaces := aggregatedCollections(make(map[seqAndDataModel]map[string][]*rwset.CollectionPvtReadWriteSet))

	// configHistoryRetriever, err := channel.Ledger().GetConfigHistoryRetriever()
	// if err != nil {
	// 	return nil, err
	// }

	// check policy for each collection and add the collection if passing the policy requirement
	// for _, item := range pvtData {
	// 	logger.Debugf("Got private data for block number %d, tx sequence %d", block.Header.Number, item.SeqInBlock)
	// 	if item.WriteSet == nil {
	// 		continue
	// 	}
	// 	for _, ns := range item.WriteSet.NsPvtRwset {
	// 		for _, col := range ns.CollectionPvtRwset {
	// 			logger.Debugf("Checking policy for namespace %s, collection %s", ns.Namespace, col.CollectionName)

	// 		}
	// 	}
	// }

	return seqs2Namespaces.asPrivateDataMap(), nil
}
