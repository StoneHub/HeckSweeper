package ui

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stonehub/hecksweeper/internal/constants"
	"github.com/stonehub/hecksweeper/internal/daily"
	"github.com/stonehub/hecksweeper/internal/game"
)

// Tick message types for animations
type deathRevealTickMsg struct{}
type transitionTickMsg struct{}

func deathRevealTickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return deathRevealTickMsg{}
	})
}

func transitionTickCmd() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg {
		return transitionTickMsg{}
	})
}

// Update handles messages and updates the model (required by bubbletea)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case deathRevealTickMsg:
		return m.handleDeathRevealTick()
	case transitionTickMsg:
		return m.handleTransitionTick()
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
	case constants.RunStatePowerUp:
		return m.handlePowerUpKeys(msg)
	case constants.RunStateDeathReveal:
		return m.handleDeathRevealKeys(msg)
	case constants.RunStateDead:
		return m.handleDeadKeys(msg)
	case constants.RunStateSummary:
		return m.handleSummaryKeys(msg)
	case constants.RunStateLeaderboard:
		return m.handleLeaderboardKeys(msg)
	case constants.RunStateTransition:
		return m.handleTransitionKeys(msg)
	}

	return m, nil
}

// handleTitleKeys handles input on the title screen
func (m Model) handleTitleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n", "enter", " ":
		m.startNewRun(false)
	case "d":
		m.startNewRun(true)
	case "l":
		m.showLeaderboard("daily")
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

	// Undo (Time Warp)
	case "u":
		if m.run.HasPowerUp(game.PowerUpTimeWarp) && m.run.UndoUsedFloor != m.run.FloorNum && m.run.LastSnapshot != nil {
			m.run.RestoreSnapshot()
			m.run.UndoUsedFloor = m.run.FloorNum
		}

	// Actions
	case " ", "enter":
		// Save snapshot for Time Warp undo
		if m.run.HasPowerUp(game.PowerUpTimeWarp) && m.run.UndoUsedFloor != m.run.FloorNum {
			m.run.SaveSnapshot()
		}
		g.RevealAtCursor()
		// Check if floor was cleared or player died
		if g.State == constants.GameStateWon {
			m.lastResult = m.run.CompleteFloor()
			m.runState = constants.RunStateCleared
		} else if g.State == constants.GameStateLost {
			// Check for death-prevention power-ups
			if game.HandleDeathPowerUps(m.run) {
				// Death was prevented, continue playing
			} else {
				m.lastResult = m.run.CompleteFloor()
				// Start death reveal cascade animation
				m.deathMonsters = g.Board.GetUnrevealedMonsterPositions(
					g.CursorPos.X, g.CursorPos.Y,
				)
				m.deathRevealIdx = 0
				m.runState = constants.RunStateDeathReveal
				m.submitScore()
				return m, deathRevealTickCmd()
			}
		}
	case "f":
		g.ToggleFlag()
	}

	return m, nil
}

// handleDeathRevealTick advances the mine cascade animation by one step
func (m Model) handleDeathRevealTick() (tea.Model, tea.Cmd) {
	if m.runState != constants.RunStateDeathReveal {
		return m, nil
	}
	if m.deathRevealIdx < len(m.deathMonsters) {
		pos := m.deathMonsters[m.deathRevealIdx]
		cell := m.run.CurrentGame.Board.GetCell(pos.X, pos.Y)
		if cell != nil {
			cell.Revealed = true
		}
		m.deathRevealIdx++
		if m.deathRevealIdx < len(m.deathMonsters) {
			return m, deathRevealTickCmd()
		}
	}
	// All revealed — wait for key press
	return m, nil
}

// handleDeathRevealKeys handles input during the death cascade animation
func (m Model) handleDeathRevealKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.deathRevealIdx < len(m.deathMonsters) {
		// Still animating — skip to end
		g := m.run.CurrentGame
		for i := m.deathRevealIdx; i < len(m.deathMonsters); i++ {
			pos := m.deathMonsters[i]
			cell := g.Board.GetCell(pos.X, pos.Y)
			if cell != nil {
				cell.Revealed = true
			}
		}
		m.deathRevealIdx = len(m.deathMonsters)
		return m, nil
	}

	// Animation done
	switch msg.String() {
	case "enter", " ":
		m.runState = constants.RunStateDead
	case "q", "esc":
		m.runState = constants.RunStateTitle
	}
	return m, nil
}

// handleTransitionTick advances the floor transition animation
func (m Model) handleTransitionTick() (tea.Model, tea.Cmd) {
	if m.runState != constants.RunStateTransition {
		return m, nil
	}
	m.transitionTicks++
	if m.transitionTicks >= 6 {
		m.runState = constants.RunStatePlaying
		return m, nil
	}
	return m, transitionTickCmd()
}

// handleTransitionKeys allows skipping the transition animation
func (m Model) handleTransitionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.runState = constants.RunStatePlaying
	return m, nil
}

// handleClearedKeys handles input on the floor cleared screen
func (m Model) handleClearedKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ", "n":
		// Generate power-up choices for the player
		existingIDs := []game.PowerUpID{}
		for _, pu := range m.run.PowerUps {
			if !pu.Consumed {
				existingIDs = append(existingIDs, pu.ID)
			}
		}
		rng := rand.New(rand.NewSource(m.run.Seed + int64(m.run.FloorNum)*1009))
		m.powerUpChoices = game.PickPowerUpChoices(rng, 3, existingIDs)
		m.runState = constants.RunStatePowerUp
	case "q", "esc":
		m.submitScore()
		m.runState = constants.RunStateSummary
	}
	return m, nil
}

// handlePowerUpKeys handles input on the power-up selection screen
func (m Model) handlePowerUpKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		if len(m.powerUpChoices) >= 1 {
			m.run.AddPowerUp(m.powerUpChoices[0])
			return m, m.advanceToNextFloor()
		}
	case "2":
		if len(m.powerUpChoices) >= 2 {
			m.run.AddPowerUp(m.powerUpChoices[1])
			return m, m.advanceToNextFloor()
		}
	case "3":
		if len(m.powerUpChoices) >= 3 {
			m.run.AddPowerUp(m.powerUpChoices[2])
			return m, m.advanceToNextFloor()
		}
	case "q", "esc":
		m.submitScore()
		m.runState = constants.RunStateSummary
	}
	return m, nil
}

// advanceToNextFloor starts the next floor with a transition animation
func (m *Model) advanceToNextFloor() tea.Cmd {
	m.run.StartNextFloor()
	m.powerUpChoices = nil
	m.transitionTicks = 0
	m.runState = constants.RunStateTransition
	return transitionTickCmd()
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
		m.seed++
		m.startNewRun(false)
	case "l":
		m.showLeaderboard("daily")
	case "q", "esc":
		m.runState = constants.RunStateTitle
	}
	return m, nil
}

// handleLeaderboardKeys handles input on the leaderboard screen
func (m Model) handleLeaderboardKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "d":
		m.showLeaderboard("daily")
	case "a":
		m.showLeaderboard("alltime")
	case "b", "q", "esc":
		m.runState = constants.RunStateTitle
	}
	return m, nil
}

// showLeaderboard loads and displays leaderboard data
func (m *Model) showLeaderboard(tab string) {
	m.leaderboardTab = tab
	m.leaderboardDate = daily.DateString()

	if m.store != nil {
		if tab == "daily" {
			m.leaderboardEntries = m.store.GetDailyLeaderboard(m.leaderboardDate, 10, m.playerID)
		} else {
			m.leaderboardEntries = m.store.GetAllTimeLeaderboard(10, m.playerID)
		}
	} else {
		m.leaderboardEntries = nil
	}

	m.runState = constants.RunStateLeaderboard
}
