package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// Color definitions for lipgloss styles
var (
	AmberColor = lipgloss.Color("214") // Amber (#FFBF00)
	BlueColor  = lipgloss.Color("27")  // Blue  (#0000FF)
	GreenColor = lipgloss.Color("82")  // Green
	RedColor   = lipgloss.Color("196") // Red
	GrayColor  = lipgloss.Color("241") // Gray
	WhiteColor = lipgloss.Color("252") // White/Light gray
)

// Shared styles for Server Rack Theme

// TitleStyle for headers and titles
var TitleStyle = lipgloss.NewStyle().
	Foreground(AmberColor).
	Bold(true).
	MarginLeft(2)

// ItemStyle for list items
var ItemStyle = lipgloss.NewStyle().
	Foreground(WhiteColor).
	PaddingLeft(4)

// SelectedItemStyle for selected items in lists
var SelectedItemStyle = lipgloss.NewStyle().
	Foreground(AmberColor).
	Bold(true).
	PaddingLeft(2)

// PaginationStyle for list pagination
var PaginationStyle = list.DefaultStyles().PaginationStyle.
	Foreground(BlueColor).
	PaddingLeft(4)

// HelpStyle for help text
var HelpStyle = list.DefaultStyles().HelpStyle.
	Foreground(GrayColor).
	PaddingLeft(4).
	PaddingBottom(1)

// RackStyle for server rack headers
var RackStyle = lipgloss.NewStyle().
	Foreground(BlueColor).
	Bold(true)

// TextAmber for amber colored text
var TextAmber = lipgloss.NewStyle().
	Foreground(AmberColor)

// SuccessStyle for success messages
var SuccessStyle = lipgloss.NewStyle().
	Foreground(GreenColor).
	Bold(true)

// ErrorStyle for error messages
var ErrorStyle = lipgloss.NewStyle().
	Foreground(RedColor).
	Bold(true)

// BorderStyle for box borders
var BorderStyle = lipgloss.NewStyle().
	Foreground(BlueColor)

// FileStyle for file paths
var FileStyle = lipgloss.NewStyle().
	Foreground(WhiteColor)

// Diff line styles

// AddedLineStyle for added lines in diffs
var AddedLineStyle = lipgloss.NewStyle().
	Foreground(GreenColor).
	Bold(true)

// RemovedLineStyle for removed lines in diffs
var RemovedLineStyle = lipgloss.NewStyle().
	Foreground(RedColor).
	Bold(true)

// ContextLineStyle for context lines in diffs
var ContextLineStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("245"))

// Status styles

// AcceptedStyle for accepted status
var AcceptedStyle = lipgloss.NewStyle().
	Foreground(GreenColor).
	Bold(true)

// RejectedStyle for rejected status
var RejectedStyle = lipgloss.NewStyle().
	Foreground(RedColor).
	Bold(true)

// PendingStyle for pending status
var PendingStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("214")).
	Bold(true)
