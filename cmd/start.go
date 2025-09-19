package cmd

import (
	"github.com/nullsec45/golang-news-api/internal/app"
	"github.com/spf13/cobra"
)

var startCmd=&cobra.Command{
	Use:"start",
	Short:"start",
	Long:"start",
	Run: func(cmd *cobra.Command, args []string) {
		app.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}