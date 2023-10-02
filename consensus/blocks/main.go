package main

import (
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/blcvn/lib-golang-test/consensus/blocks/consensus"
	cb "github.com/blcvn/lib-golang-test/consensus/blocks/types/common"
	"github.com/pkg/errors"

	"code.cloudfoundry.org/clock"
	"github.com/blcvn/lib-golang-test/consensus/blocks/app"
	"github.com/blcvn/lib-golang-test/consensus/blocks/broadcast"
	"github.com/blcvn/lib-golang-test/consensus/blocks/comm"
	"github.com/blcvn/lib-golang-test/consensus/blocks/consensus/common/cluster"
	"github.com/blcvn/lib-golang-test/consensus/blocks/consensus/common/metrics"
	"github.com/blcvn/lib-golang-test/consensus/blocks/consensus/etcdraft"
	"github.com/blcvn/lib-golang-test/consensus/blocks/consensus/protoutil"
	"github.com/blcvn/lib-golang-test/consensus/blocks/types/common/msgprocessor"
	"github.com/blcvn/lib-golang-test/consensus/blocks/types/common/types"
	"github.com/blcvn/lib-golang-test/consensus/blocks/types/orderer"
	ab "github.com/blcvn/lib-golang-test/consensus/blocks/types/orderer"
	"github.com/blcvn/lib-golang-test/log/flogging"
	"google.golang.org/protobuf/proto"

	pb_etcdraft "github.com/blcvn/lib-golang-test/consensus/blocks/types/etcdraft"
	"go.etcd.io/etcd/raft/v3"
)

var logger = flogging.MustGetLogger("blocks.main")

type Consenter struct {
	consensus.Consenter
	*etcdraft.Dispatcher
	*etcdraft.Chain
}

func (c *Consenter) HandleChain(support consensus.ConsenterSupport, metadata *cb.Metadata) (consensus.Chain, error) {
	return c.Chain, nil
}

// ReceiverByChain returns the MessageReceiver for the given channelID or nil
// if not found.
func (c *Consenter) ReceiverByChain(channelID string) etcdraft.MessageReceiver {
	logger.Infof("Consenter.ReceiverByChain: Query receiver for channel %s ", channelID)
	return c.Chain
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

type chainSupport struct {
	broadcast.ChannelSupport
}

// ClassifyMsg inspects the message header to determine which type of processing is necessary
func (s *chainSupport) ClassifyMsg(chdr *cb.ChannelHeader) msgprocessor.Classification {
	return msgprocessor.Classification(0)
}

// ProcessNormalMsg will check the validity of a message based on the current configuration.  It returns the current
// configuration sequence number and nil on success, or an error if the message is not valid
func (s *chainSupport) ProcessNormalMsg(env *cb.Envelope) (configSeq uint64, err error) {
	return uint64(0), nil
}

// ProcessConfigUpdateMsg will attempt to apply the config update to the current configuration, and if successful
// return the resulting config message and the configSeq the config was computed from.  If the config update message
// is invalid, an error is returned.
func (s *chainSupport) ProcessConfigUpdateMsg(env *cb.Envelope) (config *cb.Envelope, configSeq uint64, err error) {
	return nil, uint64(0), nil
}

// ProcessConfigMsg takes message of type `ORDERER_TX` or `CONFIG`, unpack the ConfigUpdate envelope embedded
// in it, and call `ProcessConfigUpdateMsg` to produce new Config message of the same type as original message.
// This method is used to re-validate and reproduce config message, if it's deemed not to be valid anymore.
func (s *chainSupport) ProcessConfigMsg(env *cb.Envelope) (*cb.Envelope, uint64, error) {
	return nil, uint64(0), nil
}

type Registrar struct {
	lock       sync.RWMutex
	chains     map[string]*chainSupport
	consenters map[string]consensus.Consenter
}

// GetChain retrieves the chain support for a chain if it exists.
func (r *Registrar) GetChain(chainID string) *chainSupport {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.chains[chainID]
}

// BroadcastChannelSupport returns the message channel header, whether the message is a config update
// and the channel resources for a message or an error if the message is not a message which can
// be processed directly (like CONFIG and ORDERER_TRANSACTION messages)
func (r *Registrar) BroadcastChannelSupport(msg *cb.Envelope) (*cb.ChannelHeader, bool, *chainSupport, error) {
	chdr, err := protoutil.ChannelHeader(msg)
	if err != nil {
		return nil, false, nil, errors.WithMessage(err, "could not determine channel ID")
	}

	cs := r.GetChain(chdr.ChannelId)
	// Used to be new channel creation with the system channel, but now channels are created with the channel
	// participation API only, so it is just a wrong channel name.
	if cs == nil {
		return chdr, false, nil, types.ErrChannelNotExist
	}

	isConfig := false
	switch cs.ClassifyMsg(chdr) {
	case msgprocessor.ConfigUpdateMsg:
		isConfig = true
	case msgprocessor.ConfigMsg:
		return chdr, false, nil, errors.New("message is of type that cannot be processed directly")
	// case msgprocessor.UnsupportedMsg:
	// 	return chdr, false, nil, errors.New("message is of type that is no longer supported")
	default:
	}

	return chdr, isConfig, cs, nil
}

type broadcastSupport struct {
	*Registrar
}

func (bs broadcastSupport) BroadcastChannelSupport(msg *cb.Envelope) (*cb.ChannelHeader, bool, broadcast.ChannelSupport, error) {
	return bs.Registrar.BroadcastChannelSupport(msg)
}

type broadcastMsgTracer struct {
	ab.AtomicBroadcast_BroadcastServer
}

func (bmt *broadcastMsgTracer) Recv() (*cb.Envelope, error) {
	msg, err := bmt.AtomicBroadcast_BroadcastServer.Recv()
	return msg, err
}

type server struct {
	bh *broadcast.Handler
}

// Broadcast receives a stream of messages from a client for ordering
func (s *server) Broadcast(srv ab.AtomicBroadcast_BroadcastServer) error {
	logger.Debugf("Starting new Broadcast handler")
	defer func() {
		if r := recover(); r != nil {
			logger.Criticalf("Broadcast client triggered panic: %s\n%s", r, debug.Stack())
		}
		logger.Debugf("Closing Broadcast stream")
	}()
	return s.bh.Handle(&broadcastMsgTracer{
		AtomicBroadcast_BroadcastServer: srv,
	})
}

// Deliver sends a stream of blocks to a client after ordering
func (s *server) Deliver(srv ab.AtomicBroadcast_DeliverServer) error {
	logger.Debugf("Starting new Deliver handler")
	defer func() {
		if r := recover(); r != nil {
			logger.Criticalf("Deliver client triggered panic: %s\n%s", r, debug.Stack())
		}
		logger.Debugf("Closing Deliver stream")
	}()

	// deliverServer := &deliver.Server{
	// 	Receiver: &deliverMsgTracer{
	// 		Receiver: srv,
	// 		msgTracer: msgTracer{
	// 			debug:    s.debug,
	// 			function: "Deliver",
	// 		},
	// 	},
	// 	ResponseSender: &responseSender{
	// 		AtomicBroadcast_DeliverServer: srv,
	// 	},
	// }
	// return s.dh.Handle(srv.Context(), deliverServer)
	return nil
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
	consenter.Chain = chain

	//Start grpc request
	serverConfig := comm.ServerConfig{
		ConnectionTimeout: 100 * time.Second,
		MaxRecvMsgSize:    int(100000),
		MaxSendMsgSize:    int(100000),
	}
	grpcServer := initializeGrpcServer(serverConfig)

	//binht: Create grpc transport for raft
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

	ab.RegisterClusterServer(grpcServer.Server(), clusterServer)

	broadcastMetrics := broadcast.NewMetrics(&metricProvider)

	cs := &chainSupport{}
	chains := make(map[string]*chainSupport, 0)
	chains["test"] = cs

	nconsenters := make(map[string]consensus.Consenter, 0)
	nconsenters["test"] = consenter

	r := &Registrar{
		lock:       sync.RWMutex{},
		chains:     chains,
		consenters: nconsenters,
	}

	//binht: Create grpc broadcast for client
	server := &server{
		bh: &broadcast.Handler{
			SupportRegistrar: broadcastSupport{Registrar: r},
			Metrics:          broadcastMetrics,
		},
	}
	ab.RegisterAtomicBroadcastServer(grpcServer.Server(), server)
	if err := grpcServer.Start(); err != nil {
		logger.Fatalf("Atomic Broadcast gRPC server has terminated while serving requests due to: %v", err)
	}
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
