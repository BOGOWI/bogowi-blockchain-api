const { ethers } = require("hardhat");
const fs = require("fs");

// Configuration
const CONFIG = {
    REGISTRY_ADDRESS: process.env.REGISTRY_ADDRESS || "",
    TOKEN_ADDRESS: process.env.TOKEN_ADDRESS || "",
    TREASURY_ADDRESS: process.env.TREASURY_ADDRESS || "",
    MIGRATION_HELPER_ADDRESS: process.env.MIGRATION_HELPER_ADDRESS || "",
    BATCH_SIZE: 50, // Process users in batches
    GAS_LIMIT: 5000000,
};

async function main() {
    console.log("=== Reward Distributor Migration Script ===\n");
    
    const [deployer] = await ethers.getSigners();
    console.log("Executing with account:", deployer.address);
    console.log("Account balance:", ethers.utils.formatEther(await deployer.getBalance()), "ETH\n");

    // Validate configuration
    if (!CONFIG.REGISTRY_ADDRESS || !CONFIG.TOKEN_ADDRESS || !CONFIG.TREASURY_ADDRESS) {
        throw new Error("Missing required configuration. Set environment variables.");
    }

    try {
        // Step 1: Connect to contracts
        console.log("Step 1: Connecting to existing contracts...");
        
        const registry = await ethers.getContractAt("ContractRegistry", CONFIG.REGISTRY_ADDRESS);
        const token = await ethers.getContractAt("BOGOTokenV2", CONFIG.TOKEN_ADDRESS);
        const migrationHelper = CONFIG.MIGRATION_HELPER_ADDRESS 
            ? await ethers.getContractAt("MigrationHelper", CONFIG.MIGRATION_HELPER_ADDRESS)
            : null;
        
        // Get current distributor
        const oldDistributorAddress = await registry.getContract("RewardDistributor");
        const oldDistributor = await ethers.getContractAt("BOGORewardDistributor", oldDistributorAddress);
        console.log("Current RewardDistributor:", oldDistributorAddress);
        
        // Step 2: Analyze current state
        console.log("\nStep 2: Analyzing current distributor state...");
        
        const oldBalance = await token.balanceOf(oldDistributorAddress);
        console.log("Token balance in old distributor:", ethers.utils.formatEther(oldBalance), "BOGO");
        
        // Get current configuration
        const isPaused = await oldDistributor.paused();
        console.log("Is paused:", isPaused);
        
        // Step 3: Deploy new distributor
        console.log("\nStep 3: Deploying new RewardDistributor...");
        
        const RewardDistributorV2 = await ethers.getContractFactory("BOGORewardDistributor");
        const newDistributor = await RewardDistributorV2.deploy(
            CONFIG.TOKEN_ADDRESS,
            CONFIG.TREASURY_ADDRESS,
            { gasLimit: CONFIG.GAS_LIMIT }
        );
        await newDistributor.deployed();
        console.log("New RewardDistributor deployed to:", newDistributor.address);
        
        // Step 4: Pause old distributor
        console.log("\nStep 4: Pausing old distributor...");
        
        if (!isPaused) {
            const pauseTx = await oldDistributor.pause({ gasLimit: 100000 });
            await pauseTx.wait();
            console.log("Old distributor paused");
        } else {
            console.log("Old distributor already paused");
        }
        
        // Step 5: Transfer tokens
        console.log("\nStep 5: Transferring tokens to new distributor...");
        
        if (oldBalance.gt(0)) {
            console.log("Initiating treasury sweep...");
            const sweepTx = await oldDistributor.treasurySweep(
                token.address,
                newDistributor.address,
                oldBalance,
                { gasLimit: 150000 }
            );
            await sweepTx.wait();
            console.log("Tokens transferred successfully");
        } else {
            console.log("No tokens to transfer");
        }
        
        // Step 6: Configure new distributor
        console.log("\nStep 6: Configuring new distributor...");
        
        // Copy authorized backends
        // Note: This would need to be done manually or through events analysis
        // For this example, we'll set a known backend
        if (process.env.BACKEND_ADDRESS) {
            const authTx = await newDistributor.setAuthorizedBackend(
                process.env.BACKEND_ADDRESS,
                true,
                { gasLimit: 100000 }
            );
            await authTx.wait();
            console.log("Authorized backend:", process.env.BACKEND_ADDRESS);
        }
        
        // Step 7: Update registry
        console.log("\nStep 7: Updating contract registry...");
        
        const updateTx = await registry.updateContract(
            "RewardDistributor",
            newDistributor.address,
            { gasLimit: 150000 }
        );
        await updateTx.wait();
        console.log("Registry updated");
        
        // Step 8: Verify migration
        console.log("\nStep 8: Verifying migration...");
        
        const registeredAddress = await registry.getContract("RewardDistributor");
        console.log("Registry points to:", registeredAddress);
        console.log("Migration successful:", registeredAddress === newDistributor.address);
        
        const newBalance = await token.balanceOf(newDistributor.address);
        console.log("New distributor token balance:", ethers.utils.formatEther(newBalance), "BOGO");
        
        // Step 9: Mark migration in helper (if available)
        if (migrationHelper) {
            console.log("\nStep 9: Recording migration...");
            const markTx = await migrationHelper.markMigrated(
                oldDistributorAddress,
                newDistributor.address,
                { gasLimit: 100000 }
            );
            await markTx.wait();
            console.log("Migration recorded in helper");
        }
        
        // Step 10: Save migration report
        console.log("\nStep 10: Saving migration report...");
        
        const report = {
            timestamp: new Date().toISOString(),
            network: network.name,
            migration: {
                from: oldDistributorAddress,
                to: newDistributor.address,
                version: await registry.getContractVersion("RewardDistributor"),
                tokensMigrated: ethers.utils.formatEther(oldBalance),
            },
            configuration: {
                token: CONFIG.TOKEN_ADDRESS,
                treasury: CONFIG.TREASURY_ADDRESS,
                registry: CONFIG.REGISTRY_ADDRESS,
            },
            deployer: deployer.address,
            transactionHashes: {
                deployment: newDistributor.deployTransaction.hash,
                registryUpdate: updateTx.hash,
            }
        };
        
        const reportPath = `migrations/reward-distributor-${network.name}-${Date.now()}.json`;
        fs.mkdirSync("migrations", { recursive: true });
        fs.writeFileSync(reportPath, JSON.stringify(report, null, 2));
        console.log("Migration report saved to:", reportPath);
        
        // Summary
        console.log("\n=== Migration Summary ===");
        console.log("✅ Old distributor paused");
        console.log("✅ New distributor deployed:", newDistributor.address);
        console.log("✅ Tokens transferred:", ethers.utils.formatEther(oldBalance), "BOGO");
        console.log("✅ Registry updated");
        console.log("✅ Migration complete!");
        
        // Post-migration checklist
        console.log("\n=== Post-Migration Checklist ===");
        console.log("[ ] Update frontend to use new contract address");
        console.log("[ ] Update backend services");
        console.log("[ ] Monitor new contract for 24 hours");
        console.log("[ ] Announce migration to community");
        console.log("[ ] Update documentation");
        
    } catch (error) {
        console.error("\n❌ Migration failed:", error.message);
        console.error("\nRollback procedure:");
        console.error("1. Keep old distributor paused");
        console.error("2. Investigate error");
        console.error("3. Fix issues and retry");
        console.error("4. If critical, unpause old distributor");
        throw error;
    }
}

// Migration utilities
async function getUsersToMigrate(oldDistributor, startBlock, endBlock) {
    // Get all unique users who have interacted with the contract
    // This would analyze events to find users with active states
    const filter = oldDistributor.filters.RewardClaimed();
    const events = await oldDistributor.queryFilter(filter, startBlock, endBlock);
    
    const users = new Set();
    events.forEach(event => {
        users.add(event.args.wallet);
    });
    
    return Array.from(users);
}

async function migrateUserBatch(users, oldContract, newContract, migrationHelper) {
    console.log(`Migrating batch of ${users.length} users...`);
    
    const tx = await migrationHelper.batchMarkMigrated(
        oldContract,
        users,
        { gasLimit: 2000000 }
    );
    await tx.wait();
    
    console.log(`Batch migration complete`);
}

// Execute migration
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });