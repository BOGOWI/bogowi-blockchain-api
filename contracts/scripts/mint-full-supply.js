const { ethers } = require("hardhat");
require("dotenv").config();

async function main() {
  console.log("🚀 BOGO Token Full Supply Minting\n");

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

  // Check current supply and allocations
  console.log("📊 Current Token State:");
  const totalSupply = await bogoToken.totalSupply();
  const maxSupply = await bogoToken.MAX_SUPPLY();
  const daoMinted = await bogoToken.daoMinted();
  const businessMinted = await bogoToken.businessMinted();
  const rewardsMinted = await bogoToken.rewardsMinted();
  
  console.log(`Total Supply: ${ethers.utils.formatEther(totalSupply)} BOGO`);
  console.log(`Max Supply: ${ethers.utils.formatEther(maxSupply)} BOGO`);
  console.log(`DAO Minted: ${ethers.utils.formatEther(daoMinted)} BOGO`);
  console.log(`Business Minted: ${ethers.utils.formatEther(businessMinted)} BOGO`);
  console.log(`Rewards Minted: ${ethers.utils.formatEther(rewardsMinted)} BOGO`);

  // Check wallet balance
  const walletBalance = await bogoToken.balanceOf(deployer.address);
  console.log(`\nMain Wallet Balance: ${ethers.utils.formatEther(walletBalance)} BOGO\n`);

  // Check roles
  console.log("🔐 Checking Roles:");
  const DEFAULT_ADMIN_ROLE = await bogoToken.DEFAULT_ADMIN_ROLE();
  const DAO_ROLE = await bogoToken.DAO_ROLE();
  const BUSINESS_ROLE = await bogoToken.BUSINESS_ROLE();
  const MINTER_ROLE = await bogoToken.MINTER_ROLE();

  const hasAdminRole = await bogoToken.hasRole(DEFAULT_ADMIN_ROLE, deployer.address);
  const hasDaoRole = await bogoToken.hasRole(DAO_ROLE, deployer.address);
  const hasBusinessRole = await bogoToken.hasRole(BUSINESS_ROLE, deployer.address);
  const hasMinterRole = await bogoToken.hasRole(MINTER_ROLE, deployer.address);

  console.log(`DEFAULT_ADMIN_ROLE: ${hasAdminRole ? '✅ YES' : '❌ NO'}`);
  console.log(`DAO_ROLE: ${hasDaoRole ? '✅ YES' : '❌ NO'}`);
  console.log(`BUSINESS_ROLE: ${hasBusinessRole ? '✅ YES' : '❌ NO'}`);
  console.log(`MINTER_ROLE: ${hasMinterRole ? '✅ YES' : '❌ NO'}\n`);

  // Grant missing roles if we have admin
  if (hasAdminRole) {
    console.log("🔧 Granting Missing Roles:");
    
    if (!hasDaoRole) {
      console.log("Granting DAO_ROLE...");
      const tx1 = await bogoToken.grantRole(DAO_ROLE, deployer.address);
      await tx1.wait();
      console.log("✅ DAO_ROLE granted");
    }
    
    if (!hasBusinessRole) {
      console.log("Granting BUSINESS_ROLE...");
      const tx2 = await bogoToken.grantRole(BUSINESS_ROLE, deployer.address);
      await tx2.wait();
      console.log("✅ BUSINESS_ROLE granted");
    }
    
    if (!hasMinterRole) {
      console.log("Granting MINTER_ROLE...");
      const tx3 = await bogoToken.grantRole(MINTER_ROLE, deployer.address);
      await tx3.wait();
      console.log("✅ MINTER_ROLE granted");
    }
    console.log("");
  } else {
    console.log("❌ No admin role - cannot grant missing roles\n");
    if (!hasDaoRole || !hasBusinessRole) {
      console.log("⚠️  You need DAO_ROLE and BUSINESS_ROLE to mint tokens");
      return;
    }
  }

  // Calculate remaining allocations
  const daoAllocation = await bogoToken.DAO_ALLOCATION();
  const businessAllocation = await bogoToken.BUSINESS_ALLOCATION();
  const rewardsAllocation = await bogoToken.REWARDS_ALLOCATION();
  
  const remainingDao = daoAllocation - daoMinted;
  const remainingBusiness = businessAllocation - businessMinted;
  const remainingRewards = rewardsAllocation - rewardsMinted;

  console.log("💰 Available to Mint:");
  console.log(`DAO Allocation: ${ethers.utils.formatEther(remainingDao)} BOGO`);
  console.log(`Business Allocation: ${ethers.utils.formatEther(remainingBusiness)} BOGO`);
  console.log(`Rewards Allocation: ${ethers.utils.formatEther(remainingRewards)} BOGO`);
  
  const totalAvailable = remainingDao + remainingBusiness + remainingRewards;
  console.log(`Total Available: ${ethers.utils.formatEther(totalAvailable)} BOGO\n`);

  if (totalAvailable === 0n) {
    console.log("✅ All tokens already minted!");
    return;
  }

  // Ask for confirmation
  console.log(`⚠️  About to mint ${ethers.utils.formatEther(totalAvailable)} BOGO tokens to ${deployer.address}`);
  console.log("Type 'yes' to continue or anything else to cancel:");
  
  const readline = require('readline').createInterface({
    input: process.stdin,
    output: process.stdout
  });

  const answer = await new Promise(resolve => {
    readline.question('', resolve);
  });
  readline.close();

  if (answer.toLowerCase() !== 'yes') {
    console.log("Minting cancelled");
    return;
  }

  // Mint from each allocation
  console.log("\n🔨 Minting Tokens:");

  if (remainingDao > 0n) {
    console.log(`Minting ${ethers.utils.formatEther(remainingDao)} BOGO from DAO allocation...`);
    const tx1 = await bogoToken.mintFromDAO(deployer.address, remainingDao);
    await tx1.wait();
    console.log("✅ DAO allocation minted");
  }

  if (remainingBusiness > 0n) {
    console.log(`Minting ${ethers.utils.formatEther(remainingBusiness)} BOGO from Business allocation...`);
    const tx2 = await bogoToken.mintFromBusiness(deployer.address, remainingBusiness);
    await tx2.wait();
    console.log("✅ Business allocation minted");
  }

  if (remainingRewards > 0n) {
    console.log(`Minting ${ethers.utils.formatEther(remainingRewards)} BOGO from Rewards allocation...`);
    const tx3 = await bogoToken.mintFromRewards(deployer.address, remainingRewards);
    await tx3.wait();
    console.log("✅ Rewards allocation minted");
  }

  // Check final state
  console.log("\n📈 Final State:");
  const finalTotalSupply = await bogoToken.totalSupply();
  const finalWalletBalance = await bogoToken.balanceOf(deployer.address);
  
  console.log(`Total Supply: ${ethers.utils.formatEther(finalTotalSupply)} BOGO`);
  console.log(`Main Wallet Balance: ${ethers.utils.formatEther(finalWalletBalance)} BOGO`);
  
  if (finalTotalSupply === maxSupply) {
    console.log("\n🎉 SUCCESS! Full 1 billion BOGO supply minted!");
  } else {
    console.log(`\n⚠️  Total supply is ${ethers.utils.formatEther(finalTotalSupply)} / ${ethers.utils.formatEther(maxSupply)} BOGO`);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
