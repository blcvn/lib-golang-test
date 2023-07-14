package app

import (
	"github.com/blcvn/lib-golang-test/blocks/consensus/etcdraft"
	"github.com/blcvn/lib-golang-test/blocks/types/common"
)

// CreateBlockPuller is a function to create BlockPuller on demand.
// It is passed into chain initializer so that tests could mock this.

// BlockPuller is used to pull blocks from other OSN
type AppBlockPuller struct {
	etcdraft.BlockPuller
}

func CreateBlockPuller() (etcdraft.BlockPuller, error) {
	return &AppBlockPuller{}, nil
}

func NewAppBlockPuller() *AppBlockPuller {
	return &AppBlockPuller{}
}
func (s *AppBlockPuller) PullBlock(seq uint64) *common.Block {
	return &common.Block{}
}
func (s *AppBlockPuller) HeightsByEndpoints() (map[string]uint64, error) {
	ret := make(map[string]uint64, 0)
	return ret, nil
}
func (s *AppBlockPuller) Close() {

}
