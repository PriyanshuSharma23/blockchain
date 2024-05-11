package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Cli) PrintCmd() *cobra.Command {
	var printCmd = &cobra.Command{
		Use:   "print",
		Short: "Print all the blocks in the chain",
		Long: `The "print" command displays all blocks in the blockchain, providing insights into its structure and transaction history. It iterates through each block, showing its hash and associated data. 

Example:
$ custom_blockchain print

This command outputs each block's hash and data, aiding in blockchain analysis and verification.`,
		Run: func(cmd *cobra.Command, args []string) {
			it := c.bc.Iterator()

			for {
				block, err := it.Next()

				if err != nil {
					panic(err)
				}

				if block == nil {
					break
				}

				fmt.Printf("Hash of the block %x; Value: %s\n", block.Hash, block.Data)
				// fmt.Printf("Hash of the previous Block: %x\n", block.PrevHash)
				// fmt.Printf("All the transactions: %s\n", block.Data)

				// pow := blockchain.NewProof(block)
				// fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			}
		},
	}

	return printCmd
}
