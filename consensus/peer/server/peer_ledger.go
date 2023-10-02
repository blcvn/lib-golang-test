package server

import (
	"fmt"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	commonledger "github.com/hyperledger/fabric/common/ledger"
	"github.com/hyperledger/fabric/common/ledger/blkstorage"
	"github.com/hyperledger/fabric/common/ledger/blockledger"
	"github.com/hyperledger/fabric/core/ledger"
)

//binhnt: Simulate peer ledger

type PeerLedger struct {
	BlockStore blkstorage.BlockStore
	it         blockledger.Iterator
	blockHeigh uint64
}

func (s *PeerLedger) Next() (commonledger.QueryResult, error) {
	block, status := s.it.Next()
	if status != common.Status_SUCCESS {
		return nil, fmt.Errorf("CANNOT FIND DATA %s ", status.String())
	}
	return block, nil
}

// GetBlockchainInfo returns basic info about blockchain
func (s *PeerLedger) GetBlockchainInfo() (*common.BlockchainInfo, error) {
	return &common.BlockchainInfo{}, nil
}

// GetBlockByNumber returns block at a given height
// blockNumber of  math.MaxUint64 will return last block
func (s *PeerLedger) GetBlockByNumber(blockNumber uint64) (*common.Block, error) {
	return &common.Block{}, nil
}

// GetBlocksIterator returns an iterator that starts from `startBlockNumber`(inclusive).
// The iterator is a blocking iterator i.e., it blocks till the next block gets available in the ledger
// ResultsIterator contains type BlockHolder
func (s *PeerLedger) GetBlocksIterator(c uint64) (commonledger.ResultsIterator, error) {
	// logger.Infof("PeerLedger.GetBlocksIterator: start => Return list of block start from %d ", c)
	// start := &orderer.SeekPosition{
	// 	Type: &orderer.SeekPosition_Specified{
	// 		Specified: &orderer.SeekSpecified{
	// 			Number: c,
	// 		},
	// 	},
	// }
	// it, h := s.ledger.Iterator(start)
	// s.it = it
	// s.blockHeigh = h
	return s, nil
}

// Close closes the ledger
func (s *PeerLedger) Close() {
	logger.Infof("PeerLedger.Close: start")

}
func (s *PeerLedger) GetTransactionByID(txID string) (*peer.ProcessedTransaction, error) {
	logger.Infof("PeerLedger.GetTransactionByID: start")

	return &peer.ProcessedTransaction{}, nil
}

// GetBlockByHash returns a block given it's hash
func (s *PeerLedger) GetBlockByHash(blockHash []byte) (*common.Block, error) {
	logger.Infof("PeerLedger.GetBlockByHash: start")

	return &common.Block{}, nil
}

// GetBlockByTxID returns a block which contains a transaction
func (s *PeerLedger) GetBlockByTxID(txID string) (*common.Block, error) {
	logger.Infof("PeerLedger.GetBlockByTxID: start")

	return &common.Block{}, nil
}

// GetTxValidationCodeByTxID returns reason code of transaction validation
func (s *PeerLedger) GetTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	logger.Infof("PeerLedger.GetTxValidationCodeByTxID: start")

	return peer.TxValidationCode_BAD_CHANNEL_HEADER, nil
}

// NewTxSimulator gives handle to a transaction simulator.
// A client can obtain more than one 'TxSimulator's for parallel execution.
// Any snapshoting/synchronization should be performed at the implementation level if required
func (s *PeerLedger) NewTxSimulator(txid string) (ledger.TxSimulator, error) {
	logger.Infof("PeerLedger.NewTxSimulator: start")

	return nil, nil
}

// NewQueryExecutor gives handle to a query executor.
// A client can obtain more than one 'QueryExecutor's for parallel execution.
// Any synchronization should be performed at the implementation level if required
func (s *PeerLedger) NewQueryExecutor() (ledger.QueryExecutor, error) {
	logger.Infof("PeerLedger.NewQueryExecutor: start")

	return nil, nil
}

// NewHistoryQueryExecutor gives handle to a history query executor.
// A client can obtain more than one 'HistoryQueryExecutor's for parallel execution.
// Any synchronization should be performed at the implementation level if required
func (s *PeerLedger) NewHistoryQueryExecutor() (ledger.HistoryQueryExecutor, error) {
	logger.Infof("PeerLedger.NewHistoryQueryExecutor: start")

	return nil, nil
}

// GetPvtDataAndBlockByNum returns the block and the corresponding pvt data.
// The pvt data is filtered by the list of 'ns/collections' supplied
// A nil filter does not filter any results and causes retrieving all the pvt data for the given blockNum
func (s *PeerLedger) GetPvtDataAndBlockByNum(blockNum uint64, filter ledger.PvtNsCollFilter) (*ledger.BlockAndPvtData, error) {
	logger.Infof("PeerLedger.GetPvtDataAndBlockByNum: start")

	return &ledger.BlockAndPvtData{}, nil
}

// GetPvtDataByNum returns only the pvt data  corresponding to the given block number
// The pvt data is filtered by the list of 'ns/collections' supplied in the filter
// A nil filter does not filter any results and causes retrieving all the pvt data for the given blockNum
func (s *PeerLedger) GetPvtDataByNum(blockNum uint64, filter ledger.PvtNsCollFilter) ([]*ledger.TxPvtData, error) {
	logger.Infof("PeerLedger.GetPvtDataByNum: start")

	return nil, nil
}

// CommitLegacy commits the block and the corresponding pvt data in an atomic operation following the v14 validation/commit path
// TODO: add a new Commit() path that replaces CommitLegacy() for the validation refactor described in FAB-12221
func (s *PeerLedger) CommitLegacy(blockAndPvtdata *ledger.BlockAndPvtData, commitOpts *ledger.CommitOptions) error {
	logger.Infof("PeerLedger.CommitLegacy: start")

	return nil
}

// GetConfigHistoryRetriever returns the ConfigHistoryRetriever
func (s *PeerLedger) GetConfigHistoryRetriever() (ledger.ConfigHistoryRetriever, error) {
	return nil, nil
}

// CommitPvtDataOfOldBlocks commits the private data corresponding to already committed block
// If hashes for some of the private data supplied in this function does not match
// the corresponding hash present in the block, the unmatched private data is not
// committed and instead the mismatch inforation is returned back
func (s *PeerLedger) CommitPvtDataOfOldBlocks(reconciledPvtdata []*ledger.ReconciledPvtdata) ([]*ledger.PvtdataHashMismatch, error) {
	return nil, nil
}

// GetMissingPvtDataTracker return the MissingPvtDataTracker
func (s *PeerLedger) GetMissingPvtDataTracker() (ledger.MissingPvtDataTracker, error) {
	return nil, nil
}

// DoesPvtDataInfoExist returns true when
// (1) the ledger has pvtdata associated with the given block number (or)
// (2) a few or all pvtdata associated with the given block number is missing but the
//
//	missing info is recorded in the ledger (or)
//
// (3) the block is committed and does not contain any pvtData.
func (s *PeerLedger) DoesPvtDataInfoExist(blockNum uint64) (bool, error) {
	return false, nil
}
