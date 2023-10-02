package transport

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/pkg/errors"
	"github.com/xiang90/probing"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// OperationType denotes a type of operation that the RPC can perform
// such as sending a transaction, or a consensus related message.
type OperationType int

const (
	ConsensusOperation OperationType = iota
	SubmitOperation
)

// Transport implements Transporter interface. It provides the functionality
// to send raft messages to peers, and receive raft messages from peers.
// User should call Handler method to get a handler to serve requests
// received from peerURLs.
// User needs to call Start before calling other functions, and call
// Stop when the Transport is no longer used.

// implement RPC
type Transport struct {
	Logger *zap.Logger

	DialTimeout time.Duration // maximum duration before timing out dial of the request
	// DialRetryFrequency defines the frequency of streamReader dial retrial attempts;
	// a distinct rate limiter is created per every peer (default value: 10 events/sec)
	DialRetryFrequency rate.Limit

	TLSInfo transport.TLSInfo // TLS information used when creating connection

	ID          types.ID   // local member ID
	URLs        types.URLs // local peer URLs
	ClusterID   types.ID   // raft cluster ID for request validation
	Raft        Raft       // raft state machine, to which the Transport forwards received messages and reports status
	Snapshotter *snap.Snapshotter
	ServerStats *stats.ServerStats // used to record general transportation statistics
	// used to record transportation statistics with followers when
	// performing as leader in raft protocol
	LeaderStats *stats.LeaderStats
	// ErrorC is used to report detected critical errors, e.g.,
	// the member has been permanently removed from the cluster
	// When an error is received from ErrorC, user should stop raft state
	// machine and thus stop the Transport.
	ErrorC chan error

	streamRt   http.RoundTripper // roundTripper used by streams
	pipelineRt http.RoundTripper // roundTripper used by pipelines

	mu sync.RWMutex // protect the remote and peer map
	// remotes map[types.ID]*remote // remotes map that helps newly joined member to catch up
	// peers   map[types.ID]Peer    // peers map

	pipelineProber probing.Prober
	streamProber   probing.Prober

	consensusLock sync.Mutex
	submitLock    sync.Mutex
	Timeout       time.Duration
	Channel       string
	Comm          Communicator
	lock          sync.RWMutex
	StreamsByType map[OperationType]map[uint64]*Stream
}

// NewStreamsByType returns a mapping of operation type to
// a mapping of destination to stream.
func NewStreamsByType() map[OperationType]map[uint64]*Stream {
	m := make(map[OperationType]map[uint64]*Stream)
	m[ConsensusOperation] = make(map[uint64]*Stream)
	m[SubmitOperation] = make(map[uint64]*Stream)
	return m
}

// Start starts the given Transporter.
// Start MUST be called before calling other functions in the interface.
func (s *Transport) Start() error {
	return nil
}

// Peer urls are used to connect to the remote peer.
func (s *Transport) AddPeer(id types.ID, urls []string) {

}

// RemovePeer removes the peer with given id.
func (s *Transport) RemovePeer(id types.ID) {

}

// Stop closes the connections and stops the transporter.
func (s *Transport) Stop() {

}
func (s *Transport) SendConsensus(destination uint64, msg *orderer.ConsensusRequest) error {
	stream, err := s.getOrCreateStream(destination, ConsensusOperation)
	if err != nil {
		fmt.Printf("Transport.SendConsensus: getOrCreateStream failed: %s \n", err)
		return err
	}

	req := &orderer.StepRequest{
		Payload: &orderer.StepRequest_ConsensusRequest{
			ConsensusRequest: msg,
		},
	}

	s.consensusLock.Lock()
	defer s.consensusLock.Unlock()

	fmt.Printf("Transport.SendConsensus: stream send request to %d \n", destination)
	err = stream.Send(req)
	if err != nil {
		fmt.Printf("Transport.SendConsensus: stream.Send failed %s \n", err)
		s.unMapStream(destination, ConsensusOperation, stream.ID)
	}

	fmt.Printf("Transport.SendConsensus: stream.Send to %d  ok  \n", destination)

	return err
}

// getOrCreateStream obtains a Submit stream for the given destination node
func (s *Transport) getOrCreateStream(destination uint64, operationType OperationType) (*Stream, error) {
	stream := s.getStream(destination, operationType)
	if stream != nil {
		return stream, nil
	}
	stub, err := s.Comm.Remote(s.Channel, destination)
	if err != nil {
		fmt.Printf("Transport.getOrCreateStream: Comm.Remote failed %d \n", err)
		return nil, errors.WithStack(err)
	}
	stream, err = stub.NewStream(s.Timeout)
	if err != nil {
		fmt.Printf("Transport.getOrCreateStream: stub.NewStream failed: %s \n", err)
		return nil, err
	}
	fmt.Printf("Transport.getOrCreateStream: add stream to %d in map \n", destination)
	s.mapStream(destination, stream, operationType)
	return stream, nil
}

func (s *Transport) getStream(destination uint64, operationType OperationType) *Stream {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.StreamsByType[operationType][destination]
}

func (s *Transport) mapStream(destination uint64, stream *Stream, operationType OperationType) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.StreamsByType[operationType][destination] = stream
	s.cleanCanceledStreams(operationType)
}

func (s *Transport) unMapStream(destination uint64, operationType OperationType, streamIDToUnmap uint64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	stream, exists := s.StreamsByType[operationType][destination]
	if !exists {
		fmt.Printf("No %d stream to %d found, nothing to unmap", operationType, destination)
		return
	}

	if stream.ID != streamIDToUnmap {
		fmt.Printf("Stream for %d to %d has an ID of %d, not %d", operationType, destination, stream.ID, streamIDToUnmap)
		return
	}

	delete(s.StreamsByType[operationType], destination)
}
func (s *Transport) SendSubmit(dest uint64, request *orderer.SubmitRequest, report func(err error)) error {
	return nil
}

func (s *Transport) cleanCanceledStreams(operationType OperationType) {
	for destination, stream := range s.StreamsByType[operationType] {
		if !stream.Canceled() {
			continue
		}
		fmt.Printf("Removing stream %d to %d for channel %s because it is canceled", stream.ID, destination, s.Channel)
		delete(s.StreamsByType[operationType], destination)
	}
}

func submitMsgLength(request *orderer.SubmitRequest) int {
	if request.Payload == nil {
		return 0
	}
	return len(request.Payload.Payload)
}
