package repository

import (
	"RevoluterraTest/parser/domain"
)

type ISiteRepository interface {
	GetSites(query string, page uint, limit uint) ([]domain.Site, error)
}
