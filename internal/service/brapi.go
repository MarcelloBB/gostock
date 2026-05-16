package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var brapiBaseURL = "https://brapi.dev/api/quote"

var httpClient = &http.Client{Timeout: 10 * time.Second}

type brapiResult struct {
	Symbol             string  `json:"symbol"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
}

type brapiResponse struct {
	Results []brapiResult `json:"results"`
}

func FetchPrices(tickers []string, token string) (map[string]float64, error) {
	if len(tickers) == 0 {
		return map[string]float64{}, nil
	}

	url := fmt.Sprintf("%s/%s", brapiBaseURL, strings.Join(tickers, ","))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brapi returned status %d", resp.StatusCode)
	}

	var result brapiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	prices := make(map[string]float64, len(result.Results))
	for _, r := range result.Results {
		prices[r.Symbol] = r.RegularMarketPrice
	}
	return prices, nil
}
