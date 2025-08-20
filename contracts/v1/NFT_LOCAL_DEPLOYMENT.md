# NFT Infrastructure Local Deployment Guide

## Overview
This guide walks you through deploying and testing the BOGOWI NFT infrastructure (NFTRegistry + BOGOWITickets) in a local environment before mainnet deployment.

## Prerequisites

1. **Install dependencies**:
   ```bash
   cd contracts/v1
   npm install
   ```

2. **Start local Hardhat node** (in a separate terminal):
   ```bash
   npx hardhat node
   ```
   Keep this running throughout your testing.

## Deployment Steps

### Step 1: Deploy NFT Infrastructure

Deploy all NFT contracts to your local network:

```bash
1
```

This script will:
- Check for existing core contracts (RoleManager, BOGOToken)
- Deploy NFTRegistry
- Deploy BOGOWITickets
- Configure all roles and permissions
- Register BOGOWITickets with NFTRegistry
- Save deployment info to `deployment-nft-localhost.json`

**Expected output:**
```
🎨 Starting NFT Infrastructure Local Deployment...
✅ RoleManager deployed to: 0x...
✅ NFTRegistry deployed to: 0x...
✅ BOGOWITickets deployed to: 0x...
✅ Roles granted
✅ BOGOWITickets registered in NFTRegistry
```

### Step 2: Verify Deployment

Run the verification script to ensure everything is configured correctly:

```bash
npx hardhat run scripts/verify-nft-local.js --network localhost
```

This checks:
- Contract deployments
- RoleManager registrations
- Role assignments
- NFTRegistry setup
- BOGOWITickets configuration

**All checks should show ✅**

### Step 3: Test Minting Operations

Test the NFT minting functionality:

```bash
npx hardhat run scripts/test-mint-local.js --network localhost
```

This tests:
- Single ticket minting
- Batch minting (gas optimization)
- Transfer restrictions
- Registry queries
- Redemption signatures

## NPM Scripts

Add these to your `package.json` for convenience:

```json
{
  "scripts": {
    "node": "hardhat node",
    "deploy-nft-local": "hardhat run scripts/deploy-nft-local.js --network localhost",
    "verify-nft-local": "hardhat run scripts/verify-nft-local.js --network localhost",
    "test-mint-local": "hardhat run scripts/test-mint-local.js --network localhost",
    "clean-local": "rm -f scripts/deployment-nft-localhost.json"
  }
}
```

Then you can run:
```bash
npm run deploy-nft-local
npm run verify-nft-local
npm run test-mint-local
```

## Contract Architecture

```
RoleManager (Access Control)
    ├── NFTRegistry (Central Registry)
    │   └── Manages all NFT contracts
    └── BOGOWITickets (ERC-721)
        └── Event tickets with utility flags
```

### Key Features

**NFTRegistry:**
- Central registry for all NFT contracts
- Contract type categorization (TICKET, COLLECTIBLE, BADGE)
- Version tracking
- Active/inactive status management
- Pagination support for queries

**BOGOWITickets:**
- ERC-721 compliant event tickets
- Time-locked transfers
- Expiry dates
- Redemption with backend signatures
- Utility flags (burn on redeem, non-transferable after redeem)
- Royalty support (ERC-2981)
- Batch minting for gas optimization

## Testing Checklist

- [ ] Deploy contracts locally
- [ ] Verify all contracts deployed
- [ ] Verify role assignments
- [ ] Test single ticket minting
- [ ] Test batch minting
- [ ] Test transfer restrictions
- [ ] Test registry queries
- [ ] Test redemption signatures
- [ ] Check gas usage for batch operations

## Troubleshooting

### Issue: "No deployment found"
**Solution:** Make sure you've run the deployment script first:
```bash
npx hardhat run scripts/deploy-nft-local.js --network localhost
```

### Issue: "Network connection error"
**Solution:** Ensure Hardhat node is running:
```bash
npx hardhat node
```

### Issue: "UnauthorizedRole" errors
**Solution:** Check role assignments with verify script:
```bash
npx hardhat run scripts/verify-nft-local.js --network localhost
```

### Issue: "Contract already registered"
**Solution:** Clean up and redeploy:
```bash
rm scripts/deployment-nft-localhost.json
npx hardhat run scripts/deploy-nft-local.js --network localhost
```

## Next Steps for Mainnet

After successful local testing:

1. **Update configuration** for mainnet:
   - Set proper admin addresses
   - Configure conservation DAO address
   - Set Datakyte API keys

2. **Audit considerations**:
   - Review role assignments
   - Verify time locks and expiry logic
   - Test signature verification thoroughly

3. **Deploy to testnet** first:
   ```bash
   npx hardhat run scripts/deploy-nft-testnet.js --network columbus
   ```

4. **Mainnet deployment**:
   ```bash
   npx hardhat run scripts/deploy-nft-mainnet.js --network camino
   ```

## Security Considerations

1. **Role Management**:
   - Admin roles should be multi-sig wallets in production
   - Backend role should be a secure server wallet
   - Minter role should be limited to authorized services

2. **Signature Security**:
   - Use unique nonces for redemptions
   - Implement deadline checks
   - Store used nonces to prevent replay

3. **Contract Upgrades**:
   - Contracts are not upgradeable by design
   - Plan for migration strategy if updates needed

## Gas Optimization

- Batch minting saves ~40% gas compared to individual mints
- Registry queries use pagination to avoid gas limits
- Transfer restrictions checked in-contract to save failed transaction gas

## Support

For issues or questions:
- Check test coverage: `npm run coverage`
- Run unit tests: `npm test`
- Review contract docs in `/contracts/nft/`