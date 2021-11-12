// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "https://github.com/OpenZeppelin/openzeppelin-contracts/contracts/token/ERC721/ERC721.sol";
import "https://github.com/OpenZeppelin/openzeppelin-contracts/contracts/utils/Counters.sol";
import "https://github.com/OpenZeppelin/openzeppelin-contracts/contracts/access/Ownable.sol";

contract ColorNFT is ERC721, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIDs;

    uint256 public constant TOKEN_PRICE = 50_000_000_000_000_000;
    uint public constant MAX_TOKEN_PURCHASE = 20;
    uint256 public constant MAX_TOKENS = 100;

    string public baseURI;
    string public provenanceHash;
    uint256 public startIndex;
    uint256 public startIndexBlockNumber;
    uint256 public revealTimestamp;

    constructor() ERC721("ColorNFT", "CLR") {
        revealTimestamp = block.timestamp + 1_209_600; // 14 days
    }

    function setBaseURI(string memory newBaseURI) public onlyOwner {
        baseURI = newBaseURI;
    }

    function _baseURI() internal view override returns (string memory) {
        return baseURI;
    }

    function setProvenanceHash(string memory newProvenanceHash) public onlyOwner {
        provenanceHash = newProvenanceHash;
    }

    function setRevealTimestamp(uint256 newRevealTimestamp) public onlyOwner {
        revealTimestamp = newRevealTimestamp;
    }

    function withdraw() public onlyOwner {
        uint balance = address(this).balance;
        payable(msg.sender).transfer(balance);
    }

    function ownerMint(uint numTokens) public onlyOwner {
        _freeMint(msg.sender, numTokens);
    }

    function mint(uint numTokens) public payable {
        require(numTokens <= MAX_TOKEN_PURCHASE, "Number of tokens to mint cannot exceed MAX_TOKEN_PURCHASE");
        require(msg.value >= TOKEN_PRICE * numTokens, "Insufficient ether value for mint");

        _freeMint(msg.sender, numTokens);
    }

    /**
     * Mints the requested number of tokens to the given address. Checks that minting would not exceed maximum supply.
     */
    function _freeMint(address to, uint numTokens) private {
        require(_tokenIDs.current() + numTokens <= MAX_TOKENS, "Not enough tokens left");

        for (uint i = 0; i < numTokens; i++) {
            _safeMint(to, _tokenIDs.current());
            _tokenIDs.increment();
        }

        if (startIndexBlockNumber == 0 && (_tokenIDs.current() >= MAX_TOKENS || block.timestamp >= revealTimestamp)) {
            startIndexBlockNumber = block.number;
        }
    }

    function setStartIndex() public {
        require(startIndex == 0, "Start index is already set");
        require(startIndexBlockNumber != 0, "Start index block must be set");

        // make sure contract has access to the block hash (contract only has access to the last 256 block hashes)
        if (block.number - startIndexBlockNumber < 256) {
            startIndex = uint(blockhash(startIndexBlockNumber)) % MAX_TOKENS;
            return;
        }

        startIndex = uint(blockhash(block.number - 1)) % MAX_TOKENS;
    }
}