const hre = require("hardhat");
require("dotenv").config();

async function main() {
  console.log("Starting BOGORewardDistributor deployment...");
  
  // Get configuration
  const bogoTokenAddress = process.env.BOGO_TOKEN_V2_ADDRESS;
  
  if (!bogoTokenAddress) {
    throw new Error("BOGO_TOKEN_V2_ADDRESS environment variable not set");
  }
  
  console.log("BOGO Token Address:", bogoTokenAddress);
  console.log("Deploying from account:", (await hre.ethers.getSigners())[0].address);
  
  // Get contract factory
  const RewardDistributor = await hre.ethers.getContractFactory("BOGORewardDistributor");
  
  // Deploy contract
  console.log("Deploying contract...");
  const rewardDistributor = await RewardDistributor.deploy(bogoTokenAddress);
  
  // Wait for deployment
  await rewardDistributor.deployed();
  
  console.log("BOGORewardDistributor deployed to:", rewardDistributor.address);
  console.log("Transaction hash:", rewardDistributor.deployTransaction.hash);
  
  // Wait for confirmations
  console.log("Waiting for 5 confirmations...");
  await rewardDistributor.deployTransaction.wait(5);
  
  console.log("Contract deployment confirmed!");
  console.log("\n=== DEPLOYMENT SUMMARY ===");
  console.log("Contract Address:", rewardDistributor.address);
  console.log("BOGO Token:", bogoTokenAddress);
  console.log("Network:", hre.network.name);
  console.log("========================\n");
  
  // Save deployment info
  const fs = require("fs");
  const deploymentInfo = {
    network: hre.network.name,
    contractAddress: rewardDistributor.address,
    bogoTokenAddress: bogoTokenAddress,
    deploymentTx: rewardDistributor.deployTransaction.hash,
    deployedAt: new Date().toISOString(),
    deployer: (await hre.ethers.getSigners())[0].address
  };
  
  fs.writeFileSync(
    `deployment-${hre.network.name}-${Date.now()}.json`,
    JSON.stringify(deploymentInfo, null, 2)
  );
  
  console.log("Deployment info saved to file");
  console.log("\nNEXT STEPS:");
  console.log("1. Add to .env: REWARD_DISTRIBUTOR_V2_ADDRESS=" + rewardDistributor.address);
  console.log("2. Fund the contract with BOGO tokens");
  console.log("3. Set authorized backend addresses");
  console.log("4. Add founder addresses to whitelist (if needed)");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });