package cmd

import (
	"gitea.teneshag.ru/gigabit/goauth/internal/cmd/goauth"

	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE:  goauth.StartCommand,
	}
)

func init() {
	RootCmd.AddCommand(startCmd)
}
