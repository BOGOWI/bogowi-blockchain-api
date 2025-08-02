const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("MigrationHelper", function () {
    let migrationHelper;
    let owner, migrator, user1, user2, user3;
    let oldContract, newContract;
    let mockToken;
    
    const MIGRATION_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("MIGRATION_ROLE"));
    const DEFAULT_ADMIN_ROLE = "0x0000000000000000000000000000000000000000000000000000000000000000";
    const MAX_BATCH_SIZE = 100;

    beforeEach(async function () {
        [owner, migrator, user1, user2, user3] = await ethers.getSigners();

        // Deploy MigrationHelper
        const MigrationHelper = await ethers.getContractFactory("MigrationHelper");
        migrationHelper = await MigrationHelper.deploy();
        await migrationHelper.deployed();

        // Grant migrator role
        await migrationHelper.grantRole(MIGRATION_ROLE, migrator.address);

        // Deploy mock contracts
        const MockContract = await ethers.getContractFactory("MockERC20");
        oldContract = await MockContract.deploy("Old", "OLD", ethers.utils.parseEther("1000000"));
        newContract = await MockContract.deploy("New", "NEW", ethers.utils.parseEther("1000000"));
        mockToken = await MockContract.deploy("Token", "TKN", ethers.utils.parseEther("1000000"));
        
        await oldContract.deployed();
        await newContract.deployed();
        await mockToken.deployed();

        // Fund migration helper for testing
        await owner.sendTransaction({
            to: migrationHelper.address,
            value: ethers.utils.parseEther("1")
        });
        await mockToken.transfer(migrationHelper.address, ethers.utils.parseEther("1000"));
    });

    describe("Deployment", function () {
        it("Should set correct roles", async function () {
            expect(await migrationHelper.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, owner.address)).to.be.true;
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, migrator.address)).to.be.true;
        });
    });

    describe("Single User Migration", function () {
        it("Should mark user as migrated", async function () {
            await expect(migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address))
                .to.emit(migrationHelper, "UserMigrated")
                .withArgs(user1.address, oldContract.address, migrator.address);

            expect(await migrationHelper.isMigrated(oldContract.address, user1.address)).to.be.true;
            expect(await migrationHelper.migrationCount(oldContract.address)).to.equal(1);
        });

        it("Should revert if user already migrated", async function () {
            await migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address);
            
            await expect(
                migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address)
            ).to.be.revertedWithCustomError(migrationHelper, "AlreadyMigrated");
        });

        it("Should track migrations per contract", async function () {
            await migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address);
            await migrationHelper.connect(migrator).markMigrated(oldContract.address, user2.address);
            await migrationHelper.connect(migrator).markMigrated(newContract.address, user1.address);

            expect(await migrationHelper.migrationCount(oldContract.address)).to.equal(2);
            expect(await migrationHelper.migrationCount(newContract.address)).to.equal(1);
        });

        it("Should only allow authorized migrators", async function () {
            await expect(
                migrationHelper.connect(user1).markMigrated(oldContract.address, user2.address)
            ).to.be.reverted;
        });
    });

    describe("Batch Migration", function () {
        it("Should batch mark users as migrated", async function () {
            const users = [user1.address, user2.address, user3.address];
            
            await expect(migrationHelper.connect(migrator).batchMarkMigrated(oldContract.address, users))
                .to.emit(migrationHelper, "BatchMigrationCompleted")
                .withArgs(oldContract.address, 3);

            expect(await migrationHelper.isMigrated(oldContract.address, user1.address)).to.be.true;
            expect(await migrationHelper.isMigrated(oldContract.address, user2.address)).to.be.true;
            expect(await migrationHelper.isMigrated(oldContract.address, user3.address)).to.be.true;
            expect(await migrationHelper.migrationCount(oldContract.address)).to.equal(3);
        });

        it("Should skip already migrated users in batch", async function () {
            await migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address);
            
            const users = [user1.address, user2.address, user3.address];
            await expect(migrationHelper.connect(migrator).batchMarkMigrated(oldContract.address, users))
                .to.emit(migrationHelper, "BatchMigrationCompleted")
                .withArgs(oldContract.address, 2); // Only 2 new migrations

            expect(await migrationHelper.migrationCount(oldContract.address)).to.equal(3);
        });

        it("Should enforce batch size limit", async function () {
            const users = [];
            for (let i = 0; i < MAX_BATCH_SIZE + 1; i++) {
                const wallet = ethers.Wallet.createRandom();
                users.push(wallet.address);
            }

            await expect(
                migrationHelper.connect(migrator).batchMarkMigrated(oldContract.address, users)
            ).to.be.revertedWithCustomError(migrationHelper, "BatchSizeTooLarge");
        });

        it("Should handle empty batch", async function () {
            await expect(migrationHelper.connect(migrator).batchMarkMigrated(oldContract.address, []))
                .to.emit(migrationHelper, "BatchMigrationCompleted")
                .withArgs(oldContract.address, 0);
        });
    });

    describe("Token Recovery", function () {
        it("Should recover ETH", async function () {
            const balanceBefore = await ethers.provider.getBalance(user1.address);
            
            await expect(migrationHelper.recoverTokens(ethers.constants.AddressZero, user1.address, ethers.utils.parseEther("0.5")))
                .to.emit(migrationHelper, "TokensRecovered")
                .withArgs(ethers.constants.AddressZero, user1.address, ethers.utils.parseEther("0.5"));

            const balanceAfter = await ethers.provider.getBalance(user1.address);
            expect(balanceAfter.sub(balanceBefore)).to.equal(ethers.utils.parseEther("0.5"));
        });

        it("Should recover ERC20 tokens", async function () {
            const balanceBefore = await mockToken.balanceOf(user1.address);
            
            await expect(migrationHelper.recoverTokens(mockToken.address, user1.address, ethers.utils.parseEther("100")))
                .to.emit(migrationHelper, "TokensRecovered")
                .withArgs(mockToken.address, user1.address, ethers.utils.parseEther("100"));

            const balanceAfter = await mockToken.balanceOf(user1.address);
            expect(balanceAfter.sub(balanceBefore)).to.equal(ethers.utils.parseEther("100"));
        });

        it("Should only allow admin to recover tokens", async function () {
            await expect(
                migrationHelper.connect(user1).recoverTokens(mockToken.address, user1.address, ethers.utils.parseEther("100"))
            ).to.be.reverted;
        });

        it("Should validate recovery parameters", async function () {
            await expect(
                migrationHelper.recoverTokens(mockToken.address, ethers.constants.AddressZero, ethers.utils.parseEther("100"))
            ).to.be.revertedWith("Invalid recipient");

            await expect(
                migrationHelper.recoverTokens(mockToken.address, user1.address, 0)
            ).to.be.revertedWith("Invalid amount");
        });
    });

    describe("Pausable", function () {
        it("Should pause and unpause", async function () {
            await migrationHelper.pause();
            expect(await migrationHelper.paused()).to.be.true;

            await expect(
                migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address)
            ).to.be.revertedWith("Pausable: paused");

            await migrationHelper.unpause();
            expect(await migrationHelper.paused()).to.be.false;

            await expect(
                migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address)
            ).to.not.be.reverted;
        });

        it("Should only allow admin to pause/unpause", async function () {
            await expect(
                migrationHelper.connect(user1).pause()
            ).to.be.reverted;

            await migrationHelper.pause();

            await expect(
                migrationHelper.connect(user1).unpause()
            ).to.be.reverted;
        });
    });

    describe("Access Control", function () {
        it("Should manage migration role", async function () {
            const newMigrator = user3;
            
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, newMigrator.address)).to.be.false;
            
            await migrationHelper.grantRole(MIGRATION_ROLE, newMigrator.address);
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, newMigrator.address)).to.be.true;
            
            await expect(
                migrationHelper.connect(newMigrator).markMigrated(oldContract.address, user1.address)
            ).to.not.be.reverted;
            
            await migrationHelper.revokeRole(MIGRATION_ROLE, newMigrator.address);
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, newMigrator.address)).to.be.false;
        });
    });

    describe("Gas Optimization", function () {
        it("Should efficiently handle large batches", async function () {
            const users = [];
            for (let i = 0; i < 50; i++) {
                const wallet = ethers.Wallet.createRandom();
                users.push(wallet.address);
            }

            const tx = await migrationHelper.connect(migrator).batchMarkMigrated(oldContract.address, users);
            const receipt = await tx.wait();
            
            const gasPerUser = receipt.gasUsed.div(50);
            console.log(`Gas per user in batch: ${gasPerUser}`);
            
            // Should be efficient - less than 50k gas per user
            expect(gasPerUser.toNumber()).to.be.lessThan(50000);
        });
    });

    describe("Integration Scenarios", function () {
        it("Should support complete migration flow", async function () {
            // 1. Mark individual high-value users
            await migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address);
            
            // 2. Batch migrate regular users
            const regularUsers = [user2.address, user3.address];
            await migrationHelper.connect(migrator).batchMarkMigrated(oldContract.address, regularUsers);
            
            // 3. Verify all migrations
            expect(await migrationHelper.isMigrated(oldContract.address, user1.address)).to.be.true;
            expect(await migrationHelper.isMigrated(oldContract.address, user2.address)).to.be.true;
            expect(await migrationHelper.isMigrated(oldContract.address, user3.address)).to.be.true;
            
            // 4. Check total count
            expect(await migrationHelper.migrationCount(oldContract.address)).to.equal(3);
            
            // 5. Attempt duplicate migration (should be safe)
            await expect(
                migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address)
            ).to.be.revertedWithCustomError(migrationHelper, "AlreadyMigrated");
        });

        it("Should handle multi-contract migration", async function () {
            // Migrate user1 from multiple old contracts
            await migrationHelper.connect(migrator).markMigrated(oldContract.address, user1.address);
            await migrationHelper.connect(migrator).markMigrated(newContract.address, user1.address);
            
            expect(await migrationHelper.isMigrated(oldContract.address, user1.address)).to.be.true;
            expect(await migrationHelper.isMigrated(newContract.address, user1.address)).to.be.true;
            
            // User2 only migrated from one contract
            await migrationHelper.connect(migrator).markMigrated(oldContract.address, user2.address);
            
            expect(await migrationHelper.isMigrated(oldContract.address, user2.address)).to.be.true;
            expect(await migrationHelper.isMigrated(newContract.address, user2.address)).to.be.false;
        });
    });
});