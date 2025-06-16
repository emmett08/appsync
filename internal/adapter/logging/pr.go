package logging

import (
	"context"
	"log"

	"github.com/emmett08/1ai-pr/internal/port"
)

type prLogging struct{ next port.PullRequestService }

func DecoratePR(n port.PullRequestService) port.PullRequestService { return &prLogging{next: n} }

func (p *prLogging) Open(ctx context.Context, repo, branch, title, body string) (string, error) {
	log.Printf("open pr %s branch %s", repo, branch)
	return p.next.Open(ctx, repo, branch, title, body)
}
