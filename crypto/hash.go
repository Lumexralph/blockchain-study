package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func main() {
	data := "Hello, Blockchain!"
	hash := sha256.Sum256([]byte(data))
	fmt.Printf("Data: %s\n", data)
	fmt.Printf("SHA-256 Hash: %x\n", hash)
	fmt.Printf("Hash Length: %d bytes\n", len(hash))

	data2 := "Hello, Blockchain!!"
	hash2 := sha256.Sum256([]byte(data2))
	fmt.Printf("Data: %s\n", data2)
	fmt.Printf("SHA-256 Hash: %x\n", hash2)
	fmt.Printf("Hash Length: %d bytes\n", len(hash2))

	privateKey, publicKey, err := generateKeys()
	if err != nil {
		fmt.Println("Error generating keys:", err)
		return
	}
	fmt.Printf("Private Key: %s\n", privateKey)
	fmt.Printf("Public Key: %s\n", publicKey)
}

func generateKeys() (privateKey, publicKey string, _ error) {
	pKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}
	pubKey := &pKey.PublicKey

	// Sign a message to create a digital signature.
	msg := "I need to create a L1 blockchain"
	hash := sha256.Sum256([]byte(msg))
	signature, err := ecdsa.SignASN1(rand.Reader, pKey, hash[:])
	if err != nil {
		return "", "", err
	}
	fmt.Printf("Message: %s\n", msg)
	fmt.Printf("Signature: %x\n", signature)

	// Verify the signature.
	verified := ecdsa.VerifyASN1(pubKey, hash[:], signature)
	if verified {
		fmt.Println("Signature verified successfully!")
	} else {
		fmt.Println("Signature verification failed.")
	}
	return fmt.Sprintf("%x", pKey.D), fmt.Sprintf("%x", pubKey.X), nil
}
