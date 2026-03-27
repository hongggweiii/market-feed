package domain

type Trade struct {
	EventType     string // Event type
	EventTime     int64  // Event time
	Symbol        string // Symbol
	TradeID       int64  // Trade ID
	Price         string // Price
	Quantity      string // Quantity
	TradeTime     int64  // Trade time
	IsMarketMaker bool   // Is the buyer the market maker?
}
