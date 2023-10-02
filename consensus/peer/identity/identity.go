package identity

// SignerSerializer groups the Sign and Serialize methods.
type SignerSerializer interface {
	Signer
	Serializer
}

type Signer interface {
	Sign(message []byte) ([]byte, error)
}

// Serializer is an interface which wraps the Serialize function.
//
// Serialize converts an identity to bytes.  It returns an error on failure.
type Serializer interface {
	Serialize() ([]byte, error)
}
