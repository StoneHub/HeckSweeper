package server

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	gossh "golang.org/x/crypto/ssh"

	"github.com/stonehub/hecksweeper/internal/storage"
	"github.com/stonehub/hecksweeper/internal/ui"
)

// gameStore is the shared persistent store for all sessions
var gameStore *storage.Store

// SetStore configures the shared store for all game sessions
func SetStore(s *storage.Store) {
	gameStore = s
}

// gameHandler creates a new bubbletea model for each SSH session.
// Each connection gets an independent game instance.
func gameHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	seed := time.Now().UnixNano()

	// Identify player by SSH key fingerprint (if provided)
	playerID := "anon"
	playerName := ""
	if pubKey := sess.PublicKey(); pubKey != nil {
		playerID = gossh.FingerprintSHA256(pubKey)
	}

	// Check for player name env var
	if name := getEnv(sess, "PLAYER_NAME"); name != "" {
		playerName = name
	} else if sess.User() != "" {
		playerName = sess.User()
	}

	// Record player in store
	if gameStore != nil {
		gameStore.RecordPlayer(playerID, playerName)
	}

	model := ui.NewModelWithStore(seed, true, gameStore, playerID, playerName)

	return model, []tea.ProgramOption{tea.WithAltScreen()}
}

// getEnv reads an environment variable from the SSH session
func getEnv(sess ssh.Session, key string) string {
	for _, env := range sess.Environ() {
		if len(env) > len(key)+1 && env[:len(key)+1] == key+"=" {
			return env[len(key)+1:]
		}
	}
	return ""
}
