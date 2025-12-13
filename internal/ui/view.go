package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/stonehub/hecksweeper/internal/constants"
	"github.com/stonehub/hecksweeper/internal/game"
)

// View renders the UI (required by bubbletea)
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Title
	title := "╔═══════════════════════════╗\n"
	title += "║     H E C K S W E E P E R     ║\n"
	title += "╚═══════════════════════════╝"
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Stats bar
	stats := fmt.Sprintf(
		"Moves: %d  |  Flags: %d  |  Remaining: %d  |  Seed: %d",
		m.game.MoveCount,
		m.game.FlagCount,
		m.game.Board.RemainingSafe,
		m.seed,
	)
	b.WriteString(statsStyle.Render(stats))
	b.WriteString("\n\n")

	// Render the board
	b.WriteString(m.renderBoard())
	b.WriteString("\n")

	// Game state messages
	switch m.game.State {
	case constants.GameStateWon:
		msg := "  🎉 VICTORY! All safe cells revealed! 🎉  \n"
		msg += fmt.Sprintf("  Completed in %d moves  ", m.game.MoveCount)
		b.WriteString(victoryStyle.Render(msg))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press 'r' to restart, 'n' for new game, 'q' to quit"))

	case constants.GameStateLost:
		msg := "  💀 GAME OVER! You hit a monster! 💀  \n"
		msg += fmt.Sprintf("  Survived %d moves  ", m.game.MoveCount)
		b.WriteString(gameOverStyle.Render(msg))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press 'r' to restart, 'n' for new game, 'q' to quit"))

	case constants.GameStatePlaying:
		b.WriteString(helpStyle.Render("Arrow/WASD: Move  |  Space/Enter: Reveal  |  F: Flag  |  Q: Quit"))
	}

	return boardContainerStyle.Render(b.String())
}

// renderBoard renders the game board
func (m Model) renderBoard() string {
	var b strings.Builder

	// Top border
	b.WriteString(borderStyle.Render(m.glyphs.BorderTL))
	for x := 0; x < m.game.Board.Width; x++ {
		b.WriteString(borderStyle.Render(m.glyphs.BorderHoriz + m.glyphs.BorderHoriz))
	}
	b.WriteString(borderStyle.Render(m.glyphs.BorderTR))
	b.WriteString("\n")

	// Board cells
	for y := 0; y < m.game.Board.Height; y++ {
		// Left border
		b.WriteString(borderStyle.Render(m.glyphs.BorderVert))

		// Cells in this row
		for x := 0; x < m.game.Board.Width; x++ {
			cell := m.game.Board.GetCell(x, y)
			isCursor := (x == m.game.CursorPos.X && y == m.game.CursorPos.Y)

			cellStr := m.renderCell(cell, isCursor)
			b.WriteString(cellStr)
		}

		// Right border
		b.WriteString(borderStyle.Render(m.glyphs.BorderVert))
		b.WriteString("\n")
	}

	// Bottom border
	b.WriteString(borderStyle.Render(m.glyphs.BorderBL))
	for x := 0; x < m.game.Board.Width; x++ {
		b.WriteString(borderStyle.Render(m.glyphs.BorderHoriz + m.glyphs.BorderHoriz))
	}
	b.WriteString(borderStyle.Render(m.glyphs.BorderBR))

	return b.String()
}

// renderCell renders a single cell
func (m Model) renderCell(cell *game.Cell, isCursor bool) string {
	if cell == nil {
		return "  "
	}

	var content string
	var style lipgloss.Style

	// Determine cell content and style
	if cell.Revealed {
		if cell.HasMonster {
			// Monster cell
			content = m.glyphs.Monster + " "
			style = monsterCellStyle
		} else if cell.ThreatCount > 0 {
			// Revealed cell with threat count
			content = fmt.Sprintf("%d ", cell.ThreatCount)
			style = cellStyle.Copy().Foreground(GetThreatColor(cell.ThreatCount))
		} else {
			// Safe cell (no adjacent monsters)
			content = "  "
			style = revealedCellStyle
		}
	} else if cell.Flagged {
		// Flagged cell
		content = m.glyphs.Flagged + " "
		style = flaggedCellStyle
	} else {
		// Hidden cell
		content = m.glyphs.Hidden + " "
		style = hiddenCellStyle
	}

	// Apply cursor highlight
	if isCursor && m.game.State == constants.GameStatePlaying {
		style = style.Copy().
			Background(lipgloss.Color("#003300")).
			Underline(true)
	}

	return style.Render(content)
}
