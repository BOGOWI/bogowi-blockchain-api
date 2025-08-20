const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🎨 Deploying NFT Infrastructure to MAINNET...\n");
  console.log("=" .repeat(60));
  console.log("⚠️  WARNING: THIS IS MAINNET DEPLOYMENT!");
  console.log("=" .repeat(60) + "\n");

  // Verify we're on mainnet
  const network = hre.network.name;
  const chainId = await hre.ethers.provider.getNetwork().then(n => n.chainId);
  
  if (chainId !== 500n) {
    console.error("❌ This script is for Camino MAINNET only!");
    console.error("   Current chain ID:", chainId);
    console.error("   Expected: 500 (Camino Mainnet)");
    process.exit(1);
  }

  // Safety confirmation
  console.log("🔴 MAINNET DEPLOYMENT CONFIRMATION");
  console.log("   Chain: Camino Mainnet (500)");
  console.log("   This will deploy REAL contracts with REAL CAM!");
  console.log("\n   Press Ctrl+C to cancel...");
  console.log("   Continuing in 10 seconds...\n");
  
  // Give time to cancel
  await new Promise(resolve => setTimeout(resolve, 10000));

  // Get deployer
  const [deployer] = await hre.ethers.getSigners();
  console.log("📍 Network:", network);
  console.log("🔗 Chain ID:", chainId);
  console.log("👤 Deployer:", deployer.address);
  
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("💰 Deployer balance:", hre.ethers.formatEther(balance), "CAM");
  
  // Minimum balance check
  const minBalance = hre.ethers.parseEther("5"); // Require at least 5 CAM
  if (balance < minBalance) {
    console.error("❌ Insufficient balance for deployment!");
    console.error("   Required: 5 CAM minimum");
    console.error("   Current:", hre.ethers.formatEther(balance), "CAM");
    process.exit(1);
  }

  // Load EXISTING mainnet deployment
  const coreDeploymentPath = path.join(__dirname, "deployment-camino.json");
  if (!fs.existsSync(coreDeploymentPath)) {
    console.error("❌ Core contracts deployment not found!");
    console.error("   Expected file:", coreDeploymentPath);
    console.error("   Core contracts must be deployed first!");
    process.exit(1);
  }

  const coreDeployment = JSON.parse(fs.readFileSync(coreDeploymentPath, 'utf8'));
  const roleManagerAddress = coreDeployment.contracts.RoleManager;
  const bogoTokenAddress = coreDeployment.contracts.BOGOToken;
  const adminAddress = coreDeployment.adminAddress;

  console.log("\n📂 Using EXISTING mainnet contracts:");
  console.log("  RoleManager:", roleManagerAddress);
  console.log("  BOGOToken:", bogoTokenAddress);
  console.log("  Admin:", adminAddress);

  // Conservation DAO address for mainnet
  const CONSERVATION_DAO = process.env.CONSERVATION_DAO_MAINNET || adminAddress;
  console.log("  Conservation DAO:", CONSERVATION_DAO);

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
  
  // Wait for confirmation
  console.log("   Waiting for block confirmations...");
  await nftRegistry.deploymentTransaction().wait(3);
  console.log("   Confirmed!");

  // Deploy BOGOWITickets
  console.log("\n2. Deploying BOGOWITickets...");
  const BOGOWITickets = await hre.ethers.getContractFactory("BOGOWITickets");
  const tickets = await BOGOWITickets.deploy(
    roleManagerAddress,
    CONSERVATION_DAO
  );
  await tickets.waitForDeployment();
  const ticketsAddress = await tickets.getAddress();
  console.log("✅ BOGOWITickets deployed to:", ticketsAddress);
  
  // Wait for confirmation
  console.log("   Waiting for block confirmations...");
  await tickets.deploymentTransaction().wait(3);
  console.log("   Confirmed!");

  console.log("\n" + "=" .repeat(60));
  console.log("⚙️  Configuration Status");
  console.log("=" .repeat(60) + "\n");

  // Get RoleManager instance
  const roleManager = await hre.ethers.getContractAt("RoleManager", roleManagerAddress);

  // Check if deployer has admin role
  const DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
  const deployerIsAdmin = await roleManager.hasRole(DEFAULT_ADMIN_ROLE, deployer.address);
  
  if (!deployerIsAdmin) {
    console.log("⚠️  MANUAL CONFIGURATION REQUIRED!");
    console.log("   Deployer doesn't have admin role in RoleManager");
    console.log("   The following steps must be completed by admin:", adminAddress);
    console.log("\n   Required Actions:");
    console.log("   1. Register NFTRegistry with RoleManager");
    console.log("   2. Register BOGOWITickets with RoleManager");
    console.log("   3. Grant REGISTRY_ADMIN_ROLE to appropriate address");
    console.log("   4. Grant CONTRACT_DEPLOYER_ROLE to appropriate address");
    console.log("   5. Grant NFT_MINTER_ROLE to minting service");
    console.log("   6. Grant BACKEND_ROLE to backend service");
    console.log("   7. Register BOGOWITickets in NFTRegistry");
  } else {
    console.log("✅ Deployer has admin role - attempting configuration...");
    
    // Register contracts
    try {
      console.log("\n3. Registering contracts with RoleManager...");
      await roleManager.registerContract(nftRegistryAddress, "NFTRegistry");
      console.log("✅ NFTRegistry registered");
      await roleManager.registerContract(ticketsAddress, "BOGOWITickets");
      console.log("✅ BOGOWITickets registered");
    } catch (error) {
      console.log("⚠️  Registration failed:", error.message);
    }
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
    conservationDAO: CONSERVATION_DAO,
    adminAddress: adminAddress,
    status: deployerIsAdmin ? "DEPLOYED_PARTIAL_CONFIG" : "DEPLOYED_NOT_CONFIGURED",
    gasUsed: {
      NFTRegistry: (await nftRegistry.deploymentTransaction()).gasLimit.toString(),
      BOGOWITickets: (await tickets.deploymentTransaction()).gasLimit.toString()
    }
  };

  const nftDeploymentPath = path.join(__dirname, `deployment-nft-${network}.json`);
  fs.writeFileSync(nftDeploymentPath, JSON.stringify(deploymentInfo, null, 2));
  
  // Create verification script
  const verifyScript = `
// Verification commands for Camino Explorer
// Run these after deployment is confirmed

npx hardhat verify --network ${network} ${nftRegistryAddress} "${roleManagerAddress}"

npx hardhat verify --network ${network} ${ticketsAddress} "${roleManagerAddress}" "${CONSERVATION_DAO}"
`;

  const verifyPath = path.join(__dirname, `verify-nft-${network}.sh`);
  fs.writeFileSync(verifyPath, verifyScript);

  console.log("\n" + "=" .repeat(60));
  console.log("📁 Deployment info saved to:", nftDeploymentPath);
  console.log("📜 Verification script saved to:", verifyPath);
  console.log("=" .repeat(60));

  // Print summary
  console.log("\n🎉 MAINNET NFT DEPLOYMENT COMPLETE!");
  console.log("=" .repeat(60));
  console.log("Deployed Contracts:");
  console.log("  NFTRegistry:", nftRegistryAddress);
  console.log("  BOGOWITickets:", ticketsAddress);
  console.log("\nConfiguration:");
  console.log("  Conservation DAO:", CONSERVATION_DAO);
  console.log("  Admin:", adminAddress);
  console.log("  Status:", deployerIsAdmin ? "Partial" : "Manual Required");
  console.log("=" .repeat(60));
  
  console.log("\n📝 CRITICAL Next Steps:");
  console.log("1. Save deployment addresses immediately!");
  console.log("2. Verify contracts on explorer (run verify script)");
  console.log("3. Complete role configuration if needed");
  console.log("4. Test with small amounts first");
  console.log("5. Set up monitoring and alerts");
  console.log("\n⚠️  DO NOT SHARE PRIVATE KEYS OR DEPLOYMENT INFO!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });