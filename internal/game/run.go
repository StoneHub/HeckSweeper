package game

import (
	"math"
	"time"

	"github.com/stonehub/hecksweeper/internal/constants"
)

// Floor scaling tables
var floorConfigs = []FloorConfig{
	{Width: 8, Height: 6, MonsterDensity: 0.10, FloorNumber: 1},
	{Width: 12, Height: 8, MonsterDensity: 0.12, FloorNumber: 2},
	{Width: 16, Height: 10, MonsterDensity: 0.14, FloorNumber: 3},
	{Width: 18, Height: 12, MonsterDensity: 0.16, FloorNumber: 4},
	{Width: 20, Height: 14, MonsterDensity: 0.18, FloorNumber: 5},
	{Width: 22, Height: 14, MonsterDensity: 0.20, FloorNumber: 6},
	{Width: 24, Height: 16, MonsterDensity: 0.22, FloorNumber: 7},
}

// GetFloorConfig returns the configuration for a given floor number.
// Floors beyond the table scale density incrementally.
func GetFloorConfig(floor int) FloorConfig {
	if floor <= 0 {
		floor = 1
	}

	if floor <= len(floorConfigs) {
		return floorConfigs[floor-1]
	}

	// Beyond predefined floors: max board size, increasing density
	last := floorConfigs[len(floorConfigs)-1]
	extraFloors := floor - len(floorConfigs)
	density := last.MonsterDensity + float64(extraFloors)*0.01
	if density > 0.35 {
		density = 0.35
	}

	return FloorConfig{
		Width:          last.Width,
		Height:         last.Height,
		MonsterDensity: density,
		FloorNumber:    floor,
	}
}

// NewRun creates a new roguelite run
func NewRun(seed int64, useUnicode bool, isDaily bool) *Run {
	return &Run{
		Seed:         seed,
		FloorNum:     0,
		TotalScore:   0,
		FloorHistory: []FloorResult{},
		StartedAt:    time.Now(),
		IsDaily:      isDaily,
		UseUnicode:   useUnicode,
	}
}

// StartNextFloor advances to the next floor and creates its game
func (r *Run) StartNextFloor() {
	r.FloorNum++
	cfg := GetFloorConfig(r.FloorNum)

	// Apply Lucky Charm: reduce density by 2% per instance
	density := cfg.MonsterDensity
	for _, pu := range r.PowerUps {
		if pu.ID == PowerUpLuckyCharm && !pu.Consumed {
			density -= 0.02
		}
	}
	if density < 0.05 {
		density = 0.05 // Minimum density floor
	}

	// Derive a floor-specific seed so each floor is unique but deterministic
	floorSeed := r.Seed + int64(r.FloorNum)*7919

	r.CurrentGame = NewGame(cfg.Width, cfg.Height, floorSeed, r.UseUnicode)
	// Re-generate board with adjusted density if Lucky Charm is active
	if density != cfg.MonsterDensity {
		r.CurrentGame.Board = GenerateBoard(cfg.Width, cfg.Height, density, floorSeed)
	}

	r.FloorStarted = time.Now()

	// Apply floor-start power-ups (Scout, X-Ray, Dowsing)
	ApplyFloorStartPowerUps(r)
}

// CompleteFloor records the result of the current floor and returns it
func (r *Run) CompleteFloor() FloorResult {
	g := r.CurrentGame
	cfg := GetFloorConfig(r.FloorNum)

	totalSafe := cfg.Width*cfg.Height - len(g.Board.Monsters)
	cellsRevealed := totalSafe - g.Board.RemainingSafe
	cleared := g.Board.IsWon()
	elapsed := time.Since(r.FloorStarted)

	score := CalculateFloorScore(cellsRevealed, g.MoveCount, elapsed, r.FloorNum)

	result := FloorResult{
		Floor:         r.FloorNum,
		Moves:         g.MoveCount,
		CellsRevealed: cellsRevealed,
		TotalSafe:     totalSafe,
		Cleared:       cleared,
		TimeSpent:     elapsed,
		Score:         score,
	}

	r.FloorHistory = append(r.FloorHistory, result)
	r.TotalScore += score

	return result
}

// CalculateFloorScore computes the score for a single floor
func CalculateFloorScore(cellsRevealed, moves int, elapsed time.Duration, floor int) int {
	baseScore := float64(cellsRevealed) * 10

	// Speed bonus: up to 500 points, decaying over 2 minutes
	timeLimitSec := 120.0
	elapsedSec := elapsed.Seconds()
	speedBonus := math.Max(0, (timeLimitSec-elapsedSec)*500/timeLimitSec)

	// Efficiency bonus: fewer moves = more bonus
	if moves == 0 {
		moves = 1
	}
	efficiency := float64(cellsRevealed) / float64(moves)
	efficiencyBonus := efficiency * 100

	// Floor multiplier: deeper floors are worth more
	floorMultiplier := 1.0 + float64(floor-1)*0.5

	total := (baseScore + speedBonus + efficiencyBonus) * floorMultiplier

	return int(total)
}

// TotalFloorsCleared returns how many floors were successfully cleared
func (r *Run) TotalFloorsCleared() int {
	count := 0
	for _, f := range r.FloorHistory {
		if f.Cleared {
			count++
		}
	}
	return count
}

// TotalTime returns the total run duration
func (r *Run) TotalTime() time.Duration {
	return time.Since(r.StartedAt)
}

// IsActive returns true if the run is still in progress
func (r *Run) IsActive() bool {
	if r.CurrentGame == nil {
		return true
	}
	return r.CurrentGame.State != constants.GameStateLost
}

// SaveSnapshot saves the current board state for Time Warp undo
func (r *Run) SaveSnapshot() {
	g := r.CurrentGame
	if g == nil {
		return
	}
	cellsCopy := make([]Cell, len(g.Board.Cells))
	copy(cellsCopy, g.Board.Cells)
	monstersCopy := make(map[int]bool, len(g.Board.Monsters))
	for k, v := range g.Board.Monsters {
		monstersCopy[k] = v
	}
	r.LastSnapshot = &BoardSnapshot{
		Cells:         cellsCopy,
		Monsters:      monstersCopy,
		RemainingSafe: g.Board.RemainingSafe,
		CursorPos:     g.CursorPos,
		GameState:     g.State,
		MoveCount:     g.MoveCount,
		FlagCount:     g.FlagCount,
	}
}

// RestoreSnapshot restores the board to the saved snapshot state
func (r *Run) RestoreSnapshot() {
	if r.LastSnapshot == nil || r.CurrentGame == nil {
		return
	}
	snap := r.LastSnapshot
	g := r.CurrentGame
	g.Board.Cells = snap.Cells
	g.Board.Monsters = snap.Monsters
	g.Board.RemainingSafe = snap.RemainingSafe
	g.CursorPos = snap.CursorPos
	g.State = snap.GameState
	g.MoveCount = snap.MoveCount
	g.FlagCount = snap.FlagCount
	r.LastSnapshot = nil
}
