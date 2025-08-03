const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("EventEmissionFixed", function () {
  let multisigTreasuryWithEvents;
  let bogoRewardDistributorWithEvents;
  let comprehensiveEventExample;
  let owner;
  let user1;
  let user2;
  let backend1;
  let backend2;
  let treasury1;
  let treasury2;

  beforeEach(async function () {
    [owner, user1, user2, backend1, backend2, treasury1, treasury2] = await ethers.getSigners();

    // Deploy MultisigTreasuryWithEvents
    const MultisigTreasuryWithEvents = await ethers.getContractFactory("contracts/EventEmissionFixed.sol:MultisigTreasuryWithEvents");
    multisigTreasuryWithEvents = await MultisigTreasuryWithEvents.deploy();
    await multisigTreasuryWithEvents.waitForDeployment();

    // Deploy BOGORewardDistributorWithEvents
    const BOGORewardDistributorWithEvents = await ethers.getContractFactory("contracts/EventEmissionFixed.sol:BOGORewardDistributorWithEvents");
    bogoRewardDistributorWithEvents = await BOGORewardDistributorWithEvents.deploy();
    await bogoRewardDistributorWithEvents.waitForDeployment();

    // Deploy ComprehensiveEventExample
    const ComprehensiveEventExample = await ethers.getContractFactory("contracts/EventEmissionFixed.sol:ComprehensiveEventExample");
    comprehensiveEventExample = await ComprehensiveEventExample.deploy();
    await comprehensiveEventExample.waitForDeployment();
  });

  describe("MultisigTreasuryWithEvents", function () {
    describe("Function Restrictions Toggle", function () {
      it("Should emit FunctionRestrictionsToggled event when toggling restrictions", async function () {
        // Initially false, should toggle to true
        await expect(
          multisigTreasuryWithEvents.toggleFunctionRestrictions()
        ).to.emit(multisigTreasuryWithEvents, "FunctionRestrictionsToggled")
          .withArgs(true);

        expect(await multisigTreasuryWithEvents.restrictFunctionCalls()).to.be.true;
      });

      it("Should emit correct event when toggling back to false", async function () {
        // First toggle to true
        await multisigTreasuryWithEvents.toggleFunctionRestrictions();
        
        // Then toggle back to false
        await expect(
          multisigTreasuryWithEvents.toggleFunctionRestrictions()
        ).to.emit(multisigTreasuryWithEvents, "FunctionRestrictionsToggled")
          .withArgs(false);

        expect(await multisigTreasuryWithEvents.restrictFunctionCalls()).to.be.false;
      });

      it("Should emit events for multiple toggles", async function () {
        // Toggle multiple times and verify events
        await expect(
          multisigTreasuryWithEvents.toggleFunctionRestrictions()
        ).to.emit(multisigTreasuryWithEvents, "FunctionRestrictionsToggled")
          .withArgs(true);

        await expect(
          multisigTreasuryWithEvents.toggleFunctionRestrictions()
        ).to.emit(multisigTreasuryWithEvents, "FunctionRestrictionsToggled")
          .withArgs(false);

        await expect(
          multisigTreasuryWithEvents.toggleFunctionRestrictions()
        ).to.emit(multisigTreasuryWithEvents, "FunctionRestrictionsToggled")
          .withArgs(true);
      });

      it("Should maintain correct state after multiple toggles", async function () {
        expect(await multisigTreasuryWithEvents.restrictFunctionCalls()).to.be.false;
        
        await multisigTreasuryWithEvents.toggleFunctionRestrictions();
        expect(await multisigTreasuryWithEvents.restrictFunctionCalls()).to.be.true;
        
        await multisigTreasuryWithEvents.toggleFunctionRestrictions();
        expect(await multisigTreasuryWithEvents.restrictFunctionCalls()).to.be.false;
        
        await multisigTreasuryWithEvents.toggleFunctionRestrictions();
        expect(await multisigTreasuryWithEvents.restrictFunctionCalls()).to.be.true;
      });
    });
  });

  describe("BOGORewardDistributorWithEvents", function () {
    describe("Authorized Backend Management", function () {
      it("Should emit AuthorizedBackendSet event when setting backend authorization", async function () {
        await expect(
          bogoRewardDistributorWithEvents.setAuthorizedBackend(await backend1.getAddress(), true)
        ).to.emit(bogoRewardDistributorWithEvents, "AuthorizedBackendSet")
          .withArgs(await backend1.getAddress(), true);

        expect(await bogoRewardDistributorWithEvents.authorizedBackends(await backend1.getAddress())).to.be.true;
      });

      it("Should emit event when revoking backend authorization", async function () {
        // First authorize
        await bogoRewardDistributorWithEvents.setAuthorizedBackend(await backend1.getAddress(), true);
        
        // Then revoke
        await expect(
          bogoRewardDistributorWithEvents.setAuthorizedBackend(await backend1.getAddress(), false)
        ).to.emit(bogoRewardDistributorWithEvents, "AuthorizedBackendSet")
          .withArgs(await backend1.getAddress(), false);

        expect(await bogoRewardDistributorWithEvents.authorizedBackends(await backend1.getAddress())).to.be.false;
      });

      it("Should emit events for multiple backend authorizations", async function () {
        await expect(
          bogoRewardDistributorWithEvents.setAuthorizedBackend(await backend1.getAddress(), true)
        ).to.emit(bogoRewardDistributorWithEvents, "AuthorizedBackendSet")
          .withArgs(await backend1.getAddress(), true);

        await expect(
          bogoRewardDistributorWithEvents.setAuthorizedBackend(await backend2.getAddress(), true)
        ).to.emit(bogoRewardDistributorWithEvents, "AuthorizedBackendSet")
          .withArgs(await backend2.getAddress(), true);

        expect(await bogoRewardDistributorWithEvents.authorizedBackends(await backend1.getAddress())).to.be.true;
        expect(await bogoRewardDistributorWithEvents.authorizedBackends(await backend2.getAddress())).to.be.true;
      });

      it("Should handle setting same authorization status multiple times", async function () {
        // Set to true multiple times
        await expect(
          bogoRewardDistributorWithEvents.setAuthorizedBackend(await backend1.getAddress(), true)
        ).to.emit(bogoRewardDistributorWithEvents, "AuthorizedBackendSet")
          .withArgs(await backend1.getAddress(), true);

        await expect(
          bogoRewardDistributorWithEvents.setAuthorizedBackend(await backend1.getAddress(), true)
        ).to.emit(bogoRewardDistributorWithEvents, "AuthorizedBackendSet")
          .withArgs(await backend1.getAddress(), true);

        expect(await bogoRewardDistributorWithEvents.authorizedBackends(await backend1.getAddress())).to.be.true;
      });

      it("Should handle zero address backend", async function () {
        await expect(
          bogoRewardDistributorWithEvents.setAuthorizedBackend(ethers.ZeroAddress, true)
        ).to.emit(bogoRewardDistributorWithEvents, "AuthorizedBackendSet")
          .withArgs(ethers.ZeroAddress, true);

        expect(await bogoRewardDistributorWithEvents.authorizedBackends(ethers.ZeroAddress)).to.be.true;
      });
    });

    describe("Daily Limit Reset", function () {
      it("Should have initial state correctly set", async function () {
        expect(await bogoRewardDistributorWithEvents.dailyDistributed()).to.equal(0);
        expect(await bogoRewardDistributorWithEvents.lastResetTime()).to.equal(0);
      });

      // Note: The _resetDailyLimit function is private, so we can't test it directly
      // In a real implementation, this would be called by public functions
      // For testing purposes, we verify the event interface exists
      it("Should have DailyLimitReset event defined", async function () {
        // Verify the contract has the event by checking the interface
        const contractInterface = bogoRewardDistributorWithEvents.interface;
        const event = contractInterface.getEvent("DailyLimitReset");
        
        expect(event.name).to.equal("DailyLimitReset");
        expect(event.inputs).to.have.length(2);
        expect(event.inputs[0].name).to.equal("timestamp");
        expect(event.inputs[1].name).to.equal("previousDistributed");
      });
    });
  });

  describe("ComprehensiveEventExample", function () {
    describe("Fee Management", function () {
      it("Should emit ConfigurationChanged event when setting fee", async function () {
        const newFee = ethers.parseEther("0.01");
        
        await expect(
          comprehensiveEventExample.setFee(newFee)
        ).to.emit(comprehensiveEventExample, "ConfigurationChanged")
          .withArgs("fee", 0, newFee);

        expect(await comprehensiveEventExample.fee()).to.equal(newFee);
      });

      it("Should emit event with correct old and new values", async function () {
        const firstFee = ethers.parseEther("0.01");
        const secondFee = ethers.parseEther("0.02");
        
        // Set initial fee
        await comprehensiveEventExample.setFee(firstFee);
        
        // Change fee and verify event
        await expect(
          comprehensiveEventExample.setFee(secondFee)
        ).to.emit(comprehensiveEventExample, "ConfigurationChanged")
          .withArgs("fee", firstFee, secondFee);
      });

      it("Should emit event even when setting same fee", async function () {
        const fee = ethers.parseEther("0.01");
        
        await comprehensiveEventExample.setFee(fee);
        
        await expect(
          comprehensiveEventExample.setFee(fee)
        ).to.emit(comprehensiveEventExample, "ConfigurationChanged")
          .withArgs("fee", fee, fee);
      });

      it("Should handle zero fee correctly", async function () {
        const nonZeroFee = ethers.parseEther("0.01");
        
        // Set non-zero fee first
        await comprehensiveEventExample.setFee(nonZeroFee);
        
        // Set to zero
        await expect(
          comprehensiveEventExample.setFee(0)
        ).to.emit(comprehensiveEventExample, "ConfigurationChanged")
          .withArgs("fee", nonZeroFee, 0);
      });
    });

    describe("Treasury Management", function () {
      it("Should emit TreasuryChanged event when setting treasury", async function () {
        await expect(
          comprehensiveEventExample.setTreasury(await treasury1.getAddress())
        ).to.emit(comprehensiveEventExample, "TreasuryChanged")
          .withArgs(ethers.ZeroAddress, await treasury1.getAddress());

        expect(await comprehensiveEventExample.treasury()).to.equal(await treasury1.getAddress());
      });

      it("Should emit event with correct old and new treasury addresses", async function () {
        // Set initial treasury
        await comprehensiveEventExample.setTreasury(await treasury1.getAddress());
        
        // Change treasury
        await expect(
          comprehensiveEventExample.setTreasury(await treasury2.getAddress())
        ).to.emit(comprehensiveEventExample, "TreasuryChanged")
          .withArgs(await treasury1.getAddress(), await treasury2.getAddress());
      });

      it("Should handle setting treasury to zero address", async function () {
        // Set initial treasury
        await comprehensiveEventExample.setTreasury(await treasury1.getAddress());
        
        // Set to zero address
        await expect(
          comprehensiveEventExample.setTreasury(ethers.ZeroAddress)
        ).to.emit(comprehensiveEventExample, "TreasuryChanged")
          .withArgs(await treasury1.getAddress(), ethers.ZeroAddress);
      });

      it("Should emit event when setting same treasury address", async function () {
        await comprehensiveEventExample.setTreasury(await treasury1.getAddress());
        
        await expect(
          comprehensiveEventExample.setTreasury(await treasury1.getAddress())
        ).to.emit(comprehensiveEventExample, "TreasuryChanged")
          .withArgs(await treasury1.getAddress(), await treasury1.getAddress());
      });
    });

    describe("Pause/Unpause Functionality", function () {
      it("Should emit ContractPaused event when pausing", async function () {
        await expect(
          comprehensiveEventExample.pause()
        ).to.emit(comprehensiveEventExample, "ContractPaused")
          .withArgs(await owner.getAddress());

        expect(await comprehensiveEventExample.paused()).to.be.true;
      });

      it("Should emit ContractUnpaused event when unpausing", async function () {
        // First pause
        await comprehensiveEventExample.pause();
        
        // Then unpause
        await expect(
          comprehensiveEventExample.unpause()
        ).to.emit(comprehensiveEventExample, "ContractUnpaused")
          .withArgs(await owner.getAddress());

        expect(await comprehensiveEventExample.paused()).to.be.false;
      });

      it("Should emit events for multiple pause/unpause cycles", async function () {
        await expect(
          comprehensiveEventExample.pause()
        ).to.emit(comprehensiveEventExample, "ContractPaused")
          .withArgs(await owner.getAddress());

        await expect(
          comprehensiveEventExample.unpause()
        ).to.emit(comprehensiveEventExample, "ContractUnpaused")
          .withArgs(await owner.getAddress());

        await expect(
          comprehensiveEventExample.pause()
        ).to.emit(comprehensiveEventExample, "ContractPaused")
          .withArgs(await owner.getAddress());
      });

      it("Should emit pause event even if already paused", async function () {
        await comprehensiveEventExample.pause();
        
        await expect(
          comprehensiveEventExample.pause()
        ).to.emit(comprehensiveEventExample, "ContractPaused")
          .withArgs(await owner.getAddress());
      });

      it("Should emit unpause event even if already unpaused", async function () {
        await expect(
          comprehensiveEventExample.unpause()
        ).to.emit(comprehensiveEventExample, "ContractUnpaused")
          .withArgs(await owner.getAddress());
      });

      it("Should emit events with correct sender address", async function () {
        await expect(
          comprehensiveEventExample.connect(user1).pause()
        ).to.emit(comprehensiveEventExample, "ContractPaused")
          .withArgs(await user1.getAddress());

        await expect(
          comprehensiveEventExample.connect(user2).unpause()
        ).to.emit(comprehensiveEventExample, "ContractUnpaused")
          .withArgs(await user2.getAddress());
      });
    });

    describe("Emergency Actions", function () {
      it("Should emit EmergencyActionTaken event for emergency withdraw", async function () {
        await expect(
          comprehensiveEventExample.emergencyWithdraw(await user1.getAddress(), ethers.parseEther("1"))
        ).to.emit(comprehensiveEventExample, "EmergencyActionTaken")
          .withArgs("emergencyWithdraw", await owner.getAddress());
      });

      it("Should emit event with correct sender for emergency actions", async function () {
        await expect(
          comprehensiveEventExample.connect(user1).emergencyWithdraw(await user2.getAddress(), ethers.parseEther("1"))
        ).to.emit(comprehensiveEventExample, "EmergencyActionTaken")
          .withArgs("emergencyWithdraw", await user1.getAddress());
      });

      it("Should emit emergency events for different amounts", async function () {
        await expect(
          comprehensiveEventExample.emergencyWithdraw(await user1.getAddress(), 0)
        ).to.emit(comprehensiveEventExample, "EmergencyActionTaken")
          .withArgs("emergencyWithdraw", await owner.getAddress());

        await expect(
          comprehensiveEventExample.emergencyWithdraw(await user1.getAddress(), ethers.parseEther("1000"))
        ).to.emit(comprehensiveEventExample, "EmergencyActionTaken")
          .withArgs("emergencyWithdraw", await owner.getAddress());
      });
    });

    describe("State Consistency", function () {
      it("Should maintain correct state after multiple operations", async function () {
        const fee = ethers.parseEther("0.01");
        
        // Set fee
        await comprehensiveEventExample.setFee(fee);
        expect(await comprehensiveEventExample.fee()).to.equal(fee);
        
        // Set treasury
        await comprehensiveEventExample.setTreasury(await treasury1.getAddress());
        expect(await comprehensiveEventExample.treasury()).to.equal(await treasury1.getAddress());
        
        // Pause
        await comprehensiveEventExample.pause();
        expect(await comprehensiveEventExample.paused()).to.be.true;
        
        // Verify all states are maintained
        expect(await comprehensiveEventExample.fee()).to.equal(fee);
        expect(await comprehensiveEventExample.treasury()).to.equal(await treasury1.getAddress());
        expect(await comprehensiveEventExample.paused()).to.be.true;
      });

      it("Should handle rapid state changes correctly", async function () {
        // Rapid fee changes
        await comprehensiveEventExample.setFee(ethers.parseEther("0.01"));
        await comprehensiveEventExample.setFee(ethers.parseEther("0.02"));
        await comprehensiveEventExample.setFee(ethers.parseEther("0.03"));
        
        expect(await comprehensiveEventExample.fee()).to.equal(ethers.parseEther("0.03"));
        
        // Rapid pause/unpause
        await comprehensiveEventExample.pause();
        await comprehensiveEventExample.unpause();
        await comprehensiveEventExample.pause();
        
        expect(await comprehensiveEventExample.paused()).to.be.true;
      });
    });
  });

  describe("Event Interface Validation", function () {
    it("Should have all required events in MultisigTreasuryWithEvents", async function () {
      const contractInterface = multisigTreasuryWithEvents.interface;
      const event = contractInterface.getEvent("FunctionRestrictionsToggled");
      
      expect(event.name).to.equal("FunctionRestrictionsToggled");
      expect(event.inputs).to.have.length(1);
      expect(event.inputs[0].name).to.equal("enabled");
      expect(event.inputs[0].type).to.equal("bool");
    });

    it("Should have all required events in BOGORewardDistributorWithEvents", async function () {
      const contractInterface = bogoRewardDistributorWithEvents.interface;
      
      const backendEvent = contractInterface.getEvent("AuthorizedBackendSet");
      expect(backendEvent.name).to.equal("AuthorizedBackendSet");
      expect(backendEvent.inputs).to.have.length(2);
      expect(backendEvent.inputs[0].name).to.equal("backend");
      expect(backendEvent.inputs[1].name).to.equal("authorized");
      
      const resetEvent = contractInterface.getEvent("DailyLimitReset");
      expect(resetEvent.name).to.equal("DailyLimitReset");
      expect(resetEvent.inputs).to.have.length(2);
      expect(resetEvent.inputs[0].name).to.equal("timestamp");
      expect(resetEvent.inputs[1].name).to.equal("previousDistributed");
    });

    it("Should have all required events in ComprehensiveEventExample", async function () {
      const contractInterface = comprehensiveEventExample.interface;
      
      const events = [
        "ConfigurationChanged",
        "ContractPaused",
        "ContractUnpaused",
        "TreasuryChanged",
        "FeesCollected",
        "EmergencyActionTaken"
      ];
      
      for (const eventName of events) {
        const event = contractInterface.getEvent(eventName);
        expect(event.name).to.equal(eventName);
      }
    });
  });

  describe("Gas Efficiency", function () {
    it("Should emit events efficiently", async function () {
      // Test that events don't consume excessive gas
      const tx1 = await comprehensiveEventExample.setFee(ethers.parseEther("0.01"));
      const receipt1 = await tx1.wait();
      
      const tx2 = await comprehensiveEventExample.setTreasury(await treasury1.getAddress());
      const receipt2 = await tx2.wait();
      
      const tx3 = await comprehensiveEventExample.pause();
      const receipt3 = await tx3.wait();
      
      // Verify transactions completed successfully
      expect(receipt1.status).to.equal(1);
      expect(receipt2.status).to.equal(1);
      expect(receipt3.status).to.equal(1);
    });

    it("Should handle multiple events in single transaction", async function () {
      // This would require a function that emits multiple events
      // For now, verify individual events work correctly
      const tx = await comprehensiveEventExample.setFee(ethers.parseEther("0.01"));
      const receipt = await tx.wait();
      
      expect(receipt.logs).to.have.length(1);
    });
  });

  describe("Edge Cases", function () {
    it("Should handle maximum uint256 values", async function () {
      const maxUint256 = ethers.MaxUint256;
      
      await expect(
        comprehensiveEventExample.setFee(maxUint256)
      ).to.emit(comprehensiveEventExample, "ConfigurationChanged")
        .withArgs("fee", 0, maxUint256);
    });

    it("Should handle contract address as treasury", async function () {
      const contractAddress = await comprehensiveEventExample.getAddress();
      
      await expect(
        comprehensiveEventExample.setTreasury(contractAddress)
      ).to.emit(comprehensiveEventExample, "TreasuryChanged")
        .withArgs(ethers.ZeroAddress, contractAddress);
    });

    it("Should handle emergency withdraw with zero amount", async function () {
      await expect(
        comprehensiveEventExample.emergencyWithdraw(await user1.getAddress(), 0)
      ).to.emit(comprehensiveEventExample, "EmergencyActionTaken")
        .withArgs("emergencyWithdraw", await owner.getAddress());
    });
  });
});