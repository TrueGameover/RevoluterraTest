package infrastructure

import (
	"RevoluterraTest/parser/repository"
	"fmt"
	"github.com/kkhrychikov/revo-testing"
	"io/ioutil"
	"net/http"
	"time"
)

type YandexSitesRepository struct {
	repository.ISiteRepository
}

func (rep YandexSitesRepository) GetSites(query string, page uint, limit uint) ([]string, error) {
	var sites []string

	url := fmt.Sprintf(revo.BaseYandexURL, limit, page, query)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36")

	client := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
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
