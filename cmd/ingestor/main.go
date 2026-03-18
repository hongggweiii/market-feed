package main

import (
	"log"

	"github.com/hongggweiii/market-feed/internal/broker"
)

func main() {
	err := broker.PrepareKafkaTopic("localhost:9092", "crypto.trades.raw")
	if err != nil {
		log.Fatalf("Failed to create Kafka topic: %v", err)
	}
}
