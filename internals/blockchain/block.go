package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
	Nonce     int64
}

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
