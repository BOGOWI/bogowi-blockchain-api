const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("📝 Registering contracts with RoleManager...\n");

  // Load deployment info
  const deploymentPath = path.join(__dirname, `deployment-${hre.network.name}.json`);
  if (!fs.existsSync(deploymentPath)) {
    throw new Error(`No deployment found for network: ${hre.network.name}`);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  const [deployer] = await hre.ethers.getSigners();
  
  console.log("Registering with account:", deployer.address);

  // Get RoleManager instance
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);

  // Check if admin
  const DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
  const isAdmin = await roleManager.hasRole(DEFAULT_ADMIN_ROLE, deployer.address);
  
  if (!isAdmin) {
    throw new Error("Deployer must have DEFAULT_ADMIN_ROLE to register contracts!");
  }

  // Register BOGOToken
  console.log("1. Registering BOGOToken...");
  let tx = await roleManager.registerContract(deployment.contracts.BOGOToken, "BOGOToken");
  await tx.wait();
  console.log("✅ BOGOToken registered");

  // Register BOGORewardDistributor
  console.log("\n2. Registering BOGORewardDistributor...");
  tx = await roleManager.registerContract(deployment.contracts.BOGORewardDistributor, "BOGORewardDistributor");
  await tx.wait();
  console.log("✅ BOGORewardDistributor registered");

  // Verify registrations
  console.log("\n🔍 Verifying registrations...");
  
  const isTokenRegistered = await roleManager.isRegistered(deployment.contracts.BOGOToken);
  const isDistributorRegistered = await roleManager.isRegistered(deployment.contracts.BOGORewardDistributor);
  
  console.log(`BOGOToken registered: ${isTokenRegistered ? "✅" : "❌"}`);
  console.log(`BOGORewardDistributor registered: ${isDistributorRegistered ? "✅" : "❌"}`);

  if (!isTokenRegistered || !isDistributorRegistered) {
    throw new Error("Contract registration failed!");
  }

  console.log("\n✅ All contracts registered successfully!");
  console.log("\n⚠️  IMPORTANT: This step is REQUIRED before any role-based operations!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });