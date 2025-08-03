// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./StandardErrors.sol";

/**
 * @title BOGORewardDistributor
 * @author BOGOWI Team
 * @notice Manages BOGO token reward distribution for various engagement activities
 * @dev Implements templated rewards with cooldowns, whitelists, and daily limits
 * Features:
 * - Pre-defined reward templates (welcome, founder, referral, etc.)
 * - Referral system with circular reference prevention
 * - Daily global distribution limits
 * - Whitelist support for exclusive rewards
 * - Backend authorization for custom rewards
 * @custom:security-contact security@bogowi.com
 */
contract BOGORewardDistributor is Pausable, ReentrancyGuard, StandardErrors {
    IERC20 public immutable bogoToken;
    address public immutable treasury;
    bool public immutable IS_TEST_MODE;
    
    /**
     * @dev Reward template structure
     * @param id Unique identifier for the template
     * @param fixedAmount Fixed reward amount (0 for custom rewards)
     * @param maxAmount Maximum amount for custom rewards
     * @param cooldownPeriod Time between claims in seconds
     * @param maxClaimsPerWallet Maximum claims per wallet (0 = unlimited)
     * @param requiresWhitelist Whether whitelist is required
     * @param active Whether template is currently active
     */
    struct RewardTemplate {
        string id;
        uint256 fixedAmount;
        uint256 maxAmount;
        uint256 cooldownPeriod;
        uint256 maxClaimsPerWallet;
        bool requiresWhitelist;
        bool active;
    }
    
    mapping(string => RewardTemplate) public templates;
    mapping(address => mapping(string => uint256)) public lastClaim;
    mapping(address => mapping(string => uint256)) public claimCount;
    mapping(address => bool) public founderWhitelist;
    mapping(address => bool) public authorizedBackends;
    
    // Referral tracking
    mapping(address => address) public referredBy;
    mapping(address => uint256) public referralCount;
    mapping(address => uint256) public referralDepth; // Track depth in referral chain
    uint256 public constant MAX_REFERRAL_DEPTH = 10; // Maximum depth to prevent long chains
    
    // Daily limits
    uint256 public constant DAILY_GLOBAL_LIMIT = 500000 * 10**18; // 500k BOGO
    uint256 public dailyDistributed;
    uint256 public lastResetTime;
    
    // Events
    event RewardClaimed(address indexed wallet, string templateId, uint256 amount);
    event ReferralClaimed(address indexed referrer, address indexed referred, uint256 amount);
    event TemplateUpdated(string templateId);
    event WhitelistUpdated(address indexed wallet, bool status);
    event AuthorizedBackendSet(address indexed backend, bool authorized);
    event DailyLimitReset(uint256 timestamp, uint256 previousDistributed);
    
    modifier onlyAuthorized() {
        require(authorizedBackends[msg.sender], NOT_BACKEND);
        _;
    }
    
    modifier onlyTreasury() {
        require(msg.sender == treasury, NOT_TREASURY);
        _;
    }
    
    modifier onlyTreasuryOrTest() {
        if (IS_TEST_MODE) {
            require(msg.sender == treasury || authorizedBackends[msg.sender], NOT_BACKEND);
        } else {
            require(msg.sender == treasury, NOT_TREASURY);
        }
        _;
    }
    
    /**
     * @notice Initializes the reward distributor
     * @dev Sets up BOGO token, treasury, and initializes reward templates
     * @param _bogoToken Address of the BOGO token contract
     * @param _treasury Address that funds the rewards (MultisigTreasury)
     * @param _isTestMode Whether to enable test mode for enhanced testability
     */
    constructor(address _bogoToken, address _treasury, bool _isTestMode) {
        require(_bogoToken != address(0), ZERO_ADDRESS);
        require(_treasury != address(0), ZERO_ADDRESS);
        
        bogoToken = IERC20(_bogoToken);
        treasury = _treasury;
        IS_TEST_MODE = _isTestMode;
        
        lastResetTime = block.timestamp;
        _initializeTemplates();
    }
    
    /**
     * @dev Initializes all predefined reward templates
     * @custom:templates
     * - welcome_bonus: 10 BOGO for new users
     * - founder_bonus: 100 BOGO for whitelisted founders
     * - referral_bonus: 20 BOGO per successful referral
     * - first_nft_mint: 25 BOGO for first NFT mint
     * - dao_participation: 15 BOGO with 30-day cooldown
     * - attraction_tiers: 10-50 BOGO based on tier
     * - custom_reward: Variable amount up to 1000 BOGO
     */
    function _initializeTemplates() private {
        // Onboarding rewards
        templates["welcome_bonus"] = RewardTemplate({
            id: "welcome_bonus",
            fixedAmount: 10 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 1,
            requiresWhitelist: false,
            active: true
        });
        
        templates["founder_bonus"] = RewardTemplate({
            id: "founder_bonus",
            fixedAmount: 100 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 1,
            requiresWhitelist: true,
            active: true
        });
        
        // Engagement rewards
        templates["referral_bonus"] = RewardTemplate({
            id: "referral_bonus",
            fixedAmount: 20 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 0, // unlimited
            requiresWhitelist: false,
            active: true
        });
        
        templates["first_nft_mint"] = RewardTemplate({
            id: "first_nft_mint",
            fixedAmount: 25 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 1,
            requiresWhitelist: false,
            active: true
        });
        
        templates["dao_participation"] = RewardTemplate({
            id: "dao_participation",
            fixedAmount: 15 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 30 days,
            maxClaimsPerWallet: 0, // unlimited with cooldown
            requiresWhitelist: false,
            active: true
        });
        
        // Attraction tiers
        templates["attraction_tier_1"] = RewardTemplate({
            id: "attraction_tier_1",
            fixedAmount: 10 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 0,
            requiresWhitelist: false,
            active: true
        });
        
        templates["attraction_tier_2"] = RewardTemplate({
            id: "attraction_tier_2",
            fixedAmount: 20 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 0,
            requiresWhitelist: false,
            active: true
        });
        
        templates["attraction_tier_3"] = RewardTemplate({
            id: "attraction_tier_3",
            fixedAmount: 40 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 0,
            requiresWhitelist: false,
            active: true
        });
        
        templates["attraction_tier_4"] = RewardTemplate({
            id: "attraction_tier_4",
            fixedAmount: 50 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 0,
            requiresWhitelist: false,
            active: true
        });
        
        // Custom reward
        templates["custom_reward"] = RewardTemplate({
            id: "custom_reward",
            fixedAmount: 0,
            maxAmount: 1000 * 10**18,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 0,
            requiresWhitelist: false,
            active: true
        });
    }
    
    /**
     * @dev Resets daily distribution counter if 24 hours have passed
     * @custom:emits DailyLimitReset
     */
    function _resetDailyLimit() private {
        if (block.timestamp >= lastResetTime + 1 days) {
            uint256 previousDistributed = dailyDistributed;
            dailyDistributed = 0;
            lastResetTime = block.timestamp;
            emit DailyLimitReset(block.timestamp, previousDistributed);
        }
    }
    
    /**
     * @notice Claims a fixed reward based on template ID
     * @dev Validates eligibility based on template rules
     * @param templateId ID of the reward template to claim
     * @custom:emits RewardClaimed
     * @custom:security Enforces cooldowns, limits, and whitelist requirements
     */
    function claimReward(string memory templateId) external nonReentrant whenNotPaused {
        _resetDailyLimit();
        
        RewardTemplate memory template = templates[templateId];
        require(template.active, INACTIVE);
        require(template.fixedAmount > 0, INVALID_AMOUNT);
        
        // Check eligibility
        if (template.maxClaimsPerWallet > 0) {
            require(claimCount[msg.sender][templateId] < template.maxClaimsPerWallet, 
                    MAX_REACHED);
        }
        
        if (template.cooldownPeriod > 0) {
            require(block.timestamp >= lastClaim[msg.sender][templateId] + template.cooldownPeriod,
                    COOLDOWN_ACTIVE);
        }
        
        if (template.requiresWhitelist) {
            require(founderWhitelist[msg.sender], NOT_WHITELISTED);
        }
        
        // Check daily limit
        require(dailyDistributed + template.fixedAmount <= DAILY_GLOBAL_LIMIT, 
                DAILY_LIMIT_EXCEEDED);
        
        // Update state
        claimCount[msg.sender][templateId]++;
        lastClaim[msg.sender][templateId] = block.timestamp;
        dailyDistributed += template.fixedAmount;
        
        // Transfer tokens
        require(bogoToken.transfer(msg.sender, template.fixedAmount), TRANSFER_FAILED);
        
        emit RewardClaimed(msg.sender, templateId, template.fixedAmount);
    }
    
    /**
     * @notice Distributes custom reward amounts to specified recipient
     * @dev Only authorized backends can call this function
     * @param recipient Address to receive the reward
     * @param amount Amount of BOGO tokens to distribute
     * @param reason Description of why reward is being given
     * @custom:emits RewardClaimed
     * @custom:security Requires backend authorization
     */
    function claimCustomReward(address recipient, uint256 amount, string memory reason) 
        external onlyAuthorized nonReentrant whenNotPaused {
        _resetDailyLimit();
        
        RewardTemplate memory template = templates["custom_reward"];
        require(template.active, INACTIVE);
        require(amount > 0 && amount <= template.maxAmount, INVALID_AMOUNT);
        require(dailyDistributed + amount <= DAILY_GLOBAL_LIMIT, DAILY_LIMIT_EXCEEDED);
        
        dailyDistributed += amount;
        
        require(bogoToken.transfer(recipient, amount), TRANSFER_FAILED);
        
        emit RewardClaimed(recipient, reason, amount);
    }
    
    /**
     * @notice Claims referral bonus for the referrer
     * @dev Prevents circular referrals and enforces depth limits
     * @param referrer Address of the user who referred the caller
     * @custom:emits ReferralClaimed
     * @custom:security Prevents self-referral and circular chains
     */
    function claimReferralBonus(address referrer) external nonReentrant whenNotPaused {
        require(referredBy[msg.sender] == address(0), ALREADY_EXISTS);
        require(referrer != msg.sender, SELF_REFERENCE);
        require(referrer != address(0), ZERO_ADDRESS);
        
        // Check for circular referrals
        require(!_hasCircularReferral(msg.sender, referrer), CIRCULAR_REFERENCE);
        
        // Check referral depth
        uint256 referrerDepth = referralDepth[referrer];
        require(referrerDepth < MAX_REFERRAL_DEPTH, EXCEEDS_LIMIT);
        
        RewardTemplate memory template = templates["referral_bonus"];
        require(template.active, INACTIVE);
        
        _resetDailyLimit();
        require(dailyDistributed + template.fixedAmount <= DAILY_GLOBAL_LIMIT, 
                DAILY_LIMIT_EXCEEDED);
        
        // Update state
        referredBy[msg.sender] = referrer;
        referralCount[referrer]++;
        referralDepth[msg.sender] = referrerDepth + 1;
        dailyDistributed += template.fixedAmount;
        
        // Transfer to referrer
        require(bogoToken.transfer(referrer, template.fixedAmount), TRANSFER_FAILED);
        
        emit ReferralClaimed(referrer, msg.sender, template.fixedAmount);
    }
    
    /**
     * @notice Adds multiple addresses to founder whitelist
     * @dev Only treasury can manage whitelist
     * @param wallets Array of addresses to whitelist
     * @custom:emits WhitelistUpdated
     */
    function addToWhitelist(address[] memory wallets) external onlyTreasury {
        for (uint i = 0; i < wallets.length; i++) {
            founderWhitelist[wallets[i]] = true;
            emit WhitelistUpdated(wallets[i], true);
        }
    }
    
    /**
     * @notice Removes an address from founder whitelist
     * @param wallet Address to remove from whitelist
     * @custom:emits WhitelistUpdated
     */
    function removeFromWhitelist(address wallet) external onlyTreasury {
        founderWhitelist[wallet] = false;
        emit WhitelistUpdated(wallet, false);
    }
    
    /**
     * @notice Sets backend authorization status
     * @dev Authorized backends can distribute custom rewards
     * @param backend Address of the backend system
     * @param authorized Whether to authorize or revoke
     * @custom:emits AuthorizedBackendSet
     */
    function setAuthorizedBackend(address backend, bool authorized) external onlyTreasury {
        authorizedBackends[backend] = authorized;
        emit AuthorizedBackendSet(backend, authorized);
    }
    
    /**
     * @notice Updates an existing reward template
     * @param templateId ID of template to update
     * @param newTemplate New template configuration
     * @custom:emits TemplateUpdated
     */
    function updateTemplate(string memory templateId, RewardTemplate memory newTemplate) 
        external onlyTreasury {
        templates[templateId] = newTemplate;
        emit TemplateUpdated(templateId);
    }
    
    /**
     * @notice Pauses all reward claims
     * @dev Emergency function to halt distributions
     */
    function pause() external onlyTreasuryOrTest {
        _pause();
    }

    /**
     * @notice Unpauses reward claims
     */
    function unpause() external onlyTreasuryOrTest {
        _unpause();
    }
    
    /**
     * @notice Checks if a wallet can claim a specific reward
     * @param wallet Address to check eligibility for
     * @param templateId ID of the reward template
     * @return eligible Whether the wallet can claim
     * @return reason Human-readable reason if not eligible
     */
    function canClaim(address wallet, string memory templateId) external view returns (bool, string memory) {
        RewardTemplate memory template = templates[templateId];
        
        if (!template.active) return (false, "Template not active");
        
        if (template.maxClaimsPerWallet > 0 && 
            claimCount[wallet][templateId] >= template.maxClaimsPerWallet) {
            return (false, "Max claims reached");
        }
        
        if (template.cooldownPeriod > 0 && 
            block.timestamp < lastClaim[wallet][templateId] + template.cooldownPeriod) {
            return (false, "Cooldown period active");
        }
        
        if (template.requiresWhitelist && !founderWhitelist[wallet]) {
            return (false, "Not whitelisted");
        }
        
        return (true, "Eligible");
    }
    
    /**
     * @notice Returns remaining tokens in daily distribution limit
     * @return Amount of BOGO tokens remaining for today
     */
    function getRemainingDailyLimit() external view returns (uint256) {
        if (block.timestamp >= lastResetTime + 1 days) {
            return DAILY_GLOBAL_LIMIT;
        }
        return DAILY_GLOBAL_LIMIT - dailyDistributed;
    }
    
    /**
     * @notice Emergency function to recover tokens
     * @dev Used for migrations or recovering stuck tokens
     * @param token Token address (0x0 for ETH)
     * @param to Recipient address
     * @param amount Amount to transfer
     * @custom:emits TreasurySweep
     * @custom:security Only treasury can execute
     */
    function treasurySweep(address token, address to, uint256 amount) external onlyTreasury nonReentrant {
        require(to != address(0), ZERO_ADDRESS);
        require(amount > 0, ZERO_AMOUNT);
        
        if (token == address(0)) {
            // Withdraw ETH
            (bool success, ) = to.call{value: amount}("");
            require(success, TRANSFER_FAILED);
        } else {
            // Withdraw ERC20 tokens
            IERC20(token).transfer(to, amount);
        }
        
        emit TreasurySweep(token, to, amount);
    }
    
    // Event for treasury sweep operations
    event TreasurySweep(address indexed token, address indexed to, uint256 amount);
    
    /**
     * @dev Checks if adding a referral would create a circular reference
     * @param newUser The user being referred
     * @param referrer The proposed referrer
     * @return True if circular reference detected
     */
    function _hasCircularReferral(address newUser, address referrer) private view returns (bool) {
        address current = referrer;
        uint256 depth = 0;
        
        // Traverse up the referral chain
        while (current != address(0) && depth < MAX_REFERRAL_DEPTH) {
            if (current == newUser) {
                return true; // Circular reference found
            }
            current = referredBy[current];
            depth++;
        }
        
        return false;
    }
    
    /**
     * @notice Gets the complete referral chain for a user
     * @dev Returns all referrers up to MAX_REFERRAL_DEPTH
     * @param user The user to check referral chain for
     * @return Array of referrer addresses from immediate to root
     */
    function getReferralChain(address user) external view returns (address[] memory) {
        address[] memory tempChain = new address[](MAX_REFERRAL_DEPTH);
        address current = user;
        uint256 count = 0;
        
        while (current != address(0) && count < MAX_REFERRAL_DEPTH) {
            current = referredBy[current];
            if (current != address(0)) {
                tempChain[count] = current;
                count++;
            }
        }
        
        // Create properly sized array
        address[] memory chain = new address[](count);
        for (uint256 i = 0; i < count; i++) {
            chain[i] = tempChain[i];
        }
        
        return chain;
    }
    
    // =============================================================================
    // TEST MODE FUNCTIONS
    // =============================================================================
    
    /**
     * @notice Test-only function to directly set pause state
     * @dev Only available in test mode for comprehensive testing
     * @param _paused Whether to pause or unpause the contract
     */
    function testSetPaused(bool _paused) external {
        require(IS_TEST_MODE, "Only in test mode");
        if (_paused) {
            _pause();
        } else {
            _unpause();
        }
    }
    
    /**
     * @notice Test-only function to reset daily distribution limit
     * @dev Allows testing daily limit scenarios without waiting 24 hours
     */
    function testResetDailyLimit() external {
        require(IS_TEST_MODE, "Only in test mode");
        dailyDistributed = 0;
        lastResetTime = block.timestamp;
    }
    
    /**
     * @notice Test-only function to set daily distributed amount
     * @dev Allows testing daily limit edge cases
     * @param _amount Amount to set as already distributed today
     */
    function testSetDailyDistributed(uint256 _amount) external {
        require(IS_TEST_MODE, "Only in test mode");
        require(_amount <= DAILY_GLOBAL_LIMIT, "Amount exceeds daily limit");
        dailyDistributed = _amount;
    }
    
    /**
     * @notice Test-only function to manipulate time for testing
     * @dev Allows testing time-dependent functionality
     * @param _timestamp New timestamp to set for last reset
     */
    function testSetLastResetTime(uint256 _timestamp) external {
        require(IS_TEST_MODE, "Only in test mode");
        lastResetTime = _timestamp;
    }
    
    /**
     * @notice Test-only function to directly set claim count
     * @dev Allows testing max claims scenarios
     * @param _user User address
     * @param _templateId Template ID
     * @param _count Claim count to set
     */
    function testSetClaimCount(address _user, string memory _templateId, uint256 _count) external {
        require(IS_TEST_MODE, "Only in test mode");
        claimCount[_user][_templateId] = _count;
    }
    
    /**
     * @notice Test-only function to directly set last claim time
     * @dev Allows testing cooldown scenarios
     * @param _user User address
     * @param _templateId Template ID
     * @param _timestamp Last claim timestamp to set
     */
    function testSetLastClaim(address _user, string memory _templateId, uint256 _timestamp) external {
        require(IS_TEST_MODE, "Only in test mode");
        lastClaim[_user][_templateId] = _timestamp;
    }
}