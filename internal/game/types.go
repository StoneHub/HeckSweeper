package game

import "time"

// Cell represents a single cell on the board
type Cell struct {
	HasMonster  bool
	Revealed    bool
	Flagged     bool
	ThreatCount int // Number of adjacent monsters
}

// Board represents the game board state
type Board struct {
	Width         int
	Height        int
	Cells         []Cell
	Monsters      map[int]bool // Set of monster positions
	RemainingSafe int          // Safe cells left to reveal
}

// Position represents a 2D coordinate
type Position struct {
	X int
	Y int
}

// Game represents a single floor's game state
type Game struct {
	Board      *Board
	CursorPos  Position
	State      string
	MoveCount  int
	FlagCount  int
	Seed       int64
	UseUnicode bool
}

// FloorConfig holds the parameters for a specific floor
type FloorConfig struct {
	Width          int
	Height         int
	MonsterDensity float64
	FloorNumber    int
}

// FloorResult records the outcome of a completed floor
type FloorResult struct {
	Floor         int
	Moves         int
	CellsRevealed int
	TotalSafe     int
	Cleared       bool
	TimeSpent     time.Duration
	Score         int
}

// Run represents a multi-floor roguelite run
type Run struct {
	Seed         int64
	FloorNum     int
	TotalScore   int
	PowerUps     []ActivePowerUp
	CurrentGame  *Game
	FloorHistory []FloorResult
	StartedAt    time.Time
	FloorStarted time.Time
	IsDaily      bool
	UseUnicode   bool
}
