package flow

import (
	"context"
	"crypto/tls"

	"github.com/blcvn/lib-golang-test/consensus/peer/config"
	"github.com/blcvn/lib-golang-test/consensus/peer/fabric"
	"github.com/blcvn/lib-golang-test/log/flogging"
	"github.com/hyperledger/fabric/cmd/common/signer"

	pb "github.com/hyperledger/fabric-protos-go/peer"
)

var listenerLogger = flogging.MustGetLogger("nhs.lib.flow.listener")

func NewDeliverClient(
	ctx context.Context,
	channelID string,
	peer config.FabricConfigPeer,
	configSigner config.FabricConfigSigner,
) (pb.Deliver_DeliverFilteredClient, error) {

	deliverClient, err := fabric.NewDeliverClient(peer)
	if err != nil {
		listenerLogger.Errorf("PeerDeliver failed %s", err.Error())
		return nil, err
	}
	certificate, err := tls.LoadX509KeyPair(peer.TLSCertFile, peer.TLSKeyFile)
	if err != nil {
		listenerLogger.Errorf("LoadX509KeyPair failed %s", err.Error())
		return nil, err
	}

	//2. Send delivery filter to peer
	signerConfig := signer.Config{
		MSPID:        configSigner.MSPID,
		IdentityPath: configSigner.IdentityPath,
		KeyPath:      configSigner.KeyPath,
	}

	signer, err := signer.NewSigner(signerConfig)
	if err != nil {
		listenerLogger.Errorf("NewSigner failed %s", err.Error())
		return nil, err
	}
	envelope := fabric.CreateDeliverEnvelope(channelID, certificate, signer)

	deliver, err := deliverClient.Deliver(ctx)
	if err != nil {
		listenerLogger.Errorf("DeliverFiltered failed %s: error: ", err.Error())
		return nil, err
	}

	err = deliver.Send(envelope)
	if err != nil {
		listenerLogger.Errorf("CreateDeliverEnvelope failed %s", err.Error())
	}
	defer deliver.CloseSend()
	return deliver, nil
}
