package client

import (
	"net/http"
	"time"
)

type defaultHeaders struct {
	roundTripper http.RoundTripper
	cookies      []http.Cookie
}

func (tripper defaultHeaders) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36")
	r.Header.Add("Referrer-Policy", "origin")
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	r.Header.Add("Accept-Encoding", "identity")
	r.Header.Add("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Referer", "https://www.google.com/")

	for _, cookie := range tripper.cookies {
		r.AddCookie(&cookie)
	}

	resp, err := tripper.roundTripper.RoundTrip(r)

	if len(tripper.cookies) == 0 && resp != nil {
		for _, cookie := range resp.Cookies() {
			tripper.cookies = append(tripper.cookies, *cookie)
		}
	}

	return resp, err
}

func CreateHttpClient(timeoutSeconds time.Duration) *http.Client {
	client := http.Client{
		Timeout:   timeoutSeconds * time.Second,
		Transport: defaultHeaders{roundTripper: http.DefaultTransport},
	}

	return &client
}
