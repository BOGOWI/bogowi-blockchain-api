# Zero Address Validation Guide

## Overview
This guide documents all functions that require zero address validation to prevent token burns and lost funds.

## Why Zero Address Validation is Critical

Sending tokens or granting permissions to the zero address (0x0000...0000) results in:
- **Permanent token loss** - Tokens sent to zero address are burned forever
- **Security vulnerabilities** - Roles granted to zero address cannot be revoked
- **Failed transactions** - Operations involving zero address often fail unexpectedly
- **Poor user experience** - Users may accidentally input zero address

## Implementation Pattern

### Using Custom Errors (Gas Efficient)
```solidity
error InvalidAddress();

modifier notZeroAddress(address addr) {
    if (addr == address(0)) revert InvalidAddress();
    _;
}

function transfer(address to, uint256 amount) 
    external 
    notZeroAddress(to) 
{
    // Function logic
}
```

### Traditional Require (Less Gas Efficient)
```solidity
function transfer(address to, uint256 amount) external {
    require(to != address(0), "Invalid address");
    // Function logic
}
```

## Functions Requiring Zero Address Validation

### BOGOTokenV2

#### Token Minting Functions
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `mintFromDAO(address to, uint256 amount)` | `to` | Permanent token loss |
| `mintFromBusiness(address to, uint256 amount)` | `to` | Permanent token loss |
| `mintFromRewards(address to, uint256 amount)` | `to` | Permanent token loss |

#### Token Burning Functions
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `burnFrom(address account, uint256 amount)` | `account` | Transaction failure |

#### Role Management
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `grantRole(bytes32 role, address account)` | `account` | Unrevokable permission |
| `revokeRole(bytes32 role, address account)` | `account` | Transaction failure |

#### Token Registration
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `queueRegisterFlavoredToken(string flavor, address tokenAddress)` | `tokenAddress` | Invalid configuration |

### BOGORewardDistributor

#### Constructor
| Parameter | Risk if Not Validated |
|-----------|----------------------|
| `_bogoToken` | Contract deployment failure |
| `_treasury` | No treasury control |

#### Reward Distribution
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `claimCustomReward(address recipient, uint256 amount, string reason)` | `recipient` | Permanent token loss |
| `claimReferralBonus(address referrer)` | `referrer` | Invalid referral chain |

#### Admin Functions
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `addToWhitelist(address[] wallets)` | Each `wallet` in array | Invalid whitelist entry |
| `removeFromWhitelist(address wallet)` | `wallet` | Transaction failure |
| `setAuthorizedBackend(address backend, bool authorized)` | `backend` | Invalid authorization |
| `treasurySweep(address token, address to, uint256 amount)` | `to` | Permanent fund loss |

### CommercialNFT & ConservationNFT

#### Minting Functions
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `mint(address to, uint256 tokenType, uint256 amount)` | `to` | NFT permanently lost |
| `mintBatch(address to, uint256[] tokenTypes, uint256[] amounts)` | `to` | Multiple NFTs lost |

#### Configuration
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `setTreasuryAddress(address newTreasuryAddress)` | `newTreasuryAddress` | No treasury control |

### MultisigTreasury

#### Transaction Management
| Function | Parameters to Validate | Risk if Not Validated |
|----------|----------------------|----------------------|
| `submitTransaction(address to, uint256 value, bytes data)` | `to` | Failed transaction |
| `addSigner(address signer)` | `signer` | Invalid signer |
| `emergencyWithdrawETH(address payable to, uint256 amount)` | `to` | Permanent ETH loss |

## Amount Validation

In addition to address validation, validate amounts are greater than zero:

```solidity
error InvalidAmount();

modifier notZeroAmount(uint256 amount) {
    if (amount == 0) revert InvalidAmount();
    _;
}
```

Functions requiring amount validation:
- All minting functions
- All burning functions
- All transfer functions
- Treasury sweep operations

## Testing Zero Address Validation

### Essential Test Cases

```javascript
// Test zero address rejection
await expect(
    contract.transfer(ZERO_ADDRESS, amount)
).to.be.revertedWithCustomError(contract, "InvalidAddress");

// Test zero amount rejection
await expect(
    contract.transfer(validAddress, 0)
).to.be.revertedWithCustomError(contract, "InvalidAmount");

// Test valid operation succeeds
await expect(
    contract.transfer(validAddress, validAmount)
).to.not.be.reverted;
```

### Edge Cases to Test
1. Arrays containing zero address
2. Batch operations with mixed valid/invalid addresses
3. Zero address in different parameter positions
4. Combined zero address and zero amount

## Best Practices

1. **Use modifiers** for consistent validation across functions
2. **Validate early** in function execution to save gas
3. **Use custom errors** instead of require strings for gas efficiency
4. **Document validation** in function comments
5. **Test thoroughly** including edge cases

## Migration Checklist

When updating existing contracts:

- [ ] Identify all functions accepting address parameters
- [ ] Add `notZeroAddress` modifier to critical functions
- [ ] Add `notZeroAmount` modifier where applicable
- [ ] Replace require statements with custom errors
- [ ] Update tests to check for custom errors
- [ ] Document changes in contract comments
- [ ] Verify no legitimate use cases for zero address

## Gas Optimization

Custom errors save approximately 200-300 gas per revert:

```solidity
// Gas efficient
if (to == address(0)) revert InvalidAddress();

// Less efficient
require(to != address(0), "Invalid address");
```

## Security Considerations

1. **Constructor validation**: Always validate in constructors to prevent deployment issues
2. **Admin functions**: Critical for role and permission management
3. **Token operations**: Essential for preventing permanent loss
4. **Upgrade safety**: Ensure validation in upgradeable contracts

## Conclusion

Zero address validation is a critical security measure that prevents:
- Permanent token/ETH loss
- Invalid contract states
- Failed transactions
- Poor user experience

Always validate address parameters in functions that:
- Transfer value (tokens, ETH, NFTs)
- Grant permissions or roles
- Set important contract addresses
- Modify user balances or states