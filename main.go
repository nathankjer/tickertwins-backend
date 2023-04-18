package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nathankjer/tickertwins-backend/controllers"
)

func main() {
	router := gin.Default()

	router.GET("/tickers", controllers.GetTickers)
	router.GET("/tickers/:symbol/similar", controllers.GetSimilarTickers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(fmt.Sprintf(":%s", port))
}
