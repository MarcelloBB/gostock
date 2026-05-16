package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiToken string

func init() {
	configCmd.Flags().StringVar(&apiToken, "token", "", "Brapi API token")
	_ = configCmd.MarkFlagRequired("token")
	RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Brapi API token",
	RunE:  runConfigCmd,
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(home, ".gostock.yaml")
	viper.SetConfigFile(configFile)
	viper.Set("token", apiToken)

	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("falha ao salvar configuração: %w", err)
	}

	fmt.Printf("Token salvo em %s\n", configFile)
	return nil
}
