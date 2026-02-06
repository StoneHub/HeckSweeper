package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"github.com/stonehub/hecksweeper/internal/storage"
)

// Config holds the SSH server configuration
type Config struct {
	Host        string
	Port        int
	KeyPath     string
	DataPath    string
	IdleTimeout time.Duration
	MaxTimeout  time.Duration
}

// DefaultConfig returns sensible defaults for the SSH server
func DefaultConfig() Config {
	return Config{
		Host:        "0.0.0.0",
		Port:        2222,
		KeyPath:     ".ssh/hecksweeper_ed25519",
		DataPath:    "hecksweeper_data.json",
		IdleTimeout: 10 * time.Minute,
		MaxTimeout:  60 * time.Minute,
	}
}

// Start creates and runs the Wish SSH server with graceful shutdown
func Start(cfg Config) error {
	// Initialize persistent storage
	store, err := storage.NewStore(cfg.DataPath)
	if err != nil {
		return fmt.Errorf("could not initialize storage: %w", err)
	}
	SetStore(store)

	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(cfg.KeyPath),
		wish.WithIdleTimeout(cfg.IdleTimeout),
		wish.WithMaxTimeout(cfg.MaxTimeout),
		wish.WithMiddleware(
			bm.Middleware(gameHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
		// Accept all connections (anonymous play)
		wish.WithPublicKeyAuth(func(_ ssh.Context, _ ssh.PublicKey) bool {
			return true
		}),
		wish.WithPasswordAuth(func(_ ssh.Context, _ string) bool {
			return true
		}),
	)
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	// Graceful shutdown on SIGINT/SIGTERM
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	log.Info("Starting HeckSweeper SSH server", "addr", addr)
	fmt.Printf("HeckSweeper SSH server listening on %s\n", addr)
	fmt.Printf("Connect with: ssh -p %d localhost\n", cfg.Port)
	fmt.Printf("Data file: %s\n", cfg.DataPath)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Error("Server error", "error", err)
		}
	}()

	<-done
	log.Info("Shutting down server...")
	fmt.Println("\nShutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.Shutdown(ctx)
}
