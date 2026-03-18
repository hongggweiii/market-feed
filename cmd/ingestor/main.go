package main

import (
	"log"

	"github.com/hongggweiii/market-feed/internal/broker"
	"github.com/hongggweiii/market-feed/internal/exchange"
)

func main() {
	err := broker.PrepareKafkaTopic("localhost:9092", "crypto.trades.raw")
	if err != nil {
		log.Fatalf("Failed to create Kafka topic: %v", err)
	}

	err = exchange.StreamBinanceTrades("BTCUSDT")
	if err != nil {
		log.Fatalf("Stream stopped: %v", err)
	}
}
