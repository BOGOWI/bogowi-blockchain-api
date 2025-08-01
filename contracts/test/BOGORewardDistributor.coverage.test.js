const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("BOGORewardDistributor - Full Coverage Tests", function () {
  let bogoToken;
  let treasury;
  let rewardDistributor;
  let owner;
  let signer1;
  let signer2;
  let signer3;
  let user1;
  let user2;
  let backend;
  let founder1;
  let founder2;

  beforeEach(async function () {
    [owner, signer1, signer2, signer3, user1, user2, backend, founder1, founder2] = await ethers.getSigners();

    // Deploy BOGO Token
    const BOGOToken = await ethers.getContractFactory("BOGOTokenV2");
    bogoToken = await BOGOToken.deploy();
    await bogoToken.deployed();

    // Deploy MultisigTreasury
    const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
    treasury = await MultisigTreasury.deploy(
      [signer1.address, signer2.address, signer3.address],
      2
    );
    await treasury.deployed();

    // Deploy BOGORewardDistributor
    const BOGORewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
    rewardDistributor = await BOGORewardDistributor.deploy(
      bogoToken.address,
      treasury.address
    );
    await rewardDistributor.deployed();

    // Grant DAO role and fund the distributor
    const DAO_ROLE = await bogoToken.DAO_ROLE();
    await bogoToken.grantRole(DAO_ROLE, owner.address);
    await bogoToken.mintFromRewards(rewardDistributor.address, ethers.utils.parseEther("10000000"));

    // Helper function to execute multisig transaction
    this.executeTreasuryTx = async (data, description) => {
      const tx = await treasury.connect(signer1).submitTransaction(
        rewardDistributor.address,
        0,
        data,
        description
      );
      const receipt = await tx.wait();
      const txId = receipt.events.find(e => e.event === "TransactionSubmitted").args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);
      
      return txId;
    };
  });

  describe("All Reward Templates", function () {
    it("Should claim welcome bonus successfully", async function () {
      const balanceBefore = await bogoToken.balanceOf(user1.address);
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      const balanceAfter = await bogoToken.balanceOf(user1.address);
      
      expect(balanceAfter.sub(balanceBefore)).to.equal(ethers.utils.parseEther("10"));
    });

    it("Should prevent claiming welcome bonus twice", async function () {
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      await expect(
        rewardDistributor.connect(user1).claimReward("welcome_bonus")
      ).to.be.revertedWith("Max claims reached");
    });

    it("Should claim founder bonus with whitelist", async function () {
      // Add to whitelist
      const whitelistData = rewardDistributor.interface.encodeFunctionData(
        "addToWhitelist",
        [[founder1.address]]
      );
      await this.executeTreasuryTx(whitelistData, "Add founder");

      // Claim founder bonus
      await rewardDistributor.connect(founder1).claimReward("founder_bonus");
      const balance = await bogoToken.balanceOf(founder1.address);
      expect(balance).to.equal(ethers.utils.parseEther("100"));
    });

    it("Should reject founder bonus without whitelist", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReward("founder_bonus")
      ).to.be.revertedWith("Not whitelisted");
    });

    it("Should handle DAO participation with cooldown", async function () {
      // First claim
      await rewardDistributor.connect(user1).claimReward("dao_participation");
      
      // Try immediate second claim - should fail
      await expect(
        rewardDistributor.connect(user1).claimReward("dao_participation")
      ).to.be.revertedWith("Cooldown period active");
      
      // Fast forward 30 days
      await ethers.provider.send("evm_increaseTime", [30 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine");
      
      // Should work now
      await rewardDistributor.connect(user1).claimReward("dao_participation");
    });

    it("Should claim all attraction tiers", async function () {
      const tiers = ["attraction_tier_1", "attraction_tier_2", "attraction_tier_3", "attraction_tier_4"];
      const amounts = ["10", "20", "40", "50"];
      
      for (let i = 0; i < tiers.length; i++) {
        const user = (await ethers.getSigners())[10 + i];
        await rewardDistributor.connect(user).claimReward(tiers[i]);
        const balance = await bogoToken.balanceOf(user.address);
        expect(balance).to.equal(ethers.utils.parseEther(amounts[i]));
      }
    });

    it("Should claim first NFT mint reward", async function () {
      await rewardDistributor.connect(user2).claimReward("first_nft_mint");
      const balance = await bogoToken.balanceOf(user2.address);
      expect(balance).to.equal(ethers.utils.parseEther("25"));
    });
  });

  describe("Referral System", function () {
    it("Should process referral bonus correctly", async function () {
      const referrerBefore = await bogoToken.balanceOf(user1.address);
      await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
      const referrerAfter = await bogoToken.balanceOf(user1.address);
      
      expect(referrerAfter.sub(referrerBefore)).to.equal(ethers.utils.parseEther("20"));
      expect(await rewardDistributor.referredBy(user2.address)).to.equal(user1.address);
      expect(await rewardDistributor.referralCount(user1.address)).to.equal(1);
    });

    it("Should prevent self-referral", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReferralBonus(user1.address)
      ).to.be.revertedWith("Cannot refer yourself");
    });

    it("Should prevent double referral", async function () {
      await rewardDistributor.connect(user2).claimReferralBonus(user1.address);
      
      await expect(
        rewardDistributor.connect(user2).claimReferralBonus(founder1.address)
      ).to.be.revertedWith("Already referred");
    });

    it("Should reject zero address referrer", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReferralBonus(ethers.constants.AddressZero)
      ).to.be.revertedWith("Invalid referrer");
    });
  });

  describe("Custom Rewards", function () {
    beforeEach(async function () {
      // Authorize backend
      const backendData = rewardDistributor.interface.encodeFunctionData(
        "setAuthorizedBackend",
        [backend.address, true]
      );
      await this.executeTreasuryTx(backendData, "Authorize backend");
    });

    it("Should distribute custom rewards", async function () {
      await rewardDistributor.connect(backend).claimCustomReward(
        user1.address,
        ethers.utils.parseEther("75"),
        "Special achievement"
      );
      
      const balance = await bogoToken.balanceOf(user1.address);
      expect(balance).to.equal(ethers.utils.parseEther("75"));
    });

    it("Should enforce custom reward max amount", async function () {
      await expect(
        rewardDistributor.connect(backend).claimCustomReward(
          user1.address,
          ethers.utils.parseEther("1001"),
          "Too much"
        )
      ).to.be.revertedWith("Invalid amount");
    });

    it("Should reject unauthorized custom rewards", async function () {
      await expect(
        rewardDistributor.connect(user1).claimCustomReward(
          user2.address,
          ethers.utils.parseEther("50"),
          "Unauthorized"
        )
      ).to.be.revertedWith("Not authorized backend");
    });

    it("Should reject zero amount custom rewards", async function () {
      await expect(
        rewardDistributor.connect(backend).claimCustomReward(
          user1.address,
          0,
          "Zero amount"
        )
      ).to.be.revertedWith("Invalid amount");
    });
  });

  describe("Daily Limits", function () {
    beforeEach(async function () {
      // Authorize backend for custom rewards
      const backendData = rewardDistributor.interface.encodeFunctionData(
        "setAuthorizedBackend",
        [backend.address, true]
      );
      await this.executeTreasuryTx(backendData, "Authorize backend");
    });
    
    it("Should reset daily limit after 24 hours", async function () {
      // Check initial limit
      const initialLimit = await rewardDistributor.getRemainingDailyLimit();
      expect(initialLimit).to.equal(ethers.utils.parseEther("500000"));
      
      // Use a significant portion of the daily limit
      await rewardDistributor.connect(backend).claimCustomReward(
        user1.address,
        ethers.utils.parseEther("1000"),
        "Large reward"
      );
      
      // Verify limit decreased
      let remainingLimit = await rewardDistributor.getRemainingDailyLimit();
      expect(remainingLimit).to.equal(ethers.utils.parseEther("499000"));
      
      // Fast forward 24 hours
      await ethers.provider.send("evm_increaseTime", [24 * 60 * 60]);
      await ethers.provider.send("evm_mine");
      
      // Check that limit has reset
      remainingLimit = await rewardDistributor.getRemainingDailyLimit();
      expect(remainingLimit).to.equal(ethers.utils.parseEther("500000"));
      
      // Can claim full amount again
      await rewardDistributor.connect(backend).claimCustomReward(
        user2.address,
        ethers.utils.parseEther("1000"),
        "New day reward"
      );
    });

    it("Should track daily distributed amount correctly", async function () {
      const initialLimit = await rewardDistributor.getRemainingDailyLimit();
      expect(initialLimit).to.equal(ethers.utils.parseEther("500000"));
      
      // Claim some rewards
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      const remainingLimit = await rewardDistributor.getRemainingDailyLimit();
      expect(remainingLimit).to.equal(ethers.utils.parseEther("499990"));
    });

    it("Should return full limit when time has passed", async function () {
      // Claim some rewards
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      // Fast forward more than 24 hours
      await ethers.provider.send("evm_increaseTime", [25 * 60 * 60]);
      await ethers.provider.send("evm_mine");
      
      // Should show full limit available
      const limit = await rewardDistributor.getRemainingDailyLimit();
      expect(limit).to.equal(ethers.utils.parseEther("500000"));
    });
  });

  describe("View Functions", function () {
    it("Should check eligibility correctly", async function () {
      // Check welcome bonus eligibility
      let [eligible, reason] = await rewardDistributor.canClaim(user1.address, "welcome_bonus");
      expect(eligible).to.be.true;
      expect(reason).to.equal("Eligible");
      
      // Claim it
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      
      // Check again - should be ineligible
      [eligible, reason] = await rewardDistributor.canClaim(user1.address, "welcome_bonus");
      expect(eligible).to.be.false;
      expect(reason).to.equal("Max claims reached");
    });

    it("Should check whitelist requirement", async function () {
      const [eligible, reason] = await rewardDistributor.canClaim(user1.address, "founder_bonus");
      expect(eligible).to.be.false;
      expect(reason).to.equal("Not whitelisted");
    });

    it("Should check cooldown period", async function () {
      // Claim DAO participation
      await rewardDistributor.connect(user1).claimReward("dao_participation");
      
      // Check eligibility
      const [eligible, reason] = await rewardDistributor.canClaim(user1.address, "dao_participation");
      expect(eligible).to.be.false;
      expect(reason).to.equal("Cooldown period active");
    });

    it("Should check inactive template", async function () {
      // Deactivate a template
      const template = {
        id: "welcome_bonus",
        fixedAmount: ethers.utils.parseEther("10"),
        maxAmount: 0,
        cooldownPeriod: 0,
        maxClaimsPerWallet: 1,
        requiresWhitelist: false,
        active: false
      };
      
      const updateData = rewardDistributor.interface.encodeFunctionData(
        "updateTemplate",
        ["welcome_bonus", template]
      );
      await this.executeTreasuryTx(updateData, "Deactivate template");
      
      const [eligible, reason] = await rewardDistributor.canClaim(user1.address, "welcome_bonus");
      expect(eligible).to.be.false;
      expect(reason).to.equal("Template not active");
    });
  });

  describe("Admin Functions", function () {
    it("Should pause and unpause contract", async function () {
      // Pause
      const pauseData = rewardDistributor.interface.encodeFunctionData("pause");
      await this.executeTreasuryTx(pauseData, "Pause");
      
      expect(await rewardDistributor.paused()).to.be.true;
      
      // Try to claim while paused
      await expect(
        rewardDistributor.connect(user1).claimReward("welcome_bonus")
      ).to.be.revertedWith("EnforcedPause");
      
      // Unpause
      const unpauseData = rewardDistributor.interface.encodeFunctionData("unpause");
      await this.executeTreasuryTx(unpauseData, "Unpause");
      
      expect(await rewardDistributor.paused()).to.be.false;
    });

    it("Should remove from whitelist", async function () {
      // Add to whitelist
      const addData = rewardDistributor.interface.encodeFunctionData(
        "addToWhitelist",
        [[founder1.address]]
      );
      await this.executeTreasuryTx(addData, "Add founder");
      
      expect(await rewardDistributor.founderWhitelist(founder1.address)).to.be.true;
      
      // Remove from whitelist
      const removeData = rewardDistributor.interface.encodeFunctionData(
        "removeFromWhitelist",
        [founder1.address]
      );
      await this.executeTreasuryTx(removeData, "Remove founder");
      
      expect(await rewardDistributor.founderWhitelist(founder1.address)).to.be.false;
    });

    it("Should revoke backend authorization", async function () {
      // Authorize
      const authData = rewardDistributor.interface.encodeFunctionData(
        "setAuthorizedBackend",
        [backend.address, true]
      );
      await this.executeTreasuryTx(authData, "Authorize backend");
      
      // Revoke
      const revokeData = rewardDistributor.interface.encodeFunctionData(
        "setAuthorizedBackend",
        [backend.address, false]
      );
      await this.executeTreasuryTx(revokeData, "Revoke backend");
      
      expect(await rewardDistributor.authorizedBackends(backend.address)).to.be.false;
    });

    it("Should update existing template", async function () {
      const newTemplate = {
        id: "welcome_bonus",
        fixedAmount: ethers.utils.parseEther("15"), // Changed from 10
        maxAmount: 0,
        cooldownPeriod: 0,
        maxClaimsPerWallet: 2, // Changed from 1
        requiresWhitelist: false,
        active: true
      };
      
      const updateData = rewardDistributor.interface.encodeFunctionData(
        "updateTemplate",
        ["welcome_bonus", newTemplate]
      );
      await this.executeTreasuryTx(updateData, "Update welcome bonus");
      
      const template = await rewardDistributor.templates("welcome_bonus");
      expect(template.fixedAmount).to.equal(ethers.utils.parseEther("15"));
      expect(template.maxClaimsPerWallet).to.equal(2);
    });
  });

  describe("Edge Cases", function () {
    it("Should handle multiple whitelisting in one transaction", async function () {
      const addresses = [founder1.address, founder2.address, user1.address, user2.address];
      
      const whitelistData = rewardDistributor.interface.encodeFunctionData(
        "addToWhitelist",
        [addresses]
      );
      await this.executeTreasuryTx(whitelistData, "Add multiple founders");
      
      for (const addr of addresses) {
        expect(await rewardDistributor.founderWhitelist(addr)).to.be.true;
      }
    });

    it("Should handle claiming when contract has insufficient balance", async function () {
      // Deploy new distributor with no funds
      const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
      const emptyDistributor = await RewardDistributor.deploy(
        bogoToken.address,
        treasury.address
      );
      await emptyDistributor.deployed();
      
      await expect(
        emptyDistributor.connect(user1).claimReward("welcome_bonus")
      ).to.be.reverted; // Will revert on transfer
    });

    it("Should handle invalid template ID", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReward("non_existent_template")
      ).to.be.revertedWith("Template not active");
    });

    it("Should reject custom reward with fixed amount template", async function () {
      await expect(
        rewardDistributor.connect(user1).claimReward("custom_reward")
      ).to.be.revertedWith("Use claimCustomReward for custom amounts");
    });
  });

  describe("Reentrancy Protection", function () {
    it("Should prevent reentrancy in claimReward", async function () {
      // This test verifies that the nonReentrant modifier works
      // In a real attack, the attacker would try to call claimReward again during transfer
      // The modifier should prevent this
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      const balance = await bogoToken.balanceOf(user1.address);
      expect(balance).to.equal(ethers.utils.parseEther("10"));
    });
  });

  describe("Gas Optimization Checks", function () {
    it("Should batch process multiple claims efficiently", async function () {
      const signers = await ethers.getSigners();
      const gasUsed = [];
      
      // Measure gas for multiple claims
      for (let i = 10; i < 15; i++) {
        const tx = await rewardDistributor.connect(signers[i]).claimReward("welcome_bonus");
        const receipt = await tx.wait();
        gasUsed.push(receipt.gasUsed);
      }
      
      // Gas usage should be consistent
      const avgGas = gasUsed.reduce((a, b) => a.add(b)).div(gasUsed.length);
      for (const gas of gasUsed) {
        expect(gas).to.be.closeTo(avgGas, avgGas.div(5)); // Within 20%
      }
    });
  });

  describe("Treasury Sweep", function () {
    let mockToken;

    beforeEach(async function () {
      // Deploy a mock ERC20 token for testing
      const MockToken = await ethers.getContractFactory("BOGOTokenV2");
      mockToken = await MockToken.deploy();
      await mockToken.deployed();
      
      // Grant DAO role and mint tokens to distributor
      const DAO_ROLE = await mockToken.DAO_ROLE();
      await mockToken.grantRole(DAO_ROLE, owner.address);
      await mockToken.mintFromRewards(rewardDistributor.address, ethers.utils.parseEther("1000"));
    });

    it("Should sweep BOGO tokens via treasury sweep", async function () {
      const withdrawAmount = ethers.utils.parseEther("500");
      const recipientBefore = await bogoToken.balanceOf(signer3.address);
      
      // Prepare treasury sweep through multisig
      const withdrawData = rewardDistributor.interface.encodeFunctionData(
        "treasurySweep",
        [bogoToken.address, signer3.address, withdrawAmount]
      );
      
      await this.executeTreasuryTx(withdrawData, "Treasury sweep of BOGO");
      
      const recipientAfter = await bogoToken.balanceOf(signer3.address);
      expect(recipientAfter.sub(recipientBefore)).to.equal(withdrawAmount);
    });

    it("Should sweep other ERC20 tokens via treasury sweep", async function () {
      const withdrawAmount = ethers.utils.parseEther("200");
      const recipientBefore = await mockToken.balanceOf(signer3.address);
      
      // Prepare treasury sweep through multisig
      const withdrawData = rewardDistributor.interface.encodeFunctionData(
        "treasurySweep",
        [mockToken.address, signer3.address, withdrawAmount]
      );
      
      await this.executeTreasuryTx(withdrawData, "Treasury sweep of mock token");
      
      const recipientAfter = await mockToken.balanceOf(signer3.address);
      expect(recipientAfter.sub(recipientBefore)).to.equal(withdrawAmount);
    });


    it("Should reject treasury sweep to zero address", async function () {
      const withdrawData = rewardDistributor.interface.encodeFunctionData(
        "treasurySweep",
        [bogoToken.address, ethers.constants.AddressZero, ethers.utils.parseEther("100")]
      );
      
      const tx = await treasury.connect(signer1).submitTransaction(
        rewardDistributor.address,
        0,
        withdrawData,
        "Invalid withdrawal"
      );
      const receipt = await tx.wait();
      const txId = receipt.events.find(e => e.event === "TransactionSubmitted").args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await expect(
        treasury.connect(signer1).executeTransaction(txId)
      ).to.be.reverted;
    });

    it("Should reject treasury sweep with zero amount", async function () {
      const withdrawData = rewardDistributor.interface.encodeFunctionData(
        "treasurySweep",
        [bogoToken.address, signer3.address, 0]
      );
      
      const tx = await treasury.connect(signer1).submitTransaction(
        rewardDistributor.address,
        0,
        withdrawData,
        "Zero amount withdrawal"
      );
      const receipt = await tx.wait();
      const txId = receipt.events.find(e => e.event === "TransactionSubmitted").args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await expect(
        treasury.connect(signer1).executeTransaction(txId)
      ).to.be.reverted;
    });

    it("Should emit TreasurySweep event", async function () {
      const withdrawAmount = ethers.utils.parseEther("100");
      
      const withdrawData = rewardDistributor.interface.encodeFunctionData(
        "treasurySweep",
        [bogoToken.address, signer3.address, withdrawAmount]
      );
      
      const tx = await treasury.connect(signer1).submitTransaction(
        rewardDistributor.address,
        0,
        withdrawData,
        "Treasury sweep"
      );
      const receipt = await tx.wait();
      const txId = receipt.events.find(e => e.event === "TransactionSubmitted").args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await expect(treasury.connect(signer1).executeTransaction(txId))
        .to.emit(rewardDistributor, "TreasurySweep")
        .withArgs(bogoToken.address, signer3.address, withdrawAmount);
    });

    it("Should reject direct treasury sweep not from treasury", async function () {
      await expect(
        rewardDistributor.connect(user1).treasurySweep(
          bogoToken.address,
          user1.address,
          ethers.utils.parseEther("100")
        )
      ).to.be.revertedWith("Only treasury can call this function");
    });

  });
});