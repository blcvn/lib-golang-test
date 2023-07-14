package app

import (
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/msp"
)

// ChannelCapabilities defines the capabilities for a channel
type AppChannelCapability struct {
	channelconfig.ChannelCapabilities
}

func (s *AppChannelCapability) Supported() error {
	return nil
}

// MSPVersion specifies the version of the MSP this channel must understand, including the MSP types
// and MSP principal types.
func (s *AppChannelCapability) MSPVersion() msp.MSPVersion {
	return msp.MSPv1_0
}

// ConsensusTypeMigration return true if consensus-type migration is permitted in both orderer and peer.
func (s *AppChannelCapability) ConsensusTypeMigration() bool {
	return false
}

// OrgSpecificOrdererEndpoints return true if the channel config processing allows orderer orgs to specify their own endpoints
func (s *AppChannelCapability) OrgSpecificOrdererEndpoints() bool {
	return false
}
