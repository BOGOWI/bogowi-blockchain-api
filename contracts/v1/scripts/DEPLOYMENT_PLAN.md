# V1 Contracts Deployment Plan

## Overview
This document outlines the deployment process for BOGOWI v1 contracts on both Columbus testnet and Camino mainnet.

## Contracts to Deploy
1. **RoleManager.sol** - Centralized role management
2. **BOGOToken.sol** - Main BOGO token (ERC20)
3. **BOGORewardDistributor.sol** - Reward distribution system

## Deployment Order (CRITICAL)
1. Deploy RoleManager first (no dependencies)
2. Deploy BOGOToken with RoleManager address
3. Deploy BOGORewardDistributor with RoleManager and BOGOToken addresses

## Pre-Deployment Requirements

### Critical Preflight Checklist
**⚠️ MANDATORY**: Complete the `DEPLOYMENT_PREFLIGHT_CHECKLIST.md` before any deployment!
This aerospace-grade checklist ensures zero-defect deployment.

### Environment Variables Needed
- `RPC_URL` - Network RPC endpoint
- `PRIVATE_KEY` - Deployer wallet private key  
- `ADMIN_ADDRESS` - Admin address for role management
- `MULTISIG_SIGNER_1`, `MULTISIG_SIGNER_2`, `MULTISIG_SIGNER_3` - For future multisig

### Addresses We Have (from .env)
- Deployer wallet (from PRIVATE_KEY)
- Admin wallet (needs to be set)
- Multisig signers (need to be set)

### Addresses To Be Created
- `ROLE_MANAGER_ADDRESS` - Will be created during deployment
- `BOGO_TOKEN_ADDRESS` - Will be created during deployment
- `REWARD_DISTRIBUTOR_ADDRESS` - Will be created during deployment

## Deployment Steps

### For Testnet:
```bash
# 1. Deploy contracts
npm run deploy:testnet

# 2. CRITICAL: Register contracts with RoleManager
npm run register:testnet

# 3. Setup roles
npm run setup-roles:testnet

# 4. Mint initial supply (10M from rewards allocation)
npm run mint-supply:testnet

# 5. Fund distributor with 10M BOGO
npm run fund-distributor:testnet

# 6. Verify deployment
npm run verify-deployment:testnet
```

### For Mainnet:
```bash
# Same steps but with :mainnet suffix
npm run deploy:mainnet
npm run register:mainnet
npm run setup-roles:mainnet
npm run mint-supply:mainnet
npm run fund-distributor:mainnet
npm run verify-deployment:mainnet
```

## Critical Notes

### Contract Registration (REQUIRED!)
- Contracts MUST be registered with RoleManager after deployment
- Without registration, all role checks will fail with "Not a registered contract"
- This is done by the `register-contracts.js` script

### Token Allocation
- Total BOGO supply: 1 billion
- Rewards allocation: 50M (5%)
- Initial distributor funding: 10M from rewards
- Remaining rewards: 40M (can be minted on-demand by distributor)

### Role Requirements
- Deployer needs BUSINESS_ROLE to mint from rewards
- RewardDistributor gets BUSINESS_ROLE to mint more rewards as needed
- Admin gets PAUSER_ROLE for emergency stops

## Scripts Overview
1. `deploy-all.js` - Deploys all three contracts
2. `register-contracts.js` - Registers contracts with RoleManager (CRITICAL!)
3. `setup-roles.js` - Grants necessary roles
4. `mint-initial-supply.js` - Mints 10M BOGO from rewards allocation
5. `fund-distributor.js` - Transfers 10M to distributor
6. `verify-deployment.js` - Verifies everything is configured correctly

## Important Notes
- Keep deployment transaction hashes
- Save all contract addresses immediately (auto-saved to deployment-{network}.json)
- Test on Columbus before mainnet
- Use same deployment order on both networks
- Update .env with new addresses after deployment