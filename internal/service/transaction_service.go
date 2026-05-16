package service

import (
	"github.com/MarcelloBB/gostock/internal/models"
	"github.com/MarcelloBB/gostock/internal/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) RecordBuy(ticker string, qty, price float64) error {
	return s.repo.Create(&models.Transaction{
		Ticker:   ticker,
		Quantity: qty,
		Price:    price,
		Type:     "BUY",
	})
}

func (s *TransactionService) RecordSell(ticker string, qty, price float64) error {
	return s.repo.Create(&models.Transaction{
		Ticker:   ticker,
		Quantity: qty,
		Price:    price,
		Type:     "SELL",
	})
}

func (s *TransactionService) GetHistory() ([]models.Transaction, error) {
	return s.repo.FindAll()
}

func (s *TransactionService) RemoveTransaction(id uint) error {
	return s.repo.Delete(id)
}

// GetPortfolioForTicker retorna quantidade líquida e preço médio atual de um ativo.
func (s *TransactionService) GetPortfolioForTicker(ticker string) (qty, avgPrice float64, err error) {
	transactions, err := s.repo.FindByTicker(ticker)
	if err != nil {
		return 0, 0, err
	}
	assets := CalculatePortfolio(transactions)
	if asset, ok := assets[ticker]; ok {
		return asset.TotalQty, asset.AveragePrice, nil
	}
	return 0, 0, nil
}

func (s *TransactionService) HasSufficientBalance(ticker string, qty float64) (bool, error) {
	current, _, err := s.GetPortfolioForTicker(ticker)
	if err != nil {
		return false, err
	}
	return current >= qty, nil
}
