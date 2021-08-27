package normalizer

import (
	"github.com/TrueGameover/RestN/rest"
	"github.com/TrueGameover/RevoluterraTest/benchmark/domain"
)

type HostBenchmarkStatisticNormalizer struct {
	rest.IResponseNormalizer
}

func (norm HostBenchmarkStatisticNormalizer) Normalize(object interface{}, normalize rest.NormalizeMethod, options rest.Options, depth int) interface{} {
	statistic, _ := object.(domain.HostBenchmarkStatistic)

	return statistic.Rps.AverageRps
}

func (norm HostBenchmarkStatisticNormalizer) Support(object interface{}) (ok bool) {
	_, ok = object.(domain.HostBenchmarkStatistic)
	return
}
