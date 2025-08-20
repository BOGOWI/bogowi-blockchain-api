const axios = require('axios');

// API base URL
const API_URL = 'http://localhost:3000';

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

// Test functions
async function testSystemStatus() {
  log('\n📊 Testing System Status...', 'blue');
  try {
    const response = await axios.get(`${API_URL}/status`);
    log('✅ System Status:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return false;
  }
}

async function testMintSingleTicket() {
  log('\n🎫 Testing Single Ticket Mint...', 'blue');
  try {
    const ticketId = Math.floor(Math.random() * 1000000);
    const response = await axios.post(`${API_URL}/mint`, {
      recipient: testUsers.user1,
      ticketId: ticketId,
      rewardBasisPoints: 500, // 5% reward
      metadataURI: `https://api.bogowi.com/metadata/${ticketId}`
    });
    log('✅ Ticket Minted:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return ticketId;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return null;
  }
}

async function testBatchMint() {
  log('\n🎫 Testing Batch Mint...', 'blue');
  try {
    const baseId = Math.floor(Math.random() * 1000000);
    const response = await axios.post(`${API_URL}/mint-batch`, {
      recipients: [testUsers.user1, testUsers.user2, testUsers.user1],
      ticketIds: [baseId, baseId + 1, baseId + 2],
      rewardBasisPoints: [300, 500, 1000], // 3%, 5%, 10%
      metadataURIs: [
        `https://api.bogowi.com/metadata/${baseId}`,
        `https://api.bogowi.com/metadata/${baseId + 1}`,
        `https://api.bogowi.com/metadata/${baseId + 2}`
      ]
    });
    log('✅ Batch Minted:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return [baseId, baseId + 1, baseId + 2];
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return [];
  }
}

async function testGetTicket(ticketId) {
  log(`\n🔍 Getting Ticket #${ticketId}...`, 'blue');
  try {
    const response = await axios.get(`${API_URL}/ticket/${ticketId}`);
    log('✅ Ticket Details:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return false;
  }
}

async function testGetUserTickets(address) {
  log(`\n👤 Getting tickets for ${address}...`, 'blue');
  try {
    const response = await axios.get(`${API_URL}/user/${address}/tickets`);
    log('✅ User Tickets:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return response.data.tickets;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return [];
  }
}

async function testRedeemTicket(ticketId) {
  log(`\n🔥 Redeeming Ticket #${ticketId}...`, 'blue');
  try {
    const response = await axios.post(`${API_URL}/redeem`, {
      ticketId: ticketId,
      recipient: testUsers.user1,
      rewardBasisPoints: 500,
      metadataURI: `https://api.bogowi.com/redeemed/${ticketId}`
    });
    log('✅ Ticket Redeemed (Burned):', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return false;
  }
}

async function testRegistry() {
  log('\n📚 Testing Registry...', 'blue');
  try {
    const response = await axios.get(`${API_URL}/registry`);
    log('✅ Registry Contents:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return false;
  }
}

async function testPauseUnpause() {
  log('\n⏸️  Testing Pause/Unpause...', 'blue');
  try {
    // Pause
    log('Pausing contract...', 'yellow');
    await axios.post(`${API_URL}/pause`);
    log('✅ Contract paused', 'green');
    
    // Try to mint while paused (should fail)
    log('Trying to mint while paused...', 'yellow');
    try {
      await axios.post(`${API_URL}/mint`, {
        recipient: testUsers.user1,
        ticketId: 999999,
        rewardBasisPoints: 500,
        metadataURI: 'test'
      });
      log('❌ Should have failed!', 'red');
    } catch (error) {
      log('✅ Mint correctly failed while paused', 'green');
    }
    
    // Unpause
    log('Unpausing contract...', 'yellow');
    await axios.post(`${API_URL}/unpause`);
    log('✅ Contract unpaused', 'green');
    
    return true;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return false;
  }
}

async function testStats() {
  log('\n📈 Getting Statistics...', 'blue');
  try {
    const response = await axios.get(`${API_URL}/stats`);
    log('✅ Statistics:', 'green');
    console.log(JSON.stringify(response.data, null, 2));
    return true;
  } catch (error) {
    log(`❌ Error: ${error.message}`, 'red');
    return false;
  }
}

// Main test runner
async function runAllTests() {
  log('\n🚀 Starting NFT API Tests', 'green');
  log('=' .repeat(50), 'blue');
  
  // Wait for API to be ready
  log('\nWaiting for API server...', 'yellow');
  let retries = 5;
  while (retries > 0) {
    try {
      await axios.get(`${API_URL}/status`);
      break;
    } catch (error) {
      retries--;
      if (retries === 0) {
        log('❌ API server not responding. Make sure to run: node scripts/test-nft-api-local.js', 'red');
        process.exit(1);
      }
      await new Promise(r => setTimeout(r, 2000));
    }
  }
  
  const results = [];
  
  // Run tests
  results.push(['System Status', await testSystemStatus()]);
  
  const ticketId = await testMintSingleTicket();
  results.push(['Single Mint', ticketId !== null]);
  
  const batchIds = await testBatchMint();
  results.push(['Batch Mint', batchIds.length > 0]);
  
  if (ticketId) {
    results.push(['Get Ticket', await testGetTicket(ticketId)]);
  }
  
  results.push(['Get User Tickets', (await testGetUserTickets(testUsers.user1)).length > 0]);
  
  results.push(['Registry', await testRegistry()]);
  
  results.push(['Pause/Unpause', await testPauseUnpause()]);
  
  // Test redeem (will burn the ticket)
  if (ticketId) {
    results.push(['Redeem Ticket', await testRedeemTicket(ticketId)]);
    
    // Verify it's burned
    log('\n🔍 Verifying ticket was burned...', 'yellow');
    try {
      await axios.get(`${API_URL}/ticket/${ticketId}`);
      log('❌ Ticket still exists!', 'red');
      results.push(['Burn Verification', false]);
    } catch (error) {
      log('✅ Ticket successfully burned', 'green');
      results.push(['Burn Verification', true]);
    }
  }
  
  results.push(['Statistics', await testStats()]);
  
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
  }
}

// Run tests
runAllTests().catch(console.error);