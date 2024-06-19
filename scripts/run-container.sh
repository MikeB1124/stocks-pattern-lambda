#!/bin/bash

docker stop stock-pattern-lambda
docker rm stock-pattern-lambda
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
docker build -t stock-pattern-lambda-image .
docker run --name stock-pattern-lambda -p 9000:8080 --env-file .env stock-pattern-lambda-image