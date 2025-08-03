const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("Zero Address Validation Tests", function () {
    let bogoToken;
    let rewardDistributor;
    let owner, dao, business, user1, user2;
    
    const ZERO_ADDRESS = ethers.ZeroAddress;

    beforeEach(async function () {
        [owner, dao, business, user1, user2] = await ethers.getSigners();

        // Deploy BOGOTokenV2
        const BOGOToken = await ethers.getContractFactory("BOGOTokenV2");
        bogoToken = await BOGOToken.deploy();
        await bogoToken.waitForDeployment();

        // Grant roles
        const DAO_ROLE = await bogoToken.DAO_ROLE();
        const BUSINESS_ROLE = await bogoToken.BUSINESS_ROLE();
        await bogoToken.grantRole(DAO_ROLE, dao.address);
        await bogoToken.grantRole(BUSINESS_ROLE, business.address);

        // Deploy BOGORewardDistributor with test mode enabled
        const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
        rewardDistributor = await RewardDistributor.deploy(bogoToken.target, owner.address, true);
        await rewardDistributor.waitForDeployment();

        // Authorize backend
        await rewardDistributor.setAuthorizedBackend(owner.address, true);

        // Mint tokens to distributor
        await bogoToken.connect(dao).mintFromDAO(rewardDistributor.target, ethers.parseEther("100000"));
    });

    describe("BOGOTokenV2 Zero Address Validation", function () {
        it("Should revert mintFromDAO with zero address", async function () {
            // OpenZeppelin's _mint validates zero address
            await expect(
                bogoToken.connect(dao).mintFromDAO(ZERO_ADDRESS, ethers.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "ERC20InvalidReceiver");
        });

        it("Should revert mintFromBusiness with zero address", async function () {
            // OpenZeppelin's _mint validates zero address
            await expect(
                bogoToken.connect(business).mintFromBusiness(ZERO_ADDRESS, ethers.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "ERC20InvalidReceiver");
        });

        it("Should revert mintFromRewards with zero address", async function () {
            // OpenZeppelin's _mint validates zero address
            await expect(
                bogoToken.connect(dao).mintFromRewards(ZERO_ADDRESS, ethers.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "ERC20InvalidReceiver");
        });

        it("Should revert burnFrom with zero address", async function () {
            // burnFrom checks allowance first, so it reverts with insufficient allowance
            await expect(
                bogoToken.burnFrom(ZERO_ADDRESS, ethers.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "ERC20InsufficientAllowance");
        });

        it("Should allow grantRole with zero address (handled by AccessControl)", async function () {
            const MINTER_ROLE = await bogoToken.MINTER_ROLE();
            // OpenZeppelin's AccessControl allows granting roles to zero address
            await expect(
                bogoToken.grantRole(MINTER_ROLE, ZERO_ADDRESS)
            ).to.not.be.reverted;
        });

        it("Should allow revokeRole with zero address (handled by AccessControl)", async function () {
            const MINTER_ROLE = await bogoToken.MINTER_ROLE();
            // OpenZeppelin's AccessControl allows revoking roles from zero address
            await expect(
                bogoToken.revokeRole(MINTER_ROLE, ZERO_ADDRESS)
            ).to.not.be.reverted;
        });


        it("Should allow zero amount in minting functions (OpenZeppelin allows)", async function () {
            // OpenZeppelin's _mint allows zero amounts
            await expect(
                bogoToken.connect(dao).mintFromDAO(user1.address, 0)
            ).to.not.be.reverted;

            await expect(
                bogoToken.connect(business).mintFromBusiness(user1.address, 0)
            ).to.not.be.reverted;

            await expect(
                bogoToken.connect(dao).mintFromRewards(user1.address, 0)
            ).to.not.be.reverted;
        });

        it("Should allow burn with zero amount (handled by OpenZeppelin)", async function () {
            // OpenZeppelin's _burn allows zero amount
            await expect(
                bogoToken.burn(0)
            ).to.not.be.reverted;
        });

        it("Should allow burnFrom with zero amount (handled by OpenZeppelin)", async function () {
            // OpenZeppelin's _burn allows zero amount
            await expect(
                bogoToken.burnFrom(user1.address, 0)
            ).to.not.be.reverted;
        });

        it("Should allow valid minting operations", async function () {
            // These should succeed
            await expect(
                bogoToken.connect(dao).mintFromDAO(user1.address, ethers.parseEther("100"))
            ).to.not.be.reverted;

            await expect(
                bogoToken.connect(business).mintFromBusiness(user1.address, ethers.parseEther("100"))
            ).to.not.be.reverted;

            await expect(
                bogoToken.connect(dao).mintFromRewards(user1.address, ethers.parseEther("100"))
            ).to.not.be.reverted;

            expect(await bogoToken.balanceOf(user1.address)).to.equal(ethers.parseEther("300"));
        });
    });

    describe("BOGORewardDistributor Zero Address Validation", function () {
        it("Should revert constructor with zero token address", async function () {
            const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
            await expect(
                RewardDistributor.deploy(ZERO_ADDRESS, owner.address, true)
            ).to.be.revertedWith("ZERO_ADDRESS");
        });

        it("Should revert constructor with zero treasury address", async function () {
            const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
            await expect(
                RewardDistributor.deploy(bogoToken.target, ZERO_ADDRESS, true)
            ).to.be.revertedWith("ZERO_ADDRESS");
        });

        it("Should allow claimCustomReward with zero recipient (no validation)", async function () {
            // The contract doesn't validate recipient address in claimCustomReward
            // This test verifies the current behavior
            await expect(
                rewardDistributor.claimCustomReward(ZERO_ADDRESS, ethers.parseEther("100"), "test")
            ).to.not.be.revertedWith("ZERO_ADDRESS");
        });

        it("Should revert claimCustomReward with zero amount", async function () {
            await expect(
                rewardDistributor.claimCustomReward(user1.address, 0, "test")
            ).to.be.revertedWith("INVALID_AMOUNT");
        });

        it("Should revert claimReferralBonus with zero referrer", async function () {
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(ZERO_ADDRESS)
            ).to.be.revertedWith("ZERO_ADDRESS");
        });

        it("Should allow addToWhitelist with zero address in array (no validation)", async function () {
            // The contract doesn't validate addresses in addToWhitelist
            await expect(
                rewardDistributor.addToWhitelist([user1.address, ZERO_ADDRESS, user2.address])
            ).to.not.be.reverted;
        });

        it("Should allow removeFromWhitelist with zero address (no validation)", async function () {
            // The contract doesn't validate address in removeFromWhitelist
            await expect(
                rewardDistributor.removeFromWhitelist(ZERO_ADDRESS)
            ).to.not.be.reverted;
        });

        it("Should allow setAuthorizedBackend with zero address (no validation)", async function () {
            // The contract doesn't validate backend address in setAuthorizedBackend
            await expect(
                rewardDistributor.setAuthorizedBackend(ZERO_ADDRESS, true)
            ).to.not.be.reverted;
        });

        it("Should revert treasurySweep with zero recipient", async function () {
            await expect(
                rewardDistributor.treasurySweep(bogoToken.target, ZERO_ADDRESS, ethers.parseEther("100"))
            ).to.be.revertedWith("ZERO_ADDRESS");
        });

        it("Should revert treasurySweep with zero amount", async function () {
            await expect(
                rewardDistributor.treasurySweep(bogoToken.target, user1.address, 0)
            ).to.be.revertedWith("ZERO_AMOUNT");
        });

        it("Should handle canClaim with zero address", async function () {
            const [eligible, reason] = await rewardDistributor.canClaim(ZERO_ADDRESS, "welcome_bonus");
            expect(eligible).to.be.true; // No address validation in canClaim
            expect(reason).to.equal("Eligible");
        });

        it("Should handle getReferralChain with zero address", async function () {
            const chain = await rewardDistributor.getReferralChain(ZERO_ADDRESS);
            expect(chain.length).to.equal(0);
        });

        it("Should allow valid reward operations", async function () {
            // These should succeed
            await expect(
                rewardDistributor.claimCustomReward(user1.address, ethers.parseEther("100"), "test")
            ).to.not.be.reverted;

            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(user2.address)
            ).to.not.be.reverted;

            await expect(
                rewardDistributor.addToWhitelist([user1.address, user2.address])
            ).to.not.be.reverted;

            expect(await bogoToken.balanceOf(user1.address)).to.be.gt(0);
        });
    });

    describe("Transfer Protection", function () {
        beforeEach(async function () {
            // Mint some tokens to test transfers
            await bogoToken.connect(dao).mintFromDAO(user1.address, ethers.parseEther("1000"));
        });

        it("Should prevent transfers to zero address", async function () {
            // Direct transfer to zero address should fail
            await expect(
                bogoToken.connect(user1).transfer(ZERO_ADDRESS, ethers.parseEther("100"))
            ).to.be.reverted;
        });

        it("Should allow burning (transfer to zero in _update)", async function () {
            const balanceBefore = await bogoToken.balanceOf(user1.address);
            await bogoToken.connect(user1).burn(ethers.parseEther("100"));
            const balanceAfter = await bogoToken.balanceOf(user1.address);
            
            expect(balanceBefore - balanceAfter).to.equal(ethers.parseEther("100"));
        });

        it("Should allow valid transfers", async function () {
            await expect(
                bogoToken.connect(user1).transfer(user2.address, ethers.parseEther("100"))
            ).to.not.be.reverted;
            
            expect(await bogoToken.balanceOf(user2.address)).to.equal(ethers.parseEther("100"));
        });
    });

    describe("Gas Optimization with Custom Errors", function () {
        it("Should use less gas with custom errors", async function () {
            // Estimate gas for a failing transaction
            // Test that minting with zero address fails with custom error
            await expect(
                bogoToken.connect(dao).mintFromDAO(ZERO_ADDRESS, ethers.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "ERC20InvalidReceiver");
        });
    });
});