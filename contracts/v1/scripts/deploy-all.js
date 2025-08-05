const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🚀 Starting BOGOWI V1 deployment...\n");

  // Get deployer account
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with account:", deployer.address);
  
  const balance = await hre.ethers.provider.getBalance(deployer.address);
  console.log("Account balance:", hre.ethers.formatEther(balance), "CAM\n");

  // Check required environment variables
  // For testnet, use the testnet admin address
  let adminAddress = process.env.ADMIN_ADDRESS;
  
  // If on testnet network, use testnet admin
  if (hre.network.name === "columbus" || hre.network.name === "testnet") {
    adminAddress = "0xB34A822F735CDE477cbB39a06118267D00948ef7"; // testnet admin
  } else if (hre.network.name === "camino" || hre.network.name === "mainnet") {
    adminAddress = "0x444ddA4cA50765D3c0c0c662aAecF3b5D49761Ea"; // mainnet admin
  }
  
  if (!adminAddress || adminAddress === "YOUR_ADMIN_ADDRESS") {
    throw new Error("Admin address not configured for network: " + hre.network.name);
  }
  
  console.log("Admin address:", adminAddress);

  // Deploy RoleManager
  console.log("1. Deploying RoleManager...");
  const RoleManager = await hre.ethers.getContractFactory("RoleManager");
  const roleManager = await RoleManager.deploy(); // No parameters - constructor grants admin to deployer
  await roleManager.waitForDeployment();
  const roleManagerAddress = await roleManager.getAddress();
  console.log("✅ RoleManager deployed to::", roleManagerAddress);
  
  // Grant admin role to the specified admin if different from deployer
  if (adminAddress.toLowerCase() !== deployer.address.toLowerCase()) {
    console.log("Granting admin role to:", adminAddress);
    const DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
    const tx = await roleManager.grantRole(DEFAULT_ADMIN_ROLE, adminAddress);
    await tx.wait();
    console.log("✅ Admin role granted to:", adminAddress);
  }

  // Deploy BOGOToken
  console.log("\n2. Deploying BOGOToken...");
  const BOGOToken = await hre.ethers.getContractFactory("BOGOToken");
  const bogoToken = await BOGOToken.deploy(
    roleManagerAddress,
    "BOGOWI Token",  // Token name
    "BOGO"          // Token symbol
  );
  await bogoToken.waitForDeployment();
  const bogoTokenAddress = await bogoToken.getAddress();
  console.log("✅ BOGOToken deployed to:", bogoTokenAddress);

  // Deploy BOGORewardDistributor
  console.log("\n3. Deploying BOGORewardDistributor...");
  const BOGORewardDistributor = await hre.ethers.getContractFactory("BOGORewardDistributor");
  const rewardDistributor = await BOGORewardDistributor.deploy(
    roleManagerAddress,
    bogoTokenAddress
  );
  await rewardDistributor.waitForDeployment();
  const rewardDistributorAddress = await rewardDistributor.getAddress();
  console.log("✅ BOGORewardDistributor deployed to:", rewardDistributorAddress);

  // Save deployment info
  const deploymentInfo = {
    network: hre.network.name,
    deployer: deployer.address,
    timestamp: new Date().toISOString(),
    contracts: {
      RoleManager: roleManagerAddress,
      BOGOToken: bogoTokenAddress,
      BOGORewardDistributor: rewardDistributorAddress
    },
    adminAddress: adminAddress
  };

  const deploymentPath = path.join(__dirname, `deployment-${hre.network.name}.json`);
  fs.writeFileSync(deploymentPath, JSON.stringify(deploymentInfo, null, 2));
  console.log("\n📁 Deployment info saved to:", deploymentPath);

  // Print summary
  console.log("\n" + "=".repeat(50));
  console.log("DEPLOYMENT SUMMARY");
  console.log("=".repeat(50));
  console.log(`Network: ${hre.network.name}`);
  console.log(`RoleManager: ${roleManagerAddress}`);
  console.log(`BOGOToken: ${bogoTokenAddress}`);
  console.log(`RewardDistributor: ${rewardDistributorAddress}`);
  console.log("=".repeat(50));

  console.log("\n⚠️  IMPORTANT: Update your .env file with these addresses!");
  console.log("ROLE_MANAGER_ADDRESS=" + roleManagerAddress);
  console.log("BOGO_TOKEN_ADDRESS=" + bogoTokenAddress);
  console.log("REWARD_DISTRIBUTOR_ADDRESS=" + rewardDistributorAddress);

  console.log("\n✅ Deployment complete! Next steps:");
  console.log("1. Run 'npm run setup-roles' to configure roles");
  console.log("2. Run 'npm run mint-supply' to mint initial tokens");
  console.log("3. Run 'npm run verify-deployment' to verify contracts");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });