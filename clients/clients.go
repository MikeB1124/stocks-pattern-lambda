package clients

import (
	"context"
	"fmt"
	"log"

	stockslambdautils "github.com/MikeB1124/stocks-lambda-utils"
	"github.com/MikeB1124/stocks-pattern-lambda/configuration"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient stockslambdautils.MongoClient
var AlpacaClient stockslambdautils.AlpacaClient

func init() {
	log.Println("Loading clients...")
	config := configuration.GetConfig()
	initMongoDB(config)
	initAlpaca(config)
}

func initMongoDB(config configuration.Configration) {
	log.Println("Initializing MongoDB client...")
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.du0vf.mongodb.net", config.MongoDB.Username, config.MongoDB.Password))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	MongoClient = stockslambdautils.MongoClient{Client: client}
}

func initAlpaca(config configuration.Configration) {
	log.Println("Initializing Alpaca client...")
	AlpacaClient = stockslambdautils.AlpacaClient{Client: alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    config.Alpaca.ApiKey,
		APISecret: config.Alpaca.ApiSecret,
		BaseURL:   config.Alpaca.PaperApiUrl,
	})}
}
