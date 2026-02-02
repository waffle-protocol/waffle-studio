package patch

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Hunk represents a single diff hunk to be applied
type Hunk struct {
	FilePath  string   // File to modify
	StartLine int      // Starting line number (1-indexed)
	OldLines  []string // Lines to find and replace
	NewLines  []string // Replacement lines
}

// Result represents the result of applying a patch
type Result struct {
	FilePath     string
	Success      bool
	Error        error
	LinesChanged int
}

// ApplyAll applies all hunks and returns results
func ApplyAll(hunks []Hunk) []Result {
	var results []Result

	// Group hunks by file (apply in reverse line order to avoid offset issues)
	fileHunks := make(map[string][]Hunk)
	for _, h := range hunks {
		fileHunks[h.FilePath] = append(fileHunks[h.FilePath], h)
	}

	for filePath, hs := range fileHunks {
		// Sort hunks by StartLine in descending order
		// This ensures we apply from bottom to top, avoiding line offset issues
		sortHunksDescending(hs)

		result := Result{FilePath: filePath}

		err := applyHunksToFile(filePath, hs)
		if err != nil {
			result.Success = false
			result.Error = err
		} else {
			result.Success = true
			result.LinesChanged = countLinesChanged(hs)
		}

		results = append(results, result)
	}

	return results
}

// sortHunksDescending sorts hunks by StartLine in descending order
func sortHunksDescending(hunks []Hunk) {
	sort.Slice(hunks, func(i, j int) bool {
		return hunks[i].StartLine > hunks[j].StartLine
	})
}

// countLinesChanged counts total lines affected by hunks
func countLinesChanged(hunks []Hunk) int {
	count := 0
	for _, h := range hunks {
		count += len(h.OldLines) + len(h.NewLines)
	}
	return count
}

// applyHunksToFile applies multiple hunks to a single file
func applyHunksToFile(filePath string, hunks []Hunk) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Read file
	lines, err := readFileLines(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Apply each hunk
	for _, hunk := range hunks {
		lines, err = applyHunkToLines(lines, hunk)
		if err != nil {
			return fmt.Errorf("failed to apply hunk at line %d: %w", hunk.StartLine, err)
		}
	}

	// Write file back
	err = writeFileLines(absPath, lines)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// applyHunkToLines applies a single hunk to lines slice
func applyHunkToLines(lines []string, hunk Hunk) ([]string, error) {
	startIdx := hunk.StartLine - 1 // Convert to 0-indexed

	// Validate start index
	if startIdx < 0 || startIdx >= len(lines) {
		return nil, fmt.Errorf("start line %d out of range (file has %d lines)", hunk.StartLine, len(lines))
	}

	// Find matching old lines
	matchIdx := findMatchingLines(lines, hunk.OldLines, startIdx)
	if matchIdx == -1 {
		// Try searching around the expected position
		matchIdx = findMatchingLines(lines, hunk.OldLines, 0)
		if matchIdx == -1 {
			return nil, fmt.Errorf("could not find matching lines in file")
		}
	}

	// Replace old lines with new lines
	endIdx := matchIdx + len(hunk.OldLines)

	// Build new lines slice
	newLines := make([]string, 0, len(lines)-len(hunk.OldLines)+len(hunk.NewLines))
	newLines = append(newLines, lines[:matchIdx]...)
	newLines = append(newLines, hunk.NewLines...)
	newLines = append(newLines, lines[endIdx:]...)

	return newLines, nil
}

// findMatchingLines finds the index where oldLines match in lines
func findMatchingLines(lines, oldLines []string, startFrom int) int {
	if len(oldLines) == 0 {
		return -1
	}

	for i := startFrom; i <= len(lines)-len(oldLines); i++ {
		match := true
		for j, oldLine := range oldLines {
			// Normalize whitespace for comparison
			if normalizeWhitespace(lines[i+j]) != normalizeWhitespace(oldLine) {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// normalizeWhitespace trims and normalizes whitespace for comparison
func normalizeWhitespace(s string) string {
	return strings.TrimSpace(s)
}

// readFileLines reads a file and returns lines
func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// writeFileLines writes lines to a file
func writeFileLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, line := range lines {
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}
		// Add newline except for last line (unless file was originally multi-line)
		if i < len(lines)-1 {
			_, err = writer.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
	// Always end with newline
	_, err = writer.WriteString("\n")
	if err != nil {
		return err
	}

	return writer.Flush()
}
