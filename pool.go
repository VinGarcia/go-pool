package pool

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Pool struct {
	jobsCh      chan func()
	basePool    *errgroup.Group
	currentPool *sync.WaitGroup
	cancelFn    func()
}

func New(ctx context.Context, numWorkers int) *Pool {
	if numWorkers < 1 {
		numWorkers = 1
	}

	jobsCh := make(chan func())
	ctx, cancel := context.WithCancel(ctx)
	basePool, ctx := errgroup.WithContext(ctx)
	currentPool := &sync.WaitGroup{}

	for range numWorkers {
		basePool.Go(func() error {
			for {
				select {
				case job := <-jobsCh:
					processJob(job, currentPool)
				case <-ctx.Done():
					return nil
				}
			}
		})
	}

	return &Pool{
		basePool:    basePool,
		currentPool: currentPool,
		jobsCh:      jobsCh,
		cancelFn:    cancel,
	}
}

func processJob(job func(), wg *sync.WaitGroup) {
	defer wg.Done()
	job()
}

func (p *Pool) Go(job func()) {
	p.currentPool.Add(1)
	p.jobsCh <- job
}

func (p *Pool) Wait() error {
	p.currentPool.Wait()
	return nil
}

func (p *Pool) Close() error {
	p.cancelFn()
	return p.basePool.Wait()
}
