package main

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/config"
	"github.com/waffle-studio/waffle/internal/contracts"
)

// ============================================================================
// Styling - Server Rack Theme (Blue/Amber)
// ============================================================================

var (
	// Color definitions
	amberColor = lipgloss.Color("214") // Amber (#FFBF00)
	blueColor  = lipgloss.Color("27")  // Blue  (#0000FF)

	// Styles for the list
	titleStyle = lipgloss.NewStyle().
			Foreground(amberColor).
			Bold(true).
			MarginLeft(2)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(amberColor).
				Bold(true).
				PaddingLeft(2)

	paginationStyle = list.DefaultStyles().PaginationStyle.
			Foreground(blueColor).
			PaddingLeft(4)

	helpStyle = list.DefaultStyles().HelpStyle.
			Foreground(lipgloss.Color("241")).
			PaddingLeft(4).
			PaddingBottom(1)

	// Styles for the baking result
	rackStyle = lipgloss.NewStyle().
			Foreground(blueColor).
			Bold(true)

	textAmber = lipgloss.NewStyle().
			Foreground(amberColor)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Bold(true)
)

// ============================================================================
// Cobra Command Definition
// ============================================================================

var bakeReward float64

func init() {
	bakeCmd.Flags().Float64VarP(&bakeReward, "reward", "r", 10.0, "SYRUP reward for bakers")
	rootCmd.AddCommand(bakeCmd)
}

var bakeCmd = &cobra.Command{
	Use:   "bake",
	Short: "Select and bake a code recipe",
	Long:  `Select a file from the current directory and process it as a Waffle recipe.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the Bubble Tea program
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())

		// Run the TUI
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running bake TUI: %v\n", err)
			os.Exit(1)
		}

		// Handle blockchain flow if file was selected
		if m, ok := finalModel.(model); ok && m.baking && m.selectedFile != "" {
			handleBakeRequest(m.selectedFile, bakeReward)
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
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("▶ " + strings.Join(s, " "))
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
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

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
		b.WriteString(rackStyle.Render("[ 🧇 BAKING ]"))
		b.WriteString(textAmber.Render(fmt.Sprintf("  Processing '%s'...\n", m.selectedFile)))
		b.WriteString(rackStyle.Render("[ 🍯 GLAZING ]"))
		b.WriteString(textAmber.Render("  Applying patches...\n"))
		b.WriteString(successStyle.Render("✨ Waffle is ready!\n"))
		b.WriteString("\n")
		return b.String()
	}

	// Normal view - show the list
	return "\n" + m.list.View()
}

// ============================================================================
// Blockchain Integration
// ============================================================================

// handleBakeRequest handles the blockchain flow after file selection
func handleBakeRequest(filePath string, reward float64) {
	fmt.Println()
	fmt.Println(rackStyle.Render("[ 🔗 BLOCKCHAIN ]"))
	fmt.Printf("  %s Connecting to network...\n", textAmber.Render("⏳"))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("  %s Failed to load config: %s\n", "\033[31m❌\033[0m", err)
		return
	}

	if cfg.PrivateKey == "" || cfg.SyrupToken == "" || cfg.BakeRegistry == "" {
		fmt.Printf("  %s Missing configuration\n", "\033[31m❌\033[0m")
		fmt.Println("  Please set PRIVATE_KEY, SYRUP_TOKEN, BAKE_REGISTRY in ~/.waffle/config.yaml or env vars.")
		return
	}

	// Connect to RPC
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		fmt.Printf("  %s Failed to connect: %s\n", "\033[31m❌\033[0m", err)
		return
	}
	defer client.Close()

	// Create registry instance
	registry, err := contracts.NewRegistryClient(cfg.BakeRegistry, client, cfg.PrivateKey)
	if err != nil {
		fmt.Printf("  %s Failed to setup registry: %s\n", "\033[31m", err)
		return
	}

	// Read file content for hashing
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("  %s Failed to read file: %s\n", "\033[31m❌\033[0m", err)
		return
	}

	codeHash := contracts.HashCode(string(content))

	// Convert reward to wei (18 decimals)
	rewardWei := new(big.Int)
	rewardFloat := new(big.Float).SetFloat64(reward)
	decimals := new(big.Float).SetFloat64(1e18)
	rewardFloat.Mul(rewardFloat, decimals)
	rewardFloat.Int(rewardWei)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Step 1: Approve SYRUP
	fmt.Printf("  %s Approving %g SYRUP...\n", textAmber.Render("⏳"), reward)

	approveReceipt, err := registry.ApproveToken(ctx, cfg.SyrupToken, rewardWei)
	if err != nil {
		fmt.Printf("  %s Approve failed: %s\n", "\033[31m❌\033[0m", err)
		return
	}
	fmt.Printf("  %s Approved! Tx: %s\n", successStyle.Render("✅"), approveReceipt.TxHash.Hex()[:18]+"...")

	// Step 2: Create Request
	fmt.Printf("  %s Creating request on BakeRegistry...\n", textAmber.Render("⏳"))

	receipt, requestID, err := registry.CreateRequest(ctx, codeHash, rewardWei)
	if err != nil {
		fmt.Printf("  %s CreateRequest failed: %s\n", "\033[31m❌\033[0m", err)
		return
	}

	fmt.Println()
	fmt.Println(rackStyle.Render("[ 🧇 REQUEST CREATED ]"))
	fmt.Printf("  %s Request ID: %s#%d%s\n", successStyle.Render("✅"), "\033[1m", requestID, "\033[0m")
	fmt.Printf("  💰 Reward: %g SYRUP (escrowed)\n", reward)
	fmt.Printf("  📄 File: %s\n", filePath)
	fmt.Printf("  🔗 Tx: %s\n", receipt.TxHash.Hex())
	fmt.Println()
	fmt.Println(successStyle.Render("  ✨ Waiting for a Baker to pick up your request!"))
	fmt.Println()
}
