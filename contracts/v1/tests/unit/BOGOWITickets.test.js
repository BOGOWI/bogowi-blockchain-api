const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture, time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGOWITickets", function () {
    // Constants for testing
    const DEFAULT_ROYALTY_BPS = 500; // 5%
    const ONE_DAY = 24 * 60 * 60;
    const ONE_WEEK = 7 * ONE_DAY;
    
    // Test fixture for deployment
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
        await roleManager.grantRole(ADMIN_ROLE, backend.address); // Backend can sign redemptions
        
        // Deploy BOGOWITickets
        const BOGOWITickets = await ethers.getContractFactory("BOGOWITickets");
        const tickets = await BOGOWITickets.deploy(
            await roleManager.getAddress(),
            conservationDAO.address
        );
        await tickets.waitForDeployment();
        
        // Register tickets contract with RoleManager
        await roleManager.registerContract(await tickets.getAddress(), "BOGOWITickets");
        
        // Deploy NFTRegistry for integration tests
        const NFTRegistry = await ethers.getContractFactory("NFTRegistry");
        const registry = await NFTRegistry.deploy(await roleManager.getAddress());
        await registry.waitForDeployment();
        await roleManager.registerContract(await registry.getAddress(), "NFTRegistry");
        
        // Grant registry roles
        const REGISTRY_ADMIN_ROLE = ethers.keccak256(ethers.toUtf8Bytes("REGISTRY_ADMIN_ROLE"));
        const CONTRACT_DEPLOYER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("CONTRACT_DEPLOYER_ROLE"));
        await roleManager.grantRole(REGISTRY_ADMIN_ROLE, admin.address);
        await roleManager.grantRole(CONTRACT_DEPLOYER_ROLE, admin.address);
        
        return {
            tickets,
            registry,
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
            ADMIN_ROLE,
            DEFAULT_ADMIN_ROLE
        };
    }
    
    // Helper function to create valid mint params
    async function createMintParams(to, eventId = "EVENT001", bookingId = null) {
        const now = await time.latest();
        const transferUnlock = now + ONE_DAY; // Unlock after 1 day
        const expiry = now + ONE_WEEK; // Expire after 1 week
        
        return {
            to: to,
            bookingId: bookingId || ethers.keccak256(ethers.toUtf8Bytes(Math.random().toString())),
            eventId: ethers.keccak256(ethers.toUtf8Bytes(eventId)),
            utilityFlags: 0,
            transferUnlockAt: transferUnlock,
            expiresAt: expiry,
            metadataURI: "ipfs://QmTest123",
            rewardBasisPoints: 100 // 1% BOGO rewards
        };
    }
    
    // Helper to create EIP-712 signature for redemption
    async function createRedemptionSignature(tickets, tokenId, redeemer, nonce, deadline, signer) {
        const chainId = await signer.provider.getNetwork().then(n => n.chainId);
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
    
    describe("Deployment", function () {
        it("Should deploy with correct parameters", async function () {
            const { tickets, conservationDAO } = await loadFixture(deployTicketsFixture);
            
            expect(await tickets.name()).to.equal("BOGOWI Tickets");
            expect(await tickets.symbol()).to.equal("BWTIX");
            expect(await tickets.conservationDAO()).to.equal(conservationDAO.address);
        });
        
        it("Should set default royalty correctly", async function () {
            const { tickets, conservationDAO } = await loadFixture(deployTicketsFixture);
            
            // Mint a ticket to test royalty
            const { minter, user1 } = await loadFixture(deployTicketsFixture);
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            const tokenId = 10001;
            const salePrice = ethers.parseEther("1");
            const [receiver, royaltyAmount] = await tickets.royaltyInfo(tokenId, salePrice);
            
            expect(receiver).to.equal(conservationDAO.address);
            expect(royaltyAmount).to.equal(salePrice * BigInt(DEFAULT_ROYALTY_BPS) / BigInt(10000));
        });
    });
    
    describe("Minting", function () {
        it("Should mint a ticket with valid parameters", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            
            await expect(tickets.connect(minter).mintTicket(params))
                .to.emit(tickets, "TicketMinted")
                .withArgs(
                    10001, // First token ID
                    params.bookingId,
                    params.eventId,
                    user1.address,
                    params.rewardBasisPoints
                );
            
            expect(await tickets.ownerOf(10001)).to.equal(user1.address);
            expect(await tickets.tokenURI(10001)).to.equal(params.metadataURI);
        });
        
        it("Should prevent minting with duplicate booking ID", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployTicketsFixture);
            
            const bookingId = ethers.keccak256(ethers.toUtf8Bytes("BOOKING001"));
            const params1 = await createMintParams(user1.address, "EVENT001", bookingId);
            const params2 = await createMintParams(user1.address, "EVENT002", bookingId);
            
            await tickets.connect(minter).mintTicket(params1);
            
            await expect(tickets.connect(minter).mintTicket(params2))
                .to.be.revertedWith("Booking ID already used");
        });
        
        it("Should prevent minting with invalid parameters", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployTicketsFixture);
            
            // Zero address
            let params = await createMintParams(ethers.ZeroAddress);
            await expect(tickets.connect(minter).mintTicket(params))
                .to.be.revertedWith("Cannot mint to zero address");
            
            // Expired timestamp
            params = await createMintParams(user1.address);
            params.expiresAt = (await time.latest()) - 1000;
            await expect(tickets.connect(minter).mintTicket(params))
                .to.be.revertedWith("Expiry must be in future");
            
            // Unlock after expiry
            params = await createMintParams(user1.address);
            params.transferUnlockAt = params.expiresAt + 1000;
            await expect(tickets.connect(minter).mintTicket(params))
                .to.be.revertedWith("Unlock must be before expiry");
        });
        
        it("Should only allow NFT_MINTER_ROLE to mint", async function () {
            const { tickets, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            
            await expect(tickets.connect(user1).mintTicket(params))
                .to.be.revertedWithCustomError(tickets, "UnauthorizedRole");
        });
        
        it("Should mint batch of tickets", async function () {
            const { tickets, minter, user1, user2 } = await loadFixture(deployTicketsFixture);
            
            const params = [
                await createMintParams(user1.address, "EVENT001"),
                await createMintParams(user2.address, "EVENT002"),
                await createMintParams(user1.address, "EVENT003")
            ];
            
            const tx = await tickets.connect(minter).mintBatch(params);
            const receipt = await tx.wait();
            
            // Check all tickets were minted
            expect(await tickets.ownerOf(10001)).to.equal(user1.address);
            expect(await tickets.ownerOf(10002)).to.equal(user2.address);
            expect(await tickets.ownerOf(10003)).to.equal(user1.address);
        });
    });
    
    describe("Transfer Rules", function () {
        it("Should prevent transfer before unlock time", async function () {
            const { tickets, minter, user1, user2 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            await expect(
                tickets.connect(user1).transferFrom(user1.address, user2.address, 10001)
            ).to.be.revertedWith("Transfer locked until unlock time");
        });
        
        it("Should allow transfer after unlock time", async function () {
            const { tickets, minter, user1, user2 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            params.transferUnlockAt = (await time.latest()) - 1000; // Already unlocked
            await tickets.connect(minter).mintTicket(params);
            
            await tickets.connect(user1).transferFrom(user1.address, user2.address, 10001);
            expect(await tickets.ownerOf(10001)).to.equal(user2.address);
        });
        
        it("Should prevent transfer of expired ticket", async function () {
            const { tickets, minter, user1, user2 } = await loadFixture(deployTicketsFixture);
            
            const currentTime = await time.latest();
            const params = await createMintParams(user1.address);
            params.transferUnlockAt = currentTime - 2000;
            params.expiresAt = currentTime + 10; // Expires in 10 seconds
            await tickets.connect(minter).mintTicket(params);
            
            // Wait for expiry
            await time.increase(11);
            
            await expect(
                tickets.connect(user1).transferFrom(user1.address, user2.address, 10001)
            ).to.be.revertedWith("Cannot transfer expired ticket");
        });
        
        it("Should update transfer unlock time", async function () {
            const { tickets, minter, admin, user1, user2 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            const newUnlockTime = (await time.latest()) - 1000;
            
            await expect(tickets.connect(admin).updateTransferUnlock(10001, newUnlockTime))
                .to.emit(tickets, "TransferUnlockUpdated")
                .withArgs(10001, newUnlockTime);
            
            // Should now be transferable
            await tickets.connect(user1).transferFrom(user1.address, user2.address, 10001);
            expect(await tickets.ownerOf(10001)).to.equal(user2.address);
        });
    });
    
    describe("Redemption", function () {
        it("Should redeem ticket with valid signature", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            // Mint ticket
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            // Create redemption signature
            const tokenId = 10001;
            const nonce = 1;
            const deadline = (await time.latest()) + 3600; // 1 hour from now
            
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            const redemptionData = {
                chainId: 501,
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                signature: signature
            };
            
            await expect(tickets.connect(user1).redeemTicket(redemptionData))
                .to.emit(tickets, "TicketRedeemed")
                .withArgs(tokenId, user1.address, await time.latest() + 1);
            
            // Check ticket is redeemed
            expect(await tickets.isRedeemed(tokenId)).to.be.true;
        });
        
        it("Should prevent double redemption", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            const tokenId = 10001;
            const nonce = 1;
            const deadline = (await time.latest()) + 3600;
            
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            const redemptionData = {
                chainId: 501,
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                signature: signature
            };
            
            await tickets.connect(user1).redeemTicket(redemptionData);
            
            // Try to redeem again
            await expect(tickets.connect(user1).redeemTicket(redemptionData))
                .to.be.revertedWith("Ticket not redeemable");
        });
        
        it("Should prevent nonce reuse", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            // Mint two tickets
            await tickets.connect(minter).mintTicket(await createMintParams(user1.address));
            await tickets.connect(minter).mintTicket(await createMintParams(user1.address));
            
            const nonce = 1;
            const deadline = (await time.latest()) + 3600;
            
            // Redeem first ticket with nonce 1
            const signature1 = await createRedemptionSignature(
                tickets,
                10001,
                user1,
                nonce,
                deadline,
                backend
            );
            
            await tickets.connect(user1).redeemTicket({
                tokenId: 10001,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                chainId: 501,
                signature: signature1
            });
            
            // Try to use same nonce for second ticket
            const signature2 = await createRedemptionSignature(
                tickets,
                10002,
                user1,
                nonce,
                deadline,
                backend
            );
            
            await expect(tickets.connect(user1).redeemTicket({
                tokenId: 10002,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                chainId: 501,
                signature: signature2
            })).to.be.revertedWith("Nonce already used");
        });
        
        it("Should reject invalid signature", async function () {
            const { tickets, minter, user1, user2 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            const tokenId = 10001;
            const nonce = 1;
            const deadline = (await time.latest()) + 3600;
            
            // Create signature from unauthorized signer (user2)
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                user2 // Wrong signer
            );
            
            const redemptionData = {
                chainId: 501,
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                signature: signature
            };
            
            await expect(tickets.connect(user1).redeemTicket(redemptionData))
                .to.be.revertedWith("Invalid signature");
        });
        
        it("Should reject expired signature", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            const tokenId = 10001;
            const nonce = 1;
            const deadline = (await time.latest()) - 1000; // Already expired
            
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            const redemptionData = {
                chainId: 501,
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                signature: signature
            };
            
            await expect(tickets.connect(user1).redeemTicket(redemptionData))
                .to.be.revertedWith("Signature expired");
        });
        
        it("Should prevent transfer after redemption", async function () {
            const { tickets, minter, backend, user1, user2 } = await loadFixture(deployTicketsFixture);
            
            // Mint with immediate unlock
            const params = await createMintParams(user1.address);
            params.transferUnlockAt = (await time.latest()) - 1000;
            await tickets.connect(minter).mintTicket(params);
            
            // Redeem ticket
            const tokenId = 10001;
            const nonce = 1;
            const deadline = (await time.latest()) + 3600;
            
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
            
            // Try to transfer redeemed ticket
            await expect(
                tickets.connect(user1).transferFrom(user1.address, user2.address, tokenId)
            ).to.be.revertedWith("Cannot transfer redeemed ticket");
        });
    });
    
    describe("Expiry", function () {
        it("Should mark ticket as expired", async function () {
            const { tickets, minter, admin, user1 } = await loadFixture(deployTicketsFixture);
            
            const currentTime = await time.latest();
            const params = await createMintParams(user1.address);
            params.transferUnlockAt = currentTime + 2; // Unlock in 2 seconds
            params.expiresAt = currentTime + 5; // Expires in 5 seconds
            await tickets.connect(minter).mintTicket(params);
            
            // Wait for expiry plus grace period (5 minutes)
            await time.increase(306);
            
            await expect(tickets.connect(admin).expireTicket(10001))
                .to.emit(tickets, "TicketExpired")
                .withArgs(10001);
            
            expect(await tickets.isExpired(10001)).to.be.true;
        });
        
        it("Should prevent marking non-expired ticket as expired", async function () {
            const { tickets, minter, admin, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            await expect(tickets.connect(admin).expireTicket(10001))
                .to.be.revertedWith("Ticket not yet expired");
        });
        
        it("Should prevent redeeming expired ticket", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            const currentTime = await time.latest();
            const params = await createMintParams(user1.address);
            params.transferUnlockAt = currentTime + 2; // Unlock in 2 seconds
            params.expiresAt = currentTime + 5; // Expires in 5 seconds
            await tickets.connect(minter).mintTicket(params);
            
            // Wait for expiry
            await time.increase(6);
            
            const tokenId = 10001;
            const nonce = 1;
            const deadline = (await time.latest()) + 3600;
            
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            await expect(tickets.connect(user1).redeemTicket({
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                chainId: 501,
                signature: signature
            })).to.be.revertedWith("Ticket expired");
        });
    });
    
    describe("View Functions", function () {
        it("Should return correct ticket data", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            const ticketData = await tickets.getTicketData(10001);
            
            expect(ticketData.bookingId).to.equal(params.bookingId);
            expect(ticketData.eventId).to.equal(params.eventId);
            expect(ticketData.utilityFlags).to.equal(params.utilityFlags);
            expect(ticketData.transferUnlockAt).to.equal(params.transferUnlockAt);
            expect(ticketData.expiresAt).to.equal(params.expiresAt);
            expect(ticketData.state).to.equal(0); // ISSUED
            expect(ticketData.nonTransferableAfterRedeem).to.be.true;
            expect(ticketData.burnOnRedeem).to.be.false;
        });
        
        it("Should correctly report transferability", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            // Not transferable before unlock
            expect(await tickets.isTransferable(10001)).to.be.false;
            
            // Fast forward past unlock time
            await time.increase(ONE_DAY + 1);
            
            // Now transferable
            expect(await tickets.isTransferable(10001)).to.be.true;
        });
        
        it("Should verify redemption signature", async function () {
            const { tickets, minter, backend, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            await tickets.connect(minter).mintTicket(params);
            
            const tokenId = 10001;
            const nonce = 1;
            const deadline = (await time.latest()) + 3600;
            
            const signature = await createRedemptionSignature(
                tickets,
                tokenId,
                user1,
                nonce,
                deadline,
                backend
            );
            
            const redemptionData = {
                chainId: 501,
                tokenId: tokenId,
                redeemer: user1.address,
                nonce: nonce,
                deadline: deadline,
                signature: signature
            };
            
            expect(await tickets.verifyRedemptionSignature(redemptionData)).to.be.true;
        });
    });
    
    describe("Admin Functions", function () {
        it("Should update royalty info", async function () {
            const { tickets, admin, user1 } = await loadFixture(deployTicketsFixture);
            
            const newReceiver = user1.address;
            const newBps = 1000; // 10%
            
            await expect(tickets.connect(admin).setRoyaltyInfo(newReceiver, newBps))
                .to.emit(tickets, "RoyaltyInfoUpdated")
                .withArgs(newReceiver, newBps);
        });
        
        it("Should prevent setting royalty too high", async function () {
            const { tickets, admin, user1 } = await loadFixture(deployTicketsFixture);
            
            await expect(tickets.connect(admin).setRoyaltyInfo(user1.address, 1001))
                .to.be.revertedWith("Royalty too high");
        });
        
        it("Should pause and unpause contract", async function () {
            const { tickets, minter, pauser, user1 } = await loadFixture(deployTicketsFixture);
            
            // Pause
            await tickets.connect(pauser).pause();
            expect(await tickets.paused()).to.be.true;
            
            // Cannot mint while paused
            const params = await createMintParams(user1.address);
            await expect(tickets.connect(minter).mintTicket(params))
                .to.be.revertedWithCustomError(tickets, "EnforcedPause");
            
            // Unpause
            await tickets.connect(pauser).unpause();
            expect(await tickets.paused()).to.be.false;
            
            // Can mint again
            await tickets.connect(minter).mintTicket(params);
        });
    });
    
    describe("Registry Integration", function () {
        it("Should register with NFTRegistry", async function () {
            const { tickets, registry, admin } = await loadFixture(deployTicketsFixture);
            
            const ticketsAddress = await tickets.getAddress();
            
            await expect(
                registry.connect(admin).registerContract(
                    ticketsAddress,
                    0, // TICKET type
                    "BOGOWI Tickets",
                    "1.0.0"
                )
            ).to.emit(registry, "ContractRegistered")
                .withArgs(ticketsAddress, 0, "BOGOWI Tickets", "1.0.0", admin.address);
            
            expect(await registry.isRegistered(ticketsAddress)).to.be.true;
            
            const info = await registry.getContractInfo(ticketsAddress);
            expect(info.name).to.equal("BOGOWI Tickets");
            expect(info.contractType).to.equal(0); // TICKET
        });
    });
    
    describe("Gas Optimization", function () {
        it("Should batch mint efficiently", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployTicketsFixture);
            
            const batchSize = 10;
            const params = [];
            
            for (let i = 0; i < batchSize; i++) {
                params.push(await createMintParams(user1.address, `EVENT${i}`));
            }
            
            const tx = await tickets.connect(minter).mintBatch(params);
            const receipt = await tx.wait();
            
            // Check gas used per ticket
            const gasPerTicket = receipt.gasUsed / BigInt(batchSize);
            
            // Gas per ticket should be reasonable (less than 200k)
            expect(gasPerTicket).to.be.lt(200000);
        });
    });
    
    describe("Edge Cases", function () {
        it("Should handle maximum values correctly", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            params.utilityFlags = 0xFFFFFFFF; // Max uint32
            // Use reasonable future values that fit within uint64
            const now = await time.latest();
            params.transferUnlockAt = now + ONE_DAY;
            params.expiresAt = now + ONE_WEEK;
            
            // Should succeed with max utility flags
            await tickets.connect(minter).mintTicket(params);
            
            const ticketData = await tickets.getTicketData(10001);
            expect(ticketData.utilityFlags).to.equal(0xFFFFFFFF);
        });
        
        it("Should handle empty metadata URI", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployTicketsFixture);
            
            const params = await createMintParams(user1.address);
            params.metadataURI = "";
            
            await tickets.connect(minter).mintTicket(params);
            
            expect(await tickets.tokenURI(10001)).to.equal("");
        });
        
        it("Should revert on non-existent token queries", async function () {
            const { tickets } = await loadFixture(deployTicketsFixture);
            
            await expect(tickets.getTicketData(99999))
                .to.be.revertedWith("Token does not exist");
            
            await expect(tickets.ownerOf(99999))
                .to.be.revertedWithCustomError(tickets, "ERC721NonexistentToken");
        });
    });
});