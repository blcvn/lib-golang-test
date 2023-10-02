package main

import (
	"fmt"
	"os"

	"github.com/blcvn/lib-golang-test/consensus/ledger/fileledger"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var RootCmd = &cobra.Command{
	Use:   "main",
	Short: "",
	Long:  "",
}

func main() {
	RootCmd.AddCommand(fileledger.Cmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
