package cmd

import (
	"catalogServiceGit/internal/cmd/coreapp"

	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE:  coreapp.StartCommand,
	}
)

func init() {
	RootCmd.AddCommand(startCmd)
}
