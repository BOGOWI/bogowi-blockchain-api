const hre = require("hardhat");
const { ethers } = hre;

async function main() {
  console.log("\n🔍 CHECKING FOR WITHDRAWAL FUNCTIONS IN LOCKED DISTRIBUTOR");
  console.log("=".repeat(70));

  const LOCKED_DISTRIBUTOR = "0x95C5Be6Ed592C663fF2953C683dBc5E2C257eA9f";
  const BOGO_TOKEN = "0xD394c8fEe6dC8b25DD423AE2f6e68191BD379c0C";
  const [signer] = await ethers.getSigners();
  
  console.log(`Distributor: ${LOCKED_DISTRIBUTOR}`);
  console.log(`Your address: ${signer.address}`);
  
  try {
    // First, verify ownership
    const ownableABI = ["function owner() view returns (address)"];
    const ownableContract = new ethers.Contract(LOCKED_DISTRIBUTOR, ownableABI, signer);
    const owner = await ownableContract.owner();
    console.log(`Contract owner: ${owner}`);
    console.log(`You are owner: ${owner.toLowerCase() === signer.address.toLowerCase() ? "✅ YES" : "❌ NO"}`);
    
    // Common withdrawal function signatures to try
    const withdrawalFunctions = [
      {
        name: "emergencyWithdraw",
        abi: "function emergencyWithdraw(address token, address to, uint256 amount)",
        args: (token, recipient, amount) => [token, recipient, amount]
      },
      {
        name: "withdrawToken",
        abi: "function withdrawToken(address token, uint256 amount)",
        args: (token, recipient, amount) => [token, amount]
      },
      {
        name: "withdraw",
        abi: "function withdraw(address token)",
        args: (token, recipient, amount) => [token]
      },
      {
        name: "rescueTokens",
        abi: "function rescueTokens(address token, address to, uint256 amount)",
        args: (token, recipient, amount) => [token, recipient, amount]
      },
      {
        name: "recoverERC20",
        abi: "function recoverERC20(address token)",
        args: (token, recipient, amount) => [token]
      },
      {
        name: "sweep",
        abi: "function sweep(address token, address to)",
        args: (token, recipient, amount) => [token, recipient]
      }
    ];
    
    console.log("\n🔧 Testing withdrawal functions...\n");
    
    let foundFunction = false;
    
    for (const func of withdrawalFunctions) {
      try {
        console.log(`Testing ${func.name}...`);
        
        const contract = new ethers.Contract(LOCKED_DISTRIBUTOR, [func.abi], signer);
        
        // Try to estimate gas for the function call
        const args = func.args(BOGO_TOKEN, signer.address, ethers.utils.parseEther("1"));
        const gasEstimate = await contract.estimateGas[func.name](...args);
        
        console.log(`✅ ${func.name} EXISTS! Gas estimate: ${gasEstimate.toString()}`);
        foundFunction = true;
        
        // Ask user if they want to execute
        console.log(`\n💡 Found working function: ${func.name}`);
        console.log(`To withdraw tokens, you can use:`);
        console.log(`\nnpx hardhat run scripts/withdraw-from-old-distributor.js --network camino`);
        break;
        
      } catch (error) {
        if (error.message.includes("not a function") || 
            error.message.includes("no matching function") ||
            error.message.includes("call revert exception")) {
          console.log(`❌ ${func.name} - Not found`);
        } else {
          console.log(`❌ ${func.name} - Error: ${error.reason || error.message.slice(0, 50)}...`);
        }
      }
    }
    
    if (!foundFunction) {
      console.log("\n😞 No standard withdrawal functions found");
      console.log("\n🔍 Checking for other possibilities...");
      
      // Check if contract might be upgradeable
      const upgradeableABI = [
        "function upgradeTo(address newImplementation)",
        "function implementation() view returns (address)"
      ];
      
      try {
        const upgradeableContract = new ethers.Contract(LOCKED_DISTRIBUTOR, upgradeableABI, signer);
        const impl = await upgradeableContract.implementation();
        console.log("\n✅ Contract appears to be upgradeable!");
        console.log(`Current implementation: ${impl}`);
      } catch (e) {
        console.log("❌ Contract is not upgradeable");
      }
      
      // Check if it might be a proxy
      const proxyABI = ["function admin() view returns (address)"];
      try {
        const proxyContract = new ethers.Contract(LOCKED_DISTRIBUTOR, proxyABI, signer);
        const admin = await proxyContract.admin();
        console.log("\n✅ Contract might be a proxy!");
        console.log(`Admin: ${admin}`);
      } catch (e) {
        // Not a proxy
      }
    }
    
  } catch (error) {
    console.error("\nError:", error.message);
  }
  
  console.log("\n" + "=".repeat(70));
  console.log("SUMMARY:");
  console.log("Contract: 0x95C5Be6Ed592C663fF2953C683dBc5E2C257eA9f");
  console.log("Locked BOGO: 99,810 tokens");
  console.log("Recovery options depend on the contract's functions");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });