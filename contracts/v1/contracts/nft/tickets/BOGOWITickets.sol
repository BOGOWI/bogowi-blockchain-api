// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IBOGOWITickets.sol";
import "../../base/RoleManaged.sol";
import "../../utils/Roles.sol";
import "../common/NFTErrors.sol";
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/token/common/ERC2981.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/EIP712.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/**
 * @title BOGOWITickets
 * @notice BOGOWI Tickets (BWTIX) - Phase 1 NFT ticketing system on Camino Network
 * @dev ERC-721 implementation with time-locked transfers, QR redemption, and BOGO rewards
 */
contract BOGOWITickets is 
    IBOGOWITickets,
    ERC721,
    ERC721URIStorage,
    ERC2981,
    EIP712,
    Pausable,
    ReentrancyGuard,
    RoleManaged 
{
    using ECDSA for bytes32;
    
    // Role definitions - use centralized constants from RoleManager
    // These should match the definitions in RoleManager contract
    
    // EIP-712 type hash for redemption
    bytes32 public constant REDEMPTION_TYPEHASH = keccak256(
        "RedeemTicket(uint256 tokenId,address redeemer,uint256 nonce,uint256 deadline,uint256 chainId)"
    );
    
    // Constants
    uint96 public constant DEFAULT_ROYALTY_BPS = 500; // 5% royalty
    uint256 public constant CAMINO_MAINNET_CHAIN_ID = 500;
    uint256 public constant CAMINO_TESTNET_CHAIN_ID = 501;
    uint256 public constant MAX_BATCH_SIZE = 100; // Reduced for gas safety (was 500)
    uint256 public constant BASE_GAS_PER_MINT = 150000; // Base gas per mint operation
    uint256 public constant ADDITIONAL_GAS_PER_MINT = 30000; // Additional gas buffer for metadata operations
    uint256 public constant INITIAL_TOKEN_ID = 10001; // Configurable start token ID
    address public conservationDAO; // Royalty receiver
    
    // Datakyte integration
    string private _baseTokenURI;
    
    // State variables
    uint256 private _nextTokenId;
    mapping(uint256 => TicketData) private _tickets;
    mapping(uint256 => bool) private _usedNonces;
    mapping(bytes32 => bool) private _usedBookingIds;
    
    // MEV protection - commit-reveal for high-value operations
    mapping(bytes32 => uint256) private _commitments;
    uint256 public constant COMMIT_DELAY = 1 minutes; // Minimum delay between commit and reveal
    uint256 public expiryGracePeriod = 5 minutes; // Configurable grace period for expired tickets
    
    // Constructor
    constructor(
        address _roleManager,
        address _conservationDAO
    ) 
        ERC721("BOGOWI Tickets", "BWTIX")
        EIP712("BOGOWITickets", "1")
        RoleManaged(_roleManager)
    {
        require(_conservationDAO != address(0), "Invalid DAO address");
        
        // Validate deployment on Camino network
        uint256 chainId = block.chainid;
        require(
            chainId == CAMINO_MAINNET_CHAIN_ID || chainId == CAMINO_TESTNET_CHAIN_ID,
            "Must deploy on Camino network"
        );
        
        conservationDAO = _conservationDAO;
        
        // Set default royalty to Conservation DAO at 5%
        _setDefaultRoyalty(_conservationDAO, DEFAULT_ROYALTY_BPS);
        
        _nextTokenId = INITIAL_TOKEN_ID; // Start token IDs at configured value
        
        // Set default Datakyte base URI
        _baseTokenURI = "https://dklnk.to/api/nfts/";
    }
    
    /**
     * @notice Mint a new ticket NFT
     * @dev Emits TicketMinted event with BOGO reward basis points
     */
    function mintTicket(
        MintParams memory params
    ) external override onlyRole(Roles.NFT_MINTER_ROLE) whenNotPaused nonReentrant returns (uint256) {
        // Validations
        require(params.to != address(0), "Cannot mint to zero address");
        require(!_usedBookingIds[params.bookingId], "Booking ID already used");
        require(params.expiresAt > block.timestamp, "Expiry must be in future");
        require(params.transferUnlockAt < params.expiresAt, "Unlock must be before expiry");
        
        // Additional timestamp validations
        require(params.expiresAt <= block.timestamp + 365 days, "Expiry too far in future");
        
        uint256 tokenId = _nextTokenId++;
        
        // Store ticket data - optimized struct packing order
        _tickets[tokenId] = TicketData({
            bookingId: params.bookingId,
            eventId: params.eventId,
            transferUnlockAt: params.transferUnlockAt,
            expiresAt: params.expiresAt,
            utilityFlags: params.utilityFlags,
            state: TicketState.ISSUED,
            nonTransferableAfterRedeem: params.utilityFlags & 0x01 == 0, // Bit 0: if set, allow transfer after redeem
            burnOnRedeem: params.utilityFlags & 0x02 != 0 // Bit 1: if set, burn on redeem
        });
        
        // Mark booking ID as used
        _usedBookingIds[params.bookingId] = true;
        
        // Mint the NFT
        _safeMint(params.to, tokenId);
        
        // Set metadata URI if provided
        if (bytes(params.metadataURI).length > 0) {
            _setTokenURI(tokenId, params.metadataURI);
        }
        
        // Emit event with BOGO reward basis points
        emit TicketMinted(
            tokenId,
            params.bookingId,
            params.eventId,
            params.to,
            params.rewardBasisPoints
        );
        
        return tokenId;
    }
    
    /**
     * @notice Mint multiple tickets in batch
     * @dev Limited to MAX_BATCH_SIZE to prevent DoS attacks
     */
    function mintBatch(
        MintParams[] memory params
    ) external override onlyRole(Roles.NFT_MINTER_ROLE) whenNotPaused nonReentrant returns (uint256[] memory) {
        require(params.length > 0, "Empty batch");
        require(params.length <= MAX_BATCH_SIZE, "Batch size exceeds maximum");
        
        // Enhanced gas estimation considering metadata operations
        uint256 estimatedGas = params.length * BASE_GAS_PER_MINT;
        
        // Add extra gas for metadata operations if present
        for (uint256 i = 0; i < params.length; i++) {
            if (bytes(params[i].metadataURI).length > 0) {
                estimatedGas += ADDITIONAL_GAS_PER_MINT;
            }
        }
        
        // Safety buffer: require 20% more gas than estimated
        uint256 requiredGas = (estimatedGas * 120) / 100;
        require(gasleft() >= requiredGas, "Insufficient gas for batch");
        
        emit BatchMintStarted(params.length, msg.sender);
        
        uint256[] memory tokenIds = new uint256[](params.length);
        
        for (uint256 i = 0; i < params.length; i++) {
            // Call internal mint logic directly to avoid external call
            tokenIds[i] = _mintTicketInternal(params[i]);
        }
        
        return tokenIds;
    }
    
    /**
     * @dev Internal mint function without role check for batch operations
     */
    function _mintTicketInternal(
        MintParams memory params
    ) internal returns (uint256) {
        // Validations
        require(params.to != address(0), "Cannot mint to zero address");
        require(!_usedBookingIds[params.bookingId], "Booking ID already used");
        require(params.expiresAt > block.timestamp, "Expiry must be in future");
        require(params.transferUnlockAt < params.expiresAt, "Unlock must be before expiry");
        
        // Additional timestamp validations
        require(params.expiresAt <= block.timestamp + 365 days, "Expiry too far in future");
        
        uint256 tokenId = _nextTokenId++;
        
        // Store ticket data - optimized struct packing order
        _tickets[tokenId] = TicketData({
            bookingId: params.bookingId,
            eventId: params.eventId,
            transferUnlockAt: params.transferUnlockAt,
            expiresAt: params.expiresAt,
            utilityFlags: params.utilityFlags,
            state: TicketState.ISSUED,
            nonTransferableAfterRedeem: params.utilityFlags & 0x01 == 0, // Bit 0: if set, allow transfer after redeem
            burnOnRedeem: params.utilityFlags & 0x02 != 0 // Bit 1: if set, burn on redeem
        });
        
        // Mark booking ID as used
        _usedBookingIds[params.bookingId] = true;
        
        // Mint the NFT
        _safeMint(params.to, tokenId);
        
        // Set metadata URI if provided
        if (bytes(params.metadataURI).length > 0) {
            _setTokenURI(tokenId, params.metadataURI);
        }
        
        // Emit event with BOGO reward basis points
        emit TicketMinted(
            tokenId,
            params.bookingId,
            params.eventId,
            params.to,
            params.rewardBasisPoints
        );
        
        return tokenId;
    }
    
    /**
     * @notice Redeem a ticket using QR code and EIP-712 signature
     */
    function redeemTicket(
        RedemptionData memory redemptionData
    ) external override whenNotPaused nonReentrant {
        uint256 tokenId = redemptionData.tokenId;
        
        // Validations
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        require(_tickets[tokenId].state == TicketState.ISSUED, "Ticket not redeemable");
        require(block.timestamp < _tickets[tokenId].expiresAt, "Ticket expired");
        require(!_usedNonces[redemptionData.nonce], "Nonce already used");
        require(block.timestamp <= redemptionData.deadline, "Signature expired");
        require(redemptionData.chainId == block.chainid, "Invalid chain ID");
        
        // Verify signature
        require(verifyRedemptionSignature(redemptionData), "Invalid signature");
        
        // Mark nonce as used
        _usedNonces[redemptionData.nonce] = true;
        emit NonceUsed(redemptionData.nonce, redemptionData.redeemer);
        
        // Update ticket state
        _tickets[tokenId].state = TicketState.REDEEMED;
        
        // Handle burn on redeem if configured
        if (_tickets[tokenId].burnOnRedeem) {
            emit TicketBurned(tokenId, ownerOf(tokenId));
            _burn(tokenId);
        }
        
        emit TicketRedeemed(tokenId, redemptionData.redeemer, block.timestamp);
    }
    
    /**
     * @notice Update transfer unlock time for a ticket
     */
    function updateTransferUnlock(
        uint256 tokenId,
        uint64 newUnlockTime
    ) external override onlyRole(Roles.ADMIN_ROLE) {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        require(newUnlockTime < _tickets[tokenId].expiresAt, "Unlock must be before expiry");
        
        _tickets[tokenId].transferUnlockAt = newUnlockTime;
        
        emit TransferUnlockUpdated(tokenId, newUnlockTime);
    }
    
    /**
     * @notice Mark a ticket as expired
     * @dev Restricted to ADMIN_ROLE to prevent griefing attacks
     */
    function expireTicket(uint256 tokenId) external override onlyRole(Roles.ADMIN_ROLE) {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        require(block.timestamp >= _tickets[tokenId].expiresAt, "Ticket not yet expired");
        require(_tickets[tokenId].state == TicketState.ISSUED, "Ticket already processed");
        
        // Add grace period to prevent frontrunning - ticket must be expired for configured grace period
        require(block.timestamp >= _tickets[tokenId].expiresAt + expiryGracePeriod, "Grace period not met");
        
        _tickets[tokenId].state = TicketState.EXPIRED;
        
        emit TicketExpired(tokenId);
    }
    
    /**
     * @notice Verify a redemption signature using EIP-712
     */
    function verifyRedemptionSignature(
        RedemptionData memory redemptionData
    ) public view override returns (bool) {
        bytes32 structHash = keccak256(
            abi.encode(
                REDEMPTION_TYPEHASH,
                redemptionData.tokenId,
                redemptionData.redeemer,
                redemptionData.nonce,
                redemptionData.deadline,
                redemptionData.chainId
            )
        );
        
        bytes32 hash = _hashTypedDataV4(structHash);
        address signer = hash.recover(redemptionData.signature);
        
        // Signer must have BACKEND_ROLE (backend service)
        return roleManager.hasRole(Roles.BACKEND_ROLE, signer);
    }
    
    // View Functions
    
    function getTicketData(uint256 tokenId) external view override returns (TicketData memory) {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        return _tickets[tokenId];
    }
    
    function isTransferable(uint256 tokenId) external view override returns (bool) {
        if (_ownerOf(tokenId) == address(0)) return false;
        
        TicketData memory ticket = _tickets[tokenId];
        
        // Check if transfer is unlocked
        if (block.timestamp < ticket.transferUnlockAt) return false;
        
        // Check if expired
        if (block.timestamp >= ticket.expiresAt) return false;
        
        // Check if redeemed and non-transferable after redemption
        if (ticket.state == TicketState.REDEEMED && ticket.nonTransferableAfterRedeem) return false;
        
        return true;
    }
    
    function isRedeemed(uint256 tokenId) external view override returns (bool) {
        return _tickets[tokenId].state == TicketState.REDEEMED;
    }
    
    function isExpired(uint256 tokenId) external view override returns (bool) {
        return block.timestamp >= _tickets[tokenId].expiresAt || 
               _tickets[tokenId].state == TicketState.EXPIRED;
    }
    
    // Transfer Override
    
    /**
     * @notice Override transfer to enforce time-lock and redemption rules
     */
    function _update(
        address to,
        uint256 tokenId,
        address auth
    ) internal override(ERC721) returns (address) {
        address from = _ownerOf(tokenId);
        
        // Allow minting and burning
        if (from != address(0) && to != address(0)) {
            TicketData memory ticket = _tickets[tokenId];
            
            // Check transfer lock
            require(
                block.timestamp >= ticket.transferUnlockAt,
                "Transfer locked until unlock time"
            );
            
            // Check expiry
            require(
                block.timestamp < ticket.expiresAt,
                "Cannot transfer expired ticket"
            );
            
            // Check redemption status
            require(
                ticket.state != TicketState.REDEEMED || !ticket.nonTransferableAfterRedeem,
                "Cannot transfer redeemed ticket"
            );
        }
        
        return super._update(to, tokenId, auth);
    }
    
    // Admin Functions
    
    /**
     * @notice Update royalty information
     */
    function setRoyaltyInfo(
        address receiver,
        uint96 feeBasisPoints
    ) external onlyRole(Roles.ADMIN_ROLE) {
        require(receiver != address(0), "Invalid receiver");
        require(feeBasisPoints <= 1000, "Royalty too high"); // Max 10%
        
        _setDefaultRoyalty(receiver, feeBasisPoints);
        
        emit RoyaltyInfoUpdated(receiver, feeBasisPoints);
    }
    
    /**
     * @notice Burn a ticket token
     * @dev Only the owner can burn their own ticket
     */
    function burn(uint256 tokenId) external {
        require(_ownerOf(tokenId) == msg.sender, "Not token owner");
        _burn(tokenId);
    }
    
    /**
     * @notice Set the expiry grace period
     * @param newGracePeriod The new grace period in seconds (minimum 1 minute, maximum 1 day)
     */
    function setExpiryGracePeriod(uint256 newGracePeriod) external onlyRole(Roles.ADMIN_ROLE) {
        require(newGracePeriod >= 1 minutes, "Grace period too short");
        require(newGracePeriod <= 1 days, "Grace period too long");
        
        uint256 oldGracePeriod = expiryGracePeriod;
        expiryGracePeriod = newGracePeriod;
        
        emit ExpiryGracePeriodUpdated(oldGracePeriod, newGracePeriod);
    }
    
    /**
     * @notice Set the base URI for Datakyte metadata
     * @param newBaseURI The new base URI (e.g., "https://dklnk.to/api/nfts/")
     */
    function setBaseURI(string memory newBaseURI) external onlyRole(Roles.ADMIN_ROLE) {
        _baseTokenURI = newBaseURI;
        emit BaseURIUpdated(newBaseURI);
    }
    
    /**
     * @notice Get the current base URI
     */
    function baseURI() external view returns (string memory) {
        return _baseTokenURI;
    }
    
    /**
     * @notice Pause all minting and transfers
     */
    function pause() external onlyRole(Roles.PAUSER_ROLE) {
        _pause();
    }
    
    /**
     * @notice Unpause all minting and transfers
     */
    function unpause() external onlyRole(Roles.PAUSER_ROLE) {
        _unpause();
    }
    
    // Required Overrides
    
    function tokenURI(uint256 tokenId)
        public
        view
        virtual
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        
        // Check if individual URI is set
        string memory individualURI = super.tokenURI(tokenId);
        if (bytes(individualURI).length > 0) {
            return individualURI;
        }
        
        // Otherwise, use Datakyte format
        return string(
            abi.encodePacked(
                _baseTokenURI,
                Strings.toHexString(uint160(address(this)), 20),
                "/",
                Strings.toString(tokenId),
                "/metadata"
            )
        );
    }
    
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage, ERC2981, IERC165)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}