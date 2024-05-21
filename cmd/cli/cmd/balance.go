/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
	"github.com/PriyanshuSharma23/custom_blockchain/internals/wallet"
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance [address]",
	Short: "Retrieve the balance of a account",
	Long:  `Retrieve the balance of a account`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]

		// if !wallet.ValidateAddress(address) {
		// 	handleErr(fmt.Errorf("the address is invalid"), nil)
		// }

		bc, err := blockchain.ContinueBlockchain(true, nil)
		handleErr(err, bc)

		pubKeyHash, err := wallet.PublicKeyHashFromAddr(address)
		handleErr(err, bc)

		txnOuts, err := bc.FindUTXO(pubKeyHash)
		handleErr(err, bc)

		balance := 0
		for _, out := range txnOuts {
			balance += out.Value
		}

		log.Printf("Balance in %s is %d\n", address, balance)
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
