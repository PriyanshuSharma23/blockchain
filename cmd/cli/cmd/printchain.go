/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
	"github.com/spf13/cobra"
)

// printchainCmd represents the printchain command
var printchainCmd = &cobra.Command{
	Use:   "printchain",
	Short: "Print the blocks in a blockchain",
	Long:  `Print the blocks in a blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		bc, err := blockchain.ContinueBlockchain(true, nil)

		handleErr(err, bc)

		it := bc.Iterator()

		for {
			block, err := it.Next()
			handleErr(err, bc)

			if block == nil {
				break
			}

			fmt.Printf("Hash of the block %x\n", block.Hash)
			fmt.Printf("Hash of the previous Block: %x\n", block.PrevHash)

			pow := blockchain.NewProof(block)
			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(printchainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// printchainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// printchainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
