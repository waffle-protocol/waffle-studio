package ui

import (
	"fmt"
	"strings"
)

// BoxWidth is the default width for boxes
const BoxWidth = 56

// Box represents a styled box for terminal output
type Box struct {
	width int
}

// NewBox creates a new box with default width
func NewBox() *Box {
	return &Box{width: BoxWidth}
}

// NewBoxWithWidth creates a new box with specified width
func NewBoxWithWidth(width int) *Box {
	return &Box{width: width}
}

// Header prints the top of a box with a title
func (b *Box) Header(title string) {
	borderH := strings.Repeat("═", b.width-2)
	fmt.Printf("%s╔%s╗%s\n", Blue, borderH, Reset)
	fmt.Printf("%s║%s %s%s%s%-*s%s║%s\n", Blue, Reset, Amber+Bold, title, Reset, b.width-len(title)-4, "", Blue, Reset)
	fmt.Printf("%s╠%s╣%s\n", Blue, borderH, Reset)
}

// Row prints a row in the box
func (b *Box) Row(content string) {
	// Calculate visible length (without ANSI codes)
	visibleLen := visibleLength(content)
	padding := b.width - visibleLen - 4
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s║%s %s%*s %s║%s\n", Blue, Reset, content, padding, "", Blue, Reset)
}

// RowWithLabel prints a row with an icon/label and value
func (b *Box) RowWithLabel(icon, label, value string) {
	content := fmt.Sprintf("%s %s: %s%s%s", icon, label, Amber, value, Reset)
	b.Row(content)
}

// EmptyRow prints an empty row
func (b *Box) EmptyRow() {
	fmt.Printf("%s║%s%*s%s║%s\n", Blue, Reset, b.width-2, "", Blue, Reset)
}

// Footer prints the bottom of a box
func (b *Box) Footer() {
	borderH := strings.Repeat("═", b.width-2)
	fmt.Printf("%s╚%s╝%s\n", Blue, borderH, Reset)
}

// Separator prints a separator line within the box
func (b *Box) Separator() {
	borderH := strings.Repeat("═", b.width-2)
	fmt.Printf("%s╠%s╣%s\n", Blue, borderH, Reset)
}

// WalletBox prints a formatted wallet display box
func WalletBox(address string, balance float64, tokenName string) {
	box := NewBox()
	fmt.Println()
	box.Header("WAFFLE WALLET")
	box.RowWithLabel("💳", "Address", address)
	box.RowWithLabel("🍯", "Balance", fmt.Sprintf("%.4f", balance))
	box.Row(fmt.Sprintf("            %s%s%s", Bold, tokenName, Reset))
	box.Footer()
	fmt.Println()
}

// ResultBox prints a transaction result box
func ResultBox(title, requestID string, fields map[string]string, txHash string) {
	box := NewBox()
	fmt.Println()
	box.Header(title)
	if requestID != "" {
		box.Row(fmt.Sprintf("   Request ID: #%-37s", requestID))
	}
	for label, value := range fields {
		box.Row(fmt.Sprintf("   %s: %s%s%s", label, Amber, value, Reset))
	}
	if txHash != "" {
		box.Row(fmt.Sprintf("   🔗 Tx: %-42s", truncateHash(txHash, 42)))
	}
	box.Footer()
	fmt.Println()
}

// RequestsHeader prints the header for requests list
func RequestsHeader(title string, width int) {
	fmt.Println()
	borderH := strings.Repeat("═", width-2)
	fmt.Printf("%s╔%s╗%s\n", Blue, borderH, Reset)
	fmt.Printf("%s║%s %s%s%s%-*s%s║%s\n", Blue, Reset, Amber+Bold, title, Reset, width-len(title)-4, "", Blue, Reset)
	fmt.Printf("%s╠%s╣%s\n", Blue, borderH, Reset)
}

// RequestsFooter prints the footer for requests list
func RequestsFooter(width int) {
	borderH := strings.Repeat("═", width-2)
	fmt.Printf("%s╚%s╝%s\n", Blue, borderH, Reset)
	fmt.Println()
}

// RequestsRow prints a row in the requests list
func RequestsRow(width int, content string) {
	visibleLen := visibleLength(content)
	padding := width - visibleLen - 4
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s║%s %s%*s %s║%s\n", Blue, Reset, content, padding, "", Blue, Reset)
}

// visibleLength calculates the visible length of a string (excluding ANSI codes)
func visibleLength(s string) int {
	// Simple ANSI stripping - remove escape sequences
	inEscape := false
	length := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}

// truncateHash truncates a hash to specified length
func truncateHash(hash string, maxLen int) string {
	if len(hash) <= maxLen {
		return hash
	}
	return hash[:maxLen]
}
