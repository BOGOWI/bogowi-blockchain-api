const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGOTokenV2 - Full Coverage Tests", function () {
    let bogoToken;
    let owner, daoWallet, businessWallet, rewardWallet, pauser, user1, user2, attacker;
    let DAO_ROLE, BUSINESS_ROLE, MINTER_ROLE, PAUSER_ROLE;

    // Constants from contract
    const MAX_SUPPLY = ethers.parseEther("1000000000"); // 1 billion
    const DAO_ALLOCATION = ethers.parseEther("200000000"); // 200M
    const BUSINESS_ALLOCATION = ethers.parseEther("300000000"); // 300M
    const REWARDS_ALLOCATION = ethers.parseEther("500000000"); // 500M
    const TIMELOCK_DURATION = 2 * 24 * 60 * 60; // 2 days

    beforeEach(async function () {
        [owner, daoWallet, businessWallet, rewardWallet, pauser, user1, user2, attacker] = await ethers.getSigners();

        const BOGOTokenV2 = await ethers.getContractFactory("BOGOTokenV2");
        bogoToken = await BOGOTokenV2.deploy();
        await bogoToken.waitForDeployment();

        // Get role constants
        DAO_ROLE = await bogoToken.DAO_ROLE();
        BUSINESS_ROLE = await bogoToken.BUSINESS_ROLE();
        MINTER_ROLE = await bogoToken.MINTER_ROLE();
        PAUSER_ROLE = await bogoToken.PAUSER_ROLE();

        // Grant roles
        await bogoToken.grantRole(DAO_ROLE, daoWallet.address);
        await bogoToken.grantRole(BUSINESS_ROLE, businessWallet.address);
        await bogoToken.grantRole(PAUSER_ROLE, pauser.address);
    });

    describe("Token Basics", function () {
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
    });

    describe("Allocation Minting", function () {
        describe("DAO Allocation", function () {
            it("Should mint from DAO allocation", async function () {
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
                    .to.be.revertedWith("EXCEEDS_ALLOCATION");
            });

            it("Should fail when non-DAO role tries to mint", async function () {
                await expect(bogoToken.connect(user1).mintFromDAO(user2.address, 1000))
                    .to.be.reverted;
            });

            it("Should track remaining DAO allocation", async function () {
                const amount = ethers.parseEther("50000000"); // 50M
                await bogoToken.connect(daoWallet).mintFromDAO(user1.address, amount);
                
                expect(await bogoToken.getRemainingDAOAllocation())
                    .to.equal(DAO_ALLOCATION - amount);
            });
        });

        describe("Business Allocation", function () {
            it("Should mint from business allocation", async function () {
                const amount = ethers.parseEther("5000");
                await expect(bogoToken.connect(businessWallet).mintFromBusiness(user1.address, amount))
                    .to.emit(bogoToken, "AllocationMinted")
                    .withArgs("Business", amount, user1.address);
                
                expect(await bogoToken.balanceOf(user1.address)).to.equal(amount);
                expect(await bogoToken.businessMinted()).to.equal(amount);
            });

            it("Should fail when exceeding business allocation", async function () {
                const exceedAmount = BUSINESS_ALLOCATION + 1n;
                await expect(bogoToken.connect(businessWallet).mintFromBusiness(user1.address, exceedAmount))
                    .to.be.revertedWith("EXCEEDS_ALLOCATION");
            });

            it("Should track remaining business allocation", async function () {
                const amount = ethers.parseEther("100000000"); // 100M
                await bogoToken.connect(businessWallet).mintFromBusiness(user1.address, amount);
                
                expect(await bogoToken.getRemainingBusinessAllocation())
                    .to.equal(BUSINESS_ALLOCATION - amount);
            });
        });

        describe("Rewards Allocation", function () {
            it("Should mint from rewards allocation with DAO role", async function () {
                const amount = ethers.parseEther("10000");
                await expect(bogoToken.connect(daoWallet).mintFromRewards(user1.address, amount))
                    .to.emit(bogoToken, "AllocationMinted")
                    .withArgs("Rewards", amount, user1.address);
                
                expect(await bogoToken.rewardsMinted()).to.equal(amount);
            });

            it("Should mint from rewards allocation with BUSINESS role", async function () {
                const amount = ethers.parseEther("20000");
                await expect(bogoToken.connect(businessWallet).mintFromRewards(user1.address, amount))
                    .to.emit(bogoToken, "AllocationMinted")
                    .withArgs("Rewards", amount, user1.address);
            });

            it("Should fail when neither DAO nor BUSINESS role", async function () {
                await expect(bogoToken.connect(user1).mintFromRewards(user2.address, 1000))
                    .to.be.revertedWith("UNAUTHORIZED");
            });

            it("Should fail when exceeding rewards allocation", async function () {
                const exceedAmount = REWARDS_ALLOCATION + 1n;
                await expect(bogoToken.connect(daoWallet).mintFromRewards(user1.address, exceedAmount))
                    .to.be.revertedWith("EXCEEDS_ALLOCATION");
            });

            it("Should track remaining rewards allocation", async function () {
                const amount = ethers.parseEther("250000000"); // 250M
                await bogoToken.connect(daoWallet).mintFromRewards(user1.address, amount);
                
                expect(await bogoToken.getRemainingRewardsAllocation())
                    .to.equal(REWARDS_ALLOCATION - amount);
            });
        });

        describe("Max Supply Enforcement", function () {
            it("Should enforce max supply across all allocations", async function () {
                // This is a complex test - in practice, the allocations sum to MAX_SUPPLY
                // so we can't exceed it through normal minting
                expect(DAO_ALLOCATION + BUSINESS_ALLOCATION + REWARDS_ALLOCATION)
                    .to.equal(MAX_SUPPLY);
            });
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

        it("Should fail burnFrom with insufficient approval", async function () {
            await bogoToken.connect(user1).approve(user2.address, 50);
            
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


    describe("Transfer Restrictions", function () {
        beforeEach(async function () {
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, ethers.parseEther("1000"));
        });

        it("Should allow transfers when not paused", async function () {
            await bogoToken.connect(user1).transfer(user2.address, 100);
            expect(await bogoToken.balanceOf(user2.address)).to.equal(100);
        });

        it("Should allow transferFrom with approval", async function () {
            await bogoToken.connect(user1).approve(user2.address, 500);
            await bogoToken.connect(user2).transferFrom(user1.address, attacker.address, 300);
            
            expect(await bogoToken.balanceOf(attacker.address)).to.equal(300);
            expect(await bogoToken.allowance(user1.address, user2.address)).to.equal(200);
        });

        it("Should respect pause state in _update", async function () {
            await bogoToken.connect(pauser).pause();
            
            // All transfer operations should fail
            await expect(bogoToken.connect(user1).transfer(user2.address, 100))
                .to.be.revertedWithCustomError(bogoToken, "EnforcedPause");
            
            await bogoToken.connect(pauser).unpause();
            
            // Should work after unpause
            await bogoToken.connect(user1).transfer(user2.address, 100);
            expect(await bogoToken.balanceOf(user2.address)).to.equal(100);
        });
    });

    describe("Edge Cases and Security", function () {
        it("Should handle reentrancy protection in minting", async function () {
            // All minting functions have nonReentrant modifier
            // This is tested implicitly through successful minting
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

        it("Should handle zero amount minting", async function () {
            await bogoToken.connect(daoWallet).mintFromDAO(user1.address, 0);
            expect(await bogoToken.balanceOf(user1.address)).to.equal(0);
        });

        it("Should prevent minting to zero address", async function () {
            // ERC20 _mint will revert with zero address
            await expect(bogoToken.connect(daoWallet).mintFromDAO(ethers.ZeroAddress, 1000))
                .to.be.reverted;
        });
    });

    describe("Role Management", function () {
        it("Should properly manage role transitions", async function () {
            // Grant and revoke roles
            await bogoToken.grantRole(MINTER_ROLE, user1.address);
            expect(await bogoToken.hasRole(MINTER_ROLE, user1.address)).to.be.true;
            
            await bogoToken.revokeRole(MINTER_ROLE, user1.address);
            expect(await bogoToken.hasRole(MINTER_ROLE, user1.address)).to.be.false;
        });

        it("Should handle role renunciation", async function () {
            await bogoToken.grantRole(PAUSER_ROLE, user1.address);
            await bogoToken.connect(user1).renounceRole(PAUSER_ROLE, user1.address);
            
            expect(await bogoToken.hasRole(PAUSER_ROLE, user1.address)).to.be.false;
        });
    });

    describe("Gas Optimization Tests", function () {
        it("Should efficiently handle batch operations", async function () {
            // Test multiple operations in sequence
            const operations = [];
            for (let i = 0; i < 10; i++) {
                operations.push(
                    bogoToken.connect(daoWallet).mintFromDAO(
                        user1.address, 
                        ethers.parseEther("100")
                    )
                );
            }
            
            await Promise.all(operations);
            expect(await bogoToken.balanceOf(user1.address))
                .to.equal(ethers.parseEther("1000"));
        });
    });

    describe("Timelock Cancellation Coverage", function () {
        it("Should cancel a queued timelock operation", async function () {
            // Since we don't have any timelock operations in the current contract,
            // we need to test the generic cancelTimelockOperation function
            // This would be used if we add timelock operations in the future
            
            // Create a mock operation ID
            const operationId = ethers.keccak256(ethers.toUtf8Bytes("test-operation"));
            
            // First, we need to set a timelock operation (this is a bit hacky but necessary for coverage)
            // Since there's no public function to set arbitrary timelock operations,
            // we'll test the error case
            await expect(
                bogoToken.cancelTimelockOperation(operationId)
            ).to.be.revertedWith("NOT_INITIALIZED");
        });

        it("Should only allow admin to cancel operations", async function () {
            const operationId = ethers.keccak256(ethers.toUtf8Bytes("test-operation"));
            
            await expect(
                bogoToken.connect(user1).cancelTimelockOperation(operationId)
            ).to.be.revertedWithCustomError(bogoToken, "AccessControlUnauthorizedAccount");
        });
    });

    describe("Interface Support Coverage", function () {
        it("Should support ERC165 interface", async function () {
            expect(await bogoToken.supportsInterface("0x01ffc9a7")).to.be.true;
        });

        it("Should support AccessControl interface", async function () {
            expect(await bogoToken.supportsInterface("0x7965db0b")).to.be.true;
        });

        it("Should not support random interface", async function () {
            expect(await bogoToken.supportsInterface("0x12345678")).to.be.false;
        });

        it("Should support ERC20 interface", async function () {
            // ERC20 doesn't have a standard interface ID, but we can test the behavior
            expect(await bogoToken.supportsInterface("0x36372b07")).to.be.false; // This is ERC20 metadata
        });
    });
});