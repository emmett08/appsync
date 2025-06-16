package logging

import (
	"context"
	"log"

	"github.com/emmett08/1ai-pr/internal/port"
)

type repoLogging struct{ next port.Repository }

func DecorateRepo(n port.Repository) port.Repository { return &repoLogging{next: n} }

func (r *repoLogging) Clone(ctx context.Context, url, base, dst string) (port.Worktree, error) {
	log.Printf("clone %s base %s", url, base)
	return r.next.Clone(ctx, url, base, dst)
}
