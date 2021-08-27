package service

import (
	"context"
	"github.com/TrueGameover/RevoluterraTest/benchmark/domain"
	"time"
)

type WorkerService struct {
	queue          chan domain.IJob
	ctx            context.Context
	resultsChannel chan<- byte
	cancelContext  context.CancelFunc
	startTime      time.Time
}

func (w *WorkerService) Init(workersCount uint, queueSize uint, completeChannel chan<- byte) {
	w.queue = make(chan domain.IJob, queueSize)
	w.resultsChannel = completeChannel
	w.ctx, w.cancelContext = context.WithCancel(context.Background())

	for i := uint(0); i < workersCount; i++ {
		w.runWorker()
	}
}

func (w *WorkerService) Start() {
	w.startTime = time.Now()
}

func (w *WorkerService) runWorker() {
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			default:
			}

			select {
			case <-w.ctx.Done():
				return
			case job, ok := <-w.queue:
				if !ok {
					return
				}

				targetStart := job.GetCreationTime().Add(job.GetWaitTime())

				if targetStart.Before(time.Now()) {
					w.resultsChannel <- job.Do()

				} else {
					w.queue <- job
				}
			}
		}
	}()
}

func (w *WorkerService) Destroy() {
	w.cancelContext()
	close(w.queue)
}

func (w *WorkerService) AddJob(job domain.IJob) {
	w.queue <- job
}
