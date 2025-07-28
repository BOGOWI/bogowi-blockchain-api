const { ethers } = require("hardhat");
require("dotenv").config();

async function main() {
  console.log("💰 Minting Remaining BOGO Tokens\n");

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
  console.log(`Rewards Minted: ${ethers.utils.formatEther(rewardsMinted)} BOGO\n`);

  // Get allocation constants
  const daoAllocation = await bogoToken.DAO_ALLOCATION();
  const businessAllocation = await bogoToken.BUSINESS_ALLOCATION();
  const rewardsAllocation = await bogoToken.REWARDS_ALLOCATION();
  
  const remainingDao = daoAllocation - daoMinted;
  const remainingBusiness = businessAllocation - businessMinted;
  const remainingRewards = rewardsAllocation - rewardsMinted;

  console.log("💰 Available to Mint:");
  console.log(`DAO Allocation: ${ethers.utils.formatEther(remainingDao)} BOGO`);
  console.log(`Business Allocation: ${ethers.utils.formatEther(remainingBusiness)} BOGO`);
  console.log(`Rewards Allocation: ${ethers.utils.formatEther(remainingRewards)} BOGO\n`);
  
  const totalAvailable = remainingDao + remainingBusiness + remainingRewards;
  console.log(`Total Available: ${ethers.utils.formatEther(totalAvailable)} BOGO\n`);

  if (totalAvailable === 0n) {
    console.log("✅ All tokens already minted!");
    return;
  }

  // Check required roles
  console.log("🔐 Checking Required Roles:");
  const DAO_ROLE = await bogoToken.DAO_ROLE();
  const BUSINESS_ROLE = await bogoToken.BUSINESS_ROLE();
  
  const hasDaoRole = await bogoToken.hasRole(DAO_ROLE, deployer.address);
  const hasBusinessRole = await bogoToken.hasRole(BUSINESS_ROLE, deployer.address);

  console.log(`DAO_ROLE: ${hasDaoRole ? '✅ YES' : '❌ NO'}`);
  console.log(`BUSINESS_ROLE: ${hasBusinessRole ? '✅ YES' : '❌ NO'}\n`);

  // Mint from each allocation
  console.log("🔨 Minting Tokens:\n");

  try {
    if (remainingDao > 0n) {
      if (!hasDaoRole) {
        console.log("❌ Cannot mint DAO allocation - missing DAO_ROLE");
      } else {
        console.log(`Minting ${ethers.utils.formatEther(remainingDao)} BOGO from DAO allocation...`);
        const tx1 = await bogoToken.mintFromDAO(deployer.address, remainingDao);
        const receipt1 = await tx1.wait();
        console.log(`✅ DAO allocation minted - Gas used: ${receipt1.gasUsed}`);
      }
    }

    if (remainingBusiness > 0n) {
      if (!hasBusinessRole) {
        console.log("❌ Cannot mint Business allocation - missing BUSINESS_ROLE");
      } else {
        console.log(`Minting ${ethers.utils.formatEther(remainingBusiness)} BOGO from Business allocation...`);
        const tx2 = await bogoToken.mintFromBusiness(deployer.address, remainingBusiness);
        const receipt2 = await tx2.wait();
        console.log(`✅ Business allocation minted - Gas used: ${receipt2.gasUsed}`);
      }
    }

    if (remainingRewards > 0n) {
      if (!hasDaoRole && !hasBusinessRole) {
        console.log("❌ Cannot mint Rewards allocation - missing DAO_ROLE or BUSINESS_ROLE");
      } else {
        console.log(`Minting ${ethers.utils.formatEther(remainingRewards)} BOGO from Rewards allocation...`);
        const tx3 = await bogoToken.mintFromRewards(deployer.address, remainingRewards);
        const receipt3 = await tx3.wait();
        console.log(`✅ Rewards allocation minted - Gas used: ${receipt3.gasUsed}`);
      }
    }

  } catch (error) {
    console.error("❌ Error during minting:", error.message);
    if (error.reason) {
      console.error("Reason:", error.reason);
    }
    return;
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
    const remaining = maxSupply - finalTotalSupply;
    console.log(`\n⚠️  ${ethers.utils.formatEther(remaining)} BOGO tokens still available to mint`);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
