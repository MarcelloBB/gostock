package repository

import (
	"github.com/MarcelloBB/gostock/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(t *models.Transaction) error {
	return r.db.Create(t).Error
}

func (r *TransactionRepository) FindAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Order("created_at asc").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) FindByTicker(ticker string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("ticker = ?", ticker).Order("created_at asc").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}
