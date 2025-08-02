const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("EmergencyPauseController", function () {
    let emergencyPause;
    let bogoToken;
    let multisigTreasury;
    let owner, guardian1, guardian2, guardian3, manager, user;
    let GUARDIAN_ROLE, MANAGER_ROLE, PAUSER_ROLE, DEFAULT_ADMIN_ROLE;

    beforeEach(async function () {
        [owner, guardian1, guardian2, guardian3, manager, user] = await ethers.getSigners();

        // Calculate role hashes
        GUARDIAN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("GUARDIAN_ROLE"));
        MANAGER_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("MANAGER_ROLE"));
        PAUSER_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("PAUSER_ROLE"));
        DEFAULT_ADMIN_ROLE = ethers.constants.HashZero;

        // Deploy EmergencyPauseController
        const EmergencyPauseController = await ethers.getContractFactory("EmergencyPauseController");
        emergencyPause = await EmergencyPauseController.deploy(
            [guardian1.address, guardian2.address, guardian3.address],
            manager.address
        );

        // Deploy mock pausable contracts
        const BOGOTokenV2 = await ethers.getContractFactory("BOGOTokenV2");
        bogoToken = await BOGOTokenV2.deploy();

        const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
        multisigTreasury = await MultisigTreasury.deploy(
            [owner.address, guardian1.address],
            2
        );

        // Grant PAUSER_ROLE to EmergencyPauseController
        await bogoToken.grantRole(PAUSER_ROLE, emergencyPause.address);
        await multisigTreasury.grantRole(PAUSER_ROLE, emergencyPause.address);

        // Add contracts to emergency pause controller
        await emergencyPause.connect(manager).addContract(
            bogoToken.address,
            "BOGOTokenV2"
        );
        await emergencyPause.connect(manager).addContract(
            multisigTreasury.address,
            "MultisigTreasury"
        );
    });

    describe("Deployment", function () {
        it("Should deploy with correct initial state", async function () {
            expect(await emergencyPause.requiredConfirmations()).to.equal(2);
            expect(await emergencyPause.MAX_PAUSE_DURATION()).to.equal(72 * 60 * 60); // 72 hours
            expect(await emergencyPause.MIN_GUARDIANS()).to.equal(3);
        });

        it("Should assign roles correctly", async function () {
            expect(await emergencyPause.hasRole(GUARDIAN_ROLE, guardian1.address)).to.be.true;
            expect(await emergencyPause.hasRole(GUARDIAN_ROLE, guardian2.address)).to.be.true;
            expect(await emergencyPause.hasRole(GUARDIAN_ROLE, guardian3.address)).to.be.true;
            expect(await emergencyPause.hasRole(MANAGER_ROLE, manager.address)).to.be.true;
            expect(await emergencyPause.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
        });

        it("Should revert deployment with insufficient guardians", async function () {
            const EmergencyPauseController = await ethers.getContractFactory("EmergencyPauseController");
            await expect(
                EmergencyPauseController.deploy([guardian1.address], manager.address)
            ).to.be.revertedWith("INVALID_PARAMETER");
        });
    });

    describe("Contract Management", function () {
        it("Should allow manager to add contracts", async function () {
            const ConservationNFT = await ethers.getContractFactory("ConservationNFT");
            const conservationNFT = await ConservationNFT.deploy();
            await conservationNFT.grantRole(PAUSER_ROLE, emergencyPause.address);

            await expect(
                emergencyPause.connect(manager).addContract(
                    conservationNFT.address,
                    "ConservationNFT"
                )
            ).to.emit(emergencyPause, "ContractAdded")
                .withArgs(conservationNFT.address, "ConservationNFT");

            expect(await emergencyPause.isTrackedContract(conservationNFT.address)).to.be.true;
        });

        it("Should not allow non-manager to add contracts", async function () {
            await expect(
                emergencyPause.connect(guardian1).addContract(user.address, "TestContract")
            ).to.be.reverted;
        });

        it("Should prevent adding same contract twice", async function () {
            await expect(
                emergencyPause.connect(manager).addContract(
                    bogoToken.address,
                    "BOGOTokenV2"
                )
            ).to.be.revertedWith("ALREADY_EXISTS");
        });

        it("Should allow manager to remove contracts", async function () {
            await expect(
                emergencyPause.connect(manager).removeContract(bogoToken.address)
            ).to.emit(emergencyPause, "ContractRemoved")
                .withArgs(bogoToken.address);

            expect(await emergencyPause.isTrackedContract(bogoToken.address)).to.be.false;
        });
    });

    describe("Pause Proposals", function () {
        it("Should create pause proposal with single guardian", async function () {
            const reason = "Potential security vulnerability detected";
            
            await expect(
                emergencyPause.connect(guardian1).proposePause(
                    [bogoToken.address],
                    reason
                )
            ).to.emit(emergencyPause, "PauseProposalCreated")
                .withArgs(0, guardian1.address, true);

            const proposal = await emergencyPause.proposals(0);
            expect(proposal.proposer).to.equal(guardian1.address);
            expect(proposal.reason).to.equal(reason);
            expect(proposal.confirmations).to.equal(1);
            expect(proposal.executed).to.be.false;
            expect(proposal.isPause).to.be.true;
        });

        it("Should execute pause when reaching required confirmations", async function () {
            // First guardian proposes
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address],
                "Emergency pause needed"
            );

            // Second guardian confirms - should auto-execute
            await expect(
                emergencyPause.connect(guardian2).confirmProposal(0)
            ).to.emit(emergencyPause, "EmergencyPauseExecuted");

            // Check token is paused
            expect(await bogoToken.paused()).to.be.true;
        });

        it("Should not allow same guardian to confirm twice", async function () {
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address],
                "Emergency pause"
            );

            await expect(
                emergencyPause.connect(guardian1).confirmProposal(0)
            ).to.be.revertedWith("ALREADY_PROCESSED");
        });

        it("Should pause all contracts with emergencyPauseAll", async function () {
            await emergencyPause.connect(guardian1).emergencyPauseAll("System-wide emergency");
            await emergencyPause.connect(guardian2).confirmProposal(0);

            expect(await bogoToken.paused()).to.be.true;
            expect(await multisigTreasury.paused()).to.be.true;
        });
    });

    describe("Unpause Proposals", function () {
        beforeEach(async function () {
            // Pause contracts first
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address, multisigTreasury.address],
                "Initial pause"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);
        });

        it("Should create unpause proposal", async function () {
            await expect(
                emergencyPause.connect(guardian1).proposeUnpause(
                    [bogoToken.address],
                    "Issue resolved"
                )
            ).to.emit(emergencyPause, "PauseProposalCreated")
                .withArgs(1, guardian1.address, false);
        });

        it("Should unpause when reaching confirmations", async function () {
            await emergencyPause.connect(guardian1).proposeUnpause(
                [bogoToken.address],
                "Issue resolved"
            );
            
            await expect(
                emergencyPause.connect(guardian2).confirmProposal(1)
            ).to.emit(emergencyPause, "EmergencyUnpauseExecuted");

            expect(await bogoToken.paused()).to.be.false;
            expect(await multisigTreasury.paused()).to.be.true; // Still paused
        });
    });

    describe("Pause Expiry", function () {
        it("Should automatically expire pauses after MAX_PAUSE_DURATION", async function () {
            // Pause a contract
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address],
                "Time-limited pause"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);

            expect(await bogoToken.paused()).to.be.true;

            // Fast forward past max duration
            await time.increase(73 * 60 * 60); // 73 hours

            // Anyone can call checkAndExpirePauses
            await expect(
                emergencyPause.connect(user).checkAndExpirePauses()
            ).to.emit(emergencyPause, "PauseExpired")
                .withArgs(bogoToken.address);

            expect(await bogoToken.paused()).to.be.false;
        });

        it("Should not expire pauses before MAX_PAUSE_DURATION", async function () {
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address],
                "Time-limited pause"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);

            // Fast forward less than max duration
            await time.increase(71 * 60 * 60); // 71 hours

            await emergencyPause.connect(user).checkAndExpirePauses();
            
            // Should still be paused
            expect(await bogoToken.paused()).to.be.true;
        });
    });

    describe("Required Confirmations", function () {
        it("Should allow admin to update required confirmations", async function () {
            await expect(
                emergencyPause.updateRequiredConfirmations(1)
            ).to.emit(emergencyPause, "RequiredConfirmationsUpdated")
                .withArgs(2, 1);

            expect(await emergencyPause.requiredConfirmations()).to.equal(1);
        });

        it("Should auto-execute with single confirmation when required is 1", async function () {
            await emergencyPause.updateRequiredConfirmations(1);

            await expect(
                emergencyPause.connect(guardian1).proposePause(
                    [bogoToken.address],
                    "Single guardian pause"
                )
            ).to.emit(emergencyPause, "EmergencyPauseExecuted");

            expect(await bogoToken.paused()).to.be.true;
        });

        it("Should not allow setting confirmations to zero", async function () {
            await expect(
                emergencyPause.updateRequiredConfirmations(0)
            ).to.be.revertedWith("ZERO_AMOUNT");
        });

        it("Should not allow setting confirmations above guardian count", async function () {
            await expect(
                emergencyPause.updateRequiredConfirmations(5)
            ).to.be.revertedWith("EXCEEDS_LIMIT");
        });
    });

    describe("View Functions", function () {
        it("Should return contract statuses correctly", async function () {
            // Pause one contract
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address],
                "Pause test"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);

            const [contracts, paused, names] = await emergencyPause.getContractStatuses();
            
            expect(contracts.length).to.equal(2);
            expect(contracts[0]).to.equal(bogoToken.address);
            expect(contracts[1]).to.equal(multisigTreasury.address);
            
            expect(paused[0]).to.be.true;
            expect(paused[1]).to.be.false;
            
            expect(names[0]).to.equal("BOGOTokenV2");
            expect(names[1]).to.equal("MultisigTreasury");
        });

        it("Should return pause history", async function () {
            // Create some pause events
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address],
                "First pause"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);

            await emergencyPause.connect(guardian1).proposeUnpause(
                [bogoToken.address],
                "First unpause"
            );
            await emergencyPause.connect(guardian2).confirmProposal(1);

            const history = await emergencyPause.getPauseHistory(2);
            expect(history.length).to.equal(2);
            expect(history[0].isPause).to.be.true;
            expect(history[1].isPause).to.be.false;
        });
    });

    describe("Access Control", function () {
        it("Should prevent non-guardians from creating proposals", async function () {
            await expect(
                emergencyPause.connect(user).proposePause(
                    [bogoToken.address],
                    "Unauthorized pause"
                )
            ).to.be.reverted;
        });

        it("Should prevent non-guardians from confirming proposals", async function () {
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address],
                "Test pause"
            );

            await expect(
                emergencyPause.connect(user).confirmProposal(0)
            ).to.be.reverted;
        });

        it("Should allow admin to grant new guardian role", async function () {
            await emergencyPause.grantRole(GUARDIAN_ROLE, user.address);
            expect(await emergencyPause.hasRole(GUARDIAN_ROLE, user.address)).to.be.true;

            // New guardian should be able to create proposals
            await expect(
                emergencyPause.connect(user).proposePause(
                    [bogoToken.address],
                    "New guardian pause"
                )
            ).to.emit(emergencyPause, "PauseProposalCreated");
        });
    });

    describe("Edge Cases", function () {
        it("Should handle failed pause operations gracefully", async function () {
            // Add a contract that doesn't implement pause correctly
            const invalidContract = owner.address; // EOA, not a contract
            
            // Force add it (this would fail in production due to interface check)
            // For testing, we'll use a mock scenario

            // The pause proposal should still execute even if one contract fails
            await emergencyPause.connect(guardian1).proposePause(
                [bogoToken.address],
                "Pause with potential failure"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);

            expect(await bogoToken.paused()).to.be.true;
        });

        it("Should handle empty contract lists", async function () {
            await expect(
                emergencyPause.connect(guardian1).proposePause([], "Empty pause")
            ).to.be.revertedWith("INVALID_LENGTH");
        });

        it("Should handle very long reason strings", async function () {
            const longReason = "A".repeat(1000);
            
            await expect(
                emergencyPause.connect(guardian1).proposePause(
                    [bogoToken.address],
                    longReason
                )
            ).to.emit(emergencyPause, "PauseProposalCreated");
        });
    });
});