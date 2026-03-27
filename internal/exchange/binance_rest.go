package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hongggweiii/market-feed/internal/domain"
	"github.com/shopspring/decimal"
)

type binanceDepthSnapshotDTO struct {
	LastUpdateID int64               `json:"lastUpdateId"`
	Bids         [][]decimal.Decimal `json:"bids"`
	Asks         [][]decimal.Decimal `json:"asks"`
}

func FetchDepthSnapshot(symbol string) (*domain.DepthSnapshot, error) {
	const limit = 1000
	baseUrl := "https://api.binance.com"
	lowercaseSymbol := strings.ToUpper(symbol)
	restUrl := fmt.Sprintf("%s/api/v3/depth?symbol=%s&limit=%d", baseUrl, lowercaseSymbol, limit)

	resp, err := http.Get(restUrl)
	if err != nil {
		fmt.Println("Error fetching depth:", err)
		return nil, err
	}
	defer resp.Body.Close() // Prevent resource leaks

	dto := new(binanceDepthSnapshotDTO)
	if resp.StatusCode == http.StatusOK {
		// io.ReadAll() take sup lots of memory
		if err := json.NewDecoder(resp.Body).Decode(dto); err != nil {
			return nil, fmt.Errorf("Failed to decode response: %w", err)
		}
	} else {
		return nil, fmt.Errorf("Request failed with status: %d", resp.StatusCode)

	}

	// Map DTO to Domain model
	snapshot := &domain.DepthSnapshot{
		LastUpdateID: dto.LastUpdateID,
		Bids:         dto.Bids,
		Asks:         dto.Asks,
	}

	fmt.Println("Successful fetch!")
	return snapshot, nil
}
