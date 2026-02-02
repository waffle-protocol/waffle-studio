package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/patch"
	"github.com/waffle-studio/waffle/internal/ui"
)

// ============================================================================
// Cobra Command Definition
// ============================================================================

func init() {
	rootCmd.AddCommand(reviewCmd)
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review and apply code patches interactively",
	Long:  `Interactive diff viewer for reviewing Waffle bake results. Accept or reject hunks with y/n keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the Bubble Tea program
		p := tea.NewProgram(initialReviewModel(), tea.WithAltScreen())

		// Run the TUI
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running review TUI: %v\n", err)
			os.Exit(1)
		}

		// Show summary after TUI exits
		if m, ok := finalModel.(reviewModel); ok {
			showReviewSummary(m)
		}
	},
}

// ============================================================================
// Hunk Data Structure
// ============================================================================

// Hunk represents a single diff hunk (a block of changes)
type Hunk struct {
	FilePath  string   // File being modified
	StartLine int      // Starting line number in original file
	OldLines  []string // Lines to be removed (prefixed with -)
	NewLines  []string // Lines to be added (prefixed with +)
	Context   []string // Surrounding context lines
	Accepted  bool     // User accepted this hunk
	Rejected  bool     // User rejected this hunk
	Decided   bool     // User has made a decision
}

// ============================================================================
// Bubble Tea Model
// ============================================================================

type reviewModel struct {
	hunks      []Hunk // All hunks to review
	currentIdx int    // Currently selected hunk index
	quitting   bool   // User is quitting
	width      int    // Terminal width
	height     int    // Terminal height
}

// initialReviewModel creates the initial state with mock diff data
func initialReviewModel() reviewModel {
	// Mock hunks for demonstration
	// In production, this would come from P2P Baker responses
	hunks := []Hunk{
		{
			FilePath:  "cmd/waffle/main.go",
			StartLine: 15,
			OldLines:  []string{"func oldHandler() {", "    return nil", "}"},
			NewLines:  []string{"func newHandler(ctx context.Context) error {", "    return ctx.Err()", "}"},
			Context:   []string{"// Handler for processing requests"},
		},
		{
			FilePath:  "cmd/waffle/main.go",
			StartLine: 42,
			OldLines:  []string{"fmt.Println(\"debug\")"},
			NewLines:  []string{"log.Debug(\"processing request\")"},
			Context:   []string{"// Log the request status"},
		},
		{
			FilePath:  "internal/engine/processor.go",
			StartLine: 78,
			OldLines:  []string{"time.Sleep(100 * time.Millisecond)"},
			NewLines:  []string{"select {", "case <-ctx.Done():", "    return ctx.Err()", "case <-time.After(100 * time.Millisecond):", "}"},
			Context:   []string{"// Wait for processing to complete"},
		},
	}

	return reviewModel{
		hunks:      hunks,
		currentIdx: 0,
	}
}

// ============================================================================
// Bubble Tea Lifecycle Methods
// ============================================================================

// Init is called when the program starts
func (m reviewModel) Init() tea.Cmd {
	return nil
}

// Update handles all messages/events
func (m reviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		// Quit
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		// Navigation - next hunk
		case "j", "down", "tab":
			if m.currentIdx < len(m.hunks)-1 {
				m.currentIdx++
			}
			return m, nil

		// Navigation - previous hunk
		case "k", "up", "shift+tab":
			if m.currentIdx > 0 {
				m.currentIdx--
			}
			return m, nil

		// Accept hunk
		case "y":
			if m.currentIdx < len(m.hunks) {
				m.hunks[m.currentIdx].Accepted = true
				m.hunks[m.currentIdx].Rejected = false
				m.hunks[m.currentIdx].Decided = true
				// Auto-advance to next undecided hunk
				m.advanceToNextUndecided()
			}
			return m, nil

		// Reject hunk
		case "n":
			if m.currentIdx < len(m.hunks) {
				m.hunks[m.currentIdx].Accepted = false
				m.hunks[m.currentIdx].Rejected = true
				m.hunks[m.currentIdx].Decided = true
				// Auto-advance to next undecided hunk
				m.advanceToNextUndecided()
			}
			return m, nil

		// Skip (move without deciding)
		case "s", " ":
			if m.currentIdx < len(m.hunks)-1 {
				m.currentIdx++
			}
			return m, nil

		// Apply all accepted and quit
		case "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

// advanceToNextUndecided moves to the next hunk that hasn't been decided
func (m *reviewModel) advanceToNextUndecided() {
	// First, try to find undecided hunk after current position
	for i := m.currentIdx + 1; i < len(m.hunks); i++ {
		if !m.hunks[i].Decided {
			m.currentIdx = i
			return
		}
	}
	// If none found, stay at current (or move to next if possible)
	if m.currentIdx < len(m.hunks)-1 {
		m.currentIdx++
	}
}

// View renders the UI
func (m reviewModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	width := m.width
	if width < 60 {
		width = 60
	}
	if width > 80 {
		width = 80
	}

	// Get current hunk
	if len(m.hunks) == 0 {
		return "\n  No hunks to review.\n\n"
	}
	hunk := m.hunks[m.currentIdx]

	// Build the box
	borderH := strings.Repeat("═", width-2)

	// Header
	b.WriteString("\n")
	b.WriteString(ui.BorderStyle.Render("╔" + borderH + "╗"))
	b.WriteString("\n")

	// Title with hunk counter
	title := fmt.Sprintf("📜 WAFFLE REVIEW                        [%d/%d hunks]",
		m.currentIdx+1, len(m.hunks))
	title = padRight(title, width-4)
	b.WriteString(ui.BorderStyle.Render("║ "))
	b.WriteString(ui.TitleStyle.Render(title))
	b.WriteString(ui.BorderStyle.Render(" ║"))
	b.WriteString("\n")

	// Separator
	b.WriteString(ui.BorderStyle.Render("╠" + borderH + "╣"))
	b.WriteString("\n")

	// File info
	fileInfo := fmt.Sprintf("File: %s (line %d)", hunk.FilePath, hunk.StartLine)
	fileInfo = padRight(fileInfo, width-4)
	b.WriteString(ui.BorderStyle.Render("║ "))
	b.WriteString(ui.FileStyle.Render(fileInfo))
	b.WriteString(ui.BorderStyle.Render(" ║"))
	b.WriteString("\n")

	// Status indicator
	var statusText string
	if hunk.Decided {
		if hunk.Accepted {
			statusText = ui.AcceptedStyle.Render("✅ ACCEPTED")
		} else {
			statusText = ui.RejectedStyle.Render("❌ REJECTED")
		}
	} else {
		statusText = ui.PendingStyle.Render("⏳ PENDING")
	}
	statusLine := padRight("Status: "+statusText, width-4)
	b.WriteString(ui.BorderStyle.Render("║ "))
	b.WriteString(statusLine)
	b.WriteString(ui.BorderStyle.Render(" ║"))
	b.WriteString("\n")

	// Separator
	b.WriteString(ui.BorderStyle.Render("╠" + borderH + "╣"))
	b.WriteString("\n")

	// Context lines
	for _, line := range hunk.Context {
		ctxLine := padRight("  "+line, width-4)
		b.WriteString(ui.BorderStyle.Render("║ "))
		b.WriteString(ui.ContextLineStyle.Render(ctxLine))
		b.WriteString(ui.BorderStyle.Render(" ║"))
		b.WriteString("\n")
	}

	// Removed lines (old)
	for _, line := range hunk.OldLines {
		diffLine := padRight("- "+line, width-4)
		b.WriteString(ui.BorderStyle.Render("║ "))
		b.WriteString(ui.RemovedLineStyle.Render(diffLine))
		b.WriteString(ui.BorderStyle.Render(" ║"))
		b.WriteString("\n")
	}

	// Added lines (new)
	for _, line := range hunk.NewLines {
		diffLine := padRight("+ "+line, width-4)
		b.WriteString(ui.BorderStyle.Render("║ "))
		b.WriteString(ui.AddedLineStyle.Render(diffLine))
		b.WriteString(ui.BorderStyle.Render(" ║"))
		b.WriteString("\n")
	}

	// Separator
	b.WriteString(ui.BorderStyle.Render("╠" + borderH + "╣"))
	b.WriteString("\n")

	// Help bar
	helpText := "[y] Accept  [n] Reject  [s] Skip  [j/k] Navigate  [q] Quit & Apply"
	helpText = padRight(helpText, width-4)
	b.WriteString(ui.BorderStyle.Render("║ "))
	b.WriteString(ui.HelpStyle.Render(helpText))
	b.WriteString(ui.BorderStyle.Render(" ║"))
	b.WriteString("\n")

	// Footer
	b.WriteString(ui.BorderStyle.Render("╚" + borderH + "╝"))
	b.WriteString("\n")

	// Progress summary
	accepted, rejected, pending := m.countDecisions()
	progressText := fmt.Sprintf("\n  Progress: %s %d accepted  %s %d rejected  %s %d pending\n",
		ui.AcceptedStyle.Render("●"), accepted,
		ui.RejectedStyle.Render("●"), rejected,
		ui.PendingStyle.Render("●"), pending)
	b.WriteString(progressText)

	return b.String()
}

// countDecisions returns counts of accepted, rejected, and pending hunks
func (m reviewModel) countDecisions() (accepted, rejected, pending int) {
	for _, h := range m.hunks {
		if h.Decided {
			if h.Accepted {
				accepted++
			} else {
				rejected++
			}
		} else {
			pending++
		}
	}
	return
}

// padRight pads a string to the specified width
func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

// ============================================================================
// Post-Review Summary
// ============================================================================

// showReviewSummary displays the final summary and applies patches
func showReviewSummary(m reviewModel) {
	accepted, rejected, pending := m.countDecisions()

	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 🧇 REVIEW COMPLETE ]"))
	fmt.Println()

	if accepted > 0 {
		fmt.Printf("  %s %d hunk(s) accepted\n", ui.AcceptedStyle.Render("✅"), accepted)
		for _, h := range m.hunks {
			if h.Accepted {
				fmt.Printf("     └─ %s (line %d)\n", ui.TextAmber.Render(h.FilePath), h.StartLine)
			}
		}
	}

	if rejected > 0 {
		fmt.Printf("  %s %d hunk(s) rejected\n", ui.RejectedStyle.Render("❌"), rejected)
	}

	if pending > 0 {
		fmt.Printf("  %s %d hunk(s) skipped\n", ui.PendingStyle.Render("⏳"), pending)
	}

	fmt.Println()

	// Apply patches if any accepted
	if accepted > 0 {
		// Ask for confirmation
		fmt.Printf("%s Apply %d patch(es) to files? [y/n]: %s", ui.Bold, accepted, ui.Reset)

		var response string
		fmt.Scanln(&response)
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" || response == "yes" {
			applyAcceptedPatches(m.hunks)
		} else {
			fmt.Println(ui.HelpStyle.Render("  Patches not applied."))
		}
	} else {
		fmt.Println(ui.HelpStyle.Render("  No patches to apply."))
	}
	fmt.Println()
}

// applyAcceptedPatches applies all accepted hunks to files
func applyAcceptedPatches(hunks []Hunk) {
	// Convert review hunks to patch hunks
	var patchHunks []patch.Hunk
	for _, h := range hunks {
		if h.Accepted {
			patchHunks = append(patchHunks, patch.Hunk{
				FilePath:  h.FilePath,
				StartLine: h.StartLine,
				OldLines:  h.OldLines,
				NewLines:  h.NewLines,
			})
		}
	}

	if len(patchHunks) == 0 {
		return
	}

	fmt.Println()
	fmt.Println(ui.RackStyle.Render("[ 🍯 APPLYING PATCHES ]"))
	fmt.Println()

	results := patch.ApplyAll(patchHunks)

	for _, r := range results {
		if r.Success {
			fmt.Printf("  %s %s\n", ui.AcceptedStyle.Render("✅"), r.FilePath)
		} else {
			fmt.Printf("  %s %s: %s\n", ui.RejectedStyle.Render("❌"), r.FilePath, r.Error)
		}
	}

	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("  ✨ Patches applied successfully!"))
}
