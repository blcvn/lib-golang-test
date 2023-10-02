package deliveryservice

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/blcvn/lib-golang-test/consensus/peer/blocksprovider"
	"github.com/blcvn/lib-golang-test/log/flogging"
	"github.com/hyperledger/fabric/common/channelconfig"
)

var logger = flogging.MustGetLogger("deliveryClient")

// deliverServiceImpl the implementation of the delivery service
// maintains connection to the ordering service and maps of
// blocks providers
type deliverServiceImpl struct {
	conf           *Config
	channelID      string
	blockDeliverer BlockDeliverer
	lock           sync.Mutex
	stopping       bool
}

// StartDeliverForChannel starts blocks delivery for channel
// initializes the grpc stream for given chainID, creates blocks provider instance
// that spawns in go routine to read new blocks starting from the position provided by ledger
// info instance.
func (d *deliverServiceImpl) StartDeliverForChannel(chainID string, ledgerInfo blocksprovider.LedgerInfo, finalizer func()) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.stopping {
		errMsg := fmt.Sprintf("block deliverer for channel `%s` is stopping", chainID)
		logger.Errorf("Delivery service: %s", errMsg)
		return errors.New(errMsg)
	}

	if d.blockDeliverer != nil {
		errMsg := fmt.Sprintf("block deliverer for channel `%s` already exists", chainID)
		logger.Errorf("Delivery service: %s", errMsg)
		return errors.New(errMsg)
	}

	// TODO save the initial bundle in the block deliverer in order to maintain a stand alone BlockVerifier that gets updated
	// immediately after a config block is pulled and verified.
	bundle, err := channelconfig.NewBundle(chainID, d.conf.ChannelConfig, d.conf.CryptoProvider)
	if err != nil {
		return errors.WithMessagef(err, "failed to create block deliverer for channel `%s`", chainID)
	}
	oc, ok := bundle.OrdererConfig()
	if !ok {
		// This should never happen because it is checked in peer.createChannel()
		return errors.Errorf("failed to create block deliverer for channel `%s`, missing OrdererConfig", chainID)
	}

	switch ct := oc.ConsensusType(); ct {
	case "etcdraft":
		d.blockDeliverer, err = d.createBlockDelivererCFT(chainID, ledgerInfo)
	case "BFT":
		d.blockDeliverer, err = d.createBlockDelivererBFT(chainID, ledgerInfo)
	default:
		err = errors.Errorf("unexpected consensus type: `%s`", ct)
	}

	if err != nil {
		return err
	}

	if !d.conf.DeliverServiceConfig.BlockGossipEnabled {
		logger.Infow("This peer will retrieve blocks from ordering service (will not disseminate them to other peers in the organization)", "channel", chainID)
	} else {
		logger.Infow("This peer will retrieve blocks from ordering service and disseminate to other peers in the organization", "channel", chainID)
	}

	d.channelID = chainID

	go func() {
		d.blockDeliverer.DeliverBlocks()
		finalizer()
	}()
	return nil
}

// StopDeliverForChannel stops blocks delivery for channel by stopping channel block provider
func (d *deliverServiceImpl) StopDeliverForChannel() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.stopping {
		errMsg := fmt.Sprintf("block deliverer for channel `%s` is already stopped", d.channelID)
		logger.Errorf("Delivery service: %s", errMsg)
		return errors.New(errMsg)
	}

	if d.blockDeliverer == nil {
		errMsg := fmt.Sprintf("block deliverer for channel `%s` is <nil>, can't stop delivery", d.channelID)
		logger.Errorf("Delivery service: %s", errMsg)
		return errors.New(errMsg)
	}
	d.blockDeliverer.Stop()
	d.blockDeliverer = nil

	logger.Debugf("This peer will stop passing blocks from orderer service to other peers on channel: %s", d.channelID)
	return nil
}

// Stop all service and release resources
func (d *deliverServiceImpl) Stop() {
	d.lock.Lock()
	defer d.lock.Unlock()
	// Marking flag to indicate the shutdown of the delivery service
	d.stopping = true

	if d.blockDeliverer != nil {
		d.blockDeliverer.Stop()
		d.blockDeliverer = nil
	}
}
