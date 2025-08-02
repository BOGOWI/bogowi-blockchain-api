# Registry Pattern Implementation Guide

## Overview
This guide explains how to implement the registry pattern for BOGOWI's non-upgradeable architecture.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     ContractRegistry                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Mappings:                                          │   │
│  │  - "BOGOToken" → 0x123... (immutable)              │   │
│  │  - "RewardDistributor" → 0x456... (v2)             │   │
│  │  - "CommercialNFT" → 0x789... (v1)                 │   │
│  │  - "MigrationHelper" → 0xABC...                    │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                               ↑
                               │ Lookup
                               │
┌──────────────────┐   ┌──────────────────┐   ┌──────────────┐
│   Frontend/DApp  │   │  Other Contracts │   │   Scripts    │
└──────────────────┘   └──────────────────┘   └──────────────┘
```

## Implementation Steps

### Step 1: Deploy Core Infrastructure

```javascript
// scripts/deploy-core-infrastructure.js
const { ethers } = require("hardhat");

async function main() {
    console.log("Deploying core infrastructure...");
    
    // 1. Deploy ContractRegistry
    const ContractRegistry = await ethers.getContractFactory("ContractRegistry");
    const registry = await ContractRegistry.deploy();
    await registry.deployed();
    console.log("ContractRegistry deployed to:", registry.address);
    
    // 2. Deploy MigrationHelper
    const MigrationHelper = await ethers.getContractFactory("MigrationHelper");
    const migrationHelper = await MigrationHelper.deploy();
    await migrationHelper.deployed();
    console.log("MigrationHelper deployed to:", migrationHelper.address);
    
    // 3. Register MigrationHelper in registry
    await registry.registerContract("MigrationHelper", migrationHelper.address);
    
    return { registry, migrationHelper };
}
```

### Step 2: Deploy Immutable Core Contracts

```javascript
// scripts/deploy-core-contracts.js
async function deployCore(registryAddress) {
    // Deploy immutable contracts
    const BOGOToken = await ethers.getContractFactory("BOGOTokenV2");
    const token = await BOGOToken.deploy();
    await token.deployed();
    
    const MultisigTreasury = await ethers.getContractFactory("MultisigTreasury");
    const treasury = await MultisigTreasury.deploy(signers, threshold);
    await treasury.deployed();
    
    const RoleManager = await ethers.getContractFactory("RoleManager");
    const roleManager = await RoleManager.deploy();
    await roleManager.deployed();
    
    // Register in registry (these won't change)
    const registry = await ethers.getContractAt("ContractRegistry", registryAddress);
    await registry.registerContract("BOGOToken", token.address);
    await registry.registerContract("MultisigTreasury", treasury.address);
    await registry.registerContract("RoleManager", roleManager.address);
    
    return { token, treasury, roleManager };
}
```

### Step 3: Deploy Replaceable Contracts

```javascript
// scripts/deploy-peripherals.js
async function deployPeripherals(registryAddress, tokenAddress, treasuryAddress) {
    const registry = await ethers.getContractAt("ContractRegistry", registryAddress);
    
    // Deploy replaceable contracts
    const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor");
    const distributor = await RewardDistributor.deploy(tokenAddress, treasuryAddress);
    await distributor.deployed();
    
    // Register in registry
    await registry.registerContract("RewardDistributor", distributor.address);
    
    return { distributor };
}
```

## Using the Registry in Contracts

### Example: Contract that Uses Registry

```solidity
// contracts/ExampleConsumer.sol
pragma solidity ^0.8.20;

import "./utils/ContractRegistry.sol";
import "./interfaces/IRewardDistributor.sol";

contract ExampleConsumer {
    ContractRegistry public immutable registry;
    
    constructor(address _registry) {
        registry = ContractRegistry(_registry);
    }
    
    function claimReward() external {
        // Get current distributor from registry
        address distributorAddr = registry.getContract("RewardDistributor");
        IRewardDistributor distributor = IRewardDistributor(distributorAddr);
        
        // Use the distributor
        distributor.claimReward("welcome_bonus");
    }
}
```

## Frontend Integration

### JavaScript/TypeScript Integration

```typescript
// frontend/src/contracts/registry.ts
import { ethers } from 'ethers';
import ContractRegistryABI from '../abi/ContractRegistry.json';
import BOGOTokenABI from '../abi/BOGOTokenV2.json';
import RewardDistributorABI from '../abi/BOGORewardDistributor.json';

class ContractManager {
    private registry: ethers.Contract;
    private provider: ethers.Provider;
    private cache: Map<string, string> = new Map();
    
    constructor(registryAddress: string, provider: ethers.Provider) {
        this.provider = provider;
        this.registry = new ethers.Contract(
            registryAddress,
            ContractRegistryABI,
            provider
        );
    }
    
    async getContract(name: string): Promise<string> {
        // Check cache first
        if (this.cache.has(name)) {
            return this.cache.get(name)!;
        }
        
        // Fetch from registry
        const address = await this.registry.getContract(name);
        this.cache.set(name, address);
        return address;
    }
    
    async getToken(): Promise<ethers.Contract> {
        const address = await this.getContract("BOGOToken");
        return new ethers.Contract(address, BOGOTokenABI, this.provider);
    }
    
    async getRewardDistributor(): Promise<ethers.Contract> {
        const address = await this.getContract("RewardDistributor");
        return new ethers.Contract(address, RewardDistributorABI, this.provider);
    }
    
    // Clear cache when contracts might have been updated
    clearCache() {
        this.cache.clear();
    }
}

// Usage
const contractManager = new ContractManager(REGISTRY_ADDRESS, provider);
const token = await contractManager.getToken();
const balance = await token.balanceOf(userAddress);
```

## Migration Process

### When to Migrate

1. **Critical Bug Found**: Deploy fix immediately
2. **New Features**: Deploy enhanced version
3. **Gas Optimizations**: Deploy optimized version
4. **Regulatory Changes**: Deploy compliant version

### Migration Workflow

```javascript
// scripts/migrate-reward-distributor.js
async function migrateRewardDistributor() {
    const [deployer] = await ethers.getSigners();
    
    // 1. Get registry
    const registry = await ethers.getContractAt(
        "ContractRegistry", 
        REGISTRY_ADDRESS
    );
    
    // 2. Get old distributor
    const oldDistributorAddr = await registry.getContract("RewardDistributor");
    const oldDistributor = await ethers.getContractAt(
        "BOGORewardDistributor",
        oldDistributorAddr
    );
    
    // 3. Deploy new distributor
    const NewDistributor = await ethers.getContractFactory("BOGORewardDistributorV2");
    const newDistributor = await NewDistributor.deploy(
        TOKEN_ADDRESS,
        TREASURY_ADDRESS
    );
    await newDistributor.deployed();
    
    // 4. Pause old distributor
    await oldDistributor.pause();
    
    // 5. Transfer any remaining tokens
    const balance = await token.balanceOf(oldDistributorAddr);
    if (balance.gt(0)) {
        await oldDistributor.treasurySweep(
            token.address,
            newDistributor.address,
            balance
        );
    }
    
    // 6. Update registry
    await registry.updateContract("RewardDistributor", newDistributor.address);
    
    // 7. Initialize new distributor
    await newDistributor.setAuthorizedBackend(BACKEND_ADDRESS, true);
    
    console.log("Migration complete!");
    console.log("Old distributor:", oldDistributorAddr);
    console.log("New distributor:", newDistributor.address);
}
```

## Best Practices

### 1. Contract Naming Convention
```solidity
// Use consistent naming in registry
"BOGOToken"           // Core token (never changes)
"RewardDistributor"   // Current reward distributor
"CommercialNFT"       // Current NFT contract
"MigrationHelper"     // Migration utilities
```

### 2. Version Management
```solidity
contract RewardDistributorV2 {
    uint256 public constant VERSION = 2;
    
    // Include version in events
    event Initialized(uint256 version);
}
```

### 3. Interface Stability
```solidity
// Define stable interfaces
interface IRewardDistributor {
    function claimReward(string memory templateId) external;
    function pause() external;
    function unpause() external;
}

// All versions must implement the interface
contract RewardDistributorV1 is IRewardDistributor { }
contract RewardDistributorV2 is IRewardDistributor { }
```

### 4. Graceful Deprecation
```solidity
contract RewardDistributorV1 {
    bool public deprecated;
    address public newVersion;
    
    modifier notDeprecated() {
        require(!deprecated, "Contract deprecated, use new version");
        _;
    }
    
    function deprecate(address _newVersion) external onlyOwner {
        deprecated = true;
        newVersion = _newVersion;
        emit Deprecated(_newVersion);
    }
}
```

## Testing Registry Updates

```javascript
// test/registry-updates.test.js
describe("Registry Updates", function () {
    it("Should update contract address", async function () {
        // Deploy v1
        const v1 = await DistributorV1.deploy();
        await registry.registerContract("RewardDistributor", v1.address);
        
        // Deploy v2
        const v2 = await DistributorV2.deploy();
        await registry.updateContract("RewardDistributor", v2.address);
        
        // Verify update
        expect(await registry.getContract("RewardDistributor"))
            .to.equal(v2.address);
        expect(await registry.getContractVersion("RewardDistributor"))
            .to.equal(2);
    });
});
```

## Monitoring and Alerts

### Key Events to Monitor
```javascript
// Set up event listeners
registry.on("ContractUpdated", (name, oldAddr, newAddr, version) => {
    console.log(`Contract ${name} updated from ${oldAddr} to ${newAddr}`);
    // Send alerts
    // Update frontend cache
    // Log for audit
});
```

### Health Checks
```javascript
async function healthCheck() {
    const contracts = [
        "BOGOToken",
        "RewardDistributor",
        "MultisigTreasury"
    ];
    
    for (const name of contracts) {
        try {
            const addr = await registry.getContract(name);
            const code = await provider.getCode(addr);
            if (code === "0x") {
                console.error(`${name} has no code at ${addr}`);
            }
        } catch (e) {
            console.error(`${name} not found in registry`);
        }
    }
}
```

## Security Considerations

1. **Registry Access Control**: Only multisig should update
2. **Contract Verification**: Always verify new contracts
3. **Migration Testing**: Test on testnet first
4. **Rollback Plan**: Keep old contract addresses
5. **Communication**: Notify users of changes

## Conclusion

The registry pattern provides a clean separation between immutable core contracts and replaceable peripherals. This maintains the security benefits of non-upgradeable contracts while allowing necessary flexibility for bug fixes and feature additions.