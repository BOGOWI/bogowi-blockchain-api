const { ethers } = require("hardhat");

/**
 * Script to update pause permissions on existing contracts
 * This grants the EmergencyPauseController the ability to pause contracts
 * Must be executed by an account with DEFAULT_ADMIN_ROLE
 */
async function main() {
    console.log("🔐 Updating pause permissions for EmergencyPauseController...");
    
    const [executor] = await ethers.getSigners();
    console.log("Executing with account:", executor.address);
    
    // Load addresses from environment or config
    const EMERGENCY_PAUSE_ADDRESS = process.env.EMERGENCY_PAUSE_ADDRESS;
    const PAUSER_ROLE = ethers.keccak256(ethers.toUtf8Bytes("PAUSER_ROLE"));
    
    if (!EMERGENCY_PAUSE_ADDRESS) {
        throw new Error("EMERGENCY_PAUSE_ADDRESS not set in environment");
    }
    
    const contracts = [
        {
            name: "BOGOTokenV2",
            address: process.env.BOGO_TOKEN_ADDRESS,
            abi: "BOGOTokenV2"
        },
        {
            name: "MultisigTreasury", 
            address: process.env.MULTISIG_TREASURY_ADDRESS,
            abi: "MultisigTreasury"
        },
        {
            name: "BOGORewardDistributor",
            address: process.env.REWARD_DISTRIBUTOR_ADDRESS,
            abi: "BOGORewardDistributor"
        },
        {
            name: "CommercialNFT",
            address: process.env.COMMERCIAL_NFT_ADDRESS,
            abi: "CommercialNFT"
        },
        {
            name: "ConservationNFT",
            address: process.env.CONSERVATION_NFT_ADDRESS,
            abi: "ConservationNFT"
        }
    ];
    
    for (const contractInfo of contracts) {
        if (!contractInfo.address) {
            console.log(`⚠️  Skipping ${contractInfo.name} - address not provided`);
            continue;
        }
        
        try {
            console.log(`\n📝 Processing ${contractInfo.name}...`);
            const contract = await ethers.getContractAt(contractInfo.abi, contractInfo.address);
            
            // Check if EmergencyPauseController already has PAUSER_ROLE
            const hasPauserRole = await contract.hasRole(PAUSER_ROLE, EMERGENCY_PAUSE_ADDRESS);
            
            if (hasPauserRole) {
                console.log(`✅ ${contractInfo.name} already has PAUSER_ROLE granted to EmergencyPauseController`);
            } else {
                // Grant PAUSER_ROLE
                console.log(`🔄 Granting PAUSER_ROLE to EmergencyPauseController...`);
                const tx = await contract.grantRole(PAUSER_ROLE, EMERGENCY_PAUSE_ADDRESS);
                await tx.wait();
                console.log(`✅ PAUSER_ROLE granted to EmergencyPauseController`);
            }
            
            // Verify the role was granted
            const verified = await contract.hasRole(PAUSER_ROLE, EMERGENCY_PAUSE_ADDRESS);
            console.log(`✓ Verification: PAUSER_ROLE granted = ${verified}`);
            
        } catch (error) {
            console.error(`❌ Error processing ${contractInfo.name}:`, error.message);
        }
    }
    
    console.log("\n🎉 Permission update complete!");
    
    // Display summary
    console.log("\n📊 Summary:");
    console.log("EmergencyPauseController:", EMERGENCY_PAUSE_ADDRESS);
    console.log("Can now pause/unpause all BOGOWI contracts in case of emergency");
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });