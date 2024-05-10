package blockchain

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
)

type Blockchain struct {
	LastHash []byte
	db       *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	db          *badger.DB
}

const (
	dbPath = "./tmp/blocks"
)

var (
	lhKey = []byte("lh")
)

func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte

	err := bc.db.View(func(txn *badger.Txn) error {
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
		return err
	}

	newBlock := NewBlock(lastHash, []byte(data))

	err = bc.db.Update(func(txn *badger.Txn) error {
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

func NewBlockchain() (*Blockchain, error) {
	var lastHash []byte

	dbOpts := badger.DefaultOptions(dbPath)
	dbOpts.Dir = dbPath
	dbOpts.ValueDir = dbPath

	db, err := badger.Open(dbOpts)

	if err != nil {
		return nil, err
	}

	err = db.Update(func(txn *badger.Txn) error {
		val, err := txn.Get(lhKey)

		if errors.Is(err, badger.ErrKeyNotFound) {
			b := NewGenesisBlock()

			err = txn.Set(b.Hash, b.Serialize())

			if err != nil {
				return err
			}

			err = txn.Set(lhKey, b.Hash)
			lastHash = b.Hash

			return err
		}

		err = val.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return &Blockchain{
		db:       db,
		LastHash: lastHash,
	}, nil
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: bc.LastHash,
		db:          bc.db,
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
