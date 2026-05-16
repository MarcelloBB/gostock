package cmd

import (
	"fmt"

	"github.com/MarcelloBB/gostock/internal/repository"
	"github.com/MarcelloBB/gostock/internal/service"
	"github.com/MarcelloBB/gostock/internal/ui"
	"github.com/spf13/cobra"
)

var removeID uint

func init() {
	historyCmd.Flags().UintVar(&removeID, "remove", 0, "ID da transação a remover")
	RootCmd.AddCommand(historyCmd)
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "List transaction history",
	RunE:  runHistoryCmd,
}

func runHistoryCmd(cmd *cobra.Command, args []string) error {
	repo := repository.NewTransactionRepository(db)
	svc := service.NewTransactionService(repo)

	if removeID > 0 {
		if err := svc.RemoveTransaction(removeID); err != nil {
			return fmt.Errorf("erro ao remover transação #%d: %w", removeID, err)
		}
		fmt.Printf("Transação #%d removida. O monitor recalculará o portfólio automaticamente.\n", removeID)
		return nil
	}

	transactions, err := svc.GetHistory()
	if err != nil {
		return err
	}

	ui.RenderHistory(transactions)
	return nil
}
