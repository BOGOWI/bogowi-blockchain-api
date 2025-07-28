const { ethers } = require("hardhat");

async function main() {
  const rewardDistributorAddress = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  const bogoTokenAddress = process.env.BOGO_TOKEN_V2_ADDRESS;
  const fundingAmount = process.env.FUNDING_AMOUNT || "1000000"; // Default 1M BOGO
  
  if (!rewardDistributorAddress || !bogoTokenAddress) {
    throw new Error("Required environment variables not set");
  }
  
  console.log("Funding Reward Distributor...");
  console.log("Reward Distributor:", rewardDistributorAddress);
  console.log("BOGO Token:", bogoTokenAddress);
  console.log("Amount:", fundingAmount, "BOGO");
  
  const [signer] = await ethers.getSigners();
  console.log("Funding from:", signer.address);
  
  // Get BOGO token contract
  const bogoToken = await ethers.getContractAt("IERC20", bogoTokenAddress);
  
  // Check balance
  const balance = await bogoToken.balanceOf(signer.address);
  console.log("Signer BOGO Balance:", ethers.utils.formatEther(balance), "BOGO");
  
  const amountWei = ethers.utils.parseEther(fundingAmount);
  
  if (balance.lt(amountWei)) {
    throw new Error("Insufficient BOGO balance");
  }
  
  // Transfer tokens
  console.log("\nTransferring tokens...");
  const tx = await bogoToken.transfer(rewardDistributorAddress, amountWei);
  console.log("Transaction hash:", tx.hash);
  
  // Wait for confirmation
  await tx.wait();
  console.log("✅ Transfer confirmed!");
  
  // Check new balance
  const newBalance = await bogoToken.balanceOf(rewardDistributorAddress);
  console.log("\nReward Distributor Balance:", ethers.utils.formatEther(newBalance), "BOGO");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });