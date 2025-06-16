package rollout

import (
	"context"
	"github.com/emmett08/1ai-pr/internal/app/plan"
	"runtime"

	"github.com/emmett08/1ai-pr/internal/port"
	"github.com/emmett08/1ai-pr/internal/shared/worker"
)

type UseCase struct {
	repo    port.Repository
	pr      port.PullRequestService
	clock   port.Clock
	planner *plan.Service
}

func NewUseCase(r port.Repository, p port.PullRequestService, c port.Clock) *UseCase {
	return &UseCase{
		repo:    r,
		pr:      p,
		clock:   c,
		planner: plan.New(c),
	}
}

func (u *UseCase) Execute(ctx context.Context, cfg string) error {
	jobs := u.planner.Plan(cfg)
	pool := worker.New(runtime.NumCPU())
	for _, j := range jobs {
		job := j
		pool.Do(func() error { return u.handle(ctx, job) })
	}
	return pool.Wait()
}

func (u *UseCase) handle(ctx context.Context, j plan.Job) error {
	return nil
}
