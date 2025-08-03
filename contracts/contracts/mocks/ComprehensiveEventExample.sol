// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

/**
 * @title ComprehensiveEventExample
 * @dev Example contract demonstrating comprehensive event emission patterns
 * This shows additional events that could enhance transparency across contracts
 */
contract ComprehensiveEventExample is Ownable, Pausable {
    uint256 public fee;
    address public treasury;
    
    // Comprehensive events for enhanced transparency
    event ConfigurationChanged(string indexed parameter, uint256 oldValue, uint256 newValue);
    event TreasuryChanged(address indexed oldTreasury, address indexed newTreasury);
    event ContractPaused(address indexed by);
    event ContractUnpaused(address indexed by);
    event EmergencyActionTaken(string indexed action, address indexed by);
    event FundsReceived(address indexed from, uint256 amount, string purpose);
    event ProposalCreated(uint256 indexed proposalId, address indexed proposer, string description);
    
    constructor() Ownable(msg.sender) {
        fee = 0;
        treasury = address(0);
    }
    
    /**
     * @dev Set fee with event emission
     */
    function setFee(uint256 newFee) external onlyOwner {
        uint256 oldFee = fee;
        fee = newFee;
        emit ConfigurationChanged("fee", oldFee, newFee);
    }
    
    /**
     * @dev Set treasury with event emission
     */
    function setTreasury(address newTreasury) external onlyOwner {
        address oldTreasury = treasury;
        treasury = newTreasury;
        emit TreasuryChanged(oldTreasury, newTreasury);
    }
    
    /**
     * @dev Pause contract with event emission
     */
    function pause() external onlyOwner {
        _pause();
        emit ContractPaused(msg.sender);
    }
    
    /**
     * @dev Unpause contract with event emission
     */
    function unpause() external onlyOwner {
        _unpause();
        emit ContractUnpaused(msg.sender);
    }
    
    /**
     * @dev Emergency withdrawal with event emission
     */
    function emergencyWithdraw(address to, uint256 amount) external onlyOwner {
        emit EmergencyActionTaken("emergencyWithdraw", msg.sender);
        // In a real contract, this would perform the actual withdrawal
        // For testing purposes, we just emit the event
    }
    
    /**
     * @dev Receive funds with event emission
     */
    function receiveFunds(string memory purpose) external payable {
        emit FundsReceived(msg.sender, msg.value, purpose);
    }
    
    /**
     * @dev Create proposal with event emission
     */
    function createProposal(uint256 proposalId, string memory description) external {
        emit ProposalCreated(proposalId, msg.sender, description);
    }
    
    /**
     * @dev Batch configuration changes with multiple events
     */
    function batchConfigurationChange(
        uint256 newFee,
        address newTreasury
    ) external onlyOwner {
        if (newFee != fee) {
            uint256 oldFee = fee;
            fee = newFee;
            emit ConfigurationChanged("fee", oldFee, newFee);
        }
        
        if (newTreasury != treasury) {
            address oldTreasury = treasury;
            treasury = newTreasury;
            emit TreasuryChanged(oldTreasury, newTreasury);
        }
    }
    
    /**
     * @dev Fallback function to receive ETH with event
     */
    receive() external payable {
        emit FundsReceived(msg.sender, msg.value, "direct_transfer");
    }
}