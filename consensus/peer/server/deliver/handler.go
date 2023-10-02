package deliver

import (
	"context"
	"io"
	"math"
	"time"

	cb "github.com/hyperledger/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"

	"github.com/blcvn/lib-golang-test/log/flogging"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("deliver")

type Filtered interface {
	IsFiltered() bool
}

// Handler handles server requests.
type Handler struct {
	ChainManager ChainManager
	TimeWindow   time.Duration
}

// NewHandler creates an implementation of the Handler interface.
func NewHandler(cm ChainManager, timeWindow time.Duration) *Handler {
	return &Handler{
		ChainManager: cm,
		TimeWindow:   timeWindow,
	}
}

// Handle receives incoming deliver requests.
func (h *Handler) Handle(ctx context.Context, srv *Server) error {
	addr := util.ExtractRemoteAddress(ctx)
	logger.Debugf("Starting new deliver loop for %s", addr)
	for {
		logger.Debugf("Attempting to read seek info message from %s", addr)
		envelope, err := srv.Recv()
		if err == io.EOF {
			logger.Debugf("Received EOF from %s, hangup", addr)
			return nil
		}
		if err != nil {
			logger.Warningf("Error reading from %s: %s", addr, err)
			return err
		}

		status, err := h.deliverBlocks(ctx, srv, envelope)
		if err != nil {
			return err
		}

		err = srv.SendStatusResponse(status)
		if status != cb.Status_SUCCESS {
			return err
		}
		if err != nil {
			logger.Warningf("Error sending to %s: %s", addr, err)
			return err
		}

		logger.Debugf("Waiting for new SeekInfo from %s", addr)
	}
}

func isFiltered(srv *Server) bool {
	if filtered, ok := srv.ResponseSender.(Filtered); ok {
		return filtered.IsFiltered()
	}
	return false
}

func (h *Handler) deliverBlocks(ctx context.Context, srv *Server, envelope *cb.Envelope) (status cb.Status, err error) {
	addr := util.ExtractRemoteAddress(ctx)
	payload, chdr, shdr, err := h.parseEnvelope(ctx, envelope)
	if err != nil {
		logger.Warningf("error parsing envelope from %s: %s", addr, err)
		return cb.Status_BAD_REQUEST, nil
	}

	chain := h.ChainManager.GetChain(chdr.ChannelId)
	if chain == nil {
		// Note, we log this at DEBUG because SDKs will poll waiting for channels to be created
		// So we would expect our log to be somewhat flooded with these
		logger.Debugf("Rejecting deliver for %s because channel %s not found", addr, chdr.ChannelId)
		return cb.Status_NOT_FOUND, nil
	}

	seekInfo := &ab.SeekInfo{}
	if err = proto.Unmarshal(payload.Data, seekInfo); err != nil {
		logger.Warningf("[channel: %s] Received a signed deliver request from %s with malformed seekInfo payload: %s", chdr.ChannelId, addr, err)
		return cb.Status_BAD_REQUEST, nil
	}

	erroredChan := chain.Errored()
	if seekInfo.ErrorResponse == ab.SeekInfo_BEST_EFFORT {
		// In a 'best effort' delivery of blocks, we should ignore consenter errors
		// and continue to deliver blocks according to the client's request.
		erroredChan = nil
	}
	select {
	case <-erroredChan:
		logger.Warningf("[channel: %s] Rejecting deliver request for %s because of consenter error", chdr.ChannelId, addr)
		return cb.Status_SERVICE_UNAVAILABLE, nil
	default:
	}

	if seekInfo.Start == nil || seekInfo.Stop == nil {
		logger.Warningf("[channel: %s] Received seekInfo message from %s with missing start or stop %v, %v", chdr.ChannelId, addr, seekInfo.Start, seekInfo.Stop)
		return cb.Status_BAD_REQUEST, nil
	}

	logger.Debugf("[channel: %s] Received seekInfo (%p) %v from %s", chdr.ChannelId, seekInfo, seekInfo, addr)

	cursor, number := chain.Reader().Iterator(seekInfo.Start)
	logger.Debugf("[channel: %s] current block number: %d ", chdr.ChannelId, number)

	defer cursor.Close()
	var stopNum uint64
	switch stop := seekInfo.Stop.Type.(type) {
	case *ab.SeekPosition_Oldest:
		stopNum = number
	case *ab.SeekPosition_Newest:
		// when seeking only the newest block (i.e. starting
		// and stopping at newest), don't reevaluate the ledger
		// height as this can lead to multiple blocks being
		// sent when only one is expected
		if proto.Equal(seekInfo.Start, seekInfo.Stop) {
			stopNum = number
			break
		}
		stopNum = chain.Reader().Height() - 1
	case *ab.SeekPosition_Specified:
		stopNum = stop.Specified.Number
		if stopNum < number {
			logger.Warningf("[channel: %s] Received invalid seekInfo message from %s: start number %d greater than stop number %d", chdr.ChannelId, addr, number, stopNum)
			return cb.Status_BAD_REQUEST, nil
		}
	}

	for {
		if seekInfo.Behavior == ab.SeekInfo_FAIL_IF_NOT_READY {
			if number > chain.Reader().Height()-1 {
				logger.Warningf("[channel: %s] Block %d not found, block number greater than chain length bounds", chdr.ChannelId, number)
				return cb.Status_NOT_FOUND, nil
			}
		}

		var block *cb.Block
		var status cb.Status

		iterCh := make(chan struct{})
		go func() {
			block, status = cursor.Next()
			close(iterCh)
		}()

		select {
		case <-ctx.Done():
			logger.Debugf("Context canceled, aborting wait for next block")
			return cb.Status_INTERNAL_SERVER_ERROR, errors.Wrapf(ctx.Err(), "context finished before block retrieved")
		case <-erroredChan:
			// TODO, today, the only user of the errorChan is the orderer consensus implementations.  If the peer ever reports
			// this error, we will need to update this error message, possibly finding a way to signal what error text to return.
			logger.Warningf("Aborting deliver for request because the backing consensus implementation indicates an error")
			return cb.Status_SERVICE_UNAVAILABLE, nil
		case <-iterCh:
			// Iterator has set the block and status vars
		}

		if status != cb.Status_SUCCESS {
			logger.Warningf("[channel: %s] Error reading from channel, cause was: %v", chdr.ChannelId, status)
			return status, nil
		}

		// increment block number to support FAIL_IF_NOT_READY deliver behavior
		number++

		logger.Debugf("[channel: %s] Delivering block [%d] for (%p) for %s", chdr.ChannelId, block.Header.Number, seekInfo, addr)

		if seekInfo.ContentType == ab.SeekInfo_HEADER_WITH_SIG {
			block.Data = nil
		}

		signedData := &protoutil.SignedData{Data: envelope.Payload, Identity: shdr.Creator, Signature: envelope.Signature}
		if err := srv.SendBlockResponse(block, chdr.ChannelId, chain, signedData); err != nil {
			logger.Warningf("[channel: %s] Error sending to %s: %s", chdr.ChannelId, addr, err)
			return cb.Status_INTERNAL_SERVER_ERROR, err
		}

		if stopNum == block.Header.Number {
			break
		}
	}

	logger.Debugf("[channel: %s] Done delivering to %s for (%p)", chdr.ChannelId, addr, seekInfo)

	return cb.Status_SUCCESS, nil
}

func (h *Handler) parseEnvelope(ctx context.Context, envelope *cb.Envelope) (*cb.Payload, *cb.ChannelHeader, *cb.SignatureHeader, error) {
	payload, err := protoutil.UnmarshalPayload(envelope.Payload)
	if err != nil {
		return nil, nil, nil, err
	}

	if payload.Header == nil {
		return nil, nil, nil, errors.New("envelope has no header")
	}

	chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, nil, nil, err
	}

	shdr, err := protoutil.UnmarshalSignatureHeader(payload.Header.SignatureHeader)
	if err != nil {
		return nil, nil, nil, err
	}

	err = h.validateChannelHeader(ctx, chdr)
	if err != nil {
		return nil, nil, nil, err
	}

	return payload, chdr, shdr, nil
}

func (h *Handler) validateChannelHeader(ctx context.Context, chdr *cb.ChannelHeader) error {
	if chdr.GetTimestamp() == nil {
		err := errors.New("channel header in envelope must contain timestamp")
		return err
	}

	envTime := time.Unix(chdr.GetTimestamp().Seconds, int64(chdr.GetTimestamp().Nanos)).UTC()
	serverTime := time.Now()

	if math.Abs(float64(serverTime.UnixNano()-envTime.UnixNano())) > float64(h.TimeWindow.Nanoseconds()) {
		err := errors.Errorf("envelope timestamp %s is more than %s apart from current server time %s", envTime, h.TimeWindow, serverTime)
		return err
	}

	return nil
}

func noExpiration(_ []byte) time.Time {
	return time.Time{}
}

func (h *Handler) HandleAttestation(ctx context.Context, srv *Server, env *cb.Envelope) error {
	status, err := h.deliverBlocks(ctx, srv, env)
	if err != nil {
		return err
	}

	err = srv.SendStatusResponse(status)
	if status != cb.Status_SUCCESS {
		return err
	}
	if err != nil {
		addr := util.ExtractRemoteAddress(ctx)
		logger.Warningf("Error sending to %s: %s", addr, err)
		return err
	}
	return nil
}
