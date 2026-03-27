package domain

import "github.com/shopspring/decimal"

type Trade struct {
	EventType     string          // Event type
	EventTime     int64           // Event time
	Symbol        string          // Symbol
	TradeID       int64           // Trade ID
	Price         decimal.Decimal // Price
	Quantity      decimal.Decimal // Quantity
	TradeTime     int64           // Trade time
	IsMarketMaker bool            // Is the buyer the market maker?
}
