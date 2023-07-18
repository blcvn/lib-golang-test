package main

import (
	"flag"
	"strings"
	"sync"

	"github.com/blcvn/lib-golang-test/raft-grpc/app"
	"github.com/blcvn/lib-golang-test/raft-grpc/transport"
	raft "go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	kvport := flag.Int("port", 9121, "key-value server port")
	raftport := flag.Int("raftport", 8080, "grpc server port")
	join := flag.Bool("join", false, "join an existing cluster")

	flag.Parse()

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	peers := strings.Split(*cluster, ",")

	channel := "test"
	mapping := transport.NewMemberMapping()

	rpeers := make([]raft.Peer, len(peers))
	for i, endpoint := range peers {
		stub := &transport.Stub{
			ID:       uint64(i + 1),
			Endpoint: endpoint,
		}
		mapping.Put(stub)
		rpeers[i] = raft.Peer{ID: uint64(i + 1)}
	}

	dialer := &transport.SDialer{}
	connectionStore := transport.NewConnectionStore(dialer)
	communicator := &transport.Comm{
		Shutdown:    false,
		Lock:        sync.RWMutex{},
		Connections: connectionStore,
		Chan2Members: map[string]transport.MemberMapping{
			channel: *mapping,
		},
	}

	// raft provides a commit stream for the proposals from the http api
	var kvs *app.Kvstore
	getSnapshot := func() ([]byte, error) { return kvs.GetSnapshot() }

	rnode, commitC, errorC, snapshotterReady := app.NewRaftNode(*id, channel, communicator, rpeers, *join, getSnapshot, proposeC, confChangeC)
	kvs = app.NewKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	//Start grpc server
	cs := &app.ChainSupport{
		MessageReceiver: rnode,
	}
	//Handle request from client to grpcServer
	dpatch := &transport.ServerHandle{
		Handler: &transport.Dispatcher{
			ChainSelector: cs,
		},
	}
	clusterServer := &transport.ClusterService{
		Dispatcher: dpatch,
	}
	go transport.StartGrpcServer(*raftport, clusterServer)

	// the key-value http handler will propose updates to raft
	app.ServeHttpKVAPI(kvs, *kvport, confChangeC, errorC)
}
