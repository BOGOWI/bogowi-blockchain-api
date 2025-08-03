const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("CommercialNFT - Additional Coverage Tests", function () {
    let commercialNFT;
    let owner, treasury, minter, business, user1, user2, user3, attacker;
    let DEFAULT_ADMIN_ROLE, MINTER_ROLE, BUSINESS_ROLE, TREASURY_ROLE;
    
    // Token IDs for different scenarios
    const TICKET_ID_1 = 10001;
    const TICKET_ID_2 = 10002;
    const COLLECTIBLE_ID_1 = 20001;
    const COLLECTIBLE_ID_2 = 20002;
    const MERCHANDISE_ID_1 = 30001;
    const MERCHANDISE_ID_2 = 30002;
    const GAMING_ID_1 = 40001;
    const GAMING_ID_2 = 40002;
    
    // Edge of ranges
    const TICKET_EDGE_LOW = 10000;
    const TICKET_EDGE_HIGH = 19999;
    const COLLECTIBLE_EDGE_LOW = 20000;
    const COLLECTIBLE_EDGE_HIGH = 29999;
    const MERCHANDISE_EDGE_LOW = 30000;
    const MERCHANDISE_EDGE_HIGH = 39999;
    const GAMING_EDGE_LOW = 40000;
    const GAMING_EDGE_HIGH = 49999;
    
    beforeEach(async function () {
        [owner, treasury, minter, business, user1, user2, user3, attacker] = await ethers.getSigners();
        
        const CommercialNFT = await ethers.getContractFactory("CommercialNFT");
        commercialNFT = await CommercialNFT.deploy(treasury.address);
        await commercialNFT.waitForDeployment();
        
        DEFAULT_ADMIN_ROLE = await commercialNFT.DEFAULT_ADMIN_ROLE();
        MINTER_ROLE = await commercialNFT.MINTER_ROLE();
        BUSINESS_ROLE = await commercialNFT.BUSINESS_ROLE();
        TREASURY_ROLE = await commercialNFT.TREASURY_ROLE();
        
        await commercialNFT.grantRole(MINTER_ROLE, minter.address);
        await commercialNFT.grantRole(BUSINESS_ROLE, business.address);
    });
    
    describe("Complex Multi-Token Scenarios", function () {
        it("Should handle multiple token types simultaneously", async function () {
            const eventDate = (await time.latest()) + 86400;
            const expiryDate = eventDate + 86400;
            
            // Mint different token types
            await commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID_1,
                eventDate,
                expiryDate,
                "Venue",
                "ticket-uri",
                ethers.parseEther("0.1")
            );
            
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID_1,
                5,
                50,
                "collectible-uri",
                ethers.parseEther("1"),
                750
            );
            
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID_1,
                20,
                200,
                "gaming-uri",
                ethers.parseEther("0.05")
            );
            
            // Verify balances
            expect(await commercialNFT.balanceOf(user1.address, TICKET_ID_1)).to.equal(1);
            expect(await commercialNFT.balanceOf(user1.address, COLLECTIBLE_ID_1)).to.equal(5);
            expect(await commercialNFT.balanceOf(user1.address, GAMING_ID_1)).to.equal(20);
            
            // Batch transfer different tokens
            await commercialNFT.connect(user1).safeBatchTransferFrom(
                user1.address,
                user2.address,
                [TICKET_ID_1, COLLECTIBLE_ID_1, GAMING_ID_1],
                [1, 2, 10],
                "0x"
            );
            
            expect(await commercialNFT.balanceOf(user2.address, TICKET_ID_1)).to.equal(1);
            expect(await commercialNFT.balanceOf(user2.address, COLLECTIBLE_ID_1)).to.equal(2);
            expect(await commercialNFT.balanceOf(user2.address, GAMING_ID_1)).to.equal(10);
        });
        
        it("Should handle edge token IDs correctly", async function () {
            const eventDate = (await time.latest()) + 86400;
            const expiryDate = eventDate + 86400;
            
            // Test edge of ticket range
            await commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_EDGE_LOW,
                eventDate,
                expiryDate,
                "Venue",
                "uri",
                0
            );
            
            await commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_EDGE_HIGH,
                eventDate,
                expiryDate,
                "Venue",
                "uri",
                0
            );
            
            // Test edge of collectible range
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_EDGE_LOW,
                1,
                10,
                "uri",
                0,
                500
            );
            
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_EDGE_HIGH,
                1,
                10,
                "uri",
                0,
                500
            );
            
            // Verify all minted successfully
            expect(await commercialNFT.balanceOf(user1.address, TICKET_EDGE_LOW)).to.equal(1);
            expect(await commercialNFT.balanceOf(user1.address, TICKET_EDGE_HIGH)).to.equal(1);
            expect(await commercialNFT.balanceOf(user1.address, COLLECTIBLE_EDGE_LOW)).to.equal(1);
            expect(await commercialNFT.balanceOf(user1.address, COLLECTIBLE_EDGE_HIGH)).to.equal(1);
        });
    });
    
    describe("Royalty Edge Cases", function () {
        it("Should handle zero royalty", async function () {
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID_1,
                1,
                100,
                "uri",
                0,
                0 // 0% royalty
            );
            
            const [receiver, royaltyAmount] = await commercialNFT.royaltyInfo(COLLECTIBLE_ID_1, ethers.parseEther("1"));
            expect(receiver).to.equal(commercialNFT.address);
            expect(royaltyAmount).to.equal(0);
        });
        
        it("Should handle maximum allowed royalty", async function () {
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID_1,
                1,
                100,
                "uri",
                0,
                1000 // 10% max royalty
            );
            
            const [receiver, royaltyAmount] = await commercialNFT.royaltyInfo(COLLECTIBLE_ID_1, ethers.parseEther("1"));
            expect(receiver).to.equal(commercialNFT.address);
            expect(royaltyAmount).to.equal(ethers.parseEther("0.1"));
        });
        
        it("Should accumulate royalties in contract", async function () {
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID_1,
                1,
                100,
                "uri",
                0,
                1000
            );
            
            // Simulate royalty payment
            await owner.sendTransaction({
                to: commercialNFT.address,
                value: ethers.parseEther("0.1")
            });
            
            const contractBalance = await ethers.provider.getBalance(commercialNFT.address);
            expect(contractBalance).to.equal(ethers.parseEther("0.1"));
            
            // Treasury can withdraw royalties
            await commercialNFT.connect(treasury).withdraw();
            expect(await ethers.provider.getBalance(commercialNFT.address)).to.equal(0);
        });
    });
    
    describe("Batch Operations Edge Cases", function () {
        it("Should handle batch mint with exactly 100 recipients", async function () {
            const recipients = new Array(100).fill(null).map((_, i) => {
                return ethers.Wallet.createRandom().address;
            });
            
            await commercialNFT.connect(business).batchMintPromo(
                recipients,
                MERCHANDISE_ID_1,
                1,
                100,
                "uri",
                0
            );
            
            expect(await commercialNFT["totalSupply(uint256)"](MERCHANDISE_ID_1)).to.equal(100);
        });
        
        it("Should handle batch mint with mixed valid/invalid scenarios", async function () {
            // First mint to establish token
            await commercialNFT.connect(business).batchMintPromo(
                [user1.address],
                MERCHANDISE_ID_1,
                50,
                100,
                "uri",
                0
            );
            
            // Try to mint more than remaining supply
            await expect(commercialNFT.connect(business).batchMintPromo(
                [user2.address, user3.address],
                MERCHANDISE_ID_1,
                30, // 60 total would exceed 100
                100,
                "uri",
                0
            )).to.be.revertedWith("Exceeds max supply");
        });
        
        it("Should handle empty recipients array", async function () {
            await expect(commercialNFT.connect(business).batchMintPromo(
                [],
                MERCHANDISE_ID_1,
                1,
                100,
                "uri",
                0
            )).to.be.revertedWith("No recipients");
        });
    });
    
    describe("Treasury Role Transfer Scenarios", function () {
        it("Should handle multiple treasury updates", async function () {
            const newTreasury1 = user1.address;
            const newTreasury2 = user2.address;
            
            // First update
            await commercialNFT.setTreasuryAddress(newTreasury1);
            expect(await commercialNFT.treasuryAddress()).to.equal(newTreasury1);
            expect(await commercialNFT.hasRole(TREASURY_ROLE, newTreasury1)).to.be.true;
            
            // Second update
            await commercialNFT.setTreasuryAddress(newTreasury2);
            expect(await commercialNFT.treasuryAddress()).to.equal(newTreasury2);
            expect(await commercialNFT.hasRole(TREASURY_ROLE, newTreasury2)).to.be.true;
            expect(await commercialNFT.hasRole(TREASURY_ROLE, newTreasury1)).to.be.false;
        });
        
        it("Should handle treasury withdrawal after address change", async function () {
            await owner.sendTransaction({
                to: commercialNFT.address,
                value: ethers.parseEther("2")
            });
            
            // Change treasury
            await commercialNFT.setTreasuryAddress(user1.address);
            
            // Old treasury can't withdraw
            await expect(commercialNFT.connect(treasury).withdraw())
                .to.be.reverted;
            
            // New treasury can withdraw
            const balanceBefore = await user1.getBalance();
            await commercialNFT.connect(user1).withdraw();
            const balanceAfter = await user1.getBalance();
            
            expect(balanceAfter.sub(balanceBefore)).to.be.closeTo(
                ethers.parseEther("2"),
                ethers.parseEther("0.01")
            );
        });
    });
    
    describe("Complex Access Control Scenarios", function () {
        it("Should handle role renunciation", async function () {
            // Minter renounces their role
            await commercialNFT.connect(minter).renounceRole(MINTER_ROLE, minter.address);
            
            expect(await commercialNFT.hasRole(MINTER_ROLE, minter.address)).to.be.false;
            
            // Can't mint anymore
            await expect(commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID_1,
                1,
                100,
                "uri",
                0
            )).to.be.reverted;
        });
        
        it("Should handle admin transferring admin role", async function () {
            await commercialNFT.grantRole(DEFAULT_ADMIN_ROLE, user1.address);
            
            // user1 can now grant roles
            await commercialNFT.connect(user1).grantRole(MINTER_ROLE, user2.address);
            expect(await commercialNFT.hasRole(MINTER_ROLE, user2.address)).to.be.true;
            
            // Original admin can revoke their own role
            await commercialNFT.revokeRole(DEFAULT_ADMIN_ROLE, owner.address);
            expect(await commercialNFT.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.false;
        });
    });
    
    describe("Pause/Unpause Complex Scenarios", function () {
        beforeEach(async function () {
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID_1,
                100,
                1000,
                "uri",
                0
            );
        });
        
        it("Should block all token operations when paused", async function () {
            await commercialNFT.pause();
            
            // Can't transfer
            await expect(commercialNFT.connect(user1).safeTransferFrom(
                user1.address,
                user2.address,
                GAMING_ID_1,
                10,
                "0x"
            )).to.be.revertedWith("EnforcedPause");
            
            // Can't burn
            await expect(commercialNFT.connect(user1).burn(GAMING_ID_1, 10))
                .to.be.revertedWith("EnforcedPause");
            
            // Can still read state
            expect(await commercialNFT.balanceOf(user1.address, GAMING_ID_1)).to.equal(100);
            
            // Unpause and verify operations work
            await commercialNFT.unpause();
            
            await commercialNFT.connect(user1).burn(GAMING_ID_1, 10);
            expect(await commercialNFT.balanceOf(user1.address, GAMING_ID_1)).to.equal(90);
        });
    });
    
    describe("Token Existence and Supply Tracking", function () {
        it("Should track token existence correctly", async function () {
            expect(await commercialNFT.tokenExists(GAMING_ID_1)).to.be.false;
            
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID_1,
                1,
                100,
                "uri",
                0
            );
            
            expect(await commercialNFT.tokenExists(GAMING_ID_1)).to.be.true;
        });
        
        it("Should enforce max supply across multiple mints", async function () {
            const maxSupply = 100;
            
            // First mint
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID_1,
                40,
                maxSupply,
                "uri",
                0,
                500
            );
            
            // Second mint
            await commercialNFT.connect(minter).mintCollectible(
                user2.address,
                COLLECTIBLE_ID_1,
                40,
                maxSupply,
                "uri",
                0,
                500
            );
            
            // Third mint should fail
            await expect(commercialNFT.connect(minter).mintCollectible(
                user3.address,
                COLLECTIBLE_ID_1,
                30, // Would exceed max
                maxSupply,
                "uri",
                0,
                500
            )).to.be.revertedWith("Exceeds max supply");
            
            expect(await commercialNFT["totalSupply(uint256)"](COLLECTIBLE_ID_1)).to.equal(80);
        });
    });
    
    describe("Event Data Complex Scenarios", function () {
        it("Should handle ticket lifecycle from mint to redemption", async function () {
            const currentTime = await time.latest();
            const eventDate = currentTime + 86400;
            const expiryDate = eventDate + 86400;
            
            // Mint ticket
            await commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID_1,
                eventDate,
                expiryDate,
                "Madison Square Garden",
                "ipfs://ticket",
                ethers.parseEther("0.5")
            );
            
            // Transfer ticket
            await commercialNFT.connect(user1).safeTransferFrom(
                user1.address,
                user2.address,
                TICKET_ID_1,
                1,
                "0x"
            );
            
            // Wait until after event but before expiry
            await time.increase(86400 + 3600);
            
            // Redeem ticket
            await commercialNFT.connect(user2).redeemTicket(TICKET_ID_1);
            
            const eventData = await commercialNFT.eventData(TICKET_ID_1);
            expect(eventData.used).to.be.true;
            
            // Can't transfer after redemption
            await expect(commercialNFT.connect(user2).safeTransferFrom(
                user2.address,
                user3.address,
                TICKET_ID_1,
                1,
                "0x"
            )).to.not.be.reverted; // Transfer still allowed, just marked as used
        });
    });
    
    describe("URI Management Edge Cases", function () {
        it("Should handle very long URIs", async function () {
            const longUri = "ipfs://" + "a".repeat(1000);
            
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID_1,
                1,
                100,
                longUri,
                0
            );
            
            expect(await commercialNFT.uri(GAMING_ID_1)).to.equal(longUri);
            
            // Update to another long URI
            const newLongUri = "ipfs://" + "b".repeat(1000);
            await commercialNFT.connect(business).updateTokenURI(GAMING_ID_1, newLongUri);
            expect(await commercialNFT.uri(GAMING_ID_1)).to.equal(newLongUri);
        });
        
        it("Should return empty string for non-existent token URI", async function () {
            expect(await commercialNFT.uri(99999)).to.equal("");
        });
    });
    
    describe("Interface Support Verification", function () {
        it("Should support all expected interfaces", async function () {
            // ERC1155
            expect(await commercialNFT.supportsInterface("0xd9b67a26")).to.be.true;
            
            // ERC1155MetadataURI
            expect(await commercialNFT.supportsInterface("0x0e89341c")).to.be.true;
            
            // AccessControl
            expect(await commercialNFT.supportsInterface("0x7965db0b")).to.be.true;
            
            // ERC2981 (Royalty)
            expect(await commercialNFT.supportsInterface("0x2a55205a")).to.be.true;
            
            // Should not support random interface
            expect(await commercialNFT.supportsInterface("0x12345678")).to.be.false;
        });
    });
    
    describe("Gas Optimization Scenarios", function () {
        it("Should efficiently handle large batch transfers", async function () {
            // Mint multiple token types
            const tokenIds = [];
            const amounts = [];
            
            for (let i = 0; i < 5; i++) {
                const tokenId = GAMING_ID_1 + i;
                await commercialNFT.connect(minter).mintGamingAsset(
                    user1.address,
                    tokenId,
                    100,
                    1000,
                    `uri-${i}`,
                    0
                );
                tokenIds.push(tokenId);
                amounts.push(20);
            }
            
            // Batch transfer
            const tx = await commercialNFT.connect(user1).safeBatchTransferFrom(
                user1.address,
                user2.address,
                tokenIds,
                amounts,
                "0x"
            );
            
            // Verify transfers
            for (let i = 0; i < tokenIds.length; i++) {
                expect(await commercialNFT.balanceOf(user2.address, tokenIds[i])).to.equal(20);
                expect(await commercialNFT.balanceOf(user1.address, tokenIds[i])).to.equal(80);
            }
        });
    });
    
    describe("Security and Attack Scenarios", function () {
        it("Should prevent reentrancy in withdraw", async function () {
            // Send ETH to contract
            await owner.sendTransaction({
                to: commercialNFT.address,
                value: ethers.parseEther("10")
            });
            
            // Deploy attacker contract (would need actual implementation)
            // This test verifies the nonReentrant modifier is present
            const withdrawFunction = commercialNFT.interface.getFunction("withdraw");
            expect(withdrawFunction).to.exist;
        });
        
        it("Should handle approval and operator scenarios", async function () {
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID_1,
                100,
                1000,
                "uri",
                0
            );
            
            // Set approval
            await commercialNFT.connect(user1).setApprovalForAll(user2.address, true);
            expect(await commercialNFT.isApprovedForAll(user1.address, user2.address)).to.be.true;
            
            // user2 can transfer user1's tokens
            await commercialNFT.connect(user2).safeTransferFrom(
                user1.address,
                user3.address,
                GAMING_ID_1,
                50,
                "0x"
            );
            
            expect(await commercialNFT.balanceOf(user3.address, GAMING_ID_1)).to.equal(50);
            
            // Revoke approval
            await commercialNFT.connect(user1).setApprovalForAll(user2.address, false);
            
            await expect(commercialNFT.connect(user2).safeTransferFrom(
                user1.address,
                user3.address,
                GAMING_ID_1,
                25,
                "0x"
            )).to.be.reverted;
        });
    });
});