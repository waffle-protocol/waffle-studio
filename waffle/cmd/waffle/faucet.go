package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/relay"
	"github.com/waffle-studio/waffle/internal/ui"
)

func init() {
	rootCmd.AddCommand(faucetCmd)
}

var faucetCmd = &cobra.Command{
	Use:   "faucet",
	Short: "Claim free ETH + SYRUP tokens via relay (testnet only)",
	Long:  `Claim free testnet ETH and SYRUP tokens from the gasless relay server.`,
	Run: func(cmd *cobra.Command, args []string) {
		claimFaucet()
	},
}

func claimFaucet() {
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ FAUCET ]"))
	fmt.Printf("  %s Requesting tokens from relay...\n", ui.TextAmber.Render("⏳"))

	// Bootstrap to get wallet address
	ctx, err := cli.BootstrapWithWallet()
	if err != nil {
		cli.PrintError(err.Error())
		cli.PrintConfigHint()
		return
	}
	defer ctx.Close()

	// Get wallet address
	address := ctx.Config.GetAddress()
	if address == "" {
		cli.PrintError("Failed to get wallet address")
		return
	}

	fmt.Printf("  📬 Address: %s\n", address)

	// Create relay client
	relayClient := relay.NewClient(ctx.Config.RelayURL)

	// Check relay health
	if err := relayClient.Health(); err != nil {
		cli.PrintErrorf("Relay server unavailable: %s", err)
		fmt.Println()
		fmt.Println("  💡 The relay server might be down or unreachable.")
		fmt.Println("     You can get testnet ETH from: https://www.alchemy.com/faucets/base-sepolia")
		return
	}

	fmt.Printf("  %s Relay server connected\n", ui.SuccessStyle.Render("✅"))
	fmt.Printf("  %s Sending ETH + SYRUP...\n", ui.TextAmber.Render("⏳"))

	// Request faucet
	resp, err := relayClient.Faucet(address)
	if err != nil {
		cli.PrintErrorf("Faucet request failed: %s", err)
		return
	}

	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ SUCCESS ]"))
	fmt.Printf("  %s Received tokens!\n", ui.SuccessStyle.Render("✅"))
	fmt.Printf("  ⛽ ETH:   %s (for gas)\n", resp.EthAmount)
	fmt.Printf("  🍯 SYRUP: %s\n", resp.SyrupAmount)
	fmt.Printf("  🔗 Tx:    %s\n", resp.TxHash[:20]+"...")
	fmt.Println()
	fmt.Println("  Run 'waffle balance' to check your new balance.")
	fmt.Println()
}
