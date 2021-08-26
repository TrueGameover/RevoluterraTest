package repository

import "context"

type ISiteRequesterRepository interface {
	GetSiteResponseCode(context *context.Context, host string) (int, error)
}
