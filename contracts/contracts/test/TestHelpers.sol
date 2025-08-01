// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract InfiniteGas {
    function infiniteLoop() external pure {
        uint256 i = 0;
        while(i < 1) {
            i = 0;
        }
    }
}

contract ReentrancyAttacker {
    address public treasury;
    
    constructor(address _treasury) {
        treasury = _treasury;
    }
    
    receive() external payable {
        // Try to re-enter
        (bool success,) = treasury.call(
            abi.encodeWithSignature("executeTransaction(uint256)", 0)
        );
    }
    
    function attack() external payable {}
}

contract SelfDestructContract {
    function destruct(address payable recipient) external {
        selfdestruct(recipient);
    }
}