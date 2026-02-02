// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title SyrupToken
 * @dev ERC-20 token for the Waffle P2P LLM Token Marketplace
 * @notice This token is used as the primary currency in the Waffle ecosystem
 */
contract SyrupToken is ERC20, Ownable {
    /**
     * @dev Constructor that mints 1,000,000 SYRUP tokens to the deployer
     * @notice Tokens are minted with 18 decimals (standard ERC-20)
     */
    constructor() ERC20("SyrupToken", "SYRUP") Ownable(msg.sender) {
        // Mint 1,000,000 tokens to the deployer (with 18 decimals)
        _mint(msg.sender, 1_000_000 * 10 ** decimals());
    }

    /**
     * @dev Allows the owner to mint additional tokens
     * @param to The address to receive the minted tokens
     * @param amount The amount of tokens to mint (without decimals)
     */
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount * 10 ** decimals());
    }

    /**
     * @dev Allows token holders to burn their tokens
     * @param amount The amount of tokens to burn
     */
    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }

    /**
     * @dev Testnet faucet - anyone can claim 100 SYRUP once per address
     * @notice Only for testnet use! Remove in production.
     */
    mapping(address => bool) public hasClaimed;

    function faucet() external {
        require(!hasClaimed[msg.sender], "Already claimed");
        hasClaimed[msg.sender] = true;
        _mint(msg.sender, 100 * 10 ** decimals());
    }
}
