package main

import (
	"fmt"
	"math/big"
	"math/rand"
)

type Validator struct {
	Address   string
	Stake     *big.Int
	IsSlashed bool
}

type PoSConsensus struct {
	Validators      []*Validator
	TotalStake      *big.Int
	MinStake        *big.Int
	SlashingPenalty float64 // Percentage of stake to slash
}

func NewPoSConsensus(minStake *big.Int) *PoSConsensus {
	return &PoSConsensus{
		Validators:      make([]*Validator, 0),
		TotalStake:      big.NewInt(0),
		MinStake:        minStake,
		SlashingPenalty: 0.05, // 5% penalty
	}
}

func (pos *PoSConsensus) AddValidator(address string, stake *big.Int) error {
	if stake.Cmp(pos.MinStake) < 0 {
		return fmt.Errorf("stake below minimum required: %s < %s", stake.String(), pos.MinStake.String())
	}

	validator := &Validator{
		Address:   address,
		Stake:     new(big.Int).Set(stake),
		IsSlashed: false,
	}

	pos.Validators = append(pos.Validators, validator)
	pos.TotalStake.Add(pos.TotalStake, stake)

	fmt.Printf("Validator %s joined with stake %s\n", address, stake.String())
	return nil
}

func (pos *PoSConsensus) SelectValidator() *Validator {
	// Weighted random selection based on stake.
	// Generate a random value between 0 and TotalStake.
	randValue := rand.Int63n(pos.TotalStake.Int64())
	var currentSum int64 = 0

	for _, validator := range pos.Validators {
		if validator.IsSlashed {
			continue // Skip slashed validators
		}
		currentSum += validator.Stake.Int64()
		if randValue < currentSum {
			return validator
		}
	}
	return pos.Validators[0] // Fallback to the first validator if something goes wrong.
}

func (pos *PoSConsensus) SlashValidator(addr string) {
	for _, validator := range pos.Validators {
		if validator.Address != addr {
			continue
		}

		// Calculate the slashing amount.
		slashAmount := new(big.Int).Set(validator.Stake)
		slashAmount.Mul(slashAmount, big.NewInt(int64(pos.SlashingPenalty*100)))
		slashAmount.Div(slashAmount, big.NewInt(100))

		// Reduce the stake.
		validator.Stake.Sub(validator.Stake, slashAmount)
		pos.TotalStake.Sub(pos.TotalStake, slashAmount)
		validator.IsSlashed = true

		fmt.Printf("Validator %s slashed! Lost %s coins \n", addr, slashAmount.String())
		return
	}
}

func (pos *PoSConsensus) ProposeBlock() string {
	validator := pos.SelectValidator()
	fmt.Printf("Validator %s selected to propose block (stake: %s)\n", validator.Address, validator.Stake.String())
	return validator.Address
}

func main() {
	minStake := big.NewInt(1000)
	pos := NewPoSConsensus(minStake)

	// Add validators with different stakes.
	if err := pos.AddValidator("Alice", big.NewInt(5000)); err != nil {
		fmt.Println("Error adding validator:", err)
	}
	if err := pos.AddValidator("Bob", big.NewInt(3000)); err != nil {
		fmt.Println("Error adding validator:", err)
	}
	if err := pos.AddValidator("Charlie", big.NewInt(2000)); err != nil {
		fmt.Println("Error adding validator:", err)
	}

	fmt.Printf("\nTotal Stake: %s\n", pos.TotalStake.String())

	// Simulate block proposals.
	fmt.Println("Block Proposals:")
	for i := 0; i < 5; i++ {
		fmt.Printf("Block %d: ", i)
		pos.ProposeBlock()
	}

	// Simulate slashing a validator.
	fmt.Println("\nSlashing Bob for misbehavior...")
	pos.SlashValidator("Bob")

	fmt.Printf("\nTotal Stake after slashing: %s\n", pos.TotalStake.String())
}
