package config

import (
	"time"

	"github.com/blcvn/lib-golang-test/log/flogging"
	"google.golang.org/grpc"
)

// ServerConfig defines the parameters for configuring a GRPCServer instance
type ServerConfig struct {
	// ConnectionTimeout specifies the timeout for connection establishment
	// for all new connections
	ConnectionTimeout time.Duration
	// SecOpts defines the security parameters
	SecOpts SecureOptions
	// KaOpts defines the keepalive parameters
	KaOpts KeepaliveOptions
	// StreamInterceptors specifies a list of interceptors to apply to
	// streaming RPCs.  They are executed in order.
	StreamInterceptors []grpc.StreamServerInterceptor
	// UnaryInterceptors specifies a list of interceptors to apply to unary
	// RPCs.  They are executed in order.
	UnaryInterceptors []grpc.UnaryServerInterceptor
	// Logger specifies the logger the server will use
	Logger *flogging.FabricLogger
	// HealthCheckEnabled enables the gRPC Health Checking Protocol for the server
	HealthCheckEnabled bool
	// ServerStatsHandler should be set if metrics on connections are to be reported.
	// ServerStatsHandler *ServerStatsHandler
	// Maximum message size the server can receive
	MaxRecvMsgSize int
	// Maximum message size the server can send
	MaxSendMsgSize int
}
