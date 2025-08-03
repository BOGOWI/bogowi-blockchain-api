# Event Emission Patches for BOGOWI Contracts

## Overview
This document provides the specific patches needed to add missing event emissions to existing contracts.

## 1. MultisigTreasury.sol Patches

### Add Event Declaration
```solidity
// Add after line 68 with other events
event FunctionRestrictionsToggled(bool enabled);
```

### Update toggleFunctionRestrictions Function
```solidity
// Replace function at line 487
function toggleFunctionRestrictions() external onlyMultisig {
    restrictFunctionCalls = !restrictFunctionCalls;
    emit FunctionRestrictionsToggled(restrictFunctionCalls);
}
```

## 2. BOGORewardDistributor.sol Patches

### Add Event Declarations
```solidity
// Add after line 46 with other events
event AuthorizedBackendSet(address indexed backend, bool authorized);
event DailyLimitReset(uint256 timestamp, uint256 previousDistributed);
```

### Update setAuthorizedBackend Function
```solidity
// Replace function at line 273
function setAuthorizedBackend(address backend, bool authorized) external onlyTreasury {
    authorizedBackends[backend] = authorized;
    emit AuthorizedBackendSet(backend, authorized);
}
```

### Update _resetDailyLimit Function
```solidity
// Replace function at line 169
function _resetDailyLimit() private {
    if (block.timestamp >= lastResetTime + 1 days) {
        uint256 previousDistributed = dailyDistributed;
        dailyDistributed = 0;
        lastResetTime = block.timestamp;
        emit DailyLimitReset(block.timestamp, previousDistributed);
    }
}
```

## 3. Enhanced Event Tracking (Optional)

### For Critical Operations
```solidity
// Add to pause/unpause functions if you want more detail
function pause() external onlyTreasury {
    _pause();
    emit ContractPaused(msg.sender);
}

function unpause() external onlyTreasury {
    _unpause();
    emit ContractUnpaused(msg.sender);
}
```

## 4. Complete List of Functions That Should Emit Events

### Always Emit Events For:
1. **State Changes**: Any function that modifies contract state
2. **Financial Operations**: Transfers, mints, burns, fees
3. **Access Control**: Role grants, revokes, ownership transfers
4. **Configuration Changes**: Parameter updates, threshold changes
5. **Emergency Actions**: Pauses, emergency withdrawals
6. **User Actions**: Claims, stakes, votes

### Current Status by Contract:

#### BOGOTokenV2
- ✅ Minting functions (emit AllocationMinted)
- ✅ Timelock operations (emit TimelockQueued/Executed/Cancelled)
- ✅ Pause/Unpause (inherited from Pausable)
- ✅ Burns (ERC20 Transfer event to 0x0)

#### MultisigTreasury
- ✅ Transaction lifecycle (Submit/Confirm/Execute/Cancel)
- ✅ Signer management (Add/Remove)
- ✅ Threshold changes
- ✅ Emergency withdrawals
- ❌ **Missing**: Function restrictions toggle

#### BOGORewardDistributor
- ✅ Reward claims
- ✅ Referral claims
- ✅ Template updates
- ✅ Whitelist updates
- ❌ **Missing**: Backend authorization
- ❌ **Missing**: Daily limit resets

#### CommercialNFT/ConservationNFT
- ✅ Minting operations
- ✅ Metadata updates
- ✅ Role changes
- ✅ Pause/Unpause

## 5. Implementation Priority

### High Priority (Security/Financial)
1. Authorization changes (setAuthorizedBackend)
2. Configuration toggles (toggleFunctionRestrictions)

### Medium Priority (Monitoring)
1. Daily limit resets
2. Enhanced pause/unpause tracking

### Low Priority (Nice to Have)
1. Detailed parameter change tracking
2. User action analytics

## 6. Gas Considerations

Events are relatively cheap but do add gas cost:
- Event with no indexed parameters: ~375 gas
- Each indexed parameter: +375 gas
- Each non-indexed parameter: ~8 gas per byte

Balance between transparency and gas efficiency.

## 7. Event Monitoring Best Practices

### Frontend Integration
```javascript
// Listen for authorization changes
rewardDistributor.on("AuthorizedBackendSet", (backend, authorized) => {
    console.log(`Backend ${backend} authorization: ${authorized}`);
});

// Listen for daily resets
rewardDistributor.on("DailyLimitReset", (timestamp, previousAmount) => {
    console.log(`Daily limit reset at ${timestamp}, distributed: ${previousAmount}`);
});
```

### Indexing Strategy
- Index addresses for filtering by user/contract
- Index roles/types for filtering by category
- Don't index amounts (use logs for analytics)

## 8. Testing Event Emissions

```javascript
// Example test
it("Should emit FunctionRestrictionsToggled event", async function () {
    await expect(multisig.toggleFunctionRestrictions())
        .to.emit(multisig, "FunctionRestrictionsToggled")
        .withArgs(true);
});

it("Should emit AuthorizedBackendSet event", async function () {
    await expect(distributor.setAuthorizedBackend(backend.address, true))
        .to.emit(distributor, "AuthorizedBackendSet")
        .withArgs(backend.address, true);
});
```

## Summary

The main missing events are:
1. `FunctionRestrictionsToggled` in MultisigTreasury
2. `AuthorizedBackendSet` in BOGORewardDistributor
3. `DailyLimitReset` in BOGORewardDistributor

These are all simple additions that improve transparency without breaking changes.