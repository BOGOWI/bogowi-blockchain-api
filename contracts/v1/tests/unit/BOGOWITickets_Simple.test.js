const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture, time } = require("@nomicfoundation/hardhat-network-helpers");

describe("BOGOWITickets", function () {
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
    
    describe("Basic Functionality", function () {
        it("Should deploy with correct parameters", async function () {
            const { tickets } = await loadFixture(deployFixture);
            
            expect(await tickets.name()).to.equal("BOGOWI Tickets");
            expect(await tickets.symbol()).to.equal("BWTIX");
        });
        
        it("Should mint a ticket", async function () {
            const { tickets, minter, user1 } = await loadFixture(deployFixture);
            
            const now = await time.latest();
            const params = {
                to: user1.address,
                bookingId: ethers.keccak256(ethers.toUtf8Bytes("BOOKING001")),
                eventId: ethers.keccak256(ethers.toUtf8Bytes("EVENT001")),
                utilityFlags: 0,
                transferUnlockAt: now + 100,
                expiresAt: now + 1000,
                metadataURI: "ipfs://test",
                rewardBasisPoints: 100
            };
            
            await expect(tickets.connect(minter).mintTicket(params))
                .to.emit(tickets, "TicketMinted");
            
            expect(await tickets.ownerOf(10001)).to.equal(user1.address);
        });
    });
});