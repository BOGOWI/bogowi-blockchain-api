const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("ContractRegistry", function () {
    let registry;
    let owner, admin, user1, user2;
    let mockContract1, mockContract2, mockContract3;
    
    const REGISTRY_ADMIN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("REGISTRY_ADMIN_ROLE"));
    const DEFAULT_ADMIN_ROLE = "0x0000000000000000000000000000000000000000000000000000000000000000";

    beforeEach(async function () {
        [owner, admin, user1, user2] = await ethers.getSigners();

        // Deploy ContractRegistry
        const ContractRegistry = await ethers.getContractFactory("ContractRegistry");
        registry = await ContractRegistry.deploy();
        await registry.deployed();

        // Deploy mock contracts
        const MockContract = await ethers.getContractFactory("MockERC20");
        mockContract1 = await MockContract.deploy("Mock1", "M1", ethers.utils.parseEther("1000000"));
        mockContract2 = await MockContract.deploy("Mock2", "M2", ethers.utils.parseEther("1000000"));
        mockContract3 = await MockContract.deploy("Mock3", "M3", ethers.utils.parseEther("1000000"));
        await mockContract1.deployed();
        await mockContract2.deployed();
        await mockContract3.deployed();

        // Grant admin role
        await registry.grantRole(REGISTRY_ADMIN_ROLE, admin.address);
    });

    describe("Deployment", function () {
        it("Should set the correct deployer roles", async function () {
            expect(await registry.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
            expect(await registry.hasRole(REGISTRY_ADMIN_ROLE, owner.address)).to.be.true;
        });
    });

    describe("Contract Registration", function () {
        it("Should register a new contract", async function () {
            await expect(registry.registerContract("BOGOToken", mockContract1.address))
                .to.emit(registry, "ContractRegistered")
                .withArgs("BOGOToken", mockContract1.address, 1);

            expect(await registry.getContract("BOGOToken")).to.equal(mockContract1.address);
            expect(await registry.getContractVersion("BOGOToken")).to.equal(1);
            expect(await registry.isRegistered("BOGOToken")).to.be.true;
        });

        it("Should revert when registering with zero address", async function () {
            await expect(
                registry.registerContract("BOGOToken", ethers.constants.AddressZero)
            ).to.be.revertedWithCustomError(registry, "InvalidAddress");
        });

        it("Should revert when registering already registered contract", async function () {
            await registry.registerContract("BOGOToken", mockContract1.address);
            
            await expect(
                registry.registerContract("BOGOToken", mockContract2.address)
            ).to.be.revertedWithCustomError(registry, "ContractAlreadyRegistered");
        });

        it("Should not allow non-admin to register contract", async function () {
            await expect(
                registry.connect(user1).registerContract("BOGOToken", mockContract1.address)
            ).to.be.reverted;
        });

        it("Should register multiple different contracts", async function () {
            await registry.registerContract("BOGOToken", mockContract1.address);
            await registry.registerContract("RewardDistributor", mockContract2.address);
            await registry.registerContract("Treasury", mockContract3.address);

            expect(await registry.getContract("BOGOToken")).to.equal(mockContract1.address);
            expect(await registry.getContract("RewardDistributor")).to.equal(mockContract2.address);
            expect(await registry.getContract("Treasury")).to.equal(mockContract3.address);
        });
    });

    describe("Contract Updates", function () {
        beforeEach(async function () {
            await registry.registerContract("RewardDistributor", mockContract1.address);
        });

        it("Should update an existing contract", async function () {
            await expect(registry.updateContract("RewardDistributor", mockContract2.address))
                .to.emit(registry, "ContractUpdated")
                .withArgs("RewardDistributor", mockContract1.address, mockContract2.address, 2);

            expect(await registry.getContract("RewardDistributor")).to.equal(mockContract2.address);
            expect(await registry.getContractVersion("RewardDistributor")).to.equal(2);
        });

        it("Should maintain contract history", async function () {
            await registry.updateContract("RewardDistributor", mockContract2.address);
            await registry.updateContract("RewardDistributor", mockContract3.address);

            const history = await registry.getContractHistory("RewardDistributor");
            expect(history).to.have.lengthOf(3);
            expect(history[0]).to.equal(mockContract1.address);
            expect(history[1]).to.equal(mockContract2.address);
            expect(history[2]).to.equal(mockContract3.address);
        });

        it("Should revert when updating with zero address", async function () {
            await expect(
                registry.updateContract("RewardDistributor", ethers.constants.AddressZero)
            ).to.be.revertedWithCustomError(registry, "InvalidAddress");
        });

        it("Should revert when updating non-existent contract", async function () {
            await expect(
                registry.updateContract("NonExistent", mockContract2.address)
            ).to.be.revertedWithCustomError(registry, "ContractNotFound");
        });

        it("Should not allow non-admin to update contract", async function () {
            await expect(
                registry.connect(user1).updateContract("RewardDistributor", mockContract2.address)
            ).to.be.reverted;
        });
    });

    describe("Contract Deprecation", function () {
        beforeEach(async function () {
            await registry.registerContract("OldContract", mockContract1.address);
        });

        it("Should deprecate a contract", async function () {
            await expect(registry.deprecateContract("OldContract"))
                .to.emit(registry, "ContractDeprecated")
                .withArgs("OldContract", mockContract1.address);

            await expect(
                registry.getContract("OldContract")
            ).to.be.revertedWithCustomError(registry, "ContractNotFound");

            expect(await registry.isRegistered("OldContract")).to.be.false;
        });

        it("Should maintain history after deprecation", async function () {
            await registry.updateContract("OldContract", mockContract2.address);
            await registry.deprecateContract("OldContract");

            // History should still be accessible
            const history = await registry.getContractHistory("OldContract");
            expect(history).to.have.lengthOf(2);
        });

        it("Should revert when deprecating non-existent contract", async function () {
            await expect(
                registry.deprecateContract("NonExistent")
            ).to.be.revertedWithCustomError(registry, "ContractNotFound");
        });
    });

    describe("Queries", function () {
        beforeEach(async function () {
            await registry.registerContract("BOGOToken", mockContract1.address);
            await registry.registerContract("RewardDistributor", mockContract2.address);
        });

        it("Should return correct contract address", async function () {
            expect(await registry.getContract("BOGOToken")).to.equal(mockContract1.address);
            expect(await registry.getContract("RewardDistributor")).to.equal(mockContract2.address);
        });

        it("Should revert when querying non-existent contract", async function () {
            await expect(
                registry.getContract("NonExistent")
            ).to.be.revertedWithCustomError(registry, "ContractNotFound");
        });

        it("Should return correct registration status", async function () {
            expect(await registry.isRegistered("BOGOToken")).to.be.true;
            expect(await registry.isRegistered("RewardDistributor")).to.be.true;
            expect(await registry.isRegistered("NonExistent")).to.be.false;
        });

        it("Should return correct version", async function () {
            expect(await registry.getContractVersion("BOGOToken")).to.equal(1);
            
            await registry.updateContract("BOGOToken", mockContract3.address);
            expect(await registry.getContractVersion("BOGOToken")).to.equal(2);
        });
    });

    describe("Pausable", function () {
        beforeEach(async function () {
            await registry.registerContract("BOGOToken", mockContract1.address);
        });

        it("Should pause and unpause", async function () {
            await registry.pause();
            expect(await registry.paused()).to.be.true;

            // Operations should fail when paused
            await expect(
                registry.registerContract("NewContract", mockContract2.address)
            ).to.be.revertedWith("Pausable: paused");

            await expect(
                registry.updateContract("BOGOToken", mockContract2.address)
            ).to.be.revertedWith("Pausable: paused");

            await registry.unpause();
            expect(await registry.paused()).to.be.false;

            // Operations should work after unpause
            await expect(
                registry.registerContract("NewContract", mockContract2.address)
            ).to.not.be.reverted;
        });

        it("Should only allow admin to pause/unpause", async function () {
            await expect(
                registry.connect(user1).pause()
            ).to.be.reverted;

            await registry.pause();

            await expect(
                registry.connect(user1).unpause()
            ).to.be.reverted;
        });
    });

    describe("Access Control", function () {
        it("Should grant and revoke admin roles", async function () {
            expect(await registry.hasRole(REGISTRY_ADMIN_ROLE, user1.address)).to.be.false;
            
            await registry.grantRole(REGISTRY_ADMIN_ROLE, user1.address);
            expect(await registry.hasRole(REGISTRY_ADMIN_ROLE, user1.address)).to.be.true;
            
            // User1 should now be able to register contracts
            await expect(
                registry.connect(user1).registerContract("UserContract", mockContract1.address)
            ).to.not.be.reverted;
            
            await registry.revokeRole(REGISTRY_ADMIN_ROLE, user1.address);
            expect(await registry.hasRole(REGISTRY_ADMIN_ROLE, user1.address)).to.be.false;
        });
    });

    describe("Gas Optimization", function () {
        it("Should handle multiple updates efficiently", async function () {
            await registry.registerContract("GasTest", mockContract1.address);
            
            const updates = 10;
            let totalGas = 0;
            
            for (let i = 0; i < updates; i++) {
                const tx = await registry.updateContract("GasTest", mockContract2.address);
                const receipt = await tx.wait();
                totalGas += receipt.gasUsed.toNumber();
                
                // Alternate between contracts
                if (i % 2 === 0) {
                    await registry.updateContract("GasTest", mockContract1.address);
                }
            }
            
            const avgGas = totalGas / updates;
            console.log(`Average gas per update: ${avgGas}`);
            expect(avgGas).to.be.lessThan(100000); // Should be efficient
        });
    });

    describe("Integration Scenarios", function () {
        it("Should support complete contract migration flow", async function () {
            // 1. Register initial contract
            await registry.registerContract("RewardDistributor", mockContract1.address);
            
            // 2. Simulate time passing and need for update
            await ethers.provider.send("evm_increaseTime", [86400]); // 1 day
            await ethers.provider.send("evm_mine");
            
            // 3. Deploy and register new version
            await registry.updateContract("RewardDistributor", mockContract2.address);
            
            // 4. Verify migration
            expect(await registry.getContract("RewardDistributor")).to.equal(mockContract2.address);
            expect(await registry.getContractVersion("RewardDistributor")).to.equal(2);
            
            // 5. Check history
            const history = await registry.getContractHistory("RewardDistributor");
            expect(history[0]).to.equal(mockContract1.address);
            expect(history[1]).to.equal(mockContract2.address);
        });

        it("Should handle emergency deprecation", async function () {
            // Register multiple contracts
            await registry.registerContract("Vulnerable", mockContract1.address);
            await registry.registerContract("Safe1", mockContract2.address);
            await registry.registerContract("Safe2", mockContract3.address);
            
            // Emergency: deprecate vulnerable contract
            await registry.pause();
            await registry.unpause(); // Can still deprecate after unpause
            await registry.deprecateContract("Vulnerable");
            
            // Verify other contracts still accessible
            expect(await registry.getContract("Safe1")).to.equal(mockContract2.address);
            expect(await registry.getContract("Safe2")).to.equal(mockContract3.address);
            
            // Vulnerable contract should not be accessible
            await expect(
                registry.getContract("Vulnerable")
            ).to.be.revertedWithCustomError(registry, "ContractNotFound");
        });
    });
});