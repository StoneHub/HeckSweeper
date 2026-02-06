package constants

const (
	// Default board dimensions
	DefaultWidth  = 24
	DefaultHeight = 16

	// Monster density (percentage)
	DefaultMonsterDensity = 0.18

	// Default seed for reproducible games
	DefaultSeed = int64(42)

	// Game states (single floor)
	GameStatePlaying = "playing"
	GameStateWon     = "won"
	GameStateLost    = "lost"
	GameStateMenu    = "menu"

	// Run states (roguelite progression)
	RunStateTitle   = "title"
	RunStatePlaying = "run_playing" // Active floor
	RunStateCleared = "run_cleared" // Floor cleared, showing stats
	RunStatePowerUp = "run_powerup" // Choosing a power-up
	RunStateDead    = "run_dead"    // Hit a monster, run over
	RunStateSummary = "run_summary" // Final run stats
)

// Glyphs for rendering
type GlyphSet struct {
	Hidden      string
	Flagged     string
	Monster     string
	Cursor      string
	BorderHoriz string
	BorderVert  string
	BorderTL    string
	BorderTR    string
	BorderBL    string
	BorderBR    string
}

var (
	// Unicode glyph set (fancy)
	UnicodeGlyphs = GlyphSet{
		Hidden:      "·",
		Flagged:     "⚑",
		Monster:     "☠",
		Cursor:      "▓",
		BorderHoriz: "─",
		BorderVert:  "│",
		BorderTL:    "┌",
		BorderTR:    "┐",
		BorderBL:    "└",
		BorderBR:    "┘",
	}

	// ASCII glyph set (compatible)
	ASCIIGlyphs = GlyphSet{
		Hidden:      ".",
		Flagged:     "F",
		Monster:     "X",
		Cursor:      "#",
		BorderHoriz: "-",
		BorderVert:  "|",
		BorderTL:    "+",
		BorderTR:    "+",
		BorderBL:    "+",
		BorderBR:    "+",
	}
)
