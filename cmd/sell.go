package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MarcelloBB/gostock/internal/repository"
	"github.com/MarcelloBB/gostock/internal/service"
	"github.com/MarcelloBB/gostock/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(sellCmd)
}

var sellCmd = &cobra.Command{
	Use:   "sell <ticker> <quantidade>",
	Short: "Register a stock sale",
	Args:  cobra.ExactArgs(2),
	RunE:  runSellCmd,
}

func runSellCmd(cmd *cobra.Command, args []string) error {
	ticker := strings.ToUpper(args[0])
	qty, err := strconv.ParseFloat(args[1], 64)
	if err != nil || qty <= 0 {
		return fmt.Errorf("quantidade inválida: %s", args[1])
	}

	repo := repository.NewTransactionRepository(db)
	svc := service.NewTransactionService(repo)

	sufficient, err := svc.HasSufficientBalance(ticker, qty)
	if err != nil {
		return err
	}
	if !sufficient {
		return fmt.Errorf("saldo insuficiente para venda de %.0f %s", qty, ticker)
	}

	token := viper.GetString("token")
	prices, err := service.FetchPrices([]string{ticker}, token)
	if err != nil {
		return fmt.Errorf("erro ao buscar cotação de %s: %w", ticker, err)
	}
	currentPrice, ok := prices[ticker]
	if !ok {
		return fmt.Errorf("ticker %s não encontrado", ticker)
	}

	_, avgPrice, err := svc.GetPortfolioForTicker(ticker)
	if err != nil {
		return err
	}

	profitLoss := ((currentPrice - avgPrice) / avgPrice) * 100

	card := ui.CardStyle.Render(fmt.Sprintf(
		"Ativo: %-10s | Qtd: %.0f\n"+
			"Preço de Mercado:    R$ %.2f\n"+
			"Seu PM Atual:        R$ %.2f\n"+
			"Resultado Previsto:  %.2f%%\n"+
			"Total a Receber:     R$ %.2f",
		ticker, qty, currentPrice, avgPrice, profitLoss, currentPrice*qty,
	))
	fmt.Println(card)

	fmt.Print("Confirmar transação? (y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
		fmt.Println("Operação cancelada.")
		return nil
	}

	if err := svc.RecordSell(ticker, qty, currentPrice); err != nil {
		return fmt.Errorf("erro ao registrar venda: %w", err)
	}

	fmt.Printf("Venda de %.0f %s registrada com sucesso!\n", qty, ticker)
	return nil
}
