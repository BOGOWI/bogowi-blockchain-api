const hre = require("hardhat");
const { ethers } = hre;

async function main() {
  console.log("\n🔍 BOGO TOKEN ALLOCATION ANALYSIS");
  console.log("=".repeat(50));

  // Get network info
  const network = await ethers.provider.getNetwork();
  console.log(`Network: ${network.name} (chainId: ${network.chainId})`);
  console.log(`Block: ${await ethers.provider.getBlockNumber()}`);
  console.log("=".repeat(50));

  // Contract addresses - update these with your deployed addresses
  const addresses = {
    bogoToken: process.env.BOGO_TOKEN_ADDRESS || "",
    rewardDistributor: process.env.REWARD_DISTRIBUTOR_ADDRESS || "",
    multisigTreasury: process.env.MULTISIG_TREASURY_ADDRESS || "",
    conservationNFT: process.env.CONSERVATION_NFT_ADDRESS || "",
    commercialNFT: process.env.COMMERCIAL_NFT_ADDRESS || "",
    // Add any other contract addresses here
  };

  // Check if BOGO token address is provided
  if (!addresses.bogoToken) {
    console.error("❌ BOGO_TOKEN_ADDRESS not set in environment variables!");
    console.log("\nPlease set the following in your .env file:");
    console.log("BOGO_TOKEN_ADDRESS=0x...");
    console.log("REWARD_DISTRIBUTOR_ADDRESS=0x...");
    console.log("MULTISIG_TREASURY_ADDRESS=0x...");
    return;
  }

  // Get BOGO token contract
  const bogoToken = await ethers.getContractAt("BOGOTokenV2", addresses.bogoToken);
  
  // Get token info
  const name = await bogoToken.name();
  const symbol = await bogoToken.symbol();
  const decimals = await bogoToken.decimals();
  const totalSupply = await bogoToken.totalSupply();
  
  console.log(`\n📊 TOKEN INFO`);
  console.log(`Name: ${name}`);
  console.log(`Symbol: ${symbol}`);
  console.log(`Decimals: ${decimals}`);
  console.log(`Total Supply: ${ethers.utils.formatEther(totalSupply)} ${symbol}`);
  console.log("=".repeat(50));

  // Function to check balance
  async function checkBalance(address, label) {
    if (!address || address === "") {
      return { label, address: "Not deployed", balance: "0", percentage: "0" };
    }
    try {
      const balance = await bogoToken.balanceOf(address);
      const percentage = totalSupply.gt(0) 
        ? balance.mul(10000).div(totalSupply).toNumber() / 100 
        : 0;
      return {
        label,
        address,
        balance: ethers.utils.formatEther(balance),
        percentage: percentage.toFixed(2)
      };
    } catch (error) {
      return { label, address, balance: "Error", percentage: "0" };
    }
  }

  console.log("\n💰 TOKEN ALLOCATION BY CONTRACT");
  console.log("-".repeat(100));
  console.log(`${"Contract".padEnd(25)} | ${"Address".padEnd(42)} | ${"Balance".padEnd(20)} | %`);
  console.log("-".repeat(100));

  // Check all contract balances
  const contractBalances = [];
  for (const [key, address] of Object.entries(addresses)) {
    if (key !== "bogoToken" && address) {
      const result = await checkBalance(address, key);
      contractBalances.push(result);
      console.log(`${result.label.padEnd(25)} | ${result.address.padEnd(42)} | ${result.balance.padEnd(20)} | ${result.percentage}%`);
    }
  }

  // Check special addresses
  console.log("\n💎 SPECIAL ALLOCATIONS");
  console.log("-".repeat(100));

  // Get role-based addresses
  const DEFAULT_ADMIN_ROLE = await bogoToken.DEFAULT_ADMIN_ROLE();
  const DAO_ROLE = await bogoToken.DAO_ROLE();
  const admins = [];
  const daoMembers = [];

  // Check for role holders (this is a simplified check - you might need to adjust based on your contract)
  try {
    // Get signers to check their roles
    const signers = await ethers.getSigners();
    for (let i = 0; i < Math.min(signers.length, 10); i++) {
      const address = signers[i].address;
      const hasAdminRole = await bogoToken.hasRole(DEFAULT_ADMIN_ROLE, address);
      const hasDaoRole = await bogoToken.hasRole(DAO_ROLE, address);
      
      if (hasAdminRole) admins.push(address);
      if (hasDaoRole) daoMembers.push(address);
    }
  } catch (error) {
    console.log("Could not check roles:", error.message);
  }

  // Check deployer balance
  const [deployer] = await ethers.getSigners();
  const deployerBalance = await checkBalance(deployer.address, "Deployer");
  console.log(`${deployerBalance.label.padEnd(25)} | ${deployerBalance.address.padEnd(42)} | ${deployerBalance.balance.padEnd(20)} | ${deployerBalance.percentage}%`);

  // Check admin balances
  for (const admin of admins) {
    const adminBalance = await checkBalance(admin, "Admin");
    console.log(`${adminBalance.label.padEnd(25)} | ${adminBalance.address.padEnd(42)} | ${adminBalance.balance.padEnd(20)} | ${adminBalance.percentage}%`);
  }

  // Check reward allocations
  console.log("\n🎁 REWARD POOL ALLOCATIONS");
  console.log("-".repeat(100));
  
  try {
    const rewardSupply = await bogoToken.rewardSupply();
    const mintedRewards = await bogoToken.mintedRewards();
    const remainingRewards = rewardSupply.sub(mintedRewards);
    
    console.log(`Total Reward Supply: ${ethers.utils.formatEther(rewardSupply)} ${symbol}`);
    console.log(`Minted Rewards: ${ethers.utils.formatEther(mintedRewards)} ${symbol}`);
    console.log(`Remaining Rewards: ${ethers.utils.formatEther(remainingRewards)} ${symbol}`);
    console.log(`Reward Pool Usage: ${mintedRewards.mul(10000).div(rewardSupply).toNumber() / 100}%`);
  } catch (error) {
    console.log("Could not fetch reward pool data");
  }

  // Summary
  console.log("\n📈 SUMMARY");
  console.log("=".repeat(50));
  
  let totalInContracts = ethers.BigNumber.from(0);
  for (const balance of contractBalances) {
    if (balance.balance !== "Error" && balance.balance !== "0") {
      totalInContracts = totalInContracts.add(ethers.utils.parseEther(balance.balance));
    }
  }
  
  const totalInContractsFormatted = ethers.utils.formatEther(totalInContracts);
  const contractsPercentage = totalSupply.gt(0) 
    ? totalInContracts.mul(10000).div(totalSupply).toNumber() / 100 
    : 0;
  
  console.log(`Total in Contracts: ${totalInContractsFormatted} ${symbol} (${contractsPercentage.toFixed(2)}%)`);
  console.log(`Circulating Supply: ${ethers.utils.formatEther(totalSupply.sub(totalInContracts))} ${symbol} (${(100 - contractsPercentage).toFixed(2)}%)`);

  // Warnings
  console.log("\n⚠️  IMPORTANT NOTES");
  console.log("-".repeat(50));
  
  if (addresses.rewardDistributor && addresses.rewardDistributor !== "") {
    const distributorBalance = await bogoToken.balanceOf(addresses.rewardDistributor);
    if (distributorBalance.gt(0)) {
      console.log(`✅ Reward Distributor has ${ethers.utils.formatEther(distributorBalance)} ${symbol} available`);
    } else {
      console.log(`⚠️  Reward Distributor has NO tokens - needs funding!`);
    }
  }
  
  console.log("\n💡 To get more detailed analysis:");
  console.log("1. Set all contract addresses in your .env file");
  console.log("2. Run specific analysis scripts for each contract");
  console.log("3. Use Etherscan to trace large holder addresses");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });