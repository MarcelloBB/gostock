package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/MarcelloBB/gostock/internal/database"
	"github.com/MarcelloBB/gostock/internal/models"
	"github.com/MarcelloBB/gostock/internal/service"
	"github.com/MarcelloBB/gostock/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

var RootCmd = &cobra.Command{
	Use:   "gostock",
	Short: "BR stock management in terminal",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initDB()
	},
	RunE: runMonitor,
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	home, _ := os.UserHomeDir()
	configFile := filepath.Join(home, ".gostock.yaml")
	viper.SetConfigFile(configFile)
	_ = viper.ReadInConfig()
}

func initDB() error {
	var err error
	db, err = database.Connect()
	return err
}

func runMonitor(cmd *cobra.Command, args []string) error {
	token := viper.GetString("token")

	for {
		fmt.Print("\033[H\033[2J")

		var transactions []models.Transaction
		if err := db.Order("created_at asc").Find(&transactions).Error; err != nil {
			return err
		}

		assets := service.CalculatePortfolio(transactions)

		if len(assets) > 0 {
			tickers := make([]string, 0, len(assets))
			for t := range assets {
				tickers = append(tickers, t)
			}
			sort.Strings(tickers)

			prices, err := service.FetchPrices(tickers, token)
			if err == nil {
				for ticker, price := range prices {
					if asset, ok := assets[ticker]; ok {
						asset.CurrentPrice = price
						asset.TotalValue = asset.TotalQty * price
						if asset.AveragePrice > 0 {
							asset.ProfitLoss = ((price - asset.AveragePrice) / asset.AveragePrice) * 100
						}
					}
				}
			}
		}

		ui.RenderTable(assets)
		time.Sleep(30 * time.Second)
	}
}
