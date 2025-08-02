const { expect } = require("chai");
const { ethers } = require("hardhat");
const { time } = require("@nomicfoundation/hardhat-network-helpers");

describe("CommercialNFT", function () {
    let commercialNFT;
    let owner, treasury, minter, business, user1, user2, user3;
    let DEFAULT_ADMIN_ROLE, MINTER_ROLE, BUSINESS_ROLE, TREASURY_ROLE;
    
    // Token type constants
    const EVENT_TICKET = 1;
    const COLLECTIBLE = 2;
    const MERCHANDISE = 3;
    const GAMING_ASSET = 4;
    
    // Token ID ranges
    const TICKET_ID = 10001;
    const COLLECTIBLE_ID = 20001;
    const MERCHANDISE_ID = 30001;
    const GAMING_ID = 40001;
    
    beforeEach(async function () {
        [owner, treasury, minter, business, user1, user2, user3] = await ethers.getSigners();
        
        const CommercialNFT = await ethers.getContractFactory("CommercialNFT");
        commercialNFT = await CommercialNFT.deploy(treasury.address);
        await commercialNFT.deployed();
        
        // Get role constants
        DEFAULT_ADMIN_ROLE = await commercialNFT.DEFAULT_ADMIN_ROLE();
        MINTER_ROLE = await commercialNFT.MINTER_ROLE();
        BUSINESS_ROLE = await commercialNFT.BUSINESS_ROLE();
        TREASURY_ROLE = await commercialNFT.TREASURY_ROLE();
        
        // Grant roles
        await commercialNFT.grantRole(MINTER_ROLE, minter.address);
        await commercialNFT.grantRole(BUSINESS_ROLE, business.address);
    });
    
    describe("Deployment and Initialization", function () {
        it("Should set the correct treasury address", async function () {
            expect(await commercialNFT.treasuryAddress()).to.equal(treasury.address);
        });
        
        it("Should grant treasury role to treasury address", async function () {
            expect(await commercialNFT.hasRole(TREASURY_ROLE, treasury.address)).to.be.true;
        });
        
        it("Should grant all roles to deployer", async function () {
            expect(await commercialNFT.hasRole(DEFAULT_ADMIN_ROLE, owner.address)).to.be.true;
            expect(await commercialNFT.hasRole(MINTER_ROLE, owner.address)).to.be.true;
            expect(await commercialNFT.hasRole(BUSINESS_ROLE, owner.address)).to.be.true;
        });
        
        it("Should set default royalty", async function () {
            // Mint a token and check royalty info
            await commercialNFT.mintGamingAsset(user1.address, GAMING_ID, 1, 100, "uri", 0);
            const [receiver, royaltyAmount] = await commercialNFT.royaltyInfo(GAMING_ID, 10000);
            expect(receiver).to.equal(commercialNFT.address);
            expect(royaltyAmount).to.equal(500); // 5% of 10000
        });
        
        it("Should fail deployment with zero treasury address", async function () {
            const CommercialNFT = await ethers.getContractFactory("CommercialNFT");
            await expect(CommercialNFT.deploy(ethers.constants.AddressZero))
                .to.be.revertedWith("Invalid treasury address");
        });
    });
    
    describe("Event Ticket Minting", function () {
        let eventDate, expiryDate;
        
        beforeEach(async function () {
            const currentTime = await time.latest();
            eventDate = currentTime + 86400; // 1 day from now
            expiryDate = eventDate + 86400; // 2 days from now
        });
        
        it("Should mint event ticket successfully", async function () {
            await expect(commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID,
                eventDate,
                expiryDate,
                "Madison Square Garden",
                "ipfs://ticket-uri",
                ethers.utils.parseEther("0.1")
            )).to.emit(commercialNFT, "CommercialNFTMinted")
                .withArgs(user1.address, TICKET_ID, EVENT_TICKET, ethers.utils.parseEther("0.1"));
            
            expect(await commercialNFT.balanceOf(user1.address, TICKET_ID)).to.equal(1);
            
            const eventData = await commercialNFT.eventData(TICKET_ID);
            expect(eventData.venue).to.equal("Madison Square Garden");
            expect(eventData.used).to.be.false;
        });
        
        it("Should fail with invalid ticket ID range", async function () {
            await expect(commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                5000, // Invalid ID
                eventDate,
                expiryDate,
                "Venue",
                "uri",
                0
            )).to.be.revertedWith("Invalid ticket token ID range");
        });
        
        it("Should fail with past event date", async function () {
            const pastDate = (await time.latest()) - 86400;
            await expect(commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID,
                pastDate,
                expiryDate,
                "Venue",
                "uri",
                0
            )).to.be.revertedWith("Event date must be in future");
        });
        
        it("Should fail with expiry before event", async function () {
            await expect(commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID,
                eventDate,
                eventDate - 1,
                "Venue",
                "uri",
                0
            )).to.be.revertedWith("Expiry must be after event date");
        });
        
        it("Should fail minting duplicate token ID", async function () {
            await commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID,
                eventDate,
                expiryDate,
                "Venue",
                "uri",
                0
            );
            
            await expect(commercialNFT.connect(minter).mintEventTicket(
                user2.address,
                TICKET_ID,
                eventDate,
                expiryDate,
                "Venue",
                "uri",
                0
            )).to.be.revertedWith("Token ID already exists");
        });
    });
    
    describe("Ticket Redemption", function () {
        let eventDate, expiryDate;
        
        beforeEach(async function () {
            const currentTime = await time.latest();
            eventDate = currentTime + 86400;
            expiryDate = eventDate + 86400;
            
            await commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID,
                eventDate,
                expiryDate,
                "Venue",
                "uri",
                0
            );
        });
        
        it("Should redeem ticket successfully", async function () {
            await expect(commercialNFT.connect(user1).redeemTicket(TICKET_ID))
                .to.emit(commercialNFT, "TicketRedeemed")
                .withArgs(TICKET_ID, user1.address);
            
            const eventData = await commercialNFT.eventData(TICKET_ID);
            expect(eventData.used).to.be.true;
        });
        
        it("Should fail redemption by non-holder", async function () {
            await expect(commercialNFT.connect(user2).redeemTicket(TICKET_ID))
                .to.be.revertedWith("Not ticket holder");
        });
        
        it("Should fail double redemption", async function () {
            await commercialNFT.connect(user1).redeemTicket(TICKET_ID);
            await expect(commercialNFT.connect(user1).redeemTicket(TICKET_ID))
                .to.be.revertedWith("Ticket already used");
        });
        
        it("Should fail redemption after expiry", async function () {
            await time.increase(86400 * 3); // 3 days
            await expect(commercialNFT.connect(user1).redeemTicket(TICKET_ID))
                .to.be.revertedWith("Ticket expired");
        });
        
        it("Should fail redemption for non-ticket token", async function () {
            await expect(commercialNFT.connect(user1).redeemTicket(COLLECTIBLE_ID))
                .to.be.revertedWith("Not ticket holder");
        });
    });
    
    describe("Collectible Minting", function () {
        it("Should mint collectible with custom royalty", async function () {
            await expect(commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID,
                10,
                100,
                "ipfs://collectible-uri",
                ethers.utils.parseEther("1"),
                750 // 7.5% royalty
            )).to.emit(commercialNFT, "CommercialNFTMinted")
                .withArgs(user1.address, COLLECTIBLE_ID, COLLECTIBLE, ethers.utils.parseEther("1"));
            
            expect(await commercialNFT.balanceOf(user1.address, COLLECTIBLE_ID)).to.equal(10);
            expect(await commercialNFT["totalSupply(uint256)"](COLLECTIBLE_ID)).to.equal(10);
            
            // Check royalty
            const [receiver, royaltyAmount] = await commercialNFT.royaltyInfo(COLLECTIBLE_ID, 10000);
            expect(receiver).to.equal(commercialNFT.address);
            expect(royaltyAmount).to.equal(750); // 7.5% of 10000
        });
        
        it("Should enforce max supply", async function () {
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID,
                90,
                100,
                "uri",
                0,
                500
            );
            
            await expect(commercialNFT.connect(minter).mintCollectible(
                user2.address,
                COLLECTIBLE_ID,
                20, // Would exceed max supply
                100,
                "uri",
                0,
                500
            )).to.be.revertedWith("Exceeds max supply");
        });
        
        it("Should fail with royalty exceeding maximum", async function () {
            await expect(commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID,
                1,
                100,
                "uri",
                0,
                1001 // >10%
            )).to.be.revertedWith("Royalty exceeds maximum");
        });
        
        it("Should fail with invalid collectible ID range", async function () {
            await expect(commercialNFT.connect(minter).mintCollectible(
                user1.address,
                TICKET_ID, // Wrong range
                1,
                100,
                "uri",
                0,
                500
            )).to.be.revertedWith("Invalid collectible token ID range");
        });
    });
    
    describe("Batch Minting for Promotions", function () {
        it("Should batch mint to multiple recipients", async function () {
            const recipients = [user1.address, user2.address, user3.address];
            
            await commercialNFT.connect(business).batchMintPromo(
                recipients,
                MERCHANDISE_ID,
                5,
                1000,
                "ipfs://promo-uri",
                0
            );
            
            for (const recipient of recipients) {
                expect(await commercialNFT.balanceOf(recipient, MERCHANDISE_ID)).to.equal(5);
            }
            expect(await commercialNFT["totalSupply(uint256)"](MERCHANDISE_ID)).to.equal(15);
        });
        
        it("Should fail with too many recipients", async function () {
            const recipients = new Array(101).fill(user1.address);
            
            await expect(commercialNFT.connect(business).batchMintPromo(
                recipients,
                MERCHANDISE_ID,
                1,
                1000,
                "uri",
                0
            )).to.be.revertedWith("Too many recipients");
        });
        
        it("Should fail with invalid recipient address", async function () {
            const recipients = [user1.address, ethers.constants.AddressZero];
            
            await expect(commercialNFT.connect(business).batchMintPromo(
                recipients,
                MERCHANDISE_ID,
                1,
                100,
                "uri",
                0
            )).to.be.revertedWith("Invalid recipient");
        });
        
        it("Should fail when exceeding max supply", async function () {
            const recipients = [user1.address, user2.address];
            
            await expect(commercialNFT.connect(business).batchMintPromo(
                recipients,
                MERCHANDISE_ID,
                60,
                100, // Max supply 100, trying to mint 120
                "uri",
                0
            )).to.be.revertedWith("Total mint exceeds max supply");
        });
        
        it("Should only be callable by BUSINESS_ROLE", async function () {
            await expect(commercialNFT.connect(user1).batchMintPromo(
                [user2.address],
                MERCHANDISE_ID,
                1,
                100,
                "uri",
                0
            )).to.be.reverted;
        });
    });
    
    describe("Gaming Asset Minting", function () {
        it("Should mint gaming assets", async function () {
            await expect(commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID,
                50,
                1000,
                "ipfs://gaming-uri",
                ethers.utils.parseEther("0.05")
            )).to.emit(commercialNFT, "CommercialNFTMinted")
                .withArgs(user1.address, GAMING_ID, GAMING_ASSET, ethers.utils.parseEther("0.05"));
            
            expect(await commercialNFT.balanceOf(user1.address, GAMING_ID)).to.equal(50);
            
            const tokenInfo = await commercialNFT.tokenInfo(GAMING_ID);
            expect(tokenInfo.burnable).to.be.true;
            expect(tokenInfo.tradeable).to.be.true;
        });
        
        it("Should fail with invalid gaming ID range", async function () {
            await expect(commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                MERCHANDISE_ID, // Wrong range
                1,
                100,
                "uri",
                0
            )).to.be.revertedWith("Invalid gaming token ID range");
        });
    });
    
    describe("Token Burning", function () {
        beforeEach(async function () {
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID,
                10,
                100,
                "uri",
                0
            );
        });
        
        it("Should burn burnable tokens", async function () {
            await commercialNFT.connect(user1).burn(GAMING_ID, 5);
            expect(await commercialNFT.balanceOf(user1.address, GAMING_ID)).to.equal(5);
            expect(await commercialNFT["totalSupply(uint256)"](GAMING_ID)).to.equal(5);
        });
        
        it("Should fail burning non-burnable tokens", async function () {
            // Mint a non-burnable collectible
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID,
                1,
                100,
                "uri",
                0,
                500
            );
            
            await expect(commercialNFT.connect(user1).burn(COLLECTIBLE_ID, 1))
                .to.be.revertedWith("Token not burnable");
        });
    });
    
    describe("Token Transfers and Tradeability", function () {
        beforeEach(async function () {
            // Mint tradeable token
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID,
                10,
                100,
                "uri",
                0
            );
        });
        
        it("Should transfer tradeable tokens", async function () {
            await commercialNFT.connect(user1).safeTransferFrom(
                user1.address,
                user2.address,
                GAMING_ID,
                5,
                "0x"
            );
            
            expect(await commercialNFT.balanceOf(user1.address, GAMING_ID)).to.equal(5);
            expect(await commercialNFT.balanceOf(user2.address, GAMING_ID)).to.equal(5);
        });
        
        it("Should fail transfer when paused", async function () {
            await commercialNFT.pause();
            
            await expect(commercialNFT.connect(user1).safeTransferFrom(
                user1.address,
                user2.address,
                GAMING_ID,
                5,
                "0x"
            )).to.be.revertedWith("EnforcedPause");
        });
    });
    
    describe("Token URI Management", function () {
        beforeEach(async function () {
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID,
                1,
                100,
                "ipfs://original-uri",
                0
            );
        });
        
        it("Should update token URI", async function () {
            await expect(commercialNFT.connect(business).updateTokenURI(GAMING_ID, "ipfs://new-uri"))
                .to.emit(commercialNFT, "TokenURIUpdated")
                .withArgs(GAMING_ID, "ipfs://new-uri");
            
            expect(await commercialNFT.uri(GAMING_ID)).to.equal("ipfs://new-uri");
        });
        
        it("Should fail updating non-existent token URI", async function () {
            await expect(commercialNFT.connect(business).updateTokenURI(99999, "uri"))
                .to.be.revertedWith("Token does not exist");
        });
        
        it("Should fail updating with empty URI", async function () {
            await expect(commercialNFT.connect(business).updateTokenURI(GAMING_ID, ""))
                .to.be.revertedWith("URI cannot be empty");
        });
        
        it("Should only be updatable by BUSINESS_ROLE", async function () {
            await expect(commercialNFT.connect(user1).updateTokenURI(GAMING_ID, "uri"))
                .to.be.reverted;
        });
    });
    
    describe("Treasury Management", function () {
        beforeEach(async function () {
            // Send some ETH to the contract
            await owner.sendTransaction({
                to: commercialNFT.address,
                value: ethers.utils.parseEther("5")
            });
        });
        
        it("Should withdraw funds to treasury", async function () {
            const initialBalance = await treasury.getBalance();
            
            await expect(commercialNFT.connect(treasury).withdraw())
                .to.emit(commercialNFT, "FundsWithdrawn")
                .withArgs(treasury.address, ethers.utils.parseEther("5"));
            
            const finalBalance = await treasury.getBalance();
            expect(finalBalance.sub(initialBalance)).to.be.closeTo(
                ethers.utils.parseEther("5"),
                ethers.utils.parseEther("0.01") // Account for gas
            );
        });
        
        it("Should fail withdrawal with no funds", async function () {
            await commercialNFT.connect(treasury).withdraw();
            
            await expect(commercialNFT.connect(treasury).withdraw())
                .to.be.revertedWith("No funds to withdraw");
        });
        
        it("Should only allow TREASURY_ROLE to withdraw", async function () {
            await expect(commercialNFT.connect(user1).withdraw())
                .to.be.reverted;
        });
        
        it("Should update treasury address", async function () {
            await expect(commercialNFT.setTreasuryAddress(user3.address))
                .to.emit(commercialNFT, "TreasuryAddressUpdated")
                .withArgs(treasury.address, user3.address);
            
            expect(await commercialNFT.treasuryAddress()).to.equal(user3.address);
            expect(await commercialNFT.hasRole(TREASURY_ROLE, user3.address)).to.be.true;
            expect(await commercialNFT.hasRole(TREASURY_ROLE, treasury.address)).to.be.false;
            
            // New treasury should be able to withdraw
            await expect(commercialNFT.connect(user3).withdraw())
                .to.emit(commercialNFT, "FundsWithdrawn");
        });
        
        it("Should fail setting zero treasury address", async function () {
            await expect(commercialNFT.setTreasuryAddress(ethers.constants.AddressZero))
                .to.be.revertedWith("Invalid treasury address");
        });
        
        it("Should only allow admin to update treasury", async function () {
            await expect(commercialNFT.connect(user1).setTreasuryAddress(user3.address))
                .to.be.reverted;
        });
    });
    
    describe("Access Control", function () {
        it("Should grant and revoke roles correctly", async function () {
            expect(await commercialNFT.hasRole(MINTER_ROLE, user1.address)).to.be.false;
            
            await commercialNFT.grantRole(MINTER_ROLE, user1.address);
            expect(await commercialNFT.hasRole(MINTER_ROLE, user1.address)).to.be.true;
            
            await commercialNFT.revokeRole(MINTER_ROLE, user1.address);
            expect(await commercialNFT.hasRole(MINTER_ROLE, user1.address)).to.be.false;
        });
        
        it("Should enforce role requirements", async function () {
            // Try minting without MINTER_ROLE
            await expect(commercialNFT.connect(user1).mintGamingAsset(
                user2.address,
                GAMING_ID,
                1,
                100,
                "uri",
                0
            )).to.be.reverted;
            
            // Try pausing without ADMIN_ROLE
            await expect(commercialNFT.connect(user1).pause())
                .to.be.reverted;
        });
    });
    
    describe("Pause Functionality", function () {
        it("Should pause and unpause contract", async function () {
            await commercialNFT.pause();
            expect(await commercialNFT.paused()).to.be.true;
            
            await commercialNFT.unpause();
            expect(await commercialNFT.paused()).to.be.false;
        });
        
        it("Should block transfers when paused", async function () {
            await commercialNFT.connect(minter).mintGamingAsset(
                user1.address,
                GAMING_ID,
                10,
                100,
                "uri",
                0
            );
            
            await commercialNFT.pause();
            
            await expect(commercialNFT.connect(user1).safeTransferFrom(
                user1.address,
                user2.address,
                GAMING_ID,
                5,
                "0x"
            )).to.be.revertedWith("EnforcedPause");
        });
    });
    
    describe("ERC2981 Royalty Support", function () {
        it("Should support ERC2981 interface", async function () {
            expect(await commercialNFT.supportsInterface("0x2a55205a")).to.be.true; // ERC2981
        });
        
        it("Should return correct royalty info", async function () {
            await commercialNFT.connect(minter).mintCollectible(
                user1.address,
                COLLECTIBLE_ID,
                1,
                100,
                "uri",
                0,
                1000 // 10% royalty
            );
            
            const [receiver, royaltyAmount] = await commercialNFT.royaltyInfo(COLLECTIBLE_ID, ethers.utils.parseEther("1"));
            expect(receiver).to.equal(commercialNFT.address);
            expect(royaltyAmount).to.equal(ethers.utils.parseEther("0.1")); // 10%
        });
    });
    
    describe("Edge Cases and Error Handling", function () {
        it("Should handle zero address checks", async function () {
            await expect(commercialNFT.connect(minter).mintEventTicket(
                ethers.constants.AddressZero,
                TICKET_ID,
                Math.floor(Date.now() / 1000) + 86400,
                Math.floor(Date.now() / 1000) + 172800,
                "Venue",
                "uri",
                0
            )).to.be.revertedWith("Invalid recipient address");
        });
        
        it("Should handle empty string validations", async function () {
            const currentTime = await time.latest();
            const eventDate = currentTime + 86400;
            const expiryDate = eventDate + 86400;
            
            await expect(commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID,
                eventDate,
                expiryDate,
                "", // Empty venue
                "uri",
                0
            )).to.be.revertedWith("Venue cannot be empty");
            
            await expect(commercialNFT.connect(minter).mintEventTicket(
                user1.address,
                TICKET_ID,
                eventDate,
                expiryDate,
                "Venue",
                "", // Empty URI
                0
            )).to.be.revertedWith("URI cannot be empty");
        });
        
        it("Should handle reentrancy protection", async function () {
            // Withdraw is protected by nonReentrant modifier
            // This is enforced at the contract level
            // Proper reentrancy testing would require a malicious contract
            
            // Send ETH to the contract first
            await owner.sendTransaction({
                to: commercialNFT.address,
                value: ethers.utils.parseEther("1")
            });
            
            // Normal withdrawal should work
            await expect(commercialNFT.connect(treasury).withdraw())
                .to.not.be.reverted;
        });
        
        it("Should accept ETH via receive function", async function () {
            const initialBalance = await ethers.provider.getBalance(commercialNFT.address);
            
            await owner.sendTransaction({
                to: commercialNFT.address,
                value: ethers.utils.parseEther("1")
            });
            
            const finalBalance = await ethers.provider.getBalance(commercialNFT.address);
            expect(finalBalance.sub(initialBalance)).to.equal(ethers.utils.parseEther("1"));
        });
    });
});