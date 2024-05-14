package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
)

type Blockchain struct {
	LastHash []byte
	DB       *badger.DB
	Debug    bool
	Logger   *log.Logger
}

type BlockchainIterator struct {
	CurrentHash []byte
	db          *badger.DB
}

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "GENESIS"
)

var (
	lhKey                       = []byte("lh")
	ErrBlockchainAlreadyExists  = fmt.Errorf("blockchain: already exists, db not empty")
	ErrBlockchainNotInitialized = fmt.Errorf("blockchain: db does not exists, please initialize")
	ErrInsufficientFunds        = fmt.Errorf("blockchain: sender has insufficient funds")
)

// Used to create the blockchain from scratch, with the genesis block and the coinbase transaction
func CreateBlockchain(address string, debug bool, logger *log.Logger) (*Blockchain, error) {
	if dbExists() {
		return nil, ErrBlockchainAlreadyExists
	}

	bc, err := initializeBlockchain(debug, logger)
	if err != nil {
		return nil, err
	}

	err = bc.DB.Update(func(txn *badger.Txn) error {
		coinbaseTxn := CoinbaseTxn(address, genesisData)
		block := NewGenesisBlock(coinbaseTxn)

		if bc.Debug {
			bc.Logger.Println("Genesis Created")
		}

		err := txn.Set(block.Hash, block.Serialize())
		if err != nil {
			return err
		}

		err = txn.Set(lhKey, block.Hash)
		if err != nil {
			return err
		}

		bc.LastHash = block.Hash
		return nil
	})

	if err != nil {
		return nil, err
	}

	return bc, nil
}

// Used to intialize an existing blockchain instance in the database
func ContinueBlockchain(debug bool, logger *log.Logger) (*Blockchain, error) {
	if !dbExists() {
		return nil, ErrBlockchainNotInitialized
	}

	bc, err := initializeBlockchain(debug, logger)
	if err != nil {
		return nil, err
	}

	lastHash, err := bc.fetchLastHash()

	if err != nil {
		return nil, err
	}

	bc.LastHash = lastHash
	return bc, nil
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: bc.LastHash,
		db:          bc.DB,
	}
}

func (it *BlockchainIterator) Next() (*Block, error) {
	var block *Block

	if len(it.CurrentHash) == 0 {
		return nil, nil
	}

	err := it.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(it.CurrentHash)

		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	it.CurrentHash = block.PrevHash
	return block, nil
}

func (bc *Blockchain) fetchLastHash() ([]byte, error) {
	var lastHash []byte

	err := bc.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(lhKey)

		if err != nil {
			return err
		}

		if err = item.Value(func(val []byte) error {
			lastHash = val

			return nil
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return lastHash, nil
}

func (bc *Blockchain) AddBlock(txns []Transaction) error {
	lastHash, err := bc.fetchLastHash()

	if err != nil {
		return err
	}

	newBlock := NewBlock(lastHash, txns)

	err = bc.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())

		if err != nil {
			return err
		}

		err = txn.Set(lhKey, newBlock.Hash)
		if err != nil {
			return err
		}

		bc.LastHash = newBlock.Hash
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (bc *Blockchain) FindUnspentTransactions(address string) ([]Transaction, error) {
	var unspentTxs []Transaction
	var spentTxnOutputs = make(map[string][]int)

	iter := bc.Iterator()

	for {
		blk, err := iter.Next()

		if err != nil {
			return nil, err
		}

		if nil == blk {
			break
		}

		// Iterate over the transactions
		for _, txn := range blk.Transactions {
			txID := hex.EncodeToString(txn.ID)

			// loop over the ouputs in txn
			for outIdx, out := range txn.Outputs {
				if spentTxnOutputs[txID] != nil && linearSearch(spentTxnOutputs[txID], outIdx) != -1 {
					continue
				}

				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, txn)
				}
			}

			if txn.IsCoinbase() {
				continue
			}

			// loop over all the inputs in txn
			for _, in := range txn.Inputs {
				if in.CanUnlock(address) {
					inTxnID := hex.EncodeToString(in.ID)
					spentTxnOutputs[inTxnID] = append(spentTxnOutputs[inTxnID], in.Out)
				}
			}
		}

	}

	return unspentTxs, nil
}

func (bc *Blockchain) FindUTXO(address string) ([]TxnOutput, error) {
	var utxos []TxnOutput

	txns, err := bc.FindUnspentTransactions(address)

	if err != nil {
		return nil, err
	}

	for _, tx := range txns {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				utxos = append(utxos, out)
			}
		}
	}

	return utxos, nil
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (accumulated int, unspentOuts map[string][]int, err error) {
	unspentOuts = make(map[string][]int)

	txns, err := bc.FindUnspentTransactions(address)
	if err != nil {
		return 0, nil, err
	}

	for _, tx := range txns {
		txID := hex.EncodeToString(tx.ID)

		overflown := false

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					overflown = true
					break
				}
			}
		}

		if overflown {
			break
		}
	}

	return accumulated, unspentOuts, nil
}
