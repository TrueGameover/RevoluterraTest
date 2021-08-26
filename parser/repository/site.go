package repository

type ISiteRepository interface {
	GetSites(query string, page uint, limit uint) ([]string, error)
}
