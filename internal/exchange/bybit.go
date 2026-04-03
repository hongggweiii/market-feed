package exchange

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hongggweiii/market-nodes/internal/domain"
	"github.com/shopspring/decimal"
)

type ByBitClient struct{}

type byBitWSSubscribeDTO struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

type bybitDepthUpdateDTO struct {
	Type string `json:"type"` // "snapshot" or "delta"
	Data struct {
		Bids [][]string `json:"b"` // [price, size]
		Asks [][]string `json:"a"` // [price, size]
	} `json:"data"`
}

func (c *ByBitClient) Name() string {
	return "ByBit"
}

func (c *ByBitClient) FetchDepthSnapshot(symbol string) (*domain.DepthSnapshot, error) {
	// Bybit's WS sends a full snapshot upon connecting, no need to fetch separately via REST
	return &domain.DepthSnapshot{LastUpdateID: 0}, nil
}

func (c *ByBitClient) StreamOrderBookDepthUpdates(symbol string, updates chan<- *domain.DepthUpdate) error {
	wsUrl := "wss://stream.bybit.com/v5/public/spot"

	// Initialise Websocket connection
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		fmt.Printf("Error while initialising Websocket connection: %v", err)
		return err
	}
	defer conn.Close()

	// Send subscription message to Websocket
	subMsg := byBitWSSubscribeDTO{
		Op:   "subscribe",
		Args: []string{fmt.Sprintf("orderbook.50.%s", symbol)}, // 50 levels deep
	}
	if err := conn.WriteJSON(subMsg); err != nil {
		fmt.Printf("Error while subscribing to channel: %v", err)
		return err
	}

	fmt.Printf("[ByBit] Connected and subscribed to %s\n", symbol)

	// Infinite loop for Websocket
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		dto := new(bybitDepthUpdateDTO)
		if err := json.Unmarshal(p, dto); err != nil {
			continue
		}

		// Ignore non-orderbook messages (like the connection success message)
		if dto.Type != "snapshot" && dto.Type != "delta" {
			continue
		}

		// Parse the depth update changes
		bids := make(map[string]decimal.Decimal)
		asks := make(map[string]decimal.Decimal)

		for _, b := range dto.Data.Bids {
			price := b[0]
			size, _ := decimal.NewFromString(b[1])
			bids[price] = size
		}
		for _, a := range dto.Data.Asks {
			price := a[0]
			size, _ := decimal.NewFromString(a[1])
			asks[price] = size
		}

		updates <- &domain.DepthUpdate{
			EventType:     dto.Type,
			EventTime:     time.Now().UnixMilli(),
			Symbol:        symbol,
			FirstUpdateID: 0, // Sets 0 for bypass at engine.go since Bybit doesn't provide update IDs
			FinalUpdateID: 0,
			Bids:          bids,
			Asks:          asks,
		}
	}

	return nil
}
