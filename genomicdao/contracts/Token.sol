// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract PostCovidStrokePrevention is ERC20, ERC20Burnable, Ownable {

    mapping(uint256 => uint256) riskScoreToAward;

    constructor() ERC20("Post-Covid Stroke Prevention", "PCSP") {
        _mint(msg.sender, 1000000000 * 10 ** decimals());

        riskScoreToAward[1] = 15000 * 10 ** decimals();
        riskScoreToAward[2] = 3000 * 10 ** decimals();
        riskScoreToAward[3] = 225 * 10 ** decimals();
        riskScoreToAward[4] = 30 * 10 ** decimals();
    }

    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }

    function reward(address to, uint256 riskScore) public onlyOwner {
        // TODO: Implement this method: Award PCSP to the user based on his/her risk score
        require(to != address(0), "Invalid address");
        require(riskScoreToAward[riskScore] > 0, "No reward for the risk score");

        uint256 awardAmount = riskScoreToAward[riskScore];

        _mint(to, awardAmount);
    }
}