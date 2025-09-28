package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func (w *Wallet) GetAddress() string {
	// In Bitcoin, address is derived from public key hash.
	pubKeyBytes := append(w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes()...)
	hash := sha256.Sum256(pubKeyBytes)
	return fmt.Sprintf("%x", hash[:25]) // Simplified address, subset of the hash.
}

func (w *Wallet) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return ecdsa.SignASN1(rand.Reader, w.PrivateKey, hash[:])
}

// verifySignature verifies the signature of the data using the provided public key.
// Other peers on the network can use this function to verify the authenticity of the data.
func verifySignature(pubKey *ecdsa.PublicKey, data, signature []byte) bool {
	hash := sha256.Sum256(data)
	return ecdsa.VerifyASN1(pubKey, hash[:], signature)
}

func main() {
	alice, err := NewWallet()
	if err != nil {
		log.Fatal(err)
	}
	bob, err := NewWallet()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Alice's address:", alice.GetAddress())
	fmt.Println("Bob's address:", bob.GetAddress())

	// Alice signs a transaction to send 10 units to Bob.
	txData := []byte("Alice pays Bob 10 units")
	signature, err := alice.Sign(txData)
	if err != nil {
		log.Fatal(err)
	}

	// Bob (or any other peer) verifies the transaction signature.
	txValid := verifySignature(alice.PublicKey, txData, signature)
	fmt.Println("The transaction is valid from Alice? ", txValid)

	// Bob can't forge Alice's signature.
	forgedSig, err := bob.Sign(txData)
	if err != nil {
		log.Fatal(err)
	}
	forgedValid := verifySignature(alice.PublicKey, txData, forgedSig)
	fmt.Println("The forged transaction by Bob is valid? ", forgedValid)

	fmt.Println("Your private key proves ownership. Your address is where people send money.")
}
