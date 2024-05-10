package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 20

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-Difficulty)

	return &ProofOfWork{
		Block:  b,
		Target: target,
	}
}

func (p *ProofOfWork) InitData(nonce int64) []byte {
	data := bytes.Join([][]byte{
		p.Block.PrevHash,
		p.Block.Data,
		toHex(int64(nonce)),
		toHex(int64(Difficulty)),
	}, []byte{})

	return data
}

func (p *ProofOfWork) Run() (int64, []byte) {
	var intHash big.Int
	var hash [32]byte

	var nonce int64 = 0

	for nonce < math.MaxInt64 {
		data := p.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(p.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println()
	return nonce, hash[:]
}

func (p *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := p.InitData(p.Block.Nonce)
	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(p.Target) == -1
}

func toHex(n int64) []byte {
	buff := bytes.NewBuffer([]byte{})
	err := binary.Write(buff, binary.BigEndian, n)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
