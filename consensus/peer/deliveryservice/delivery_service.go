package deliveryservice

import (
	"github.com/blcvn/lib-golang-test/consensus/peer/blocksprovider"
)

type DeliverService interface {
	// StartDeliverForChannel dynamically starts delivery of new blocks from ordering service
	// to channel peers.
	// When the delivery finishes, the finalizer func is called
	StartDeliverForChannel(chainID string, ledgerInfo blocksprovider.LedgerInfo, finalizer func()) error

	// StopDeliverForChannel dynamically stops delivery of new blocks from ordering service
	// to channel peers. StartDeliverForChannel can be called again, and delivery will resume.
	StopDeliverForChannel() error

	// Stop terminates delivery service and closes the connection. Marks the service as stopped, meaning that
	// StartDeliverForChannel cannot be called again.
	Stop()
}

// BlockDeliverer communicates with orderers to obtain new blocks and send them to the committer service, for a
// specific channel. It can be implemented using different protocols depending on the ordering service consensus type,
// e.g CFT (etcdraft) or BFT (SmartBFT).
type BlockDeliverer interface {
	Stop()
	DeliverBlocks()
}

func NewDeliverService(conf *Config) DeliverService {
	ds := &deliverServiceImpl{
		conf: conf,
	}
	return ds
}
