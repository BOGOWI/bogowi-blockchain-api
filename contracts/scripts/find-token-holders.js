const hre = require("hardhat");
const { ethers } = hre;

async function main() {
  console.log("\n🔍 FINDING BOGO TOKEN HOLDERS");
  console.log("=".repeat(50));

  const BOGO_TOKEN_ADDRESS = process.env.BOGO_TOKEN_ADDRESS || "0xD394c8fEe6dC8b25DD423AE2f6e68191BD379c0C";
  const bogoToken = await ethers.getContractAt("BOGOTokenV2", BOGO_TOKEN_ADDRESS);
  
  const totalSupply = await bogoToken.totalSupply();
  console.log(`Total Supply: ${ethers.utils.formatEther(totalSupply)} BOGO`);

  // Get all Transfer events to build holder list
  console.log("\n📊 Analyzing Transfer events...");
  
  const filter = bogoToken.filters.Transfer();
  const currentBlock = await ethers.provider.getBlockNumber();
  const BLOCKS_PER_QUERY = 10000;
  
  // Map to track all addresses that have received tokens
  const holders = new Map();
  
  console.log(`Scanning blocks 0 to ${currentBlock}...`);
  
  for (let fromBlock = 0; fromBlock <= currentBlock; fromBlock += BLOCKS_PER_QUERY) {
    const toBlock = Math.min(fromBlock + BLOCKS_PER_QUERY - 1, currentBlock);
    
    try {
      const events = await bogoToken.queryFilter(filter, fromBlock, toBlock);
      
      for (const event of events) {
        // Track recipient
        holders.set(event.args.to, true);
        // Track sender (in case they still have balance)
        holders.set(event.args.from, true);
      }
      
      if (events.length > 0) {
        console.log(`Found ${events.length} transfers in blocks ${fromBlock}-${toBlock}`);
      }
    } catch (error) {
      console.log(`Error querying blocks ${fromBlock}-${toBlock}:`, error.message);
    }
  }
  
  // Remove zero address
  holders.delete(ethers.constants.AddressZero);
  
  console.log(`\nFound ${holders.size} unique addresses that have interacted with BOGO`);
  
  // Check balances of all holders
  console.log("\n💰 CHECKING BALANCES...");
  const balances = [];
  
  for (const address of holders.keys()) {
    try {
      const balance = await bogoToken.balanceOf(address);
      if (balance.gt(0)) {
        balances.push({
          address,
          balance: ethers.utils.formatEther(balance),
          balanceWei: balance
        });
      }
    } catch (error) {
      console.log(`Error checking balance for ${address}`);
    }
  }
  
  // Sort by balance descending
  balances.sort((a, b) => b.balanceWei.gt(a.balanceWei) ? 1 : -1);
  
  console.log("\n📋 TOP TOKEN HOLDERS");
  console.log("-".repeat(80));
  console.log(`${"Rank".padEnd(6)} | ${"Address".padEnd(42)} | ${"Balance".padEnd(25)}`);
  console.log("-".repeat(80));
  
  let totalTracked = ethers.BigNumber.from(0);
  
  for (let i = 0; i < Math.min(balances.length, 20); i++) {
    const holder = balances[i];
    totalTracked = totalTracked.add(holder.balanceWei);
    console.log(`${(i + 1).toString().padEnd(6)} | ${holder.address.padEnd(42)} | ${holder.balance.padEnd(25)} BOGO`);
  }
  
  console.log("\n📊 DISTRIBUTION SUMMARY");
  console.log("-".repeat(50));
  console.log(`Total holders with balance: ${balances.length}`);
  console.log(`Total BOGO tracked: ${ethers.utils.formatEther(totalTracked)}`);
  console.log(`Untracked BOGO: ${ethers.utils.formatEther(totalSupply.sub(totalTracked))}`);
  
  // Categorize addresses
  const deployer = "0xB34A822F735CDE477cbB39a06118267D00948ef7";
  const rewardDistributor = "0xe8431D35b02A913EC32E797C22135F352AD790df";
  
  let deployerBalance = ethers.BigNumber.from(0);
  let distributorBalance = ethers.BigNumber.from(0);
  let userBalances = ethers.BigNumber.from(0);
  
  for (const holder of balances) {
    if (holder.address.toLowerCase() === deployer.toLowerCase()) {
      deployerBalance = holder.balanceWei;
    } else if (holder.address.toLowerCase() === rewardDistributor.toLowerCase()) {
      distributorBalance = holder.balanceWei;
    } else {
      userBalances = userBalances.add(holder.balanceWei);
    }
  }
  
  console.log("\n🎯 ALLOCATION BREAKDOWN");
  console.log("-".repeat(50));
  console.log(`Deployer/Treasury: ${ethers.utils.formatEther(deployerBalance)} BOGO (${deployerBalance.mul(10000).div(totalSupply).toNumber() / 100}%)`);
  console.log(`Reward Distributor: ${ethers.utils.formatEther(distributorBalance)} BOGO (${distributorBalance.mul(10000).div(totalSupply).toNumber() / 100}%)`);
  console.log(`User Wallets: ${ethers.utils.formatEther(userBalances)} BOGO (${userBalances.mul(10000).div(totalSupply).toNumber() / 100}%)`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });