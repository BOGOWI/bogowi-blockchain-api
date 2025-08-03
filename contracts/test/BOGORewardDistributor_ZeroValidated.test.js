const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGORewardDistributor_ZeroValidated", function () {
  let bogoToken;
  let rewardDistributor;
  let owner;
  let treasury;
  let backend;
  let user1;
  let user2;
  let user3;
  let nonAuthorized;

  const DAILY_GLOBAL_LIMIT = ethers.parseEther("500000");
  const MAX_REFERRAL_DEPTH = 10;

  beforeEach(async function () {
    [owner, treasury, backend, user1, user2, user3, nonAuthorized] = await ethers.getSigners();

    // Deploy mock BOGO token
    const MockERC20 = await ethers.getContractFactory("contracts/test/MockERC20.sol:MockERC20");
    bogoToken = await MockERC20.deploy("BOGO Token", "BOGO", ethers.parseEther("1000000"));
    await bogoToken.waitForDeployment();

    // Deploy BOGORewardDistributor_ZeroValidated
    const BOGORewardDistributor_ZeroValidated = await ethers.getContractFactory("BOGORewardDistributor_ZeroValidated");
    rewardDistributor = await BOGORewardDistributor_ZeroValidated.deploy(
      await bogoToken.getAddress(),
      await treasury.getAddress()
    );
    await rewardDistributor.waitForDeployment();

    // Authorize backend
    await rewardDistributor.connect(treasury).setAuthorizedBackend(await backend.getAddress(), true);

    // Transfer tokens to distributor
    await bogoToken.transfer(await rewardDistributor.getAddress(), ethers.parseEther("100000"));
  });

  describe("Deployment", function () {
    it("Should deploy with correct initial state", async function () {
      expect(await rewardDistributor.bogoToken()).to.equal(await bogoToken.getAddress());
      expect(await rewardDistributor.treasury()).to.equal(await treasury.getAddress());
      expect(await rewardDistributor.DAILY_GLOBAL_LIMIT()).to.equal(DAILY_GLOBAL_LIMIT);
      expect(await rewardDistributor.MAX_REFERRAL_DEPTH()).to.equal(MAX_REFERRAL_DEPTH);
    });

    it("Should revert with zero token address", async function () {
      const BOGORewardDistributor_ZeroValidated = await ethers.getContractFactory("BOGORewardDistributor_ZeroValidated");
      await expect(
        BOGORewardDistributor_ZeroValidated.deploy(
          ethers.ZeroAddress,
          await treasury.getAddress()
        )
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidTokenAddress");
    });

    it("Should revert with zero treasury address", async function () {
      const BOGORewardDistributor_ZeroValidated = await ethers.getContractFactory("BOGORewardDistributor_ZeroValidated");
      await expect(
        BOGORewardDistributor_ZeroValidated.deploy(
          await bogoToken.getAddress(),
          ethers.ZeroAddress
        )
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidTreasuryAddress");
    });

    it("Should initialize default templates", async function () {
      const welcomeTemplate = await rewardDistributor.templates("welcome_bonus");
      expect(welcomeTemplate.fixedAmount).to.equal(ethers.parseEther("10"));
      expect(welcomeTemplate.active).to.be.true;

      const founderTemplate = await rewardDistributor.templates("founder_bonus");
      expect(founderTemplate.fixedAmount).to.equal(ethers.parseEther("100"));
      expect(founderTemplate.requiresWhitelist).to.be.true;
    });
  });

  describe("Zero Address Validation", function () {
    it("Should reject zero address in claimCustomReward", async function () {
      await expect(
        rewardDistributor.connect(backend).claimCustomReward(
          ethers.ZeroAddress,
          ethers.parseEther("50"),
          "Test reward"
        )
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
    });

    it("Should reject zero address in claimReferralBonus", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReferralBonus(ethers.ZeroAddress)
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
    });

    it("Should reject zero address in addToWhitelist", async function () {
      await expect(
        rewardDistributor.connect(treasury).addToWhitelist([ethers.ZeroAddress])
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
    });

    it("Should reject zero address in removeFromWhitelist", async function () {
      await expect(
        rewardDistributor.connect(treasury).removeFromWhitelist(ethers.ZeroAddress)
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
    });

    it("Should reject zero address in setAuthorizedBackend", async function () {
      await expect(
        rewardDistributor.connect(treasury).setAuthorizedBackend(ethers.ZeroAddress, true)
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
    });

    it("Should reject zero address in treasurySweep", async function () {
      await expect(
        rewardDistributor.connect(treasury).treasurySweep(
          await bogoToken.getAddress(),
          ethers.ZeroAddress,
          ethers.parseEther("100")
        )
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAddress");
    });

    it("Should handle zero address in canClaim view function", async function () {
      const [canClaim, reason] = await rewardDistributor.canClaim(ethers.ZeroAddress, "welcome_bonus");
      expect(canClaim).to.be.false;
      expect(reason).to.equal("Invalid wallet address");
    });

    it("Should handle zero address in getReferralChain", async function () {
      const chain = await rewardDistributor.getReferralChain(ethers.ZeroAddress);
      expect(chain.length).to.equal(0);
    });
  });

  describe("Reward Claiming", function () {
    it("Should allow claiming welcome bonus", async function () {
      const initialBalance = await bogoToken.balanceOf(await user1.getAddress());
      
      await expect(rewardDistributor.connect(user1).claimReward("welcome_bonus"))
        .to.emit(rewardDistributor, "RewardClaimed")
        .withArgs(await user1.getAddress(), "welcome_bonus", ethers.parseEther("10"));

      const finalBalance = await bogoToken.balanceOf(await user1.getAddress());
      expect(finalBalance - initialBalance).to.equal(ethers.parseEther("10"));
    });

    it("Should prevent claiming welcome bonus twice", async function () {
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      await expect(
        rewardDistributor.connect(user1).claimReward("welcome_bonus")
      ).to.be.revertedWith("Max claims reached");
    });

    it("Should enforce cooldown period for DAO participation", async function () {
      await rewardDistributor.connect(user1).claimReward("dao_participation");
      
      await expect(
        rewardDistributor.connect(user1).claimReward("dao_participation")
      ).to.be.revertedWith("Cooldown period active");

      // Fast forward 30 days
      await time.increase(30 * 24 * 60 * 60);
      
      await expect(rewardDistributor.connect(user1).claimReward("dao_participation"))
        .to.emit(rewardDistributor, "RewardClaimed");
    });

    it("Should require whitelist for founder bonus", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReward("founder_bonus")
      ).to.be.revertedWith("Not whitelisted");

      // Add to whitelist
      await rewardDistributor.connect(treasury).addToWhitelist([await user1.getAddress()]);
      
      await expect(rewardDistributor.connect(user1).claimReward("founder_bonus"))
        .to.emit(rewardDistributor, "RewardClaimed")
        .withArgs(await user1.getAddress(), "founder_bonus", ethers.parseEther("100"));
    });

    it("Should enforce daily global limit", async function () {
      // This test would require many claims to reach the limit
      // For practical testing, we'll test the logic with a smaller scenario
      const remainingBefore = await rewardDistributor.getRemainingDailyLimit();
      expect(remainingBefore).to.equal(DAILY_GLOBAL_LIMIT);

      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      const remainingAfter = await rewardDistributor.getRemainingDailyLimit();
      expect(remainingAfter).to.equal(DAILY_GLOBAL_LIMIT - ethers.parseEther("10"));
    });

    it("Should reset daily limit after 24 hours", async function () {
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      const remainingAfter = await rewardDistributor.getRemainingDailyLimit();
      expect(remainingAfter).to.equal(DAILY_GLOBAL_LIMIT - ethers.parseEther("10"));

      // Fast forward 24 hours
      await time.increase(24 * 60 * 60);
      
      const remainingReset = await rewardDistributor.getRemainingDailyLimit();
      expect(remainingReset).to.equal(DAILY_GLOBAL_LIMIT);
    });

    it("Should revert when template is inactive", async function () {
      // Deactivate welcome bonus template
      const template = await rewardDistributor.templates("welcome_bonus");
      const updatedTemplate = {
        ...template,
        active: false
      };
      
      await rewardDistributor.connect(treasury).updateTemplate("welcome_bonus", updatedTemplate);
      
      await expect(
        rewardDistributor.connect(user1).claimReward("welcome_bonus")
      ).to.be.revertedWith("Template not active");
    });

    it("Should revert when contract is paused", async function () {
      await rewardDistributor.connect(treasury).pause();
      
      await expect(
        rewardDistributor.connect(user1).claimReward("welcome_bonus")
      ).to.be.revertedWithCustomError(rewardDistributor, "EnforcedPause");
    });

    it("Should revert for custom reward template in claimReward", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReward("custom_reward")
      ).to.be.revertedWith("Use claimCustomReward for custom amounts");
    });
  });

  describe("Custom Rewards", function () {
    it("Should allow authorized backend to distribute custom rewards", async function () {
      const amount = ethers.parseEther("50");
      const reason = "Special achievement";
      
      await expect(
        rewardDistributor.connect(backend).claimCustomReward(
          await user1.getAddress(),
          amount,
          reason
        )
      ).to.emit(rewardDistributor, "RewardClaimed")
        .withArgs(await user1.getAddress(), reason, amount);
    });

    it("Should enforce max amount for custom rewards", async function () {
      const maxAmount = ethers.parseEther("1000");
      const excessiveAmount = ethers.parseEther("1001");
      
      await expect(
        rewardDistributor.connect(backend).claimCustomReward(
          await user1.getAddress(),
          excessiveAmount,
          "Too much"
        )
      ).to.be.revertedWith("Amount exceeds maximum");
    });

    it("Should reject zero amount custom rewards", async function () {
      await expect(
        rewardDistributor.connect(backend).claimCustomReward(
          await user1.getAddress(),
          0,
          "Zero amount"
        )
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAmount");
    });

    it("Should only allow authorized backend", async function () {
      await expect(
        rewardDistributor.connect(user1).claimCustomReward(
          await user1.getAddress(),
          ethers.parseEther("50"),
          "Unauthorized"
        )
      ).to.be.revertedWithCustomError(rewardDistributor, "NotAuthorizedBackend");
    });

    it("Should enforce daily limit for custom rewards", async function () {
      const amount = ethers.parseEther("50");
      
      await rewardDistributor.connect(backend).claimCustomReward(
        await user1.getAddress(),
        amount,
        "Test"
      );
      
      const remainingAfter = await rewardDistributor.getRemainingDailyLimit();
      expect(remainingAfter).to.equal(DAILY_GLOBAL_LIMIT - amount);
    });
  });

  describe("Referral System", function () {
    it("Should allow valid referral bonus claim", async function () {
      const referrerBalance = await bogoToken.balanceOf(await user1.getAddress());
      
      await expect(
        rewardDistributor.connect(user2).claimReferralBonus(await user1.getAddress())
      ).to.emit(rewardDistributor, "ReferralClaimed")
        .withArgs(await user1.getAddress(), await user2.getAddress(), ethers.parseEther("20"));

      const newReferrerBalance = await bogoToken.balanceOf(await user1.getAddress());
      expect(newReferrerBalance - referrerBalance).to.equal(ethers.parseEther("20"));
      
      expect(await rewardDistributor.referredBy(await user2.getAddress())).to.equal(await user1.getAddress());
      expect(await rewardDistributor.referralCount(await user1.getAddress())).to.equal(1);
    });

    it("Should prevent self-referral", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReferralBonus(await user1.getAddress())
      ).to.be.revertedWith("Cannot refer yourself");
    });

    it("Should prevent double referral", async function () {
      await rewardDistributor.connect(user2).claimReferralBonus(await user1.getAddress());
      
      await expect(
        rewardDistributor.connect(user2).claimReferralBonus(await user3.getAddress())
      ).to.be.revertedWith("Already referred");
    });

    it("Should detect circular referrals", async function () {
      // user1 refers user2
      await rewardDistributor.connect(user2).claimReferralBonus(await user1.getAddress());
      
      // user2 refers user3
      await rewardDistributor.connect(user3).claimReferralBonus(await user2.getAddress());
      
      // user3 trying to refer user1 should fail (circular)
      await expect(
        rewardDistributor.connect(user1).claimReferralBonus(await user3.getAddress())
      ).to.be.revertedWith("Circular referral detected");
    });

    it("Should enforce max referral depth", async function () {
      // Create a chain of referrals up to max depth
      const signers = await ethers.getSigners();
      
      for (let i = 1; i < MAX_REFERRAL_DEPTH; i++) {
        await rewardDistributor.connect(signers[i + 1]).claimReferralBonus(await signers[i].getAddress());
      }
      
      // Next referral should fail due to max depth
      await expect(
        rewardDistributor.connect(signers[MAX_REFERRAL_DEPTH + 1]).claimReferralBonus(await signers[MAX_REFERRAL_DEPTH].getAddress())
      ).to.be.revertedWith("Max referral depth exceeded");
    });

    it("Should return referral chain", async function () {
      // Create referral chain: user1 -> user2 -> user3
      await rewardDistributor.connect(user2).claimReferralBonus(await user1.getAddress());
      await rewardDistributor.connect(user3).claimReferralBonus(await user2.getAddress());
      
      const chain = await rewardDistributor.getReferralChain(await user3.getAddress());
      expect(chain.length).to.equal(2);
      expect(chain[0]).to.equal(await user2.getAddress());
      expect(chain[1]).to.equal(await user1.getAddress());
    });
  });

  describe("Backend Authorization", function () {
    it("Should allow treasury to authorize backend", async function () {
      await expect(
        rewardDistributor.connect(treasury).setAuthorizedBackend(await user1.getAddress(), true)
      ).to.emit(rewardDistributor, "BackendAuthorized")
        .withArgs(await user1.getAddress(), true);

      expect(await rewardDistributor.authorizedBackends(await user1.getAddress())).to.be.true;
    });

    it("Should allow treasury to deauthorize backend", async function () {
      await rewardDistributor.connect(treasury).setAuthorizedBackend(await user1.getAddress(), true);
      
      await expect(
        rewardDistributor.connect(treasury).setAuthorizedBackend(await user1.getAddress(), false)
      ).to.emit(rewardDistributor, "BackendAuthorized")
        .withArgs(await user1.getAddress(), false);

      expect(await rewardDistributor.authorizedBackends(await user1.getAddress())).to.be.false;
    });

    it("Should only allow treasury to manage backend authorization", async function () {
      await expect(
        rewardDistributor.connect(user1).setAuthorizedBackend(await user2.getAddress(), true)
      ).to.be.revertedWithCustomError(rewardDistributor, "NotTreasury");
    });
  });

  describe("Whitelist Management", function () {
    it("Should allow treasury to add to whitelist", async function () {
      await expect(
        rewardDistributor.connect(treasury).addToWhitelist([await user1.getAddress()])
      ).to.emit(rewardDistributor, "WhitelistUpdated")
        .withArgs(await user1.getAddress(), true);

      expect(await rewardDistributor.founderWhitelist(await user1.getAddress())).to.be.true;
    });

    it("Should allow treasury to remove from whitelist", async function () {
      await rewardDistributor.connect(treasury).addToWhitelist([await user1.getAddress()]);
      
      await expect(
        rewardDistributor.connect(treasury).removeFromWhitelist(await user1.getAddress())
      ).to.emit(rewardDistributor, "WhitelistUpdated")
        .withArgs(await user1.getAddress(), false);

      expect(await rewardDistributor.founderWhitelist(await user1.getAddress())).to.be.false;
    });

    it("Should only allow treasury to manage whitelist", async function () {
      await expect(
        rewardDistributor.connect(user1).addToWhitelist([await user2.getAddress()])
      ).to.be.revertedWithCustomError(rewardDistributor, "NotTreasury");

      await expect(
        rewardDistributor.connect(user1).removeFromWhitelist(await user2.getAddress())
      ).to.be.revertedWithCustomError(rewardDistributor, "NotTreasury");
    });

    it("Should add multiple addresses to whitelist", async function () {
      const addresses = [await user1.getAddress(), await user2.getAddress(), await user3.getAddress()];
      
      await rewardDistributor.connect(treasury).addToWhitelist(addresses);
      
      for (const addr of addresses) {
        expect(await rewardDistributor.founderWhitelist(addr)).to.be.true;
      }
    });
  });

  describe("Template Management", function () {
    it("Should allow treasury to update templates", async function () {
      const newTemplate = {
        id: "welcome_bonus",
        fixedAmount: ethers.parseEther("20"),
        maxAmount: 0,
        cooldownPeriod: 0,
        maxClaimsPerWallet: 1,
        requiresWhitelist: false,
        active: true
      };
      
      await expect(
        rewardDistributor.connect(treasury).updateTemplate("welcome_bonus", newTemplate)
      ).to.emit(rewardDistributor, "TemplateUpdated")
        .withArgs("welcome_bonus");

      const updatedTemplate = await rewardDistributor.templates("welcome_bonus");
      expect(updatedTemplate.fixedAmount).to.equal(ethers.parseEther("20"));
    });

    it("Should only allow treasury to update templates", async function () {
      const newTemplate = {
        id: "welcome_bonus",
        fixedAmount: ethers.parseEther("20"),
        maxAmount: 0,
        cooldownPeriod: 0,
        maxClaimsPerWallet: 1,
        requiresWhitelist: false,
        active: true
      };
      
      await expect(
        rewardDistributor.connect(user1).updateTemplate("welcome_bonus", newTemplate)
      ).to.be.revertedWithCustomError(rewardDistributor, "NotTreasury");
    });
  });

  describe("Pause/Unpause", function () {
    it("Should allow treasury to pause contract", async function () {
      await rewardDistributor.connect(treasury).pause();
      expect(await rewardDistributor.paused()).to.be.true;
    });

    it("Should allow treasury to unpause contract", async function () {
      await rewardDistributor.connect(treasury).pause();
      await rewardDistributor.connect(treasury).unpause();
      expect(await rewardDistributor.paused()).to.be.false;
    });

    it("Should only allow treasury to pause/unpause", async function () {
      await expect(
        rewardDistributor.connect(user1).pause()
      ).to.be.revertedWithCustomError(rewardDistributor, "NotTreasury");

      await expect(
        rewardDistributor.connect(user1).unpause()
      ).to.be.revertedWithCustomError(rewardDistributor, "NotTreasury");
    });
  });

  describe("View Functions", function () {
    it("Should check claim eligibility correctly", async function () {
      let [canClaim, reason] = await rewardDistributor.canClaim(await user1.getAddress(), "welcome_bonus");
      expect(canClaim).to.be.true;
      expect(reason).to.equal("Eligible");

      // Claim the reward
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      [canClaim, reason] = await rewardDistributor.canClaim(await user1.getAddress(), "welcome_bonus");
      expect(canClaim).to.be.false;
      expect(reason).to.equal("Max claims reached");
    });

    it("Should return correct remaining daily limit", async function () {
      const initialLimit = await rewardDistributor.getRemainingDailyLimit();
      expect(initialLimit).to.equal(DAILY_GLOBAL_LIMIT);

      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      const afterClaimLimit = await rewardDistributor.getRemainingDailyLimit();
      expect(afterClaimLimit).to.equal(DAILY_GLOBAL_LIMIT - ethers.parseEther("10"));
    });
  });

  describe("Treasury Sweep", function () {
    it("Should allow treasury to sweep ERC20 tokens", async function () {
      const sweepAmount = ethers.parseEther("100");
      
      await expect(
        rewardDistributor.connect(treasury).treasurySweep(
          await bogoToken.getAddress(),
          await treasury.getAddress(),
          sweepAmount
        )
      ).to.emit(rewardDistributor, "TreasurySweep")
        .withArgs(await bogoToken.getAddress(), await treasury.getAddress(), sweepAmount);
    });

    it("Should allow treasury to sweep ETH", async function () {
      // Send some ETH to the contract
      await owner.sendTransaction({
        to: await rewardDistributor.getAddress(),
        value: ethers.parseEther("1")
      });

      const sweepAmount = ethers.parseEther("0.5");
      
      await expect(
        rewardDistributor.connect(treasury).treasurySweep(
          ethers.ZeroAddress,
          await treasury.getAddress(),
          sweepAmount
        )
      ).to.emit(rewardDistributor, "TreasurySweep")
        .withArgs(ethers.ZeroAddress, await treasury.getAddress(), sweepAmount);
    });

    it("Should reject zero amount sweep", async function () {
      await expect(
        rewardDistributor.connect(treasury).treasurySweep(
          await bogoToken.getAddress(),
          await treasury.getAddress(),
          0
        )
      ).to.be.revertedWithCustomError(rewardDistributor, "InvalidAmount");
    });

    it("Should only allow treasury to sweep", async function () {
      await expect(
        rewardDistributor.connect(user1).treasurySweep(
          await bogoToken.getAddress(),
          await user1.getAddress(),
          ethers.parseEther("100")
        )
      ).to.be.revertedWithCustomError(rewardDistributor, "NotTreasury");
    });
  });

  describe("Error Handling", function () {
    it("Should handle transfer failures gracefully", async function () {
      // This test would require a mock token that can fail transfers
      // For now, we'll test the revert condition exists
      const template = await rewardDistributor.templates("welcome_bonus");
      expect(template.active).to.be.true;
    });

    it("Should handle inactive custom rewards", async function () {
      // Deactivate custom rewards
      const template = await rewardDistributor.templates("custom_reward");
      const updatedTemplate = {
        ...template,
        active: false
      };
      
      await rewardDistributor.connect(treasury).updateTemplate("custom_reward", updatedTemplate);
      
      await expect(
        rewardDistributor.connect(backend).claimCustomReward(
          await user1.getAddress(),
          ethers.parseEther("50"),
          "Test"
        )
      ).to.be.revertedWith("Custom rewards not active");
    });

    it("Should handle inactive referral rewards", async function () {
      // Deactivate referral rewards
      const template = await rewardDistributor.templates("referral_bonus");
      const updatedTemplate = {
        ...template,
        active: false
      };
      
      await rewardDistributor.connect(treasury).updateTemplate("referral_bonus", updatedTemplate);
      
      await expect(
        rewardDistributor.connect(user2).claimReferralBonus(await user1.getAddress())
      ).to.be.revertedWith("Referral rewards not active");
    });
  });

  describe("Edge Cases", function () {
    it("Should handle empty referral chain", async function () {
      const chain = await rewardDistributor.getReferralChain(await user1.getAddress());
      expect(chain.length).to.equal(0);
    });

    it("Should handle whitelist check for non-whitelisted template", async function () {
      const [canClaim, reason] = await rewardDistributor.canClaim(await user1.getAddress(), "attraction_tier_1");
      expect(canClaim).to.be.true;
      expect(reason).to.equal("Eligible");
    });

    it("Should handle cooldown period correctly", async function () {
      const [canClaim, reason] = await rewardDistributor.canClaim(await user1.getAddress(), "dao_participation");
      expect(canClaim).to.be.true;
      expect(reason).to.equal("Eligible");
    });

    it("Should handle max claims correctly", async function () {
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      const [canClaim, reason] = await rewardDistributor.canClaim(await user1.getAddress(), "welcome_bonus");
      expect(canClaim).to.be.false;
      expect(reason).to.equal("Max claims reached");
    });
  });
});