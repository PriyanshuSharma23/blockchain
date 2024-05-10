package blockchain

import (
	"time"
)

type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
	Nonce     int64
}

// func setHash(b *Block) {
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	headers := bytes.Join([][]byte{timestamp, b.PrevHash, b.Data}, []byte{})
// 	hash := sha256.Sum256(headers)
// 	b.Hash = hash[:]
// }

func NewBlock(prevBlockHash []byte, data []byte) *Block {
	b := &Block{
		Timestamp: time.Now().Unix(),
		PrevHash:  prevBlockHash,
		Data:      data,
	}

	pow := NewProof(b)
	nonce, hash := pow.Run()

	b.Nonce = nonce
	b.Hash = hash

	return b
}

func NewGenesisBlock() *Block {
	return NewBlock([]byte{}, []byte("Genesis"))
}
