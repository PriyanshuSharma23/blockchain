// Save wallets in disk
package wallet

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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

func (ws *Wallets) encode() ([]byte, error) {
	wsSlice := []any{}

	for _, w := range ws.m {
		mapStringAny := map[string]any{
			"PrivateKey": map[string]any{
				"D": w.PrivateKey.D,
				"PublicKey": map[string]any{
					"X": w.PrivateKey.PublicKey.X,
					"Y": w.PrivateKey.PublicKey.Y,
				},
				"X": w.PrivateKey.X,
				"Y": w.PrivateKey.Y,
			},
			"PublicKey": w.PublicKey,
		}

		wsSlice = append(wsSlice, mapStringAny)
	}

	return json.Marshal(wsSlice)
}

func decodeWallets(r io.Reader) (Wallets, error) {
	var wSlice []Wallet

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&wSlice)

	var wm = make(walletsMap)
	for _, w := range wSlice {
		wm[w.Address()] = &w
	}

	if err != nil {
		log.Panic(err)
	}

	return Wallets{
		m: wm,
	}, nil
}
