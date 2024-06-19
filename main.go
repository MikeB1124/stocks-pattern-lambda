package main

import (
	"github.com/MikeB1124/stocks-pattern-lambda/controllers"
	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/lambda"
)

var router *lmdrouter.Router

func init() {
	router = lmdrouter.NewRouter("")
	router.Route("POST", "/webhook/harmonic-pattern", controllers.HarmonicPatternWebhook)
}

func main() {
	lambda.Start(router.Handler)
}
