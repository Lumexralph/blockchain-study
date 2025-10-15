// SPDX-License-Identifier: CC-BY-SA-4.0

// Version of Solidity compiler this program was written for.
pragma solidity ^0.8.0;

contract Owned {
    address private owner;

    // Initialize the contract: set the owner
    // It is executed only once when the contract is deployed.
    constructor() {
        // The person (EOA) who deploys or created the contract is the owner.
        owner = msg.sender;
    }

    // Access control modifier.
    // Modifiers are most often used to create conditions that apply to many functions within a contract.
    // This is the basic design pattern for access control, allowing only the owner of a contract to execute
    // any function that has the onlyOwner modifier.
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }
}

contract Mortal is Owned {
    // We’ve chosen to make the addresses indexed, to allow searching and
    // filtering in any user interface built to access our contract.
    event Withdrawal(address indexed to, uint amount);
    event Deposit(address indexed from, uint amount);

    // Destroy the contract and send the remaining funds to the owner.
    function destroy() public onlyOwner {
        selfdestruct(payable(owner));
    }
}

// Our first contract is a faucet!
contract Faucet is Mortal {
    // Accept any incoming amount of ether.
    receive() external payable {
        emit Deposit(msg.sender, msg.value);
    }

    // Give out ether to anyone who asks.
    function withdraw(uint withdraw_amount) public {
        // Limit withdrawal amount.
        // 100000000000000000 wei = 0.1 ether
        // Additional error-checking code like this will increase gas consumption slightly,
        // but it offers better error reporting than if omitted.
        // You will need to find the right balance between gas consumption and verbose error checking
        // based on the expected use of your contract. In the case of a Faucet contract intended for a testnet,
        // we’d probably err on the side of extra reporting even if it costs more gas.
        // Perhaps for a mainnet contract we’d choose to be frugal with our gas usage instead.
        require(withdraw_amount <= 0.1 ether, "Withdrawal amount exceeds limit");
        require(address(this).balance >= withdraw_amount, "Insufficient contract balance");

        // Send the amount to the address that requested it.
        payable(msg.sender).transfer(withdraw_amount);
        // Incorporate the event data in the transaction logs.
        emit Withdrawal(msg.sender, withdraw_amount);
    }
}
