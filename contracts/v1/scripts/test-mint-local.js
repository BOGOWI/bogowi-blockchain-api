const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🎫 Testing NFT Minting Operations...\n");
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

  // Test 1: Single ticket minting
  console.log("\n" + "=" .repeat(60));
  console.log("1️⃣  Test: Single Ticket Minting");
  console.log("-" .repeat(40));

  const currentTime = Math.floor(Date.now() / 1000);
  // Use timestamp to make booking IDs unique
  const uniqueId = Date.now().toString();
  const mintParams1 = {
    to: user1.address,
    bookingId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes(`BOOKING_${uniqueId}_001`)),
    eventId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("EVENT_CONCERT_2025")),
    utilityFlags: 0, // No special flags
    transferUnlockAt: currentTime + 86400, // Locked for 1 day
    expiresAt: currentTime + 2592000, // Expires in 30 days
    metadataURI: "ipfs://QmTest123/metadata.json",
    rewardBasisPoints: 500 // 5% rewards
  };

  try {
    console.log("🎟️  Minting ticket for User1...");
    const tx1 = await tickets.connect(minter).mintTicket(mintParams1);
    const receipt1 = await tx1.wait();
    
    // Get token ID from event
    const mintEvent = receipt1.logs.find(log => {
      try {
        const parsed = tickets.interface.parseLog(log);
        return parsed.name === 'TicketMinted';
      } catch {
        return false;
      }
    });
    
    const tokenId1 = mintEvent ? tickets.interface.parseLog(mintEvent).args[0] : 10001;
    console.log(`✅ Ticket minted! Token ID: ${tokenId1}`);
    
    // Verify ownership
    const owner1 = await tickets.ownerOf(tokenId1);
    console.log(`   Owner: ${owner1}`);
    console.log(`   Metadata URI: ${await tickets.tokenURI(tokenId1)}`);
    
    // Get ticket data
    const ticketData = await tickets.getTicketData(tokenId1);
    console.log(`   Event ID: ${ticketData.eventId.substring(0, 10)}...`);
    console.log(`   Expires: ${new Date(Number(ticketData.expiresAt) * 1000).toLocaleString()}`);
    console.log(`   Transferable after: ${new Date(Number(ticketData.transferUnlockAt) * 1000).toLocaleString()}`);
  } catch (error) {
    console.log(`❌ Error minting single ticket: ${error.message}`);
  }

  // Test 2: Batch minting
  console.log("\n" + "=" .repeat(60));
  console.log("2️⃣  Test: Batch Ticket Minting");
  console.log("-" .repeat(40));

  const batchParams = [
    {
      to: user1.address,
      bookingId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes(`BATCH_${uniqueId}_001`)),
      eventId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("EVENT_FESTIVAL_2025")),
      utilityFlags: 1, // Non-transferable after redeem
      transferUnlockAt: currentTime,
      expiresAt: currentTime + 604800, // 7 days
      metadataURI: "ipfs://QmBatch001/1.json",
      rewardBasisPoints: 250
    },
    {
      to: user2.address,
      bookingId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes(`BATCH_${uniqueId}_002`)),
      eventId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("EVENT_FESTIVAL_2025")),
      utilityFlags: 2, // Burn on redeem
      transferUnlockAt: currentTime,
      expiresAt: currentTime + 604800,
      metadataURI: "ipfs://QmBatch001/2.json",
      rewardBasisPoints: 250
    },
    {
      to: user2.address,
      bookingId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes(`BATCH_${uniqueId}_003`)),
      eventId: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("EVENT_FESTIVAL_2025")),
      utilityFlags: 3, // Both flags set
      transferUnlockAt: currentTime,
      expiresAt: currentTime + 604800,
      metadataURI: "ipfs://QmBatch001/3.json",
      rewardBasisPoints: 250
    }
  ];

  try {
    console.log(`🎟️  Batch minting ${batchParams.length} tickets...`);
    const tx2 = await tickets.connect(minter).mintBatch(batchParams);
    const receipt2 = await tx2.wait();
    
    console.log(`✅ Batch mint successful!`);
    console.log(`   Gas used: ${receipt2.gasUsed}`);
    
    // Check balances
    const balance1 = await tickets.balanceOf(user1.address);
    const balance2 = await tickets.balanceOf(user2.address);
    console.log(`   User1 balance: ${balance1} tickets`);
    console.log(`   User2 balance: ${balance2} tickets`);
  } catch (error) {
    console.log(`❌ Error batch minting: ${error.message}`);
  }

  // Test 3: Transfer restrictions
  console.log("\n" + "=" .repeat(60));
  console.log("3️⃣  Test: Transfer Restrictions");
  console.log("-" .repeat(40));

  try {
    // Try to transfer a locked ticket (should fail)
    console.log("🔒 Testing transfer of locked ticket...");
    const tokenId = 10001; // First minted ticket
    
    // Check if transferable
    const isTransferable = await tickets.isTransferable(tokenId);
    console.log(`   Is transferable: ${isTransferable}`);
    
    if (!isTransferable) {
      console.log("   Attempting transfer (should fail)...");
      try {
        await tickets.connect(user1).transferFrom(user1.address, user2.address, tokenId);
        console.log("   ❌ Transfer succeeded when it shouldn't!");
      } catch (error) {
        console.log("   ✅ Transfer correctly blocked:", error.message.substring(0, 50) + "...");
      }
    }
  } catch (error) {
    console.log(`❌ Error testing transfers: ${error.message}`);
  }

  // Test 4: Registry queries
  console.log("\n" + "=" .repeat(60));
  console.log("4️⃣  Test: Registry Queries");
  console.log("-" .repeat(40));

  try {
    // Get all ticket contracts
    const ticketContracts = await nftRegistry.getContractsByType(0); // TICKET type
    console.log(`📊 Ticket contracts in registry: ${ticketContracts.length}`);
    
    for (const contractAddr of ticketContracts) {
      if (contractAddr !== hre.ethers.ZeroAddress) {
        const info = await nftRegistry.getContractInfo(contractAddr);
        console.log(`   - ${info.name} (v${info.version})`);
        console.log(`     Address: ${contractAddr}`);
        console.log(`     Active: ${info.isActive}`);
      }
    }
    
    // Get paginated active contracts
    const [contracts, hasMore] = await nftRegistry.getActiveContractsPaginated(0, 10);
    console.log(`\n📄 Paginated query: ${contracts.length} active contracts`);
    console.log(`   Has more: ${hasMore}`);
  } catch (error) {
    console.log(`❌ Error querying registry: ${error.message}`);
  }

  // Test 5: Ticket redemption (signature required)
  console.log("\n" + "=" .repeat(60));
  console.log("5️⃣  Test: Ticket Redemption Signature");
  console.log("-" .repeat(40));

  let testTokenId = 10001; // Default token ID
  
  try {
    // Try to get an actual token ID that was minted to user1
    try {
      const user1Balance = await tickets.balanceOf(user1.address);
      if (user1Balance > 0) {
        // Get the first token owned by user1
        testTokenId = await tickets.tokenOfOwnerByIndex(user1.address, 0);
        console.log(`   Using actual minted token ID: ${testTokenId}`);
      }
    } catch (e) {
      // If tokenOfOwnerByIndex doesn't work, use the token ID from first mint
      console.log(`   Using expected token ID: ${testTokenId}`);
    }
    
    const tokenId = testTokenId;
    const nonce = 1;
    const deadline = currentTime + 3600; // 1 hour from now
    
    // Create EIP-712 signature
    const domain = {
      name: "BOGOWITickets",
      version: "1",
      chainId: 501, // Camino testnet chain ID (configured in hardhat.config.js)
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
      redeemer: user1.address,
      nonce: nonce,
      deadline: deadline,
      chainId: 501
    };
    
    const signature = await backend.signTypedData(domain, types, value);
    
    const redemptionData = {
      tokenId: tokenId,
      redeemer: user1.address,
      nonce: nonce,
      deadline: deadline,
      chainId: 501,
      signature: signature
    };
    
    // Verify signature
    const isValid = await tickets.verifyRedemptionSignature(redemptionData);
    console.log(`🔐 Signature valid: ${isValid}`);
    
    if (isValid) {
      console.log("   ✅ Redemption signature created successfully");
      console.log("   Note: Actual redemption would mark ticket as used");
    }
  } catch (error) {
    console.log(`❌ Error testing redemption: ${error.message}`);
  }

  // Summary
  console.log("\n" + "=" .repeat(60));
  console.log("📊 Test Summary");
  console.log("-" .repeat(40));
  
  try {
    // BOGOWITickets doesn't have totalSupply function, use balances instead
    const user1Balance = await tickets.balanceOf(user1.address);
    const user2Balance = await tickets.balanceOf(user2.address);
    
    console.log(`User1 owns: ${user1Balance} tickets`);
    console.log(`User2 owns: ${user2Balance} tickets`);
    console.log(`Total tickets minted: ${Number(user1Balance) + Number(user2Balance)}`);
    
    const registryTotal = await nftRegistry.getContractCount();
    console.log(`Total contracts in registry: ${registryTotal}`);
  } catch (error) {
    console.log(`Error getting summary: ${error.message}`);
  }
  
  console.log("\n✅ NFT minting tests complete!");
  console.log("=" .repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });