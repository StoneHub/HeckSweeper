package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/stonehub/hecksweeper/internal/game"
)

var (
	rarityColors = map[game.Rarity]lipgloss.Color{
		game.Common:    lipgloss.Color("#aaaaaa"),
		game.Uncommon:  lipgloss.Color("#00ccff"),
		game.Rare:      lipgloss.Color("#ffaa00"),
		game.Legendary: lipgloss.Color("#ff44ff"),
	}

	powerUpBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4a4a4a")).
			Padding(1, 2).
			Width(50)
)

// renderPowerUpSelection renders the power-up choice screen
func (m Model) renderPowerUpSelection() string {
	var b strings.Builder

	header := fmt.Sprintf("  FLOOR %d CLEARED — CHOOSE A POWER-UP  ", m.lastResult.Floor)
	b.WriteString(victoryStyle.Render(header))
	b.WriteString("\n\n")

	// Score context
	scoreInfo := fmt.Sprintf("  Floor score: %d  |  Total: %d", m.lastResult.Score, m.run.TotalScore)
	b.WriteString(statsStyle.Render(scoreInfo))
	b.WriteString("\n\n")

	// Render each power-up choice
	for i, pu := range m.powerUpChoices {
		color := rarityColors[pu.Rarity]
		rarityStr := game.RarityName(pu.Rarity)

		optStyle := powerUpBoxStyle.BorderForeground(color)

		content := fmt.Sprintf("[%d] %s (%s)\n    %s", i+1, pu.Name, rarityStr, pu.Description)
		b.WriteString(optStyle.Render(content))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Show current power-ups
	if names := m.run.PowerUpNames(); len(names) > 0 {
		b.WriteString(statsStyle.Render(fmt.Sprintf("  Current: %s", strings.Join(names, ", "))))
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("Press 1, 2, or 3 to choose  |  Esc to end run"))

	return boardContainerStyle.Render(b.String())
}
