// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

contract MockPausable is AccessControl, Pausable {
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    
    string public name;
    
    constructor(string memory _name) {
        name = _name;
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
    }
    
    function pause() public {
        require(
            hasRole(DEFAULT_ADMIN_ROLE, msg.sender) || hasRole(PAUSER_ROLE, msg.sender),
            "Unauthorized"
        );
        _pause();
    }
    
    function unpause() public {
        require(
            hasRole(DEFAULT_ADMIN_ROLE, msg.sender) || hasRole(PAUSER_ROLE, msg.sender),
            "Unauthorized"
        );
        _unpause();
    }
}