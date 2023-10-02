package deliver

import (
	"github.com/hyperledger/fabric/common/ledger/blockledger"
	"github.com/hyperledger/fabric/common/policies"
)

type Channel struct {
	// Ledger     ledger.PeerLedger
	BlockLedger blockledger.ReadWriter
	PManager    policies.Manager
}

// Sequence returns the current config sequence number of the channel.
func (c *Channel) Sequence() uint64 {
	return uint64(0)
}

// PolicyManager returns the policies.Manager for the channel that reflects the
// current channel configuration. Users should not memoize references to this object.
func (c *Channel) PolicyManager() policies.Manager {
	return c.PManager
}
func (c *Channel) Reader() blockledger.Reader {
	// return fileledger.NewFileLedger(fileLedgerBlockStore{c.Ledger})
	// return fileledger.NewFileLedger(c.BlockStore)
	return c.BlockLedger
}

func (c *Channel) Errored() <-chan struct{} {
	// If this is ever updated to return a real channel, the error message
	// in deliver.go around this channel closing should be updated.
	return nil
}
