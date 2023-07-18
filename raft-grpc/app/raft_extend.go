package app

import (
	"context"
	"fmt"

	"github.com/blcvn/lib-golang-test/blocks/consensus/protoutil"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/orderer"
	raft "go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

// Send sends out the given messages to the remote peers.
// Each message has a To field, which is an id that maps
// to an existing peer in the transport.
// If the id cannot be found in the transport, the message
// will be ignored.

func (s *RaftNode) Send(msgs []raftpb.Message) {
	s.unreachableLock.RLock()
	defer s.unreachableLock.RUnlock()

	for _, msg := range msgs {
		if msg.To == 0 {
			continue
		}

		status := raft.SnapshotFinish

		// Replace node list in snapshot with CURRENT node list in cluster.
		if msg.Type == raftpb.MsgSnap {
			msg.Snapshot.Metadata.ConfState = s.confState
		}

		msgBytes := protoutil.MarshalOrPanic(&msg)
		err := s.transport.SendConsensus(msg.To, &orderer.ConsensusRequest{Channel: s.chainID, Payload: msgBytes})
		if err != nil {
			s.node.ReportUnreachable(msg.To)
			s.logSendFailure(msg.To, err)

			status = raft.SnapshotFailure
		} else if _, ok := s.unreachable[msg.To]; ok {
			delete(s.unreachable, msg.To)
		}

		if msg.Type == raftpb.MsgSnap {
			s.node.ReportSnapshot(msg.To, status)
		}
	}
}

func (n *RaftNode) logSendFailure(dest uint64, err error) {
	if _, ok := n.unreachable[dest]; ok {
		fmt.Printf("Failed to send StepRequest to %d, because: %s", dest, err)
		return
	}

	fmt.Printf("Failed to send StepRequest to %d, because: %s", dest, err)
	n.unreachable[dest] = struct{}{}
}

// Consensus passes the given ConsensusRequest message to the raft.Node instance
func (c *RaftNode) Consensus(req *orderer.ConsensusRequest, sender uint64) error {
	stepMsg := &raftpb.Message{}
	if err := proto.Unmarshal(req.Payload, stepMsg); err != nil {
		return fmt.Errorf("failed to unmarshal StepRequest payload to Raft Message: %s", err)
	}

	if stepMsg.To != uint64(c.id) {
		fmt.Printf("halt......")
		return nil
	}

	if err := c.node.Step(context.TODO(), *stepMsg); err != nil {
		return fmt.Errorf("failed to process Raft Step message: %s", err)
	}

	// if len(req.Metadata) == 0 || atomic.LoadUint64(&c.lastKnownLeader) != sender { // ignore metadata from non-leader
	// 	return nil
	// }

	// clusterMetadata := &orderer.ClusterMetadata{}
	// if err := proto.Unmarshal(req.Metadata, clusterMetadata); err != nil {
	// 	return errors.Errorf("failed to unmarshal ClusterMetadata: %s", err)
	// }

	return nil
}

// Submit passes the given SubmitRequest message to the MessageReceiver
func (c *RaftNode) Submit(req *orderer.SubmitRequest, sender uint64) error {
	return nil
}
