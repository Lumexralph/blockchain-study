# Proof of Stake (PoS)

Though some blockchain networks like Bitcoin use Proof of Work (PoW) as their consensus mechanism, others utilize Proof of Stake (PoS)
to achieve distributed consensus. PoS is designed to be more energy-efficient and scalable compared to PoW.

Instead of computing power securing the network (PoW), validators "stake" their cryptocurrency holdings as collateral to propose and validate new blocks.
Every computer in the network must agree upon each new block and the chain as a whole. These computers are known as "nodes".
Nodes ensure everyone interacting with the blockchain has the same data. To accomplish this distributed agreement, blockchains need a consensus mechanism.

Ethereum uses a proof-of-stake-based consensus mechanism. Anyone who wants to add new blocks to the chain must stake
ETH - the native currency in Ethereum - as collateral and run validator software. These "validators" can then be randomly
selected to propose blocks that other validators check and add to the blockchain. There is a system of rewards and penalties
that strongly incentivize participants to be honest and available online as much as possible.

## Validtors
To participate as a validator on Ethereum network, a user must deposit 32 ETH into the deposit contract and run three separate pieces of software: an execution client,
a consensus client, and a validator client. On depositing their ETH, the user joins an activation queue that limits the rate of new validators joining the network.
Once activated, validators receive new blocks from peers on the Ethereum network. The transactions delivered in the block are re-executed to check that the proposed
changes to Ethereum's state are valid, and the block signature is checked. The validator then sends a vote (called an attestation) in favor of that block across the network.

## Fork Choice
When the network performs optimally and honestly, there is only ever one new block at the head of the chain, and all validators attest to it. However,
it is possible for validators to have different views of the head of the chain due to network latency or because a block proposer has equivocated.
Therefore, consensus clients require an algorithm to decide which one to favor. The algorithm used in proof-of-stake Ethereum is called LMD-GHOSTopens
in a new tab, and it works by identifying the fork that has the greatest weight of attestations in its history.

In PoS-based public blockchains, a set of validators take turns proposing and voting on the next block, and the weight of each validator's vote depends on the size of its deposit (i.e. stake).

## PoW vs PoS

Proof of Work:
- Security from computational cost
- Energy intensive
- 51% attack = 51% of hash power
- No finality (probabilistic)

Proof of Stake:
- Security from economic cost
- Energy efficient
- 51% attack = 51% of staked value
- Can have finality (with proper design)

Key Insight: PoS replaces computational cost with economic cost, making attacks expensive without burning energy.

### PoS Variants:

### Delegated Proof of Stake \(DPoS\)

- Token holders vote for validators
- Fixed number of validators \(typically 21-100\)
- Fast but more centralized
- Used by: EOS, TRON

### Bonded Proof of Stake

- Validators lock up stake for set periods
- Stake cannot be withdrawn immediately \(unbonding period\)
- Used by: Cosmos, Polkadot

### Liquid Proof of Stake

- Delegators can change validators at any time
- No unbonding period
- More flexible

Comparison Table:

Algorithm    | Validators | Speed  | Decentralization | Finality
-------------|-----------|--------|------------------|----------
PoW          | Anyone    | Slow   | High            | Probabilistic
PoS          | Stakers   | Fast   | Medium          | Depends
DPoS         | Elected   | Faster | Lower           | Fast
PBFT         | Known     | Fast   | Low             | Deterministic

#### The Blockchain Trilemma:
\`\`\`text
        Decentralization
              /\\
             /  \\
            /    \\
           /      \\
          /________\\
Security                Scalability
\`\`\`

You can typically optimize for only 2 of 3:
- Bitcoin: Decentralization + Security (sacrifices scalability)
- DPoS: Security + Scalability (sacrifices decentralization)
- Avalanche: Claims all 3 (we'll explore soon!)

