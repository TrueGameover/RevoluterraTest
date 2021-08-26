package infrastructure

import (
	"RevoluterraTest/benchmark/repository"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
)

type SiteRequester struct {
	repository.ISiteRequesterRepository
}

func (s SiteRequester) GetSiteResponseCode(context *context.Context, host string) (int, error) {
	req, err := http.NewRequestWithContext(*context, http.MethodGet, host, nil)

	if err != nil {
		println(err)
	}

	resp, _ := http.DefaultClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)

		return resp.StatusCode, err
	}

	return -1, errors.New("no response")
}
