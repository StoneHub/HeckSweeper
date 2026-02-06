package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stonehub/hecksweeper/internal/constants"
	"github.com/stonehub/hecksweeper/internal/daily"
	"github.com/stonehub/hecksweeper/internal/game"
	"github.com/stonehub/hecksweeper/internal/storage"
)

// Model represents the UI state for bubbletea
type Model struct {
	run            *game.Run
	runState       string // Current run state (title, playing, cleared, dead, summary, leaderboard)
	lastResult     game.FloorResult
	powerUpChoices []game.PowerUpDef // Current power-up options being offered
	glyphs         constants.GlyphSet
	seed           int64
	useUnicode     bool
	quitting       bool

	// Leaderboard and persistence
	store              *storage.Store
	playerID           string
	playerName         string
	leaderboardEntries []storage.LeaderboardEntry
	leaderboardDate    string
	leaderboardTab     string // "daily" or "alltime"
}

// NewModel creates a new UI model starting at the title screen
func NewModel(seed int64, useUnicode bool) Model {
	glyphs := constants.ASCIIGlyphs
	if useUnicode {
		glyphs = constants.UnicodeGlyphs
	}

	return Model{
		runState:   constants.RunStateTitle,
		glyphs:     glyphs,
		seed:       seed,
		useUnicode: useUnicode,
		quitting:   false,
	}
}

// NewModelWithStore creates a model with persistent storage (for server mode)
func NewModelWithStore(seed int64, useUnicode bool, store *storage.Store, playerID, playerName string) Model {
	m := NewModel(seed, useUnicode)
	m.store = store
	m.playerID = playerID
	m.playerName = playerName
	return m
}

// Init initializes the model (required by bubbletea)
func (m Model) Init() tea.Cmd {
	return nil
}

// startNewRun begins a new roguelite run and starts floor 1
func (m *Model) startNewRun(isDaily bool) {
	seed := m.seed
	if isDaily {
		seed = daily.Seed()
	}
	m.run = game.NewRun(seed, m.useUnicode, isDaily)
	m.run.StartNextFloor()
	m.runState = constants.RunStatePlaying
}

// submitScore saves the run score to the store (if available)
func (m *Model) submitScore() {
	if m.store == nil || m.run == nil {
		return
	}

	score := storage.Score{
		PlayerID:   m.playerID,
		PlayerName: m.playerName,
		Date:       daily.DateString(),
		Score:      m.run.TotalScore,
		Floors:     m.run.TotalFloorsCleared(),
		DurationMs: m.run.TotalTime().Milliseconds(),
		Seed:       m.run.Seed,
		IsDaily:    m.run.IsDaily,
	}
	m.store.SubmitScore(score)
}
