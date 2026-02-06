package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stonehub/hecksweeper/internal/constants"
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
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	}

	switch m.runState {
	case constants.RunStateTitle:
		return m.handleTitleKeys(msg)
	case constants.RunStatePlaying:
		return m.handlePlayingKeys(msg)
	case constants.RunStateCleared:
		return m.handleClearedKeys(msg)
	case constants.RunStateDead:
		return m.handleDeadKeys(msg)
	case constants.RunStateSummary:
		return m.handleSummaryKeys(msg)
	}

	return m, nil
}

// handleTitleKeys handles input on the title screen
func (m Model) handleTitleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n", "enter", " ":
		m.startNewRun(false)
	case "q", "esc":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

// handlePlayingKeys handles input during active gameplay
func (m Model) handlePlayingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	g := m.run.CurrentGame

	switch msg.String() {
	case "q", "esc":
		// Quit to title (abandon run)
		m.runState = constants.RunStateTitle
		return m, nil
	}

	if g.State != constants.GameStatePlaying {
		return m, nil
	}

	switch msg.String() {
	// Movement
	case "up", "k", "w":
		g.MoveCursor(0, -1)
	case "down", "j", "s":
		g.MoveCursor(0, 1)
	case "left", "h", "a":
		g.MoveCursor(-1, 0)
	case "right", "l", "d":
		g.MoveCursor(1, 0)

	// Actions
	case " ", "enter":
		g.RevealAtCursor()
		// Check if floor was cleared or player died
		if g.State == constants.GameStateWon {
			m.lastResult = m.run.CompleteFloor()
			m.runState = constants.RunStateCleared
		} else if g.State == constants.GameStateLost {
			m.lastResult = m.run.CompleteFloor()
			m.runState = constants.RunStateDead
		}
	case "f":
		g.ToggleFlag()
	}

	return m, nil
}

// handleClearedKeys handles input on the floor cleared screen
func (m Model) handleClearedKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ", "n":
		// Advance to next floor
		m.run.StartNextFloor()
		m.runState = constants.RunStatePlaying
	case "q", "esc":
		m.runState = constants.RunStateSummary
	}
	return m, nil
}

// handleDeadKeys handles input on the death screen
func (m Model) handleDeadKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		m.runState = constants.RunStateSummary
	case "q", "esc":
		m.runState = constants.RunStateTitle
	}
	return m, nil
}

// handleSummaryKeys handles input on the run summary screen
func (m Model) handleSummaryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n", "enter", " ":
		// Start a new run with incremented seed
		m.seed++
		m.startNewRun(false)
	case "q", "esc":
		m.runState = constants.RunStateTitle
	}
	return m, nil
}
