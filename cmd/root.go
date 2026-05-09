package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gostock",
	Short: "br stock management in terminal",
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}
