const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("RoleManager", function () {
    let roleManager;
    let owner, admin, dao, business, minter, pauser, treasury, backend, user1, user2;
    let mockContract;

    // Role constants
    let DEFAULT_ADMIN_ROLE, DAO_ROLE, BUSINESS_ROLE, MINTER_ROLE, PAUSER_ROLE, TREASURY_ROLE, DISTRIBUTOR_BACKEND_ROLE;

    beforeEach(async function () {
        [owner, admin, dao, business, minter, pauser, treasury, backend, user1, user2] = await ethers.getSigners();

        // Deploy RoleManager
        const RoleManager = await ethers.getContractFactory("RoleManager");
        roleManager = await RoleManager.deploy();
        await roleManager.deployed();

        // Get role constants
        DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
        DAO_ROLE = await roleManager.DAO_ROLE();
        BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
        MINTER_ROLE = await roleManager.MINTER_ROLE();
        PAUSER_ROLE = await roleManager.PAUSER_ROLE();
        TREASURY_ROLE = await roleManager.TREASURY_ROLE();
        DISTRIBUTOR_BACKEND_ROLE = await roleManager.DISTRIBUTOR_BACKEND_ROLE();

        // Deploy a mock contract for testing
        const MockContract = await ethers.getContractFactory("RoleManager");
        mockContract = await MockContract.deploy();
        await mockContract.deployed();
    });

    describe("Deployment", function () {
        it("Should set the correct deployer as DEFAULT_ADMIN", async function () {
            expect(await roleManager.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
        });

        it("Should have correct role admin relationships", async function () {
            expect(await roleManager.getRoleAdmin(DAO_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
            expect(await roleManager.getRoleAdmin(BUSINESS_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
            expect(await roleManager.getRoleAdmin(MINTER_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
            expect(await roleManager.getRoleAdmin(PAUSER_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
            expect(await roleManager.getRoleAdmin(TREASURY_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
            expect(await roleManager.getRoleAdmin(DISTRIBUTOR_BACKEND_ROLE)).to.equal(DEFAULT_ADMIN_ROLE);
        });
    });

    describe("Contract Registration", function () {
        it("Should allow admin to register a contract", async function () {
            await expect(roleManager.registerContract(mockContract.address, "MockContract"))
                .to.emit(roleManager, "ContractRegistered")
                .withArgs(mockContract.address, "MockContract");

            expect(await roleManager.registeredContracts(mockContract.address)).to.be.true;
            expect(await roleManager.contractNames(mockContract.address)).to.equal("MockContract");
        });

        it("Should not allow non-admin to register a contract", async function () {
            await expect(
                roleManager.connect(user1).registerContract(mockContract.address, "MockContract")
            ).to.be.reverted;
        });

        it("Should not allow registering zero address", async function () {
            await expect(
                roleManager.registerContract(ethers.constants.AddressZero, "MockContract")
            ).to.be.revertedWith("Invalid contract address");
        });

        it("Should not allow registering same contract twice", async function () {
            await roleManager.registerContract(mockContract.address, "MockContract");
            await expect(
                roleManager.registerContract(mockContract.address, "MockContract2")
            ).to.be.revertedWith("Contract already registered");
        });

        it("Should allow admin to deregister a contract", async function () {
            await roleManager.registerContract(mockContract.address, "MockContract");
            
            await expect(roleManager.deregisterContract(mockContract.address))
                .to.emit(roleManager, "ContractDeregistered")
                .withArgs(mockContract.address);

            expect(await roleManager.registeredContracts(mockContract.address)).to.be.false;
            expect(await roleManager.contractNames(mockContract.address)).to.equal("");
        });
    });

    describe("Role Management", function () {
        it("Should allow admin to grant roles", async function () {
            await expect(roleManager.grantRole(DAO_ROLE, dao.address))
                .to.emit(roleManager, "RoleGrantedGlobally")
                .withArgs(DAO_ROLE, dao.address, owner.address);

            expect(await roleManager.hasRole(DAO_ROLE, dao.address)).to.be.true;
        });

        it("Should allow admin to revoke roles", async function () {
            await roleManager.grantRole(DAO_ROLE, dao.address);
            
            await expect(roleManager.revokeRole(DAO_ROLE, dao.address))
                .to.emit(roleManager, "RoleRevokedGlobally")
                .withArgs(DAO_ROLE, dao.address, owner.address);

            expect(await roleManager.hasRole(DAO_ROLE, dao.address)).to.be.false;
        });

        it("Should not allow non-admin to grant roles", async function () {
            await expect(
                roleManager.connect(user1).grantRole(DAO_ROLE, dao.address)
            ).to.be.reverted;
        });

        it("Should allow batch role granting", async function () {
            const accounts = [dao.address, user1.address, user2.address];
            
            await roleManager.batchGrantRole(DAO_ROLE, accounts);

            for (const account of accounts) {
                expect(await roleManager.hasRole(DAO_ROLE, account)).to.be.true;
            }
        });

        it("Should allow batch role revoking", async function () {
            const accounts = [dao.address, user1.address, user2.address];
            await roleManager.batchGrantRole(DAO_ROLE, accounts);
            
            await roleManager.batchRevokeRole(DAO_ROLE, accounts);

            for (const account of accounts) {
                expect(await roleManager.hasRole(DAO_ROLE, account)).to.be.false;
            }
        });
    });

    describe("Role Checking", function () {
        beforeEach(async function () {
            // Register the mock contract
            await roleManager.registerContract(owner.address, "TestContract");
        });

        it("Should allow registered contracts to check roles", async function () {
            await roleManager.grantRole(DAO_ROLE, dao.address);
            
            // Simulate call from registered contract (owner in this case)
            expect(await roleManager.checkRole(DAO_ROLE, dao.address)).to.be.true;
            expect(await roleManager.checkRole(DAO_ROLE, user1.address)).to.be.false;
        });

        it("Should not allow unregistered contracts to check roles", async function () {
            // Deregister the contract
            await roleManager.deregisterContract(owner.address);
            
            await expect(
                roleManager.checkRole(DAO_ROLE, dao.address)
            ).to.be.revertedWith("Not a registered contract");
        });
    });

    describe("Admin Transfer", function () {
        it("Should allow admin to transfer admin role", async function () {
            await roleManager.transferAdmin(admin.address);

            expect(await roleManager.hasRole(DEFAULT_ADMIN_ROLE, admin.address)).to.be.true;
            expect(await roleManager.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.false;
        });

        it("Should not allow transfer to zero address", async function () {
            await expect(
                roleManager.transferAdmin(ethers.constants.AddressZero)
            ).to.be.revertedWith("Invalid new admin");
        });

        it("Should not allow non-admin to transfer admin role", async function () {
            await expect(
                roleManager.connect(user1).transferAdmin(admin.address)
            ).to.be.reverted;
        });
    });

    describe("Pausable", function () {
        it("Should allow admin to pause and unpause", async function () {
            await roleManager.pause();
            expect(await roleManager.paused()).to.be.true;

            await roleManager.unpause();
            expect(await roleManager.paused()).to.be.false;
        });

        it("Should not allow non-admin to pause", async function () {
            await expect(
                roleManager.connect(user1).pause()
            ).to.be.reverted;
        });
    });

    describe("Multiple Role Scenarios", function () {
        it("Should handle multiple roles for single account", async function () {
            await roleManager.grantRole(DAO_ROLE, user1.address);
            await roleManager.grantRole(MINTER_ROLE, user1.address);
            await roleManager.grantRole(PAUSER_ROLE, user1.address);

            expect(await roleManager.hasRole(DAO_ROLE, user1.address)).to.be.true;
            expect(await roleManager.hasRole(MINTER_ROLE, user1.address)).to.be.true;
            expect(await roleManager.hasRole(PAUSER_ROLE, user1.address)).to.be.true;
            expect(await roleManager.hasRole(BUSINESS_ROLE, user1.address)).to.be.false;
        });

        it("Should handle role renouncing", async function () {
            await roleManager.grantRole(DAO_ROLE, owner.address);
            
            await roleManager.renounceRole(DAO_ROLE, owner.address);
            
            expect(await roleManager.hasRole(DAO_ROLE, owner.address)).to.be.false;
        });

        it("Should not allow renouncing someone else's role", async function () {
            await roleManager.grantRole(DAO_ROLE, user1.address);
            
            await expect(
                roleManager.connect(user2).renounceRole(DAO_ROLE, user1.address)
            ).to.be.reverted;
        });
    });

    describe("Gas Optimization Tests", function () {
        it("Should efficiently handle batch operations", async function () {
            const accounts = [];
            for (let i = 0; i < 10; i++) {
                const wallet = ethers.Wallet.createRandom();
                accounts.push(wallet.address);
            }

            const tx = await roleManager.batchGrantRole(DAO_ROLE, accounts);
            const receipt = await tx.wait();
            
            // Check gas usage is reasonable
            expect(receipt.gasUsed.toNumber()).to.be.lessThan(500000);
        });
    });
});