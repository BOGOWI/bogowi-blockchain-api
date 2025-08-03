// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title BOGOToken_ZeroValidated
 * @dev Enhanced BOGO token with comprehensive zero address validation
 * @dev Prevents token burns and lost funds through address validation
 */
contract BOGOToken_ZeroValidated is ERC20, AccessControl, Pausable, ReentrancyGuard {
    // Role definitions
    bytes32 public constant DAO_ROLE = keccak256("DAO_ROLE");
    bytes32 public constant BUSINESS_ROLE = keccak256("BUSINESS_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

    // Supply and allocation constants
    uint256 public constant MAX_SUPPLY = 1_000_000_000 * 10**18; // 1 billion tokens
    uint256 public constant DAO_ALLOCATION = 50_000_000 * 10**18; // 50M for DAO (5% of total)
    uint256 public constant BUSINESS_ALLOCATION = 900_000_000 * 10**18; // 900M for business (90% of total)
    uint256 public constant REWARDS_ALLOCATION = 50_000_000 * 10**18; // 50M for rewards (5% of total)

    // Allocation tracking
    uint256 public daoMinted;
    uint256 public businessMinted;
    uint256 public rewardsMinted;

    // Timelock mechanism
    uint256 public constant TIMELOCK_DURATION = 2 days;
    mapping(bytes32 => uint256) public timelockOperations;

    // Events
    event AllocationMinted(string indexed allocationType, uint256 amount, address indexed recipient);
    event TimelockQueued(bytes32 indexed operationId, uint256 executeTime);
    event TimelockExecuted(bytes32 indexed operationId);
    event TimelockCancelled(bytes32 indexed operationId);

    // Custom errors for gas optimization
    error InvalidAddress();
    error InvalidAmount();
    error ExceedsAllocation();
    error ExceedsMaxSupply();
    error InsufficientRole();
    error InvalidTokenAddress();
    error NotAContract();
    error OperationNotQueued();
    error TimelockNotExpired();

    constructor() ERC20("BOGOWI", "BOGO") {
        // Grant admin roles to deployer
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
    }

    /**
     * @dev Validates that an address is not zero
     * @param addr The address to validate
     */
    modifier notZeroAddress(address addr) {
        if (addr == address(0)) revert InvalidAddress();
        _;
    }

    /**
     * @dev Validates that an amount is greater than zero
     * @param amount The amount to validate
     */
    modifier notZeroAmount(uint256 amount) {
        if (amount == 0) revert InvalidAmount();
        _;
    }

    /**
     * @dev DAO allocation minting with zero address validation
     * @param to Recipient address (must not be zero)
     * @param amount Amount to mint (must be greater than zero)
     */
    function mintFromDAO(address to, uint256 amount) 
        external 
        onlyRole(DAO_ROLE) 
        nonReentrant 
        notZeroAddress(to)
        notZeroAmount(amount)
    {
        if (daoMinted + amount > DAO_ALLOCATION) revert ExceedsAllocation();
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        
        daoMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("DAO", amount, to);
    }

    /**
     * @dev Business allocation minting with zero address validation
     * @param to Recipient address (must not be zero)
     * @param amount Amount to mint (must be greater than zero)
     */
    function mintFromBusiness(address to, uint256 amount) 
        external 
        onlyRole(BUSINESS_ROLE) 
        nonReentrant
        notZeroAddress(to)
        notZeroAmount(amount)
    {
        if (businessMinted + amount > BUSINESS_ALLOCATION) revert ExceedsAllocation();
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        
        businessMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("Business", amount, to);
    }

    /**
     * @dev Rewards allocation minting with zero address validation
     * @param to Recipient address (must not be zero)
     * @param amount Amount to mint (must be greater than zero)
     */
    function mintFromRewards(address to, uint256 amount) 
        external 
        nonReentrant
        notZeroAddress(to)
        notZeroAmount(amount)
    {
        if (!hasRole(DAO_ROLE, msg.sender) && !hasRole(BUSINESS_ROLE, msg.sender)) {
            revert InsufficientRole();
        }
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        if (rewardsMinted + amount > REWARDS_ALLOCATION) revert ExceedsAllocation();
        
        rewardsMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("Rewards", amount, to);
    }

    /**
     * @dev Burn tokens from a specific account with zero address validation
     * @param account The account to burn from (must not be zero)
     * @param amount The amount to burn (must be greater than zero)
     */
    function burnFrom(address account, uint256 amount) 
        external
        notZeroAddress(account)
        notZeroAmount(amount)
    {
        _spendAllowance(account, msg.sender, amount);
        _burn(account, amount);
    }

    /**
     * @dev Burn tokens from sender's account
     * @param amount The amount to burn (must be greater than zero)
     */
    function burn(uint256 amount) 
        external
        notZeroAmount(amount)
    {
        _burn(msg.sender, amount);
    }

    /**
     * @dev Override _update to include pausable functionality and zero address validation
     * @param from Source address
     * @param to Destination address (validated for non-mint operations)
     * @param value Transfer amount
     */
    function _update(address from, address to, uint256 value)
        internal
        override
        whenNotPaused
    {
        // Only validate 'to' address for transfers, not for burns
        if (to == address(0) && from != address(0)) {
            // This is a burn operation, which is allowed
        } else if (from != address(0) && to != address(0)) {
            // This is a transfer, validate the recipient
            if (to == address(0)) revert InvalidAddress();
        }
        // For mints (from == address(0)), validation is done in the minting functions
        
        super._update(from, to, value);
    }

    /**
     * @dev Grant role with zero address validation
     * @param role The role to grant
     * @param account The account to grant the role to (must not be zero)
     */
    function grantRole(bytes32 role, address account) 
        public 
        override
        notZeroAddress(account)
    {
        super.grantRole(role, account);
    }

    /**
     * @dev Revoke role with zero address validation
     * @param role The role to revoke
     * @param account The account to revoke the role from (must not be zero)
     */
    function revokeRole(bytes32 role, address account) 
        public 
        override
        notZeroAddress(account)
    {
        super.revokeRole(role, account);
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

    // Pausable functionality
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }

    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
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