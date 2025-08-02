# RoleManager Migration Guide

## Overview
This guide explains how to migrate from individual contract role management to the centralized RoleManager system.

## Benefits of Centralized Role Management

1. **Single Source of Truth**: All role assignments are managed in one contract
2. **Consistent Permissions**: Ensures role consistency across all contracts
3. **Reduced Gas Costs**: Batch operations for role management
4. **Simplified Administration**: One interface for all role management
5. **Better Security**: Centralized audit trail and access control

## Architecture

### Core Components

1. **RoleManager.sol**: Central contract managing all roles
2. **IRoleManager.sol**: Interface for contracts to interact with RoleManager
3. **RoleManaged.sol**: Base contract providing role checking functionality
4. **Updated Contracts**: Modified versions using RoleManager

### Role Definitions

- `DEFAULT_ADMIN_ROLE`: System administrator (manages other roles)
- `DAO_ROLE`: DAO governance operations
- `BUSINESS_ROLE`: Business operations and treasury management
- `MINTER_ROLE`: Token minting permissions
- `PAUSER_ROLE`: Emergency pause capabilities
- `TREASURY_ROLE`: Treasury and financial operations
- `DISTRIBUTOR_BACKEND_ROLE`: Backend services for reward distribution

## Migration Steps

### 1. Deploy RoleManager
```bash
npx hardhat run scripts/deploy-role-manager.js --network <network>
```

### 2. Update Existing Contract References

Replace existing contracts with their RoleManaged versions:
- `BOGOTokenV2` → `BOGOTokenV2_RoleManaged`
- `BOGORewardDistributor` → `BOGORewardDistributor_RoleManaged`
- `CommercialNFT` → `CommercialNFT_RoleManaged` (to be implemented)
- `ConservationNFT` → `ConservationNFT_RoleManaged` (to be implemented)

### 3. Register Contracts
```javascript
// Register each contract with RoleManager
await roleManager.registerContract(contractAddress, "ContractName");
```

### 4. Migrate Role Assignments
```javascript
// Example: Migrate existing role holders
const existingHolders = await oldContract.getRoleMemberCount(ROLE);
for (let i = 0; i < existingHolders; i++) {
    const holder = await oldContract.getRoleMember(ROLE, i);
    await roleManager.grantRole(ROLE, holder);
}
```

### 5. Update Integration Points

#### For Smart Contracts
```solidity
// Old approach
contract MyContract is AccessControl {
    function someFunction() external onlyRole(MINTER_ROLE) {
        // ...
    }
}

// New approach
contract MyContract is RoleManaged {
    constructor(address _roleManager) RoleManaged(_roleManager) {}
    
    function someFunction() external onlyRole(roleManager.MINTER_ROLE()) {
        // ...
    }
}
```

#### For Backend Services
```javascript
// Check roles via RoleManager
const hasRole = await roleManager.hasRole(ROLE, userAddress);
```

## Security Considerations

1. **Admin Transfer**: Use `transferAdmin()` carefully - it's irreversible
2. **Contract Registration**: Only register trusted contracts
3. **Role Hierarchy**: DEFAULT_ADMIN_ROLE can manage all other roles
4. **Emergency Response**: PAUSER_ROLE can pause RoleManager itself

## Testing Checklist

- [ ] All contracts can check roles via RoleManager
- [ ] Role assignments work correctly
- [ ] Batch operations function properly
- [ ] Emergency pause mechanisms work
- [ ] Role revocation prevents access
- [ ] Gas costs are acceptable

## Rollback Plan

If issues arise:
1. Pause affected contracts
2. Deploy fixed versions
3. Re-register contracts
4. Restore role assignments from backup

## Monitoring

Track these events:
- `ContractRegistered`
- `RoleGrantedGlobally`
- `RoleRevokedGlobally`

## FAQ

**Q: What happens to existing role assignments?**
A: They need to be manually migrated to RoleManager.

**Q: Can contracts have local roles?**
A: No, all roles must be managed through RoleManager.

**Q: How do we add new roles?**
A: Add them to RoleManager and update contracts to use them.

**Q: What if RoleManager is compromised?**
A: Use multisig for DEFAULT_ADMIN_ROLE and implement timelock for critical operations.