// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title NFTErrors
 * @notice Standardized error definitions for BOGOWI NFT contracts
 * @dev Used across all NFT contracts for consistent error handling
 */
library NFTErrors {
    // General errors
    error InvalidAddress(address addr);
    error InvalidAmount(uint256 amount);
    error InvalidTokenId(uint256 tokenId);
    error Unauthorized(address caller);
    error ContractPaused();
    
    // Registry errors
    error ContractAlreadyRegistered(address contractAddress);
    error ContractNotRegistered(address contractAddress);
    error ContractNotActive(address contractAddress);
    error InvalidContractType(uint8 contractType);
    
    // Ticket errors
    error TicketAlreadyRedeemed(uint256 tokenId);
    error TicketExpired(uint256 tokenId, uint256 expiredAt);
    error TicketTransferLocked(uint256 tokenId, uint256 unlockAt);
    error InvalidSignature(bytes signature);
    error InvalidBookingId(bytes32 bookingId);
    error InvalidEventId(bytes32 eventId);
    error NonceAlreadyUsed(uint256 nonce);
    error RedemptionWindowClosed(uint256 tokenId);
    error InvalidRedemptionTime(uint256 currentTime, uint256 eventStart, uint256 eventEnd);
    
    // Transfer errors
    error TransferNotAllowed(uint256 tokenId);
    error TransferToZeroAddress();
    error TransferFromZeroAddress();
    error TokenDoesNotExist(uint256 tokenId);
    error NotTokenOwner(address caller, uint256 tokenId);
    
    // Minting errors
    error MintingPaused();
    error ExceedsMaxSupply(uint256 requested, uint256 maxSupply);
    error InvalidMintParameters();
    error MintToZeroAddress();
    
    // Royalty errors
    error InvalidRoyaltyPercentage(uint256 percentage);
    error InvalidRoyaltyReceiver(address receiver);
    
    // Access control errors
    error MissingRole(bytes32 role, address account);
    error RoleAlreadyGranted(bytes32 role, address account);
    error CannotRevokeRole(bytes32 role, address account);
}