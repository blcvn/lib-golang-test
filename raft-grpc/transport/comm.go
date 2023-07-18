package transport

import (
	"context"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// Communicator defines communication for a consenter
type Communicator interface {
	// Remote returns a RemoteContext for the given RemoteNode ID in the context
	// of the given channel, or error if connection cannot be established, or
	// the channel wasn't configured
	Remote(channel string, id uint64) (*RemoteContext, error)
	// Configure configures the communication to connect to all
	// given members, and disconnect from any members not among the given
	// members.
	// Configure(channel string, members []RemoteNode)
	// // Shutdown shuts down the communicator
	// Shutdown()
}

type CommClientStream struct {
	StepClient orderer.Cluster_StepClient
}

func (cs *CommClientStream) Send(request *orderer.StepRequest) error {
	return cs.StepClient.Send(request)
}

func (cs *CommClientStream) Recv() (*orderer.StepResponse, error) {
	return cs.StepClient.Recv()
}

func (cs *CommClientStream) Auth() error {
	return nil
}

func (cs *CommClientStream) Context() context.Context {
	return cs.StepClient.Context()
}

// Implement Communicator, ServerDispatch
type Comm struct {
	ShutdownSignal chan struct{}
	Shutdown       bool
	SendBufferSize int
	Lock           sync.RWMutex
	Connections    *ConnectionStore
	Chan2Members   map[string]MemberMapping
}

// Remote obtains a RemoteContext linked to the destination node on the context
// of a given channel
func (c *Comm) Remote(channel string, id uint64) (*RemoteContext, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	fmt.Printf("Comm.Remote: get Sub of %d \n", id)
	// if c.Shutdown {
	// 	return nil, errors.New("communication has been shut down")
	// }

	mapping, exists := c.Chan2Members[channel]
	if !exists {
		return nil, errors.Errorf("channel %s doesn't exist", channel)
	}
	stub := mapping.ByID(id)
	if stub == nil {
		return nil, errors.Errorf("node %d doesn't exist in channel %s's membership", id, channel)
	}

	//Check stub if active or not
	if stub.Active() {
		fmt.Printf("Comm.Remote: Stub of %d active \n", id)
		return stub.RemoteContext, nil
	}

	//Create connection in order to active
	fmt.Printf("Comm.Remote: call stub.Activate \n")

	err := stub.Activate(c.createRemoteContext(stub, channel))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return stub.RemoteContext, nil
}

// createRemoteStub returns a function that creates a RemoteContext.
// It is used as a parameter to Stub.Activate() in order to activate
// a stub atomically.
func (c *Comm) createRemoteContext(stub *Stub, channel string) func() (*RemoteContext, error) {
	fmt.Printf("Comm.createRemoteContext: return createRemteContext Function \n")
	return func() (*RemoteContext, error) {
		fmt.Printf("createRemteContext: call Connections.Connection \n")

		conn, err := c.Connections.Connection(stub.Endpoint, []byte(stub.Endpoint))
		if err != nil {
			fmt.Printf("createRemteContext: Unable to obtain connection to %d(%s) (channel %s): %v", stub.ID, stub.Endpoint, channel, err)
			return nil, err
		}

		probeConnection := func(conn *grpc.ClientConn) error {
			connState := conn.GetState()
			if connState == connectivity.Connecting {
				return errors.Errorf("connection to %d(%s) is in state %s \n", stub.ID, stub.Endpoint, connState)
			}
			return nil
		}

		clusterClient := orderer.NewClusterClient(conn)

		getStream := func(ctx context.Context) (StepClientStream, error) {
			fmt.Printf("getStream: call clusterClient.Step \n")
			stream, err := clusterClient.Step(ctx)
			if err != nil {
				fmt.Printf("getStream: clusterClient.Step failed %s \n", err)
				return nil, err
			}
			stepClientStream := &CommClientStream{
				StepClient: stream,
			}
			return stepClientStream, nil
		}

		rc := &RemoteContext{
			Channel:        channel,
			endpoint:       stub.Endpoint,
			conn:           conn,
			GetStreamFunc:  getStream,
			SendBuffSize:   c.SendBufferSize,
			shutdownSignal: make(chan struct{}),
			ProbeConn:      probeConnection,
		}
		return rc, nil
	}
}
