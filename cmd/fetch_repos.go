package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	yamlv3 "gopkg.in/yaml.v3"
	"net"
	"net/http"
	urlpkg "net/url"
	"os"
	"regexp"
	"strings"
)

const defaultRegex = `(?P<team>[^_]+)_(?P<owner>[^_]+)_(?P<repo>[^_]+)$`

func init() {
	RootCmd.AddCommand(newFetchCmd())
}

func newFetchCmd() *cobra.Command {
	var (
		org, repo, path, output, token string
		regexExpr                      string
		apiURL, ref                    string
	)

	cmd := &cobra.Command{
		Use:   "fetch-repos",
		Short: "Generate repos.yaml from a GitHub repo’s top-level directories",
		RunE: func(_ *cobra.Command, _ []string) error {
			re := regexp.MustCompile(regexExpr)

			type entry struct {
				Type  string `json:"type"`
				Name  string `json:"name,omitempty"`
				Team  string `json:"team,omitempty"`
				Owner string `json:"owner,omitempty"`
				Repo  string `json:"repo,omitempty"`
			}
			var entries []entry

			if apiURL != "" {
				// validate host to satisfy G107
				u, err := urlpkg.Parse(apiURL)
				if err != nil {
					return fmt.Errorf("invalid api-url: %w", err)
				}
				host := u.Hostname()
				// allow real GitHub or local loopback for tests
				ip := net.ParseIP(host)
				isLoopback := ip != nil && ip.IsLoopback()
				if host != "api.github.com" && !isLoopback {
					return fmt.Errorf("disallowed api-url host: %s", host)
				}
				url := apiURL
				if ref != "" {
					sep := "?"
					if strings.Contains(url, "?") {
						sep = "&"
					}
					url += sep + "ref=" + ref
				}

				resp, err := http.Get(url)
				if err != nil {
					return fmt.Errorf("fetch listing: %w", err)
				}
				defer func() {
					if cErr := resp.Body.Close(); cErr != nil {
						_, err := fmt.Fprintf(os.Stderr, "warning: error closing body: %v\n", cErr)
						if err != nil {
							return
						}
					}
				}()
				if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
					return fmt.Errorf("decode JSON: %w", err)
				}
			} else {
				ctx := context.Background()
				tc := oauth2.NewClient(ctx,
					oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
				cli := github.NewClient(tc)

				opt := &github.RepositoryContentGetOptions{}
				if ref != "" {
					opt.Ref = ref
				}

				_, dir, _, err := cli.Repositories.GetContents(ctx, org, repo, path, opt)
				if err != nil {
					return fmt.Errorf("listing content: %w", err)
				}
				for _, d := range dir {
					entries = append(entries, entry{
						Type: d.GetType(),
						Name: d.GetName(),
					})
				}
			}

			var repos []map[string]string
			for _, e := range entries {
				if e.Type != "dir" {
					continue
				}
				if e.Team != "" && e.Owner != "" && e.Repo != "" {
					repos = append(repos, map[string]string{
						"team":  e.Team,
						"owner": e.Owner,
						"repo":  e.Repo,
					})
					continue
				}
				groups := extractGroups(re, e.Name)
				team, okT := groups["team"]
				repoName, okR := groups["repo"]
				if !okT || !okR {
					continue
				}
				item := map[string]string{
					"team": team,
					"repo": repoName,
				}
				if owner := groups["owner"]; owner != "" {
					item["owner"] = owner
				}
				repos = append(repos, item)
			}

			out := map[string]any{"repos": repos}
			b, err := yamlv3.Marshal(out)
			if err != nil {
				return err
			}
			// G306: use 0600 so that gosec is happy
			return os.WriteFile(output, b, 0o600)
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "GitHub token (required)")
	cmd.Flags().StringVar(&org, "owner", "", "GitHub owner/org (required)")
	cmd.Flags().StringVar(&repo, "repo", "", "GitHub repo (required)")
	cmd.Flags().StringVar(&path, "path", "", "Path inside the repo")
	cmd.Flags().StringVar(&regexExpr, "regex", defaultRegex, "Named-capture regex")
	cmd.Flags().StringVar(&output, "output", "repos.yaml", "Output file")
	cmd.Flags().StringVar(&apiURL, "api-url", "", "Override API base-URL (for tests)")
	cmd.Flags().StringVar(&ref, "ref", "", "Branch or commit SHA")

	_ = cmd.MarkFlagRequired("token")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")

	return cmd
}
