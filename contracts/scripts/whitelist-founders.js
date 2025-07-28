const { ethers } = require("hardhat");

async function main() {
  const rewardDistributorAddress = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  const founderAddresses = process.env.FOUNDER_ADDRESSES?.split(",") || [];
  
  if (!rewardDistributorAddress) {
    throw new Error("REWARD_DISTRIBUTOR_V2_ADDRESS not set");
  }
  
  if (founderAddresses.length === 0) {
    console.log("No founder addresses provided. Set FOUNDER_ADDRESSES env var with comma-separated addresses.");
    return;
  }
  
  console.log("Adding founders to whitelist...");
  console.log("Reward Distributor:", rewardDistributorAddress);
  console.log("Founder Addresses:", founderAddresses);
  
  const [signer] = await ethers.getSigners();
  console.log("Transaction from:", signer.address);
  
  // Get contract instance
  const contract = await ethers.getContractAt("BOGORewardDistributor", rewardDistributorAddress);
  
  // Check current whitelist status
  console.log("\nChecking current status...");
  for (const founder of founderAddresses) {
    const isWhitelisted = await contract.founderWhitelist(founder);
    console.log(`${founder}: ${isWhitelisted ? "Already whitelisted" : "Not whitelisted"}`);
  }
  
  // Filter only non-whitelisted addresses
  const toWhitelist = [];
  for (const founder of founderAddresses) {
    const isWhitelisted = await contract.founderWhitelist(founder);
    if (!isWhitelisted) {
      toWhitelist.push(founder);
    }
  }
  
  if (toWhitelist.length === 0) {
    console.log("\n✅ All founders already whitelisted");
    return;
  }
  
  // Add to whitelist
  console.log(`\nWhitelisting ${toWhitelist.length} addresses...`);
  const tx = await contract.addToWhitelist(toWhitelist);
  console.log("Transaction hash:", tx.hash);
  
  // Wait for confirmation
  await tx.wait();
  console.log("✅ Whitelist updated!");
  
  // Verify
  console.log("\nVerifying whitelist status...");
  for (const founder of founderAddresses) {
    const isWhitelisted = await contract.founderWhitelist(founder);
    console.log(`${founder}: ${isWhitelisted ? "✅ Whitelisted" : "❌ Not whitelisted"}`);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });