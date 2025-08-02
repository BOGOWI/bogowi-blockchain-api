// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./StandardErrors.sol";

/**
 * @title BOGOTokenV2
 * @author BOGOWI Team
 * @notice Enhanced BOGO token with role-based access control, supply management, and timelock governance
 * @dev Implements ERC20 with additional features:
 * - Role-based minting allocations (DAO, Business, Rewards)
 * - Timelock mechanism for critical operations
 * - Flavored token registration system
 * - Pausable transfers for emergency situations
 * - Burn functionality for deflationary mechanics
 * @custom:security-contact security@bogowi.com
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

    /**
     * @notice Initializes the BOGO token with name "BOGOWI" and symbol "BOGO"
     * @dev Grants DEFAULT_ADMIN_ROLE, MINTER_ROLE, and PAUSER_ROLE to deployer
     */
    constructor() ERC20("BOGOWI", "BOGO") {
        // Grant admin roles to deployer
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
    }

    /**
     * @notice Mints tokens from the DAO allocation
     * @dev Requires DAO_ROLE, enforces allocation and supply limits
     * @param to Address to receive the minted tokens
     * @param amount Amount of tokens to mint (with 18 decimals)
     * @custom:emits AllocationMinted
     */
    function mintFromDAO(address to, uint256 amount) external onlyRole(DAO_ROLE) nonReentrant {
        require(daoMinted + amount <= DAO_ALLOCATION, EXCEEDS_ALLOCATION);
        require(totalSupply() + amount <= MAX_SUPPLY, EXCEEDS_SUPPLY);
        
        daoMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("DAO", amount, to);
    }

    /**
     * @notice Mints tokens from the Business allocation
     * @dev Requires BUSINESS_ROLE, enforces allocation and supply limits
     * @param to Address to receive the minted tokens
     * @param amount Amount of tokens to mint (with 18 decimals)
     * @custom:emits AllocationMinted
     */
    function mintFromBusiness(address to, uint256 amount) external onlyRole(BUSINESS_ROLE) nonReentrant {
        require(businessMinted + amount <= BUSINESS_ALLOCATION, EXCEEDS_ALLOCATION);
        require(totalSupply() + amount <= MAX_SUPPLY, EXCEEDS_SUPPLY);
        
        businessMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("Business", amount, to);
    }

    /**
     * @notice Mints tokens from the Rewards allocation
     * @dev Requires either DAO_ROLE or BUSINESS_ROLE, enforces allocation and supply limits
     * @param to Address to receive the minted tokens
     * @param amount Amount of tokens to mint (with 18 decimals)
     * @custom:emits AllocationMinted
     */
    function mintFromRewards(address to, uint256 amount) external nonReentrant {
        require(hasRole(DAO_ROLE, msg.sender) || hasRole(BUSINESS_ROLE, msg.sender), 
                UNAUTHORIZED);
        require(rewardsMinted + amount <= REWARDS_ALLOCATION, EXCEEDS_ALLOCATION);
        require(totalSupply() + amount <= MAX_SUPPLY, EXCEEDS_SUPPLY);
        
        rewardsMinted += amount;
        _mint(to, amount);
        
        emit AllocationMinted("Rewards", amount, to);
    }

    /**
     * @notice Queues a flavored token registration for timelock execution
     * @dev Requires DEFAULT_ADMIN_ROLE, validates token address is a contract
     * @param flavor Unique identifier for the flavored token
     * @param tokenAddress Contract address of the flavored token
     * @custom:emits TimelockQueued
     */
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

    /**
     * @notice Executes a previously queued flavored token registration
     * @dev Requires DEFAULT_ADMIN_ROLE and timelock period to have passed
     * @param flavor Unique identifier for the flavored token (must match queued operation)
     * @param tokenAddress Contract address of the flavored token (must match queued operation)
     * @custom:emits FlavoredTokenRegistered, TimelockExecuted
     */
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

    /**
     * @notice Cancels a queued timelock operation
     * @dev Requires DEFAULT_ADMIN_ROLE
     * @param operationId Keccak256 hash of the operation to cancel
     * @custom:emits TimelockCancelled
     */
    function cancelTimelockOperation(bytes32 operationId) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(timelockOperations[operationId] != 0, NOT_INITIALIZED);
        
        delete timelockOperations[operationId];
        emit TimelockCancelled(operationId);
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
     * @notice Gas-optimized lookup of flavored token by hash
     * @dev L2 optimization: Uses hash instead of string for cheaper lookups
     * @param flavorHash Keccak256 hash of the flavor string
     * @return Address of the flavored token contract
     */
    function getFlavoredTokenByHash(bytes32 flavorHash) external view returns (address) {
        return flavoredTokensByHash[flavorHash];
    }

    /**
     * @notice Pauses all token transfers
     * @dev Requires PAUSER_ROLE, affects transfer(), transferFrom(), mint(), and burn()
     * @custom:security Emergency function
     */
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }

    /**
     * @notice Unpauses token transfers
     * @dev Requires PAUSER_ROLE
     */
    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
    }

    /**
     * @notice Burns tokens from the caller's balance
     * @dev Reduces total supply, cannot be reversed
     * @param amount Amount of tokens to burn (with 18 decimals)
     */
    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }

    /**
     * @notice Burns tokens from a specified account
     * @dev Requires approval from the account, reduces total supply
     * @param account Address to burn tokens from
     * @param amount Amount of tokens to burn (with 18 decimals)
     */
    function burnFrom(address account, uint256 amount) external {
        _spendAllowance(account, msg.sender, amount);
        _burn(account, amount);
    }

    /**
     * @notice Internal function to update balances
     * @dev Overrides ERC20 to add pause functionality
     * @param from Address sending tokens
     * @param to Address receiving tokens
     * @param value Amount of tokens to transfer
     */
    function _update(address from, address to, uint256 value)
        internal
        override
        whenNotPaused
    {
        super._update(from, to, value);
    }

    /**
     * @notice Checks if an address is a contract
     * @dev Uses extcodesize opcode for verification
     * @param account Address to check
     * @return bool True if address contains contract code
     */
    function _isContract(address account) internal view returns (bool) {
        uint256 size;
        assembly {
            size := extcodesize(account)
        }
        return size > 0;
    }

    /**
     * @notice Query if a contract implements an interface
     * @dev Required by AccessControl inheritance
     * @param interfaceId The interface identifier, as specified in ERC-165
     * @return bool True if the contract implements interfaceId
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
