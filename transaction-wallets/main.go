/*
UTXO (Unspent Transaction Output) model implementation in Go.
This code defines a simple UTXO structure and provides functions to create, spend, and query UTXOs.

There are two ways to track money:
 1. Account-based model (like a bank account):
    Alice: 100 naira
    Bob: 50 naira
    Charlie: 25 naira
 2. UTXO-based model (like cash in your wallet, like physical cash):
    Unspent Transaction Output #1: 20 coins → Alice
    Unspent Transaction Output #2: 30 coins → Alice
    Unspent Transaction Output #3: 50 coins → Alice
    Unspent Transaction Output #4: 50 coins → Bob

In the UTXO model, each transaction consumes some UTXOs and creates new ones.
This allows for better privacy and parallel processing of transactions.
*/
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// UTXO represents an unspent transaction output.
type UTXO struct {
	TxID   string  // Unique identifier for the UTXO
	Index  int     // Index of the output in the transaction
	Owner  string  // Owner of the UTXO
	Amount float64 // Amount of currency in the UTXO
}

type TxInput struct {
	TxID      string // Which transaction output we are spending.
	Index     int    // Which output in that transaction.
	Signature string // Proof we own the UTXO.
}

type TxOutput struct {
	Amount    float64
	Recipient string
}

type Transaction struct {
	ID        string // Unique identifier for the transaction
	Inputs    []TxInput
	Outputs   []TxOutput
	Timestamp time.Time
}

func calculateTxID(tx *Transaction) (string, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

// UTXOSet represents a collection of UTXOs.
// In a real-world application, this would be stored in a database.
// Tracks all unspent transaction outputs.
type UTXOSet struct {
	UTXOs map[string]*UTXO // Map of UTXO ID to UTXO
}

func NewUTXOSet() *UTXOSet {
	return &UTXOSet{UTXOs: make(map[string]*UTXO)}
}

// AddUTXO adds a new UTXO to the set.
func (u *UTXOSet) AddUTXO(utxo *UTXO) {
	key := fmt.Sprintf("%s:%d", utxo.TxID, utxo.Index)
	u.UTXOs[key] = utxo
}

// SpendUTXO marks a UTXO as spent by removing it from the set.
func (u *UTXOSet) SpendUTXO(txID string, index int) {
	key := fmt.Sprintf("%s:%d", txID, index)
	delete(u.UTXOs, key)
}

// GetBalance calculates the total balance for a given owner.
// The owner is identified by their address derived from the public address.
func (u *UTXOSet) GetBalance(addr string) float64 {
	var balance float64
	for _, utxo := range u.UTXOs {
		if utxo.Owner == addr {
			balance += utxo.Amount
		}
	}
	return balance
}

// FindSpendableUTXOs finds UTXOs that can be used to cover the specified amount.
// Returns a list of UTXOs and the total amount they cover.
func (u *UTXOSet) FindSpendableUTXOs(addr string, amount float64) ([]*UTXO, float64) {
	var selectedUTXOs []*UTXO
	var total float64

	for _, utxo := range u.UTXOs {
		if utxo.Owner == addr {
			selectedUTXOs = append(selectedUTXOs, utxo)
			total += utxo.Amount
			if total >= amount {
				break
			}
		}
	}

	if total < amount {
		return nil, 0 // Not enough funds
	}
	return selectedUTXOs, total
}

type Mempool struct {
	PendingTransactions []*Transaction
	SpentUTXOs          map[string]bool // Track UTXOs being spent
}

func (mp *Mempool) AddTransaction(tx *Transaction) error {
	// Check for double spending in mempool
	for _, input := range tx.Inputs {
		key := fmt.Sprintf("%s:%d", input.TxID, input.Index)
		if mp.SpentUTXOs[key] {
			return fmt.Errorf("UTXO %s already being spent", key)
		}
	}

	// Mark UTXOs as being spent
	for _, input := range tx.Inputs {
		key := fmt.Sprintf("%s:%d", input.TxID, input.Index)
		mp.SpentUTXOs[key] = true
	}

	mp.PendingTransactions = append(mp.PendingTransactions, tx)
	return nil
}

// Block represents a block in the blockchain.
type Block struct {
	Index        int
	Timestamp    time.Time
	Transactions []*Transaction
	PreviousHash string
	Nonce        int
	Hash         string
}

// CalculateHash calculates the hash of the block.
func (b *Block) CalculateHash() string {
	data := fmt.Sprintf("%d%s%s%d",
		b.Index,
		b.Timestamp.String(),
		b.PreviousHash,
		b.Nonce,
	)

	// Also include all transaction IDs
	for _, tx := range b.Transactions {
		data += tx.ID
	}

	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// MineBlock performs the proof of work algorithm to mine the block.
func (b *Block) MineBlock(difficulty int) {
	// Create a target string with 'difficulty' leading zeros
	target := strings.Repeat("0", difficulty)

	for {
		b.Hash = b.CalculateHash()

		// Check if hash has the required number of leading zeros
		if b.Hash[:difficulty] == target {
			break
		}
		b.Nonce++
	}

	fmt.Printf("Block mined: %s\n", b.Hash)
}

type Miner struct {
	Address string
	Mempool *Mempool
}

func (m *Miner) SelectTransactions() []*Transaction {
	// Select transactions with highest fees
	// Simplified: just take first few transactions
	maxTxs := 5
	if len(m.Mempool.PendingTransactions) < maxTxs {
		return m.Mempool.PendingTransactions
	}
	return m.Mempool.PendingTransactions[:maxTxs]
}

func (m *Miner) CreateCoinbaseTransaction(blockReward, fees int) *Transaction {
	// Coinbase transaction: creates new coins for the miner
	return &Transaction{
		ID:     "coinbase",
		Inputs: []TxInput{}, // No inputs - creates money
		Outputs: []TxOutput{
			{
				Amount:    float64(blockReward + fees),
				Recipient: m.Address,
			},
		},
		Timestamp: time.Now(),
	}
}

func (m *Miner) MineBlock(prevBlock *Block, difficulty int) *Block {
	transactions := m.SelectTransactions()

	// Calculate total fees
	fees := 0
	for _, tx := range transactions {
		// Calculate fee for each transaction
		// fees += calculateTransactionFee(tx)
		fees += 1
		_ = tx
	}

	// Create coinbase transaction
	coinbase := m.CreateCoinbaseTransaction(50, fees) // 50 coin block reward
	allTransactions := append([]*Transaction{coinbase}, transactions...)

	// Create and mine block
	block := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: allTransactions,
		PreviousHash: prevBlock.Hash,
		Nonce:        0,
	}

	// Mine the block (proof of work)
	block.MineBlock(difficulty)

	return block
}

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

func NewWallet() (*Wallet, error) {
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	public := &private.PublicKey

	// Generate address from public key
	pubKeyBytes := append(public.X.Bytes(), public.Y.Bytes()...)
	hash := sha256.Sum256(pubKeyBytes)
	address := fmt.Sprintf("%x", hash[:20])

	return &Wallet{
		PrivateKey: private,
		PublicKey:  public,
		Address:    address,
	}, nil
}

func (w *Wallet) GetBalance(utxoSet *UTXOSet) float64 {
	return utxoSet.GetBalance(w.Address)
}

func (w *Wallet) CreateTransaction(to string, amount float64, utxoSet *UTXOSet) (*Transaction, error) {
	// Find UTXOs to spend
	utxosToSpend, _ := utxoSet.FindSpendableUTXOs(w.Address, amount)
	if utxosToSpend == nil {
		return nil, fmt.Errorf("insufficient funds")
	}

	// Calculate total input
	totalInput := 0.0
	var inputs []TxInput
	for _, utxo := range utxosToSpend {
		totalInput += utxo.Amount
		inputs = append(inputs, TxInput{
			TxID:  utxo.TxID,
			Index: utxo.Index,
			// Signature will be added after transaction is created
		})
	}

	// Create outputs
	var outputs []TxOutput
	outputs = append(outputs, TxOutput{
		Amount:    amount,
		Recipient: to,
	})

	// Add change output if necessary
	change := totalInput - amount
	if change > 0 {
		outputs = append(outputs, TxOutput{
			Amount:    change,
			Recipient: w.Address,
		})
	}

	// Create transaction
	tx := &Transaction{
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: time.Now(),
	}

	txID, err := calculateTxID(tx)
	if err != nil {
		return nil, err
	}
	tx.ID = txID

	// Sign transaction
	for i := range tx.Inputs {
		signature, err := w.signTransaction(tx)
		if err != nil {
			return nil, err
		}
		tx.Inputs[i].Signature = string(signature)
	}

	return tx, nil
}

func (w *Wallet) signTransaction(tx *Transaction) ([]byte, error) {
	txData, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return w.Sign(txData)
}

// Sign signs the provided data with the wallet's private key
func (w *Wallet) Sign(data []byte) ([]byte, error) {
	// Create a hash of the data
	hash := sha256.Sum256(data)

	// Sign the hash with the private key
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	// Combine r and s into a signature
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// VerifySignature verifies a signature against the provided data using the public key
func (w *Wallet) VerifySignature(data []byte, signature []byte) bool {
	// Create a hash of the data
	hash := sha256.Sum256(data)

	// Split signature into r and s components
	sigLen := len(signature) / 2
	r := new(big.Int).SetBytes(signature[:sigLen])
	s := new(big.Int).SetBytes(signature[sigLen:])

	// Verify the signature
	return ecdsa.Verify(w.PublicKey, hash[:], r, s)
}

type CryptoCurrency struct {
	Blockchain []*Block
	UTXOSet    *UTXOSet
	Mempool    *Mempool
	Miners     []*Miner
	Wallets    []*Wallet
}

func (cc *CryptoCurrency) Initialize() {
	cc.UTXOSet = NewUTXOSet()
	cc.Mempool = &Mempool{
		PendingTransactions: []*Transaction{},
		SpentUTXOs:          make(map[string]bool),
	}

	// Create genesis block with initial coin distribution
	cc.createGenesisBlock()
}

func (cc *CryptoCurrency) SendMoney(from *Wallet, to string, amount float64) error {
	// Create transaction
	tx, err := from.CreateTransaction(to, amount, cc.UTXOSet)
	if err != nil {
		return err
	}

	// Add to mempool
	return cc.Mempool.AddTransaction(tx)
}

func (cc *CryptoCurrency) createGenesisBlock() {
	// Create a genesis wallet that will receive the initial coin distribution
	genesisWallet, err := NewWallet()
	if err != nil {
		fmt.Println("Error creating genesis wallet:", err)
		return
	}

	// Add the genesis wallet to the list of wallets
	cc.Wallets = append(cc.Wallets, genesisWallet)

	// Create a coinbase transaction that creates the initial coin supply
	coinbase := &Transaction{
		ID:     "genesis_coinbase",
		Inputs: []TxInput{}, // No inputs in the coinbase transaction
		Outputs: []TxOutput{
			{
				Amount:    1000000, // Initial coin supply (e.g., 1 million coins)
				Recipient: genesisWallet.Address,
			},
		},
		Timestamp: time.Now(),
	}

	// Create the genesis block
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: []*Transaction{coinbase},
		PreviousHash: "0", // First block has no previous hash
		Nonce:        0,
		Hash:         "", // Will be calculated below
	}

	// Mine the genesis block with a low difficulty
	genesisBlock.MineBlock(2)

	// Add the block to the blockchain
	cc.Blockchain = append(cc.Blockchain, genesisBlock)

	// Update the UTXO set with the coinbase transaction output
	for i, output := range coinbase.Outputs {
		utxo := &UTXO{
			TxID:   coinbase.ID,
			Index:  i,
			Owner:  output.Recipient,
			Amount: output.Amount,
		}
		cc.UTXOSet.AddUTXO(utxo)
	}

	fmt.Printf("Genesis block created with %f coins sent to %s\n", coinbase.Outputs[0].Amount, genesisWallet.Address)
}

func (cc *CryptoCurrency) MineNextBlock() *Block {
	// Miner creates new block with transactions from mempool
	miner := cc.Miners[0] // Use first miner
	prevBlock := cc.Blockchain[len(cc.Blockchain)-1]

	newBlock := miner.MineBlock(prevBlock, 4) // difficulty 4

	// Update UTXO set with new block's transactions
	cc.processBlock(newBlock)

	// Add to blockchain
	cc.Blockchain = append(cc.Blockchain, newBlock)

	return newBlock
}

// processBlock updates the UTXO set with the transactions in the provided block
func (cc *CryptoCurrency) processBlock(block *Block) {
	// Process each transaction in the block
	for _, tx := range block.Transactions {
		// Remove spent UTXOs
		for _, input := range tx.Inputs {
			// Skip coinbase transactions (they don't have inputs)
			if tx.ID == "coinbase" || tx.ID == "genesis_coinbase" {
				continue
			}

			cc.UTXOSet.SpendUTXO(input.TxID, input.Index)
		}

		// Create new UTXOs from transaction outputs
		for i, output := range tx.Outputs {
			newUTXO := &UTXO{
				TxID:   tx.ID,
				Index:  i,
				Owner:  output.Recipient,
				Amount: output.Amount,
			}
			cc.UTXOSet.AddUTXO(newUTXO)
		}

		// Remove transaction from mempool after processing
		for i, pendingTx := range cc.Mempool.PendingTransactions {
			if pendingTx.ID == tx.ID {
				// Remove this transaction from pending transactions
				cc.Mempool.PendingTransactions = append(
					cc.Mempool.PendingTransactions[:i],
					cc.Mempool.PendingTransactions[i+1:]...,
				)
				break
			}
		}
	}
}

func main() {
	utxoSet := NewUTXOSet()

	// Genesis transaction: 500 coins to Alice.
	// Maybe funding their wallet from an exchange.
	// In a real blockchain, this would be part of the genesis block.
	genesisUTXO := &UTXO{
		TxID:   "genesis",
		Index:  0,
		Amount: 500,
		Owner:  "alice_address",
	}
	utxoSet.AddUTXO(genesisUTXO)

	fmt.Printf("Alice's balance: %f\n", utxoSet.GetBalance("alice_address"))

	// Alice wants to send 30 coins to Bob.
	utxosToSpend, amount := utxoSet.FindSpendableUTXOs("alice_address", 30)
	if len(utxosToSpend) == 0 || amount == 0 {
		fmt.Println("Alice does not have enough funds to send 30 coins to Bob.")
		return
	}

	tx := &Transaction{
		Inputs: []TxInput{
			{TxID: "genesis", Index: 0, Signature: "alice_signature"},
		},
		Outputs: []TxOutput{
			{Amount: 30, Recipient: "bob_address"},    // 30 to Bob
			{Amount: 470, Recipient: "alice_address"}, // 470 change to Alice
		},
		Timestamp: time.Now(),
	}

	txID, err := calculateTxID(tx)
	if err != nil {
		fmt.Println("Error calculating transaction ID:", err)
		return
	}
	tx.ID = txID

	fmt.Printf("\nBefore transaction: ")
	fmt.Printf("Alice's balance: %f\n", utxoSet.GetBalance("alice_address"))
	fmt.Printf("Bob's balance: %f\n", utxoSet.GetBalance("bob_address"))

	// Process transaction
	// 1. Remove spent UTXOs
	for _, input := range tx.Inputs {
		utxoSet.SpendUTXO(input.TxID, input.Index)
	}

	// 2. Create new UTXOs
	for i, output := range tx.Outputs {
		newUTXO := &UTXO{
			TxID:   tx.ID,
			Index:  i,
			Amount: output.Amount,
			Owner:  output.Recipient,
		}
		utxoSet.AddUTXO(newUTXO)
	}

	fmt.Printf("\nAfter transaction:")
	fmt.Printf("Alice's balance: %f\n", utxoSet.GetBalance("alice_address"))
	fmt.Printf("Bob's balance: %f\n", utxoSet.GetBalance("bob_address"))

	// Run the full cryptocurrency simulation
	fmt.Println("\n=== Running Cryptocurrency Simulation ===")
	runCryptoCurrency()
}

func runCryptoCurrency() {
	// Create cryptocurrency system
	crypto := &CryptoCurrency{}
	crypto.Initialize()

	// Create wallets
	alice, err := NewWallet()
	if err != nil {
		fmt.Println("Error creating Alice's wallet:", err)
		return
	}

	bob, err := NewWallet()
	if err != nil {
		fmt.Println("Error creating Bob's wallet:", err)
		return
	}

	// Fund Alice's wallet from genesis
	genesisUTXO := &UTXO{
		TxID:   "genesis",
		Index:  0,
		Amount: 100,
		Owner:  alice.Address,
	}
	crypto.UTXOSet.AddUTXO(genesisUTXO)
	crypto.Wallets = []*Wallet{alice, bob}

	// Create miner
	miner := &Miner{Address: "miner_address", Mempool: crypto.Mempool}
	crypto.Miners = []*Miner{miner}

	// Send some transactions
	fmt.Println("Initial balances:")
	fmt.Printf("Alice: %f\n", alice.GetBalance(crypto.UTXOSet))
	fmt.Printf("Bob: %f\n", bob.GetBalance(crypto.UTXOSet))

	// Alice sends money to Bob
	if err := crypto.SendMoney(alice, bob.Address, 25); err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
		return
	}

	// Mine a block to process the transaction
	fmt.Println("\nMining block...")
	crypto.MineNextBlock()

	fmt.Println("Final balances:")
	fmt.Printf("Alice: %f\n", alice.GetBalance(crypto.UTXOSet))
	fmt.Printf("Bob: %f\n", bob.GetBalance(crypto.UTXOSet))
}
