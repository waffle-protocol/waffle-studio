package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/wallet"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check wallet balance",
	Long:  `Display the current SYRUP token balance for the connected wallet.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show loading indicator
		fmt.Printf("%s⏳ Connecting to blockchain...%s\n", ColorBlue, ColorReset)

		// Load configuration from environment
		config, err := wallet.LoadConfig()
		if err != nil {
			fmt.Printf("%s❌ Error: %s%s\n", "\033[31m", err.Error(), ColorReset)
			fmt.Println()
			fmt.Println("Please set environment variables:")
			fmt.Println("  export PRIVATE_KEY=0x...")
			fmt.Println("  export SYRUP_TOKEN=0x...")
			fmt.Println("  export RPC_URL=http://127.0.0.1:8545  (optional)")
			os.Exit(1)
		}

		// Create wallet
		w, err := wallet.New(config)
		if err != nil {
			fmt.Printf("%s❌ Failed to connect: %s%s\n", "\033[31m", err.Error(), ColorReset)
			os.Exit(1)
		}
		defer w.Close()

		// Query token balance
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tokenAddr := common.HexToAddress(config.SyrupToken)
		balance, err := w.GetTokenBalance(ctx, tokenAddr)
		if err != nil {
			fmt.Printf("%s❌ Failed to get balance: %s%s\n", "\033[31m", err.Error(), ColorReset)
			os.Exit(1)
		}

		// Format balance (18 decimals for SYRUP)
		balanceFormatted := wallet.FormatTokenBalance(balance, 18)
		address := w.AddressHex()

		// Server Rack style wallet display with box border
		fmt.Println()
		fmt.Printf("%s╔══════════════════════════════════════════════════════╗%s\n", ColorBlue, ColorReset)
		fmt.Printf("%s║%s %sWAFFLE WALLET%s                                        %s║%s\n", ColorBlue, ColorReset, ColorAmber+ColorBold, ColorReset, ColorBlue, ColorReset)
		fmt.Printf("%s╠══════════════════════════════════════════════════════╣%s\n", ColorBlue, ColorReset)
		fmt.Printf("%s║%s 💳 Address: %s%-42s%s %s║%s\n", ColorBlue, ColorReset, ColorAmber, address, ColorReset, ColorBlue, ColorReset)
		fmt.Printf("%s║%s 🍯 Balance: %s%-42.4f%s %s║%s\n", ColorBlue, ColorReset, ColorAmber, balanceFormatted, ColorReset, ColorBlue, ColorReset)
		fmt.Printf("%s║%s            %sSYRUP%s                                      %s║%s\n", ColorBlue, ColorReset, ColorBold, ColorReset, ColorBlue, ColorReset)
		fmt.Printf("%s╚══════════════════════════════════════════════════════╝%s\n", ColorBlue, ColorReset)
		fmt.Println()
	},
}
