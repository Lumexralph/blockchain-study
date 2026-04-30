package main

import "fmt"

func main() {
	// Create 4 nodes (can tolerate 1 fault)
	nodes := []*node{
		{id: 0, value: 1, isFaulty: false},
		{id: 1, value: 1, isFaulty: false},
		{id: 2, value: 1, isFaulty: false},
		{id: 3, value: 0, isFaulty: true}, // Byzantine node
	}

	f := 1 // Tolerate 1 fault
	result := simpleByzantineAgreement(nodes, f)

	fmt.Printf("Consensus value: %d\n", result)
	fmt.Printf("Despite 1 Byzantine node, 3 honest nodes reached agreement!\n")
}

type node struct {
	id       int
	isFaulty bool
	value    int
}

type message struct {
	From  int
	Value int
	Round int
}

// SimpleByzantineAgreement implements a basic BFT algorithm.
func simpleByzantineAgreement(nodes []*node, f int) int {
	n := len(nodes)

	if n < 3*f+1 {
		panic("Not enough nodes for Byzantine fault tolerance")
	}
	// Round 1: Each node broadcasts its value
	votes := make(map[int]int) // value -> count
	for _, nod := range nodes {
		value := nod.value

		// Byzantine node might send wrong value.
		if nod.isFaulty {
			value = -1 // Send malicious value.
		}
		votes[value]++ // Faulty nodes will always vote for -1 with their count incremented.
	}

	// Find majority value (needs > 2f+1 votes)
	threshold := 2*f + 1
	for value, count := range votes {
		if count >= threshold {
			return value // Consensus reached with this value.
		}
	}

	return -1 // No consensus
}
