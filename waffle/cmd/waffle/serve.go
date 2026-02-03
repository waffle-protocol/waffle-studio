package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/ai"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/p2p"
	"github.com/waffle-studio/waffle/internal/ui"
)

var servePort int

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 0, "Port to listen on (0 for random)")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as a Baker provider",
	Long:  `Start the Waffle node in provider mode to process bake requests from peers.`,
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ BAKER PROVIDER ]"))
	fmt.Printf("  %s Initializing...\n", ui.TextAmber.Render(""))

	// Bootstrap wallet for encryption key
	ctx, err := cli.BootstrapWithWallet()
	if err != nil {
		cli.PrintError(err.Error())
		return
	}
	defer ctx.Close()

	// Get Gemini API key
	apiKey, err := ai.GetAPIKey()
	if err != nil {
		cli.PrintErrorf("Failed to get Gemini API key: %s", err)
		return
	}

	// Initialize Gemini AI provider
	bgCtx := context.Background()
	aiProvider, err := ai.NewGeminiProvider(bgCtx, apiKey)
	if err != nil {
		cli.PrintErrorf("Failed to initialize AI provider: %s", err)
		return
	}
	defer aiProvider.Close()

	fmt.Printf("  %s AI provider ready\n", ui.SuccessStyle.Render(""))

	// Create P2P node
	node, err := p2p.NewNode(p2p.NodeConfig{
		PrivateKey: ctx.Config.PrivateKey,
		ListenPort: servePort,
	})
	if err != nil {
		cli.PrintErrorf("Failed to create P2P node: %s", err)
		return
	}
	defer node.Close()

	// Setup provider handler
	node.SetupProvider(func(prompt string, fileData []byte) ([]byte, error) {
		fmt.Printf("\n  %s Received request\n", ui.TextAmber.Render(""))
		fmt.Printf("    Prompt: %s\n", prompt)
		fmt.Printf("    File size: %d bytes\n", len(fileData))

		result, err := aiProvider.ProcessCode(bgCtx, prompt, fileData)
		if err != nil {
			fmt.Printf("  %s Error: %s\n", ui.ErrorStyle.Render(""), err)
			return nil, err
		}

		fmt.Printf("  %s Request processed\n", ui.SuccessStyle.Render(""))
		return result, nil
	})

	// Print node info
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ NODE INFO ]"))
	fmt.Printf("  Peer ID: %s\n", node.ID())
	fmt.Println("  Addresses:")
	for _, addr := range node.Addrs() {
		fmt.Printf("    %s\n", addr)
	}
	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("  Baker provider is running!"))
	fmt.Println("  Waiting for bake requests... (Ctrl+C to stop)")
	fmt.Println()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println()
	fmt.Println(ui.TextAmber.Render("  Shutting down..."))
}
