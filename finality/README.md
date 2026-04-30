# Finality

In the blockchain setting, finality is the affirmation that all well-formed blocks will not be revoked once committed to the blockchain.

Ref: https://medium.com/mechanism-labs/finality-in-blockchain-consensus-d1f83c120a9a

Types of Finality:

Probabilistic Finality (Bitcoin, PoW):

Becomes more certain over time
Never 100% guaranteed
Wait for N confirmations


Deterministic Finality (PBFT, PoS):

Instant finality
Cannot be reversed
Requires 2/3+ agreement


Economic Finality (PoS with slashing):

Reversing costs money
Attackers lose stake
Practical finality
