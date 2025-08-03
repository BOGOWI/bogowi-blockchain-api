const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("MultisigTreasury_GasOptimized", function () {
  let multisigTreasury;
  let mockERC20;
  let mockERC721;
  let mockERC1155;
  let owner;
  let signer1;
  let signer2;
  let signer3;
  let signer4;
  let signer5;
  let nonSigner;
  let recipient;
  let treasury;

  const THRESHOLD = 2;
  const MAX_SIGNERS = 20;
  const MAX_BATCH_SIZE = 100;
  const DEFAULT_PAGE_SIZE = 50;
  const MAX_PAGE_SIZE = 100;
  const TRANSACTION_EXPIRY = 7 * 24 * 60 * 60; // 7 days
  const EXECUTION_DELAY = 60 * 60; // 1 hour
  const MAX_GAS_LIMIT = 5000000;

  beforeEach(async function () {
    [owner, signer1, signer2, signer3, signer4, signer5, nonSigner, recipient, treasury] = await ethers.getSigners();

    // Deploy mock tokens
    const MockERC20 = await ethers.getContractFactory("contracts/mocks/MockERC20.sol:MockERC20");
    mockERC20 = await MockERC20.deploy("Test Token", "TEST", ethers.parseEther("1000000"));
    await mockERC20.waitForDeployment();

    const MockERC721 = await ethers.getContractFactory("contracts/mocks/MockERC721.sol:MockERC721");
    mockERC721 = await MockERC721.deploy("Test NFT", "TNFT");
    await mockERC721.waitForDeployment();

    const MockERC1155 = await ethers.getContractFactory("contracts/mocks/MockERC1155.sol:MockERC1155");
    mockERC1155 = await MockERC1155.deploy("https://test.com/{id}");
    await mockERC1155.waitForDeployment();

    // Deploy MultisigTreasury_GasOptimized
    const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury_GasOptimized");
    const signers = [await signer1.getAddress(), await signer2.getAddress(), await signer3.getAddress()];
    multisigTreasury = await MultisigTreasury.deploy(signers, THRESHOLD);
    await multisigTreasury.waitForDeployment();

    // Fund the treasury
    await owner.sendTransaction({
      to: await multisigTreasury.getAddress(),
      value: ethers.parseEther("10")
    });

    // Transfer some tokens to treasury
    await mockERC20.transfer(await multisigTreasury.getAddress(), ethers.parseEther("1000"));
  });

  describe("Deployment", function () {
    it("Should deploy with correct initial state", async function () {
      expect(await multisigTreasury.threshold()).to.equal(THRESHOLD);
      expect(await multisigTreasury.signerCount()).to.equal(3);
      expect(await multisigTreasury.transactionCount()).to.equal(0);
      expect(await multisigTreasury.autoExecuteEnabled()).to.be.true;
      expect(await multisigTreasury.restrictFunctionCalls()).to.be.false;
    });

    it("Should set signers correctly", async function () {
      const signers = await multisigTreasury.getSigners();
      expect(signers).to.have.length(3);
      expect(signers).to.include(await signer1.getAddress());
      expect(signers).to.include(await signer2.getAddress());
      expect(signers).to.include(await signer3.getAddress());
    });

    it("Should revert with invalid constructor parameters", async function () {
      const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury_GasOptimized");
      
      // Empty signers array
      await expect(
        MultisigTreasury.deploy([], 1)
      ).to.be.revertedWith("Signers required");

      // Zero threshold
      await expect(
        MultisigTreasury.deploy([await signer1.getAddress()], 0)
      ).to.be.revertedWith("Invalid threshold");

      // Threshold greater than signers
      await expect(
        MultisigTreasury.deploy([await signer1.getAddress()], 2)
      ).to.be.revertedWith("Invalid threshold");

      // Duplicate signers
      await expect(
        MultisigTreasury.deploy([await signer1.getAddress(), await signer1.getAddress()], 1)
      ).to.be.revertedWith("Duplicate signer");

      // Zero address signer
      await expect(
        MultisigTreasury.deploy([ethers.ZeroAddress], 1)
      ).to.be.revertedWith("Invalid signer");
    });

    it("Should handle maximum signers in constructor", async function () {
      const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury_GasOptimized");
      const maxSigners = Array.from({ length: MAX_SIGNERS }, (_, i) => 
        ethers.Wallet.createRandom().address
      );
      
      const treasury = await MultisigTreasury.deploy(maxSigners, MAX_SIGNERS);
      await treasury.waitForDeployment();
      
      expect(await treasury.signerCount()).to.equal(MAX_SIGNERS);
    });

    it("Should revert when exceeding max signers in constructor", async function () {
      const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury_GasOptimized");
      const tooManySigners = Array.from({ length: MAX_SIGNERS + 1 }, (_, i) => 
        ethers.Wallet.createRandom().address
      );
      
      await expect(
        MultisigTreasury.deploy(tooManySigners, 1)
      ).to.be.revertedWith("Too many signers");
    });
  });

  describe("Transaction Management", function () {
    describe("Submit Transaction", function () {
      it("Should submit transaction successfully", async function () {
        const tx = await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("1"),
          "0x",
          "Test transaction"
        );

        await expect(tx)
          .to.emit(multisigTreasury, "TransactionSubmitted")
          .withArgs(0, await signer1.getAddress(), await recipient.getAddress(), ethers.parseEther("1"));

        await expect(tx)
          .to.emit(multisigTreasury, "TransactionConfirmed")
          .withArgs(0, await signer1.getAddress());

        expect(await multisigTreasury.transactionCount()).to.equal(1);
      });

      it("Should revert when non-signer submits transaction", async function () {
        await expect(
          multisigTreasury.connect(nonSigner).submitTransaction(
            await recipient.getAddress(),
            ethers.parseEther("1"),
            "0x",
            "Test transaction"
          )
        ).to.be.revertedWith("Not a signer");
      });

      it("Should revert with zero address recipient", async function () {
        await expect(
          multisigTreasury.connect(signer1).submitTransaction(
            ethers.ZeroAddress,
            ethers.parseEther("1"),
            "0x",
            "Test transaction"
          )
        ).to.be.revertedWith("Invalid recipient");
      });

      it("Should revert when paused", async function () {
        // Pause the contract
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("pause"),
          "Pause contract"
        );
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        
        // Wait for execution delay
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        await multisigTreasury.connect(signer1).executeTransaction(0);

        await expect(
          multisigTreasury.connect(signer1).submitTransaction(
            await recipient.getAddress(),
            ethers.parseEther("1"),
            "0x",
            "Test transaction"
          )
        ).to.be.revertedWithCustomError(multisigTreasury, "EnforcedPause");
      });
    });

    describe("Confirm Transaction", function () {
      beforeEach(async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("1"),
          "0x",
          "Test transaction"
        );
      });

      it("Should confirm transaction successfully", async function () {
        await expect(
          multisigTreasury.connect(signer2).confirmTransaction(0)
        ).to.emit(multisigTreasury, "TransactionConfirmed")
          .withArgs(0, await signer2.getAddress());

        expect(await multisigTreasury.hasConfirmed(0, await signer2.getAddress())).to.be.true;
        expect(await multisigTreasury.getConfirmationCount(0)).to.equal(2);
      });

      it("Should revert when already confirmed", async function () {
        await expect(
          multisigTreasury.connect(signer1).confirmTransaction(0)
        ).to.be.revertedWith("Already confirmed");
      });

      it("Should revert for non-existent transaction", async function () {
        await expect(
          multisigTreasury.connect(signer2).confirmTransaction(999)
        ).to.be.revertedWith("Transaction does not exist");
      });

      it("Should auto-execute when threshold reached and delay passed", async function () {
        // Wait for execution delay
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        const initialBalance = await ethers.provider.getBalance(await recipient.getAddress());
        
        await expect(
          multisigTreasury.connect(signer2).confirmTransaction(0)
        ).to.emit(multisigTreasury, "TransactionExecuted")
          .withArgs(0, await signer2.getAddress());

        const finalBalance = await ethers.provider.getBalance(await recipient.getAddress());
        expect(finalBalance - initialBalance).to.equal(ethers.parseEther("1"));
      });

      it("Should not auto-execute before delay", async function () {
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        
        const transaction = await multisigTreasury.getTransaction(0);
        expect(transaction.executed).to.be.false;
      });
    });

    describe("Batch Confirm Transactions", function () {
      beforeEach(async function () {
        // Submit multiple transactions
        for (let i = 0; i < 5; i++) {
          await multisigTreasury.connect(signer1).submitTransaction(
            await recipient.getAddress(),
            ethers.parseEther("0.1"),
            "0x",
            `Test transaction ${i}`
          );
        }
      });

      it("Should batch confirm multiple transactions", async function () {
        const txIds = [0, 1, 2, 3, 4];
        
        await expect(
          multisigTreasury.connect(signer2).batchConfirmTransactions(txIds)
        ).to.emit(multisigTreasury, "BatchOperationExecuted")
          .withArgs("batchConfirm", 5);

        for (let i = 0; i < 5; i++) {
          expect(await multisigTreasury.hasConfirmed(i, await signer2.getAddress())).to.be.true;
        }
      });

      it("Should skip invalid transactions in batch", async function () {
        const txIds = [0, 999, 1, 2]; // 999 doesn't exist
        
        await expect(
          multisigTreasury.connect(signer2).batchConfirmTransactions(txIds)
        ).to.emit(multisigTreasury, "BatchOperationExecuted")
          .withArgs("batchConfirm", 3); // Only 3 valid transactions
      });

      it("Should revert when batch size exceeded", async function () {
        const largeBatch = Array.from({ length: MAX_BATCH_SIZE + 1 }, (_, i) => i);
        
        await expect(
          multisigTreasury.connect(signer2).batchConfirmTransactions(largeBatch)
        ).to.be.revertedWith("Batch size exceeded");
      });

      it("Should handle empty batch", async function () {
        await expect(
          multisigTreasury.connect(signer2).batchConfirmTransactions([])
        ).to.emit(multisigTreasury, "BatchOperationExecuted")
          .withArgs("batchConfirm", 0);
      });
    });

    describe("Revoke Confirmation", function () {
      beforeEach(async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("1"),
          "0x",
          "Test transaction"
        );
        await multisigTreasury.connect(signer2).confirmTransaction(0);
      });

      it("Should revoke confirmation successfully", async function () {
        await expect(
          multisigTreasury.connect(signer2).revokeConfirmation(0)
        ).to.emit(multisigTreasury, "ConfirmationRevoked")
          .withArgs(0, await signer2.getAddress());

        expect(await multisigTreasury.hasConfirmed(0, await signer2.getAddress())).to.be.false;
        expect(await multisigTreasury.getConfirmationCount(0)).to.equal(1);
      });

      it("Should revert when not confirmed", async function () {
        await expect(
          multisigTreasury.connect(signer3).revokeConfirmation(0)
        ).to.be.revertedWith("Not confirmed");
      });
    });

    describe("Execute Transaction", function () {
      beforeEach(async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("1"),
          "0x",
          "Test transaction"
        );
        await multisigTreasury.connect(signer2).confirmTransaction(0);
      });

      it("Should execute transaction successfully", async function () {
        // Wait for execution delay
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        const initialBalance = await ethers.provider.getBalance(await recipient.getAddress());
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.emit(multisigTreasury, "TransactionExecuted")
          .withArgs(0, await signer1.getAddress());

        const finalBalance = await ethers.provider.getBalance(await recipient.getAddress());
        expect(finalBalance - initialBalance).to.equal(ethers.parseEther("1"));
        
        const transaction = await multisigTreasury.getTransaction(0);
        expect(transaction.executed).to.be.true;
      });

      it("Should revert with insufficient confirmations", async function () {
        // Revoke one confirmation
        await multisigTreasury.connect(signer2).revokeConfirmation(0);
        
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.be.revertedWith("Insufficient confirmations");
      });

      it("Should revert before execution delay", async function () {
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.be.revertedWith("Execution delay not met");
      });

      it("Should revert when already executed", async function () {
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        await multisigTreasury.connect(signer1).executeTransaction(0);
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.be.revertedWith("Transaction already executed");
      });
    });

    describe("Transaction Expiry", function () {
      it("Should cancel expired transaction", async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("1"),
          "0x",
          "Test transaction"
        );
        
        // Fast forward past expiry
        await ethers.provider.send("evm_increaseTime", [TRANSACTION_EXPIRY + 1]);
        
        await expect(
          multisigTreasury.connect(signer1).cancelExpiredTransaction(0)
        ).to.emit(multisigTreasury, "TransactionCancelled")
          .withArgs(0);

        const transaction = await multisigTreasury.getTransaction(0);
        expect(transaction.executed).to.be.true;
      });

      it("Should revert when cancelling non-expired transaction", async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("1"),
          "0x",
          "Test transaction"
        );
        
        await expect(
          multisigTreasury.connect(signer1).cancelExpiredTransaction(0)
        ).to.be.revertedWith("Transaction not expired");
      });

      it("Should batch cancel expired transactions", async function () {
        // Submit transactions
        for (let i = 0; i < 5; i++) {
          await multisigTreasury.connect(signer1).submitTransaction(
            await recipient.getAddress(),
            ethers.parseEther("0.1"),
            "0x",
            `Test transaction ${i}`
          );
        }
        
        // Fast forward past expiry
        await ethers.provider.send("evm_increaseTime", [TRANSACTION_EXPIRY + 1]);
        
        const txIds = [0, 1, 2, 3, 4];
        await expect(
          multisigTreasury.connect(signer1).batchCancelExpiredTransactions(txIds)
        ).to.emit(multisigTreasury, "BatchOperationExecuted")
          .withArgs("batchCancel", 5);
      });

      it("Should check if transaction is expired", async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("1"),
          "0x",
          "Test transaction"
        );
        
        expect(await multisigTreasury.isTransactionExpired(0)).to.be.false;
        
        await ethers.provider.send("evm_increaseTime", [TRANSACTION_EXPIRY + 1]);
        await ethers.provider.send("evm_mine");
        
        expect(await multisigTreasury.isTransactionExpired(0)).to.be.true;
      });
    });
  });

  describe("Signer Management", function () {
    describe("Add Signer", function () {
      it("Should add signer through multisig", async function () {
        const newSigner = await signer4.getAddress();
        
        // Submit transaction to add signer
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("addSigner", [newSigner]),
          "Add new signer"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.emit(multisigTreasury, "SignerAdded")
          .withArgs(newSigner);

        expect(await multisigTreasury.signerCount()).to.equal(4);
        const signers = await multisigTreasury.getSigners();
        expect(signers).to.include(newSigner);
      });

      it("Should revert when adding duplicate signer", async function () {
        const existingSigner = await signer1.getAddress();
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("addSigner", [existingSigner]),
          "Add duplicate signer"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        try {
          await multisigTreasury.connect(signer1).executeTransaction(0);
          expect.fail("Expected transaction to revert");
        } catch (error) {
          expect(error.message).to.include("Already a signer");
        }
      });

      it("Should revert when adding zero address", async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("addSigner", [ethers.ZeroAddress]),
          "Add zero address signer"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        try {
          await multisigTreasury.connect(signer1).executeTransaction(0);
          expect.fail("Expected transaction to revert");
        } catch (error) {
          expect(error.message).to.include("Invalid signer");
        }
      });
    });

    describe("Batch Add Signers", function () {
      it("Should batch add multiple signers", async function () {
        const newSigners = [await signer4.getAddress(), await signer5.getAddress()];
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("batchAddSigners", [newSigners]),
          "Batch add signers"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.emit(multisigTreasury, "BatchOperationExecuted")
          .withArgs("batchAddSigners", 2);

        expect(await multisigTreasury.signerCount()).to.equal(5);
      });

      it("Should revert when batch size exceeded", async function () {
        const largeBatch = Array.from({ length: MAX_BATCH_SIZE + 1 }, () => 
          ethers.Wallet.createRandom().address
        );
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("batchAddSigners", [largeBatch]),
          "Large batch add"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        try {
          await multisigTreasury.connect(signer1).executeTransaction(0);
          expect.fail("Expected transaction to revert");
        } catch (error) {
          expect(error.message).to.include("Batch size exceeded");
        }
      });
    });

    describe("Remove Signer", function () {
      it("Should remove signer through multisig", async function () {
        const signerToRemove = await signer3.getAddress();
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("removeSigner", [signerToRemove]),
          "Remove signer"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.emit(multisigTreasury, "SignerRemoved")
          .withArgs(signerToRemove);

        expect(await multisigTreasury.signerCount()).to.equal(2);
        const signers = await multisigTreasury.getSigners();
        expect(signers).to.not.include(signerToRemove);
      });

      it("Should revert when removing would break threshold", async function () {
        // Try to remove signer when only 3 signers and threshold is 2
        // Removing one would leave 2 signers, which equals threshold (still valid)
        // But removing another would break it
        
        // First remove one signer
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("removeSigner", [await signer3.getAddress()]),
          "Remove signer"
        );
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        await multisigTreasury.connect(signer1).executeTransaction(0);
        
        // Now try to remove another (would leave 1 signer with threshold 2)
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("removeSigner", [await signer2.getAddress()]),
          "Remove another signer"
        );
        await multisigTreasury.connect(signer2).confirmTransaction(1);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        try {
          await multisigTreasury.connect(signer1).executeTransaction(1);
          expect.fail("Expected transaction to revert");
        } catch (error) {
          expect(error.message).to.include("Would break threshold");
        }
      });
    });

    describe("Change Threshold", function () {
      it("Should change threshold through multisig", async function () {
        const newThreshold = 3;
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("changeThreshold", [newThreshold]),
          "Change threshold"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.emit(multisigTreasury, "ThresholdChanged")
          .withArgs(THRESHOLD, newThreshold);

        expect(await multisigTreasury.threshold()).to.equal(newThreshold);
      });

      it("Should revert with invalid threshold", async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("changeThreshold", [0]),
          "Invalid threshold"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        try {
          await multisigTreasury.connect(signer1).executeTransaction(0);
          expect.fail("Expected transaction to revert");
        } catch (error) {
          expect(error.message).to.include("Invalid threshold");
        }
      });
    });
  });

  describe("Pagination", function () {
    beforeEach(async function () {
      // Add more signers for pagination testing
      const newSigners = [];
      for (let i = 0; i < 10; i++) {
        newSigners.push(ethers.Wallet.createRandom().address);
      }
      
      await multisigTreasury.connect(signer1).submitTransaction(
        await multisigTreasury.getAddress(),
        0,
        multisigTreasury.interface.encodeFunctionData("batchAddSigners", [newSigners]),
        "Add signers for pagination"
      );
      await multisigTreasury.connect(signer2).confirmTransaction(0);
      await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
      await multisigTreasury.connect(signer1).executeTransaction(0);
    });

    describe("Signer Pagination", function () {
      it("Should paginate signers correctly", async function () {
        const pageSize = 5;
        const page0 = await multisigTreasury.getSignersPaginated(0, pageSize);
        
        expect(page0.signerAddresses).to.have.length(pageSize);
        expect(page0.pagination.totalCount).to.equal(13); // 3 original + 10 added
        expect(page0.pagination.pageSize).to.equal(pageSize);
        expect(page0.pagination.currentPage).to.equal(0);
        expect(page0.pagination.totalPages).to.equal(3); // ceil(13/5)
      });

      it("Should handle last page with fewer items", async function () {
        const pageSize = 5;
        const lastPage = await multisigTreasury.getSignersPaginated(2, pageSize);
        
        expect(lastPage.signerAddresses).to.have.length(3); // 13 % 5 = 3
        expect(lastPage.pagination.currentPage).to.equal(2);
      });

      it("Should revert with invalid page size", async function () {
        await expect(
          multisigTreasury.getSignersPaginated(0, 0)
        ).to.be.revertedWith("Invalid page size");

        await expect(
          multisigTreasury.getSignersPaginated(0, MAX_PAGE_SIZE + 1)
        ).to.be.revertedWith("Invalid page size");
      });

      it("Should revert with page out of bounds", async function () {
        await expect(
          multisigTreasury.getSignersPaginated(999, 5)
        ).to.be.revertedWith("Page out of bounds");
      });
    });

    describe("Transaction Pagination", function () {
      beforeEach(async function () {
        // Submit multiple transactions (starting from ID 1 since 0 is used in parent beforeEach)
        for (let i = 0; i < 15; i++) {
          await multisigTreasury.connect(signer1).submitTransaction(
            await recipient.getAddress(),
            ethers.parseEther("0.01"),
            "0x",
            `Transaction ${i + 1}`
          );
        }
      });

      it("Should paginate pending transactions", async function () {
        const pageSize = 5;
        const page0 = await multisigTreasury.getPendingTransactionsPaginated(0, pageSize);
        
        expect(page0.txIds).to.have.length(pageSize);
        expect(page0.pagination.totalCount).to.equal(15); // Only the 15 from this test (parent tx is executed)
        expect(page0.pagination.totalPages).to.equal(3); // ceil(15/5)
      });

      it("Should return correct pending transaction count", async function () {
        expect(await multisigTreasury.getPendingTransactionCount()).to.equal(15); // Only the 15 from this test (parent tx is executed)
        
        // Execute one transaction (use ID 1 to avoid conflict with parent)
        await multisigTreasury.connect(signer2).confirmTransaction(1);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        await multisigTreasury.connect(signer1).executeTransaction(1);
        
        expect(await multisigTreasury.getPendingTransactionCount()).to.equal(14);
      });

      it("Should handle empty pending transactions", async function () {
        // Deploy a fresh contract to test empty state
        const MultisigTreasury = await ethers.getContractFactory("contracts/MultisigTreasury_GasOptimized.sol:MultisigTreasury_GasOptimized");
        const freshTreasury = await MultisigTreasury.deploy(
          [await signer1.getAddress(), await signer2.getAddress(), await signer3.getAddress()],
          2
        );
        
        const result = await freshTreasury.getPendingTransactionsPaginated(0, 5);
        expect(result.txIds).to.have.length(0);
        expect(result.pagination.totalCount).to.equal(0);
      });
    });
  });

  describe("Emergency Functions", function () {
    beforeEach(async function () {
      // Pause the contract for emergency functions
      await multisigTreasury.connect(signer1).submitTransaction(
        await multisigTreasury.getAddress(),
        0,
        multisigTreasury.interface.encodeFunctionData("pause"),
        "Pause for emergency"
      );
      await multisigTreasury.connect(signer2).confirmTransaction(0);
      await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
      await multisigTreasury.connect(signer1).executeTransaction(0);
    });

    describe("Emergency Withdraw ETH", function () {
      it("Should require multiple approvals for emergency withdraw", async function () {
        const withdrawAmount = ethers.parseEther("1");
        
        await expect(
          multisigTreasury.connect(signer1).emergencyWithdrawETH(
            await recipient.getAddress(),
            withdrawAmount
          )
        ).to.emit(multisigTreasury, "EmergencyApprovalGranted")
          .withArgs(await signer1.getAddress());

        expect(await multisigTreasury.hasEmergencyApproval(await signer1.getAddress())).to.be.true;
        expect(await multisigTreasury.getEmergencyApprovalCount()).to.equal(1);
      });

      it("Should execute emergency withdraw when threshold reached", async function () {
        const withdrawAmount = ethers.parseEther("1");
        const initialBalance = await ethers.provider.getBalance(await recipient.getAddress());
        
        // First approval
        await multisigTreasury.connect(signer1).emergencyWithdrawETH(
          await recipient.getAddress(),
          withdrawAmount
        );
        
        // Second approval (reaches threshold)
        await expect(
          multisigTreasury.connect(signer2).emergencyWithdrawETH(
            await recipient.getAddress(),
            withdrawAmount
          )
        ).to.emit(multisigTreasury, "EmergencyWithdraw")
          .withArgs(ethers.ZeroAddress, await recipient.getAddress(), withdrawAmount);

        const finalBalance = await ethers.provider.getBalance(await recipient.getAddress());
        expect(finalBalance - initialBalance).to.equal(withdrawAmount);
        
        // Approvals should be reset
        expect(await multisigTreasury.getEmergencyApprovalCount()).to.equal(0);
      });

      it("Should revert when exceeding 50% balance limit", async function () {
        const balance = await ethers.provider.getBalance(await multisigTreasury.getAddress());
        const excessiveAmount = balance / 2n + 1n;
        
        await expect(
          multisigTreasury.connect(signer1).emergencyWithdrawETH(
            await recipient.getAddress(),
            excessiveAmount
          )
        ).to.be.revertedWith("Amount exceeds 50% of balance");
      });

      it("Should revert when already approved", async function () {
        const withdrawAmount = ethers.parseEther("1");
        
        await multisigTreasury.connect(signer1).emergencyWithdrawETH(
          await recipient.getAddress(),
          withdrawAmount
        );
        
        await expect(
          multisigTreasury.connect(signer1).emergencyWithdrawETH(
            await recipient.getAddress(),
            withdrawAmount
          )
        ).to.be.revertedWith("Already approved emergency");
      });

      it("Should revert with zero address recipient", async function () {
        await expect(
          multisigTreasury.connect(signer1).emergencyWithdrawETH(
            ethers.ZeroAddress,
            ethers.parseEther("1")
          )
        ).to.be.revertedWith("Invalid recipient");
      });

      it("Should revert when not paused", async function () {
        // First, we need to unpause the contract through a multisig transaction
        // But we need to do this while the contract is still paused
        
        // Create a fresh multisig for this test that starts unpaused
        const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury_GasOptimized");
        const signers = [await signer1.getAddress(), await signer2.getAddress(), await signer3.getAddress()];
        const unpausedTreasury = await MultisigTreasury.deploy(signers, THRESHOLD);
        await unpausedTreasury.waitForDeployment();
        
        // Fund the unpaused treasury
        await owner.sendTransaction({
          to: await unpausedTreasury.getAddress(),
          value: ethers.parseEther("10")
        });
        
        // Emergency withdraw should revert when not paused
        await expect(
          unpausedTreasury.connect(signer1).emergencyWithdrawETH(
            await recipient.getAddress(),
            ethers.parseEther("1")
          )
        ).to.be.revertedWithCustomError(unpausedTreasury, "ExpectedPause");
      });
    });
  });

  describe("Token Management", function () {
    describe("ERC20 Transfers", function () {
      it("Should transfer ERC20 tokens through multisig", async function () {
        const transferAmount = ethers.parseEther("100");
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("transferERC20", [
            await mockERC20.getAddress(),
            await recipient.getAddress(),
            transferAmount
          ]),
          "Transfer ERC20"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        await multisigTreasury.connect(signer1).executeTransaction(0);
        
        expect(await mockERC20.balanceOf(await recipient.getAddress())).to.equal(transferAmount);
      });
    });

    describe("ERC721 Transfers", function () {
      beforeEach(async function () {
        // Mint NFT to treasury
        await mockERC721.mint(await multisigTreasury.getAddress(), 1);
      });

      it("Should transfer ERC721 tokens through multisig", async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("transferERC721", [
            await mockERC721.getAddress(),
            await recipient.getAddress(),
            1
          ]),
          "Transfer ERC721"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        await multisigTreasury.connect(signer1).executeTransaction(0);
        
        expect(await mockERC721.ownerOf(1)).to.equal(await recipient.getAddress());
      });
    });

    describe("ERC1155 Transfers", function () {
      beforeEach(async function () {
        // Mint ERC1155 to treasury
        await mockERC1155.mint(await multisigTreasury.getAddress(), 1, 100, "0x");
      });

      it("Should transfer ERC1155 tokens through multisig", async function () {
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("transferERC1155", [
            await mockERC1155.getAddress(),
            await recipient.getAddress(),
            1,
            50,
            "0x"
          ]),
          "Transfer ERC1155"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        await multisigTreasury.connect(signer1).executeTransaction(0);
        
        expect(await mockERC1155.balanceOf(await recipient.getAddress(), 1)).to.equal(50);
      });
    });
  });

  describe("Configuration", function () {
    describe("Auto Execute Toggle", function () {
      it("Should toggle auto execute through multisig", async function () {
        expect(await multisigTreasury.autoExecuteEnabled()).to.be.true;
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("toggleAutoExecute"),
          "Toggle auto execute"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.emit(multisigTreasury, "AutoExecuteToggled")
          .withArgs(false);

        expect(await multisigTreasury.autoExecuteEnabled()).to.be.false;
      });
    });

    describe("Function Restrictions", function () {
      it("Should toggle function restrictions", async function () {
        expect(await multisigTreasury.restrictFunctionCalls()).to.be.false;
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("toggleFunctionRestrictions"),
          "Toggle function restrictions"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        await multisigTreasury.connect(signer1).executeTransaction(0);
        
        expect(await multisigTreasury.restrictFunctionCalls()).to.be.true;
      });

      it("Should set function allowance", async function () {
        const targetContract = await mockERC20.getAddress();
        const selector = mockERC20.interface.getFunction("transfer").selector;
        
        await multisigTreasury.connect(signer1).submitTransaction(
          await multisigTreasury.getAddress(),
          0,
          multisigTreasury.interface.encodeFunctionData("setFunctionAllowance", [
            targetContract,
            selector,
            true
          ]),
          "Set function allowance"
        );
        
        await multisigTreasury.connect(signer2).confirmTransaction(0);
        await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
        
        await expect(
          multisigTreasury.connect(signer1).executeTransaction(0)
        ).to.emit(multisigTreasury, "FunctionAllowanceSet")
          .withArgs(targetContract, selector, true);

        expect(await multisigTreasury.allowedFunctions(targetContract, selector)).to.be.true;
      });
    });
  });

  describe("Token Receipt", function () {
    it("Should handle ETH deposits", async function () {
      const depositAmount = ethers.parseEther("1");
      
      await expect(
        owner.sendTransaction({
          to: await multisigTreasury.getAddress(),
          value: depositAmount
        })
      ).to.emit(multisigTreasury, "Deposit")
        .withArgs(await owner.getAddress(), depositAmount);
    });

    it("Should handle ERC721 receipt", async function () {
      // Mint NFT to treasury - this will trigger the onERC721Received callback
      await mockERC721.mint(await multisigTreasury.getAddress(), 1);
      
      // Verify the NFT was received
      expect(await mockERC721.ownerOf(1)).to.equal(await multisigTreasury.getAddress());
    });

    it("Should handle ERC1155 receipt", async function () {
      // Mint ERC1155 to treasury - this will trigger the onERC1155Received callback
      await mockERC1155.mint(await multisigTreasury.getAddress(), 1, 100, "0x");
      
      // Verify the tokens were received
      expect(await mockERC1155.balanceOf(await multisigTreasury.getAddress(), 1)).to.equal(100);
    });
  });

  describe("Gas Optimization Features", function () {
    it("Should respect batch size limits", async function () {
      // Test with exactly MAX_BATCH_SIZE
      const maxBatch = Array.from({ length: MAX_BATCH_SIZE }, (_, i) => i);
      
      await expect(
        multisigTreasury.connect(signer2).batchConfirmTransactions(maxBatch)
      ).to.not.be.reverted;
    });

    it("Should handle gas limit protection in batch operations", async function () {
      // Submit many transactions
      for (let i = 0; i < 50; i++) {
        await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("0.001"),
          "0x",
          `Transaction ${i}`
        );
      }
      
      const txIds = Array.from({ length: 50 }, (_, i) => i);
      
      // Should complete without running out of gas
      await expect(
        multisigTreasury.connect(signer2).batchConfirmTransactions(txIds)
      ).to.not.be.reverted;
    });

    it("Should maintain efficient pending transaction tracking", async function () {
      // Submit transactions
      for (let i = 0; i < 10; i++) {
        await multisigTreasury.connect(signer1).submitTransaction(
          await recipient.getAddress(),
          ethers.parseEther("0.001"),
          "0x",
          `Transaction ${i}`
        );
      }
      
      expect(await multisigTreasury.getPendingTransactionCount()).to.equal(10);
      
      // Execute one transaction
      await multisigTreasury.connect(signer2).confirmTransaction(5);
      await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
      await multisigTreasury.connect(signer1).executeTransaction(5);
      
      expect(await multisigTreasury.getPendingTransactionCount()).to.equal(9);
    });
  });

  describe("Edge Cases and Security", function () {
    it("Should handle reentrancy protection", async function () {
      // Deploy a malicious contract that tries to reenter
      const MaliciousReentrant = await ethers.getContractFactory("contracts/test/TestHelpers.sol:ReentrancyAttacker");
      const attacker = await MaliciousReentrant.deploy(await multisigTreasury.getAddress());
      await attacker.waitForDeployment();
      
      // Submit transaction to send ETH to attacker (which will try to reenter)
      await multisigTreasury.connect(signer1).submitTransaction(
        await attacker.getAddress(),
        ethers.parseEther("1"),
        "0x",
        "Send to attacker"
      );
      
      await multisigTreasury.connect(signer2).confirmTransaction(0);
      await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
      
      // The transaction should execute successfully but the attacker's reentrancy attempt should fail
      // The ReentrancyGuard prevents the nested call from succeeding
      await expect(
        multisigTreasury.connect(signer1).executeTransaction(0)
      ).to.emit(multisigTreasury, "TransactionExecuted");
    });

    it("Should handle large data payloads", async function () {
      const largeData = "0x" + "00".repeat(1000); // 1000 bytes of data
      
      await multisigTreasury.connect(signer1).submitTransaction(
        await recipient.getAddress(),
        0,
        largeData,
        "Large data transaction"
      );
      
      const transaction = await multisigTreasury.getTransaction(0);
      expect(transaction.data).to.equal(largeData);
    });

    it("Should handle zero value transactions", async function () {
      await multisigTreasury.connect(signer1).submitTransaction(
        await recipient.getAddress(),
        0,
        "0x",
        "Zero value transaction"
      );
      
      await multisigTreasury.connect(signer2).confirmTransaction(0);
      await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
      
      await expect(
        multisigTreasury.connect(signer1).executeTransaction(0)
      ).to.emit(multisigTreasury, "TransactionExecuted");
    });

    it("Should handle contract self-calls", async function () {
      // Submit transaction to call a function on the treasury itself
      await multisigTreasury.connect(signer1).submitTransaction(
        await multisigTreasury.getAddress(),
        0,
        multisigTreasury.interface.encodeFunctionData("toggleAutoExecute"),
        "Self-call transaction"
      );
      
      await multisigTreasury.connect(signer2).confirmTransaction(0);
      await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
      
      await expect(
        multisigTreasury.connect(signer1).executeTransaction(0)
      ).to.emit(multisigTreasury, "AutoExecuteToggled");
    });

    it("Should handle maximum uint256 values", async function () {
      const maxUint256 = ethers.MaxUint256;
      
      // This should revert due to insufficient balance, but not due to overflow
      await multisigTreasury.connect(signer1).submitTransaction(
        await recipient.getAddress(),
        maxUint256,
        "0x",
        "Max value transaction"
      );
      
      await multisigTreasury.connect(signer2).confirmTransaction(0);
      await ethers.provider.send("evm_increaseTime", [EXECUTION_DELAY + 1]);
      
      // Should fail due to insufficient balance, not overflow
      await expect(
        multisigTreasury.connect(signer1).executeTransaction(0)
      ).to.be.reverted;
    });
  });

  describe("View Functions", function () {
    beforeEach(async function () {
      await multisigTreasury.connect(signer1).submitTransaction(
        await recipient.getAddress(),
        ethers.parseEther("1"),
        "0x",
        "Test transaction"
      );
    });

    it("Should return correct transaction details", async function () {
      const transaction = await multisigTreasury.getTransaction(0);
      
      expect(transaction.to).to.equal(await recipient.getAddress());
      expect(transaction.value).to.equal(ethers.parseEther("1"));
      expect(transaction.data).to.equal("0x");
      expect(transaction.description).to.equal("Test transaction");
      expect(transaction.executed).to.be.false;
      expect(transaction.confirmationCount).to.equal(1);
    });

    it("Should return correct confirmation status", async function () {
      expect(await multisigTreasury.hasConfirmed(0, await signer1.getAddress())).to.be.true;
      expect(await multisigTreasury.hasConfirmed(0, await signer2.getAddress())).to.be.false;
      
      await multisigTreasury.connect(signer2).confirmTransaction(0);
      
      expect(await multisigTreasury.hasConfirmed(0, await signer2.getAddress())).to.be.true;
    });

    it("Should return correct confirmation count", async function () {
      expect(await multisigTreasury.getConfirmationCount(0)).to.equal(1);
      
      await multisigTreasury.connect(signer2).confirmTransaction(0);
      
      expect(await multisigTreasury.getConfirmationCount(0)).to.equal(2);
    });

    it("Should return correct expiry status", async function () {
      expect(await multisigTreasury.isTransactionExpired(0)).to.be.false;
      
      await ethers.provider.send("evm_increaseTime", [TRANSACTION_EXPIRY + 1]);
      await ethers.provider.send("evm_mine"); // Mine a block to update timestamp
      
      expect(await multisigTreasury.isTransactionExpired(0)).to.be.true;
    });
  });
});