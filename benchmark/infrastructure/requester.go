package infrastructure

import (
	"github.com/TrueGameover/RevoluterraTest/benchmark/repository"
	"github.com/hashicorp/go-multierror"
	"io/ioutil"
	"net/http"
)

type SiteRequester struct {
	repository.ISiteRequesterRepository
}

const Scheme = "https://"

func (s SiteRequester) GetSiteResponseCode(client *http.Client, host string) (int, error) {
	var errMain *multierror.Error

	req, err := http.NewRequest(http.MethodGet, Scheme+host, nil)

	if err != nil {
		errMain = multierror.Append(errMain, err)
	}

	resp, err := client.Do(req)

	if err != nil {
		errMain = multierror.Append(errMain, err)
	}

	if resp != nil {
		defer resp.Body.Close()
		_, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			errMain = multierror.Append(errMain, err)
		}

		return resp.StatusCode, errMain.ErrorOrNil()
	}

	return 0, errMain.ErrorOrNil()
}
