# Blockchain Infrastructure Study

A blockchain is a public database that is updated and shared across many computers in a network.
"Block" refers to data and state being stored in consecutive groups known as "blocks". If you send ETH to someone else, the transaction data needs to be added to a block to be successful.
"Chain" refers to the fact that each block cryptographically references its parent. In other words, blocks get chained together. The data in a block cannot change without changing all subsequent blocks, which would require the consensus of the entire network.
Every computer in the network must agree upon each new block and the chain as a whole. These computers are known as "nodes". Nodes ensure everyone interacting with the blockchain has the same data. To accomplish this distributed agreement, blockchains need a consensus mechanism.

This repository contains a study project focused on blockchain infrastructure, including basic blockchain implementation, proof-of-work mining, and cryptographic utilities. The project is divided into three main components:

1. **blockchain-v1**: A simple blockchain implementation with transaction handling and validation.
2. **blockchain-v2**: An enhanced version of the blockchain that includes proof-of-work mining and handles blockchain forks.
3. **crypto**: A set of cryptographic utilities demonstrating SHA-256 hashing and ECDSA key generation and signing.
4. **digital-identity**: A basic digital identity management system using public/private key pairs.
5. **mining-system**: A simple mining system that simulates mining operations and block creation.
6. **utxo**: A basic implementation of the UTXO (Unspent Transaction Output) model used in cryptocurrencies. See [transaction-wallets](./transaction-wallets/README.md).
