const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  const deploymentPath = path.join(__dirname, "scripts/deployment-camino.json");
  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  
  const token = await hre.ethers.getContractAt("BOGOToken", deployment.contracts.BOGOToken);
  const address = "0x18c4ACE2cbD28e9F31b55Aef1A3BFC4EC12cE956";
  
  console.log("Checking balance for:", address);
  console.log("Token contract:", deployment.contracts.BOGOToken);
  console.log("Network:", hre.network.name);
  
  const balance = await token.balanceOf(address);
  console.log("Balance:", hre.ethers.formatEther(balance), "BOGO");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
