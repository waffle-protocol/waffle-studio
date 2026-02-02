package ui

// ANSI color codes for terminal output
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Amber   = "\033[38;5;214m"
	Blue    = "\033[34m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"
)

// Error returns a red colored error message with ❌ prefix
func Error(text string) string {
	return Red + "❌ " + text + Reset
}

// Success returns a green colored success message with ✅ prefix
func Success(text string) string {
	return Green + "✅ " + text + Reset
}

// Warning returns a yellow colored warning message with ⚠️ prefix
func Warning(text string) string {
	return Yellow + "⚠️  " + text + Reset
}

// Info returns a blue colored info message with ⏳ prefix
func Info(text string) string {
	return Blue + "⏳ " + text + Reset
}

// Cancelled returns a gray colored cancelled message with 🚫 prefix
func Cancelled(text string) string {
	return Gray + "🚫 " + text + Reset
}
