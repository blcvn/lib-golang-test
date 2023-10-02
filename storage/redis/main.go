package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

func main() {
	args := os.Args[1:]

	if len(args) != 4 {
		fmt.Println("Please run as: ./redis-test  master  sentinels  key value ")
		return
	}
	master := args[0]
	sentinelList := args[1]

	key := args[2]
	value := args[3]
	sentinels := strings.Split(sentinelList, ",")
	TestSentinel(master, sentinels, key, value)
}

func TestSingleSetGet() {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}

func TestClusterManual() {
	clusterSlots := func(ctx context.Context) ([]redis.ClusterSlot, error) {
		slots := []redis.ClusterSlot{
			// First node with 1 master and 1 slave.
			{
				Start: 0,
				End:   8191,
				Nodes: []redis.ClusterNode{{
					Addr: ":16379", // master
				}, {
					Addr: ":26379", // 1st slave
				}},
			},
			// Second node with 1 master and 1 slave.
			// {
			// 	Start: 8192,
			// 	End:   16383,
			// 	Nodes: []redis.ClusterNode{{
			// 		Addr: ":26379", // master
			// 	}, {
			// 		Addr: ":16379", // 1st slave
			// 	}},
			// },
		}
		return slots, nil
	}

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		ClusterSlots:  clusterSlots,
		RouteRandomly: true,
	})

	var ctx = context.Background()
	rdb.Ping(ctx)

	// err := rdb.Set(ctx, "key", "value12", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}

func TestSentinel(master string, sentinels []string, key string, value string) {
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "redismaster",
		SentinelAddrs: []string{"172.19.0.7:26379", "172.19.0.6:26379", "172.19.0.4:26379"},
	})
	var ctx = context.Background()
	rdb.Ping(ctx)

	err := rdb.Set(ctx, key, "value9", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf(" key= %s, val = %s \n", key, val)
}
