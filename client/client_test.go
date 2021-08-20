package main_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	pb "github.com/binhnt-teko/sharding_admin/schema/accounting"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/stretchr/testify/assert"
)

type RequestData struct {
	Name       string
	Cmd        string
	AccountId  string
	ReceiverId string
	Amount     int64
	Type       pb.Account_TYPE
	State      pb.Account_STATE
}

var API_URL string

func setupClient(t *testing.T) func(t *testing.T) {
	t.Log("setup Client")

	API_URL = "http://localhost:8001"
	return func(t *testing.T) {
		t.Log("teardown connection")
	}
}

func GetEndPoint(cmd string) string {
	switch cmd {
	case "account_create":
		return "/accounts/create"
	case "account_info":
		return "/accounts/info"
	case "account_update":
		return "/accounts/update"
	case "credit":
		return "/credit"
	case "debit":
		return "/debit"
	case "transfer":
		return "/transfer"
	}
	return ""
}
func ConvertRequest(index int, tc RequestData) *pb.Request {
	TraceNo := fmt.Sprintf("Trace_%d", index)
	in := &pb.Request{
		TraceNo: TraceNo,
		ReqTime: 0,
	}
	switch tc.Cmd {
	case "account_create":
		account := &pb.Account{
			AccountId: tc.AccountId,
			Type:      pb.Account_NETWORK,
			Level:     pb.Account_LEVEL1,
			State:     pb.Account_ACTIVE,
		}
		in.Type = pb.RequestType_ACCOUNT_CREATE
		in.Data = &pb.Request_Account{
			Account: account,
		}
		return in
	case "account_info":
		account := &pb.Account{
			AccountId: tc.AccountId,
		}
		in.Type = pb.RequestType_ACCOUNT_INFO
		in.Data = &pb.Request_Account{
			Account: account,
		}
		return in

	case "debit":
		transaction := &pb.Transaction{
			AccountId: tc.AccountId,
			Amount:    tc.Amount,
		}
		in.Type = pb.RequestType_TRANS_DEBIT
		in.Data = &pb.Request_Transaction{
			Transaction: transaction,
		}
		return in
	case "credit":
		transaction := &pb.Transaction{
			AccountId: tc.AccountId,
			Amount:    tc.Amount,
		}
		in.Type = pb.RequestType_TRANS_CREDIT
		in.Data = &pb.Request_Transaction{
			Transaction: transaction,
		}
		return in

	case "transfer":
		transaction := &pb.Transaction{
			AccountId:  tc.AccountId,
			ReceiverId: tc.ReceiverId,
			Amount:     tc.Amount,
		}
		in.Type = pb.RequestType_TRANS_TRANSFER
		in.Data = &pb.Request_Transaction{
			Transaction: transaction,
		}
		return in
	}
	return in
}
func Request(t *testing.T, cmd string, in *pb.Request) (*pb.Response, error) {
	endpoint := GetEndPoint(cmd)
	data, err := protojson.Marshal(in)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", API_URL+endpoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	// fmt.Println("Request: Start send request ")
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println("Request: Try decode json ")
	res := &pb.Response{}
	if err := protojson.Unmarshal(body, res); err != nil {

		return nil, err
	}
	fmt.Printf("Request: Result:  %+v  \n", res)

	return res, nil
}

func CheckOneRequest(t *testing.T, index int, tc RequestData) {
	in := ConvertRequest(index, tc)
	res, err := Request(t, tc.Cmd, in)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	switch tc.Cmd {
	case "account_create":
		switch res.Code {
		case 6:
			msg := fmt.Sprintf("Account %s is already existed.", tc.AccountId)
			assert.Equal(t, res.Message, msg)
			break
		case 0:
			acc := res.GetAccount()
			assert.Equal(t, acc.AccountId, tc.AccountId)
			break
		}
		break

	case "account_info":
		switch res.Code {
		case 0:
			acc := res.GetAccount()
			assert.Equal(t, acc.AccountId, tc.AccountId)
			break
		}

	case "account_update":
		acc := res.GetAccount()
		assert.Equal(t, acc.AccountId, tc.AccountId)
	case "credit":
		switch res.Code {
		case 6:
			msg := fmt.Sprintf("traceNo Trace_%d exists!", index)
			assert.Equal(t, res.Message, msg)
			break
		case 0:
			tx := res.GetTransaction()
			assert.Equal(t, tx.AccountId, tc.AccountId)
			assert.Equal(t, tx.ReceiverId, tc.ReceiverId)
			break
		}

	case "debit":
		switch res.Code {
		case 6:
			msg := fmt.Sprintf("traceNo Trace_%d exists!", index)
			assert.Equal(t, res.Message, msg)
			break
		case 0:
			tx := res.GetTransaction()
			assert.Equal(t, tx.AccountId, tc.AccountId)
			assert.Equal(t, tx.ReceiverId, tc.ReceiverId)
			break
		}
	case "transfer":
		switch res.Code {
		case 2:
			msg := "Cannot enough money"
			assert.Equal(t, res.Message, msg)
			break
		case 6:
			msg := fmt.Sprintf("traceNo Trace_%d exists!", index)
			assert.Equal(t, res.Message, msg)
			break
		case 0:
			tx := res.GetTransaction()
			assert.Equal(t, tx.AccountId, tc.AccountId+"_channel1")
			assert.Equal(t, tx.ReceiverId, tc.ReceiverId)
			break
		}

	}
}
func Test_Transfer_Network(t *testing.T) {
	cases := []RequestData{
		{
			"create account 0001_347",
			"account_create",
			"0001_347",
			"",
			0,
			pb.Account_NETWORK,
			pb.Account_ACTIVE,
		},
		{
			"create account 0001_348",
			"account_create",
			"0001_348",
			"",
			0,
			pb.Account_MEMBER,
			pb.Account_ACTIVE,
		},
		{
			"Credit Account 0001_347",
			"credit",
			"0001_347",
			"",
			1000,
			pb.Account_NO_USE_TYPE,
			pb.Account_ACTIVE,
		},
		{
			"Query account  0001_347",
			"account_info",
			"0001_347",
			"",
			0,
			pb.Account_NO_USE_TYPE,
			pb.Account_ACTIVE,
		},
		{
			"Debit account  0001_347",
			"debit",
			"0001_347",
			"",
			1,
			pb.Account_NO_USE_TYPE,
			pb.Account_ACTIVE,
		},
		{
			"Transfer account  0001_347 to 0001_348",
			"transfer",
			"0001_347",
			"0001_348",
			1,
			pb.Account_NO_USE_TYPE,
			pb.Account_ACTIVE,
		},
		{
			"Get account  0001_348",
			"account_info",
			"0001_348",
			"",
			0,
			pb.Account_NO_USE_TYPE,
			pb.Account_ACTIVE,
		},
	}

	teardownClient := setupClient(t)
	defer teardownClient(t)

	for index, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Logf("Start task %d.%s  ", index+1, tc.Name)
			CheckOneRequest(t, index, tc)
		})
	}
}
