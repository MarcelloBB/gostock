package model

type Transaction struct {
    ID        uint      `gorm:"primaryKey"`
    Ticker    string    `gorm:"index"`
    Quantity  float64
    Price     float64
    Type      string    // "BUY" / "SELL"
    CreatedAt time.Time
}

type Asset struct {
    Ticker       string
    TotalQty     float64
    AveragePrice float64
    CurrentPrice float64
    TotalValue   float64
    ProfitLoss   float64
}
