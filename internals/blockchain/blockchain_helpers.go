package blockchain

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/dgraph-io/badger/v4"
)

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func initDB() (*badger.DB, error) {
	dbOpts := badger.DefaultOptions(dbPath)
	dbOpts.Dir = dbPath
	dbOpts.ValueDir = dbPath

	return badger.Open(dbOpts)
}

func initializeBlockchain(debug bool, logger *log.Logger) (*Blockchain, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}

	if debug {
		logger.Println("Successfully connected to the database")
	}

	bc := &Blockchain{
		DB:     db,
		Debug:  debug,
		Logger: logger,
	}

	return bc, nil
}

// LinearSearch searches for a target element in a slice and returns its index if found, or -1 if not found.
func linearSearch(slice []int, target int) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

func (bc *Blockchain) getPrevTxns(tx *Transaction) map[string]Transaction {
	var m = map[string]Transaction{}

	for _, in := range tx.Inputs {
		tx, err := bc.FindTransaction(in.ID)
		txKey := hex.EncodeToString(in.ID)

		if err != nil {
			panic("transaction not found " + txKey)
		}

		m[txKey] = *tx
	}

	return m
}
