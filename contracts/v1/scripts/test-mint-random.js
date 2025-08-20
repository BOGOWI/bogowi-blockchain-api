const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🎫 Testing NFT Minting with Random IDs...\n");
  console.log("=" .repeat(60));

  // Load deployment info
  const network = hre.network.name;
  const deploymentPath = path.join(__dirname, `deployment-nft-${network}.json`);
  
  if (!fs.existsSync(deploymentPath)) {
    console.error("❌ No deployment found for network:", network);
    console.error("   Run deploy-nft-local.js first!");
    process.exit(1);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  const [deployer, admin, minter, backend, user1, user2] = await hre.ethers.getSigners();

  // Get contract instances
  const tickets = await hre.ethers.getContractAt("BOGOWITickets", deployment.contracts.BOGOWITickets);
  const nftRegistry = await hre.ethers.getContractAt("NFTRegistry", deployment.contracts.NFTRegistry);

  console.log("📍 Using contracts:");
  console.log("  BOGOWITickets:", deployment.contracts.BOGOWITickets);
  console.log("  NFTRegistry:", deployment.contracts.NFTRegistry);
  console.log("\n👥 Test accounts:");
  console.log("  Minter:", minter.address);
  console.log("  User1:", user1.address);
  console.log("  User2:", user2.address);

  // Generate random booking IDs to avoid conflicts
  const randomSuffix = Math.floor(Math.random() * 1000000);

  // Test 1: Single ticket minting
  console.log("\n" + "=" .repeat(60));
  console.log("1️⃣  Test: Single Ticket Minting");
  console.log("-" .repeat(40));

  const currentTime = Math.floor(Date.now() / 1000);
  const mintParams1 = {
    to: user1.address,
    bookingId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes(`BOOKING_${randomSuffix}`)),
    eventId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("EVENT_CONCERT_2025")),
    utilityFlags: 0, // No special flags
    transferUnlockAt: currentTime + 86400, // Locked for 1 day
    expiresAt: currentTime + 2592000, // Expires in 30 days
    metadataURI: `ipfs://QmTest${randomSuffix}/metadata.json`,
    rewardBasisPoints: 500 // 5% rewards
  };

  try {
    console.log("🎟️  Minting ticket for User1...");
    const tx1 = await tickets.connect(minter).mintTicket(mintParams1);
    const receipt1 = await tx1.wait();
    
    // Get the minted token ID from events
    const event = receipt1.logs.find(log => {
      try {
        const parsed = tickets.interface.parseLog(log);
        return parsed.name === 'TicketMinted';
      } catch {
        return false;
      }
    });

    if (event) {
      const parsedEvent = tickets.interface.parseLog(event);
      const tokenId = parsedEvent.args.tokenId;
      const owner = await tickets.ownerOf(tokenId);
      const ticketData = await tickets.tickets(tokenId);
      const uri = await tickets.tokenURI(tokenId);
      
      console.log(`✅ Ticket minted! Token ID: ${tokenId}`);
      console.log(`   Owner: ${owner}`);
      console.log(`   Metadata URI: ${uri}`);
      console.log(`   Event ID: ${ticketData.eventId.slice(0, 10)}...`);
      console.log(`   Expires: ${new Date(Number(ticketData.expiresAt) * 1000).toLocaleString()}`);
      console.log(`   Transferable after: ${new Date(Number(ticketData.transferUnlockAt) * 1000).toLocaleString()}`);
      
      // Store for later tests
      global.lastMintedTokenId = tokenId;
    }
  } catch (error) {
    console.log("❌ Error minting single ticket:", error.message);
  }

  // Test 2: Batch minting
  console.log("\n" + "=" .repeat(60));
  console.log("2️⃣  Test: Batch Ticket Minting");
  console.log("-" .repeat(40));

  const batchMintParams = [];
  for (let i = 0; i < 3; i++) {
    batchMintParams.push({
      to: i === 1 ? user2.address : user1.address,
      bookingId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes(`BATCH_BOOKING_${randomSuffix + i}`)),
      eventId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("EVENT_FESTIVAL_2025")),
      utilityFlags: i === 0 ? 1 : 0, // First ticket has VIP flag
      transferUnlockAt: currentTime + 172800, // Locked for 2 days
      expiresAt: currentTime + 5184000, // Expires in 60 days
      metadataURI: `ipfs://QmBatch${randomSuffix + i}/metadata.json`,
      rewardBasisPoints: [300, 500, 1000][i] // 3%, 5%, 10%
    });
  }

  try {
    console.log("🎟️  Batch minting 3 tickets...");
    const tx2 = await tickets.connect(minter).mintTickets(batchMintParams);
    const receipt2 = await tx2.wait();
    
    const user1Balance = await tickets.balanceOf(user1.address);
    const user2Balance = await tickets.balanceOf(user2.address);
    
    console.log("✅ Batch mint successful!");
    console.log(`   Gas used: ${receipt2.gasUsed}`);
    console.log(`   User1 balance: ${user1Balance} tickets`);
    console.log(`   User2 balance: ${user2Balance} tickets`);
  } catch (error) {
    console.log("❌ Error batch minting:", error.message);
  }

  // Test 3: Transfer restrictions
  console.log("\n" + "=" .repeat(60));
  console.log("3️⃣  Test: Transfer Restrictions");
  console.log("-" .repeat(40));

  if (global.lastMintedTokenId) {
    try {
      console.log("🔒 Testing transfer of locked ticket...");
      const isTransferable = await tickets.isTransferable(global.lastMintedTokenId);
      console.log(`   Is transferable: ${isTransferable}`);
      
      console.log("   Attempting transfer (should fail)...");
      try {
        await tickets.connect(user1).transferFrom(
          user1.address,
          user2.address,
          global.lastMintedTokenId
        );
        console.log("   ❌ Transfer should have been blocked!");
      } catch (transferError) {
        console.log(`   ✅ Transfer correctly blocked: Error: ${transferError.message.split('\n')[0]}...`);
      }
    } catch (error) {
      console.log("❌ Error testing transfer:", error.message);
    }
  }

  // Test 4: Registry queries
  console.log("\n" + "=" .repeat(60));
  console.log("4️⃣  Test: Registry Queries");
  console.log("-" .repeat(40));

  try {
    const ticketContracts = await nftRegistry.getContractsByType(0); // TICKET type
    console.log(`📊 Ticket contracts in registry: ${ticketContracts.length}`);
    
    for (const addr of ticketContracts) {
      const info = await nftRegistry.getContractInfo(addr);
      console.log(`   - ${info.name} (v${info.version})`);
      console.log(`     Address: ${addr}`);
      console.log(`     Active: ${info.isActive}`);
    }

    // Test pagination
    const paginatedContracts = await nftRegistry.getActiveContracts(0, 10);
    console.log(`\n📄 Paginated query: ${paginatedContracts[0].length} active contracts`);
    console.log(`   Has more: ${paginatedContracts[1]}`);
  } catch (error) {
    console.log("❌ Error querying registry:", error.message);
  }

  // Test 5: Ticket redemption signature
  console.log("\n" + "=" .repeat(60));
  console.log("5️⃣  Test: Ticket Redemption Signature");
  console.log("-" .repeat(40));

  try {
    // Create redemption data
    const ticketId = Math.floor(Math.random() * 1000000);
    const redemptionData = {
      ticketId: ticketId,
      recipient: user1.address,
      rewardBasisPoints: 500,
      metadataURI: `ipfs://redeemed/${ticketId}`,
      deadline: currentTime + 3600 // 1 hour deadline
    };

    // Sign with backend
    const domain = {
      name: "BOGOWITickets",
      version: "1",
      chainId: 501, // Camino testnet chain ID
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

    const signature = await backend.signTypedData(domain, types, redemptionData);
    
    // Verify signature
    const isValid = await tickets.isValidRedemptionSignature({
      ...redemptionData,
      signature
    });
    
    console.log(`🔐 Signature valid: ${isValid}`);
  } catch (error) {
    console.log("❌ Error testing signature:", error.message);
  }

  // Summary
  console.log("\n" + "=" .repeat(60));
  console.log("📊 Test Summary");
  console.log("-" .repeat(40));
  
  try {
    const totalSupply = await tickets.nextTokenId() - 1n;
    const user1Balance = await tickets.balanceOf(user1.address);
    const user2Balance = await tickets.balanceOf(user2.address);
    
    console.log(`Total tickets minted: ${totalSupply}`);
    console.log(`User1 balance: ${user1Balance}`);
    console.log(`User2 balance: ${user2Balance}`);
  } catch (error) {
    console.log("Error getting summary:", error.message);
  }

  console.log("\n✅ NFT minting tests with random IDs complete!");
  console.log("=" .repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });