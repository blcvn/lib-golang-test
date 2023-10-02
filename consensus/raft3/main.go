package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/blcvn/lib-golang-test/consensus/raft3/app"
	"github.com/blcvn/lib-golang-test/consensus/raft3/v3/raftpb"
)

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	isLearner := flag.Bool("learn", false, "learner")
	id := flag.Int("id", 1, "node ID")
	kvport := flag.Int("port", 9121, "key-value server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	// raft provides a commit stream for the proposals from the http api
	var kvs *app.Kvstore

	getSnapshot := func() ([]byte, error) { return kvs.GetSnapshot() }
	peers := strings.Split(*cluster, ",")

	if *isLearner {
		fmt.Println("start learner")
	}
	commitC, errorC, snapshotterReady := app.NewRaftNode(*id, *isLearner, peers, *join, getSnapshot, proposeC, confChangeC)

	kvs = app.NewKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	// the key-value http handler will propose updates to raft
	app.ServeHttpKVAPI(kvs, *kvport, confChangeC, errorC)
}
