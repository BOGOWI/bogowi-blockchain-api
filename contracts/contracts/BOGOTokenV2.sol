// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./StandardErrors.sol";

/**
 * @title BOGOTokenV2
 * @dev Enhanced BOGO token with role-based access control, supply management, and timelock governance
 * @dev Replaces the basic Token.sol contract with enterprise-grade features
 */
contract BOGOTokenV2 is ERC20, AccessControl, Pausable, ReentrancyGuard, StandardErrors {
    // Role definitions
    bytes32 public constant DAO_ROLE = keccak256("DAO_ROLE");
    bytes32 public constant BUSINESS_ROLE = keccak256("BUSINESS_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

    // Supply and allocation constants
    uint256 public constant MAX_SUPPLY = 1_000_000_000 * 10**18; // 1 billion tokens
    uint256 public constant DAO_ALLOCATION = 200_000_000 * 10**18; // 200M for DAO
    uint256 public constant BUSINESS_ALLOCATION = 300_000_000 * 10**18; // 300M for business
    uint256 public constant REWARDS_ALLOCATION = 500_000_000 * 10**18; // 500M for rewards

    // Allocation tracking
    uint256 public daoMinted;
    uint256 public businessMinted;
    uint256 public rewardsMinted;

    // Flavored token system
    mapping(string => address) public flavoredTokens;
    mapping(bytes32 => address) private flavoredTokensByHash; // L2 fix: gas optimization
    
    // Timelock mechanism (M1 + L1 fixes)
    uint256 public constant TIMELOCK_DURATION = 2 days;
    mapping(bytes32 => uint256) public timelockOperations;

    // Events
    event AllocationMinted(string indexed allocationType, uint256 amount, address indexed recipient);
    event FlavoredTokenRegistered(string indexed flavor, address indexed tokenAddress);
    event TimelockQueued(bytes32 indexed operationId, uint256 executeTime);
    event TimelockExecuted(bytes32 indexed operationId);
    event TimelockCancelled(bytes32 indexed operationId);

    constructor() ERC20("BOGOWI", "BOGO") {
        // Grant admin roles to deployer
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
    }

    // DAO allocation minting
    function mintFromDAO(address to, uint256 amount) external onlyRole(DAO_ROLE) nonReentrant {
        require(daoMinted + amount <= DAO_ALLOCATION, EXCEEDS_ALLOCATION);
        require(totalSupply() + amount <= MAX_SUPPLY, EXCEEDS_SUPPLY);
        
        daoMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("DAO", amount, to);
    }

    // Business allocation minting
    function mintFromBusiness(address to, uint256 amount) external onlyRole(BUSINESS_ROLE) nonReentrant {
        require(businessMinted + amount <= BUSINESS_ALLOCATION, EXCEEDS_ALLOCATION);
        require(totalSupply() + amount <= MAX_SUPPLY, EXCEEDS_SUPPLY);
        
        businessMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("Business", amount, to);
    }

    // Rewards allocation minting (can be used by DAO or Business roles)
    function mintFromRewards(address to, uint256 amount) external nonReentrant {
        require(hasRole(DAO_ROLE, msg.sender) || hasRole(BUSINESS_ROLE, msg.sender), 
                UNAUTHORIZED);
        require(rewardsMinted + amount <= REWARDS_ALLOCATION, EXCEEDS_ALLOCATION);
        require(totalSupply() + amount <= MAX_SUPPLY, EXCEEDS_SUPPLY);
        
        rewardsMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("Rewards", amount, to);
    }

    // Timelock governance for flavored token registration (M1 + L1 fixes)
    function queueRegisterFlavoredToken(string memory flavor, address tokenAddress) 
        external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(tokenAddress != address(0), ZERO_ADDRESS); // L1 fix
        require(_isContract(tokenAddress), INVALID_ADDRESS); // L1 fix
        
        bytes32 operationId = keccak256(
            abi.encodePacked("registerFlavoredToken", flavor, tokenAddress)
        );
        
        uint256 executeTime = block.timestamp + TIMELOCK_DURATION;
        timelockOperations[operationId] = executeTime;
        
        emit TimelockQueued(operationId, executeTime);
    }

    function executeRegisterFlavoredToken(string memory flavor, address tokenAddress) 
        external onlyRole(DEFAULT_ADMIN_ROLE) {
        bytes32 operationId = keccak256(
            abi.encodePacked("registerFlavoredToken", flavor, tokenAddress)
        );
        
        uint256 executeTime = timelockOperations[operationId];
        require(executeTime != 0, NOT_INITIALIZED);
        require(block.timestamp >= executeTime, NOT_EXPIRED);
        
        // Execute the registration
        flavoredTokens[flavor] = tokenAddress;
        
        // L2 fix: Store hash mapping for gas optimization
        bytes32 flavorHash = keccak256(bytes(flavor));
        flavoredTokensByHash[flavorHash] = tokenAddress;
        
        // Clear the timelock
        delete timelockOperations[operationId];
        
        emit FlavoredTokenRegistered(flavor, tokenAddress);
        emit TimelockExecuted(operationId);
    }

    function cancelTimelockOperation(bytes32 operationId) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(timelockOperations[operationId] != 0, NOT_INITIALIZED);
        
        delete timelockOperations[operationId];
        emit TimelockCancelled(operationId);
    }

    // Utility functions
    function getRemainingDAOAllocation() external view returns (uint256) {
        return DAO_ALLOCATION - daoMinted;
    }

    function getRemainingBusinessAllocation() external view returns (uint256) {
        return BUSINESS_ALLOCATION - businessMinted;
    }

    function getRemainingRewardsAllocation() external view returns (uint256) {
        return REWARDS_ALLOCATION - rewardsMinted;
    }

    // L2 fix: Gas-optimized flavor lookup by hash
    function getFlavoredTokenByHash(bytes32 flavorHash) external view returns (address) {
        return flavoredTokensByHash[flavorHash];
    }

    // Pausable functionality
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }

    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
    }

    // Token burning for deflationary mechanics
    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }

    function burnFrom(address account, uint256 amount) external {
        _spendAllowance(account, msg.sender, amount);
        _burn(account, amount);
    }

    // Override transfer functions to respect pause state
    function _update(address from, address to, uint256 value)
        internal
        override
        whenNotPaused
    {
        super._update(from, to, value);
    }

    // L1 fix: Contract address validation
    function _isContract(address account) internal view returns (bool) {
        uint256 size;
        assembly {
            size := extcodesize(account)
        }
        return size > 0;
    }

    // Required by AccessControl
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
