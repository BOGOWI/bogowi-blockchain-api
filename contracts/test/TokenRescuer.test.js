const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("TokenRescuer", function () {
  let tokenRescuer;
  let mockToken;
  let mockContract;
  let owner, user1, user2;

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();
    
    // Deploy TokenRescuer
    const TokenRescuer = await ethers.getContractFactory("TokenRescuer");
    tokenRescuer = await TokenRescuer.deploy();
    
    // Deploy mock ERC20 token for testing
    const MockERC20 = await ethers.getContractFactory("contracts/mocks/MockERC20.sol:MockERC20");
    mockToken = await MockERC20.deploy("MockToken", "MTK", ethers.parseEther("1000000"));
    
    // Deploy mock contract that can hold tokens
    const MockContract = await ethers.getContractFactory("MockContractWithTokens");
    mockContract = await MockContract.deploy();
  });

  describe("Deployment", function () {
    it("Should set the correct owner", async function () {
      expect(await tokenRescuer.owner()).to.equal(owner.address);
    });
    
    it("Should be immutable owner", async function () {
      // Owner should be immutable (no setter function)
      expect(tokenRescuer.setOwner).to.be.undefined;
    });
  });

  describe("Access Control", function () {
    it("Should only allow owner to call rescue functions", async function () {
      const data = "0x";
      
      await expect(
        tokenRescuer.connect(user1).rescue(mockContract.target, data)
      ).to.be.revertedWith("Not owner");
      
      await expect(
        tokenRescuer.connect(user1).rescueTokens(
          mockToken.target,
          mockContract.target,
          user1.address,
          100
        )
      ).to.be.revertedWith("Not owner");
    });
  });

  describe("Generic Rescue Function", function () {
    it("Should execute arbitrary calls successfully", async function () {
      // Give the tokenRescuer some tokens to transfer
      await mockToken.transfer(tokenRescuer.target, ethers.parseEther("1000"));
      
      // Encode a function call to transfer tokens from the rescuer
      const transferData = mockToken.interface.encodeFunctionData(
        "transfer",
        [user1.address, ethers.parseEther("100")]
      );
      
      // Execute the rescue call
      await tokenRescuer.rescue(mockToken.target, transferData);
      
      expect(await mockToken.balanceOf(user1.address)).to.equal(ethers.parseEther("100"));
    });
    
    it("Should revert on failed calls", async function () {
      // Try to call a non-existent function
      const invalidData = "0x12345678";
      
      await expect(
        tokenRescuer.rescue(mockContract.target, invalidData)
      ).to.be.revertedWith("Rescue call failed");
    });
    
    it("Should handle calls with return data", async function () {
      // Call a view function that returns data
      const balanceOfData = mockToken.interface.encodeFunctionData(
        "balanceOf",
        [owner.address]
      );
      
      await tokenRescuer.rescue(mockToken.target, balanceOfData);
      
      // The rescue function will succeed if the call succeeds
      // We can't directly access return data in this test setup
    });
  });

  describe("Token Rescue Function", function () {
    beforeEach(async function () {
      // Give the mock contract some tokens
      await mockToken.transfer(mockContract.target, ethers.parseEther("1000"));
      
      // Set up the mock contract to respond to rescue calls
      await mockContract.setTokenAddress(mockToken.target);
    });
    
    it("Should rescue tokens using transferToken function", async function () {
      const rescueAmount = ethers.parseEther("100");
      
      await tokenRescuer.rescueTokens(
        mockToken.target,
        mockContract.target,
        user1.address,
        rescueAmount
      );
      
      expect(await mockToken.balanceOf(user1.address)).to.equal(rescueAmount);
    });
    
    it("Should try alternative function names on failure", async function () {
      // Configure mock to only respond to withdrawToken
      await mockContract.setRespondToTransferToken(false);
      await mockContract.setRespondToWithdrawToken(true);
      
      const rescueAmount = ethers.parseEther("100");
      
      await tokenRescuer.rescueTokens(
        mockToken.target,
        mockContract.target,
        user1.address,
        rescueAmount
      );
      
      expect(await mockToken.balanceOf(user1.address)).to.equal(rescueAmount);
    });
    
    it("Should handle emergencyWithdraw pattern", async function () {
      // Configure mock to only respond to emergencyWithdraw
      await mockContract.setRespondToTransferToken(false);
      await mockContract.setRespondToWithdrawToken(false);
      await mockContract.setRespondToEmergencyWithdraw(true);
      
      const rescueAmount = ethers.parseEther("100");
      
      // First transfer tokens to the rescuer contract
      await mockToken.transfer(tokenRescuer.target, rescueAmount);
      
      await tokenRescuer.rescueTokens(
        mockToken.target,
        mockContract.target,
        user1.address,
        rescueAmount
      );
      
      expect(await mockToken.balanceOf(user1.address)).to.equal(rescueAmount);
    });
    
    it("Should revert when no working function is found", async function () {
      // Configure mock to not respond to any function
      await mockContract.setRespondToTransferToken(false);
      await mockContract.setRespondToWithdrawToken(false);
      await mockContract.setRespondToEmergencyWithdraw(false);
      
      await expect(
        tokenRescuer.rescueTokens(
          mockToken.target,
          mockContract.target,
          user1.address,
          ethers.parseEther("100")
        )
      ).to.be.revertedWith("No working withdrawal function found");
    });
  });

  describe("Reentrancy Protection", function () {
    it("Should prevent reentrancy in rescue function", async function () {
      // Deploy a malicious contract that tries to reenter
      const MaliciousContract = await ethers.getContractFactory("MaliciousReentrant");
      const maliciousContract = await MaliciousContract.deploy(tokenRescuer.target);
      
      const data = "0x";
      
      await expect(
        tokenRescuer.rescue(maliciousContract.target, data)
      ).to.be.revertedWith("Rescue call failed");
    });
    
    it("Should prevent reentrancy in rescueTokens function", async function () {
      const MaliciousContract = await ethers.getContractFactory("MaliciousReentrant");
      const maliciousContract = await MaliciousContract.deploy(tokenRescuer.target);
      
      await expect(
        tokenRescuer.rescueTokens(
          mockToken.target,
          maliciousContract.target,
          user1.address,
          100
        )
      ).to.be.revertedWith("No working withdrawal function found");
    });
  });

  describe("ETH Handling", function () {
    it("Should receive ETH", async function () {
      const amount = ethers.parseEther("1");
      
      await expect(
        owner.sendTransaction({
          to: tokenRescuer.target,
          value: amount
        })
      ).to.not.be.reverted;
      
      expect(await ethers.provider.getBalance(tokenRescuer.target)).to.equal(amount);
    });
    
    it("Should rescue ETH using generic rescue function", async function () {
      const amount = ethers.parseEther("1");
      
      // Send ETH to the rescuer
      await owner.sendTransaction({
        to: tokenRescuer.target,
        value: amount
      });
      
      const initialBalance = await ethers.provider.getBalance(user1.address);
      
      // Rescue ETH by calling a transfer (without sending value to non-payable function)
      const transferData = "0x";
      await tokenRescuer.rescue(user1.address, transferData);
      
      // Note: This is a simplified test. In practice, you'd need a more sophisticated setup
      // to test ETH rescue scenarios
    });
  });

  describe("Edge Cases", function () {
    it("Should handle zero amount rescues", async function () {
      await tokenRescuer.rescueTokens(
        mockToken.target,
        mockContract.target,
        user1.address,
        0
      );
      
      expect(await mockToken.balanceOf(user1.address)).to.equal(0);
    });
    
    it("Should handle rescue to zero address", async function () {
      // This should be handled by the target contract's validation
      await expect(
        tokenRescuer.rescueTokens(
          mockToken.target,
          mockContract.target,
          ethers.ZeroAddress,
          100
        )
      ).to.be.reverted; // Will revert from the mock contract
    });
    
    it("Should handle invalid token addresses", async function () {
      await expect(
        tokenRescuer.rescueTokens(
          ethers.ZeroAddress,
          mockContract.target,
          user1.address,
          100
        )
      ).to.be.reverted;
    });
    
    it("Should handle invalid contract addresses", async function () {
      // Calling functions on zero address might succeed but do nothing
      // This is expected behavior in Solidity - calls to zero address succeed
      await tokenRescuer.rescueTokens(
        mockToken.target,
        ethers.ZeroAddress,
        user1.address,
        100
      );
      
      // Verify no tokens were actually transferred
      expect(await mockToken.balanceOf(user1.address)).to.equal(0);
    });
  });

  describe("Complex Rescue Scenarios", function () {
    it("Should rescue multiple token types", async function () {
      // Deploy another token
      const MockERC20_2 = await ethers.getContractFactory("contracts/mocks/MockERC20.sol:MockERC20");
      const mockToken2 = await MockERC20_2.deploy("MockToken2", "MTK2", ethers.parseEther("1000000"));
      
      // Send tokens to mock contract
      await mockToken.transfer(mockContract.target, ethers.parseEther("100"));
      await mockToken2.transfer(mockContract.target, ethers.parseEther("200"));
      
      // Set up mock contract for both tokens
      await mockContract.setTokenAddress(mockToken.target);
      
      // Rescue first token
      await tokenRescuer.rescueTokens(
        mockToken.target,
        mockContract.target,
        user1.address,
        ethers.parseEther("100")
      );
      
      // Update mock contract for second token
      await mockContract.setTokenAddress(mockToken2.target);
      
      // Rescue second token
      await tokenRescuer.rescueTokens(
        mockToken2.target,
        mockContract.target,
        user1.address,
        ethers.parseEther("200")
      );
      
      expect(await mockToken.balanceOf(user1.address)).to.equal(ethers.parseEther("100"));
      expect(await mockToken2.balanceOf(user1.address)).to.equal(ethers.parseEther("200"));
    });
    
    it("Should handle partial rescues", async function () {
      const totalAmount = ethers.parseEther("1000");
      const rescueAmount = ethers.parseEther("300");
      
      await mockToken.transfer(mockContract.target, totalAmount);
      await mockContract.setTokenAddress(mockToken.target);
      
      await tokenRescuer.rescueTokens(
        mockToken.target,
        mockContract.target,
        user1.address,
        rescueAmount
      );
      
      expect(await mockToken.balanceOf(user1.address)).to.equal(rescueAmount);
      expect(await mockToken.balanceOf(mockContract.target)).to.equal(totalAmount - rescueAmount);
    });
  });

  describe("Gas Optimization", function () {
    it("Should be gas efficient for simple rescues", async function () {
      await mockToken.transfer(mockContract.target, ethers.parseEther("100"));
      await mockContract.setTokenAddress(mockToken.target);
      
      const tx = await tokenRescuer.rescueTokens(
        mockToken.target,
        mockContract.target,
        user1.address,
        ethers.parseEther("100")
      );
      
      const receipt = await tx.wait();
      
      // Gas usage should be reasonable (adjust threshold as needed)
      expect(receipt.gasUsed).to.be.lessThan(200000);
    });
  });
});

// Mock contracts for testing

// Mock ERC20 token
const MockERC20Source = `
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MockERC20 is ERC20 {
    constructor(string memory name, string memory symbol, uint256 totalSupply) ERC20(name, symbol) {
        _mint(msg.sender, totalSupply);
    }
}
`;

// Mock contract that can hold and transfer tokens
const MockContractWithTokensSource = `
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
}
`;

// Malicious contract for reentrancy testing
const MaliciousReentrantSource = `
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface ITokenRescuer {
    function rescue(address target, bytes calldata data) external returns (bool success, bytes memory result);
}

contract MaliciousReentrant {
    ITokenRescuer public rescuer;
    bool public hasReentered = false;
    
    constructor(address _rescuer) {
        rescuer = ITokenRescuer(_rescuer);
    }
    
    fallback() external {
        if (!hasReentered) {
            hasReentered = true;
            rescuer.rescue(address(this), "");
        }
    }
    
    function transferToken(address, address, uint256) external {
        if (!hasReentered) {
            hasReentered = true;
            rescuer.rescue(address(this), "");
        }
    }
}
`;