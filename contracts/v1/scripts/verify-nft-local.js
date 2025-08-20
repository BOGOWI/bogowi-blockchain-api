const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🔍 Verifying NFT Infrastructure Deployment...\n");
  console.log("=" .repeat(60));

  // Load deployment info
  const network = hre.network.name;
  const deploymentPath = path.join(__dirname, `deployment-nft-${network}.json`);
  
  if (!fs.existsSync(deploymentPath)) {
    console.error("❌ No deployment found for network:", network);
    console.error("   Run deploy-nft-local.js first!");
    process.exit(1);
  }

  const deployment = JSON.parse(fs.readFileSync(deploymentPath, 'utf8'));
  console.log("📂 Loaded deployment from:", deploymentPath);
  console.log("📅 Deployed at:", deployment.timestamp);
  console.log();

  const [deployer] = await hre.ethers.getSigners();
  let allChecks = true;

  // 1. Verify contract deployments
  console.log("1️⃣  Verifying Contract Deployments");
  console.log("-" .repeat(40));
  
  for (const [name, address] of Object.entries(deployment.contracts)) {
    const code = await hre.ethers.provider.getCode(address);
    const isDeployed = code !== "0x";
    console.log(`  ${isDeployed ? '✅' : '❌'} ${name}: ${address}`);
    if (!isDeployed) allChecks = false;
  }

  // 2. Verify RoleManager setup
  console.log("\n2️⃣  Verifying RoleManager Setup");
  console.log("-" .repeat(40));
  
  const roleManager = await hre.ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);
  
  // Note: RoleManager doesn't expose a public function to check if contracts are registered
  // But we can verify that contracts are working by checking role assignments
  console.log("  ℹ️  RoleManager deployed at:", deployment.contracts.RoleManager);
  console.log("  ℹ️  Contracts registered during deployment:");
  console.log("     - NFTRegistry:", deployment.contracts.NFTRegistry);
  console.log("     - BOGOWITickets:", deployment.contracts.BOGOWITickets);

  // 3. Verify role assignments
  console.log("\n3️⃣  Verifying Role Assignments");
  console.log("-" .repeat(40));
  
  const roles = [
    { name: "REGISTRY_ADMIN_ROLE", hash: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("REGISTRY_ADMIN_ROLE")), address: deployment.roles.admin },
    { name: "CONTRACT_DEPLOYER_ROLE", hash: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("CONTRACT_DEPLOYER_ROLE")), address: deployment.roles.admin },
    { name: "NFT_MINTER_ROLE", hash: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("NFT_MINTER_ROLE")), address: deployment.roles.minter },
    { name: "BACKEND_ROLE", hash: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("BACKEND_ROLE")), address: deployment.roles.backend },
    { name: "ADMIN_ROLE", hash: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("ADMIN_ROLE")), address: deployment.roles.admin },
    { name: "PAUSER_ROLE", hash: hre.ethers.keccak256(hre.ethers.toUtf8Bytes("PAUSER_ROLE")), address: deployment.roles.admin }
  ];

  for (const role of roles) {
    const hasRole = await roleManager.hasRole(role.hash, role.address);
    const shortAddr = role.address.substring(0, 6) + "..." + role.address.substring(38);
    console.log(`  ${hasRole ? '✅' : '❌'} ${role.name} -> ${shortAddr}`);
    if (!hasRole) allChecks = false;
  }

  // 4. Verify NFTRegistry
  console.log("\n4️⃣  Verifying NFTRegistry");
  console.log("-" .repeat(40));
  
  const nftRegistry = await hre.ethers.getContractAt("NFTRegistry", deployment.contracts.NFTRegistry);
  
  // Check if BOGOWITickets is registered
  const isTicketsRegistered = await nftRegistry.isRegistered(deployment.contracts.BOGOWITickets);
  console.log(`  ${isTicketsRegistered ? '✅' : '❌'} BOGOWITickets registered in registry`);
  if (!isTicketsRegistered) allChecks = false;
  
  if (isTicketsRegistered) {
    const ticketInfo = await nftRegistry.getContractInfo(deployment.contracts.BOGOWITickets);
    console.log(`     Name: ${ticketInfo.name}`);
    console.log(`     Version: ${ticketInfo.version}`);
    console.log(`     Type: ${ticketInfo.contractType} (TICKET)`);
    console.log(`     Active: ${ticketInfo.isActive}`);
  }
  
  // Check total contracts
  const totalContracts = await nftRegistry.getContractCount();
  console.log(`  📊 Total contracts in registry: ${totalContracts}`);

  // 5. Verify BOGOWITickets
  console.log("\n5️⃣  Verifying BOGOWITickets");
  console.log("-" .repeat(40));
  
  const tickets = await hre.ethers.getContractAt("BOGOWITickets", deployment.contracts.BOGOWITickets);
  
  // Check basic properties
  const name = await tickets.name();
  const symbol = await tickets.symbol();
  const paused = await tickets.paused();
  const conservationDAO = await tickets.conservationDAO();
  
  console.log(`  ✅ Name: ${name}`);
  console.log(`  ✅ Symbol: ${symbol}`);
  console.log(`  ${!paused ? '✅' : '⚠️ '} Paused: ${paused}`);
  console.log(`  ✅ Conservation DAO: ${conservationDAO}`);
  
  // Check if minter can mint
  const minterRole = hre.ethers.keccak256(hre.ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
  const canMint = await roleManager.hasRole(minterRole, deployment.roles.minter);
  console.log(`  ${canMint ? '✅' : '❌'} Minter has minting permission`);
  if (!canMint) allChecks = false;

  // 6. Test contract interactions
  console.log("\n6️⃣  Testing Contract Interactions");
  console.log("-" .repeat(40));
  
  try {
    // Try to get contracts by type from registry
    const ticketContracts = await nftRegistry.getContractsByType(0); // TICKET type
    console.log(`  ✅ Can query contracts by type`);
    console.log(`     Found ${ticketContracts.length} ticket contract(s)`);
    
    // Try to get active contracts
    const activeContracts = await nftRegistry.getActiveContracts();
    console.log(`  ✅ Can query active contracts`);
    console.log(`     Found ${activeContracts.length} active contract(s)`);
  } catch (error) {
    console.log(`  ❌ Error testing interactions: ${error.message}`);
    allChecks = false;
  }

  // Final summary
  console.log("\n" + "=" .repeat(60));
  if (allChecks) {
    console.log("✅ All verification checks passed!");
    console.log("\n🎉 NFT Infrastructure is ready for testing!");
    console.log("\nYou can now:");
    console.log("  1. Mint test tickets: npm run test-mint-local");
    console.log("  2. Test registry operations: npm run test-registry-local");
  } else {
    console.log("❌ Some checks failed. Please review and fix issues.");
  }
  console.log("=" .repeat(60));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });