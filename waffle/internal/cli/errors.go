package cli

import (
	"fmt"
	"os"

	"github.com/waffle-studio/waffle/internal/ui"
)

// PrintError prints a formatted error message
func PrintError(msg string) {
	fmt.Println(ui.Error(msg))
}

// PrintErrorf prints a formatted error message with format string
func PrintErrorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(ui.Error(msg))
}

// PrintWarning prints a formatted warning message
func PrintWarning(msg string) {
	fmt.Println(ui.Warning(msg))
}

// PrintWarningf prints a formatted warning message with format string
func PrintWarningf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(ui.Warning(msg))
}

// PrintSuccess prints a formatted success message
func PrintSuccess(msg string) {
	fmt.Println(ui.Success(msg))
}

// PrintSuccessf prints a formatted success message with format string
func PrintSuccessf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(ui.Success(msg))
}

// PrintConfigHint prints a hint about configuration
func PrintConfigHint() {
	fmt.Println()
	fmt.Println("Please set PRIVATE_KEY and BAKE_REGISTRY in ~/.waffle/config.yaml or env vars.")
}

// PrintConfigHintFull prints a full configuration hint
func PrintConfigHintFull() {
	fmt.Println()
	fmt.Println("Please set PRIVATE_KEY, SYRUP_TOKEN, BAKE_REGISTRY in ~/.waffle/config.yaml or env vars.")
}

// FatalError prints error and exits
func FatalError(msg string) {
	PrintError(msg)
	os.Exit(1)
}

// FatalErrorf prints formatted error and exits
func FatalErrorf(format string, args ...interface{}) {
	PrintErrorf(format, args...)
	os.Exit(1)
}

// FatalErrorWithHint prints error with config hint and exits
func FatalErrorWithHint(msg string) {
	PrintError(msg)
	PrintConfigHint()
	os.Exit(1)
}
