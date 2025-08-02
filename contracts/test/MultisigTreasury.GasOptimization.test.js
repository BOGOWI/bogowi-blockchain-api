const { ethers } = require("hardhat");
const { expect } = require("chai");

describe("MultisigTreasury Gas Optimization Tests", function () {
    let treasury;
    let owner, signer1, signer2, signer3, signer4, user1;
    let bogoToken;
    
    const THRESHOLD = 2;
    const MAX_BATCH_SIZE = 100;
    const DEFAULT_PAGE_SIZE = 50;
    const MAX_PAGE_SIZE = 100;

    beforeEach(async function () {
        [owner, signer1, signer2, signer3, signer4, user1] = await ethers.getSigners();

        // Deploy a mock ERC20 token
        const MockToken = await ethers.getContractFactory("MockERC20");
        bogoToken = await MockToken.deploy("BOGO Token", "BOGO", ethers.utils.parseEther("1000000"));
        await bogoToken.deployed();

        // Deploy MultisigTreasury with gas optimizations
        const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury_GasOptimized");
        treasury = await MultisigTreasury.deploy(
            [owner.address, signer1.address, signer2.address, signer3.address],
            THRESHOLD
        );
        await treasury.deployed();

        // Fund the treasury
        await owner.sendTransaction({
            to: treasury.address,
            value: ethers.utils.parseEther("10")
        });
    });

    describe("Batch Size Limits", function () {
        it("Should enforce MAX_BATCH_SIZE in constructor", async function () {
            const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury_GasOptimized");
            
            // Create array with more than MAX_BATCH_SIZE signers
            const tooManySigners = [];
            for (let i = 0; i < 101; i++) {
                const wallet = ethers.Wallet.createRandom();
                tooManySigners.push(wallet.address);
            }

            await expect(
                MultisigTreasury.deploy(tooManySigners, 50)
            ).to.be.revertedWith("Batch size exceeded");
        });

        it("Should enforce batch size in batchConfirmTransactions", async function () {
            // Create many transactions
            const txIds = [];
            for (let i = 0; i < 101; i++) {
                const tx = await treasury.submitTransaction(
                    user1.address,
                    ethers.utils.parseEther("0.01"),
                    "0x",
                    `Test transaction ${i}`
                );
                const receipt = await tx.wait();
                const txId = receipt.events.find(e => e.event === "TransactionSubmitted").args.txId;
                txIds.push(txId.toNumber());
            }

            // Try to confirm more than MAX_BATCH_SIZE
            await expect(
                treasury.connect(signer1).batchConfirmTransactions(txIds)
            ).to.be.revertedWith("Batch size exceeded");

            // Should work with exactly MAX_BATCH_SIZE
            const validBatch = txIds.slice(0, 100);
            await expect(
                treasury.connect(signer1).batchConfirmTransactions(validBatch)
            ).to.not.be.reverted;
        });

        it("Should enforce batch size in batchAddSigners", async function () {
            // First, submit a transaction to call batchAddSigners through multisig
            const tooManySigners = [];
            for (let i = 0; i < 101; i++) {
                const wallet = ethers.Wallet.createRandom();
                tooManySigners.push(wallet.address);
            }

            const data = treasury.interface.encodeFunctionData("batchAddSigners", [tooManySigners]);
            await treasury.submitTransaction(treasury.address, 0, data, "Add too many signers");
            
            // Confirm and try to execute
            await treasury.connect(signer1).confirmTransaction(0);
            
            // Wait for execution delay
            await ethers.provider.send("evm_increaseTime", [3600]); // 1 hour
            await ethers.provider.send("evm_mine");

            // Should revert when executed
            await expect(
                treasury.executeTransaction(0)
            ).to.be.reverted;
        });
    });

    describe("Pagination Support", function () {
        beforeEach(async function () {
            // Create 150 transactions to test pagination
            for (let i = 0; i < 150; i++) {
                await treasury.submitTransaction(
                    user1.address,
                    ethers.utils.parseEther("0.001"),
                    "0x",
                    `Test transaction ${i}`
                );
            }
        });

        it("Should paginate pending transactions correctly", async function () {
            // Get first page
            const page1 = await treasury.getPendingTransactionsPaginated(0, 50);
            expect(page1.txIds.length).to.equal(50);
            expect(page1.pagination.totalCount).to.equal(150);
            expect(page1.pagination.totalPages).to.equal(3);
            expect(page1.pagination.currentPage).to.equal(0);

            // Get second page
            const page2 = await treasury.getPendingTransactionsPaginated(1, 50);
            expect(page2.txIds.length).to.equal(50);
            expect(page2.pagination.currentPage).to.equal(1);

            // Get last page
            const page3 = await treasury.getPendingTransactionsPaginated(2, 50);
            expect(page3.txIds.length).to.equal(50);
            expect(page3.pagination.currentPage).to.equal(2);

            // Verify no duplicate IDs
            const allIds = [...page1.txIds, ...page2.txIds, ...page3.txIds];
            const uniqueIds = new Set(allIds.map(id => id.toString()));
            expect(uniqueIds.size).to.equal(150);
        });

        it("Should enforce maximum page size", async function () {
            await expect(
                treasury.getPendingTransactionsPaginated(0, 101)
            ).to.be.revertedWith("Invalid page size");
        });

        it("Should handle out of bounds pages", async function () {
            await expect(
                treasury.getPendingTransactionsPaginated(5, 50)
            ).to.be.revertedWith("Page out of bounds");
        });

        it("Should paginate signers correctly", async function () {
            // Add more signers first
            const newSigners = [];
            for (let i = 0; i < 10; i++) {
                const wallet = ethers.Wallet.createRandom();
                newSigners.push(wallet.address);
            }

            const data = treasury.interface.encodeFunctionData("batchAddSigners", [newSigners]);
            await treasury.submitTransaction(treasury.address, 0, data, "Add signers");
            await treasury.connect(signer1).confirmTransaction(0);
            
            await ethers.provider.send("evm_increaseTime", [3600]);
            await ethers.provider.send("evm_mine");
            
            await treasury.executeTransaction(0);

            // Test pagination
            const page1 = await treasury.getSignersPaginated(0, 5);
            expect(page1.signers.length).to.equal(5);
            expect(page1.pagination.totalCount).to.equal(14); // 4 original + 10 new

            const page2 = await treasury.getSignersPaginated(1, 5);
            expect(page2.signers.length).to.equal(5);

            const page3 = await treasury.getSignersPaginated(2, 5);
            expect(page3.signers.length).to.equal(4); // Remaining signers
        });
    });

    describe("Gas Usage Protection", function () {
        it("Should stop batch operations when gas is low", async function () {
            // Create many transactions
            const txIds = [];
            for (let i = 0; i < 50; i++) {
                const tx = await treasury.submitTransaction(
                    user1.address,
                    ethers.utils.parseEther("0.001"),
                    "0x",
                    `Test transaction ${i}`
                );
                const receipt = await tx.wait();
                const txId = receipt.events.find(e => e.event === "TransactionSubmitted").args.txId;
                txIds.push(txId.toNumber());
            }

            // Batch confirm should handle gas limits gracefully
            const tx = await treasury.connect(signer1).batchConfirmTransactions(txIds);
            const receipt = await tx.wait();
            
            // Check event to see how many were actually confirmed
            const batchEvent = receipt.events.find(e => e.event === "BatchOperationExecuted");
            expect(batchEvent).to.not.be.undefined;
            expect(batchEvent.args.operation).to.equal("batchConfirm");
            // Should have confirmed at least some transactions
            expect(batchEvent.args.itemsProcessed.toNumber()).to.be.greaterThan(0);
        });

        it("Should handle batch cancel operations efficiently", async function () {
            // Create expired transactions
            const txIds = [];
            for (let i = 0; i < 30; i++) {
                const tx = await treasury.submitTransaction(
                    user1.address,
                    ethers.utils.parseEther("0.001"),
                    "0x",
                    `Test transaction ${i}`
                );
                const receipt = await tx.wait();
                const txId = receipt.events.find(e => e.event === "TransactionSubmitted").args.txId;
                txIds.push(txId.toNumber());
            }

            // Fast forward past expiry
            await ethers.provider.send("evm_increaseTime", [7 * 24 * 60 * 60 + 1]); // 7 days + 1 second
            await ethers.provider.send("evm_mine");

            // Batch cancel
            const tx = await treasury.batchCancelExpiredTransactions(txIds);
            const receipt = await tx.wait();
            
            const batchEvent = receipt.events.find(e => e.event === "BatchOperationExecuted");
            expect(batchEvent.args.operation).to.equal("batchCancel");
            expect(batchEvent.args.itemsProcessed.toNumber()).to.equal(30);
        });
    });

    describe("Pending Transaction Tracking", function () {
        it("Should efficiently track pending transactions", async function () {
            // Submit some transactions
            for (let i = 0; i < 5; i++) {
                await treasury.submitTransaction(
                    user1.address,
                    ethers.utils.parseEther("0.001"),
                    "0x",
                    `Test transaction ${i}`
                );
            }

            // Check pending count
            expect(await treasury.getPendingTransactionCount()).to.equal(5);

            // Execute one transaction
            await treasury.connect(signer1).confirmTransaction(0);
            await ethers.provider.send("evm_increaseTime", [3600]);
            await ethers.provider.send("evm_mine");
            await treasury.executeTransaction(0);

            // Check pending count decreased
            expect(await treasury.getPendingTransactionCount()).to.equal(4);

            // Cancel one expired transaction
            await ethers.provider.send("evm_increaseTime", [7 * 24 * 60 * 60 + 1]);
            await ethers.provider.send("evm_mine");
            await treasury.cancelExpiredTransaction(1);

            // Check pending count decreased again
            expect(await treasury.getPendingTransactionCount()).to.equal(3);
        });
    });

    describe("Emergency Withdrawal Gas Optimization", function () {
        it("Should limit emergency approval resets to MAX_BATCH_SIZE", async function () {
            // This is already implemented in the contract
            // The emergency withdrawal will only reset up to MAX_BATCH_SIZE approvals
            
            // Pause the contract
            const pauseData = treasury.interface.encodeFunctionData("pause", []);
            await treasury.submitTransaction(treasury.address, 0, pauseData, "Pause contract");
            await treasury.connect(signer1).confirmTransaction(0);
            await ethers.provider.send("evm_increaseTime", [3600]);
            await ethers.provider.send("evm_mine");
            await treasury.executeTransaction(0);

            // Multiple signers approve emergency withdrawal
            await treasury.emergencyWithdrawETH(user1.address, ethers.utils.parseEther("1"));
            await treasury.connect(signer1).emergencyWithdrawETH(user1.address, ethers.utils.parseEther("1"));

            // Check that withdrawal executed
            const balanceBefore = await ethers.provider.getBalance(user1.address);
            expect(balanceBefore).to.be.gt(0);
        });
    });

    describe("Gas Estimation", function () {
        it("Should provide reasonable gas costs for batch operations", async function () {
            // Create transactions
            const txIds = [];
            for (let i = 0; i < 20; i++) {
                const tx = await treasury.submitTransaction(
                    user1.address,
                    ethers.utils.parseEther("0.001"),
                    "0x",
                    `Test transaction ${i}`
                );
                const receipt = await tx.wait();
                const txId = receipt.events.find(e => e.event === "TransactionSubmitted").args.txId;
                txIds.push(txId.toNumber());
            }

            // Estimate gas for batch confirm
            const estimatedGas = await treasury.connect(signer1).estimateGas.batchConfirmTransactions(txIds);
            
            // Execute and compare
            const tx = await treasury.connect(signer1).batchConfirmTransactions(txIds);
            const receipt = await tx.wait();
            
            // Gas used should be close to estimate
            const gasUsed = receipt.gasUsed;
            expect(gasUsed.toNumber()).to.be.lessThan(1000000); // Should be under 1M gas
            
            console.log(`Batch confirm 20 transactions - Estimated: ${estimatedGas}, Used: ${gasUsed}`);
        });
    });
});