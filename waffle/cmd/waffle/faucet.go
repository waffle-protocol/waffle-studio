package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
)

func init() {
	rootCmd.AddCommand(faucetCmd)
}

var faucetCmd = &cobra.Command{
	Use:   "faucet",
	Short: "Claim 100 free SYRUP tokens (testnet only)",
	Long:  `Claim 100 SYRUP tokens from the testnet faucet. Each address can only claim once.`,
	Run: func(cmd *cobra.Command, args []string) {
		claimFaucet()
	},
}

const faucetABI = `[{"inputs":[],"name":"faucet","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

func claimFaucet() {
	cli.PrintLoading("Claiming SYRUP from faucet...")

	ctx, err := cli.BootstrapFull()
	if err != nil {
		cli.PrintError(err.Error())
		cli.PrintConfigHint()
		return
	}
	defer ctx.Close()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(faucetABI))
	if err != nil {
		cli.PrintErrorf("Failed to parse ABI: %s", err)
		return
	}

	// Pack faucet call
	data, err := parsedABI.Pack("faucet")
	if err != nil {
		cli.PrintErrorf("Failed to pack faucet data: %s", err)
		return
	}

	// Get private key
	keyHex := strings.TrimPrefix(ctx.Config.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		cli.PrintErrorf("Invalid private key: %s", err)
		return
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	tokenAddress := common.HexToAddress(ctx.Config.SyrupToken)

	// Get nonce
	nonce, err := ctx.Client.PendingNonceAt(timeoutCtx, fromAddress)
	if err != nil {
		cli.PrintErrorf("Failed to get nonce: %s", err)
		return
	}

	// Get gas price
	gasPrice, err := ctx.Client.SuggestGasPrice(timeoutCtx)
	if err != nil {
		cli.PrintErrorf("Failed to get gas price: %s", err)
		return
	}

	// Estimate gas
	gasLimit, err := ctx.Client.EstimateGas(timeoutCtx, ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		cli.PrintErrorf("Failed (already claimed or contract issue): %s", err)
		return
	}

	// Get chain ID
	chainID, err := ctx.Client.ChainID(timeoutCtx)
	if err != nil {
		cli.PrintErrorf("Failed to get chain ID: %s", err)
		return
	}

	// Create and sign transaction
	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		cli.PrintErrorf("Failed to sign transaction: %s", err)
		return
	}

	// Send transaction
	err = ctx.Client.SendTransaction(timeoutCtx, signedTx)
	if err != nil {
		cli.PrintErrorf("Failed to send transaction: %s", err)
		return
	}

	fmt.Println()
	fmt.Printf("  ⏳ Transaction sent: %s\n", signedTx.Hash().Hex()[:20]+"...")
	fmt.Println("  ⏳ Waiting for confirmation...")

	// Wait for receipt
	for {
		receipt, err := ctx.Client.TransactionReceipt(timeoutCtx, signedTx.Hash())
		if err == nil {
			if receipt.Status == 1 {
				fmt.Println()
				fmt.Println("  ✅ Successfully claimed 100 SYRUP!")
				fmt.Printf("  🔗 Tx: %s\n", signedTx.Hash().Hex())
				fmt.Println()
				fmt.Println("  Run 'waffle balance' to check your new balance.")
				fmt.Println()
			} else {
				cli.PrintError("Transaction failed")
			}
			return
		}

		select {
		case <-timeoutCtx.Done():
			cli.PrintError("Timeout waiting for transaction")
			return
		case <-time.After(time.Second):
			continue
		}
	}
}
