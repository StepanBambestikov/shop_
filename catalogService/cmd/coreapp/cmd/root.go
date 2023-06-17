package cmd

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "crudcore",
		Short: "Simple crud core server",
		Long:  "Simple crud core server",
	}
)

func Execute() error {
	return RootCmd.Execute()
}
