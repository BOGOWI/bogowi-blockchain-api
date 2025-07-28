const { ethers } = require("hardhat");

async function main() {
  const rewardDistributorAddress = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  const backendWalletAddress = process.env.BACKEND_WALLET_ADDRESS;
  
  if (!rewardDistributorAddress || !backendWalletAddress) {
    throw new Error("Required environment variables not set");
  }
  
  console.log("Setting up authorized backend...");
  console.log("Reward Distributor:", rewardDistributorAddress);
  console.log("Backend Wallet:", backendWalletAddress);
  
  const [signer] = await ethers.getSigners();
  console.log("Transaction from:", signer.address);
  
  // Get contract instance
  const contract = await ethers.getContractAt("BOGORewardDistributor", rewardDistributorAddress);
  
  // Check if already authorized
  const isAuthorized = await contract.authorizedBackends(backendWalletAddress);
  if (isAuthorized) {
    console.log("✅ Backend already authorized");
    return;
  }
  
  // Authorize backend
  console.log("\nAuthorizing backend...");
  const tx = await contract.setAuthorizedBackend(backendWalletAddress, true);
  console.log("Transaction hash:", tx.hash);
  
  // Wait for confirmation
  await tx.wait();
  console.log("✅ Backend authorized!");
  
  // Verify
  const newStatus = await contract.authorizedBackends(backendWalletAddress);
  console.log("Authorization status:", newStatus);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });