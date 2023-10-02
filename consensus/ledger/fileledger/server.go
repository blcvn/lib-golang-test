package fileledger

import (
	"github.com/blcvn/lib-golang-test/consensus/ledger/metrics"
	"github.com/blcvn/lib-golang-test/consensus/ledger/txflags"
	"github.com/blcvn/lib-golang-test/log/flogging"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/ledger/blockledger/fileledger"
	"github.com/hyperledger/fabric/protoutil"
)

var logger = flogging.MustGetLogger("common.deliverevents")

func Start() error {
	directory := "data/ledger"
	channelName := "test"

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

	previousHash := []byte{}
	for i := 0; i < 1000; i++ {
		logger.Infof("Generate block: %d ", i)
		blk, err := genBlock(uint64(i), previousHash)
		if err != nil {
			logger.Errorf("Cannot genBlock ledger: ", err)
			return err
		}
		err = ledger.Append(blk)
		if err != nil {
			logger.Infof("Cannot add block to ledger: ", err)
			return err
		}
		previousHash = protoutil.BlockHeaderHash(blk.Header)
	}

	return nil
}

func genBlock(number uint64, previousHash []byte) (*cb.Block, error) {
	transactions := [][]byte{}
	for i := 0; i < 100; i++ {
		tx, err := getTransaction()
		if err != nil {
			logger.Errorf("Cannot generate transaction")
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	blkData := &cb.BlockData{
		Data: transactions,
	}

	header := &cb.BlockHeader{
		Number:       number,
		PreviousHash: previousHash,
		DataHash:     protoutil.BlockDataHash(blkData),
	}

	block := &cb.Block{
		Header: header,
		Data:   blkData,
	}

	protoutil.InitBlockMetadata(block)
	txsFilter := txflags.NewWithValues(len(block.Data.Data), peer.TxValidationCode_VALID)
	block.Metadata.Metadata[cb.BlockMetadataIndex_TRANSACTIONS_FILTER] = txsFilter
	return block, nil
}

func getProposalWithType(ccID string, ccVersion string, pType cb.HeaderType, signerSerialized []byte) (*peer.Proposal, error) {
	cis := &peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{Name: ccID, Version: ccVersion},
			Input:       &peer.ChaincodeInput{Args: [][]byte{[]byte("func")}},
			Type:        peer.ChaincodeSpec_GOLANG,
		},
	}

	proposal, _, err := protoutil.CreateProposalFromCIS(pType, "testchannelid", cis, signerSerialized)
	return proposal, err
}

func getTransaction() ([]byte, error) {
	ccID := ""
	ccVersion := ""
	pType := cb.HeaderType_ENDORSER_TRANSACTION
	event := []byte{}
	res := []byte{}
	tx, err := genEnvelop(ccID, ccVersion, pType, event, res)

	if err != nil {
		logger.Errorf("Canootn gen block: ", err)
		return nil, err
	}
	data := protoutil.MarshalOrPanic(tx)
	return data, err
}

func genEnvelop(ccID string, ccVersion string, pType cb.HeaderType, event []byte, res []byte) (*cb.Envelope, error) {
	response := &peer.Response{Status: 200}

	signer := NewSigner()

	signerSerialized, err := signer.Serialize()
	if err != nil {
		logger.Errorf("Could not serialize identity")
		return nil, err
	}

	proposal, err := getProposalWithType(ccID, ccVersion, pType, signerSerialized)
	// endorse it to get a proposal response
	presp, err := protoutil.CreateProposalResponse(proposal.Header, proposal.Payload, response, res, event, &peer.ChaincodeID{Name: ccID, Version: ccVersion}, signer)

	// assemble a transaction from that proposal and endorsement
	tx, err := protoutil.CreateSignedTx(proposal, signer, presp)
	return tx, err
}
