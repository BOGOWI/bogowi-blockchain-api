const hre = require("hardhat");
require("dotenv").config({ path: "../.env" });

async function main() {
  console.log("=== Whitelist and Claim Strategy ===\n");
  
  const [owner] = await hre.ethers.getSigners();
  const oldDistributor = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  
  // Get contract
  const distributor = await hre.ethers.getContractAt("BOGORewardDistributor", oldDistributor);
  
  // Per wallet amounts
  const amounts = {
    welcome_bonus: 10,
    founder_bonus: 100,  // If whitelisted
    referral_bonus: 20,  // If referred
    first_nft_mint: 25,
    dao_participation: 15,
    attraction_tier_1: 10,
    attraction_tier_2: 20,
    attraction_tier_3: 40,
    attraction_tier_4: 50
  };
  
  const totalPerWallet = Object.values(amounts).reduce((a, b) => a + b, 0);
  const totalPerWalletWithoutFounder = totalPerWallet - 100;
  
  console.log("Per wallet (with founder whitelist):", totalPerWallet, "BOGO");
  console.log("Per wallet (without whitelist):", totalPerWalletWithoutFounder, "BOGO");
  
  const remaining = 99810;
  console.log("\nTo drain", remaining, "BOGO:");
  console.log("- With whitelisting:", Math.ceil(remaining / totalPerWallet), "wallets needed");
  console.log("- Without whitelisting:", Math.ceil(remaining / totalPerWalletWithoutFounder), "wallets needed");
  
  console.log("\n=== Fastest Drainage Plan ===");
  console.log("1. Create 20 wallets");
  console.log("2. Whitelist them all for founder bonus");
  console.log("3. Each wallet claims all rewards");
  console.log("4. Total per batch: 20 × 290 = 5,800 BOGO");
  console.log("5. Repeat 18 times (360 wallets total)");
  
  console.log("\n=== Let's Start ===");
  console.log("Creating wallets and whitelisting...\n");
  
  // Create 10 wallets as example
  const wallets = [];
  for (let i = 0; i < 10; i++) {
    const wallet = hre.ethers.Wallet.createRandom();
    wallets.push(wallet.address);
  }
  
  console.log("Generated wallets:");
  wallets.forEach((w, i) => console.log(`${i + 1}. ${w}`));
  
  try {
    // Whitelist all wallets
    console.log("\nWhitelisting all wallets...");
    const tx = await distributor.addToWhitelist(wallets);
    await tx.wait();
    console.log("✅ Whitelisted!");
    
    console.log("\n⚠️  NEXT STEPS:");
    console.log("1. Fund each wallet with small CAM for gas");
    console.log("2. Run claim script for each wallet");
    console.log("3. They'll each claim 290 BOGO");
    console.log("4. Transfer BOGO back to your main wallet");
    
  } catch (error) {
    console.log("❌ Whitelisting failed:", error.message);
    console.log("\nThe contract might have different whitelisting mechanism");
  }
  
  console.log("\n=== Time Estimate ===");
  console.log("Manual process: 360 wallets × 5 min each = 30 hours");
  console.log("With automation: 2-3 hours");
  console.log("Cost: Gas fees for ~2,000 transactions");
  
  console.log("\n=== Alternative ===");
  console.log("Just use the new V3 contract and let this one drain naturally.");
  console.log("The locked tokens effectively reduce supply, increasing value!");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });