// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title MigrationHelper
 * @dev Assists with migrating data between contract versions
 * @notice Use this for migrating user data when replacing peripheral contracts
 */
contract MigrationHelper is AccessControl, Pausable {
    bytes32 public constant MIGRATION_ROLE = keccak256("MIGRATION_ROLE");
    
    // Migration tracking
    mapping(address => mapping(address => bool)) public migrated; // oldContract => user => migrated
    mapping(address => uint256) public migrationCount;
    
    // Events
    event MigrationStarted(address indexed oldContract, address indexed newContract);
    event UserMigrated(address indexed user, address indexed oldContract, address indexed newContract);
    event BatchMigrationCompleted(address indexed oldContract, uint256 count);
    event TokensRecovered(address indexed token, address indexed to, uint256 amount);
    
    // Custom errors
    error AlreadyMigrated();
    error MigrationFailed();
    error InvalidContracts();
    error BatchSizeTooLarge();
    
    uint256 public constant MAX_BATCH_SIZE = 100;
    
    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MIGRATION_ROLE, msg.sender);
    }
    
    /**
     * @dev Check if a user has been migrated
     * @param oldContract The old contract address
     * @param user The user address
     * @return True if the user has been migrated
     */
    function isMigrated(address oldContract, address user) external view returns (bool) {
        return migrated[oldContract][user];
    }
    
    /**
     * @dev Mark user as migrated
     * @param oldContract The old contract address
     * @param user The user address
     */
    function markMigrated(address oldContract, address user) 
        external 
        onlyRole(MIGRATION_ROLE) 
        whenNotPaused 
    {
        if (migrated[oldContract][user]) revert AlreadyMigrated();
        
        migrated[oldContract][user] = true;
        migrationCount[oldContract]++;
        
        emit UserMigrated(user, oldContract, msg.sender);
    }
    
    /**
     * @dev Batch mark users as migrated
     * @param oldContract The old contract address
     * @param users Array of user addresses
     */
    function batchMarkMigrated(address oldContract, address[] calldata users) 
        external 
        onlyRole(MIGRATION_ROLE) 
        whenNotPaused 
    {
        if (users.length > MAX_BATCH_SIZE) revert BatchSizeTooLarge();
        
        uint256 count = 0;
        for (uint256 i = 0; i < users.length; i++) {
            if (!migrated[oldContract][users[i]]) {
                migrated[oldContract][users[i]] = true;
                count++;
                emit UserMigrated(users[i], oldContract, msg.sender);
            }
        }
        
        migrationCount[oldContract] += count;
        emit BatchMigrationCompleted(oldContract, count);
    }
    
    /**
     * @dev Emergency token recovery
     * @param token The token address (0 for ETH)
     * @param to The recipient address
     * @param amount The amount to recover
     */
    function recoverTokens(address token, address to, uint256 amount) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(to != address(0), "Invalid recipient");
        require(amount > 0, "Invalid amount");
        
        if (token == address(0)) {
            (bool success, ) = to.call{value: amount}("");
            require(success, "ETH transfer failed");
        } else {
            IERC20(token).transfer(to, amount);
        }
        
        emit TokensRecovered(token, to, amount);
    }
    
    /**
     * @dev Pause migrations
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }
    
    /**
     * @dev Unpause migrations
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }
    
    /**
     * @dev Receive ETH
     */
    receive() external payable {}
}

/**
 * @title IMigratable
 * @dev Interface for contracts that support migration
 */
interface IMigratable {
    function migrate(address newContract) external;
    function acceptMigration(address user, bytes calldata data) external;
    function exportUserData(address user) external view returns (bytes memory);
}

/**
 * @title MigratableContract
 * @dev Base contract for contracts that support migration
 */
abstract contract MigratableContract is IMigratable, Pausable, AccessControl {
    address public newVersion;
    bool public deprecated;
    
    event Deprecated(address indexed newVersion);
    event UserDataExported(address indexed user);
    event UserDataImported(address indexed user);
    
    /**
     * @dev Deprecate this contract and point to new version
     * @param _newVersion Address of the new contract version
     */
    function deprecate(address _newVersion) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(_newVersion != address(0), "Invalid new version");
        require(!deprecated, "Already deprecated");
        
        deprecated = true;
        newVersion = _newVersion;
        _pause(); // Pause operations
        
        emit Deprecated(_newVersion);
    }
    
    /**
     * @dev Export user data for migration
     * @param user The user address
     * @return Encoded user data
     */
    function exportUserData(address user) external view virtual returns (bytes memory);
    
    /**
     * @dev Accept migrated user data
     * @param user The user address
     * @param data The encoded user data
     */
    function acceptMigration(address user, bytes calldata data) external virtual;
    
    /**
     * @dev Migrate specific user to new contract
     * @param newContract The new contract address
     */
    function migrate(address newContract) external virtual {
        require(deprecated, "Contract not deprecated");
        require(newContract == newVersion, "Invalid new contract");
        
        bytes memory userData = this.exportUserData(msg.sender);
        IMigratable(newContract).acceptMigration(msg.sender, userData);
        
        emit UserDataExported(msg.sender);
    }
}