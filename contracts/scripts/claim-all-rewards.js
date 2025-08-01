const hre = require("hardhat");
require("dotenv").config({ path: "../.env" });

async function main() {
  console.log("=== Claim All Available Rewards from Old Distributor ===\n");
  
  const [signer] = await hre.ethers.getSigners();
  const oldDistributor = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  const bogoToken = process.env.BOGO_TOKEN_V2_ADDRESS;
  
  console.log("Your address:", signer.address);
  console.log("Old distributor:", oldDistributor);
  
  // Get contract instance
  const distributor = await hre.ethers.getContractAt("BOGORewardDistributor", oldDistributor);
  const token = await hre.ethers.getContractAt("IERC20", bogoToken);
  
  // Check your current balance
  const balanceBefore = await token.balanceOf(signer.address);
  console.log("Your BOGO balance:", hre.ethers.utils.formatEther(balanceBefore));
  
  // List of all reward templates to try
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
  
  console.log("\nChecking eligibility for all rewards...\n");
  
  let totalClaimed = hre.ethers.BigNumber.from(0);
  
  for (const templateId of templates) {
    try {
      // Check if eligible
      const [eligible, reason] = await distributor.canClaim(signer.address, templateId);
      
      if (eligible) {
        console.log(`✅ ${templateId}: Eligible! Claiming...`);
        
        try {
          // Get template details to know the amount
          const template = await distributor.templates(templateId);
          
          if (template.fixedAmount.gt(0)) {
            // Claim the reward
            const tx = await distributor.claimReward(templateId);
            await tx.wait();
            
            console.log(`   Claimed: ${hre.ethers.utils.formatEther(template.fixedAmount)} BOGO`);
            totalClaimed = totalClaimed.add(template.fixedAmount);
          } else {
            console.log(`   Skipped: Custom amount template`);
          }
          
        } catch (claimError) {
          console.log(`   ❌ Claim failed: ${claimError.reason || claimError.message}`);
        }
        
      } else {
        console.log(`❌ ${templateId}: ${reason}`);
      }
      
    } catch (error) {
      console.log(`❌ ${templateId}: Error checking - ${error.message}`);
    }
  }
  
  // Check if you're whitelisted for founder bonus
  try {
    const isWhitelisted = await distributor.founderWhitelist(signer.address);
    if (!isWhitelisted) {
      console.log("\n💡 TIP: You're not on the founder whitelist.");
      console.log("   As the owner, you could have whitelisted yourself for founder_bonus (100 BOGO)");
    }
  } catch (e) {}
  
  // Check referral options
  console.log("\n=== Referral Bonus ===");
  try {
    const referredBy = await distributor.referredBy(signer.address);
    if (referredBy === hre.ethers.constants.AddressZero) {
      console.log("❌ You haven't been referred by anyone");
      console.log("💡 TIP: You could have someone refer you for 20 BOGO bonus to them");
    }
  } catch (e) {}
  
  // Final balance
  const balanceAfter = await token.balanceOf(signer.address);
  const actualClaimed = balanceAfter.sub(balanceBefore);
  
  console.log("\n=== SUMMARY ===");
  console.log("Total claimed:", hre.ethers.utils.formatEther(actualClaimed), "BOGO");
  console.log("New balance:", hre.ethers.utils.formatEther(balanceAfter), "BOGO");
  
  // Check distributor balance
  const distributorBalance = await token.balanceOf(oldDistributor);
  console.log("Remaining in distributor:", hre.ethers.utils.formatEther(distributorBalance), "BOGO");
  
  if (actualClaimed.eq(0)) {
    console.log("\n💡 TIPS TO CLAIM MORE:");
    console.log("1. Create multiple wallets and claim welcome_bonus (10 BOGO each)");
    console.log("2. Have those wallets claim referral bonus with you as referrer (20 BOGO each to you)");
    console.log("3. Whitelist addresses for founder_bonus (100 BOGO each)");
    console.log("4. Wait for cooldown periods to expire on repeatable rewards");
    console.log("\nNote: This would be the same as what regular users can do.");
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });