package app

import (
	"github.com/blcvn/lib-golang-test/blocks/consensus/etcdraft"
	pb "go.etcd.io/etcd/raft/v3/raftpb"
)

type AppMemoryStorage struct {
	etcdraft.MemoryStorage
}

/*** Raft Storage is an interface that may be implemented by the application  to retrieve log entries from storage. ***/
// If any Storage method returns an error, the raft instance will
// become inoperable and refuse to participate in elections; the
// application is responsible for cleanup and recovery in this case.
// InitialState returns the saved HardState and ConfState information.
func (s *AppMemoryStorage) InitialState() (pb.HardState, pb.ConfState, error) {
	return pb.HardState{}, pb.ConfState{}, nil
}

// Entries returns a slice of log entries in the range [lo,hi).
// MaxSize limits the total size of the log entries returned, but
// Entries returns at least one entry if any.
func (s *AppMemoryStorage) Entries(lo, hi, maxSize uint64) ([]pb.Entry, error) {
	return []pb.Entry{}, nil
}

// Term returns the term of entry i, which must be in the range
// [FirstIndex()-1, LastIndex()]. The term of the entry before
// FirstIndex is retained for matching purposes even though the
// rest of that entry may not be available.
func (s *AppMemoryStorage) Term(i uint64) (uint64, error) {
	return uint64(0), nil
}

// LastIndex returns the index of the last entry in the log.
func (s *AppMemoryStorage) LastIndex() (uint64, error) {
	return uint64(0), nil
}

// FirstIndex returns the index of the first log entry that is
// possibly available via Entries (older entries have been incorporated
// into the latest Snapshot; if storage only contains the dummy entry the
// first log entry is not available).
func (s *AppMemoryStorage) FirstIndex() (uint64, error) {
	return uint64(0), nil
}

// Snapshot returns the most recent snapshot.
// If snapshot is temporarily unavailable, it should return ErrSnapshotTemporarilyUnavailable,
// so raft state machine could know that Storage needs some time to prepare
// snapshot and call Snapshot later.
func (s *AppMemoryStorage) Snapshot() (pb.Snapshot, error) {
	return pb.Snapshot{}, nil
}

/*
**
// MemoryStorage is currently backed by etcd/raft.MemoryStorage. This interface is
// defined to expose dependencies of fsm so that it may be swapped in the
// future. TODO(jay) Add other necessary methods to this interface once we need
// them in implementation, e.g. ApplySnapshot.
*/
func (s *AppMemoryStorage) Append(entries []pb.Entry) error {
	return nil
}

func (s *AppMemoryStorage) SetHardState(st pb.HardState) error {
	return nil
}
func (s *AppMemoryStorage) CreateSnapshot(i uint64, cs *pb.ConfState, data []byte) (pb.Snapshot, error) {
	return pb.Snapshot{}, nil
}
func (s *AppMemoryStorage) Compact(compactIndex uint64) error {
	return nil

}
func (s *AppMemoryStorage) ApplySnapshot(snap pb.Snapshot) error {
	return nil
}
