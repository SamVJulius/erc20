import { network } from "hardhat";

const { ethers } = await network.connect();

async function main() {
  const contractAddress = process.env.MTK_CONTRACT_ADDRESS || "0xEfc2DfbFd0E03F80B8af07c73431C18E6D0f6E4b";

  console.log("=== Explorer Network Debug ===\n");

  // Get network info from connected RPC
  const chainId = await ethers.provider.getNetwork();
  const blockNumber = await ethers.provider.getBlockNumber();

  console.log("Terminal/Hardhat (via ethers):");
  console.log(`  RPC: http://127.0.0.1:8545`);
  console.log(`  Chain ID: ${chainId.chainId}`);
  console.log(`  Latest Block: ${blockNumber}`);
  console.log("");

  console.log("⚠️  MetaMask should match BOTH:");
  console.log(`  - RPC URL: http://127.0.0.1:8545`);
  console.log(`  - Chain ID: ${chainId.chainId}`);
  console.log("");

  console.log("=== Token Transfers ===\n");
  console.log("Contract Address:", contractAddress);
  console.log("");

  try {
    const token = await ethers.getContractAt("MyToken", contractAddress);
    const name = await token.name();
    const symbol = await token.symbol();
    console.log("✅ Token found:", name, `(${symbol})`);
    console.log("");

    // Get Transfer events from last 5000 blocks (RPC limit safety)
    const transferFilter = await token.filters.Transfer();
    const fromBlock = Math.max(0, blockNumber - 5000);
    const events = await token.queryFilter(transferFilter, fromBlock, blockNumber);

    console.log(`Total Transfer events: ${events.length}`);
    console.log("");

    if (events.length > 0) {
      console.log("Recent transfers (excluding mint from 0x0):");
      const nonMint = events.filter((e) => e.args.from !== "0x0000000000000000000000000000000000000000");
      
      if (nonMint.length === 0) {
        console.log("  (only mint event found - no wallet-to-wallet transfers yet)");
      } else {
        for (const event of nonMint.slice(-5)) {
          const from = event.args.from;
          const to = event.args.to;
          const value = ethers.formatEther(event.args.value);
          console.log(`  ${from} → ${to}: ${value} ${symbol}`);
        }
      }
    } else {
      console.log("❌ No transfers found!");
    }
    console.log("");
    console.log("If MetaMask transfers don't appear:");
    console.log("1. Check MetaMask Network Settings match chain ID above");
    console.log("2. Verify RPC URL is http://127.0.0.1:8545");
    console.log("3. Check MetaMask's transaction history for errors");
  } catch (error: any) {
    console.error("❌ Error:", error.message);
  }
}

main().catch((error) => {
  console.error("❌ Debug failed:", error);
  process.exitCode = 1;
});
