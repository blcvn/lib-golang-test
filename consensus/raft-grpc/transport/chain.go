package transport

type ChainSupport struct {
	MessageReceiver
}

func (s *ChainSupport) ReceiverByChain(channelID string) MessageReceiver {
	return s.MessageReceiver
}
