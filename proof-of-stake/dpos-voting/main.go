package main

import (
	"fmt"
	"sort"
)

type Candidate struct {
	Addr  string
	Votes int
}

type Delegator struct {
	Addr     string
	Tokens   int
	VotedFor string
}

type DPoS struct {
	Candidates    []*Candidate
	Delegators    []*Delegator
	NumValidators int
}

func NewDPoS(numValidators int) *DPoS {
	return &DPoS{
		Candidates:    make([]*Candidate, 0),
		Delegators:    make([]*Delegator, 0),
		NumValidators: numValidators,
	}
}

func (dpos *DPoS) RegisterCandidate(addr string) {
	candidate := &Candidate{
		Addr:  addr,
		Votes: 0,
	}
	dpos.Candidates = append(dpos.Candidates, candidate)
	fmt.Printf("%s registered as validator candidate\n", addr)
}

func (dpos *DPoS) Vote(delegator, candidate string, tokens int) {
	// Add delegator.
	del := &Delegator{
		Addr:     delegator,
		Tokens:   tokens,
		VotedFor: candidate,
	}
	dpos.Delegators = append(dpos.Delegators, del)

	// Add votes to candidate.
	for _, cand := range dpos.Candidates {
		if cand.Addr != candidate {
			continue
		}

		cand.Votes += tokens
		fmt.Printf("%s voted for %s with %d tokens\n", delegator, candidate, tokens)
		return
	}
}

func (dpos *DPoS) ElectValidators() []*Candidate {
	// Sort candidates by votes.
	sort.Slice(dpos.Candidates, func(i, j int) bool {
		return dpos.Candidates[i].Votes > dpos.Candidates[j].Votes
	})

	// Select top N candidates as validators.
	elected := dpos.Candidates
	if len(elected) > dpos.NumValidators {
		elected = elected[:dpos.NumValidators]
	}

	fmt.Println("\nElected Validators:")
	for i, validator := range elected {
		fmt.Printf("%d. %s (votes: %d)\n", i+1, validator.Addr, validator.Votes)
	}

	return elected
}

func main() {
	// Create DPoS system with 3 validators.
	dpos := NewDPoS(3)

	// Register candidates.
	dpos.RegisterCandidate("Alice")
	dpos.RegisterCandidate("Bob")
	dpos.RegisterCandidate("Charlie")
	dpos.RegisterCandidate("Dave")
	dpos.RegisterCandidate("T")
	dpos.RegisterCandidate("Olu")

	fmt.Println("--------------------------------------------------------------------------------")

	// Delegators vote for candidates.
	dpos.Vote("Delegator1", "Alice", 100)
	dpos.Vote("Delegator2", "Bob", 50)
	dpos.Vote("Delegator3", "Charlie", 30)
	dpos.Vote("Delegator4", "Dave", 20)
	dpos.Vote("Delegator5", "T", 10)
	dpos.Vote("Delegator6", "Olu", 5)

	// Elect validators.
	dpos.ElectValidators()
}
