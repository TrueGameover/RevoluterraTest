package domain

import (
	"context"
	"time"
)

type IWorker interface {
	Run(jobsChannel <-chan IJob, ctx context.Context)
}

type IJob interface {
	Do() byte
	GetWaitTime() time.Duration
	GetCreationTime() time.Time
}

type IWorkersProvider interface {
	Init(workersCount uint, creator func() IWorker, completeChannel chan<- IJob)
	Destroy()
	AddJob(job IJob)
}
