package repository

import (
	"testing"
	"time"

	"github.com/MarcelloBB/gostock/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestTransactionRepository_Create(t *testing.T) {
	repo := NewTransactionRepository(newTestDB(t))
	tx := &models.Transaction{Ticker: "VALE3", Quantity: 100, Price: 65.00, Type: "BUY", CreatedAt: time.Now()}

	if err := repo.Create(tx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.ID == 0 {
		t.Error("expected ID to be set after create")
	}
}

func TestTransactionRepository_FindAll(t *testing.T) {
	repo := NewTransactionRepository(newTestDB(t))
	now := time.Now()
	_ = repo.Create(&models.Transaction{Ticker: "VALE3", Quantity: 100, Price: 65.00, Type: "BUY", CreatedAt: now})
	_ = repo.Create(&models.Transaction{Ticker: "MXRF11", Quantity: 50, Price: 10.20, Type: "BUY", CreatedAt: now.Add(time.Hour)})

	transactions, err := repo.FindAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(transactions) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(transactions))
	}
}

func TestTransactionRepository_FindByTicker(t *testing.T) {
	repo := NewTransactionRepository(newTestDB(t))
	now := time.Now()
	_ = repo.Create(&models.Transaction{Ticker: "VALE3", Quantity: 100, Price: 65.00, Type: "BUY", CreatedAt: now})
	_ = repo.Create(&models.Transaction{Ticker: "MXRF11", Quantity: 50, Price: 10.20, Type: "BUY", CreatedAt: now.Add(time.Hour)})

	transactions, err := repo.FindByTicker("VALE3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(transactions) != 1 {
		t.Errorf("expected 1 transaction for VALE3, got %d", len(transactions))
	}
	if transactions[0].Ticker != "VALE3" {
		t.Errorf("expected ticker VALE3, got %s", transactions[0].Ticker)
	}
}

func TestTransactionRepository_Delete(t *testing.T) {
	repo := NewTransactionRepository(newTestDB(t))
	tx := &models.Transaction{Ticker: "VALE3", Quantity: 100, Price: 65.00, Type: "BUY", CreatedAt: time.Now()}
	_ = repo.Create(tx)

	if err := repo.Delete(tx.ID); err != nil {
		t.Fatalf("unexpected error on delete: %v", err)
	}

	transactions, _ := repo.FindAll()
	if len(transactions) != 0 {
		t.Errorf("expected 0 transactions after delete, got %d", len(transactions))
	}
}
