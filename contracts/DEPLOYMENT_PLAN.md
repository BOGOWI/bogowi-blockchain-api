# BOGOWI Mainnet Deployment Plan

## Overview
Phased deployment strategy for BOGOWI token ecosystem on Camino mainnet.

## Token Economics
- **Total Supply**: 1,000,000,000 BOGO (1 Billion)
- **Initial Rewards Allocation**: 50,000,000 BOGO (5%)
- **Unallocated Reserve**: 950,000,000 BOGO (95%)

## Phase 1: Core Infrastructure (Day 1)

### Contracts to Deploy
1. **ContractRegistry_v1_0_0**
   - Deploy first to track all contracts
   - No dependencies

2. **MultisigTreasury_v1_0_0**
   - Configure with initial signers
   - Set appropriate threshold (e.g., 3/5)
   - No token dependency yet

3. **BOGOToken_v1_0_0**
   - Simplified version without:
     - Complex allocation buckets
   - Features to include:
     - Basic ERC20 functionality
     - Minting controls
     - Pause mechanism
     - Zero address validation

### Deployment Steps
```bash
1. Deploy ContractRegistry
2. Deploy MultisigTreasury with signers
3. Deploy BOGOToken
4. Grant MINTER_ROLE to MultisigTreasury
5. Mint initial supply through MultisigTreasury
   - 50M to rewards allocation address
   - 950M to treasury for future allocation
6. Register all contracts in ContractRegistry
7. Verify all contracts on explorer
```

## Phase 2: Reward System (Week 1-2)

### Contracts to Deploy
1. **RewardDistributor_v1_0_0**
   - Configure reward templates
   - Set authorized backends
   - Fund with initial rewards allocation

### Configuration
- Transfer 50M BOGO from treasury to RewardDistributor
- Configure daily limits
- Set up reward templates
- Authorize backend systems

## Phase 3: NFT Ecosystem (Month 1-2)

### Contracts to Deploy
1. **CommercialNFT_v1_0_0**
2. **ConservationNFT_v1_0_0**
3. **RoleManager_v1_0_0** (if needed)

## Pre-Deployment Checklist

### Security
- [ ] Security audit completed
- [ ] Test coverage > 95%
- [ ] Slither analysis passed
- [ ] Manual code review done

### Testing
- [ ] Deployed and tested on Camino testnet (Columbus)
- [ ] Integration tests passed
- [ ] Gas optimization verified
- [ ] Migration scenarios tested

### Documentation
- [ ] Deployment runbook created
- [ ] Emergency procedures documented
- [ ] Public documentation updated
- [ ] Team trained on multisig operations

### Infrastructure
- [ ] Multisig signers confirmed and tested
- [ ] Monitoring systems ready
- [ ] Backup and recovery procedures
- [ ] Contract verification scripts ready

## Contract Simplifications for v1.0.0

### BOGOToken_v1_0_0 Changes
Remove:
- Complex allocation buckets (DAO/Business)
- Timelock mechanisms (can add later)

Simplify to:
```solidity
contract BOGOToken_v1_0_0 {
    uint256 public constant MAX_SUPPLY = 1_000_000_000 * 10**18;
    uint256 public constant INITIAL_REWARDS = 50_000_000 * 10**18;
    
    mapping(address => bool) public minters;
    uint256 public totalMinted;
    
    // Basic minting with supply cap
    function mint(address to, uint256 amount) external onlyMinter {
        require(totalMinted + amount <= MAX_SUPPLY, "Exceeds max supply");
        totalMinted += amount;
        _mint(to, amount);
    }
}
```

### MultisigTreasury_v1_0_0 Focus
- Keep core multisig functionality
- Remove complex features for v1:
  - Function call restrictions (add later)
  - Emergency mechanisms (add with governance)
  
## Risk Mitigation

1. **Start Small**: Deploy minimal contracts first
2. **Test Thoroughly**: Use testnet for full simulation
3. **Monitor Closely**: Watch first 24-48 hours carefully
4. **Have Rollback Plan**: Document recovery procedures
5. **Gradual Rollout**: Don't mint full supply immediately

## Post-Deployment

### Immediate Actions
1. Verify all contracts on explorer
2. Update documentation with addresses
3. Monitor initial transactions
4. Test all critical functions

### First Week
1. Deploy reward system
2. Begin small reward distributions
3. Monitor gas usage and optimize
4. Gather community feedback

### First Month
1. Deploy NFT contracts
2. Integrate with frontend
3. Plan next phase features
4. Conduct security review

## Emergency Contacts

- Technical Lead: [REDACTED]
- Security Team: security@bogowi.com
- Multisig Signers: [Configure in deployment]

## Version Control

All deployment transactions and addresses must be documented in:
- `deployments/mainnet/v1_0_0.json`
- ContractRegistry on-chain
- Public documentation