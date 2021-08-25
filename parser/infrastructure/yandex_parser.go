package infrastructure

import (
	"RevoluterraTest/parser/domain"
	"RevoluterraTest/parser/repository"
	"fmt"
	"github.com/kkhrychikov/revo-testing"
	"io/ioutil"
	"net/http"
)

type YandexSitesRepository struct {
	repository.ISiteRepository
}

func (rep YandexSitesRepository) GetSites(query string, page uint, limit uint) ([]domain.Site, error) {
	var sites []domain.Site
	url := fmt.Sprintf(revo.BaseYandexURL, limit, page, query)
	resp, err := http.Get(url)
	defer func(r *http.Response) {
		if err := r.Body.Close(); err != nil {
			fmt.Println(err)
		}

	}(resp)

	if err != nil {
		return sites, err
	}

	raw, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return sites, err
	}

	parsed := revo.ParseYandexResponse(raw)
	for _, item := range parsed.Items {
		sites = append(sites, domain.Site{Host: item.Host})
	}

	return sites, parsed.Error
}
