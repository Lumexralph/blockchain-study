# Byzantine Agreement

Byzantine Agreement is a fundamental concept in distributed systems and fault-tolerant computing. It addresses the challenge
of achieving consensus among distributed nodes, even when some nodes may behave arbitrarily or maliciously (known as Byzantine faults).

Key Concepts:

- Byzantine Fault: Node behaves arbitrarily (crash, send wrong data, malicious)
- Byzantine Agreement: Reaching consensus despite Byzantine faults
- Safety: Bad things never happen (no conflicting decisions)
- Liveness: Good things eventually happen (system makes progress)
- `f` Byzantine nodes: Can tolerate up to `f` faults with `3f+1` total nodes
- Byzantine agreement requires a supermajority (2/3 at least) to overcome malicious nodes.

The Math:

To tolerate `f` Byzantine faults:
- Need at least 3f + 1 nodes total
- Need at least 2f + 1 honest nodes

Example: 4 nodes can tolerate 1 Byzantine fault
- 3(1) + 1 = 4 total nodes needed
- If 1 is Byzantine, 3 remain honest (majority)

Common Algorithms:
- Practical Byzantine Fault Tolerance (PBFT)
- Raft (not Byzantine, but related consensus algorithm)
- Paxos (not Byzantine, but related consensus algorithm)

## Understanding Practical Byzantine Fault Tolerance (PBFT)

It's a replication algorithm that is able to tolerate Byzantine faults. Faults resulting from malicious attacks and software errors are examples of Byzantine faults.
These faults can cause nodes to behave arbitrarily, including lying, colluding, or failing to respond.
Since malicious attacks and software errors can cause faulty nodes to exhibit Byzantine (i.e., arbitrary) behavior, Byzantine-fault-tolerant algorithms are increasingly important.

The algorithm is based on a state machine replication approach, where each node maintains a copy of the system state and processes client requests in a coordinated manner.

### Insights from the PBFT Paper

How you ensure a message is authentic: signing a digest of a message and appending it to the plaintext of the message
rather than signing the full message.

- All replicas know the others’ public keys to verify signatures.
