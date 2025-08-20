const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture, time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGOWITickets - BurnOnRedeem Coverage", function () {
    const ONE_DAY = 24 * 60 * 60;
    const ONE_WEEK = 7 * ONE_DAY;
    
    async function deployTicketsFixture() {
        const [owner, admin, minter, pauser, conservationDAO, user1, user2, backend] = 
            await ethers.getSigners();
        
        // Deploy RoleManager
        const RoleManager = await ethers.getContractFactory("RoleManager");
        const roleManager = await RoleManager.deploy();
        await roleManager.waitForDeployment();
        
        // Setup roles
        const NFT_MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
        const PAUSER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("PAUSER_ROLE"));
        const ADMIN_ROLE = ethers.keccak256(ethers.toUtf8Bytes("ADMIN_ROLE"));
        const DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
        
        // Grant roles
        await roleManager.grantRole(DEFAULT_ADMIN_ROLE, admin.address);
        await roleManager.grantRole(NFT_MINTER_ROLE, minter.address);
        await roleManager.grantRole(PAUSER_ROLE, pauser.address);
        await roleManager.grantRole(ADMIN_ROLE, admin.address);
        await roleManager.grantRole(ethers.keccak256(ethers.toUtf8Bytes("BACKEND_ROLE")), backend.address);
        
        // Deploy BOGOWITickets
        const BOGOWITickets = await ethers.getContractFactory("BOGOWITickets");
        const tickets = await BOGOWITickets.deploy(
            await roleManager.getAddress(),
            conservationDAO.address
        );
        await tickets.waitForDeployment();
        
        // Register tickets contract with RoleManager
        await roleManager.registerContract(await tickets.getAddress(), "BOGOWITickets");
        
        return {
            tickets,
            roleManager,
            owner,
            admin,
            minter,
            pauser,
            conservationDAO,
            user1,
            user2,
            backend,
            NFT_MINTER_ROLE,
            PAUSER_ROLE,
            ADMIN_ROLE
        };
    }
    
    // Helper to create EIP-712 signature for redemption
    async function createRedemptionSignature(tickets, tokenId, redeemer, nonce, deadline, signer) {
        const chainId = 501; // Use the chainId that will be passed in redemptionData
        const domain = {
            name: "BOGOWITickets",
            version: "1",
            chainId: chainId,
            verifyingContract: await tickets.getAddress()
        };
        
        const types = {
            RedeemTicket: [
                { name: "tokenId", type: "uint256" },
                { name: "redeemer", type: "address" },
                { name: "nonce", type: "uint256" },
                { name: "deadline", type: "uint256" },
                { name: "chainId", type: "uint256" }
            ]
        };
        
        const value = {
            tokenId: tokenId,
            redeemer: redeemer.address,
            nonce: nonce,
            deadline: deadline,
            chainId: chainId
        };
        
        const signature = await signer.signTypedData(domain, types, value);
        return signature;
    }
    
    describe("Burn On Redeem Feature", function () {
        
        it("Should burn ticket when burnOnRedeem flag is set", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            const currentTime = await time.latest();
            const params = {
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("BURN_TEST")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT001")),
                utilityFlags: 0x02, // Bit 1 set = burn on redeem
                transferUnlockAt: currentTime + ONE_DAY,
                expiresAt: currentTime + ONE_WEEK,
                metadataURI: "ipfs://QmBurn",
                rewardBasisPoints: 100
            };
            
            await tickets.connect(minter).mintTicket(params);
            const tokenId = 10001;
            
            // Verify ticket exists
            expect(await tickets.ownerOf(tokenId)).to.equal(user1.address);
            
            // Create redemption signature
            const nonce = 1001;
            const deadline = currentTime + 3600;
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            // Redeem - should burn the ticket
            await expect(
                tickets.connect(user1).redeemTicket({
                    tokenId: tokenId,
                    redeemer: user1.address,
                    nonce: nonce,
                    deadline: deadline,
                    chainId: 501,
                    signature: signature
                })
            ).to.emit(tickets, "TicketRedeemed")
                .withArgs(tokenId, user1.address, await time.latest() + 1);
            
            // Verify ticket was burned
            await expect(tickets.ownerOf(tokenId))
                .to.be.revertedWithCustomError(tickets, "ERC721NonexistentToken");
        });
        
        it("Should NOT burn ticket when burnOnRedeem flag is not set", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            const currentTime = await time.latest();
            const params = {
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("NO_BURN_TEST")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT002")),
                utilityFlags: 0x00, // No flags set = don't burn on redeem
                transferUnlockAt: currentTime + ONE_DAY,
                expiresAt: currentTime + ONE_WEEK,
                metadataURI: "ipfs://QmNoBurn",
                rewardBasisPoints: 100
            };
            
            await tickets.connect(minter).mintTicket(params);
            const tokenId = 10001;
            
            // Create redemption signature
            const nonce = 2001;
            const deadline = currentTime + 3600;
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            // Redeem - should NOT burn the ticket
            await tickets.connect(user1).redeemTicket({
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                chainId: 501,
                signature: signature
            });
            
            // Verify ticket still exists but is redeemed
            expect(await tickets.ownerOf(tokenId)).to.equal(user1.address);
            expect(await tickets.isRedeemed(tokenId)).to.be.true;
        });
        
        it("Should allow transfer after redeem when flag is set", async function () {
            const { tickets, minter, backend, user1, user2 } = await loadFixture(deployTicketsFixture);
            
            const currentTime = await time.latest();
            const params = {
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("TRANSFER_AFTER_REDEEM")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT003")),
                utilityFlags: 0x01, // Bit 0 set = allow transfer after redeem
                transferUnlockAt: currentTime - 1000, // Already unlocked
                expiresAt: currentTime + ONE_WEEK,
                metadataURI: "",
                rewardBasisPoints: 0
            };
            
            await tickets.connect(minter).mintTicket(params);
            const tokenId = 10001;
            
            // Redeem the ticket
            const nonce = 3001;
            const deadline = currentTime + 3600;
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            await tickets.connect(user1).redeemTicket({
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                chainId: 501,
                signature: signature
            });
            
            // Should be able to transfer after redemption
            await tickets.connect(user1).transferFrom(user1.address, user2.address, tokenId);
            expect(await tickets.ownerOf(tokenId)).to.equal(user2.address);
        });
        
        it("Should test combined flags (burn on redeem + allow transfer)", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            const currentTime = await time.latest();
            const params = {
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("COMBINED_FLAGS")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT004")),
                utilityFlags: 0x03, // Both flags set
                transferUnlockAt: currentTime + ONE_DAY,
                expiresAt: currentTime + ONE_WEEK,
                metadataURI: "ipfs://QmCombined",
                rewardBasisPoints: 10000 // Max rewards
            };
            
            await tickets.connect(minter).mintTicket(params);
            const tokenId = 10001;
            
            // Get ticket data to verify flags were interpreted correctly
            const ticketData = await tickets.getTicketData(tokenId);
            expect(ticketData.utilityFlags).to.equal(0x03);
            
            // When bit 0 is set, nonTransferableAfterRedeem should be false
            // When bit 1 is set, burnOnRedeem should be true
            expect(ticketData.nonTransferableAfterRedeem).to.be.false;
            expect(ticketData.burnOnRedeem).to.be.true;
            
            // Redeem - should burn due to bit 1
            const nonce = 4001;
            const deadline = currentTime + 3600;
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            await tickets.connect(user1).redeemTicket({
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                chainId: 501,
                signature: signature
            });
            
            // Token should be burned
            await expect(tickets.ownerOf(tokenId))
                .to.be.revertedWithCustomError(tickets, "ERC721NonexistentToken");
        });
    });
});