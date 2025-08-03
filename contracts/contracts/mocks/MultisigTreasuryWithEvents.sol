// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title MultisigTreasuryWithEvents
 * @dev Example implementation showing how MultisigTreasury should emit events
 * This demonstrates the missing events from EventEmissionFixed.sol
 */
contract MultisigTreasuryWithEvents {
    bool public restrictFunctionCalls;
    
    // Events that should be added to the actual MultisigTreasury contract
    event FunctionRestrictionsToggled(bool enabled);
    
    constructor() {
        restrictFunctionCalls = false;
    }
    
    /**
     * @dev Toggle function call restrictions and emit event
     */
    function toggleFunctionRestrictions() external {
        restrictFunctionCalls = !restrictFunctionCalls;
        emit FunctionRestrictionsToggled(restrictFunctionCalls);
    }
}