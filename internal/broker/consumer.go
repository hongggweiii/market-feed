package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hongggweiii/market-feed/internal/domain"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokerAddress string, topic string, groupId string) *KafkaConsumer {
	return &KafkaConsumer{
		reader: &kafka.Reader(kafka.ReaderConfig{
			Brokers: []string{brokerAddress},
			Topic:   topic,
			GroupID: groupId,
		}),
	}
}

func (c *KafkaConsumer) ConsumeTrade(ctx context.Context) (domain.Trade, error) {
	trade := new(domain.Trade)

	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(msg.Value, trade)
	if err != nil {
		fmt.Printf("Error while unserializing: %v", err)
		return nil, err
	}

	return *trade, nil
}
