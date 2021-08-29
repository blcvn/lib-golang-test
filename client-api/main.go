package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	pb "github.com/binhnt-teko/sharding_admin/schema/accounting"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	TestLimitRate()
}
func TestLimitRate() {
	fmt.Println("Start test limit rate: ")

	var wg sync.WaitGroup

	numWorker := 5

	queue := make(chan int)
	fmt.Println("Start total of ", numWorker)
	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go func(id int, queue chan int) {
			fmt.Println("Start worker number ", id)
			wg.Done()
			for msg := range queue {
				fmt.Printf("Worker %d send request %d\n", id, msg)

				// time.Sleep(1 * time.Second)
				accountId := fmt.Sprintf("Account.%d", msg)
				traceId := fmt.Sprintf("Trace.%d", time.Now())
				createAccount(msg, traceId, accountId)
			}

		}(i, queue)
	}
	wg.Wait()

	//Send message job
	maxRequest := 20
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Printf("\nStart send %d requests\n", maxRequest)
		for i := 0; i < maxRequest; i++ {
			// time.Sleep(500 * time.Millisecond)
			queue <- i
		}

		time.Sleep(20 * time.Second)
	}()

	// time.Sleep(20 * time.Second)
	wg.Wait()
	fmt.Println("Main: Completed")
}
func createAccount(requestId int, trace_no string, account_id string) {

	url := "http://localhost:8000/accounts/create"
	method := "POST"
	account := &pb.Account{
		AccountId: account_id,
		Balance:   0,
		State:     pb.Account_ACTIVE,
		Type:      pb.Account_NETWORK,
		Level:     pb.Account_LEVEL1,
		BranchId:  "123",
	}
	in := &pb.Request{
		TraceNo: trace_no,
		Type:    pb.RequestType_ACCOUNT_CREATE,
		ReqTime: 0,
		Data: &pb.Request_Account{
			Account: account,
		},
	}
	data, err := protojson.Marshal(in)
	if err != nil {
		fmt.Println("Request: ", requestId, err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))

	if err != nil {
		fmt.Println("Request: ", requestId, err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer {{jwt_token}}")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Request: ", requestId, err)

		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Request: ", requestId, err)
		return
	}
	fmt.Println("Request: ", requestId, string(body))
}
