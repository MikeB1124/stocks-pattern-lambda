#!/bin/bash


curl -X POST "http://localhost:9000/2015-03-31/functions/function/invocations" \
  --data '{"httpMethod": "POST", "path": "/harmonic-pattern", "body": "{\"msg_type\": \"pattern.notification\", \"data\": [{\"patterntype\": \"bearish\", \"patternname\": \"deep crab\", \"profit1\": 107.1265, \"displaySymbol\": \"USD/JPY\", \"symbol\": \"OANDA:USD_JPY\", \"profit2\": 106.9315, \"stoploss\": 107.5794, \"url\": \"https://harmonicpattern.com/pattern#noti/5\", \"timeframe\": \"H1\", \"status\": \"complete\", \"entry\": \"107.4054_107.4060\", \"patternclass\": \"harmonic\"}]}"}' \
