const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🔐 Setting up roles for BOGOWI V1...\n");

  // Load deployment info
  const deploymentPath = path.join(__dirname, `deployment-${hre.network.name}.json`);
  if (!fs.existsSync(deploymentPath)) {
    throw new Error(`No deployment found for network: ${hre.network.name}`);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  const [deployer] = await hre.ethers.getSigners();
  
  console.log("Setting up roles with account:", deployer.address);

  // Get contract instances
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);
  const bogoToken = await hre.ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  const rewardDistributor = await hre.ethers.getContractAt("BOGORewardDistributor", deployment.contracts.BOGORewardDistributor);

  // Get role constants
  const BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
  const DAO_ROLE = await roleManager.DAO_ROLE();
  const PAUSER_ROLE = await roleManager.PAUSER_ROLE();
  const OPERATOR_ROLE = await roleManager.OPERATOR_ROLE();
  const DISTRIBUTOR_ROLE = await roleManager.DISTRIBUTOR_ROLE();

  console.log("📋 Role IDs:");
  console.log("BUSINESS_ROLE:", BUSINESS_ROLE);
  console.log("DAO_ROLE:", DAO_ROLE);
  console.log("PAUSER_ROLE:", PAUSER_ROLE);
  console.log("OPERATOR_ROLE:", OPERATOR_ROLE);
  console.log("DISTRIBUTOR_ROLE:", DISTRIBUTOR_ROLE);

  // Setup roles for BOGOToken
  console.log("\n🪙 Setting up BOGOToken roles...");
  
  // Grant BUSINESS_ROLE to deployer (temporary for initial minting)
  console.log("Granting BUSINESS_ROLE to deployer...");
  let tx = await roleManager.grantRole(BUSINESS_ROLE, deployer.address);
  await tx.wait();
  console.log("✅ BUSINESS_ROLE granted to deployer");

  // Grant BUSINESS_ROLE to RewardDistributor (for minting rewards on-demand)
  console.log("Granting BUSINESS_ROLE to RewardDistributor...");
  tx = await roleManager.grantRole(BUSINESS_ROLE, deployment.contracts.BOGORewardDistributor);
  await tx.wait();
  console.log("✅ BUSINESS_ROLE granted to RewardDistributor");

  // Setup roles for RewardDistributor
  console.log("\n🎁 Setting up RewardDistributor roles...");
  
  // Grant DISTRIBUTOR_ROLE to deployer (for testing)
  console.log("Granting DISTRIBUTOR_ROLE to deployer...");
  tx = await roleManager.grantRole(DISTRIBUTOR_ROLE, deployer.address);
  await tx.wait();
  console.log("✅ DISTRIBUTOR_ROLE granted to deployer");

  // Grant OPERATOR_ROLE to deployer
  console.log("Granting OPERATOR_ROLE to deployer...");
  tx = await roleManager.grantRole(OPERATOR_ROLE, deployer.address);
  await tx.wait();
  console.log("✅ OPERATOR_ROLE granted to deployer");

  // Grant PAUSER_ROLE to admin
  console.log("\n🛑 Setting up emergency roles...");
  console.log("Granting PAUSER_ROLE to admin...");
  tx = await roleManager.grantRole(PAUSER_ROLE, deployment.adminAddress);
  await tx.wait();
  console.log("✅ PAUSER_ROLE granted to admin");

  // Backend wallet setup (if provided)
  const backendAddress = process.env.BACKEND_WALLET_ADDRESS;
  if (backendAddress && backendAddress !== "YOUR_BACKEND_WALLET_ADDRESS") {
    console.log("\n🔧 Setting up backend wallet roles...");
    console.log("Backend wallet:", backendAddress);
    
    // Grant DISTRIBUTOR_ROLE to backend
    console.log("Granting DISTRIBUTOR_ROLE to backend...");
    tx = await roleManager.grantRole(DISTRIBUTOR_ROLE, backendAddress);
    await tx.wait();
    console.log("✅ DISTRIBUTOR_ROLE granted to backend");

    // Grant OPERATOR_ROLE to backend
    console.log("Granting OPERATOR_ROLE to backend...");
    tx = await roleManager.grantRole(OPERATOR_ROLE, backendAddress);
    await tx.wait();
    console.log("✅ OPERATOR_ROLE granted to backend");
  }

  // Verify roles
  console.log("\n🔍 Verifying role assignments...");
  
  const roles = [
    { role: BUSINESS_ROLE, name: "BUSINESS", address: deployer.address },
    { role: BUSINESS_ROLE, name: "BUSINESS", address: deployment.contracts.BOGORewardDistributor },
    { role: DISTRIBUTOR_ROLE, name: "DISTRIBUTOR", address: deployer.address },
    { role: OPERATOR_ROLE, name: "OPERATOR", address: deployer.address },
    { role: PAUSER_ROLE, name: "PAUSER", address: deployment.adminAddress }
  ];

  if (backendAddress && backendAddress !== "YOUR_BACKEND_WALLET_ADDRESS") {
    roles.push({ role: DISTRIBUTOR_ROLE, name: "DISTRIBUTOR", address: backendAddress });
    roles.push({ role: OPERATOR_ROLE, name: "OPERATOR", address: backendAddress });
  }

  for (const { role, name, address } of roles) {
    const hasRole = await roleManager.hasRole(role, address);
    console.log(`${hasRole ? "✅" : "❌"} ${address} has ${name}_ROLE`);
  }

  console.log("\n✅ Role setup complete!");
  console.log("\n⚠️  IMPORTANT SECURITY NOTE:");
  console.log("- Remove BUSINESS_ROLE from deployer after initial rewards are minted");
  console.log("- RewardDistributor can mint more rewards on-demand from the 50M allocation");
  console.log("- Consider using multisig for admin roles in production");
  console.log("- Regularly audit role assignments");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });