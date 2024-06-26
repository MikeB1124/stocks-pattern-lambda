package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	stockslambdautils "github.com/MikeB1124/stocks-lambda-utils/v2"
	"github.com/MikeB1124/stocks-pattern-lambda/configuration"
	"github.com/MikeB1124/stocks-pattern-lambda/stockutils"
	"github.com/aws/aws-lambda-go/events"
)

var pacificTime, _ = time.LoadLocation("America/Los_Angeles")

func HarmonicPatternWebhook(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing harmonic pattern %+v\n", event)

	// Do not allow trading on weekends
	dateTimeNow := time.Now().In(pacificTime)
	if dateTimeNow.Weekday() == time.Saturday || dateTimeNow.Weekday() == time.Sunday {
		log.Println("Weekend, no trading allowed")
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "Weekend, no trading allowed", StatusCode: 400})
	}

	// Do not allow trading after 1:00 PM PST
	if dateTimeNow.Hour() >= 13 {
		log.Println("This lambda does not allow trading after 1:00 PM PST")
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "This lambda does not allow trading after 1:00 PM PST", StatusCode: 400})
	}

	var webhookRequest stockslambdautils.PatternWebhookRequest
	err := json.Unmarshal([]byte(event.Body), &webhookRequest)
	if err != nil {
		log.Printf("Failed to unmarshal request %+v\n", err)
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: err.Error(), StatusCode: 400})
	}

	if webhookRequest.MsgType != "pattern.notification" {
		log.Printf("Invalid msg_type %+v\n", webhookRequest.MsgType)
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "Invalid msg_type", StatusCode: 400})
	}

	failedCount := 0
	log.Printf("Start Creating %d Order(s)\n", len(webhookRequest.Data))
	for _, pattern := range webhookRequest.Data {
		log.Printf("Processing pattern %+v\n", pattern)

		// Only support US equities
		if !strings.Contains(pattern.Symbol, ".US") {
			log.Printf("Non US equities are currently not supported  %+v\n", pattern.Symbol)
			continue
		}

		// Only support bullish trend
		if pattern.PatternType != "bullish" {
			log.Println("We only support bullish trent at the moment.")
			continue
		}

		// Check if open orders exist for the symbol
		openOrders, err := configuration.AlpacaClient.GetAlpacaOrders("open", []string{pattern.DisplaySymbol}, true)
		if err != nil {
			log.Printf("Failed to get open orders %+v\n", err)
			failedCount++
			continue
		}

		if len(openOrders) > 0 {
			log.Printf("Open orders already exists for %+v\n", pattern.DisplaySymbol)
			continue
		}

		var entryPrice float64
		if _, ok := pattern.Entry.(float64); ok {
			entryPrice = pattern.Entry.(float64)
		} else {
			// Parse entry price to float
			entryPrice, err = strconv.ParseFloat(strings.Split(pattern.Entry.(string), "_")[0], 64)
			if err != nil {
				log.Printf("Failed to parse entry price %+v\n", pattern.Entry)
				failedCount++
				continue
			}
		}

		// Round to 2 decimal places
		entryPrice = float64(int(entryPrice*100)) / 100
		stopPrice := float64(int(pattern.StopLoss*100)) / 100
		takeProfitPrice := float64(int(pattern.ProfitOne*100)) / 100

		// Calculate the number of shares to buy
		qtyToBuy, err := stockutils.SharesToBuy(entryPrice)
		if err != nil {
			log.Println(err)
			failedCount++
			continue
		}

		log.Printf("Buy %d shares at %f\n", qtyToBuy, entryPrice)

		// Create order
		order, err := configuration.AlpacaClient.CreateAlpacaOrder(
			pattern.DisplaySymbol,
			entryPrice,
			stopPrice,
			takeProfitPrice,
			qtyToBuy,
			"buy",
			"limit",
			"bracket",
			"gtc",
		)
		if err != nil {
			log.Printf("Failed to create order %+v\n", err)
			failedCount++
			continue
		}
		log.Printf("Order created %+v\n", order)

		// Insert entry order to database
		var alpacaTradeOrder stockslambdautils.AlpacaTrade
		alpacaTradeOrder.Order = stockslambdautils.FormatAlpacaOrderForDB(order)
		alpacaTradeOrder.PatternData = pattern

		timeNow := time.Now().UTC()
		alpacaTradeOrder.RecordUpdatedAt = &timeNow

		if err := configuration.MongoClient.InsertEntryOrder(alpacaTradeOrder); err != nil {
			log.Printf("FAILED inserting order %+v\n", alpacaTradeOrder)

			// Cancel the order that was created
			if err := configuration.AlpacaClient.CancelAlpacaOrder(alpacaTradeOrder.Order.ID); err != nil {
				log.Printf("Failed to cancel order %+v\n", err)
			}
			failedCount++
			continue
		}
	}
	log.Printf("Orders Created, Successful Orders: %d   Failed Orders: %d\n", len(webhookRequest.Data)-failedCount, failedCount)
	return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: fmt.Sprintf("Orders Created, Successful Orders: %d Failed Orders: %d", len(webhookRequest.Data)-failedCount, failedCount), StatusCode: 200})
}
