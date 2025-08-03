// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../interfaces/IRoleManager.sol";

/**
 * @title RoleManaged
 * @dev Base contract for contracts that use centralized role management
 */
abstract contract RoleManaged {
    IRoleManager public immutable roleManager;
    
    // Events
    event RoleManagerSet(address indexed roleManagerAddress);
    
    // Errors
    error UnauthorizedRole(bytes32 role, address account);
    error RoleManagerNotSet();
    
    constructor(address _roleManager) {
        require(_roleManager != address(0), "Invalid RoleManager address");
        roleManager = IRoleManager(_roleManager);
        emit RoleManagerSet(_roleManager);
    }
    
    /**
     * @dev Modifier to check if an account has a specific role
     * @param role The role to check
     */
    modifier onlyRole(bytes32 role) {
        if (!roleManager.checkRole(role, msg.sender)) {
            revert UnauthorizedRole(role, msg.sender);
        }
        _;
    }
    
    /**
     * @dev Check if an account has a specific role
     * @param role The role to check
     * @param account The account to check
     * @return bool Whether the account has the role
     */
    function hasRole(bytes32 role, address account) public view returns (bool) {
        return roleManager.hasRole(role, account);
    }
    
    /**
     * @dev Get the address of the RoleManager contract
     * @return address The RoleManager contract address
     */
    function getRoleManager() public view returns (address) {
        return address(roleManager);
    }
}