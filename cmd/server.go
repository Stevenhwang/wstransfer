package cmd

import (
	"wstransfer/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start transfer server",

	Run: func(cmd *cobra.Command, args []string) {
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
