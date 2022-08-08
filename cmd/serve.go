package cmd

import (
	"github.com/spf13/cobra"
	"ratatoskr/internal/models/resolver"
	"ratatoskr/internal/servers"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as daemon",
	Run: func(cmd *cobra.Command, args []string) {
		resolver.BlacklistLoad()
		resolver.LocalLoad()
		servers.DNS()
		waitChan := make(chan int)
		<-waitChan
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags()
}
