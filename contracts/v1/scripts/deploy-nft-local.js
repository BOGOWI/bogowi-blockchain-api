const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🎨 Starting NFT Infrastructure Local Deployment...\n");
  console.log("=" .repeat(60));

  // Verify we're on local network
  const network = hre.network.name;
  if (network !== "hardhat" && network !== "localhost") {
    console.error("❌ This script is for local deployment only!");
    console.error("   Current network:", network);
    console.error("   Run with: npx hardhat run scripts/deploy-nft-local.js --network localhost");
    process.exit(1);
  }

  // Get signers
  const [deployer, admin, minter, backend, user1, user2] = await hre.ethers.getSigners();
  
  console.log("📍 Network:", network);
  console.log("👤 Deployer:", deployer.address);
  console.log("👤 Admin:", admin.address);
  console.log("👤 Minter:", minter.address);
  console.log("👤 Backend:", backend.address);
  
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("💰 Deployer balance:", hre.ethers.formatEther(balance), "CAM (simulated)\n");

  // Check if core contracts are already deployed
  let roleManagerAddress;
  let bogoTokenAddress;
  
  const existingDeploymentPath = path.join(__dirname, `deployment-${network}.json`);
  if (fs.existsSync(existingDeploymentPath)) {
    console.log("📂 Found existing deployment, reusing core contracts...");
    const existingDeployment = JSON.parse(fs.readFileSync(existingDeploymentPath, 'utf8'));
    roleManagerAddress = existingDeployment.contracts.RoleManager;
    bogoTokenAddress = existingDeployment.contracts.BOGOToken;
    console.log("✅ Using existing RoleManager:", roleManagerAddress);
    console.log("✅ Using existing BOGOToken:", bogoTokenAddress);
  } else {
    console.log("🚀 No existing deployment found, deploying core contracts first...\n");
    
    // Deploy RoleManager
    console.log("1. Deploying RoleManager...");
    const RoleManager = await hre.ethers.getContractFactory("RoleManager");
    const roleManager = await RoleManager.deploy();
    await roleManager.waitForDeployment();
    roleManagerAddress = await roleManager.getAddress();
    console.log("✅ RoleManager deployed to:", roleManagerAddress);
    
    // Grant admin role
    const DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
    await roleManager.grantRole(DEFAULT_ADMIN_ROLE, admin.address);
    console.log("✅ Admin role granted to:", admin.address);
    
    // Deploy BOGOToken
    console.log("\n2. Deploying BOGOToken...");
    const BOGOToken = await hre.ethers.getContractFactory("BOGOToken");
    const bogoToken = await BOGOToken.deploy(
      roleManagerAddress,
      "BOGOWI Token",
      "BOGO"
    );
    await bogoToken.waitForDeployment();
    bogoTokenAddress = await bogoToken.getAddress();
    console.log("✅ BOGOToken deployed to:", bogoTokenAddress);
    
    // Register BOGOToken with RoleManager
    await roleManager.registerContract(bogoTokenAddress, "BOGOToken");
    console.log("✅ BOGOToken registered with RoleManager");
  }

  console.log("\n" + "=" .repeat(60));
  console.log("🎫 Deploying NFT Infrastructure...");
  console.log("=" .repeat(60) + "\n");

  // Deploy NFTRegistry
  console.log("3. Deploying NFTRegistry...");
  const NFTRegistry = await hre.ethers.getContractFactory("NFTRegistry");
  const nftRegistry = await NFTRegistry.deploy(roleManagerAddress);
  await nftRegistry.waitForDeployment();
  const nftRegistryAddress = await nftRegistry.getAddress();
  console.log("✅ NFTRegistry deployed to:", nftRegistryAddress);

  // Deploy BOGOWITickets
  console.log("\n4. Deploying BOGOWITickets...");
  const BOGOWITickets = await hre.ethers.getContractFactory("BOGOWITickets");
  const tickets = await BOGOWITickets.deploy(
    roleManagerAddress,
    admin.address // Conservation DAO address
  );
  await tickets.waitForDeployment();
  const ticketsAddress = await tickets.getAddress();
  console.log("✅ BOGOWITickets deployed to:", ticketsAddress);

  console.log("\n" + "=" .repeat(60));
  console.log("⚙️  Configuring Contracts...");
  console.log("=" .repeat(60) + "\n");

  // Get RoleManager instance
  const roleManager = await hre.ethers.getContractAt("RoleManager", roleManagerAddress);

  // Register contracts with RoleManager
  console.log("5. Registering contracts with RoleManager...");
  await roleManager.registerContract(nftRegistryAddress, "NFTRegistry");
  console.log("✅ NFTRegistry registered");
  await roleManager.registerContract(ticketsAddress, "BOGOWITickets");
  console.log("✅ BOGOWITickets registered");

  // Setup roles
  console.log("\n6. Setting up roles...");
  const REGISTRY_ADMIN_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("REGISTRY_ADMIN_ROLE"));
  const CONTRACT_DEPLOYER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("CONTRACT_DEPLOYER_ROLE"));
  const NFT_MINTER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
  const BACKEND_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("BACKEND_ROLE"));
  const ADMIN_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("ADMIN_ROLE"));
  const PAUSER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("PAUSER_ROLE"));

  // Grant roles
  await roleManager.grantRole(REGISTRY_ADMIN_ROLE, admin.address);
  await roleManager.grantRole(CONTRACT_DEPLOYER_ROLE, admin.address);
  await roleManager.grantRole(NFT_MINTER_ROLE, minter.address);
  await roleManager.grantRole(BACKEND_ROLE, backend.address);
  await roleManager.grantRole(ADMIN_ROLE, admin.address);
  await roleManager.grantRole(PAUSER_ROLE, admin.address);
  
  console.log("✅ Roles granted:");
  console.log("   - REGISTRY_ADMIN_ROLE -> admin");
  console.log("   - CONTRACT_DEPLOYER_ROLE -> admin");
  console.log("   - NFT_MINTER_ROLE -> minter");
  console.log("   - BACKEND_ROLE -> backend");
  console.log("   - ADMIN_ROLE -> admin");
  console.log("   - PAUSER_ROLE -> admin");

  // Register BOGOWITickets with NFTRegistry
  console.log("\n7. Registering BOGOWITickets with NFTRegistry...");
  await nftRegistry.connect(admin).registerContract(
    ticketsAddress,
    0, // ContractType.TICKET
    "BOGOWI Event Tickets",
    "1.0.0"
  );
  console.log("✅ BOGOWITickets registered in NFTRegistry");

  // Verify deployment
  console.log("\n8. Verifying deployment...");
  const isTicketsRegistered = await nftRegistry.isRegistered(ticketsAddress);
  const ticketsInfo = await nftRegistry.getContractInfo(ticketsAddress);
  console.log(`✅ BOGOWITickets registration verified: ${isTicketsRegistered}`);
  console.log(`   Contract Type: ${["TICKET", "COLLECTIBLE", "BADGE"][ticketsInfo.contractType]}`);
  console.log(`   Active: ${ticketsInfo.isActive}`);
  
  const totalContracts = await nftRegistry.getContractCount();
  console.log(`✅ Total contracts in registry: ${totalContracts}`);
  
  // Save deployment info
  const deploymentInfo = {
    network: network,
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    contracts: {
      RoleManager: roleManagerAddress,
      BOGOToken: bogoTokenAddress,
      NFTRegistry: nftRegistryAddress,
      BOGOWITickets: ticketsAddress
    },
    roles: {
      admin: admin.address,
      contractDeployer: admin.address,  // Added for clarity
      minter: minter.address,
      backend: backend.address
    },
    testUsers: {
      user1: user1.address,
      user2: user2.address
    },
    chainId: 501  // Document the chain ID being used
  };

  const nftDeploymentPath = path.join(__dirname, `deployment-nft-${network}.json`);
  fs.writeFileSync(nftDeploymentPath, JSON.stringify(deploymentInfo, null, 2));
  
  console.log("\n" + "=" .repeat(60));
  console.log("📁 Deployment info saved to:", nftDeploymentPath);
  console.log("=" .repeat(60));

  // Print summary
  console.log("\n🎉 NFT INFRASTRUCTURE DEPLOYMENT COMPLETE!");
  console.log("=" .repeat(60));
  console.log("📋 Network Configuration:");
  console.log("  Network:", network);
  console.log("  Chain ID:", 501, "(Camino Testnet)");
  console.log("\n📍 Contract Addresses:");
  console.log("  RoleManager:", roleManagerAddress);
  console.log("  BOGOToken:", bogoTokenAddress);
  console.log("  NFTRegistry:", nftRegistryAddress);
  console.log("  BOGOWITickets:", ticketsAddress);
  console.log("\n👥 Role Assignments:");
  console.log("  Admin/Contract Deployer:", admin.address);
  console.log("  NFT Minter:", minter.address);
  console.log("  Backend:", backend.address);
  console.log("\n📊 Registry Status:");
  console.log("  Total Contracts:", totalContracts);
  console.log("  BOGOWITickets Active:", ticketsInfo.isActive);
  console.log("=" .repeat(60));
  
  console.log("\n📝 Next Steps:");
  console.log("1. Run verification: npm run verify-nft-local");
  console.log("2. Test minting: npm run test-mint-local");
  console.log("3. Test registry: npm run test-registry-local");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });