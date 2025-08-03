// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../base/RoleManaged.sol";
import "../interfaces/IRoleManager.sol";

/**
 * @title BOGOTokenV2_RoleManaged
 * @dev Enhanced BOGO token with centralized role management, supply management, and timelock governance
 */
contract BOGOTokenV2_RoleManaged is ERC20, RoleManaged, Pausable, ReentrancyGuard {
    // Supply constants
    uint256 public constant MAX_SUPPLY = 21_000_000 * 10**18; // 21 million tokens
    uint256 public constant DAO_ALLOCATION = 11_550_000 * 10**18; // 55% (11.55M)
    uint256 public constant BUSINESS_ALLOCATION = 9_450_000 * 10**18; // 45% (9.45M)
    
    // Minted amounts tracking
    uint256 public daoMinted;
    uint256 public businessMinted;
    uint256 public rewardsMinted;
    
    // Timelock for governance operations
    mapping(bytes32 => TimelockOperation) public timelockOperations;
    uint256 public constant TIMELOCK_DELAY = 48 hours;
    
    struct TimelockOperation {
        address target;
        uint256 value;
        bytes data;
        uint256 executeTime;
        bool executed;
    }
    
    // Events
    event DAOAllocationMinted(address indexed to, uint256 amount);
    event BusinessAllocationMinted(address indexed to, uint256 amount);
    event RewardsAllocated(address indexed to, uint256 amount);
    event TimelockOperationScheduled(bytes32 indexed operationId, uint256 executeTime);
    event TimelockOperationExecuted(bytes32 indexed operationId);
    event TimelockOperationCancelled(bytes32 indexed operationId);
    
    constructor(
        address _roleManager,
        string memory _name,
        string memory _symbol
    ) ERC20(_name, _symbol) RoleManaged(_roleManager) {
        // Initial setup handled by RoleManager
    }
    
    /**
     * @dev Mint tokens from DAO allocation
     */
    function mintFromDAO(address to, uint256 amount) external onlyRole(roleManager.DAO_ROLE()) nonReentrant {
        require(to != address(0), "Invalid recipient");
        require(amount > 0, "Amount must be greater than 0");
        require(daoMinted + amount <= DAO_ALLOCATION, "Exceeds DAO allocation");
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        
        daoMinted += amount;
        _mint(to, amount);
        
        emit DAOAllocationMinted(to, amount);
    }
    
    /**
     * @dev Mint tokens from Business allocation
     */
    function mintFromBusiness(address to, uint256 amount) external onlyRole(roleManager.BUSINESS_ROLE()) nonReentrant {
        require(to != address(0), "Invalid recipient");
        require(amount > 0, "Amount must be greater than 0");
        require(businessMinted + amount <= BUSINESS_ALLOCATION, "Exceeds business allocation");
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        
        businessMinted += amount;
        _mint(to, amount);
        
        emit BusinessAllocationMinted(to, amount);
    }
    
    /**
     * @dev Allocate rewards (can be used by DAO or Business roles)
     */
    function allocateRewards(address to, uint256 amount) external nonReentrant {
        require(hasRole(roleManager.DAO_ROLE(), msg.sender) || hasRole(roleManager.BUSINESS_ROLE(), msg.sender), 
                "Must have DAO or BUSINESS role");
        require(to != address(0), "Invalid recipient");
        require(amount > 0, "Amount must be greater than 0");
        
        // Check which allocation to use based on caller's role
        if (hasRole(roleManager.DAO_ROLE(), msg.sender)) {
            require(daoMinted + amount <= DAO_ALLOCATION, "Exceeds DAO allocation");
            daoMinted += amount;
        } else {
            require(businessMinted + amount <= BUSINESS_ALLOCATION, "Exceeds business allocation");
            businessMinted += amount;
        }
        
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        
        rewardsMinted += amount;
        _mint(to, amount);
        
        emit RewardsAllocated(to, amount);
    }
    
    /**
     * @dev Schedule a timelock operation
     */
    function scheduleTimelockOperation(
        address target,
        uint256 value,
        bytes calldata data
    ) external onlyRole(roleManager.DEFAULT_ADMIN_ROLE()) {
        bytes32 operationId = keccak256(abi.encode(target, value, data, block.timestamp));
        require(timelockOperations[operationId].executeTime == 0, "Operation already scheduled");
        
        timelockOperations[operationId] = TimelockOperation({
            target: target,
            value: value,
            data: data,
            executeTime: block.timestamp + TIMELOCK_DELAY,
            executed: false
        });
        
        emit TimelockOperationScheduled(operationId, block.timestamp + TIMELOCK_DELAY);
    }
    
    /**
     * @dev Execute a timelock operation
     */
    function executeTimelockOperation(bytes32 operationId) 
        external onlyRole(roleManager.DEFAULT_ADMIN_ROLE()) {
        TimelockOperation storage operation = timelockOperations[operationId];
        require(operation.executeTime != 0, "Operation not found");
        require(block.timestamp >= operation.executeTime, "Timelock not expired");
        require(!operation.executed, "Operation already executed");
        
        operation.executed = true;
        
        (bool success, ) = operation.target.call{value: operation.value}(operation.data);
        require(success, "Operation failed");
        
        emit TimelockOperationExecuted(operationId);
    }
    
    /**
     * @dev Cancel a timelock operation
     */
    function cancelTimelockOperation(bytes32 operationId) external onlyRole(roleManager.DEFAULT_ADMIN_ROLE()) {
        require(timelockOperations[operationId].executeTime != 0, "Operation not found");
        require(!timelockOperations[operationId].executed, "Operation already executed");
        
        delete timelockOperations[operationId];
        
        emit TimelockOperationCancelled(operationId);
    }
    
    /**
     * @dev Pause token transfers
     */
    function pause() external onlyRole(roleManager.PAUSER_ROLE()) {
        _pause();
    }
    
    /**
     * @dev Unpause token transfers
     */
    function unpause() external onlyRole(roleManager.PAUSER_ROLE()) {
        _unpause();
    }
    
    /**
     * @dev Override _update to include pausable functionality
     */
    function _update(address from, address to, uint256 amount) internal override whenNotPaused {
        super._update(from, to, amount);
    }
    
    /**
     * @dev View function to check remaining allocations
     */
    function getRemainingAllocations() external view returns (
        uint256 daoRemaining,
        uint256 businessRemaining,
        uint256 totalRemaining
    ) {
        daoRemaining = DAO_ALLOCATION - daoMinted;
        businessRemaining = BUSINESS_ALLOCATION - businessMinted;
        totalRemaining = MAX_SUPPLY - totalSupply();
    }
}