// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface ITokenRescuer {
    function rescue(address target, bytes calldata data) external returns (bool success, bytes memory result);
}

contract MaliciousReentrant {
    ITokenRescuer public rescuer;
    bool public hasReentered = false;
    
    constructor(address _rescuer) {
        rescuer = ITokenRescuer(_rescuer);
    }
    
    fallback() external {
        if (!hasReentered) {
            hasReentered = true;
            rescuer.rescue(address(this), "");
        }
    }
    
    function transferToken(address, address, uint256) external {
        if (!hasReentered) {
            hasReentered = true;
            rescuer.rescue(address(this), "");
        }
    }
}