package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
// 	clientv3 "go.etcd.io/etcd/client/v3"
// )

// func main() {
// 	cli, err := clientv3.New(clientv3.Config{
// 		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
// 		DialTimeout: 5 * time.Second,
// 	})
// 	if err != nil {
// 		// handle error!
// 	}

// 	requestTimeout := time.Duration(1)
// 	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
// 	_, err = cli.Put(ctx, "", "sample_value")
// 	cancel()
// 	if err != nil {
// 		switch err {
// 		case context.Canceled:
// 			fmt.Printf("ctx is canceled by another routine: %v\n", err)
// 		case context.DeadlineExceeded:
// 			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
// 		case rpctypes.ErrEmptyKey:
// 			fmt.Printf("client-side error: %v\n", err)
// 		default:
// 			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
// 		}
// 	}

// 	defer cli.Close()
// }
