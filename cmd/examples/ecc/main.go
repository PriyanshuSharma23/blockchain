package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func main() {
	// Generate a new private key using the P-256 elliptic curve
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}

	// Get the corresponding public key from the private key
	publicKey := privateKey.PublicKey

	// Serialize the private key to a hex string
	privateKeyBytes := privateKey.D.Bytes()
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// Serialize the public key to a hex string
	publicKeyXBytes := publicKey.X.Bytes()
	publicKeyXHex := hex.EncodeToString(publicKeyXBytes)
	publicKeyYBytes := publicKey.Y.Bytes()
	publicKeyYHex := hex.EncodeToString(publicKeyYBytes)

	fmt.Println("Private Key:", privateKeyHex)
	fmt.Println("Public Key X:", publicKeyXHex)
	fmt.Println("Public Key Y:", publicKeyYHex)
}
