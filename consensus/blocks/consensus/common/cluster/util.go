package cluster

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/blcvn/lib-golang-test/consensus/blocks/types/orderer"
	"github.com/hyperledger/fabric/common/util"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// ConnByCertMap maps certificates represented as strings
// to gRPC connections
type ConnByCertMap map[string]*grpc.ClientConn

// Lable used for TLS Export Keying Material call
const KeyingMaterialLabel = "orderer v3 authentication label"

// Lookup looks up a certificate and returns the connection that was mapped
// to the certificate, and whether it was found or not
func (cbc ConnByCertMap) Lookup(cert []byte) (*grpc.ClientConn, bool) {
	conn, ok := cbc[string(cert)]
	return conn, ok
}

// Put associates the given connection to the certificate
func (cbc ConnByCertMap) Put(cert []byte, conn *grpc.ClientConn) {
	cbc[string(cert)] = conn
}

// Remove removes the connection that is associated to the given certificate
func (cbc ConnByCertMap) Remove(cert []byte) {
	delete(cbc, string(cert))
}

// Size returns the size of the connections by certificate mapping
func (cbc ConnByCertMap) Size() int {
	return len(cbc)
}

func requestAsString(request *orderer.StepRequest) string {
	switch t := request.GetPayload().(type) {
	case *orderer.StepRequest_SubmitRequest:
		if t.SubmitRequest == nil || t.SubmitRequest.Payload == nil {
			return fmt.Sprintf("Empty SubmitRequest: %v", t.SubmitRequest)
		}
		return fmt.Sprintf("SubmitRequest for channel %s with payload of size %d",
			t.SubmitRequest.Channel, len(t.SubmitRequest.Payload.Payload))
	case *orderer.StepRequest_ConsensusRequest:
		return fmt.Sprintf("ConsensusRequest for channel %s with payload of size %d",
			t.ConsensusRequest.Channel, len(t.ConsensusRequest.Payload))
	default:
		return fmt.Sprintf("unknown type: %v", request)
	}
}

// DERtoPEM returns a PEM representation of the DER
// encoded certificate
func DERtoPEM(der []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	}))
}

// ExtractCertificateFromContext returns the TLS certificate (if applicable)
// from the given context of a gRPC stream
func ExtractCertificateFromContext(ctx context.Context) *x509.Certificate {
	pr, extracted := peer.FromContext(ctx)
	if !extracted {
		return nil
	}

	authInfo := pr.AuthInfo
	if authInfo == nil {
		return nil
	}

	tlsInfo, isTLSConn := authInfo.(credentials.TLSInfo)
	if !isTLSConn {
		return nil
	}
	certs := tlsInfo.State.PeerCertificates
	if len(certs) == 0 {
		return nil
	}
	return certs[0]
}

// ExtractRawCertificateFromContext returns the raw TLS certificate (if applicable)
// from the given context of a gRPC stream
func ExtractRawCertificateFromContext(ctx context.Context) []byte {
	cert := ExtractCertificateFromContext(ctx)
	if cert == nil {
		return nil
	}
	return cert.Raw
}

//binhnt add

func GetSessionBindingHash(authReq *orderer.NodeAuthRequest) []byte {
	return util.ComputeSHA256(util.ConcatenateBytes(
		[]byte(strconv.FormatUint(uint64(authReq.Version), 10)),
		[]byte(authReq.Timestamp.String()),
		[]byte(strconv.FormatUint(authReq.FromId, 10)),
		[]byte(strconv.FormatUint(authReq.ToId, 10)),
		[]byte(authReq.Channel),
	))
}

func GetTLSSessionBinding(ctx context.Context, bindingPayload []byte) ([]byte, error) {
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("failed extracting stream context")
	}
	connState := peerInfo.AuthInfo.(credentials.TLSInfo).State

	tlsBinding, err := exportKM(connState, KeyingMaterialLabel, bindingPayload)
	if err != nil {
		return nil, errors.Wrap(err, "failed exporting keying material")
	}

	return tlsBinding, nil
}

func exportKM(cs tls.ConnectionState, label string, context []byte) ([]byte, error) {
	tlsBinding, err := cs.ExportKeyingMaterial(label, context, 32)
	if err != nil {
		return nil, errors.Wrap(err, "failed generating TLS Binding material")
	}
	return tlsBinding, nil
}

func VerifySignature(identity, msgHash, signature []byte) error {
	block, _ := pem.Decode(identity)
	if block == nil {
		return errors.New("pem decoding failed")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.Wrap(err, "key extraction failed")
	}

	pubKey, isECDSA := cert.PublicKey.(*ecdsa.PublicKey)
	if !isECDSA {
		return errors.New("not valid public key")
	}

	validSignature := ecdsa.VerifyASN1(pubKey, msgHash, signature)

	if !validSignature {
		return errors.New("signature invalid")
	}
	return nil
}

func SHA256Digest(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

type certificateExpirationCheck struct {
	minimumExpirationWarningInterval time.Duration
	expiresAt                        time.Time
	expirationWarningThreshold       time.Duration
	lastWarning                      time.Time
	nodeName                         string
	endpoint                         string
	alert                            func(string, ...interface{})
}

func (exp *certificateExpirationCheck) checkExpiration(currentTime time.Time, channel string) {
	timeLeft := exp.expiresAt.Sub(currentTime)
	if timeLeft > exp.expirationWarningThreshold {
		return
	}

	timeSinceLastWarning := currentTime.Sub(exp.lastWarning)
	if timeSinceLastWarning < exp.minimumExpirationWarningInterval {
		return
	}

	exp.alert("Certificate of %s from %s for channel %s expires in less than %v",
		exp.nodeName, exp.endpoint, channel, timeLeft)
	exp.lastWarning = currentTime
}

// ExtractCertificateFromContext returns the TLS certificate (if applicable)
// from the given context of a gRPC stream
func Util_ExtractCertificateFromContext(ctx context.Context) *x509.Certificate {
	pr, extracted := peer.FromContext(ctx)
	if !extracted {
		return nil
	}

	authInfo := pr.AuthInfo
	if authInfo == nil {
		return nil
	}

	tlsInfo, isTLSConn := authInfo.(credentials.TLSInfo)
	if !isTLSConn {
		return nil
	}
	certs := tlsInfo.State.PeerCertificates
	if len(certs) == 0 {
		return nil
	}
	return certs[0]
}

// StreamCountReporter reports the number of streams currently connected to this node
type StreamCountReporter struct {
	Metrics *Metrics
	count   uint32
}

func (scr *StreamCountReporter) Increment() {
	count := atomic.AddUint32(&scr.count, 1)
	scr.Metrics.reportStreamCount(count)
}

func (scr *StreamCountReporter) Decrement() {
	count := atomic.AddUint32(&scr.count, ^uint32(0))
	scr.Metrics.reportStreamCount(count)
}
