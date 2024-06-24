package configuration

import (
	"encoding/json"
	"log"
	"os"

	stockslambdautils "github.com/MikeB1124/stocks-lambda-utils/v2"
)

type Configration struct {
	MongoDB MongoDB `json:"mongodb"`
	Alpaca  Alpaca  `json:"alpaca"`
}

type MongoDB struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Alpaca struct {
	PaperApiUrl string `json:"paperApiUrl"`
	ApiKey      string `json:"apiKey"`
	ApiSecret   string `json:"apiSecret"`
}

var Config Configration
var MongoClient stockslambdautils.MongoClient
var AlpacaClient stockslambdautils.AlpacaClient

func init() {
	log.Println("Loading configuration...")
	sharedSecretName := os.Getenv("SHARED_SECRETS")
	if sharedSecretName == "" {
		log.Fatal("SHARED_SECRETS environment variable is required")
	}

	secrets, err := stockslambdautils.AWSGetSecrets(sharedSecretName)
	if err != nil {
		log.Fatal(err)
	}

	var lambdaConfig Configration
	err = json.Unmarshal([]byte(secrets), &lambdaConfig)
	if err != nil {
		log.Fatal(err)
	}

	MongoClient, err = stockslambdautils.NewMongoClient(lambdaConfig.MongoDB.Username, lambdaConfig.MongoDB.Password)
	if err != nil {
		log.Fatal(err)
	}
	AlpacaClient = stockslambdautils.NewAlpacaClient(lambdaConfig.Alpaca.ApiKey, lambdaConfig.Alpaca.ApiSecret, lambdaConfig.Alpaca.PaperApiUrl)
	Config = lambdaConfig
}

func GetConfig() Configration {
	return Config
}
