package git

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5/config"

	"github.com/emmett08/1ai-pr/internal/port"

	gitv5 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type repo struct{}

func New() port.Repository { return &repo{} }

func (r *repo) Clone(ctx context.Context, url, base, dst string) (port.Worktree, error) {
	repository, err := gitv5.PlainCloneContext(ctx, dst, false, &gitv5.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(base),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return nil, err
	}
	return &worktree{repository}, nil
}

type worktree struct{ repo *gitv5.Repository }

func (w *worktree) Commit(msg string) error {
	wt, err := w.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = wt.Commit(msg, &gitv5.CommitOptions{All: true})
	return err
}

func (w *worktree) Push(ctx context.Context, branch string) error {
	return w.repo.PushContext(ctx, &gitv5.PushOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/heads/%[1]s:refs/heads/%[1]s", branch)),
		},
	})
}
