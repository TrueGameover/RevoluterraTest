package sites

import (
	infrastructure2 "RevoluterraTest/benchmark/infrastructure"
	repository2 "RevoluterraTest/benchmark/repository"
	service2 "RevoluterraTest/benchmark/service"
	"RevoluterraTest/parser/infrastructure"
	"RevoluterraTest/parser/repository"
	"RevoluterraTest/parser/service"
	context2 "context"
	"github.com/TrueGameover/RestN/rest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"os"
	"strconv"
	"sync"
)

type queryInput struct {
	Query string `form:"query" binding:"required,min=2"`
}

var queueSize uint
var pagesDeep uint
var perPage uint

func Init() {
	size, err := strconv.Atoi(os.Getenv("SITES_QUEUE_SIZE"))
	if err != nil {
		panic(err)
	}
	queueSize = uint(size)

	deep, err := strconv.Atoi(os.Getenv("REQUEST_PAGES_DEEP"))
	if err != nil {
		panic(err)
	}
	pagesDeep = uint(deep)

	limit, err := strconv.Atoi(os.Getenv("REQUEST_PER_PAGE_SITES"))
	if err != nil {
		panic(err)
	}
	perPage = uint(limit)
}

// @Router /sites [get]
// @Summary Calculate sites rps by query
// @Tags SitesController
// @Param query query string true "Your query"
// @Produce json
// @Success 200
func HandleRequest(context *gin.Context) {
	input := queryInput{}
	response := rest.RestResponse{
		Status: 0,
		Locale: rest.Locale{},
		Error: rest.RestError{
			Validation: rest.Validation{},
		},
	}

	if err := context.ShouldBindQuery(&input); err != nil {
		validation := rest.Validation{FieldErrors: []rest.FieldValidationError{}}

		for _, fieldErr := range err.(validator.ValidationErrors) {
			validation.FieldErrors = append(validation.FieldErrors, rest.FieldValidationError{
				Field:   fieldErr.Field(),
				Message: fieldErr.Error(),
			})
		}

		response.Error = rest.RestError{Validation: validation}
		context.JSON(400, response.NormalizeResponse())
	}

	parserService := service.ParserService{SitesRepository: repository.ISiteRepository(infrastructure.YandexSitesRepository{})}
	sitesWaitGroup := sync.WaitGroup{}
	sitesChannel := make(chan string, queueSize)

	parserService.StreamSites(&sitesWaitGroup, sitesChannel, input.Query, pagesDeep, perPage)

	rpsService := service2.RpsService{
		Hosts:       new(sync.Map),
		SiteTracker: repository2.ISiteRequesterRepository(infrastructure2.SiteRequester{}),
	}

	emptyCtx, cancelCtx := context2.WithCancel(context2.Background())
	rpsService.ListenForHosts(&emptyCtx, &sitesWaitGroup, sitesChannel)

	sitesWaitGroup.Wait()
	cancelCtx()
	close(sitesChannel)

	response.SetBody(rpsService.Hosts)
	r := response.NormalizeResponse()

	context.JSON(200, r)
}
