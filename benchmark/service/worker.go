package service

import (
	"context"
	"github.com/TrueGameover/RevoluterraTest/benchmark/domain"
	"sync/atomic"
	"time"
)

type WorkerService struct {
	queue          chan domain.IJob
	ctx            context.Context
	resultsChannel chan<- byte
	cancelContext  context.CancelFunc
	startTime      time.Time
	needClose      uint32
	closed         uint32
}

func (w *WorkerService) Init(workersCount uint, queueSize uint, completeChannel chan<- byte) {
	w.queue = make(chan domain.IJob, queueSize)
	w.resultsChannel = completeChannel
	w.ctx, w.cancelContext = context.WithCancel(context.Background())
	w.needClose = 0
	w.closed = 0

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
				now := time.Now()

				if targetStart.After(now) {
					time.Sleep(targetStart.Sub(now))
				}

				w.resultsChannel <- job.Do()
			}
		}
	}()
}

func (w *WorkerService) Destroy(force bool) {
	w.cancelContext()
	atomic.StoreUint32(&w.needClose, 1)

	if force {
		go func() {
			val := atomic.LoadUint32(&w.needClose)

			if val == 1 {
				atomic.StoreUint32(&w.closed, 1)
				close(w.queue)
			}
		}()
	}
}

func (w *WorkerService) AddJob(job domain.IJob, canClose bool) {
	needClose := atomic.LoadUint32(&w.needClose)
	closed := atomic.LoadUint32(&w.closed)

	if needClose == 0 && closed == 0 {
		w.queue <- job

	} else {
		if canClose && closed == 0 {
			atomic.StoreUint32(&w.closed, 1)
			close(w.queue)
		}
	}
}
