package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "appsync",
	Short: "Synchronise 1AI application skeletons into tenant repos",
}
