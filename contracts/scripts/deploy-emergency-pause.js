const { ethers, network } = require("hardhat");

async function main() {
    console.log("🚨 Deploying Emergency Pause Controller...");
    
    const [deployer, guardian1, guardian2, guardian3, manager] = await ethers.getSigners();
    
    console.log("Deployer address:", deployer.address);
    console.log("Guardian addresses:", [guardian1.address, guardian2.address, guardian3.address]);
    console.log("Manager address:", manager.address);
    
    // Deploy EmergencyPauseController
    const EmergencyPauseController = await ethers.getContractFactory("EmergencyPauseController");
    const emergencyPause = await EmergencyPauseController.deploy(
        [guardian1.address, guardian2.address, guardian3.address], // Initial guardians
        manager.address // Contract manager
    );
    
    await emergencyPause.waitForDeployment();
    console.log("✅ EmergencyPauseController deployed to:", await emergencyPause.getAddress());
    
    // Get deployed contract addresses (these should be loaded from your deployment config)
    const contracts = {
        bogoToken: process.env.BOGO_TOKEN_ADDRESS,
        multisigTreasury: process.env.MULTISIG_TREASURY_ADDRESS,
        rewardDistributor: process.env.REWARD_DISTRIBUTOR_ADDRESS,
        commercialNFT: process.env.COMMERCIAL_NFT_ADDRESS,
        conservationNFT: process.env.CONSERVATION_NFT_ADDRESS
    };
    
    // Add all pausable contracts to the emergency controller
    console.log("\n📝 Adding pausable contracts...");
    
    if (contracts.bogoToken) {
        await emergencyPause.connect(manager).addContract(contracts.bogoToken, "BOGOTokenV2");
        console.log("✅ Added BOGOTokenV2");
    }
    
    if (contracts.multisigTreasury) {
        await emergencyPause.connect(manager).addContract(contracts.multisigTreasury, "MultisigTreasury");
        console.log("✅ Added MultisigTreasury");
    }
    
    if (contracts.rewardDistributor) {
        await emergencyPause.connect(manager).addContract(contracts.rewardDistributor, "BOGORewardDistributor");
        console.log("✅ Added BOGORewardDistributor");
    }
    
    if (contracts.commercialNFT) {
        await emergencyPause.connect(manager).addContract(contracts.commercialNFT, "CommercialNFT");
        console.log("✅ Added CommercialNFT");
    }
    
    if (contracts.conservationNFT) {
        await emergencyPause.connect(manager).addContract(contracts.conservationNFT, "ConservationNFT");
        console.log("✅ Added ConservationNFT");
    }
    
    // Grant PAUSER_ROLE to EmergencyPauseController in each contract
    console.log("\n🔐 Granting pause permissions...");
    
    // Note: This requires the deployer to have DEFAULT_ADMIN_ROLE in each contract
    // In production, this would be done through the MultisigTreasury
    
    const PAUSER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("PAUSER_ROLE"));
    const DEFAULT_ADMIN_ROLE = ethers.ZeroHash;
    
    // Example for BOGOTokenV2 (repeat for other contracts)
    if (contracts.bogoToken) {
        const bogoToken = await ethers.getContractAt("BOGOTokenV2", contracts.bogoToken);
        
        // Grant PAUSER_ROLE to emergency controller
        await bogoToken.grantRole(PAUSER_ROLE, await emergencyPause.getAddress());
        console.log("✅ Granted PAUSER_ROLE to EmergencyPauseController in BOGOTokenV2");
        
        // Optionally revoke PAUSER_ROLE from deployer for security
        // await bogoToken.revokeRole(PAUSER_ROLE, deployer.address);
    }
    
    console.log("\n🎉 Emergency Pause Controller deployment complete!");
    console.log("\nDeployment Summary:");
    console.log("- EmergencyPauseController:", await emergencyPause.getAddress());
    console.log("- Required confirmations:", await emergencyPause.requiredConfirmations());
    console.log("- Max pause duration:", await emergencyPause.MAX_PAUSE_DURATION(), "seconds");
    
    // Save deployment info
    const deploymentInfo = {
        network: network.name,
        emergencyPauseController: await emergencyPause.getAddress(),
        guardians: [guardian1.address, guardian2.address, guardian3.address],
        manager: manager.address,
        timestamp: new Date().toISOString()
    };
    
    const fs = require("fs");
    const path = require("path");
    
    // Create deployments directory if it doesn't exist
    const deploymentsDir = path.join(__dirname, "..", "deployments");
    if (!fs.existsSync(deploymentsDir)) {
        fs.mkdirSync(deploymentsDir);
    }
    
    fs.writeFileSync(
        path.join(deploymentsDir, `emergency-pause-${network.name}.json`),
        JSON.stringify(deploymentInfo, null, 2)
    );
    
    console.log("\n📄 Deployment info saved to deployments/emergency-pause-" + network.name + ".json");
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });