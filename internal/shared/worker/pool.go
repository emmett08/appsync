package worker

import (
	"errors"
	"sync"
)

type Pool struct {
	wg   sync.WaitGroup
	sem  chan struct{}
	errs []error
	mu   sync.Mutex
}

func New(concurrency int) *Pool {
	return &Pool{
		sem: make(chan struct{}, concurrency),
	}
}

func (p *Pool) Do(fn func() error) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.sem <- struct{}{}
		if err := fn(); err != nil {
			p.mu.Lock()
			p.errs = append(p.errs, err)
			p.mu.Unlock()
		}
		<-p.sem
	}()
}

func (p *Pool) Wait() error {
	p.wg.Wait()
	if len(p.errs) == 0 {
		return nil
	}
	return errors.Join(p.errs...)
}
