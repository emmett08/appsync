package worker_test

import (
	"errors"
	"testing"

	"github.com/emmett08/1ai-pr/internal/shared/worker"
)

func TestPool(t *testing.T) {
	p := worker.New(4)
	for i := 0; i < 10; i++ {
		n := i
		p.Do(func() error {
			if n%3 == 0 {
				return errors.New("fail")
			}
			return nil
		})
	}
	err := p.Wait()
	if err == nil {
		t.Fatalf("expected aggregated error")
	}
}
