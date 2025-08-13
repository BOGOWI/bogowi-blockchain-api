// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

/**
 * @title MockERC721
 * @notice Mock ERC721 contract for testing NFTRegistry
 */
contract MockERC721 is ERC721 {
    uint256 private _tokenIdCounter;
    
    constructor() ERC721("MockERC721", "MOCK") {}
    
    function mint(address to) public returns (uint256) {
        uint256 tokenId = _tokenIdCounter++;
        _mint(to, tokenId);
        return tokenId;
    }
}