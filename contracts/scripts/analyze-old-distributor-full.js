const hre = require("hardhat");
const { ethers } = hre;

async function main() {
  console.log("\n🔍 COMPREHENSIVE ANALYSIS OF OLD REWARD DISTRIBUTOR");
  console.log("=".repeat(70));

  const OLD_DISTRIBUTOR = process.env.ACTUAL_DISTRIBUTOR || "0x95C5Be6Ed592C663fF2953C683dBc5E2C257eA9f";
  const BOGO_TOKEN = "0xD394c8fEe6dC8b25DD423AE2f6e68191BD379c0C";
  
  console.log(`Old Distributor: ${OLD_DISTRIBUTOR}`);
  console.log(`BOGO Token: ${BOGO_TOKEN}`);

  try {
    // Get token contract
    const bogoToken = await ethers.getContractAt("BOGOTokenV2", BOGO_TOKEN);
    
    // Check balance
    const balance = await bogoToken.balanceOf(OLD_DISTRIBUTOR);
    console.log(`\n💰 Current Token Balance: ${ethers.utils.formatEther(balance)} BOGO`);
    
    if (balance.eq(0)) {
      console.log("✅ No tokens are locked in the old distributor");
    } else {
      console.log("⚠️  WARNING: Tokens are locked in the old distributor!");
    }

    // Get contract code to verify it exists
    const code = await ethers.provider.getCode(OLD_DISTRIBUTOR);
    console.log(`\n📄 Contract Status:`);
    console.log(`- Code size: ${code.length} bytes`);
    console.log(`- Is contract: ${code.length > 2 ? 'Yes' : 'No'}`);

    // Check all historical transfers
    console.log("\n📊 TRANSFER HISTORY ANALYSIS");
    console.log("-".repeat(70));
    
    // Incoming transfers
    const filterIn = bogoToken.filters.Transfer(null, OLD_DISTRIBUTOR);
    const eventsIn = await bogoToken.queryFilter(filterIn);
    
    console.log(`\n📥 Incoming Transfers: ${eventsIn.length}`);
    if (eventsIn.length > 0) {
      let totalIn = ethers.BigNumber.from(0);
      const uniqueSenders = new Set();
      
      for (const event of eventsIn) {
        totalIn = totalIn.add(event.args.value);
        uniqueSenders.add(event.args.from);
      }
      
      console.log(`- Total received: ${ethers.utils.formatEther(totalIn)} BOGO`);
      console.log(`- Unique senders: ${uniqueSenders.size}`);
      console.log(`- First transfer block: ${eventsIn[0].blockNumber}`);
      console.log(`- Last transfer block: ${eventsIn[eventsIn.length - 1].blockNumber}`);
      
      // Show last few incoming transfers
      console.log("\nLast 3 incoming transfers:");
      for (let i = Math.max(0, eventsIn.length - 3); i < eventsIn.length; i++) {
        const event = eventsIn[i];
        console.log(`  - ${ethers.utils.formatEther(event.args.value)} BOGO from ${event.args.from.slice(0, 10)}...`);
      }
    }
    
    // Outgoing transfers
    const filterOut = bogoToken.filters.Transfer(OLD_DISTRIBUTOR, null);
    const eventsOut = await bogoToken.queryFilter(filterOut);
    
    console.log(`\n📤 Outgoing Transfers: ${eventsOut.length}`);
    if (eventsOut.length > 0) {
      let totalOut = ethers.BigNumber.from(0);
      const uniqueRecipients = new Set();
      
      for (const event of eventsOut) {
        totalOut = totalOut.add(event.args.value);
        uniqueRecipients.add(event.args.to);
      }
      
      console.log(`- Total distributed: ${ethers.utils.formatEther(totalOut)} BOGO`);
      console.log(`- Unique recipients: ${uniqueRecipients.size}`);
      console.log(`- First distribution block: ${eventsOut[0].blockNumber}`);
      console.log(`- Last distribution block: ${eventsOut[eventsOut.length - 1].blockNumber}`);
      
      // Show last few distributions
      console.log("\nLast 3 distributions:");
      for (let i = Math.max(0, eventsOut.length - 3); i < eventsOut.length; i++) {
        const event = eventsOut[i];
        const block = await event.getBlock();
        const date = new Date(block.timestamp * 1000).toLocaleDateString();
        console.log(`  - ${ethers.utils.formatEther(event.args.value)} BOGO to ${event.args.to.slice(0, 10)}... on ${date}`);
      }
    }
    
    // Calculate net flow
    console.log("\n💵 NET FLOW ANALYSIS");
    console.log("-".repeat(70));
    
    const totalReceived = eventsIn.reduce((sum, e) => sum.add(e.args.value), ethers.BigNumber.from(0));
    const totalSent = eventsOut.reduce((sum, e) => sum.add(e.args.value), ethers.BigNumber.from(0));
    const shouldHaveBalance = totalReceived.sub(totalSent);
    
    console.log(`Total Received:     ${ethers.utils.formatEther(totalReceived)} BOGO`);
    console.log(`Total Distributed:  ${ethers.utils.formatEther(totalSent)} BOGO`);
    console.log(`Expected Balance:   ${ethers.utils.formatEther(shouldHaveBalance)} BOGO`);
    console.log(`Actual Balance:     ${ethers.utils.formatEther(balance)} BOGO`);
    
    if (!shouldHaveBalance.eq(balance)) {
      console.log(`\n⚠️  DISCREPANCY: ${ethers.utils.formatEther(shouldHaveBalance.sub(balance))} BOGO mismatch!`);
    } else {
      console.log("\n✅ Balance matches expected amount");
    }
    
    // Summary
    console.log("\n📋 SUMMARY");
    console.log("-".repeat(70));
    
    if (balance.gt(0)) {
      console.log("❌ CRITICAL: This contract has NO treasurySweep or emergencyWithdraw function!");
      console.log(`❌ ${ethers.utils.formatEther(balance)} BOGO tokens are PERMANENTLY LOCKED!`);
      console.log("\n🔴 These tokens cannot be recovered without:");
      console.log("   1. The contract having a withdrawal function (it doesn't)");
      console.log("   2. A critical vulnerability that allows extraction (not recommended)");
      console.log("   3. The owner having special privileges (check ownership)");
    } else {
      console.log("✅ No tokens are currently locked in this contract");
      console.log("✅ Safe to deploy new distributor without token loss");
    }
    
  } catch (error) {
    console.error("\n❌ Error:", error.message);
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });