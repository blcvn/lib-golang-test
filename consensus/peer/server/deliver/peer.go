package deliver

import "sync"

type Peer struct {
	mutex    sync.RWMutex
	Channels map[string]*Channel
}

func NewPeer() *Peer {
	channels := make(map[string]*Channel, 0)
	return &Peer{
		mutex:    sync.RWMutex{},
		Channels: channels,
	}
}
func (p *Peer) AddChannel(cid string, channel *Channel) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	p.Channels[cid] = channel
}
func (p *Peer) Channel(cid string) *Channel {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if c, ok := p.Channels[cid]; ok {
		return c
	}
	return nil
}
