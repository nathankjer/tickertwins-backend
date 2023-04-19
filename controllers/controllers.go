package controllers

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nathankjer/tickertwins-backend/db"
	"github.com/nathankjer/tickertwins-backend/models"
)

func GetTickers(c *gin.Context) {
	query := strings.TrimSpace(strings.ToUpper(c.Query("q")))
	if query == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "Missing required parameter 'q'."})
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
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, tickers)
}

func GetSimilarTickers(c *gin.Context) {
	symbol := strings.TrimSpace(strings.ToUpper(c.Param("symbol")))
	if symbol == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "Missing required parameter 'symbol'."})
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
		Limit(50).
		Find(&relatedTickers).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Error retrieving similar tickers."})
		return
	}

	response := models.SimilarTickerResponse{
		Ticker:         ticker,
		SimilarTickers: relatedTickers,
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, response)
}

var sp500 = []string{"AAPL", "MSFT", "AMZN", "NVDA", "GOOGL", "BRK.B", "GOOG", "TSLA", "META", "XOM", "UNH", "JNJ", "JPM", "V", "PG", "MA", "CVX", "HD", "LLY", "MRK", "ABBV", "AVGO", "PEP", "KO", "PFE", "TMO", "COST", "BAC", "MCD", "WMT", "CSCO", "CRM", "DIS", "ABT", "LIN", "ACN", "ADBE", "DHR", "VZ", "TXN", "CMCSA", "WFC", "NKE", "NEE", "PM", "RTX", "BMY", "NFLX", "ORCL", "AMD", "UPS", "T", "QCOM", "INTC", "AMGN", "HON", "COP", "SBUX", "LOW", "INTU", "UNP", "CAT", "MS", "ELV", "IBM", "BA", "GS", "SPGI", "PLD", "LMT", "MDT", "DE", "GE", "BLK", "GILD", "BKNG", "SYK", "AXP", "CVS", "AMT", "ADI", "C", "MDLZ", "NOW", "AMAT", "ISRG", "ADP", "TJX", "TMUS", "REGN", "PYPL", "MMC", "VRTX", "CB", "MO", "ZTS", "PGR", "SCHW", "SO", "CI", "DUK", "TGT", "FISV", "BSX", "SLB", "BDX", "EOG", "CME", "NOC", "MU", "AON", "LRCX", "EQIX", "ITW", "ETN", "HUM", "CSX", "APD", "CL", "WM", "ATVI", "ICE", "FCX", "MMM", "MPC", "EL", "CDNS", "SNPS", "HCA", "CCI", "ORLY", "SHW", "PXD", "FDX", "EW", "GD", "GIS", "KLAC", "PNC", "AZO", "F", "MCK", "USB", "EMR", "VLO", "CMG", "GM", "D", "MSI", "SRE", "PSX", "AEP", "NSC", "DG", "MCO", "MRNA", "ROP", "KMB", "APH", "PSA", "DXCM", "OXY", "MAR", "TFC", "NXPI", "ADM", "CTVA", "MCHP", "FTNT", "AJG", "MSCI", "ADSK", "EXC", "BIIB", "PH", "A", "ECL", "TT", "MET", "ANET", "HES", "TEL", "MNST", "DOW", "CTAS", "JCI", "IDXX", "TRV", "TDG", "HLT", "YUM", "O", "LHX", "AIG", "NEM", "XEL", "SYY", "PCAR", "HSY", "CNC", "AFL", "IQV", "CARR", "COF", "NUE", "STZ", "WMB", "CHTR", "SPG", "ROST", "ILMN", "DVN", "WELL", "MTD", "PAYX", "KMI", "ED", "OTIS", "FIS", "ON", "EA", "CMI", "CPRT", "AMP", "VICI", "RMD", "PPG", "DD", "BK", "WBD", "PRU", "AME", "ROK", "PEG", "KHC", "CTSH", "KR", "DHI", "ODFL", "DLTR", "ENPH", "FAST", "ALL", "WEC", "HAL", "VRSK", "GEHC", "KDP", "GWW", "OKE", "BKR", "APTV", "AWK", "GPN", "SBAC", "RSG", "CSGP", "ZBH", "ANSS", "ES", "DLR", "EIX", "DFS", "KEYS", "ULTA", "PCG", "ABC", "LEN", "WST", "HPQ", "TSCO", "FANG", "URI", "GLW", "ACGL", "WBA", "WTW", "CDW", "TROW", "STT", "ALGN", "IT", "LYB", "IFF", "CEG", "AVB", "ALB", "EFX", "PWR", "FTV", "EBAY", "GPC", "WY", "AEE", "VMC", "CBRE", "IR", "PODD", "DAL", "ETR", "FE", "HIG", "MLM", "CHD", "DTE", "BAX", "FSLR", "MPWR", "MKC", "MTB", "PPL", "CAH", "EXR", "EQR", "HOLX", "DOV", "TDY", "LH", "HPE", "CTRA", "VRSN", "TTWO", "CLX", "OMC", "ARE", "CNP", "LUV", "INVH", "LVS", "XYL", "NDAQ", "FITB", "STE", "DRI", "RJF", "WAT", "COO", "WAB", "CMS", "NTRS", "TSN", "VTR", "RF", "EXPD", "CAG", "SWKS", "SEDG", "STLD", "FICO", "PFG", "MAA", "K", "TRGP", "PKI", "BR", "NVR", "MOH", "CINF", "EPAM", "HBAN", "AMCR", "IEX", "SJM", "FLT", "ATO", "DGX", "AES", "MOS", "BALL", "FDS", "HWM", "MRO", "LW", "FMC", "ZBRA", "IRM", "TER", "CF", "GRMN", "TYL", "PAYC", "J", "CFG", "IPG", "BBY", "NTAP", "JBHT", "AVY", "TXT", "CBOE", "BG", "RE", "EVRG", "LKQ", "BRO", "MGM", "PHM", "INCY", "EXPE", "UAL", "RCL", "PTC", "LNT", "ESS", "TECH", "PKG", "POOL", "AKAM", "SYF", "IP", "ETSY", "APA", "SNA", "MKTX", "WRB", "LDOS", "UDR", "STX", "TFX", "TRMB", "VTRS", "EQT", "HST", "DPZ", "PEAK", "CPT", "NDSN", "SWK", "KIM", "WYNN", "WDC", "KEY", "BWA", "BF.B", "HRL", "JKHY", "NI", "CHRW", "HSIC", "KMX", "CPB", "L", "PARA", "MAS", "CE", "JNPR", "TAP", "CRL", "CDAY", "FOXA", "GEN", "BIO", "MTCH", "EMN", "TPR", "GL", "CCL", "LYV", "QRVO", "CZR", "REG", "ALLE", "ROL", "PNW", "XRAY", "UHS", "PNR", "AOS", "FFIV", "AAL", "HII", "RHI", "NRG", "CTLT", "BBWI", "IVZ", "WRK", "BEN", "AAP", "BXP", "WHR", "VFC", "FRT", "SEE", "HAS", "NWSA", "GNRC", "AIZ", "OGN", "CMA", "DXC", "NCLH", "ALK", "MHK", "RL", "NWL", "DVA", "ZION", "FOX", "LNC", "FRC", "NWS", "DISH"}

func GetRandomTickers(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	selectedSymbols := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		index := rand.Intn(len(sp500))
		selectedSymbols = append(selectedSymbols, sp500[index])
		sp500[index] = sp500[len(sp500)-1]
		sp500 = sp500[:len(sp500)-1]
	}

	var tickers []models.Ticker
	err := db.DB.Where("symbol IN (?) AND enabled = ?", selectedSymbols, true).
		Find(&tickers).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Error retrieving random tickers."})
		return
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(200, tickers)
}
