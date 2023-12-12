package smallstep

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"

	"github.com/pkg/errors"
)

type Verifier struct {
	rootCABytes []byte
}

func InitVerifier(rootCABytes []byte) *Verifier {
	return &Verifier{
		rootCABytes: rootCABytes,
	}
}

// VerifyCert make sure input cert is valid and issued by the rootCA
func (v *Verifier) VerifyCert(cert *x509.Certificate, ipems []byte) error {
	var (
		host             = ""
		intermediatePool = x509.NewCertPool()
		rootPool         *x509.CertPool
	)

	if cert == nil {
		return errors.Errorf("input certificate bytes contains no PEM certificate blocks")
	}
	if len(ipems) > 0 && !intermediatePool.AppendCertsFromPEM(ipems) {
		return errors.Errorf("failure creating intermediate list from input certificate")
	}

	rootPool, err := certPoolFromBytes(v.rootCABytes)
	if err != nil {
		errors.Wrapf(err, "failure to load root certificate pool from input root certificate")
	}

	opts := x509.VerifyOptions{
		DNSName:       host,
		Roots:         rootPool,
		Intermediates: intermediatePool,
		// Support verification of any type of cert.
		//
		// TODO: add something like --purpose client,server,... and configure
		// this property accordingly.
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	if _, err := cert.Verify(opts); err != nil {
		return errors.Wrapf(err, "failed to verify certificate")
	}

	return nil
}

// VerifyCertSubject make sure input cert belongs to subject
func (v *Verifier) VerifyCertSubject(cert *x509.Certificate, subjectName string) error {
	if cert == nil {
		return errors.Errorf("input certificate bytes contains no PEM certificate blocks")
	}

	if cert.Subject.CommonName != subjectName {
		return errors.Errorf("cert belongs to %s, not %s", cert.Subject.CommonName, subjectName)
	}

	return nil
}

// VerifySignature make sure message is signed by the identity own this public key
func (v *Verifier) VerifySignature(pubKey *ecdsa.PublicKey, msg string, sig []byte) bool {
	hash := sha256.Sum256([]byte(msg))

	return ecdsa.VerifyASN1(pubKey, hash[:], sig)
}
