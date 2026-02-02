package main

import (
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
	"github.com/waffle-studio/waffle/internal/ui"
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
// Blockchain Integration
// ============================================================================

// handleBakeRequest handles the blockchain flow after file selection
func handleBakeRequest(filePath string, reward float64) {
	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 🔗 BLOCKCHAIN ]"))
	fmt.Printf("  %s Connecting to network...\n", ui.TextAmber.Render("⏳"))

	// Bootstrap with full context
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

	// Read file content for hashing
	content, err := os.ReadFile(filePath)
	if err != nil {
		cli.PrintErrorf("Failed to read file: %s", err)
		return
	}

	codeHash := contracts.HashCode(string(content))

	// Convert reward to wei (18 decimals)
	rewardWei := new(big.Int)
	rewardFloat := new(big.Float).SetFloat64(reward)
	decimals := new(big.Float).SetFloat64(1e18)
	rewardFloat.Mul(rewardFloat, decimals)
	rewardFloat.Int(rewardWei)

	timeoutCtx, cancel := cli.WithTimeout()
	defer cancel()

	// Step 1: Approve SYRUP
	fmt.Printf("  %s Approving %g SYRUP...\n", ui.TextAmber.Render("⏳"), reward)

	approveReceipt, err := ctx.Registry.ApproveToken(timeoutCtx, ctx.Config.SyrupToken, rewardWei)
	if err != nil {
		cli.PrintErrorf("Approve failed: %s", err)
		return
	}
	fmt.Printf("  %s Approved! Tx: %s\n", ui.SuccessStyle.Render("✅"), approveReceipt.TxHash.Hex()[:18]+"...")

	// Step 2: Create Request
	fmt.Printf("  %s Creating request on BakeRegistry...\n", ui.TextAmber.Render("⏳"))

	receipt, requestID, err := ctx.Registry.CreateRequest(timeoutCtx, codeHash, rewardWei)
	if err != nil {
		cli.PrintErrorf("CreateRequest failed: %s", err)
		return
	}

	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 🧇 REQUEST CREATED ]"))
	fmt.Printf("  %s Request ID: %s#%d%s\n", ui.SuccessStyle.Render("✅"), ui.Bold, requestID, ui.Reset)
	fmt.Printf("  💰 Reward: %g SYRUP (escrowed)\n", reward)
	fmt.Printf("  📄 File: %s\n", filePath)
	fmt.Printf("  🔗 Tx: %s\n", receipt.TxHash.Hex())
	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("  ✨ Waiting for a Baker to pick up your request!"))
	fmt.Println()
}
