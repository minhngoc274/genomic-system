require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.19",
  networks: {
    "lifenetwork": {
      url: "http://127.0.0.1:9650/ext/bc/2bGAh54yzGQ3nj4txNDDSzoKZTNfghang9z92Cgz315ch1nAsA/rpc ",
      chainId: 9999,
      accounts: [process.env.PRIVATE_KEY],
    }
  }
};
