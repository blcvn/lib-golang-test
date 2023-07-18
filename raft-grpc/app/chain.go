package app

import "github.com/blcvn/lib-golang-test/raft-grpc/transport"

type ChainSupport struct {
	transport.MessageReceiver
}

func (s *ChainSupport) ReceiverByChain(channelID string) transport.MessageReceiver {
	return s.MessageReceiver
}
