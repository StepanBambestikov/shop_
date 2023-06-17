package cmd

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "sgate",
		Short: "Simple crud gate server",
		Long:  "Simple crud gate server",
	}
)

func Execute() error {
	return RootCmd.Execute()
}
