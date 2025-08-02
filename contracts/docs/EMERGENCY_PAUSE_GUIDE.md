# Emergency Pause Controller Guide

## Overview

The EmergencyPauseController is a critical security component of the BOGOWI ecosystem that allows authorized guardians to pause all contracts simultaneously in case of a security vulnerability or critical issue.

## Key Features

- **Multi-Signature Guardian System**: Requires multiple guardians to confirm pause/unpause actions
- **Time-Limited Pauses**: Automatic expiry after 72 hours to prevent indefinite lockups
- **Batch Operations**: Pause/unpause multiple contracts in a single transaction
- **Audit Trail**: Complete history of all pause actions for transparency
- **Role-Based Access**: Separate roles for guardians and contract managers

## Architecture

### Roles

1. **GUARDIAN_ROLE**: Can propose and confirm pause/unpause actions
2. **MANAGER_ROLE**: Can add/remove contracts from monitoring
3. **DEFAULT_ADMIN_ROLE**: Can update system parameters and grant roles
4. **PAUSER_ROLE**: Granted to EmergencyPauseController on each pausable contract

### Contracts Integration

All BOGOWI contracts implement the `IPausableContract` interface:
- BOGOTokenV2
- MultisigTreasury
- BOGORewardDistributor
- CommercialNFT
- ConservationNFT

## Deployment

### 1. Deploy EmergencyPauseController

```bash
npx hardhat run scripts/deploy-emergency-pause.js --network <network>
```

### 2. Grant Permissions

```bash
npx hardhat run scripts/update-pause-permissions.js --network <network>
```

### 3. Verify Deployment

Check the deployment file in `deployments/emergency-pause-<network>.json`

## Usage

### Emergency Pause Process

1. **Guardian Detects Issue**
   - Guardian identifies potential vulnerability
   - Initiates pause proposal with reason

2. **Proposal Creation**
   ```javascript
   // Pause specific contracts
   await emergencyPause.proposePause(
     [bogoTokenAddress, treasuryAddress],
     "Potential reentrancy vulnerability in treasury"
   );
   
   // Or pause all contracts
   await emergencyPause.emergencyPauseAll(
     "System-wide security issue detected"
   );
   ```

3. **Guardian Confirmation**
   - Other guardians review the proposal
   - Confirm if they agree with the action
   ```javascript
   await emergencyPause.confirmProposal(proposalId);
   ```

4. **Automatic Execution**
   - Once required confirmations reached (default: 2)
   - Contracts are automatically paused

### Unpause Process

1. **Issue Resolution**
   - Development team fixes the vulnerability
   - Deploys patches if necessary

2. **Unpause Proposal**
   ```javascript
   await emergencyPause.proposeUnpause(
     [bogoTokenAddress, treasuryAddress],
     "Vulnerability patched in commit abc123"
   );
   ```

3. **Confirmation & Execution**
   - Same process as pause proposals

### Automatic Expiry

- Pauses automatically expire after 72 hours
- Anyone can call `checkAndExpirePauses()` to trigger expiry
- Prevents indefinite contract lockup

## Configuration

### Required Confirmations

Update the number of guardian confirmations required:
```javascript
await emergencyPause.updateRequiredConfirmations(3);
```

### Adding New Contracts

Manager adds new contracts to monitoring:
```javascript
await emergencyPause.addContract(
  newContractAddress,
  "NewContractName"
);
```

### Guardian Management

Add new guardian:
```javascript
await emergencyPause.grantRole(GUARDIAN_ROLE, newGuardianAddress);
```

Remove guardian:
```javascript
await emergencyPause.revokeRole(GUARDIAN_ROLE, oldGuardianAddress);
```

## Security Considerations

1. **Guardian Keys**: Store guardian private keys in hardware wallets
2. **Multi-Location**: Distribute guardians across different geographic locations
3. **Response Time**: Establish communication channels for rapid response
4. **Regular Drills**: Practice emergency pause procedures quarterly
5. **Monitoring**: Set up alerts for pause events

## Emergency Response Checklist

- [ ] Identify the security issue
- [ ] Document the vulnerability details
- [ ] Create pause proposal with clear reason
- [ ] Notify other guardians via secure channel
- [ ] Confirm proposal from multiple guardians
- [ ] Verify all affected contracts are paused
- [ ] Begin vulnerability remediation
- [ ] Test fixes thoroughly
- [ ] Create unpause proposal when safe
- [ ] Monitor system after unpause

## Integration Example

For new contracts to be pausable:

```solidity
contract NewContract is Pausable, AccessControl {
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    
    function pause() public {
        require(
            hasRole(DEFAULT_ADMIN_ROLE, msg.sender) || 
            hasRole(PAUSER_ROLE, msg.sender),
            "UNAUTHORIZED"
        );
        _pause();
    }
    
    function unpause() public {
        require(
            hasRole(DEFAULT_ADMIN_ROLE, msg.sender) || 
            hasRole(PAUSER_ROLE, msg.sender),
            "UNAUTHORIZED"
        );
        _unpause();
    }
}
```

## Monitoring

### Events to Monitor

- `PauseProposalCreated`: New pause/unpause proposal
- `ProposalConfirmed`: Guardian confirmation
- `EmergencyPauseExecuted`: Contracts paused
- `EmergencyUnpauseExecuted`: Contracts unpaused
- `PauseExpired`: Automatic expiry triggered

### Status Queries

Get current contract statuses:
```javascript
const [contracts, paused, names] = await emergencyPause.getContractStatuses();
```

Get pause history:
```javascript
const history = await emergencyPause.getPauseHistory(10); // Last 10 events
```

## Testing

Run emergency pause tests:
```bash
npx hardhat test test/EmergencyPauseController.test.js
```

## Contact

- Security Team: security@bogowi.com
- Emergency Hotline: [Establish 24/7 contact]
- Guardian Coordination: [Secure communication channel]