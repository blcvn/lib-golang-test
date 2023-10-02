package server

import (
	"github.com/blcvn/lib-golang-test/consensus/peer/server/deliver"
	"github.com/hyperledger/fabric/msp"
)

type PolicyCheckerProvider func(resourceName string) deliver.PolicyCheckerFunc

// IdentityDeserializerManager returns instances of Deserializer
type IdentityDeserializerManager interface {
	// Deserializer returns an instance of transaction.Deserializer for the passed channel
	// if the channel exists
	Deserializer(channel string) (msp.IdentityDeserializer, error)
}
