package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nathankjer/tickertwins-backend/controllers"
)

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://tickertwins.com", "https://tickertwins-frontend-vct4oolaqq-ue.a.run.app/"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	router.Use(cors.New(config))

	router.GET("/tickers", controllers.GetTickers)
	router.GET("/tickers/:symbol/similar", controllers.GetSimilarTickers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(fmt.Sprintf(":%s", port))
}
