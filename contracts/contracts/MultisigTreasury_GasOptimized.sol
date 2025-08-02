// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC1155/IERC1155.sol";
import "@openzeppelin/contracts/access/IAccessControl.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/**
 * @title MultisigTreasury_GasOptimized
 * @dev Multisig treasury with gas optimizations for batch operations
 */
contract MultisigTreasury_GasOptimized is ReentrancyGuard, Pausable {
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

    // Pagination support
    struct PaginationInfo {
        uint256 totalCount;
        uint256 pageSize;
        uint256 currentPage;
        uint256 totalPages;
    }

    uint256 public threshold;
    uint256 public signerCount;
    uint256 public transactionCount;
    uint256 public constant MAX_SIGNERS = 20;
    uint256 public constant TRANSACTION_EXPIRY = 7 days;
    uint256 public constant EXECUTION_DELAY = 1 hours;
    uint256 public constant MAX_GAS_LIMIT = 5000000;
    
    // Gas optimization constants
    uint256 public constant MAX_BATCH_SIZE = 100; // Maximum items per batch operation
    uint256 public constant DEFAULT_PAGE_SIZE = 50; // Default pagination size
    uint256 public constant MAX_PAGE_SIZE = 100; // Maximum allowed page size
    
    bool public autoExecuteEnabled = true;
    
    mapping(address => Signer) public signers;
    mapping(uint256 => Transaction) public transactions;
    mapping(uint256 => mapping(address => bool)) public confirmations;
    
    EnumerableSet.AddressSet private signerSet;
    
    mapping(address => bool) public emergencyApprovals;
    uint256 public emergencyApprovalCount;
    
    mapping(address => mapping(bytes4 => bool)) public allowedFunctions;
    bool public restrictFunctionCalls = false;
    
    // Track active transactions for efficient querying
    EnumerableSet.UintSet private pendingTransactionIds;
    
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
    event BatchOperationExecuted(string operation, uint256 itemsProcessed);

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

        // Batch size check for constructor
        require(_signers.length <= MAX_BATCH_SIZE, "Batch size exceeded");

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

        // Add to pending transactions set
        pendingTransactionIds.add(txId);

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
            if (block.timestamp >= transactions[_txId].timestamp + EXECUTION_DELAY) {
                executeTransaction(_txId);
            }
        }
    }

    /**
     * @dev Batch confirm multiple transactions with gas limit protection
     * @param _txIds Array of transaction IDs to confirm
     */
    function batchConfirmTransactions(uint256[] calldata _txIds) 
        external 
        onlySigner 
        whenNotPaused 
    {
        require(_txIds.length <= MAX_BATCH_SIZE, "Batch size exceeded");
        
        uint256 confirmed = 0;
        for (uint256 i = 0; i < _txIds.length; i++) {
            uint256 txId = _txIds[i];
            
            // Skip if already confirmed or invalid
            if (txId >= transactionCount || 
                transactions[txId].executed ||
                confirmations[txId][msg.sender] ||
                block.timestamp > transactions[txId].timestamp + TRANSACTION_EXPIRY) {
                continue;
            }
            
            confirmations[txId][msg.sender] = true;
            transactions[txId].confirmations++;
            emit TransactionConfirmed(txId, msg.sender);
            confirmed++;
            
            // Check gas usage periodically
            if (i % 10 == 0 && gasleft() < 100000) {
                break;
            }
        }
        
        emit BatchOperationExecuted("batchConfirm", confirmed);
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
        require(block.timestamp >= txn.timestamp + EXECUTION_DELAY, "Execution delay not met");
        require(txn.to != address(0), "Invalid recipient");
        require(gasleft() >= MAX_GAS_LIMIT / 2, "Insufficient gas");
        
        if (restrictFunctionCalls && txn.data.length >= 4) {
            bytes memory data = txn.data;
            bytes4 selector;
            assembly {
                selector := mload(add(data, 32))
            }
            require(allowedFunctions[txn.to][selector], "Function not allowed");
        }
        
        txn.executed = true;
        
        // Remove from pending transactions
        pendingTransactionIds.remove(_txId);
        
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
        
        transactions[_txId].executed = true;
        pendingTransactionIds.remove(_txId);
        
        emit TransactionCancelled(_txId);
    }

    /**
     * @dev Batch cancel expired transactions with gas limit protection
     * @param _txIds Array of transaction IDs to cancel
     */
    function batchCancelExpiredTransactions(uint256[] calldata _txIds) 
        external 
        onlySigner 
    {
        require(_txIds.length <= MAX_BATCH_SIZE, "Batch size exceeded");
        
        uint256 cancelled = 0;
        for (uint256 i = 0; i < _txIds.length; i++) {
            uint256 txId = _txIds[i];
            
            if (txId < transactionCount && 
                !transactions[txId].executed &&
                block.timestamp > transactions[txId].timestamp + TRANSACTION_EXPIRY) {
                
                transactions[txId].executed = true;
                pendingTransactionIds.remove(txId);
                emit TransactionCancelled(txId);
                cancelled++;
            }
            
            // Check gas usage periodically
            if (i % 10 == 0 && gasleft() < 100000) {
                break;
            }
        }
        
        emit BatchOperationExecuted("batchCancel", cancelled);
    }

    // Signer Management Functions
    
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

    /**
     * @dev Batch add multiple signers with gas limit protection
     * @param _signers Array of addresses to add as signers
     */
    function batchAddSigners(address[] calldata _signers) external onlyMultisig {
        require(_signers.length <= MAX_BATCH_SIZE, "Batch size exceeded");
        require(signerCount + _signers.length <= MAX_SIGNERS, "Would exceed max signers");
        
        for (uint256 i = 0; i < _signers.length; i++) {
            address signer = _signers[i];
            require(signer != address(0), "Invalid signer");
            require(!signers[signer].isSigner, "Duplicate signer");
            
            signers[signer] = Signer({
                isSigner: true,
                addedAt: block.timestamp
            });
            signerSet.add(signer);
            signerCount++;
            
            emit SignerAdded(signer);
        }
        
        emit BatchOperationExecuted("batchAddSigners", _signers.length);
    }

    function removeSigner(address _signer) external onlyMultisig {
        require(signers[_signer].isSigner, "Not a signer");
        require(signerCount - 1 >= threshold, "Would break threshold");
        
        signers[_signer].isSigner = false;
        signerCount--;
        signerSet.remove(_signer);
        
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

    // Token Management Functions (unchanged)
    
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
        require(!emergencyApprovals[msg.sender], "Already approved emergency");
        require(_amount <= address(this).balance / 2, "Amount exceeds 50% of balance");
        require(_to != address(0), "Invalid recipient");
        
        emergencyApprovals[msg.sender] = true;
        emergencyApprovalCount++;
        
        emit EmergencyApprovalGranted(msg.sender);
        
        if (emergencyApprovalCount >= threshold) {
            // Batch reset with gas limit protection
            address[] memory signerAddresses = signerSet.values();
            uint256 maxReset = signerAddresses.length > MAX_BATCH_SIZE ? MAX_BATCH_SIZE : signerAddresses.length;
            
            for (uint256 i = 0; i < maxReset; i++) {
                if (emergencyApprovals[signerAddresses[i]]) {
                    emergencyApprovals[signerAddresses[i]] = false;
                }
            }
            emergencyApprovalCount = 0;
            
            _to.transfer(_amount);
            emit EmergencyWithdraw(address(0), _to, _amount);
        }
    }

    // Optimized View Functions with Pagination
    
    function getSigners() external view returns (address[] memory) {
        return signerSet.values();
    }

    /**
     * @dev Get paginated list of signers
     * @param _page Page number (0-indexed)
     * @param _pageSize Number of items per page (max MAX_PAGE_SIZE)
     */
    function getSignersPaginated(uint256 _page, uint256 _pageSize) 
        external 
        view 
        returns (
            address[] memory signers,
            PaginationInfo memory pagination
        ) 
    {
        require(_pageSize > 0 && _pageSize <= MAX_PAGE_SIZE, "Invalid page size");
        
        address[] memory allSigners = signerSet.values();
        uint256 totalCount = allSigners.length;
        uint256 totalPages = (totalCount + _pageSize - 1) / _pageSize;
        
        require(_page < totalPages, "Page out of bounds");
        
        uint256 start = _page * _pageSize;
        uint256 end = start + _pageSize;
        if (end > totalCount) {
            end = totalCount;
        }
        
        signers = new address[](end - start);
        for (uint256 i = 0; i < end - start; i++) {
            signers[i] = allSigners[start + i];
        }
        
        pagination = PaginationInfo({
            totalCount: totalCount,
            pageSize: _pageSize,
            currentPage: _page,
            totalPages: totalPages
        });
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

    /**
     * @dev Get paginated list of pending transactions
     * @param _page Page number (0-indexed)
     * @param _pageSize Number of items per page (max MAX_PAGE_SIZE)
     */
    function getPendingTransactionsPaginated(uint256 _page, uint256 _pageSize) 
        external 
        view 
        returns (
            uint256[] memory txIds,
            PaginationInfo memory pagination
        ) 
    {
        require(_pageSize > 0 && _pageSize <= MAX_PAGE_SIZE, "Invalid page size");
        
        uint256[] memory allPending = pendingTransactionIds.values();
        uint256 totalCount = allPending.length;
        uint256 totalPages = (totalCount + _pageSize - 1) / _pageSize;
        
        if (_page >= totalPages && totalPages > 0) {
            revert("Page out of bounds");
        }
        
        uint256 start = _page * _pageSize;
        uint256 end = start + _pageSize;
        if (end > totalCount) {
            end = totalCount;
        }
        
        txIds = new uint256[](end - start);
        for (uint256 i = 0; i < end - start; i++) {
            txIds[i] = allPending[start + i];
        }
        
        pagination = PaginationInfo({
            totalCount: totalCount,
            pageSize: _pageSize,
            currentPage: _page,
            totalPages: totalPages
        });
    }

    /**
     * @dev Get count of pending transactions (gas efficient)
     */
    function getPendingTransactionCount() external view returns (uint256) {
        return pendingTransactionIds.length();
    }

    function getConfirmationCount(uint256 _txId) external view returns (uint256) {
        return transactions[_txId].confirmations;
    }

    function hasConfirmed(uint256 _txId, address _signer) external view returns (bool) {
        return confirmations[_txId][_signer];
    }

    function isTransactionExpired(uint256 _txId) external view returns (bool) {
        return block.timestamp > transactions[_txId].timestamp + TRANSACTION_EXPIRY;
    }

    // Token receipt functions (unchanged)
    
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
        uint256,
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
    
    // Configuration functions
    
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