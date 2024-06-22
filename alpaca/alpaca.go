package alpaca

import (
	"log"

	"github.com/MikeB1124/stocks-pattern-lambda/configuration"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
)

var client *alpaca.Client

func init() {
	log.Println("Initializing Alpaca client...")
	configuration := configuration.GetConfig()
	client = alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    configuration.Alpaca.ApiKey,
		APISecret: configuration.Alpaca.ApiSecret,
		BaseURL:   configuration.Alpaca.PaperApiUrl,
	})
}

func GetAlpacaAccount() *alpaca.Account {
	acct, err := client.GetAccount()
	if err != nil {
		panic(err)
	}
	return acct
}

func CreateBracketOrder(symbol string, entryPrice float64, qty int, stopLoss float64, takeProfit float64) (*alpaca.Order, error) {
	entryPriceDecimal := decimal.NewFromFloat(entryPrice)
	stopLossDecimal := decimal.NewFromFloat(stopLoss)
	takeProfitDecimal := decimal.NewFromFloat(takeProfit)
	qtyDecimal := decimal.NewFromInt(int64(qty))

	order, err := client.PlaceOrder(alpaca.PlaceOrderRequest{
		OrderClass:  alpaca.Bracket,
		Symbol:      symbol,
		Qty:         &qtyDecimal,
		Side:        alpaca.Buy,
		Type:        alpaca.Limit,
		TimeInForce: alpaca.Day,
		LimitPrice:  &entryPriceDecimal,
		TakeProfit: &alpaca.TakeProfit{
			LimitPrice: &takeProfitDecimal,
		},
		StopLoss: &alpaca.StopLoss{
			StopPrice: &stopLossDecimal,
		},
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func GetAlpacaOrders(status string, symbols []string) ([]alpaca.Order, error) {
	orders, err := client.GetOrders(alpaca.GetOrdersRequest{
		Status:  status,
		Nested:  true,
		Symbols: symbols,
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}
