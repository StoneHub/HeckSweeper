package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stonehub/hecksweeper/internal/constants"
	"github.com/stonehub/hecksweeper/internal/game"
)

// Update handles messages and updates the model (required by bubbletea)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys (work in any state)
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.quitting = true
		return m, tea.Quit
	}

	// Menu state keys
	if m.game.State == constants.GameStateMenu {
		switch msg.String() {
		case "n":
			// New game
			m.game = game.NewGame(m.width, m.height, m.seed, m.useUnicode)
		}
		return m, nil
	}

	// Game over state keys
	if m.game.State == constants.GameStateWon || m.game.State == constants.GameStateLost {
		switch msg.String() {
		case "r":
			// Restart game
			m.game = game.NewGame(m.width, m.height, m.seed, m.useUnicode)
		case "n":
			// New game with different seed
			m.seed++
			m.game = game.NewGame(m.width, m.height, m.seed, m.useUnicode)
		}
		return m, nil
	}

	// Playing state keys
	if m.game.State == constants.GameStatePlaying {
		switch msg.String() {
		// Movement - Arrow keys
		case "up", "k", "w":
			m.game.MoveCursor(0, -1)
		case "down", "j", "s":
			m.game.MoveCursor(0, 1)
		case "left", "h", "a":
			m.game.MoveCursor(-1, 0)
		case "right", "l", "d":
			m.game.MoveCursor(1, 0)

		// Actions
		case " ", "enter":
			// Reveal cell
			m.game.RevealAtCursor()
		case "f":
			// Toggle flag
			m.game.ToggleFlag()
		}
	}

	return m, nil
}
