// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../base/RoleManaged.sol";
import "../interfaces/IRoleManager.sol";

/**
 * @title BOGOToken
 * @author BOGOWI Team
 * @notice Production BOGO token with external role management and comprehensive security
 * @dev Implements ERC20 with:
 * - External RoleManager for unified access control
 * - Zero address validation for safety
 * - Role-based minting allocations (DAO, Business, Rewards)
 * - Pausable transfers for emergency situations
 * - Burn functionality for deflationary mechanics
 * - Reentrancy protection
 * @custom:security-contact hugo@kode.zone
 */
contract BOGOToken is ERC20, RoleManaged, Pausable, ReentrancyGuard {
    // Supply and allocation constants
    uint256 public constant MAX_SUPPLY = 1_000_000_000 * 10**18; // 1 billion tokens
    uint256 public constant DAO_ALLOCATION = 50_000_000 * 10**18; // 50M for DAO (5% of total)
    uint256 public constant BUSINESS_ALLOCATION = 900_000_000 * 10**18; // 900M for business (90% of total)
    uint256 public constant REWARDS_ALLOCATION = 50_000_000 * 10**18; // 50M for rewards (5% of total)

    // Allocation tracking
    uint256 public daoMinted;
    uint256 public businessMinted;
    uint256 public rewardsMinted;

    // Events
    event AllocationMinted(string indexed allocationType, uint256 amount, address indexed recipient);

    // Custom errors for gas optimization
    error InvalidAddress();
    error InvalidAmount();
    error ExceedsAllocation();
    error ExceedsMaxSupply();
    error InsufficientRole();

    /**
     * @notice Initializes the BOGO token
     * @param _roleManager Address of the RoleManager contract
     * @param _name Token name
     * @param _symbol Token symbol
     */
    constructor(
        address _roleManager,
        string memory _name,
        string memory _symbol
    ) ERC20(_name, _symbol) RoleManaged(_roleManager) {
        // RoleManaged handles roleManager validation
        // No roles granted here - managed by RoleManager
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
     * @notice Mints tokens from the DAO allocation
     * @dev Requires DAO_ROLE, enforces allocation and supply limits
     * @param to Address to receive the minted tokens
     * @param amount Amount of tokens to mint (with 18 decimals)
     */
    function mintFromDAO(address to, uint256 amount) 
        external 
        onlyRole(roleManager.DAO_ROLE()) 
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
     * @notice Mints tokens from the Business allocation
     * @dev Requires BUSINESS_ROLE, enforces allocation and supply limits
     * @param to Address to receive the minted tokens
     * @param amount Amount of tokens to mint (with 18 decimals)
     */
    function mintFromBusiness(address to, uint256 amount) 
        external 
        onlyRole(roleManager.BUSINESS_ROLE()) 
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
     * @notice Mints tokens from the Rewards allocation
     * @dev Requires either DAO_ROLE or BUSINESS_ROLE, enforces allocation and supply limits
     * @param to Address to receive the minted tokens
     * @param amount Amount of tokens to mint (with 18 decimals)
     */
    function mintFromRewards(address to, uint256 amount) 
        external 
        nonReentrant
        notZeroAddress(to)
        notZeroAmount(amount)
    {
        if (!roleManager.checkRole(roleManager.DAO_ROLE(), msg.sender) && 
            !roleManager.checkRole(roleManager.BUSINESS_ROLE(), msg.sender)) {
            revert InsufficientRole();
        }
        if (rewardsMinted + amount > REWARDS_ALLOCATION) revert ExceedsAllocation();
        if (totalSupply() + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        
        rewardsMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("Rewards", amount, to);
    }

    /**
     * @notice Burns tokens from the caller's balance
     * @dev Reduces total supply, cannot be reversed
     * @param amount Amount of tokens to burn (with 18 decimals)
     */
    function burn(uint256 amount) 
        external
        notZeroAmount(amount)
    {
        _burn(msg.sender, amount);
    }

    /**
     * @notice Burns tokens from a specified account
     * @dev Requires approval from the account, reduces total supply
     * @param account Address to burn tokens from
     * @param amount Amount of tokens to burn (with 18 decimals)
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
     * @notice Returns the remaining unminted DAO allocation
     * @return Remaining DAO allocation in tokens (with 18 decimals)
     */
    function getRemainingDAOAllocation() external view returns (uint256) {
        return DAO_ALLOCATION - daoMinted;
    }

    /**
     * @notice Returns the remaining unminted Business allocation
     * @return Remaining Business allocation in tokens (with 18 decimals)
     */
    function getRemainingBusinessAllocation() external view returns (uint256) {
        return BUSINESS_ALLOCATION - businessMinted;
    }

    /**
     * @notice Returns the remaining unminted Rewards allocation
     * @return Remaining Rewards allocation in tokens (with 18 decimals)
     */
    function getRemainingRewardsAllocation() external view returns (uint256) {
        return REWARDS_ALLOCATION - rewardsMinted;
    }

    /**
     * @notice Pauses all token transfers
     * @dev Requires PAUSER_ROLE from RoleManager
     */
    function pause() external onlyRole(roleManager.PAUSER_ROLE()) {
        _pause();
    }

    /**
     * @notice Unpauses token transfers
     * @dev Requires PAUSER_ROLE from RoleManager
     */
    function unpause() external onlyRole(roleManager.PAUSER_ROLE()) {
        _unpause();
    }

    /**
     * @notice Internal function to update balances
     * @dev Overrides ERC20 to add pause functionality and zero address validation
     * @param from Address sending tokens
     * @param to Address receiving tokens
     * @param value Amount of tokens to transfer
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
     * @notice Grant role with zero address validation
     * @dev This function is not used since roles are managed by RoleManager
     * @param account The account to grant the role to
     */
    function grantRole(bytes32, address account) 
        public 
        pure
        notZeroAddress(account)
    {
        // Roles are managed by RoleManager, not here
        revert("Use RoleManager to manage roles");
    }

    /**
     * @notice Revoke role with zero address validation
     * @dev This function is not used since roles are managed by RoleManager
     * @param account The account to revoke the role from
     */
    function revokeRole(bytes32, address account) 
        public 
        pure
        notZeroAddress(account)
    {
        // Roles are managed by RoleManager, not here
        revert("Use RoleManager to manage roles");
    }
}