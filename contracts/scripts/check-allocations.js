const { ethers } = require("hardhat");
require("dotenv").config();

async function main() {
  console.log("🔍 Checking Allocation Constants\n");

  const [deployer] = await ethers.getSigners();
  const bogoTokenAddress = process.env.BOGO_TOKEN_V2_ADDRESS;
  
  if (!bogoTokenAddress) {
    throw new Error("BOGO_TOKEN_V2_ADDRESS not set in .env");
  }

  console.log("Deployer/Main Wallet:", deployer.address);
  console.log("BOGO Token Address:", bogoTokenAddress);
  console.log("=====================================\n");

  // Get contract instance
  const bogoToken = await ethers.getContractAt("BOGOTokenV2", bogoTokenAddress);

  try {
    console.log("📊 Basic Token Info:");
    const name = await bogoToken.name();
    const symbol = await bogoToken.symbol();
    const totalSupply = await bogoToken.totalSupply();
    const balance = await bogoToken.balanceOf(deployer.address);
    
    console.log(`Name: ${name}`);
    console.log(`Symbol: ${symbol}`);
    console.log(`Total Supply: ${ethers.utils.formatEther(totalSupply)} BOGO`);
    console.log(`Your Balance: ${ethers.utils.formatEther(balance)} BOGO\n`);
  } catch (error) {
    console.error("❌ Error getting basic info:", error.message);
  }

  // Test each allocation constant individually
  console.log("🔧 Testing Allocation Constants:");
  
  try {
    const maxSupply = await bogoToken.MAX_SUPPLY();
    console.log(`✅ MAX_SUPPLY: ${ethers.utils.formatEther(maxSupply)} BOGO`);
  } catch (error) {
    console.log(`❌ MAX_SUPPLY: ${error.message}`);
  }

  try {
    const daoAllocation = await bogoToken.DAO_ALLOCATION();
    console.log(`✅ DAO_ALLOCATION: ${ethers.utils.formatEther(daoAllocation)} BOGO`);
  } catch (error) {
    console.log(`❌ DAO_ALLOCATION: ${error.message}`);
  }

  try {
    const businessAllocation = await bogoToken.BUSINESS_ALLOCATION();
    console.log(`✅ BUSINESS_ALLOCATION: ${ethers.utils.formatEther(businessAllocation)} BOGO`);
  } catch (error) {
    console.log(`❌ BUSINESS_ALLOCATION: ${error.message}`);
  }

  try {
    const rewardsAllocation = await bogoToken.REWARDS_ALLOCATION();
    console.log(`✅ REWARDS_ALLOCATION: ${ethers.utils.formatEther(rewardsAllocation)} BOGO`);
  } catch (error) {
    console.log(`❌ REWARDS_ALLOCATION: ${error.message}`);
  }

  console.log("\n🔧 Testing Minted Amounts:");
  
  try {
    const daoMinted = await bogoToken.daoMinted();
    console.log(`✅ daoMinted: ${ethers.utils.formatEther(daoMinted)} BOGO`);
  } catch (error) {
    console.log(`❌ daoMinted: ${error.message}`);
  }

  try {
    const businessMinted = await bogoToken.businessMinted();
    console.log(`✅ businessMinted: ${ethers.utils.formatEther(businessMinted)} BOGO`);
  } catch (error) {
    console.log(`❌ businessMinted: ${error.message}`);
  }

  try {
    const rewardsMinted = await bogoToken.rewardsMinted();
    console.log(`✅ rewardsMinted: ${ethers.utils.formatEther(rewardsMinted)} BOGO`);
  } catch (error) {
    console.log(`❌ rewardsMinted: ${error.message}`);
  }

  console.log("\n🔧 Testing Remaining Allocation Functions:");
  
  try {
    const remainingDao = await bogoToken.getRemainingDAOAllocation();
    console.log(`✅ getRemainingDAOAllocation: ${ethers.utils.formatEther(remainingDao)} BOGO`);
  } catch (error) {
    console.log(`❌ getRemainingDAOAllocation: ${error.message}`);
  }

  try {
    const remainingBusiness = await bogoToken.getRemainingBusinessAllocation();
    console.log(`✅ getRemainingBusinessAllocation: ${ethers.utils.formatEther(remainingBusiness)} BOGO`);
  } catch (error) {
    console.log(`❌ getRemainingBusinessAllocation: ${error.message}`);
  }

  try {
    const remainingRewards = await bogoToken.getRemainingRewardsAllocation();
    console.log(`✅ getRemainingRewardsAllocation: ${ethers.utils.formatEther(remainingRewards)} BOGO`);
  } catch (error) {
    console.log(`❌ getRemainingRewardsAllocation: ${error.message}`);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
