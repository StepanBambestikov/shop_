package cmd

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "goauth",
		Short: "Gigatest auth server",
		Long:  "Gigatest auth server",
	}
)

func Execute() error {
	return RootCmd.Execute()
}
