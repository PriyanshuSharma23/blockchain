// Save wallets in disk
package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"io"
	"math/big"
	"os"
)

const (
	filePath = "./tmp/wallets/WALLETS.json"
)

type walletsMap map[string]*Wallet

type Wallets struct {
	m walletsMap
}

func NewWallets() *Wallets {
	ws, err := loadFromFile()

	if err != nil {
		switch {
		case os.IsNotExist(err):
			return &Wallets{
				m: make(map[string]*Wallet),
			}
		default:
			panic(err)
		}
	}

	return ws
}

func (ws *Wallets) GetAllWallets() []*Wallet {
	var wSlice = make([]*Wallet, 0, len(ws.m))
	for _, w := range ws.m {
		wSlice = append(wSlice, w)
	}

	return wSlice
}

func (ws *Wallets) AddWallet(w *Wallet) {
	ws.m[w.Address()] = w
	ws.saveToFile()
}

func (ws *Wallets) GetWallet(addr string) (w *Wallet, ok bool) {
	w, ok = ws.m[addr]
	return
}

func (ws *Wallets) saveToFile() error {
	data, err := ws.encode()
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func loadFromFile() (*Wallets, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	ws, err := decodeWallets(bytes.NewReader(data))
	return &ws, err
}

// PrivateKeyJSON struct to represent the JSON structure of the private key
type privateKeyJSON struct {
	D         string        `json:"D"`
	PublicKey publicKeyJSON `json:"PublicKey"`
}

// PublicKeyJSON struct to represent the JSON structure of the public key
type publicKeyJSON struct {
	X string `json:"X"`
	Y string `json:"Y"`
}

// WalletJSON struct to represent the JSON structure of a wallet
type walletJSON struct {
	PrivateKey privateKeyJSON `json:"PrivateKey"`
	PublicKey  []byte         `json:"PublicKey"`
}

func (ws *Wallets) encode() ([]byte, error) {
	wsSlice := []any{}

	for _, w := range ws.m {
		walletJSON := walletJSON{
			PrivateKey: privateKeyJSON{
				D: w.PrivateKey.D.String(),
				PublicKey: publicKeyJSON{
					X: w.PrivateKey.PublicKey.X.String(),
					Y: w.PrivateKey.PublicKey.Y.String(),
				},
			},
			PublicKey: w.PublicKey,
		}

		wsSlice = append(wsSlice, walletJSON)
	}

	return json.Marshal(wsSlice)
}

func decodeWallets(r io.Reader) (Wallets, error) {
	var wsSlice []walletJSON
	var wallets Wallets
	wallets.m = make(map[string]*Wallet)

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&wsSlice); err != nil {
		return wallets, err
	}

	for _, w := range wsSlice {
		D := new(big.Int)
		D.SetString(w.PrivateKey.D, 10)

		X := new(big.Int)
		X.SetString(w.PrivateKey.PublicKey.X, 10)

		Y := new(big.Int)
		Y.SetString(w.PrivateKey.PublicKey.Y, 10)

		privateKey := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: elliptic.P256(), // Assuming the curve is P256, modify as necessary
				X:     X,
				Y:     Y,
			},
			D: D,
		}

		wallet := &Wallet{
			PrivateKey: privateKey,
			PublicKey:  w.PublicKey,
		}

		wallets.m[wallet.Address()] = wallet
	}

	return wallets, nil
}
