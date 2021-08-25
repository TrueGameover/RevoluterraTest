package router

import (
	_ "RevoluterraTest/docs"
	"RevoluterraTest/router/sites"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(server *gin.Engine) {
	url := ginSwagger.URL("/swagger/doc.json")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	sites.Init()
	server.GET("/sites", sites.HandleRequest)
}
