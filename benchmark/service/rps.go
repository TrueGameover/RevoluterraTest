package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/TrueGameover/RevoluterraTest/benchmark/domain"
	"github.com/TrueGameover/RevoluterraTest/benchmark/repository"
	client2 "github.com/TrueGameover/RevoluterraTest/client"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type RpsService struct {
	Hosts       *sync.Map
	SiteTracker repository.ISiteRequesterRepository
}

const bufferSize = 100

type SiteRequestJob struct {
	domain.IJob
	Index       byte
	SiteTracker *repository.ISiteRequesterRepository
	Host        string
	Created     time.Time
	WaitTime    time.Duration
	Client      *http.Client
}

func (s SiteRequestJob) Do() byte {
	code, _ := (*s.SiteTracker).GetSiteResponseCode(s.Client, s.Host)

	if code >= 200 && code <= 300 {
		return s.Index
	}

	return 0
}

func (s SiteRequestJob) GetCreationTime() time.Time {
	return s.Created
}

func (s SiteRequestJob) GetWaitTime() time.Duration {
	return s.WaitTime
}

func (service *RpsService) BenchHost(host string, maxRps uint, requestTimeoutSeconds int, workersCount uint) (domain.HostBenchmarkStatistic, error) {
	val, ok := service.Hosts.Load(host)

	if !ok {
		return domain.HostBenchmarkStatistic{}, errors.New("host not found but bench has been started")
	}

	stats, ok := val.(domain.HostBenchmarkStatistic)

	if !ok {
		return domain.HostBenchmarkStatistic{}, errors.New("host has a wrong type")
	}

	requestsTotalCount := maxRps * uint(requestTimeoutSeconds)
	sample := make([]uint32, requestTimeoutSeconds)
	waitGroup := sync.WaitGroup{}
	siteResponseChannel := make(chan byte, requestsTotalCount)
	workerService := WorkerService{}
	workerService.Init(workersCount, requestsTotalCount, siteResponseChannel)
	workerService.Start()

	timeout := time.Duration(requestTimeoutSeconds) * time.Second

	// for example server can respond 5 seconds, but response on all requests
	for i := 1; i <= requestTimeoutSeconds; i++ {
		client := client2.CreateHttpClient(timeout)

		for j := uint(0); j < maxRps; j++ {
			waitGroup.Add(1)
			workerService.AddJob(SiteRequestJob{
				Index:       byte(i),
				SiteTracker: &service.SiteTracker,
				Client:      client,
				Host:        host,
				Created:     time.Now(),
				WaitTime:    time.Duration(i) * time.Second,
			})
		}
	}

	go func() {
		waitGroup.Wait()
		close(siteResponseChannel)
	}()

	var completed = false
	for !completed {
		select {
		case job, status := <-siteResponseChannel:
			if status {
				index := uint(job)

				if index > 0 {
					sample[index-1] += 1
				}
				waitGroup.Done()

			} else {
				completed = true
				break
			}
		}
	}

	workerService.Destroy()

	sum := uint32(0)
	for _, rps := range sample {
		sum += rps
	}

	averageRps := sum / uint32(len(sample))
	stats.Rps.AverageRps = uint(averageRps)
	stats.InProgress = false

	return stats, nil
}

func (service *RpsService) ListenForHosts(context *context.Context, waitGroup *sync.WaitGroup, stream <-chan string) {
	maxRps, err := strconv.Atoi(os.Getenv("BENCHMARK_MAX_REQUESTS_PER_SECOND"))
	if err != nil {
		panic(err)
	}

	timeout, err := strconv.Atoi(os.Getenv("BENCHMARK_MAX_TIMEOUT_SECONDS"))
	if err != nil {
		panic(err)
	}

	workersCount, err := strconv.Atoi(os.Getenv("BENCHMARK_WORKERS_PER_HOST_COUNT"))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-(*context).Done():
				return
			default:

			}

			select {
			case host, ok := <-stream:
				if !ok {
					return
				}

				_, ok = service.Hosts.Load(host)

				if !ok {
					service.Hosts.Store(host, domain.HostBenchmarkStatistic{
						InProgress: true,
					})

					go func(host string, waitGroup *sync.WaitGroup) {
						result, err := service.BenchHost(host, uint(maxRps), timeout, uint(workersCount))

						if err != nil {
							fmt.Println(err)
						}

						service.Hosts.Store(host, result)
						waitGroup.Done()

					}(host, waitGroup)

				} else {
					waitGroup.Done()
				}
			case <-(*context).Done():
				return
			}
		}
	}()
}
