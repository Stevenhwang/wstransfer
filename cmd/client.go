package cmd

import (
	"wstransfer/client"

	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "start transfer client",

	Run: func(cmd *cobra.Command, args []string) {
		client.Start()
	},
}

func init() {
	RootCmd.AddCommand(clientCmd)
}
