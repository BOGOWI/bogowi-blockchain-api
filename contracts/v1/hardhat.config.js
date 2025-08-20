require("@nomicfoundation/hardhat-toolbox");
require("solidity-coverage");

// Load appropriate .env file based on NODE_ENV or default to .env.dev for non-production
const envFile = process.env.NODE_ENV === 'production' ? "../../.env" : "../../.env.dev";
require("dotenv").config({ path: envFile });

// Get private key from environment or .env file
const PRIVATE_KEY = process.env.TESTNET_PRIVATE_KEY || process.env.MAINNET_PRIVATE_KEY || process.env.PRIVATE_KEY || "";

// Network configurations
const networks = {};

// Hardcoded network configurations
if (PRIVATE_KEY) {
  // Columbus testnet
  networks.columbus = {
    url: "https://columbus.camino.network/ext/bc/C/rpc",
    accounts: [PRIVATE_KEY],
    chainId: 501
  };
  networks.testnet = networks.columbus; // alias
  
  // Camino mainnet
  networks.camino = {
    url: "https://api.camino.network/ext/bc/C/rpc",
    accounts: [PRIVATE_KEY],
    chainId: 500
  };
  networks.mainnet = networks.camino; // alias
}

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    version: "0.8.20",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  },
  paths: {
    sources: "./contracts",
    tests: "./tests",
    cache: "./cache",
    artifacts: "./artifacts"
  },
  networks: {
    hardhat: {
      allowUnlimitedContractSize: true,
      chainId: 501 // Use Camino testnet chain ID for testing
    },
    ...networks
  }
};