package server

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/stonehub/hecksweeper/internal/ui"
)

// gameHandler creates a new bubbletea model for each SSH session.
// Each connection gets an independent game instance.
func gameHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	// Use current timestamp as seed for each new connection
	seed := time.Now().UnixNano()

	// Check if the client sent a PLAYER_NAME env var
	// (via: ssh -o SetEnv=PLAYER_NAME=myname ...)
	// This is stored for potential future leaderboard use
	_ = getEnv(sess, "PLAYER_NAME")

	model := ui.NewModel(seed, true)

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
