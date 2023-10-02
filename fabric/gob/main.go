package main

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"google.golang.org/protobuf/runtime/protoimpl"
)

type ChaincodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TxID            string             `protobuf:"bytes,1,opt,name=TxID,proto3" json:"TxID,omitempty"`
	AccountTxHash   string             `protobuf:"bytes,2,opt,name=AccountTxHash,proto3" json:"AccountTxHash,omitempty"`
	ReceiverTxHash  string             `protobuf:"bytes,3,opt,name=ReceiverTxHash,proto3" json:"ReceiverTxHash,omitempty"`
	Message         string             `protobuf:"bytes,4,opt,name=Message,proto3" json:"Message,omitempty"`
	Signer          string             `protobuf:"bytes,5,opt,name=Signer,proto3" json:"Signer,omitempty"`
	BatchJobWrapper []*BatchJobAccount `protobuf:"bytes,6,rep,name=BatchJobWrapper,proto3" json:"BatchJobWrapper,omitempty"`
}

//Define models
type BatchJobAccount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountID     string `protobuf:"bytes,1,opt,name=AccountID,proto3" json:"AccountID,omitempty"`   // for revoke
	ReceiverID    string `protobuf:"bytes,2,opt,name=ReceiverID,proto3" json:"ReceiverID,omitempty"` // for batch transfer
	ErrMsg        string `protobuf:"bytes,3,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	ChaincodeHash string `protobuf:"bytes,4,opt,name=ChaincodeHash,proto3" json:"ChaincodeHash,omitempty"` // will be senderChaincodeHash in revoke case, receiverChaincodeHash in batchTransfer case
}

func main() {
	txID := "ff604c46f9f8269d384fb4853f4958ce2d03158330bb6b640fdfaff4050f1459"
	accountTxHash := "a5d5729b80e9ff458f4ac1115581a99e994ab2d72f991e3a7cdb59ee9189ff45"
	receiverTxHash := ""
	message := ""
	// batchJobWrapper := nil
	channelID := "vnpay-channel-1"
	chaincodeResponse := ChaincodeResponse{
		TxID:            txID,
		AccountTxHash:   accountTxHash,
		ReceiverTxHash:  receiverTxHash,
		Message:         message,
		BatchJobWrapper: nil,
		Signer:          channelID,
	}

	var buf bytes.Buffer

	chaincodeResponseEnc := gob.NewEncoder(&buf)
	if err := chaincodeResponseEnc.Encode(chaincodeResponse); err != nil {
		fmt.Println("error")
	}
	fmt.Printf("chaincodeResponse: %+v \n buf: %+v \n", chaincodeResponse, buf)

}
