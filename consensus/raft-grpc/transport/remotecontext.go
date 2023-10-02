/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transport

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type streamsMapperReporter struct {
	size uint32
	sync.Map
}

func (smr *streamsMapperReporter) Delete(key interface{}) {
	smr.Map.Delete(key)
	atomic.AddUint32(&smr.size, ^uint32(0))
}

func (smr *streamsMapperReporter) Store(key, value interface{}) {
	smr.Map.Store(key, value)
	atomic.AddUint32(&smr.size, 1)
}

// RemoteContext interacts with remote cluster
// nodes. Every call can be aborted via call to Abort()
type RemoteContext struct {
	Channel       string
	endpoint      string
	conn          *grpc.ClientConn
	GetStreamFunc func(context.Context) (StepClientStream, error) // interface{}

	// expiresAt                        time.Time
	// minimumExpirationWarningInterval time.Duration
	// certExpWarningThreshold          time.Duration
	SendBuffSize   int
	shutdownSignal chan struct{}
	// Logger         *flogging.FabricLogger
	ProbeConn    func(conn *grpc.ClientConn) error
	nextStreamID uint64
}

// NewStream creates a new stream.
// It is not thread safe, and Send() or Recv() block only until the timeout expires.
func (rc *RemoteContext) NewStream(timeout time.Duration) (*Stream, error) {
	if rc.ProbeConn != nil {
		if err := rc.ProbeConn(rc.conn); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())

	if rc.GetStreamFunc == nil {
		fmt.Printf("RemoteContext.NewStream: GetStreamFunc empty \n")
		return nil, fmt.Errorf("GetStreamFunc empty")
	}
	stream, err := rc.GetStreamFunc(ctx)
	if err != nil {
		fmt.Printf("RemoteContext.NewStream: GetStreamFunc failed \n")
		cancel()
		return nil, errors.WithStack(err)
	}

	// nodeName := commonNameFromContext(stream.Context())
	streamID := atomic.AddUint64(&rc.nextStreamID, 1)
	nodeName := "123"

	var canceled uint32

	abortChan := make(chan struct{})
	abortReason := &atomic.Value{}

	once := &sync.Once{}

	cancelWithReason := func(err error) {
		once.Do(func() {
			abortReason.Store(err.Error())
			cancel()
			// rc.streamsByID.Delete(streamID)
			// rc.Metrics.reportEgressStreamCount(rc.Channel, atomic.LoadUint32(&rc.streamsByID.size))
			// rc.Logger.Debugf("Stream %d to %s(%s) is aborted", streamID, nodeName, rc.endpoint)
			atomic.StoreUint32(&canceled, 1)
			close(abortChan)
		})
	}

	logger := flogging.MustGetLogger("orderer.common.cluster.step")
	stepLogger := logger.WithOptions(zap.AddCallerSkip(1))

	s := &Stream{
		Channel:     rc.Channel,
		abortReason: abortReason,
		abortChan:   abortChan,
		sendBuff: make(chan struct {
			request *orderer.StepRequest
			report  func(error)
		}, rc.SendBuffSize),
		commShutdown: rc.shutdownSignal,
		NodeName:     nodeName,
		Logger:       stepLogger,
		ID:           streamID,
		Endpoint:     rc.endpoint,
		Timeout:      timeout,
		StepClient:   stream,
		Cancel:       cancelWithReason,
		canceled:     &canceled,
	}

	err = stream.Auth()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new stream")
	}

	fmt.Printf("RemoteContext.NewStream: Created new stream to %s with ID of %d and buffer size of %d \n",
		rc.endpoint, streamID, cap(s.sendBuff))

	go func() {
		// rc.workerCountReporter.increment(s.metrics)
		s.serviceStream()
		// rc.workerCountReporter.decrement(s.metrics)
	}()

	return s, nil
}
