const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("📚 Testing NFT Registry Operations...\n");
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
  const [deployer, admin, contractDeployer, backend, user1, user2] = await hre.ethers.getSigners();

  // Get contract instances
  const nftRegistry = await hre.ethers.getContractAt("NFTRegistry", deployment.contracts.NFTRegistry);
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);

  console.log("📍 Using contracts:");
  console.log("  NFTRegistry:", deployment.contracts.NFTRegistry);
  console.log("  RoleManager:", deployment.contracts.RoleManager);
  console.log("\n👥 Test accounts:");
  console.log("  Admin:", admin.address);
  console.log("  Contract Deployer:", contractDeployer.address);

  // Test 1: Query existing registrations
  console.log("\n" + "=" .repeat(60));
  console.log("1️⃣  Test: Query Existing Registrations");
  console.log("-" .repeat(40));

  try {
    const totalContracts = await nftRegistry.getContractCount();
    console.log(`📊 Total registered contracts: ${totalContracts}`);
    
    // Get contracts by type
    const ticketContracts = await nftRegistry.getContractsByType(0); // TICKET
    const collectibleContracts = await nftRegistry.getContractsByType(1); // COLLECTIBLE
    const badgeContracts = await nftRegistry.getContractsByType(2); // BADGE
    
    console.log(`   Ticket contracts: ${ticketContracts.length}`);
    console.log(`   Collectible contracts: ${collectibleContracts.length}`);
    console.log(`   Badge contracts: ${badgeContracts.length}`);
    
    // Get active contracts
    const activeContracts = await nftRegistry.getActiveContracts();
    console.log(`   Active contracts: ${activeContracts.length}`);
    
    // Check BOGOWITickets registration
    if (deployment.contracts.BOGOWITickets) {
      const isRegistered = await nftRegistry.isRegistered(deployment.contracts.BOGOWITickets);
      console.log(`\n✅ BOGOWITickets registered: ${isRegistered}`);
      
      if (isRegistered) {
        const info = await nftRegistry.getContractInfo(deployment.contracts.BOGOWITickets);
        console.log(`   Name: ${info.name}`);
        console.log(`   Version: ${info.version}`);
        console.log(`   Type: ${["TICKET", "COLLECTIBLE", "BADGE"][info.contractType]}`);
        console.log(`   Active: ${info.isActive}`);
      }
    }
  } catch (error) {
    console.log(`❌ Error querying registrations: ${error.message}`);
  }

  // Test 2: Register a new mock contract
  console.log("\n" + "=" .repeat(60));
  console.log("2️⃣  Test: Register New Contract");
  console.log("-" .repeat(40));

  try {
    // Deploy a simple mock ERC721 contract
    const MockNFT = await hre.ethers.getContractFactory("BOGOWITickets");
    const mockNFT = await MockNFT.connect(deployer).deploy(
      deployment.contracts.RoleManager,
      admin.address // conservation DAO
    );
    await mockNFT.waitForDeployment();
    const mockAddress = await mockNFT.getAddress();
    
    console.log(`🎭 Deployed mock NFT at: ${mockAddress}`);
    
    // Register it in the registry - use admin who has CONTRACT_DEPLOYER_ROLE
    const tx = await nftRegistry.connect(admin).registerContract(
      mockAddress,
      2, // BADGE type (contractType comes second)
      "Test Badge NFT",
      "1.0.0"
    );
    await tx.wait();
    
    console.log(`✅ Successfully registered new contract`);
    
    // Verify registration
    const isRegistered = await nftRegistry.isRegistered(mockAddress);
    console.log(`   Registration verified: ${isRegistered}`);
    
    const newTotal = await nftRegistry.getContractCount();
    console.log(`   Total contracts now: ${newTotal}`);
  } catch (error) {
    console.log(`⚠️  Could not register new contract: ${error.message}`);
  }

  // Test 3: Pagination
  console.log("\n" + "=" .repeat(60));
  console.log("3️⃣  Test: Pagination");
  console.log("-" .repeat(40));

  try {
    const pageSize = 10;
    const offset = 0;
    
    const [contracts, hasMore] = await nftRegistry.getActiveContractsPaginated(offset, pageSize);
    console.log(`📄 Retrieved ${contracts.length} contracts (page size: ${pageSize})`);
    console.log(`   Has more pages: ${hasMore}`);
    
    if (contracts.length > 0) {
      console.log("\n   First contract in page:");
      const info = await nftRegistry.getContractInfo(contracts[0]);
      console.log(`     Address: ${contracts[0]}`);
      console.log(`     Name: ${info.name}`);
    }
  } catch (error) {
    console.log(`❌ Error testing pagination: ${error.message}`);
  }

  // Test 4: Deactivate/Reactivate
  console.log("\n" + "=" .repeat(60));
  console.log("4️⃣  Test: Deactivate/Reactivate Contract");
  console.log("-" .repeat(40));

  try {
    const activeContracts = await nftRegistry.getActiveContracts();
    
    if (activeContracts.length > 0) {
      const targetContract = activeContracts[0];
      console.log(`🎯 Testing with contract: ${targetContract}`);
      
      // Deactivate (set status to false)
      const deactivateTx = await nftRegistry.connect(admin).setContractStatus(targetContract, false);
      await deactivateTx.wait();
      console.log(`   ✅ Contract deactivated`);
      
      // Check active status
      const info1 = await nftRegistry.getContractInfo(targetContract);
      console.log(`   Active status: ${info1.isActive}`);
      
      // Reactivate (set status to true)
      const reactivateTx = await nftRegistry.connect(admin).setContractStatus(targetContract, true);
      await reactivateTx.wait();
      console.log(`   ✅ Contract reactivated`);
      
      // Check active status again
      const info2 = await nftRegistry.getContractInfo(targetContract);
      console.log(`   Active status: ${info2.isActive}`);
    } else {
      console.log("   No active contracts to test with");
    }
  } catch (error) {
    console.log(`❌ Error testing deactivation: ${error.message}`);
  }

  // Test 5: Access Control
  console.log("\n" + "=" .repeat(60));
  console.log("5️⃣  Test: Access Control");
  console.log("-" .repeat(40));

  try {
    const activeContracts = await nftRegistry.getActiveContracts();
    
    if (activeContracts.length > 0) {
      const targetContract = activeContracts[0];
      
      // Try to deactivate without permission (should fail)
      try {
        await nftRegistry.connect(user1).setContractStatus(targetContract, false);
        console.log("   ❌ Unauthorized deactivation succeeded (should have failed)");
      } catch (error) {
        console.log("   ✅ Unauthorized deactivation correctly blocked");
      }
      
      // Try to register without permission (should fail)
      try {
        await nftRegistry.connect(user1).registerContract(
          "0x0000000000000000000000000000000000000001",
          0, // TICKET type
          "Fake NFT",
          "1.0.0"
        );
        console.log("   ❌ Unauthorized registration succeeded (should have failed)");
      } catch (error) {
        console.log("   ✅ Unauthorized registration correctly blocked");
      }
    }
  } catch (error) {
    console.log(`❌ Error testing access control: ${error.message}`);
  }

  // Summary
  console.log("\n" + "=" .repeat(60));
  console.log("📊 Registry Test Summary");
  console.log("-" .repeat(40));
  
  try {
    const finalTotal = await nftRegistry.getContractCount();
    const activeContracts = await nftRegistry.getActiveContracts();
    
    console.log(`Total contracts registered: ${finalTotal}`);
    console.log(`Active contracts: ${activeContracts.length}`);
    
    // List all types
    const tickets = await nftRegistry.getContractsByType(0);
    const collectibles = await nftRegistry.getContractsByType(1);
    const badges = await nftRegistry.getContractsByType(2);
    
    console.log(`\nBy type:`);
    console.log(`  Tickets: ${tickets.length}`);
    console.log(`  Collectibles: ${collectibles.length}`);
    console.log(`  Badges: ${badges.length}`);
  } catch (error) {
    console.log(`Error getting summary: ${error.message}`);
  }
  
  console.log("\n✅ NFT Registry tests complete!");
  console.log("=" .repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });