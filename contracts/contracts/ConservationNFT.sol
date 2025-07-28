// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/token/ERC1155/extensions/ERC1155Supply.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title ConservationNFT
 * @dev NFT contract for conservation-related assets, controlled by DAO
 */
contract ConservationNFT is ERC1155, AccessControl, Pausable, ERC1155Supply, ReentrancyGuard {
    // Role constants
    bytes32 public constant DAO_ROLE = keccak256("DAO_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    
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
     * @dev Mint wildlife adoption certificate
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
        require(tokenId >= ADOPTION_ID_START && tokenId <= ADOPTION_ID_END, "Invalid adoption token ID range");
        require(!tokenExists[tokenId], "Token ID already exists");
        require(bytes(species).length > 0, "Species cannot be empty");
        require(bytes(location).length > 0, "Location cannot be empty");
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(impactScore > 0 && impactScore <= 1000, "Invalid impact score");
        require(mintedPerType[ADOPTION_CERTIFICATE] < maxSupplyPerType[ADOPTION_CERTIFICATE], "Max supply reached");
        
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
     * @dev Mint conservation achievement badge
     */
    function mintConservationBadge(
        address to,
        uint256 tokenId,
        uint256 impactScore,
        string memory tokenUri,
        string memory achievementType
    ) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= BADGE_ID_START && tokenId <= BADGE_ID_END, "Invalid badge token ID range");
        require(!tokenExists[tokenId], "Token ID already exists");
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(bytes(achievementType).length > 0, "Achievement type cannot be empty");
        require(impactScore > 0 && impactScore <= 1000, "Invalid impact score");
        require(mintedPerType[CONSERVATION_BADGE] < maxSupplyPerType[CONSERVATION_BADGE], "Max supply reached");
        
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
     * @dev Mint donor recognition NFT
     */
    function mintDonorNFT(
        address to,
        uint256 tokenId,
        uint256 donationAmount,
        string memory tokenUri,
        string memory campaign
    ) external onlyRole(MINTER_ROLE) {
        require(to != address(0), "Invalid recipient address");
        require(tokenId >= DONOR_ID_START && tokenId <= DONOR_ID_END, "Invalid donor token ID range");
        require(!tokenExists[tokenId], "Token ID already exists");
        require(donationAmount > 0, "Donation amount must be greater than 0");
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(bytes(campaign).length > 0, "Campaign cannot be empty");
        require(mintedPerType[DONOR_RECOGNITION] < maxSupplyPerType[DONOR_RECOGNITION], "Max supply reached");
        
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
     * @dev Override to prevent transfers of soulbound tokens
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

    function uri(uint256 tokenId) public view override returns (string memory) {
        return tokenURIs[tokenId];
    }

    /**
     * @dev Mint impact certificate for verified conservation outcomes
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
        require(tokenId >= IMPACT_ID_START && tokenId <= IMPACT_ID_END, "Invalid impact token ID range");
        require(!tokenExists[tokenId], "Token ID already exists");
        require(bytes(project).length > 0, "Project cannot be empty");
        require(bytes(location).length > 0, "Location cannot be empty");
        require(bytes(tokenUri).length > 0, "URI cannot be empty");
        require(impactScore > 0 && impactScore <= 1000, "Invalid impact score");
        require(mintedPerType[IMPACT_CERTIFICATE] < maxSupplyPerType[IMPACT_CERTIFICATE], "Max supply reached");
        
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
     * @dev Set maximum supply for a token type
     */
    function setMaxSupply(uint256 tokenType, uint256 maxSupply) external onlyRole(DAO_ROLE) {
        require(tokenType >= ADOPTION_CERTIFICATE && tokenType <= IMPACT_CERTIFICATE, "Invalid token type");
        require(maxSupply >= mintedPerType[tokenType], "Max supply less than already minted");
        
        maxSupplyPerType[tokenType] = maxSupply;
        emit MaxSupplySet(tokenType, maxSupply);
    }
    
    /**
     * @dev Update token URI
     */
    function updateTokenURI(uint256 tokenId, string memory newUri) external onlyRole(DAO_ROLE) {
        require(tokenExists[tokenId], "Token does not exist");
        require(bytes(newUri).length > 0, "URI cannot be empty");
        
        tokenURIs[tokenId] = newUri;
        emit TokenURIUpdated(tokenId, newUri);
    }
    
    /**
     * @dev Calculate impact score based on donation amount
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
     * @dev Get token type from token ID
     */
    function getTokenType(uint256 tokenId) public pure returns (uint256) {
        if (tokenId >= ADOPTION_ID_START && tokenId <= ADOPTION_ID_END) return ADOPTION_CERTIFICATE;
        if (tokenId >= BADGE_ID_START && tokenId <= BADGE_ID_END) return CONSERVATION_BADGE;
        if (tokenId >= DONOR_ID_START && tokenId <= DONOR_ID_END) return DONOR_RECOGNITION;
        if (tokenId >= IMPACT_ID_START && tokenId <= IMPACT_ID_END) return IMPACT_CERTIFICATE;
        revert("Invalid token ID");
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
        override(ERC1155, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
