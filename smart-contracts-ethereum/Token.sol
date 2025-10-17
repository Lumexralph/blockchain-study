// SPDX-License-Identifier: CC-BY-SA-4.0

// Version of Solidity compiler this program was written for.
pragma solidity ^0.8.0;

import "Faucet.sol";

contract Token is Mortal {
    Faucet private faucet;

    constructor() {
        // new will create the contract on the blockchain and return an object that you can use to reference it.
        // You can optionally specify the value of ether transfer on creation,
        // and pass arguments to the new contract’s constructor.
        faucet = (new Faucet){value: 1 ether}();
    }

    // Note: While you are the owner of the Token contract, the Token contract itself owns the new Faucet contract,
    // so only the Token contract can destroy it.
    function destroy() public onlyOwner {
        faucet.destroy();
    }
}
