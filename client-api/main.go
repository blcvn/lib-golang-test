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
	queue := make(chan int, numWorker)
	fmt.Println("Start total of ", numWorker)
	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go func(id int, queue chan int) {
			defer wg.Done()
			fmt.Println("Start worker number ", id)
			for msg := range queue {
				fmt.Println("Send request from worker: ", id, msg)
				accountId := fmt.Sprintf("Account.%d", msg)
				traceId := fmt.Sprintf("Trace.%d", time.Now())
				createAccount(traceId, accountId)
			}

		}(i, queue)
	}

	//Send message job
	wg.Add(1)
	go func() {
		defer wg.Done()
		maxRequest := 100
		for i := 0; i < maxRequest; i++ {
			queue <- i
		}
	}()

	wg.Wait()
	fmt.Println("Main: Completed")
}
func createAccount(trace_no string, account_id string) {

	url := "http://localhost:8000/accounts/create"
	method := "POST"
	// request :=
	// 	`{
	// 	"trace_no": #trace_no,
	// 	"type": "ACCOUNT_CREATE",
	// 	"req_time": 1611037185225,
	// 	"account": {
	// 		"account_id": "#account_id",
	// 		"state": "ACTIVE",
	// 		"type": "NETWORK",
	// 		"level": "LEVEL1",
	// 		"branch_id": "2902"
	// 	}
	// }`
	// request = strings.ReplaceAll(request, "#trace_no", trace_no)
	// request = strings.ReplaceAll(request, "#account_id", account_id)
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
		fmt.Println(err)
		return
	}

	// data := []byte(request)
	// payload := strings.NewReader(request)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer {{jwt_token}}")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
