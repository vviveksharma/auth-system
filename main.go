package main

import (
	"log"
	"os"
	"strings"

	"github.com/vviveksharma/auth/config"
)

var ServerMode string // Set via ldflags during build

func main() {
	// Check environment variable first, then build-time flag
	mode := strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' {
			return -1 // strip newlines to prevent log injection (G706)
		}
		return r
	}, os.Getenv("SERVER_MODE"))
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
	case "Project":
		config.InitProject()
	case "Org":
		config.InitOrg()
	case "BOTH":
		config.Init()
	default:
		log.Fatalf("❌ Invalid SERVER_MODE: %s (expected: API, UI, Project, Org, or BOTH)", mode)
	}
}
