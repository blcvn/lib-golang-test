package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Start client service",
	Run: func(cmd *cobra.Command, args []string) {
		// port, _ := cmd.Flags().GetInt("port")
		// api.GatewayServer(port)
		fmt.Println("Start client ")
		Start()
	},
}
