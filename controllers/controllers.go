package controllers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MikeB1124/stocks-pattern-lambda/db"
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

	log.Printf("Start Processing %d pattern(s)\n", len(webhookRequest.Data))
	for _, pattern := range webhookRequest.Data {
		log.Printf("Processing pattern %+v\n", pattern)
		var patternData db.DBPatternData
		timeZone, _ := time.LoadLocation("America/Los_Angeles")
		patternData.CreatedAt = time.Now().UTC().In(timeZone).Format("2006-01-02T15:04:05Z")
		patternData.Data = pattern

		if err := db.InsertPattern(patternData); err != nil {
			log.Printf("FAILED inserting pattern %+v\n", pattern)
			return createResponse(Response{Message: err.Error(), StatusCode: 500})
		}
	}
	log.Printf("End Processing %d pattern(s)\n", len(webhookRequest.Data))

	return createResponse(Response{Message: "Harmonic Pattern Webhook Processing Completed", StatusCode: 200})
}
