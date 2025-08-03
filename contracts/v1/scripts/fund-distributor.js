const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("💰 Funding BOGO Reward Distributor...\n");

  // Load deployment info
  const deploymentPath = path.join(__dirname, `deployment-${hre.network.name}.json`);
  if (!fs.existsSync(deploymentPath)) {
    throw new Error(`No deployment found for network: ${hre.network.name}`);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  const [deployer] = await hre.ethers.getSigners();
  
  console.log("Funding with account:", deployer.address);

  // Get contract instances
  const bogoToken = await hre.ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  const rewardDistributor = await hre.ethers.getContractAt("BOGORewardDistributor", deployment.contracts.BOGORewardDistributor);

  // Get current balances
  const deployerBalance = await bogoToken.balanceOf(deployer.address);
  const distributorBalance = await bogoToken.balanceOf(deployment.contracts.BOGORewardDistributor);
  
  console.log("Current Balances:");
  console.log("Deployer:", hre.ethers.formatEther(deployerBalance), "BOGO");
  console.log("RewardDistributor:", hre.ethers.formatEther(distributorBalance), "BOGO");

  // Define funding amount (10M BOGO from rewards allocation)
  const FUNDING_AMOUNT = hre.ethers.parseEther("10000000"); // 10 million
  
  if (deployerBalance < FUNDING_AMOUNT) {
    throw new Error(`Insufficient balance. Need ${hre.ethers.formatEther(FUNDING_AMOUNT)} BOGO but have ${hre.ethers.formatEther(deployerBalance)}`);
  }

  // Transfer tokens to distributor
  console.log("\n💸 Transferring tokens to RewardDistributor...");
  console.log("Amount:", hre.ethers.formatEther(FUNDING_AMOUNT), "BOGO");
  
  const tx = await bogoToken.transfer(deployment.contracts.BOGORewardDistributor, FUNDING_AMOUNT);
  console.log("Transaction hash:", tx.hash);
  console.log("Waiting for confirmation...");
  
  const receipt = await tx.wait();
  console.log("✅ Transfer confirmed in block:", receipt.blockNumber);

  // Verify new balances
  const newDeployerBalance = await bogoToken.balanceOf(deployer.address);
  const newDistributorBalance = await bogoToken.balanceOf(deployment.contracts.BOGORewardDistributor);
  
  console.log("\n📊 Post-transfer balances:");
  console.log("Deployer:", hre.ethers.formatEther(newDeployerBalance), "BOGO");
  console.log("RewardDistributor:", hre.ethers.formatEther(newDistributorBalance), "BOGO");

  // Check daily limit vs balance
  const DAILY_GLOBAL_LIMIT = await rewardDistributor.DAILY_GLOBAL_LIMIT();
  const daysOfRewards = newDistributorBalance / DAILY_GLOBAL_LIMIT;
  
  console.log("\n📈 Distribution capacity:");
  console.log("Daily limit:", hre.ethers.formatEther(DAILY_GLOBAL_LIMIT), "BOGO");
  console.log("Days of rewards available:", Math.floor(Number(daysOfRewards)));

  // Check reward templates
  console.log("\n🎁 Reward templates funding status:");
  const templates = ["attraction_tier_1", "attraction_tier_2", "attraction_tier_3", "attraction_tier_4"];
  
  for (const templateId of templates) {
    const template = await rewardDistributor.templates(templateId);
    if (template.active) {
      const rewardsAvailable = newDistributorBalance / template.fixedAmount;
      console.log(`- ${templateId}: ${Number(rewardsAvailable).toLocaleString()} rewards @ ${hre.ethers.formatEther(template.fixedAmount)} BOGO each`);
    }
  }

  console.log("\n✅ RewardDistributor funded successfully!");
  console.log("\n💡 The distributor is now ready to:");
  console.log("- Process reward claims from users");
  console.log("- Handle referral rewards");
  console.log("- Distribute custom rewards via backend");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });