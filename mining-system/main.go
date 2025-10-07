package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
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
	Nonce        int
}

func (b *Block) calculateHash() string {
	record := fmt.Sprintf("%d%s%v%s%d", b.Index, b.Timestamp, b.Transactions, b.PrevHash, b.Nonce)
	h := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", h)
}

func (b *Block) mineBlock(difficulty int) {
	target := strings.Repeat("0", difficulty)
	start := time.Now()

	fmt.Printf("Mining block %d...\n", b.Index)

	for {
		b.Hash = b.calculateHash()
		if strings.HasPrefix(b.Hash, target) {
			fmt.Printf("Block mined: %s\n", b.Hash)
			fmt.Printf("Nonce: %d\n", b.Nonce)
			fmt.Printf("Time: %v\n\n", time.Since(start))
			break
		}
		b.Nonce += 1

		// Progress indicator, every 100,000 hashes.
		if b.Nonce%100_000 == 0 {
			fmt.Printf("  Tried %d hashes...\n", b.Nonce)
		}
	}
}

type BlockChain struct {
	Chain      []*Block
	Difficulty int
}

func (bc *BlockChain) createGenesisBlock() *Block {
	b := &Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: nil,
		PrevHash:     "0",
		Nonce:        0,
	}
	b.mineBlock(bc.Difficulty)
	return b
}

func (bc *BlockChain) addBlock(from, to string, amount float64) {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: []Transaction{{From: from, To: to, Amount: amount}},
		PrevHash:     prevBlock.Hash,
		Nonce:        0,
	}
	newBlock.mineBlock(bc.Difficulty)
	bc.Chain = append(bc.Chain, newBlock)
}

func (bc *BlockChain) isChainValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		prevBlock := bc.Chain[i-1]

		if currentBlock.Hash != currentBlock.calculateHash() {
			return false
		}
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}

		// Check if the block was mined correctly.
		target := strings.Repeat("0", bc.Difficulty)
		if !strings.HasPrefix(currentBlock.Hash, target) {
			fmt.Printf("Block %d not properly mined\n", i)
			return false
		}
	}
	return true
}

func runBlockChain() {
	bc := &BlockChain{Difficulty: 4}

	// Mine genesis block
	fmt.Println("Mining genesis block...")
	bc.Chain = append(bc.Chain, bc.createGenesisBlock())

	// Add more blocks
	bc.addBlock("Alice", "Bob", 10)
	bc.addBlock("Bob", "Charlie", 5)
	bc.addBlock("Charlie", "Dave", 2.5)

	// Display the blockchain
	fmt.Println("\n=== BLOCKCHAIN ===")
	for _, block := range bc.Chain {
		fmt.Printf("Block %d:\n", block.Index)
		fmt.Printf("  Timestamp: %s\n", block.Timestamp.Format("15:04:05"))
		fmt.Printf("  Data: %v\n", block.Transactions)
		fmt.Printf("  Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		fmt.Printf("  Hash: %s\n\n", block.Hash)
	}

	// Validate blockchain
	fmt.Printf("Blockchain valid: %v\n", bc.isChainValid())
}

func simulateFork() {
	bc := &BlockChain{Difficulty: 4}

	// Mine genesis block
	fmt.Println("Mining genesis block...")
	bc.Chain = append(bc.Chain, bc.createGenesisBlock())

	// Add more blocks
	bc.addBlock("Alice", "Bob", 10)
	bc.addBlock("Bob", "Charlie", 5)
	bc.addBlock("Charlie", "Dave", 2.5)
	bc.addBlock("Charlie", "Dave", 7.5)
	bc.addBlock("Charlie", "Dave", 9.5)
	bc.addBlock("T", "Dave", 2.5)

	chainLastIndex := len(bc.Chain) - 1
	fmt.Printf("Current chain length: %d blocks\n", len(bc.Chain))

	// Now two nodes create different versions of the next block (fork).
	// Both blocks point to the same previous block (index 5).
	nextIndex := chainLastIndex + 1
	bk6A := &Block{
		Index:        nextIndex,
		Timestamp:    time.Now(),
		Transactions: []Transaction{{From: "Eve", To: "Frank", Amount: 1}},
		PrevHash:     bc.Chain[chainLastIndex].Hash, // Same previous hash to create a fork.
		Nonce:        0,
	}
	bk6A.mineBlock(bc.Difficulty)

	bk6B := &Block{
		Index:        nextIndex,
		Timestamp:    time.Now(),
		Transactions: []Transaction{{From: "Eve", To: "Grace", Amount: 1.5}},
		PrevHash:     bc.Chain[chainLastIndex].Hash, // Same previous hash to create a fork.
		Nonce:        0,
	}
	bk6B.mineBlock(bc.Difficulty)

	fmt.Println("Fork created! Two valid block 6s exist.")
	baseChain := []*Block{}
	baseChain = append(baseChain, bc.Chain...) // Copy current chain

	// Network splits temporarily
	chainA := append(baseChain, bk6A)
	chainB := append(baseChain, bk6B)
	_ = chainB

	// Then chain A mines a new block
	bk7A := &Block{
		Index:        nextIndex + 1,
		Timestamp:    time.Now(),
		Transactions: []Transaction{{From: "Frank", To: "Heidi", Amount: 0.5}},
		PrevHash:     bk6A.Hash,
		Nonce:        0,
	}
	bk7A.mineBlock(bc.Difficulty)
	chainA = append(chainA, bk7A)

	fmt.Println("\n=== FORKED BLOCKCHAIN A ===")
	for _, block := range chainA {
		fmt.Printf("Block %d:\n", block.Index)
		fmt.Printf("  Timestamp: %s\n", block.Timestamp.Format("15:04:05"))
		fmt.Printf("  Data: %v\n", block.Transactions)
		fmt.Printf("  Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		fmt.Printf("  Hash: %s\n\n", block.Hash)
	}
	bc.Chain = chainA
	// Chain A is now longer (7 blocks vs 6)
	// Network consensus: Everyone switches to Chain A
	fmt.Println("Chain A is longer. Network switches to Chain A.")
	fmt.Println("Chain B block 6 becomes orphaned.")
	fmt.Println("is Chain A valid?", bc.isChainValid())
}

func main() {
	//runBlockChain()
	simulateFork()
}
