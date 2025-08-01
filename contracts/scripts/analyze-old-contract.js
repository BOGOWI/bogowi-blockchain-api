const hre = require("hardhat");
require("dotenv").config({ path: "../.env" });

async function main() {
  const oldDistributor = process.env.REWARD_DISTRIBUTOR_V2_ADDRESS;
  console.log("Analyzing contract:", oldDistributor);
  
  // Get the bytecode to check size
  const code = await hre.ethers.provider.getCode(oldDistributor);
  console.log("Contract code size:", code.length, "chars");
  
  // List of common withdrawal function signatures to check
  const functionSigs = [
    "withdraw()",
    "withdraw(address)",
    "withdraw(address,uint256)",
    "withdrawToken(address)",
    "withdrawToken(address,uint256)",
    "withdrawToken(address,address,uint256)",
    "emergencyWithdraw()",
    "emergencyWithdraw(address)",
    "emergencyWithdraw(address,uint256)",
    "emergencyWithdraw(address,address,uint256)",
    "rescueTokens(address)",
    "rescueTokens(address,uint256)",
    "rescueTokens(address,address,uint256)",
    "transferToken(address,address,uint256)",
    "adminWithdraw(address,uint256)",
    "ownerWithdraw(address,uint256)",
    "drain()",
    "drainToken(address)",
    "destroy()",
    "kill()",
    "destruct()"
  ];
  
  console.log("\nChecking for withdrawal functions...");
  
  const [signer] = await hre.ethers.getSigners();
  let foundAny = false;
  
  for (const sig of functionSigs) {
    try {
      // Try to call each function
      const funcName = sig.split("(")[0];
      const iface = new hre.ethers.utils.Interface([`function ${sig}`]);
      const selector = iface.getSighash(sig);
      
      // Check if function exists by trying to staticcall with empty params
      try {
        await hre.ethers.provider.call({
          to: oldDistributor,
          data: selector,
          from: signer.address
        });
        console.log(`✓ Found function: ${sig}`);
        foundAny = true;
      } catch (e) {
        // Function might exist but revert, check error
        if (e.error && e.error.data && e.error.data !== "0x") {
          console.log(`✓ Found function (reverts): ${sig}`);
          foundAny = true;
        }
      }
      
    } catch (e) {
      // Silently skip - function doesn't exist
    }
  }
  
  if (!foundAny) {
    console.log("\n❌ NO withdrawal functions found!");
    console.log("\nThe contract appears to be designed without any withdrawal mechanism.");
    console.log("This is a security feature - tokens can only exit through legitimate claims.");
    console.log("\n🔒 The 100,000 BOGO tokens are LOCKED in the contract.");
    console.log("\nYour options:");
    console.log("1. Let users claim from the old contract until it's empty");
    console.log("2. Announce migration and have users claim remaining rewards");
    console.log("3. Use both contracts in parallel");
    console.log("4. Consider the locked tokens as 'burned' from circulation");
  }
  
  // Check if contract is upgradeable (proxy)
  console.log("\nChecking if contract is upgradeable...");
  
  // Check for common proxy storage slots
  const implSlot = "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"; // EIP-1967
  const adminSlot = "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103";
  
  const implAddress = await hre.ethers.provider.getStorageAt(oldDistributor, implSlot);
  const adminAddress = await hre.ethers.provider.getStorageAt(oldDistributor, adminSlot);
  
  if (implAddress !== "0x0000000000000000000000000000000000000000000000000000000000000000") {
    console.log("✓ Contract appears to be a proxy!");
    console.log("Implementation:", implAddress);
    console.log("Admin:", adminAddress);
    console.log("\nYou might be able to upgrade the implementation to add withdrawal!");
  } else {
    console.log("❌ Contract is NOT upgradeable");
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });