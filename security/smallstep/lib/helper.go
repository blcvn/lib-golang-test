package smallstep

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/Hnampk/prometheuslog/flogging"
	"github.com/pkg/errors"
	"github.com/smallstep/cli/utils"
	"go.step.sm/cli-utils/errs"
	"go.step.sm/crypto/keyutil"
	"go.step.sm/crypto/pemutil"
)

var (
	helperLogger = flogging.MustGetLogger("libs.ca.smallstep.helper")
)

func certPoolFromBytes(certBytes []byte) (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	var found bool
	if ok := pool.AppendCertsFromPEM(certBytes); ok {
		found = true
	}
	if !found {
		return nil, errors.New("error reading cert pool: not certificates found")
	}

	return pool, nil
}

func PublicKeyFromCertOrPrivFile(filePath string, password []byte) (*ecdsa.PublicKey, error) {
	var b, err = utils.ReadFile(filePath)
	if err != nil {
		return nil, errs.FileError(err, filePath)
	}

	opts := []pemutil.Options{pemutil.WithFilename(filePath), pemutil.WithFirstBlock()}
	if len(password) != 0 {
		opts = append(opts, pemutil.WithPassword(password))
	}

	return PublicKeyFromCert(b, opts)
}

func PublicKeyFromCert(keyBytes []byte, opts []pemutil.Options) (*ecdsa.PublicKey, error) {
	k, err := pemutil.ParseKey(keyBytes, opts...)
	if err != nil {
		return nil, err
	}

	pub, err := keyutil.PublicKey(k)
	if err != nil {
		return nil, err
	}

	block, err := pemutil.Serialize(pub)
	if err != nil {
		return nil, err
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		helperLogger.Errorf("error while ParsePKIXPublicKey: %s", err.Error())
		return nil, err
	}
	return pubKey.(*ecdsa.PublicKey), nil
}

func LoadPrivKeyFromFile(path string) *ecdsa.PrivateKey {
	privBytes, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("error while read Priv file: %s", err.Error())
	}
	block, _ := pem.Decode(privBytes)
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return privateKey
}

func CertFromBytes(crtBytes []byte) (*x509.Certificate, []byte, error) {
	var (
		cert  *x509.Certificate
		block *pem.Block
		ipems []byte
	)
	// The first certificate PEM in the file is our leaf Certificate.
	// Any certificate after the first is added to the list of Intermediate
	// certificates used for path validation.
	for len(crtBytes) > 0 {
		block, crtBytes = pem.Decode(crtBytes)
		if block == nil {
			return nil, nil, errors.Errorf("input certificate bytes contains an invalid PEM block")
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		var err error
		if cert == nil {
			cert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, nil, errors.WithStack(err)
			}
		} else {
			ipems = append(ipems, pem.EncodeToMemory(block)...)
		}
	}
	if cert == nil {
		return nil, nil, errors.Errorf("input certificate bytes contains no PEM certificate blocks")
	}

	return cert, ipems, nil
}
