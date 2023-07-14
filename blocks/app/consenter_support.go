package app

import (
	cb "github.com/blcvn/lib-golang-test/blocks/types/common"

	"github.com/blcvn/lib-golang-test/blocks/consensus"
	"github.com/hyperledger/fabric/common/channelconfig"
)

type ConsenterSupport struct {
	consensus.ConsenterSupport
}

func NewConsenterSupport() *ConsenterSupport {
	return &ConsenterSupport{}
}

// Signer is an interface which wraps the Sign method.
func (s *ConsenterSupport) Sign(message []byte) ([]byte, error) {
	return nil, nil
}

// Serializer is an interface which wraps the Serialize function.
func (s *ConsenterSupport) Serialize() ([]byte, error) {
	return nil, nil
}

// Processor provides the methods necessary to classify and process any message which
// arrives through the Broadcast interface.

func (s *ConsenterSupport) ClassifyMsg(chdr *cb.ChannelHeader) consensus.Classification {
	return consensus.Classification(0)
}

// ProcessNormalMsg will check the validity of a message based on the current configuration.  It returns the current
// configuration sequence number and nil on success, or an error if the message is not valid
func (s *ConsenterSupport) ProcessNormalMsg(env *cb.Envelope) (configSeq uint64, err error) {
	return uint64(0), nil
}

// ProcessConfigUpdateMsg will attempt to apply the config update to the current configuration, and if successful
// return the resulting config message and the configSeq the config was computed from.  If the config update message
// is invalid, an error is returned.
func (s *ConsenterSupport) ProcessConfigUpdateMsg(env *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error) {
	return nil, uint64(0), nil
}

// ProcessConfigMsg takes message of type `ORDERER_TX` or `CONFIG`, unpack the ConfigUpdate envelope embedded
// in it, and call `ProcessConfigUpdateMsg` to produce new Config message of the same type as original message.
// This method is used to re-validate and reproduce config message, if it's deemed not to be valid anymore.
func (s *ConsenterSupport) ProcessConfigMsg(env *cb.Envelope) (*cb.Envelope, uint64, error) {
	return nil, uint64(0), nil
}

// ConsenterSupport provides the resources available to a Consenter implementation.

func BlockVerifier(header *cb.BlockHeader, metadata *cb.BlockMetadata) error {
	return nil
}

// SignatureVerifier verifies a signature of a block.
func (s *ConsenterSupport) SignatureVerifier() consensus.BlockVerifierFunc {
	return BlockVerifier
}

// BlockCutter returns the block cutting helper for this channel.
func (s *ConsenterSupport) BlockCutter() consensus.Receiver {
	return &AppReceiver{}
}

// SharedConfig provides the shared config from the channel's current config block.
func (s *ConsenterSupport) SharedConfig() channelconfig.Orderer {
	return &AppOrderer{}
}

// ChannelConfig provides the channel config from the channel's current config block.
func (s *ConsenterSupport) ChannelConfig() channelconfig.Channel {
	return &AppChannel{}
}

// CreateNextBlock takes a list of messages and creates the next block based on the block with highest block number committed to the ledger
// Note that either WriteBlock or WriteConfigBlock must be called before invoking this method a second time.
func (s *ConsenterSupport) CreateNextBlock(messages []*cb.Envelope) *cb.Block {
	return &cb.Block{}
}

// Block returns a block with the given number,
// or nil if such a block doesn't exist.
func (s *ConsenterSupport) Block(number uint64) *cb.Block {
	return &cb.Block{}
}

// WriteBlock commits a block to the ledger.
func (s *ConsenterSupport) WriteBlock(block *cb.Block, encodedMetadataValue []byte) {

}

// WriteConfigBlock commits a block to the ledger, and applies the config update inside.
func (s *ConsenterSupport) WriteConfigBlock(block *cb.Block, encodedMetadataValue []byte) {

}

// Sequence returns the current config sequence.
func (s *ConsenterSupport) Sequence() uint64 {
	return uint64(0)

}

// ChannelID returns the channel ID this support is associated with.
func (s *ConsenterSupport) ChannelID() string {
	return "channel1"
}

// Height returns the number of blocks in the chain this channel is associated with.
func (s *ConsenterSupport) Height() uint64 {
	return uint64(0)
}

// Append appends a new block to the ledger in its raw form,
// unlike WriteBlock that also mutates its metadata.
func (s *ConsenterSupport) Append(block *cb.Block) error {
	return nil
}
