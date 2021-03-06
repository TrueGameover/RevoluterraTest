package service

import (
	"fmt"
	"github.com/TrueGameover/RevoluterraTest/parser/repository"
	"sync"
)

type ParserService struct {
	SitesRepository repository.ISiteRepository
}

func (service *ParserService) StreamSites(waitGroup *sync.WaitGroup, sitesChannel chan<- string, query string, tryPagesCount uint, perPage uint) {
	for i := uint(1); i <= tryPagesCount; i++ {
		waitGroup.Add(1)
		go func(page uint) {
			sites, err := service.SitesRepository.GetSites(query, page, perPage)

			if err != nil {
				fmt.Println(err)
			}

			for _, site := range sites {
				if len(site) > 0 {
					waitGroup.Add(1)
					sitesChannel <- site
				}
			}

			fmt.Printf("Sites count = %d\n", len(sites))
			waitGroup.Done()
		}(i)
	}
}
