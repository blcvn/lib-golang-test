package comm

import (
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/pkg/errors"
)

// SecureOptions defines the TLS security parameters for a GRPCServer or
// GRPCClient instance.
type SecureOptions struct {
	// VerifyCertificate, if not nil, is called after normal
	// certificate verification by either a TLS client or server.
	// If it returns a non-nil error, the handshake is aborted and that error results.
	VerifyCertificate func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error
	// PEM-encoded X509 public key to be used for TLS communication
	Certificate []byte
	// PEM-encoded private key to be used for TLS communication
	Key []byte
	// Set of PEM-encoded X509 certificate authorities used by clients to
	// verify server certificates
	ServerRootCAs [][]byte
	// Set of PEM-encoded X509 certificate authorities used by servers to
	// verify client certificates
	ClientRootCAs [][]byte
	// Whether or not to use TLS for communication
	UseTLS bool
	// Whether or not TLS client must present certificates for authentication
	RequireClientCert bool
	// CipherSuites is a list of supported cipher suites for TLS
	CipherSuites []uint16
	// TimeShift makes TLS handshakes time sampling shift to the past by a given duration
	TimeShift time.Duration
	// ServerNameOverride is used to verify the hostname on the returned certificates. It
	// is also included in the client's handshake to support virtual hosting
	// unless it is an IP address.
	ServerNameOverride string
}

func (so SecureOptions) TLSConfig() (*tls.Config, error) {
	// if TLS is not enabled, return
	if !so.UseTLS {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion:            tls.VersionTLS12,
		ServerName:            so.ServerNameOverride,
		VerifyPeerCertificate: so.VerifyCertificate,
	}
	if len(so.ServerRootCAs) > 0 {
		tlsConfig.RootCAs = x509.NewCertPool()
		for _, certBytes := range so.ServerRootCAs {
			if !tlsConfig.RootCAs.AppendCertsFromPEM(certBytes) {
				return nil, errors.New("error adding root certificate")
			}
		}
	}

	if so.RequireClientCert {
		cert, err := so.ClientCertificate()
		if err != nil {
			return nil, errors.WithMessage(err, "failed to load client certificate")
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}

	if so.TimeShift > 0 {
		tlsConfig.Time = func() time.Time {
			return time.Now().Add((-1) * so.TimeShift)
		}
	}

	return tlsConfig, nil
}

// ClientCertificate returns the client certificate that will be used
// for mutual TLS.
func (so SecureOptions) ClientCertificate() (tls.Certificate, error) {
	if so.Key == nil || so.Certificate == nil {
		return tls.Certificate{}, errors.New("both Key and Certificate are required when using mutual TLS")
	}
	cert, err := tls.X509KeyPair(so.Certificate, so.Key)
	if err != nil {
		return tls.Certificate{}, errors.WithMessage(err, "failed to create key pair")
	}
	return cert, nil
}
