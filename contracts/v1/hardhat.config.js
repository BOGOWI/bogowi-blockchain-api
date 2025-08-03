require("@nomicfoundation/hardhat-toolbox");
require("solidity-coverage");

// Use environment variables if available, otherwise use .env file as fallback
const PRIVATE_KEY = process.env.PRIVATE_KEY || process.env.DEPLOYER_PRIVATE_KEY || "";
const RPC_URL = process.env.RPC_URL || "";

// Network configurations
const networks = {};

// Only add network config if we have the required environment variables
if (PRIVATE_KEY && RPC_URL) {
  // Determine network name based on RPC URL
  let networkName = "custom";
  if (RPC_URL.includes("columbus")) {
    networkName = "testnet";
  } else if (RPC_URL.includes("api.camino.network")) {
    networkName = "mainnet";
  }
  
  networks[networkName] = {
    url: RPC_URL,
    accounts: [PRIVATE_KEY],
    chainId: networkName === "mainnet" ? 500 : 501, // Camino mainnet: 500, Columbus testnet: 501
  };
  
  // Also add specific network aliases
  if (networkName === "testnet") {
    networks.columbus = networks.testnet;
  } else if (networkName === "mainnet") {
    networks.camino = networks.mainnet;
  }
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
      allowUnlimitedContractSize: true
    },
    ...networks
  }
};