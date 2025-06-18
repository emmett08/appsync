package cmd

import (
	"context"

	"github.com/emmett08/appsync/internal/app"
	"github.com/emmett08/appsync/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	var (
		root    string
		dest    string
		owner   string
		repo    string
		mode    string
		team    string
		appName string
	)
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Copy skeleton CRDs and raise PRs",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, _ := cmd.Flags().GetString("token")
			f := config.Filter{Team: team, App: appName}

			scanner := app.CatalogScanner{Root: root, Filter: f}
			factory := app.CRDFactory{}
			renderer := app.ManifestRenderer{}
			gateway := app.NewGitHubGateway(token, owner, repo)

			var strat app.PRStrategy
			if mode == "direct" {
				strat = app.DirectCommitStrategy{}
			} else {
				strat = app.FeatureBranchPRStrategy{}
			}

			coord := app.SyncCoordinator{
				Scanner:    scanner,
				Factory:    factory,
				Renderer:   renderer,
				Gateway:    gateway,
				PRStrategy: strat,
				TargetRoot: dest,
			}
			return coord.Sync(context.Background())
		},
	}
	cmd.Flags().StringVar(&root, "root", "", "catalog root")
	cmd.Flags().StringVar(&dest, "dest", "/tmp/appsync", "working dir")
	cmd.Flags().StringVar(&owner, "owner", "", "GitHub owner")
	cmd.Flags().StringVar(&repo, "repo", "", "GitHub repo")
	cmd.Flags().StringVar(&mode, "mode", "pr", "direct | pr")
	cmd.Flags().StringVar(&team, "team", "", "filter by team")
	cmd.Flags().StringVar(&appName, "app", "", "filter by application")
	err := cmd.MarkFlagRequired("root")
	if err != nil {
		return
	}
	err = cmd.MarkFlagRequired("owner")
	if err != nil {
		return
	}
	err = cmd.MarkFlagRequired("repo")
	if err != nil {
		return
	}

	RootCmd.AddCommand(cmd)
}
