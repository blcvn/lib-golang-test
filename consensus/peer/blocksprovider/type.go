package blocksprovider

import (
	"context"

	"github.com/blcvn/lib-golang-test/consensus/peer/orderers"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/gossip"
	"github.com/hyperledger/fabric-protos-go/orderer"
	gossipcommon "github.com/hyperledger/fabric/gossip/common"
	"google.golang.org/grpc"
)

type BlockVerifier interface {
	VerifyBlock(channelID gossipcommon.ChannelID, blockNum uint64, block *common.Block) error

	// VerifyBlockAttestation does the same as VerifyBlock, except it assumes block.Data = nil. It therefore does not
	// compute the block.Data.Hash() and compare it to the block.Header.DataHash. This is used when the orderer
	// delivers a block with header & metadata only, as an attestation of block existence.
	VerifyBlockAttestation(channelID string, block *common.Block) error
}

type GossipServiceAdapter interface {
	// AddPayload adds payload to the local state sync buffer
	AddPayload(chainID string, payload *gossip.Payload) error

	// Gossip the message across the peers
	Gossip(msg *gossip.GossipMessage)
}

type LedgerInfo interface {
	// LedgerHeight returns current local ledger height
	LedgerHeight() (uint64, error)
}

type Dialer interface {
	Dial(address string, rootCerts [][]byte) (*grpc.ClientConn, error)
}

type DeliverStreamer interface {
	Deliver(context.Context, *grpc.ClientConn) (orderer.AtomicBroadcast_DeliverClient, error)
}

type OrdererConnectionSource interface {
	RandomEndpoint() (*orderers.Endpoint, error)
	Endpoints() []*orderers.Endpoint
}
