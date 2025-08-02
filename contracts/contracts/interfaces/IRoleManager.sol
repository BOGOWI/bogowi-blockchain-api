// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IRoleManager {
    function checkRole(bytes32 role, address account) external view returns (bool);
    function hasRole(bytes32 role, address account) external view returns (bool);
    
    function DAO_ROLE() external view returns (bytes32);
    function BUSINESS_ROLE() external view returns (bytes32);
    function MINTER_ROLE() external view returns (bytes32);
    function PAUSER_ROLE() external view returns (bytes32);
    function TREASURY_ROLE() external view returns (bytes32);
    function DISTRIBUTOR_BACKEND_ROLE() external view returns (bytes32);
    function DEFAULT_ADMIN_ROLE() external view returns (bytes32);
}