package app

import (
	"time"

	ab "github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/common/channelconfig"
)

// Orderer stores the common shared orderer config
type AppOrderer struct {
	channelconfig.Orderer
}

// ConsensusType returns the configured consensus type
func (s *AppOrderer) ConsensusType() string {
	return ""
}

// ConsensusMetadata returns the metadata associated with the consensus type.
func (s *AppOrderer) ConsensusMetadata() []byte {
	return []byte{}
}

// ConsensusState returns the consensus-type state.
func (s *AppOrderer) ConsensusState() ab.ConsensusType_State {
	return ab.ConsensusType_STATE_NORMAL
}

// BatchSize returns the maximum number of messages to include in a block
func (s *AppOrderer) BatchSize() *ab.BatchSize {
	return &ab.BatchSize{}
}

// BatchTimeout returns the amount of time to wait before creating a batch
func (s *AppOrderer) BatchTimeout() time.Duration {
	return time.Second
}

// MaxChannelsCount returns the maximum count of channels to allow for an ordering network
func (s *AppOrderer) MaxChannelsCount() uint64 {
	return uint64(0)
}

// KafkaBrokers returns the addresses (IP:port notation) of a set of "bootstrap"
// Kafka brokers, i.e. this is not necessarily the entire set of Kafka brokers
// used for ordering
func (s *AppOrderer) KafkaBrokers() []string {
	return []string{}
}

// Organizations returns the organizations for the ordering service
func (s *AppOrderer) Organizations() map[string]channelconfig.OrdererOrg {
	mapp := map[string]channelconfig.OrdererOrg{}
	return mapp
}

// Capabilities defines the capabilities for the orderer portion of a channel
func (s *AppOrderer) Capabilities() channelconfig.OrdererCapabilities {
	return &AppOrdererCapabilities{}
}
