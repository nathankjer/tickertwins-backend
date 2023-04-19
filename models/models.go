package models

type Ticker struct {
	ID          uint   `gorm:"primaryKey" json:"-"`
	Symbol      string `gorm:"column:symbol" json:"symbol"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	Types       string `gorm:"column:types" json:"type"`
	Enabled     bool   `gorm:"column:enabled" json:"-"`
}

type SimilarTicker struct {
	ID              uint `gorm:"primaryKey" json:"-"`
	TickerID        uint `gorm:"column:ticker_id" json:"ticker_id"`
	RelatedTickerID uint `gorm:"column:related_ticker_id" json:"related_ticker_id"`
	Position        uint `gorm:"column:position" json:"position"`
}

type SimilarTickerResponse struct {
	Ticker         Ticker   `json:"ticker"`
	SimilarTickers []Ticker `json:"similar_tickers"`
}
