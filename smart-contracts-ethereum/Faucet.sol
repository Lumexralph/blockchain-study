// SPDX-License-Identifier: CC-BY-SA-4.0

// Version of Solidity compiler this program was written for.
pragma solidity ^0.8.0;

// Our first contract is a faucet!
contract Faucet {
    address private owner;

    // Initialize the contract: set the owner
    // It is executed only once when the contract is deployed.
    constructor() {
        // The person (EOA) who deploys or created the contract is the owner.
        owner = msg.sender;
    }

    // Accept any incoming amount of ether.
    receive() external payable {}

    function withdraw(uint withdraw_amount) public {
        // Limit withdrawal amount.
        // 100000000000000000 wei = 0.1 ether
        require(withdraw_amount <= 0.1 ether, "Withdrawal amount exceeds limit");

        // Send the amount to the address that requested it.
        payable(msg.sender).transfer(withdraw_amount);
    }

    // Modifiers are most often used to create conditions that apply to many functions within a contract.
    // This is the basic design pattern for access control, allowing only the owner of a contract to execute
    // any function that has the onlyOwner modifier.
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    function destroy() public onlyOwner {
        // Only allow the owner to destroy the contract.
        // Destroy the contract and send the remaining funds to the owner.
        selfdestruct(payable(owner));
    }
}
