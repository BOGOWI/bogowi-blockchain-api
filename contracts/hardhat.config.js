require("@nomicfoundation/hardhat-chai-matchers");
require("@nomicfoundation/hardhat-ethers");
require("@nomicfoundation/hardhat-verify");
require("solidity-coverage");
require("dotenv").config({ path: "../.env" });

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
  networks: {
    hardhat: {
      chainId: 1337
    },
    columbus: {
      url: process.env.RPC_URL || "https://columbus.camino.network/ext/bc/C/rpc",
      chainId: 501,
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
      gasPrice: 225000000000, // 225 gwei (above minimum 200 gwei)
      timeout: 60000
    },
    camino: {
      url: process.env.RPC_URL || "https://columbus.camino.network/ext/bc/C/rpc",
      chainId: 501,
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
      gasPrice: 225000000000, // 225 gwei (above minimum 200 gwei) 
      timeout: 60000
    }
  },
  etherscan: {
    apiKey: {
      columbus: "placeholder",
      camino: "placeholder"
    },
    customChains: [
      {
        network: "columbus",
        chainId: 501,
        urls: {
          apiURL: "https://explorer.camino.network/api",
          browserURL: "https://explorer.camino.network"
        }
      },
      {
        network: "camino",
        chainId: 500,
        urls: {
          apiURL: "https://explorer.camino.network/api",
          browserURL: "https://explorer.camino.network"
        }
      }
    ]
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts"
  }
};