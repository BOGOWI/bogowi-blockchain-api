const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("EmergencyPauseController", function () {
    let emergencyPause;
    let contract1, contract2;
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
        const MockPausable = await ethers.getContractFactory("MockPausable");
        contract1 = await MockPausable.deploy("Contract1");
        contract2 = await MockPausable.deploy("Contract2");

        // Grant PAUSER_ROLE to EmergencyPauseController
        await contract1.grantRole(PAUSER_ROLE, emergencyPause.address);
        await contract2.grantRole(PAUSER_ROLE, emergencyPause.address);

        // Add contracts to emergency pause controller
        await emergencyPause.connect(manager).addContract(contract1.address, "Contract1");
        await emergencyPause.connect(manager).addContract(contract2.address, "Contract2");
    });

    describe("Deployment", function () {
        it("Should deploy with correct initial state", async function () {
            expect(await emergencyPause.requiredConfirmations()).to.equal(2);
            expect(await emergencyPause.MAX_PAUSE_DURATION()).to.equal(72 * 60 * 60);
            expect(await emergencyPause.MIN_GUARDIANS()).to.equal(3);
        });

        it("Should assign roles correctly", async function () {
            expect(await emergencyPause.hasRole(GUARDIAN_ROLE, guardian1.address)).to.be.true;
            expect(await emergencyPause.hasRole(GUARDIAN_ROLE, guardian2.address)).to.be.true;
            expect(await emergencyPause.hasRole(GUARDIAN_ROLE, guardian3.address)).to.be.true;
            expect(await emergencyPause.hasRole(MANAGER_ROLE, manager.address)).to.be.true;
            expect(await emergencyPause.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
        });
    });

    describe("Contract Management", function () {
        it("Should track added contracts", async function () {
            expect(await emergencyPause.isTrackedContract(contract1.address)).to.be.true;
            expect(await emergencyPause.isTrackedContract(contract2.address)).to.be.true;
        });

        it("Should not allow non-manager to add contracts", async function () {
            const MockPausable = await ethers.getContractFactory("MockPausable");
            const contract3 = await MockPausable.deploy("Contract3");
            
            await expect(
                emergencyPause.connect(guardian1).addContract(contract3.address, "Contract3")
            ).to.be.reverted;
        });

        it("Should allow manager to remove contracts", async function () {
            await expect(
                emergencyPause.connect(manager).removeContract(contract1.address)
            ).to.emit(emergencyPause, "ContractRemoved")
                .withArgs(contract1.address);

            expect(await emergencyPause.isTrackedContract(contract1.address)).to.be.false;
        });
    });

    describe("Pause Proposals", function () {
        it("Should create pause proposal with single guardian", async function () {
            const reason = "Security vulnerability detected";
            
            await expect(
                emergencyPause.connect(guardian1).proposePause(
                    [contract1.address],
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
                [contract1.address],
                "Emergency pause needed"
            );

            // Second guardian confirms - should auto-execute
            await expect(
                emergencyPause.connect(guardian2).confirmProposal(0)
            ).to.emit(emergencyPause, "EmergencyPauseExecuted");

            // Check contract is paused
            expect(await contract1.paused()).to.be.true;
            expect(await contract2.paused()).to.be.false; // Not included in pause
        });

        it("Should pause all contracts with emergencyPauseAll", async function () {
            await emergencyPause.connect(guardian1).emergencyPauseAll("System-wide emergency");
            await emergencyPause.connect(guardian2).confirmProposal(0);

            expect(await contract1.paused()).to.be.true;
            expect(await contract2.paused()).to.be.true;
        });
    });

    describe("Unpause Proposals", function () {
        beforeEach(async function () {
            // Pause contracts first
            await emergencyPause.connect(guardian1).proposePause(
                [contract1.address, contract2.address],
                "Initial pause"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);
        });

        it("Should create unpause proposal", async function () {
            await expect(
                emergencyPause.connect(guardian1).proposeUnpause(
                    [contract1.address],
                    "Issue resolved"
                )
            ).to.emit(emergencyPause, "PauseProposalCreated")
                .withArgs(1, guardian1.address, false);
        });

        it("Should unpause when reaching confirmations", async function () {
            await emergencyPause.connect(guardian1).proposeUnpause(
                [contract1.address],
                "Issue resolved"
            );
            
            await expect(
                emergencyPause.connect(guardian2).confirmProposal(1)
            ).to.emit(emergencyPause, "EmergencyUnpauseExecuted");

            expect(await contract1.paused()).to.be.false;
            expect(await contract2.paused()).to.be.true; // Still paused
        });
    });

    describe("Pause Expiry", function () {
        it("Should automatically expire pauses after MAX_PAUSE_DURATION", async function () {
            // Pause a contract
            await emergencyPause.connect(guardian1).proposePause(
                [contract1.address],
                "Time-limited pause"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);

            expect(await contract1.paused()).to.be.true;

            // Fast forward past max duration
            await time.increase(73 * 60 * 60); // 73 hours

            // Anyone can call checkAndExpirePauses
            await expect(
                emergencyPause.connect(user).checkAndExpirePauses()
            ).to.emit(emergencyPause, "PauseExpired")
                .withArgs(contract1.address);

            expect(await contract1.paused()).to.be.false;
        });
    });

    describe("View Functions", function () {
        it("Should return contract statuses correctly", async function () {
            // Pause one contract
            await emergencyPause.connect(guardian1).proposePause(
                [contract1.address],
                "Pause test"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);

            const [contracts, paused, names] = await emergencyPause.getContractStatuses();
            
            expect(contracts.length).to.equal(2);
            expect(contracts[0]).to.equal(contract1.address);
            expect(contracts[1]).to.equal(contract2.address);
            
            expect(paused[0]).to.be.true;
            expect(paused[1]).to.be.false;
            
            expect(names[0]).to.equal("Contract1");
            expect(names[1]).to.equal("Contract2");
        });

        it("Should return pause history", async function () {
            // Create some pause events
            await emergencyPause.connect(guardian1).proposePause(
                [contract1.address],
                "First pause"
            );
            await emergencyPause.connect(guardian2).confirmProposal(0);

            await emergencyPause.connect(guardian1).proposeUnpause(
                [contract1.address],
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
                    [contract1.address],
                    "Unauthorized pause"
                )
            ).to.be.reverted;
        });

        it("Should allow admin to grant new guardian role", async function () {
            await emergencyPause.grantRole(GUARDIAN_ROLE, user.address);
            expect(await emergencyPause.hasRole(GUARDIAN_ROLE, user.address)).to.be.true;

            // New guardian should be able to create proposals
            await expect(
                emergencyPause.connect(user).proposePause(
                    [contract1.address],
                    "New guardian pause"
                )
            ).to.emit(emergencyPause, "PauseProposalCreated");
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
                    [contract1.address],
                    "Single guardian pause"
                )
            ).to.emit(emergencyPause, "EmergencyPauseExecuted");

            expect(await contract1.paused()).to.be.true;
        });
    });
});