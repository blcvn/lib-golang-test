package deliver

import (
	"github.com/hyperledger/fabric-protos-go/common"
	commonledger "github.com/hyperledger/fabric/common/ledger"
	"github.com/hyperledger/fabric/core/ledger"
)

type fileLedgerBlockStore struct {
	ledger.PeerLedger
}

func (flbs fileLedgerBlockStore) AddBlock(*common.Block) error {
	logger.Infof("fileLedgerBlockStore.AddBlock: Add block to ledger")
	return nil
}

func (flbs fileLedgerBlockStore) RetrieveBlocks(startBlockNumber uint64) (commonledger.ResultsIterator, error) {
	logger.Infof("fileLedgerBlockStore.RetrieveBlocks: Get block from %d ", startBlockNumber)

	return flbs.GetBlocksIterator(startBlockNumber)
}

func (flbs fileLedgerBlockStore) Shutdown() {
	logger.Infof("fileLedgerBlockStore.Shutdown: shutdown")

}
