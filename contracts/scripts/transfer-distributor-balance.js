const hre = require("hardhat");
require("dotenv").config({ path: "../.env" });

async function main() {
  console.log("=== Transfer BOGO Balance from Old to New Distributor ===\n");
  
  // Get configuration
  const bogoTokenAddress = process.env.BOGO_TOKEN_V2_ADDRESS;
  const oldDistributorAddress = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  const newDistributorAddress = process.env.REWARD_DISTRIBUTOR_V3_ADDRESS;
  
  if (!bogoTokenAddress || !oldDistributorAddress) {
    throw new Error("Missing required addresses in .env");
  }
  
  if (!newDistributorAddress) {
    console.log("REWARD_DISTRIBUTOR_V3_ADDRESS not set in .env");
    console.log("Please deploy the new distributor first and add its address to .env");
    process.exit(1);
  }
  
  const [signer] = await hre.ethers.getSigners();
  console.log("Operating with account:", signer.address);
  
  // Get contract instances
  const bogoToken = await hre.ethers.getContractAt("IERC20", bogoTokenAddress);
  const oldDistributor = await hre.ethers.getContractAt("BOGORewardDistributor", oldDistributorAddress);
  
  // Check current balance
  const balance = await bogoToken.balanceOf(oldDistributorAddress);
  const formattedBalance = hre.ethers.utils.formatEther(balance);
  
  console.log("\n=== Current Status ===");
  console.log("Old Distributor:", oldDistributorAddress);
  console.log("New Distributor:", newDistributorAddress);
  console.log("BOGO Token:", bogoTokenAddress);
  console.log("Current Balance:", formattedBalance, "BOGO");
  
  if (balance.eq(0)) {
    console.log("\n❌ Old distributor has no BOGO balance to transfer");
    return;
  }
  
  // Check if we're the owner of old distributor
  try {
    const owner = await oldDistributor.owner();
    console.log("\nOld Distributor Owner:", owner);
    
    if (owner.toLowerCase() !== signer.address.toLowerCase()) {
      console.log("\n❌ ERROR: You are not the owner of the old distributor");
      console.log("Only the owner can withdraw tokens");
      console.log("Your address:", signer.address);
      console.log("Owner address:", owner);
      return;
    }
  } catch (error) {
    console.log("\n❌ ERROR: Cannot read owner - contract might use different access control");
    console.log("Error:", error.message);
    return;
  }
  
  // Confirmation prompt
  console.log("\n⚠️  CONFIRMATION REQUIRED ⚠️");
  console.log(`About to transfer ${formattedBalance} BOGO tokens`);
  console.log(`FROM: ${oldDistributorAddress}`);
  console.log(`TO:   ${newDistributorAddress}`);
  console.log("\nPress Ctrl+C to cancel, or wait 5 seconds to continue...");
  
  await new Promise(resolve => setTimeout(resolve, 5000));
  
  console.log("\n📤 Initiating transfer...");
  
  try {
    // Most distributor contracts have an emergency withdraw function
    // Try common function names
    let tx;
    
    // Try withdrawToken function
    try {
      tx = await oldDistributor.withdrawToken(bogoTokenAddress, newDistributorAddress, balance);
      console.log("✓ Using withdrawToken function");
    } catch (e1) {
      // Try emergencyWithdraw
      try {
        tx = await oldDistributor.emergencyWithdraw(bogoTokenAddress, newDistributorAddress, balance);
        console.log("✓ Using emergencyWithdraw function");
      } catch (e2) {
        // Try withdraw
        try {
          tx = await oldDistributor.withdraw(bogoTokenAddress, balance);
          console.log("✓ Using withdraw function (will need manual transfer)");
        } catch (e3) {
          console.log("\n❌ ERROR: No withdraw function found on old distributor");
          console.log("The contract might not have a withdrawal mechanism");
          console.log("\nYou may need to:");
          console.log("1. Deploy a new distributor without funding it");
          console.log("2. Update all systems to use the new distributor");
          console.log("3. Let the old one drain naturally through claims");
          return;
        }
      }
    }
    
    console.log("Transaction hash:", tx.hash);
    console.log("Waiting for confirmation...");
    
    const receipt = await tx.wait();
    console.log("✅ Transaction confirmed!");
    console.log("Gas used:", receipt.gasUsed.toString());
    
    // Verify the transfer
    const newBalance = await bogoToken.balanceOf(oldDistributorAddress);
    const newDistributorBalance = await bogoToken.balanceOf(newDistributorAddress);
    
    console.log("\n=== Final Balances ===");
    console.log("Old Distributor:", hre.ethers.utils.formatEther(newBalance), "BOGO");
    console.log("New Distributor:", hre.ethers.utils.formatEther(newDistributorBalance), "BOGO");
    
    if (newBalance.eq(0)) {
      console.log("\n✅ SUCCESS: All tokens transferred!");
    } else {
      console.log("\n⚠️  WARNING: Some tokens remain in old distributor");
    }
    
  } catch (error) {
    console.log("\n❌ ERROR during transfer:");
    console.log(error.message);
    
    // If direct transfer failed, provide manual instructions
    console.log("\n📋 MANUAL TRANSFER INSTRUCTIONS:");
    console.log("1. The old distributor might not have a direct transfer function");
    console.log("2. You may need to:");
    console.log("   a. Check if there's an admin function to withdraw");
    console.log("   b. Or deploy new distributor and update all references");
    console.log("   c. Or let the old one drain through normal operations");
  }
  
  console.log("\n=== Next Steps ===");
  console.log("1. Update .env: REWARD_DISTRIBUTOR_V2_ADDRESS=" + newDistributorAddress);
  console.log("2. Update backend systems to use new address");
  console.log("3. Test the new distributor with a small claim");
  console.log("4. Monitor both distributors during transition");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });