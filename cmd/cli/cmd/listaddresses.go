/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/wallet"
	"github.com/spf13/cobra"
)

// listaddressesCmd represents the listaddresses command
var listaddressesCmd = &cobra.Command{
	Use:   "listaddresses",
	Short: "List addresses of all the wallet present in memory",

	Run: func(cmd *cobra.Command, args []string) {
		ws := wallet.NewWallets()

		wSlice := ws.GetAllWallets()

		fmt.Println("Wallets: ")
		for _, w := range wSlice {
			fmt.Printf("  %s\n", w.Address())
		}
	},
}

func init() {
	rootCmd.AddCommand(listaddressesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listaddressesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listaddressesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
