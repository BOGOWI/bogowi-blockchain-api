// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC1155/extensions/ERC1155Supply.sol";
import "@openzeppelin/contracts/token/common/ERC2981.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title CommercialNFT
 * @dev NFT contract for commercial assets, controlled by BOGOWI business
 */
contract CommercialNFT is ERC1155, AccessControl, Pausable, ERC1155Supply, ERC2981, ReentrancyGuard {
    // Role constants
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant BUSINESS_ROLE = keccak256("BUSINESS_ROLE");
    bytes32 public constant TREASURY_ROLE = keccak256("TREASURY_ROLE");
    
    // Token type categories
    uint256 public constant EVENT_TICKET = 1;
    uint256 public constant COLLECTIBLE = 2;
    uint256 public constant MERCHANDISE = 3;
    uint256 public constant GAMING_ASSET = 4;
    
    // Token ID ranges
    uint256 public constant TICKET_ID_START = 10000;
    uint256 public constant TICKET_ID_END = 19999;
    uint256 public constant COLLECTIBLE_ID_START = 20000;
    uint256 public constant COLLECTIBLE_ID_END = 29999;
    uint256 public constant MERCHANDISE_ID_START = 30000;
    uint256 public constant MERCHANDISE_ID_END = 39999;
    uint256 public constant GAMING_ID_START = 40000;
    uint256 public constant GAMING_ID_END = 49999;
    
    // Constants
    uint256 public constant MAX_ROYALTY_PERCENTAGE = 1000; // 10% max royalty
    uint256 public constant DEFAULT_ROYALTY = 500; // 5%
    
    // Treasury address for withdrawals
    address public treasuryAddress;

    // Token metadata
    mapping(uint256 => TokenInfo) public tokenInfo;
    mapping(uint256 => uint256) public maxSupply;
    mapping(uint256 => bool) public tokenExists;
    
    struct TokenInfo {
        string uri;
        uint256 price;
        bool burnable;
        bool tradeable;
        uint256 royaltyPercentage; // Basis points (100 = 1%)
    }

    // Event data for tickets
    mapping(uint256 => EventData) public eventData;
    
    struct EventData {
        uint256 eventDate;
        uint256 expiryDate;
        string venue;
        bool used;
    }

    event CommercialNFTMinted(
        address indexed recipient,
        uint256 indexed tokenId,
        uint256 tokenType,
        uint256 price
    );

    event TicketRedeemed(uint256 indexed tokenId, address indexed holder);
    event TokenURIUpdated(uint256 indexed tokenId, string newUri);
    event FundsWithdrawn(address indexed recipient, uint256 amount);
    event TreasuryAddressUpdated(address indexed oldTreasury, address indexed newTreasury);

    constructor(address _treasuryAddress) ERC1155("") {
        require(_treasuryAddress != address(0), "Invalid treasury address");
        
        // Grant admin roles to deployer
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
        _grantRole(BUSINESS_ROLE, msg.sender);
        _grantRole(TREASURY_ROLE, _treasuryAddress);
        
        // Set treasury address
        treasuryAddress = _treasuryAddress;
        
        // Set default royalty for all tokens (5%) - royalties go to contract
        _setDefaultRoyalty(address(this), uint96(DEFAULT_ROYALTY));
    }

    /**
     * @dev Mint event tickets
     */
    function mintEventTicket(
        address to,
        uint256 tokenId,
        uint256 eventDate,
        uint256 expiryDate,
        string memory venue,
        string memory tokenUri,
        uint256 price
    ) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= TICKET_ID_START && tokenId <= TICKET_ID_END, "Invalid ticket token ID range");
        require(!tokenExists[tokenId], "Token ID already exists");
        require(eventDate > block.timestamp, "Event date must be in future");
        require(expiryDate > eventDate, "Expiry must be after event date");
        require(bytes(venue).length > 0, "Venue cannot be empty");
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        
        eventData[tokenId] = EventData({
            eventDate: eventDate,
            expiryDate: expiryDate,
            venue: venue,
            used: false
        });
        
        tokenInfo[tokenId] = TokenInfo({
            uri: tokenUri,
            price: price,
            burnable: true,
            tradeable: true,
            royaltyPercentage: DEFAULT_ROYALTY
        });
        
        tokenExists[tokenId] = true;
        _mint(to, tokenId, 1, "");
        
        emit CommercialNFTMinted(to, tokenId, EVENT_TICKET, price);
    }

    /**
     * @dev Mint collectibles with rarity
     */
    function mintCollectible(
        address to,
        uint256 tokenId,
        uint256 amount,
        uint256 _maxSupply,
        string memory tokenUri,
        uint256 price,
        uint256 royaltyPercentage
    ) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= COLLECTIBLE_ID_START && tokenId <= COLLECTIBLE_ID_END, "Invalid collectible token ID range");
        require(amount > 0, "Amount must be greater than 0");
        require(_maxSupply > 0, "Max supply must be greater than 0");
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(royaltyPercentage <= MAX_ROYALTY_PERCENTAGE, "Royalty exceeds maximum");
        require(super.totalSupply(tokenId) + amount <= _maxSupply, "Exceeds max supply");
        
        if (!tokenExists[tokenId]) {
            tokenExists[tokenId] = true;
            maxSupply[tokenId] = _maxSupply;
            tokenInfo[tokenId] = TokenInfo({
                uri: tokenUri,
                price: price,
                burnable: false,
                tradeable: true,
                royaltyPercentage: royaltyPercentage
            });
            
            // Set specific royalty for this token - royalties go to the contract address
            // The TREASURY_ROLE can withdraw accumulated royalties via withdraw()
            _setTokenRoyalty(tokenId, address(this), uint96(royaltyPercentage));
        } else {
            require(maxSupply[tokenId] == _maxSupply, "Max supply mismatch");
        }
        
        _mint(to, tokenId, amount, "");
        
        emit CommercialNFTMinted(to, tokenId, COLLECTIBLE, price);
    }

    /**
     * @dev Batch mint for promotional campaigns
     */
    function batchMintPromo(
        address[] memory recipients,
        uint256 tokenId,
        uint256 amount,
        uint256 _maxSupply,
        string memory tokenUri,
        uint256 price
    ) external onlyRole(BUSINESS_ROLE) {
        require(recipients.length > 0, "No recipients");
        require(recipients.length <= 100, "Too many recipients");
        require(tokenId >= MERCHANDISE_ID_START && tokenId <= MERCHANDISE_ID_END, "Invalid promo token ID range");
        require(amount > 0, "Amount must be greater than 0");
        
        uint256 totalAmount = recipients.length * amount;
        
        // Initialize token if new
        if (!tokenExists[tokenId]) {
            require(_maxSupply > 0, "Max supply must be set for new token");
            require(bytes(tokenUri).length > 0, "URI cannot be empty");
            require(totalAmount <= _maxSupply, "Total mint exceeds max supply");
            
            tokenExists[tokenId] = true;
            maxSupply[tokenId] = _maxSupply;
            tokenInfo[tokenId] = TokenInfo({
                uri: tokenUri,
                price: price,
                burnable: true,
                tradeable: true,
                royaltyPercentage: DEFAULT_ROYALTY
            });
        } else {
            require(super.totalSupply(tokenId) + totalAmount <= maxSupply[tokenId], "Exceeds max supply");
        }
        
        for (uint256 i = 0; i < recipients.length; i++) {
            require(recipients[i] != address(0), "Invalid recipient");
            _mint(recipients[i], tokenId, amount, "");
        }
    }

    /**
     * @dev Redeem event ticket
     */
    function redeemTicket(uint256 tokenId) external {
        require(balanceOf(msg.sender, tokenId) > 0, "Not ticket holder");
        require(tokenId >= 10000 && tokenId < 20000, "Not a ticket");
        require(!eventData[tokenId].used, "Ticket already used");
        require(block.timestamp <= eventData[tokenId].expiryDate, "Ticket expired");
        
        eventData[tokenId].used = true;
        emit TicketRedeemed(tokenId, msg.sender);
    }

    /**
     * @dev Burn tokens (if burnable)
     */
    function burn(uint256 tokenId, uint256 amount) external {
        require(tokenInfo[tokenId].burnable, "Token not burnable");
        _burn(msg.sender, tokenId, amount);
    }

    /**
     * @dev Override to check tradeability
     */
    function _update(
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts
    ) internal override(ERC1155, ERC1155Supply) whenNotPaused {
        for (uint256 i = 0; i < ids.length; i++) {
            if (from != address(0) && to != address(0)) { // Not mint or burn
                require(tokenInfo[ids[i]].tradeable, "Token not tradeable");
            }
        }
        super._update(from, to, ids, amounts);
    }

    function uri(uint256 tokenId) public view override returns (string memory) {
        return tokenInfo[tokenId].uri;
    }

    /**
     * @dev Update token URI (for dynamic metadata)
     */
    function updateTokenURI(uint256 tokenId, string memory newUri) external onlyRole(BUSINESS_ROLE) {
        require(tokenExists[tokenId], "Token does not exist");
        require(bytes(newUri).length > 0, "URI cannot be empty");
        tokenInfo[tokenId].uri = newUri;
        emit TokenURIUpdated(tokenId, newUri);
    }

    /**
     * @dev Mint gaming assets
     */
    function mintGamingAsset(
        address to,
        uint256 tokenId,
        uint256 amount,
        uint256 _maxSupply,
        string memory tokenUri,
        uint256 price
    ) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= GAMING_ID_START && tokenId <= GAMING_ID_END, "Invalid gaming token ID range");
        require(amount > 0, "Amount must be greater than 0");
        require(_maxSupply > 0, "Max supply must be greater than 0");
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(super.totalSupply(tokenId) + amount <= _maxSupply, "Exceeds max supply");
        
        if (!tokenExists[tokenId]) {
            tokenExists[tokenId] = true;
            maxSupply[tokenId] = _maxSupply;
            tokenInfo[tokenId] = TokenInfo({
                uri: tokenUri,
                price: price,
                burnable: true,
                tradeable: true,
                royaltyPercentage: DEFAULT_ROYALTY
            });
        } else {
            require(maxSupply[tokenId] == _maxSupply, "Max supply mismatch");
        }
        
        _mint(to, tokenId, amount, "");
        emit CommercialNFTMinted(to, tokenId, GAMING_ASSET, price);
    }
    
    /**
     * @dev Withdraw accumulated funds to the treasury address
     */
    function withdraw() external onlyRole(TREASURY_ROLE) nonReentrant {
        require(treasuryAddress != address(0), "Treasury address not set");
        
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        (bool success, ) = payable(treasuryAddress).call{value: balance}("");
        require(success, "Withdrawal failed");
        
        emit FundsWithdrawn(treasuryAddress, balance);
    }
    
    /**
     * @dev Update treasury address (admin only)
     * @param newTreasuryAddress The new treasury address
     */
    function setTreasuryAddress(address newTreasuryAddress) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newTreasuryAddress != address(0), "Invalid treasury address");
        
        // Revoke TREASURY_ROLE from old address and grant to new
        address oldTreasury = treasuryAddress;
        _revokeRole(TREASURY_ROLE, oldTreasury);
        _grantRole(TREASURY_ROLE, newTreasuryAddress);
        
        treasuryAddress = newTreasuryAddress;
        
        emit TreasuryAddressUpdated(oldTreasury, newTreasuryAddress);
    }

    function pause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    function unpause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC1155, ERC2981, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
    
    
    /**
     * @dev Receive function to accept ETH
     */
    receive() external payable {}
}
