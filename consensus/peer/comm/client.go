package comm

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	DefaultMaxRecvMsgSize = 100 * 1024 * 1024
	DefaultMaxSendMsgSize = 100 * 1024 * 1024
)

// ClientConfig defines the parameters for configuring a GRPCClient instance
type ClientConfig struct {
	// SecOpts defines the security parameters
	SecOpts SecureOptions
	// KaOpts defines the keepalive parameters
	KaOpts KeepaliveOptions
	// DialTimeout controls how long the client can block when attempting to
	// establish a connection to a server
	DialTimeout time.Duration
	// AsyncConnect makes connection creation non blocking
	AsyncConnect bool
	// Maximum message size the client can receive
	MaxRecvMsgSize int
	// Maximum message size the client can send
	MaxSendMsgSize int
}

// Convert the ClientConfig to the approriate set of grpc.DialOptions.
func (cc ClientConfig) DialOptions() ([]grpc.DialOption, error) {
	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                cc.KaOpts.ClientInterval,
		Timeout:             cc.KaOpts.ClientTimeout,
		PermitWithoutStream: true,
	}))

	// Unless asynchronous connect is set, make connection establishment blocking.
	if !cc.AsyncConnect {
		dialOpts = append(dialOpts,
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
		)
	}
	// set send/recv message size to package defaults
	maxRecvMsgSize := DefaultMaxRecvMsgSize
	if cc.MaxRecvMsgSize != 0 {
		maxRecvMsgSize = cc.MaxRecvMsgSize
	}
	maxSendMsgSize := DefaultMaxSendMsgSize
	if cc.MaxSendMsgSize != 0 {
		maxSendMsgSize = cc.MaxSendMsgSize
	}
	dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
		grpc.MaxCallSendMsgSize(maxSendMsgSize),
	))

	tlsConfig, err := cc.SecOpts.TLSConfig()
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		// transportCreds := &DynamicClientCredentials{TLSConfig: tlsConfig}
		// dialOpts = append(dialOpts, grpc.WithTransportCredentials(transportCreds))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}

	return dialOpts, nil
}

func (cc ClientConfig) Dial(address string) (*grpc.ClientConn, error) {
	dialOpts, err := cc.DialOptions()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), cc.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, dialOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new connection")
	}
	return conn, nil
}

// Clone clones this ClientConfig
func (cc ClientConfig) Clone() ClientConfig {
	shallowClone := cc
	return shallowClone
}
