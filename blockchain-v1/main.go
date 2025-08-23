package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Block represents a single block in the blockchain.
type Block struct {
	Index     int
	Timestamp time.Time
	Data      string
	PrevHash  string
	Hash      string
}

func (b *Block) calculateHash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(
		fmt.Sprintf("%d%s%s%s", b.Index, b.Timestamp, b.Data, b.PrevHash),
	)))
}

// createGenesisBlock creates the first block in the blockchain.
func createGenesisBlock() *Block {
	b := &Block{
		Index:     0,
		Timestamp: time.Now(),
		Data:      "Genesis Block",
		PrevHash:  "0",
	}
	b.Hash = b.calculateHash()
	return b
}

func generateNewBlock(prevBlock *Block, data string) *Block {
	b := &Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
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
	blockchain = append(blockchain, generateNewBlock(blockchain[0], "First transaction data"))
	blockchain = append(blockchain, generateNewBlock(blockchain[1], "Second transaction data"))
	blockchain = append(blockchain, generateNewBlock(blockchain[2], "Third transaction data"))
	blockchain = append(blockchain, generateNewBlock(blockchain[3], "Fourth transaction data"))

	// Try to add an invalid block (tampered data).
	invalidBlock := &Block{
		Index:     7,
		Timestamp: time.Now(),
		Data:      "Tampered data",
		PrevHash:  blockchain[4].Hash,
	}
	invalidBlock.Hash = invalidBlock.calculateHash()
	blockchain = append(blockchain, invalidBlock)

	// Print the blockchain.
	for i, block := range blockchain {
		fmt.Printf("Block %d:\n", block.Index)
		fmt.Printf("  Timestamp: %s\n", block.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Data: %s\n", block.Data)
		fmt.Printf("  Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("  Hash: %s\n", block.Hash)
		fmt.Printf("  Valid: %v\n\n", isBlockValid(block, blockchain[max(0, i-1)]))
	}
}
