package daily

import (
	"hash/fnv"
	"time"
)

// Seed returns a deterministic seed for today's daily challenge.
// All players on the same UTC day get the same seed.
func Seed() int64 {
	return SeedForDate(time.Now().UTC())
}

// SeedForDate returns a deterministic seed for a specific date.
func SeedForDate(t time.Time) int64 {
	dateStr := t.Format("2006-01-02")
	h := fnv.New64a()
	h.Write([]byte("hecksweeper-daily-" + dateStr))
	return int64(h.Sum64())
}

// DateString returns today's date string in UTC (used as leaderboard key)
func DateString() string {
	return time.Now().UTC().Format("2006-01-02")
}

// DateStringFor returns the date string for a given time
func DateStringFor(t time.Time) string {
	return t.UTC().Format("2006-01-02")
}
