package controllers

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/MikeB1124/stocks-pattern-lambda/alpaca"
	"github.com/MikeB1124/stocks-pattern-lambda/db"
	"github.com/MikeB1124/stocks-pattern-lambda/stockutils"
	"github.com/MikeB1124/stocks-pattern-lambda/utils"
	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func HarmonicPatternWebhook(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing harmonic pattern %+v\n", event)

	var webhookRequest db.PatternWebhookRequest
	err := json.Unmarshal([]byte(event.Body), &webhookRequest)
	if err != nil {
		return createResponse(Response{Message: err.Error(), StatusCode: 400})
	}

	if webhookRequest.MsgType != "pattern.notification" {
		log.Printf("Invalid msg_type %+v\n", webhookRequest.MsgType)
		return createResponse(Response{Message: "Invalid msg_type", StatusCode: 400})
	}

	failedCount := 0
	log.Printf("Start Creating %d Order(s)\n", len(webhookRequest.Data))
	for _, pattern := range webhookRequest.Data {
		log.Printf("Processing pattern %+v\n", pattern)

		// Only support US equities
		if !strings.Contains(pattern.Symbol, ".US") {
			log.Printf("Non US equities are currently not supported  %+v\n", pattern.Symbol)
			failedCount++
			continue
		}

		// Only support bullish trend
		if pattern.PatternType != "bullish" {
			log.Println("We only support bullish trent at the moment.")
			failedCount++
			continue
		}

		// Parse entry price to float
		entryPrice, err := strconv.ParseFloat(strings.Split(pattern.Entry, "_")[0], 64)
		if err != nil {
			log.Printf("Failed to parse entry price %+v\n", pattern.Entry)
			failedCount++
			continue
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
		order, err := alpaca.CreateOrder(pattern.DisplaySymbol, entryPrice, qtyToBuy, pattern.StopLoss)
		if err != nil {
			log.Printf("Failed to create order %+v\n", err)
			failedCount++
			continue
		}
		log.Printf("Order created %+v\n", order)

		// Insert entry order to database
		var alpacaEntryOrder db.AlpacaEntryOrder
		alpacaEntryOrder.EntryOrderID = order.ID
		alpacaEntryOrder.EntryOrderStatus = "CREATED"
		alpacaEntryOrder.Qty = qtyToBuy
		alpacaEntryOrder.Data = pattern
		alpacaEntryOrder.CreatedAt = utils.GetCurrentTime()

		if err := db.InsertEntryOrder(alpacaEntryOrder); err != nil {
			log.Printf("FAILED inserting order %+v\n", alpacaEntryOrder)
			// Cancel the order that was created
			failedCount++
			return createResponse(Response{Message: err.Error(), StatusCode: 500})
		}
	}
	log.Printf("Orders Created, Successful Orders: %d   Failed Orders: %d\n", len(webhookRequest.Data)-failedCount, failedCount)

	return createResponse(Response{Message: "Harmonic Pattern Webhook Processing Completed", StatusCode: 200})
}
