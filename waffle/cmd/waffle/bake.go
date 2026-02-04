package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/diff"
	"github.com/waffle-studio/waffle/internal/p2p"
	"github.com/waffle-studio/waffle/internal/ui"
)

// ============================================================================
// Cobra Command Definition
// ============================================================================

var bakeReward float64
var bakePeer string

func init() {
	bakeCmd.Flags().Float64VarP(&bakeReward, "reward", "r", 10.0, "SYRUP reward for bakers")
	bakeCmd.Flags().StringVar(&bakePeer, "peer", "", "P2P provider address to connect to directly")
	rootCmd.AddCommand(bakeCmd)
}

var bakeCmd = &cobra.Command{
	Use:   "bake",
	Short: "Select and bake a code recipe",
	Long:  `Select a file from the current directory and process it as a Waffle recipe with P2P + blockchain integration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the Bubble Tea program
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())

		// Run the TUI
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running bake TUI: %v\n", err)
			os.Exit(1)
		}

		// Handle flow if file was selected
		if m, ok := finalModel.(model); ok && m.baking && m.selectedFile != "" {
			handleIntegratedBake(m.selectedFile, bakeReward, bakePeer)
		}
	},
}

// ============================================================================
// List Item Implementation
// ============================================================================

// fileItem represents a file in the list
type fileItem struct {
	name  string
	isDir bool
}

// FilterValue implements list.Item interface - used for filtering
func (f fileItem) FilterValue() string { return f.name }

// itemDelegate handles rendering of list items
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(fileItem)
	if !ok {
		return
	}

	// Format: icon + filename
	icon := "📄"
	if item.isDir {
		icon = "📁"
	}
	str := fmt.Sprintf("%s %s", icon, item.name)

	// Apply styling based on selection
	fn := ui.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return ui.SelectedItemStyle.Render("▶ " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// ============================================================================
// Bubble Tea Model
// ============================================================================

// model represents the application state
type model struct {
	list         list.Model // The file selection list
	selectedFile string     // The file chosen by user
	baking       bool       // Whether we're in baking state
	quitting     bool       // Whether the user is quitting
}

// initialModel creates the initial state
func initialModel() model {
	// Read files from current directory
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		os.Exit(1)
	}

	// Convert to list items
	var items []list.Item
	for _, f := range files {
		// Skip hidden files
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		items = append(items, fileItem{
			name:  f.Name(),
			isDir: f.IsDir(),
		})
	}

	// Create the list with custom delegate
	const listHeight = 15
	const listWidth = 50

	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = "📂 Select a Code Recipe to Bake"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = ui.TitleStyle
	l.Styles.PaginationStyle = ui.PaginationStyle
	l.Styles.HelpStyle = ui.HelpStyle

	return model{list: l}
}

// ============================================================================
// Bubble Tea Lifecycle Methods
// ============================================================================

// Init is called when the program starts - no initial commands needed
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles all messages/events
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window resize event
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	// Keyboard input
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {

		// Quit keys
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		// Selection key
		case "enter":
			// Get selected item
			if item, ok := m.list.SelectedItem().(fileItem); ok {
				m.selectedFile = item.name
				m.baking = true
				// Small delay to show the baking animation, then quit
				return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
					return tea.Quit()
				})
			}
		}
	}

	// Pass other messages to the list for handling
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the UI
func (m model) View() string {
	// User quit without selecting
	if m.quitting {
		return "\n  👋 Bake cancelled.\n\n"
	}

	// Baking complete - show Server Rack style result
	if m.baking {
		var b strings.Builder
		b.WriteString("\n")
		b.WriteString(ui.RackStyle.Render("[ 🧇 BAKING ]"))
		b.WriteString(ui.TextAmber.Render(fmt.Sprintf("  Processing '%s'...\n", m.selectedFile)))
		b.WriteString(ui.RackStyle.Render("[ 🍯 GLAZING ]"))
		b.WriteString(ui.TextAmber.Render("  Applying patches...\n"))
		b.WriteString(ui.SuccessStyle.Render("✨ Waffle is ready!\n"))
		b.WriteString("\n")
		return b.String()
	}

	// Normal view - show the list
	return "\n" + m.list.View()
}

// ============================================================================
// Integrated P2P + Blockchain Baking
// ============================================================================

// handleIntegratedBake handles the integrated P2P + blockchain flow
func handleIntegratedBake(filePath string, reward float64, peerAddr string) {
	reader := bufio.NewReader(os.Stdin)

	// ========== Phase 1: Initialize ==========
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 🔗 INITIALIZING ]"))
	fmt.Printf("  %s Bootstrapping full context...\n", ui.TextAmber.Render("⏳"))

	// Bootstrap with full context (wallet + blockchain)
	ctx, err := cli.BootstrapFull()
	if err != nil {
		cli.PrintError(err.Error())
		cli.PrintConfigHintFull()
		return
	}
	defer ctx.Close()

	if ctx.Config.SyrupToken == "" {
		cli.PrintError("SYRUP_TOKEN not configured")
		cli.PrintConfigHintFull()
		return
	}

	fmt.Printf("  %s Connected to blockchain\n", ui.SuccessStyle.Render("✅"))

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		cli.PrintErrorf("Failed to read file: %s", err)
		return
	}

	// ========== Phase 2: P2P Provider Discovery ==========
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 📡 P2P DISCOVERY ]"))
	fmt.Printf("  %s Creating P2P node...\n", ui.TextAmber.Render("⏳"))

	// Create P2P node
	node, err := p2p.NewNode(p2p.NodeConfig{
		PrivateKey: ctx.Config.PrivateKey,
		ListenPort: 0,
	})
	if err != nil {
		cli.PrintErrorf("Failed to create P2P node: %s", err)
		return
	}
	defer node.Close()

	fmt.Printf("  %s Node ID: %s\n", ui.SuccessStyle.Render("✅"), node.ID()[:16]+"...")

	// Connect to manual peer if specified
	if peerAddr != "" {
		fmt.Printf("  %s Connecting to peer: %s...\n", ui.TextAmber.Render("⏳"), peerAddr)
		if err := node.ConnectToPeer(context.Background(), peerAddr); err != nil {
			cli.PrintErrorf("Failed to connect to peer: %s", err)
			return
		}
		fmt.Printf("  %s Connected to peer!\n", ui.SuccessStyle.Render("✅"))
	}

	// Discover providers
	fmt.Printf("  %s Discovering providers...\n", ui.TextAmber.Render("⏳"))
	peers, err := node.DiscoverPeers(15 * time.Second)
	if err != nil {
		cli.PrintErrorf("Failed to discover peers: %s", err)
		return
	}

	if len(peers) == 0 {
		cli.PrintError("No providers found. Start one with: waffle serve")
		return
	}

	fmt.Printf("  %s Found %d provider(s)\n", ui.SuccessStyle.Render("✅"), len(peers))

	// ========== Phase 3: Get Prompt from User ==========
	fmt.Println()
	fmt.Print("  Enter your modification request: ")
	prompt, err := reader.ReadString('\n')
	if err != nil {
		cli.PrintErrorf("Failed to read prompt: %s", err)
		return
	}
	prompt = strings.TrimSpace(prompt)

	if prompt == "" {
		cli.PrintError("Prompt cannot be empty")
		return
	}

	// ========== Phase 4: Create Blockchain Request (SYRUP Escrow) ==========
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 💰 BLOCKCHAIN ESCROW ]"))

	codeHash := contracts.HashCode(string(content))

	// Convert reward to wei (18 decimals)
	rewardWei := new(big.Int)
	rewardFloat := new(big.Float).SetFloat64(reward)
	decimals := new(big.Float).SetFloat64(1e18)
	rewardFloat.Mul(rewardFloat, decimals)
	rewardFloat.Int(rewardWei)

	timeoutCtx, cancel := cli.WithTimeout()
	defer cancel()

	// Approve SYRUP
	fmt.Printf("  %s Approving %g SYRUP...\n", ui.TextAmber.Render("⏳"), reward)

	approveReceipt, err := ctx.Registry.ApproveToken(timeoutCtx, ctx.Config.SyrupToken, rewardWei)
	if err != nil {
		cli.PrintErrorf("Approve failed: %s", err)
		return
	}
	fmt.Printf("  %s Approved! Tx: %s\n", ui.SuccessStyle.Render("✅"), approveReceipt.TxHash.Hex()[:18]+"...")

	// Create Request
	fmt.Printf("  %s Creating request (escrowing SYRUP)...\n", ui.TextAmber.Render("⏳"))

	_, requestID, err := ctx.Registry.CreateRequest(timeoutCtx, codeHash, rewardWei)
	if err != nil {
		cli.PrintErrorf("CreateRequest failed: %s", err)
		return
	}

	fmt.Printf("  %s Request #%d created (%.2f SYRUP escrowed)\n", ui.SuccessStyle.Render("✅"), requestID, reward)

	// ========== Phase 5: Send P2P Request with RequestID ==========
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 🤖 AI PROCESSING ]"))
	fmt.Printf("  %s Sending request to provider (Request #%d)...\n", ui.TextAmber.Render("⏳"), requestID)

	bgCtx := context.Background()
	response, err := node.SendRequestWithID(bgCtx, peers[0], prompt, content, requestID.Uint64())
	if err != nil {
		cli.PrintErrorf("Failed to send request: %s", err)
		// Offer to cancel the blockchain request
		fmt.Print("\n  Cancel blockchain request and get SYRUP refund? [y/N]: ")
		cancelConfirm, _ := reader.ReadString('\n')
		cancelConfirm = strings.TrimSpace(strings.ToLower(cancelConfirm))
		if cancelConfirm == "y" || cancelConfirm == "yes" {
			handleCancel(ctx, requestID)
		}
		return
	}

	fmt.Printf("  %s Received response from provider\n", ui.SuccessStyle.Render("✅"))
	fmt.Printf("    Provider: %s\n", response.ProviderAddress)
	fmt.Printf("    Tokens used: %d\n", response.TokenUsage)

	// ========== Phase 6: Show Diff ==========
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 📊 CHANGES ]"))
	stats := diff.CalculateStats(string(content), string(response.Data))
	diff.PrintGitStyleStat(filePath, stats)
	fmt.Println()
	diff.PrintDiff(string(content), string(response.Data))
	diff.PrintSummary(stats)

	// ========== Phase 7: Accept / Reject / Cancel ==========
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ ⚖️  DECISION ]"))
	fmt.Println("  What would you like to do?")
	fmt.Println("    [a] Accept - Apply changes & pay provider")
	fmt.Println("    [r] Reject - Try another provider")
	fmt.Println("    [c] Cancel - Discard & get SYRUP refund")
	fmt.Println()
	fmt.Print("  Your choice [a/r/c]: ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	switch choice {
	case "a", "accept":
		handleAccept(ctx, requestID, response, filePath, reward)

	case "r", "reject":
		handleReject(ctx, requestID)

	case "c", "cancel":
		handleCancel(ctx, requestID)

	default:
		fmt.Println()
		fmt.Println(ui.TextAmber.Render("  Invalid choice. Request remains pending."))
		fmt.Printf("  You can manage it later with: waffle requests\n")
	}

	fmt.Println()
}

// handleAccept accepts the solution and pays the provider
func handleAccept(ctx *cli.ClientContext, requestID *big.Int, response *p2p.SyrupPayload, filePath string, reward float64) {
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ ✅ ACCEPTING ]"))

	// Calculate payment based on token usage
	// Payment = min(reward, tokenUsage * rate)
	paymentWei := new(big.Int)
	paymentFloat := new(big.Float).SetFloat64(reward)
	decimals := new(big.Float).SetFloat64(1e18)
	paymentFloat.Mul(paymentFloat, decimals)
	paymentFloat.Int(paymentWei)

	timeoutCtx, cancel := cli.WithTimeout()
	defer cancel()

	fmt.Printf("  %s Accepting solution & paying provider...\n", ui.TextAmber.Render("⏳"))

	_, err := ctx.Registry.AcceptSolution(timeoutCtx, requestID, paymentWei)
	if err != nil {
		cli.PrintErrorf("AcceptSolution failed: %s", err)
		return
	}

	fmt.Printf("  %s Payment sent to provider!\n", ui.SuccessStyle.Render("✅"))

	// Apply changes to file
	if err := os.WriteFile(filePath, response.Data, 0644); err != nil {
		cli.PrintErrorf("Failed to write file: %s", err)
		return
	}

	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("  ✨ Changes applied successfully!"))
	fmt.Printf("  💰 Paid %.2f SYRUP to %s\n", reward, response.ProviderAddress[:10]+"...")
}

// handleReject rejects the solution (allows trying another provider)
func handleReject(ctx *cli.ClientContext, requestID *big.Int) {
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ ❌ REJECTING ]"))

	timeoutCtx, cancel := cli.WithTimeout()
	defer cancel()

	fmt.Printf("  %s Rejecting solution...\n", ui.TextAmber.Render("⏳"))

	_, err := ctx.Registry.RejectSolution(timeoutCtx, requestID)
	if err != nil {
		cli.PrintErrorf("RejectSolution failed: %s", err)
		return
	}

	fmt.Printf("  %s Solution rejected\n", ui.SuccessStyle.Render("✅"))
	fmt.Println("  You can try another provider with: waffle bake")
}

// handleCancel cancels the request and refunds SYRUP
func handleCancel(ctx *cli.ClientContext, requestID *big.Int) {
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 🔙 CANCELLING ]"))

	timeoutCtx, cancel := cli.WithTimeout()
	defer cancel()

	fmt.Printf("  %s Cancelling request & refunding SYRUP...\n", ui.TextAmber.Render("⏳"))

	_, err := ctx.Registry.CancelRequest(timeoutCtx, requestID)
	if err != nil {
		cli.PrintErrorf("CancelRequest failed: %s", err)
		return
	}

	fmt.Printf("  %s Request cancelled, SYRUP refunded!\n", ui.SuccessStyle.Render("✅"))
}
