package fileledger

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "gen",
	Short: "Start gen transactions",
	Run: func(cmd *cobra.Command, args []string) {
		// port, _ := cmd.Flags().GetInt("port")
		fmt.Println("Start gen transactions ")
		Start()
	},
}
