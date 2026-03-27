package domain

import "github.com/shopspring/decimal"

type DepthSnapshot struct {
	LastUpdateID int64               // Last update ID
	Bids         [][]decimal.Decimal // Bids
	Asks         [][]decimal.Decimal // Asks

}
