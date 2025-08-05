const { ethers } = require("hardhat");

async function main() {
  const txHash = "0x9bd44a46d330772738d021b3db9eef53a460044fbed0bec3d34be4351d8474bc";
  console.log(`Checking transaction: ${txHash}\n`);

  const provider = ethers.provider;
  
  // Get transaction
  const tx = await provider.getTransaction(txHash);
  console.log("From:", tx.from);
  console.log("To:", tx.to);
  console.log("Data:", tx.data);
  
  // Get receipt
  const receipt = await provider.getTransactionReceipt(txHash);
  console.log("\nStatus:", receipt.status === 1 ? "SUCCESS" : "FAILED");
  console.log("Logs:", receipt.logs.length);

  // Decode the function call
  const iface = new ethers.Interface([
    "function claimCustomReward(address recipient, uint256 amount, string reason)"
  ]);
  
  try {
    const decoded = iface.parseTransaction({ data: tx.data });
    console.log("\nFunction:", decoded.name);
    console.log("Recipient:", decoded.args[0]);
    console.log("Amount (wei):", decoded.args[1].toString());
    console.log("Amount (BOGO):", ethers.formatEther(decoded.args[1]));
    console.log("Reason:", decoded.args[2]);
  } catch (e) {
    console.log("Error decoding:", e.message);
  }

  // Check BOGO balance of recipient
  const bogoAbi = ["function balanceOf(address) view returns (uint256)"];
  const bogoToken = new ethers.Contract("0x49fc9939D8431371dD22658a8a969Ec798A26fFB", bogoAbi, provider);
  
  const recipientFromTx = "0x9318a4c47a052F1ef4351A4d86c7C903a056209D";
  const balance = await bogoToken.balanceOf(recipientFromTx);
  console.log(`\nRecipient ${recipientFromTx} balance:`, ethers.formatEther(balance), "BOGO");
}

main().catch(console.error);