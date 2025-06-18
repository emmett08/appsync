package cmd

import (
	"context"
	"os"
	"regexp"

	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

func init() {
	var owner, repo, path, output, token, regexExpr, apiURL, ref string

	cmd := &cobra.Command{
		Use:   "fetch-repos",
		Short: "Generate repos.yaml from a GitHub repo's top-level directories",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
			httpClient := oauth2.NewClient(ctx, ts)

			var gh *github.Client
			if apiURL == "" {
				gh = github.NewClient(httpClient)
			}

			opts := &github.RepositoryContentGetOptions{Ref: ref}
			_, list, _, err := gh.Repositories.GetContents(ctx, owner, repo, path, opts)
			if err != nil {
				return err
			}

			re, err := regexp.Compile(regexExpr)
			if err != nil {
				return err
			}

			type entry struct {
				Team  string `yaml:"team,omitempty"`
				Owner string `yaml:"owner,omitempty"`
				Repo  string `yaml:"repo,omitempty"`
			}
			var out []entry

			for _, v := range list {
				if v.GetType() != "dir" {
					continue
				}
				g := extractGroups(re, v.GetName())
				team := g["team"]
				repoName := g["repo"]
				if team == "" || repoName == "" {
					continue
				}
				out = append(out, entry{
					Team:  team,
					Owner: g["owner"],
					Repo:  repoName,
				})
			}

			data, err := yaml.Marshal(map[string]any{"repos": out})
			if err != nil {
				return err
			}
			return os.WriteFile(output, data, 0644)
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "GitHub API token (required)")
	cmd.Flags().StringVar(&owner, "owner", "", "GitHub owner/org (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "GitHub repository (required)")
	cmd.Flags().StringVar(&path, "path", "", "Path inside the repo")
	cmd.Flags().StringVar(&regexExpr, "regex", `(?P<team>[^_]+)_(?P<owner>[^_]+)_(?P<repo>[^_]+)$`, "Named-capture regex")
	cmd.Flags().StringVar(&output, "output", "repos.yaml", "Output file")
	cmd.Flags().StringVar(&apiURL, "api-url", "", "Override GitHub API base URL")
	cmd.Flags().StringVar(&ref, "ref", "", "Branch name or commit SHA")

	_ = cmd.MarkFlagRequired("token")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")

	RootCmd.AddCommand(cmd)
}
