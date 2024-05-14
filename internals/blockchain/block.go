package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	PrevHash     []byte
	Transactions []Transaction
	Nonce        int64
}

func NewBlock(prevBlockHash []byte, txns []Transaction) *Block {
	b := &Block{
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevBlockHash,
		Transactions: txns,
	}

	pow := NewProof(b)
	nonce, hash := pow.Run()

	b.Nonce = nonce
	b.Hash = hash

	return b
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]byte{}, []Transaction{*coinbase})
}

func (b *Block) Serialize() []byte {
	var buff = bytes.NewBuffer([]byte{})
	err := gob.NewEncoder(buff).Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

func (b *Block) HashTransactions() []byte {
	data := [][]byte{}

	for _, tx := range b.Transactions {
		data = append(data, tx.ID)
	}

	hash := sha256.Sum256(bytes.Join(data, []byte{}))

	return hash[:]
}
