package app

import (
	pb_msp "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric/msp"
)

type AppMSP struct {
	msp.MSP
}

// DeserializeIdentity deserializes an identity.
// Deserialization will fail if the identity is associated to
// an msp that is different from this one that is performing
// the deserialization.
func (s *AppMSP) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	return &AppSigningIdentity{}, nil
}

// IsWellFormed checks if the given identity can be deserialized into its provider-specific form
func (s *AppMSP) IsWellFormed(identity *pb_msp.SerializedIdentity) error {
	return nil
}

// Setup the MSP instance according to configuration information
func (s *AppMSP) Setup(config *pb_msp.MSPConfig) error {
	return nil
}

// GetVersion returns the version of this MSP
func (s *AppMSP) GetVersion() msp.MSPVersion {
	return msp.MSPv1_0
}

// GetType returns the provider type
func (s *AppMSP) GetType() msp.ProviderType {
	return msp.ProviderType(1)
}

// GetIdentifier returns the provider identifier
func (s *AppMSP) GetIdentifier() (string, error) {
	return "test", nil
}

// GetSigningIdentity returns a signing identity corresponding to the provided identifier
func (s *AppMSP) GetSigningIdentity(identifier *msp.IdentityIdentifier) (msp.SigningIdentity, error) {
	return &AppSigningIdentity{}, nil
}

// GetDefaultSigningIdentity returns the default signing identity
func (s *AppMSP) GetDefaultSigningIdentity() (msp.SigningIdentity, error) {
	return &AppSigningIdentity{}, nil
}

// GetTLSRootCerts returns the TLS root certificates for this MSP
func (s *AppMSP) GetTLSRootCerts() [][]byte {
	return [][]byte{}
}

// GetTLSIntermediateCerts returns the TLS intermediate root certificates for this MSP
func (s *AppMSP) GetTLSIntermediateCerts() [][]byte {
	return [][]byte{}
}

// Validate checks whether the supplied identity is valid
func (s *AppMSP) Validate(id msp.Identity) error {
	return nil
}

// SatisfiesPrincipal checks whether the identity matches
// the description supplied in MSPPrincipal. The check may
// involve a byte-by-byte comparison (if the principal is
// a serialized identity) or may require MSP validation
func (s *AppMSP) SatisfiesPrincipal(id msp.Identity, principal *pb_msp.MSPPrincipal) error {
	return nil
}
