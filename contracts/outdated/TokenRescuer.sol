// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract TokenRescuer is ReentrancyGuard {
    address public immutable owner;
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    constructor() {
        owner = msg.sender;
    }
    
    // Call any function on any contract as the rescuer
    function rescue(
        address target,
        bytes calldata data
    ) external onlyOwner nonReentrant returns (bool success, bytes memory result) {
        (success, result) = target.call(data);
        require(success, "Rescue call failed");
    }
    
    // Direct token transfer for contracts that expose transfer functions
    function rescueTokens(
        address token,
        address from, 
        address to,
        uint256 amount
    ) external onlyOwner nonReentrant {
        // This only works if 'from' contract has a function to transfer tokens
        // and we are authorized to call it
        (bool success, ) = from.call(
            abi.encodeWithSignature(
                "transferToken(address,address,uint256)", 
                token, 
                to, 
                amount
            )
        );
        
        if (!success) {
            // Try alternative function names
            (success, ) = from.call(
                abi.encodeWithSignature(
                    "withdrawToken(address,address,uint256)",
                    token,
                    to, 
                    amount
                )
            );
        }
        
        if (!success) {
            // Try another alternative
            (success, ) = from.call(
                abi.encodeWithSignature(
                    "emergencyWithdraw(address,uint256)",
                    token,
                    amount
                )
            );
            
            if (success) {
                // If emergencyWithdraw worked, now transfer from this contract
                IERC20(token).transfer(to, amount);
            }
        }
        
        require(success, "No working withdrawal function found");
    }
    
    // Receive tokens
    receive() external payable {}
}