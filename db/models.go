package db

type PatternWebhookRequest struct {
	Data []PatternData `json:"data"`
}

type PatternData struct {
	PatternType   string  `json:"patterntype"`
	PatternName   string  `json:"patternname"`
	ProfitOne     float64 `json:"profit1"`
	DisplaySymbol string  `json:"displaysymbol"`
	Symbol        string  `json:"symbol"`
	StopLoss      float64 `json:"stoploss"`
	PatternUrl    string  `json:"url"`
	TimeFrame     string  `json:"timeframe"`
	Status        string  `json:"status"`
	Entry         string  `json:"entry"`
	PatternClass  string  `json:"patternclass"`
}

type DBPatternData struct {
	Data      PatternData `json:"data"`
	CreatedAt string      `json:"createdAt"`
}
