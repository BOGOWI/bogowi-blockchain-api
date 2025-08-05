const hre = require("hardhat");

async function main() {
  // Get tx hash from environment variable instead
  const txHash = process.env.TX_HASH || "0x9bd44a46d330772738d021b3db9eef53a460044fbed0bec3d34be4351d8474bc";
  if (!txHash) {
    console.log("Please set TX_HASH environment variable");
    process.exit(1);
  }

  console.log(`🔍 Checking transaction: ${txHash}\n`);

  try {
    // Get transaction details
    const tx = await hre.ethers.provider.getTransaction(txHash);
    if (!tx) {
      console.log("❌ Transaction not found!");
      return;
    }

    console.log("📝 Transaction Details:");
    console.log("From:", tx.from);
    console.log("To:", tx.to);
    console.log("Value:", hre.ethers.formatEther(tx.value), "CAM");
    console.log("Gas Price:", hre.ethers.formatUnits(tx.gasPrice, "gwei"), "gwei");
    console.log("Gas Limit:", tx.gasLimit.toString());
    console.log("Nonce:", tx.nonce);
    console.log("Block:", tx.blockNumber);

    // Get transaction receipt
    const receipt = await hre.ethers.provider.getTransactionReceipt(txHash);
    if (!receipt) {
      console.log("\n⏳ Transaction is pending...");
      return;
    }

    console.log("\n📋 Transaction Receipt:");
    console.log("Status:", receipt.status === 1 ? "✅ Success" : "❌ Failed");
    console.log("Gas Used:", receipt.gasUsed.toString());
    if (receipt.effectiveGasPrice) {
      console.log("Effective Gas Price:", hre.ethers.formatUnits(receipt.effectiveGasPrice, "gwei"), "gwei");
    }
    console.log("Block Number:", receipt.blockNumber);

    // Decode the transaction data
    if (tx.data && tx.data !== "0x") {
      console.log("\n🔧 Decoding Transaction Data:");
      
      // Load ABI
      const RewardDistributor = await hre.ethers.getContractFactory("BOGORewardDistributor");
      const iface = RewardDistributor.interface;
      
      try {
        const decoded = iface.parseTransaction({ data: tx.data });
        console.log("Function:", decoded.name);
        console.log("Arguments:");
        
        if (decoded.name === "claimCustomReward") {
          console.log("  - Recipient:", decoded.args[0]);
          console.log("  - Amount:", hre.ethers.formatEther(decoded.args[1]), "BOGO");
          console.log("  - Reason:", decoded.args[2]);
        } else {
          decoded.args.forEach((arg, i) => {
            console.log(`  - Arg ${i}:`, arg.toString());
          });
        }
      } catch (e) {
        console.log("Could not decode transaction data");
      }
    }

    // Check logs
    if (receipt.logs && receipt.logs.length > 0) {
      console.log("\n📜 Event Logs:");
      
      // Load contract to decode logs
      const fs = require("fs");
      const path = require("path");
      const deploymentPath = path.join(__dirname, `deployment-camino.json`);
      
      if (fs.existsSync(deploymentPath)) {
        const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
        const rewardDistributor = await hre.ethers.getContractAt("BOGORewardDistributor", deployment.contracts.BOGORewardDistributor);
        
        for (const log of receipt.logs) {
          try {
            const parsed = rewardDistributor.interface.parseLog(log);
            if (parsed) {
              console.log(`\nEvent: ${parsed.name}`);
              if (parsed.name === "RewardClaimed") {
                console.log("  - Wallet:", parsed.args[0]);
                console.log("  - Template/Reason:", parsed.args[1]);
                console.log("  - Amount:", hre.ethers.formatEther(parsed.args[2]), "BOGO");
              } else {
                parsed.args.forEach((arg, i) => {
                  console.log(`  - ${parsed.fragment.inputs[i].name}:`, arg.toString());
                });
              }
            }
          } catch (e) {
            // Not a RewardDistributor event
          }
        }
      }
    }

    // If failed, try to get revert reason
    if (receipt.status === 0) {
      console.log("\n❌ Transaction Failed!");
      try {
        // Try to simulate the transaction to get revert reason
        const tx2 = {
          from: tx.from,
          to: tx.to,
          data: tx.data,
          value: tx.value,
          gasLimit: tx.gasLimit,
          gasPrice: tx.gasPrice
        };
        
        await hre.ethers.provider.call(tx2, tx.blockNumber - 1);
      } catch (error) {
        console.log("Revert reason:", error.reason || error.message);
      }
    }

  } catch (error) {
    console.error("Error:", error.message);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });