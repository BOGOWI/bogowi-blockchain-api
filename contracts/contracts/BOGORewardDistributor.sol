// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract BOGORewardDistributor is Pausable, ReentrancyGuard {
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
    
    // Daily limits
    uint256 public constant DAILY_GLOBAL_LIMIT = 500000 * 10**18; // 500k BOGO
    uint256 public dailyDistributed;
    uint256 public lastResetTime;
    
    // Events
    event RewardClaimed(address indexed wallet, string templateId, uint256 amount);
    event ReferralClaimed(address indexed referrer, address indexed referred, uint256 amount);
    event TemplateUpdated(string templateId);
    event WhitelistUpdated(address indexed wallet, bool status);
    
    modifier onlyAuthorized() {
        require(authorizedBackends[msg.sender], "Not authorized backend");
        _;
    }
    
    modifier onlyTreasury() {
        require(msg.sender == treasury, "Only treasury can call this function");
        _;
    }
    
    constructor(address _bogoToken, address _treasury) {
        require(_treasury != address(0), "Invalid treasury address");
        bogoToken = IERC20(_bogoToken);
        treasury = _treasury;
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
        external onlyAuthorized nonReentrant whenNotPaused {
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
        
        RewardTemplate memory template = templates["referral_bonus"];
        require(template.active, "Referral rewards not active");
        
        _resetDailyLimit();
        require(dailyDistributed + template.fixedAmount <= DAILY_GLOBAL_LIMIT, 
                "Daily limit exceeded");
        
        // Update state
        referredBy[msg.sender] = referrer;
        referralCount[referrer]++;
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
    
    function setAuthorizedBackend(address backend, bool authorized) external onlyTreasury {
        authorizedBackends[backend] = authorized;
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
    
    // Treasury sweep function - for token migration during contract upgrades
    function treasurySweep(address token, address to, uint256 amount) external onlyTreasury {
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
    
    // Event for treasury sweep operations
    event TreasurySweep(address indexed token, address indexed to, uint256 amount);
}