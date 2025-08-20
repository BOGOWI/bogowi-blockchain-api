const hre = require("hardhat");
const { ethers } = require("hardhat");
const fs = require("fs");
const path = require("path");

// Load deployment info
const deployment = JSON.parse(
  fs.readFileSync(path.join(__dirname, "deployment-nft-localhost.json"), 'utf8')
);

// Color output helpers
const colors = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m'
};

function log(message, color = 'reset') {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

async function testDirectContractInteraction() {
  log('\n🚀 Testing Direct Contract Interaction', 'green');
  log('=' .repeat(50), 'blue');
  
  // Get signers
  const [deployer, admin, minter, backend, user1, user2] = await ethers.getSigners();
  
  // Get contract instances
  const tickets = await ethers.getContractAt("BOGOWITickets", deployment.contracts.BOGOWITickets);
  const nftRegistry = await ethers.getContractAt("NFTRegistry", deployment.contracts.NFTRegistry);
  const roleManager = await ethers.getContractAt("RoleManager", deployment.contracts.RoleManager);
  
  log('\n📊 Contract Status:', 'blue');
  log(`  BOGOWITickets: ${deployment.contracts.BOGOWITickets}`, 'yellow');
  log(`  NFTRegistry: ${deployment.contracts.NFTRegistry}`, 'yellow');
  log(`  RoleManager: ${deployment.contracts.RoleManager}`, 'yellow');
  
  // Test 1: Check roles
  log('\n🔐 Checking Roles...', 'blue');
  const NFT_MINTER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("NFT_MINTER_ROLE"));
  const BACKEND_ROLE = ethers.keccak256(ethers.toUtf8Bytes("BACKEND_ROLE"));
  
  const minterHasRole = await roleManager.hasRole(NFT_MINTER_ROLE, minter.address);
  const backendHasRole = await roleManager.hasRole(BACKEND_ROLE, backend.address);
  
  log(`  Minter has NFT_MINTER_ROLE: ${minterHasRole}`, minterHasRole ? 'green' : 'red');
  log(`  Backend has BACKEND_ROLE: ${backendHasRole}`, backendHasRole ? 'green' : 'red');
  
  // Test 2: Mint a ticket
  log('\n🎫 Minting Ticket...', 'blue');
  const ticketId = Math.floor(Math.random() * 1000000);
  const deadline = Math.floor(Date.now() / 1000) + 3600;
  
  // Create signature
  const domain = {
    name: "BOGOWITickets",
    version: "1",
    chainId: 501,
    verifyingContract: await tickets.getAddress()
  };
  
  const types = {
    RedemptionData: [
      { name: "ticketId", type: "uint256" },
      { name: "recipient", type: "address" },
      { name: "rewardBasisPoints", type: "uint256" },
      { name: "metadataURI", type: "string" },
      { name: "deadline", type: "uint256" }
    ]
  };
  
  const value = {
    ticketId,
    recipient: user1.address,
    rewardBasisPoints: 500,
    metadataURI: `https://api.bogowi.com/metadata/${ticketId}`,
    deadline
  };
  
  const signature = await backend.signTypedData(domain, types, value);
  
  try {
    const tx = await tickets.connect(minter).mint(
      user1.address,
      ticketId,
      500,
      `https://api.bogowi.com/metadata/${ticketId}`,
      deadline,
      signature
    );
    await tx.wait();
    log(`  ✅ Ticket #${ticketId} minted to ${user1.address}`, 'green');
  } catch (error) {
    log(`  ❌ Minting failed: ${error.message}`, 'red');
  }
  
  // Test 3: Check ticket data
  log('\n🔍 Checking Ticket Data...', 'blue');
  try {
    const owner = await tickets.ownerOf(ticketId);
    const ticketData = await tickets.tickets(ticketId);
    const tokenURI = await tickets.tokenURI(ticketId);
    
    log(`  Owner: ${owner}`, 'green');
    log(`  Reward: ${ticketData.rewardBasisPoints / 100}%`, 'green');
    log(`  URI: ${tokenURI}`, 'green');
  } catch (error) {
    log(`  ❌ Error reading ticket: ${error.message}`, 'red');
  }
  
  // Test 4: Check registry
  log('\n📚 Checking Registry...', 'blue');
  try {
    const count = await nftRegistry.getContractCount();
    log(`  Registered contracts: ${count}`, 'green');
    
    for (let i = 0; i < count; i++) {
      const address = await nftRegistry.getContractAtIndex(i);
      const info = await nftRegistry.getContractInfo(address);
      log(`    - ${info.name} v${info.version} at ${address}`, 'yellow');
    }
  } catch (error) {
    log(`  ❌ Error reading registry: ${error.message}`, 'red');
  }
  
  // Test 5: Batch mint
  log('\n🎫 Batch Minting...', 'blue');
  const batchSize = 3;
  const batchIds = [];
  const recipients = [];
  const rewards = [];
  const uris = [];
  const deadlines = [];
  const signatures = [];
  
  for (let i = 0; i < batchSize; i++) {
    const id = Math.floor(Math.random() * 1000000);
    batchIds.push(id);
    recipients.push(i === 1 ? user2.address : user1.address);
    rewards.push([300, 500, 1000][i]);
    uris.push(`https://api.bogowi.com/metadata/${id}`);
    deadlines.push(deadline);
    
    const sig = await backend.signTypedData(domain, types, {
      ticketId: id,
      recipient: recipients[i],
      rewardBasisPoints: rewards[i],
      metadataURI: uris[i],
      deadline: deadline
    });
    signatures.push(sig);
  }
  
  try {
    const tx = await tickets.connect(minter).mintBatch(
      recipients,
      batchIds,
      rewards,
      uris,
      deadlines,
      signatures
    );
    await tx.wait();
    log(`  ✅ Batch minted ${batchSize} tickets`, 'green');
    log(`    IDs: ${batchIds.join(', ')}`, 'yellow');
  } catch (error) {
    log(`  ❌ Batch minting failed: ${error.message}`, 'red');
  }
  
  // Test 6: Check balances
  log('\n💰 Checking Balances...', 'blue');
  const balance1 = await tickets.balanceOf(user1.address);
  const balance2 = await tickets.balanceOf(user2.address);
  log(`  User1: ${balance1} tickets`, 'green');
  log(`  User2: ${balance2} tickets`, 'green');
  
  // Test 7: Total supply
  log('\n📈 Total Supply...', 'blue');
  const totalSupply = await tickets.totalSupply();
  log(`  Total tickets minted: ${totalSupply}`, 'green');
  
  log('\n✅ Direct contract testing complete!', 'green');
  log('=' .repeat(50), 'blue');
  
  return {
    ticketId,
    batchIds,
    totalSupply: totalSupply.toString()
  };
}

// Main function
async function main() {
  try {
    // Check if local node is running
    const provider = ethers.provider;
    const network = await provider.getNetwork();
    
    if (network.chainId !== 501n) {
      log('❌ Not connected to local network!', 'red');
      log(`Current chain ID: ${network.chainId}, expected: 501`, 'yellow');
      log('Please run: npx hardhat node', 'yellow');
      process.exit(1);
    }
    
    log('✅ Connected to local Hardhat network', 'green');
    
    // Run tests
    const results = await testDirectContractInteraction();
    
    log('\n📋 Summary:', 'blue');
    log(`  Minted ticket ID: ${results.ticketId}`, 'yellow');
    log(`  Batch ticket IDs: ${results.batchIds.join(', ')}`, 'yellow');
    log(`  Total supply: ${results.totalSupply}`, 'yellow');
    
    log('\n💡 You can now test your Go API with these ticket IDs!', 'green');
    log('  Use the contract addresses from deployment-nft-localhost.json', 'yellow');
    log('  Configure your Go API to use http://localhost:8545 as RPC', 'yellow');
    
  } catch (error) {
    log(`\n❌ Error: ${error.message}`, 'red');
    process.exit(1);
  }
}

main().catch(console.error);