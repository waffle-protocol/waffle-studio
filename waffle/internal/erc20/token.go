package erc20

// BalanceOfABI is the ABI for ERC-20 balanceOf function
const BalanceOfABI = `[{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

// ApproveABI is the ABI for ERC-20 approve function
const ApproveABI = `[{"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`
