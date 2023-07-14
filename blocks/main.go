package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"code.cloudfoundry.org/clock"
	"github.com/blcvn/lib-golang-test/blocks/app"
	"github.com/blcvn/lib-golang-test/blocks/comm"
	"github.com/blcvn/lib-golang-test/blocks/consensus/common/cluster"
	"github.com/blcvn/lib-golang-test/blocks/consensus/common/metrics"
	"github.com/blcvn/lib-golang-test/blocks/consensus/etcdraft"
	"github.com/blcvn/lib-golang-test/blocks/types/orderer"
	ab "github.com/blcvn/lib-golang-test/blocks/types/orderer"
	"github.com/blcvn/lib-golang-test/flogging"
	"google.golang.org/protobuf/proto"

	pb_etcdraft "github.com/blcvn/lib-golang-test/blocks/types/etcdraft"
	"go.etcd.io/etcd/raft/v3"
)

var logger = flogging.MustGetLogger("blocks.main")

type Consenter struct {
	*etcdraft.Dispatcher
	CurrentChain *etcdraft.Chain
}

// ReceiverByChain returns the MessageReceiver for the given channelID or nil
// if not found.
func (c *Consenter) ReceiverByChain(channelID string) etcdraft.MessageReceiver {
	logger.Infof("Consenter.ReceiverByChain: Query receiver for channel %s ", channelID)
	return c.CurrentChain
}

// TargetChannel extracts the channel from the given proto.Message.
// Returns an empty string on failure.
func (c *Consenter) TargetChannel(message proto.Message) string {
	switch req := message.(type) {
	case *orderer.ConsensusRequest:
		return req.Channel
	case *orderer.SubmitRequest:
		return req.Channel
	default:
		return ""
	}
}

func main() {
	// cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	// isLearner := flag.Bool("learn", false, "learner")
	// id := flag.Int("id", 1, "node ID")
	// kvport := flag.Int("port", 9121, "key-value server port")
	// join := flag.Bool("join", false, "join an existing cluster")
	// flag.Parse()
	block_meta := pb_etcdraft.BlockMetadata{
		ConsenterIds:    []uint64{1}, //binhnt: supply raftNodeId
		NextConsenterId: 0,
		RaftIndex:       0,
	}
	consenters := make(map[uint64]*pb_etcdraft.Consenter, 0)
	consenters[1] = &pb_etcdraft.Consenter{
		Host:          "localhost",
		Port:          8080,
		ClientTlsCert: []byte{},
		ServerTlsCert: []byte{},
	}

	metricProvider := app.AppMetricProvider{}

	support := app.NewConsenterSupport()
	ms := app.AppMemoryStorage{}
	etcdRaftMetrics := etcdraft.NewMetrics(&metricProvider)
	opts := etcdraft.Options{
		RPCTimeout:             10,
		RaftID:                 1,
		Clock:                  clock.NewClock(),
		WALDir:                 "data/waldir",
		SnapDir:                "data/snapdir",
		SnapshotIntervalSize:   0,
		SnapshotCatchUpEntries: 0,
		MemoryStorage:          &ms,
		Logger:                 logger,
		TickInterval:           1 * time.Second,
		ElectionTick:           3,
		HeartbeatTick:          2,
		MaxSizePerMsg:          0,
		MaxInflightBlocks:      10,
		BlockMetadata:          &block_meta,
		Consenters:             consenters,
		MigrationInit:          false,
		Metrics:                etcdRaftMetrics,
		Cert:                   []byte{},
		EvictionSuspicion:      0,
		LeaderCheckInterval:    0,
	}
	conf := &app.Configurator{}

	connOpt := metrics.GaugeOpts{}
	tlsConnectionCount := metricProvider.NewGauge(connOpt)
	secureDialer := app.AppSecureDialer{}
	connections := cluster.NewConnectionStore(&secureDialer, tlsConnectionCount)

	consenter := &Consenter{}
	consenter.Dispatcher = &etcdraft.Dispatcher{
		Logger:        logger,
		ChainSelector: consenter,
	}

	chan2Members := make(map[string]cluster.MemberMapping)
	chan2Members["test"] = cluster.MemberMapping{
		SamePublicKey: func([]byte, []byte) bool {
			return true
		},
	}

	communicator := &cluster.Comm{
		MinimumExpirationWarningInterval: 0,
		CertExpWarningThreshold:          0,
		SendBufferSize:                   1000,
		Lock:                             sync.RWMutex{},
		Logger:                           logger,
		ChanExt:                          consenter,
		H:                                consenter,
		Connections:                      connections,
		Chan2Members:                     chan2Members,
		Metrics:                          &cluster.Metrics{},
		CompareCertificate: func([]byte, []byte) bool {
			return true
		},
	}

	streamByType := cluster.NewStreamsByType()

	rpc := &cluster.RPC{
		Logger:        logger,
		Timeout:       100 * time.Second,
		Channel:       "testing",
		Comm:          communicator,
		StreamsByType: streamByType,
	} //binhnt.Using rpc as

	f := app.CreateBlockPuller

	observeC := make(chan<- raft.SoftState, 0)

	chain, err := etcdraft.NewChain(support, opts, conf, rpc, f, app.HaltCallback, observeC)
	if err != nil {
		fmt.Printf("Errror ")
		return
	}
	chain.Start()

	//Connect Chain with Grpserver

	consenter.CurrentChain = chain

	//Start grpc request
	serverConfig := comm.ServerConfig{
		ConnectionTimeout: 100 * time.Second,
		MaxRecvMsgSize:    int(100000),
		MaxSendMsgSize:    int(100000),
	}
	srv := initializeGrpcServer(serverConfig)

	clusterMetrics := cluster.NewMetrics(&metricProvider)
	clusterServer := &cluster.Service{
		StreamCountReporter: &cluster.StreamCountReporter{
			Metrics: clusterMetrics,
		},
		Dispatcher:                       communicator,
		Logger:                           logger,
		StepLogger:                       logger,
		MinimumExpirationWarningInterval: 0,
		CertExpWarningThreshold:          0,
	}

	ab.RegisterClusterServer(srv.Server(), clusterServer)
	srv.Start()
}
func initializeGrpcServer(serverConfig comm.ServerConfig) *comm.GRPCServer {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "localhost", 8080))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	// Create GRPC server - return if an error occurs
	grpcServer, err := comm.NewGRPCServerFromListener(lis, serverConfig)
	if err != nil {
		log.Fatal("Failed to return new GRPC server:", err)
	}

	return grpcServer
}
