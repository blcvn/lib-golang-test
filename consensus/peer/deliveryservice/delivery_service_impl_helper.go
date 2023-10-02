package deliveryservice

import (
	"fmt"
	"time"

	"github.com/blcvn/lib-golang-test/consensus/peer/blocksprovider"

	"github.com/blcvn/lib-golang-test/consensus/peer/comm"
	"github.com/blcvn/lib-golang-test/log/flogging"
	"github.com/hyperledger/fabric/common/util"
)

func (d *deliverServiceImpl) createBlockDelivererCFT(chainID string, ledgerInfo blocksprovider.LedgerInfo) (*blocksprovider.Deliverer, error) {
	dc := &blocksprovider.Deliverer{
		ChannelID:     chainID,
		Gossip:        d.conf.Gossip,
		Ledger:        ledgerInfo,
		BlockVerifier: d.conf.CryptoSvc,
		Dialer: DialerAdapter{
			ClientConfig: comm.ClientConfig{
				DialTimeout: d.conf.DeliverServiceConfig.ConnectionTimeout,
				KaOpts:      d.conf.DeliverServiceConfig.KeepaliveOptions,
				SecOpts:     d.conf.DeliverServiceConfig.SecOpts,
			},
		},
		Orderers:            d.conf.OrdererSource,
		DoneC:               make(chan struct{}),
		Signer:              d.conf.Signer,
		DeliverStreamer:     DeliverAdapter{},
		Logger:              flogging.MustGetLogger("peer.blocksprovider").With("channel", chainID),
		MaxRetryDelay:       d.conf.DeliverServiceConfig.ReConnectBackoffThreshold,
		MaxRetryDuration:    d.conf.DeliverServiceConfig.ReconnectTotalTimeThreshold,
		BlockGossipDisabled: !d.conf.DeliverServiceConfig.BlockGossipEnabled,
		InitialRetryDelay:   100 * time.Millisecond,
		YieldLeadership:     !d.conf.IsStaticLeader,
	}

	if d.conf.DeliverServiceConfig.SecOpts.RequireClientCert {
		cert, err := d.conf.DeliverServiceConfig.SecOpts.ClientCertificate()
		if err != nil {
			return nil, fmt.Errorf("failed to access client TLS configuration: %w", err)
		}
		dc.TLSCertHash = util.ComputeSHA256(cert.Certificate[0])
	}
	return dc, nil
}

func (d *deliverServiceImpl) createBlockDelivererBFT(chainID string, ledgerInfo blocksprovider.LedgerInfo) (*blocksprovider.Deliverer, error) {
	// TODO create a BFT BlockDeliverer
	logger.Warning("Consensus type `BFT` BlockDeliverer not supported yet, creating a CFT one")
	return d.createBlockDelivererCFT(chainID, ledgerInfo)
}
