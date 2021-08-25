package service

import (
	"RevoluterraTest/parser/domain"
	"RevoluterraTest/parser/repository"
	"fmt"
	"sync"
)

type ParserService struct {
	SitesRepository repository.ISiteRepository
}

func (service *ParserService) StreamSites(waitGroup *sync.WaitGroup, sitesChannel chan domain.Site, query string, tryPagesCount uint, perPage uint) {
	for i := uint(1); i <= tryPagesCount; i++ {
		waitGroup.Add(1)
		go func(page uint) {
			sites, err := service.SitesRepository.GetSites(query, page, perPage)

			if err != nil {
				fmt.Println(err)
			}

			for _, site := range sites {
				sitesChannel <- site
			}

			waitGroup.Done()
		}(i)
	}
}