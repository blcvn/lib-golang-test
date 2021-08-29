package main

import (
	"context"
	"fmt"
	"math/rand"
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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// Create the command line application
	app := cli.NewApp()
	app.Name = "client-app-test"
	app.Usage = "implements client direct ClientApp"

	// Describe the commands in the app
	app.Commands = []cli.Command{
		{
			Name:   "create_account_network",
			Usage:  "create account network ",
			Action: createAccount,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
				cli.StringFlag{
					Name:  "branchId",
					Usage: "branchId",
					Value: "03002",
				},
			},
		},
		{
			Name:   "create_account_member",
			Usage:  "create account member ",
			Action: createAccountMember,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
				cli.StringFlag{
					Name:  "branchId",
					Usage: "branchId",
					Value: "03002",
				},
			},
		},
		{
			Name:   "get_account_info",
			Usage:  "get account info ",
			Action: getAccountInfo,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
			},
		},
		{
			Name:   "get_account_info_balance",
			Usage:  "get account info balance",
			Action: getAccountInfoBalance,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
			},
		},
		{
			Name:   "update_account",
			Usage:  "update account",
			Action: updateAccount,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
			},
		},
		{
			Name:   "move_account",
			Usage:  "account move",
			Action: moveAccount,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
				cli.StringFlag{
					Name:  "channelId",
					Usage: "channelId",
					Value: "channel1",
				},
			},
		},
		{
			Name:   "trans_credit",
			Usage:  "transaction credit",
			Action: transCredit,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
				cli.Int64Flag{
					Name:  "amount",
					Usage: "amount",
					Value: 1000,
				},
			},
		},

		{
			Name:   "trans_transfer",
			Usage:  "transaction transfer",
			Action: transTransfer,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "accountId",
					Usage: "accountId",
					Value: "123456",
				},
				cli.StringFlag{
					Name:  "receiveId",
					Usage: "receiveId",
					Value: "123456",
				},
				cli.Int64Flag{
					Name:  "amount",
					Usage: "amount",
					Value: 100,
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
	branchId := c.String("branchId")
	fmt.Println("createAccount:  Create account in generator: ", c.String("accountId"))
	funcName := "/types.AccountService/Acc_Create"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: randSeq(10),
		Type:    pb.RequestType_ACCOUNT_CREATE,
		ReqTime: 0,
		Data: &pb.Request_Account{
			Account: &pb.Account{
				AccountId: accountId,
				Balance:   0,
				State:     pb.Account_ACTIVE,
				Type:      pb.Account_NETWORK,
				Level:     pb.Account_LEVEL1,
				BranchId: branchId,
			},
		},
	}
	ctx := context.Background()
	res, err := InvokeGenerator(ctx, gen, funcName, in)
	fmt.Println("Result: ", res)
	return err
}

func createAccountMember(c *cli.Context) error {
	accountId := c.String("accountId")
	branchId := c.String("branchId")
	fmt.Println("createAccount:  Create account in generator: ", c.String("accountId"))
	funcName := "/types.AccountService/Acc_Create"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: randSeq(10),
		Type:    pb.RequestType_ACCOUNT_CREATE,
		ReqTime: 0,
		Data: &pb.Request_Account{
			Account: &pb.Account{
				AccountId: accountId,
				Balance:   0,
				State:     pb.Account_ACTIVE,
				Type:      pb.Account_MEMBER,
				Level:     pb.Account_LEVEL1,
				BranchId: branchId,
			},
		},
	}
	ctx := context.Background()
	res, err := InvokeGenerator(ctx, gen, funcName, in)
	fmt.Println("Result: ", res)
	return err
}

func getAccountInfo(c *cli.Context) error {
	accountId := c.String("accountId")
	fmt.Println("getAccountInfo:  Get account info in generator: ", c.String("accountId"))
	funcName := "/types.AccountService/Acc_Info"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: randSeq(10),
		Type:    pb.RequestType_ACCOUNT_INFO,
		ReqTime: 0,
		Data: &pb.Request_Account{
			Account: &pb.Account{
				AccountId: accountId,
			},
		},
	}
	ctx := context.Background()
	res, err := InvokeGenerator(ctx, gen, funcName, in)
	fmt.Println("Result: ", res)
	return err
}

func getAccountInfoBalance(c *cli.Context) error {
	accountId := c.String("accountId")
	fmt.Println("getAccountInfoBalance:  Get account info balance in generator: ", c.String("accountId"))
	funcName := "/types.AccountService/Acc_Info_Balance_List"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: randSeq(10),
		Type:    pb.RequestType_ACCOUNT_INFO_BALANCE_LIST,
		ReqTime: 0,
		Data: &pb.Request_Accounts{
			Accounts: &pb.Accounts{
				IdList: []string{accountId},
			},
		},
	}
	ctx := context.Background()
	res, err := InvokeGenerator(ctx, gen, funcName, in)
	fmt.Println("Result: ", res)
	return err
}

func updateAccount(c *cli.Context) error {
	accountId := c.String("accountId")
	fmt.Println("updateAccount:  update account in generator: ", c.String("accountId"))
	funcName := "/types.AccountService/Acc_Update"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: randSeq(10),
		Type:    pb.RequestType_ACCOUNT_UPDATE,
		ReqTime: 0,
		Data: &pb.Request_Account{
			Account: &pb.Account{
				AccountId: accountId,
				State:     pb.Account_INACTIVE,
			},
		},
	}
	ctx := context.Background()
	res, err := InvokeGenerator(ctx, gen, funcName, in)
	fmt.Println("Result: ", res)
	return err
}

func moveAccount(c *cli.Context) error {
	accountId := c.String("accountId")
	channelId := c.String("channelId")
	fmt.Println("updateAccount:  move account in generator: ", c.String("accountId"))
	funcName := "/types.AccountService/Acc_Move"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: randSeq(10),
		Type:    pb.RequestType_ACCOUNT_MOVE,
		ReqTime: 0,
		Data: &pb.Request_Account{
			Account: &pb.Account{
				AccountId: accountId,
				Channel:   channelId,
			},
		},
	}
	ctx := context.Background()
	res, err := InvokeGenerator(ctx, gen, funcName, in)
	fmt.Println("Result: ", res)
	return err
}

func transCredit(c *cli.Context) error {
	accountId := c.String("accountId")
	amount := c.Int64("amount")
	fmt.Println("transCredit:  credit to account in generator: ", c.String("accountId"))
	funcName := "/types.TransactionService/Tran_Credit"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: randSeq(10),
		Type:    pb.RequestType_TRANS_CREDIT,
		ReqTime: 0,
		Data: &pb.Request_Transaction{
			Transaction: &pb.Transaction{
				AccountId: accountId,
				Amount: amount,
			},
		},
	}
	ctx := context.Background()
	res, err := InvokeGenerator(ctx, gen, funcName, in)
	fmt.Println("Result: ", res)
	return err
}

func transTransfer(c *cli.Context) error {
	accountId := c.String("accountId")
	amount := c.Int64("amount")
	receiveId := c.String("receiveId")
	fmt.Println("transTransfer:  transfer from account : ", c.String("accountId"), "to account : ", receiveId, " amount: ", int(amount))
	funcName := "/types.TransactionService/Tran_Transfer"
	gen := &GeneratorInstance{
		Name: "TestClient",
		Url:  "127.0.0.1:5001",
	}

	in := &pb.Request{
		TraceNo: randSeq(10),
		Type:    pb.RequestType_TRANS_TRANSFER,
		ReqTime: 0,
		Data: &pb.Request_Transaction{
			Transaction: &pb.Transaction{
				AccountId: accountId,
				ReceiverId: receiveId,
				Amount: amount,
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
