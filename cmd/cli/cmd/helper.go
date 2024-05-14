package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/blockchain"
)

func handleErr(err error, chain *blockchain.Blockchain) {
	if err != nil {
		if chain != nil {
			chain.DB.Close()
		}

		log.Println(err)
		fmt.Println(string(debug.Stack()))
		// runtime.Goexit()
		os.Exit(1)
	}
}
