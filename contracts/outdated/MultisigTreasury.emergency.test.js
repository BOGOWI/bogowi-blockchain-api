const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("MultisigTreasury - Emergency Withdrawal Tests", function () {
  let treasury;
  let owner, signer1, signer2, signer3, signer4, nonSigner, recipient;

  beforeEach(async function () {
    [owner, signer1, signer2, signer3, signer4, nonSigner, recipient] = await ethers.getSigners();

    // Deploy MultisigTreasury with 3 signers, threshold 2
    const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
    treasury = await MultisigTreasury.deploy(
      [signer1.address, signer2.address, signer3.address],
      2
    );
    await treasury.waitForDeployment();

    // Fund the treasury with ETH
    await owner.sendTransaction({
      to: await treasury.getAddress(),
      value: ethers.parseEther("10")
    });
  });

  describe("Emergency Withdrawal Approval Accumulation", function () {
    beforeEach(async function () {
      // Pause the contract to enable emergency functions
      // First need to create and execute a pause transaction
      const pauseData = treasury.interface.encodeFunctionData("pause");
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        pauseData,
        "Pause for emergency"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs.find(e => e.fragment && e.fragment.name === "TransactionSubmitted").args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);
      
      expect(await treasury.paused()).to.be.true;
    });

    it("Should accumulate approvals without resetting after first signer", async function () {
      const withdrawAmount = ethers.parseEther("2");
      
      // First signer approves
      await expect(
        treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount)
      ).to.emit(treasury, "EmergencyApprovalGranted")
        .withArgs(await signer1.getAddress());
      
      // Verify approval was recorded
      expect(await treasury.emergencyApprovals(signer1.address)).to.be.true;
      expect(await treasury.emergencyApprovalCount()).to.equal(1);
      
      // Second signer should still be able to approve
      await expect(
        treasury.connect(signer2).emergencyWithdrawETH(recipient.address, withdrawAmount)
      ).to.emit(treasury, "EmergencyWithdraw")
        .withArgs(ethers.ZeroAddress, recipient.address, withdrawAmount);
      
      // Verify approvals were reset only after execution
      expect(await treasury.emergencyApprovals(signer1.address)).to.be.false;
      expect(await treasury.emergencyApprovals(signer2.address)).to.be.false;
      expect(await treasury.emergencyApprovalCount()).to.equal(0);
    });

    it("Should not allow double approval from same signer", async function () {
      const withdrawAmount = ethers.parseEther("2");
      
      // First approval
      await treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount);
      
      // Try to approve again
      await expect(
        treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount)
      ).to.be.revertedWith("ALREADY_PROCESSED");
    });

    it("Should require exact threshold approvals", async function () {
      const withdrawAmount = ethers.parseEther("2");
      const recipientBefore = await ethers.provider.getBalance(recipient.address);
      
      // First approval - should not execute
      await treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount);
      
      // Check no funds transferred yet
      const recipientAfterFirst = await ethers.provider.getBalance(recipient.address);
      expect(recipientAfterFirst).to.equal(recipientBefore);
      
      // Second approval - should execute
      await treasury.connect(signer2).emergencyWithdrawETH(recipient.address, withdrawAmount);
      
      // Check funds transferred
      const recipientAfterSecond = await ethers.provider.getBalance(recipient.address);
      expect(recipientAfterSecond - recipientBefore).to.equal(withdrawAmount);
    });

    it("Should respect current threshold for emergency approvals", async function () {
      // The treasury starts with threshold 2, which is what we'll test
      const withdrawAmount = ethers.parseEther("2");
      
      // Verify initial threshold
      expect(await treasury.threshold()).to.equal(2);
      
      // First approval
      await treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount);
      expect(await treasury.emergencyApprovalCount()).to.equal(1);
      
      // Verify withdrawal hasn't happened yet
      const balanceBefore = await ethers.provider.getBalance(recipient.address);
      
      // Second approval should trigger execution (threshold = 2)
      await expect(
        treasury.connect(signer2).emergencyWithdrawETH(recipient.address, withdrawAmount)
      ).to.emit(treasury, "EmergencyWithdraw");
      
      // Verify withdrawal happened
      const balanceAfter = await ethers.provider.getBalance(recipient.address);
      expect(balanceAfter - balanceBefore).to.equal(withdrawAmount);
      
      // Verify approvals were reset
      expect(await treasury.emergencyApprovalCount()).to.equal(0);
    });
  });

  describe("Emergency Withdrawal Security", function () {
    beforeEach(async function () {
      // Pause the contract
      const pauseData = treasury.interface.encodeFunctionData("pause");
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        pauseData,
        "Pause for emergency"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs.find(e => e.fragment && e.fragment.name === "TransactionSubmitted").args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);
    });

    it("Should only work when paused", async function () {
      // Create a new treasury instance that starts unpaused
      const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
      const unpausedTreasury = await MultisigTreasury.deploy(
        [await signer1.getAddress(), await signer2.getAddress(), await signer3.getAddress()],
        2
      );
      await unpausedTreasury.waitForDeployment();
      
      // Fund it
      await owner.sendTransaction({
        to: await unpausedTreasury.getAddress(),
        value: ethers.parseEther("10")
      });
      
      expect(await unpausedTreasury.paused()).to.be.false;
      
      // Try emergency withdrawal when not paused
      await expect(
        unpausedTreasury.connect(signer1).emergencyWithdrawETH(recipient.address, ethers.parseEther("1"))
      ).to.be.revertedWithCustomError(unpausedTreasury, "ExpectedPause");
    });

    it("Should enforce 50% balance limit", async function () {
      const balance = await ethers.provider.getBalance(await treasury.getAddress());
      const tooMuch = balance / 2n + 1n;
      
      await expect(
        treasury.connect(signer1).emergencyWithdrawETH(recipient.address, tooMuch)
      ).to.be.revertedWith("EXCEEDS_LIMIT");
    });

    it("Should reject zero address recipient", async function () {
      await expect(
        treasury.connect(signer1).emergencyWithdrawETH(
          ethers.ZeroAddress, 
          ethers.parseEther("1")
        )
      ).to.be.revertedWith("ZERO_ADDRESS");
    });

    it("Should only allow signers to approve", async function () {
      await expect(
        treasury.connect(nonSigner).emergencyWithdrawETH(
          recipient.address, 
          ethers.parseEther("1")
        )
      ).to.be.revertedWith("NOT_SIGNER");
    });
  });

  describe("Emergency Withdrawal State Management", function () {
    beforeEach(async function () {
      // Pause the contract
      const pauseData = treasury.interface.encodeFunctionData("pause");
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        pauseData,
        "Pause for emergency"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs.find(e => e.fragment && e.fragment.name === "TransactionSubmitted").args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);
    });

    it("Should reset all approvals after execution", async function () {
      const withdrawAmount = ethers.parseEther("2");
      
      // Get approvals from multiple signers
      await treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount);
      await treasury.connect(signer2).emergencyWithdrawETH(recipient.address, withdrawAmount);
      
      // All approvals should be reset
      expect(await treasury.emergencyApprovals(await signer1.getAddress())).to.be.false;
      expect(await treasury.emergencyApprovals(await signer2.getAddress())).to.be.false;
      expect(await treasury.emergencyApprovals(await signer3.getAddress())).to.be.false;
      expect(await treasury.emergencyApprovalCount()).to.equal(0);
    });

    it("Should allow new emergency withdrawal after previous one completes", async function () {
      const withdrawAmount = ethers.parseEther("2");
      
      // First emergency withdrawal
      await treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount);
      await treasury.connect(signer2).emergencyWithdrawETH(recipient.address, withdrawAmount);
      
      // Second emergency withdrawal should work
      await treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount);
      await expect(
        treasury.connect(signer3).emergencyWithdrawETH(recipient.address, withdrawAmount)
      ).to.emit(treasury, "EmergencyWithdraw");
    });

    it("Should handle pre-existing approvals correctly", async function () {
      const withdrawAmount = ethers.parseEther("2");
      
      // First signer approves
      await treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount);
      expect(await treasury.emergencyApprovalCount()).to.equal(1);
      
      // If we had removed a signer here (which we can't do while paused), 
      // their approval would still count until the emergency withdrawal completes
      
      // Second signer completes the emergency withdrawal
      await expect(
        treasury.connect(signer3).emergencyWithdrawETH(recipient.address, withdrawAmount)
      ).to.emit(treasury, "EmergencyWithdraw");
      
      // After completion, all approvals are reset
      expect(await treasury.emergencyApprovalCount()).to.equal(0);
      expect(await treasury.emergencyApprovals(await signer1.getAddress())).to.be.false;
      expect(await treasury.emergencyApprovals(await signer3.getAddress())).to.be.false;
    });
  });

  describe("Emergency Withdrawal Events", function () {
    beforeEach(async function () {
      // Pause the contract
      const pauseData = treasury.interface.encodeFunctionData("pause");
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        pauseData,
        "Pause for emergency"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs.find(e => e.fragment && e.fragment.name === "TransactionSubmitted").args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);
    });

    it("Should emit correct events in order", async function () {
      const withdrawAmount = ethers.parseEther("2");
      
      // First approval
      await expect(
        treasury.connect(signer1).emergencyWithdrawETH(recipient.address, withdrawAmount)
      ).to.emit(treasury, "EmergencyApprovalGranted")
        .withArgs(signer1.address);
      
      // Second approval and execution
      const tx = await treasury.connect(signer2).emergencyWithdrawETH(recipient.address, withdrawAmount);
      
      await expect(tx)
        .to.emit(treasury, "EmergencyApprovalGranted")
        .withArgs(await signer2.getAddress());
        
      await expect(tx)
        .to.emit(treasury, "EmergencyWithdraw")
        .withArgs(ethers.ZeroAddress, recipient.address, withdrawAmount);
    });
  });
});