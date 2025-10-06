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
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
}
