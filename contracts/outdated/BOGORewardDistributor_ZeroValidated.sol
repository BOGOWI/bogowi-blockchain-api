// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title BOGORewardDistributor_ZeroValidated
 * @dev Reward distributor with comprehensive zero address validation
 */
contract BOGORewardDistributor_ZeroValidated is Pausable, ReentrancyGuard {
    IERC20 public immutable bogoToken;
    address public immutable treasury;
    
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
    mapping(address => uint256) public referralDepth;
    uint256 public constant MAX_REFERRAL_DEPTH = 10;
    
    // Daily limits
    uint256 public constant DAILY_GLOBAL_LIMIT = 500000 * 10**18;
    uint256 public dailyDistributed;
    uint256 public lastResetTime;
    
    // Events
    event RewardClaimed(address indexed wallet, string templateId, uint256 amount);
    event ReferralClaimed(address indexed referrer, address indexed referred, uint256 amount);
    event TemplateUpdated(string templateId);
    event WhitelistUpdated(address indexed wallet, bool status);
    event BackendAuthorized(address indexed backend, bool authorized);
    event TreasurySweep(address indexed token, address indexed to, uint256 amount);
    
    // Custom errors for gas optimization
    error InvalidAddress();
    error InvalidTreasuryAddress();
    error InvalidTokenAddress();
    error NotAuthorizedBackend();
    error NotTreasury();
    error InvalidAmount();
    error InvalidRecipient();
    error TransferFailed();
    
    modifier onlyAuthorized() {
        if (!authorizedBackends[msg.sender]) revert NotAuthorizedBackend();
        _;
    }
    
    modifier onlyTreasury() {
        if (msg.sender != treasury) revert NotTreasury();
        _;
    }
    
    modifier notZeroAddress(address addr) {
        if (addr == address(0)) revert InvalidAddress();
        _;
    }
    
    constructor(address _bogoToken, address _treasury) {
        if (_bogoToken == address(0)) revert InvalidTokenAddress();
        if (_treasury == address(0)) revert InvalidTreasuryAddress();
        
        bogoToken = IERC20(_bogoToken);
        treasury = _treasury;
        lastResetTime = block.timestamp;
        _initializeTemplates();
    }
    
    function _initializeTemplates() private {
        // Initialize all reward templates
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
        
        templates["referral_bonus"] = RewardTemplate({
            id: "referral_bonus",
            fixedAmount: 20 * 10**18,
            maxAmount: 0,
            cooldownPeriod: 0,
            maxClaimsPerWallet: 0,
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
            maxClaimsPerWallet: 0,
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
    
    function _resetDailyLimit() private {
        if (block.timestamp >= lastResetTime + 1 days) {
            dailyDistributed = 0;
            lastResetTime = block.timestamp;
        }
    }
    
    function claimReward(string memory templateId) external nonReentrant whenNotPaused {
        _resetDailyLimit();
        
        RewardTemplate memory template = templates[templateId];
        require(template.active, "Template not active");
        require(template.fixedAmount > 0, "Use claimCustomReward for custom amounts");
        
        // Check eligibility
        if (template.maxClaimsPerWallet > 0) {
            require(claimCount[msg.sender][templateId] < template.maxClaimsPerWallet, 
                    "Max claims reached");
        }
        
        if (template.cooldownPeriod > 0) {
            require(block.timestamp >= lastClaim[msg.sender][templateId] + template.cooldownPeriod,
                    "Cooldown period active");
        }
        
        if (template.requiresWhitelist) {
            require(founderWhitelist[msg.sender], "Not whitelisted");
        }
        
        // Check daily limit
        require(dailyDistributed + template.fixedAmount <= DAILY_GLOBAL_LIMIT, 
                "Daily limit exceeded");
        
        // Update state
        claimCount[msg.sender][templateId]++;
        lastClaim[msg.sender][templateId] = block.timestamp;
        dailyDistributed += template.fixedAmount;
        
        // Transfer tokens
        if (!bogoToken.transfer(msg.sender, template.fixedAmount)) revert TransferFailed();
        
        emit RewardClaimed(msg.sender, templateId, template.fixedAmount);
    }
    
    /**
     * @dev Claim custom reward with zero address validation
     * @param recipient The recipient address (must not be zero)
     * @param amount The reward amount (must be greater than zero)
     * @param reason The reason for the reward
     */
    function claimCustomReward(address recipient, uint256 amount, string memory reason) 
        external 
        onlyAuthorized 
        nonReentrant 
        whenNotPaused
        notZeroAddress(recipient)
    {
        if (amount == 0) revert InvalidAmount();
        _resetDailyLimit();
        
        RewardTemplate memory template = templates["custom_reward"];
        require(template.active, "Custom rewards not active");
        require(amount <= template.maxAmount, "Amount exceeds maximum");
        require(dailyDistributed + amount <= DAILY_GLOBAL_LIMIT, "Daily limit exceeded");
        
        dailyDistributed += amount;
        
        if (!bogoToken.transfer(recipient, amount)) revert TransferFailed();
        
        emit RewardClaimed(recipient, reason, amount);
    }
    
    /**
     * @dev Claim referral bonus with zero address validation
     * @param referrer The referrer address (must not be zero)
     */
    function claimReferralBonus(address referrer) 
        external 
        nonReentrant 
        whenNotPaused
        notZeroAddress(referrer)
    {
        require(referredBy[msg.sender] == address(0), "Already referred");
        require(referrer != msg.sender, "Cannot refer yourself");
        
        // Check for circular referrals
        require(!_hasCircularReferral(msg.sender, referrer), "Circular referral detected");
        
        // Check referral depth
        uint256 referrerDepth = referralDepth[referrer];
        require(referrerDepth < MAX_REFERRAL_DEPTH, "Max referral depth exceeded");
        
        RewardTemplate memory template = templates["referral_bonus"];
        require(template.active, "Referral rewards not active");
        
        _resetDailyLimit();
        require(dailyDistributed + template.fixedAmount <= DAILY_GLOBAL_LIMIT, 
                "Daily limit exceeded");
        
        // Update state
        referredBy[msg.sender] = referrer;
        referralCount[referrer]++;
        referralDepth[msg.sender] = referrerDepth + 1;
        dailyDistributed += template.fixedAmount;
        
        // Transfer to referrer
        if (!bogoToken.transfer(referrer, template.fixedAmount)) revert TransferFailed();
        
        emit ReferralClaimed(referrer, msg.sender, template.fixedAmount);
    }
    
    /**
     * @dev Add addresses to whitelist with zero address validation
     * @param wallets Array of wallet addresses to whitelist
     */
    function addToWhitelist(address[] memory wallets) external onlyTreasury {
        for (uint i = 0; i < wallets.length; i++) {
            if (wallets[i] == address(0)) revert InvalidAddress();
            founderWhitelist[wallets[i]] = true;
            emit WhitelistUpdated(wallets[i], true);
        }
    }
    
    /**
     * @dev Remove address from whitelist with zero address validation
     * @param wallet The wallet address to remove from whitelist
     */
    function removeFromWhitelist(address wallet) 
        external 
        onlyTreasury
        notZeroAddress(wallet)
    {
        founderWhitelist[wallet] = false;
        emit WhitelistUpdated(wallet, false);
    }
    
    /**
     * @dev Set authorized backend with zero address validation
     * @param backend The backend address
     * @param authorized Whether the backend is authorized
     */
    function setAuthorizedBackend(address backend, bool authorized) 
        external 
        onlyTreasury
        notZeroAddress(backend)
    {
        authorizedBackends[backend] = authorized;
        emit BackendAuthorized(backend, authorized);
    }
    
    function updateTemplate(string memory templateId, RewardTemplate memory newTemplate) 
        external onlyTreasury {
        templates[templateId] = newTemplate;
        emit TemplateUpdated(templateId);
    }
    
    function pause() external onlyTreasury {
        _pause();
    }
    
    function unpause() external onlyTreasury {
        _unpause();
    }
    
    // View functions
    function canClaim(address wallet, string memory templateId) external view returns (bool, string memory) {
        if (wallet == address(0)) return (false, "Invalid wallet address");
        
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
    
    function getRemainingDailyLimit() external view returns (uint256) {
        if (block.timestamp >= lastResetTime + 1 days) {
            return DAILY_GLOBAL_LIMIT;
        }
        return DAILY_GLOBAL_LIMIT - dailyDistributed;
    }
    
    /**
     * @dev Treasury sweep with zero address validation
     * @param token The token address (0 for ETH)
     * @param to The recipient address (must not be zero)
     * @param amount The amount to sweep
     */
    function treasurySweep(address token, address to, uint256 amount) 
        external 
        onlyTreasury 
        nonReentrant
        notZeroAddress(to)
    {
        if (amount == 0) revert InvalidAmount();
        
        if (token == address(0)) {
            // Withdraw ETH
            (bool success, ) = to.call{value: amount}("");
            require(success, "ETH transfer failed");
        } else {
            // Withdraw ERC20 tokens
            IERC20(token).transfer(to, amount);
        }
        
        emit TreasurySweep(token, to, amount);
    }
    
    function _hasCircularReferral(address newUser, address referrer) private view returns (bool) {
        address current = referrer;
        uint256 depth = 0;
        
        while (current != address(0) && depth < MAX_REFERRAL_DEPTH) {
            if (current == newUser) {
                return true;
            }
            current = referredBy[current];
            depth++;
        }
        
        return false;
    }
    
    function getReferralChain(address user) external view returns (address[] memory) {
        if (user == address(0)) {
            return new address[](0);
        }
        
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
        
        address[] memory chain = new address[](count);
        for (uint256 i = 0; i < count; i++) {
            chain[i] = tempChain[i];
        }
        
        return chain;
    }
}