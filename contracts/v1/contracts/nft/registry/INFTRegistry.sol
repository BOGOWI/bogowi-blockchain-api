// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title INFTRegistry
 * @notice Interface for the central NFT registry that manages all NFT contract types
 * @dev Registry for BOGOWI NFT ecosystem - manages Tickets, Collectibles, and Badges
 */
interface INFTRegistry {
    enum ContractType {
        TICKET,      // Event tickets (ERC-721)
        COLLECTIBLE, // Digital collectibles (ERC-1155)
        BADGE        // Achievement badges (ERC-721)
    }

    struct ContractInfo {
        address contractAddress;
        ContractType contractType;
        string name;
        string version;
        bool isActive;
        uint256 registeredAt;
        address registeredBy;
    }

    event ContractRegistered(
        address indexed contractAddress,
        ContractType indexed contractType,
        string name,
        string version,
        address indexed registeredBy
    );

    event ContractUnregistered(
        address indexed contractAddress,
        address indexed unregisteredBy
    );

    event ContractStatusUpdated(
        address indexed contractAddress,
        bool isActive
    );

    event RegistryUpgraded(
        address indexed oldImplementation,
        address indexed newImplementation
    );

    /**
     * @notice Register a new NFT contract
     * @param contractAddress Address of the NFT contract to register
     * @param contractType Type of the contract (TICKET, COLLECTIBLE, BADGE)
     * @param name Human-readable name of the contract
     * @param version Version identifier of the contract
     */
    function registerContract(
        address contractAddress,
        ContractType contractType,
        string memory name,
        string memory version
    ) external;

    /**
     * @notice Unregister an NFT contract
     * @param contractAddress Address of the contract to unregister
     */
    function unregisterContract(address contractAddress) external;

    /**
     * @notice Update the active status of a registered contract
     * @param contractAddress Address of the contract
     * @param isActive New active status
     */
    function setContractStatus(address contractAddress, bool isActive) external;

    /**
     * @notice Get information about a registered contract
     * @param contractAddress Address of the contract
     * @return ContractInfo struct with contract details
     */
    function getContractInfo(address contractAddress) external view returns (ContractInfo memory);

    /**
     * @notice Get all contracts of a specific type
     * @param contractType Type of contracts to retrieve
     * @return Array of contract addresses
     */
    function getContractsByType(ContractType contractType) external view returns (address[] memory);

    /**
     * @notice Get all active contracts
     * @return Array of active contract addresses
     */
    function getActiveContracts() external view returns (address[] memory);

    /**
     * @notice Check if a contract is registered
     * @param contractAddress Address to check
     * @return True if registered, false otherwise
     */
    function isRegistered(address contractAddress) external view returns (bool);

    /**
     * @notice Check if a contract is active
     * @param contractAddress Address to check
     * @return True if active, false otherwise
     */
    function isActive(address contractAddress) external view returns (bool);

    /**
     * @notice Get the total number of registered contracts
     * @return Total count of registered contracts
     */
    function getContractCount() external view returns (uint256);
}