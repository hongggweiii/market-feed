package domain

type Trade struct {
	EventType     string `json:"e"` // Event type
	EventTime     int64  `json:"E"` // Event time
	Symbol        string `json:"s"` // Symbol
	TradeID       int64  `json:"t"` // Trade ID
	Price         string `json:"p"` // Price
	Quantity      string `json:"q"` // Quantity
	TradeTime     int64  `json:"T"` // Trade time
	IsMarketMaker bool   `json:"m"` // Is the buyer the market maker?
}
