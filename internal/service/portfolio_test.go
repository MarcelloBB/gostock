package service

import (
	"testing"
	"time"

	"github.com/MarcelloBB/gostock/internal/models"
)

func TestCalculatePortfolio_SingleBuy(t *testing.T) {
	transactions := []models.Transaction{
		{ID: 1, Ticker: "VALE3", Quantity: 100, Price: 65.00, Type: "BUY", CreatedAt: time.Now()},
	}

	assets := CalculatePortfolio(transactions)

	if len(assets) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(assets))
	}
	a := assets["VALE3"]
	if a.TotalQty != 100 {
		t.Errorf("expected qty=100, got %.2f", a.TotalQty)
	}
	if a.AveragePrice != 65.00 {
		t.Errorf("expected avg=65.00, got %.2f", a.AveragePrice)
	}
}

func TestCalculatePortfolio_MultipleBuysWeightedAverage(t *testing.T) {
	now := time.Now()
	transactions := []models.Transaction{
		{ID: 1, Ticker: "VALE3", Quantity: 100, Price: 60.00, Type: "BUY", CreatedAt: now},
		{ID: 2, Ticker: "VALE3", Quantity: 100, Price: 70.00, Type: "BUY", CreatedAt: now.Add(time.Hour)},
	}

	assets := CalculatePortfolio(transactions)

	a := assets["VALE3"]
	if a.TotalQty != 200 {
		t.Errorf("expected qty=200, got %.2f", a.TotalQty)
	}
	// (100*60 + 100*70) / 200 = 65.00
	if a.AveragePrice != 65.00 {
		t.Errorf("expected avg=65.00, got %.2f", a.AveragePrice)
	}
}

func TestCalculatePortfolio_SellKeepsAveragePrice(t *testing.T) {
	now := time.Now()
	transactions := []models.Transaction{
		{ID: 1, Ticker: "VALE3", Quantity: 100, Price: 65.00, Type: "BUY", CreatedAt: now},
		{ID: 2, Ticker: "VALE3", Quantity: 20, Price: 72.00, Type: "SELL", CreatedAt: now.Add(time.Hour)},
	}

	assets := CalculatePortfolio(transactions)

	a := assets["VALE3"]
	if a.TotalQty != 80 {
		t.Errorf("expected qty=80, got %.2f", a.TotalQty)
	}
	if a.AveragePrice != 65.00 {
		t.Errorf("expected avg unchanged at 65.00, got %.2f", a.AveragePrice)
	}
}

func TestCalculatePortfolio_ZeroQtyAssetRemoved(t *testing.T) {
	now := time.Now()
	transactions := []models.Transaction{
		{ID: 1, Ticker: "VALE3", Quantity: 100, Price: 65.00, Type: "BUY", CreatedAt: now},
		{ID: 2, Ticker: "VALE3", Quantity: 100, Price: 72.00, Type: "SELL", CreatedAt: now.Add(time.Hour)},
	}

	assets := CalculatePortfolio(transactions)

	if len(assets) != 0 {
		t.Errorf("expected empty portfolio after full sell, got %d assets", len(assets))
	}
}

func TestCalculatePortfolio_MultipleAssets(t *testing.T) {
	now := time.Now()
	transactions := []models.Transaction{
		{ID: 1, Ticker: "VALE3", Quantity: 100, Price: 65.00, Type: "BUY", CreatedAt: now},
		{ID: 2, Ticker: "MXRF11", Quantity: 50, Price: 10.20, Type: "BUY", CreatedAt: now.Add(time.Hour)},
	}

	assets := CalculatePortfolio(transactions)

	if len(assets) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(assets))
	}
	if assets["VALE3"].TotalQty != 100 {
		t.Errorf("VALE3: expected qty=100, got %.2f", assets["VALE3"].TotalQty)
	}
	if assets["MXRF11"].TotalQty != 50 {
		t.Errorf("MXRF11: expected qty=50, got %.2f", assets["MXRF11"].TotalQty)
	}
}
