package service

import (
	"RevoluterraTest/benchmark/domain"
	"RevoluterraTest/benchmark/repository"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type RpsService struct {
	Hosts       *sync.Map
	SiteTracker repository.ISiteRequesterRepository
}

func (service *RpsService) BenchHost(parentContext *context.Context, host string, maxRps int, requestTimeoutSeconds int) (domain.HostBenchmarkStatistic, error) {
	val, ok := service.Hosts.Load(host)

	if !ok {
		return domain.HostBenchmarkStatistic{}, errors.New("host not found but bench has been started")
	}

	stats, ok := val.(domain.HostBenchmarkStatistic)

	if !ok {
		return domain.HostBenchmarkStatistic{}, errors.New("host has a wrong type")
	}

	sample := make([]uint32, requestTimeoutSeconds)
	waitGroup := sync.WaitGroup{}

	// for example server can respond 5 seconds, but response on all requests
	for i := 0; i < requestTimeoutSeconds; i++ {
		start := time.Now()

		sample[i] = 0
		ctx, cancel := context.WithTimeout(*parentContext, time.Duration(requestTimeoutSeconds)*time.Second)
		defer cancel()

		// every second
		for j := 0; j < maxRps; j++ {
			waitGroup.Add(1)
			go func(second int) {
				code, err := service.SiteTracker.GetSiteResponseCode(&ctx, host)

				/*if err != nil {
					fmt.Println(err)
				}*/

				if code >= 200 && code <= 302 && err == nil {
					atomic.AddUint32(&sample[second], 1)
				}

				waitGroup.Done()
			}(i)
		}

		elapsed := time.Since(start)

		// try to be more accurate
		if elapsed.Milliseconds() < 1000 {
			time.Sleep(time.Duration(1000 - elapsed.Milliseconds()))
		}
	}

	waitGroup.Wait()

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

	go func() {
		for {
			select {
			case <-(*context).Done():
				break
			default:

			}

			select {
			case host, ok := <-stream:
				if !ok {
					break
				}

				_, ok = service.Hosts.Load(host)

				if !ok {
					service.Hosts.Store(host, domain.HostBenchmarkStatistic{
						InProgress: true,
					})

					go func(host string, waitGroup *sync.WaitGroup) {
						result, err := service.BenchHost(context, host, maxRps, timeout)

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
				break
			}
		}
	}()
}
