package app

import (
	"context"
	"fmt"
	"time"
)

type PRStrategy interface {
	Apply(ctx context.Context, repo RepoGateway, files map[string][]byte) error
}

type DirectCommitStrategy struct{}

func (DirectCommitStrategy) Apply(ctx context.Context, repo RepoGateway, files map[string][]byte) error {
	branch, err := repo.DefaultBranch(ctx)
	if err != nil {
		return err
	}
	for p, b := range files {
		if err := repo.WriteFile(ctx, p, b, branch); err != nil {
			return err
		}
	}
	return nil
}

type FeatureBranchPRStrategy struct{}

func (FeatureBranchPRStrategy) Apply(ctx context.Context, repo RepoGateway, files map[string][]byte) error {
	base, err := repo.DefaultBranch(ctx)
	if err != nil {
		return err
	}
	head := fmt.Sprintf("appsync/%d", time.Now().Unix())
	if err := repo.CreateBranch(ctx, base, head); err != nil {
		return err
	}
	for p, b := range files {
		if err := repo.WriteFile(ctx, p, b, head); err != nil {
			return err
		}
	}
	_, err = repo.PullRequest(ctx, "appsync sync", "", base, head)
	return err
}
