# BOGOWI Contract Versioning Strategy

## Version Format
All contracts follow semantic versioning: `vMAJOR.MINOR.PATCH`

### Version Increments
- **MAJOR**: Breaking changes, incompatible upgrades
- **MINOR**: New features, backwards compatible
- **PATCH**: Bug fixes, optimizations

## Naming Convention
```
ContractName_vX_Y_Z.sol
```

Examples:
- `BOGOToken_v1_0_0.sol`
- `MultisigTreasury_v1_0_0.sol`
- `RewardDistributor_v1_0_0.sol`

## Directory Structure
```
contracts/
├── core/                    # Essential contracts
├── nft/                     # NFT-related contracts
├── governance/              # Governance and access control
├── utils/                   # Utility contracts
├── interfaces/              # Contract interfaces
├── deprecated/              # Old versions (not deployed)
└── mocks/                   # Test contracts
```

## Version Tracking

### In-Contract Version
```solidity
string public constant VERSION = "1.0.0";
uint256 public constant VERSION_NUMBER = 1000000; // Format: MAJOR*1e6 + MINOR*1e3 + PATCH
```

### Deployment Registry
Track all deployed contracts in ContractRegistry with:
- Contract name
- Version
- Address
- Deployment block
- Active status

## Migration Rules

1. **Never modify deployed contracts** - Always deploy new versions
2. **Maintain upgrade paths** - Document migration procedures
3. **Deprecate gracefully** - Mark old versions, don't delete
4. **Test migrations** - Always test on testnet first

## Initial Mainnet Versions

### Phase 1 - Core Contracts (v1.0.0)
- `BOGOToken_v1_0_0.sol` - Simplified token with 5% rewards allocation
- `MultisigTreasury_v1_0_0.sol` - Secure fund management
- `ContractRegistry_v1_0_0.sol` - Track deployments

### Phase 2 - Rewards (v1.0.0)
- `RewardDistributor_v1_0_0.sol` - Reward distribution system

### Phase 3 - NFTs (v1.0.0)
- `CommercialNFT_v1_0_0.sol`
- `ConservationNFT_v1_0_0.sol`

## Version History Format

Each contract should maintain a version history comment:

```solidity
/**
 * @title BOGOToken
 * @custom:version 1.0.0
 * @custom:version-history
 * - 1.0.0: Initial mainnet deployment
 *   - Simplified allocations (5% rewards, 95% unallocated)
 *   - Removed flavored tokens
 *   - Basic minting and burning
 */
```

## Upgrade Procedures

1. Deploy new version contract
2. Update ContractRegistry
3. Migrate state if needed
4. Update documentation
5. Deprecate old version

## Testing Requirements

Before mainnet deployment:
1. Full test coverage on new version
2. Migration testing on testnet
3. Security audit for major versions
4. Gas optimization benchmarks