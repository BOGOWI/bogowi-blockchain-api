# BOGOWI Non-Upgradeable Deployment Sequence

## Overview
This document provides the complete deployment sequence for BOGOWI's non-upgradeable architecture with registry pattern support.

## Pre-Deployment Checklist

- [ ] All contracts compiled successfully
- [ ] All tests passing
- [ ] Security audit completed
- [ ] Deployment addresses documented
- [ ] Gas funds available in deployer wallet
- [ ] Multisig signers confirmed
- [ ] Environment variables configured
- [ ] Deployment scripts tested on testnet

## Deployment Order

The deployment must follow this specific order due to dependencies:

```
1. ContractRegistry
2. RoleManager
3. MultisigTreasury
4. BOGOTokenV2
5. MigrationHelper
6. BOGORewardDistributor
7. NFT Contracts (Commercial/Conservation)
8. Register all contracts
9. Configure roles and permissions
10. Verify all contracts
```

## Detailed Deployment Steps

### Phase 1: Core Infrastructure

#### 1.1 Deploy ContractRegistry
```bash
npx hardhat run scripts/01-deploy-registry.js --network mainnet
```

**Script: `scripts/01-deploy-registry.js`**
```javascript
const { ethers } = require("hardhat");
const fs = require("fs");

async function main() {
    console.log("Deploying ContractRegistry...");
    
    const ContractRegistry = await ethers.getContractFactory("ContractRegistry");
    const registry = await ContractRegistry.deploy();
    await registry.deployed();
    
    console.log("ContractRegistry deployed to:", registry.address);
    
    // Save address
    const deployment = {
        ContractRegistry: registry.address,
        deployer: (await ethers.getSigners())[0].address,
        timestamp: new Date().toISOString(),
        network: network.name
    };
    
    fs.writeFileSync(
        `deployments/1-registry-${network.name}.json`,
        JSON.stringify(deployment, null, 2)
    );
    
    return registry.address;
}

main().catch((error) => {
    console.error(error);
    process.exit(1);
});
```

#### 1.2 Deploy RoleManager
```bash
npx hardhat run scripts/02-deploy-role-manager.js --network mainnet
```

**Deployment includes:**
- Setting up all role definitions
- Configuring role hierarchies
- Initial admin assignment

### Phase 2: Core Contracts (Immutable)

#### 2.1 Deploy MultisigTreasury
```bash
npx hardhat run scripts/03-deploy-multisig.js --network mainnet
```

**Configuration Required:**
- Signer addresses (from secure storage)
- Threshold (e.g., 3 of 5)
- Execution delay settings

#### 2.2 Deploy BOGOTokenV2
```bash
npx hardhat run scripts/04-deploy-token.js --network mainnet
```

**Initial Setup:**
- Name: "BOGOWI"
- Symbol: "BOGO"
- Initial roles to deployer (temporary)

### Phase 3: Support Contracts

#### 3.1 Deploy MigrationHelper
```bash
npx hardhat run scripts/05-deploy-migration-helper.js --network mainnet
```

#### 3.2 Deploy BOGORewardDistributor
```bash
npx hardhat run scripts/06-deploy-reward-distributor.js --network mainnet
```

**Parameters:**
- Token address (from step 2.2)
- Treasury address (from step 2.1)

### Phase 4: Registration and Configuration

#### 4.1 Register All Contracts
```bash
npx hardhat run scripts/07-register-contracts.js --network mainnet
```

**Script: `scripts/07-register-contracts.js`**
```javascript
async function main() {
    const addresses = require("../deployments/addresses.json");
    const registry = await ethers.getContractAt("ContractRegistry", addresses.ContractRegistry);
    
    // Register immutable contracts
    await registry.registerContract("BOGOToken", addresses.BOGOTokenV2);
    await registry.registerContract("MultisigTreasury", addresses.MultisigTreasury);
    await registry.registerContract("RoleManager", addresses.RoleManager);
    
    // Register replaceable contracts
    await registry.registerContract("RewardDistributor", addresses.BOGORewardDistributor);
    await registry.registerContract("MigrationHelper", addresses.MigrationHelper);
    
    console.log("All contracts registered");
}
```

#### 4.2 Configure Roles
```bash
npx hardhat run scripts/08-configure-roles.js --network mainnet
```

**Role Assignments:**
```javascript
// Transfer roles from deployer to proper addresses
await roleManager.grantRole(DAO_ROLE, multisigTreasury.address);
await roleManager.grantRole(BUSINESS_ROLE, businessMultisig.address);
await roleManager.grantRole(TREASURY_ROLE, multisigTreasury.address);
await roleManager.grantRole(PAUSER_ROLE, emergencyPauser.address);

// Revoke temporary roles from deployer
await roleManager.revokeRole(MINTER_ROLE, deployer.address);
await roleManager.renounceRole(DEFAULT_ADMIN_ROLE, deployer.address);
```

### Phase 5: Verification and Handover

#### 5.1 Verify All Contracts
```bash
npx hardhat run scripts/09-verify-contracts.js --network mainnet
```

**Verification includes:**
- Source code verification on Etherscan
- Contract state validation
- Role assignments check
- Registry entries confirmation

#### 5.2 Final Configuration
```bash
npx hardhat run scripts/10-final-configuration.js --network mainnet
```

**Final Steps:**
- Fund reward distributor
- Set authorized backends
- Configure whitelist (if any)
- Enable emergency procedures

## Post-Deployment Verification

### Automated Checks Script
```javascript
// scripts/verify-deployment.js
async function verifyDeployment() {
    console.log("Running post-deployment verification...");
    
    const checks = [
        checkRegistryEntries(),
        checkRoleAssignments(),
        checkContractConnections(),
        checkTokenSupply(),
        checkMultisigSetup(),
        checkPausability(),
        checkEmergencyProcedures()
    ];
    
    const results = await Promise.all(checks);
    const failed = results.filter(r => !r.success);
    
    if (failed.length > 0) {
        console.error("❌ Verification failed:");
        failed.forEach(f => console.error(f.message));
        process.exit(1);
    }
    
    console.log("✅ All checks passed!");
}
```

### Manual Verification Checklist

- [ ] Registry accessible and contains all contracts
- [ ] Token minting works with proper roles
- [ ] Multisig can execute transactions
- [ ] Reward distributor can distribute tokens
- [ ] Pause functionality works
- [ ] Emergency procedures tested
- [ ] Frontend can connect to contracts
- [ ] Events being emitted correctly

## Environment Configuration

### Required Environment Variables
```bash
# .env.production
DEPLOYER_PRIVATE_KEY=
ETHERSCAN_API_KEY=
INFURA_PROJECT_ID=

# Multisig Configuration
MULTISIG_SIGNERS=0x...,0x...,0x...
MULTISIG_THRESHOLD=3

# Role Assignments
DAO_ADDRESS=
BUSINESS_ADDRESS=
EMERGENCY_PAUSER=
BACKEND_ADDRESS=

# Network Configuration
MAINNET_RPC_URL=
GAS_PRICE_GWEI=
```

## Gas Optimization Tips

1. **Deploy during low traffic**: Check gas prices
2. **Use CREATE2**: For predictable addresses
3. **Batch operations**: Register multiple contracts in one tx
4. **Optimize constructor parameters**: Minimize storage writes

## Emergency Deployment Recovery

If deployment fails at any step:

1. **Do NOT proceed** to next steps
2. **Document** the failure point
3. **Analyze** gas usage and errors
4. **Resume** from last successful step
5. **Update** deployment records

### Recovery Script Template
```javascript
async function resumeDeployment(fromStep) {
    const deploymentState = require("./deployment-state.json");
    
    switch(fromStep) {
        case 'registry':
            return deployFromRegistry();
        case 'token':
            return deployFromToken(deploymentState.registry);
        // ... other cases
    }
}
```

## Deployment Timeline

Estimated time for complete deployment:

| Phase | Duration | Gas Cost (estimated) |
|-------|----------|---------------------|
| Core Infrastructure | 15 min | ~2M gas |
| Core Contracts | 20 min | ~5M gas |
| Support Contracts | 15 min | ~3M gas |
| Configuration | 30 min | ~1M gas |
| Verification | 20 min | - |
| **Total** | **~2 hours** | **~11M gas** |

## Post-Deployment Actions

### Immediate (within 1 hour)
1. Verify all contracts on Etherscan
2. Update frontend configurations
3. Test basic operations
4. Monitor for anomalies

### Within 24 hours
1. Complete security checklist
2. Enable monitoring alerts
3. Document final addresses
4. Announce deployment

### Within 1 week
1. Conduct post-deployment audit
2. Review gas usage
3. Optimize if needed
4. Plan for next phase

## Rollback Procedure

If critical issues discovered post-deployment:

1. **Pause all contracts** immediately
2. **Assess impact** and affected users
3. **Communicate** with community
4. **Plan migration** to fixed contracts
5. **Execute migration** using MigrationHelper

## Documentation Updates

After successful deployment:

1. Update README with mainnet addresses
2. Update API documentation
3. Create user guides
4. Update developer documentation
5. Archive deployment artifacts

## Conclusion

This deployment sequence ensures a secure, verifiable deployment of BOGOWI's non-upgradeable architecture. The registry pattern provides flexibility for peripheral contracts while maintaining the security and trust of immutable core contracts.

**Remember**: Once deployed, core contracts cannot be changed. Double-check everything before mainnet deployment.