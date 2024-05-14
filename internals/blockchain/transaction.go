package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxnInput
	Outputs []TxnOutput
}

type TxnOutput struct {
	Value  int
	PubKey string
}

type TxnInput struct {
	ID  []byte
	Out int
	Sig string
}

func CoinbaseTxn(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxnInput{[]byte{}, -1, data}
	txout := TxnOutput{100, to}

	tx := Transaction{nil, []TxnInput{txin}, []TxnOutput{txout}}
	tx.SetID()

	return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	var encoder = gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxnInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxnOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

func NewTransaction(from string, to string, amount int, bc *Blockchain) (*Transaction, error) {
	acc, unspentOuts, err := bc.FindSpendableOutputs(from, amount)

	if err != nil {
		return nil, err
	}

	if acc < amount {
		return nil, ErrInsufficientFunds
	}

	txnInputs := make([]TxnInput, 0)
	txnOutputs := make([]TxnOutput, 0)

	for txID, outs := range unspentOuts {
		txIDBytes, err := hex.DecodeString(txID)

		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			txnInputs = append(txnInputs, TxnInput{
				ID:  txIDBytes,
				Out: out,
				Sig: from,
			})
		}
	}

	txnOutputs = append(txnOutputs, TxnOutput{
		Value:  amount,
		PubKey: to,
	})

	if amount < acc {
		txnOutputs = append(txnOutputs, TxnOutput{
			Value:  acc - amount,
			PubKey: from, // unused token
		})
	}

	tx := &Transaction{
		Inputs:  txnInputs,
		Outputs: txnOutputs,
	}
	tx.SetID()

	return tx, nil
}
