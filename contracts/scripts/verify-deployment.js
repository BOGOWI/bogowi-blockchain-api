const { ethers } = require("hardhat");
require("dotenv").config();

async function main() {
  const contractAddress = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  
  if (!contractAddress) {
    throw new Error("REWARD_DISTRIBUTOR_V2_ADDRESS not set");
  }
  
  console.log("Verifying BOGORewardDistributor at:", contractAddress);
  console.log("Network:", hre.network.name);
  console.log("=====================================\n");
  
  try {
    // Check if contract exists
    const code = await ethers.provider.getCode(contractAddress);
    if (code === "0x") {
      throw new Error("No contract found at this address!");
    }
    console.log("✅ Contract exists on chain");
    
    // Get contract instance
    const contract = await ethers.getContractAt("BOGORewardDistributor", contractAddress);
    
    // Check basic properties
    console.log("\n=== CONTRACT PROPERTIES ===");
    const bogoToken = await contract.bogoToken();
    console.log("BOGO Token:", bogoToken);
    
    const owner = await contract.owner();
    console.log("Owner:", owner);
    
    const paused = await contract.paused();
    console.log("Paused:", paused);
    
    const dailyLimit = await contract.DAILY_GLOBAL_LIMIT();
    console.log("Daily Limit:", ethers.utils.formatEther(dailyLimit), "BOGO");
    
    const dailyDistributed = await contract.dailyDistributed();
    console.log("Daily Distributed:", ethers.utils.formatEther(dailyDistributed), "BOGO");
    
    const remainingLimit = await contract.getRemainingDailyLimit();
    console.log("Remaining Daily Limit:", ethers.utils.formatEther(remainingLimit), "BOGO");
    
    // Check templates
    console.log("\n=== TEMPLATE STATUS ===");
    const templates = [
      "welcome_bonus",
      "founder_bonus",
      "referral_bonus",
      "first_nft_mint",
      "dao_participation",
      "attraction_tier_1",
      "attraction_tier_2",
      "attraction_tier_3",
      "attraction_tier_4",
      "custom_reward"
    ];
    
    let activeCount = 0;
    for (const templateId of templates) {
      const template = await contract.templates(templateId);
      const status = template.active ? "✅ ACTIVE" : "❌ INACTIVE";
      console.log(`${templateId}: ${status}`);
      if (template.active) {
        activeCount++;
        if (template.fixedAmount && template.fixedAmount.gt(0)) {
          console.log(`  - Amount: ${ethers.utils.formatEther(template.fixedAmount)} BOGO`);
        }
        if (template.maxAmount && template.maxAmount.gt(0)) {
          console.log(`  - Max Amount: ${ethers.utils.formatEther(template.maxAmount)} BOGO`);
        }
        if (template.maxClaimsPerWallet && template.maxClaimsPerWallet.gt(0)) {
          console.log(`  - Max Claims: ${template.maxClaimsPerWallet}`);
        }
        if (template.cooldownPeriod && template.cooldownPeriod.gt(0)) {
          console.log(`  - Cooldown: ${template.cooldownPeriod} seconds`);
        }
        if (template.requiresWhitelist) {
          console.log(`  - Requires Whitelist: YES`);
        }
      }
    }
    
    console.log(`\nActive Templates: ${activeCount}/${templates.length}`);
    
    // Check BOGO token balance
    console.log("\n=== TOKEN BALANCE ===");
    const bogoTokenContract = await ethers.getContractAt(
      "IERC20",
      bogoToken
    );
    const balance = await bogoTokenContract.balanceOf(contractAddress);
    console.log("Contract BOGO Balance:", ethers.utils.formatEther(balance), "BOGO");
    
    const allowance = await bogoTokenContract.allowance(owner, contractAddress);
    console.log("Owner Allowance:", ethers.utils.formatEther(allowance), "BOGO");
    
    // Check authorized backends
    console.log("\n=== AUTHORIZED BACKENDS ===");
    const backendAddress = process.env.BACKEND_WALLET_ADDRESS;
    if (backendAddress) {
      const isAuthorized = await contract.authorizedBackends(backendAddress);
      console.log(`${backendAddress}: ${isAuthorized ? "✅ AUTHORIZED" : "❌ NOT AUTHORIZED"}`);
    } else {
      console.log("No BACKEND_WALLET_ADDRESS set to check");
    }
    
    console.log("\n=== VERIFICATION COMPLETE ===");
    console.log("✅ Contract is deployed and initialized");
    
    if (balance.eq(0) && allowance.eq(0)) {
      console.log("\n⚠️  WARNING: Contract has no BOGO tokens and no allowance!");
      console.log("   You need to either:");
      console.log("   1. Transfer BOGO tokens to the contract, OR");
      console.log("   2. Approve the contract to spend BOGO tokens");
    }
    
  } catch (error) {
    console.error("❌ Verification failed:", error.message);
    process.exit(1);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });