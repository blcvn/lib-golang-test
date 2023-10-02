package deliveryservice

import (
	"time"

	cb "github.com/hyperledger/fabric-protos-go/common"

	"github.com/blcvn/lib-golang-test/consensus/peer/blocksprovider"
	"github.com/blcvn/lib-golang-test/consensus/peer/comm"
	"github.com/blcvn/lib-golang-test/consensus/peer/identity"
	"github.com/blcvn/lib-golang-test/consensus/peer/orderers"

	"github.com/hyperledger/fabric/bccsp"
)

// Config dictates the DeliveryService's properties,
// namely how it connects to an ordering service endpoint,
// how it verifies messages received from it,
// and how it disseminates the messages to other peers
type Config struct {
	IsStaticLeader bool
	// CryptoSvc performs cryptographic actions like message verification and signing
	// and identity validation.
	CryptoSvc blocksprovider.BlockVerifier
	// Gossip enables to enumerate peers in the channel, send a message to peers,
	// and add a block to the gossip state transfer layer.
	Gossip blocksprovider.GossipServiceAdapter
	// OrdererSource provides orderer endpoints, complete with TLS cert pools.
	OrdererSource *orderers.ConnectionSource
	// Signer is the identity used to sign requests.
	Signer identity.SignerSerializer
	// DeliverServiceConfig is the configuration object.
	DeliverServiceConfig *DeliverServiceConfig
	// ChannelConfig the initial channel config.
	ChannelConfig *cb.Config
	// CryptoProvider the crypto service provider.
	CryptoProvider bccsp.BCCSP
}

// DeliverServiceConfig is the struct that defines the deliverservice configuration.
type DeliverServiceConfig struct {
	// PeerTLSEnabled enables/disables Peer TLS.
	PeerTLSEnabled bool
	// BlockGossipEnabled enables block forwarding via gossip
	BlockGossipEnabled bool
	// ReConnectBackoffThreshold sets the delivery service maximal delay between consencutive retries.
	ReConnectBackoffThreshold time.Duration
	// ReconnectTotalTimeThreshold sets the total time the delivery service may spend in reconnection attempts
	// until its retry logic gives up and returns an error.
	ReconnectTotalTimeThreshold time.Duration
	// ConnectionTimeout sets the delivery service <-> ordering service node connection timeout
	ConnectionTimeout time.Duration
	// Keepalive option for deliveryservice
	KeepaliveOptions comm.KeepaliveOptions
	// SecOpts provides the TLS info for connections
	SecOpts comm.SecureOptions

	// OrdererEndpointOverrides is a map of orderer addresses which should be
	// re-mapped to a different orderer endpoint.
	OrdererEndpointOverrides map[string]*orderers.Endpoint
}
