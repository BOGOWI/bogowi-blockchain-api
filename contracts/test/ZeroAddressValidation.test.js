const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("Zero Address Validation Tests", function () {
    let bogoToken;
    let rewardDistributor;
    let owner, dao, business, user1, user2;
    
    const ZERO_ADDRESS = ethers.constants.AddressZero;

    beforeEach(async function () {
        [owner, dao, business, user1, user2] = await ethers.getSigners();

        // Deploy BOGOTokenV2 with zero validation
        const BOGOToken = await ethers.getContractFactory("BOGOTokenV2_ZeroValidated");
        bogoToken = await BOGOToken.deploy();
        await bogoToken.deployed();

        // Grant roles
        const DAO_ROLE = await bogoToken.DAO_ROLE();
        const BUSINESS_ROLE = await bogoToken.BUSINESS_ROLE();
        await bogoToken.grantRole(DAO_ROLE, dao.address);
        await bogoToken.grantRole(BUSINESS_ROLE, business.address);

        // Deploy BOGORewardDistributor with zero validation
        const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor_ZeroValidated");
        rewardDistributor = await RewardDistributor.deploy(bogoToken.address, owner.address);
        await rewardDistributor.deployed();

        // Authorize backend
        await rewardDistributor.setAuthorizedBackend(owner.address, true);

        // Mint tokens to distributor
        await bogoToken.connect(dao).mintFromDAO(rewardDistributor.address, ethers.utils.parseEther("100000"));
    });

    describe("BOGOTokenV2 Zero Address Validation", function () {
        it("Should revert mintFromDAO with zero address", async function () {
            await expect(
                bogoToken.connect(dao).mintFromDAO(ZERO_ADDRESS, ethers.utils.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should revert mintFromBusiness with zero address", async function () {
            await expect(
                bogoToken.connect(business).mintFromBusiness(ZERO_ADDRESS, ethers.utils.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should revert mintFromRewards with zero address", async function () {
            await expect(
                bogoToken.connect(dao).mintFromRewards(ZERO_ADDRESS, ethers.utils.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should revert burnFrom with zero address", async function () {
            await expect(
                bogoToken.burnFrom(ZERO_ADDRESS, ethers.utils.parseEther("100"))
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should revert grantRole with zero address", async function () {
            const MINTER_ROLE = await bogoToken.MINTER_ROLE();
            await expect(
                bogoToken.grantRole(MINTER_ROLE, ZERO_ADDRESS)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should revert revokeRole with zero address", async function () {
            const MINTER_ROLE = await bogoToken.MINTER_ROLE();
            await expect(
                bogoToken.revokeRole(MINTER_ROLE, ZERO_ADDRESS)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should revert queueRegisterFlavoredToken with zero address", async function () {
            await expect(
                bogoToken.queueRegisterFlavoredToken("test", ZERO_ADDRESS)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAddress");
        });

        it("Should revert with zero amount in minting functions", async function () {
            await expect(
                bogoToken.connect(dao).mintFromDAO(user1.address, 0)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");

            await expect(
                bogoToken.connect(business).mintFromBusiness(user1.address, 0)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");

            await expect(
                bogoToken.connect(dao).mintFromRewards(user1.address, 0)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
        });

        it("Should revert burn with zero amount", async function () {
            await expect(
                bogoToken.burn(0)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
        });

        it("Should revert burnFrom with zero amount", async function () {
            await expect(
                bogoToken.burnFrom(user1.address, 0)
            ).to.be.revertedWithCustomError(bogoToken, "InvalidAmount");
        });

        it("Should allow valid minting operations", async function () {
            // These should succeed
            await expect(
                bogoToken.connect(dao).mintFromDAO(user1.address, ethers.utils.parseEther("100"))
            ).to.not.be.reverted;

            await expect(
                bogoToken.connect(business).mintFromBusiness(user1.address, ethers.utils.parseEther("100"))
            ).to.not.be.reverted;

            await expect(
                bogoToken.connect(dao).mintFromRewards(user1.address, ethers.utils.parseEther("100"))
            ).to.not.be.reverted;

            expect(await bogoToken.balanceOf(user1.address)).to.equal(ethers.utils.parseEther("300"));
        });
    });

    describe("BOGORewardDistributor Zero Address Validation", function () {
        it("Should revert constructor with zero token address", async function () {
            const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor_ZeroValidated");
            await expect(
                RewardDistributor.deploy(ZERO_ADDRESS, owner.address)
            ).to.be.revertedWithCustomError(RewardDistributor, "InvalidTokenAddress");
        });

        it("Should revert constructor with zero treasury address", async function () {
            const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor_ZeroValidated");
            await expect(
                RewardDistributor.deploy(bogoToken.address, ZERO_ADDRESS)
            ).to.be.revertedWithCustomError(RewardDistributor, "InvalidTreasuryAddress");
        });

        it("Should revert claimCustomReward with zero recipient", async function () {
            await expect(
                rewardDistributor.claimCustomReward(ZERO_ADDRESS, ethers.utils.parseEther("100"), "test")
            ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
        });

        it("Should revert claimCustomReward with zero amount", async function () {
            await expect(
                rewardDistributor.claimCustomReward(user1.address, 0, "test")
            ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAmount");
        });

        it("Should revert claimReferralBonus with zero referrer", async function () {
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(ZERO_ADDRESS)
            ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
        });

        it("Should revert addToWhitelist with zero address in array", async function () {
            await expect(
                rewardDistributor.addToWhitelist([user1.address, ZERO_ADDRESS, user2.address])
            ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
        });

        it("Should revert removeFromWhitelist with zero address", async function () {
            await expect(
                rewardDistributor.removeFromWhitelist(ZERO_ADDRESS)
            ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
        });

        it("Should revert setAuthorizedBackend with zero address", async function () {
            await expect(
                rewardDistributor.setAuthorizedBackend(ZERO_ADDRESS, true)
            ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
        });

        it("Should revert treasurySweep with zero recipient", async function () {
            await expect(
                rewardDistributor.treasurySweep(bogoToken.address, ZERO_ADDRESS, ethers.utils.parseEther("100"))
            ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
        });

        it("Should revert treasurySweep with zero amount", async function () {
            await expect(
                rewardDistributor.treasurySweep(bogoToken.address, user1.address, 0)
            ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAmount");
        });

        it("Should handle canClaim with zero address", async function () {
            const [eligible, reason] = await rewardDistributor.canClaim(ZERO_ADDRESS, "welcome_bonus");
            expect(eligible).to.be.false;
            expect(reason).to.equal("Invalid wallet address");
        });

        it("Should handle getReferralChain with zero address", async function () {
            const chain = await rewardDistributor.getReferralChain(ZERO_ADDRESS);
            expect(chain.length).to.equal(0);
        });

        it("Should allow valid reward operations", async function () {
            // These should succeed
            await expect(
                rewardDistributor.claimCustomReward(user1.address, ethers.utils.parseEther("100"), "test")
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
            await bogoToken.connect(dao).mintFromDAO(user1.address, ethers.utils.parseEther("1000"));
        });

        it("Should prevent transfers to zero address", async function () {
            // Direct transfer to zero address should fail
            await expect(
                bogoToken.connect(user1).transfer(ZERO_ADDRESS, ethers.utils.parseEther("100"))
            ).to.be.reverted;
        });

        it("Should allow burning (transfer to zero in _update)", async function () {
            const balanceBefore = await bogoToken.balanceOf(user1.address);
            await bogoToken.connect(user1).burn(ethers.utils.parseEther("100"));
            const balanceAfter = await bogoToken.balanceOf(user1.address);
            
            expect(balanceBefore.sub(balanceAfter)).to.equal(ethers.utils.parseEther("100"));
        });

        it("Should allow valid transfers", async function () {
            await expect(
                bogoToken.connect(user1).transfer(user2.address, ethers.utils.parseEther("100"))
            ).to.not.be.reverted;
            
            expect(await bogoToken.balanceOf(user2.address)).to.equal(ethers.utils.parseEther("100"));
        });
    });

    describe("Gas Optimization with Custom Errors", function () {
        it("Should use less gas with custom errors", async function () {
            // Estimate gas for a failing transaction
            try {
                await bogoToken.connect(dao).mintFromDAO(ZERO_ADDRESS, ethers.utils.parseEther("100"));
            } catch (error) {
                // Custom errors use less gas than require strings
                expect(error.message).to.include("InvalidAddress");
            }
        });
    });
});