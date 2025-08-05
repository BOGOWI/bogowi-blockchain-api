const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("💰 Refunding BOGO Reward Distributor...\n");

  const [deployer] = await hre.ethers.getSigners();
  console.log("Operating with account:", deployer.address);

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
  
  // Get contract instances
  const bogoToken = await hre.ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  const rewardDistributor = await hre.ethers.getContractAt("BOGORewardDistributor", deployment.contracts.BOGORewardDistributor);
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);

  // Check current balances
  const distributorBalance = await bogoToken.balanceOf(deployment.contracts.BOGORewardDistributor);
  const deployerBalance = await bogoToken.balanceOf(deployer.address);

  console.log("Current Balances:");
  console.log("RewardDistributor:", hre.ethers.formatEther(distributorBalance), "BOGO");
  console.log("Deployer:", hre.ethers.formatEther(deployerBalance), "BOGO");

  // Check rewards allocation
  const rewardsAllocation = await bogoToken.REWARDS_ALLOCATION();
  const rewardsMinted = await bogoToken.rewardsMinted();
  const remainingRewards = rewardsAllocation - rewardsMinted;

  console.log("\n📊 Rewards Allocation Status:");
  console.log("Total Rewards Allocation:", hre.ethers.formatEther(rewardsAllocation), "BOGO");
  console.log("Already Minted:", hre.ethers.formatEther(rewardsMinted), "BOGO");
  console.log("Remaining:", hre.ethers.formatEther(remainingRewards), "BOGO");

  // Determine how much to mint
  const targetBalance = hre.ethers.parseEther("10000000"); // 10M BOGO target
  const currentBalance = distributorBalance;
  const needed = targetBalance > currentBalance ? targetBalance - currentBalance : 0n;

  if (needed === 0n) {
    console.log("\n✅ RewardDistributor already has sufficient balance!");
    return;
  }

  console.log("\n💰 Need to add:", hre.ethers.formatEther(needed), "BOGO to reach target of 10M");

  // Check if we need to mint or just transfer
  if (deployerBalance >= needed) {
    console.log("✅ Deployer has enough balance, will transfer directly");
    
    console.log("\n💸 Transferring tokens to RewardDistributor...");
    const tx = await bogoToken.transfer(deployment.contracts.BOGORewardDistributor, needed);
    console.log("Transaction hash:", tx.hash);
    await tx.wait();
    console.log("✅ Transfer confirmed!");
  } else {
    // Need to mint more tokens
    const toMint = needed - deployerBalance;
    
    if (toMint > remainingRewards) {
      console.log(`❌ Cannot mint ${hre.ethers.formatEther(toMint)} BOGO - only ${hre.ethers.formatEther(remainingRewards)} remaining in rewards allocation!`);
      return;
    }

    // Check if deployer has BUSINESS_ROLE
    const BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
    const hasBusinessRole = await roleManager.hasRole(BUSINESS_ROLE, deployer.address);
    
    if (!hasBusinessRole) {
      console.log("❌ Deployer doesn't have BUSINESS_ROLE to mint rewards!");
      console.log("Please grant BUSINESS_ROLE to deployer first");
      return;
    }

    console.log(`\n🏭 Minting ${hre.ethers.formatEther(toMint)} BOGO from rewards allocation...`);
    const mintTx = await bogoToken.mintRewards(toMint);
    console.log("Mint transaction hash:", mintTx.hash);
    await mintTx.wait();
    console.log("✅ Minting confirmed!");

    // Now transfer all to distributor
    console.log("\n💸 Transferring all tokens to RewardDistributor...");
    const transferTx = await bogoToken.transfer(deployment.contracts.BOGORewardDistributor, needed);
    console.log("Transfer transaction hash:", transferTx.hash);
    await transferTx.wait();
    console.log("✅ Transfer confirmed!");
  }

  // Verify final balances
  const finalDistributorBalance = await bogoToken.balanceOf(deployment.contracts.BOGORewardDistributor);
  const finalDeployerBalance = await bogoToken.balanceOf(deployer.address);

  console.log("\n📊 Final Balances:");
  console.log("RewardDistributor:", hre.ethers.formatEther(finalDistributorBalance), "BOGO");
  console.log("Deployer:", hre.ethers.formatEther(finalDeployerBalance), "BOGO");

  // Calculate days of rewards
  const dailyLimit = await rewardDistributor.DAILY_GLOBAL_LIMIT();
  const daysOfRewards = Number(hre.ethers.formatEther(finalDistributorBalance)) / Number(hre.ethers.formatEther(dailyLimit));
  console.log(`\n✅ RewardDistributor now has ${daysOfRewards.toFixed(1)} days of rewards!`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });