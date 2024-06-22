#!/bin/bash


curl -X POST "http://localhost:9000/2015-03-31/functions/function/invocations" \
  --data '{"httpMethod": "POST", "path": "/webhook/harmonic-pattern", "body": "{\"msg_type\": \"pattern.notification\", \"data\": [{\"patterntype\": \"bullish\", \"patternname\": \"gartley\", \"profit1\": 444.9544, \"displaySymbol\": \"MSFT\", \"symbol\": \"MSFT.US\", \"stoploss\": 440.23, \"url\": \"https://harmonicpattern.com/pattern#noti/14978988\", \"timeframe\": \"H1\", \"status\": \"complete\", \"entry\": \"442.9228_443.2655\", \"patternclass\": \"harmonic\"}]}"}' \
