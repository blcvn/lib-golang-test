package main

import (
	"fmt"
	"testing"

	pb "github.com/binhnt-teko/sharding_admin/schema/accounting"
)

var API_URL string

func setupClient(t *testing.T) func(t *testing.T) {
	t.Log("setup Client")

	API_URL = "http://localhost:8000"
	return func(t *testing.T) {
		t.Log("teardown connection")
	}
}
func Test_RateLimit(t *testing.T) {
	teardownClient := setupClient(t)
	defer teardownClient(t)

	maxRequest := 100
	for i := 0; i < maxRequest; i++ {
		t.Log("Send query: ", i)
		name := fmt.Sprintf("createAccount_%d", i)
		accountId := fmt.Sprintf("%d", i)
		req := RequestData{
			Name:       name,
			Cmd:        "account_create",
			AccountId:  accountId,
			ReceiverId: "",
			Amount:     0,
			Type:       pb.Account_MEMBER,
			State:      pb.Account_ACTIVE,
		}
		in := ConvertRequest(i, req)
		res, err := Request(t, req.Cmd, in)
		if err != nil {
			t.Log("Error in request: ", i, err)
		} else {
			t.Log("request: ", i, res)
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
		// {
		// 	"create account 0001_348",
		// 	"account_create",
		// 	"0001_348",
		// 	"",
		// 	0,
		// 	pb.Account_MEMBER,
		// 	pb.Account_ACTIVE,
		// },
		// {
		// 	"Credit Account 0001_347",
		// 	"credit",
		// 	"0001_347",
		// 	"",
		// 	1000,
		// 	pb.Account_NO_USE_TYPE,
		// 	pb.Account_ACTIVE,
		// },
		// {
		// 	"Query account  0001_347",
		// 	"account_info",
		// 	"0001_347",
		// 	"",
		// 	0,
		// 	pb.Account_NO_USE_TYPE,
		// 	pb.Account_ACTIVE,
		// },
		// {
		// 	"Debit account  0001_347",
		// 	"debit",
		// 	"0001_347",
		// 	"",
		// 	1,
		// 	pb.Account_NO_USE_TYPE,
		// 	pb.Account_ACTIVE,
		// },
		// {
		// 	"Transfer account  0001_347 to 0001_348",
		// 	"transfer",
		// 	"0001_347",
		// 	"0001_348",
		// 	1,
		// 	pb.Account_NO_USE_TYPE,
		// 	pb.Account_ACTIVE,
		// },
		// {
		// 	"Get account  0001_348",
		// 	"account_info",
		// 	"0001_348",
		// 	"",
		// 	0,
		// 	pb.Account_NO_USE_TYPE,
		// 	pb.Account_ACTIVE,
		// },
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
