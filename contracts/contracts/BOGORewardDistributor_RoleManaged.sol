// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./base/RoleManaged.sol";
import "./interfaces/IRoleManager.sol";

/**
 * @title BOGORewardDistributor_RoleManaged
 * @dev Reward distributor with centralized role management
 */
contract BOGORewardDistributor_RoleManaged is RoleManaged, Pausable, ReentrancyGuard {
    IERC20 public immutable bogoToken;
    
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
    event TreasurySweep(address indexed token, address indexed to, uint256 amount);
    
    modifier onlyTreasury() {
        require(hasRole(roleManager.TREASURY_ROLE(), msg.sender), "Only treasury can call this function");
        _;
    }
    
    modifier onlyAuthorizedBackend() {
        require(hasRole(roleManager.DISTRIBUTOR_BACKEND_ROLE(), msg.sender), "Not authorized backend");
        _;
    }
    
    constructor(
        address _roleManager,
        address _bogoToken
    ) RoleManaged(_roleManager) {
        require(_bogoToken != address(0), "Invalid token address");
        bogoToken = IERC20(_bogoToken);
        lastResetTime = block.timestamp;
        _initializeTemplates();
    }
    
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
        require(bogoToken.transfer(msg.sender, template.fixedAmount), "Transfer failed");
        
        emit RewardClaimed(msg.sender, templateId, template.fixedAmount);
    }
    
    function claimCustomReward(address recipient, uint256 amount, string memory reason) 
        external onlyAuthorizedBackend nonReentrant whenNotPaused {
        _resetDailyLimit();
        
        RewardTemplate memory template = templates["custom_reward"];
        require(template.active, "Custom rewards not active");
        require(amount > 0 && amount <= template.maxAmount, "Invalid amount");
        require(dailyDistributed + amount <= DAILY_GLOBAL_LIMIT, "Daily limit exceeded");
        
        dailyDistributed += amount;
        
        require(bogoToken.transfer(recipient, amount), "Transfer failed");
        
        emit RewardClaimed(recipient, reason, amount);
    }
    
    function claimReferralBonus(address referrer) external nonReentrant whenNotPaused {
        require(referredBy[msg.sender] == address(0), "Already referred");
        require(referrer != msg.sender, "Cannot refer yourself");
        require(referrer != address(0), "Invalid referrer");
        
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
        require(bogoToken.transfer(referrer, template.fixedAmount), "Transfer failed");
        
        emit ReferralClaimed(referrer, msg.sender, template.fixedAmount);
    }
    
    // Admin functions
    function addToWhitelist(address[] memory wallets) external onlyTreasury {
        for (uint i = 0; i < wallets.length; i++) {
            founderWhitelist[wallets[i]] = true;
            emit WhitelistUpdated(wallets[i], true);
        }
    }
    
    function removeFromWhitelist(address wallet) external onlyTreasury {
        founderWhitelist[wallet] = false;
        emit WhitelistUpdated(wallet, false);
    }
    
    function updateTemplate(string memory templateId, RewardTemplate memory newTemplate) 
        external onlyTreasury {
        templates[templateId] = newTemplate;
        emit TemplateUpdated(templateId);
    }
    
    function pause() external onlyRole(roleManager.PAUSER_ROLE()) {
        _pause();
    }
    
    function unpause() external onlyRole(roleManager.PAUSER_ROLE()) {
        _unpause();
    }
    
    // View functions
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
    
    function getRemainingDailyLimit() external view returns (uint256) {
        if (block.timestamp >= lastResetTime + 1 days) {
            return DAILY_GLOBAL_LIMIT;
        }
        return DAILY_GLOBAL_LIMIT - dailyDistributed;
    }
    
    // Treasury sweep function
    function treasurySweep(address token, address to, uint256 amount) external onlyTreasury nonReentrant {
        require(to != address(0), "Invalid recipient");
        require(amount > 0, "Invalid amount");
        
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