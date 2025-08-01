const hre = require("hardhat");
require("dotenv").config({ path: "../.env" });

async function main() {
  const oldDistributorAddress = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  const treasuryAddress = process.env.MULTISIG_TREASURY_ADDRESS;
  
  console.log("Checking old distributor:", oldDistributorAddress);
  console.log("Current treasury:", treasuryAddress);
  
  // Try to call treasury() function directly
  try {
    const oldDistributor = await hre.ethers.getContractAt(
      ["function treasury() view returns (address)"],
      oldDistributorAddress
    );
    
    const treasury = await oldDistributor.treasury();
    console.log("\n✅ Old distributor HAS treasury control!");
    console.log("Treasury address:", treasury);
    
    if (treasury.toLowerCase() === treasuryAddress.toLowerCase()) {
      console.log("✅ It uses the SAME treasury as your new setup!");
      
      // Now check if it has emergencyWithdraw
      try {
        const distributorWithWithdraw = await hre.ethers.getContractAt(
          ["function emergencyWithdraw(address token, address to, uint256 amount) external"],
          oldDistributorAddress  
        );
        
        console.log("\n✅ Old distributor HAS emergencyWithdraw function!");
        console.log("\nYou can withdraw using multisig!");
        
        // Generate the transaction data
        const bogoToken = process.env.BOGO_TOKEN_V2_ADDRESS;
        const newDistributor = process.env.REWARD_DISTRIBUTOR_V3_ADDRESS;
        const amount = hre.ethers.utils.parseEther("100000");
        
        const iface = new hre.ethers.utils.Interface([
          "function emergencyWithdraw(address token, address to, uint256 amount)"
        ]);
        
        const data = iface.encodeFunctionData("emergencyWithdraw", [
          bogoToken,
          newDistributor, 
          amount
        ]);
        
        console.log("\n📋 SUBMIT THIS TO MULTISIG:");
        console.log("================================");
        console.log("To:", oldDistributorAddress);
        console.log("Value:", "0");
        console.log("Data:", data);
        console.log("Description:", "Withdraw 100k BOGO to new distributor");
        console.log("================================");
        
      } catch (e) {
        console.log("\n❌ Old distributor does NOT have emergencyWithdraw");
        console.log("Even though it has treasury control, it cannot withdraw funds");
      }
      
    } else {
      console.log("⚠️  It uses a DIFFERENT treasury:", treasury);
      console.log("You need signers from that treasury to withdraw");
    }
    
  } catch (error) {
    console.log("\n❌ Old distributor does NOT have treasury control");
    console.log("It's likely using the old single-owner pattern");
    
    // Check for owner
    try {
      const oldDistributor = await hre.ethers.getContractAt(
        ["function owner() view returns (address)"],
        oldDistributorAddress
      );
      
      const owner = await oldDistributor.owner();
      console.log("Owner:", owner);
      
      const [signer] = await hre.ethers.getSigners();
      if (owner.toLowerCase() === signer.address.toLowerCase()) {
        console.log("✅ You are the owner!");
      } else {
        console.log("❌ You are NOT the owner");
      }
      
    } catch (e) {
      console.log("Cannot read owner either - unknown access control");
    }
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });