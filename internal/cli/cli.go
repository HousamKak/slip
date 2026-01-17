package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"slip/internal/app"
	"slip/internal/platform"
	"slip/internal/plugins"
	"slip/internal/state"
	"slip/internal/store"
	"slip/internal/term"
	"slip/internal/tmux"
)

const (
	defaultWidth  = 40
	defaultHeight = 20
	defaultFPS    = 30
)

// Run is the main entry point for the CLI
func Run(args []string) int {
	if len(args) == 0 {
		return runCmd(nil)
	}

	switch args[0] {
	case "run":
		return runCmd(args[1:])
	case "summon":
		return summonCmd(args[1:])
	case "hide":
		return hideCmd(args[1:])
	case "dock":
		return dockCmd(args[1:])
	case "play":
		return playCmd(args[1:])
	case "animate":
		return animateCmd(args[1:])
	case "games":
		return gamesCmd(args[1:])
	case "animations":
		return animationsCmd(args[1:])
	case "store":
		return storeCmd(args[1:])
	case "scores":
		return scoresCmd(args[1:])
	case "config":
		return configCmd(args[1:])
	case "status":
		return statusCmd(args[1:])
	case "version", "-v", "--version":
		fmt.Println("slip version 1.0.0")
		return 0
	case "help", "-h", "--help":
		printUsage()
		return 0
	default:
		// Check if it's a game name
		a := app.New(defaultWidth, defaultHeight, defaultFPS)
		for _, g := range a.ListGames() {
			if g.ID == args[0] {
				return playCmd(args)
			}
		}
		for _, anim := range a.ListAnimations() {
			if anim.ID == args[0] {
				return animateCmd(args)
			}
		}

		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", args[0])
		printUsage()
		return 1
	}
}

func runCmd(args []string) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	width := fs.Int("width", envInt("SLIP_WIDTH", defaultWidth), "Screen width")
	height := fs.Int("height", envInt("SLIP_HEIGHT", defaultHeight), "Screen height")
	size := fs.String("size", envString("SLIP_SIZE", ""), "Screen size (WxH)")
	fps := fs.Int("fps", envInt("SLIP_FPS", defaultFPS), "Frames per second")
	tuiMode := fs.String("tui", "", "TUI mode: custom or bubbletea (overrides config)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	// Validate and set TUI mode if provided
	if *tuiMode != "" {
		if *tuiMode != "custom" && *tuiMode != "bubbletea" {
			fmt.Fprintln(os.Stderr, "--tui must be 'custom' or 'bubbletea'")
			return 2
		}
		os.Setenv("SLIP_TUI_MODE", *tuiMode)
	}

	if *size != "" {
		w, h, err := parseSize(*size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid size: %s\n", err)
			return 2
		}
		*width = w
		*height = h
	}

	// Auto-detect terminal size if not specified
	if *width == defaultWidth && *height == defaultHeight {
		if w, h, err := term.Size(); err == nil {
			*width = w
			*height = h - 1 // Leave room for prompt
		}
	}

	a := app.New(*width, *height, *fps)
	if err := a.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}

func summonCmd(args []string) int {
	fs := flag.NewFlagSet("summon", flag.ContinueOnError)
	dock := fs.String("dock", envString("SLIP_DOCK", "right"), "Dock position (right, left, bottom)")
	size := fs.String("size", envString("SLIP_SIZE", "40x20"), "Pane size (WxH)")
	fps := fs.Int("fps", envInt("SLIP_FPS", defaultFPS), "Frames per second")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	width, height, err := parseSize(*size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid size: %s\n", err)
		return 2
	}

	// Build command to run in the pane
	exe, err := os.Executable()
	if err != nil {
		exe = "slip"
	}

	// Check if in tmux
	if tmux.InTmux() {
		if !tmux.HasTmux() {
			fmt.Fprintln(os.Stderr, "tmux not found in PATH")
			return 1
		}

		// Kill existing slip pane if any
		_ = tmux.KillSlipPane()

		// Calculate pane size based on dock position
		paneSize := width + 2
		if strings.ToLower(*dock) == "bottom" {
			paneSize = height + 2
		}

		// For tmux, don't pass explicit dimensions - let slip auto-detect the pane size
		// This prevents aspect ratio issues when the actual pane differs from requested
		tmuxCmdArgs := []string{exe, "run", "--fps", strconv.Itoa(*fps)}

		// Create new pane
		paneID, err := tmux.SplitSlip(*dock, paneSize, tmuxCmdArgs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create tmux pane: %s\n", err)
			return 1
		}

		fmt.Printf("Slip summoned in pane %s\n", paneID)
		return 0
	}

	// Not in tmux - try Windows Terminal
	if platform.IsWindows() && !platform.IsWSL() {
		if platform.HasWindowsTerminal() {
			if platform.InWindowsTerminal() {
				// Inside Windows Terminal - can do proper docking
				if err := platform.LaunchWindowsTerminal(cmdArgs, *dock); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to launch Windows Terminal pane: %s\n", err)
					return 1
				}
				fmt.Printf("Slip summoned in Windows Terminal pane (docked %s)\n", *dock)
				return 0
			} else {
				// WT available but not running inside it - will open new window
				if err := platform.LaunchWindowsTerminal(cmdArgs, *dock); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to launch Windows Terminal: %s\n", err)
					return 1
				}
				fmt.Println("Slip launched in new Windows Terminal window")
				fmt.Println("Tip: Run Slip from inside Windows Terminal for proper pane docking")
				return 0
			}
		}
		// Windows but no WT
		printWindowsHelp()
		return 1
	}

	// Not Windows, not in tmux - check if tmux is available
	if tmux.HasTmux() {
		// tmux is installed but user is not inside a tmux session
		fmt.Fprintln(os.Stderr, "Not inside a tmux session")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "To use slip summon, first start tmux:")
		fmt.Fprintln(os.Stderr, "  1. Run: tmux")
		fmt.Fprintln(os.Stderr, "  2. Then run: slip summon")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Or just run 'slip run' in current terminal")
		return 1
	}

	// tmux not installed
	fmt.Fprintln(os.Stderr, "Docking not available")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Docking requires:")
	fmt.Fprintln(os.Stderr, "  • tmux on Linux/macOS/WSL")
	fmt.Fprintln(os.Stderr, "  • Windows Terminal (wt.exe) on Windows")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Alternatives:")
	fmt.Fprintln(os.Stderr, "  • Install tmux: apt install tmux / brew install tmux")
	fmt.Fprintln(os.Stderr, "  • Run 'slip run' in current terminal")
	return 1
}

func printWindowsHelp() {
	fmt.Fprintln(os.Stderr, "╭──────────────────────────────────────────────────╮")
	fmt.Fprintln(os.Stderr, "│  Docking not available                          │")
	fmt.Fprintln(os.Stderr, "├──────────────────────────────────────────────────┤")
	fmt.Fprintln(os.Stderr, "│  Docking requires:                              │")
	fmt.Fprintln(os.Stderr, "│  • Windows Terminal (wt.exe) on Windows         │")
	fmt.Fprintln(os.Stderr, "│  • tmux on Linux/macOS                          │")
	fmt.Fprintln(os.Stderr, "│                                                 │")
	fmt.Fprintln(os.Stderr, "│  Alternatives:                                  │")
	fmt.Fprintln(os.Stderr, "│  • Install Windows Terminal from Microsoft Store│")
	fmt.Fprintln(os.Stderr, "│  • Use Windows Terminal's built-in pane split:  │")
	fmt.Fprintln(os.Stderr, "│    Alt+Shift+Plus (right) or Alt+Shift+- (down) │")
	fmt.Fprintln(os.Stderr, "│    Then run 'slip run' in the new pane          │")
	fmt.Fprintln(os.Stderr, "│  • Install tmux via WSL for full dock support   │")
	fmt.Fprintln(os.Stderr, "│  • Run 'slip run' in current terminal           │")
	fmt.Fprintln(os.Stderr, "╰──────────────────────────────────────────────────╯")
}

func hideCmd(args []string) int {
	if tmux.InTmux() && tmux.HasTmux() {
		if err := tmux.KillSlipPane(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to hide slip: %s\n", err)
			return 1
		}
		fmt.Println("Slip hidden")
		return 0
	}

	fmt.Fprintln(os.Stderr, "Not in tmux or no slip pane found")
	return 1
}

func dockCmd(args []string) int {
	fs := flag.NewFlagSet("dock", flag.ContinueOnError)
	dock := fs.String("dock", "right", "Dock position (right, left, bottom)")
	size := fs.String("size", envString("SLIP_SIZE", "40x20"), "Pane size (WxH)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	// Re-summon at new position
	return summonCmd([]string{"--dock", *dock, "--size", *size})
}

func playCmd(args []string) int {
	if len(args) == 0 {
		// List available games
		return gamesCmd([]string{"list"})
	}

	gameID := args[0]

	fs := flag.NewFlagSet("play", flag.ContinueOnError)
	width := fs.Int("width", envInt("SLIP_WIDTH", defaultWidth), "")
	height := fs.Int("height", envInt("SLIP_HEIGHT", defaultHeight), "")
	size := fs.String("size", "", "")
	fps := fs.Int("fps", envInt("SLIP_FPS", defaultFPS), "")

	_ = fs.Parse(args[1:])

	if *size != "" {
		w, h, err := parseSize(*size)
		if err == nil {
			*width = w
			*height = h
		}
	}

	// Auto-detect terminal size
	if *width == defaultWidth && *height == defaultHeight {
		if w, h, err := term.Size(); err == nil {
			*width = w
			*height = h - 1
		}
	}

	a := app.New(*width, *height, *fps)
	if err := a.RunGame(gameID); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}

func animateCmd(args []string) int {
	if len(args) == 0 {
		return animationsCmd([]string{"list"})
	}

	animID := args[0]

	fs := flag.NewFlagSet("animate", flag.ContinueOnError)
	width := fs.Int("width", envInt("SLIP_WIDTH", defaultWidth), "")
	height := fs.Int("height", envInt("SLIP_HEIGHT", defaultHeight), "")
	size := fs.String("size", "", "")
	fps := fs.Int("fps", envInt("SLIP_FPS", defaultFPS), "")

	_ = fs.Parse(args[1:])

	if *size != "" {
		w, h, err := parseSize(*size)
		if err == nil {
			*width = w
			*height = h
		}
	}

	if *width == defaultWidth && *height == defaultHeight {
		if w, h, err := term.Size(); err == nil {
			*width = w
			*height = h - 1
		}
	}

	a := app.New(*width, *height, *fps)
	if err := a.RunAnimation(animID); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}

func gamesCmd(args []string) int {
	a := app.New(defaultWidth, defaultHeight, defaultFPS)
	games := a.ListGames()

	fmt.Println("Available games:")
	fmt.Println()
	for _, g := range games {
		fmt.Printf("  %-12s %s\n", g.ID, g.Description)
	}
	fmt.Println()
	fmt.Println("Usage: slip play <game>")
	return 0
}

func animationsCmd(args []string) int {
	a := app.New(defaultWidth, defaultHeight, defaultFPS)
	anims := a.ListAnimations()

	fmt.Println("Available animations:")
	fmt.Println()
	for _, anim := range anims {
		fmt.Printf("  %-12s %s\n", anim.ID, anim.Description)
	}
	fmt.Println()
	fmt.Println("Usage: slip animate <animation>")
	return 0
}

func scoresCmd(args []string) int {
	sb, err := state.LoadScores()
	if err != nil {
		fmt.Println("No scores recorded yet.")
		return 0
	}

	gameID := ""
	if len(args) > 0 {
		gameID = args[0]
	}

	if gameID != "" {
		scores := sb.GetTopScores(gameID, 10)
		if len(scores) == 0 {
			fmt.Printf("No scores for %s yet.\n", gameID)
			return 0
		}
		fmt.Printf("High scores for %s:\n\n", gameID)
		for i, s := range scores {
			fmt.Printf("  %2d. %8d  %s\n", i+1, s.Value, s.Player)
		}
	} else {
		fmt.Println("High scores:")
		fmt.Println()
		a := app.New(defaultWidth, defaultHeight, defaultFPS)
		for _, g := range a.ListGames() {
			high := sb.GetHighScore(g.ID)
			if high > 0 {
				fmt.Printf("  %-12s %d\n", g.ID, high)
			}
		}
	}
	return 0
}

func configCmd(args []string) int {
	if len(args) == 0 {
		// Show all config
		cfg, err := state.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
			return 1
		}
		fmt.Println("Current configuration:")
		fmt.Printf("  default_dock:     %s\n", cfg.DefaultDock)
		fmt.Printf("  default_size:     %s\n", cfg.DefaultSize)
		fmt.Printf("  default_fps:      %d\n", cfg.DefaultFPS)
		fmt.Printf("  theme:            %s\n", cfg.Theme)
		fmt.Printf("  tui_mode:         %s\n", cfg.TUIMode)
		fmt.Printf("  show_fps_counter: %v\n", cfg.ShowFPSCounter)
		fmt.Printf("  player_name:      %s\n", cfg.PlayerName)
		return 0
	}

	switch args[0] {
	case "get":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: slip config get <key>")
			return 2
		}
		val, err := state.GetConfig(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return 1
		}
		fmt.Println(val)
	case "set":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: slip config set <key> <value>")
			return 2
		}
		if err := state.SetConfig(args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return 1
		}
		fmt.Printf("Set %s = %s\n", args[1], args[2])
	case "reset":
		if err := state.SaveConfig(state.Config{}); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return 1
		}
		fmt.Println("Configuration reset to defaults")
	default:
		fmt.Fprintln(os.Stderr, "Usage: slip config [get|set|reset] ...")
		return 2
	}
	return 0
}

func statusCmd(args []string) int {
	// Check tmux pane
	if tmux.InTmux() && tmux.HasTmux() {
		paneID, found, _ := tmux.FindSlipPane()
		if found {
			fmt.Printf("Slip running in tmux pane %s\n", paneID)
			return 0
		}
	}

	fmt.Println("Slip not running")
	return 1
}

func storeCmd(args []string) int {
	if len(args) == 0 {
		// Show help for store command
		fmt.Println("Usage: slip store <command>")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  list              List all available games/animations")
		fmt.Println("  search <query>    Search for games/animations")
		fmt.Println("  info <id>         Show detailed info about a game/animation")
		fmt.Println("  install <id>      Install a game/animation")
		fmt.Println("  uninstall <id>    Uninstall a game/animation")
		fmt.Println("  installed         List installed plugins")
		fmt.Println("  update            Update the registry cache")
		return 0
	}

	// Import store package
	storeClient := store.NewClient()
	pluginLoader := plugins.NewLoader()

	switch args[0] {
	case "list":
		return storeListCmd(storeClient)
	case "search":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: slip store search <query>")
			return 2
		}
		return storeSearchCmd(storeClient, args[1])
	case "info":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: slip store info <id>")
			return 2
		}
		return storeInfoCmd(storeClient, args[1])
	case "install":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: slip store install <id>")
			return 2
		}
		return storeInstallCmd(storeClient, pluginLoader, args[1])
	case "uninstall":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: slip store uninstall <id>")
			return 2
		}
		return storeUninstallCmd(pluginLoader, args[1])
	case "installed":
		return storeInstalledCmd(pluginLoader)
	case "update":
		return storeUpdateCmd(storeClient)
	default:
		fmt.Fprintf(os.Stderr, "Unknown store command: %s\n", args[0])
		return 1
	}
}

func storeListCmd(client *store.Client) int {
	games, err := client.ListGames()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch store: %s\n", err)
		return 1
	}

	if len(games) == 0 {
		fmt.Println("No games available in the store.")
		return 0
	}

	fmt.Println("Available in Store:")
	fmt.Println()
	for _, game := range games {
		featured := ""
		if game.Featured {
			featured = " [FEATURED]"
		}
		fmt.Printf("  %-15s %s%s\n", game.ID, game.Name, featured)
		fmt.Printf("  %-15s %s\n", "", game.Description)
		fmt.Printf("  %-15s Author: %s | Version: %s\n", "", game.Author, game.Version)
		fmt.Println()
	}
	fmt.Printf("Total: %d games/animations\n", len(games))
	fmt.Println()
	fmt.Println("Use 'slip store install <id>' to install")
	return 0
}

func storeSearchCmd(client *store.Client, query string) int {
	games, err := client.SearchGames(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search failed: %s\n", err)
		return 1
	}

	if len(games) == 0 {
		fmt.Printf("No results found for '%s'\n", query)
		return 0
	}

	fmt.Printf("Search results for '%s':\n\n", query)
	for _, game := range games {
		fmt.Printf("  %-15s %s\n", game.ID, game.Name)
		fmt.Printf("  %-15s %s\n", "", game.Description)
		fmt.Println()
	}
	return 0
}

func storeInfoCmd(client *store.Client, id string) int {
	game, err := client.GetGame(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get game info: %s\n", err)
		return 1
	}

	fmt.Printf("ID:          %s\n", game.ID)
	fmt.Printf("Name:        %s\n", game.Name)
	fmt.Printf("Description: %s\n", game.Description)
	fmt.Printf("Author:      %s\n", game.Author)
	fmt.Printf("Version:     %s\n", game.Version)
	fmt.Printf("Type:        %s\n", game.Type)
	fmt.Printf("Category:    %s\n", game.Category)
	if game.Featured {
		fmt.Println("Featured:    Yes")
	}
	fmt.Println()

	// Check if available for current platform
	url, err := client.GetDownloadURL(game)
	if err != nil {
		fmt.Printf("Status:      Not available for your platform\n")
	} else {
		fmt.Printf("Status:      Available\n")
		fmt.Printf("Download:    %s\n", url)
	}
	return 0
}

func storeInstallCmd(client *store.Client, loader *plugins.Loader, id string) int {
	// Check if already installed
	if _, ok := loader.Get(id); ok {
		fmt.Printf("'%s' is already installed\n", id)
		fmt.Println("Use 'slip store uninstall' to remove it first")
		return 0
	}

	// Get game info
	game, err := client.GetGame(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Game not found: %s\n", err)
		return 1
	}

	fmt.Printf("Installing '%s' (%s)...\n", game.Name, game.Version)

	// Download and install with progress
	lastPct := -1
	err = client.DownloadWithProgress(game, "", func(pct int) {
		if pct != lastPct && pct%10 == 0 {
			fmt.Printf("\rDownloading: %d%%", pct)
			lastPct = pct
		}
	})
	fmt.Println() // New line after progress

	if err != nil {
		fmt.Fprintf(os.Stderr, "Installation failed: %s\n", err)
		return 1
	}

	fmt.Printf("Successfully installed '%s'\n", game.Name)
	fmt.Println("Restart Slip to see the new game in your menu")
	return 0
}

func storeUninstallCmd(loader *plugins.Loader, id string) int {
	// Check if installed
	if _, ok := loader.Get(id); !ok {
		fmt.Printf("'%s' is not installed\n", id)
		return 0
	}

	// Uninstall
	err := loader.Uninstall(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Uninstall failed: %s\n", err)
		return 1
	}

	fmt.Printf("Successfully uninstalled '%s'\n", id)
	return 0
}

func storeInstalledCmd(loader *plugins.Loader) int {
	if err := loader.Scan(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to scan plugins: %s\n", err)
		return 1
	}

	plugins := loader.List()
	if len(plugins) == 0 {
		fmt.Println("No plugins installed.")
		fmt.Println()
		fmt.Println("Use 'slip store list' to browse available content")
		return 0
	}

	fmt.Println("Installed plugins:")
	fmt.Println()
	for _, plugin := range plugins {
		fmt.Printf("  %-15s %s (%s)\n", plugin.ID, plugin.Name, plugin.Version)
		fmt.Printf("  %-15s %s\n", "", plugin.Description)
		fmt.Println()
	}
	fmt.Printf("Total: %d plugins\n", len(plugins))
	return 0
}

func storeUpdateCmd(client *store.Client) int {
	fmt.Println("Updating registry...")
	_, err := client.FetchRegistry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %s\n", err)
		return 1
	}
	fmt.Println("Registry updated successfully")
	return 0
}

func printUsage() {
	usage := `slip - Terminal Entertainment Platform

Usage:
  slip                    Start Slip interactive shell
  slip run [options]      Run Slip in current terminal
  slip summon [options]   Summon Slip in a tmux pane
  slip hide               Hide the Slip tmux pane
  slip dock [options]     Move Slip to a different dock position

  slip play <game>        Play a specific game
  slip animate <anim>     Run a specific animation
  slip games              List available games
  slip animations         List available animations
  slip scores [game]      View high scores

  slip store              Manage game store
  slip store list         List available content in store
  slip store search <q>   Search the store
  slip store install <id> Install a game/animation
  slip store installed    List installed plugins

  slip config             Show configuration
  slip config get <key>   Get a config value
  slip config set <k> <v> Set a config value
  slip status             Check if Slip is running

  slip version            Show version
  slip help               Show this help

Options:
  --width <n>             Screen width (default: auto)
  --height <n>            Screen height (default: auto)
  --size <WxH>            Screen size (e.g., 40x20)
  --fps <n>               Frames per second (default: 30)
  --tui <mode>            TUI mode: custom or bubbletea (overrides config)
  --dock <pos>            Dock position: right, left, bottom (default: right)

Examples:
  slip                    Start interactive shell
  slip summon             Summon in tmux (right side)
  slip summon --dock bottom --size 80x10
  slip play snake         Play Snake directly
  slip animate starfield  Watch starfield animation

Environment:
  SLIP_WIDTH, SLIP_HEIGHT, SLIP_SIZE, SLIP_FPS, SLIP_DOCK
`
	fmt.Print(usage)
}

// Helper functions

func parseSize(s string) (int, int, error) {
	parts := strings.Split(strings.ToLower(s), "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("size must be WxH format")
	}
	w, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	h, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}
	if w < 10 || h < 5 {
		return 0, 0, fmt.Errorf("minimum size is 10x5")
	}
	return w, h, nil
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
