package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"log"

	smallstep "github.com/blcvn/lib-golang-test/security/smallstep/lib"

	"go.step.sm/crypto/pemutil"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type MyChaincodeService struct {
	verifier *smallstep.Verifier
}

type UTXO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             *UID   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Value          int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`                                           //Value of UTXO
	Hash           string `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`                                              //Checking field
	BurnShardId    string `protobuf:"bytes,4,opt,name=burn_shard_id,json=burnShardId,proto3" json:"burn_shard_id,omitempty"`           //UTXO only spend in this shard
	HoldPeriodTime int64  `protobuf:"varint,5,opt,name=hold_period_time,json=holdPeriodTime,proto3" json:"hold_period_time,omitempty"` //Period of time will be hold, after wallet_time + hold_period_time: UTXO can withdraw from wallet and  after nonce+hold_period_time: UTXO can spend in fabric
	WalletTime     int64  `protobuf:"varint,6,opt,name=wallet_time,json=walletTime,proto3" json:"wallet_time,omitempty"`               //Time put UTXO in wallet (used in wallet) != nonce (Created_time in fabric), default wallet_time = nonce
	Txid           string `protobuf:"bytes,7,opt,name=txid,proto3" json:"txid,omitempty"`
	Signature      []byte
	SignerCert     []byte
}

type UID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId string   `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	ShardId   string   `protobuf:"bytes,2,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"` //shard issue this UID
	Nonce     int64    `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Trace     string   `protobuf:"bytes,4,opt,name=trace,proto3" json:"trace,omitempty"` // utxo trace to print to log through modules (accountID-nonce-shardID-burnShardID)
	Type      UID_TYPE `protobuf:"varint,5,opt,name=type,proto3,enum=base.UID_TYPE" json:"type,omitempty"`
}

type UID_TYPE int32

const (
	UID_COMMON    UID_TYPE = 0
	UID_OVERDRAFT UID_TYPE = 1 // overdraft utxo
)

var (
	chaincodeService = &MyChaincodeService{}
)

func (cc *MyChaincodeService) registerRootCA(rootCACertBytes []byte) {
	cc.verifier = smallstep.InitVerifier(rootCACertBytes)
}

func (cc *MyChaincodeService) verify(utxo *UTXO, ownerName string) error {
	if err := cc.verifyCert(utxo.SignerCert, ownerName); err != nil {
		return err
	}

	log.Println("valid cert")
	if err := cc.verifySignature(utxo); err != nil {
		return err
	}
	log.Println("valid signature")

	return nil
}

func (cc *MyChaincodeService) verifyCert(certBytes []byte, ownerName string) error {
	cert, ipems, err := smallstep.CertFromBytes(certBytes)
	if err != nil {
		return err
	}

	if err := cc.verifier.VerifyCert(cert, ipems); err != nil {
		return err

	}

	if err := cc.verifier.VerifyCertSubject(cert, ownerName); err != nil {
		return err
	}

	return nil
}

func (cc *MyChaincodeService) verifySignature(utxo *UTXO) error {
	opts := []pemutil.Options{pemutil.WithFirstBlock()}
	pubKey, err := smallstep.PublicKeyFromCert(utxo.SignerCert, opts)
	if err != nil {
		return err
	}

	if !verifySignature(pubKey, utxo.Hash, utxo.Signature) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

// verifySignature
func verifySignature(pubKey *ecdsa.PublicKey, msg string, sig []byte) bool {
	hash := sha256.Sum256([]byte(msg))

	return ecdsa.VerifyASN1(pubKey, hash[:], sig)
}

func (cc *MyChaincodeService) getMintUTXO() UTXO {
	return UTXO{
		Id: &UID{
			AccountId: "123myacc456",
			ShardId:   "shard1",
			Nonce:     100,
			Trace:     "123mytrace456",
			Type:      UID_COMMON,
		},
		Value:          1000,
		Hash:           "0x123456789",
		BurnShardId:    "shard2",
		HoldPeriodTime: 0,
		WalletTime:     0,
		Txid:           "0x987654321",
	}
}
