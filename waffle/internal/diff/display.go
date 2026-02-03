package diff

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/waffle-studio/waffle/internal/ui"
)

// Stats holds statistics about changes between two texts
type Stats struct {
	Insertions int
	Deletions  int
}

// CalculateStats calculates the number of inserted and deleted lines
func CalculateStats(original, modified string) Stats {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(original, modified, false)

	var stats Stats
	for _, d := range diffs {
		lines := strings.Count(d.Text, "\n")
		if d.Text != "" && !strings.HasSuffix(d.Text, "\n") {
			lines++ // Count partial lines
		}
		switch d.Type {
		case diffmatchpatch.DiffInsert:
			stats.Insertions += lines
		case diffmatchpatch.DiffDelete:
			stats.Deletions += lines
		}
	}

	return stats
}

// PrintGitStyleStat prints git-style stats with +/- bar graph
func PrintGitStyleStat(filename string, stats Stats) {
	total := stats.Insertions + stats.Deletions
	if total == 0 {
		fmt.Printf(" %s | 0\n", filename)
		return
	}

	// Calculate bar lengths (max 50 chars for the bar)
	const maxBarWidth = 50
	scale := float64(maxBarWidth) / float64(total)
	if scale > 1 {
		scale = 1
	}

	plusCount := int(float64(stats.Insertions) * scale)
	minusCount := int(float64(stats.Deletions) * scale)

	// Ensure at least 1 char if there are changes
	if stats.Insertions > 0 && plusCount == 0 {
		plusCount = 1
	}
	if stats.Deletions > 0 && minusCount == 0 {
		minusCount = 1
	}

	// Build the bar
	plusBar := ui.AddedLineStyle.Render(strings.Repeat("+", plusCount))
	minusBar := ui.RemovedLineStyle.Render(strings.Repeat("-", minusCount))

	fmt.Printf(" %s | %d %s%s\n", filename, total, plusBar, minusBar)
}

// PrintDiff prints a unified diff between original and modified text
func PrintDiff(original, modified string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(original, modified, false)
	diffs = dmp.DiffCleanupSemantic(diffs)

	for _, d := range diffs {
		lines := strings.Split(d.Text, "\n")
		for i, line := range lines {
			// Skip empty trailing line from split
			if i == len(lines)-1 && line == "" {
				continue
			}

			switch d.Type {
			case diffmatchpatch.DiffInsert:
				fmt.Println(ui.AddedLineStyle.Render("+ " + line))
			case diffmatchpatch.DiffDelete:
				fmt.Println(ui.RemovedLineStyle.Render("- " + line))
			case diffmatchpatch.DiffEqual:
				fmt.Println(ui.ContextLineStyle.Render("  " + line))
			}
		}
	}
}

// PrintSummary prints a summary of changes
func PrintSummary(stats Stats) {
	fmt.Printf("\n %s, %s\n",
		ui.AddedLineStyle.Render(fmt.Sprintf("%d insertions(+)", stats.Insertions)),
		ui.RemovedLineStyle.Render(fmt.Sprintf("%d deletions(-)", stats.Deletions)),
	)
}
