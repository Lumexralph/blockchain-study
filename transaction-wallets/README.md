# Transactions and Wallets

- UTXO = Unspent Transaction Output
- blockchain balances aren't stored anywhere - they're calculated by examining transaction history.
- My private key IS my money - there's no account, just cryptographic proof of ownership
- Transactions consume UTXOs and create new UTXOs
- UTXOs work like physical cash - I spend specific 'bills' and get change back!
- Transaction: A message that transfers ownership from one key to another
- Addresses: Public identifiers derived from public keys
- UTXOs work like physical cash - you spend specific "bills" and get change back.

## Why Transaction Fees?

- Incentive for miners to include transactions in blocks
- Spam prevention - makes attacks expensive
- Priority system - higher fees get processed faster

Fee Calculation:
Transaction Fee = Total Inputs - Total Outputs
If Alice has a 100-coin UTXO and wants to send 30 coins to Bob:

Input: 100 coins
Outputs: 30 coins (to Bob) + 69 coins (change to Alice)
Fee: 100 - 30 - 69 = 1 coin (goes to miner)

## What Makes a Transaction Valid?

- Input Verification: Referenced UTXOs must exist and be unspent
- Signature Verification: Spender must prove ownership of input UTXOs
- Balance Check: Total inputs ≥ Total outputs
- No Double Spending: UTXOs can only be spent once

## How Transactions Spread Through Network:

- Wallet creates and signs transaction
- Broadcasts to connected peers
- Peers validate and relay to their peers
- Miners collect transactions in mempool
- Mining includes transactions in blocks

