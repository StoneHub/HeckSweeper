package game

// Cell represents a single cell on the board
type Cell struct {
	HasMonster bool
	Revealed   bool
	Flagged    bool
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

// Game represents the complete game state
type Game struct {
	Board         *Board
	CursorPos     Position
	State         string
	MoveCount     int
	FlagCount     int
	Seed          int64
	UseUnicode    bool
}
