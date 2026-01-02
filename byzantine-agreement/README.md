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
