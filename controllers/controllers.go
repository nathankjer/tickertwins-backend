package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nathankjer/tickertwins-backend/db"
	"github.com/nathankjer/tickertwins-backend/models"
)

func GetTickers(c *gin.Context) {
	query := strings.ToUpper(c.Query("q"))
	if query == "" {
		c.JSON(400, gin.H{"error": "Missing required parameter 'q'."})
		return
	}

	var tickers []models.Ticker
	db.DB.Where("UPPER(symbol) ILIKE ? AND enabled = ?", query+"%", true).
		Or("UPPER(name) ILIKE ? AND enabled = ?", "%"+query+"%", true).
		Limit(7).
		Find(&tickers)
        c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, tickers)
}

func GetSimilarTickers(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))
	if symbol == "" {
		c.JSON(400, gin.H{"error": "Missing required path parameter 'symbol'."})
		return
	}

	var ticker models.Ticker
	if err := db.DB.Where("UPPER(symbol) = ? AND enabled = ?", symbol, true).First(&ticker).Error; err != nil {
		c.JSON(404, gin.H{"error": "Ticker not found."})
		return
	}

	var relatedTickers []models.Ticker
	db.DB.Table("similar_tickers").
		Select("tickers.*").
		Joins("JOIN tickers ON similar_tickers.related_ticker_id = tickers.id").
		Where("similar_tickers.ticker_id = ? AND tickers.enabled = ?", ticker.ID, true).
		Order("similar_tickers.position").
		Limit(30).
		Find(&relatedTickers)

	response := models.SimilarTickerResponse{
		Ticker:         ticker,
		SimilarTickers: relatedTickers,
	}
        c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, response)
}
