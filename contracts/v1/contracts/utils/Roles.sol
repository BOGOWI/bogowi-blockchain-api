// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title Roles
 * @notice Centralized role definitions for the BOGOWI ecosystem
 * @dev All contracts should import and use these role definitions
 */
library Roles {
    // Core administrative roles
    bytes32 public constant DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    
    // Token roles
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");
    bytes32 public constant DAO_ROLE = keccak256("DAO_ROLE");
    bytes32 public constant BUSINESS_ROLE = keccak256("BUSINESS_ROLE");
    
    // NFT specific roles
    bytes32 public constant NFT_MINTER_ROLE = keccak256("NFT_MINTER_ROLE");
    bytes32 public constant REGISTRY_ADMIN_ROLE = keccak256("REGISTRY_ADMIN_ROLE");
    bytes32 public constant CONTRACT_DEPLOYER_ROLE = keccak256("CONTRACT_DEPLOYER_ROLE");
    
    // Reward system roles
    bytes32 public constant TREASURY_ROLE = keccak256("TREASURY_ROLE");
    bytes32 public constant BACKEND_ROLE = keccak256("BACKEND_ROLE");
    bytes32 public constant REWARD_DISTRIBUTOR_ROLE = keccak256("REWARD_DISTRIBUTOR_ROLE");
    
    // Marketplace roles (future)
    bytes32 public constant MARKETPLACE_ADMIN_ROLE = keccak256("MARKETPLACE_ADMIN_ROLE");
    bytes32 public constant FEE_COLLECTOR_ROLE = keccak256("FEE_COLLECTOR_ROLE");
}