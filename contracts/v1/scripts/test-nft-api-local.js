const hre = require("hardhat");
const { ethers } = require("hardhat");
const fs = require("fs");
const path = require("path");
const express = require("express");
const bodyParser = require("body-parser");

// Load deployment info
const deployment = JSON.parse(
  fs.readFileSync(path.join(__dirname, "deployment-nft-localhost.json"), 'utf8')
);

// Express API setup
const app = express();
app.use(bodyParser.json());
const PORT = process.env.PORT || 3000;

// Contract instances
let roleManager, bogoToken, nftRegistry, tickets;
let signers = {};

// Initialize contracts and signers
async function initializeContracts() {
  console.log("🔗 Initializing contracts...\n");
  
  // Get signers
  const allSigners = await ethers.getSigners();
  signers = {
    deployer: allSigners[0],
    admin: allSigners[1],
    minter: allSigners[2],
    backend: allSigners[3],
    user1: allSigners[4],
    user2: allSigners[5]
  };

  // Get contract instances
  roleManager = await ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);
  bogoToken = await ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  nftRegistry = await ethers.getContractAt("NFTRegistry", deployment.contracts.NFTRegistry);
  tickets = await ethers.getContractAt("BOGOWITickets", deployment.contracts.BOGOWITickets);

  console.log("✅ Contracts initialized");
  console.log("  RoleManager:", deployment.contracts.RoleManager);
  console.log("  BOGOToken:", deployment.contracts.BOGOToken);
  console.log("  NFTRegistry:", deployment.contracts.NFTRegistry);
  console.log("  BOGOWITickets:", deployment.contracts.BOGOWITickets);
  console.log("\n📝 Test accounts:");
  console.log("  Admin:", signers.admin.address);
  console.log("  Minter:", signers.minter.address);
  console.log("  Backend:", signers.backend.address);
  console.log("  User1:", signers.user1.address);
  console.log("  User2:", signers.user2.address);
}

// Helper function to create redemption signature
async function createRedemptionSignature(
  ticketId,
  recipient,
  rewardBasisPoints,
  metadataURI,
  deadline,
  signer
) {
  const domain = {
    name: "BOGOWITickets",
    version: "1",
    chainId: 31337, // Local chain ID
    verifyingContract: await tickets.getAddress()
  };

  const types = {
    RedemptionData: [
      { name: "ticketId", type: "uint256" },
      { name: "recipient", type: "address" },
      { name: "rewardBasisPoints", type: "uint256" },
      { name: "metadataURI", type: "string" },
      { name: "deadline", type: "uint256" }
    ]
  };

  const value = {
    ticketId,
    recipient,
    rewardBasisPoints,
    metadataURI,
    deadline
  };

  return await signer.signTypedData(domain, types, value);
}

// ============ API ENDPOINTS ============

// 1. GET /status - Check system status
app.get("/status", async (req, res) => {
  try {
    const [totalSupply, isPaused, registryCount] = await Promise.all([
      tickets.totalSupply(),
      tickets.paused(),
      nftRegistry.getContractCount()
    ]);

    res.json({
      status: "online",
      network: "localhost",
      contracts: {
        tickets: await tickets.getAddress(),
        registry: await nftRegistry.getAddress(),
        token: await bogoToken.getAddress()
      },
      stats: {
        totalTickets: totalSupply.toString(),
        isPaused: isPaused,
        registeredContracts: registryCount.toString()
      }
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 2. POST /mint - Mint a single ticket
app.post("/mint", async (req, res) => {
  try {
    const { recipient, ticketId, rewardBasisPoints, metadataURI, deadline } = req.body;
    
    // Validate inputs
    if (!recipient || ticketId === undefined || !rewardBasisPoints || !metadataURI) {
      return res.status(400).json({ error: "Missing required parameters" });
    }

    // Set deadline if not provided (1 hour from now)
    const finalDeadline = deadline || Math.floor(Date.now() / 1000) + 3600;

    // Create signature as backend
    const signature = await createRedemptionSignature(
      ticketId,
      recipient,
      rewardBasisPoints,
      metadataURI,
      finalDeadline,
      signers.backend
    );

    // Mint as minter
    const tx = await tickets.connect(signers.minter).mint(
      recipient,
      ticketId,
      rewardBasisPoints,
      metadataURI,
      finalDeadline,
      signature
    );

    const receipt = await tx.wait();
    
    res.json({
      success: true,
      txHash: receipt.hash,
      ticketId: ticketId.toString(),
      recipient: recipient,
      gasUsed: receipt.gasUsed.toString()
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 3. POST /mint-batch - Mint multiple tickets
app.post("/mint-batch", async (req, res) => {
  try {
    const { recipients, ticketIds, rewardBasisPoints, metadataURIs, deadlines } = req.body;
    
    // Validate inputs
    if (!recipients || !ticketIds || !rewardBasisPoints || !metadataURIs) {
      return res.status(400).json({ error: "Missing required parameters" });
    }

    if (recipients.length !== ticketIds.length || 
        recipients.length !== rewardBasisPoints.length ||
        recipients.length !== metadataURIs.length) {
      return res.status(400).json({ error: "Array lengths must match" });
    }

    // Set deadlines if not provided
    const finalDeadlines = deadlines || recipients.map(() => Math.floor(Date.now() / 1000) + 3600);

    // Create signatures for each ticket
    const signatures = await Promise.all(
      ticketIds.map((ticketId, i) => 
        createRedemptionSignature(
          ticketId,
          recipients[i],
          rewardBasisPoints[i],
          metadataURIs[i],
          finalDeadlines[i],
          signers.backend
        )
      )
    );

    // Mint batch as minter
    const tx = await tickets.connect(signers.minter).mintBatch(
      recipients,
      ticketIds,
      rewardBasisPoints,
      metadataURIs,
      finalDeadlines,
      signatures
    );

    const receipt = await tx.wait();
    
    res.json({
      success: true,
      txHash: receipt.hash,
      count: recipients.length,
      gasUsed: receipt.gasUsed.toString()
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 4. GET /ticket/:id - Get ticket details
app.get("/ticket/:id", async (req, res) => {
  try {
    const ticketId = req.params.id;
    
    // Check if ticket exists
    const owner = await tickets.ownerOf(ticketId).catch(() => null);
    if (!owner) {
      return res.status(404).json({ error: "Ticket not found" });
    }

    // Get ticket data
    const [ticketData, tokenURI] = await Promise.all([
      tickets.tickets(ticketId),
      tickets.tokenURI(ticketId)
    ]);

    res.json({
      ticketId: ticketId,
      owner: owner,
      rewardBasisPoints: ticketData.rewardBasisPoints.toString(),
      mintedAt: ticketData.mintedAt.toString(),
      metadataURI: tokenURI,
      rewardPercentage: (Number(ticketData.rewardBasisPoints) / 100).toFixed(2) + "%"
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 5. GET /user/:address/tickets - Get user's tickets
app.get("/user/:address/tickets", async (req, res) => {
  try {
    const userAddress = req.params.address;
    
    // Validate address
    if (!ethers.isAddress(userAddress)) {
      return res.status(400).json({ error: "Invalid address" });
    }

    const balance = await tickets.balanceOf(userAddress);
    const ticketIds = [];
    
    // Get all ticket IDs owned by user
    for (let i = 0; i < balance; i++) {
      const tokenId = await tickets.tokenOfOwnerByIndex(userAddress, i);
      ticketIds.push(tokenId.toString());
    }

    // Get details for each ticket
    const ticketDetails = await Promise.all(
      ticketIds.map(async (id) => {
        const data = await tickets.tickets(id);
        return {
          ticketId: id,
          rewardBasisPoints: data.rewardBasisPoints.toString(),
          mintedAt: data.mintedAt.toString(),
          rewardPercentage: (Number(data.rewardBasisPoints) / 100).toFixed(2) + "%"
        };
      })
    );

    res.json({
      address: userAddress,
      balance: balance.toString(),
      tickets: ticketDetails
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 6. POST /redeem - Redeem a ticket (burn on redeem)
app.post("/redeem", async (req, res) => {
  try {
    const { ticketId, recipient, rewardBasisPoints, metadataURI } = req.body;
    
    if (!ticketId || !recipient || !rewardBasisPoints || !metadataURI) {
      return res.status(400).json({ error: "Missing required parameters" });
    }

    // Create redemption data
    const deadline = Math.floor(Date.now() / 1000) + 3600;
    const signature = await createRedemptionSignature(
      ticketId,
      recipient,
      rewardBasisPoints,
      metadataURI,
      deadline,
      signers.backend
    );

    const redemptionData = {
      ticketId,
      recipient,
      rewardBasisPoints,
      metadataURI,
      deadline,
      signature
    };

    // Redeem as the ticket owner
    const owner = await tickets.ownerOf(ticketId);
    const ownerSigner = Object.values(signers).find(s => s.address === owner);
    
    if (!ownerSigner) {
      return res.status(400).json({ error: "Ticket owner not found in test accounts" });
    }

    const tx = await tickets.connect(ownerSigner).redeemTicket(redemptionData);
    const receipt = await tx.wait();

    res.json({
      success: true,
      txHash: receipt.hash,
      ticketId: ticketId,
      burned: true,
      gasUsed: receipt.gasUsed.toString()
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 7. GET /registry - Get all registered NFT contracts
app.get("/registry", async (req, res) => {
  try {
    const count = await nftRegistry.getContractCount();
    const contracts = [];
    
    for (let i = 0; i < count; i++) {
      const address = await nftRegistry.getContractAtIndex(i);
      const info = await nftRegistry.getContractInfo(address);
      contracts.push({
        address: address,
        name: info.name,
        version: info.version,
        contractType: ["ERC721", "ERC1155", "MIXED"][info.contractType],
        isActive: info.isActive
      });
    }

    res.json({
      count: count.toString(),
      contracts: contracts
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 8. POST /pause - Pause the contract (admin only)
app.post("/pause", async (req, res) => {
  try {
    const tx = await tickets.connect(signers.admin).pause();
    await tx.wait();
    
    res.json({
      success: true,
      paused: true,
      message: "Contract paused successfully"
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 9. POST /unpause - Unpause the contract (admin only)
app.post("/unpause", async (req, res) => {
  try {
    const tx = await tickets.connect(signers.admin).unpause();
    await tx.wait();
    
    res.json({
      success: true,
      paused: false,
      message: "Contract unpaused successfully"
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// 10. GET /stats - Get comprehensive statistics
app.get("/stats", async (req, res) => {
  try {
    const [totalSupply, isPaused, registryCount, baseURI, gracePeriod] = await Promise.all([
      tickets.totalSupply(),
      tickets.paused(),
      nftRegistry.getContractCount(),
      tickets.baseURI(),
      tickets.expiryGracePeriod()
    ]);

    // Get role counts
    const roles = {
      admins: 0,
      minters: 0,
      backends: 0
    };

    // Count role members (simplified for demo)
    const DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
    const NFT_MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
    const BACKEND_ROLE = ethers.keccak256(ethers.toUtf8Bytes("BACKEND_ROLE"));

    for (const signer of Object.values(signers)) {
      if (await roleManager.hasRole(DEFAULT_ADMIN_ROLE, signer.address)) roles.admins++;
      if (await roleManager.hasRole(NFT_MINTER_ROLE, signer.address)) roles.minters++;
      if (await roleManager.hasRole(BACKEND_ROLE, signer.address)) roles.backends++;
    }

    res.json({
      tickets: {
        totalSupply: totalSupply.toString(),
        isPaused: isPaused,
        baseURI: baseURI,
        expiryGracePeriod: gracePeriod.toString() + " seconds"
      },
      registry: {
        contractCount: registryCount.toString()
      },
      roles: roles,
      network: {
        chainId: 31337,
        name: "localhost"
      }
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// Start the server
async function startServer() {
  await initializeContracts();
  
  app.listen(PORT, () => {
    console.log(`\n🚀 NFT API Server running on http://localhost:${PORT}`);
    console.log("\n📚 Available endpoints:");
    console.log("  GET  /status                - System status");
    console.log("  GET  /stats                 - Comprehensive statistics");
    console.log("  POST /mint                  - Mint single ticket");
    console.log("  POST /mint-batch            - Mint multiple tickets");
    console.log("  GET  /ticket/:id            - Get ticket details");
    console.log("  GET  /user/:address/tickets - Get user's tickets");
    console.log("  POST /redeem                - Redeem ticket (burns it)");
    console.log("  GET  /registry              - List registered contracts");
    console.log("  POST /pause                 - Pause contract");
    console.log("  POST /unpause               - Unpause contract");
    console.log("\n💡 Use curl or Postman to test the endpoints");
  });
}

// Handle errors
process.on('unhandledRejection', (error) => {
  console.error('Unhandled promise rejection:', error);
});

// Start the server
startServer().catch(console.error);