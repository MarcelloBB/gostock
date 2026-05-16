package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/MarcelloBB/gostock/internal/models"
)

const tableWidth = 78

func RenderTable(assets map[string]*models.Asset) {
	now := time.Now().Format("15:04:05")
	fmt.Printf("\n  GoStock Monitor · Atualizado às %s\n", HeaderStyle.Render(now))

	sep := TableBorderStyle.Render(strings.Repeat("─", tableWidth))
	fmt.Println("  " + sep)
	fmt.Printf("  %-8s  %-8s  %-10s  %-10s  %-13s  %s\n",
		HeaderStyle.Render("Ticker"),
		HeaderStyle.Render("Qtd"),
		HeaderStyle.Render("PM (R$)"),
		HeaderStyle.Render("Cotação"),
		HeaderStyle.Render("Patrimônio"),
		HeaderStyle.Render("Resultado"),
	)
	fmt.Println("  " + sep)

	tickers := make([]string, 0, len(assets))
	for t := range assets {
		tickers = append(tickers, t)
	}
	sort.Strings(tickers)

	var total float64
	for _, ticker := range tickers {
		a := assets[ticker]
		total += a.TotalValue

		cotacao := "N/A"
		result := "N/A"
		if a.CurrentPrice > 0 {
			cotacao = fmt.Sprintf("R$ %.2f", a.CurrentPrice)
			pct := fmt.Sprintf("%.2f%%", a.ProfitLoss)
			if a.ProfitLoss >= 0 {
				result = ProfitStyle.Render("+" + pct)
			} else {
				result = LossStyle.Render(pct)
			}
		}

		fmt.Printf("  %-8s  %8.2f  %10.2f  %-10s  R$ %9.2f  %s\n",
			ticker, a.TotalQty, a.AveragePrice, cotacao, a.TotalValue, result)
	}

	fmt.Println("  " + sep)
	fmt.Printf("  %-8s  %-8s  %-10s  %-10s  R$ %9.2f\n",
		"Total", "", "", "", total)
	fmt.Println()
}

func RenderHistory(transactions []models.Transaction) {
	if len(transactions) == 0 {
		fmt.Println("Nenhuma transação encontrada.")
		return
	}

	sep := TableBorderStyle.Render(strings.Repeat("─", 62))
	fmt.Println("\n" + sep)
	fmt.Printf("%-4s  %-8s  %-4s  %-8s  %-12s  %s\n",
		HeaderStyle.Render("ID"),
		HeaderStyle.Render("Ticker"),
		HeaderStyle.Render("Tipo"),
		HeaderStyle.Render("Qtd"),
		HeaderStyle.Render("Preço"),
		HeaderStyle.Render("Data"),
	)
	fmt.Println(sep)

	for _, t := range transactions {
		typeStyle := lipgloss.NewStyle()
		if t.Type == "SELL" {
			typeStyle = LossStyle
		} else {
			typeStyle = ProfitStyle
		}
		fmt.Printf("%-4d  %-8s  %-4s  %8.2f  R$ %8.2f  %s\n",
			t.ID, t.Ticker, typeStyle.Render(t.Type),
			t.Quantity, t.Price,
			t.CreatedAt.Format("2006-01-02"),
		)
	}

	fmt.Println(sep)
	fmt.Println()
}
