import { network } from "hardhat";

const { ethers } = await network.connect();

async function main() {
  const contractAddress = process.env.MTK_CONTRACT_ADDRESS;
  const to = process.env.MTK_TO_ADDRESS;
  const amountHuman = process.env.MTK_AMOUNT ?? "10";

  if (!contractAddress || !to) {
    throw new Error(
      "Set env vars: MTK_CONTRACT_ADDRESS, MTK_TO_ADDRESS, optional MTK_AMOUNT"
    );
  }

  const [sender] = await ethers.getSigners();
  const token = await ethers.getContractAt("MyToken", contractAddress);
  const decimals = await token.decimals();
  const amount = ethers.parseUnits(amountHuman, decimals);

  console.log("Sending MTK transfer...");
  console.log("From:", sender.address);
  console.log("To:", to);
  console.log("Amount:", amountHuman);

  const tx = await token.transfer(to, amount);
  console.log("Tx hash:", tx.hash);

  await tx.wait();
  console.log("✅ Transfer confirmed");
}

main().catch((error) => {
  console.error("❌ MTK transfer failed:", error);
  process.exitCode = 1;
});
