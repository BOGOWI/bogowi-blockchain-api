// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "../base/RoleManaged.sol";

contract TestRoleManaged is RoleManaged {
    uint256 public value;
    
    constructor(address _roleManager) RoleManaged(_roleManager) {}
    
    function setValueDAO(uint256 _value) external onlyRole(keccak256("DAO_ROLE")) {
        value = _value;
    }
    
    function setValueBusiness(uint256 _value) external onlyRole(keccak256("BUSINESS_ROLE")) {
        value = _value;
    }
    
    function checkUserRole(bytes32 role, address account) external view returns (bool) {
        return hasRole(role, account);
    }
    
    function getRoleManagerAddress() external view returns (address) {
        return getRoleManager();
    }
}