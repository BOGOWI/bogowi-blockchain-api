const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("MigrationHelper", function () {
    let migrationHelper;
    let owner, migrator, user1, user2, user3;
    let oldContract, newContract;
    let mockToken;
    
    const MIGRATION_ROLE = ethers.keccak256(ethers.toUtf8Bytes("MIGRATION_ROLE"));
    const DEFAULT_ADMIN_ROLE = "0x0000000000000000000000000000000000000000000000000000000000000000";
    const MAX_BATCH_SIZE = 100;

    beforeEach(async function () {
        [owner, migrator, user1, user2, user3] = await ethers.getSigners();

        // Deploy MigrationHelper
        const MigrationHelper = await ethers.getContractFactory("MigrationHelper");
        migrationHelper = await MigrationHelper.deploy();
        await migrationHelper.waitForDeployment();

        // Grant migrator role
        await migrationHelper.grantRole(MIGRATION_ROLE, await migrator.getAddress());

        // Deploy mock contracts
        const MockContract = await ethers.getContractFactory("contracts/test/MockERC20.sol:MockERC20");
        oldContract = await MockContract.deploy("Old", "OLD", ethers.parseEther("1000000"));
        newContract = await MockContract.deploy("New", "NEW", ethers.parseEther("1000000"));
        mockToken = await MockContract.deploy("Token", "TKN", ethers.parseEther("1000000"));
        
        await oldContract.waitForDeployment();
        await newContract.waitForDeployment();
        await mockToken.waitForDeployment();

        // Fund migration helper for testing
        await owner.sendTransaction({
            to: await migrationHelper.getAddress(),
            value: ethers.parseEther("1")
        });
        await mockToken.transfer(await migrationHelper.getAddress(), ethers.parseEther("1000"));
    });

    describe("Deployment", function () {
        it("Should set correct roles", async function () {
            expect(await migrationHelper.hasRole(DEFAULT_ADMIN_ROLE, await owner.getAddress())).to.be.true;
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, await owner.getAddress())).to.be.true;
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, await migrator.getAddress())).to.be.true;
        });
    });

    describe("Single User Migration", function () {
        it("Should mark user as migrated", async function () {
            await expect(migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress()))
                .to.emit(migrationHelper, "UserMigrated")
                .withArgs(await user1.getAddress(), await oldContract.getAddress(), await migrator.getAddress());

            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user1.getAddress())).to.be.true;
            expect(await migrationHelper.migrationCount(await oldContract.getAddress())).to.equal(1);
        });

        it("Should revert if user already migrated", async function () {
            await migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress());
            
            await expect(
                migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress())
            ).to.be.revertedWithCustomError(migrationHelper, "AlreadyMigrated");
        });

        it("Should track migrations per contract", async function () {
            await migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress());
            await migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user2.getAddress());
            await migrationHelper.connect(migrator).markMigrated(await newContract.getAddress(), await user1.getAddress());

            expect(await migrationHelper.migrationCount(await oldContract.getAddress())).to.equal(2);
            expect(await migrationHelper.migrationCount(await newContract.getAddress())).to.equal(1);
        });

        it("Should only allow authorized migrators", async function () {
            await expect(
                migrationHelper.connect(user1).markMigrated(await oldContract.getAddress(), await user2.getAddress())
            ).to.be.reverted;
        });
    });

    describe("Batch Migration", function () {
        it("Should batch mark users as migrated", async function () {
            const users = [await user1.getAddress(), await user2.getAddress(), await user3.getAddress()];
            
            await expect(migrationHelper.connect(migrator).batchMarkMigrated(await oldContract.getAddress(), users))
                .to.emit(migrationHelper, "BatchMigrationCompleted")
                .withArgs(await oldContract.getAddress(), 3);

            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user1.getAddress())).to.be.true;
            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user2.getAddress())).to.be.true;
            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user3.getAddress())).to.be.true;
            expect(await migrationHelper.migrationCount(await oldContract.getAddress())).to.equal(3);
        });

        it("Should skip already migrated users in batch", async function () {
            await migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress());
            
            const users = [await user1.getAddress(), await user2.getAddress(), await user3.getAddress()];
            await expect(migrationHelper.connect(migrator).batchMarkMigrated(await oldContract.getAddress(), users))
                .to.emit(migrationHelper, "BatchMigrationCompleted")
                .withArgs(await oldContract.getAddress(), 2); // Only 2 new migrations

            expect(await migrationHelper.migrationCount(await oldContract.getAddress())).to.equal(3);
        });

        it("Should enforce batch size limit", async function () {
            const users = [];
            for (let i = 0; i < MAX_BATCH_SIZE + 1; i++) {
                const wallet = ethers.Wallet.createRandom();
                users.push(wallet.address);
            }

            await expect(
                migrationHelper.connect(migrator).batchMarkMigrated(await oldContract.getAddress(), users)
            ).to.be.revertedWithCustomError(migrationHelper, "BatchSizeTooLarge");
        });

        it("Should handle empty batch", async function () {
            await expect(migrationHelper.connect(migrator).batchMarkMigrated(await oldContract.getAddress(), []))
                .to.emit(migrationHelper, "BatchMigrationCompleted")
                .withArgs(await oldContract.getAddress(), 0);
        });
    });

    describe("Token Recovery", function () {
        it("Should recover ETH", async function () {
            const balanceBefore = await ethers.provider.getBalance(await user1.getAddress());
            
            await expect(migrationHelper.recoverTokens(ethers.ZeroAddress, await user1.getAddress(), ethers.parseEther("0.5")))
                .to.emit(migrationHelper, "TokensRecovered")
                .withArgs(ethers.ZeroAddress, await user1.getAddress(), ethers.parseEther("0.5"));

            const balanceAfter = await ethers.provider.getBalance(await user1.getAddress());
            expect(balanceAfter - balanceBefore).to.equal(ethers.parseEther("0.5"));
        });

        it("Should recover ERC20 tokens", async function () {
            const balanceBefore = await mockToken.balanceOf(await user1.getAddress());
            
            await expect(migrationHelper.recoverTokens(await mockToken.getAddress(), await user1.getAddress(), ethers.parseEther("100")))
                .to.emit(migrationHelper, "TokensRecovered")
                .withArgs(await mockToken.getAddress(), await user1.getAddress(), ethers.parseEther("100"));

            const balanceAfter = await mockToken.balanceOf(await user1.getAddress());
            expect(balanceAfter - balanceBefore).to.equal(ethers.parseEther("100"));
        });

        it("Should only allow admin to recover tokens", async function () {
            await expect(
                migrationHelper.connect(user1).recoverTokens(await mockToken.getAddress(), await user1.getAddress(), ethers.parseEther("100"))
            ).to.be.reverted;
        });

        it("Should validate recovery parameters", async function () {
            await expect(
                migrationHelper.recoverTokens(await mockToken.getAddress(), ethers.ZeroAddress, ethers.parseEther("100"))
            ).to.be.revertedWith("Invalid recipient");

            await expect(
                migrationHelper.recoverTokens(await mockToken.getAddress(), await user1.getAddress(), 0)
            ).to.be.revertedWith("Invalid amount");
        });
    });

    describe("Pausable", function () {
        it("Should pause and unpause", async function () {
            await migrationHelper.pause();
            expect(await migrationHelper.paused()).to.be.true;

            await expect(
                migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress())
            ).to.be.revertedWithCustomError(migrationHelper, "EnforcedPause");

            await migrationHelper.unpause();
            expect(await migrationHelper.paused()).to.be.false;

            await expect(
                migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress())
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
            
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, await newMigrator.getAddress())).to.be.false;
            
            await migrationHelper.grantRole(MIGRATION_ROLE, await newMigrator.getAddress());
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, await newMigrator.getAddress())).to.be.true;
            
            await expect(
                migrationHelper.connect(newMigrator).markMigrated(await oldContract.getAddress(), await user1.getAddress())
            ).to.not.be.reverted;
            
            await migrationHelper.revokeRole(MIGRATION_ROLE, await newMigrator.getAddress());
            expect(await migrationHelper.hasRole(MIGRATION_ROLE, await newMigrator.getAddress())).to.be.false;
        });
    });

    describe("Gas Optimization", function () {
        it("Should efficiently handle large batches", async function () {
            const users = [];
            for (let i = 0; i < 50; i++) {
                const wallet = ethers.Wallet.createRandom();
                users.push(wallet.address);
            }

            const tx = await migrationHelper.connect(migrator).batchMarkMigrated(await oldContract.getAddress(), users);
            const receipt = await tx.wait();
            
            const gasPerUser = receipt.gasUsed / 50n;
            console.log(`Gas per user in batch: ${gasPerUser}`);
            
            // Should be efficient - less than 50k gas per user
            expect(Number(gasPerUser)).to.be.lessThan(50000);
        });
    });

    describe("Integration Scenarios", function () {
        it("Should support complete migration flow", async function () {
            // 1. Mark individual high-value users
            await migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress());
            
            // 2. Batch migrate regular users
            const regularUsers = [await user2.getAddress(), await user3.getAddress()];
            await migrationHelper.connect(migrator).batchMarkMigrated(await oldContract.getAddress(), regularUsers);
            
            // 3. Verify all migrations
            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user1.getAddress())).to.be.true;
            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user2.getAddress())).to.be.true;
            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user3.getAddress())).to.be.true;
            
            // 4. Check total count
            expect(await migrationHelper.migrationCount(await oldContract.getAddress())).to.equal(3);
            
            // 5. Attempt duplicate migration (should be safe)
            await expect(
                migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress())
            ).to.be.revertedWithCustomError(migrationHelper, "AlreadyMigrated");
        });

        it("Should handle multi-contract migration", async function () {
            // Migrate user1 from multiple old contracts
            await migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user1.getAddress());
            await migrationHelper.connect(migrator).markMigrated(await newContract.getAddress(), await user1.getAddress());
            
            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user1.getAddress())).to.be.true;
            expect(await migrationHelper.isMigrated(await newContract.getAddress(), await user1.getAddress())).to.be.true;
            
            // User2 only migrated from one contract
            await migrationHelper.connect(migrator).markMigrated(await oldContract.getAddress(), await user2.getAddress());
            
            expect(await migrationHelper.isMigrated(await oldContract.getAddress(), await user2.getAddress())).to.be.true;
            expect(await migrationHelper.isMigrated(await newContract.getAddress(), await user2.getAddress())).to.be.false;
        });
    });
});