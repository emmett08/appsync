package rollout_test

import (
	"context"
	"testing"
	"time"

	"github.com/emmett08/1ai-pr/internal/app/rollout"
	"github.com/emmett08/1ai-pr/internal/port"
)

type fakeRepo struct{ cloned bool }

func (f *fakeRepo) Clone(ctx context.Context, url, base, dst string) (port.Worktree, error) {
	f.cloned = true
	return &fakeWT{}, nil
}

type fakeWT struct{}

func (*fakeWT) Commit(string) error                      { return nil }
func (*fakeWT) Push(ctx context.Context, b string) error { return nil }

type fakePR struct{ opened bool }

func (f *fakePR) Open(context.Context, string, string, string, string) (string, error) {
	f.opened = true
	return "https://example/pr/1", nil
}

type fakeClock struct{}

func (fakeClock) Now() time.Time { return time.Now() }

func TestExecute_NoJobs(t *testing.T) {
	uc := rollout.NewUseCase(&fakeRepo{}, &fakePR{}, fakeClock{})
	if err := uc.Execute(context.Background(), ""); err != nil {
		t.Fatalf("execute: %v", err)
	}
}
