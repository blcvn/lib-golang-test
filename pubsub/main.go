package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	tpb "github.com/blcvn/lib-golang-test/pubsub/types"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var (
	Rdb      *redis.Client
	channels = make(map[string]string, 0)
)

func Init(rdb *redis.Client) {
	Rdb = rdb
}

func RegisterTask(name string, channel string) {
	channels[name] = channel
}

func Task_Publish(ctx context.Context, task tpb.Task) error {
	if channel, ok := channels[task.Path]; ok {
		reqBody, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return Rdb.Publish(ctx, channel, reqBody).Err()
	}
	return fmt.Errorf("CANNOT FIND CHANNEL")
}

func main() {
	// Create a new context
	ctx := context.Background()
	// Create a new Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	Init(rdb)
	RegisterTask("iptables", "mychannel")

	// Subscribe to a channel
	pubsub := Rdb.Subscribe(ctx, "mychannel")
	defer pubsub.Close()

	// Wait for confirmation that we are subscribed to the channel
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	// Start goroutine to receive messages
	ch := pubsub.Channel()

	//Create 3 worker
	for i := 0; i < 3; i++ {
		go func(i int) {
			fmt.Printf("Worker %d listening...\n", i)
			for msg := range ch {
				createTableTask, taskId, err := excuteCreateTable(msg)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Worker %d received origin task: %s\n", i, *createTableTask)
				fmt.Printf("Worker %d received task ID: %s\n", i, taskId)
			}
		}(i)
	}

	//Create 3 task
	for i := 0; i < 3; i++ {
		iptableInfo := tpb.IptableTask{
			Cmd:       "create",
			Interface: "eth1",
			Chain:     "FORWARD",
			Table:     "filter",
			Protocol:  "tcp",
			SPort:     "80",
			DPort:     "8080",
			Action:    "accept",
			Remote:    "172.16.100.0/24",
		}
		data, err := json.Marshal(iptableInfo)
		if err != nil {
			panic(err)
		}
		taskId := uuid.New()
		task := tpb.Task{
			Id:   taskId,
			Path: "iptables",
			Data: data,
		}
		body, err := json.Marshal(task)
		if err != nil {
			panic(err)
		}
		err = rdb.Publish(ctx, "mychannel", body).Err()
		if err != nil {
			panic(err)
		}
	}

	// Block main goroutine
	select {}
}

func excuteCreateTable(msg *redis.Message) (*tpb.IptableTask, uuid.UUID, error) {
	//Unmashel msg to original IptableTask struct
	if !json.Valid([]byte(msg.Payload)) {
		log.Printf("invalid JSON message: %s", msg.Payload)
		return nil, uuid.Nil, nil
	}
	var task tpb.Task
	var iptableTask tpb.IptableTask
	err := json.Unmarshal([]byte(msg.Payload), &task)
	if err != nil {
		return nil, uuid.Nil, err
	}
	err = json.Unmarshal(task.Data, &iptableTask)
	if err != nil {
		return nil, uuid.Nil, err
	}
	//save taskId to redis
	taskId := task.Id.String()
	err = Rdb.Set(context.Background(), taskId, "pending", 0).Err()
	if err != nil {
		return nil, uuid.Nil, err
	}
	return &iptableTask, task.Id, nil

}
