package repository

import "net/http"

type ISiteRequesterRepository interface {
	GetSiteResponseCode(client *http.Client, host string) (int, error)
}
