# Smart Contract Monitoring Setup Guide

## Monitoring Options & Pricing

### 1. **Tenderly** 
Website: https://tenderly.co

#### Pricing Tiers:
- **Free Tier** ✅
  - 90 days data retention
  - Real-time monitoring
  - 10 alert rules
  - 500 simulations/month
  - 3 team members
  - **Perfect for starting out!**

- **Developer ($120/month)**
  - Unlimited data retention
  - 100 alert rules
  - 5,000 simulations/month
  - 10 team members

- **Pro ($500/month)**
  - Everything in Developer
  - Priority support
  - Custom integrations

### 2. **OpenZeppelin Defender**
Website: https://defender.openzeppelin.com

#### Pricing:
- **Free Tier** ✅
  - 10 Sentinels (monitors)
  - 120 Autotask runs/month
  - 5 team members
  - Basic monitoring
  - **Good free option!**

- **Starter ($250/month)**
  - 30 Sentinels
  - 1,500 Autotask runs
  - Advanced features

### 3. **Alchemy Notify**
Website: https://www.alchemy.com/notify

#### Pricing:
- **Free Tier** ✅
  - Unlimited webhooks
  - Address activity notifications
  - NFT activity tracking
  - **Great for basic alerts!**

### 4. **Forta**
Website: https://forta.org

#### Pricing:
- **Free** ✅
  - Open source threat detection
  - Community bots
  - Basic alerts
  - **Best for security monitoring!**

## Setting Up Tenderly (Free Tier)

### Step 1: Create Account
1. Go to https://tenderly.co
2. Sign up for free account
3. Verify email

### Step 2: Create Project
```bash
# Install Tenderly CLI
curl https://raw.githubusercontent.com/Tenderly/tenderly-cli/master/scripts/install-linux.sh | sh

# Or on macOS
brew tap tenderly/tenderly
brew install tenderly

# Login
tenderly login

# Initialize in your project
cd contracts/v1
tenderly init
```

### Step 3: Configure tenderly.yaml
Create `tenderly.yaml` in your project root:

```yaml
account_id: your-username
project_slug: bogowi-nft
contracts:
  - network: 500  # Camino Mainnet
    contracts:
      - address: "0x..."  # Your NFTRegistry address
        name: NFTRegistry
      - address: "0x..."  # Your BOGOWITickets address
        name: BOGOWITickets
  - network: 501  # Columbus Testnet
    contracts:
      - address: "0x..."
        name: NFTRegistry
      - address: "0x..."
        name: BOGOWITickets
```

### Step 4: Push Contracts
```bash
# Push your contracts to Tenderly
tenderly push

# Or verify specific contracts
tenderly verify --network 500 0xYourContractAddress ContractName
```

### Step 5: Set Up Alerts

#### Via Web Dashboard:
1. Go to your project dashboard
2. Click "Alerts" → "New Alert"
3. Choose alert type

#### Common Alert Rules:

**1. Large Value Transfers**
```javascript
// Alert when ticket minted with high reward
{
  "type": "EVENT",
  "contract": "BOGOWITickets",
  "event": "TicketMinted",
  "condition": "rewardBasisPoints > 1000" // > 10%
}
```

**2. Failed Transactions**
```javascript
{
  "type": "TRANSACTION",
  "status": "failed",
  "to": ["0xYourContractAddress"]
}
```

**3. Role Changes**
```javascript
{
  "type": "EVENT",
  "contract": "RoleManager",
  "event": "RoleGranted"
}
```

**4. Pause Events**
```javascript
{
  "type": "EVENT",
  "event": "Paused"
}
```

### Step 6: Set Up Notifications

#### Webhook Integration:
```javascript
// webhook endpoint example
app.post('/tenderly-webhook', (req, res) => {
  const alert = req.body;
  
  if (alert.type === 'ALERT') {
    // Send to Discord/Slack/Email
    notifyTeam(alert);
  }
  
  res.sendStatus(200);
});
```

#### Discord Integration:
```javascript
// Discord webhook URL
const webhookUrl = "https://discord.com/api/webhooks/...";

// Send alert to Discord
async function notifyDiscord(alert) {
  await fetch(webhookUrl, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      content: `🚨 Contract Alert: ${alert.message}`,
      embeds: [{
        title: alert.name,
        description: alert.description,
        color: 0xff0000,
        fields: [
          { name: "Contract", value: alert.contract },
          { name: "Network", value: alert.network },
          { name: "Transaction", value: alert.txHash }
        ]
      }]
    })
  });
}
```

## Quick Setup Script

Create `scripts/setup-monitoring.js`:

```javascript
const axios = require('axios');
require('dotenv').config();

async function setupTenderlyMonitoring() {
  const TENDERLY_API_KEY = process.env.TENDERLY_API_KEY;
  const TENDERLY_PROJECT = process.env.TENDERLY_PROJECT;
  const TENDERLY_ACCOUNT = process.env.TENDERLY_ACCOUNT;
  
  // Your deployed contracts
  const contracts = {
    NFTRegistry: process.env.NFT_REGISTRY_ADDRESS,
    BOGOWITickets: process.env.TICKETS_ADDRESS
  };
  
  // API endpoint
  const baseURL = `https://api.tenderly.co/api/v1/account/${TENDERLY_ACCOUNT}/project/${TENDERLY_PROJECT}`;
  
  // Add contracts to monitoring
  for (const [name, address] of Object.entries(contracts)) {
    try {
      const response = await axios.post(
        `${baseURL}/contracts`,
        {
          address: address,
          network_id: "500", // Camino mainnet
          display_name: name
        },
        {
          headers: {
            'X-Access-Key': TENDERLY_API_KEY,
            'Content-Type': 'application/json'
          }
        }
      );
      console.log(`✅ Added ${name} to monitoring`);
    } catch (error) {
      console.error(`❌ Failed to add ${name}:`, error.message);
    }
  }
  
  // Create alerts
  const alerts = [
    {
      name: "High Value Ticket Minted",
      type: "EVENT_LOG",
      contracts: [contracts.BOGOWITickets],
      conditions: {
        event: "TicketMinted",
        filter: "rewardBasisPoints > 1000"
      }
    },
    {
      name: "Contract Paused",
      type: "EVENT_LOG",
      contracts: Object.values(contracts),
      conditions: {
        event: "Paused"
      }
    },
    {
      name: "Failed Transactions",
      type: "TRANSACTION",
      contracts: Object.values(contracts),
      conditions: {
        status: "failed"
      }
    }
  ];
  
  for (const alert of alerts) {
    try {
      await axios.post(
        `${baseURL}/alerts`,
        alert,
        {
          headers: {
            'X-Access-Key': TENDERLY_API_KEY,
            'Content-Type': 'application/json'
          }
        }
      );
      console.log(`✅ Created alert: ${alert.name}`);
    } catch (error) {
      console.error(`❌ Failed to create alert:`, error.message);
    }
  }
  
  console.log("\n📊 Monitoring setup complete!");
  console.log(`View dashboard: https://dashboard.tenderly.co/${TENDERLY_ACCOUNT}/${TENDERLY_PROJECT}`);
}

// Run setup
setupTenderlyMonitoring().catch(console.error);
```

## Environment Variables

Add to `.env`:

```bash
# Tenderly
TENDERLY_API_KEY=your_api_key_here
TENDERLY_ACCOUNT=your_account
TENDERLY_PROJECT=bogowi-nft

# Contracts (after deployment)
NFT_REGISTRY_ADDRESS=0x...
TICKETS_ADDRESS=0x...

# Notifications
DISCORD_WEBHOOK=https://discord.com/api/webhooks/...
SLACK_WEBHOOK=https://hooks.slack.com/services/...
```

## What to Monitor

### Critical Events
- ✅ **Contract Paused/Unpaused**
- ✅ **Role Changes** (especially admin roles)
- ✅ **Large Batch Mints** (gas spike detection)
- ✅ **Failed Transactions**
- ✅ **Registry Updates** (new contracts added)

### Performance Metrics
- ✅ **Gas Usage Trends**
- ✅ **Transaction Volume**
- ✅ **Unique Users**
- ✅ **Error Rate**

### Security Alerts
- ✅ **Unusual Activity Patterns**
- ✅ **Repeated Failed Attempts**
- ✅ **Unexpected Contract Calls**
- ✅ **Role Escalation Attempts**

## Free Monitoring Stack Recommendation

For starting out with **zero cost**:

1. **Tenderly Free Tier**
   - Real-time monitoring
   - 10 alerts
   - Transaction debugging

2. **Alchemy Notify**
   - Webhook notifications
   - Address activity

3. **Forta**
   - Security monitoring
   - Threat detection

4. **Custom Dashboard**
   ```javascript
   // Simple monitoring script
   setInterval(async () => {
     const balance = await provider.getBalance(contract);
     const totalSupply = await contract.totalSupply();
     
     console.log(`Balance: ${balance}`);
     console.log(`Total NFTs: ${totalSupply}`);
     
     // Alert if something unusual
     if (balance < threshold) {
       sendAlert("Low balance warning!");
     }
   }, 60000); // Check every minute
   ```

## Quick Start Commands

```bash
# 1. Install Tenderly CLI
brew install tenderly

# 2. Login
tenderly login

# 3. Initialize
tenderly init

# 4. Push contracts
tenderly push

# 5. View dashboard
open https://dashboard.tenderly.co

# 6. Run monitoring setup
node scripts/setup-monitoring.js
```

## Conclusion

**Start with Tenderly's free tier** - it's more than enough for initial monitoring and has:
- No credit card required
- Real-time alerts
- Great debugging tools
- Clean dashboard

Once you exceed free limits (unlikely initially), you can:
1. Upgrade to paid tier
2. Switch to OpenZeppelin Defender
3. Build custom monitoring

The free tier should handle your needs for months!