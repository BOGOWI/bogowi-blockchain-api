// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC1155/IERC1155.sol";
import "@openzeppelin/contracts/access/IAccessControl.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "./StandardErrors.sol";

/**
 * @title MultisigTreasury
 * @author BOGOWI Team
 * @notice Multi-signature treasury contract for secure fund management
 * @dev Implements time-locked transactions, emergency withdrawals, and role management
 * Features:
 * - Multi-signature transaction approval
 * - Time-locked execution with expiry
 * - Emergency withdrawal mechanism
 * - ERC20/ERC721/ERC1155 token support
 * - Function call restrictions
 * - Auto-execution option
 * @custom:security-contact security@bogowi.com
 */
contract MultisigTreasury is ReentrancyGuard, Pausable, StandardErrors {
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.AddressSet;

    /**
     * @dev Transaction structure for queued operations
     * @param to Target address for the transaction
     * @param value ETH value to send with transaction
     * @param data Encoded function call data
     * @param description Human-readable description
     * @param executed Whether transaction has been executed
     * @param confirmations Number of signer confirmations
     * @param timestamp When transaction was submitted
     */
    struct Transaction {
        address to;
        uint256 value;
        bytes data;
        string description;
        bool executed;
        uint256 confirmations;
        uint256 timestamp;
    }

    /**
     * @dev Signer information structure
     * @param isSigner Whether address is an active signer
     * @param addedAt Timestamp when signer was added
     */
    struct Signer {
        bool isSigner;
        uint256 addedAt;
    }

    uint256 public threshold;
    uint256 public signerCount;
    uint256 public transactionCount;
    uint256 public constant MAX_SIGNERS = 20;
    uint256 public constant TRANSACTION_EXPIRY = 7 days;
    uint256 public constant EXECUTION_DELAY = 1 hours; // M1: Add execution delay
    uint256 public constant MAX_GAS_LIMIT = 5000000; // M2: Add gas limit
    bool public autoExecuteEnabled = true; // M1: Make auto-execution optional
    
    mapping(address => Signer) public signers;
    mapping(uint256 => Transaction) public transactions;
    mapping(uint256 => mapping(address => bool)) public confirmations;
    
    EnumerableSet.AddressSet private signerSet; // H2: Use EnumerableSet for signers
    
    // H1: Separate mapping for emergency approvals
    mapping(address => bool) public emergencyApprovals;
    uint256 public emergencyApprovalCount;
    
    // M3: Function call restrictions
    mapping(address => mapping(bytes4 => bool)) public allowedFunctions;
    bool public restrictFunctionCalls = false;
    
    event SignerAdded(address indexed signer);
    event SignerRemoved(address indexed signer);
    event ThresholdChanged(uint256 oldThreshold, uint256 newThreshold);
    event TransactionSubmitted(uint256 indexed txId, address indexed submitter, address to, uint256 value);
    event TransactionConfirmed(uint256 indexed txId, address indexed signer);
    event ConfirmationRevoked(uint256 indexed txId, address indexed signer);
    event TransactionExecuted(uint256 indexed txId, address indexed executor);
    event TransactionCancelled(uint256 indexed txId);
    event Deposit(address indexed sender, uint256 amount);
    event TokensReceived(address indexed token, address indexed from, uint256 amount);
    event EmergencyWithdraw(address indexed token, address indexed to, uint256 amount);
    event FunctionRestrictionsToggled(bool enabled);
    event AutoExecuteToggled(bool enabled);
    event FunctionAllowanceSet(address indexed target, bytes4 indexed selector, bool allowed);
    event EmergencyApprovalGranted(address indexed signer);
    event EmergencyApprovalRevoked(address indexed signer);

    modifier onlySigner() {
        require(signers[msg.sender].isSigner, NOT_SIGNER);
        _;
    }

    modifier onlyMultisig() {
        require(msg.sender == address(this), NOT_MULTISIG);
        _;
    }

    modifier transactionExists(uint256 _txId) {
        require(_txId < transactionCount, DOES_NOT_EXIST);
        _;
    }

    modifier notExecuted(uint256 _txId) {
        require(!transactions[_txId].executed, ALREADY_PROCESSED);
        _;
    }

    modifier notExpired(uint256 _txId) {
        require(
            block.timestamp <= transactions[_txId].timestamp + TRANSACTION_EXPIRY,
            EXPIRED
        );
        _;
    }

    /**
     * @notice Initializes the multisig treasury with signers and threshold
     * @dev Validates all parameters and sets up initial signer set
     * @param _signers Array of initial signer addresses
     * @param _threshold Number of confirmations required for execution
     */
    constructor(address[] memory _signers, uint256 _threshold) {
        require(_signers.length > 0, INVALID_PARAMETER);
        require(_signers.length <= MAX_SIGNERS, EXCEEDS_LIMIT);
        require(_threshold > 0 && _threshold <= _signers.length, INVALID_PARAMETER);

        for (uint256 i = 0; i < _signers.length; i++) {
            address signer = _signers[i];
            require(signer != address(0), ZERO_ADDRESS);
            require(!signers[signer].isSigner, ALREADY_EXISTS);
            
            signers[signer] = Signer({
                isSigner: true,
                addedAt: block.timestamp
            });
            signerSet.add(signer);
        }

        signerCount = _signers.length;
        threshold = _threshold;
    }

    /**
     * @notice Receives ETH deposits
     * @dev Emits Deposit event for tracking
     */
    receive() external payable {
        emit Deposit(msg.sender, msg.value);
    }

    /**
     * @notice Submits a new transaction for multisig approval
     * @dev Auto-confirms for the submitter, requires signer role
     * @param _to Target address for the transaction
     * @param _value ETH amount to send (can be 0 for contract calls)
     * @param _data Encoded function call data (empty for simple transfers)
     * @param _description Human-readable description of the transaction
     * @return txId The ID of the newly created transaction
     * @custom:emits TransactionSubmitted, TransactionConfirmed
     */
    function submitTransaction(
        address _to,
        uint256 _value,
        bytes memory _data,
        string memory _description
    ) external onlySigner whenNotPaused returns (uint256) {
        require(_to != address(0), ZERO_ADDRESS);
        
        uint256 txId = transactionCount++;
        
        transactions[txId] = Transaction({
            to: _to,
            value: _value,
            data: _data,
            description: _description,
            executed: false,
            confirmations: 0,
            timestamp: block.timestamp
        });

        emit TransactionSubmitted(txId, msg.sender, _to, _value);
        
        // Auto-confirm for submitter
        confirmTransaction(txId);
        
        return txId;
    }

    /**
     * @notice Confirms a pending transaction
     * @dev Auto-executes if threshold is reached and delay has passed
     * @param _txId ID of the transaction to confirm
     * @custom:emits TransactionConfirmed
     * @custom:security May trigger auto-execution if threshold met
     */
    function confirmTransaction(uint256 _txId) 
        public 
        onlySigner 
        transactionExists(_txId) 
        notExecuted(_txId)
        notExpired(_txId)
    {
        require(!confirmations[_txId][msg.sender], ALREADY_PROCESSED);
        
        confirmations[_txId][msg.sender] = true;
        transactions[_txId].confirmations++;
        
        emit TransactionConfirmed(_txId, msg.sender);
        
        // Auto-execute if enabled and threshold reached
        if (autoExecuteEnabled && transactions[_txId].confirmations >= threshold) {
            // Check if execution delay has passed (M1)
            if (block.timestamp >= transactions[_txId].timestamp + EXECUTION_DELAY) {
                executeTransaction(_txId);
            }
        }
    }

    /**
     * @notice Revokes a previously given confirmation
     * @dev Can only revoke own confirmations on unexecuted transactions
     * @param _txId ID of the transaction to revoke confirmation for
     * @custom:emits ConfirmationRevoked
     */
    function revokeConfirmation(uint256 _txId)
        external
        onlySigner
        transactionExists(_txId)
        notExecuted(_txId)
    {
        require(confirmations[_txId][msg.sender], NOT_INITIALIZED);
        
        confirmations[_txId][msg.sender] = false;
        transactions[_txId].confirmations--;
        
        emit ConfirmationRevoked(_txId, msg.sender);
    }

    /**
     * @notice Executes a confirmed transaction
     * @dev Requires threshold confirmations and execution delay
     * @param _txId ID of the transaction to execute
     * @custom:emits TransactionExecuted
     * @custom:security Enforces execution delay and gas limits
     */
    function executeTransaction(uint256 _txId)
        public
        onlySigner
        transactionExists(_txId)
        notExecuted(_txId)
        notExpired(_txId)
        nonReentrant
    {
        Transaction storage txn = transactions[_txId];
        require(txn.confirmations >= threshold, CONDITIONS_NOT_MET);
        
        // M1: Ensure execution delay has passed
        require(block.timestamp >= txn.timestamp + EXECUTION_DELAY, NOT_READY);
        
        // M2: Validate transaction parameters
        require(txn.to != address(0), "Invalid recipient");
        require(gasleft() >= MAX_GAS_LIMIT / 2, EXCEEDS_LIMIT);
        
        // M3: Check function call restrictions if enabled
        if (restrictFunctionCalls && txn.data.length >= 4) {
            bytes memory data = txn.data;
            bytes4 selector;
            assembly {
                selector := mload(add(data, 32))
            }
            require(allowedFunctions[txn.to][selector], UNAUTHORIZED);
        }
        
        txn.executed = true;
        
        (bool success, bytes memory returnData) = txn.to.call{value: txn.value, gas: MAX_GAS_LIMIT}(txn.data);
        require(success, string(returnData));
        
        emit TransactionExecuted(_txId, msg.sender);
    }

    /**
     * @notice Cancels an expired transaction
     * @dev Can only cancel after TRANSACTION_EXPIRY period
     * @param _txId ID of the expired transaction to cancel
     * @custom:emits TransactionCancelled
     */
    function cancelExpiredTransaction(uint256 _txId)
        external
        onlySigner
        transactionExists(_txId)
        notExecuted(_txId)
    {
        require(
            block.timestamp > transactions[_txId].timestamp + TRANSACTION_EXPIRY,
            NOT_EXPIRED
        );
        
        transactions[_txId].executed = true; // Mark as executed to prevent future execution
        emit TransactionCancelled(_txId);
    }

    /**
     * @notice Adds a new signer to the multisig
     * @dev Can only be called by the multisig itself through a transaction
     * @param _signer Address to add as a new signer
     * @custom:emits SignerAdded
     */
    function addSigner(address _signer) external onlyMultisig {
        require(_signer != address(0), ZERO_ADDRESS);
        require(!signers[_signer].isSigner, ALREADY_EXISTS);
        require(signerCount < MAX_SIGNERS, MAX_REACHED);
        
        signers[_signer] = Signer({
            isSigner: true,
            addedAt: block.timestamp
        });
        signerSet.add(_signer);
        signerCount++;
        
        emit SignerAdded(_signer);
    }

    /**
     * @notice Removes an existing signer from the multisig
     * @dev Ensures threshold remains valid after removal
     * @param _signer Address to remove from signers
     * @custom:emits SignerRemoved
     */
    function removeSigner(address _signer) external onlyMultisig {
        require(signers[_signer].isSigner, NOT_SIGNER);
        require(signerCount - 1 >= threshold, CONDITIONS_NOT_MET);
        
        signers[_signer].isSigner = false;
        signerCount--;
        
        signerSet.remove(_signer); // H2: Use EnumerableSet for safe removal
        
        // Remove any emergency approvals for this signer
        if (emergencyApprovals[_signer]) {
            emergencyApprovals[_signer] = false;
            emergencyApprovalCount--;
        }
        
        emit SignerRemoved(_signer);
    }

    /**
     * @notice Changes the confirmation threshold
     * @dev New threshold must be valid for current signer count
     * @param _threshold New number of confirmations required
     * @custom:emits ThresholdChanged
     */
    function changeThreshold(uint256 _threshold) external onlyMultisig {
        require(_threshold > 0 && _threshold <= signerCount, INVALID_PARAMETER);
        
        uint256 oldThreshold = threshold;
        threshold = _threshold;
        
        emit ThresholdChanged(oldThreshold, _threshold);
    }

    /**
     * @notice Grants a role in an external AccessControl contract
     * @dev Used to manage roles in other BOGOWI contracts
     * @param _contract Address of the AccessControl contract
     * @param _role Role identifier to grant
     * @param _account Address to grant the role to
     */
    function grantRole(
        address _contract,
        bytes32 _role,
        address _account
    ) external onlyMultisig {
        IAccessControl(_contract).grantRole(_role, _account);
    }

    /**
     * @notice Revokes a role in an external AccessControl contract
     * @dev Used to manage roles in other BOGOWI contracts
     * @param _contract Address of the AccessControl contract
     * @param _role Role identifier to revoke
     * @param _account Address to revoke the role from
     */
    function revokeRole(
        address _contract,
        bytes32 _role,
        address _account
    ) external onlyMultisig {
        IAccessControl(_contract).revokeRole(_role, _account);
    }

    /**
     * @notice Transfers ERC20 tokens from treasury
     * @dev Uses SafeERC20 for secure transfers
     * @param _token ERC20 token contract address
     * @param _to Recipient address
     * @param _amount Amount of tokens to transfer
     */
    function transferERC20(
        address _token,
        address _to,
        uint256 _amount
    ) external onlyMultisig {
        IERC20(_token).safeTransfer(_to, _amount);
    }

    /**
     * @notice Transfers an ERC721 NFT from treasury
     * @param _token ERC721 token contract address
     * @param _to Recipient address
     * @param _tokenId ID of the NFT to transfer
     */
    function transferERC721(
        address _token,
        address _to,
        uint256 _tokenId
    ) external onlyMultisig {
        IERC721(_token).transferFrom(address(this), _to, _tokenId);
    }

    /**
     * @notice Transfers ERC1155 tokens from treasury
     * @param _token ERC1155 token contract address
     * @param _to Recipient address
     * @param _id Token ID to transfer
     * @param _amount Amount of tokens to transfer
     * @param _data Additional data for the transfer
     */
    function transferERC1155(
        address _token,
        address _to,
        uint256 _id,
        uint256 _amount,
        bytes memory _data
    ) external onlyMultisig {
        IERC1155(_token).safeTransferFrom(address(this), _to, _id, _amount, _data);
    }

    /**
     * @notice Pauses all treasury operations
     * @dev Emergency function, affects submissions and confirmations
     * @custom:security Only for emergency situations
     */
    function pause() external onlyMultisig {
        _pause();
    }

    /**
     * @notice Unpauses treasury operations
     * @dev Resumes normal operations after emergency
     */
    function unpause() external onlyMultisig {
        _unpause();
    }

    /**
     * @notice Emergency ETH withdrawal during pause
     * @dev Requires threshold approvals, max 50% of balance
     * @param _to Recipient address for emergency withdrawal
     * @param _amount Amount of ETH to withdraw (max 50% of balance)
     * @custom:emits EmergencyApprovalGranted, EmergencyWithdraw
     * @custom:security Critical function with special approval process
     */
    function emergencyWithdrawETH(address payable _to, uint256 _amount) 
        external 
        whenPaused 
        onlySigner
        nonReentrant
    {
        // H1: Use separate mapping for emergency approvals
        require(!emergencyApprovals[msg.sender], ALREADY_PROCESSED);
        
        // M4: Additional authorization check
        require(_amount <= address(this).balance / 2, EXCEEDS_LIMIT);
        require(_to != address(0), ZERO_ADDRESS);
        
        // Mark emergency approval
        emergencyApprovals[msg.sender] = true;
        emergencyApprovalCount++;
        
        emit EmergencyApprovalGranted(msg.sender);
        
        // Check if threshold reached
        if (emergencyApprovalCount >= threshold) {
            // Reset emergency approvals
            address[] memory signerAddresses = signerSet.values();
            for (uint256 i = 0; i < signerAddresses.length; i++) {
                if (emergencyApprovals[signerAddresses[i]]) {
                    emergencyApprovals[signerAddresses[i]] = false;
                }
            }
            emergencyApprovalCount = 0;
            
            // Execute withdrawal
            _to.transfer(_amount);
            emit EmergencyWithdraw(address(0), _to, _amount);
        }
    }

    /**
     * @notice Returns all current signers
     * @return Array of signer addresses
     */
    function getSigners() external view returns (address[] memory) {
        return signerSet.values();
    }

    /**
     * @notice Returns detailed transaction information
     * @param _txId ID of the transaction to query
     * @return to Target address
     * @return value ETH value
     * @return data Encoded function data
     * @return description Transaction description
     * @return executed Execution status
     * @return confirmationCount Number of confirmations
     * @return timestamp Submission timestamp
     */
    function getTransaction(uint256 _txId) 
        external 
        view 
        returns (
            address to,
            uint256 value,
            bytes memory data,
            string memory description,
            bool executed,
            uint256 confirmationCount,
            uint256 timestamp
        ) 
    {
        Transaction memory txn = transactions[_txId];
        return (
            txn.to,
            txn.value,
            txn.data,
            txn.description,
            txn.executed,
            txn.confirmations,
            txn.timestamp
        );
    }

    /**
     * @notice Returns the number of confirmations for a transaction
     * @param _txId ID of the transaction
     * @return Number of confirmations
     */
    function getConfirmationCount(uint256 _txId) external view returns (uint256) {
        return transactions[_txId].confirmations;
    }

    /**
     * @notice Checks if a signer has confirmed a transaction
     * @param _txId ID of the transaction
     * @param _signer Address of the signer to check
     * @return True if the signer has confirmed
     */
    function hasConfirmed(uint256 _txId, address _signer) external view returns (bool) {
        return confirmations[_txId][_signer];
    }

    /**
     * @notice Returns all pending (unexecuted, non-expired) transactions
     * @return Array of pending transaction IDs
     */
    function getPendingTransactions() external view returns (uint256[] memory) {
        uint256 pendingCount = 0;
        
        // Count pending transactions
        for (uint256 i = 0; i < transactionCount; i++) {
            if (!transactions[i].executed && 
                block.timestamp <= transactions[i].timestamp + TRANSACTION_EXPIRY) {
                pendingCount++;
            }
        }
        
        // Collect pending transaction IDs
        uint256[] memory pendingTxs = new uint256[](pendingCount);
        uint256 index = 0;
        
        for (uint256 i = 0; i < transactionCount; i++) {
            if (!transactions[i].executed && 
                block.timestamp <= transactions[i].timestamp + TRANSACTION_EXPIRY) {
                pendingTxs[index++] = i;
            }
        }
        
        return pendingTxs;
    }

    /**
     * @notice Checks if a transaction has expired
     * @param _txId ID of the transaction
     * @return True if transaction has expired
     */
    function isTransactionExpired(uint256 _txId) external view returns (bool) {
        return block.timestamp > transactions[_txId].timestamp + TRANSACTION_EXPIRY;
    }

    /**
     * @notice Handles receipt of ERC721 tokens
     * @dev Required for receiving NFTs
     * @param from Address sending the NFT
     * @param tokenId ID of the received NFT
     * @return ERC721 receiver interface selector
     */
    function onERC721Received(
        address,
        address from,
        uint256 tokenId,
        bytes memory
    ) external returns (bytes4) {
        emit TokensReceived(msg.sender, from, tokenId);
        return this.onERC721Received.selector;
    }

    /**
     * @notice Handles receipt of single ERC1155 token
     * @param from Address sending the token
     * @param value Amount received
     * @return ERC1155 receiver interface selector
     */
    function onERC1155Received(
        address,
        address from,
        uint256 /* id */,
        uint256 value,
        bytes memory
    ) external returns (bytes4) {
        emit TokensReceived(msg.sender, from, value);
        return this.onERC1155Received.selector;
    }

    /**
     * @notice Handles receipt of multiple ERC1155 tokens
     * @param from Address sending the tokens
     * @return ERC1155 batch receiver interface selector
     */
    function onERC1155BatchReceived(
        address,
        address from,
        uint256[] memory,
        uint256[] memory,
        bytes memory
    ) external returns (bytes4) {
        emit TokensReceived(msg.sender, from, 0);
        return this.onERC1155BatchReceived.selector;
    }
    
    /**
     * @notice Toggles automatic execution of transactions
     * @dev When enabled, transactions auto-execute when threshold is met
     * @custom:emits AutoExecuteToggled
     */
    function toggleAutoExecute() external onlyMultisig {
        autoExecuteEnabled = !autoExecuteEnabled;
        emit AutoExecuteToggled(autoExecuteEnabled);
    }
    
    /**
     * @notice Sets whether a function selector is allowed for a target
     * @dev Used when function call restrictions are enabled
     * @param _target Contract address
     * @param _selector Function selector (4 bytes)
     * @param _allowed Whether the function is allowed
     * @custom:emits FunctionAllowanceSet
     */
    function setFunctionAllowance(address _target, bytes4 _selector, bool _allowed) external onlyMultisig {
        allowedFunctions[_target][_selector] = _allowed;
        emit FunctionAllowanceSet(_target, _selector, _allowed);
    }
    
    /**
     * @notice Toggles function call restrictions
     * @dev When enabled, only allowed function selectors can be called
     * @custom:emits FunctionRestrictionsToggled
     */
    function toggleFunctionRestrictions() external onlyMultisig {
        restrictFunctionCalls = !restrictFunctionCalls;
        emit FunctionRestrictionsToggled(restrictFunctionCalls);
    }
    
    /**
     * @notice Returns current emergency approval count
     * @return Number of signers who approved emergency withdrawal
     */
    function getEmergencyApprovalCount() external view returns (uint256) {
        return emergencyApprovalCount;
    }
    
    /**
     * @notice Checks if a signer has approved emergency withdrawal
     * @param _signer Address to check
     * @return True if signer has approved emergency withdrawal
     */
    function hasEmergencyApproval(address _signer) external view returns (bool) {
        return emergencyApprovals[_signer];
    }
}
