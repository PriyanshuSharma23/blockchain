package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
}

func setHash(b *Block) {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{timestamp, b.PrevHash, b.Data}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(prevBlockHash []byte, data []byte) *Block {
	b := &Block{
		Timestamp: time.Now().Unix(),
		PrevHash:  prevBlockHash,
		Data:      data,
	}

	setHash(b)

	return b
}

func NewGenesisBlock() *Block {
	return NewBlock([]byte{}, []byte("Genesis"))
}
