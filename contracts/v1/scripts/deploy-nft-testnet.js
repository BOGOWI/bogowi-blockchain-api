const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🎨 Deploying NFT Infrastructure to TESTNET...\n");
  console.log("=" .repeat(60));

  // Verify we're on testnet
  const network = hre.network.name;
  const chainId = await hre.ethers.provider.getNetwork().then(n => n.chainId);
  
  if (chainId !== 501n) {
    console.error("❌ This script is for Camino Columbus testnet only!");
    console.error("   Current chain ID:", chainId);
    console.error("   Expected: 501 (Columbus)");
    process.exit(1);
  }

  // Get deployer
  const [deployer] = await hre.ethers.getSigners();
  console.log("📍 Network:", network);
  console.log("🔗 Chain ID:", chainId);
  console.log("👤 Deployer:", deployer.address);
  
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("💰 Deployer balance:", hre.ethers.formatEther(balance), "CAM\n");

  // Load EXISTING testnet deployment
  const coreDeploymentPath = path.join(__dirname, "deployment-columbus.json");
  if (!fs.existsSync(coreDeploymentPath)) {
    console.error("❌ Core contracts deployment not found!");
    console.error("   Expected file:", coreDeploymentPath);
    console.error("   Run core deployment first!");
    process.exit(1);
  }

  const coreDeployment = JSON.parse(fs.readFileSync(coreDeploymentPath, 'utf8'));
  const roleManagerAddress = coreDeployment.contracts.RoleManager;
  const bogoTokenAddress = coreDeployment.contracts.BOGOToken;
  const adminAddress = coreDeployment.adminAddress;

  console.log("📂 Using EXISTING core contracts:");
  console.log("  RoleManager:", roleManagerAddress);
  console.log("  BOGOToken:", bogoTokenAddress);
  console.log("  Admin:", adminAddress);

  console.log("\n" + "=" .repeat(60));
  console.log("🎫 Deploying NEW NFT Contracts...");
  console.log("=" .repeat(60) + "\n");

  // Deploy NFTRegistry
  console.log("1. Deploying NFTRegistry...");
  const NFTRegistry = await hre.ethers.getContractFactory("NFTRegistry");
  const nftRegistry = await NFTRegistry.deploy(roleManagerAddress);
  await nftRegistry.waitForDeployment();
  const nftRegistryAddress = await nftRegistry.getAddress();
  console.log("✅ NFTRegistry deployed to:", nftRegistryAddress);

  // Deploy BOGOWITickets
  console.log("\n2. Deploying BOGOWITickets...");
  const BOGOWITickets = await hre.ethers.getContractFactory("BOGOWITickets");
  const tickets = await BOGOWITickets.deploy(
    roleManagerAddress,
    adminAddress // Conservation DAO (using admin for testnet)
  );
  await tickets.waitForDeployment();
  const ticketsAddress = await tickets.getAddress();
  console.log("✅ BOGOWITickets deployed to:", ticketsAddress);

  console.log("\n" + "=" .repeat(60));
  console.log("⚙️  Configuring Contracts...");
  console.log("=" .repeat(60) + "\n");

  // Get RoleManager instance
  const roleManager = await hre.ethers.getContractAt("RoleManager", roleManagerAddress);

  // Check if deployer has admin role
  const DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
  const deployerIsAdmin = await roleManager.hasRole(DEFAULT_ADMIN_ROLE, deployer.address);
  
  if (!deployerIsAdmin) {
    console.error("❌ Deployer doesn't have admin role in RoleManager!");
    console.error("   Contact admin to grant roles and register contracts");
    console.error("   Admin address:", adminAddress);
    
    // Save deployment info anyway
    const deploymentInfo = {
      network: network,
      chainId: Number(chainId),
      deployer: deployer.address,
      timestamp: new Date().toISOString(),
      contracts: {
        RoleManager: roleManagerAddress,
        BOGOToken: bogoTokenAddress,
        NFTRegistry: nftRegistryAddress,
        BOGOWITickets: ticketsAddress
      },
      adminAddress: adminAddress,
      status: "DEPLOYED_NOT_CONFIGURED"
    };

    const nftDeploymentPath = path.join(__dirname, `deployment-nft-${network}.json`);
    fs.writeFileSync(nftDeploymentPath, JSON.stringify(deploymentInfo, null, 2));
    
    console.log("\n⚠️  Contracts deployed but not configured!");
    console.log("📁 Deployment info saved to:", nftDeploymentPath);
    console.log("\nManual steps required:");
    console.log("1. Admin must register contracts with RoleManager");
    console.log("2. Admin must grant necessary roles");
    console.log("3. Admin must register BOGOWITickets with NFTRegistry");
    process.exit(0);
  }

  // If deployer is admin, continue with configuration
  console.log("3. Registering contracts with RoleManager...");
  try {
    await roleManager.registerContract(nftRegistryAddress, "NFTRegistry");
    console.log("✅ NFTRegistry registered");
  } catch (error) {
    console.log("⚠️  NFTRegistry registration failed:", error.message);
  }

  try {
    await roleManager.registerContract(ticketsAddress, "BOGOWITickets");
    console.log("✅ BOGOWITickets registered");
  } catch (error) {
    console.log("⚠️  BOGOWITickets registration failed:", error.message);
  }

  // Setup roles
  console.log("\n4. Setting up roles...");
  const REGISTRY_ADMIN_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("REGISTRY_ADMIN_ROLE"));
  const CONTRACT_DEPLOYER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("CONTRACT_DEPLOYER_ROLE"));
  const NFT_MINTER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
  const BACKEND_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("BACKEND_ROLE"));

  // Grant roles to admin
  try {
    await roleManager.grantRole(REGISTRY_ADMIN_ROLE, adminAddress);
    await roleManager.grantRole(CONTRACT_DEPLOYER_ROLE, adminAddress);
    console.log("✅ Registry roles granted to admin");
  } catch (error) {
    console.log("⚠️  Role granting failed:", error.message);
  }

  // Register BOGOWITickets with NFTRegistry
  console.log("\n5. Registering BOGOWITickets with NFTRegistry...");
  try {
    // Need CONTRACT_DEPLOYER_ROLE for this
    const hasDeployerRole = await roleManager.hasRole(CONTRACT_DEPLOYER_ROLE, deployer.address);
    if (!hasDeployerRole) {
      await roleManager.grantRole(CONTRACT_DEPLOYER_ROLE, deployer.address);
    }
    
    await nftRegistry.registerContract(
      ticketsAddress,
      0, // ContractType.TICKET
      "BOGOWI Event Tickets",
      "1.0.0"
    );
    console.log("✅ BOGOWITickets registered in NFTRegistry");
  } catch (error) {
    console.log("⚠️  NFTRegistry registration failed:", error.message);
  }

  // Save deployment info
  const deploymentInfo = {
    network: network,
    chainId: Number(chainId),
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    contracts: {
      RoleManager: roleManagerAddress,
      BOGOToken: bogoTokenAddress,
      NFTRegistry: nftRegistryAddress,
      BOGOWITickets: ticketsAddress
    },
    adminAddress: adminAddress,
    status: "DEPLOYED_AND_CONFIGURED"
  };

  const nftDeploymentPath = path.join(__dirname, `deployment-nft-${network}.json`);
  fs.writeFileSync(nftDeploymentPath, JSON.stringify(deploymentInfo, null, 2));
  
  console.log("\n" + "=" .repeat(60));
  console.log("📁 Deployment info saved to:", nftDeploymentPath);
  console.log("=" .repeat(60));

  // Print summary
  console.log("\n🎉 NFT INFRASTRUCTURE DEPLOYMENT COMPLETE!");
  console.log("=" .repeat(60));
  console.log("Existing Contracts (reused):");
  console.log("  RoleManager:", roleManagerAddress);
  console.log("  BOGOToken:", bogoTokenAddress);
  console.log("\nNew Contracts (deployed):");
  console.log("  NFTRegistry:", nftRegistryAddress);
  console.log("  BOGOWITickets:", ticketsAddress);
  console.log("\nAdmin:", adminAddress);
  console.log("=" .repeat(60));
  
  console.log("\n📝 Next Steps:");
  console.log("1. Verify contracts on explorer");
  console.log("2. Grant NFT_MINTER_ROLE to authorized minters");
  console.log("3. Grant BACKEND_ROLE to backend service");
  console.log("4. Test minting operations");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });