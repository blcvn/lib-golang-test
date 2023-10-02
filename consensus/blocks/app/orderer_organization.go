package app

import (
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/msp"
)

// OrdererOrg stores the per org orderer config.
type AppOrdererOrg struct {
	channelconfig.OrdererOrg
}

// Name returns the name this org is referred to in config
func (s *AppOrdererOrg) Name() string {
	return "Test"
}

// MSPID returns the MSP ID associated with this org
func (s *AppOrdererOrg) MSPID() string {
	return "Test"
}

// MSP returns the MSP implementation for this org.
func (s *AppOrdererOrg) MSP() msp.MSP {
	return &AppMSP{}
}

// Endpoints returns the endpoints of orderer nodes.
func (s *AppOrdererOrg) Endpoints() []string {
	return []string{}
}
