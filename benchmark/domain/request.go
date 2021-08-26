package domain

type HostBenchmarkStatistic struct {
	Time       TimeStatistic
	Rps        RpsStatistic
	InProgress bool
}
