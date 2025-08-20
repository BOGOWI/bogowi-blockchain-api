const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture, time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGOWITickets Security Fixes", function () {
    async function deployFixture() {
        const [owner, admin, minter, user1] = await ethers.getSigners();
        
        // Deploy RoleManager
        const RoleManager = await ethers.getContractFactory("RoleManager");
        const roleManager = await RoleManager.deploy();
        await roleManager.waitForDeployment();
        
        // Setup roles
        const NFT_MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
        const ADMIN_ROLE = ethers.keccak256(ethers.toUtf8Bytes("ADMIN_ROLE"));
        
        // Grant roles
        await roleManager.grantRole(ADMIN_ROLE, admin.address);
        await roleManager.grantRole(NFT_MINTER_ROLE, minter.address);
        
        // Deploy BOGOWITickets
        const BOGOWITickets = await ethers.getContractFactory("BOGOWITickets");
        const tickets = await BOGOWITickets.deploy(
            await roleManager.getAddress(),
            owner.address // conservationDAO
        );
        await tickets.waitForDeployment();
        
        // Register tickets contract with RoleManager
        await roleManager.registerContract(await tickets.getAddress(), "BOGOWITickets");
        
        return { tickets, roleManager, owner, admin, minter, user1 };
    }
    
    describe("Security Fix Tests", function () {
        it("Should enforce MAX_BATCH_SIZE limit", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployFixture);
            
            const now = await time.latest();
            const params = [];
            
            // Try to exceed MAX_BATCH_SIZE (100)
            for (let i = 0; i < 101; i++) {
                params.push({
                    to: user1.address,
                    bookingId: ethers.keccak256(ethers.toUtf8Bytes(`BOOKING${i}`)),
                    eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT001")),
                    utilityFlags: 0,
                    transferUnlockAt: now + 100,
                    expiresAt: now + 1000,
                    metadataURI: "",
                    rewardBasisPoints: 0
                });
            }
            
            await expect(tickets.connect(minter).mintBatch(params))
                .to.be.revertedWith("Batch size exceeds maximum");
        });
        
        it("Should enforce grace period for expireTicket", async function () {
            const { tickets, minter, admin, user1 } = await loadFixture(deployFixture);
            
            const now = await time.latest();
            const params = {
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("GRACE_TEST")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT001")),
                utilityFlags: 0,
                transferUnlockAt: now + 5,
                expiresAt: now + 10,
                metadataURI: "",
                rewardBasisPoints: 0
            };
            
            await tickets.connect(minter).mintTicket(params);
            
            // Wait for expiry
            await time.increase(11);
            
            // Should fail - grace period not met
            await expect(tickets.connect(admin).expireTicket(10001))
                .to.be.revertedWith("Grace period not met");
            
            // Wait for grace period (5 minutes)
            await time.increase(300);
            
            // Should succeed now
            await tickets.connect(admin).expireTicket(10001);
        });
        
        it("Should validate timestamp ranges", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployFixture);
            
            const now = await time.latest();
            
            // Test expiry too far in future
            const params1 = {
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("FAR_FUTURE")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT001")),
                utilityFlags: 0,
                transferUnlockAt: now + 100,
                expiresAt: now + (366 * 24 * 60 * 60), // > 365 days
                metadataURI: "",
                rewardBasisPoints: 0
            };
            
            await expect(tickets.connect(minter).mintTicket(params1))
                .to.be.revertedWith("Expiry too far in future");
            
        });
        
        it("Should emit new events for state changes", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployFixture);
            
            const now = await time.latest();
            const params = [{
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("EVENT_TEST")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT001")),
                utilityFlags: 0,
                transferUnlockAt: now + 100,
                expiresAt: now + 1000,
                metadataURI: "",
                rewardBasisPoints: 0
            }];
            
            // Should emit BatchMintStarted event
            await expect(tickets.connect(minter).mintBatch(params))
                .to.emit(tickets, "BatchMintStarted")
                .withArgs(1, minter.address);
        });
        
        it("Should use INITIAL_TOKEN_ID constant", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployFixture);
            
            const now = await time.latest();
            const params = {
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("TOKEN_ID_TEST")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT001")),
                utilityFlags: 0,
                transferUnlockAt: now + 100,
                expiresAt: now + 1000,
                metadataURI: "",
                rewardBasisPoints: 0
            };
            
            await tickets.connect(minter).mintTicket(params);
            
            // First token should be INITIAL_TOKEN_ID (10001)
            expect(await tickets.ownerOf(10001)).to.equal(user1.address);
        });
    });
});