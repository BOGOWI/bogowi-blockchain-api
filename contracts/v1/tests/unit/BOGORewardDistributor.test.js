const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGORewardDistributor", function () {
    let bogoToken;
    let roleManager;
    let distributor;
    let owner;
    let treasury;
    let backend;
    let user1;
    let user2;
    let user3;
    let referrer;
    let addrs;

    beforeEach(async function () {
        [owner, treasury, backend, user1, user2, user3, referrer, ...addrs] = await ethers.getSigners();

        // Deploy RoleManager first
        const RoleManager = await ethers.getContractFactory("RoleManager");
        roleManager = await RoleManager.deploy();
        await roleManager.waitForDeployment();

        // Deploy BOGOToken with roleManager
        const BOGOToken = await ethers.getContractFactory("BOGOToken");
        bogoToken = await BOGOToken.deploy(await roleManager.getAddress(), "BOGO Token", "BOGO");
        await bogoToken.waitForDeployment();

        // Deploy BOGORewardDistributor
        const BOGORewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
        distributor = await BOGORewardDistributor.deploy(await roleManager.getAddress(), await bogoToken.getAddress());
        await distributor.waitForDeployment();

        // Register contracts in RoleManager
        await roleManager.registerContract(await bogoToken.getAddress(), "BOGOToken");
        await roleManager.registerContract(await distributor.getAddress(), "BOGORewardDistributor");

        // Setup roles
        const TREASURY_ROLE = await roleManager.TREASURY_ROLE();
        const DISTRIBUTOR_BACKEND_ROLE = await roleManager.DISTRIBUTOR_BACKEND_ROLE();
        const BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();

        await roleManager.grantRole(TREASURY_ROLE, treasury.address);
        await roleManager.grantRole(DISTRIBUTOR_BACKEND_ROLE, backend.address);
        await roleManager.grantRole(BUSINESS_ROLE, owner.address);

        // Mint tokens to owner from business allocation
        await bogoToken.mintFromBusiness(owner.address, ethers.parseEther("10000000"));
        
        // Transfer tokens to distributor
        await bogoToken.transfer(await distributor.getAddress(), ethers.parseEther("1000000"));
    });

    describe("Constructor", function () {
        it("Should revert with zero token address", async function () {
            const BOGORewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
            await expect(
                BOGORewardDistributor.deploy(await roleManager.getAddress(), ethers.ZeroAddress)
            ).to.be.revertedWithCustomError(distributor, "InvalidTokenAddress");
        });

        it("Should initialize with correct parameters", async function () {
            expect(await distributor.bogoToken()).to.equal(await bogoToken.getAddress());
            expect(await distributor.roleManager()).to.equal(await roleManager.getAddress());
        });

        it("Should initialize all reward templates", async function () {
            const welcomeBonus = await distributor.templates("welcome_bonus");
            expect(welcomeBonus.fixedAmount).to.equal(ethers.parseEther("10"));
            expect(welcomeBonus.active).to.be.true;

            const founderBonus = await distributor.templates("founder_bonus");
            expect(founderBonus.fixedAmount).to.equal(ethers.parseEther("100"));
            expect(founderBonus.requiresWhitelist).to.be.true;
        });
    });

    describe("Claim Reward", function () {
        it("Should allow claiming welcome bonus", async function () {
            await expect(distributor.connect(user1).claimReward("welcome_bonus"))
                .to.emit(distributor, "RewardClaimed")
                .withArgs(user1.address, "welcome_bonus", ethers.parseEther("10"));
        });

        it("Should prevent claiming inactive template", async function () {
            await distributor.connect(treasury).updateTemplate("welcome_bonus", {
                id: "welcome_bonus",
                fixedAmount: ethers.parseEther("10"),
                maxAmount: 0,
                cooldownPeriod: 0,
                maxClaimsPerWallet: 1,
                requiresWhitelist: false,
                active: false
            });

            await expect(
                distributor.connect(user1).claimReward("welcome_bonus")
            ).to.be.revertedWithCustomError(distributor, "TemplateNotActive");
        });

        it("Should enforce max claims per wallet", async function () {
            await distributor.connect(user1).claimReward("welcome_bonus");
            
            await expect(
                distributor.connect(user1).claimReward("welcome_bonus")
            ).to.be.revertedWithCustomError(distributor, "MaxClaimsReached");
        });

        it("Should enforce cooldown period", async function () {
            await distributor.connect(user1).claimReward("dao_participation");
            
            await expect(
                distributor.connect(user1).claimReward("dao_participation")
            ).to.be.revertedWithCustomError(distributor, "CooldownActive");

            // Fast forward time
            await time.increase(30 * 24 * 60 * 60); // 30 days

            // Should work after cooldown
            await expect(distributor.connect(user1).claimReward("dao_participation"))
                .to.emit(distributor, "RewardClaimed");
        });

        it("Should enforce whitelist requirement", async function () {
            await expect(
                distributor.connect(user1).claimReward("founder_bonus")
            ).to.be.revertedWithCustomError(distributor, "NotWhitelisted");

            // Add to whitelist
            await distributor.connect(treasury).addToWhitelist([user1.address]);

            // Should work after whitelisting
            await expect(distributor.connect(user1).claimReward("founder_bonus"))
                .to.emit(distributor, "RewardClaimed");
        });

        it("Should enforce daily limit", async function () {
            // Update a template to have a large reward
            await distributor.connect(treasury).updateTemplate("test_large", {
                id: "test_large",
                fixedAmount: ethers.parseEther("600000"), // More than daily limit
                maxAmount: 0,
                cooldownPeriod: 0,
                maxClaimsPerWallet: 1,
                requiresWhitelist: false,
                active: true
            });

            await expect(
                distributor.connect(user1).claimReward("test_large")
            ).to.be.revertedWithCustomError(distributor, "DailyLimitExceeded");
        });

        it("Should revert when claiming reward with zero fixed amount", async function () {
            // Create a template with zero fixed amount (like custom_reward)
            await distributor.connect(treasury).updateTemplate("zero_amount", {
                id: "zero_amount",
                fixedAmount: 0,
                maxAmount: ethers.parseEther("100"),
                cooldownPeriod: 0,
                maxClaimsPerWallet: 0,
                requiresWhitelist: false,
                active: true
            });

            await expect(
                distributor.connect(user1).claimReward("zero_amount")
            ).to.be.revertedWithCustomError(distributor, "InvalidTemplateAmount");
        });
    });

    describe("Custom Rewards", function () {
        it("Should allow backend to distribute custom rewards", async function () {
            const amount = ethers.parseEther("50");
            await expect(
                distributor.connect(backend).claimCustomReward(user1.address, amount, "Test reward")
            )
                .to.emit(distributor, "RewardClaimed")
                .withArgs(user1.address, "Test reward", amount);
        });

        it("Should reject custom rewards from non-backend", async function () {
            await expect(
                distributor.connect(user1).claimCustomReward(user2.address, ethers.parseEther("50"), "Test")
            ).to.be.revertedWithCustomError(distributor, "NotAuthorizedBackend");
        });

        it("Should reject zero address recipient", async function () {
            await expect(
                distributor.connect(backend).claimCustomReward(ethers.ZeroAddress, ethers.parseEther("50"), "Test")
            ).to.be.revertedWithCustomError(distributor, "InvalidAddress");
        });

        it("Should reject zero amount", async function () {
            await expect(
                distributor.connect(backend).claimCustomReward(user1.address, 0, "Test")
            ).to.be.revertedWithCustomError(distributor, "InvalidAmount");
        });

        it("Should enforce max custom reward amount", async function () {
            await expect(
                distributor.connect(backend).claimCustomReward(user1.address, ethers.parseEther("1001"), "Test")
            ).to.be.revertedWithCustomError(distributor, "InvalidAmount");
        });

        it("Should reject custom reward when template is inactive", async function () {
            // Deactivate custom_reward template
            await distributor.connect(treasury).updateTemplate("custom_reward", {
                id: "custom_reward",
                fixedAmount: 0,
                maxAmount: ethers.parseEther("1000"),
                cooldownPeriod: 0,
                maxClaimsPerWallet: 0,
                requiresWhitelist: false,
                active: false
            });

            await expect(
                distributor.connect(backend).claimCustomReward(user1.address, ethers.parseEther("100"), "Test")
            ).to.be.revertedWithCustomError(distributor, "TemplateNotActive");
        });

        it("Should enforce daily limit for custom rewards", async function () {
            // First increase the max custom reward amount
            await distributor.connect(treasury).updateTemplate("custom_reward", {
                id: "custom_reward",
                fixedAmount: 0,
                maxAmount: ethers.parseEther("500000"),
                cooldownPeriod: 0,
                maxClaimsPerWallet: 0,
                requiresWhitelist: false,
                active: true
            });

            // Use up most of the daily limit
            await distributor.connect(backend).claimCustomReward(user1.address, ethers.parseEther("499995"), "Large reward");

            // Try to exceed daily limit
            await expect(
                distributor.connect(backend).claimCustomReward(user2.address, ethers.parseEther("10"), "Exceeds limit")
            ).to.be.revertedWithCustomError(distributor, "DailyLimitExceeded");
        });
    });

    describe("Referral System", function () {
        it("Should allow claiming referral bonus", async function () {
            await expect(distributor.connect(user1).claimReferralBonus(referrer.address))
                .to.emit(distributor, "ReferralClaimed")
                .withArgs(referrer.address, user1.address, ethers.parseEther("20"));

            expect(await distributor.referredBy(user1.address)).to.equal(referrer.address);
            expect(await distributor.referralCount(referrer.address)).to.equal(1);
        });

        it("Should prevent self-referral", async function () {
            await expect(
                distributor.connect(user1).claimReferralBonus(user1.address)
            ).to.be.revertedWithCustomError(distributor, "SelfReferral");
        });

        it("Should prevent double referral", async function () {
            await distributor.connect(user1).claimReferralBonus(referrer.address);
            
            await expect(
                distributor.connect(user1).claimReferralBonus(user2.address)
            ).to.be.revertedWithCustomError(distributor, "AlreadyReferred");
        });

        it("Should prevent zero address referrer", async function () {
            await expect(
                distributor.connect(user1).claimReferralBonus(ethers.ZeroAddress)
            ).to.be.revertedWithCustomError(distributor, "InvalidAddress");
        });

        it("Should prevent circular referrals", async function () {
            // A refers B
            await distributor.connect(user2).claimReferralBonus(user1.address);
            
            // B cannot refer A
            await expect(
                distributor.connect(user1).claimReferralBonus(user2.address)
            ).to.be.revertedWithCustomError(distributor, "CircularReferral");
        });

        it("Should track referral depth", async function () {
            // Create a chain: referrer -> user1 -> user2
            await distributor.connect(user1).claimReferralBonus(referrer.address);
            await distributor.connect(user2).claimReferralBonus(user1.address);

            expect(await distributor.referralDepth(user1.address)).to.equal(1);
            expect(await distributor.referralDepth(user2.address)).to.equal(2);
        });

        it("Should enforce max referral depth", async function () {
            // Create a chain up to max depth
            let prev = referrer.address;
            for (let i = 0; i < 10; i++) {
                const user = addrs[i];
                await distributor.connect(user).claimReferralBonus(prev);
                prev = user.address;
            }

            // Try to exceed max depth
            await expect(
                distributor.connect(addrs[10]).claimReferralBonus(addrs[9].address)
            ).to.be.revertedWithCustomError(distributor, "MaxReferralDepthExceeded");
        });

        it("Should return referral chain correctly", async function () {
            // Create chain: referrer -> user1 -> user2 -> user3
            await distributor.connect(user1).claimReferralBonus(referrer.address);
            await distributor.connect(user2).claimReferralBonus(user1.address);
            await distributor.connect(user3).claimReferralBonus(user2.address);

            const chain = await distributor.getReferralChain(user3.address);
            expect(chain).to.have.lengthOf(3);
            expect(chain[0]).to.equal(user2.address);
            expect(chain[1]).to.equal(user1.address);
            expect(chain[2]).to.equal(referrer.address);
        });
    });

    describe("Whitelist Management", function () {
        it("Should allow treasury to add to whitelist", async function () {
            await expect(distributor.connect(treasury).addToWhitelist([user1.address, user2.address]))
                .to.emit(distributor, "WhitelistUpdated")
                .withArgs(user1.address, true)
                .to.emit(distributor, "WhitelistUpdated")
                .withArgs(user2.address, true);

            expect(await distributor.founderWhitelist(user1.address)).to.be.true;
            expect(await distributor.founderWhitelist(user2.address)).to.be.true;
        });

        it("Should reject zero address in whitelist", async function () {
            await expect(
                distributor.connect(treasury).addToWhitelist([ethers.ZeroAddress])
            ).to.be.revertedWithCustomError(distributor, "InvalidAddress");
        });

        it("Should allow treasury to remove from whitelist", async function () {
            await distributor.connect(treasury).addToWhitelist([user1.address]);
            
            await expect(distributor.connect(treasury).removeFromWhitelist(user1.address))
                .to.emit(distributor, "WhitelistUpdated")
                .withArgs(user1.address, false);

            expect(await distributor.founderWhitelist(user1.address)).to.be.false;
        });

        it("Should reject removing zero address from whitelist", async function () {
            await expect(
                distributor.connect(treasury).removeFromWhitelist(ethers.ZeroAddress)
            ).to.be.revertedWithCustomError(distributor, "InvalidAddress");
        });

        it("Should reject non-treasury whitelist management", async function () {
            await expect(
                distributor.connect(user1).addToWhitelist([user2.address])
            ).to.be.revertedWithCustomError(distributor, "UnauthorizedAccess");
        });
    });

    describe("Template Management", function () {
        it("Should allow treasury to update templates", async function () {
            const newTemplate = {
                id: "new_reward",
                fixedAmount: ethers.parseEther("5"),
                maxAmount: 0,
                cooldownPeriod: 60,
                maxClaimsPerWallet: 2,
                requiresWhitelist: false,
                active: true
            };

            await expect(distributor.connect(treasury).updateTemplate("new_reward", newTemplate))
                .to.emit(distributor, "TemplateUpdated")
                .withArgs("new_reward");

            const template = await distributor.templates("new_reward");
            expect(template.fixedAmount).to.equal(newTemplate.fixedAmount);
        });

        it("Should reject non-treasury template updates", async function () {
            await expect(
                distributor.connect(user1).updateTemplate("test", {
                    id: "test",
                    fixedAmount: 0,
                    maxAmount: 0,
                    cooldownPeriod: 0,
                    maxClaimsPerWallet: 0,
                    requiresWhitelist: false,
                    active: false
                })
            ).to.be.revertedWithCustomError(distributor, "UnauthorizedAccess");
        });
    });

    describe("Pause Functionality", function () {
        it("Should allow pauser to pause", async function () {
            const PAUSER_ROLE = await roleManager.PAUSER_ROLE();
            await roleManager.grantRole(PAUSER_ROLE, owner.address);

            await distributor.pause();
            expect(await distributor.paused()).to.be.true;
        });

        it("Should prevent operations when paused", async function () {
            const PAUSER_ROLE = await roleManager.PAUSER_ROLE();
            await roleManager.grantRole(PAUSER_ROLE, owner.address);
            await distributor.pause();

            await expect(
                distributor.connect(user1).claimReward("welcome_bonus")
            ).to.be.revertedWithCustomError(distributor, "EnforcedPause");
        });

        it("Should allow pauser to unpause", async function () {
            const PAUSER_ROLE = await roleManager.PAUSER_ROLE();
            await roleManager.grantRole(PAUSER_ROLE, owner.address);
            
            // First pause
            await distributor.pause();
            expect(await distributor.paused()).to.be.true;
            
            // Then unpause
            await distributor.unpause();
            expect(await distributor.paused()).to.be.false;
            
            // Operations should work again
            await expect(distributor.connect(user1).claimReward("welcome_bonus"))
                .to.emit(distributor, "RewardClaimed");
        });
    });

    describe("Daily Limit", function () {
        it("Should reset daily limit after 24 hours", async function () {
            // Update template to allow larger custom rewards temporarily for testing
            await distributor.connect(treasury).updateTemplate("custom_reward", {
                id: "custom_reward",
                fixedAmount: 0,
                maxAmount: ethers.parseEther("500000"), // Increase max for testing
                cooldownPeriod: 0,
                maxClaimsPerWallet: 0,
                requiresWhitelist: false,
                active: true
            });
            
            // Claim near the limit
            const almostLimit = ethers.parseEther("499995");
            await distributor.connect(backend).claimCustomReward(user1.address, almostLimit, "Near limit");
            
            // Next claim should fail (will exceed 500k limit) - 20 BOGO referral would exceed
            const newUser = (await ethers.getSigners())[15];
            await expect(
                distributor.connect(newUser).claimReferralBonus(referrer.address) // 20 BOGO
            ).to.be.revertedWithCustomError(distributor, "DailyLimitExceeded");

            // Fast forward 24 hours
            await time.increase(24 * 60 * 60);

            // Should work after reset
            await expect(
                distributor.connect(newUser).claimReferralBonus(referrer.address)
            ).to.emit(distributor, "ReferralClaimed");
        });

        it("Should return correct remaining daily limit", async function () {
            const initialLimit = await distributor.getRemainingDailyLimit();
            expect(initialLimit).to.equal(ethers.parseEther("500000"));

            // Claim some rewards
            await distributor.connect(user1).claimReward("welcome_bonus");
            
            const remainingLimit = await distributor.getRemainingDailyLimit();
            expect(remainingLimit).to.equal(ethers.parseEther("499990"));
        });
    });

    describe("Treasury Sweep", function () {
        it("Should allow treasury to sweep tokens", async function () {
            const sweepAmount = ethers.parseEther("100");
            
            await expect(
                distributor.connect(treasury).treasurySweep(await bogoToken.getAddress(), treasury.address, sweepAmount)
            )
                .to.emit(distributor, "TreasurySweep")
                .withArgs(await bogoToken.getAddress(), treasury.address, sweepAmount);
        });

        it("Should allow treasury to sweep ETH", async function () {
            // Send ETH to distributor
            await owner.sendTransaction({
                to: await distributor.getAddress(),
                value: ethers.parseEther("1")
            });

            const initialBalance = await ethers.provider.getBalance(treasury.address);
            
            await distributor.connect(treasury).treasurySweep(
                ethers.ZeroAddress,
                treasury.address,
                ethers.parseEther("1")
            );

            const finalBalance = await ethers.provider.getBalance(treasury.address);
            expect(finalBalance - initialBalance).to.be.closeTo(
                ethers.parseEther("1"),
                ethers.parseEther("0.01") // Account for gas
            );
        });

        it("Should reject sweep to zero address", async function () {
            await expect(
                distributor.connect(treasury).treasurySweep(
                    await bogoToken.getAddress(),
                    ethers.ZeroAddress,
                    ethers.parseEther("100")
                )
            ).to.be.revertedWithCustomError(distributor, "InvalidAddress");
        });

        it("Should reject sweep with zero amount", async function () {
            await expect(
                distributor.connect(treasury).treasurySweep(
                    await bogoToken.getAddress(),
                    treasury.address,
                    0
                )
            ).to.be.revertedWithCustomError(distributor, "InvalidAmount");
        });
    });

    describe("View Functions", function () {
        it("Should correctly report eligibility", async function () {
            // Check welcome bonus eligibility
            let [eligible, reason] = await distributor.canClaim(user1.address, "welcome_bonus");
            expect(eligible).to.be.true;
            expect(reason).to.equal("Eligible");

            // Claim and check again
            await distributor.connect(user1).claimReward("welcome_bonus");
            [eligible, reason] = await distributor.canClaim(user1.address, "welcome_bonus");
            expect(eligible).to.be.false;
            expect(reason).to.equal("Max claims reached");

            // Check with zero address
            [eligible, reason] = await distributor.canClaim(ethers.ZeroAddress, "welcome_bonus");
            expect(eligible).to.be.false;
            expect(reason).to.equal("Invalid wallet address");
        });

        it("Should report not whitelisted for founder bonus", async function () {
            // Check founder bonus without whitelist
            const [eligible, reason] = await distributor.canClaim(user1.address, "founder_bonus");
            expect(eligible).to.be.false;
            expect(reason).to.equal("Not whitelisted");
        });

        it("Should return full daily limit when 24 hours have passed", async function () {
            // First make a claim to use some limit
            await distributor.connect(user1).claimReward("welcome_bonus");
            
            // Check remaining is less than full
            let remaining = await distributor.getRemainingDailyLimit();
            expect(remaining).to.be.lt(ethers.parseEther("500000"));
            
            // Fast forward 24 hours
            await time.increase(24 * 60 * 60);
            
            // Check remaining is back to full limit
            remaining = await distributor.getRemainingDailyLimit();
            expect(remaining).to.equal(ethers.parseEther("500000"));
        });

        it("Should return empty array for zero address in getReferralChain", async function () {
            const chain = await distributor.getReferralChain(ethers.ZeroAddress);
            expect(chain).to.have.lengthOf(0);
        });

        it("Should report cooldown period active for DAO participation", async function () {
            // First claim DAO participation reward
            await distributor.connect(user1).claimReward("dao_participation");
            
            // Check immediately after - should be in cooldown
            const [eligible, reason] = await distributor.canClaim(user1.address, "dao_participation");
            expect(eligible).to.be.false;
            expect(reason).to.equal("Cooldown period active");
        });

        it("Should report template not active", async function () {
            // Deactivate a template
            await distributor.connect(treasury).updateTemplate("welcome_bonus", {
                id: "welcome_bonus",
                fixedAmount: ethers.parseEther("10"),
                maxAmount: 0,
                cooldownPeriod: 0,
                maxClaimsPerWallet: 1,
                requiresWhitelist: false,
                active: false
            });
            
            // Check eligibility for inactive template
            const [eligible, reason] = await distributor.canClaim(user1.address, "welcome_bonus");
            expect(eligible).to.be.false;
            expect(reason).to.equal("Template not active");
        });
    });
});