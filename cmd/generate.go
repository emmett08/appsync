package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/emmett08/appsync/internal/app"
	"github.com/emmett08/appsync/internal/config"
	"github.com/spf13/cobra"
	yamlv3 "gopkg.in/yaml.v3"
)

func init() {
	var (
		root      string
		dest      string
		reposFile string
		team      string
		appName   string
	)
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate Application CRs into a local sample directory",
		RunE: func(_ *cobra.Command, _ []string) error {
			data, err := os.ReadFile(reposFile)
			if err != nil {
				return fmt.Errorf("read repos file: %w", err)
			}
			var wrapper struct {
				Repos config.RepoConfigs `yaml:"repos"`
			}
			if err := yamlv3.Unmarshal(data, &wrapper); err != nil {
				return fmt.Errorf("unmarshal repos file: %w", err)
			}

			f := config.Filter{Team: team, App: appName}
			scanner := app.CatalogScanner{Root: root, Filter: f}
			descs, err := scanner.Scan(context.Background())
			if err != nil {
				return fmt.Errorf("scan catalogue: %w", err)
			}

			factory := app.CRDFactory{}
			renderer := app.ManifestRenderer{}

			for _, d := range descs {
				owner, repo, ok := wrapper.Repos.ForTeam(d.Team)
				if !ok {
					return fmt.Errorf("no repo mapping for team %s", d.Team)
				}
				repoLoc := owner + "/" + repo

				crds := factory.Create(d, repoLoc)
				appDir := filepath.Join(dest, d.App)
				// appDir := filepath.Join(dest, ".applications", d.Team, d.App)
				files, err := renderer.Render(crds, appDir)
				if err != nil {
					return fmt.Errorf("render app %s: %w", d.App, err)
				}
				for path, content := range files {
					// G306: restrict to 0600
					if err := os.WriteFile(path, content, 0o600); err != nil {
						return fmt.Errorf("writing file %q: %w", path, err)
					}
					fmt.Println("wrote", path)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&root, "root", "", "catalogue root (required)")
	cmd.Flags().StringVar(&dest, "dest", "./sample/appsync", "output directory for generated YAMLs")
	cmd.Flags().StringVar(&reposFile, "repos-file", "", "path to YAML file listing team→owner/repo mappings (required)")
	cmd.Flags().StringVar(&team, "team", "", "filter by team")
	cmd.Flags().StringVar(&appName, "app", "", "filter by application")
	_ = cmd.MarkFlagRequired("root")
	_ = cmd.MarkFlagRequired("repos-file")

	RootCmd.AddCommand(cmd)
}
