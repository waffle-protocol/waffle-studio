package ui

import "fmt"

// Status constants matching contracts.StatusX values
const (
	StatusPending   uint8 = 0
	StatusSubmitted uint8 = 1
	StatusAccepted  uint8 = 2
	StatusRejected  uint8 = 3
	StatusCancelled uint8 = 4
)

// StatusInfo contains display information for a status
type StatusInfo struct {
	Icon  string
	Text  string
	Color string
}

// StatusMap maps status codes to display information
var StatusMap = map[uint8]StatusInfo{
	StatusPending:   {Icon: "⏳", Text: "PENDING", Color: Yellow},
	StatusSubmitted: {Icon: "📝", Text: "SUBMITTED", Color: Cyan},
	StatusAccepted:  {Icon: "✅", Text: "ACCEPTED", Color: Green},
	StatusRejected:  {Icon: "❌", Text: "REJECTED", Color: Red},
	StatusCancelled: {Icon: "🚫", Text: "CANCELLED", Color: Gray},
}

// FormatStatus returns a formatted status string with icon and color
func FormatStatus(status uint8) string {
	info, ok := StatusMap[status]
	if !ok {
		return fmt.Sprintf("%s%s%s", Gray, "UNKNOWN", Reset)
	}
	return fmt.Sprintf("%s%s%s", info.Color, info.Text, Reset)
}

// GetStatusIcon returns the icon for a status
func GetStatusIcon(status uint8) string {
	info, ok := StatusMap[status]
	if !ok {
		return "❓"
	}
	return info.Icon
}

// FormatStatusWithIcon returns status icon and colored text
func FormatStatusWithIcon(status uint8) (string, string) {
	info, ok := StatusMap[status]
	if !ok {
		return "❓", fmt.Sprintf("%sUNKNOWN%s", Gray, Reset)
	}
	return info.Icon, fmt.Sprintf("%s%-9s%s", info.Color, info.Text, Reset)
}
