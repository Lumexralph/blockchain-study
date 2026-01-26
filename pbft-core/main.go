package main

import (
	"crypto/sha256"
	"fmt"
	"log"
)

type MessageType int

const (
	PrePrepare MessageType = iota
	Prepare
	Commit
)

type PBFTMessage struct {
	Type      MessageType
	View      int
	Sequence  int
	BlockHash string
	NodeID    int
}

type PBFTNode struct {
	ID              int
	View            int
	PrepareMessages map[string][]*PBFTMessage // BlockHash -> Messages mapping/records
	CommitMessages  map[string][]*PBFTMessage
	IsPrimary       bool
	MaxNumFaults    int // f
}

func NewPBFTNode(id, f int, isPrimary bool) *PBFTNode {
	return &PBFTNode{
		ID:              id,
		View:            0,
		PrepareMessages: make(map[string][]*PBFTMessage),
		CommitMessages:  make(map[string][]*PBFTMessage),
		IsPrimary:       isPrimary,
		MaxNumFaults:    f,
	}
}

func (n *PBFTNode) ProposeBlock(data string) *PBFTMessage {
	if !n.IsPrimary {
		log.Fatal("Not a primary, cannot propose block")
	}

	blockHash := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))

	return &PBFTMessage{
		Type:      PrePrepare,
		View:      n.View,
		Sequence:  0,
		BlockHash: blockHash,
		NodeID:    n.ID,
	}
}

func (n *PBFTNode) ReceivePrePrepare(msg *PBFTMessage) *PBFTMessage {
	// Validate a pre-prepare message.
	if msg.View != n.View {
		// My view is different, ignore.
		return nil
	}

	// Send a prepare message.
	return &PBFTMessage{
		Type:      Prepare,
		View:      n.View,
		Sequence:  msg.Sequence,
		BlockHash: msg.BlockHash,
		NodeID:    n.ID,
	}
}

func (n *PBFTNode) ReceivePrepare(msg *PBFTMessage) *PBFTMessage {
	// Collect PREPARE messages.
	n.PrepareMessages[msg.BlockHash] = append(n.PrepareMessages[msg.BlockHash], msg)

	// Check if we have 2f + 1 PREPARE messages.
	threshold := 2*n.MaxNumFaults + 1
	if len(n.PrepareMessages[msg.BlockHash]) >= threshold {
		// Prepared! Send a COMMIT message.
		return &PBFTMessage{
			Type:      Commit,
			View:      n.View,
			Sequence:  msg.Sequence,
			BlockHash: msg.BlockHash,
			NodeID:    n.ID,
		}
	}
	return nil
}

func (n *PBFTNode) ReceiveCommit(msg *PBFTMessage) bool {
	// Collect COMMIT messages.
	n.CommitMessages[msg.BlockHash] = append(n.CommitMessages[msg.BlockHash], msg)

	// Check if we have 2f + 1 COMMIT messages.
	threshold := 2*n.MaxNumFaults + 1
	if len(n.CommitMessages[msg.BlockHash]) >= threshold {
		// Committed!
		return true
	}
	return false
}

func main() {
	// Create a PBFT network with 4 nodes, tolerate 1 fault (f=1).
	f := 1
	nodes := []*PBFTNode{
		NewPBFTNode(0, f, true),  // Primary
		NewPBFTNode(1, f, false), // Replica
		NewPBFTNode(2, f, false), // Replica
		NewPBFTNode(3, f, false), // Replica
	}

	// Phase 1: PRE-PREPARE.
	primary := nodes[0]
	prePrepareMsg := primary.ProposeBlock("Block Data: transaction XYZ")
	fmt.Printf("Primary %d: Sent PRE-PREPARE for block %s\n", primary.ID, prePrepareMsg.BlockHash[:8])

	// Phase 2: PREPARE.
	var prepareMsgs []*PBFTMessage
	for _, node := range nodes[1:] {
		prepareMsg := node.ReceivePrePrepare(prePrepareMsg)
		if prepareMsg != nil {
			prepareMsgs = append(prepareMsgs, prepareMsg)
			fmt.Printf("Node %d: Sent PREPARE for block %s\n", node.ID, prepareMsg.BlockHash[:8])
		}
	}

	// Broadcast PREPARE messages to all nodes.
	var commitMsgs []*PBFTMessage
	for _, node := range nodes {
		for _, prepMsg := range prepareMsgs {
			commitMsg := node.ReceivePrepare(prepMsg)
			if commitMsg != nil {
				commitMsgs = append(commitMsgs, commitMsg)
				fmt.Printf("Node %d: Sent COMMIT for block %s\n", node.ID, commitMsg.BlockHash[:8])
			}
		}
	}

	// Phase 3: COMMIT.
	for _, node := range nodes {
		for _, comMsg := range commitMsgs {
			committed := node.ReceiveCommit(comMsg)
			if committed {
				fmt.Printf("Node %d: Committed block %s\n", node.ID, comMsg.BlockHash[:8])
			}
		}
	}

	fmt.Println("\nConsensus achieved! Block finalized.")
}

// Key Insight: PBFT provides finality through multiple rounds of voting, but requires O(n²) messages.
