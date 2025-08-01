const hre = require("hardhat");
require("dotenv").config();

async function main() {
  console.log("Starting MultisigTreasury deployment...");
  
  // Get configuration
  const signer1 = process.env.MULTISIG_SIGNER_1;
  const signer2 = process.env.MULTISIG_SIGNER_2;
  const signer3 = process.env.MULTISIG_SIGNER_3;
  const threshold = process.env.MULTISIG_THRESHOLD || "2";
  
  // Validate inputs
  if (!signer1 || !signer2 || !signer3) {
    throw new Error("MULTISIG_SIGNER_1, MULTISIG_SIGNER_2, and MULTISIG_SIGNER_3 must be set in .env");
  }
  
  if (signer1 === "0x..." || signer2 === "0x..." || signer3 === "0x...") {
    throw new Error("Please set valid signer addresses in .env file");
  }
  
  const signers = [signer1, signer2, signer3];
  const uniqueSigners = [...new Set(signers)];
  
  if (uniqueSigners.length !== signers.length) {
    throw new Error("Duplicate signer addresses detected");
  }
  
  console.log("Configuration:");
  console.log("- Signers:", signers);
  console.log("- Threshold:", threshold);
  console.log("- Deploying from:", (await hre.ethers.getSigners())[0].address);
  
  // Get contract factory
  const MultisigTreasury = await hre.ethers.getContractFactory("MultisigTreasury");
  
  // Deploy contract
  console.log("\nDeploying MultisigTreasury...");
  const treasury = await MultisigTreasury.deploy(signers, parseInt(threshold));
  
  // Wait for deployment
  await treasury.deployed();
  
  console.log("MultisigTreasury deployed to:", treasury.address);
  console.log("Transaction hash:", treasury.deployTransaction.hash);
  
  // Wait for confirmations
  console.log("Waiting for 5 confirmations...");
  await treasury.deployTransaction.wait(5);
  
  console.log("Contract deployment confirmed!");
  
  // Verify deployment
  const deployedThreshold = await treasury.threshold();
  const deployedSignerCount = await treasury.signerCount();
  
  console.log("\n=== DEPLOYMENT VERIFICATION ===");
  console.log("Threshold:", deployedThreshold.toString());
  console.log("Signer count:", deployedSignerCount.toString());
  
  for (let i = 0; i < signers.length; i++) {
    const isSigner = await treasury.signers(signers[i]);
    console.log(`Signer ${i + 1} (${signers[i]}):`, isSigner.isSigner ? "✓" : "✗");
  }
  
  console.log("\n=== DEPLOYMENT SUMMARY ===");
  console.log("Contract Address:", treasury.address);
  console.log("Network:", hre.network.name);
  console.log("Signers:", signers);
  console.log("Threshold:", threshold, "of", signers.length);
  console.log("========================\n");
  
  // Save deployment info
  const fs = require("fs");
  const deploymentInfo = {
    network: hre.network.name,
    contractAddress: treasury.address,
    signers: signers,
    threshold: parseInt(threshold),
    deploymentTx: treasury.deployTransaction.hash,
    deployedAt: new Date().toISOString(),
    deployer: (await hre.ethers.getSigners())[0].address,
    contractParams: {
      executionDelay: "1 hour",
      transactionExpiry: "7 days",
      maxSigners: 20,
      maxGasLimit: 5000000
    }
  };
  
  const filename = `treasury-deployment-${hre.network.name}-${Date.now()}.json`;
  fs.writeFileSync(filename, JSON.stringify(deploymentInfo, null, 2));
  
  console.log("Deployment info saved to:", filename);
  console.log("\nNEXT STEPS:");
  console.log("1. Add to .env: MULTISIG_TREASURY_ADDRESS=" + treasury.address);
  console.log("2. Deploy BOGORewardDistributor with this treasury address");
  console.log("3. Each signer should verify they can connect to the treasury");
  console.log("4. Test a multisig transaction before using in production");
  console.log("\nIMPORTANT NOTES:");
  console.log("- Transactions require " + threshold + " confirmations to execute");
  console.log("- There is a 1 hour delay after confirmation before execution");
  console.log("- Transactions expire after 7 days if not executed");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });