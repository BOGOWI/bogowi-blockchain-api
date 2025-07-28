const { ethers } = require("hardhat");
require("dotenv").config();

async function main() {
  const bogoTokenAddress = process.env.BOGO_TOKEN_V2_ADDRESS;
  
  if (!bogoTokenAddress) {
    throw new Error("BOGO_TOKEN_V2_ADDRESS not set");
  }
  
  console.log("Checking BOGO Token at:", bogoTokenAddress);
  console.log("=====================================\n");
  
  try {
    // Get contract instance with minimal ABI
    const bogoToken = await ethers.getContractAt(
      [
        "function totalSupply() view returns (uint256)",
        "function decimals() view returns (uint8)",
        "function name() view returns (string)",
        "function symbol() view returns (string)",
        "function balanceOf(address) view returns (uint256)"
      ],
      bogoTokenAddress
    );
    
    // Get basic token info
    const name = await bogoToken.name();
    const symbol = await bogoToken.symbol();
    const decimals = await bogoToken.decimals();
    const totalSupply = await bogoToken.totalSupply();
    
    console.log("Token Name:", name);
    console.log("Token Symbol:", symbol);
    console.log("Decimals:", decimals);
    console.log("Total Supply (raw):", totalSupply.toString());
    console.log("Total Supply (formatted):", ethers.utils.formatUnits(totalSupply, decimals), symbol);
    
    // Check some key addresses
    const addresses = [
      { name: "Reward Distributor", address: process.env.REWARD_DISTRIBUTOR_V2_ADDRESS },
      { name: "Multisig", address: process.env.MULTISIG_ADDRESS },
      { name: "Backend Wallet", address: process.env.BACKEND_WALLET_ADDRESS }
    ];
    
    console.log("\n=== KEY WALLET BALANCES ===");
    for (const { name, address } of addresses) {
      if (address) {
        const balance = await bogoToken.balanceOf(address);
        console.log(`${name} (${address}):`);
        console.log(`  Balance: ${ethers.utils.formatUnits(balance, decimals)} ${symbol}`);
      }
    }
    
  } catch (error) {
    console.error("❌ Error:", error.message);
    process.exit(1);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });