import { network } from "hardhat";

const { ethers } = await network.connect();

async function main() {
  const initialSupply = 1000;

  console.log("Deploying MyToken...");

  const Token = await ethers.getContractFactory("MyToken");

  const token = await Token.deploy(initialSupply);

  await token.waitForDeployment();

  const address = await token.getAddress();

  console.log("✅ Token deployed at:", address);
  console.log(`Update explorer config: CONTRACT_ADDRESS = "${address}"`);
}

main().catch((error) => {
  console.error("❌ Deployment failed:", error);
  process.exitCode = 1;
});