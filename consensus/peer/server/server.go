package server

import (
	"crypto/x509"
	"fmt"
	"time"

	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/common/ledger/blockledger/fileledger"
	"google.golang.org/grpc"

	"github.com/blcvn/lib-golang-test/consensus/peer/comm"
	"github.com/blcvn/lib-golang-test/consensus/peer/config"
	"github.com/blcvn/lib-golang-test/consensus/peer/metrics"
	"github.com/blcvn/lib-golang-test/consensus/peer/server/deliver"
	"github.com/blcvn/lib-golang-test/log/flogging"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/policies"
	"github.com/hyperledger/fabric/protoutil"
)

// var (
// 	logger = flogging.MustGetLogger("client")
// )

func Start() error {
	listenAddr := "0.0.0.0:8080"
	serverConfig := config.ServerConfig{
		ConnectionTimeout: 0,
		SecOpts: config.SecureOptions{
			VerifyCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				return nil
			},
			Certificate:        []byte{},
			Key:                []byte{},
			ServerRootCAs:      [][]byte{},
			ClientRootCAs:      [][]byte{},
			UseTLS:             false,
			RequireClientCert:  false,
			CipherSuites:       []uint16{},
			TimeShift:          0,
			ServerNameOverride: "",
		},
		KaOpts:             config.KeepaliveOptions{},
		StreamInterceptors: []grpc.StreamServerInterceptor{},
		UnaryInterceptors:  []grpc.UnaryServerInterceptor{},
		Logger:             &flogging.FabricLogger{},
		HealthCheckEnabled: false,
		MaxRecvMsgSize:     0,
		MaxSendMsgSize:     0,
	}
	peerServer, err := comm.NewGRPCServer(listenAddr, serverConfig)
	if err != nil {
		logger.Fatalf("Failed to create peer server (%s)", err)
	}

	policyProviderMap := make(map[int32]policies.Provider)
	channelGroup := &cb.ConfigGroup{
		Policies: map[string]*cb.ConfigPolicy{},
	}
	policyManager, err := policies.NewManagerImpl("", policyProviderMap, channelGroup)
	if err != nil {
		return err
	}
	channelName := "test"
	directory := "data/ledger"

	metricsProvider := &metrics.AppMetricProvider{}
	ledger_factory, err := fileledger.New(directory, metricsProvider)
	if err != nil {
		logger.Fatalf("Failed to create file ledger factory (%s)", err)
		return err
	}
	ledger, err := ledger_factory.GetOrCreate(channelName)
	if err != nil {
		logger.Fatalf("Failed to create file ledger  (%s)", err)
		return err
	}
	chain := &deliver.Channel{
		BlockLedger: ledger,
		PManager:    policyManager,
	}

	go func() {
		currentHeight := ledger.Height()
		it, cnum := ledger.Iterator(&orderer.SeekPosition{
			Type: &orderer.SeekPosition_Specified{
				Specified: &orderer.SeekSpecified{
					Number: currentHeight - 1,
				},
			},
		})
		if cnum != currentHeight-1 {
			logger.Errorf("Start add block to ledger from %d <> %d  \n", currentHeight, cnum)
			return
		}
		blk, _ := it.Next()
		fmt.Printf("Start add block to ledger from %d  \n", currentHeight)
		previousHash := protoutil.BlockHeaderHash(blk.Header)
		for i := 0; i < 1000; i++ {
			time.Sleep(1 * time.Second)
			blkNum := uint64(i) + currentHeight
			logger.Infof("Generate block: %d ", blkNum)
			blk, err := genBlock(blkNum, previousHash)
			if err != nil {
				logger.Errorf("Cannot genBlock ledger: ", err)
				return
			}
			err = ledger.Append(blk)
			if err != nil {
				logger.Infof("Cannot add block to ledger: ", err)
				return
			}
			previousHash = protoutil.BlockHeaderHash(blk.Header)
		}
		currentHeight = ledger.Height()
		fmt.Printf("End add block to ledger from %d  \n", currentHeight)

	}()
	peerInstance := deliver.NewPeer()
	peerInstance.AddChannel(channelName, chain)
	authenticationTimeWindow := time.Duration(1 * time.Second)
	abServer := &DeliverServer{
		DeliverHandler: deliver.NewHandler(
			&deliver.DeliverChainManager{Peer: peerInstance},
			authenticationTimeWindow,
		),
	}
	pb.RegisterDeliverServer(peerServer.Server(), abServer)

	if err = peerServer.Start(); err != nil {
		logger.Fatalf("Failed to create peer server (%s)", err)
	}
	return nil
}
