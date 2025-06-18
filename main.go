package main

import (
	"os"

	"github.com/emmett08/appsync/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
