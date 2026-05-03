package ledger

import (
	"context"
	"sync"

	"github.com/rphmauriciodev/gopher-ledger/internal/domain"
)

type WorkerPool struct {
	service *LedgerService
	queue   <-chan domain.Transaction
	workers int
}

func NewWorkerPool(s *LedgerService, queue <-chan domain.Transaction, workers int) *WorkerPool {
	return &WorkerPool{
		service: s,
		queue:   queue,
		workers: workers,
	}
}

func (w *WorkerPool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < w.workers; i++ {
		wg.Add(1)
		go w.worker(i, &wg, ctx)
	}

	wg.Wait()
}

func (w *WorkerPool) worker(_ int, wg *sync.WaitGroup, ctx context.Context) {
	wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case tx, ok := <-w.queue:
			if !ok {
				return
			}

			if err := w.service.ExecuteTransaction(ctx, tx); err != nil {
				continue
			}
		}
	}
}
