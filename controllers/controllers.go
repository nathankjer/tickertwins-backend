package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nathankjer/tickertwins-backend/db"
	"github.com/nathankjer/tickertwins-backend/models"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func validateParam(c *gin.Context, param, paramName, errorMsg string) bool {
	if param == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": errorMsg})
		return false
	}
	return true
}

func GetTickers(c *gin.Context) {
	query := strings.TrimSpace(strings.ToUpper(c.Query("q")))
	if !validateParam(c, query, "q", "Missing required parameter 'q'.") {
		return
	}

	var tickers []models.Ticker
	err := db.DB.Where("UPPER(symbol) ILIKE ? AND enabled = ?", query+"%", true).
		Or("UPPER(name) ILIKE ? AND enabled = ?", "%"+query+"%", true).
		Limit(7).
		Find(&tickers).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Error retrieving tickers."})
		return
	}
	c.JSON(200, tickers)
}

func GetSimilarTickers(c *gin.Context) {
	symbol := strings.TrimSpace(strings.ToUpper(c.Param("symbol")))
	if !validateParam(c, symbol, "symbol", "Missing required path parameter 'symbol'.") {
		return
	}

	var ticker models.Ticker
	err := db.DB.Where("UPPER(symbol) = ? AND enabled = ?", symbol, true).First(&ticker).Error
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "Ticker not found."})
		return
	}

	var relatedTickers []models.Ticker
	err = db.DB.Table("similar_tickers").
		Select("tickers.*").
		Joins("JOIN tickers ON similar_tickers.related_ticker_id = tickers.id").
		Where("similar_tickers.ticker_id = ? AND tickers.enabled = ?", ticker.ID, true).
		Order("similar_tickers.position").
		Limit(30).
		Find(&relatedTickers).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Error retrieving similar tickers."})
		return
	}

	response := models.SimilarTickerResponse{
		Ticker:         ticker,
		SimilarTickers: relatedTickers,
	}
	c.JSON(200, response)
}
