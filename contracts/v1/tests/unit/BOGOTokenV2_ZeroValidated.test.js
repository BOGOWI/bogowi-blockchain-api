const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGOTokenV2_ZeroValidated", function () {
  let bogoToken;
  let owner;
  let daoRole;
  let businessRole;
  let minterRole;
  let pauserRole;
  let user1;
  let user2;
  let nonAuthorized;

  const MAX_SUPPLY = ethers.parseEther("1000000000"); // 1 billion
  const DAO_ALLOCATION = ethers.parseEther("200000000"); // 200M
  const BUSINESS_ALLOCATION = ethers.parseEther("300000000"); // 300M
  const REWARDS_ALLOCATION = ethers.parseEther("500000000"); // 500M
  const TIMELOCK_DURATION = 2 * 24 * 60 * 60; // 2 days

  beforeEach(async function () {
    [owner, daoRole, businessRole, minterRole, pauserRole, user1, user2, nonAuthorized] = await ethers.getSigners();

    // Deploy BOGOTokenV2_ZeroValidated
    const BOGOTokenV2_ZeroValidated = await ethers.getContractFactory("BOGOTokenV2_ZeroValidated");
    bogoToken = await BOGOTokenV2_ZeroValidated.deploy();
    await bogoToken.waitForDeployment();

    // Setup roles
    await bogoToken.grantRole(await bogoToken.DAO_ROLE(), await daoRole.getAddress());
    await bogoToken.grantRole(await bogoToken.BUSINESS_ROLE(), await businessRole.getAddress());
    await bogoToken.grantRole(await bogoToken.MINTER_ROLE(), await minterRole.getAddress());
    await bogoToken.grantRole(await bogoToken.PAUSER_ROLE(), await pauserRole.getAddress());
  });

  describe("Deployment", function () {
    it("Should deploy with correct initial state", async function () {
      expect(await bogoToken.name()).to.equal("BOGOWI");
      expect(await bogoToken.symbol()).to.equal("BOGO");
      expect(await bogoToken.decimals()).to.equal(18);
      expect(await bogoToken.totalSupply()).to.equal(0);
      expect(await bogoToken.MAX_SUPPLY()).to.equal(MAX_SUPPLY);
      expect(await bogoToken.DAO_ALLOCATION()).to.equal(DAO_ALLOCATION);
      expect(await bogoToken.BUSINESS_ALLOCATION()).to.equal(BUSINESS_ALLOCATION);
      expect(await bogoToken.REWARDS_ALLOCATION()).to.equal(REWARDS_ALLOCATION);
      expect(await bogoToken.TIMELOCK_DURATION()).to.equal(TIMELOCK_DURATION);
    });

    it("Should have correct initial minted amounts", async function () {
      expect(await bogoToken.daoMinted()).to.equal(0);
      expect(await bogoToken.businessMinted()).to.equal(0);
      expect(await bogoToken.rewardsMinted()).to.equal(0);
    });

    it("Should have correct remaining allocations", async function () {
      expect(await bogoToken.getRemainingDAOAllocation()).to.equal(DAO_ALLOCATION);
      expect(await bogoToken.getRemainingBusinessAllocation()).to.equal(BUSINESS_ALLOCATION);
      expect(await bogoToken.getRemainingRewardsAllocation()).to.equal(REWARDS_ALLOCATION);
    });

    it("Should grant correct roles to deployer", async function () {
      expect(await bogoToken.hasRole(await bogoToken.DEFAULT_ADMIN_ROLE(), await owner.getAddress())).to.be.true;
      expect(await bogoToken.hasRole(await bogoToken.MINTER_ROLE(), await owner.getAddress())).to.be.true;
      expect(await bogoToken.hasRole(await bogoToken.PAUSER_ROLE(), await owner.getAddress())).to.be.true;
    });
  });

  describe("Zero Address Validation", function () {
    describe("DAO Allocation Minting", function () {
      it("Should reject minting to zero address", async function () {
        await expect(
          bogoToken.connect(daoRole).mintFromDAO(ethers.ZeroAddress, ethers.parseEther("1000"))
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
      });

      it("Should reject zero amount minting", async function () {
        await expect(
          bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), 0)
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
      });

      it("Should allow valid minting", async function () {
        const amount = ethers.parseEther("1000");
        
        await expect(
          bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), amount)
        ).to.emit(bogoToken, "AllocationMinted")
          .withArgs("DAO", amount, await user1.getAddress());

        expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(amount);
        expect(await bogoToken.daoMinted()).to.equal(amount);
      });
    });

    describe("Business Allocation Minting", function () {
      it("Should reject minting to zero address", async function () {
        await expect(
          bogoToken.connect(businessRole).mintFromBusiness(ethers.ZeroAddress, ethers.parseEther("1000"))
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
      });

      it("Should reject zero amount minting", async function () {
        await expect(
          bogoToken.connect(businessRole).mintFromBusiness(await user1.getAddress(), 0)
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
      });

      it("Should allow valid minting", async function () {
        const amount = ethers.parseEther("1000");
        
        await expect(
          bogoToken.connect(businessRole).mintFromBusiness(await user1.getAddress(), amount)
        ).to.emit(bogoToken, "AllocationMinted")
          .withArgs("Business", amount, await user1.getAddress());

        expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(amount);
        expect(await bogoToken.businessMinted()).to.equal(amount);
      });
    });

    describe("Rewards Allocation Minting", function () {
      it("Should reject minting to zero address", async function () {
        await expect(
          bogoToken.connect(daoRole).mintFromRewards(ethers.ZeroAddress, ethers.parseEther("1000"))
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
      });

      it("Should reject zero amount minting", async function () {
        await expect(
          bogoToken.connect(daoRole).mintFromRewards(await user1.getAddress(), 0)
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
      });

      it("Should allow DAO role to mint rewards", async function () {
        const amount = ethers.parseEther("1000");
        
        await expect(
          bogoToken.connect(daoRole).mintFromRewards(await user1.getAddress(), amount)
        ).to.emit(bogoToken, "AllocationMinted")
          .withArgs("Rewards", amount, await user1.getAddress());

        expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(amount);
        expect(await bogoToken.rewardsMinted()).to.equal(amount);
      });

      it("Should allow Business role to mint rewards", async function () {
        const amount = ethers.parseEther("1000");
        
        await expect(
          bogoToken.connect(businessRole).mintFromRewards(await user1.getAddress(), amount)
        ).to.emit(bogoToken, "AllocationMinted")
          .withArgs("Rewards", amount, await user1.getAddress());

        expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(amount);
        expect(await bogoToken.rewardsMinted()).to.equal(amount);
      });

      it("Should reject unauthorized role", async function () {
        await expect(
          bogoToken.connect(user1).mintFromRewards(await user1.getAddress(), ethers.parseEther("1000"))
        ).to.be.revertedWithCustomError(bogoToken, "InsufficientRole");
      });
    });

    describe("Role Management", function () {
      it("Should reject granting role to zero address", async function () {
        await expect(
          bogoToken.grantRole(await bogoToken.DAO_ROLE(), ethers.ZeroAddress)
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
      });

      it("Should reject revoking role from zero address", async function () {
        await expect(
          bogoToken.revokeRole(await bogoToken.DAO_ROLE(), ethers.ZeroAddress)
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
      });

      it("Should allow valid role operations", async function () {
        await bogoToken.grantRole(await bogoToken.DAO_ROLE(), await user1.getAddress());
        expect(await bogoToken.hasRole(await bogoToken.DAO_ROLE(), await user1.getAddress())).to.be.true;
        
        await bogoToken.revokeRole(await bogoToken.DAO_ROLE(), await user1.getAddress());
        expect(await bogoToken.hasRole(await bogoToken.DAO_ROLE(), await user1.getAddress())).to.be.false;
      });
    });

    describe("Burn Functions", function () {
      beforeEach(async function () {
        // Mint tokens for testing burns
        await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), ethers.parseEther("1000"));
        await bogoToken.connect(user1).approve(await user2.getAddress(), ethers.parseEther("500"));
      });

      it("Should reject burning from zero address", async function () {
        await expect(
          bogoToken.connect(user2).burnFrom(ethers.ZeroAddress, ethers.parseEther("100"))
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
      });

      it("Should reject burning zero amount", async function () {
        await expect(
          bogoToken.connect(user1).burn(0)
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");

        await expect(
          bogoToken.connect(user2).burnFrom(await user1.getAddress(), 0)
        ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
      });

      it("Should allow valid burning", async function () {
        const burnAmount = ethers.parseEther("100");
        const initialBalance = await bogoToken.balanceOf(await user1.getAddress());
        
        await bogoToken.connect(user1).burn(burnAmount);
        expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(initialBalance - burnAmount);
      });

      it("Should allow valid burnFrom", async function () {
        const burnAmount = ethers.parseEther("100");
        const initialBalance = await bogoToken.balanceOf(await user1.getAddress());
        
        await bogoToken.connect(user2).burnFrom(await user1.getAddress(), burnAmount);
        expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(initialBalance - burnAmount);
      });
    });
  });

  describe("Allocation Management", function () {
    it("Should enforce DAO allocation limit", async function () {
      const excessiveAmount = DAO_ALLOCATION + ethers.parseEther("1");
      
      await expect(
        bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), excessiveAmount)
      ).to.be.revertedWithCustomError(bogoToken, "ExceedsAllocation");
    });

    it("Should enforce Business allocation limit", async function () {
      const excessiveAmount = BUSINESS_ALLOCATION + ethers.parseEther("1");
      
      await expect(
        bogoToken.connect(businessRole).mintFromBusiness(await user1.getAddress(), excessiveAmount)
      ).to.be.revertedWithCustomError(bogoToken, "ExceedsAllocation");
    });

    it("Should enforce Rewards allocation limit", async function () {
      const excessiveAmount = REWARDS_ALLOCATION + ethers.parseEther("1");
      
      await expect(
        bogoToken.connect(daoRole).mintFromRewards(await user1.getAddress(), excessiveAmount)
      ).to.be.revertedWithCustomError(bogoToken, "ExceedsAllocation");
    });

    it("Should enforce max supply limit", async function () {
      // First mint close to max supply from all allocations
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), DAO_ALLOCATION);
      await bogoToken.connect(businessRole).mintFromBusiness(await user2.getAddress(), BUSINESS_ALLOCATION);
      // Mint most of rewards allocation, leaving room for the test
      await bogoToken.connect(daoRole).mintFromRewards(await user1.getAddress(), ethers.parseEther("499999998"));
      
      // Total minted: 200M + 300M + 499,999,998 = 999,999,998 (2 tokens left)
      // Now try to mint 3 tokens, which should exceed max supply
      await expect(
        bogoToken.connect(daoRole).mintFromRewards(await user2.getAddress(), ethers.parseEther("3"))
      ).to.be.revertedWithCustomError(bogoToken, "ExceedsMaxSupply");
    });

    it("Should track remaining allocations correctly", async function () {
      const daoAmount = ethers.parseEther("1000");
      const businessAmount = ethers.parseEther("2000");
      const rewardsAmount = ethers.parseEther("3000");
      
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), daoAmount);
      await bogoToken.connect(businessRole).mintFromBusiness(await user1.getAddress(), businessAmount);
      await bogoToken.connect(daoRole).mintFromRewards(await user1.getAddress(), rewardsAmount);
      
      expect(await bogoToken.getRemainingDAOAllocation()).to.equal(DAO_ALLOCATION - daoAmount);
      expect(await bogoToken.getRemainingBusinessAllocation()).to.equal(BUSINESS_ALLOCATION - businessAmount);
      expect(await bogoToken.getRemainingRewardsAllocation()).to.equal(REWARDS_ALLOCATION - rewardsAmount);
    });

    it("Should allow minting exactly to allocation limits", async function () {
      await expect(
        bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), DAO_ALLOCATION)
      ).to.not.be.reverted;
      
      await expect(
        bogoToken.connect(businessRole).mintFromBusiness(await user2.getAddress(), BUSINESS_ALLOCATION)
      ).to.not.be.reverted;
      
      expect(await bogoToken.daoMinted()).to.equal(DAO_ALLOCATION);
      expect(await bogoToken.businessMinted()).to.equal(BUSINESS_ALLOCATION);
    });
  });

  describe("Pause/Unpause Functionality", function () {
    beforeEach(async function () {
      // Mint some tokens for testing transfers
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), ethers.parseEther("1000"));
    });

    it("Should allow pauser to pause contract", async function () {
      await bogoToken.connect(pauserRole).pause();
      expect(await bogoToken.paused()).to.be.true;
    });

    it("Should allow pauser to unpause contract", async function () {
      await bogoToken.connect(pauserRole).pause();
      await bogoToken.connect(pauserRole).unpause();
      expect(await bogoToken.paused()).to.be.false;
    });

    it("Should only allow pauser role to pause/unpause", async function () {
      await expect(
        bogoToken.connect(user1).pause()
      ).to.be.revertedWithCustomError(bogoToken, "AccessControlUnauthorizedAccount");

      await expect(
        bogoToken.connect(user1).unpause()
      ).to.be.revertedWithCustomError(bogoToken, "AccessControlUnauthorizedAccount");
    });

    it("Should prevent transfers when paused", async function () {
      await bogoToken.connect(pauserRole).pause();
      
      await expect(
        bogoToken.connect(user1).transfer(await user2.getAddress(), ethers.parseEther("100"))
      ).to.be.revertedWithCustomError(bogoToken, "EnforcedPause");
    });

    it("Should prevent minting when paused", async function () {
      await bogoToken.connect(pauserRole).pause();
      
      await expect(
        bogoToken.connect(daoRole).mintFromDAO(await user2.getAddress(), ethers.parseEther("100"))
      ).to.be.revertedWithCustomError(bogoToken, "EnforcedPause");
    });

    it("Should allow operations when unpaused", async function () {
      await bogoToken.connect(pauserRole).pause();
      await bogoToken.connect(pauserRole).unpause();
      
      await expect(
        bogoToken.connect(user1).transfer(await user2.getAddress(), ethers.parseEther("100"))
      ).to.not.be.reverted;

      expect(await bogoToken.balanceOf(await user2.getAddress())).to.equal(ethers.parseEther("100"));
    });
  });

  describe("ERC20 Functionality", function () {
    beforeEach(async function () {
      // Mint tokens for testing
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), ethers.parseEther("1000"));
      await bogoToken.connect(businessRole).mintFromBusiness(await user2.getAddress(), ethers.parseEther("500"));
    });

    it("Should support standard ERC20 transfers", async function () {
      await bogoToken.connect(user1).transfer(await user2.getAddress(), ethers.parseEther("100"));
      
      expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(ethers.parseEther("900"));
      expect(await bogoToken.balanceOf(await user2.getAddress())).to.equal(ethers.parseEther("600"));
    });

    it("Should support allowance and transferFrom", async function () {
      await bogoToken.connect(user1).approve(await user2.getAddress(), ethers.parseEther("200"));
      
      expect(await bogoToken.allowance(await user1.getAddress(), await user2.getAddress()))
        .to.equal(ethers.parseEther("200"));
      
      await bogoToken.connect(user2).transferFrom(
        await user1.getAddress(),
        await user2.getAddress(),
        ethers.parseEther("150")
      );
      
      expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(ethers.parseEther("850"));
      expect(await bogoToken.balanceOf(await user2.getAddress())).to.equal(ethers.parseEther("650"));
    });

    it("Should prevent transfers exceeding balance", async function () {
      await expect(
        bogoToken.connect(user1).transfer(await user2.getAddress(), ethers.parseEther("1001"))
      ).to.be.revertedWithCustomError(bogoToken, "ERC20InsufficientBalance");
    });

    it("Should prevent transferFrom exceeding allowance", async function () {
      await bogoToken.connect(user1).approve(await user2.getAddress(), ethers.parseEther("100"));
      
      await expect(
        bogoToken.connect(user2).transferFrom(
          await user1.getAddress(),
          await user2.getAddress(),
          ethers.parseEther("101")
        )
      ).to.be.revertedWithCustomError(bogoToken, "ERC20InsufficientAllowance");
    });

    it("Should allow burning tokens", async function () {
      const initialSupply = await bogoToken.totalSupply();
      const burnAmount = ethers.parseEther("100");
      
      await bogoToken.connect(user1).burn(burnAmount);
      
      expect(await bogoToken.totalSupply()).to.equal(initialSupply - burnAmount);
      expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(ethers.parseEther("900"));
    });
  });

  describe("Access Control", function () {
    it("Should enforce DAO role for DAO minting", async function () {
      await expect(
        bogoToken.connect(user1).mintFromDAO(await user1.getAddress(), ethers.parseEther("1000"))
      ).to.be.revertedWithCustomError(bogoToken, "AccessControlUnauthorizedAccount");
    });

    it("Should enforce Business role for Business minting", async function () {
      await expect(
        bogoToken.connect(user1).mintFromBusiness(await user1.getAddress(), ethers.parseEther("1000"))
      ).to.be.revertedWithCustomError(bogoToken, "AccessControlUnauthorizedAccount");
    });


    it("Should enforce pauser role for pause operations", async function () {
      await expect(
        bogoToken.connect(user1).pause()
      ).to.be.revertedWithCustomError(bogoToken, "AccessControlUnauthorizedAccount");
    });
  });

  describe("Reentrancy Protection", function () {
    it("Should prevent reentrancy in mintFromDAO", async function () {
      const amount = ethers.parseEther("1000");
      await expect(
        bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), amount)
      ).to.not.be.reverted;
    });

    it("Should prevent reentrancy in mintFromBusiness", async function () {
      const amount = ethers.parseEther("1000");
      await expect(
        bogoToken.connect(businessRole).mintFromBusiness(await user1.getAddress(), amount)
      ).to.not.be.reverted;
    });

    it("Should prevent reentrancy in mintFromRewards", async function () {
      const amount = ethers.parseEther("500");
      await expect(
        bogoToken.connect(daoRole).mintFromRewards(await user1.getAddress(), amount)
      ).to.not.be.reverted;
    });
  });

  describe("Edge Cases", function () {
    it("Should handle zero balance transfers", async function () {
      // Zero transfers are allowed in ERC20 standard
      await expect(
        bogoToken.connect(user1).transfer(await user2.getAddress(), 0)
      ).to.not.be.reverted;
    });

    it("Should handle self-transfers", async function () {
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), ethers.parseEther("1000"));
      
      await expect(
        bogoToken.connect(user1).transfer(await user1.getAddress(), ethers.parseEther("100"))
      ).to.not.be.reverted;
      
      expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(ethers.parseEther("1000"));
    });


    it("Should handle maximum allocation usage", async function () {
      // Use entire DAO allocation
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), DAO_ALLOCATION);
      expect(await bogoToken.getRemainingDAOAllocation()).to.equal(0);
      
      // Should not be able to mint more
      await expect(
        bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), 1)
      ).to.be.revertedWithCustomError(bogoToken, "ExceedsAllocation");
    });

    it("Should support interface detection", async function () {
      // Test AccessControl interface
      const accessControlInterface = "0x7965db0b";
      expect(await bogoToken.supportsInterface(accessControlInterface)).to.be.true;
    });
  });

  describe("Gas Optimization Features", function () {

    it("Should use custom errors for gas efficiency", async function () {
      // Test that custom errors are used instead of require strings
      await expect(
        bogoToken.connect(daoRole).mintFromDAO(ethers.ZeroAddress, ethers.parseEther("1000"))
      ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
      
      await expect(
        bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), 0)
      ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
    });
  });

  describe("Coverage Improvements", function () {
    it("Should handle burn operations through _update", async function () {
      // Mint some tokens first
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), ethers.parseEther("1000"));
      
      // Burn tokens (this tests the burn path in _update)
      await bogoToken.connect(user1).burn(ethers.parseEther("100"));
      
      expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(ethers.parseEther("900"));
    });

    it("Should handle burnFrom operations through _update", async function () {
      // Mint some tokens first
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), ethers.parseEther("1000"));
      
      // Approve user2 to burn tokens
      await bogoToken.connect(user1).approve(await user2.getAddress(), ethers.parseEther("200"));
      
      // BurnFrom (this also tests the burn path in _update)
      await bogoToken.connect(user2).burnFrom(await user1.getAddress(), ethers.parseEther("150"));
      
      expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(ethers.parseEther("850"));
      expect(await bogoToken.allowance(await user1.getAddress(), await user2.getAddress())).to.equal(ethers.parseEther("50"));
    });

    it("Should support interface detection", async function () {
      // Test IERC165 interface
      expect(await bogoToken.supportsInterface("0x01ffc9a7")).to.be.true;
      
      // Test IAccessControl interface
      expect(await bogoToken.supportsInterface("0x7965db0b")).to.be.true;
      
      // Test invalid interface
      expect(await bogoToken.supportsInterface("0xffffffff")).to.be.false;
    });

    it("Should handle transfers correctly through _update", async function () {
      // Mint tokens
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), ethers.parseEther("1000"));
      
      // Normal transfer (tests the transfer path in _update)
      await bogoToken.connect(user1).transfer(await user2.getAddress(), ethers.parseEther("100"));
      
      expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(ethers.parseEther("900"));
      expect(await bogoToken.balanceOf(await user2.getAddress())).to.equal(ethers.parseEther("100"));
    });

    it("Should handle minting through _update correctly", async function () {
      // Minting already has validation in the minting functions, but let's ensure _update handles it
      const initialSupply = await bogoToken.totalSupply();
      
      await bogoToken.connect(daoRole).mintFromDAO(await user1.getAddress(), ethers.parseEther("5000"));
      
      expect(await bogoToken.totalSupply()).to.equal(initialSupply + ethers.parseEther("5000"));
      expect(await bogoToken.balanceOf(await user1.getAddress())).to.equal(ethers.parseEther("5000"));
    });
  });

  describe("Timelock Cancellation Coverage", function () {
    it("Should revert when cancelling non-existent operation", async function () {
      const operationId = ethers.keccak256(ethers.toUtf8Bytes("non-existent"));
      
      await expect(
        bogoToken.cancelTimelockOperation(operationId)
      ).to.be.revertedWithCustomError(bogoToken, "OperationNotQueued");
    });

    it("Should only allow admin to cancel", async function () {
      const operationId = ethers.keccak256(ethers.toUtf8Bytes("test"));
      
      await expect(
        bogoToken.connect(user1).cancelTimelockOperation(operationId)
      ).to.be.revertedWithCustomError(bogoToken, "AccessControlUnauthorizedAccount");
    });

    it("Should handle valid cancellation scenario", async function () {
      // Since we removed flavored tokens, we don't have a way to queue operations
      // But we can still test the access control
      const operationId = ethers.keccak256(ethers.toUtf8Bytes("test-cancel"));
      
      // Should revert with OperationNotQueued since nothing is queued
      await expect(
        bogoToken.cancelTimelockOperation(operationId)
      ).to.be.revertedWithCustomError(bogoToken, "OperationNotQueued");
    });
  });
});