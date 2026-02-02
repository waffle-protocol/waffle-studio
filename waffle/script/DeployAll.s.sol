// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../contracts/SyrupToken.sol";
import "../contracts/BakeRegistry.sol";

/**
 * @title DeployAll
 * @notice Deploys all Waffle contracts to the network
 * @dev Run with: forge script script/DeployAll.s.sol --rpc-url anvil --broadcast
 */
contract DeployAll is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerPrivateKey);

        // 1. Deploy SyrupToken
        SyrupToken syrupToken = new SyrupToken();
        console.log("===========================================");
        console.log("SyrupToken deployed at:", address(syrupToken));

        // 2. Deploy BakeRegistry with SyrupToken address
        BakeRegistry bakeRegistry = new BakeRegistry(address(syrupToken));
        console.log("BakeRegistry deployed at:", address(bakeRegistry));
        console.log("===========================================");

        // Log summary
        console.log("");
        console.log("DEPLOYMENT COMPLETE");
        console.log("-------------------");
        console.log("SYRUP Token:", address(syrupToken));
        console.log("Bake Registry:", address(bakeRegistry));
        console.log("Total Supply:", syrupToken.totalSupply() / 1e18, "SYRUP");
        console.log("===========================================");

        vm.stopBroadcast();
    }
}
