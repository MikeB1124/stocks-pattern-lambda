#!/bin/bash

docker stop stocks-pattern-lambda
docker rm stocks-pattern-lambda
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
docker build -t stocks-lambda-image .
docker run --name stocks-pattern-lambda -p 9000:8080 --env-file .env stocks-lambda-image