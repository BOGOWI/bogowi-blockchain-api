// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/extensions/AccessControlEnumerable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "./StandardErrors.sol";

/**
 * @title IPausableContract
 * @notice Interface for contracts that can be paused
 */
interface IPausableContract {
    function pause() external;
    function unpause() external;
    function paused() external view returns (bool);
}

/**
 * @title EmergencyPauseController
 * @author BOGOWI Team
 * @notice Centralized emergency pause controller for all BOGOWI contracts
 * @dev Allows authorized guardians to pause/unpause multiple contracts simultaneously
 * Features:
 * - Multi-signature guardian system
 * - Batch pause/unpause operations
 * - Individual contract management
 * - Emergency pause with time limit
 * - Audit trail of all pause actions
 * @custom:security-contact security@bogowi.com
 */
contract EmergencyPauseController is AccessControlEnumerable, StandardErrors {
    /// @notice Role for emergency pause guardians
    bytes32 public constant GUARDIAN_ROLE = keccak256("GUARDIAN_ROLE");
    
    /// @notice Role for contract managers who can add/remove contracts
    bytes32 public constant MANAGER_ROLE = keccak256("MANAGER_ROLE");
    
    /// @notice Maximum duration for emergency pause (72 hours)
    uint256 public constant MAX_PAUSE_DURATION = 72 hours;
    
    /// @notice Minimum guardians required for consensus
    uint256 public constant MIN_GUARDIANS = 3;
    
    /// @notice Required confirmations for emergency pause
    uint256 public requiredConfirmations = 2;
    
    /// @notice Tracked pausable contracts
    address[] public pausableContracts;
    mapping(address => bool) public isTrackedContract;
    mapping(address => string) public contractNames;
    
    /// @notice Emergency pause proposals
    struct PauseProposal {
        address proposer;
        string reason;
        uint256 timestamp;
        uint256 confirmations;
        bool executed;
        bool isPause; // true for pause, false for unpause
        address[] targetContracts;
    }
    
    /// @notice Active proposals
    mapping(uint256 => PauseProposal) public proposals;
    mapping(uint256 => mapping(address => bool)) public hasConfirmed;
    uint256 public proposalCount;
    
    /// @notice Pause history for audit trail
    struct PauseEvent {
        address guardian;
        address[] contracts;
        bool isPause;
        string reason;
        uint256 timestamp;
    }
    
    PauseEvent[] public pauseHistory;
    
    /// @notice Emergency pause expiry times
    mapping(address => uint256) public pauseExpiry;
    
    // Events
    event ContractAdded(address indexed contractAddress, string name);
    event ContractRemoved(address indexed contractAddress);
    event PauseProposalCreated(uint256 indexed proposalId, address indexed proposer, bool isPause);
    event ProposalConfirmed(uint256 indexed proposalId, address indexed guardian);
    event EmergencyPauseExecuted(address indexed executor, address[] contracts, string reason);
    event EmergencyUnpauseExecuted(address indexed executor, address[] contracts, string reason);
    event RequiredConfirmationsUpdated(uint256 oldValue, uint256 newValue);
    event PauseExpired(address indexed contractAddress);
    
    /**
     * @notice Initializes the emergency pause controller
     * @dev Sets up initial roles and guardians
     * @param _guardians Initial list of guardian addresses
     * @param _manager Address that can manage tracked contracts
     */
    constructor(address[] memory _guardians, address _manager) {
        require(_guardians.length >= MIN_GUARDIANS, INVALID_PARAMETER);
        require(_manager != address(0), ZERO_ADDRESS);
        
        // Grant roles
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MANAGER_ROLE, _manager);
        
        // Add initial guardians
        for (uint256 i = 0; i < _guardians.length; i++) {
            require(_guardians[i] != address(0), ZERO_ADDRESS);
            _grantRole(GUARDIAN_ROLE, _guardians[i]);
        }
    }
    
    /**
     * @notice Adds a contract to be monitored
     * @dev Only MANAGER_ROLE can add contracts
     * @param contractAddress Address of the pausable contract
     * @param name Human-readable name for the contract
     */
    function addContract(address contractAddress, string memory name) 
        external 
        onlyRole(MANAGER_ROLE) 
    {
        require(contractAddress != address(0), ZERO_ADDRESS);
        require(!isTrackedContract[contractAddress], ALREADY_EXISTS);
        require(bytes(name).length > 0, EMPTY_STRING);
        
        // Verify contract implements pause interface
        try IPausableContract(contractAddress).paused() returns (bool) {
            pausableContracts.push(contractAddress);
            isTrackedContract[contractAddress] = true;
            contractNames[contractAddress] = name;
            
            emit ContractAdded(contractAddress, name);
        } catch {
            revert(INVALID_ADDRESS);
        }
    }
    
    /**
     * @notice Removes a contract from monitoring
     * @dev Only MANAGER_ROLE can remove contracts
     * @param contractAddress Address of the contract to remove
     */
    function removeContract(address contractAddress) 
        external 
        onlyRole(MANAGER_ROLE) 
    {
        require(isTrackedContract[contractAddress], DOES_NOT_EXIST);
        
        isTrackedContract[contractAddress] = false;
        
        // Remove from array
        for (uint256 i = 0; i < pausableContracts.length; i++) {
            if (pausableContracts[i] == contractAddress) {
                pausableContracts[i] = pausableContracts[pausableContracts.length - 1];
                pausableContracts.pop();
                break;
            }
        }
        
        emit ContractRemoved(contractAddress);
    }
    
    /**
     * @notice Proposes an emergency pause for specific contracts
     * @dev Requires GUARDIAN_ROLE
     * @param contractAddresses Addresses of contracts to pause
     * @param reason Reason for the emergency pause
     */
    function proposePause(
        address[] memory contractAddresses,
        string memory reason
    ) external onlyRole(GUARDIAN_ROLE) returns (uint256) {
        require(contractAddresses.length > 0, INVALID_LENGTH);
        require(bytes(reason).length > 0, EMPTY_STRING);
        
        // Verify all contracts are tracked
        for (uint256 i = 0; i < contractAddresses.length; i++) {
            require(isTrackedContract[contractAddresses[i]], DOES_NOT_EXIST);
        }
        
        uint256 proposalId = proposalCount++;
        
        proposals[proposalId] = PauseProposal({
            proposer: msg.sender,
            reason: reason,
            timestamp: block.timestamp,
            confirmations: 1,
            executed: false,
            isPause: true,
            targetContracts: contractAddresses
        });
        
        hasConfirmed[proposalId][msg.sender] = true;
        
        emit PauseProposalCreated(proposalId, msg.sender, true);
        
        // Auto-execute if only one confirmation required
        if (requiredConfirmations == 1) {
            _executePause(proposalId);
        }
        
        return proposalId;
    }
    
    /**
     * @notice Proposes to unpause specific contracts
     * @dev Requires GUARDIAN_ROLE
     * @param contractAddresses Addresses of contracts to unpause
     * @param reason Reason for unpausing
     */
    function proposeUnpause(
        address[] memory contractAddresses,
        string memory reason
    ) external onlyRole(GUARDIAN_ROLE) returns (uint256) {
        require(contractAddresses.length > 0, INVALID_LENGTH);
        require(bytes(reason).length > 0, EMPTY_STRING);
        
        uint256 proposalId = proposalCount++;
        
        proposals[proposalId] = PauseProposal({
            proposer: msg.sender,
            reason: reason,
            timestamp: block.timestamp,
            confirmations: 1,
            executed: false,
            isPause: false,
            targetContracts: contractAddresses
        });
        
        hasConfirmed[proposalId][msg.sender] = true;
        
        emit PauseProposalCreated(proposalId, msg.sender, false);
        
        if (requiredConfirmations == 1) {
            _executeUnpause(proposalId);
        }
        
        return proposalId;
    }
    
    /**
     * @notice Confirms a pause/unpause proposal
     * @dev Auto-executes when threshold is reached
     * @param proposalId ID of the proposal to confirm
     */
    function confirmProposal(uint256 proposalId) 
        external 
        onlyRole(GUARDIAN_ROLE) 
    {
        require(proposalId < proposalCount, DOES_NOT_EXIST);
        require(!proposals[proposalId].executed, ALREADY_PROCESSED);
        require(!hasConfirmed[proposalId][msg.sender], ALREADY_PROCESSED);
        
        PauseProposal storage proposal = proposals[proposalId];
        hasConfirmed[proposalId][msg.sender] = true;
        proposal.confirmations++;
        
        emit ProposalConfirmed(proposalId, msg.sender);
        
        // Execute if threshold reached
        if (proposal.confirmations >= requiredConfirmations) {
            if (proposal.isPause) {
                _executePause(proposalId);
            } else {
                _executeUnpause(proposalId);
            }
        }
    }
    
    /**
     * @notice Emergency pause all tracked contracts
     * @dev Requires multiple guardian confirmations
     * @param reason Reason for emergency pause
     */
    function emergencyPauseAll(string memory reason) 
        external 
        onlyRole(GUARDIAN_ROLE) 
        returns (uint256) 
    {
        require(pausableContracts.length > 0, INVALID_STATE);
        
        // Copy to memory for the function call
        address[] memory allContracts = new address[](pausableContracts.length);
        for (uint256 i = 0; i < pausableContracts.length; i++) {
            allContracts[i] = pausableContracts[i];
        }
        
        return proposePause(allContracts, reason);
    }
    
    /**
     * @notice Executes the pause proposal
     * @dev Internal function called when threshold is met
     */
    function _executePause(uint256 proposalId) private {
        PauseProposal storage proposal = proposals[proposalId];
        require(!proposal.executed, ALREADY_PROCESSED);
        
        proposal.executed = true;
        
        // Execute pause on all target contracts
        for (uint256 i = 0; i < proposal.targetContracts.length; i++) {
            address target = proposal.targetContracts[i];
            
            try IPausableContract(target).pause() {
                // Set expiry time
                pauseExpiry[target] = block.timestamp + MAX_PAUSE_DURATION;
            } catch {
                // Continue with other contracts even if one fails
            }
        }
        
        // Record in history
        pauseHistory.push(PauseEvent({
            guardian: msg.sender,
            contracts: proposal.targetContracts,
            isPause: true,
            reason: proposal.reason,
            timestamp: block.timestamp
        }));
        
        emit EmergencyPauseExecuted(msg.sender, proposal.targetContracts, proposal.reason);
    }
    
    /**
     * @notice Executes the unpause proposal
     * @dev Internal function called when threshold is met
     */
    function _executeUnpause(uint256 proposalId) private {
        PauseProposal storage proposal = proposals[proposalId];
        require(!proposal.executed, ALREADY_PROCESSED);
        
        proposal.executed = true;
        
        // Execute unpause on all target contracts
        for (uint256 i = 0; i < proposal.targetContracts.length; i++) {
            address target = proposal.targetContracts[i];
            
            try IPausableContract(target).unpause() {
                // Clear expiry time
                pauseExpiry[target] = 0;
            } catch {
                // Continue with other contracts
            }
        }
        
        // Record in history
        pauseHistory.push(PauseEvent({
            guardian: msg.sender,
            contracts: proposal.targetContracts,
            isPause: false,
            reason: proposal.reason,
            timestamp: block.timestamp
        }));
        
        emit EmergencyUnpauseExecuted(msg.sender, proposal.targetContracts, proposal.reason);
    }
    
    /**
     * @notice Checks and expires any pauses that exceed MAX_PAUSE_DURATION
     * @dev Can be called by anyone to enforce time limits
     */
    function checkAndExpirePauses() external {
        for (uint256 i = 0; i < pausableContracts.length; i++) {
            address target = pausableContracts[i];
            
            if (pauseExpiry[target] > 0 && block.timestamp > pauseExpiry[target]) {
                try IPausableContract(target).unpause() {
                    pauseExpiry[target] = 0;
                    emit PauseExpired(target);
                } catch {
                    // Contract might have been unpaused manually
                }
            }
        }
    }
    
    /**
     * @notice Updates the required confirmations for proposals
     * @dev Only DEFAULT_ADMIN_ROLE can update
     * @param newRequired New number of required confirmations
     */
    function updateRequiredConfirmations(uint256 newRequired) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(newRequired > 0, ZERO_AMOUNT);
        require(newRequired <= getRoleCount(GUARDIAN_ROLE), EXCEEDS_LIMIT);
        
        uint256 oldRequired = requiredConfirmations;
        requiredConfirmations = newRequired;
        
        emit RequiredConfirmationsUpdated(oldRequired, newRequired);
    }
    
    /**
     * @notice Returns the current pause status of all tracked contracts
     * @return contracts Array of contract addresses
     * @return paused Array of pause statuses
     * @return names Array of contract names
     */
    function getContractStatuses() 
        external 
        view 
        returns (
            address[] memory contracts,
            bool[] memory paused,
            string[] memory names
        ) 
    {
        uint256 length = pausableContracts.length;
        contracts = new address[](length);
        paused = new bool[](length);
        names = new string[](length);
        
        for (uint256 i = 0; i < length; i++) {
            contracts[i] = pausableContracts[i];
            names[i] = contractNames[pausableContracts[i]];
            
            try IPausableContract(pausableContracts[i]).paused() returns (bool isPaused) {
                paused[i] = isPaused;
            } catch {
                paused[i] = false;
            }
        }
    }
    
    /**
     * @notice Returns the number of addresses with a specific role
     * @param role Role to count
     * @return Number of addresses with the role
     */
    function getRoleCount(bytes32 role) public view returns (uint256) {
        return getRoleMemberCount(role);
    }
    
    /**
     * @notice Returns pause history for audit purposes
     * @param limit Maximum number of events to return
     * @return Array of recent pause events
     */
    function getPauseHistory(uint256 limit) 
        external 
        view 
        returns (PauseEvent[] memory) 
    {
        uint256 length = pauseHistory.length;
        if (limit > length) limit = length;
        
        PauseEvent[] memory recent = new PauseEvent[](limit);
        
        for (uint256 i = 0; i < limit; i++) {
            recent[i] = pauseHistory[length - limit + i];
        }
        
        return recent;
    }
}