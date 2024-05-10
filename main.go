package main

import (
	"fmt"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	bc.AddBlock("first transaction")
	bc.AddBlock("Second transaction")

	for _, block := range bc.Blocks {
		fmt.Printf("Hash of the block %x\n", block.Hash)
		fmt.Printf("Hash of the previous Block: %x\n", block.PrevHash)
		fmt.Printf("All the transactions: %s\n", block.Data)
	}
}

