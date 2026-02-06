package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/stonehub/hecksweeper/internal/server"
)

func main() {
	cfg := server.DefaultConfig()

	flag.StringVar(&cfg.Host, "host", cfg.Host, "Listen address")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "SSH port")
	flag.StringVar(&cfg.KeyPath, "key", cfg.KeyPath, "Host key path (auto-generated if missing)")
	flag.Parse()

	if err := server.Start(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
