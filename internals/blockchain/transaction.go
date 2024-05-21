package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/PriyanshuSharma23/custom_blockchain/internals/wallet"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxnInput
	Outputs []TxnOutput
}

func (txn *Transaction) Hash() []byte {
	txnCpy := *txn
	txnCpy.ID = []byte{}
	hash := sha256.Sum256(txnCpy.Serialize())
	return hash[:]
}

func (txn *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	if err := encoder.Encode(txn); err != nil {
		panic(err)
	}

	return encoded.Bytes()
}

func (txn *Transaction) TrimmedCopy() *Transaction {
	txnCpy := &Transaction{}

	inputs := []TxnInput{}
	outputs := []TxnOutput{}

	for _, in := range txn.Inputs {
		inputs = append(inputs, TxnInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range txn.Outputs {
		outputs = append(outputs, TxnOutput{out.Value, out.PubKeyHash})
	}

	txnCpy.Inputs = inputs
	txnCpy.Outputs = outputs

	return txnCpy
}

func (txn *Transaction) Sign(privK ecdsa.PrivateKey, prevTxns map[string]Transaction) {
	if txn.IsCoinbase() {
		return
	}

	if err := txn.validatePrevTxns(prevTxns); err != nil {
		panic(err)
	}

	txnCpy := txn.TrimmedCopy()

	for i, in := range txnCpy.Inputs {
		prevTxn := prevTxns[hex.EncodeToString(in.ID)] // fetch the prev transaction, for its PubKeyHash and ID

		txnCpy.Inputs[i].PubKey = prevTxn.Outputs[in.Out].PubKeyHash
		txnCpy.ID = txnCpy.Hash()
		txnCpy.Inputs[i].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privK, txnCpy.ID)
		if err != nil {
			panic(err)
		}

		sig := bytes.Join([][]byte{
			r.Bytes(),
			s.Bytes(),
		}, []byte{})

		txn.Inputs[i].Sig = sig
	}
}

func (tx *Transaction) Verify(prevTxns map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	if err := tx.validatePrevTxns(prevTxns); err != nil {
		panic(err)
	}

	txCpy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for i, in := range tx.Inputs {
		prevTxn := prevTxns[hex.EncodeToString(in.ID)]

		txCpy.Inputs[i].PubKey = prevTxn.Outputs[in.Out].PubKeyHash
		txCpy.ID = txCpy.Hash()

		txCpy.Inputs[i].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		r.SetBytes(in.Sig[:len(in.Sig)/2])
		s.SetBytes(in.Sig[len(in.Sig)/2:])

		x := big.Int{}
		y := big.Int{}
		x.SetBytes(in.PubKey[:len(in.PubKey)/2])
		y.SetBytes(in.PubKey[len(in.PubKey)/2:])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, txCpy.ID, &r, &s) {
			return false
		}
	}

	return true
}

func CoinbaseTxn(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxnInput{[]byte{}, -1, []byte{}, []byte(data)}
	txout := *NewTxnOutput(100, to)

	tx := Transaction{nil, []TxnInput{txin}, []TxnOutput{txout}}
	tx.ID = tx.Hash()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func NewTransaction(from string, to string, amount int, bc *Blockchain) (*Transaction, error) {
	txnInputs := make([]TxnInput, 0)
	txnOutputs := make([]TxnOutput, 0)

	wallets := wallet.NewWallets()
	w, ok := wallets.GetWallet(from)

	if !ok {
		return nil, fmt.Errorf("wallet not found: " + from)
	}

	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, unspentOuts, err := bc.FindSpendableOutputs(pubKeyHash, amount)

	if err != nil {
		return nil, err
	}

	if acc < amount {
		return nil, ErrInsufficientFunds
	}

	for txID, outs := range unspentOuts {
		txIDBytes, err := hex.DecodeString(txID)

		if err != nil {
			return nil, err
		}

		for _, out := range outs {
			txnInputs = append(txnInputs, TxnInput{
				ID:     txIDBytes,
				Out:    out,
				PubKey: w.PublicKey,
				Sig:    nil,
			})
		}
	}

	txnOutputs = append(txnOutputs, *NewTxnOutput(amount, to))

	if amount < acc {
		txnOutputs = append(txnOutputs, *NewTxnOutput(acc-amount, from))
	}

	tx := &Transaction{
		Inputs:  txnInputs,
		Outputs: txnOutputs,
	}
	tx.ID = tx.Hash()
	bc.SignTransaction(tx, *w.PrivateKey)

	return tx, nil
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Out))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Sig))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
