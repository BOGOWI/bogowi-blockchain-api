const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🔧 Configuring NFT Contracts on MAINNET...\n");
  console.log("=" .repeat(60));

  // Verify we're on mainnet
  const network = hre.network.name;
  const chainId = await hre.ethers.provider.getNetwork().then(n => n.chainId);
  
  if (chainId !== 500n) {
    console.error("❌ This script is for Camino MAINNET only!");
    console.error("   Current chain ID:", chainId);
    console.error("   Expected: 500 (Camino Mainnet)");
    process.exit(1);
  }

  // Load deployment info
  const deploymentPath = path.join(__dirname, `deployment-nft-${network}.json`);
  if (!fs.existsSync(deploymentPath)) {
    console.error("❌ No deployment found for mainnet");
    process.exit(1);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  
  // Use ADMIN_PRIVATE_KEY from environment
  const ADMIN_PRIVATE_KEY = process.env.ADMIN_PRIVATE_KEY;
  if (!ADMIN_PRIVATE_KEY) {
    console.error("❌ ADMIN_PRIVATE_KEY not found in environment");
    console.error("   Make sure NODE_ENV=production is set");
    process.exit(1);
  }

  // Create admin signer
  const adminWallet = new hre.ethers.Wallet(ADMIN_PRIVATE_KEY, hre.ethers.provider);
  console.log("👤 Admin:", adminWallet.address);
  
  const balance = await hre.ethers.provider.getBalance(adminWallet.address);
  console.log("💰 Admin balance:", hre.ethers.formatEther(balance), "CAM");
  
  if (balance < hre.ethers.parseEther("0.5")) {
    console.error("❌ Insufficient balance for configuration");
    process.exit(1);
  }

  console.log("\n📂 Using contracts:");
  console.log("  RoleManager:", deployment.contracts.RoleManager);
  console.log("  NFTRegistry:", deployment.contracts.NFTRegistry);
  console.log("  BOGOWITickets:", deployment.contracts.BOGOWITickets);

  // Get contract instances with admin signer
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager, adminWallet);
  const nftRegistry = await hre.ethers.getContractAt("NFTRegistry", deployment.contracts.NFTRegistry, adminWallet);

  console.log("\n" + "=" .repeat(60));
  console.log("⚙️  Starting Configuration...");
  console.log("=" .repeat(60) + "\n");

  // 1. Register contracts with RoleManager
  console.log("1️⃣  Registering contracts with RoleManager...");
  
  try {
    console.log("   Registering NFTRegistry...");
    const tx1 = await roleManager.registerContract(deployment.contracts.NFTRegistry, "NFTRegistry");
    await tx1.wait();
    console.log("   ✅ NFTRegistry registered");
  } catch (error) {
    if (error.message.includes("already registered")) {
      console.log("   ℹ️  NFTRegistry already registered");
    } else {
      console.log("   ⚠️  NFTRegistry registration failed:", error.message);
    }
  }

  try {
    console.log("   Registering BOGOWITickets...");
    const tx2 = await roleManager.registerContract(deployment.contracts.BOGOWITickets, "BOGOWITickets");
    await tx2.wait();
    console.log("   ✅ BOGOWITickets registered");
  } catch (error) {
    if (error.message.includes("already registered")) {
      console.log("   ℹ️  BOGOWITickets already registered");
    } else {
      console.log("   ⚠️  BOGOWITickets registration failed:", error.message);
    }
  }

  // 2. Setup roles
  console.log("\n2️⃣  Setting up roles...");
  const REGISTRY_ADMIN_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("REGISTRY_ADMIN_ROLE"));
  const CONTRACT_DEPLOYER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("CONTRACT_DEPLOYER_ROLE"));

  try {
    console.log("   Granting REGISTRY_ADMIN_ROLE to admin...");
    const tx3 = await roleManager.grantRole(REGISTRY_ADMIN_ROLE, adminWallet.address);
    await tx3.wait();
    console.log("   ✅ REGISTRY_ADMIN_ROLE granted");
  } catch (error) {
    if (error.message.includes("already has role")) {
      console.log("   ℹ️  Admin already has REGISTRY_ADMIN_ROLE");
    } else {
      console.log("   ⚠️  Failed:", error.message);
    }
  }

  try {
    console.log("   Granting CONTRACT_DEPLOYER_ROLE to admin...");
    const tx4 = await roleManager.grantRole(CONTRACT_DEPLOYER_ROLE, adminWallet.address);
    await tx4.wait();
    console.log("   ✅ CONTRACT_DEPLOYER_ROLE granted");
  } catch (error) {
    if (error.message.includes("already has role")) {
      console.log("   ℹ️  Admin already has CONTRACT_DEPLOYER_ROLE");
    } else {
      console.log("   ⚠️  Failed:", error.message);
    }
  }

  // 3. Register BOGOWITickets in NFTRegistry
  console.log("\n3️⃣  Registering BOGOWITickets in NFTRegistry...");
  
  // Check if already registered
  const isRegistered = await nftRegistry.isRegistered(deployment.contracts.BOGOWITickets);
  
  if (isRegistered) {
    console.log("   ℹ️  BOGOWITickets already registered");
    const info = await nftRegistry.getContractInfo(deployment.contracts.BOGOWITickets);
    console.log(`   Name: ${info.name}`);
    console.log(`   Version: ${info.version}`);
    console.log(`   Type: ${["TICKET", "COLLECTIBLE", "BADGE"][info.contractType]}`);
    console.log(`   Active: ${info.isActive}`);
  } else {
    try {
      const tx5 = await nftRegistry.registerContract(
        deployment.contracts.BOGOWITickets,
        0, // ContractType.TICKET
        "BOGOWI Event Tickets",
        "1.0.0"
      );
      await tx5.wait();
      console.log("   ✅ BOGOWITickets registered successfully!");
      
      // Verify registration
      const info = await nftRegistry.getContractInfo(deployment.contracts.BOGOWITickets);
      console.log(`   Name: ${info.name}`);
      console.log(`   Version: ${info.version}`);
      console.log(`   Type: ${["TICKET", "COLLECTIBLE", "BADGE"][info.contractType]}`);
      console.log(`   Active: ${info.isActive}`);
    } catch (error) {
      console.log(`   ❌ Registration failed: ${error.message}`);
    }
  }

  // 4. Final verification
  console.log("\n4️⃣  Final Verification...");
  const totalContracts = await nftRegistry.getContractCount();
  console.log(`   Total contracts in registry: ${totalContracts}`);
  
  const activeContracts = await nftRegistry.getActiveContracts();
  console.log(`   Active contracts: ${activeContracts.length}`);
  
  // Update deployment status
  deployment.status = "DEPLOYED_AND_CONFIGURED";
  deployment.configuredAt = new Date().toISOString();
  deployment.configuredBy = adminWallet.address;
  
  fs.writeFileSync(deploymentPath, JSON.stringify(deployment, null, 2));
  
  console.log("\n" + "=" .repeat(60));
  console.log("✅ MAINNET NFT Configuration Complete!");
  console.log("=" .repeat(60));
  console.log("\n📍 Deployed Contracts:");
  console.log("  NFTRegistry:", deployment.contracts.NFTRegistry);
  console.log("  BOGOWITickets:", deployment.contracts.BOGOWITickets);
  console.log("\n🔒 MAINNET Next Steps:");
  console.log("  1. Grant NFT_MINTER_ROLE to authorized minting addresses");
  console.log("  2. Grant BACKEND_ROLE to backend service addresses");
  console.log("  3. Test with small amounts before production use");
  console.log("  4. Monitor contract activity");
  console.log("=" .repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });