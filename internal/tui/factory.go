package tui

import (
	"os"

	"slip/internal/state"
)

// NewTUI creates a TUI based on configuration
func NewTUI(cfg Config) (TUI, error) {
	// Check env var first (for CLI override)
	mode := os.Getenv("SLIP_TUI_MODE")
	if mode == "" {
		// Load user preference
		userCfg, _ := state.LoadConfig()
		mode = userCfg.TUIMode
	}
	if mode == "" {
		mode = "custom" // Default to custom for backwards compatibility
	}

	switch mode {
	case "bubbletea":
		return NewBubbleTeaTUI(cfg)
	case "custom":
		return NewCustomTUI(cfg)
	default:
		return NewCustomTUI(cfg)
	}
}
