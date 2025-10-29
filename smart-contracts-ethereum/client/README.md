# Client

The client is the entry point to the Ethereum network. The client is required to broadcast transactions and read blockchain data.

## Accounts

Accounts on Ethereum are either wallet addresses or smart contract addresses. They look like 0x71c7656ec7ab88b098defb751b7401b5f6d8976f,
and they're what you use for sending ETH to another user and also are used for referring to a smart contract on the blockchain when needing to interact with it.
They are unique and are derived from a private key.

Each Externally-Owned Account (EOA) is a public-private key pair, where the public key is used to derive a unique address for the user and the private key is used to protect the account and securely sign messages.
Therefore, in order to use Ethereum, it is first necessary to generate an EOA (hereafter, "account").

Wallet:  A file that stores your private keys. You can have multiple accounts in a wallet.
Addresses: Used to receive and send transactions on the network. They are derived from public/private ECDSA keys.

In order to make some transactions, the user must fund their account with ether. On Ethereum mainnet, ether can only be obtained in three ways:

1) by receiving it as a reward for mining/validating;
2) receiving it in a transfer from another Ethereum user or contract;
3) receiving it from an exchange, having paid for it with fiat money. On Ethereum testnets, the ether has no real world value so it
4) can be made freely available via faucets. Faucets allow users to request a transfer of testnet ether to their account.

## Interacting with the blockchain
Geth provides JSON-RPC APIs. JSON-RPC is a way to execute specific tasks by sending instructions to Geth in the form of JSON objects. RPC stands for "Remote Procedure Call" and it refers to the ability to send these JSON-encoded instructions from locations outside of those managed by Geth.
There are three transport protocols that can be used to connect to Geth:

- IPC (Inter-Process Communication): Provides unrestricted access to all APIs, but only works when the client is run on the same host as the geth node.
- HTTP: By default provides access to the **eth**, **web3** and **net** method namespaces.
- Websocket: By default provides access to the **eth**, **web3** and **net** method namespaces.


Key Insight: Ethereum tracks account state globally, not just transaction history.