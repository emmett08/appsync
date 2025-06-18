package app

import "context"

type RepoGateway interface {
	WriteFile(ctx context.Context, path string, content []byte, branch string) error
	CreateBranch(ctx context.Context, from, to string) error
	PullRequest(ctx context.Context, title, body, base, head string) (int, error)
	DefaultBranch(ctx context.Context) (string, error)
}
