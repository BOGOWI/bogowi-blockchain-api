const { ethers } = require("hardhat");

async function main() {
  // Check the existing reward distributor address from main .env
  const existingAddress = "0xe8431D35b02A913EC32E797C22135F352AD790df";
  
  console.log("Checking existing deployment at:", existingAddress);
  console.log("Network:", hre.network.name);
  console.log("=====================================\n");
  
  try {
    // Check if contract exists
    const code = await ethers.provider.getCode(existingAddress);
    if (code === "0x") {
      console.log("❌ No contract found at this address!");
      console.log("   You need to deploy a new BOGORewardDistributor contract.");
      return;
    }
    
    console.log("✅ Contract exists at this address");
    console.log("   Checking if it's a BOGORewardDistributor...\n");
    
    // Try to interact with it as BOGORewardDistributor
    try {
      const contract = await ethers.getContractAt("BOGORewardDistributor", existingAddress);
      
      // Check if it has the expected functions
      const bogoToken = await contract.bogoToken();
      console.log("BOGO Token Address:", bogoToken);
      
      // Check if it matches our expected BOGO token
      const expectedBogo = "0xD394c8fEe6dC8b25DD423AE2f6e68191BD379c0C";
      if (bogoToken.toLowerCase() === expectedBogo.toLowerCase()) {
        console.log("✅ BOGO token address matches expected!");
      } else {
        console.log("⚠️  BOGO token address doesn't match expected!");
        console.log("   Expected:", expectedBogo);
      }
      
      // Check templates
      console.log("\nChecking reward templates...");
      const templatesExist = await checkTemplates(contract);
      
      if (templatesExist) {
        console.log("\n✅ This appears to be a valid BOGORewardDistributor!");
        console.log("   You can use this existing deployment.");
        console.log("\n   Next steps:");
        console.log("   1. Run 'npm run verify:columbus' to check full status");
        console.log("   2. Ensure contract is funded");
        console.log("   3. Set up authorized backends");
      } else {
        console.log("\n⚠️  This contract doesn't have the expected templates.");
        console.log("   It might be an older version or different contract.");
        console.log("   Consider deploying a new BOGORewardDistributor.");
      }
      
    } catch (error) {
      console.log("❌ This doesn't appear to be a BOGORewardDistributor contract");
      console.log("   Error:", error.message);
      console.log("\n   You need to deploy a new BOGORewardDistributor contract.");
    }
    
  } catch (error) {
    console.error("Error checking contract:", error.message);
  }
}

async function checkTemplates(contract) {
  const expectedTemplates = [
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
  
  let foundCount = 0;
  
  for (const templateId of expectedTemplates) {
    try {
      const template = await contract.templates(templateId);
      if (template && template.id === templateId) {
        foundCount++;
        console.log(`  ✅ Found template: ${templateId}`);
      }
    } catch (error) {
      console.log(`  ❌ Missing template: ${templateId}`);
    }
  }
  
  return foundCount === expectedTemplates.length;
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });