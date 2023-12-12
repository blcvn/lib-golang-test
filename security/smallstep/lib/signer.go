package smallstep

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
)

type Signer struct {
	privKey *ecdsa.PrivateKey
	cert    []byte
}

// InitSigner with a private key
func InitSigner(privKey *ecdsa.PrivateKey, certBytes []byte) *Signer {
	return &Signer{
		privKey: privKey,
		cert:    certBytes,
	}
}

// Sign
func (s *Signer) Sign(msg string) ([]byte, error) {
	hash := sha256.Sum256([]byte(msg))

	signature, err := ecdsa.SignASN1(rand.Reader, s.privKey, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (s *Signer) GetCertBytes() []byte {
	return s.cert
}
