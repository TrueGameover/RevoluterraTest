package infrastructure

import (
	"fmt"
	client "github.com/TrueGameover/RevoluterraTest/client"
	"github.com/TrueGameover/RevoluterraTest/parser/repository"
	"github.com/kkhrychikov/revo-testing"
	"io/ioutil"
	"net/http"
	url2 "net/url"
)

type YandexSitesRepository struct {
	repository.ISiteRepository
}

func (rep YandexSitesRepository) GetSites(query string, page uint, limit uint) ([]string, error) {
	var sites []string

	url := fmt.Sprintf(revo.BaseYandexURL, limit, page, url2.QueryEscape(query))
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	httpClient := client.CreateHttpClient(15)
	resp, err := httpClient.Do(req)
	defer func(r *http.Response) {
		if r != nil {
			if err := r.Body.Close(); err != nil {
				fmt.Println(err)
			}
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
		sites = append(sites, item.Host)
	}

	return sites, parsed.Error
}
