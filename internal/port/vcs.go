
package port

import "context"

type PullRequestService interface {
    Open(ctx context.Context, repoURL string, branch string, title string, body string) (string, error)
}
