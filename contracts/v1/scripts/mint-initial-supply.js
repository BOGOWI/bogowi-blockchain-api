const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🪙 Minting initial BOGO token supply from rewards allocation...\n");

  // Load deployment info
  const deploymentPath = path.join(__dirname, `deployment-${hre.network.name}.json`);
  if (!fs.existsSync(deploymentPath)) {
    throw new Error(`No deployment found for network: ${hre.network.name}`);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  const [deployer] = await hre.ethers.getSigners();
  
  console.log("Minting with account:", deployer.address);

  // Get contract instance
  const bogoToken = await hre.ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);

  // Check roles - need DAO_ROLE or BUSINESS_ROLE for mintFromRewards
  const BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
  const DAO_ROLE = await roleManager.DAO_ROLE();
  const hasBusinessRole = await roleManager.hasRole(BUSINESS_ROLE, deployer.address);
  const hasDAORole = await roleManager.hasRole(DAO_ROLE, deployer.address);
  
  if (!hasBusinessRole && !hasDAORole) {
    throw new Error("Deployer needs BUSINESS_ROLE or DAO_ROLE to mint from rewards. Run setup-roles.js first!");
  }

  // Token details
  const name = await bogoToken.name();
  const symbol = await bogoToken.symbol();
  const decimals = await bogoToken.decimals();
  
  console.log("Token:", name, `(${symbol})`);
  console.log("Decimals:", decimals);

  // Check allocations
  const REWARDS_ALLOCATION = await bogoToken.REWARDS_ALLOCATION();
  const rewardsMinted = await bogoToken.rewardsMinted();
  const remainingRewards = await bogoToken.getRemainingRewardsAllocation();
  
  console.log("\n📊 Rewards Allocation Status:");
  console.log("Total Rewards Allocation:", hre.ethers.formatEther(REWARDS_ALLOCATION), symbol);
  console.log("Already Minted:", hre.ethers.formatEther(rewardsMinted), symbol);
  console.log("Remaining:", hre.ethers.formatEther(remainingRewards), symbol);

  // Define amount to mint for initial distribution (10M BOGO)
  const MINT_AMOUNT = hre.ethers.parseEther("10000000"); // 10 million
  
  if (remainingRewards < MINT_AMOUNT) {
    throw new Error(`Insufficient rewards allocation. Need ${hre.ethers.formatEther(MINT_AMOUNT)} but only ${hre.ethers.formatEther(remainingRewards)} remaining`);
  }

  // Mint tokens from rewards allocation
  console.log("\n🏭 Minting from rewards allocation...");
  console.log("Amount to mint:", hre.ethers.formatEther(MINT_AMOUNT), symbol);
  
  const tx = await bogoToken.mintFromRewards(deployer.address, MINT_AMOUNT);
  console.log("Transaction hash:", tx.hash);
  console.log("Waiting for confirmation...");
  
  const receipt = await tx.wait();
  console.log("✅ Minting confirmed in block:", receipt.blockNumber);

  // Verify new balances
  const newRewardsMinted = await bogoToken.rewardsMinted();
  const deployerBalance = await bogoToken.balanceOf(deployer.address);
  const totalSupply = await bogoToken.totalSupply();
  
  console.log("\n📊 Post-mint status:");
  console.log("Total Supply:", hre.ethers.formatEther(totalSupply), symbol);
  console.log("Rewards Minted:", hre.ethers.formatEther(newRewardsMinted), symbol);
  console.log("Deployer Balance:", hre.ethers.formatEther(deployerBalance), symbol);

  // Contract allocations overview
  console.log("\n📋 BOGO Token Allocations (from contract):");
  console.log("- DAO Allocation: 50M BOGO (5%)");
  console.log("- Business Allocation: 900M BOGO (90%)");
  console.log("- Rewards Allocation: 50M BOGO (5%)");
  console.log("  └─ 10M minted now for distributor");
  console.log("  └─ 40M remaining for future rewards");

  console.log("\n✅ Initial rewards supply minted successfully!");
  console.log("\n⚠️  IMPORTANT NEXT STEPS:");
  console.log("1. Transfer 10M tokens to RewardDistributor using fund-distributor.js");
  console.log("2. Grant BUSINESS_ROLE to RewardDistributor so it can mint more rewards as needed");
  console.log("3. The remaining 40M rewards can be minted on-demand by the distributor");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });