package app

import "github.com/hyperledger/fabric/common/channelconfig"

type AppChannel struct {
	channelconfig.Channel
}

// HashingAlgorithm returns the default algorithm to be used when hashing
// such as computing block hashes, and CreationPolicy digests
func (s *AppChannel) HashingAlgorithm() func(input []byte) []byte {
	return func(input []byte) []byte {
		return []byte{}
	}
}

// BlockDataHashingStructureWidth returns the width to use when constructing the
// Merkle tree to compute the BlockData hash
func (s *AppChannel) BlockDataHashingStructureWidth() uint32 {
	return uint32(0)
}

// OrdererAddresses returns the list of valid orderer addresses to connect to to invoke Broadcast/Deliver
func (s *AppChannel) OrdererAddresses() []string {
	return []string{}
}

// Capabilities defines the capabilities for a channel
func (s *AppChannel) Capabilities() channelconfig.ChannelCapabilities {
	return &AppChannelCapability{}
}
