package ui

import (
	"fmt"
	"strings"
	"time"
)

// renderCleared renders the floor cleared screen
func (m Model) renderCleared() string {
	var b strings.Builder

	r := m.lastResult

	header := fmt.Sprintf("  FLOOR %d CLEARED!  ", r.Floor)
	b.WriteString(victoryStyle.Render(header))
	b.WriteString("\n\n")

	stats := strings.Join([]string{
		fmt.Sprintf("  Cells revealed: %d / %d", r.CellsRevealed, r.TotalSafe),
		fmt.Sprintf("  Moves:          %d", r.Moves),
		fmt.Sprintf("  Time:           %s", formatDuration(r.TimeSpent)),
		fmt.Sprintf("  Floor score:    %d", r.Score),
		"",
		fmt.Sprintf("  Total score:    %d", m.run.TotalScore),
	}, "\n")
	b.WriteString(statsStyle.Render(stats))
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("Press Enter to descend to next floor  |  Esc to end run"))

	return boardContainerStyle.Render(b.String())
}

// renderDead renders the death screen
func (m Model) renderDead() string {
	var b strings.Builder

	header := fmt.Sprintf("  YOU DIED ON FLOOR %d  ", m.run.FloorNum)
	b.WriteString(gameOverStyle.Render(header))
	b.WriteString("\n\n")

	r := m.lastResult
	stats := strings.Join([]string{
		fmt.Sprintf("  Cells revealed: %d / %d", r.CellsRevealed, r.TotalSafe),
		fmt.Sprintf("  Moves:          %d", r.Moves),
		fmt.Sprintf("  Time:           %s", formatDuration(r.TimeSpent)),
		fmt.Sprintf("  Floor score:    %d", r.Score),
	}, "\n")
	b.WriteString(statsStyle.Render(stats))
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("Press Enter for run summary  |  Esc to quit"))

	return boardContainerStyle.Render(b.String())
}

// renderSummary renders the final run summary screen
func (m Model) renderSummary() string {
	var b strings.Builder

	header := "  RUN SUMMARY  "
	b.WriteString(titleStyle.Render(header))
	b.WriteString("\n\n")

	// Overall stats
	floorsCleared := m.run.TotalFloorsCleared()
	totalTime := m.run.TotalTime()

	overall := strings.Join([]string{
		fmt.Sprintf("  Floors cleared: %d", floorsCleared),
		fmt.Sprintf("  Total score:    %d", m.run.TotalScore),
		fmt.Sprintf("  Total time:     %s", formatDuration(totalTime)),
		fmt.Sprintf("  Seed:           %d", m.run.Seed),
	}, "\n")
	b.WriteString(statsStyle.Render(overall))
	b.WriteString("\n\n")

	// Per-floor breakdown
	if len(m.run.FloorHistory) > 0 {
		b.WriteString(statsStyle.Render("  Floor breakdown:"))
		b.WriteString("\n")

		for _, f := range m.run.FloorHistory {
			status := "Cleared"
			if !f.Cleared {
				status = "Died"
			}
			line := fmt.Sprintf("    Floor %d: %s  |  Score: %d  |  Moves: %d  |  %s",
				f.Floor, status, f.Score, f.Moves, formatDuration(f.TimeSpent))
			b.WriteString(statsStyle.Render(line))
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	b.WriteString(helpStyle.Render("Press N for new run  |  Esc to return to title"))

	return boardContainerStyle.Render(b.String())
}

// formatDuration formats a duration as M:SS
func formatDuration(d time.Duration) string {
	secs := int(d.Seconds())
	mins := secs / 60
	secs = secs % 60
	return fmt.Sprintf("%d:%02d", mins, secs)
}
