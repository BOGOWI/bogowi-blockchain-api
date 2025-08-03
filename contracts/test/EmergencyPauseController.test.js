const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture, time } = require("@nomicfoundation/hardhat-network-helpers");

describe("EmergencyPauseController", function () {
  async function deployEmergencyPauseControllerFixture() {
    const [owner, guardian1, guardian2, guardian3, guardian4, manager, nonGuardian] = await ethers.getSigners();
    
    // Deploy mock pausable contracts
    const MockPausable = await ethers.getContractFactory("contracts/mocks/MockPausable.sol:MockPausable");
    const pausable1 = await MockPausable.deploy("Pausable1");
    await pausable1.waitForDeployment();
    const pausable2 = await MockPausable.deploy("Pausable2");
    await pausable2.waitForDeployment();
    const pausable3 = await MockPausable.deploy("Pausable3");
    await pausable3.waitForDeployment();
    
    // Deploy EmergencyPauseController
    const EmergencyPauseController = await ethers.getContractFactory("EmergencyPauseController");
    const guardians = [await guardian1.getAddress(), await guardian2.getAddress(), await guardian3.getAddress()];
    const controller = await EmergencyPauseController.deploy(guardians, await manager.getAddress());
    await controller.waitForDeployment();
    
    // Grant PAUSER_ROLE to controller on pausable contracts
    const PAUSER_ROLE = await pausable1.PAUSER_ROLE();
    await pausable1.grantRole(PAUSER_ROLE, await controller.getAddress());
    await pausable2.grantRole(PAUSER_ROLE, await controller.getAddress());
    await pausable3.grantRole(PAUSER_ROLE, await controller.getAddress());
    
    // Add contracts to controller
    await controller.connect(manager).addContract(await pausable1.getAddress(), "Pausable Contract 1");
    await controller.connect(manager).addContract(await pausable2.getAddress(), "Pausable Contract 2");
    await controller.connect(manager).addContract(await pausable3.getAddress(), "Pausable Contract 3");
    
    return {
      controller,
      pausable1,
      pausable2,
      pausable3,
      owner,
      guardian1,
      guardian2,
      guardian3,
      guardian4,
      manager,
      nonGuardian,
      guardians
    };
  }
  
  describe("Deployment", function () {
    it("Should deploy with correct initial state", async function () {
      const { controller, guardians, manager } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      expect(await controller.requiredConfirmations()).to.equal(2);
      expect(await controller.proposalCount()).to.equal(0);
      
      // Check guardian roles
      const GUARDIAN_ROLE = await controller.GUARDIAN_ROLE();
      for (const guardian of guardians) {
        expect(await controller.hasRole(GUARDIAN_ROLE, guardian)).to.be.true;
      }
      
      // Check manager role
      const MANAGER_ROLE = await controller.MANAGER_ROLE();
      expect(await controller.hasRole(MANAGER_ROLE, await manager.getAddress())).to.be.true;
      
      // Check contract tracking - contracts should be added during fixture setup
      const [contracts, paused, names] = await controller.getContractStatuses();
      expect(contracts).to.have.lengthOf(3);
    });
    
    it("Should revert with invalid parameters", async function () {
      const EmergencyPauseController = await ethers.getContractFactory("EmergencyPauseController");
      const [, guardian1, guardian2, manager] = await ethers.getSigners();
      
      // Too few guardians
      await expect(
        EmergencyPauseController.deploy([await guardian1.getAddress(), await guardian2.getAddress()], await manager.getAddress())
      ).to.be.revertedWith("INVALID_PARAMETER");
      
      // Zero address manager
      await expect(
        EmergencyPauseController.deploy([await guardian1.getAddress(), await guardian2.getAddress(), await manager.getAddress()], ethers.ZeroAddress)
      ).to.be.revertedWith("ZERO_ADDRESS");
      
      // Zero address guardian
      await expect(
        EmergencyPauseController.deploy([await guardian1.getAddress(), ethers.ZeroAddress, await manager.getAddress()], await manager.getAddress())
      ).to.be.revertedWith("ZERO_ADDRESS");
    });
  });
  
  describe("Contract Management", function () {
    it("Should add contract successfully", async function () {
      const { controller, manager } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const MockPausable = await ethers.getContractFactory("contracts/mocks/MockPausable.sol:MockPausable");
      const newPausable = await MockPausable.deploy("NewPausable");
      await newPausable.waitForDeployment();
      
      // Grant PAUSER_ROLE to controller
      const PAUSER_ROLE = await newPausable.PAUSER_ROLE();
      await newPausable.grantRole(PAUSER_ROLE, await controller.getAddress());
      
      await expect(
        controller.connect(manager).addContract(await newPausable.getAddress(), "New Pausable Contract")
      ).to.emit(controller, "ContractAdded")
        .withArgs(await newPausable.getAddress(), "New Pausable Contract");
      
      expect(await controller.isTrackedContract(await newPausable.getAddress())).to.be.true;
      expect(await controller.contractNames(await newPausable.getAddress())).to.equal("New Pausable Contract");
    });
    
    it("Should revert adding contract with invalid parameters", async function () {
      const { controller, manager, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Zero address
      await expect(
        controller.connect(manager).addContract(ethers.ZeroAddress, "Test")
      ).to.be.revertedWith("ZERO_ADDRESS");
      
      // Already exists
      await expect(
        controller.connect(manager).addContract(await pausable1.getAddress(), "Duplicate")
      ).to.be.revertedWith("ALREADY_EXISTS");
      
      // Empty name
      const MockPausable = await ethers.getContractFactory("contracts/mocks/MockPausable.sol:MockPausable");
      const newPausable = await MockPausable.deploy("TestPausable");
      await newPausable.waitForDeployment();
      // Grant PAUSER_ROLE to controller
      const PAUSER_ROLE = await newPausable.PAUSER_ROLE();
      await newPausable.grantRole(PAUSER_ROLE, await controller.getAddress());
      
      await expect(
        controller.connect(manager).addContract(await newPausable.getAddress(), "")
      ).to.be.revertedWith("EMPTY_STRING");
      
      // Invalid contract (doesn't implement pause interface)
      const MockERC20 = await ethers.getContractFactory("contracts/mocks/MockERC20.sol:MockERC20");
      const token = await MockERC20.deploy("Test", "TEST", 1000);
      await token.waitForDeployment();
      await expect(
        controller.connect(manager).addContract(await token.getAddress(), "Invalid Contract")
      ).to.be.revertedWith("INVALID_ADDRESS");
    });
    
    it("Should remove contract successfully", async function () {
      const { controller, manager, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      await expect(
        controller.connect(manager).removeContract(await pausable1.getAddress())
      ).to.emit(controller, "ContractRemoved")
        .withArgs(await pausable1.getAddress());
      
      expect(await controller.isTrackedContract(await pausable1.getAddress())).to.be.false;
    });
    
    it("Should revert removing non-existent contract", async function () {
      const { controller, manager } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const randomAddress = await ethers.Wallet.createRandom().getAddress();
      await expect(
        controller.connect(manager).removeContract(randomAddress)
      ).to.be.revertedWith("DOES_NOT_EXIST");
    });
    
    it("Should revert when not manager", async function () {
      const { controller, guardian1, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const MANAGER_ROLE = await controller.MANAGER_ROLE();
      await expect(
        controller.connect(guardian1).addContract(await pausable1.getAddress(), "Test")
      ).to.be.revertedWithCustomError(controller, "AccessControlUnauthorizedAccount")
        .withArgs(await guardian1.getAddress(), MANAGER_ROLE);
    });
  });
  
  describe("Pause Proposals", function () {
    it("Should create pause proposal successfully", async function () {
      const { controller, guardian1, pausable1, pausable2 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const contracts = [await pausable1.getAddress(), await pausable2.getAddress()];
      const reason = "Security vulnerability detected";
      
      await expect(
        controller.connect(guardian1).proposePause(contracts, reason)
      ).to.emit(controller, "PauseProposalCreated")
        .withArgs(0, await guardian1.getAddress(), true);
      
      const proposal = await controller.proposals(0);
      expect(proposal.proposer).to.equal(await guardian1.getAddress());
      expect(proposal.reason).to.equal(reason);
      expect(proposal.confirmations).to.equal(1);
      expect(proposal.executed).to.be.false;
      expect(proposal.isPause).to.be.true;
      
      expect(await controller.hasConfirmed(0, await guardian1.getAddress())).to.be.true;
    });
    
    it("Should create unpause proposal successfully", async function () {
      const { controller, guardian1, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const contracts = [await pausable1.getAddress()];
      const reason = "Issue resolved";
      
      await expect(
        controller.connect(guardian1).proposeUnpause(contracts, reason)
      ).to.emit(controller, "PauseProposalCreated")
        .withArgs(0, await guardian1.getAddress(), false);
      
      const proposal = await controller.proposals(0);
      expect(proposal.isPause).to.be.false;
    });
    
    it("Should auto-execute with single confirmation requirement", async function () {
      const { controller, guardian1, pausable1, owner } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Set required confirmations to 1
      await controller.connect(owner).updateRequiredConfirmations(1);
      
      const contracts = [await pausable1.getAddress()];
      const reason = "Emergency pause";
      
      await expect(
        controller.connect(guardian1).proposePause(contracts, reason)
      ).to.emit(controller, "EmergencyPauseExecuted")
        .withArgs(await guardian1.getAddress(), contracts, reason);
      
      expect(await pausable1.paused()).to.be.true;
      
      const proposal = await controller.proposals(0);
      expect(proposal.executed).to.be.true;
    });
    
    it("Should revert with invalid parameters", async function () {
      const { controller, guardian1, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Empty contracts array
      await expect(
        controller.connect(guardian1).proposePause([], "Test reason")
      ).to.be.revertedWith("INVALID_LENGTH");
      
      // Empty reason
      await expect(
        controller.connect(guardian1).proposePause([await pausable1.getAddress()], "")
      ).to.be.revertedWith("EMPTY_STRING");
      
      // Non-tracked contract
      const randomAddress = await ethers.Wallet.createRandom().getAddress();
      await expect(
        controller.connect(guardian1).proposePause([randomAddress], "Test reason")
      ).to.be.revertedWith("DOES_NOT_EXIST");
    });
    
    it("Should revert when not guardian", async function () {
      const { controller, nonGuardian, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const GUARDIAN_ROLE = await controller.GUARDIAN_ROLE();
      await expect(
        controller.connect(nonGuardian).proposePause([await pausable1.getAddress()], "Test reason")
      ).to.be.revertedWithCustomError(controller, "AccessControlUnauthorizedAccount")
        .withArgs(await nonGuardian.getAddress(), GUARDIAN_ROLE);
    });
  });
  
  describe("Proposal Confirmation", function () {
    it("Should confirm proposal and execute when threshold reached", async function () {
      const { controller, guardian1, guardian2, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const contracts = [await pausable1.getAddress()];
      const reason = "Security issue";
      
      // Create proposal
      await controller.connect(guardian1).proposePause(contracts, reason);
      
      // Confirm and execute
      await expect(
        controller.connect(guardian2).confirmProposal(0)
      ).to.emit(controller, "ProposalConfirmed")
        .withArgs(0, await guardian2.getAddress())
        .and.to.emit(controller, "EmergencyPauseExecuted")
        .withArgs(await guardian2.getAddress(), contracts, reason);
      
      expect(await pausable1.paused()).to.be.true;
      
      const proposal = await controller.proposals(0);
      expect(proposal.confirmations).to.equal(2);
      expect(proposal.executed).to.be.true;
    });
    
    it("Should confirm without executing when threshold not reached", async function () {
      const { controller, guardian1, guardian2, guardian3, pausable1, owner } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Set required confirmations to 3
      await controller.connect(owner).updateRequiredConfirmations(3);
      
      const contracts = [await pausable1.getAddress()];
      const reason = "Security issue";
      
      // Create proposal
      await controller.connect(guardian1).proposePause(contracts, reason);
      
      // Confirm but don't execute
      await expect(
        controller.connect(guardian2).confirmProposal(0)
      ).to.emit(controller, "ProposalConfirmed")
        .withArgs(0, await guardian2.getAddress())
        .and.to.not.emit(controller, "EmergencyPauseExecuted");
      
      expect(await pausable1.paused()).to.be.false;
      
      const proposal = await controller.proposals(0);
      expect(proposal.confirmations).to.equal(2);
      expect(proposal.executed).to.be.false;
    });
    
    it("Should revert double confirmation", async function () {
      const { controller, guardian1, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      await controller.connect(guardian1).proposePause([await pausable1.getAddress()], "Test reason");
      
      await expect(
        controller.connect(guardian1).confirmProposal(0)
      ).to.be.revertedWith("ALREADY_PROCESSED");
    });
    
    it("Should revert confirming executed proposal", async function () {
      const { controller, guardian1, guardian2, guardian3, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      await controller.connect(guardian1).proposePause([await pausable1.getAddress()], "Test reason");
      await controller.connect(guardian2).confirmProposal(0);
      
      await expect(
        controller.connect(guardian3).confirmProposal(0)
      ).to.be.revertedWith("ALREADY_PROCESSED");
    });
    
    it("Should revert confirming non-existent proposal", async function () {
      const { controller, guardian1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      await expect(
        controller.connect(guardian1).confirmProposal(999)
      ).to.be.revertedWith("DOES_NOT_EXIST");
    });
  });
  
  describe("Emergency Pause All", function () {
    it("Should create proposal to pause all contracts", async function () {
      const { controller, guardian1, pausable1, pausable2, pausable3 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const reason = "Critical system vulnerability";
      
      await expect(
        controller.connect(guardian1).emergencyPauseAll(reason)
      ).to.emit(controller, "PauseProposalCreated")
        .withArgs(0, await guardian1.getAddress(), true);
      
      const proposal = await controller.proposals(0);
      expect(proposal.reason).to.equal(reason);
      expect(proposal.isPause).to.be.true;
      expect(proposal.executed).to.be.false;
      expect(proposal.proposer).to.equal(await guardian1.getAddress());
    });
    
    it("Should revert when no contracts tracked", async function () {
      const { controller, guardian1, manager, pausable1, pausable2, pausable3 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Remove all contracts
      await controller.connect(manager).removeContract(await pausable1.getAddress());
      await controller.connect(manager).removeContract(await pausable2.getAddress());
      await controller.connect(manager).removeContract(await pausable3.getAddress());
      
      await expect(
        controller.connect(guardian1).emergencyPauseAll("Test reason")
      ).to.be.revertedWith("INVALID_STATE");
    });
  });
  
  describe("Pause Expiry", function () {
    it("Should set pause expiry when pausing", async function () {
      const { controller, guardian1, guardian2, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      await controller.connect(guardian1).proposePause([await pausable1.getAddress()], "Test reason");
      await controller.connect(guardian2).confirmProposal(0);
      
      const expiry = await controller.pauseExpiry(await pausable1.getAddress());
      const currentTime = await time.latest();
      const maxDuration = await controller.MAX_PAUSE_DURATION();
      
      expect(expiry).to.be.closeTo(currentTime + Number(maxDuration), 5);
    });
    
    it("Should clear pause expiry when unpausing", async function () {
      const { controller, guardian1, guardian2, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Pause first
      await controller.connect(guardian1).proposePause([await pausable1.getAddress()], "Test reason");
      await controller.connect(guardian2).confirmProposal(0);
      
      // Then unpause
      await controller.connect(guardian1).proposeUnpause([await pausable1.getAddress()], "Issue resolved");
      await controller.connect(guardian2).confirmProposal(1);
      
      expect(await controller.pauseExpiry(await pausable1.getAddress())).to.equal(0);
      expect(await pausable1.paused()).to.be.false;
    });
    
    it("Should auto-expire pauses after max duration", async function () {
      const { controller, guardian1, guardian2, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Pause contract
      await controller.connect(guardian1).proposePause([await pausable1.getAddress()], "Test reason");
      await controller.connect(guardian2).confirmProposal(0);
      
      expect(await pausable1.paused()).to.be.true;
      
      // Fast forward past max duration
      const maxDuration = await controller.MAX_PAUSE_DURATION();
      await time.increase(Number(maxDuration) + 1);
      
      // Check and expire pauses
      await expect(
        controller.checkAndExpirePauses()
      ).to.emit(controller, "PauseExpired")
        .withArgs(await pausable1.getAddress());
      
      expect(await pausable1.paused()).to.be.false;
      expect(await controller.pauseExpiry(await pausable1.getAddress())).to.equal(0);
    });
  });
  
  describe("Configuration Updates", function () {
    it("Should update required confirmations", async function () {
      const { controller, owner } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      await expect(
        controller.connect(owner).updateRequiredConfirmations(3)
      ).to.emit(controller, "RequiredConfirmationsUpdated")
        .withArgs(2, 3);
      
      expect(await controller.requiredConfirmations()).to.equal(3);
    });
    
    it("Should revert invalid confirmation updates", async function () {
      const { controller, owner } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Zero confirmations
      await expect(
        controller.connect(owner).updateRequiredConfirmations(0)
      ).to.be.revertedWith("ZERO_AMOUNT");
      
      // More than available guardians
      await expect(
        controller.connect(owner).updateRequiredConfirmations(5)
      ).to.be.revertedWith("EXCEEDS_LIMIT");
    });
    
    it("Should revert when not admin", async function () {
      const { controller, guardian1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const DEFAULT_ADMIN_ROLE = await controller.DEFAULT_ADMIN_ROLE();
      await expect(
        controller.connect(guardian1).updateRequiredConfirmations(3)
      ).to.be.revertedWithCustomError(controller, "AccessControlUnauthorizedAccount")
        .withArgs(await guardian1.getAddress(), DEFAULT_ADMIN_ROLE);
    });
  });
  
  describe("View Functions", function () {
    it("Should return contract statuses", async function () {
      const { controller, guardian1, guardian2, pausable1, pausable2, pausable3 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Pause one contract
      await controller.connect(guardian1).proposePause([await pausable1.getAddress()], "Test reason");
      await controller.connect(guardian2).confirmProposal(0);
      
      const [contracts, paused, names] = await controller.getContractStatuses();
      
      expect(contracts).to.have.lengthOf(3);
      expect(paused).to.have.lengthOf(3);
      expect(names).to.have.lengthOf(3);
      
      const pausable1Address = await pausable1.getAddress();
      const pausable1Index = contracts.findIndex(addr => addr === pausable1Address);
      expect(paused[pausable1Index]).to.be.true;
      expect(names[pausable1Index]).to.equal("Pausable Contract 1");
    });
    
    it("Should return role count", async function () {
      const { controller } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const GUARDIAN_ROLE = await controller.GUARDIAN_ROLE();
      const MANAGER_ROLE = await controller.MANAGER_ROLE();
      
      expect(await controller.getRoleCount(GUARDIAN_ROLE)).to.equal(3);
      expect(await controller.getRoleCount(MANAGER_ROLE)).to.equal(1);
    });
    
    it("Should return pause history", async function () {
      const { controller, guardian1, guardian2, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Create some pause events
      await controller.connect(guardian1).proposePause([await pausable1.getAddress()], "First pause");
      await controller.connect(guardian2).confirmProposal(0);
      
      await controller.connect(guardian1).proposeUnpause([await pausable1.getAddress()], "First unpause");
      await controller.connect(guardian2).confirmProposal(1);
      
      const history = await controller.getPauseHistory(10);
      expect(history).to.have.lengthOf(2);
      
      expect(history[0].isPause).to.be.true;
      expect(history[0].reason).to.equal("First pause");
      expect(history[1].isPause).to.be.false;
      expect(history[1].reason).to.equal("First unpause");
    });
    
    it("Should limit pause history results", async function () {
      const { controller, guardian1, guardian2, pausable1 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Create multiple pause events
      for (let i = 0; i < 5; i++) {
        await controller.connect(guardian1).proposePause([await pausable1.getAddress()], `Pause ${i}`);
        await controller.connect(guardian2).confirmProposal(i * 2);
        
        await controller.connect(guardian1).proposeUnpause([await pausable1.getAddress()], `Unpause ${i}`);
        await controller.connect(guardian2).confirmProposal(i * 2 + 1);
      }
      
      const limitedHistory = await controller.getPauseHistory(3);
      expect(limitedHistory).to.have.lengthOf(3);
    });
  });
  
  describe("Role Management", function () {
    it("Should add new guardian", async function () {
      const { controller, owner, guardian4 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const GUARDIAN_ROLE = await controller.GUARDIAN_ROLE();
      
      await controller.connect(owner).grantRole(GUARDIAN_ROLE, await guardian4.getAddress());
      
      expect(await controller.hasRole(GUARDIAN_ROLE, await guardian4.getAddress())).to.be.true;
      expect(await controller.getRoleCount(GUARDIAN_ROLE)).to.equal(4);
    });
    
    it("Should remove guardian", async function () {
      const { controller, owner, guardian3 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const GUARDIAN_ROLE = await controller.GUARDIAN_ROLE();
      
      await controller.connect(owner).revokeRole(GUARDIAN_ROLE, await guardian3.getAddress());
      
      expect(await controller.hasRole(GUARDIAN_ROLE, await guardian3.getAddress())).to.be.false;
      expect(await controller.getRoleCount(GUARDIAN_ROLE)).to.equal(2);
    });
    
    it("Should update required confirmations when guardians change", async function () {
      const { controller, owner, guardian3, guardian4 } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const GUARDIAN_ROLE = await controller.GUARDIAN_ROLE();
      
      // Add guardian
      await controller.connect(owner).grantRole(GUARDIAN_ROLE, await guardian4.getAddress());
      await controller.connect(owner).updateRequiredConfirmations(3);
      
      // Remove guardian - should still work with 3 guardians
      await controller.connect(owner).revokeRole(GUARDIAN_ROLE, await guardian3.getAddress());
      
      // But can't require more confirmations than guardians
      await expect(
        controller.connect(owner).updateRequiredConfirmations(4)
      ).to.be.revertedWith("EXCEEDS_LIMIT");
    });
  });
  
  describe("Edge Cases and Error Handling", function () {
    it("Should handle contract pause failures gracefully", async function () {
      const { controller, guardian1, guardian2, manager } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      // Deploy a contract that will fail to pause
      const MockFailingPausable = await ethers.getContractFactory("contracts/mocks/MockFailingPausable.sol:MockFailingPausable");
      const failingContract = await MockFailingPausable.deploy();
      await failingContract.waitForDeployment();
      
      await controller.connect(manager).addContract(await failingContract.getAddress(), "Failing Contract");
      
      // Should not revert even if one contract fails
      await controller.connect(guardian1).proposePause([await failingContract.getAddress()], "Test reason");
      await expect(
        controller.connect(guardian2).confirmProposal(0)
      ).to.emit(controller, "EmergencyPauseExecuted");
    });
    
    it("Should handle mixed success/failure in batch operations", async function () {
      const { controller, guardian1, guardian2, pausable1, manager } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const MockFailingPausable = await ethers.getContractFactory("contracts/mocks/MockFailingPausable.sol:MockFailingPausable");
      const failingContract = await MockFailingPausable.deploy();
      
      await controller.connect(manager).addContract(await failingContract.getAddress(), "Failing Contract");
      
      // Mix of working and failing contracts
      const contracts = [await pausable1.getAddress(), await failingContract.getAddress()];
      
      await controller.connect(guardian1).proposePause(contracts, "Mixed test");
      await controller.connect(guardian2).confirmProposal(0);
      
      // Working contract should be paused
      expect(await pausable1.paused()).to.be.true;
    });
    
    it("Should handle pause expiry check with failing contracts", async function () {
      const { controller, guardian1, guardian2, pausable1, manager } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const MockFailingPausable = await ethers.getContractFactory("contracts/mocks/MockFailingPausable.sol:MockFailingPausable");
      const failingContract = await MockFailingPausable.deploy();
      
      await controller.connect(manager).addContract(await failingContract.getAddress(), "Failing Contract");
      
      // Pause both contracts
      await controller.connect(guardian1).proposePause([await pausable1.getAddress(), await failingContract.getAddress()], "Test");
      await controller.connect(guardian2).confirmProposal(0);
      
      // Fast forward past expiry
      const maxDuration = await controller.MAX_PAUSE_DURATION();
      await time.increase(Number(maxDuration) + 1);
      
      // Should not revert even if one contract fails to unpause
      await controller.checkAndExpirePauses();
      
      expect(await pausable1.paused()).to.be.false;
    });
    
    it("Should handle contract status check with failing contracts", async function () {
      const { controller, manager } = await loadFixture(deployEmergencyPauseControllerFixture);
      
      const MockFailingPausable = await ethers.getContractFactory("contracts/mocks/MockFailingPausable.sol:MockFailingPausable");
      const failingContract = await MockFailingPausable.deploy();
      
      await controller.connect(manager).addContract(await failingContract.getAddress(), "Failing Contract");
      
      const [contracts, paused, names] = await controller.getContractStatuses();
      
      // Should return false for failing contracts
      const failingContractAddress = await failingContract.getAddress();
      const failingIndex = contracts.findIndex(addr => addr === failingContractAddress);
      expect(paused[failingIndex]).to.be.false;
      expect(names[failingIndex]).to.equal("Failing Contract");
    });
  });
});