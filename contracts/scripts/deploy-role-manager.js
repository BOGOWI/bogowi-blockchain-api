const { ethers } = require("hardhat");

async function main() {
    console.log("Starting RoleManager deployment and migration...");
    
    const [deployer] = await ethers.getSigners();
    console.log("Deploying with account:", deployer.address);
    
    // 1. Deploy RoleManager
    console.log("\n1. Deploying RoleManager...");
    const RoleManager = await ethers.getContractFactory("RoleManager");
    const roleManager = await RoleManager.deploy();
    await roleManager.deployed();
    const roleManagerAddress = roleManager.address;
    console.log("RoleManager deployed to:", roleManagerAddress);
    
    // 2. Get role constants
    const DEFAULT_ADMIN_ROLE = await roleManager.DEFAULT_ADMIN_ROLE();
    const DAO_ROLE = await roleManager.DAO_ROLE();
    const BUSINESS_ROLE = await roleManager.BUSINESS_ROLE();
    const MINTER_ROLE = await roleManager.MINTER_ROLE();
    const PAUSER_ROLE = await roleManager.PAUSER_ROLE();
    const TREASURY_ROLE = await roleManager.TREASURY_ROLE();
    const DISTRIBUTOR_BACKEND_ROLE = await roleManager.DISTRIBUTOR_BACKEND_ROLE();
    
    console.log("\n2. Role Constants:");
    console.log("DEFAULT_ADMIN_ROLE:", DEFAULT_ADMIN_ROLE);
    console.log("DAO_ROLE:", DAO_ROLE);
    console.log("BUSINESS_ROLE:", BUSINESS_ROLE);
    console.log("MINTER_ROLE:", MINTER_ROLE);
    console.log("PAUSER_ROLE:", PAUSER_ROLE);
    console.log("TREASURY_ROLE:", TREASURY_ROLE);
    console.log("DISTRIBUTOR_BACKEND_ROLE:", DISTRIBUTOR_BACKEND_ROLE);
    
    // 3. Deploy updated contracts with RoleManager
    console.log("\n3. Deploying contracts with RoleManager integration...");
    
    // Deploy BOGOTokenV2_RoleManaged
    const BOGOTokenV2 = await ethers.getContractFactory("BOGOTokenV2_RoleManaged");
    const bogoToken = await BOGOTokenV2.deploy(
        roleManagerAddress,
        "BOGO Token",
        "BOGO"
    );
    await bogoToken.deployed();
    const bogoTokenAddress = bogoToken.address;
    console.log("BOGOTokenV2_RoleManaged deployed to:", bogoTokenAddress);
    
    // Deploy BOGORewardDistributor_RoleManaged
    const RewardDistributor = await ethers.getContractFactory("BOGORewardDistributor_RoleManaged");
    const rewardDistributor = await RewardDistributor.deploy(
        roleManagerAddress,
        bogoTokenAddress
    );
    await rewardDistributor.deployed();
    const rewardDistributorAddress = rewardDistributor.address;
    console.log("BOGORewardDistributor_RoleManaged deployed to:", rewardDistributorAddress);
    
    // 4. Register contracts with RoleManager
    console.log("\n4. Registering contracts with RoleManager...");
    
    await roleManager.registerContract(bogoTokenAddress, "BOGOTokenV2");
    console.log("Registered BOGOTokenV2");
    
    await roleManager.registerContract(rewardDistributorAddress, "BOGORewardDistributor");
    console.log("Registered BOGORewardDistributor");
    
    // 5. Configure initial roles (example)
    console.log("\n5. Configuring initial roles...");
    
    // Example role assignments - adjust these based on your needs
    const roleAssignments = {
        // Treasury multisig address (update with actual address)
        treasury: "0x0000000000000000000000000000000000000000",
        // DAO multisig address (update with actual address)
        dao: "0x0000000000000000000000000000000000000000",
        // Business operations address (update with actual address)
        business: "0x0000000000000000000000000000000000000000",
        // Backend service address (update with actual address)
        backend: "0x0000000000000000000000000000000000000000",
    };
    
    // Only assign roles if addresses are set
    if (roleAssignments.treasury !== ethers.ZeroAddress) {
        await roleManager.grantRole(TREASURY_ROLE, roleAssignments.treasury);
        console.log("Granted TREASURY_ROLE to:", roleAssignments.treasury);
    }
    
    if (roleAssignments.dao !== ethers.ZeroAddress) {
        await roleManager.grantRole(DAO_ROLE, roleAssignments.dao);
        console.log("Granted DAO_ROLE to:", roleAssignments.dao);
    }
    
    if (roleAssignments.business !== ethers.ZeroAddress) {
        await roleManager.grantRole(BUSINESS_ROLE, roleAssignments.business);
        console.log("Granted BUSINESS_ROLE to:", roleAssignments.business);
    }
    
    if (roleAssignments.backend !== ethers.ZeroAddress) {
        await roleManager.grantRole(DISTRIBUTOR_BACKEND_ROLE, roleAssignments.backend);
        console.log("Granted DISTRIBUTOR_BACKEND_ROLE to:", roleAssignments.backend);
    }
    
    // Grant deployer some initial roles for testing
    await roleManager.grantRole(MINTER_ROLE, deployer.address);
    await roleManager.grantRole(PAUSER_ROLE, deployer.address);
    console.log("Granted MINTER_ROLE and PAUSER_ROLE to deployer:", deployer.address);
    
    // 6. Save deployment addresses
    console.log("\n6. Deployment Summary:");
    console.log("========================");
    console.log("RoleManager:", roleManagerAddress);
    console.log("BOGOTokenV2:", bogoTokenAddress);
    console.log("RewardDistributor:", rewardDistributorAddress);
    console.log("========================");
    
    // Save to file
    const fs = require("fs");
    const deploymentInfo = {
        network: network.name,
        timestamp: new Date().toISOString(),
        contracts: {
            RoleManager: roleManagerAddress,
            BOGOTokenV2_RoleManaged: bogoTokenAddress,
            BOGORewardDistributor_RoleManaged: rewardDistributorAddress
        },
        roles: {
            DEFAULT_ADMIN_ROLE,
            DAO_ROLE,
            BUSINESS_ROLE,
            MINTER_ROLE,
            PAUSER_ROLE,
            TREASURY_ROLE,
            DISTRIBUTOR_BACKEND_ROLE
        }
    };
    
    fs.writeFileSync(
        `deployments/rolemanager-${network.name}-${Date.now()}.json`,
        JSON.stringify(deploymentInfo, null, 2)
    );
    
    console.log("\nDeployment complete! Check deployments/ folder for details.");
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });