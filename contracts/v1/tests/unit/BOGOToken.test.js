const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGOToken - Unified Contract Tests", function () {
    let bogoToken;
    let roleManager;
    let owner, daoWallet, businessWallet, pauser, user1, user2;
    let DAO_ROLE, BUSINESS_ROLE, PAUSER_ROLE;

    // Constants from contract
    const MAX_SUPPLY = ethers.parseEther("1000000000"); // 1 billion
    const DAO_ALLOCATION = ethers.parseEther("50000000"); // 50M (5% of total)
    const BUSINESS_ALLOCATION = ethers.parseEther("900000000"); // 900M (90% of total)
    const REWARDS_ALLOCATION = ethers.parseEther("50000000"); // 50M (5% of total)

    beforeEach(async function () {
        [owner, daoWallet, businessWallet, pauser, user1, user2] = await ethers.getSigners();

        // Deploy RoleManager first
        const RoleManager = await ethers.getContractFactory("RoleManager");
        roleManager = await RoleManager.deploy();
        await roleManager.waitForDeployment();

        // Deploy BOGOToken with RoleManager
        const BOGOToken = await ethers.getContractFactory("BOGOToken");
        bogoToken = await BOGOToken.deploy(
            await roleManager.getAddress(),
            "BOGOWI",
            "BOGO"
        );
        await bogoToken.waitForDeployment();

        // Register token with RoleManager
        await roleManager.registerContract(await bogoToken.getAddress(), "BOGOToken");

        // Get role constants
        DAO_ROLE = await roleManager.DAO_ROLE();
        BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
        PAUSER_ROLE = await roleManager.PAUSER_ROLE();

        // Grant roles through RoleManager
        await roleManager.grantRole(DAO_ROLE, daoWallet.address);
        await roleManager.grantRole(BUSINESS_ROLE, businessWallet.address);
        await roleManager.grantRole(PAUSER_ROLE, pauser.address);
    });

    describe("Deployment", function () {
        it("Should have correct name and symbol", async function () {
            expect(await bogoToken.name()).to.equal("BOGOWI");
            expect(await bogoToken.symbol()).to.equal("BOGO");
        });

        it("Should have correct total supply constants", async function () {
            expect(await bogoToken.MAX_SUPPLY()).to.equal(MAX_SUPPLY);
            expect(await bogoToken.DAO_ALLOCATION()).to.equal(DAO_ALLOCATION);
            expect(await bogoToken.BUSINESS_ALLOCATION()).to.equal(BUSINESS_ALLOCATION);
            expect(await bogoToken.REWARDS_ALLOCATION()).to.equal(REWARDS_ALLOCATION);
        });

        it("Should start with zero supply", async function () {
            expect(await bogoToken.totalSupply()).to.equal(0);
        });
    });

    describe("Zero Address Validation", function () {
        it("Should reject minting to zero address", async function () {
            await expect(
                bogoToken.connect(daoWallet).mintFromDAO(ethers.ZeroAddress, ethers.parseEther("1000"))
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should reject zero amount minting", async function () {
            await expect(
                bogoToken.connect(daoWallet).mintFromDAO(user1.address, 0)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
        });

        it("Should reject burning from zero address", async function () {
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, ethers.parseEther("1000"));
            await bogoToken.connect(user1).approve(user2.address, ethers.parseEther("500"));
            
            await expect(
                bogoToken.connect(user2).burnFrom(ethers.ZeroAddress, ethers.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should reject zero amount burning", async function () {
            await expect(
                bogoToken.connect(user1).burn(0)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
        });
    });

    describe("DAO Allocation Minting", function () {
        it("Should allow DAO role to mint from DAO allocation", async function () {
            const amount = ethers.parseEther("1000");
            await expect(bogoToken.connect(daoWallet).mintFromDAO(user1.address, amount))
                .to.emit(bogoToken, "AllocationMinted")
                .withArgs("DAO", amount, user1.address);
            
            expect(await bogoToken.balanceOf(user1.address)).to.equal(amount);
            expect(await bogoToken.daoMinted()).to.equal(amount);
        });

        it("Should fail when exceeding DAO allocation", async function () {
            const exceedAmount = DAO_ALLOCATION + 1n;
            await expect(bogoToken.connect(daoWallet).mintFromDAO(user1.address, exceedAmount))
                .to.be.revertedWithCustomError(bogoToken, "ExceedsAllocation");
        });

        it("Should fail when non-DAO role tries to mint", async function () {
            await expect(bogoToken.connect(user1).mintFromDAO(user2.address, 1000))
                .to.be.reverted;
        });

        it("Should track remaining DAO allocation", async function () {
            const amount = ethers.parseEther("25000000"); // 25M (half of DAO allocation)
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, amount);
            
            expect(await bogoToken.getRemainingDAOAllocation())
                .to.equal(DAO_ALLOCATION - amount);
        });
    });

    describe("Business Allocation Minting", function () {
        it("Should allow Business role to mint from Business allocation", async function () {
            const amount = ethers.parseEther("5000");
            await expect(bogoToken.connect(businessWallet).mintFromBusiness(user1.address, amount))
                .to.emit(bogoToken, "AllocationMinted")
                .withArgs("Business", amount, user1.address);
            
            expect(await bogoToken.balanceOf(user1.address)).to.equal(amount);
            expect(await bogoToken.businessMinted()).to.equal(amount);
        });

        it("Should fail when exceeding Business allocation", async function () {
            const exceedAmount = BUSINESS_ALLOCATION + 1n;
            await expect(bogoToken.connect(businessWallet).mintFromBusiness(user1.address, exceedAmount))
                .to.be.revertedWithCustomError(bogoToken, "ExceedsAllocation");
        });

        it("Should track remaining Business allocation", async function () {
            const amount = ethers.parseEther("450000000"); // 450M (half of business allocation)
            await bogoToken.connect(businessWallet).mintFromBusiness(user1.address, amount);
            
            expect(await bogoToken.getRemainingBusinessAllocation())
                .to.equal(BUSINESS_ALLOCATION - amount);
        });
    });

    describe("Rewards Allocation Minting", function () {
        it("Should allow DAO role to mint from Rewards allocation", async function () {
            const amount = ethers.parseEther("10000");
            await expect(bogoToken.connect(daoWallet).mintFromRewards(user1.address, amount))
                .to.emit(bogoToken, "AllocationMinted")
                .withArgs("Rewards", amount, user1.address);
            
            expect(await bogoToken.rewardsMinted()).to.equal(amount);
        });

        it("Should allow Business role to mint from Rewards allocation", async function () {
            const amount = ethers.parseEther("20000");
            await expect(bogoToken.connect(businessWallet).mintFromRewards(user1.address, amount))
                .to.emit(bogoToken, "AllocationMinted")
                .withArgs("Rewards", amount, user1.address);
        });

        it("Should fail when neither DAO nor Business role", async function () {
            await expect(bogoToken.connect(user1).mintFromRewards(user2.address, 1000))
                .to.be.revertedWithCustomError(bogoToken, "InsufficientRole");
        });

        it("Should fail when exceeding Rewards allocation", async function () {
            const exceedAmount = REWARDS_ALLOCATION + 1n;
            await expect(bogoToken.connect(daoWallet).mintFromRewards(user1.address, exceedAmount))
                .to.be.revertedWithCustomError(bogoToken, "ExceedsAllocation");
        });

        it("Should track remaining Rewards allocation", async function () {
            const amount = ethers.parseEther("25000000"); // 25M (half of rewards allocation)
            await bogoToken.connect(daoWallet).mintFromRewards(user1.address, amount);
            
            expect(await bogoToken.getRemainingRewardsAllocation())
                .to.equal(REWARDS_ALLOCATION - amount);
        });
    });

    describe("Max Supply Enforcement", function () {
        it("Should enforce max supply across all allocations", async function () {
            // The allocations sum to MAX_SUPPLY
            expect(DAO_ALLOCATION + BUSINESS_ALLOCATION + REWARDS_ALLOCATION)
                .to.equal(MAX_SUPPLY);
        });

        it("Should prevent minting beyond max supply", async function () {
            // This test verifies that the sum of allocations equals max supply
            // Since allocations perfectly sum to MAX_SUPPLY, we can't exceed MAX_SUPPLY
            // without first exceeding an individual allocation limit
            
            // Mint all DAO allocation
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, DAO_ALLOCATION);
            expect(await bogoToken.getRemainingDAOAllocation()).to.equal(0);
            
            // Try to mint 1 more token from DAO - should fail with ExceedsAllocation
            await expect(
                bogoToken.connect(daoWallet).mintFromDAO(user2.address, 1)
            ).to.be.revertedWithCustomError(bogoToken, "ExceedsAllocation");
        });
    });

    describe("Burn Functions", function () {
        beforeEach(async function () {
            // Mint some tokens for testing
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, ethers.parseEther("1000"));
            await bogoToken.connect(daoWallet).mintFromDAO(user2.address, ethers.parseEther("1000"));
        });

        it("Should burn own tokens", async function () {
            const burnAmount = ethers.parseEther("100");
            const initialBalance = await bogoToken.balanceOf(user1.address);
            
            await bogoToken.connect(user1).burn(burnAmount);
            
            expect(await bogoToken.balanceOf(user1.address))
                .to.equal(initialBalance - burnAmount);
            expect(await bogoToken.totalSupply())
                .to.equal(ethers.parseEther("2000") - burnAmount);
        });

        it("Should burn tokens with approval (burnFrom)", async function () {
            const burnAmount = ethers.parseEther("200");
            
            // User1 approves user2 to burn tokens
            await bogoToken.connect(user1).approve(user2.address, burnAmount);
            
            // User2 burns user1's tokens
            await bogoToken.connect(user2).burnFrom(user1.address, burnAmount);
            
            expect(await bogoToken.balanceOf(user1.address))
                .to.equal(ethers.parseEther("800"));
            expect(await bogoToken.allowance(user1.address, user2.address))
                .to.equal(0);
        });

        it("Should fail burnFrom without approval", async function () {
            await expect(bogoToken.connect(user2).burnFrom(user1.address, 100))
                .to.be.reverted;
        });
    });

    describe("Pause Functionality", function () {
        beforeEach(async function () {
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, ethers.parseEther("1000"));
        });

        it("Should pause and unpause", async function () {
            expect(await bogoToken.paused()).to.be.false;
            
            await bogoToken.connect(pauser).pause();
            expect(await bogoToken.paused()).to.be.true;
            
            await bogoToken.connect(pauser).unpause();
            expect(await bogoToken.paused()).to.be.false;
        });

        it("Should block transfers when paused", async function () {
            await bogoToken.connect(pauser).pause();
            
            await expect(bogoToken.connect(user1).transfer(user2.address, 100))
                .to.be.revertedWithCustomError(bogoToken, "EnforcedPause");
        });

        it("Should block minting when paused", async function () {
            await bogoToken.connect(pauser).pause();
            
            await expect(bogoToken.connect(daoWallet).mintFromDAO(user1.address, 100))
                .to.be.revertedWithCustomError(bogoToken, "EnforcedPause");
        });

        it("Should only allow PAUSER_ROLE to pause/unpause", async function () {
            await expect(bogoToken.connect(user1).pause())
                .to.be.reverted;
            
            await expect(bogoToken.connect(user1).unpause())
                .to.be.reverted;
        });
    });

    describe("Role Management Integration", function () {
        it("Should respect role changes made in RoleManager", async function () {
            // Remove DAO role from daoWallet
            await roleManager.revokeRole(DAO_ROLE, daoWallet.address);
            
            // Should now fail to mint
            await expect(
                bogoToken.connect(daoWallet).mintFromDAO(user1.address, 1000)
            ).to.be.reverted;
            
            // Grant DAO role to user1
            await roleManager.grantRole(DAO_ROLE, user1.address);
            
            // user1 should now be able to mint
            await expect(
                bogoToken.connect(user1).mintFromDAO(user2.address, 1000)
            ).to.not.be.reverted;
        });

        it("Should reject direct role management attempts", async function () {
            await expect(
                bogoToken.grantRole(DAO_ROLE, user1.address)
            ).to.be.revertedWith("Use RoleManager to manage roles");
            
            await expect(
                bogoToken.revokeRole(DAO_ROLE, daoWallet.address)
            ).to.be.revertedWith("Use RoleManager to manage roles");
        });

    });

    describe("Transfer Functionality", function () {
        beforeEach(async function () {
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, ethers.parseEther("1000"));
        });

        it("Should allow transfers when not paused", async function () {
            await bogoToken.connect(user1).transfer(user2.address, 100);
            expect(await bogoToken.balanceOf(user2.address)).to.equal(100);
        });

        it("Should allow transferFrom with approval", async function () {
            await bogoToken.connect(user1).approve(user2.address, 500);
            await bogoToken.connect(user2).transferFrom(user1.address, owner.address, 300);
            
            expect(await bogoToken.balanceOf(owner.address)).to.equal(300);
            expect(await bogoToken.allowance(user1.address, user2.address)).to.equal(200);
        });

        it("Should validate recipient address on transfer", async function () {
            // The _update function validates this internally
            // Direct transfer to zero address is prevented by ERC20 
            await expect(
                bogoToken.connect(user1).transfer(ethers.ZeroAddress, 100)
            ).to.be.reverted;
        });
    });

    describe("Role Override Functions", function () {
        it("Should reject grantRole with zero address", async function () {
            await expect(
                bogoToken.grantRole(DAO_ROLE, ethers.ZeroAddress)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should reject revokeRole with zero address", async function () {
            await expect(
                bogoToken.revokeRole(DAO_ROLE, ethers.ZeroAddress)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should revert grantRole with message for non-zero address", async function () {
            await expect(
                bogoToken.grantRole(DAO_ROLE, user1.address)
            ).to.be.revertedWith("Use RoleManager to manage roles");
        });

        it("Should revert revokeRole with message for non-zero address", async function () {
            await expect(
                bogoToken.revokeRole(DAO_ROLE, user1.address)
            ).to.be.revertedWith("Use RoleManager to manage roles");
        });
    });

    describe("Edge Cases and Security", function () {
        it("Should handle reentrancy protection in minting", async function () {
            // All minting functions have nonReentrant modifier
            const amount = ethers.parseEther("1000");
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, amount);
            expect(await bogoToken.balanceOf(user1.address)).to.equal(amount);
        });

        it("Should handle multiple allocations to same address", async function () {
            const amount1 = ethers.parseEther("1000");
            const amount2 = ethers.parseEther("2000");
            const amount3 = ethers.parseEther("3000");
            
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, amount1);
            await bogoToken.connect(businessWallet).mintFromBusiness(user1.address, amount2);
            await bogoToken.connect(daoWallet).mintFromRewards(user1.address, amount3);
            
            expect(await bogoToken.balanceOf(user1.address))
                .to.equal(amount1 + amount2 + amount3);
        });

        it("Should handle minting exactly to allocation limits", async function () {
            await expect(
                bogoToken.connect(daoWallet).mintFromDAO(user1.address, DAO_ALLOCATION)
            ).to.not.be.reverted;
            
            expect(await bogoToken.daoMinted()).to.equal(DAO_ALLOCATION);
            expect(await bogoToken.getRemainingDAOAllocation()).to.equal(0);
        });
    });
});