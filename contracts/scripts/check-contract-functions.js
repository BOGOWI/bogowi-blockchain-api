const { ethers } = require("hardhat");
require("dotenv").config();

async function main() {
  console.log("🔍 Checking Contract Functions\n");

  const [deployer] = await ethers.getSigners();
  const bogoTokenAddress = process.env.BOGO_TOKEN_V2_ADDRESS;
  
  console.log("Deployer:", deployer.address);
  console.log("Contract:", bogoTokenAddress);
  console.log("=====================================\n");

  // Try with minimal interface first
  const minimalABI = [
    "function totalSupply() view returns (uint256)",
    "function balanceOf(address) view returns (uint256)",
    "function name() view returns (string)",
    "function symbol() view returns (string)",
    "function decimals() view returns (uint8)"
  ];

  const bogoToken = new ethers.Contract(bogoTokenAddress, minimalABI, deployer);

  try {
    console.log("📊 Basic Token Info:");
    const name = await bogoToken.name();
    const symbol = await bogoToken.symbol();
    const decimals = await bogoToken.decimals();
    const totalSupply = await bogoToken.totalSupply();
    const balance = await bogoToken.balanceOf(deployer.address);
    
    console.log(`Name: ${name}`);
    console.log(`Symbol: ${symbol}`);
    console.log(`Decimals: ${decimals}`);
    console.log(`Total Supply: ${ethers.utils.formatEther(totalSupply)} ${symbol}`);
    console.log(`Your Balance: ${ethers.utils.formatEther(balance)} ${symbol}\n`);
  } catch (error) {
    console.log("❌ Error reading basic info:", error.message);
  }

  // Try to check if it has allocation functions
  const extendedABI = [
    ...minimalABI,
    "function DAO_ALLOCATION() view returns (uint256)",
    "function BUSINESS_ALLOCATION() view returns (uint256)", 
    "function REWARDS_ALLOCATION() view returns (uint256)",
    "function daoMinted() view returns (uint256)",
    "function businessMinted() view returns (uint256)",
    "function rewardsMinted() view returns (uint256)"
  ];

  const extendedToken = new ethers.Contract(bogoTokenAddress, extendedABI, deployer);

  console.log("🔧 Checking Allocation Functions:");
  
  const functions = [
    "DAO_ALLOCATION",
    "BUSINESS_ALLOCATION", 
    "REWARDS_ALLOCATION",
    "daoMinted",
    "businessMinted", 
    "rewardsMinted"
  ];

  for (const func of functions) {
    try {
      const result = await extendedToken[func]();
      console.log(`✅ ${func}: ${ethers.utils.formatEther(result)} BOGO`);
    } catch (error) {
      console.log(`❌ ${func}: Function not available`);
    }
  }

  // Check if it has minting functions
  console.log("\n🪙 Checking Minting Functions:");
  const mintingABI = [
    ...minimalABI,
    "function mintFromDAO(address,uint256) external",
    "function mintFromBusiness(address,uint256) external",
    "function mintFromRewards(address,uint256) external"
  ];

  const mintingToken = new ethers.Contract(bogoTokenAddress, mintingABI, deployer);
  
  const mintFunctions = ["mintFromDAO", "mintFromBusiness", "mintFromRewards"];
  
  for (const func of mintFunctions) {
    try {
      // Just check if function exists by getting its fragment
      const fragment = mintingToken.interface.getFunction(func);
      console.log(`✅ ${func}: Available`);
    } catch (error) {
      console.log(`❌ ${func}: Not available`);
    }
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
