const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🔍 Verifying BOGOWI V1 deployment...\n");

  // Load deployment info
  const deploymentPath = path.join(__dirname, `deployment-${hre.network.name}.json`);
  if (!fs.existsSync(deploymentPath)) {
    throw new Error(`No deployment found for network: ${hre.network.name}`);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  console.log("📁 Loaded deployment from:", deployment.timestamp);
  console.log("Network:", deployment.network);
  console.log("Deployer:", deployment.deployer, "\n");

  // Get contract instances
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);
  const bogoToken = await hre.ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  const rewardDistributor = await hre.ethers.getContractAt("BOGORewardDistributor", deployment.contracts.BOGORewardDistributor);

  console.log("🔍 Verifying RoleManager...");
  console.log("Address:", deployment.contracts.RoleManager);
  
  // Check admin role (0x00 is OpenZeppelin's DEFAULT_ADMIN_ROLE)
  const DEFAULT_ADMIN_ROLE = "0x0000000000000000000000000000000000000000000000000000000000000000";
  const hasAdminRole = await roleManager.hasRole(DEFAULT_ADMIN_ROLE, deployment.adminAddress);
  console.log("Admin has DEFAULT_ADMIN_ROLE:", hasAdminRole ? "✅" : "❌");
  
  // Check actual roles that exist
  const BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
  const PAUSER_ROLE = await roleManager.PAUSER_ROLE();
  const DISTRIBUTOR_BACKEND_ROLE = await roleManager.DISTRIBUTOR_BACKEND_ROLE();
  
  console.log("\nRole Constants:");
  console.log("BUSINESS_ROLE:", BUSINESS_ROLE);
  console.log("PAUSER_ROLE:", PAUSER_ROLE);
  console.log("DISTRIBUTOR_BACKEND_ROLE:", DISTRIBUTOR_BACKEND_ROLE);

  console.log("\n🔍 Verifying BOGOToken...");
  console.log("Address:", deployment.contracts.BOGOToken);
  
  // Check token details
  const name = await bogoToken.name();
  const symbol = await bogoToken.symbol();
  const decimals = await bogoToken.decimals();
  const totalSupply = await bogoToken.totalSupply();
  console.log("Name:", name);
  console.log("Symbol:", symbol);
  console.log("Decimals:", decimals);
  console.log("Total Supply:", hre.ethers.formatEther(totalSupply), symbol);
  
  // Check role manager
  const tokenRoleManager = await bogoToken.roleManager();
  console.log("RoleManager matches:", tokenRoleManager === deployment.contracts.RoleManager ? "✅" : "❌");

  console.log("\n🔍 Verifying BOGORewardDistributor...");
  console.log("Address:", deployment.contracts.BOGORewardDistributor);
  
  // Check configuration
  const distributorToken = await rewardDistributor.bogoToken();
  const distributorRoleManager = await rewardDistributor.roleManager();
  console.log("Token address matches:", distributorToken === deployment.contracts.BOGOToken ? "✅" : "❌");
  console.log("RoleManager matches:", distributorRoleManager === deployment.contracts.RoleManager ? "✅" : "❌");
  
  // Check paused state
  const isPaused = await rewardDistributor.paused();
  console.log("Contract paused:", isPaused);
  
  // Check reward templates
  const templates = ["attraction_tier_1", "attraction_tier_2", "attraction_tier_3", "attraction_tier_4", "custom_reward"];
  console.log("\nReward Templates:");
  for (const templateId of templates) {
    try {
      const template = await rewardDistributor.templates(templateId);
      if (template.active) {
        console.log(`- ${templateId}: ✅ Active (${hre.ethers.formatEther(template.fixedAmount)} BOGO)`);
      } else {
        console.log(`- ${templateId}: ❌ Inactive`);
      }
    } catch (e) {
      console.log(`- ${templateId}: ❌ Not found`);
    }
  }

  // Check daily limit
  const DAILY_GLOBAL_LIMIT = await rewardDistributor.DAILY_GLOBAL_LIMIT();
  console.log("\nDaily Global Limit:", hre.ethers.formatEther(DAILY_GLOBAL_LIMIT), "BOGO");

  console.log("\n" + "=".repeat(50));
  console.log("VERIFICATION SUMMARY");
  console.log("=".repeat(50));
  
  // Check if key roles are assigned
  const deployerHasBusinessRole = await roleManager.hasRole(BUSINESS_ROLE, deployment.deployer);
  const distributorHasBusinessRole = await roleManager.hasRole(BUSINESS_ROLE, deployment.contracts.BOGORewardDistributor);
  
  const checks = [
    { name: "RoleManager deployed", pass: deployment.contracts.RoleManager.startsWith("0x") },
    { name: "BOGOToken deployed", pass: deployment.contracts.BOGOToken.startsWith("0x") },
    { name: "RewardDistributor deployed", pass: deployment.contracts.BOGORewardDistributor.startsWith("0x") },
    { name: "Admin has DEFAULT_ADMIN_ROLE", pass: hasAdminRole },
    { name: "Token configured correctly", pass: tokenRoleManager === deployment.contracts.RoleManager },
    { name: "Distributor configured correctly", pass: distributorToken === deployment.contracts.BOGOToken && distributorRoleManager === deployment.contracts.RoleManager },
    { name: "Deployer has BUSINESS_ROLE", pass: deployerHasBusinessRole },
    { name: "Distributor has BUSINESS_ROLE", pass: distributorHasBusinessRole }
  ];

  let allPassed = true;
  for (const check of checks) {
    console.log(`${check.pass ? "✅" : "❌"} ${check.name}`);
    if (!check.pass) allPassed = false;
  }

  console.log("=".repeat(50));
  console.log(allPassed ? "✅ All checks passed!" : "❌ Some checks failed!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });