// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC1155/IERC1155.sol";
import "@openzeppelin/contracts/access/IAccessControl.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

contract MultisigTreasury is ReentrancyGuard, Pausable {
    using SafeERC20 for IERC20;
    using EnumerableSet for EnumerableSet.AddressSet;

    struct Transaction {
        address to;
        uint256 value;
        bytes data;
        string description;
        bool executed;
        uint256 confirmations;
        uint256 timestamp;
    }

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
    event AutoExecuteToggled(bool enabled);
    event FunctionAllowanceSet(address indexed target, bytes4 indexed selector, bool allowed);
    event EmergencyApprovalGranted(address indexed signer);
    event EmergencyApprovalRevoked(address indexed signer);

    modifier onlySigner() {
        require(signers[msg.sender].isSigner, "Not a signer");
        _;
    }

    modifier onlyMultisig() {
        require(msg.sender == address(this), "Only multisig");
        _;
    }

    modifier transactionExists(uint256 _txId) {
        require(_txId < transactionCount, "Transaction does not exist");
        _;
    }

    modifier notExecuted(uint256 _txId) {
        require(!transactions[_txId].executed, "Transaction already executed");
        _;
    }

    modifier notExpired(uint256 _txId) {
        require(
            block.timestamp <= transactions[_txId].timestamp + TRANSACTION_EXPIRY,
            "Transaction expired"
        );
        _;
    }

    constructor(address[] memory _signers, uint256 _threshold) {
        require(_signers.length > 0, "Signers required");
        require(_signers.length <= MAX_SIGNERS, "Too many signers");
        require(_threshold > 0 && _threshold <= _signers.length, "Invalid threshold");

        for (uint256 i = 0; i < _signers.length; i++) {
            address signer = _signers[i];
            require(signer != address(0), "Invalid signer");
            require(!signers[signer].isSigner, "Duplicate signer");
            
            signers[signer] = Signer({
                isSigner: true,
                addedAt: block.timestamp
            });
            signerSet.add(signer);
        }

        signerCount = _signers.length;
        threshold = _threshold;
    }

    receive() external payable {
        emit Deposit(msg.sender, msg.value);
    }

    function submitTransaction(
        address _to,
        uint256 _value,
        bytes memory _data,
        string memory _description
    ) external onlySigner whenNotPaused returns (uint256) {
        require(_to != address(0), "Invalid recipient");
        
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

    function confirmTransaction(uint256 _txId) 
        public 
        onlySigner 
        transactionExists(_txId) 
        notExecuted(_txId)
        notExpired(_txId)
    {
        require(!confirmations[_txId][msg.sender], "Already confirmed");
        
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

    function revokeConfirmation(uint256 _txId)
        external
        onlySigner
        transactionExists(_txId)
        notExecuted(_txId)
    {
        require(confirmations[_txId][msg.sender], "Not confirmed");
        
        confirmations[_txId][msg.sender] = false;
        transactions[_txId].confirmations--;
        
        emit ConfirmationRevoked(_txId, msg.sender);
    }

    function executeTransaction(uint256 _txId)
        public
        onlySigner
        transactionExists(_txId)
        notExecuted(_txId)
        notExpired(_txId)
        nonReentrant
    {
        Transaction storage txn = transactions[_txId];
        require(txn.confirmations >= threshold, "Insufficient confirmations");
        
        // M1: Ensure execution delay has passed
        require(block.timestamp >= txn.timestamp + EXECUTION_DELAY, "Execution delay not met");
        
        // M2: Validate transaction parameters
        require(txn.to != address(0), "Invalid recipient");
        require(gasleft() >= MAX_GAS_LIMIT / 2, "Insufficient gas");
        
        // M3: Check function call restrictions if enabled
        if (restrictFunctionCalls && txn.data.length >= 4) {
            bytes memory data = txn.data;
            bytes4 selector;
            assembly {
                selector := mload(add(data, 32))
            }
            require(allowedFunctions[txn.to][selector], "Function not allowed");
        }
        
        txn.executed = true;
        
        (bool success, bytes memory returnData) = txn.to.call{value: txn.value, gas: MAX_GAS_LIMIT}(txn.data);
        require(success, string(returnData));
        
        emit TransactionExecuted(_txId, msg.sender);
    }

    function cancelExpiredTransaction(uint256 _txId)
        external
        onlySigner
        transactionExists(_txId)
        notExecuted(_txId)
    {
        require(
            block.timestamp > transactions[_txId].timestamp + TRANSACTION_EXPIRY,
            "Transaction not expired"
        );
        
        transactions[_txId].executed = true; // Mark as executed to prevent future execution
        emit TransactionCancelled(_txId);
    }

    // Signer Management Functions (only callable by multisig itself)
    
    function addSigner(address _signer) external onlyMultisig {
        require(_signer != address(0), "Invalid signer");
        require(!signers[_signer].isSigner, "Already a signer");
        require(signerCount < MAX_SIGNERS, "Max signers reached");
        
        signers[_signer] = Signer({
            isSigner: true,
            addedAt: block.timestamp
        });
        signerSet.add(_signer);
        signerCount++;
        
        emit SignerAdded(_signer);
    }

    function removeSigner(address _signer) external onlyMultisig {
        require(signers[_signer].isSigner, "Not a signer");
        require(signerCount - 1 >= threshold, "Would break threshold");
        
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

    function changeThreshold(uint256 _threshold) external onlyMultisig {
        require(_threshold > 0 && _threshold <= signerCount, "Invalid threshold");
        
        uint256 oldThreshold = threshold;
        threshold = _threshold;
        
        emit ThresholdChanged(oldThreshold, _threshold);
    }

    // Token Management Functions
    
    function grantRole(
        address _contract,
        bytes32 _role,
        address _account
    ) external onlyMultisig {
        IAccessControl(_contract).grantRole(_role, _account);
    }

    function revokeRole(
        address _contract,
        bytes32 _role,
        address _account
    ) external onlyMultisig {
        IAccessControl(_contract).revokeRole(_role, _account);
    }

    function transferERC20(
        address _token,
        address _to,
        uint256 _amount
    ) external onlyMultisig {
        IERC20(_token).safeTransfer(_to, _amount);
    }

    function transferERC721(
        address _token,
        address _to,
        uint256 _tokenId
    ) external onlyMultisig {
        IERC721(_token).transferFrom(address(this), _to, _tokenId);
    }

    function transferERC1155(
        address _token,
        address _to,
        uint256 _id,
        uint256 _amount,
        bytes memory _data
    ) external onlyMultisig {
        IERC1155(_token).safeTransferFrom(address(this), _to, _id, _amount, _data);
    }

    // Emergency Functions
    
    function pause() external onlyMultisig {
        _pause();
    }

    function unpause() external onlyMultisig {
        _unpause();
    }

    function emergencyWithdrawETH(address payable _to, uint256 _amount) 
        external 
        whenPaused 
        onlySigner
        nonReentrant
    {
        // H1: Use separate mapping for emergency approvals
        require(!emergencyApprovals[msg.sender], "Already approved emergency");
        
        // M4: Additional authorization check
        require(_amount <= address(this).balance / 2, "Amount exceeds 50% of balance");
        require(_to != address(0), "Invalid recipient");
        
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

    // View Functions
    
    function getSigners() external view returns (address[] memory) {
        return signerSet.values();
    }

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

    function getConfirmationCount(uint256 _txId) external view returns (uint256) {
        return transactions[_txId].confirmations;
    }

    function hasConfirmed(uint256 _txId, address _signer) external view returns (bool) {
        return confirmations[_txId][_signer];
    }

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

    function isTransactionExpired(uint256 _txId) external view returns (bool) {
        return block.timestamp > transactions[_txId].timestamp + TRANSACTION_EXPIRY;
    }

    // Token receipt functions
    
    function onERC721Received(
        address,
        address from,
        uint256 tokenId,
        bytes memory
    ) external returns (bytes4) {
        emit TokensReceived(msg.sender, from, tokenId);
        return this.onERC721Received.selector;
    }

    function onERC1155Received(
        address,
        address from,
        uint256 id,
        uint256 value,
        bytes memory
    ) external returns (bytes4) {
        emit TokensReceived(msg.sender, from, value);
        return this.onERC1155Received.selector;
    }

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
    
    // New view functions for enhanced functionality
    
    function toggleAutoExecute() external onlyMultisig {
        autoExecuteEnabled = !autoExecuteEnabled;
        emit AutoExecuteToggled(autoExecuteEnabled);
    }
    
    function setFunctionAllowance(address _target, bytes4 _selector, bool _allowed) external onlyMultisig {
        allowedFunctions[_target][_selector] = _allowed;
        emit FunctionAllowanceSet(_target, _selector, _allowed);
    }
    
    function toggleFunctionRestrictions() external onlyMultisig {
        restrictFunctionCalls = !restrictFunctionCalls;
    }
    
    function getEmergencyApprovalCount() external view returns (uint256) {
        return emergencyApprovalCount;
    }
    
    function hasEmergencyApproval(address _signer) external view returns (bool) {
        return emergencyApprovals[_signer];
    }
}
