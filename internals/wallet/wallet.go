package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/mr-tron/base58/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	checkSumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func newPair() (*ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()

	pk, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	pub := bytes.Join([][]byte{
		pk.X.Bytes(),
		pk.Y.Bytes(),
	}, []byte{})

	return pk, pub, nil
}

func NewWallet() (*Wallet, error) {
	priv, pub, err := newPair()
	if err != nil {
		return nil, err
	}

	return &Wallet{priv, pub}, nil
}

func (w *Wallet) Address() string {
	pbkHash := publicKeyHash(w.PublicKey)

	c := checksum(pbkHash)

	data := bytes.Join([][]byte{
		{version},
		pbkHash,
		c[:],
	}, []byte{})

	address := base58.Encode(data)
	return address
}

func publicKeyHash(publicKey []byte) []byte {
	pbHash := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pbHash[:])
	if err != nil {
		panic(err)
	}

	pbHashRipeMD := hasher.Sum(nil)
	return pbHashRipeMD
}

func checksum(pbHash []byte) [checkSumLength]byte {
	sha1 := sha256.Sum256(pbHash)
	sha2 := sha256.Sum256(sha1[:])

	return [checkSumLength]byte((sha2[:checkSumLength]))
}
