package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchPrices_EmptyTickers(t *testing.T) {
	prices, err := FetchPrices([]string{}, "")
	if err != nil {
		t.Fatalf("expected no error for empty tickers, got %v", err)
	}
	if len(prices) != 0 {
		t.Errorf("expected empty map, got %v", prices)
	}
}

func TestFetchPrices_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(brapiResponse{
			Results: []brapiResult{
				{Symbol: "VALE3", RegularMarketPrice: 71.40},
				{Symbol: "MXRF11", RegularMarketPrice: 9.85},
			},
		})
	}))
	defer server.Close()

	origURL := brapiBaseURL
	brapiBaseURL = server.URL
	defer func() { brapiBaseURL = origURL }()

	prices, err := FetchPrices([]string{"VALE3", "MXRF11"}, "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prices["VALE3"] != 71.40 {
		t.Errorf("expected VALE3=71.40, got %.2f", prices["VALE3"])
	}
	if prices["MXRF11"] != 9.85 {
		t.Errorf("expected MXRF11=9.85, got %.2f", prices["MXRF11"])
	}
}

func TestFetchPrices_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer server.Close()

	origURL := brapiBaseURL
	brapiBaseURL = server.URL
	defer func() { brapiBaseURL = origURL }()

	_, err := FetchPrices([]string{"VALE3"}, "")
	if err == nil {
		t.Error("expected error for 500 response, got nil")
	}
}
