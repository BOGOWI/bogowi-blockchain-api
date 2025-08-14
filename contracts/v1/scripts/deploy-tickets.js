const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🎫 Deploying BOGOWI Tickets NFT Contract...\n");

  // Get the network we're deploying to
  const network = hre.network.name;
  const chainId = await hre.ethers.provider.getNetwork().then(n => n.chainId);
  
  console.log(`📍 Network: ${network}`);
  console.log(`🔗 Chain ID: ${chainId}`);
  
  // Validate we're on Camino
  if (chainId !== 500n && chainId !== 501n) {
    console.error("❌ Error: Must deploy on Camino network (chainId 500 or 501)");
    process.exit(1);
  }
  
  const isMainnet = chainId === 500n;
  console.log(`🌐 Environment: ${isMainnet ? 'MAINNET' : 'TESTNET'}\n`);

  // Get signers
  const [deployer] = await hre.ethers.getSigners();
  console.log(`👤 Deployer address: ${deployer.address}`);
  
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log(`💰 Deployer balance: ${hre.ethers.formatEther(balance)} CAM\n`);

  // Configuration
  const config = {
    roleManager: process.env.ROLE_MANAGER_ADDRESS || "",
    conservationDAO: process.env.CONSERVATION_DAO_ADDRESS || deployer.address,
    datakyteAPIKey: isMainnet 
      ? "dk_e2aad94de12a2a7e7865a70b369e1eab69e2b5e2896577a5fbcbbb50d709bd3d"
      : "dk_d707e26c919e72ab2bb3b81897566c393f4e2eba54d07ff680d765ee03d6cc5d"
  };

  // Deploy RoleManager if not provided
  let roleManagerAddress = config.roleManager;
  if (!roleManagerAddress) {
    console.log("📝 Deploying RoleManager...");
    const RoleManager = await hre.ethers.getContractFactory("RoleManager");
    const roleManager = await RoleManager.deploy();
    await roleManager.waitForDeployment();
    roleManagerAddress = await roleManager.getAddress();
    console.log(`✅ RoleManager deployed to: ${roleManagerAddress}\n`);
  } else {
    console.log(`📌 Using existing RoleManager: ${roleManagerAddress}\n`);
  }

  // Deploy BOGOWITickets (with Datakyte integration)
  console.log("📝 Deploying BOGOWITickets...");
  const BOGOWITickets = await hre.ethers.getContractFactory("BOGOWITickets");
  const tickets = await BOGOWITickets.deploy(
    roleManagerAddress,
    config.conservationDAO
  );
  await tickets.waitForDeployment();
  const ticketsAddress = await tickets.getAddress();
  console.log(`✅ BOGOWITickets deployed to: ${ticketsAddress}\n`);

  // Register with RoleManager
  if (!config.roleManager) {
    console.log("📝 Registering contract with RoleManager...");
    const roleManager = await hre.ethers.getContractAt("RoleManager", roleManagerAddress);
    await roleManager.registerContract(ticketsAddress, "BOGOWITickets");
    console.log("✅ Contract registered\n");

    // Grant roles
    console.log("📝 Setting up roles...");
    const NFT_MINTER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
    const ADMIN_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("ADMIN_ROLE"));
    const PAUSER_ROLE = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("PAUSER_ROLE"));
    
    await roleManager.grantRole(ADMIN_ROLE, deployer.address);
    await roleManager.grantRole(NFT_MINTER_ROLE, deployer.address);
    await roleManager.grantRole(PAUSER_ROLE, deployer.address);
    console.log("✅ Roles granted to deployer\n");
  }

  // Save deployment info
  const deploymentInfo = {
    network: network,
    chainId: chainId.toString(),
    deployedAt: new Date().toISOString(),
    contracts: {
      roleManager: roleManagerAddress,
      tickets: ticketsAddress,
      conservationDAO: config.conservationDAO
    },
    datakyte: {
      apiKey: config.datakyteAPIKey,
      metadataBaseURL: `https://dklnk.to/api/nfts/${ticketsAddress}/{tokenId}/metadata`
    },
    deployer: deployer.address,
    blockNumber: await hre.ethers.provider.getBlockNumber()
  };

  // Create deployments directory if it doesn't exist
  const deploymentsDir = path.join(__dirname, '..', 'deployments');
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }

  // Save deployment file
  const filename = `tickets-${network}-${Date.now()}.json`;
  const filepath = path.join(deploymentsDir, filename);
  fs.writeFileSync(filepath, JSON.stringify(deploymentInfo, null, 2));
  console.log(`💾 Deployment info saved to: ${filepath}\n`);

  // Print summary
  console.log("=" .repeat(60));
  console.log("🎉 DEPLOYMENT SUCCESSFUL!");
  console.log("=" .repeat(60));
  console.log("\n📋 Summary:");
  console.log(`  • Network: ${network} (Chain ID: ${chainId})`);
  console.log(`  • RoleManager: ${roleManagerAddress}`);
  console.log(`  • BOGOWITickets: ${ticketsAddress}`);
  console.log(`  • Conservation DAO: ${config.conservationDAO}`);
  console.log("\n🔗 Datakyte Integration:");
  console.log(`  • API Key: ${config.datakyteAPIKey.substring(0, 10)}...`);
  console.log(`  • Metadata URL Pattern:`);
  console.log(`    https://dklnk.to/api/nfts/${ticketsAddress}/{tokenId}/metadata`);
  console.log("\n📝 Next Steps:");
  console.log("  1. Update .env with contract addresses");
  console.log("  2. Configure backend API with contract address");
  console.log("  3. Test minting a ticket");
  console.log("  4. Verify metadata on Datakyte");
  console.log("=" .repeat(60));

  // Verify contracts on explorer (if available)
  if (network !== "hardhat" && network !== "localhost") {
    console.log("\n🔍 Verifying contracts on explorer...");
    try {
      await hre.run("verify:verify", {
        address: ticketsAddress,
        constructorArguments: [roleManagerAddress, config.conservationDAO],
      });
      console.log("✅ Contract verified on explorer");
    } catch (error) {
      console.log("⚠️  Contract verification failed:", error.message);
    }
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });