# Smart Contracts and Ethereum

Bitcoin lets you send money without banks. But what if the blockchain could run programs? What if you could create self-executing contracts, decentralized exchanges, or digital organizations - all without any central authority?
We will understand how Ethereum transformed blockchain from a payment system into a world computer that can run any program.

## What Are Smart Contracts?

A smart contract is a program that runs on the blockchain. Once deployed, it executes automatically when conditions are met, and no one can stop or change it.

Smart contracts are the fundamental building blocks of Ethereum's application layer. They are computer programs stored on the blockchain that follow "if this then that" logic, and are guaranteed to execute according to the rules defined by its code, which cannot be changed once created.

- https://ethereum.org/smart-contracts/
- https://www.ibm.com/think/topics/smart-contracts
- See: https://www.youtube.com/watch?v=pWGLtjG-F5c

Traditional Contract:
"Alice will pay Bob $1000 on January 1st, 2025"
- Requires trust
- Requires enforcement
- Can be disputed

Smart Contract:
if (date == "2025-01-01") {
transfer(alice, bob, 1000)
}

- Executes automatically
- No trust needed
- No disputes possible

### Key Concepts:

- Deterministic Execution: Same input always produces same output
- Immutable Code: Once deployed, code cannot be changed
- Transparent: Everyone can verify the code
- Trustless: No intermediary needed

### Oracle Services

Provide external data to smart contracts. For example, a weather oracle can provide weather data to a smart contract for insurance payouts.

### Smart Contract

there are two different types of accounts in Ethereum: externally owned accounts (EOAs) and contract accounts.
EOAs are controlled by users, often via software such as a wallet application that is external to the Ethereum platform.
In contrast, contract accounts are controlled by program code (also commonly referred to as “smart contracts”) that
is executed by the Ethereum Virtual Machine. In short, EOAs are simple accounts without any associated code or data storage,
whereas contract accounts have both associated code and data storage.
EOAs are controlled by transactions created and cryptographically signed with a private key in the "real world"
external to and independent of the protocol, whereas contract accounts do not have private keys and so "control themselves"
in the predetermined way prescribed by their smart contract code. Both types of accounts are identified by an Ethereum address.

“smart contracts” to refer to immutable computer programs that run deterministically in the context of an Ethereum Virtual Machine (EVM)
as part of the Ethereum network protocol—i.e., on the decentralized Ethereum world computer.

- Smart contracts are simply computer programs. The word “contract” has no legal meaning in this context.
- Once deployed, the code of a smart contract cannot change. Unlike with traditional software, the only way to modify a smart contract is to deploy a new instance.
- The EVM runs as a local instance on every Ethereum node, but because all instances of the EVM operate on the same initial state and produce the same final state, the system as a whole operates as a single "world computer."
This means they have one single, consistent view of the entire Ethereum state at any point in time.

### Life Cycle of a Smart Contract

1. **Development**: Write the smart contract code using a programming language like Solidity or Vyper.
2. But in order to run, they must be compiled to the low-level bytecode that runs in the EVM. Once compiled, they are deployed on the Ethereum platform using a special contract creation transaction, which is identified as such by being sent to the special contract creation address, namely 0x0.
3. Each contract is identified by an Ethereum address, which is derived from the **contract creation transaction** as a function of the originating account and nonce.
4. The Ethereum address of a contract can be used in a transaction as the recipient, sending funds to the contract or calling one of the contract’s functions. Note that, unlike with EOAs, there are no keys associated with an account created for a new smart contract.
5. Importantly, contracts only run if they are called by a transaction. All smart contracts in Ethereum are executed, ultimately, because of a transaction initiated from an EOA. A contract can call another contract that can call another contract, and so on, but the first contract in such a chain of execution will have always been called by a transaction from an EOA.
6. It is also worth noting that smart contracts are not executed "in parallel" in any sense—the Ethereum world computer can be considered to be a single-threaded machine.
7. Transactions are atomic, they are either successfully terminated or reverted. An all-or-nothing operation.
8. A failed transaction is still recorded as having been attempted, and the ether spent on gas for the execution is deducted from the originating account, but it otherwise has no other effects on contract or account state.
9. it is important to remember that a contract’s code cannot be changed. However, a contract can be “deleted,” removing the code and its internal state (storage) from its address, leaving a blank account. Any transactions sent to that account address after the contract has been deleted do not result in any code execution, because there is no longer any code there to execute. To delete a contract, you execute an EVM opcode called SELFDESTRUCT (previously called SUICIDE).
10. Deleting a contract in this way does not remove the transaction history (past) of the contract, since the blockchain itself is immutable.
11. It is also important to note that the SELFDESTRUCT capability will only be available if the contract author programmed the smart contract to have that functionality. If the contract’s code does not have a SELFDESTRUCT opcode, or it is inaccessible, the smart contract cannot be deleted.

### Example of writing Smart Contract in Solidity

This can be found in the [Faucet.sol](smart-contracts-ethereum/Faucet.sol) contract code. We will use the Solidity compiler on the command line to compile our contract directly:

```
solc --optimize --bin smart-contracts-ethereum/Faucet.sol
```

The result that solc produces is a hex-serialized binary that can be submitted to the Ethereum blockchain.

To summarize, a contract’s life cycle starts with a creation transaction from an EOA or contract account. If there is a constructor, it is executed as part of contract creation, to initialize the state of the contract as it is being created, and is then discarded.
Contracts are destroyed by a special EVM opcode called SELFDESTRUCT. n Solidity, this opcode is exposed as a high-level built-in function called selfdestruct, which takes one argument: the address to receive any ether balance remaining in the contract account.
Note that you must explicitly add this command to your contract if you want it to be deletable—this is the only way a contract can be deleted, and it is not present by default. In this way, users of a contract who might rely on a contract being there forever can be certain that a contract can’t be deleted if it doesn’t contain a SELFDESTRUCT opcode.

When a contract terminates with an error, all the state changes (changes to variables, balances, etc.) are reverted, all the way up the chain of contract calls if more than one contract was called. This ensures that transactions are atomic, meaning they either complete successfully or have no effect on state and are reverted entirely.

### The Ethereum Contract ABI

In computer software, an application binary interface is an interface between two program modules; often, between the operating system and user programs. An ABI defines how data structures and functions are accessed in machine code.
In Ethereum, the ABI is used to encode contract calls for the EVM and to read data out of transactions. The purpose of an ABI is to define the functions in the contract that can be invoked and describe how each function will accept arguments and return its result.

We use the solc command-line Solidity compiler to produce the ABI for our Faucet.sol example contract:
```
solc --abi smart-contracts-ethereum/Faucet.sol

======= smart-contracts-ethereum/Faucet.sol:Faucet =======
Contract JSON ABI
[{"inputs":[{"internalType":"uint256","name":"withdraw_amount","type":"uint256"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}
```
This JSON can be used by any application that wants to access the Faucet contract once it is deployed. Using the ABI, an application such as a wallet or DApp browser can construct transactions that call the functions in Faucet with the correct arguments and argument types.

All that is needed for an application to interact with a contract is an **ABI** and the **address** where the contract has been deployed.

### Events

Events are a way for your smart contract to communicate that something happened on the blockchain to your app front-end, which can be 'listening' for certain events and take action when they happen.
When you call a function on a smart contract, you can emit an event. The event is stored in the transaction's log, which is a special data structure in the Ethereum blockchain. Your front-end application can then listen for these events and update the UI accordingly.

When a transaction completes (successfully or not), it produces a transaction receipt. The transaction receipt contains log entries that provide information about the actions that occurred during the execution of the transaction.
Events are the Solidity high-level objects that are used to construct these logs.

Events are a very useful mechanism, not only for intra-contract communication, but also for debugging during development.

### Calling Other Contracts (send, call, callcode, delegatecall)

Calling other contracts from within your contract is a very useful but potentially dangerous operation.
We’ll examine the various ways you can achieve this and evaluate the risks of each method.
In short, the risks arise from the fact that you may not know much about a contract you are calling into or that
is calling into your contract. When writing smart contracts, you must keep in mind that while you may mostly expect
to be dealing with EOAs, there is nothing to stop arbitrarily complex and perhaps malign contracts from calling into and
being called by your code.

You mostly don't have control over who calls your contract, so you must be very careful about trusting any data that comes from external contracts.

### Gas and Fees

Every operation that a smart contract performs costs a certain amount of gas. When you deploy a contract or call a function on a contract,
you need to specify a gas limit and pay for the gas used.

Gas is a resource constraining the maximum amount of computation that Ethereum will allow a transaction to consume. If the gas limit is exceeded during computation, the following series of events occurs:

- An "out of gas" exception is thrown.

- The state of the contract prior to execution is restored (reverted).

- All ether used to pay for the gas is taken as a transaction fee; it is not refunded.

Because gas is paid by the user who initiates the transaction, users are discouraged from calling functions that have a high gas cost.
It is thus in the programmer’s best interest to minimize the gas cost of a contract’s functions.
To this end, there are certain practices that are recommended when constructing smart contracts, to minimize the gas cost of a function call.

It is recommended that you evaluate the gas cost of functions as part of your development workflow, to avoid any surprises when deploying contracts to the mainnet.

Futher Reading: https://github.com/ethereumbook/ethereumbook/blob/develop/07smart-contracts-solidity.asciidoc#gas-considerations

### Block Explorers

Block explorers are web applications that allow you to view information about blocks, transactions, addresses, and smart contracts on the Ethereum or any blockchain.
The block contains the data of all the transactions that have been mined within the allocated block time. The block explorer also allows you to view events that were emitted during the execution of the smart contract as well as things such as how much was paid for the gas and amount of ether was transacted, etc.
Some popular Ethereum block explorers include:
- [Etherscan](https://etherscan.io/)
- [Ethplorer](https://ethplorer.io/)
- [Blockscout](https://blockscout.com/)

For Avalanche, you can use:
- [SnowTrace](https://snowtrace.io/)
- [Avalanche Explorer](https://subnets.avax.network/)

### Clients

See: [client](./client/README.md )