package stockutils

import (
	"fmt"

	"github.com/MikeB1124/stocks-pattern-lambda/alpaca"
	"github.com/shopspring/decimal"
)

var percentOfPortfolio = decimal.NewFromFloat(0.01)

func SharesToBuy(entryPrice float64) (int, error) {
	account := alpaca.GetAlpacaAccount()

	if account.NonMarginBuyingPower.IsZero() {
		return 0, fmt.Errorf("Alpaca account NonMarginBuyingPower is 0")
	}

	amountToSpent := account.NonMarginBuyingPower.Mul(percentOfPortfolio)
	quantity := amountToSpent.DivRound(decimal.NewFromFloat(entryPrice), 0)

	if quantity.IsZero() {
		return 0, fmt.Errorf("Can not buy any shares with the current price of %f", entryPrice)
	}

	return int(quantity.IntPart()), nil
}
