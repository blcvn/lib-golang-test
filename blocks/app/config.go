package app

import (
	"github.com/blcvn/lib-golang-test/blocks/consensus/common/cluster"
	"github.com/blcvn/lib-golang-test/blocks/consensus/etcdraft"
)

type Configurator struct {
	etcdraft.Configurator
}

func (s *Configurator) Configure(channel string, newNodes []cluster.RemoteNode) {

}
