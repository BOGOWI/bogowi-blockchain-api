// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./INFTRegistry.sol";
import "../../base/RoleManaged.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/introspection/IERC165.sol";

/**
 * @title NFTRegistry
 * @notice Central registry for all BOGOWI NFT contracts
 * @dev Manages registration and tracking of Tickets, Collectibles, and Badges
 * @dev Future versions will implement UUPS proxy pattern for upgradeability
 */
contract NFTRegistry is 
    INFTRegistry, 
    Pausable,
    ReentrancyGuard,
    RoleManaged 
{
    bytes32 public constant REGISTRY_ADMIN_ROLE = keccak256("REGISTRY_ADMIN_ROLE");
    bytes32 public constant CONTRACT_DEPLOYER_ROLE = keccak256("CONTRACT_DEPLOYER_ROLE");
    
    uint256 public constant MAX_CONTRACTS = 1000; // Prevent unbounded growth
    uint256 public constant MAX_CONTRACTS_PER_TYPE = 500; // Limit per type
    
    mapping(address => ContractInfo) private _contracts;
    mapping(ContractType => address[]) private _contractsByType;
    address[] private _allContracts;
    mapping(address => uint256) private _contractIndex;
    
    uint256 private _totalContracts;

    /**
     * @notice Constructor to set up the registry
     * @param _roleManager Address of the RoleManager contract
     */
    constructor(address _roleManager) RoleManaged(_roleManager) {
        // Constructor body can be empty as RoleManaged handles initialization
    }

    /**
     * @notice Register a new NFT contract
     * @dev Only callable by CONTRACT_DEPLOYER_ROLE
     */
    function registerContract(
        address contractAddress,
        ContractType contractType,
        string memory name,
        string memory version
    ) external override onlyRole(CONTRACT_DEPLOYER_ROLE) whenNotPaused nonReentrant {
        require(contractAddress != address(0), "Invalid contract address");
        require(!isRegistered(contractAddress), "Contract already registered");
        require(bytes(name).length > 0, "Name cannot be empty");
        require(bytes(version).length > 0, "Version cannot be empty");
        
        // Check array size limits
        require(_totalContracts < MAX_CONTRACTS, "Maximum contracts limit reached");
        require(_contractsByType[contractType].length < MAX_CONTRACTS_PER_TYPE, "Maximum contracts per type reached");
        
        // Validate contract interface based on type
        _validateContractInterface(contractAddress, contractType);

        ContractInfo memory info = ContractInfo({
            contractAddress: contractAddress,
            contractType: contractType,
            name: name,
            version: version,
            isActive: true,
            registeredAt: block.timestamp,
            registeredBy: msg.sender
        });

        _contracts[contractAddress] = info;
        _contractsByType[contractType].push(contractAddress);
        _allContracts.push(contractAddress);
        _contractIndex[contractAddress] = _allContracts.length - 1;
        _totalContracts++;

        emit ContractRegistered(
            contractAddress,
            contractType,
            name,
            version,
            msg.sender
        );
    }

    /**
     * @notice Unregister an NFT contract
     * @dev Only callable by REGISTRY_ADMIN_ROLE
     */
    function unregisterContract(
        address contractAddress
    ) external override onlyRole(REGISTRY_ADMIN_ROLE) nonReentrant {
        require(isRegistered(contractAddress), "Contract not registered");

        ContractInfo memory info = _contracts[contractAddress];
        
        // Store array lengths before modifications
        uint256 allContractsLength = _allContracts.length;
        require(allContractsLength > 0, "No contracts to remove");
        
        // Remove from type array first (before state changes)
        _removeFromTypeArray(info.contractType, contractAddress);
        
        // Remove from all contracts array
        uint256 index = _contractIndex[contractAddress];
        uint256 lastIndex = allContractsLength - 1;
        
        // Prevent out-of-bounds access
        require(index < allContractsLength, "Invalid contract index");
        
        if (index != lastIndex) {
            address lastContract = _allContracts[lastIndex];
            _allContracts[index] = lastContract;
            _contractIndex[lastContract] = index;
        }
        
        _allContracts.pop();
        
        // Clear all mappings atomically
        delete _contractIndex[contractAddress];
        delete _contracts[contractAddress];
        
        // Decrement counter last
        _totalContracts--;

        emit ContractUnregistered(contractAddress, msg.sender);
    }

    /**
     * @notice Update the active status of a registered contract
     * @dev Only callable by REGISTRY_ADMIN_ROLE
     */
    function setContractStatus(
        address contractAddress,
        bool _isActive
    ) external override onlyRole(REGISTRY_ADMIN_ROLE) nonReentrant {
        require(isRegistered(contractAddress), "Contract not registered");
        
        _contracts[contractAddress].isActive = _isActive;
        
        emit ContractStatusUpdated(contractAddress, _isActive);
    }

    /**
     * @notice Get information about a registered contract
     */
    function getContractInfo(
        address contractAddress
    ) external view override returns (ContractInfo memory) {
        require(isRegistered(contractAddress), "Contract not registered");
        return _contracts[contractAddress];
    }

    /**
     * @notice Get all contracts of a specific type
     */
    function getContractsByType(
        ContractType contractType
    ) external view override returns (address[] memory) {
        return _contractsByType[contractType];
    }

    /**
     * @notice Get all active contracts
     * @dev WARNING: This can be gas-intensive with many contracts. Consider using getActiveContractsPaginated instead
     */
    function getActiveContracts() external view override returns (address[] memory) {
        uint256 activeCount = 0;
        uint256 length = _allContracts.length;
        
        // Count active contracts
        for (uint256 i = 0; i < length; i++) {
            if (_contracts[_allContracts[i]].isActive) {
                activeCount++;
            }
        }
        
        // Build array of active contracts
        address[] memory activeContracts = new address[](activeCount);
        uint256 currentIndex = 0;
        
        for (uint256 i = 0; i < length; i++) {
            if (_contracts[_allContracts[i]].isActive) {
                activeContracts[currentIndex] = _allContracts[i];
                currentIndex++;
            }
        }
        
        return activeContracts;
    }
    
    /**
     * @notice Get active contracts with pagination
     * @param offset Starting index
     * @param limit Maximum number of contracts to return (max: 100)
     * @return contracts Array of active contract addresses
     * @return hasMore Whether there are more contracts after this batch
     */
    function getActiveContractsPaginated(
        uint256 offset,
        uint256 limit
    ) external view returns (address[] memory contracts, bool hasMore) {
        // Bounds validation
        require(limit > 0 && limit <= 100, "Limit must be between 1 and 100");
        
        uint256 length = _allContracts.length;
        require(offset < length || length == 0, "Offset out of bounds");
        
        uint256 remaining = length > offset ? length - offset : 0;
        uint256 returnSize = remaining > limit ? limit : remaining;
        
        contracts = new address[](returnSize);
        uint256 currentIndex = 0;
        
        for (uint256 i = offset; i < offset + returnSize && i < length; i++) {
            if (_contracts[_allContracts[i]].isActive) {
                contracts[currentIndex] = _allContracts[i];
                currentIndex++;
            }
        }
        
        // Resize array if needed - using safe approach without assembly
        if (currentIndex < returnSize) {
            // Create properly sized array
            address[] memory resizedContracts = new address[](currentIndex);
            for (uint256 j = 0; j < currentIndex; j++) {
                resizedContracts[j] = contracts[j];
            }
            contracts = resizedContracts;
        }
        
        hasMore = (offset + returnSize) < length;
        return (contracts, hasMore);
    }

    /**
     * @notice Check if a contract is registered
     */
    function isRegistered(address contractAddress) public view override returns (bool) {
        return _contracts[contractAddress].registeredAt > 0;
    }

    /**
     * @notice Check if a contract is active
     */
    function isActive(address contractAddress) external view override returns (bool) {
        return isRegistered(contractAddress) && _contracts[contractAddress].isActive;
    }

    /**
     * @notice Get the total number of registered contracts
     */
    function getContractCount() external view override returns (uint256) {
        return _totalContracts;
    }

    /**
     * @notice Pause the registry
     * @dev Only callable by REGISTRY_ADMIN_ROLE
     */
    function pause() external onlyRole(REGISTRY_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @notice Unpause the registry
     * @dev Only callable by REGISTRY_ADMIN_ROLE
     */
    function unpause() external onlyRole(REGISTRY_ADMIN_ROLE) {
        _unpause();
    }


    /**
     * @notice Validate contract interface based on type
     * @dev Checks for proper ERC-721, ERC-1155, or custom badge interface
     */
    function _validateContractInterface(
        address contractAddress,
        ContractType contractType
    ) private view {
        // First check ERC165 support
        try IERC165(contractAddress).supportsInterface(0x01ffc9a7) returns (bool supportsERC165) {
            require(supportsERC165, "Contract must support ERC165");
        } catch {
            revert("ERC165 check failed");
        }
        
        // Check specific interfaces based on contract type
        if (contractType == ContractType.TICKET || contractType == ContractType.BADGE) {
            // ERC721 interface ID: 0x80ac58cd
            try IERC165(contractAddress).supportsInterface(0x80ac58cd) returns (bool supportsERC721) {
                require(supportsERC721, "Ticket/Badge must support ERC721");
            } catch {
                revert("ERC721 interface check failed");
            }
        } else if (contractType == ContractType.COLLECTIBLE) {
            // ERC1155 interface ID: 0xd9b67a26
            try IERC165(contractAddress).supportsInterface(0xd9b67a26) returns (bool supportsERC1155) {
                require(supportsERC1155, "Collectible must support ERC1155");
            } catch {
                revert("ERC1155 interface check failed");
            }
        }
    }
    
    /**
     * @notice Remove contract from type array
     * @dev Internal helper function with bounds checking
     */
    function _removeFromTypeArray(
        ContractType contractType,
        address contractAddress
    ) private {
        address[] storage typeArray = _contractsByType[contractType];
        uint256 arrayLength = typeArray.length;
        
        // Early return if array is empty
        if (arrayLength == 0) {
            return;
        }
        
        // Find and remove the contract
        for (uint256 i = 0; i < arrayLength; i++) {
            if (typeArray[i] == contractAddress) {
                // Special case: removing last element
                if (i == arrayLength - 1) {
                    typeArray.pop();
                } else {
                    // Move last element to this position and pop
                    typeArray[i] = typeArray[arrayLength - 1];
                    typeArray.pop();
                }
                break;
            }
        }
    }
}