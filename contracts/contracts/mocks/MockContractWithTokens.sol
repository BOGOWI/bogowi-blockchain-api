// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MockContractWithTokens {
    address public tokenAddress;
    bool public respondToTransferToken = true;
    bool public respondToWithdrawToken = true;
    bool public respondToEmergencyWithdraw = true;
    
    function setTokenAddress(address _token) external {
        tokenAddress = _token;
    }
    
    function setRespondToTransferToken(bool _respond) external {
        respondToTransferToken = _respond;
    }
    
    function setRespondToWithdrawToken(bool _respond) external {
        respondToWithdrawToken = _respond;
    }
    
    function setRespondToEmergencyWithdraw(bool _respond) external {
        respondToEmergencyWithdraw = _respond;
    }
    
    function transferToken(address token, address to, uint256 amount) external {
        require(respondToTransferToken, "Not responding to transferToken");
        require(to != address(0), "Invalid recipient");
        IERC20(token).transfer(to, amount);
    }
    
    function withdrawToken(address token, address to, uint256 amount) external {
        require(respondToWithdrawToken, "Not responding to withdrawToken");
        require(to != address(0), "Invalid recipient");
        IERC20(token).transfer(to, amount);
    }
    
    function emergencyWithdraw(address token, uint256 amount) external {
        require(respondToEmergencyWithdraw, "Not responding to emergencyWithdraw");
        IERC20(token).transfer(msg.sender, amount);
    }
    
    // Allow contract to receive tokens
    receive() external payable {}
}