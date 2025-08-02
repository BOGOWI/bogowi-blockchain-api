const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("EmergencyPauseController Basic", function () {
    let emergencyPause;
    let owner, guardian1, guardian2, guardian3, manager;

    beforeEach(async function () {
        [owner, guardian1, guardian2, guardian3, manager] = await ethers.getSigners();

        // Deploy EmergencyPauseController
        const EmergencyPauseController = await ethers.getContractFactory("EmergencyPauseController");
        emergencyPause = await EmergencyPauseController.deploy(
            [guardian1.address, guardian2.address, guardian3.address],
            manager.address
        );
        await emergencyPause.deployed();
    });

    it("Should deploy successfully", async function () {
        expect(emergencyPause.address).to.not.equal(ethers.constants.AddressZero);
        expect(await emergencyPause.requiredConfirmations()).to.equal(2);
        expect(await emergencyPause.MAX_PAUSE_DURATION()).to.equal(72 * 60 * 60);
    });
});