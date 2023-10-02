package app

import (
	"time"

	pb_msp "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric/msp"
)

type AppSigningIdentity struct {
	msp.SigningIdentity
}

// ExpiresAt returns the time at which the Identity expires.
// If the returned time is the zero value, it implies
// the Identity does not expire, or that its expiration
// time is unknown
func (s *AppSigningIdentity) ExpiresAt() time.Time {
	return time.Time{}
}

// GetIdentifier returns the identifier of that identity
func (s *AppSigningIdentity) GetIdentifier() *msp.IdentityIdentifier {
	return &msp.IdentityIdentifier{}
}

// GetMSPIdentifier returns the MSP Id for this instance
func (s *AppSigningIdentity) GetMSPIdentifier() string {
	return "test"
}

// Validate uses the rules that govern this identity to validate it.
// E.g., if it is a fabric TCert implemented as identity, validate
// will check the TCert signature against the assumed root certificate
// authority.
func (s *AppSigningIdentity) Validate() error {
	return nil
}

// GetOrganizationalUnits returns zero or more organization units or
// divisions this identity is related to as long as this is public
// information. Certain MSP implementations may use attributes
// that are publicly associated to this identity, or the identifier of
// the root certificate authority that has provided signatures on this
// certificate.
// Examples:
//   - if the identity is an x.509 certificate, this function returns one
//     or more string which is encoded in the Subject's Distinguished Name
//     of the type OU
//
// TODO: For X.509 based identities, check if we need a dedicated type
//
//	for OU where the Certificate OU is properly namespaced by the
//	signer's identity
func (s *AppSigningIdentity) GetOrganizationalUnits() []*msp.OUIdentifier {
	return []*msp.OUIdentifier{}
}

// Anonymous returns true if this is an anonymous identity, false otherwise
func (s *AppSigningIdentity) Anonymous() bool {
	return false
}

// Verify a signature over some message using this identity as reference
func (s *AppSigningIdentity) Verify(msg []byte, sig []byte) error {
	return nil
}

// Serialize converts an identity to bytes
func (s *AppSigningIdentity) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// SatisfiesPrincipal checks whether this instance matches
// the description supplied in MSPPrincipal. The check may
// involve a byte-by-byte comparison (if the principal is
// a serialized identity) or may require MSP validation
func (s *AppSigningIdentity) SatisfiesPrincipal(principal *pb_msp.MSPPrincipal) error {
	return nil
}

// Sign the message
func (s *AppSigningIdentity) Sign(msg []byte) ([]byte, error) {
	return []byte{}, nil
}

// GetPublicVersion returns the public parts of this identity
func (s *AppSigningIdentity) GetPublicVersion() msp.Identity {
	return s
}
