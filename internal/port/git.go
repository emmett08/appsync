
package port

import "context"

type Repository interface {
    Clone(ctx context.Context, url string, base string, dst string) (Worktree, error)
}

type Worktree interface {
    Commit(msg string) error
    Push(ctx context.Context, branch string) error
}
