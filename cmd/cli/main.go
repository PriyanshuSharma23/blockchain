package main

import (
	"flag"
	"fmt"
	// "strconv"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
)

func main() {
	data := flag.String("add", "", "add block, pass data")
	printBlocks := flag.Bool("print", false, "print all blocks")

	flag.Parse()

	if *data == "" && *printBlocks == false {
		fmt.Println("Select an operation")
		return
	}

	if *data != "" && *printBlocks == true {
		fmt.Println("Select one operation")
		return
	}

	bc, err := blockchain.NewBlockchain()

	if err != nil {
		panic(err)
	}

	if *data != "" {
		err = bc.AddBlock(*data)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Added block: %x\n", bc.LastHash)
	} else {

		it := bc.Iterator()

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

	}

}
