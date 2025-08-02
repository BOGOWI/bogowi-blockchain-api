// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

/**
 * @title ContractRegistry
 * @dev Registry for managing peripheral contract addresses without upgrading core contracts
 * @notice This allows updating references to non-critical contracts while keeping core contracts immutable
 */
contract ContractRegistry is AccessControl, Pausable {
    bytes32 public constant REGISTRY_ADMIN_ROLE = keccak256("REGISTRY_ADMIN_ROLE");
    
    // Contract name to address mapping
    mapping(string => address) private contracts;
    mapping(string => uint256) private contractVersions;
    mapping(string => address[]) private contractHistory;
    
    // Events
    event ContractRegistered(string indexed name, address indexed contractAddress, uint256 version);
    event ContractUpdated(string indexed name, address indexed oldAddress, address indexed newAddress, uint256 version);
    event ContractDeprecated(string indexed name, address indexed contractAddress);
    
    // Custom errors
    error InvalidAddress();
    error ContractNotFound();
    error ContractAlreadyRegistered();
    error InvalidVersion();
    
    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(REGISTRY_ADMIN_ROLE, msg.sender);
    }
    
    /**
     * @dev Register a new contract
     * @param name The identifier for the contract
     * @param contractAddress The address of the contract
     */
    function registerContract(string memory name, address contractAddress) 
        external 
        onlyRole(REGISTRY_ADMIN_ROLE) 
        whenNotPaused 
    {
        if (contractAddress == address(0)) revert InvalidAddress();
        if (contracts[name] != address(0)) revert ContractAlreadyRegistered();
        
        contracts[name] = contractAddress;
        contractVersions[name] = 1;
        contractHistory[name].push(contractAddress);
        
        emit ContractRegistered(name, contractAddress, 1);
    }
    
    /**
     * @dev Update an existing contract address
     * @param name The identifier for the contract
     * @param newAddress The new address for the contract
     */
    function updateContract(string memory name, address newAddress) 
        external 
        onlyRole(REGISTRY_ADMIN_ROLE) 
        whenNotPaused 
    {
        if (newAddress == address(0)) revert InvalidAddress();
        
        address oldAddress = contracts[name];
        if (oldAddress == address(0)) revert ContractNotFound();
        
        contracts[name] = newAddress;
        contractVersions[name]++;
        contractHistory[name].push(newAddress);
        
        emit ContractUpdated(name, oldAddress, newAddress, contractVersions[name]);
    }
    
    /**
     * @dev Get contract address by name
     * @param name The identifier for the contract
     * @return The current address of the contract
     */
    function getContract(string memory name) external view returns (address) {
        address contractAddress = contracts[name];
        if (contractAddress == address(0)) revert ContractNotFound();
        return contractAddress;
    }
    
    /**
     * @dev Get contract version
     * @param name The identifier for the contract
     * @return The current version number
     */
    function getContractVersion(string memory name) external view returns (uint256) {
        if (contracts[name] == address(0)) revert ContractNotFound();
        return contractVersions[name];
    }
    
    /**
     * @dev Get contract deployment history
     * @param name The identifier for the contract
     * @return Array of all contract addresses in chronological order
     */
    function getContractHistory(string memory name) external view returns (address[] memory) {
        if (contracts[name] == address(0)) revert ContractNotFound();
        return contractHistory[name];
    }
    
    /**
     * @dev Check if a contract is registered
     * @param name The identifier for the contract
     * @return True if the contract is registered
     */
    function isRegistered(string memory name) external view returns (bool) {
        return contracts[name] != address(0);
    }
    
    /**
     * @dev Deprecate a contract (remove from registry)
     * @param name The identifier for the contract
     */
    function deprecateContract(string memory name) 
        external 
        onlyRole(REGISTRY_ADMIN_ROLE) 
    {
        address contractAddress = contracts[name];
        if (contractAddress == address(0)) revert ContractNotFound();
        
        delete contracts[name];
        // Keep history and version for audit trail
        
        emit ContractDeprecated(name, contractAddress);
    }
    
    /**
     * @dev Pause the registry
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }
    
    /**
     * @dev Unpause the registry
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }
}