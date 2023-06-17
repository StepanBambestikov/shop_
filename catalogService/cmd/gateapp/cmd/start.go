package cmd

import (
	"catalogServiceGit/internal/cmd/gateapp"

	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE:  gateapp.StartCommand,
	}
)

func init() {
	RootCmd.AddCommand(startCmd)
}
