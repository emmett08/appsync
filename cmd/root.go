package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "appsync",
	Short: "Synchronise 1AI application skeletons into tenant repos",
}

func Execute() { _ = RootCmd.Execute() }

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().StringP("token", "t", "", "GitHub access token")
	err := RootCmd.MarkPersistentFlagRequired("token")
	if err != nil {
		return
	}
}
