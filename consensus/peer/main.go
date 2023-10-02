package main

import (
	"fmt"
	"os"

	"github.com/blcvn/lib-golang-test/consensus/peer/client"
	"github.com/blcvn/lib-golang-test/consensus/peer/server"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var RootCmd = &cobra.Command{
	Use:   "main",
	Short: "",
	Long:  "",
}

func main() {
	RootCmd.AddCommand(client.Cmd)
	RootCmd.AddCommand(server.Cmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
