package configuration

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Configration struct {
	MongoDB MongoDB `json:"mongodb"`
}

type MongoDB struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Config Configration

func init() {
	log.Println("Loading configuration...")
	sharedSecretName := os.Getenv("SHARED_SECRETS")
	if sharedSecretName == "" {
		log.Fatal("SHARED_SECRETS environment variable is required")
	}

	secrets, err := getSecrets(sharedSecretName)
	if err != nil {
		log.Fatal(err)
	}

	var lambdaConfig Configration
	err = json.Unmarshal([]byte(secrets), &lambdaConfig)
	if err != nil {
		log.Fatal(err)
	}

	Config = lambdaConfig
}

func getSecrets(secretName string) (string, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	return *result.SecretString, nil
}

func GetConfig() Configration {
	return Config
}
