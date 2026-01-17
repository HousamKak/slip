package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"slip/internal/animations"
	"slip/internal/engine"
	"slip/internal/games"
	"slip/internal/state"
	"slip/internal/store"
	"slip/internal/term"
	"slip/internal/ui"
)

// CustomTUI is the custom terminal UI implementation
type CustomTUI struct {
	width       int
	height      int
	fps         int
	screen      Screen
	prevScreen  Screen
	running     bool

	// Menus
	homeMenu       *ui.Menu
	gamesMenu      *ui.Menu
	animationsMenu *ui.Menu

	// Registries
	gameRegistry *games.Registry
	animRegistry *animations.Registry

	// Store
	storeClient      *store.Client
	storeGames       []store.GameEntry
	storeLoading     bool
	storeError       string
	storeSelectedIdx int

	// Callbacks
	gameLauncher func(string) error
	animLauncher func(string) error

	// Theme
	theme ui.Theme

	// Command system
	commandMode    bool
	commandInput   string
	commandCursor  int
	commandRegistry *ui.CommandRegistry
	commandHistory *ui.CommandHistory
	commandMessage string
	commandError   error
}

// NewCustomTUI creates a new custom TUI
func NewCustomTUI(cfg Config) (*CustomTUI, error) {
	// Load theme
	theme := ui.GetTheme(cfg.Theme)
	if cfg.Theme == "" {
		userCfg, _ := state.LoadConfig()
		theme = ui.GetTheme(userCfg.Theme)
	}

	tui := &CustomTUI{
		width:           cfg.Width,
		height:          cfg.Height,
		fps:             cfg.FPS,
		screen:          ScreenHome,
		gameRegistry:    games.NewRegistry(),
		animRegistry:    animations.NewRegistry(),
		storeClient:     store.NewClient(),
		theme:           theme,
		commandRegistry: ui.NewCommandRegistry(),
		commandHistory:  ui.NewCommandHistory(100),
	}

	tui.initMenus()
	tui.initCommands()
	return tui, nil
}

func (t *CustomTUI) initMenus() {
	// Home menu
	t.homeMenu = ui.NewMenu("S L I P", []ui.MenuItem{
		{ID: "games", Label: "Games", Description: "Play classic arcade games"},
		{ID: "animations", Label: "Animations", Description: "Watch relaxing visualizations"},
		{ID: "store", Label: "Store", Description: "Browse and install games/animations"},
		{ID: "settings", Label: "Settings", Description: "Configure Slip"},
		{ID: "help", Label: "Help", Description: "How to use Slip"},
		{ID: "quit", Label: "Quit", Description: "Exit Slip"},
	})
	t.homeMenu.ShowNumbers = true

	// Build games menu from registry
	var gameItems []ui.MenuItem
	for _, info := range t.gameRegistry.List() {
		gameItems = append(gameItems, ui.MenuItem{
			ID:          info.ID,
			Label:       info.Name,
			Description: info.Description,
		})
	}
	gameItems = append(gameItems, ui.MenuItem{
		ID:    "back",
		Label: "← Back",
	})
	t.gamesMenu = ui.NewMenu("Games", gameItems)
	t.gamesMenu.ShowNumbers = true

	// Build animations menu from registry
	var animItems []ui.MenuItem
	for _, info := range t.animRegistry.List() {
		animItems = append(animItems, ui.MenuItem{
			ID:          info.ID,
			Label:       info.Name,
			Description: info.Description,
		})
	}
	animItems = append(animItems, ui.MenuItem{
		ID:    "back",
		Label: "← Back",
	})
	t.animationsMenu = ui.NewMenu("Animations", animItems)
	t.animationsMenu.ShowNumbers = true
}

func (t *CustomTUI) initCommands() {
	// Register navigation commands
	t.commandRegistry.Register(ui.NewNavigationCommand("home", "home", "Go to home screen"))
	t.commandRegistry.Register(ui.NewNavigationCommand("games", "games", "Go to games menu"))
	t.commandRegistry.Register(ui.NewNavigationCommand("animations", "animations", "Go to animations menu"))
	t.commandRegistry.Register(ui.NewNavigationCommand("store", "store", "Go to store"))
	t.commandRegistry.Register(ui.NewNavigationCommand("settings", "settings", "Go to settings"))
	t.commandRegistry.Register(ui.NewNavigationCommand("back", "back", "Go back"))

	// Register theme command
	themes := ui.ListThemes()
	t.commandRegistry.Register(ui.NewThemeCommand(themes))

	// Register game command
	var gameIDs []string
	for _, info := range t.gameRegistry.List() {
		gameIDs = append(gameIDs, info.ID)
	}
	t.commandRegistry.Register(ui.NewPlayCommand(gameIDs))

	// Register animation command
	var animIDs []string
	for _, info := range t.animRegistry.List() {
		animIDs = append(animIDs, info.ID)
	}
	t.commandRegistry.Register(ui.NewAnimateCommand(animIDs))

	// Register system commands
	t.commandRegistry.Register(ui.NewQuitCommand())
	t.commandRegistry.Register(ui.NewClearCommand())
	t.commandRegistry.Register(ui.NewHelpCommand(t.commandRegistry))
}

// Run starts the TUI
func (t *CustomTUI) Run() error {
	// Set up terminal
	term.EnableAlternateScreen(os.Stdout)
	term.HideCursor(os.Stdout)
	term.ClearScreenFull(os.Stdout)

	defer func() {
		term.ShowCursor(os.Stdout)
		term.DisableAlternateScreen(os.Stdout)
		term.ResetStyle(os.Stdout)
	}()

	// Create input reader
	input := engine.NewInputReader()
	if err := input.Start(); err != nil {
		return fmt.Errorf("failed to start input reader: %w", err)
	}
	defer input.Stop()

	// Create screen buffer
	screen := engine.NewScreen(t.width, t.height, os.Stdout)

	t.running = true

	// Resize check ticker
	resizeTicker := time.NewTicker(500 * time.Millisecond)
	defer resizeTicker.Stop()

	// Render ticker for idle updates (store loading, etc.)
	renderTicker := time.NewTicker(100 * time.Millisecond)
	defer renderTicker.Stop()

	// Initial render
	screen.Clear()
	t.render(screen)
	screen.Flush()

	for t.running {
		// Handle input and events - blocking select for better responsiveness
		select {
		case inp := <-input.Channel():
			t.handleInput(inp)
			// Render immediately after input for instant feedback
			screen.Clear()
			t.render(screen)
			screen.Flush()
		case <-resizeTicker.C:
			// Check for terminal resize
			w, h, err := term.Size()
			if err == nil && (w != t.width || h != t.height) {
				t.handleResize(w, h, screen)
				screen.Clear()
				t.render(screen)
				screen.Flush()
			}
		case <-renderTicker.C:
			// Periodic render for animations/loading states
			screen.Clear()
			t.render(screen)
			screen.Flush()
		}
	}

	return nil
}

// Stop gracefully stops the TUI
func (t *CustomTUI) Stop() {
	t.running = false
}

// SetScreen changes the current screen
func (t *CustomTUI) SetScreen(screen Screen) {
	t.prevScreen = t.screen
	t.screen = screen
}

// GetScreen returns the current screen
func (t *CustomTUI) GetScreen() Screen {
	return t.screen
}

// RegisterGameLauncher registers a game launcher callback
func (t *CustomTUI) RegisterGameLauncher(fn func(gameID string) error) {
	t.gameLauncher = fn
}

// RegisterAnimationLauncher registers an animation launcher callback
func (t *CustomTUI) RegisterAnimationLauncher(fn func(animID string) error) {
	t.animLauncher = fn
}

func (t *CustomTUI) handleInput(input engine.Input) {
	// Handle command mode
	if t.commandMode {
		t.handleCommandInput(input)
		return
	}

	// Global command mode trigger
	if input.Key == engine.KeyRune && input.Rune == ':' {
		t.commandMode = true
		t.commandInput = ""
		t.commandCursor = 0
		t.commandMessage = ""
		t.commandError = nil
		return
	}

	// Global quit
	if input.IsQuit() {
		t.running = false
		return
	}

	switch t.screen {
	case ScreenHome:
		t.handleHomeInput(input)
	case ScreenGames:
		t.handleGamesInput(input)
	case ScreenAnimations:
		t.handleAnimationsInput(input)
	case ScreenStore:
		t.handleStoreInput(input)
	case ScreenHelp:
		t.handleHelpInput(input)
	case ScreenSettings:
		t.handleSettingsInput(input)
	}
}

func (t *CustomTUI) handleHomeInput(input engine.Input) {
	if t.homeMenu.HandleInput(input) {
		selected := t.homeMenu.GetSelected()
		if selected != nil && input.Key == engine.KeyEnter {
			t.handleHomeSelection(selected.ID)
		}
	}
}

func (t *CustomTUI) handleHomeSelection(id string) {
	switch id {
	case "games":
		t.screen = ScreenGames
	case "animations":
		t.screen = ScreenAnimations
	case "store":
		t.screen = ScreenStore
		// Load store games if not loaded yet
		if len(t.storeGames) == 0 && !t.storeLoading && t.storeError == "" {
			go t.loadStore()
		}
	case "settings":
		t.screen = ScreenSettings
	case "help":
		t.screen = ScreenHelp
	case "quit":
		t.running = false
	}
}

func (t *CustomTUI) handleGamesInput(input engine.Input) {
	if input.Key == engine.KeyEscape {
		t.screen = ScreenHome
		return
	}

	if t.gamesMenu.HandleInput(input) {
		selected := t.gamesMenu.GetSelected()
		if selected != nil && input.Key == engine.KeyEnter {
			if selected.ID == "back" {
				t.screen = ScreenHome
			} else if t.gameLauncher != nil {
				_ = t.gameLauncher(selected.ID)
			}
		}
	}
}

func (t *CustomTUI) handleAnimationsInput(input engine.Input) {
	if input.Key == engine.KeyEscape {
		t.screen = ScreenHome
		return
	}

	if t.animationsMenu.HandleInput(input) {
		selected := t.animationsMenu.GetSelected()
		if selected != nil && input.Key == engine.KeyEnter {
			if selected.ID == "back" {
				t.screen = ScreenHome
			} else if t.animLauncher != nil {
				_ = t.animLauncher(selected.ID)
			}
		}
	}
}

func (t *CustomTUI) handleHelpInput(input engine.Input) {
	if input.Key == engine.KeyEscape || input.Key == engine.KeyEnter {
		t.screen = ScreenHome
	}
}

func (t *CustomTUI) handleStoreInput(input engine.Input) {
	if input.Key == engine.KeyEscape {
		t.screen = ScreenHome
		return
	}

	// Handle loading state
	if t.storeLoading {
		return
	}

	// Handle error state - Enter retries
	if t.storeError != "" {
		if input.Key == engine.KeyEnter {
			t.loadStore()
		}
		return
	}

	// Navigate store list
	switch input.Key {
	case engine.KeyUp:
		if t.storeSelectedIdx > 0 {
			t.storeSelectedIdx--
		}
	case engine.KeyDown:
		if t.storeSelectedIdx < len(t.storeGames)-1 {
			t.storeSelectedIdx++
		}
	case engine.KeyEnter:
		// Install selected game (placeholder - would need confirmation dialog)
		// For now, this is not implemented
		// In production, this would trigger installation
	}
}

func (t *CustomTUI) loadStore() {
	t.storeLoading = true
	t.storeError = ""

	// This should be done in background, but for simplicity we'll block
	games, err := t.storeClient.ListGames()
	if err != nil {
		t.storeError = err.Error()
		t.storeLoading = false
		return
	}

	t.storeGames = games
	t.storeLoading = false
	t.storeSelectedIdx = 0
}

func (t *CustomTUI) handleCommandInput(input engine.Input) {
	switch input.Key {
	case engine.KeyEscape:
		// Exit command mode
		t.commandMode = false
		t.commandInput = ""
		t.commandMessage = ""
		t.commandError = nil

	case engine.KeyEnter:
		// Execute command
		t.executeCommand(t.commandInput)
		t.commandHistory.Add(t.commandInput)
		t.commandMode = false
		t.commandInput = ""

	case engine.KeyTab:
		// Autocomplete
		suggestions := t.getCommandSuggestions()
		if len(suggestions) == 1 {
			// Single match - complete it
			cmdName, args := ui.ParseCommandLine(t.commandInput)
			if len(args) == 0 && !strings.Contains(t.commandInput, " ") {
				// Completing command name
				t.commandInput = suggestions[0] + " "
			} else if len(args) > 0 {
				// Completing argument
				t.commandInput = cmdName + " " + suggestions[0]
			}
		}

	case engine.KeyBackspace:
		// Delete character
		if len(t.commandInput) > 0 {
			t.commandInput = t.commandInput[:len(t.commandInput)-1]
		}

	case engine.KeyUp:
		// Previous command in history
		if prev := t.commandHistory.Previous(); prev != "" {
			t.commandInput = prev
		}

	case engine.KeyDown:
		// Next command in history
		if next := t.commandHistory.Next(); next != "" {
			t.commandInput = next
		} else {
			t.commandInput = ""
		}

	case engine.KeyRune:
		// Add character to input
		t.commandInput += string(input.Rune)
	}
}

func (t *CustomTUI) executeCommand(cmdLine string) {
	// Clear previous messages
	t.commandMessage = ""
	t.commandError = nil

	// Parse command
	cmdName, args := ui.ParseCommandLine(cmdLine)
	if cmdName == "" {
		return
	}

	// Get command
	cmd, ok := t.commandRegistry.Get(cmdName)
	if !ok {
		t.commandError = fmt.Errorf("unknown command: %s", cmdName)
		return
	}

	// Create command context
	ctx := &ui.CommandContext{
		TUI: t,
		Output: func(msg string, err error) {
			t.commandMessage = msg
			t.commandError = err
		},
		SetScreen: func(screen string) {
			switch screen {
			case "quit":
				t.running = false
			case "home":
				t.screen = ScreenHome
			case "games":
				t.screen = ScreenGames
			case "animations":
				t.screen = ScreenAnimations
			case "store":
				t.screen = ScreenStore
				if len(t.storeGames) == 0 && !t.storeLoading && t.storeError == "" {
					go t.loadStore()
				}
			case "settings":
				t.screen = ScreenSettings
			case "back":
				t.screen = t.prevScreen
			}
		},
		SetTheme: func(themeName string) error {
			theme := ui.GetTheme(themeName)
			t.theme = theme
			// Also save to config
			cfg, _ := state.LoadConfig()
			cfg.Theme = themeName
			_ = state.SaveConfig(cfg)
			return nil
		},
		LaunchGame: func(gameID string) error {
			if t.gameLauncher != nil {
				return t.gameLauncher(gameID)
			}
			return fmt.Errorf("game launcher not available")
		},
		LaunchAnim: func(animID string) error {
			if t.animLauncher != nil {
				return t.animLauncher(animID)
			}
			return fmt.Errorf("animation launcher not available")
		},
	}

	// Execute command
	err := cmd.Execute(ctx, args)
	if err != nil {
		t.commandError = err
	}
}

func (t *CustomTUI) handleSettingsInput(input engine.Input) {
	if input.Key == engine.KeyEscape || input.Key == engine.KeyEnter {
		t.screen = ScreenHome
	}
}

func (t *CustomTUI) render(screen *engine.Screen) {
	switch t.screen {
	case ScreenHome:
		t.renderHome(screen)
	case ScreenGames:
		t.renderGames(screen)
	case ScreenAnimations:
		t.renderAnimations(screen)
	case ScreenStore:
		t.renderStore(screen)
	case ScreenHelp:
		t.renderHelp(screen)
	case ScreenSettings:
		t.renderSettings(screen)
	}

	// Render command bar if in command mode
	if t.commandMode {
		t.renderCommandBar(screen)
	}

	// Render command message/error if any
	if t.commandMessage != "" || t.commandError != nil {
		t.renderCommandMessage(screen)
	}
}

func (t *CustomTUI) renderCommandBar(screen *engine.Screen) {
	// Get autocomplete suggestions
	suggestions := t.getCommandSuggestions()

	// Draw separator - move up if we have suggestions
	sepY := t.height - 3
	if len(suggestions) > 0 {
		sepY = t.height - 4
	}

	borderStyle := engine.Style{FG: t.theme.Accent, Bold: true}
	for x := 0; x < t.width; x++ {
		screen.Set(x, sepY, '═', borderStyle)
	}

	// Draw suggestions line if we have any
	if len(suggestions) > 0 {
		sugY := t.height - 3
		sugText := "  " + strings.Join(suggestions, "  ")
		if len(sugText) > t.width-2 {
			sugText = sugText[:t.width-5] + "..."
		}
		screen.DrawText(0, sugY, sugText, engine.Style{FG: t.theme.Secondary})
	}

	// Draw command input
	inputY := t.height - 2
	inputStyle := engine.Style{FG: t.theme.Primary, Bold: true}
	screen.DrawText(0, inputY, ":"+t.commandInput+"█", inputStyle)

	// Draw hint on the right
	hint := "Tab Complete │ Esc Cancel │ Enter Run"
	hintX := t.width - len(hint)
	if hintX > len(t.commandInput)+3 {
		screen.DrawText(hintX, inputY, hint, engine.Style{FG: t.theme.TextDim})
	}
}

func (t *CustomTUI) getCommandSuggestions() []string {
	if t.commandInput == "" {
		// Show all available commands when input is empty
		return []string{"quit", "play", "animate", "theme", "help", "games", "store"}
	}

	// Parse what we have so far
	cmdName, args := ui.ParseCommandLine(t.commandInput)

	// If we're still typing the command name
	if len(args) == 0 && !strings.Contains(t.commandInput, " ") {
		return t.commandRegistry.Complete(cmdName)
	}

	// If we have a command, get argument completions
	if cmd, ok := t.commandRegistry.Get(cmdName); ok {
		return cmd.Complete(args)
	}

	return nil
}

func (t *CustomTUI) renderCommandMessage(screen *engine.Screen) {
	// Show message or error temporarily at the top
	y := 0
	if t.commandError != nil {
		msg := "Error: " + t.commandError.Error()
		screen.DrawTextCentered(y, msg, engine.Style{FG: t.theme.Error, Bold: true})
	} else if t.commandMessage != "" {
		screen.DrawTextCentered(y, t.commandMessage, engine.Style{FG: t.theme.Success, Bold: true})
	}
}

func (t *CustomTUI) renderHome(screen *engine.Screen) {
	// Draw logo with style (24 chars wide, 5 lines tall)
	logoWidth := 24
	logoY := 2
	ui.DrawLogo(screen, (t.width-logoWidth)/2, logoY, engine.Style{FG: t.theme.Primary, Bold: true})

	// Draw subtitle
	subtitle := "Terminal Entertainment"
	screen.DrawTextCentered(8, subtitle, engine.Style{FG: t.theme.TextDim})

	// Draw menu with extra padding
	t.homeMenu.CenterOn(t.width, t.height)
	t.homeMenu.Y = 10
	t.homeMenu.Render(screen)

	// Draw footer bar with double line border
	footerY := t.height - 3
	borderStyle := engine.Style{FG: t.theme.Border}

	// Draw footer border
	for x := 0; x < t.width; x++ {
		screen.Set(x, footerY, '═', borderStyle)
	}

	// Draw footer text
	footer := "↑↓ Navigate │ Enter Select │ : Command │ Q Quit"
	screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})

	// Version in bottom right
	version := "v1.0.0"
	screen.DrawText(t.width-len(version)-2, t.height-1, version,
		engine.Style{FG: t.theme.TextDim})
}

func (t *CustomTUI) renderGames(screen *engine.Screen) {
	// Draw title with decorative border
	title := "╔═══ Games ═══╗"
	screen.DrawTextCentered(1, title, engine.Style{FG: t.theme.Primary, Bold: true})

	// Count available games
	gameCount := len(t.gamesMenu.Items) - 1 // Exclude "back" item
	countText := fmt.Sprintf("[%d available]", gameCount)
	screen.DrawTextCentered(2, countText, engine.Style{FG: t.theme.TextDim})

	// Draw menu
	t.gamesMenu.CenterOn(t.width, t.height)
	t.gamesMenu.Y = 4
	t.gamesMenu.Render(screen)

	// Draw footer bar
	footerY := t.height - 3
	borderStyle := engine.Style{FG: t.theme.Border}
	for x := 0; x < t.width; x++ {
		screen.Set(x, footerY, '═', borderStyle)
	}

	footer := "↑↓ Navigate │ Enter Play │ Esc Back"
	screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})
}

func (t *CustomTUI) renderAnimations(screen *engine.Screen) {
	// Draw title with decorative border
	title := "╔═══ Animations ═══╗"
	screen.DrawTextCentered(1, title, engine.Style{FG: t.theme.Primary, Bold: true})

	// Count available animations
	animCount := len(t.animationsMenu.Items) - 1 // Exclude "back" item
	countText := fmt.Sprintf("[%d available]", animCount)
	screen.DrawTextCentered(2, countText, engine.Style{FG: t.theme.TextDim})

	// Draw menu
	t.animationsMenu.CenterOn(t.width, t.height)
	t.animationsMenu.Y = 4
	t.animationsMenu.Render(screen)

	// Draw footer bar
	footerY := t.height - 3
	borderStyle := engine.Style{FG: t.theme.Border}
	for x := 0; x < t.width; x++ {
		screen.Set(x, footerY, '═', borderStyle)
	}

	footer := "↑↓ Navigate │ Enter Watch │ Esc Back"
	screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})
}

func (t *CustomTUI) renderStore(screen *engine.Screen) {
	// Draw title with decorative border
	title := "╔═══ Store ═══╗"
	screen.DrawTextCentered(1, title, engine.Style{FG: t.theme.Primary, Bold: true})

	// Loading state with spinner
	if t.storeLoading {
		// Draw a loading spinner
		spinnerChars := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
		spinnerIdx := (int(time.Now().UnixNano()/100000000) % len(spinnerChars))
		spinnerText := string(spinnerChars[spinnerIdx]) + " Loading store..."
		screen.DrawTextCentered(t.height/2, spinnerText, engine.Style{FG: t.theme.Warning})

		// Footer bar
		footerY := t.height - 3
		borderStyle := engine.Style{FG: t.theme.Border}
		for x := 0; x < t.width; x++ {
			screen.Set(x, footerY, '═', borderStyle)
		}
		footer := "Esc Back"
		screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})
		return
	}

	// Error state
	if t.storeError != "" {
		screen.DrawTextCentered(t.height/2-1, "╔══ Error loading store ══╗", engine.Style{FG: t.theme.Error, Bold: true})
		screen.DrawTextCentered(t.height/2, t.storeError, engine.Style{FG: t.theme.Error})
		screen.DrawTextCentered(t.height/2+2, "Press Enter to retry", engine.Style{FG: t.theme.TextDim})

		// Footer bar
		footerY := t.height - 3
		borderStyle := engine.Style{FG: t.theme.Border}
		for x := 0; x < t.width; x++ {
			screen.Set(x, footerY, '═', borderStyle)
		}
		footer := "Enter Retry │ Esc Back"
		screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})
		return
	}

	// Empty state
	if len(t.storeGames) == 0 {
		screen.DrawTextCentered(t.height/2, "No games available in store", engine.Style{FG: t.theme.TextDim})

		// Footer bar
		footerY := t.height - 3
		borderStyle := engine.Style{FG: t.theme.Border}
		for x := 0; x < t.width; x++ {
			screen.Set(x, footerY, '═', borderStyle)
		}
		footer := "Esc Back"
		screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})
		return
	}

	// Draw game list
	startY := 4
	maxVisible := (t.height - 7) / 3 // Each item takes 3 lines

	// Calculate scroll window
	scrollOffset := 0
	if t.storeSelectedIdx >= maxVisible {
		scrollOffset = t.storeSelectedIdx - maxVisible + 1
	}

	for i := 0; i < maxVisible && (scrollOffset+i) < len(t.storeGames); i++ {
		idx := scrollOffset + i
		game := t.storeGames[idx]
		y := startY + (i * 3)

		if y+2 >= t.height-3 {
			break
		}

		// Selection indicator with improved highlighting
		prefix := "  "
		style := engine.Style{FG: engine.ColorWhite}
		if idx == t.storeSelectedIdx {
			prefix = "▸ "
			style = engine.Style{FG: t.theme.Accent, Bold: true}
		}

		// Game name
		screen.DrawText(2, y, prefix+game.Name, style)

		// Description
		descStyle := engine.Style{FG: t.theme.TextDim}
		desc := game.Description
		if len(desc) > t.width-6 {
			desc = desc[:t.width-9] + "..."
		}
		screen.DrawText(4, y+1, desc, descStyle)

		// Author and version
		info := fmt.Sprintf("by %s • v%s", game.Author, game.Version)
		screen.DrawText(4, y+2, info, engine.Style{FG: t.theme.TextDim})
	}

	// Scroll indicator
	if len(t.storeGames) > maxVisible {
		indicator := fmt.Sprintf("(%d/%d)", t.storeSelectedIdx+1, len(t.storeGames))
		screen.DrawText(t.width-len(indicator)-2, 2, indicator, engine.Style{FG: t.theme.TextDim})
	}

	// Draw footer bar
	footerY := t.height - 3
	borderStyle := engine.Style{FG: t.theme.Border}
	for x := 0; x < t.width; x++ {
		screen.Set(x, footerY, '═', borderStyle)
	}

	footer := "↑↓ Navigate │ Enter Install │ Esc Back"
	screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})
}

func (t *CustomTUI) renderHelp(screen *engine.Screen) {
	// Draw title with decorative border
	title := "╔═══ Help ═══╗"
	screen.DrawTextCentered(1, title, engine.Style{FG: t.theme.Primary, Bold: true})

	help := []string{
		"Welcome to Slip!",
		"",
		"Slip is a terminal entertainment platform designed to",
		"keep you entertained while waiting for long-running",
		"commands like AI assistants, builds, or deploys.",
		"",
		"Navigation:",
		"  ↑/↓       Move selection up/down",
		"  Enter     Select item",
		"  Esc       Go back",
		"  :         Command mode",
		"  Q         Quit",
		"",
		"In Games:",
		"  P         Pause",
		"  R         Restart",
		"  Esc/Q     Return to menu",
		"",
		"Commands:  :play snake  :theme gruvbox  :help",
	}

	startY := 3
	for i, line := range help {
		if startY+i >= t.height-4 {
			break
		}
		if len(line) > 0 {
			x := (t.width - len(line)) / 2
			if x < 2 {
				x = 2
			}
			screen.DrawText(x, startY+i, line, engine.Style{FG: engine.ColorWhite})
		}
	}

	// Draw footer bar
	footerY := t.height - 3
	borderStyle := engine.Style{FG: t.theme.Border}
	for x := 0; x < t.width; x++ {
		screen.Set(x, footerY, '═', borderStyle)
	}

	footer := "Esc Back │ : Command"
	screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})
}

func (t *CustomTUI) renderSettings(screen *engine.Screen) {
	// Draw title with decorative border
	title := "╔═══ Settings ═══╗"
	screen.DrawTextCentered(1, title, engine.Style{FG: t.theme.Primary, Bold: true})

	// Load current config
	cfg, _ := state.LoadConfig()

	// Settings display
	settings := []string{
		"Current Settings:",
		"",
		fmt.Sprintf("  Theme:     %s", cfg.Theme),
		fmt.Sprintf("  TUI Mode:  %s", cfg.TUIMode),
		fmt.Sprintf("  Dock:      %s", cfg.DefaultDock),
		fmt.Sprintf("  Size:      %s", cfg.DefaultSize),
		fmt.Sprintf("  FPS:       %d", cfg.DefaultFPS),
		fmt.Sprintf("  Player:    %s", cfg.PlayerName),
		fmt.Sprintf("  Show FPS:  %v", cfg.ShowFPSCounter),
		"",
		"Use ':theme <name>' or 'slip config set <key> <value>'",
	}

	startY := 4
	for i, line := range settings {
		if startY+i >= t.height-4 {
			break
		}
		x := (t.width - len(line)) / 2
		if x < 2 {
			x = 2
		}
		screen.DrawText(x, startY+i, line, engine.Style{FG: engine.ColorWhite})
	}

	// Draw footer bar
	footerY := t.height - 3
	borderStyle := engine.Style{FG: t.theme.Border}
	for x := 0; x < t.width; x++ {
		screen.Set(x, footerY, '═', borderStyle)
	}

	footer := "Esc Back │ : Command"
	screen.DrawTextCentered(t.height-2, footer, engine.Style{FG: t.theme.TextDim})
}

func (t *CustomTUI) handleResize(width, height int, screen *engine.Screen) {
	t.width = width
	t.height = height
	screen.Resize(width, height)
	term.ClearScreenFull(os.Stdout)
}
