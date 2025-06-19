package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emmett08/appsync/internal/app"
	"github.com/emmett08/appsync/internal/config"
	"github.com/spf13/cobra"
	yamlv3 "gopkg.in/yaml.v3"
)

func init() {
	var (
		root      string
		reposFile string
		mode      string
		token     string
	)
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Push local manifests under <root> into each team's repo (opens PRs)",
		RunE: func(_ *cobra.Command, _ []string) error {
			// sanitise reposFile path and forbid upward traversal
			cleanRepos := filepath.Clean(reposFile)
			if strings.Contains(cleanRepos, "..") {
				return fmt.Errorf("invalid repos file path: %q", reposFile)
			}
			//nolint:gosec // path is sanitised above
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

			filesByTeam := map[string]map[string][]byte{}
			err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				// TODO: implicit business rule
				parts := strings.Split(rel, string(os.PathSeparator))
				team := parts[0]
				if _, ok := filesByTeam[team]; !ok {
					filesByTeam[team] = map[string][]byte{}
				}
				// sanitise path and forbid upward traversal
				cleanPath := filepath.Clean(path)
				if strings.Contains(cleanPath, "..") {
					return fmt.Errorf("invalid walk path: %q", path)
				}
				//nolint:gosec // cleaned path from local directory
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				// repoPath := filepath.ToSlash(filepath.Join(".applications", rel))
				// filesByTeam[team][repoPath] = content
				filesByTeam[team][rel] = content
				return nil
			})
			if err != nil {
				return fmt.Errorf("walk root: %w", err)
			}

			for team, files := range filesByTeam {
				owner, repoName, ok := wrapper.Repos.ForTeam(team)
				if !ok {
					return fmt.Errorf("no repository configured for team %s", team)
				}
				gw := app.GitHubGatewayFactory{}.New(token, owner, repoName)

				var strat app.PRStrategy
				if mode == "direct" {
					strat = app.DirectCommitStrategy{}
				} else {
					strat = app.FeatureBranchPRStrategy{}
				}

				if err := strat.Apply(context.Background(), gw, files); err != nil {
					return fmt.Errorf("applying PR strategy for team %s: %w", team, err)
				}
				fmt.Printf("sync complete for team %s → %s/%s\n", team, owner, repoName)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&root, "root", "", "root directory containing generated manifests (required)")
	cmd.Flags().StringVar(&reposFile, "repos-file", "", "path to YAML file listing team→owner/repo mappings (required)")
	cmd.Flags().StringVar(&mode, "mode", "feature", `PR strategy: "direct" | "feature"`)
	cmd.Flags().StringVar(&token, "token", "", "GitHub API token (or set GITHUB_TOKEN)")
	_ = cmd.MarkFlagRequired("root")
	_ = cmd.MarkFlagRequired("repos-file")

	RootCmd.AddCommand(cmd)
}
