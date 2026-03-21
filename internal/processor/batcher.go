package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/hongggweiii/market-feed/internal/broker"
	"github.com/hongggweiii/market-feed/internal/domain"
)

func insertBatch(batch []domain.Trade) error {
	fmt.Printf("Flushing batch of size %d to ClickHouse...", len(batch))
	return nil
}

func StartBatchingEngine(consumer *broker.KafkaConsumer, ctx context.Context) error {
	// Channel to pass trades between threads
	tradeChan := make(chan domain.Trade, 1000)

	go func() {
		for {
			trade, err := consumer.ConsumeTrade(ctx)
			if err != nil {
				fmt.Printf("Error while consuming trade: %v", err)
				continue
			}
			tradeChan <- trade // Send trade into channel
		}
	}() // Immediately executes, anonymous func

	var batch []domain.Trade
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		// Wait on multiple channels. Proceeds with case whose channel operation is ready first
		select {
		case trade := <-tradeChan: // Append trade to batch
			batch = append(batch, trade)
			if len(batch) >= 1000 {
				err := insertBatch(batch)
				if err != nil {
					fmt.Printf("Error while inserting batch to ClickHouse: %v", err)
				}
				batch = nil
			}
		case <-ticker.C: // Insert existing batch to database every 2s
			if len(batch) > 0 {
				err := insertBatch(batch)
				if err != nil {
					fmt.Printf("Error while inserting batch to ClickHouse: %v", err)
				}
				batch = nil
			}
		}
	}
}
