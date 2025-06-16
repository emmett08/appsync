package plan

import "github.com/emmett08/1ai-pr/internal/port"

type Service struct {
	clock port.Clock
}

func New(clock port.Clock) *Service { return &Service{clock: clock} }

func (s *Service) Plan(cfg string) []Job {
	// TODO: Parse cfg
	return nil
}
