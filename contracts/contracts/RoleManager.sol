// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

/**
 * @title RoleManager
 * @dev Centralized role management contract for the BOGOWI ecosystem
 * Provides a single source of truth for all role assignments across multiple contracts
 */
contract RoleManager is AccessControl, Pausable {
    // Global roles used across the ecosystem
    bytes32 public constant DAO_ROLE = keccak256("DAO_ROLE");
    bytes32 public constant BUSINESS_ROLE = keccak256("BUSINESS_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    bytes32 public constant TREASURY_ROLE = keccak256("TREASURY_ROLE");
    bytes32 public constant DISTRIBUTOR_BACKEND_ROLE = keccak256("DISTRIBUTOR_BACKEND_ROLE");
    
    // Contract registry
    mapping(address => bool) public registeredContracts;
    mapping(address => string) public contractNames;
    
    // Events
    event ContractRegistered(address indexed contractAddress, string name);
    event ContractDeregistered(address indexed contractAddress);
    event RoleGrantedGlobally(bytes32 indexed role, address indexed account, address indexed sender);
    event RoleRevokedGlobally(bytes32 indexed role, address indexed account, address indexed sender);
    
    modifier onlyRegisteredContract() {
        require(registeredContracts[msg.sender], "Not a registered contract");
        _;
    }
    
    constructor() {
        // Grant DEFAULT_ADMIN_ROLE to deployer
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        
        // Set up role admin relationships
        _setRoleAdmin(DAO_ROLE, DEFAULT_ADMIN_ROLE);
        _setRoleAdmin(BUSINESS_ROLE, DEFAULT_ADMIN_ROLE);
        _setRoleAdmin(MINTER_ROLE, DEFAULT_ADMIN_ROLE);
        _setRoleAdmin(PAUSER_ROLE, DEFAULT_ADMIN_ROLE);
        _setRoleAdmin(TREASURY_ROLE, DEFAULT_ADMIN_ROLE);
        _setRoleAdmin(DISTRIBUTOR_BACKEND_ROLE, DEFAULT_ADMIN_ROLE);
    }
    
    /**
     * @dev Register a contract to use this RoleManager
     * @param contractAddress The address of the contract to register
     * @param name A human-readable name for the contract
     */
    function registerContract(address contractAddress, string memory name) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(contractAddress != address(0), "Invalid contract address");
        require(!registeredContracts[contractAddress], "Contract already registered");
        
        registeredContracts[contractAddress] = true;
        contractNames[contractAddress] = name;
        
        emit ContractRegistered(contractAddress, name);
    }
    
    /**
     * @dev Deregister a contract
     * @param contractAddress The address of the contract to deregister
     */
    function deregisterContract(address contractAddress) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(registeredContracts[contractAddress], "Contract not registered");
        
        registeredContracts[contractAddress] = false;
        delete contractNames[contractAddress];
        
        emit ContractDeregistered(contractAddress);
    }
    
    /**
     * @dev Check if an account has a specific role
     * This can be called by registered contracts to verify permissions
     * @param role The role to check
     * @param account The account to check
     * @return bool Whether the account has the role
     */
    function checkRole(bytes32 role, address account) 
        external 
        view 
        onlyRegisteredContract 
        returns (bool) 
    {
        return hasRole(role, account);
    }
    
    /**
     * @dev Grant a role to an account (with additional logging)
     * @param role The role to grant
     * @param account The account to grant the role to
     */
    function grantRole(bytes32 role, address account) 
        public 
        override(AccessControl) 
        onlyRole(getRoleAdmin(role)) 
    {
        super.grantRole(role, account);
        emit RoleGrantedGlobally(role, account, msg.sender);
    }
    
    /**
     * @dev Revoke a role from an account (with additional logging)
     * @param role The role to revoke
     * @param account The account to revoke the role from
     */
    function revokeRole(bytes32 role, address account) 
        public 
        override(AccessControl) 
        onlyRole(getRoleAdmin(role)) 
    {
        super.revokeRole(role, account);
        emit RoleRevokedGlobally(role, account, msg.sender);
    }
    
    /**
     * @dev Get all registered contracts
     * @return addresses Array of registered contract addresses
     * @return names Array of contract names
     */
    function getRegisteredContracts() 
        external 
        pure
        returns (address[] memory addresses, string[] memory names) 
    {
        // This is a placeholder implementation
        // In production, consider maintaining a separate array of registered addresses
        // for efficient enumeration
        
        // For now, return empty arrays
        // TODO: Implement proper enumeration if needed
        addresses = new address[](0);
        names = new string[](0);
    }
    
    /**
     * @dev Batch grant roles
     * @param role The role to grant
     * @param accounts Array of accounts to grant the role to
     */
    function batchGrantRole(bytes32 role, address[] calldata accounts) 
        external 
        onlyRole(getRoleAdmin(role)) 
    {
        for (uint256 i = 0; i < accounts.length; i++) {
            grantRole(role, accounts[i]);
        }
    }
    
    /**
     * @dev Batch revoke roles
     * @param role The role to revoke
     * @param accounts Array of accounts to revoke the role from
     */
    function batchRevokeRole(bytes32 role, address[] calldata accounts) 
        external 
        onlyRole(getRoleAdmin(role)) 
    {
        for (uint256 i = 0; i < accounts.length; i++) {
            revokeRole(role, accounts[i]);
        }
    }
    
    /**
     * @dev Transfer admin role to a new address (careful operation)
     * @param newAdmin The address to transfer admin role to
     */
    function transferAdmin(address newAdmin) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(newAdmin != address(0), "Invalid new admin");
        grantRole(DEFAULT_ADMIN_ROLE, newAdmin);
        renounceRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }
    
    /**
     * @dev Pause the contract (prevents role checking)
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }
    
    /**
     * @dev Unpause the contract
     */
    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }
    
    /**
     * @dev Override supportsInterface to include AccessControl
     */
    function supportsInterface(bytes4 interfaceId) 
        public 
        view 
        virtual 
        override(AccessControl) 
        returns (bool) 
    {
        return super.supportsInterface(interfaceId);
    }
}