const { expect } = require("chai");
const { ethers } = require("hardhat");

describe.skip("MultisigTreasury - Full Coverage", function () {
  let treasury;
  let owner, signer1, signer2, signer3, nonSigner;
  let nft721, nft1155;

  beforeEach(async function () {
    [owner, signer1, signer2, signer3, nonSigner] = await ethers.getSigners();

    // Deploy MultisigTreasury
    const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
    treasury = await MultisigTreasury.deploy(
      [signer1.address, signer2.address, signer3.address],
      2
    );
    await treasury.waitForDeployment();

    // Deploy NFT contracts
    const ConservationNFT = await ethers.getContractFactory("ConservationNFT");
    nft721 = await ConservationNFT.deploy(
      "Test NFT",
      "TNFT",
      "https://test.com/",
      treasury.address
    );
    await nft721.waitForDeployment();

    const CommercialNFT = await ethers.getContractFactory("CommercialNFT");
    nft1155 = await CommercialNFT.deploy(
      "Test 1155",
      "T1155",
      "https://test.com/",
      treasury.address
    );
    await nft1155.waitForDeployment();
  });

  describe("NFT Receiver Functions", function () {
    it("Should receive ERC721 tokens", async function () {
      // Mint NFT to owner
      await nft721.connect(owner).mint(owner.address, "tokenURI");
      
      // Transfer to treasury
      await expect(
        nft721.connect(owner)["safeTransferFrom(address,address,uint256)"](
          owner.address,
          treasury.address,
          1
        )
      ).to.emit(treasury, "TokensReceived")
        .withArgs(nft721.address, owner.address, 1);
      
      expect(await nft721.ownerOf(1)).to.equal(treasury.address);
    });

    it("Should transfer ERC721 via multisig", async function () {
      // Mint to treasury
      await nft721.connect(owner).mint(treasury.address, "tokenURI");
      
      // Transfer out via multisig
      const data = treasury.interface.encodeFunctionData(
        "transferERC721",
        [nft721.address, owner.address, 1]
      );
      
      const tx = await treasury.connect(signer1).submitTransaction(
        treasury.address,
        0,
        data,
        "Transfer NFT"
      );
      const receipt = await tx.wait();
      const txId = receipt.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);
      
      expect(await nft721.ownerOf(1)).to.equal(owner.address);
    });

    it("Should receive ERC1155 tokens", async function () {
      // Create token type
      await nft1155.createTokenType(100, "uri");
      
      // Mint to owner
      await nft1155.mint(owner.address, 1, 50, "0x");
      
      // Transfer to treasury
      await expect(
        nft1155.connect(owner).safeTransferFrom(
          owner.address,
          treasury.address,
          1,
          25,
          "0x"
        )
      ).to.emit(treasury, "TokensReceived")
        .withArgs(nft1155.address, owner.address, 25);
      
      expect(await nft1155.balanceOf(treasury.address, 1)).to.equal(25);
    });

    it("Should receive ERC1155 batch tokens", async function () {
      // Create token types
      await nft1155.createTokenType(100, "uri1");
      await nft1155.createTokenType(200, "uri2");
      
      // Mint to owner
      await nft1155.mint(owner.address, 1, 50, "0x");
      await nft1155.mint(owner.address, 2, 100, "0x");
      
      // Batch transfer to treasury
      await expect(
        nft1155.connect(owner).safeBatchTransferFrom(
          owner.address,
          treasury.address,
          [1, 2],
          [25, 50],
          "0x"
        )
      ).to.emit(treasury, "TokensReceived")
        .withArgs(nft1155.address, owner.address, 0);
      
      expect(await nft1155.balanceOf(treasury.address, 1)).to.equal(25);
      expect(await nft1155.balanceOf(treasury.address, 2)).to.equal(50);
    });

    it("Should transfer ERC1155 via multisig", async function () {
      // Create and mint to treasury
      await nft1155.createTokenType(100, "uri");
      await nft1155.mint(treasury.address, 1, 50, "0x");
      
      // Transfer out via multisig
      const data = treasury.interface.encodeFunctionData(
        "transferERC1155",
        [nft1155.address, owner.address, 1, 25, "0x"]
      );
      
      const tx = await treasury.connect(signer1).submitTransaction(
        treasury.address,
        0,
        data,
        "Transfer 1155"
      );
      const receipt = await tx.wait();
      const txId = receipt.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);
      
      expect(await nft1155.balanceOf(owner.address, 1)).to.equal(25);
    });
  });

  describe("Error Handling", function () {
    it("Should handle duplicate signers in constructor", async function () {
      const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
      await expect(
        MultisigTreasury.deploy(
          [signer1.address, signer1.address],
          1
        )
      ).to.be.revertedWith("Duplicate signer");
    });

    it("Should handle zero address in signers", async function () {
      const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
      await expect(
        MultisigTreasury.deploy(
          [ethers.ZeroAddress],
          1
        )
      ).to.be.revertedWith("Invalid signer");
    });

    it("Should handle adding duplicate signer", async function () {
      const data = treasury.interface.encodeFunctionData("addSigner", [signer1.address]);
      
      const tx = await treasury.connect(signer1).submitTransaction(
        treasury.address,
        0,
        data,
        "Add duplicate"
      );
      const receipt = await tx.wait();
      const txId = receipt.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await expect(
        treasury.connect(signer1).executeTransaction(txId)
      ).to.be.reverted;
    });

    it("Should handle removing non-signer", async function () {
      const data = treasury.interface.encodeFunctionData("removeSigner", [nonSigner.address]);
      
      const tx = await treasury.connect(signer1).submitTransaction(
        treasury.address,
        0,
        data,
        "Remove non-signer"
      );
      const receipt = await tx.wait();
      const txId = receipt.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await expect(
        treasury.connect(signer1).executeTransaction(txId)
      ).to.be.reverted;
    });

    it("Should handle function restrictions when enabled", async function () {
      // Enable function restrictions
      const toggleData = treasury.interface.encodeFunctionData("toggleFunctionRestrictions");
      const tx1 = await treasury.connect(signer1).submitTransaction(
        treasury.address,
        0,
        toggleData,
        "Enable restrictions"
      );
      const receipt1 = await tx1.wait();
      const txId1 = receipt1.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId1);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId1);
      
      expect(await treasury.functionRestrictions()).to.be.true;
      
      // Try to call a non-allowed function
      const transferData = "0xa9059cbb" + "0".repeat(64); // transfer function
      const tx2 = await treasury.connect(signer1).submitTransaction(
        owner.address,
        0,
        transferData,
        "Restricted call"
      );
      const receipt2 = await tx2.wait();
      const txId2 = receipt2.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId2);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await expect(
        treasury.connect(signer1).executeTransaction(txId2)
      ).to.be.revertedWith("Function not allowed");
    });

    it("Should handle auto-execute when enabled", async function () {
      // Enable auto-execute
      const toggleData = treasury.interface.encodeFunctionData("toggleAutoExecute");
      const tx1 = await treasury.connect(signer1).submitTransaction(
        treasury.address,
        0,
        toggleData,
        "Enable auto-execute"
      );
      const receipt1 = await tx1.wait();
      const txId1 = receipt1.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId1);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId1);
      
      expect(await treasury.autoExecute()).to.be.true;
      
      // Submit a new transaction - should auto-execute
      const balanceBefore = await ethers.provider.getBalance(owner.address);
      
      // Need to wait for delay even with auto-execute
      const tx2 = await treasury.connect(signer1).submitTransaction(
        owner.address,
        ethers.parseEther("0.1"),
        "0x",
        "Auto-execute test"
      );
      await treasury.connect(signer2).confirmTransaction(1);
      
      // Auto-execute should still respect delay
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      // Should execute automatically when someone interacts
      await treasury.connect(signer1).submitTransaction(
        owner.address,
        0,
        "0x",
        "Trigger auto"
      );
      
      const balanceAfter = await ethers.provider.getBalance(owner.address);
      expect(balanceAfter.sub(balanceBefore)).to.be.gte(ethers.parseEther("0.09")); // Account for gas
    });
  });

  describe("Access Control", function () {
    it("Should reject direct calls to onlyMultisig functions", async function () {
      await expect(
        treasury.connect(signer1).addSigner(nonSigner.address)
      ).to.be.revertedWith("Must be called through multisig");
      
      await expect(
        treasury.connect(signer1).removeSigner(signer2.address)
      ).to.be.revertedWith("Must be called through multisig");
      
      await expect(
        treasury.connect(signer1).changeThreshold(1)
      ).to.be.revertedWith("Must be called through multisig");
      
      await expect(
        treasury.connect(signer1).pause()
      ).to.be.revertedWith("Must be called through multisig");
      
      await expect(
        treasury.connect(signer1).unpause()
      ).to.be.revertedWith("Must be called through multisig");
      
      await expect(
        treasury.connect(signer1).toggleAutoExecute()
      ).to.be.revertedWith("Must be called through multisig");
      
      await expect(
        treasury.connect(signer1).toggleFunctionRestrictions()
      ).to.be.revertedWith("Must be called through multisig");
      
      await expect(
        treasury.connect(signer1).setFunctionAllowance(owner.address, "0x12345678", true)
      ).to.be.revertedWith("Must be called through multisig");
    });

    it("Should allow only signers to perform signer actions", async function () {
      await expect(
        treasury.connect(nonSigner).confirmTransaction(0)
      ).to.be.revertedWith("Not a signer");
      
      await expect(
        treasury.connect(nonSigner).revokeConfirmation(0)
      ).to.be.revertedWith("Not a signer");
      
      await expect(
        treasury.connect(nonSigner).executeTransaction(0)
      ).to.be.revertedWith("Not a signer");
      
      await expect(
        treasury.connect(nonSigner).cancelExpiredTransaction(0)
      ).to.be.revertedWith("Not a signer");
    });
  });

  describe("Threshold Edge Cases", function () {
    it("Should handle threshold equal to signer count", async function () {
      // Change threshold to 3 (equal to signer count)
      const data = treasury.interface.encodeFunctionData("changeThreshold", [3]);
      
      const tx = await treasury.connect(signer1).submitTransaction(
        treasury.address,
        0,
        data,
        "Max threshold"
      );
      const receipt = await tx.wait();
      const txId = receipt.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId);
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      await treasury.connect(signer1).executeTransaction(txId);
      
      expect(await treasury.threshold()).to.equal(3);
      
      // Now need all 3 signers for any transaction
      const tx2 = await treasury.connect(signer1).submitTransaction(
        owner.address,
        ethers.parseEther("0.1"),
        "0x",
        "Need all signers"
      );
      const receipt2 = await tx2.wait();
      const txId2 = receipt2.events[0].args.txId;
      
      await treasury.connect(signer2).confirmTransaction(txId2);
      
      // Still not enough
      await ethers.provider.send("evm_increaseTime", [3600]);
      await ethers.provider.send("evm_mine");
      
      await expect(
        treasury.connect(signer1).executeTransaction(txId2)
      ).to.be.revertedWith("Not enough confirmations");
      
      // Third confirmation allows execution
      await treasury.connect(signer3).confirmTransaction(txId2);
      await treasury.connect(signer1).executeTransaction(txId2);
      
      const tx2Data = await treasury.getTransaction(txId2);
      expect(tx2Data.executed).to.be.true;
    });
  });
});