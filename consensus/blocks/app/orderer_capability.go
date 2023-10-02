package app

import "github.com/hyperledger/fabric/common/channelconfig"

type AppOrdererCapabilities struct {
	channelconfig.OrdererCapabilities
}

// PredictableChannelTemplate specifies whether the v1.0 undesirable behavior of setting the /Channel
// group's mod_policy to "" and copy versions from the orderer system channel config should be fixed or not.
func (s *AppOrdererCapabilities) PredictableChannelTemplate() bool {
	return false
}

// Resubmission specifies whether the v1.0 non-deterministic commitment of tx should be fixed by re-submitting
// the re-validated tx.
func (s *AppOrdererCapabilities) Resubmission() bool {
	return false
}

// Supported returns an error if there are unknown capabilities in this channel which are required
func (s *AppOrdererCapabilities) Supported() error {
	return nil
}

// ExpirationCheck specifies whether the orderer checks for identity expiration checks
// when validating messages
func (s *AppOrdererCapabilities) ExpirationCheck() bool {
	return false
}

// ConsensusTypeMigration checks whether the orderer permits a consensus-type migration.
func (s *AppOrdererCapabilities) ConsensusTypeMigration() bool {
	return false
}

// UseChannelCreationPolicyAsAdmins checks whether the orderer should use more sophisticated
// channel creation logic using channel creation policy as the Admins policy if
// the creation transaction appears to support it.
func (s *AppOrdererCapabilities) UseChannelCreationPolicyAsAdmins() bool {
	return false
}
