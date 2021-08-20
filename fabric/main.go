package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sharding_accounting/common"
	"sharding_accounting/config"
	"sharding_accounting/server/generator/fabricservice/fabric"
	loyalty_fabric "sharding_accounting/server/generator/fabricservice/fabric"
	"strings"

	pb "github.com/binhnt-teko/sharding_admin/schema/accounting"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc/codes"
)

func main() {

	// Create the command line application
	app := cli.NewApp()
	app.Name = "test_fabric"
	app.Usage = "implements loyalty direct to fabric"

	// Describe the commands in the app
	app.Commands = []cli.Command{
		{
			Name:   "proposal",
			Usage:  "run send proposal and compare result",
			Action: sendProposal,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "channel",
					Usage: "specify the port to listen on",
					Value: "vnpay-channel-1",
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

//Send proposal to peer
func sendProposal(c *cli.Context) error {
	fmt.Println("sendProposal:  Start server with channel: ", c.String("channel"))
	fabricConfigFile := config.GetFabricConfigFile()

	fmt.Println("sendProposal: fabricConfigFile :  ", fabricConfigFile)

	fabricConfig, err := loyalty_fabric.GetConfigFromFile(fabricConfigFile)
	if err != nil {
		log.Fatal("Cannot load fabric config")
	}

	lf := &loyalty_fabric.LoyaltyFabricClient{}
	lf.InitPeerClients(fabricConfig)
	lf.InitSigners(fabricConfig)

	fmt.Println("sendProposal:  Create proposal ")
	// fservice.CreateProposal()
	traceNo := "123"
	fabricType := pb.RequestType_ACCOUNT_INFO
	accountID := "8000"
	channel := "channel2"
	chaincodeName, chaincodeLang, channelID := config.GetChannelInfo(channel)

	args := []string{accountID}
	dataBytes := [][]byte{}

	signedProposal, err := FabricConsensus(lf, traceNo,
		fabricType,
		chaincodeName,
		chaincodeLang,
		channelID,
		args,
		dataBytes...,
	)
	if err != nil {
		fmt.Println("Error in call FabricConsensus err: ", err)
		return err
	}
	txID := signedProposal.TxID
	fmt.Printf("sendProposal: signedProposal.Content[0].Response =  %+v   \n\n", signedProposal.Content[0].Response)

	fmt.Println("sendProposal: message ", signedProposal.Content[0].Response.Status)
	chaincodeResponse, err := DecodeChaincodeResponse(signedProposal.Content[0].Response.Payload, txID)
	if err != nil {
		fmt.Println("Error in delcode response: ", chaincodeResponse)
		return err
	}

	fmt.Printf("sendProposal: Finish proposal: %+v   \n\n", chaincodeResponse)

	return nil
}

///####### Copy funnction ########

func DecodeChaincodeResponse(payload []byte, txID string) (chaincodeResponse *fabric.ChaincodeResponse, err error) {

	buf := bytes.NewBuffer(payload)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&chaincodeResponse); err != nil {
		errMsg := "failed to dec.Decode " + err.Error()
		err = common.ErrorWithTxID(codes.Internal, errMsg, txID)
	}
	fmt.Println("DecodeChaincodeResponse: ", chaincodeResponse)

	return
}

func getChaincodeFunctionName(consensusType pb.RequestType) string {
	cfg := config.GetConfig()
	mapFunction := cfg.GetStringMapString("loyalty_chaincode")
	funcType := strings.ToLower(consensusType.String())
	funcName := mapFunction[funcType]

	return funcName
}

func FabricConsensus(
	lc loyalty_fabric.LoyaltyFabricAPI,
	traceNo string,
	consensusType pb.RequestType,
	chaincodeName string,
	chaincodeLang string,
	channelID string,
	args []string,
	dataBytes ...[]byte,
) (*loyalty_fabric.ProposalResponse, error) {
	var err error
	if config.BypassFabric {
		return nil, nil
	}
	//0. Get type
	fabricType := getChaincodeFunctionName(consensusType)

	// 1. Propose to peers
	signedProposal, err := lc.VerifyDirect(traceNo, fabricType, chaincodeName, chaincodeLang, channelID, args, dataBytes...)
	if err != nil {
		return nil, err
	}
	if len(signedProposal.Content) == 0 {
		return signedProposal, errors.Errorf("Receive no response from peers")
	}
	if len(signedProposal.Content) < config.ConsensusMin {
		return signedProposal, errors.Errorf("Receive only %d response from peers, can't satisfy endorsement policy", len(signedProposal.Content))
	}
	if signedProposal.Error != nil {
		return signedProposal, errors.Errorf("FabricConsensus", traceNo, "Got error while VerifyDirect: %s", nil, nil, signedProposal.Error)
	}

	err = lc.CheckProposalResponses(signedProposal.Content, channelID)
	if err != nil {
		return signedProposal, err
	}

	return signedProposal, nil
}
