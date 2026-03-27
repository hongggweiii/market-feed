package exchange

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/hongggweiii/market-feed/internal/domain"
	"github.com/hongggweiii/market-feed/internal/ingestor/broker"
)

type binanceTradeDTO struct {
	EventType     string `json:"e"`
	EventTime     int64  `json:"E"`
	Symbol        string `json:"s"`
	TradeID       int64  `json:"t"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	TradeTime     int64  `json:"T"`
	IsMarketMaker bool   `json:"m"`
}

func StreamBinanceTrades(symbol string, broker *broker.KafkaProducer) error {
	baseUrl := "wss://stream.binance.com:9443/ws"
	lowercaseSymbol := strings.ToLower(symbol)
	websocketUrl := fmt.Sprintf("%s/%s@trade", baseUrl, lowercaseSymbol)

	// Initialise Websocket connection
	conn, _, err := websocket.DefaultDialer.Dial(websocketUrl, nil)
	if err != nil {
		fmt.Printf("Error while initialising Websocket connection: %v", err)
		return err
	}
	defer conn.Close()

	// Infinite loop for Websocket
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		dto := new(binanceTradeDTO)
		err = json.Unmarshal(p, dto)
		if err != nil {
			fmt.Println("Error while unserializing:", err)
			continue
		}

		trade := &domain.Trade{
			EventType:     dto.EventType,
			EventTime:     dto.EventTime,
			Symbol:        dto.Symbol,
			TradeID:       dto.TradeID,
			Price:         dto.Price,
			Quantity:      dto.Quantity,
			TradeTime:     dto.TradeTime,
			IsMarketMaker: dto.IsMarketMaker,
		}

		err = broker.PublishTrade(*trade)
		if err != nil {
			fmt.Printf("Error while publishing trade: %v", err)
		}

		fmt.Println("Published trade:", trade)
	}

	return nil
}
