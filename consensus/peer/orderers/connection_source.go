package orderers

import (
	"math/rand"
	"sync"

	"github.com/blcvn/lib-golang-test/log/flogging"
	"github.com/pkg/errors"
)

type ConnectionSource struct {
	mutex              sync.RWMutex
	allEndpoints       []*Endpoint
	orgToEndpointsHash map[string][]byte
	logger             *flogging.FabricLogger
	overrides          map[string]*Endpoint
}

func NewConnectionSource(logger *flogging.FabricLogger, overrides map[string]*Endpoint) *ConnectionSource {
	return &ConnectionSource{
		orgToEndpointsHash: map[string][]byte{},
		logger:             logger,
		overrides:          overrides,
	}
}

func (cs *ConnectionSource) RandomEndpoint() (*Endpoint, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	if len(cs.allEndpoints) == 0 {
		return nil, errors.Errorf("no endpoints currently defined")
	}
	return cs.allEndpoints[rand.Intn(len(cs.allEndpoints))], nil
}

func (cs *ConnectionSource) Endpoints() []*Endpoint {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	return cs.allEndpoints
}
