package db

type PatternWebhookRequest struct {
	MsgType string        `json:"msg_type" bson:"msg_type"`
	Data    []PatternData `json:"data" bson:"data"`
}

type PatternData struct {
	PatternType   string  `json:"patterntype" bson:"patterntype"`
	PatternName   string  `json:"patternname" bson:"patternname"`
	ProfitOne     float64 `json:"profit1" bson:"profit1"`
	DisplaySymbol string  `json:"displaysymbol" bson:"displaysymbol"`
	Symbol        string  `json:"symbol" bson:"symbol"`
	StopLoss      float64 `json:"stoploss" bson:"stoploss"`
	PatternUrl    string  `json:"url" bson:"url"`
	TimeFrame     string  `json:"timeframe" bson:"timeframe"`
	Status        string  `json:"status" bson:"status"`
	Entry         string  `json:"entry" bson:"entry"`
	PatternClass  string  `json:"patternclass" bsom:"patternclass"`
}

type AlpacaEntryOrder struct {
	EntryOrderID     string      `json:"entryOrderID" bson:"entryOrderID"`
	EntryOrderStatus string      `json:"entryOrderStatus" bson:"entryOrderStatus"`
	Qty              int         `json:"qty" bson:"qty"`
	ExitOrderID      string      `json:"exitOrderID" bson:"exitOrderID"`
	ExitOrderStatus  string      `json:"exitOrderStatus" bson:"exitOrderStatus"`
	TradeProfit      float64     `json:"tradeProfit" bson:"tradeProfit"`
	Data             PatternData `json:"data" bson:"data"`
	CreatedAt        string      `json:"createdAt" bson:"createdAt"`
	EntryUpdatedAt   string      `json:"entryUpdatedAt" bson:"entryUpdatedAt"`
	ExitUpdatedAt    string      `json:"exitUpdatedAt" bson:"exitUpdatedAt"`
}
