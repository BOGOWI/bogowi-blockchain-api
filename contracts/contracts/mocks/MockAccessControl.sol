// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";

contract MockAccessControl is AccessControl {
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant USER_ROLE = keccak256("USER_ROLE");
    bytes32 public constant MANAGER_ROLE = keccak256("MANAGER_ROLE");
    
    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
    }
    
    function adminFunction() external onlyRole(ADMIN_ROLE) {
        // Admin only function
    }
    
    function userFunction() external onlyRole(USER_ROLE) {
        // User only function
    }
    
    function managerFunction() external onlyRole(MANAGER_ROLE) {
        // Manager only function
    }
    
    function publicFunction() external pure returns (string memory) {
        return "public";
    }
}