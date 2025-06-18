package app

import (
	"context"
	"fmt"
	"path"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

type GitHubGateway struct {
	client *github.Client
	owner  string
	repo   string
}

func NewGitHubGateway(token, owner, repo string) *GitHubGateway {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	client := github.NewClient(oauth2.NewClient(context.Background(), ts))
	return &GitHubGateway{client: client, owner: owner, repo: repo}
}

func (g *GitHubGateway) DefaultBranch(ctx context.Context) (string, error) {
	r, _, err := g.client.Repositories.Get(ctx, g.owner, g.repo)
	if err != nil {
		return "", err
	}
	if r.DefaultBranch == nil {
		return "main", nil
	}
	return *r.DefaultBranch, nil
}

func (g *GitHubGateway) CreateBranch(ctx context.Context, from, to string) error {
	refName := fmt.Sprintf("refs/heads/%s", from)
	rf, _, err := g.client.Git.GetRef(ctx, g.owner, g.repo, refName)
	if err != nil {
		return err
	}
	newRef := &github.Reference{
		Ref:    github.String(fmt.Sprintf("refs/heads/%s", to)),
		Object: &github.GitObject{SHA: rf.Object.SHA},
	}
	_, _, err = g.client.Git.CreateRef(ctx, g.owner, g.repo, newRef)
	return err
}

func (g *GitHubGateway) WriteFile(ctx context.Context, filePath string, content []byte, branch string) error {
	opts := &github.RepositoryContentGetOptions{Ref: branch}
	filePath = path.Clean(filePath)
	existing, _, resp, err := g.client.Repositories.GetContents(ctx, g.owner, g.repo, filePath, opts)
	msg := fmt.Sprintf("appsync: update %s", filePath)
	if err != nil && resp.StatusCode != 404 {
		return err
	}
	options := &github.RepositoryContentFileOptions{
		Message: github.String(msg),
		Content: content,
		Branch:  github.String(branch),
	}
	if existing != nil {
		options.SHA = existing.SHA
		_, _, err = g.client.Repositories.UpdateFile(ctx, g.owner, g.repo, filePath, options)
	} else {
		_, _, err = g.client.Repositories.CreateFile(ctx, g.owner, g.repo, filePath, options)
	}
	return err
}

func (g *GitHubGateway) PullRequest(ctx context.Context, title, body, base, head string) (int, error) {
	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(head),
		Base:                github.String(base),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}
	r, _, err := g.client.PullRequests.Create(ctx, g.owner, g.repo, newPR)
	if err != nil {
		return 0, err
	}
	return r.GetNumber(), nil
}
