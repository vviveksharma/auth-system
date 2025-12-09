package main

import (
	"log"
	"os"

	"github.com/vviveksharma/auth/config"
)

var ServerMode string // Set via ldflags during build

func main() {
	// Check environment variable first, then build-time flag
	mode := os.Getenv("SERVER_MODE")
	if mode == "" {
		mode = ServerMode // From ldflags
	}
	if mode == "" {
		mode = "BOTH" // Default: run both servers
	}

	log.Printf("🚀 Starting in %s mode", mode)

	switch mode {
	case "API":
		config.InitAPIOnly()
	case "UI":
		config.InitUIOnly()
	case "BOTH":
		config.Init()
	default:
		log.Fatalf("❌ Invalid SERVER_MODE: %s (expected: API, UI, or BOTH)", mode)
	}
}
