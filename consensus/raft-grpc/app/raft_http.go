package app

// func (rc *raftNode) startRaftTransport() {
// 	rc.transport = &rafthttp.Transport{
// 		Logger:      rc.logger,
// 		ID:          types.ID(rc.id),
// 		ClusterID:   0x1000,
// 		Raft:        rc,
// 		ServerStats: stats.NewServerStats("", ""),
// 		LeaderStats: stats.NewLeaderStats(zap.NewExample(), strconv.Itoa(rc.id)),
// 		ErrorC:      make(chan error),
// 	}

// 	rc.transport.Start()
// 	for i := range rc.peers {
// 		if i+1 != rc.id {
// 			rc.transport.AddPeer(types.ID(i+1), []string{rc.peers[i]})
// 		}
// 	}

// 	go rc.serveRaft()
// }

// func (rc *raftNode) serveRaft() {
// 	url, err := url.Parse(rc.peers[rc.id-1])
// 	if err != nil {
// 		log.Fatalf("raftexample: Failed parsing URL (%v)", err)
// 	}

// 	ln, err := newStoppableListener(url.Host, rc.httpstopc)
// 	if err != nil {
// 		log.Fatalf("raftexample: Failed to listen rafthttp (%v)", err)
// 	}

// 	err = (&http.Server{Handler: rc.transport.Handler()}).Serve(ln)
// 	select {
// 	case <-rc.httpstopc:
// 	default:
// 		log.Fatalf("raftexample: Failed to serve rafthttp (%v)", err)
// 	}
// 	close(rc.httpdonec)
// }

// // Support transport
// func (rc *raftNode) Process(ctx context.Context, m raftpb.Message) error {
// 	return rc.node.Step(ctx, m)
// }
// func (rc *raftNode) IsIDRemoved(id uint64) bool  { return false }
// func (rc *raftNode) ReportUnreachable(id uint64) { rc.node.ReportUnreachable(id) }
// func (rc *raftNode) ReportSnapshot(id uint64, status raft.SnapshotStatus) {
// 	rc.node.ReportSnapshot(id, status)
// }
