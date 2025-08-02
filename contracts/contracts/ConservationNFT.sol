// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC1155/extensions/ERC1155Supply.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./StandardErrors.sol";

/**
 * @title ConservationNFT
 * @author BOGOWI Team
 * @notice ERC1155 NFT contract for conservation achievements and wildlife adoption certificates
 * @dev Implements soulbound tokens for permanent recognition and DAO-controlled minting
 * Features:
 * - Wildlife adoption certificates with impact tracking
 * - Conservation achievement badges
 * - Donor recognition tokens
 * - Environmental impact certificates
 * - Soulbound token support for non-transferable achievements
 * @custom:security-contact security@bogowi.com
 */
contract ConservationNFT is ERC1155, AccessControl, Pausable, ERC1155Supply, ReentrancyGuard, StandardErrors {
    // Role constants
    bytes32 public constant DAO_ROLE = keccak256("DAO_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    
    // Token type categories
    uint256 public constant ADOPTION_CERTIFICATE = 1;
    uint256 public constant CONSERVATION_BADGE = 2;
    uint256 public constant DONOR_RECOGNITION = 3;
    uint256 public constant IMPACT_CERTIFICATE = 4;
    
    // Token ID ranges
    uint256 public constant ADOPTION_ID_START = 1000;
    uint256 public constant ADOPTION_ID_END = 1999;
    uint256 public constant BADGE_ID_START = 2000;
    uint256 public constant BADGE_ID_END = 2999;
    uint256 public constant DONOR_ID_START = 3000;
    uint256 public constant DONOR_ID_END = 3999;
    uint256 public constant IMPACT_ID_START = 4000;
    uint256 public constant IMPACT_ID_END = 4999;
    
    // Maximum supplies per type
    mapping(uint256 => uint256) public maxSupplyPerType;
    mapping(uint256 => uint256) public mintedPerType;

    // Token metadata
    mapping(uint256 => string) public tokenURIs;
    mapping(uint256 => bool) public isSoulbound;
    mapping(uint256 => uint256) public donationAmounts;
    mapping(uint256 => bool) public tokenExists;
    
    // Conservation data
    mapping(uint256 => ConservationData) public conservationData;
    
    /**
     * @dev Conservation-specific data structure
     * @param species Wildlife species name
     * @param location Conservation location
     * @param impactScore Environmental impact score (1-1000)
     * @param date Unix timestamp of conservation activity
     * @param verified Whether data has been verified by DAO
     */
    struct ConservationData {
        string species;
        string location;
        uint256 impactScore;
        uint256 date;
        bool verified;
    }

    event ConservationNFTMinted(
        address indexed recipient,
        uint256 indexed tokenId,
        uint256 tokenType,
        string species,
        string location,
        uint256 impactScore
    );
    
    event MaxSupplySet(uint256 tokenType, uint256 maxSupply);
    event TokenURIUpdated(uint256 indexed tokenId, string newUri);

    /**
     * @notice Initializes the ConservationNFT contract
     * @dev Sets up roles and default max supplies for each token type
     */
    constructor() ERC1155("") {
        // Grant admin roles to deployer
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(DAO_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
        
        // Set default max supplies (can be updated by DAO)
        maxSupplyPerType[ADOPTION_CERTIFICATE] = 10000;
        maxSupplyPerType[CONSERVATION_BADGE] = 5000;
        maxSupplyPerType[DONOR_RECOGNITION] = 10000;
        maxSupplyPerType[IMPACT_CERTIFICATE] = 5000;
    }

    /**
     * @notice Mints a wildlife adoption certificate NFT
     * @dev Token IDs must be in range 1000-1999, soulbound by default
     * @param to Recipient address
     * @param tokenId Unique ID for the certificate (1000-1999)
     * @param species Wildlife species being adopted
     * @param location Geographic location of conservation
     * @param tokenUri Metadata URI
     * @param impactScore Environmental impact score (1-1000)
     * @custom:emits ConservationNFTMinted
     */
    function mintAdoptionCertificate(
        address to,
        uint256 tokenId,
        string memory species,
        string memory location,
        string memory tokenUri,
        uint256 impactScore
    ) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= ADOPTION_ID_START && tokenId <= ADOPTION_ID_END, INVALID_PARAMETER);
        require(!tokenExists[tokenId], "Token ID already exists");
        require(bytes(species).length > 0, EMPTY_STRING);
        require(bytes(location).length > 0, EMPTY_STRING);
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(impactScore > 0 && impactScore <= 1000, "Invalid impact score");
        require(mintedPerType[ADOPTION_CERTIFICATE] < maxSupplyPerType[ADOPTION_CERTIFICATE], MAX_REACHED);
        
        conservationData[tokenId] = ConservationData({
            species: species,
            location: location,
            impactScore: impactScore,
            date: block.timestamp,
            verified: true
        });
        
        tokenURIs[tokenId] = tokenUri;
        isSoulbound[tokenId] = true; // Adoptions are non-transferable
        tokenExists[tokenId] = true;
        mintedPerType[ADOPTION_CERTIFICATE]++;
        
        _mint(to, tokenId, 1, "");
        
        emit ConservationNFTMinted(to, tokenId, ADOPTION_CERTIFICATE, species, location, impactScore);
    }

    /**
     * @notice Mints a conservation achievement badge
     * @dev Token IDs must be in range 2000-2999, soulbound by default
     * @param to Recipient address
     * @param tokenId Unique ID for the badge (2000-2999)
     * @param impactScore Environmental impact score (1-1000)
     * @param tokenUri Metadata URI
     * @param achievementType Type of conservation achievement
     * @custom:emits ConservationNFTMinted
     */
    function mintConservationBadge(
        address to,
        uint256 tokenId,
        uint256 impactScore,
        string memory tokenUri,
        string memory achievementType
    ) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= BADGE_ID_START && tokenId <= BADGE_ID_END, INVALID_PARAMETER);
        require(!tokenExists[tokenId], "Token ID already exists");
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(bytes(achievementType).length > 0, EMPTY_STRING);
        require(impactScore > 0 && impactScore <= 1000, "Invalid impact score");
        require(mintedPerType[CONSERVATION_BADGE] < maxSupplyPerType[CONSERVATION_BADGE], MAX_REACHED);
        
        conservationData[tokenId] = ConservationData({
            species: achievementType,
            location: "Global",
            impactScore: impactScore,
            date: block.timestamp,
            verified: true
        });
        
        tokenURIs[tokenId] = tokenUri;
        isSoulbound[tokenId] = true;
        tokenExists[tokenId] = true;
        mintedPerType[CONSERVATION_BADGE]++;
        
        _mint(to, tokenId, 1, "");
        
        emit ConservationNFTMinted(to, tokenId, CONSERVATION_BADGE, achievementType, "Global", impactScore);
    }

    /**
     * @notice Mints a donor recognition NFT
     * @dev Token IDs must be in range 3000-3999, tracks donation amount
     * @param to Recipient address
     * @param tokenId Unique ID for the recognition (3000-3999)
     * @param donationAmount Donation amount in wei
     * @param tokenUri Metadata URI
     * @param campaign Conservation campaign name
     * @custom:emits ConservationNFTMinted
     */
    function mintDonorNFT(
        address to,
        uint256 tokenId,
        uint256 donationAmount,
        string memory tokenUri,
        string memory campaign
    ) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= DONOR_ID_START && tokenId <= DONOR_ID_END, INVALID_PARAMETER);
        require(!tokenExists[tokenId], "Token ID already exists");
        require(donationAmount > 0, ZERO_AMOUNT);
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(bytes(campaign).length > 0, EMPTY_STRING);
        require(mintedPerType[DONOR_RECOGNITION] < maxSupplyPerType[DONOR_RECOGNITION], MAX_REACHED);
        
        // Calculate impact score based on donation amount
        uint256 impactScore = calculateDonationImpact(donationAmount);
        
        conservationData[tokenId] = ConservationData({
            species: campaign,
            location: "Global",
            impactScore: impactScore,
            date: block.timestamp,
            verified: true
        });
        
        donationAmounts[tokenId] = donationAmount;
        tokenURIs[tokenId] = tokenUri;
        isSoulbound[tokenId] = true;
        tokenExists[tokenId] = true;
        mintedPerType[DONOR_RECOGNITION]++;
        
        _mint(to, tokenId, 1, "");
        
        emit ConservationNFTMinted(to, tokenId, DONOR_RECOGNITION, campaign, "Global", impactScore);
    }

    /**
     * @dev Internal transfer function with soulbound enforcement
     * @notice Overrides ERC1155 to prevent transfers of soulbound tokens
     * @param from Sender address
     * @param to Recipient address
     * @param ids Array of token IDs
     * @param amounts Array of amounts
     */
    function _update(
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts
    ) internal override(ERC1155, ERC1155Supply) whenNotPaused {
        for (uint256 i = 0; i < ids.length; i++) {
            if (isSoulbound[ids[i]] && from != address(0)) {
                revert("Soulbound token cannot be transferred");
            }
        }
        super._update(from, to, ids, amounts);
    }

    /**
     * @notice Returns the metadata URI for a token
     * @param tokenId ID of the token
     * @return Metadata URI string
     */
    function uri(uint256 tokenId) public view override returns (string memory) {
        return tokenURIs[tokenId];
    }

    /**
     * @notice Mints an environmental impact certificate
     * @dev Token IDs must be in range 4000-4999, requires DAO_ROLE
     * @param to Recipient address
     * @param tokenId Unique ID for the certificate (4000-4999)
     * @param project Conservation project name
     * @param location Project location
     * @param impactScore Environmental impact score (1-1000)
     * @param tokenUri Metadata URI
     * @custom:emits ConservationNFTMinted
     */
    function mintImpactCertificate(
        address to,
        uint256 tokenId,
        string memory project,
        string memory location,
        uint256 impactScore,
        string memory tokenUri
    ) external onlyRole(DAO_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= IMPACT_ID_START && tokenId <= IMPACT_ID_END, INVALID_PARAMETER);
        require(!tokenExists[tokenId], "Token ID already exists");
        require(bytes(project).length > 0, EMPTY_STRING);
        require(bytes(location).length > 0, EMPTY_STRING);
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(impactScore > 0 && impactScore <= 1000, "Invalid impact score");
        require(mintedPerType[IMPACT_CERTIFICATE] < maxSupplyPerType[IMPACT_CERTIFICATE], MAX_REACHED);
        
        conservationData[tokenId] = ConservationData({
            species: project,
            location: location,
            impactScore: impactScore,
            date: block.timestamp,
            verified: true
        });
        
        tokenURIs[tokenId] = tokenUri;
        isSoulbound[tokenId] = true;
        tokenExists[tokenId] = true;
        mintedPerType[IMPACT_CERTIFICATE]++;
        
        _mint(to, tokenId, 1, "");
        
        emit ConservationNFTMinted(to, tokenId, IMPACT_CERTIFICATE, project, location, impactScore);
    }
    
    /**
     * @notice Sets maximum supply for a token type
     * @dev Only DAO_ROLE can update max supplies
     * @param tokenType Type of token (1-4)
     * @param maxSupply New maximum supply
     * @custom:emits MaxSupplySet
     */
    function setMaxSupply(uint256 tokenType, uint256 maxSupply) external onlyRole(DAO_ROLE) {
        require(tokenType >= ADOPTION_CERTIFICATE && tokenType <= IMPACT_CERTIFICATE, INVALID_PARAMETER);
        require(maxSupply >= mintedPerType[tokenType], INVALID_PARAMETER);
        
        maxSupplyPerType[tokenType] = maxSupply;
        emit MaxSupplySet(tokenType, maxSupply);
    }
    
    /**
     * @notice Updates the metadata URI for a token
     * @dev Only DAO_ROLE can update URIs
     * @param tokenId ID of token to update
     * @param newUri New metadata URI
     * @custom:emits TokenURIUpdated
     */
    function updateTokenURI(uint256 tokenId, string memory newUri) external onlyRole(DAO_ROLE) {
        require(tokenExists[tokenId], DOES_NOT_EXIST);
        require(bytes(newUri).length > 0, EMPTY_STRING);
        
        tokenURIs[tokenId] = newUri;
        emit TokenURIUpdated(tokenId, newUri);
    }
    
    /**
     * @dev Calculates impact score based on donation amount
     * @param amount Donation amount in wei
     * @return Impact score (50-1000)
     */
    function calculateDonationImpact(uint256 amount) internal pure returns (uint256) {
        if (amount >= 10 ether) return 1000;
        if (amount >= 5 ether) return 750;
        if (amount >= 1 ether) return 500;
        if (amount >= 0.5 ether) return 250;
        if (amount >= 0.1 ether) return 100;
        return 50;
    }
    
    /**
     * @notice Returns the token type based on token ID
     * @param tokenId ID of the token
     * @return Token type (1-4)
     */
    function getTokenType(uint256 tokenId) public pure returns (uint256) {
        if (tokenId >= ADOPTION_ID_START && tokenId <= ADOPTION_ID_END) return ADOPTION_CERTIFICATE;
        if (tokenId >= BADGE_ID_START && tokenId <= BADGE_ID_END) return CONSERVATION_BADGE;
        if (tokenId >= DONOR_ID_START && tokenId <= DONOR_ID_END) return DONOR_RECOGNITION;
        if (tokenId >= IMPACT_ID_START && tokenId <= IMPACT_ID_END) return IMPACT_CERTIFICATE;
        revert("Invalid token ID");
    }
    
    /**
     * @notice Pauses all token transfers and minting
     * @dev Emergency function for DEFAULT_ADMIN_ROLE or PAUSER_ROLE
     */
    function pause() public {
        require(
            hasRole(DEFAULT_ADMIN_ROLE, msg.sender) || hasRole(PAUSER_ROLE, msg.sender),
            UNAUTHORIZED
        );
        _pause();
    }

    /**
     * @notice Unpauses token transfers and minting
     */
    function unpause() public {
        require(
            hasRole(DEFAULT_ADMIN_ROLE, msg.sender) || hasRole(PAUSER_ROLE, msg.sender),
            UNAUTHORIZED
        );
        _unpause();
    }
    
    /**
     * @notice Checks if contract supports an interface
     * @dev Combines ERC1155 and AccessControl interfaces
     * @param interfaceId Interface identifier to check
     * @return True if interface is supported
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC1155, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
