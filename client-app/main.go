package main

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/binhnt-teko/loyalty-proto-go/accounting"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

const (
	TIMEOUT_GENERATOR  = 10
	PERIOD_MONITORRING = 10
	PERIOD_RETRY       = 1
	RETRY_MAX          = 3
	IsDisableMonitor   = false
)

type GeneratorInstance struct {
	Name string
	Url  string
}

func main() {
	// Create the command line application
	app := cli.NewApp()
	app.Name = "client-app-test"
	app.Usage = "implements client direct ClientApp"

	// Describe the commands in the app
	app.Commands = []cli.Command{
		{
			Name:   "create_account",
			Usage:  "all account info ",
			Action: createAccount,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
			},
		},
		// {
		// 	Name:   "echo",
		// 	Usage:  "run the sping client",
		// 	Action: startClient,
		// 	Flags: []cli.Flag{
		// 		cli.StringFlag{
		// 			Name:  "n, name",
		// 			Usage: "specify the name of the client",
		// 		},
		// 		cli.UintFlag{
		// 			Name:  "p, port",
		// 			Usage: "specify the port to ping to",
		// 			Value: DefaultPort,
		// 		},
		// 		cli.UintFlag{
		// 			Name:  "l, limit",
		// 			Usage: "specify the max number of pings to send",
		// 			Value: DefaultPings,
		// 		},
		// 		cli.Int64Flag{
		// 			Name:  "d, delay",
		// 			Usage: "the delay between pings in milliseconds",
		// 			Value: DefaultDelay,
		// 		},
		// 	},
		// },
	}

	// Run the application
	app.Run(os.Args)
}
func createAccount(c *cli.Context) error {
	accountId := c.String("accountId")
	fmt.Println("createAccount:  Create account in generator: ", c.String("accountId"))
	funcName := "/types.AccountService/Acc_Create"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: "TraceNo101",
		Type:    pb.RequestType_ACCOUNT_CREATE,
		ReqTime: 0,
		Data: &pb.Request_Account{
			Account: &pb.Account{
				AccountId: accountId,
				Balance:   0,
				State:     0,
				Type:      pb.Account_NETWORK,
				Level:     pb.Account_LEVEL1,
			},
		},
	}
	ctx := context.Background()
	res, err := InvokeGenerator(ctx, gen, funcName, in)
	fmt.Println("Result: ", res)
	return err
}
func InvokeGenerator(ctx context.Context, generator *GeneratorInstance, funcName string, in *pb.Request) (*pb.Response, error) {
	fmt.Printf("grpcclient.InvokeGenerator: start call generator %s \n\n ", generator.Name)
	genName := generator.Name
	conn := NewConnection(generator.Name, generator.Url)
	if conn == nil {
		msg := fmt.Sprintf("GeneratorClient.InvokeGenerator; Cannot get active connection to generator %s ", genName)
		return nil, errors.New(msg)
	}

	opts1 := []grpc.CallOption{}
	out := new(pb.Response)

	err := conn.Invoke(ctx, funcName, in, out, opts1...)
	if err != nil {
		return nil, err
	}

	select {
	case <-time.After(time.Second * 30):
		return nil, errors.Errorf("timeout while InvokeGenerator")
	default:
		return out, nil
	}
}
func NewConnection(Name string, Url string) *grpc.ClientConn {
	fmt.Printf("Endpoint.NewConnection:: Endpoint: %s, url: %s => Recreate connection to \n\n", Name, Url)

	opts := []grpc.DialOption{
		// grpc.WithInsecure(),
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Duration(TIMEOUT_GENERATOR) * time.Second),
	}
	opts = append(opts, grpc.WithInsecure())
	// if TLSEnable || TLSMutualEnable {
	// 	tlsCredentials, err := LoadTLSCredentialsForClient()
	// 	if err != nil {
	// 		log.Fatalf("Endpoint %s, cannot load TLS credentials:  %s", Name, err)
	// 		return nil
	// 	}
	// 	fmt.Printf("Endpoint.NewConnection: Endpoint: %s,  url: %s  =>  Load tlsCredentials: %+v  \n\n", Name, Url, tlsCredentials)
	// 	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	// } else {
	// 	opts = append(opts, grpc.WithInsecure())
	// }

	conn, err := grpc.Dial(Url, opts...)
	if err != nil {
		fmt.Printf("Endpoint.NewConnection: Endpoint: %s, url  %s => Error %s \n", Name, Url, err)
		return nil
	}
	fmt.Printf("Endpoint.NewConnection: Endpoint: %s , url: %s => update connection status \n", Name, Url)

	return conn
}
