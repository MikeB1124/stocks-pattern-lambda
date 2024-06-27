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

func HarmonicPatternWebhook(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing harmonic pattern %+v\n", event)

	var webhookRequest stockslambdautils.PatternWebhookRequest
	err := json.Unmarshal([]byte(event.Body), &webhookRequest)
	if err != nil {
		log.Printf("Failed to unmarshal request %+v\n", err)
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: err.Error(), StatusCode: 400})
	}
	log.Printf("Unmarshaled Webhook data %+v", webhookRequest)

	if webhookRequest.MsgType != "pattern.notification" {
		log.Printf("Invalid msg_type %+v\n", webhookRequest.MsgType)
		return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: "Invalid msg_type", StatusCode: 400})
	}
	log.Printf("Harmonic pattern message type: %s", webhookRequest.MsgType)

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
		openOrders, err := configuration.AlpacaClient.GetAlpacaOrders("open", []string{pattern.DisplaySymbol})
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
			qtyToBuy,
			"buy",
			"limit",
		)
		if err != nil {
			log.Printf("Failed to create order %+v\n", err)
			failedCount++
			continue
		}
		log.Printf("Order created %+v\n", order)

		// Insert entry order to database
		var alpacaEntryOrder stockslambdautils.AlpacaTrade
		alpacaEntryOrder.EntryOrder = stockslambdautils.FormatAlpacaOrderForDB(order)
		alpacaEntryOrder.ExitOrder = nil
		alpacaEntryOrder.PatternData = pattern

		timeNow := time.Now().UTC()
		alpacaEntryOrder.RecordUpdatedAt = &timeNow

		if err := configuration.MongoClient.InsertEntryOrder(alpacaEntryOrder); err != nil {
			log.Printf("FAILED inserting order %+v\n", alpacaEntryOrder)

			// Cancel the order that was created
			if err := configuration.AlpacaClient.CancelAlpacaOrder(alpacaEntryOrder.EntryOrder.ID); err != nil {
				log.Printf("Failed to cancel order %+v\n", err)
			}
			failedCount++
			continue
		}
	}
	log.Printf("Orders Created, Successful Orders: %d   Failed Orders: %d\n", len(webhookRequest.Data)-failedCount, failedCount)
	return stockslambdautils.CreateResponse(stockslambdautils.Response{Message: fmt.Sprintf("Orders Created, Successful Orders: %d Failed Orders: %d", len(webhookRequest.Data)-failedCount, failedCount), StatusCode: 200})
}
