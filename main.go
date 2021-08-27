package main

import (
	"github.com/TrueGameover/RestN/normalizer"
	"github.com/TrueGameover/RestN/rest"
	normalizer2 "github.com/TrueGameover/RevoluterraTest/benchmark/infrastructure/normalizer"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
import "github.com/TrueGameover/RevoluterraTest/router"

// @title RpsChecked
// @version 0.01
func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	server := gin.Default()
	router.RegisterRoutes(server)
	initNormalizers()

	if err := server.Run(":5000"); err != nil {
		panic(err)
	}
}

func initNormalizers() {
	normalizer.Init()
	rest.RegisterNormalizer(normalizer2.HostBenchmarkStatisticNormalizer{})
}
