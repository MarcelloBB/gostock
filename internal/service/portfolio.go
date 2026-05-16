package service

import "github.com/MarcelloBB/gostock/internal/models"

// CalculatePortfolio percorre transações em ordem cronológica e calcula
// o preço médio ponderado por ativo. Segue padrão fiscal B3:
// vendas reduzem apenas quantidade, PM das posições restantes não muda.
func CalculatePortfolio(transactions []models.Transaction) map[string]*models.Asset {
	assets := make(map[string]*models.Asset)

	for _, t := range transactions {
		if _, ok := assets[t.Ticker]; !ok {
			assets[t.Ticker] = &models.Asset{Ticker: t.Ticker}
		}
		asset := assets[t.Ticker]

		if t.Type == "BUY" {
			totalCost := asset.AveragePrice*asset.TotalQty + t.Quantity*t.Price
			asset.TotalQty += t.Quantity
			if asset.TotalQty > 0 {
				asset.AveragePrice = totalCost / asset.TotalQty
			}
		} else if t.Type == "SELL" {
			asset.TotalQty -= t.Quantity
		}
	}

	for ticker, asset := range assets {
		if asset.TotalQty <= 0 {
			delete(assets, ticker)
		}
	}

	return assets
}
