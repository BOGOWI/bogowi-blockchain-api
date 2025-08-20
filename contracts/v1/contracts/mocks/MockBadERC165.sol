// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title MockBadERC165
 * @notice Mock contract that fails ERC165 checks for testing
 */
contract MockBadERC165 {
    // This contract intentionally doesn't implement ERC165 properly
    
    // Fails the supportsInterface check
    function supportsInterface(bytes4) external pure returns (bool) {
        revert("Mock ERC165 failure");
    }
}

/**
 * @title MockBadERC721
 * @notice Mock contract that supports ERC165 but fails ERC721 check
 */
contract MockBadERC721 {
    // Supports ERC165
    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        if (interfaceId == 0x01ffc9a7) {
            return true; // ERC165
        }
        if (interfaceId == 0x80ac58cd) {
            revert("Mock ERC721 failure");
        }
        return false;
    }
}

/**
 * @title MockBadERC1155
 * @notice Mock contract that supports ERC165 but fails ERC1155 check
 */
contract MockBadERC1155 {
    // Supports ERC165
    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        if (interfaceId == 0x01ffc9a7) {
            return true; // ERC165
        }
        if (interfaceId == 0xd9b67a26) {
            revert("Mock ERC1155 failure");
        }
        return false;
    }
}