package domain

import "github.com/shopspring/decimal"

type DepthSnapshot struct {
	LastUpdateID int64               `json:"lastUpdateId"` // Last update ID
	Bids         [][]decimal.Decimal `json:"bids"`         // Bids
	Asks         [][]decimal.Decimal `json:"asks"`         // Asks

}
