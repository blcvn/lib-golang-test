package deliver

import (
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"
)

type ResponseSender interface {
	// SendStatusResponse sends completion status to the client.
	SendStatusResponse(status cb.Status) error
	// SendBlockResponse sends the block and optionally private data to the client.
	SendBlockResponse(data *cb.Block, channelID string, chain Chain, signedData *protoutil.SignedData) error
	// DataType returns the data type sent by the sender
	DataType() string
}

type PolicyChecker interface {
	CheckPolicy(envelope *cb.Envelope, channelID string) error
}

type Receiver interface {
	Recv() (*cb.Envelope, error)
}

type Server struct {
	Receiver
	PolicyChecker
	ResponseSender
}
