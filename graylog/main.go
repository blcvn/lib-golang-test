package main

import (
	"time"

	"github.com/blcvn/lib-golang-test/graylog/gabs"
	gelf "github.com/blcvn/lib-golang-test/graylog/geft"
)

func main() {
	g, err := gabs.NewGraylog(gabs.Endpoint{
		Transport: gabs.UDP,
		Address:   "localhost",
		Port:      12201,
	})
	if err != nil {
		panic(err)
	}

	// Send a message
	err = g.Send(gabs.Message{
		Version:      "1.1",
		Host:         "localhost",
		ShortMessage: "Test from data ",
		FullMessage:  "Stacktrace",
		Timestamp:    time.Now().Unix(),
		Level:        1,
		Extra: map[string]string{
			"MY-EXTRA-FIELD": "extra_value",
		},
	})
	if err != nil {
		panic(err)
	}

	// Close the graylog connection
	if err := g.Close(); err != nil {
		panic(err)
	}
}

func TestGelf() {
	cfg := gelf.Config{
		GraylogPort:     12201,
		GraylogHostname: "127.0.0.1",
		Connection:      "lan",
		MaxChunkSizeWan: 14200,
		MaxChunkSizeLan: 81540,
	}

	g := gelf.New(cfg)

	g.Log(`{
      "version": "1.1",
      "host": "test.vnpay.vn",
      "timestamp": 1356262644,
	  "level": 1,
      "short_message": "Hello From Golang!",
	  "_user_id": 9001,
	  "_some_info": "foo",
	  "_some_env_var": "bar"
  }`)
}
