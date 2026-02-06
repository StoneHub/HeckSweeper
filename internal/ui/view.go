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

	switch m.runState {
	case constants.RunStateTitle:
		return m.renderTitle()
	case constants.RunStatePlaying:
		return m.renderPlaying()
	case constants.RunStateCleared:
		return m.renderCleared()
	case constants.RunStatePowerUp:
		return m.renderPowerUpSelection()
	case constants.RunStateDeathReveal:
		return m.renderDeathReveal()
	case constants.RunStateDead:
		return m.renderDead()
	case constants.RunStateSummary:
		return m.renderSummary()
	case constants.RunStateLeaderboard:
		return m.renderLeaderboard()
	case constants.RunStateTransition:
		return m.renderTransition()
	}

	return ""
}

// renderPlaying renders the active gameplay screen
func (m Model) renderPlaying() string {
	var b strings.Builder

	g := m.run.CurrentGame
	cfg := game.GetFloorConfig(m.run.FloorNum)

	// Floor header
	floorHeader := fmt.Sprintf(
		"Floor %d  |  Board: %d×%d  |  Score: %d",
		m.run.FloorNum, cfg.Width, cfg.Height, m.run.TotalScore,
	)
	b.WriteString(titleStyle.Render(floorHeader))
	b.WriteString("\n\n")

	// Stats bar
	stats := fmt.Sprintf(
		"Moves: %d  |  Flags: %d  |  Remaining: %d  |  Seed: %d",
		g.MoveCount,
		g.FlagCount,
		g.Board.RemainingSafe,
		m.seed,
	)
	b.WriteString(statsStyle.Render(stats))
	b.WriteString("\n\n")

	// Render the board
	b.WriteString(m.renderBoard())
	b.WriteString("\n")

	// Active power-ups
	if names := m.run.PowerUpNames(); len(names) > 0 {
		puStr := fmt.Sprintf("Power-ups: %s", strings.Join(names, ", "))
		b.WriteString(statsStyle.Render(puStr))
		b.WriteString("\n")
	}

	// Help text (show undo key when Time Warp is available)
	helpText := "Arrow/WASD: Move  |  Space/Enter: Reveal  |  F: Flag  |  Esc: Quit run"
	if m.run.HasPowerUp(game.PowerUpTimeWarp) && m.run.UndoUsedFloor != m.run.FloorNum {
		helpText = "Arrow/WASD: Move  |  Space/Enter: Reveal  |  F: Flag  |  U: Undo  |  Esc: Quit"
	}
	b.WriteString(helpStyle.Render(helpText))

	return boardContainerStyle.Render(b.String())
}

// renderDeathReveal renders the board during mine cascade animation
func (m Model) renderDeathReveal() string {
	var b strings.Builder

	header := "  YOU HIT A MINE!  "
	b.WriteString(gameOverStyle.Render(header))
	b.WriteString("\n\n")

	// Show the board with mines being revealed
	b.WriteString(m.renderBoard())
	b.WriteString("\n")

	remaining := len(m.deathMonsters) - m.deathRevealIdx
	if remaining > 0 {
		b.WriteString(statsStyle.Render(fmt.Sprintf("  Revealing mines... %d remaining", remaining)))
	} else {
		b.WriteString(helpStyle.Render("Press Enter for run summary  |  Esc to quit"))
	}

	return boardContainerStyle.Render(b.String())
}

// renderTransition renders the floor transition animation
func (m Model) renderTransition() string {
	var b strings.Builder

	dots := strings.Repeat(".", (m.transitionTicks%3)+1)
	header := fmt.Sprintf("  DESCENDING%s  ", dots)
	b.WriteString(titleStyle.Render(header))
	b.WriteString("\n\n")

	cfg := game.GetFloorConfig(m.run.FloorNum)
	info := strings.Join([]string{
		fmt.Sprintf("  Floor %d", m.run.FloorNum),
		fmt.Sprintf("  Board: %d×%d", cfg.Width, cfg.Height),
		fmt.Sprintf("  Score: %d", m.run.TotalScore),
	}, "\n")
	b.WriteString(statsStyle.Render(info))
	b.WriteString("\n\n")

	if names := m.run.PowerUpNames(); len(names) > 0 {
		puStr := fmt.Sprintf("  Power-ups: %s", strings.Join(names, ", "))
		b.WriteString(statsStyle.Render(puStr))
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("Press any key to skip"))

	return boardContainerStyle.Render(b.String())
}

// renderBoard renders the game board
func (m Model) renderBoard() string {
	var b strings.Builder

	g := m.run.CurrentGame

	// Top border
	b.WriteString(borderStyle.Render(m.glyphs.BorderTL))
	for x := 0; x < g.Board.Width; x++ {
		b.WriteString(borderStyle.Render(m.glyphs.BorderHoriz + m.glyphs.BorderHoriz))
	}
	b.WriteString(borderStyle.Render(m.glyphs.BorderTR))
	b.WriteString("\n")

	// Board cells
	for y := 0; y < g.Board.Height; y++ {
		b.WriteString(borderStyle.Render(m.glyphs.BorderVert))

		for x := 0; x < g.Board.Width; x++ {
			cell := g.Board.GetCell(x, y)
			isCursor := (x == g.CursorPos.X && y == g.CursorPos.Y)
			cellStr := m.renderCell(cell, isCursor)
			b.WriteString(cellStr)
		}

		b.WriteString(borderStyle.Render(m.glyphs.BorderVert))
		b.WriteString("\n")
	}

	// Bottom border
	b.WriteString(borderStyle.Render(m.glyphs.BorderBL))
	for x := 0; x < g.Board.Width; x++ {
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

	if cell.Revealed {
		if cell.HasMonster {
			content = m.glyphs.Monster + " "
			style = monsterCellStyle
		} else if cell.ThreatCount > 0 {
			content = fmt.Sprintf("%d ", cell.ThreatCount)
			style = cellStyle.Foreground(GetThreatColor(cell.ThreatCount))
		} else {
			content = "  "
			style = revealedCellStyle
		}
	} else if cell.Flagged {
		content = m.glyphs.Flagged + " "
		style = flaggedCellStyle
	} else {
		content = m.glyphs.Hidden + " "
		style = hiddenCellStyle
	}

	if isCursor && m.run != nil && m.run.CurrentGame.State == constants.GameStatePlaying {
		style = style.
			Background(lipgloss.Color("#003300")).
			Underline(true)
	}

	return style.Render(content)
}
