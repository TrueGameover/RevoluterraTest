package client

import (
	"net/http"
	"time"
)

type defaultHeaders struct {
	roundTripper http.RoundTripper
}

func (tripper defaultHeaders) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36")
	r.Header.Add("Referrer-Policy", "origin")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Accept-Encoding", "identity")
	r.Header.Add("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	r.Header.Add("Cache-Control", "no-cache")

	return tripper.roundTripper.RoundTrip(r)
}

func CreateHttpClient(timeoutSeconds time.Duration) *http.Client {
	client := http.Client{
		Timeout:   timeoutSeconds * time.Second,
		Transport: defaultHeaders{roundTripper: http.DefaultTransport},
	}

	return &client
}
