package main

import (
	"github.com/TrueGameover/RestN/normalizer"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
import "RevoluterraTest/router"

// @title RpsChecked
// @version 0.01
func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	server := gin.Default()
	router.RegisterRoutes(server)
	normalizer.Init()

	if err := server.Run(":5000"); err != nil {
		panic(err)
	}
}
