package benchmark

import (
	"fmt"
	"github.com/TrueGameover/RevoluterraTest/benchmark/domain"
	infrastructure2 "github.com/TrueGameover/RevoluterraTest/benchmark/infrastructure"
	repository2 "github.com/TrueGameover/RevoluterraTest/benchmark/repository"
	service2 "github.com/TrueGameover/RevoluterraTest/benchmark/service"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestBenchHost(t *testing.T) {
	rpsService := service2.RpsService{
		Hosts:       new(sync.Map),
		SiteTracker: repository2.ISiteRequesterRepository(infrastructure2.SiteRequester{}),
	}

	const host = "google.com"
	rpsService.Hosts.Store(host, domain.HostBenchmarkStatistic{})

	statistic, _ := rpsService.BenchHost(host, 10000, 5, 100)

	fmt.Println(statistic.Rps.AverageRps)

	require.Equal(t, false, statistic.InProgress)
}
