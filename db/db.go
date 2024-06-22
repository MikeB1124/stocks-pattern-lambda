package db

import (
	"context"
	"fmt"

	"github.com/MikeB1124/stocks-pattern-lambda/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	config := configuration.GetConfig()
	// Connect to MongoDB
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.du0vf.mongodb.net", config.MongoDB.Username, config.MongoDB.Password))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	mongoClient = client
}

func InsertEntryOrder(order AlpacaEntryOrder) error {
	collection := mongoClient.Database("Stocks").Collection("orders")
	_, err := collection.InsertOne(context.TODO(), order)
	return err
}
