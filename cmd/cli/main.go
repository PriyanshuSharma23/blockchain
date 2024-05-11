package main

import (
	"github.com/PriyanshuSharma23/custom_blockchain/cmd/cli/cmd"
	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
)

func main() {
	bc, err := blockchain.NewBlockchain()

	if err != nil {
		panic(err)
	}

	cli := cmd.NewCli(bc)
	cli.Execute()
}
