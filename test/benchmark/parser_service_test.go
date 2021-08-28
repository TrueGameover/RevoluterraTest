package benchmark

import (
	"github.com/TrueGameover/RevoluterraTest/parser/infrastructure"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestYandex(t *testing.T) {
	sitesRepository := infrastructure.YandexSitesRepository{}
	sites, err := sitesRepository.GetSites("playstation купить", 1, 10)

	require.Nil(t, err)
	require.NotEmpty(t, sites)
}
