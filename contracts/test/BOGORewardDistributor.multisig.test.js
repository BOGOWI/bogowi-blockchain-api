const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("BOGORewardDistributor - Multisig Integration", function () {
  let bogoToken;
  let treasury;
  let rewardDistributor;
  let owner;
  let signer1;
  let signer2;
  let signer3;
  let user1;
  let user2;

  beforeEach(async function () {
    [owner, signer1, signer2, signer3, user1, user2] = await ethers.getSigners();

    // Deploy BOGO Token
    const BOGOToken = await ethers.getContractFactory("BOGOTokenV2");
    bogoToken = await BOGOToken.deploy();
    await bogoToken.deployed();

    // Deploy MultisigTreasury with 3 signers and threshold of 2
    const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
    treasury = await MultisigTreasury.deploy(
      [signer1.address, signer2.address, signer3.address],
      2
    );
    await treasury.deployed();

    // Deploy BOGORewardDistributor with treasury as controller
    const BOGORewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
    rewardDistributor = await BOGORewardDistributor.deploy(
      bogoToken.address,
      treasury.address
    );
    await rewardDistributor.deployed();

    // Grant DAO role to owner for minting
    const DAO_ROLE = await bogoToken.DAO_ROLE();
    await bogoToken.grantRole(DAO_ROLE, owner.address);
    
    // Fund the reward distributor using rewards allocation
    await bogoToken.mintFromRewards(rewardDistributor.address, ethers.utils.parseEther("1000000"));
  });

  describe("Access Control", function () {
    it("Should have treasury as the controller", async function () {
      expect(await rewardDistributor.treasury()).to.equal(treasury.address);
    });

    it("Should reject direct calls to admin functions", async function () {
      await expect(
        rewardDistributor.connect(owner).pause()
      ).to.be.revertedWith("Only treasury can call this function");

      await expect(
        rewardDistributor.connect(signer1).pause()
      ).to.be.revertedWith("Only treasury can call this function");
    });

    it("Should not have an owner() function", async function () {
      // This should throw because owner() doesn't exist anymore
      try {
        await rewardDistributor.owner();
        expect.fail("owner() should not exist");
      } catch (error) {
        expect(error.message).to.include("is not a function");
      }
    });
  });

  describe("Multisig Admin Operations", function () {
    it("Should allow treasury to pause/unpause through multisig", async function () {
      // Create pause transaction
      const pauseData = rewardDistributor.interface.encodeFunctionData("pause");
      await treasury.connect(signer1).submitTransaction(
        rewardDistributor.address,
        0,
        pauseData,
        "Pause reward distributor"
      );

      // Confirm by second signer
      await treasury.connect(signer2).confirmTransaction(0);
      
      // Fast forward time to pass execution delay
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      // Execute the transaction (by a signer)
      await treasury.connect(signer1).executeTransaction(0);

      // Check if paused
      expect(await rewardDistributor.paused()).to.be.true;

      // Create unpause transaction
      const unpauseData = rewardDistributor.interface.encodeFunctionData("unpause");
      await treasury.connect(signer1).submitTransaction(
        rewardDistributor.address,
        0,
        unpauseData,
        "Unpause reward distributor"
      );

      // Confirm
      await treasury.connect(signer3).confirmTransaction(1);
      
      // Fast forward time to pass execution delay
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      // Execute the transaction (by a signer)
      await treasury.connect(signer2).executeTransaction(1);

      // Check if unpaused
      expect(await rewardDistributor.paused()).to.be.false;
    });

    it("Should allow treasury to add to whitelist through multisig", async function () {
      const addresses = [user1.address, user2.address];
      
      // Create addToWhitelist transaction
      const whitelistData = rewardDistributor.interface.encodeFunctionData(
        "addToWhitelist",
        [addresses]
      );
      
      await treasury.connect(signer1).submitTransaction(
        rewardDistributor.address,
        0,
        whitelistData,
        "Add users to founder whitelist"
      );

      // Confirm by second signer
      await treasury.connect(signer2).confirmTransaction(0);
      
      // Fast forward time and execute
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(0);

      // Check if users are whitelisted
      expect(await rewardDistributor.founderWhitelist(user1.address)).to.be.true;
      expect(await rewardDistributor.founderWhitelist(user2.address)).to.be.true;
    });

    it("Should allow treasury to set authorized backend through multisig", async function () {
      const backendAddress = user1.address;
      
      // Create setAuthorizedBackend transaction
      const backendData = rewardDistributor.interface.encodeFunctionData(
        "setAuthorizedBackend",
        [backendAddress, true]
      );
      
      await treasury.connect(signer2).submitTransaction(
        rewardDistributor.address,
        0,
        backendData,
        "Authorize backend address"
      );

      // Confirm by another signer
      await treasury.connect(signer3).confirmTransaction(0);
      
      // Fast forward time and execute
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(0);

      // Check if backend is authorized
      expect(await rewardDistributor.authorizedBackends(backendAddress)).to.be.true;
    });

    it("Should allow treasury to update reward template through multisig", async function () {
      const newTemplate = {
        id: "test_reward",
        fixedAmount: ethers.utils.parseEther("50"),
        maxAmount: 0,
        cooldownPeriod: 3600, // 1 hour
        maxClaimsPerWallet: 5,
        requiresWhitelist: false,
        active: true
      };
      
      // Create updateTemplate transaction
      const updateData = rewardDistributor.interface.encodeFunctionData(
        "updateTemplate",
        ["test_reward", newTemplate]
      );
      
      await treasury.connect(signer1).submitTransaction(
        rewardDistributor.address,
        0,
        updateData,
        "Update test reward template"
      );

      // Confirm by second signer
      await treasury.connect(signer2).confirmTransaction(0);
      
      // Fast forward time and execute
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(0);

      // Verify template was updated
      const template = await rewardDistributor.templates("test_reward");
      expect(template.fixedAmount).to.equal(ethers.utils.parseEther("50"));
      expect(template.cooldownPeriod).to.equal(3600);
      expect(template.maxClaimsPerWallet).to.equal(5);
      expect(template.active).to.be.true;
    });
  });

  describe("Normal Operations", function () {
    it("Should allow users to claim rewards normally", async function () {
      // User claims welcome bonus
      const balanceBefore = await bogoToken.balanceOf(user1.address);
      await rewardDistributor.connect(user1).claimReward("welcome_bonus");
      const balanceAfter = await bogoToken.balanceOf(user1.address);
      
      expect(balanceAfter.sub(balanceBefore)).to.equal(ethers.utils.parseEther("10"));
    });

    it("Should allow authorized backends to distribute custom rewards", async function () {
      // First authorize the backend through multisig
      const backendData = rewardDistributor.interface.encodeFunctionData(
        "setAuthorizedBackend",
        [signer1.address, true]
      );
      
      await treasury.connect(signer2).submitTransaction(
        rewardDistributor.address,
        0,
        backendData,
        "Authorize backend"
      );
      await treasury.connect(signer3).confirmTransaction(0);
      
      // Fast forward time and execute
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(0);

      // Now backend can distribute custom rewards
      await rewardDistributor.connect(signer1).claimCustomReward(
        user2.address,
        ethers.utils.parseEther("25"),
        "Achievement reward"
      );

      expect(await bogoToken.balanceOf(user2.address)).to.equal(ethers.utils.parseEther("25"));
    });
  });

  describe("Security", function () {
    it("Should not allow initialization with zero treasury address", async function () {
      const BOGORewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
      await expect(
        BOGORewardDistributor.deploy(bogoToken.address, ethers.constants.AddressZero)
      ).to.be.revertedWith("Invalid treasury address");
    });

    it("Should maintain all existing security features", async function () {
      // Reentrancy protection still works
      expect(await rewardDistributor.connect(user1).claimReward("welcome_bonus")).to.emit(
        rewardDistributor,
        "RewardClaimed"
      );
      
      // Cannot claim same one-time reward twice
      await expect(
        rewardDistributor.connect(user1).claimReward("welcome_bonus")
      ).to.be.revertedWith("Max claims reached");
    });
  });
});