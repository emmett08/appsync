package main

import (
	"context"
	"fmt"
	"os"

	"github.com/emmett08/dpe-dx-appsync/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	log := logrus.New()

	if err := cmd.RootCmd.ExecuteContext(ctx); err != nil {
		log.WithError(err).Error("command failed")
		_, err := fmt.Fprintln(os.Stderr, "ERROR:", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
