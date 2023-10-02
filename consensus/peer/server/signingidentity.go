package server

import (
	"crypto"
	"time"

	pb "github.com/hyperledger/fabric-protos-go/msp"

	"github.com/hyperledger/fabric/msp"
)

type SigningIdentity struct {
	msp.SigningIdentity

	// The identifier of the associated membership service provider
	Mspid string
	// The identifier for an identity within a provider
	Id     string
	signer crypto.Signer
}

func NewSigner() *SigningIdentity {
	return &SigningIdentity{
		SigningIdentity: nil,
		signer:          nil,
	}
}

// Sign the message
func (s SigningIdentity) Sign(msg []byte) ([]byte, error) {
	return nil, nil
}

// GetPublicVersion returns the public parts of this identity
func (s SigningIdentity) GetPublicVersion() msp.Identity {
	return s
}

func (s SigningIdentity) DeserializeIdentity(serializedIdentity []byte) (msp.Identity, error) {
	return s, nil
}

// IsWellFormed checks if the given identity can be deserialized into its provider-specific form
func (s SigningIdentity) IsWellFormed(identity *pb.SerializedIdentity) error {
	return nil
}

// ExpiresAt returns the time at which the Identity expires.
// If the returned time is the zero value, it implies
// the Identity does not expire, or that its expiration
// time is unknown
func (s SigningIdentity) ExpiresAt() time.Time {
	return time.Now()
}

// GetIdentifier returns the identifier of that identity
func (s SigningIdentity) GetIdentifier() *msp.IdentityIdentifier {
	return &msp.IdentityIdentifier{
		Mspid: s.Mspid,
		Id:    s.Id,
	}
}

// GetMSPIdentifier returns the MSP Id for this instance
func (s SigningIdentity) GetMSPIdentifier() string {
	return s.Mspid
}

// Validate uses the rules that govern this identity to validate it.
// E.g., if it is a fabric TCert implemented as identity, validate
// will check the TCert signature against the assumed root certificate
// authority.
func (s SigningIdentity) Validate() error {
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
func (s SigningIdentity) GetOrganizationalUnits() []*msp.OUIdentifier {
	return []*msp.OUIdentifier{}
}

// Anonymous returns true if this is an anonymous identity, false otherwise
func (s SigningIdentity) Anonymous() bool {
	return false
}

// Verify a signature over some message using this identity as reference
func (s SigningIdentity) Verify(msg []byte, sig []byte) error {
	return nil
}

// Serialize converts an identity to bytes
func (s SigningIdentity) Serialize() ([]byte, error) {
	return nil, nil
}

// SatisfiesPrincipal checks whether this instance matches
// the description supplied in MSPPrincipal. The check may
// involve a byte-by-byte comparison (if the principal is
// a serialized identity) or may require MSP validation
func (s SigningIdentity) SatisfiesPrincipal(principal *pb.MSPPrincipal) error {
	return nil
}
