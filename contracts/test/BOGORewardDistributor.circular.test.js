const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("BOGORewardDistributor - Circular Referral Prevention", function () {
    let bogoToken;
    let rewardDistributor;
    let treasury;
    let backend;
    let user1, user2, user3, user4, user5;
    let users;

    beforeEach(async function () {
        [treasury, backend, user1, user2, user3, user4, user5, ...users] = await ethers.getSigners();

        // Deploy BOGO token
        const BOGOToken = await ethers.getContractFactory("BOGOTokenV2");
        bogoToken = await BOGOToken.deploy();
        await bogoToken.deployed();

        // Deploy reward distributor
        const BOGORewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
        rewardDistributor = await BOGORewardDistributor.deploy(bogoToken.address, treasury.address);
        await rewardDistributor.deployed();

        // Setup
        await rewardDistributor.connect(treasury).setAuthorizedBackend(backend.address, true);
        
        // Grant DAO role to treasury and mint tokens for rewards
        await bogoToken.grantRole(await bogoToken.DAO_ROLE(), treasury.address);
        await bogoToken.connect(treasury).mintFromRewards(rewardDistributor.address, ethers.utils.parseEther("1000000"));
    });

    describe("Circular Referral Detection", function () {
        it("Should prevent direct circular referral (A→B, B→A)", async function () {
            // User1 refers User2
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            expect(await rewardDistributor.referredBy(user2.address)).to.equal(user1.address);

            // User2 tries to refer User1 (circular)
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(user2.address)
            ).to.be.revertedWith("Circular referral detected");
        });

        it("Should prevent 3-way circular referral (A→B→C→A)", async function () {
            // Create chain: User1 → User2 → User3
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            await rewardDistributor.connect(user3).claimReferralBonus(user2.address);

            // User3 tries to refer User1 (would create circle)
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(user3.address)
            ).to.be.revertedWith("Circular referral detected");
        });

        it("Should prevent complex circular referral (A→B→C→D→A)", async function () {
            // Create chain: User1 → User2 → User3 → User4
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            await rewardDistributor.connect(user3).claimReferralBonus(user2.address);
            await rewardDistributor.connect(user4).claimReferralBonus(user3.address);

            // User4 tries to refer User1 (would create circle)
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(user4.address)
            ).to.be.revertedWith("Circular referral detected");
        });

        it("Should allow valid non-circular referrals", async function () {
            // Create valid chain: User1 → User2 → User3
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            await rewardDistributor.connect(user3).claimReferralBonus(user2.address);

            // User4 can refer User3 (no circle)
            await expect(
                rewardDistributor.connect(user4).claimReferralBonus(user3.address)
            ).to.not.be.reverted;

            expect(await rewardDistributor.referredBy(user4.address)).to.equal(user3.address);
        });
    });

    describe("Referral Depth Tracking", function () {
        it("Should track referral depth correctly", async function () {
            // User1 has depth 0 (no referrer)
            expect(await rewardDistributor.referralDepth(user1.address)).to.equal(0);

            // User2 referred by User1 (depth 1)
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            expect(await rewardDistributor.referralDepth(user2.address)).to.equal(1);

            // User3 referred by User2 (depth 2)
            await rewardDistributor.connect(user3).claimReferralBonus(user2.address);
            expect(await rewardDistributor.referralDepth(user3.address)).to.equal(2);

            // User4 referred by User3 (depth 3)
            await rewardDistributor.connect(user4).claimReferralBonus(user3.address);
            expect(await rewardDistributor.referralDepth(user4.address)).to.equal(3);
        });

        it("Should prevent referrals beyond maximum depth", async function () {
            const MAX_DEPTH = await rewardDistributor.MAX_REFERRAL_DEPTH();
            
            // Create a chain up to MAX_DEPTH
            let prevUser = user1;
            for (let i = 0; i < MAX_DEPTH; i++) {
                const currentUser = users[i];
                await rewardDistributor.connect(currentUser).claimReferralBonus(prevUser.address);
                expect(await rewardDistributor.referralDepth(currentUser.address)).to.equal(i + 1);
                prevUser = currentUser;
            }

            // Try to add one more level (should fail)
            const extraUser = users[MAX_DEPTH];
            await expect(
                rewardDistributor.connect(extraUser).claimReferralBonus(prevUser.address)
            ).to.be.revertedWith("Max referral depth exceeded");
        });
    });

    describe("Get Referral Chain", function () {
        beforeEach(async function () {
            // Create chain: User1 → User2 → User3 → User4
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            await rewardDistributor.connect(user3).claimReferralBonus(user2.address);
            await rewardDistributor.connect(user4).claimReferralBonus(user3.address);
        });

        it("Should return correct referral chain for users", async function () {
            // User4's chain should be [User3, User2, User1]
            const chain4 = await rewardDistributor.getReferralChain(user4.address);
            expect(chain4.length).to.equal(3);
            expect(chain4[0]).to.equal(user3.address);
            expect(chain4[1]).to.equal(user2.address);
            expect(chain4[2]).to.equal(user1.address);

            // User3's chain should be [User2, User1]
            const chain3 = await rewardDistributor.getReferralChain(user3.address);
            expect(chain3.length).to.equal(2);
            expect(chain3[0]).to.equal(user2.address);
            expect(chain3[1]).to.equal(user1.address);

            // User2's chain should be [User1]
            const chain2 = await rewardDistributor.getReferralChain(user2.address);
            expect(chain2.length).to.equal(1);
            expect(chain2[0]).to.equal(user1.address);

            // User1's chain should be empty
            const chain1 = await rewardDistributor.getReferralChain(user1.address);
            expect(chain1.length).to.equal(0);
        });

        it("Should handle users with no referrals", async function () {
            const chain = await rewardDistributor.getReferralChain(user5.address);
            expect(chain.length).to.equal(0);
        });
    });

    describe("Edge Cases", function () {
        it("Should prevent self-referral", async function () {
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(user1.address)
            ).to.be.revertedWith("Cannot refer yourself");
        });

        it("Should prevent referring zero address", async function () {
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(ethers.constants.AddressZero)
            ).to.be.revertedWith("Invalid referrer");
        });

        it("Should prevent double referral", async function () {
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            
            await expect(
                rewardDistributor.connect(user2).claimReferralBonus(user3.address)
            ).to.be.revertedWith("Already referred");
        });

        it("Should handle circular detection at maximum depth", async function () {
            // Create a long chain close to MAX_DEPTH
            const MAX_DEPTH = await rewardDistributor.MAX_REFERRAL_DEPTH();
            let prevUser = user1;
            
            // Create chain of MAX_DEPTH - 1 users
            for (let i = 0; i < MAX_DEPTH - 1; i++) {
                const currentUser = users[i];
                await rewardDistributor.connect(currentUser).claimReferralBonus(prevUser.address);
                prevUser = currentUser;
            }

            // Last user in chain tries to refer the first user (should detect circle)
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(prevUser.address)
            ).to.be.revertedWith("Circular referral detected");
        });
    });

    describe("Referral Rewards with Circular Prevention", function () {
        it("Should correctly distribute rewards in valid referral chains", async function () {
            const referralBonus = ethers.utils.parseEther("20"); // From template
            
            const initialBalance1 = await bogoToken.balanceOf(user1.address);
            const initialBalance2 = await bogoToken.balanceOf(user2.address);

            // User2 refers User1
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            expect(await bogoToken.balanceOf(user1.address)).to.equal(initialBalance1.add(referralBonus));

            // User3 refers User2
            await rewardDistributor.connect(user3).claimReferralBonus(user2.address);
            expect(await bogoToken.balanceOf(user2.address)).to.equal(initialBalance2.add(referralBonus));

            // Verify referral counts
            expect(await rewardDistributor.referralCount(user1.address)).to.equal(1);
            expect(await rewardDistributor.referralCount(user2.address)).to.equal(1);
        });

        it("Should track referral count correctly even with prevented circular attempts", async function () {
            // User2 refers User1
            await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
            expect(await rewardDistributor.referralCount(user1.address)).to.equal(1);

            // User1 tries to refer User2 (circular - should fail)
            await expect(
                rewardDistributor.connect(user1).claimReferralBonus(user2.address)
            ).to.be.revertedWith("Circular referral detected");

            // Referral count should not increase for failed attempt
            expect(await rewardDistributor.referralCount(user2.address)).to.equal(0);

            // User3 successfully refers User1
            await rewardDistributor.connect(user3).claimReferralBonus(user1.address);
            expect(await rewardDistributor.referralCount(user1.address)).to.equal(2);
        });
    });
});