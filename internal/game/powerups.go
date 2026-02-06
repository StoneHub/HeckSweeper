package game

import (
	"math/rand"
)

// Rarity represents power-up rarity tiers
type Rarity int

const (
	Common    Rarity = iota
	Uncommon
	Rare
	Legendary
)

// PowerUpID identifies a specific power-up type
type PowerUpID string

const (
	PowerUpScout      PowerUpID = "scout"
	PowerUpShield     PowerUpID = "shield"
	PowerUpDowsing    PowerUpID = "dowsing"
	PowerUpXRay       PowerUpID = "xray"
	PowerUpLuckyCharm PowerUpID = "lucky_charm"
	PowerUpSwiftBoots PowerUpID = "swift_boots"
	PowerUpFlagMaster PowerUpID = "flag_master"
	PowerUpSecondWind PowerUpID = "second_wind"
	PowerUpTimeWarp   PowerUpID = "time_warp"
)

// PowerUpDef defines a power-up's properties
type PowerUpDef struct {
	ID          PowerUpID
	Name        string
	Description string
	Rarity      Rarity
}

// ActivePowerUp tracks a collected power-up in a run
type ActivePowerUp struct {
	ID       PowerUpID
	Name     string
	Uses     int // -1 = unlimited, 0 = consumed
	Consumed bool
}

// PowerUp catalog — all available power-ups
var powerUpCatalog = []PowerUpDef{
	{PowerUpScout, "Scout", "Reveal a random 3x3 safe area at floor start", Common},
	{PowerUpShield, "Shield", "Survive one monster hit (consumed on use)", Rare},
	{PowerUpDowsing, "Dowsing Rod", "First reveal on each floor is guaranteed safe", Common},
	{PowerUpXRay, "X-Ray", "Reveal all edge cells at floor start", Uncommon},
	{PowerUpLuckyCharm, "Lucky Charm", "Reduce monster density by 2% on future floors", Uncommon},
	{PowerUpSwiftBoots, "Swift Boots", "+50% speed bonus multiplier", Common},
	{PowerUpFlagMaster, "Flag Master", "Auto-flag obvious monster cells on reveal", Common},
	{PowerUpSecondWind, "Second Wind", "Revive once per run on death (consumed)", Legendary},
	{PowerUpTimeWarp, "Time Warp", "Rewind your last move once per floor (auto-prevents death)", Rare},
}

// rarityWeights controls how likely each rarity is to be offered
var rarityWeights = map[Rarity]int{
	Common:    40,
	Uncommon:  30,
	Rare:      20,
	Legendary: 10,
}

// RarityName returns a display string for a rarity tier
func RarityName(r Rarity) string {
	switch r {
	case Common:
		return "Common"
	case Uncommon:
		return "Uncommon"
	case Rare:
		return "Rare"
	case Legendary:
		return "Legendary"
	default:
		return "Unknown"
	}
}

// GetPowerUpCatalog returns all available power-up definitions
func GetPowerUpCatalog() []PowerUpDef {
	return powerUpCatalog
}

// PickPowerUpChoices selects n random power-ups weighted by rarity.
// Uses the provided RNG for deterministic selection.
func PickPowerUpChoices(rng *rand.Rand, n int, alreadyHave []PowerUpID) []PowerUpDef {
	// Filter out power-ups the player already has (for non-stackable ones)
	haveSet := make(map[PowerUpID]bool)
	for _, id := range alreadyHave {
		haveSet[id] = true
	}

	available := []PowerUpDef{}
	for _, p := range powerUpCatalog {
		// Shield and SecondWind can be re-acquired if consumed
		if haveSet[p.ID] && p.ID != PowerUpShield && p.ID != PowerUpSecondWind {
			continue
		}
		available = append(available, p)
	}

	if len(available) == 0 {
		return nil
	}
	if len(available) <= n {
		return available
	}

	// Weighted random selection without replacement
	chosen := []PowerUpDef{}
	remaining := make([]PowerUpDef, len(available))
	copy(remaining, available)

	for i := 0; i < n && len(remaining) > 0; i++ {
		// Build weighted index
		totalWeight := 0
		for _, p := range remaining {
			totalWeight += rarityWeights[p.Rarity]
		}

		roll := rng.Intn(totalWeight)
		cumulative := 0
		selectedIdx := 0

		for idx, p := range remaining {
			cumulative += rarityWeights[p.Rarity]
			if roll < cumulative {
				selectedIdx = idx
				break
			}
		}

		chosen = append(chosen, remaining[selectedIdx])
		// Remove selected from remaining
		remaining = append(remaining[:selectedIdx], remaining[selectedIdx+1:]...)
	}

	return chosen
}

// ApplyFloorStartPowerUps applies power-ups that activate at the start of each floor
func ApplyFloorStartPowerUps(run *Run) {
	g := run.CurrentGame
	if g == nil {
		return
	}

	for i := range run.PowerUps {
		pu := &run.PowerUps[i]
		if pu.Consumed {
			continue
		}

		switch pu.ID {
		case PowerUpScout:
			applyScout(g, run.Seed+int64(run.FloorNum)*31)

		case PowerUpXRay:
			applyXRay(g)

		case PowerUpLuckyCharm:
			// Lucky Charm effect is applied during board generation (density reduction)
			// Handled in Run.StartNextFloor

		case PowerUpDowsing:
			// Dowsing Rod makes the cursor start on a safe cell
			applyDowsing(g)
		}
	}
}

// applyScout reveals a random 3x3 safe area
func applyScout(g *Game, seed int64) {
	rng := rand.New(rand.NewSource(seed))
	b := g.Board

	// Try to find a 3x3 area with no monsters
	for attempts := 0; attempts < 100; attempts++ {
		cx := rng.Intn(b.Width-2) + 1
		cy := rng.Intn(b.Height-2) + 1

		safe := true
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				cell := b.GetCell(cx+dx, cy+dy)
				if cell != nil && cell.HasMonster {
					safe = false
					break
				}
			}
			if !safe {
				break
			}
		}

		if safe {
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					cell := b.GetCell(cx+dx, cy+dy)
					if cell != nil && !cell.Revealed {
						cell.Revealed = true
						b.RemainingSafe--
					}
				}
			}
			return
		}
	}
}

// applyXRay reveals all edge cells that don't have monsters
func applyXRay(g *Game) {
	b := g.Board
	for x := 0; x < b.Width; x++ {
		for _, y := range []int{0, b.Height - 1} {
			cell := b.GetCell(x, y)
			if cell != nil && !cell.HasMonster && !cell.Revealed {
				cell.Revealed = true
				b.RemainingSafe--
			}
		}
	}
	for y := 1; y < b.Height-1; y++ {
		for _, x := range []int{0, b.Width - 1} {
			cell := b.GetCell(x, y)
			if cell != nil && !cell.HasMonster && !cell.Revealed {
				cell.Revealed = true
				b.RemainingSafe--
			}
		}
	}
}

// applyDowsing moves cursor to a guaranteed safe cell
func applyDowsing(g *Game) {
	b := g.Board
	// Find a safe cell near the center
	cx, cy := b.Width/2, b.Height/2
	for r := 0; r < max(b.Width, b.Height); r++ {
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				x, y := cx+dx, cy+dy
				cell := b.GetCell(x, y)
				if cell != nil && !cell.HasMonster && !cell.Revealed {
					g.CursorPos = Position{X: x, Y: y}
					return
				}
			}
		}
	}
}

// HandleDeathPowerUps checks if any power-up prevents death.
// Returns true if death was prevented.
// Priority: Time Warp (renewable) > Shield/SecondWind (consumable)
func HandleDeathPowerUps(run *Run) bool {
	// First: check Time Warp (renewable, recharges each floor)
	for i := range run.PowerUps {
		pu := &run.PowerUps[i]
		if pu.Consumed || pu.ID != PowerUpTimeWarp {
			continue
		}
		if run.LastSnapshot != nil && run.UndoUsedFloor != run.FloorNum {
			run.RestoreSnapshot()
			run.UndoUsedFloor = run.FloorNum
			return true
		}
	}

	// Second: check consumable death prevention
	for i := range run.PowerUps {
		pu := &run.PowerUps[i]
		if pu.Consumed {
			continue
		}

		switch pu.ID {
		case PowerUpShield:
			pu.Consumed = true
			g := run.CurrentGame
			cell := g.Board.GetCell(g.CursorPos.X, g.CursorPos.Y)
			if cell != nil {
				cell.Revealed = false
			}
			g.State = "playing"
			return true

		case PowerUpSecondWind:
			pu.Consumed = true
			g := run.CurrentGame
			for j := range g.Board.Cells {
				if g.Board.Cells[j].HasMonster {
					g.Board.Cells[j].Revealed = false
				}
			}
			g.State = "playing"
			return true
		}
	}
	return false
}

// HasPowerUp checks if the run has a specific (non-consumed) power-up
func (r *Run) HasPowerUp(id PowerUpID) bool {
	for _, pu := range r.PowerUps {
		if pu.ID == id && !pu.Consumed {
			return true
		}
	}
	return false
}

// PowerUpNames returns display names of all active (non-consumed) power-ups
func (r *Run) PowerUpNames() []string {
	names := []string{}
	for _, pu := range r.PowerUps {
		if !pu.Consumed {
			names = append(names, pu.Name)
		}
	}
	return names
}

// AddPowerUp adds a power-up to the run
func (r *Run) AddPowerUp(def PowerUpDef) {
	uses := -1 // Unlimited by default
	if def.ID == PowerUpShield || def.ID == PowerUpSecondWind {
		uses = 1 // Single-use
	}

	r.PowerUps = append(r.PowerUps, ActivePowerUp{
		ID:   def.ID,
		Name: def.Name,
		Uses: uses,
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
