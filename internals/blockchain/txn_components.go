package blockchain

import (
	"bytes"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/wallet"
)

type TxnOutput struct {
	Value      int
	PubKeyHash []byte
}

type TxnInput struct {
	ID     []byte // refers to the previous transaction being used
	Out    int    // refers to the output index of the previous transaction
	Sig    []byte
	PubKey []byte
}

func (in *TxnInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxnOutput) Lock(addr string) {
	pubKeyHash, err := wallet.PublicKeyHashFromAddr(addr)
	if err != nil {
		panic(err)
	}
	out.PubKeyHash = pubKeyHash
}

func (out *TxnOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTxnOutput(value int, addr string) *TxnOutput {
	txo := &TxnOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	txo.Lock(addr)

	return txo
}
