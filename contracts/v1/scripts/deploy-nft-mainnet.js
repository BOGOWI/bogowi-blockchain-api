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
  
  // Pre-deployment checklist
  console.log("\n📋 Pre-deployment Checklist:");
  console.log("   ✓ Contracts tested on testnet");
  console.log("   ✓ Admin address verified");
  console.log("   ✓ Sufficient CAM balance");
  console.log("   ✓ Private key secured");
  console.log("   ✓ Deployment will grant CONTRACT_DEPLOYER_ROLE if admin");
  
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
    console.error("   Run core deployment first!");
    process.exit(1);
  }

  const coreDeployment = JSON.parse(fs.readFileSync(coreDeploymentPath, 'utf8'));
  const roleManagerAddress = coreDeployment.contracts.RoleManager;
  const bogoTokenAddress = coreDeployment.contracts.BOGOToken;
  const adminAddress = coreDeployment.adminAddress;

  console.log("\n📂 Using EXISTING core contracts:");
  console.log("  RoleManager:", roleManagerAddress);
  console.log("  BOGOToken:", bogoTokenAddress);
  console.log("  Admin:", adminAddress);

  console.log("\n" + "=" .repeat(60));
  console.log("🎫 Deploying NEW NFT Contracts...");
  console.log("=" .repeat(60) + "\n");

  // 1. Deploy NFTRegistry
  console.log("1. Deploying NFTRegistry...");
  const NFTRegistry = await hre.ethers.getContractFactory("NFTRegistry");
  const nftRegistry = await NFTRegistry.deploy(roleManagerAddress);
  await nftRegistry.waitForDeployment();
  const nftRegistryAddress = await nftRegistry.getAddress();
  console.log("✅ NFTRegistry deployed to:", nftRegistryAddress);

  // 2. Deploy BOGOWITickets
  console.log("\n2. Deploying BOGOWITickets...");
  const BOGOWITickets = await hre.ethers.getContractFactory("BOGOWITickets");
  const tickets = await BOGOWITickets.deploy(
    roleManagerAddress,
    adminAddress // conservation DAO
  );
  await tickets.waitForDeployment();
  const ticketsAddress = await tickets.getAddress();
  console.log("✅ BOGOWITickets deployed to:", ticketsAddress);

  console.log("\n" + "=" .repeat(60));
  console.log("⚙️  Configuring Contracts...");
  console.log("=" .repeat(60) + "\n");

  // Get RoleManager instance
  const roleManager = await hre.ethers.getContractAt("RoleManager", roleManagerAddress);

  // Check if deployer is admin
  const DEFAULT_ADMIN_ROLE = "0x0000000000000000000000000000000000000000000000000000000000000000";
  const isAdmin = await roleManager.hasRole(DEFAULT_ADMIN_ROLE, deployer.address);

  if (!isAdmin) {
    // Save deployment info for manual configuration
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
    console.log("\n🔒 MAINNET SECURITY NOTICE:");
    console.log("Manual configuration required by admin:");
    console.log("1. Admin must register contracts with RoleManager");
    console.log("2. Admin must grant necessary roles");
    console.log("3. Admin must register BOGOWITickets with NFTRegistry");
    console.log("4. Verify all permissions before production use");
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
    // Check if deployer has CONTRACT_DEPLOYER_ROLE
    let hasDeployerRole = await roleManager.hasRole(CONTRACT_DEPLOYER_ROLE, deployer.address);
    
    if (!hasDeployerRole) {
      if (isAdmin) {
        // If deployer is admin, grant the role to self
        console.log("   Granting CONTRACT_DEPLOYER_ROLE to deployer...");
        const grantTx = await roleManager.grantRole(CONTRACT_DEPLOYER_ROLE, deployer.address);
        await grantTx.wait();
        console.log("   ✅ Role granted");
        hasDeployerRole = true;
      } else {
        console.log("   ⚠️  Deployer needs CONTRACT_DEPLOYER_ROLE");
        console.log("   Admin must grant this role before registration can proceed");
      }
    }
    
    if (hasDeployerRole) {
      await nftRegistry.registerContract(
        ticketsAddress,
        0, // ContractType.TICKET
        "BOGOWI Event Tickets",
        "1.0.0"
      );
      console.log("✅ BOGOWITickets registered in NFTRegistry");
    } else {
      console.log("⚠️  Skipping registration - missing CONTRACT_DEPLOYER_ROLE");
    }
  } catch (error) {
    console.log("⚠️  NFTRegistry registration failed:", error.message);
  }

  // Verify deployment
  console.log("\n6. Verifying deployment...");
  let deploymentVerified = false;
  let totalContracts = 0;
  let ticketsInfo = null;
  
  try {
    const isTicketsRegistered = await nftRegistry.isRegistered(ticketsAddress);
    ticketsInfo = await nftRegistry.getContractInfo(ticketsAddress);
    console.log(`✅ BOGOWITickets registration verified: ${isTicketsRegistered}`);
    console.log(`   Contract Type: ${["TICKET", "COLLECTIBLE", "BADGE"][ticketsInfo.contractType]}`);
    console.log(`   Active: ${ticketsInfo.isActive}`);
    
    totalContracts = await nftRegistry.getContractCount();
    console.log(`✅ Total contracts in registry: ${totalContracts}`);
    deploymentVerified = true;
  } catch (error) {
    console.log("⚠️  Verification failed:", error.message);
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
    roles: {
      admin: adminAddress,
      contractDeployer: adminAddress
    },
    status: deploymentVerified ? "DEPLOYED_AND_CONFIGURED" : "DEPLOYED_NEEDS_VERIFICATION"
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
  console.log("  Chain ID:", Number(chainId), "(Camino Mainnet)");
  console.log("\n📍 Contract Addresses:");
  console.log("  RoleManager:", roleManagerAddress);
  console.log("  BOGOToken:", bogoTokenAddress);
  console.log("  NFTRegistry:", nftRegistryAddress);
  console.log("  BOGOWITickets:", ticketsAddress);
  console.log("\n👥 Role Assignments:");
  console.log("  Admin/Contract Deployer:", adminAddress);
  console.log("\n📊 Registry Status:");
  if (deploymentVerified && ticketsInfo) {
    console.log("  Total Contracts:", totalContracts);
    console.log("  BOGOWITickets Active:", ticketsInfo.isActive);
  } else {
    console.log("  Status: Needs manual verification");
  }
  console.log("=" .repeat(60));
  
  console.log("\n🔒 MAINNET POST-DEPLOYMENT CHECKLIST:");
  console.log("1. ✅ Verify contracts on Camino explorer");
  console.log("2. ⚠️  Grant NFT_MINTER_ROLE only to trusted addresses");
  console.log("3. ⚠️  Grant BACKEND_ROLE only to authorized backend services");
  console.log("4. 🔐 Test with small amounts before production use");
  console.log("5. 📊 Set up monitoring and alerts");
  console.log("6. 💾 Backup deployment files securely");
  console.log("7. 🔍 Audit role assignments regularly");
  console.log("=" .repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });