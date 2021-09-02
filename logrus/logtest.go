package main

import (
	"fmt"

	"github.com/binhnt-teko/lib-golang-test/graylog"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
// var log = logrus.New()

func main() {
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	// log.Out = os.Stdout
	Environment := "production"
	if Environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{})
	}
	// You could set this to any `io.Writer` such as a file
	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	//  log.Out = file
	// } else {
	//  log.Info("Failed to log to file, using default stderr")
	// }

	//1. Test 1
	graylog_ip := "localhost"
	graylog_port := "12201"
	url := fmt.Sprintf("%s:%s", graylog_ip, graylog_port)
	// hook := graylog.NewGraylogHook(url, map[string]interface{}{})
	hook := graylog.NewAsyncGraylogHook(url, map[string]interface{}{})
	log.AddHook(hook)
	log.WithFields(logrus.Fields{
		"animal": "thu nghiem du lieu 1",
		"size":   12312,
		"data":   "Binhnt",
	}).Info("binhnt test")
	hook.Flush()
	//3. Test 3
	// graylogToken := "token-test"
	// hook := grayloghook.NewGraylogHook(url, graylogToken, "example.org", &tls.Config{})

	// log.AddHook(hook)
	// log.WithFields(logrus.Fields{
	// 	"foo": "example",
	// }).Printf("This is an example")
}
