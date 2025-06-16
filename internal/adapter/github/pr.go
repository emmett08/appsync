package github

import (
	"context"
	"strings"

	"github.com/emmett08/1ai-pr/internal/port"
	gh "github.com/google/go-github/v60/github"
)

type pr struct{ c *gh.Client }

func New(token string) port.PullRequestService {
	return &pr{c: gh.NewClient(nil).WithAuthToken(token)}
}

func (p *pr) Open(ctx context.Context, repoURL, branch, title, body string) (string, error) {
	orgRepo := strings.TrimSuffix(strings.TrimPrefix(repoURL, "https://github.com/"), ".git")
	org, repo, _ := strings.Cut(orgRepo, "/")
	in := &gh.NewPullRequest{
		Title: &title,
		Head:  &branch,
		Base:  gh.String("main"),
		Body:  &body,
	}
	out, _, err := p.c.PullRequests.Create(ctx, org, repo, in)
	if err != nil {
		return "", err
	}
	return out.GetHTMLURL(), nil
}
