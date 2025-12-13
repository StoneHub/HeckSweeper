package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette - dark dungeon theme
	colorBackground = lipgloss.Color("#1a1a1a")
	colorBorder     = lipgloss.Color("#4a4a4a")
	colorCursor     = lipgloss.Color("#00ff00")
	colorHidden     = lipgloss.Color("#3a3a3a")
	colorRevealed   = lipgloss.Color("#7a7a7a")
	colorFlag       = lipgloss.Color("#ffaa00")
	colorMonster    = lipgloss.Color("#ff0000")
	colorSafe       = lipgloss.Color("#00aa00")
	colorThreat     = lipgloss.Color("#ffff00")
	colorTitle      = lipgloss.Color("#00ffff")
	colorStats      = lipgloss.Color("#aaaaaa")

	// Threat level colors (0-8 monsters nearby)
	threatColors = []lipgloss.Color{
		lipgloss.Color("#00aa00"), // 0 - green
		lipgloss.Color("#00cccc"), // 1 - cyan
		lipgloss.Color("#00ff00"), // 2 - bright green
		lipgloss.Color("#ffff00"), // 3 - yellow
		lipgloss.Color("#ffaa00"), // 4 - orange
		lipgloss.Color("#ff6600"), // 5 - dark orange
		lipgloss.Color("#ff3300"), // 6 - red-orange
		lipgloss.Color("#ff0000"), // 7 - red
		lipgloss.Color("#aa0000"), // 8 - dark red
	}

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(colorTitle).
			Bold(true).
			Padding(0, 1)

	statsStyle = lipgloss.NewStyle().
			Foreground(colorStats).
			Padding(0, 1)

	boardContainerStyle = lipgloss.NewStyle().
				Padding(1, 2)

	cellStyle = lipgloss.NewStyle().
			Width(2).
			Align(lipgloss.Center)

	hiddenCellStyle = cellStyle.Copy().
			Foreground(colorHidden)

	revealedCellStyle = cellStyle.Copy().
				Foreground(colorRevealed)

	flaggedCellStyle = cellStyle.Copy().
				Foreground(colorFlag).
				Bold(true)

	monsterCellStyle = cellStyle.Copy().
				Foreground(colorMonster).
				Bold(true)

	cursorStyle = cellStyle.Copy().
			Foreground(colorCursor).
			Bold(true).
			Background(lipgloss.Color("#003300"))

	borderStyle = lipgloss.NewStyle().
			Foreground(colorBorder)

	gameOverStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Bold(true).
			Padding(1).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#ff0000"))

	victoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true).
			Padding(1).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#00ff00"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true).
			Padding(1, 0)
)

// GetThreatColor returns the color for a given threat level
func GetThreatColor(threat int) lipgloss.Color {
	if threat < 0 || threat >= len(threatColors) {
		return colorRevealed
	}
	return threatColors[threat]
}

// StyledCell returns a styled cell based on its state
func StyledCell(content string, style lipgloss.Style) string {
	return style.Render(content)
}
