// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

/**
 * @title IBOGOWITickets
 * @notice Interface for BOGOWI Tickets (BWTIX) - Phase 1 NFT ticketing system
 * @dev ERC-721 compliant event tickets with time-locked transfers and QR redemption
 */
interface IBOGOWITickets is IERC721 {
    
    /**
     * @notice Ticket lifecycle states
     */
    enum TicketState {
        ISSUED,    // Default state when minted
        REDEEMED,  // Terminal state after redemption
        EXPIRED    // Terminal state after expiry
    }
    
    /**
     * @notice On-chain ticket data structure - optimized for gas efficiency
     * @dev Packed to minimize storage slots (3 slots total)
     * Slot 1: bookingId (32 bytes)
     * Slot 2: eventId (32 bytes)  
     * Slot 3: transferUnlockAt (8) + expiresAt (8) + utilityFlags (4) + state (1) + flags (1) = 22 bytes
     */
    struct TicketData {
        bytes32 bookingId;               // Slot 1: Hash of booking UUID
        bytes32 eventId;                 // Slot 2: Hash of event identifier
        uint64 transferUnlockAt;         // Slot 3: Unix timestamp when transfers unlock
        uint64 expiresAt;                // Slot 3: Unix timestamp when ticket expires
        uint32 utilityFlags;             // Slot 3: Bitmask for perks/utilities
        TicketState state;               // Slot 3: Current ticket state (uint8)
        bool nonTransferableAfterRedeem; // Slot 3: Lock transfers after redemption
        bool burnOnRedeem;               // Slot 3: Burn token upon redemption
    }
    
    /**
     * @notice Ticket metadata for minting
     */
    struct MintParams {
        address to;                      // Recipient address
        bytes32 bookingId;               // Booking identifier hash
        bytes32 eventId;                 // Event identifier hash
        uint32 utilityFlags;             // Perks bitmask
        uint64 transferUnlockAt;         // Transfer unlock timestamp
        uint64 expiresAt;                // Expiry timestamp
        string metadataURI;              // IPFS URI for metadata
        uint256 rewardBasisPoints;       // BOGO reward basis points
    }
    
    /**
     * @notice Redemption signature data for EIP-712
     */
    struct RedemptionData {
        uint256 tokenId;
        address redeemer;
        uint256 nonce;
        uint256 deadline;
        uint256 chainId;  // Added to prevent cross-chain replay
        bytes signature;
    }
    
    // Events
    event TicketMinted(
        uint256 indexed tokenId,
        bytes32 indexed bookingIdHash,
        bytes32 indexed eventIdHash,
        address buyer,
        uint256 rewardBasisPoints
    );
    
    event TicketRedeemed(
        uint256 indexed tokenId,
        address indexed redeemedBy,
        uint256 timestamp
    );
    
    event TicketExpired(
        uint256 indexed tokenId
    );
    
    event TransferUnlockUpdated(
        uint256 indexed tokenId,
        uint64 newUnlockTime
    );
    
    event RoyaltyInfoUpdated(
        address indexed receiver,
        uint96 feeBasisPoints
    );
    
    event BatchMintStarted(
        uint256 batchSize,
        address indexed minter
    );
    
    event NonceUsed(
        uint256 indexed nonce,
        address indexed user
    );
    
    event TicketBurned(
        uint256 indexed tokenId,
        address indexed owner
    );
    
    event BaseURIUpdated(
        string newBaseURI
    );
    
    event ExpiryGracePeriodUpdated(
        uint256 oldGracePeriod,
        uint256 newGracePeriod
    );
    
    event DatakyteMetadataLinked(
        uint256 indexed tokenId,
        string datakyteNftId
    );
    
    // Core Functions
    
    /**
     * @notice Mint a new ticket NFT
     * @dev Only callable by MINTER_ROLE
     * @param params Minting parameters including recipient and ticket data
     * @return tokenId The ID of the newly minted ticket
     */
    function mintTicket(MintParams memory params) external returns (uint256 tokenId);
    
    /**
     * @notice Mint multiple tickets in batch
     * @dev Only callable by MINTER_ROLE
     * @param params Array of minting parameters
     * @return tokenIds Array of newly minted token IDs
     */
    function mintBatch(MintParams[] memory params) external returns (uint256[] memory tokenIds);
    
    /**
     * @notice Redeem a ticket using QR code and signature
     * @param redemptionData Signed redemption data including nonce
     */
    function redeemTicket(RedemptionData memory redemptionData) external;
    
    /**
     * @notice Update transfer unlock time for a ticket
     * @dev Only callable by admin
     * @param tokenId Token to update
     * @param newUnlockTime New unlock timestamp
     */
    function updateTransferUnlock(uint256 tokenId, uint64 newUnlockTime) external;
    
    /**
     * @notice Mark a ticket as expired
     * @param tokenId Token to expire
     */
    function expireTicket(uint256 tokenId) external;
    
    // View Functions
    
    /**
     * @notice Get ticket data for a token
     * @param tokenId Token to query
     * @return Ticket data structure
     */
    function getTicketData(uint256 tokenId) external view returns (TicketData memory);
    
    /**
     * @notice Check if a ticket can be transferred
     * @param tokenId Token to check
     * @return bool True if transferable
     */
    function isTransferable(uint256 tokenId) external view returns (bool);
    
    /**
     * @notice Check if a ticket is redeemed
     * @param tokenId Token to check
     * @return bool True if redeemed
     */
    function isRedeemed(uint256 tokenId) external view returns (bool);
    
    /**
     * @notice Check if a ticket is expired
     * @param tokenId Token to check
     * @return bool True if expired
     */
    function isExpired(uint256 tokenId) external view returns (bool);
    
    /**
     * @notice Verify a redemption signature
     * @param redemptionData Data to verify
     * @return bool True if signature is valid
     */
    function verifyRedemptionSignature(RedemptionData memory redemptionData) external view returns (bool);
}