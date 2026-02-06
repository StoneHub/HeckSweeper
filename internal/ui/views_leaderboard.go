package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/stonehub/hecksweeper/internal/storage"
)

var (
	leaderboardHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00ffff")).
				Bold(true).
				Padding(1).
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#00ffff"))

	rankStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffaa00")).
			Bold(true).
			Width(4)

	nameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#aaaaaa")).
			Width(18)

	scoreStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Width(8).
			Align(lipgloss.Right)

	floorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00cccc")).
			Width(8).
			Align(lipgloss.Right)

	youStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00")).
			Bold(true)
)

// renderLeaderboard renders the leaderboard screen
func (m Model) renderLeaderboard() string {
	var b strings.Builder

	// Header
	date := m.leaderboardDate
	headerText := fmt.Sprintf("  DAILY LEADERBOARD  %s  ", date)
	if m.leaderboardTab == "alltime" {
		headerText = "  ALL-TIME LEADERBOARD  "
	}
	b.WriteString(leaderboardHeaderStyle.Render(headerText))
	b.WriteString("\n\n")

	// Column headers
	colHeader := fmt.Sprintf("  %s%s%s%s",
		rankStyle.Render("#"),
		nameStyle.Render("Player"),
		scoreStyle.Render("Score"),
		floorStyle.Render("Floors"),
	)
	b.WriteString(statsStyle.Render(colHeader))
	b.WriteString("\n")
	b.WriteString(statsStyle.Render("  "+strings.Repeat("─", 38)))
	b.WriteString("\n")

	// Entries
	entries := m.leaderboardEntries
	if len(entries) == 0 {
		b.WriteString(helpStyle.Render("  No scores yet. Be the first!"))
		b.WriteString("\n")
	} else {
		for _, e := range entries {
			line := renderLeaderboardEntry(e)
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	// Tab bar
	tabs := "  [D] Daily  [A] All-time  [B] Back"
	b.WriteString(helpStyle.Render(tabs))

	return boardContainerStyle.Render(b.String())
}

func renderLeaderboardEntry(e storage.LeaderboardEntry) string {
	rankStr := fmt.Sprintf("%d.", e.Rank)
	name := e.PlayerName

	if e.IsYou {
		name = youStyle.Render(name + " (you)")
	}

	return fmt.Sprintf("  %s%s%s%s",
		rankStyle.Render(rankStr),
		nameStyle.Render(name),
		scoreStyle.Render(fmt.Sprintf("%d", e.Score)),
		floorStyle.Render(fmt.Sprintf("%d", e.Floors)),
	)
}
