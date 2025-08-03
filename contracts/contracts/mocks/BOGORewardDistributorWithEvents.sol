// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title BOGORewardDistributorWithEvents
 * @dev Example implementation showing how BOGORewardDistributor should emit events
 * This demonstrates the missing events from EventEmissionFixed.sol
 */
contract BOGORewardDistributorWithEvents {
    mapping(address => bool) public authorizedBackends;
    uint256 public dailyDistributed;
    uint256 public lastResetTimestamp;
    
    // Events that should be added to the actual BOGORewardDistributor contract
    event AuthorizedBackendSet(address indexed backend, bool authorized);
    event DailyLimitReset(uint256 timestamp, uint256 previousDistributed);
    
    constructor() {
        lastResetTimestamp = block.timestamp;
    }
    
    /**
     * @dev Set authorized backend and emit event
     */
    function setAuthorizedBackend(address backend, bool authorized) external {
        authorizedBackends[backend] = authorized;
        emit AuthorizedBackendSet(backend, authorized);
    }
    
    /**
     * @dev Reset daily limit and emit event (internal function made public for testing)
     */
    function resetDailyLimit() external {
        uint256 previousDistributed = dailyDistributed;
        dailyDistributed = 0;
        lastResetTimestamp = block.timestamp;
        emit DailyLimitReset(block.timestamp, previousDistributed);
    }
    
    /**
     * @dev Simulate daily limit reset when time passes
     */
    function _checkAndResetDailyLimit() internal {
        if (block.timestamp >= lastResetTimestamp + 24 hours) {
            uint256 previousDistributed = dailyDistributed;
            dailyDistributed = 0;
            lastResetTimestamp = block.timestamp;
            emit DailyLimitReset(block.timestamp, previousDistributed);
        }
    }
}