package comm

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// KeepaliveOptions is used to set the gRPC keepalive settings for both
// clients and servers
type KeepaliveOptions struct {
	// ClientInterval is the duration after which if the client does not see
	// any activity from the server it pings the server to see if it is alive
	ClientInterval time.Duration
	// ClientTimeout is the duration the client waits for a response
	// from the server after sending a ping before closing the connection
	ClientTimeout time.Duration
	// ServerInterval is the duration after which if the server does not see
	// any activity from the client it pings the client to see if it is alive
	ServerInterval time.Duration
	// ServerTimeout is the duration the server waits for a response
	// from the client after sending a ping before closing the connection
	ServerTimeout time.Duration
	// ServerMinInterval is the minimum permitted time between client pings.
	// If clients send pings more frequently, the server will disconnect them
	ServerMinInterval time.Duration
}

// ServerKeepaliveOptions returns gRPC keepalive options for a server.
func (ka KeepaliveOptions) ServerKeepaliveOptions() []grpc.ServerOption {
	var serverOpts []grpc.ServerOption
	kap := keepalive.ServerParameters{
		Time:    ka.ServerInterval,
		Timeout: ka.ServerTimeout,
	}
	serverOpts = append(serverOpts, grpc.KeepaliveParams(kap))
	kep := keepalive.EnforcementPolicy{
		MinTime: ka.ServerMinInterval,
		// allow keepalive w/o rpc
		PermitWithoutStream: true,
	}
	serverOpts = append(serverOpts, grpc.KeepaliveEnforcementPolicy(kep))
	return serverOpts
}

// ClientKeepaliveOptions returns gRPC keepalive dial options for clients.
func (ka KeepaliveOptions) ClientKeepaliveOptions() []grpc.DialOption {
	var dialOpts []grpc.DialOption
	kap := keepalive.ClientParameters{
		Time:                ka.ClientInterval,
		Timeout:             ka.ClientTimeout,
		PermitWithoutStream: true,
	}
	dialOpts = append(dialOpts, grpc.WithKeepaliveParams(kap))
	return dialOpts
}
