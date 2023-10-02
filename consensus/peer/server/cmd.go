package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start server service",
	Run: func(cmd *cobra.Command, args []string) {
		// port, _ := cmd.Flags().GetInt("port")
		fmt.Println("Start server ")
		Start()
	},
}
