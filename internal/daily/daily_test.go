package daily

import (
	"testing"
	"time"
)

func TestSeedDeterministic(t *testing.T) {
	date := time.Date(2026, 2, 6, 0, 0, 0, 0, time.UTC)

	seed1 := SeedForDate(date)
	seed2 := SeedForDate(date)

	if seed1 != seed2 {
		t.Errorf("Same date should produce same seed: %d vs %d", seed1, seed2)
	}
}

func TestSeedDifferentDates(t *testing.T) {
	date1 := time.Date(2026, 2, 6, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2026, 2, 7, 0, 0, 0, 0, time.UTC)

	seed1 := SeedForDate(date1)
	seed2 := SeedForDate(date2)

	if seed1 == seed2 {
		t.Error("Different dates should produce different seeds")
	}
}

func TestSeedIgnoresTime(t *testing.T) {
	morning := time.Date(2026, 2, 6, 8, 30, 0, 0, time.UTC)
	evening := time.Date(2026, 2, 6, 22, 45, 0, 0, time.UTC)

	if SeedForDate(morning) != SeedForDate(evening) {
		t.Error("Same date at different times should produce same seed")
	}
}

func TestDateString(t *testing.T) {
	date := time.Date(2026, 2, 6, 0, 0, 0, 0, time.UTC)
	result := DateStringFor(date)

	if result != "2026-02-06" {
		t.Errorf("Expected 2026-02-06, got %s", result)
	}
}
