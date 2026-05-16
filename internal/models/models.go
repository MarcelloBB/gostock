package models

import "time"

// Transaction é a única entidade gravada no banco.
// Saldos e preços médios são calculados dinamicamente a cada ciclo.
type Transaction struct {
	ID        uint      `gorm:"primaryKey"`
	Ticker    string    `gorm:"index;not null"`
	Quantity  float64   `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Type      string    `gorm:"not null"` // "BUY" ou "SELL"
	CreatedAt time.Time
}

// Asset é calculado em memória — nunca persiste no banco.
type Asset struct {
	Ticker       string
	TotalQty     float64
	AveragePrice float64
	CurrentPrice float64
	TotalValue   float64
	ProfitLoss   float64 // percentual em relação ao PM
}
