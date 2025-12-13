package game

import (
	"math/rand"
)

// GenerateBoard creates a new board with randomly placed monsters
func GenerateBoard(width, height int, monsterDensity float64, seed int64) *Board {
	board := NewBoard(width, height)

	// Create seeded RNG for reproducible generation
	rng := rand.New(rand.NewSource(seed))

	// Calculate number of monsters
	totalCells := width * height
	monsterCount := int(float64(totalCells) * monsterDensity)

	// Place monsters randomly
	placedMonsters := 0
	attempts := 0
	maxAttempts := totalCells * 10 // Prevent infinite loops

	for placedMonsters < monsterCount && attempts < maxAttempts {
		attempts++
		idx := rng.Intn(totalCells)

		// Skip if already has a monster
		if board.Monsters[idx] {
			continue
		}

		// Place monster
		board.Monsters[idx] = true
		board.Cells[idx].HasMonster = true
		board.RemainingSafe--
		placedMonsters++
	}

	// Calculate threat counts for all cells
	board.CalculateThreatCounts()

	return board
}

// GenerateBoardWithExclusion generates a board ensuring the first click is safe
func GenerateBoardWithExclusion(width, height int, monsterDensity float64, seed int64, excludeX, excludeY int) *Board {
	board := NewBoard(width, height)
	rng := rand.New(rand.NewSource(seed))

	totalCells := width * height
	monsterCount := int(float64(totalCells) * monsterDensity)
	excludeIdx := board.GetIndex(excludeX, excludeY)

	// Get neighbor indices to exclude (including the clicked cell)
	excludeSet := make(map[int]bool)
	excludeSet[excludeIdx] = true
	for _, neighborIdx := range board.GetNeighborIndices(excludeIdx) {
		excludeSet[neighborIdx] = true
	}

	// Place monsters, excluding the safe zone
	placedMonsters := 0
	attempts := 0
	maxAttempts := totalCells * 10

	for placedMonsters < monsterCount && attempts < maxAttempts {
		attempts++
		idx := rng.Intn(totalCells)

		// Skip if already has a monster or in exclusion zone
		if board.Monsters[idx] || excludeSet[idx] {
			continue
		}

		// Place monster
		board.Monsters[idx] = true
		board.Cells[idx].HasMonster = true
		board.RemainingSafe--
		placedMonsters++
	}

	// Calculate threat counts
	board.CalculateThreatCounts()

	return board
}
