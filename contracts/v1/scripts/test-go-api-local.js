const axios = require('axios');
const { ethers } = require('hardhat');
const fs = require('fs');
const path = require('path');

// Load deployment info
const deployment = JSON.parse(
  fs.readFileSync(path.join(__dirname, 'deployment-nft-localhost.json'), 'utf8')
);

// Go API configuration
const GO_API_URL = process.env.API_URL || 'http://localhost:8080';
const GO_API_KEY = process.env.API_KEY || 'test-api-key';

// Test data
const testUsers = {
  user1: "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65",
  user2: "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc"
};

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

// Helper to create redemption signature for local testing
async function createLocalRedemptionSignature(ticketId, recipient, rewardBasisPoints, metadataURI, deadline) {
  const [, , , backend] = await ethers.getSigners(); // Get backend signer
  
  const domain = {
    name: "BOGOWITickets",
    version: "1",
    chainId: 501, // Camino testnet chain ID for local testing
    verifyingContract: deployment.contracts.BOGOWITickets
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
    recipient,
    rewardBasisPoints,
    metadataURI,
    deadline
  };

  return await backend.signTypedData(domain, types, value);
}

// Test API endpoints
async function testHealthCheck() {
  log('\n📊 Testing Health Check...', 'blue');
  try {
    const response = await axios.get(`${GO_API_URL}/health`);
    log('✅ Health Check:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return false;
  }
}

async function testMintTicket() {
  log('\n🎫 Testing Mint Ticket...', 'blue');
  try {
    const ticketId = Math.floor(Math.random() * 1000000);
    const deadline = Math.floor(Date.now() / 1000) + 3600; // 1 hour
    
    // Create signature locally for testing
    const signature = await createLocalRedemptionSignature(
      ticketId,
      testUsers.user1,
      500, // 5% reward
      `https://api.bogowi.com/metadata/${ticketId}`,
      deadline
    );
    
    const response = await axios.post(
      `${GO_API_URL}/api/v1/nft/tickets/mint`,
      {
        recipient: testUsers.user1,
        ticketId: ticketId.toString(),
        rewardBasisPoints: 500,
        metadataURI: `https://api.bogowi.com/metadata/${ticketId}`,
        deadline: deadline,
        signature: signature,
        network: "localhost" // Specify we're using local network
      },
      {
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': GO_API_KEY
        }
      }
    );
    
    log('✅ Ticket Minted via Go API:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return ticketId;
  } catch (error) {
    log(`❌ Error: ${error.response?.data || error.message}`, 'red');
    return null;
  }
}

async function testBatchMint() {
  log('\n🎫 Testing Batch Mint...', 'blue');
  try {
    const baseId = Math.floor(Math.random() * 1000000);
    const deadline = Math.floor(Date.now() / 1000) + 3600;
    
    const mintRequests = [];
    for (let i = 0; i < 3; i++) {
      const ticketId = baseId + i;
      const recipient = i === 1 ? testUsers.user2 : testUsers.user1;
      const rewardBasisPoints = [300, 500, 1000][i];
      
      const signature = await createLocalRedemptionSignature(
        ticketId,
        recipient,
        rewardBasisPoints,
        `https://api.bogowi.com/metadata/${ticketId}`,
        deadline
      );
      
      mintRequests.push({
        recipient: recipient,
        ticketId: ticketId.toString(),
        rewardBasisPoints: rewardBasisPoints,
        metadataURI: `https://api.bogowi.com/metadata/${ticketId}`,
        deadline: deadline,
        signature: signature
      });
    }
    
    const response = await axios.post(
      `${GO_API_URL}/api/v1/nft/tickets/mint-batch`,
      {
        requests: mintRequests,
        network: "localhost"
      },
      {
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': GO_API_KEY
        }
      }
    );
    
    log('✅ Batch Minted via Go API:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return [baseId, baseId + 1, baseId + 2];
  } catch (error) {
    log(`❌ Error: ${error.response?.data || error.message}`, 'red');
    return [];
  }
}

async function testGetTicket(ticketId) {
  log(`\n🔍 Getting Ticket #${ticketId}...`, 'blue');
  try {
    const response = await axios.get(
      `${GO_API_URL}/api/v1/nft/tickets/${ticketId}`,
      {
        params: { network: "localhost" },
        headers: { 'X-API-Key': GO_API_KEY }
      }
    );
    
    log('✅ Ticket Details from Go API:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.response?.data || error.message}`, 'red');
    return false;
  }
}

async function testGetUserTickets(address) {
  log(`\n👤 Getting tickets for ${address}...`, 'blue');
  try {
    const response = await axios.get(
      `${GO_API_URL}/api/v1/nft/tickets/user/${address}`,
      {
        params: { network: "localhost" },
        headers: { 'X-API-Key': GO_API_KEY }
      }
    );
    
    log('✅ User Tickets from Go API:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return response.data.tickets || [];
  } catch (error) {
    log(`❌ Error: ${error.response?.data || error.message}`, 'red');
    return [];
  }
}

async function testRedeemTicket(ticketId) {
  log(`\n🔥 Redeeming Ticket #${ticketId}...`, 'blue');
  try {
    const deadline = Math.floor(Date.now() / 1000) + 3600;
    const signature = await createLocalRedemptionSignature(
      ticketId,
      testUsers.user1,
      500,
      `https://api.bogowi.com/redeemed/${ticketId}`,
      deadline
    );
    
    const response = await axios.post(
      `${GO_API_URL}/api/v1/nft/tickets/redeem`,
      {
        ticketId: ticketId.toString(),
        recipient: testUsers.user1,
        rewardBasisPoints: 500,
        metadataURI: `https://api.bogowi.com/redeemed/${ticketId}`,
        deadline: deadline,
        signature: signature,
        network: "localhost"
      },
      {
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': GO_API_KEY
        }
      }
    );
    
    log('✅ Ticket Redeemed via Go API:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.response?.data || error.message}`, 'red');
    return false;
  }
}

async function testUploadImage() {
  log('\n🖼️  Testing Image Upload...', 'blue');
  try {
    // Create a test image buffer
    const testImageData = Buffer.from('Test image data for BOGOWI ticket');
    
    const FormData = require('form-data');
    const form = new FormData();
    form.append('image', testImageData, {
      filename: 'test-ticket.png',
      contentType: 'image/png'
    });
    form.append('ticketId', '12345');
    
    const response = await axios.post(
      `${GO_API_URL}/api/v1/nft/tickets/upload-image`,
      form,
      {
        headers: {
          ...form.getHeaders(),
          'X-API-Key': GO_API_KEY
        }
      }
    );
    
    log('✅ Image Uploaded:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.response?.data || error.message}`, 'red');
    return false;
  }
}

// Main test runner
async function runAllTests() {
  log('\n🚀 Testing Go API with Local Contracts', 'green');
  log('=' .repeat(50), 'blue');
  log(`API URL: ${GO_API_URL}`, 'yellow');
  log(`Local Contracts:`, 'yellow');
  log(`  NFTRegistry: ${deployment.contracts.NFTRegistry}`, 'yellow');
  log(`  BOGOWITickets: ${deployment.contracts.BOGOWITickets}`, 'yellow');
  log('=' .repeat(50), 'blue');
  
  const results = [];
  
  // Check if Go API is running
  log('\nChecking if Go API is running...', 'yellow');
  try {
    await axios.get(`${GO_API_URL}/health`);
    log('✅ Go API is running', 'green');
  } catch (error) {
    log('❌ Go API is not running!', 'red');
    log('Please start the Go API server first:', 'yellow');
    log('  cd ../.. && go run cmd/api/main.go', 'yellow');
    process.exit(1);
  }
  
  // Run tests
  results.push(['Health Check', await testHealthCheck()]);
  
  // Note: These tests assume your Go API can interact with the local blockchain
  // You may need to configure the Go API to use the local RPC endpoint
  log('\n⚠️  Note: Make sure your Go API is configured to use local blockchain:', 'yellow');
  log('  RPC URL: http://localhost:8545', 'yellow');
  log('  Contract addresses from deployment-nft-localhost.json', 'yellow');
  
  const ticketId = await testMintTicket();
  results.push(['Single Mint', ticketId !== null]);
  
  const batchIds = await testBatchMint();
  results.push(['Batch Mint', batchIds.length > 0]);
  
  if (ticketId) {
    results.push(['Get Ticket', await testGetTicket(ticketId)]);
  }
  
  results.push(['Get User Tickets', (await testGetUserTickets(testUsers.user1)).length > 0]);
  
  results.push(['Upload Image', await testUploadImage()]);
  
  if (ticketId) {
    results.push(['Redeem Ticket', await testRedeemTicket(ticketId)]);
  }
  
  // Print summary
  log('\n' + '=' .repeat(50), 'blue');
  log('📊 Test Summary', 'green');
  log('=' .repeat(50), 'blue');
  
  let passed = 0;
  let failed = 0;
  
  results.forEach(([test, result]) => {
    if (result) {
      log(`✅ ${test}`, 'green');
      passed++;
    } else {
      log(`❌ ${test}`, 'red');
      failed++;
    }
  });
  
  log('\n' + '=' .repeat(50), 'blue');
  log(`Total: ${passed} passed, ${failed} failed`, passed > failed ? 'green' : 'red');
  
  if (failed === 0) {
    log('\n🎉 All tests passed!', 'green');
  } else {
    log('\n⚠️  Some tests failed. Check configuration:', 'yellow');
    log('1. Ensure Go API is configured for local blockchain', 'yellow');
    log('2. Update contract addresses in Go API config', 'yellow');
    log('3. Make sure local blockchain is running', 'yellow');
  }
}

// Check if form-data is installed
const checkDependencies = async () => {
  try {
    require('form-data');
    require('axios');
  } catch (error) {
    log('📦 Installing required dependencies...', 'yellow');
    const { execSync } = require('child_process');
    execSync('npm install axios form-data', { stdio: 'inherit' });
  }
};

// Run tests
checkDependencies().then(() => {
  runAllTests().catch(console.error);
});