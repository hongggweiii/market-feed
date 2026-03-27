package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hongggweiii/market-feed/internal/domain"
)

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

	orderBook := new(domain.DepthSnapshot)
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading from response body: %v", err)
		}

		err = json.Unmarshal(bodyBytes, orderBook)
		fmt.Println("Successful fetch!")
	} else {
		return nil, fmt.Errorf("Request failed with status: %d", resp.StatusCode)

	}

	return orderBook, nil
}
