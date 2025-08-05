const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("💰 Checking BOGO Reward Distributor Balance...\n");

  // Load deployment info
  let networkName = hre.network.name;
  // Handle network name aliases
  if (networkName === 'mainnet') networkName = 'camino';
  if (networkName === 'testnet') networkName = 'columbus';
  
  const deploymentPath = path.join(__dirname, `deployment-${networkName}.json`);
  if (!fs.existsSync(deploymentPath)) {
    throw new Error(`No deployment found for network: ${networkName}`);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  console.log("Network:", deployment.network);
  console.log("Checking balance at:", new Date().toISOString());

  // Get contract instances
  const bogoToken = await hre.ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  const rewardDistributor = await hre.ethers.getContractAt("BOGORewardDistributor", deployment.contracts.BOGORewardDistributor);

  // Check balances
  const distributorBalance = await bogoToken.balanceOf(deployment.contracts.BOGORewardDistributor);
  const deployerBalance = await bogoToken.balanceOf(deployment.deployer);
  const totalSupply = await bogoToken.totalSupply();

  console.log("\n📊 Token Balances:");
  console.log("RewardDistributor:", hre.ethers.formatEther(distributorBalance), "BOGO");
  console.log("Deployer:", hre.ethers.formatEther(deployerBalance), "BOGO");
  console.log("Total Supply:", hre.ethers.formatEther(totalSupply), "BOGO");

  // Check daily distribution status
  const dailyLimit = await rewardDistributor.DAILY_GLOBAL_LIMIT();
  const dailyDistributed = await rewardDistributor.dailyDistributed();
  const remainingDaily = await rewardDistributor.getRemainingDailyLimit();

  console.log("\n📅 Daily Distribution Status:");
  console.log("Daily Limit:", hre.ethers.formatEther(dailyLimit), "BOGO");
  console.log("Already Distributed Today:", hre.ethers.formatEther(dailyDistributed), "BOGO");
  console.log("Remaining Today:", hre.ethers.formatEther(remainingDaily), "BOGO");

  // Check rewards allocation
  const rewardsAllocation = await bogoToken.REWARDS_ALLOCATION();
  const rewardsMinted = await bogoToken.rewardsMinted();
  const remainingRewards = rewardsAllocation - rewardsMinted;

  console.log("\n🎁 Rewards Allocation Status:");
  console.log("Total Rewards Allocation:", hre.ethers.formatEther(rewardsAllocation), "BOGO");
  console.log("Already Minted:", hre.ethers.formatEther(rewardsMinted), "BOGO");
  console.log("Remaining to Mint:", hre.ethers.formatEther(remainingRewards), "BOGO");

  // Analysis
  console.log("\n📈 Analysis:");
  const distributorBalanceNum = Number(hre.ethers.formatEther(distributorBalance));
  const dailyLimitNum = Number(hre.ethers.formatEther(dailyLimit));
  const daysOfRewards = distributorBalanceNum / dailyLimitNum;

  if (distributorBalanceNum === 0) {
    console.log("❌ RewardDistributor is EMPTY! Need to mint and transfer more tokens.");
  } else if (distributorBalanceNum < dailyLimitNum) {
    console.log(`⚠️  RewardDistributor has less than 1 day of rewards (${daysOfRewards.toFixed(2)} days)`);
  } else {
    console.log(`✅ RewardDistributor has ${daysOfRewards.toFixed(1)} days of rewards at current daily limit`);
  }

  if (remainingRewards > 0n) {
    console.log(`💡 Can mint up to ${hre.ethers.formatEther(remainingRewards)} more BOGO from rewards allocation`);
  }

  // Check if distributor has BUSINESS_ROLE
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);
  const BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
  const hasBusinessRole = await roleManager.hasRole(BUSINESS_ROLE, deployment.contracts.BOGORewardDistributor);
  
  console.log("\n🔐 Permissions:");
  console.log("RewardDistributor has BUSINESS_ROLE:", hasBusinessRole ? "✅" : "❌");
  
  if (!hasBusinessRole && remainingRewards > 0n) {
    console.log("⚠️  RewardDistributor cannot mint more rewards without BUSINESS_ROLE!");
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });