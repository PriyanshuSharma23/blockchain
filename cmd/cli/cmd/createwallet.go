/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/wallet"
	"github.com/spf13/cobra"
)

// createwalletCmd represents the createwallet command
var createwalletCmd = &cobra.Command{
	Use:   "createwallet",
	Short: "Create a new wallet",

	Run: func(cmd *cobra.Command, args []string) {
		ws := wallet.NewWallets()

		w, err := wallet.NewWallet()
		handleErr(err, nil)

		ws.AddWallet(w)
		fmt.Printf("Successfully created wallet: %s\n", w.Address())
	},
}

func init() {
	rootCmd.AddCommand(createwalletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createwalletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createwalletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
