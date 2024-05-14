package cmd

import (
	"log"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
	"github.com/spf13/cobra"
)

// createblockchainCmd represents the createblockchain command
var createblockchainCmd = &cobra.Command{
	Use:   "createblockchain",
	Short: "Create a new blockchain",
	Long: `Create a new blockchain
	
	Creates the genesis block and the coinbase transaction. It recieves an address
	which determines who can mine the coinbsae transaction tokens.`,
	Run: func(cmd *cobra.Command, args []string) {
		address := cmd.PersistentFlags().Lookup("address").Value.String()
		log.Println("Address provided:", address)

		bc, err := blockchain.CreateBlockchain(address, true, nil)
		handleErr(err, bc)

		log.Println("Blockchain created")
		log.Println("Blockchain lasthash:", bc.LastHash)
	},
}

func init() {
	rootCmd.AddCommand(createblockchainCmd)

	createblockchainCmd.PersistentFlags().StringP("address", "A", "", "address which mines the coinbase transaction tokens")
	createblockchainCmd.MarkPersistentFlagRequired("address")
}
