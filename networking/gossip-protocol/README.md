# Gossip Protocols

Gossip (epidemic) protocols spread information probabilistically: each node that learns something new forwards it to a random subset of peers.
This is the backbone of how both Bitcoin and Ethereum propagate transactions and blocks.

There are two fundamental models to understand:

Anti-entropy (pull/push-pull): Nodes periodically compare state with a random peer and reconcile differences. Guarantees eventual consistency but at a higher per-cycle cost.
Rumor mongering (push): A node that receives new information actively pushes it to random peers for a limited number of rounds. Lower cost, but some messages may fail to reach every node.

## Nuggets

We should point out that extensive replication of a database is expensive. It should be avoided
whenever possible by hierarchical decomposition of the database or by caching.

- Knowing the network capacity is important, as the network grows by maybe more hosts, you might need to check traffic intensive applications.
- Lood on the network can lead to congestion and packet loss, which can cause delays and reduce the overall performance of the network. Thereby, unresponsvie systes or delayed systems.


