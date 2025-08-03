// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title StandardErrors
 * @dev Standardized error messages for BOGOWI contracts
 */
contract StandardErrors {
    // Access Control Errors
    string constant UNAUTHORIZED = "UNAUTHORIZED";
    string constant NOT_OWNER = "NOT_OWNER";
    string constant NOT_ADMIN = "NOT_ADMIN";
    string constant NOT_SIGNER = "NOT_SIGNER";
    string constant NOT_TREASURY = "NOT_TREASURY";
    string constant NOT_BACKEND = "NOT_BACKEND";
    string constant NOT_MULTISIG = "NOT_MULTISIG";
    
    // Input Validation Errors
    string constant ZERO_ADDRESS = "ZERO_ADDRESS";
    string constant ZERO_AMOUNT = "ZERO_AMOUNT";
    string constant INVALID_ADDRESS = "INVALID_ADDRESS";
    string constant INVALID_AMOUNT = "INVALID_AMOUNT";
    string constant INVALID_PARAMETER = "INVALID_PARAMETER";
    string constant INVALID_LENGTH = "INVALID_LENGTH";
    string constant EMPTY_STRING = "EMPTY_STRING";
    
    // State Errors
    string constant ALREADY_EXISTS = "ALREADY_EXISTS";
    string constant DOES_NOT_EXIST = "DOES_NOT_EXIST";
    string constant ALREADY_INITIALIZED = "ALREADY_INITIALIZED";
    string constant NOT_INITIALIZED = "NOT_INITIALIZED";
    string constant PAUSED = "PAUSED";
    string constant NOT_PAUSED = "NOT_PAUSED";
    string constant EXPIRED = "EXPIRED";
    string constant NOT_EXPIRED = "NOT_EXPIRED";
    string constant ACTIVE = "ACTIVE";
    string constant INACTIVE = "INACTIVE";
    
    // Limit Errors
    string constant EXCEEDS_LIMIT = "EXCEEDS_LIMIT";
    string constant EXCEEDS_SUPPLY = "EXCEEDS_SUPPLY";
    string constant EXCEEDS_ALLOCATION = "EXCEEDS_ALLOCATION";
    string constant EXCEEDS_BALANCE = "EXCEEDS_BALANCE";
    string constant DAILY_LIMIT_EXCEEDED = "DAILY_LIMIT_EXCEEDED";
    string constant MAX_REACHED = "MAX_REACHED";
    string constant INSUFFICIENT_BALANCE = "INSUFFICIENT_BALANCE";
    
    // Operation Errors
    string constant TRANSFER_FAILED = "TRANSFER_FAILED";
    string constant OPERATION_FAILED = "OPERATION_FAILED";
    string constant NOT_READY = "NOT_READY";
    string constant ALREADY_PROCESSED = "ALREADY_PROCESSED";
    string constant COOLDOWN_ACTIVE = "COOLDOWN_ACTIVE";
    string constant NOT_WHITELISTED = "NOT_WHITELISTED";
    
    // Logic Errors
    string constant CIRCULAR_REFERENCE = "CIRCULAR_REFERENCE";
    string constant SELF_REFERENCE = "SELF_REFERENCE";
    string constant INVALID_STATE = "INVALID_STATE";
    string constant CONDITIONS_NOT_MET = "CONDITIONS_NOT_MET";
}

/**
 * @dev Custom errors for gas optimization
 */
interface IStandardErrors {
    // Access Control
    error Unauthorized();
    error NotOwner();
    error NotAdmin();
    error NotSigner();
    error NotTreasury();
    error NotBackend();
    error NotMultisig();
    
    // Input Validation
    error ZeroAddress();
    error ZeroAmount();
    error InvalidAddress();
    error InvalidAmount();
    error InvalidParameter();
    error InvalidLength();
    error EmptyString();
    
    // State
    error AlreadyExists();
    error DoesNotExist();
    error AlreadyInitialized();
    error NotInitialized();
    error ContractPaused();
    error ContractNotPaused();
    error Expired();
    error NotExpired();
    error Active();
    error Inactive();
    
    // Limits
    error ExceedsLimit();
    error ExceedsSupply();
    error ExceedsAllocation();
    error ExceedsBalance();
    error DailyLimitExceeded();
    error MaxReached();
    error InsufficientBalance();
    
    // Operations
    error TransferFailed();
    error OperationFailed();
    error NotReady();
    error AlreadyProcessed();
    error CooldownActive();
    error NotWhitelisted();
    
    // Logic
    error CircularReference();
    error SelfReference();
    error InvalidState();
    error ConditionsNotMet();
}