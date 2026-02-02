package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "waffle",
	Short: "Waffle - P2P LLM Token Marketplace CLI",
	Long:  `Waffle is a decentralized marketplace for LLM tokens, built on blockchain technology.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Server Rack style ASCII art
		fmt.Println()
		fmt.Printf("%s[ 🧇🧇🧇 ]%s  %sWAFFLE STACK%s online.\n", ui.Amber, ui.Reset, ui.Bold, ui.Reset)
		fmt.Printf("%s[ 🧇🍯🧇 ]%s  %sConnecting to P2P Grid...%s\n", ui.Amber, ui.Reset, ui.Blue, ui.Reset)
		fmt.Printf("%s[ 🧇🧇🧇 ]%s  v0.1.0 ready.\n", ui.Amber, ui.Reset)
		fmt.Println()
		fmt.Printf("%s🥞 Waffle 스택이 활성화되었습니다.%s\n", ui.Bold, ui.Reset)
		fmt.Println()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
