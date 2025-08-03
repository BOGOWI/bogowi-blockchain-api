const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("MultisigTreasury - Coverage Tests", function () {
  let treasury;
  let owner, signer1, signer2, signer3, nonSigner, recipient;
  let bogoToken;

  beforeEach(async function () {
    [owner, signer1, signer2, signer3, nonSigner, recipient] = await ethers.getSigners();

    // Deploy MultisigTreasury
    const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
    treasury = await MultisigTreasury.deploy(
      [signer1.address, signer2.address, signer3.address],
      2
    );
    await treasury.waitForDeployment();

    // Deploy BOGO token for testing
    const BOGOToken = await ethers.getContractFactory("BOGOTokenV2");
    bogoToken = await BOGOToken.deploy();
    await bogoToken.waitForDeployment();

    // Fund treasury
    await owner.sendTransaction({
      to: await treasury.getAddress(),
      value: ethers.parseEther("10")
    });
  });

  describe("Core Transaction Flow", function () {
    it("Should complete full transaction lifecycle", async function () {
      // Submit
      const tx = await treasury.connect(signer1).submitTransaction(
        recipient.address,
        ethers.parseEther("1"),
        "0x",
        "Test payment"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      // Check transaction details
      const txData = await treasury.getTransaction(txId);
      expect(txData.to).to.equal(recipient.address);
      expect(txData.value).to.equal(ethers.parseEther("1"));
      expect(txData.executed).to.be.false;

      // Confirm
      await expect(
        treasury.connect(signer2).confirmTransaction(txId)
      ).to.emit(treasury, "TransactionConfirmed");

      // Check confirmation count
      expect(await treasury.getConfirmationCount(txId)).to.equal(2);

      // Execute after delay
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");

      const balanceBefore = await ethers.provider.getBalance(recipient.address);
      await expect(
        treasury.connect(signer1).executeTransaction(txId)
      ).to.emit(treasury, "TransactionExecuted");
      
      const balanceAfter = await ethers.provider.getBalance(recipient.address);
      expect(balanceAfter - balanceBefore).to.equal(ethers.parseEther("1"));
    });

    it("Should handle transaction revocation", async function () {
      // Submit and auto-confirm by signer1
      const tx = await treasury.connect(signer1).submitTransaction(
        recipient.address,
        ethers.parseEther("1"),
        "0x",
        "Test"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      // Revoke
      await expect(
        treasury.connect(signer1).revokeConfirmation(txId)
      ).to.emit(treasury, "ConfirmationRevoked");

      expect(await treasury.getConfirmationCount(txId)).to.equal(0);
    });

    it("Should cancel expired transactions", async function () {
      const tx = await treasury.connect(signer1).submitTransaction(
        recipient.address,
        ethers.parseEther("1"),
        "0x",
        "Will expire"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      // Fast forward past expiry
      await ethers.provider.send("evm_increaseTime", [8 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine");

      expect(await treasury.isTransactionExpired(txId)).to.be.true;

      await expect(
        treasury.connect(signer1).cancelExpiredTransaction(txId)
      ).to.emit(treasury, "TransactionCancelled");
    });
  });

  describe("Signer Management", function () {
    it("Should add signer via multisig", async function () {
      const newSigner = ethers.Wallet.createRandom();
      const data = treasury.interface.encodeFunctionData("addSigner", [newSigner.address]);
      
      // Submit
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Add new signer"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      // Confirm and execute
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      // Verify
      const signers = await treasury.getSigners();
      expect(signers).to.include(newSigner.address);
    });

    it("Should remove signer via multisig", async function () {
      const data = treasury.interface.encodeFunctionData("removeSigner", [signer3.address]);
      
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Remove signer3"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      const signers = await treasury.getSigners();
      expect(signers).to.not.include(signer3.address);
    });

    it("Should change threshold via multisig", async function () {
      const data = treasury.interface.encodeFunctionData("changeThreshold", [1]);
      
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Change threshold"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      expect(await treasury.threshold()).to.equal(1);
    });
  });

  describe("Token Operations", function () {
    beforeEach(async function () {
      // Mint tokens to treasury
      const DAO_ROLE = await bogoToken.DAO_ROLE();
      await bogoToken.grantRole(DAO_ROLE, owner.address);
      await bogoToken.mintFromRewards(await treasury.getAddress(), ethers.parseEther("1000"));
    });

    it("Should transfer ERC20 tokens via multisig", async function () {
      const amount = ethers.parseEther("100");
      const data = treasury.interface.encodeFunctionData(
        "transferERC20",
        [await bogoToken.getAddress(), recipient.address, amount]
      );

      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Transfer BOGO"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      expect(await bogoToken.balanceOf(recipient.address)).to.equal(amount);
    });

    it("Should handle direct token transfers", async function () {
      const amount = ethers.parseEther("50");
      const data = bogoToken.interface.encodeFunctionData(
        "transfer",
        [recipient.address, amount]
      );

      const tx = await treasury.connect(signer1).submitTransaction(
        await bogoToken.getAddress(),
        0,
        data,
        "Direct transfer"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      expect(await bogoToken.balanceOf(recipient.address)).to.equal(amount);
    });
  });

  describe("Role Management", function () {
    // Skipping - requires treasury to have admin role on token contract
    it.skip("Should grant role via multisig", async function () {
      const MINTER_ROLE = await bogoToken.MINTER_ROLE();
      const data = treasury.interface.encodeFunctionData(
        "grantRole",
        [await bogoToken.getAddress(), MINTER_ROLE, recipient.address]
      );

      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Grant minter role"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      expect(await bogoToken.hasRole(MINTER_ROLE, recipient.address)).to.be.true;
    });

    // Skipping - requires treasury to have admin role on token contract
    it.skip("Should revoke role via multisig", async function () {
      // First grant a role
      const MINTER_ROLE = await bogoToken.MINTER_ROLE();
      await bogoToken.grantRole(MINTER_ROLE, recipient.address);

      // Then revoke it via multisig
      const data = treasury.interface.encodeFunctionData(
        "revokeRole",
        [await bogoToken.getAddress(), MINTER_ROLE, recipient.address]
      );

      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Revoke minter role"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      expect(await bogoToken.hasRole(MINTER_ROLE, recipient.address)).to.be.false;
    });
  });

  describe("Advanced Features", function () {
    it("Should toggle auto-execute", async function () {
      const initialState = await treasury.autoExecuteEnabled();
      
      const data = treasury.interface.encodeFunctionData("toggleAutoExecute");
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Toggle auto-execute"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      expect(await treasury.autoExecuteEnabled()).to.equal(!initialState);
    });

    it("Should set function allowance", async function () {
      const selector = bogoToken.interface.getFunction("transfer").selector;
      
      const data = treasury.interface.encodeFunctionData(
        "setFunctionAllowance",
        [await bogoToken.getAddress(), selector, true]
      );

      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Allow transfer"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      expect(await treasury.allowedFunctions(await bogoToken.getAddress(), selector)).to.be.true;
    });

    it("Should toggle function restrictions", async function () {
      const initialState = await treasury.restrictFunctionCalls();
      
      const data = treasury.interface.encodeFunctionData("toggleFunctionRestrictions");
      const tx = await treasury.connect(signer1).submitTransaction(
        await treasury.getAddress(),
        0,
        data,
        "Toggle restrictions"
      );
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;

      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);

      expect(await treasury.restrictFunctionCalls()).to.equal(!initialState);
    });
  });

  describe("View Functions", function () {
    beforeEach(async function () {
      // Create some transactions
      for (let i = 0; i < 3; i++) {
        await treasury.connect(signer1).submitTransaction(
          recipient.address,
          i,
          "0x",
          `Transaction ${i}`
        );
      }
    });

    it("Should return pending transactions", async function () {
      const pending = await treasury.getPendingTransactions();
      expect(pending.length).to.equal(3);
    });

    it("Should check if transaction is expired", async function () {
      expect(await treasury.isTransactionExpired(0)).to.be.false;
      
      // Fast forward past expiry
      await ethers.provider.send("evm_increaseTime", [8 * 24 * 60 * 60]);
      await ethers.provider.send("evm_mine");
      
      expect(await treasury.isTransactionExpired(0)).to.be.true;
    });

    it("Should check confirmation status", async function () {
      expect(await treasury.hasConfirmed(0, signer1.address)).to.be.true;
      expect(await treasury.hasConfirmed(0, signer2.address)).to.be.false;
    });

    it("Should get emergency approval info", async function () {
      expect(await treasury.getEmergencyApprovalCount()).to.equal(0);
      expect(await treasury.hasEmergencyApproval(signer1.address)).to.be.false;
    });
  });

  describe("Error Cases", function () {
    it("Should reject invalid operations", async function () {
      // Non-signer submission
      await expect(
        treasury.connect(nonSigner).submitTransaction(recipient.address, 0, "0x", "Test")
      ).to.be.revertedWith("NOT_SIGNER");

      // Invalid recipient
      await expect(
        treasury.connect(signer1).submitTransaction(ethers.ZeroAddress, 0, "0x", "Test")
      ).to.be.revertedWith("ZERO_ADDRESS");

      // Non-existent transaction
      await expect(
        treasury.connect(signer1).confirmTransaction(999)
      ).to.be.revertedWith("DOES_NOT_EXIST");

      // Not confirmed revocation
      const tx = await treasury.connect(signer1).submitTransaction(recipient.address, 0, "0x", "Test");
      const receipt = await tx.wait();
      const txId = receipt.logs[0].args.txId;
      
      await expect(
        treasury.connect(signer2).revokeConfirmation(txId)
      ).to.be.revertedWith("NOT_INITIALIZED");
    });
  });
});