const hre = require("hardhat");
require("dotenv").config({ path: "../.env" });

async function main() {
  console.log("=== Withdraw from Old Distributor ===\n");
  
  const [signer] = await hre.ethers.getSigners();
  const bogoTokenAddress = process.env.BOGO_TOKEN_V2_ADDRESS;
  const oldDistributorAddress = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  const newDistributorAddress = process.env.REWARD_DISTRIBUTOR_V3_ADDRESS;
  
  console.log("Signer:", signer.address);
  console.log("Old Distributor (V2):", oldDistributorAddress);
  console.log("New Distributor (V3):", newDistributorAddress);
  console.log("BOGO Token:", bogoTokenAddress);
  
  // Get token contract
  const bogoToken = await hre.ethers.getContractAt("IERC20", bogoTokenAddress);
  
  // Check balance
  const balance = await bogoToken.balanceOf(oldDistributorAddress);
  console.log("\nCurrent balance in old distributor:", hre.ethers.utils.formatEther(balance), "BOGO");
  
  if (balance.eq(0)) {
    console.log("No balance to transfer!");
    return;
  }
  
  // Try to get the old distributor contract
  // First, let's check if it has treasury function (new version)
  try {
    const oldDistributor = await hre.ethers.getContractAt("BOGORewardDistributor", oldDistributorAddress);
    
    // Check if it has treasury function
    try {
      const treasury = await oldDistributor.treasury();
      console.log("\n✓ Old distributor has treasury control:", treasury);
      
      // Check if it has emergencyWithdraw
      try {
        console.log("\nAttempting emergency withdrawal via treasury...");
        console.log("This requires multisig approval from the treasury");
        
        // Create the withdrawal data
        const withdrawData = oldDistributor.interface.encodeFunctionData(
          "emergencyWithdraw",
          [bogoTokenAddress, newDistributorAddress, balance]
        );
        
        console.log("\n📋 MULTISIG INSTRUCTIONS:");
        console.log("1. Submit this transaction to treasury:", treasury);
        console.log("2. Target contract:", oldDistributorAddress);
        console.log("3. Value: 0");
        console.log("4. Data:", withdrawData);
        console.log("5. Description: 'Withdraw BOGO to new distributor'");
        
        console.log("\nOr use this script to submit:");
        console.log(`
const treasury = await ethers.getContractAt("MultisigTreasury", "${treasury}");
await treasury.submitTransaction(
  "${oldDistributorAddress}",
  0,
  "${withdrawData}",
  "Withdraw ${hre.ethers.utils.formatEther(balance)} BOGO to new distributor"
);`);
        
      } catch (e) {
        console.log("\n❌ Old distributor doesn't have emergencyWithdraw function");
        console.log("This contract cannot be drained - funds are locked");
      }
      
    } catch (e) {
      // No treasury function, check for owner
      try {
        const owner = await oldDistributor.owner();
        console.log("\n✓ Old distributor has single owner:", owner);
        
        if (owner.toLowerCase() === signer.address.toLowerCase()) {
          console.log("✓ You are the owner!");
          
          // Check for withdraw functions
          console.log("\nChecking for withdrawal functions...");
          
          // This is an old contract without withdrawal functions
          console.log("\n❌ This appears to be an old contract without withdrawal functions");
          console.log("The funds are locked in this contract");
          console.log("\nOptions:");
          console.log("1. Let it drain naturally through claims");
          console.log("2. Deploy a new contract and migrate users");
          
        } else {
          console.log("❌ You are not the owner");
          console.log("Only the owner can perform admin functions");
        }
        
      } catch (e2) {
        console.log("\n❌ Cannot determine contract ownership model");
        console.log("This might be a custom implementation");
      }
    }
    
  } catch (error) {
    console.log("\n❌ Error accessing old distributor contract:");
    console.log(error.message);
  }
  
  // Check new distributor balance
  const newBalance = await bogoToken.balanceOf(newDistributorAddress);
  console.log("\n=== Current Status ===");
  console.log("Old Distributor balance:", hre.ethers.utils.formatEther(balance), "BOGO");
  console.log("New Distributor balance:", hre.ethers.utils.formatEther(newBalance), "BOGO");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });