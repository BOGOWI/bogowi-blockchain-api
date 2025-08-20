const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🔐 Granting NFT Roles on MAINNET...\n");
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
  
  // Backend/Minter address (from your .env BACKEND_WALLET_ADDRESS)
  const BACKEND_ADDRESS = "0xB34A822F735CDE477cbB39a06118267D00948ef7";
  console.log("🔧 Backend/Minter Address:", BACKEND_ADDRESS);

  const balance = await hre.ethers.provider.getBalance(adminWallet.address);
  console.log("💰 Admin balance:", hre.ethers.formatEther(balance), "CAM");

  console.log("\n📂 Using RoleManager:", deployment.contracts.RoleManager);

  // Get RoleManager instance with admin signer
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager, adminWallet);

  // Define roles
  const NFT_MINTER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
  const BACKEND_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("BACKEND_ROLE"));

  console.log("\n" + "=" .repeat(60));
  console.log("⚙️  Granting Roles...");
  console.log("=" .repeat(60) + "\n");

  // 1. Grant NFT_MINTER_ROLE
  console.log("1️⃣  Granting NFT_MINTER_ROLE...");
  try {
    const hasMinterRole = await roleManager.hasRole(NFT_MINTER_ROLE, BACKEND_ADDRESS);
    if (hasMinterRole) {
      console.log("   ℹ️  Address already has NFT_MINTER_ROLE");
    } else {
      const tx1 = await roleManager.grantRole(NFT_MINTER_ROLE, BACKEND_ADDRESS);
      await tx1.wait();
      console.log("   ✅ NFT_MINTER_ROLE granted to", BACKEND_ADDRESS);
    }
  } catch (error) {
    console.log("   ❌ Failed to grant NFT_MINTER_ROLE:", error.message);
  }

  // 2. Grant BACKEND_ROLE
  console.log("\n2️⃣  Granting BACKEND_ROLE...");
  try {
    const hasBackendRole = await roleManager.hasRole(BACKEND_ROLE, BACKEND_ADDRESS);
    if (hasBackendRole) {
      console.log("   ℹ️  Address already has BACKEND_ROLE");
    } else {
      const tx2 = await roleManager.grantRole(BACKEND_ROLE, BACKEND_ADDRESS);
      await tx2.wait();
      console.log("   ✅ BACKEND_ROLE granted to", BACKEND_ADDRESS);
    }
  } catch (error) {
    console.log("   ❌ Failed to grant BACKEND_ROLE:", error.message);
  }

  // 3. Verify all roles
  console.log("\n3️⃣  Verifying Role Assignments...");
  
  const roles = [
    { name: "NFT_MINTER_ROLE", hash: NFT_MINTER_ROLE },
    { name: "BACKEND_ROLE", hash: BACKEND_ROLE }
  ];

  for (const role of roles) {
    const hasRole = await roleManager.hasRole(role.hash, BACKEND_ADDRESS);
    console.log(`   ${hasRole ? '✅' : '❌'} ${role.name}: ${hasRole}`);
  }

  // Also check admin roles
  console.log("\n4️⃣  Admin Role Status:");
  const REGISTRY_ADMIN_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("REGISTRY_ADMIN_ROLE"));
  const CONTRACT_DEPLOYER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("CONTRACT_DEPLOYER_ROLE"));
  
  const adminRoles = [
    { name: "REGISTRY_ADMIN_ROLE", hash: REGISTRY_ADMIN_ROLE },
    { name: "CONTRACT_DEPLOYER_ROLE", hash: CONTRACT_DEPLOYER_ROLE }
  ];

  for (const role of adminRoles) {
    const hasRole = await roleManager.hasRole(role.hash, adminWallet.address);
    console.log(`   ${hasRole ? '✅' : '❌'} Admin has ${role.name}: ${hasRole}`);
  }

  console.log("\n" + "=" .repeat(60));
  console.log("✅ Role Configuration Complete!");
  console.log("=" .repeat(60));
  console.log("\n📝 Summary:");
  console.log(`  Backend/Minter Address: ${BACKEND_ADDRESS}`);
  console.log("  Granted Roles:");
  console.log("    - NFT_MINTER_ROLE (can mint NFTs)");
  console.log("    - BACKEND_ROLE (backend service access)");
  console.log("\n🔒 Security Notes:");
  console.log("  - This address can now mint NFTs");
  console.log("  - Keep the private key secure");
  console.log("  - Monitor minting activity");
  console.log("  - Set up rate limiting in your API");
  console.log("=" .repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });