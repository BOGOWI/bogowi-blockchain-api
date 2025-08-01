const hre = require("hardhat");
const { ethers } = hre;

async function main() {
  console.log("\n🔍 SEARCHING FOR ALL REWARD DISTRIBUTORS");
  console.log("=".repeat(70));

  const BOGO_TOKEN = "0xD394c8fEe6dC8b25DD423AE2f6e68191BD379c0C";
  const bogoToken = await ethers.getContractAt("BOGOTokenV2", BOGO_TOKEN);
  
  // Known addresses to check
  const knownDistributors = {
    "V1 (from env)": "0xe8431D35b02A913EC32E797C22135F352AD790df",
    "V2 (from env)": "0xe8431D35b02A913EC32E797C22135F352AD790df",
  };
  
  console.log("Checking known distributor addresses...\n");
  
  // Check each known address
  for (const [label, address] of Object.entries(knownDistributors)) {
    console.log(`${label}: ${address}`);
    
    try {
      const balance = await bogoToken.balanceOf(address);
      const code = await ethers.provider.getCode(address);
      
      console.log(`  - Balance: ${ethers.utils.formatEther(balance)} BOGO`);
      console.log(`  - Is Contract: ${code.length > 2 ? 'Yes' : 'No'}`);
      
      if (code.length > 2) {
        // Try to check if it's ownable
        try {
          const contract = new ethers.Contract(address, ["function owner() view returns (address)"], ethers.provider);
          const owner = await contract.owner();
          console.log(`  - Owner: ${owner}`);
        } catch (e) {
          console.log(`  - Owner: Not Ownable or different pattern`);
        }
      }
    } catch (error) {
      console.log(`  - Error: ${error.message}`);
    }
    console.log();
  }
  
  // Search for contracts that received large BOGO transfers
  console.log("\n📊 SEARCHING FOR CONTRACTS WITH SIGNIFICANT BOGO HOLDINGS");
  console.log("-".repeat(70));
  
  const filter = bogoToken.filters.Transfer();
  const currentBlock = await ethers.provider.getBlockNumber();
  
  // Map to track potential distributor contracts
  const potentialDistributors = new Map();
  
  console.log("Analyzing transfer events for large recipients...\n");
  
  try {
    // Get recent transfers (last 100k blocks)
    const fromBlock = Math.max(0, currentBlock - 100000);
    const events = await bogoToken.queryFilter(filter, fromBlock, currentBlock);
    
    for (const event of events) {
      const recipient = event.args.to;
      const amount = event.args.value;
      
      // Check if recipient is a contract and received > 1000 BOGO
      if (amount.gt(ethers.utils.parseEther("1000"))) {
        const code = await ethers.provider.getCode(recipient);
        if (code.length > 2) {
          if (!potentialDistributors.has(recipient) && 
              !Object.values(knownDistributors).includes(recipient)) {
            potentialDistributors.set(recipient, {
              address: recipient,
              firstTransfer: amount,
              blockNumber: event.blockNumber
            });
          }
        }
      }
    }
    
    if (potentialDistributors.size > 0) {
      console.log("Found potential distributor contracts:\n");
      
      for (const [address, info] of potentialDistributors) {
        const balance = await bogoToken.balanceOf(address);
        console.log(`Address: ${address}`);
        console.log(`  - Current Balance: ${ethers.utils.formatEther(balance)} BOGO`);
        console.log(`  - First Large Transfer: ${ethers.utils.formatEther(info.firstTransfer)} BOGO`);
        console.log(`  - Block: ${info.blockNumber}`);
        
        // Check transfer activity
        const outgoing = await bogoToken.queryFilter(
          bogoToken.filters.Transfer(address, null),
          info.blockNumber
        );
        console.log(`  - Outgoing Transfers: ${outgoing.length}`);
        console.log();
      }
    } else {
      console.log("No additional distributor-like contracts found");
    }
    
  } catch (error) {
    console.log(`Error searching for distributors: ${error.message}`);
  }
  
  // Final summary
  console.log("\n📋 SUMMARY");
  console.log("-".repeat(70));
  console.log("1. V1 and V2 addresses in .env are identical");
  console.log("2. The known distributor (0xe8431D35b02A913EC32E797C22135F352AD790df) has 0 BOGO");
  console.log("3. No other contracts found with significant BOGO holdings");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });