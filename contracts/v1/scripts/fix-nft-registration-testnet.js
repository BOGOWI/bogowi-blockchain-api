const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🔧 Fixing NFT Registration on TESTNET...\n");
  console.log("=" .repeat(60));

  // Load deployment info
  const network = hre.network.name;
  const deploymentPath = path.join(__dirname, `deployment-nft-${network}.json`);
  
  if (!fs.existsSync(deploymentPath)) {
    console.error("❌ No deployment found for network:", network);
    process.exit(1);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  const [deployer] = await hre.ethers.getSigners();

  console.log("📍 Network:", network);
  console.log("👤 Deployer:", deployer.address);
  console.log("\n📂 Using contracts:");
  console.log("  RoleManager:", deployment.contracts.RoleManager);
  console.log("  NFTRegistry:", deployment.contracts.NFTRegistry);
  console.log("  BOGOWITickets:", deployment.contracts.BOGOWITickets);

  // Get contract instances
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);
  const nftRegistry = await hre.ethers.getContractAt("NFTRegistry", deployment.contracts.NFTRegistry);

  // Setup role
  const CONTRACT_DEPLOYER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("CONTRACT_DEPLOYER_ROLE"));

  console.log("\n" + "=" .repeat(60));
  console.log("1️⃣  Checking CONTRACT_DEPLOYER_ROLE...");
  
  const hasRole = await roleManager.hasRole(CONTRACT_DEPLOYER_ROLE, deployer.address);
  console.log(`   Deployer has CONTRACT_DEPLOYER_ROLE: ${hasRole}`);
  
  if (!hasRole) {
    console.log("   ❌ Deployer needs CONTRACT_DEPLOYER_ROLE");
    console.log("   ⚠️  Admin must grant this role manually");
    console.log("\n   Run this command as admin:");
    console.log(`   await roleManager.grantRole("${CONTRACT_DEPLOYER_ROLE}", "${deployer.address}")`);
    process.exit(1);
  }

  console.log("\n2️⃣  Registering BOGOWITickets in NFTRegistry...");
  
  // Check if already registered
  const isRegistered = await nftRegistry.isRegistered(deployment.contracts.BOGOWITickets);
  
  if (isRegistered) {
    console.log("   ✅ BOGOWITickets is already registered!");
    const info = await nftRegistry.getContractInfo(deployment.contracts.BOGOWITickets);
    console.log(`   Name: ${info.name}`);
    console.log(`   Version: ${info.version}`);
    console.log(`   Type: ${["TICKET", "COLLECTIBLE", "BADGE"][info.contractType]}`);
    console.log(`   Active: ${info.isActive}`);
  } else {
    try {
      const tx = await nftRegistry.registerContract(
        deployment.contracts.BOGOWITickets,
        0, // ContractType.TICKET
        "BOGOWI Event Tickets",
        "1.0.0"
      );
      await tx.wait();
      console.log("   ✅ BOGOWITickets successfully registered!");
      
      // Verify registration
      const info = await nftRegistry.getContractInfo(deployment.contracts.BOGOWITickets);
      console.log(`   Name: ${info.name}`);
      console.log(`   Version: ${info.version}`);
      console.log(`   Type: ${["TICKET", "COLLECTIBLE", "BADGE"][info.contractType]}`);
      console.log(`   Active: ${info.isActive}`);
    } catch (error) {
      console.log(`   ❌ Registration failed: ${error.message}`);
      process.exit(1);
    }
  }

  // Final verification
  console.log("\n3️⃣  Final Verification...");
  const totalContracts = await nftRegistry.getContractCount();
  console.log(`   Total contracts in registry: ${totalContracts}`);
  
  const activeContracts = await nftRegistry.getActiveContracts();
  console.log(`   Active contracts: ${activeContracts.length}`);
  
  console.log("\n" + "=" .repeat(60));
  console.log("✅ NFT Registration Complete!");
  console.log("=" .repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });