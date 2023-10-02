package client

import (
	"context"
	"fmt"
	"strings"

	peerpb "github.com/hyperledger/fabric-protos-go/peer"

	"github.com/blcvn/lib-golang-test/consensus/peer/config"
	"github.com/blcvn/lib-golang-test/consensus/peer/flow"
	"github.com/blcvn/lib-golang-test/log/flogging"
)

var (
	fabricClientLogger = flogging.MustGetLogger("client")
)

func Start() error {

	channelID := "test"
	peerCfg := config.FabricConfigPeer{
		PeerName:        "test",
		PeerAddress:     "localhost:8080",
		TLSRootCertFile: "testdata/ca.pem",
		TLSKeyFile:      "testdata/key.pem",
		TLSCertFile:     "testdata/cert.pem",
	}
	signerCfg := config.FabricConfigSigner{
		MSPID:        "SIGNER1",
		IdentityPath: "testdata/cert.pem",
		KeyPath:      "testdata/key.pem",
	}

	deliver, err := flow.NewDeliverClient(context.Background(), channelID, peerCfg, signerCfg)
	if err != nil {
		fabricClientLogger.Errorf("Error while create NewDeliverClient on channel %s: %s", channelID, err.Error())
		return err
	}

	fabricClientLogger.Infof("created listener on channel %s, peer %s", channelID, peerCfg.PeerAddress)

	if deliver == nil {
		fabricClientLogger.Errorf("nil deliver client of peer %s ", peerCfg.PeerAddress)
		return fmt.Errorf("DELIVERY NULL")
	}
	fabricClientLogger.Infof("Start get data from server ")
	for {
		resp, err := deliver.Recv()
		if err != nil {
			fabricClientLogger.Errorf("catch me, received error message from peer %s: %s", peerCfg.PeerAddress, err.Error())

			if strings.Contains(err.Error(), "EOF") ||
				strings.Contains(err.Error(), "transport is closing") ||
				strings.Contains(err.Error(), "connection timed out") {
				return err
			}
			continue
		}
		switch r := resp.Type.(type) {
		case *peerpb.DeliverResponse_Block:
			fabricClientLogger.Debugf("StartListenBlockEvent.DeliverResponse_Block from peer %s: shard %s, blockNumber %d", peerCfg.PeerAddress, channelID, resp.GetBlock().Header.Number)

		case *peerpb.DeliverResponse_Status:
			fabricClientLogger.Errorf("catch me, deliver completed with status (%s) before txid received", r.Status)
		default:
			fabricClientLogger.Errorf("catch me, unexpected response type (%T) from %s", r, peerCfg.PeerAddress)
		}
	}
}
