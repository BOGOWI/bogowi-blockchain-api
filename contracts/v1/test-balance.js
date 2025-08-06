const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  const deploymentPath = path.join(__dirname, "scripts/deployment-columbus.json");
  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  
  const token = await hre.ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  const recipient = "0xFC95F34a265b81b49bFae997F9146eE3dc821bfB";
  
  console.log("Checking balance for:", recipient);
  console.log("Token contract:", deployment.contracts.BOGOToken);
  
  const balance = await token.balanceOf(recipient);
  console.log("Balance:", hre.ethers.formatEther(balance), "BOGO");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
