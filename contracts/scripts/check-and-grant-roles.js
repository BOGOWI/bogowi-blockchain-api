const { ethers } = require("hardhat");
require("dotenv").config();

async function main() {
  console.log("🔐 Checking and Granting Roles\n");

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

  // Get role constants
  console.log("📋 Getting Role Constants:");
  const DEFAULT_ADMIN_ROLE = await bogoToken.DEFAULT_ADMIN_ROLE();
  const DAO_ROLE = await bogoToken.DAO_ROLE();
  const BUSINESS_ROLE = await bogoToken.BUSINESS_ROLE();
  const MINTER_ROLE = await bogoToken.MINTER_ROLE();
  const PAUSER_ROLE = await bogoToken.PAUSER_ROLE();

  console.log(`DEFAULT_ADMIN_ROLE: ${DEFAULT_ADMIN_ROLE}`);
  console.log(`DAO_ROLE: ${DAO_ROLE}`);
  console.log(`BUSINESS_ROLE: ${BUSINESS_ROLE}`);
  console.log(`MINTER_ROLE: ${MINTER_ROLE}`);
  console.log(`PAUSER_ROLE: ${PAUSER_ROLE}\n`);

  // Check current roles
  console.log("🔍 Current Roles:");
  const hasAdminRole = await bogoToken.hasRole(DEFAULT_ADMIN_ROLE, deployer.address);
  const hasDaoRole = await bogoToken.hasRole(DAO_ROLE, deployer.address);
  const hasBusinessRole = await bogoToken.hasRole(BUSINESS_ROLE, deployer.address);
  const hasMinterRole = await bogoToken.hasRole(MINTER_ROLE, deployer.address);
  const hasPauserRole = await bogoToken.hasRole(PAUSER_ROLE, deployer.address);

  console.log(`DEFAULT_ADMIN_ROLE: ${hasAdminRole ? '✅ YES' : '❌ NO'}`);
  console.log(`DAO_ROLE: ${hasDaoRole ? '✅ YES' : '❌ NO'}`);
  console.log(`BUSINESS_ROLE: ${hasBusinessRole ? '✅ YES' : '❌ NO'}`);
  console.log(`MINTER_ROLE: ${hasMinterRole ? '✅ YES' : '❌ NO'}`);
  console.log(`PAUSER_ROLE: ${hasPauserRole ? '✅ YES' : '❌ NO'}\n`);

  if (!hasAdminRole) {
    console.log("❌ No admin role - cannot grant other roles");
    return;
  }

  // Grant missing roles
  console.log("🔧 Granting Missing Roles:");
  
  try {
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

    if (!hasPauserRole) {
      console.log("Granting PAUSER_ROLE...");
      const tx4 = await bogoToken.grantRole(PAUSER_ROLE, deployer.address);
      await tx4.wait();
      console.log("✅ PAUSER_ROLE granted");
    }

    console.log("\n🎉 All roles granted successfully!");
    
  } catch (error) {
    console.error("❌ Error granting roles:", error.message);
    if (error.reason) {
      console.error("Reason:", error.reason);
    }
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
