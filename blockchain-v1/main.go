package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Transaction represents a simple transaction structure.
type Transaction struct {
	From   string
	To     string
	Amount float64
}

// Block represents a single block in the blockchain.
type Block struct {
	Index        int
	Timestamp    time.Time
	Transactions []Transaction
	PrevHash     string
	Hash         string
}

func (b *Block) calculateHash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(
		fmt.Sprintf("%d%s%v%s", b.Index, b.Timestamp, b.Transactions, b.PrevHash),
	)))
}

// createGenesisBlock creates the first block in the blockchain.
func createGenesisBlock() *Block {
	b := &Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: nil,
		PrevHash:     "0",
	}
	b.Hash = b.calculateHash()
	return b
}

func generateNewBlock(prevBlock *Block, tx []Transaction) *Block {
	b := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: tx,
		PrevHash:     prevBlock.Hash,
	}
	b.Hash = b.calculateHash()
	return b
}

func isBlockValid(newBlock, prevBlock *Block) bool {
	// Handle the case where prevBlock is nil (genesis block).
	if newBlock.PrevHash == "0" && newBlock.Index == 0 {
		return true
	}

	if newBlock.Index != prevBlock.Index+1 {
		return false
	}
	if newBlock.PrevHash != prevBlock.Hash {
		return false
	}
	if newBlock.Hash != newBlock.calculateHash() {
		return false
	}
	return true
}

func main() {
	// Create the blockchain.
	var blockchain []*Block

	// Add genesis block.
	blockchain = append(blockchain, createGenesisBlock())

	// Mint new blocks.
	blockchain = append(blockchain, generateNewBlock(blockchain[0], []Transaction{
		{From: "Alice", To: "Bob", Amount: 10},
		{From: "Bob", To: "Charlie", Amount: 5},
	}))
	blockchain = append(blockchain, generateNewBlock(blockchain[1], []Transaction{
		{From: "Charlie", To: "Dave", Amount: 2},
		{From: "Dave", To: "Eve", Amount: 1},
		{From: "Eve", To: "Frank", Amount: 0.5},
	}))
	blockchain = append(blockchain, generateNewBlock(blockchain[2], []Transaction{
		{From: "Frank", To: "Grace", Amount: 0.25},
		{From: "Grace", To: "Heidi", Amount: 0.1},
		{From: "Heidi", To: "Ivan", Amount: 0.05},
		{From: "Ivan", To: "Judy", Amount: 0.01},
	}))
	blockchain = append(blockchain, generateNewBlock(blockchain[3], []Transaction{
		{From: "Judy", To: "Mallory", Amount: 0.005},
		{From: "Mallory", To: "Niaj", Amount: 0.0025},
		{From: "Niaj", To: "Olivia", Amount: 0.001},
		{From: "Olivia", To: "Peggy", Amount: 0.0005},
		{From: "Peggy", To: "Quentin", Amount: 0.0001},
	}))

	// Try to add an invalid block (tampered data).
	invalidBlock := &Block{
		Index:        7,
		Timestamp:    time.Now(),
		Transactions: []Transaction{{From: "Eve", To: "Frank", Amount: 1000}}, // Tampered transaction
		PrevHash:     blockchain[4].Hash,
	}
	invalidBlock.Hash = invalidBlock.calculateHash()
	blockchain = append(blockchain, invalidBlock)

	// Print the blockchain.
	for i, block := range blockchain {
		fmt.Printf("Block %d:\n", block.Index)
		fmt.Printf("  Timestamp: %s\n", block.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Data: %v\n", block.Transactions)
		fmt.Printf("  Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("  Hash: %s\n", block.Hash)
		fmt.Printf("  Valid: %v\n\n", isBlockValid(block, blockchain[max(0, i-1)]))
	}
}
