/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Transfer tokens from one account to another",
	Long: `Transfer tokens from one account to another

	custom_blockchain send --from [sender-address] --to [reciever-address] --amount [number of tokens]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		from := cmd.PersistentFlags().Lookup("from").Value.String()
		to := cmd.PersistentFlags().Lookup("to").Value.String()
		amount, err := cmd.PersistentFlags().GetInt("amount")
		handleErr(err, nil)

		bc, err := blockchain.ContinueBlockchain(true, nil)
		handleErr(err, bc)

		txn, err := blockchain.NewTransaction(from, to, amount, bc)
		handleErr(err, bc)

		err = bc.AddBlock([]blockchain.Transaction{*txn})
		handleErr(err, bc)

		log.Println("Successful!")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.PersistentFlags().String("from", "", "sender's private address")
	sendCmd.MarkPersistentFlagRequired("from")

	sendCmd.PersistentFlags().String("to", "", "reciever's public address")
	sendCmd.MarkPersistentFlagRequired("to")

	sendCmd.PersistentFlags().Int("amount", 0, "number of tokens to be transferred")
	sendCmd.MarkPersistentFlagRequired("amount")
}
