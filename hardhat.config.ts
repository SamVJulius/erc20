import hardhatToolboxMochaEthersPlugin from "@nomicfoundation/hardhat-toolbox-mocha-ethers";
import { defineConfig } from "hardhat/config";

export default defineConfig({
  plugins: [hardhatToolboxMochaEthersPlugin],

  solidity: {
    profiles: {
      default: {
        version: "0.8.20", // safer for OpenZeppelin
      },
      production: {
        version: "0.8.20",
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
        },
      },
    },
  },

  networks: {
    localevm: {
      type: "http",
      chainType: "l1",
      url: "http://127.0.0.1:8545",   // your Cosmos EVM RPC
      chainId: 262144,                  // ⚠️ MUST match your chain
      accounts: [
        "0x01804a6a954c20323fcda061eebabf520fb8d554dc1b46fd37380b40624e618d"          // ⚠️ replace this
      ],
    },
  },
});