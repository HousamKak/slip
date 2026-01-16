package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"slip/internal/animations"
	"slip/internal/games"
	"slip/internal/state"
	"slip/internal/store"
)

// BubbleTeaTUI implements the TUI interface using Bubble Tea
type BubbleTeaTUI struct {
	width  int
	height int
	screen Screen

	// Bubble Tea specific
	program *tea.Program

	// Menu state
	menuIndex int
	menuItems []MenuItem

	// Registries
	gameRegistry *games.Registry
	animRegistry *animations.Registry

	// Store
	storeClient      *store.Client
	storeGames       []store.GameEntry
	storeLoading     bool
	storeError       string

	// Callbacks
	gameLauncher func(string) error
	animLauncher func(string) error

	// Styles
	titleStyle      lipgloss.Style
	itemStyle       lipgloss.Style
	selectedStyle   lipgloss.Style
	descStyle       lipgloss.Style
	footerStyle     lipgloss.Style
	versionStyle    lipgloss.Style
	boxStyle        lipgloss.Style
	highlightStyle  lipgloss.Style
}

// Model for Bubble Tea
type model struct {
	tui *BubbleTeaTUI
}

// Message types
type storeLoadedMsg struct {
	games []store.GameEntry
	err   error
}

func loadStoreCmd(client *store.Client) tea.Cmd {
	return func() tea.Msg {
		games, err := client.ListGames()
		return storeLoadedMsg{games: games, err: err}
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case storeLoadedMsg:
		m.tui.storeLoading = false
		if msg.err != nil {
			m.tui.storeError = msg.err.Error()
		} else {
			m.tui.storeGames = msg.games
			m.tui.storeError = ""
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.tui.screen == ScreenHome {
				return m, tea.Quit
			}
			// Q from other screens goes back home
			m.tui.screen = ScreenHome
			m.tui.menuIndex = 0
			m.tui.updateMenuItems()

		case "up", "k":
			// In store, navigate store games directly
			if m.tui.screen == ScreenStore && len(m.tui.storeGames) > 0 {
				if m.tui.menuIndex > 0 {
					m.tui.menuIndex--
				}
			} else if m.tui.menuIndex > 0 {
				m.tui.menuIndex--
			}

		case "down", "j":
			// In store, navigate store games directly
			if m.tui.screen == ScreenStore && len(m.tui.storeGames) > 0 {
				if m.tui.menuIndex < len(m.tui.storeGames)-1 {
					m.tui.menuIndex++
				}
			} else if m.tui.menuIndex < len(m.tui.menuItems)-1 {
				m.tui.menuIndex++
			}

		case "enter":
			// Handle store screen separately
			if m.tui.screen == ScreenStore {
				// Retry on error
				if m.tui.storeError != "" {
					m.tui.storeLoading = true
					m.tui.storeError = ""
					return m, loadStoreCmd(m.tui.storeClient)
				}
				// Install game (placeholder for now)
				// In production, this would trigger installation
			} else {
				return m, m.tui.handleSelection()
			}

		case "esc":
			if m.tui.screen != ScreenHome {
				m.tui.screen = ScreenHome
				m.tui.menuIndex = 0
				m.tui.updateMenuItems()
			}
		}

	case tea.WindowSizeMsg:
		m.tui.width = msg.Width
		m.tui.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	switch m.tui.screen {
	case ScreenHome:
		return m.tui.renderHome()
	case ScreenGames:
		return m.tui.renderGames()
	case ScreenAnimations:
		return m.tui.renderAnimations()
	case ScreenStore:
		return m.tui.renderStore()
	case ScreenSettings:
		return m.tui.renderSettings()
	case ScreenHelp:
		return m.tui.renderHelp()
	default:
		return m.tui.renderHome()
	}
}

// NewBubbleTeaTUI creates a new Bubble Tea TUI
func NewBubbleTeaTUI(cfg Config) (*BubbleTeaTUI, error) {
	t := &BubbleTeaTUI{
		width:        cfg.Width,
		height:       cfg.Height,
		gameRegistry: games.NewRegistry(),
		animRegistry: animations.NewRegistry(),
		storeClient:  store.NewClient(),
	}

	// Initialize styles with cyan/bright colors (default theme)
	t.titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("51")). // Bright cyan
		MarginBottom(1).
		Align(lipgloss.Center)

	t.itemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("255")) // White

	t.selectedStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("226")). // Yellow
		Bold(true)

	t.descStyle = lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.Color("240")) // Dim gray

	t.footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	t.versionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	t.boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("51")). // Cyan
		Padding(1, 2)

	t.highlightStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")). // Cyan
		Bold(true)

	t.initMenus()
	return t, nil
}

func (t *BubbleTeaTUI) initMenus() {
	t.updateMenuItems()
}

func (t *BubbleTeaTUI) updateMenuItems() {
	switch t.screen {
	case ScreenHome:
		t.menuItems = []MenuItem{
			{ID: "games", Label: "Games", Description: "Play classic arcade games"},
			{ID: "animations", Label: "Animations", Description: "Watch relaxing visualizations"},
			{ID: "store", Label: "Store", Description: "Browse and install games/animations"},
			{ID: "settings", Label: "Settings", Description: "Configure Slip"},
			{ID: "help", Label: "Help", Description: "How to use Slip"},
			{ID: "quit", Label: "Quit", Description: "Exit Slip"},
		}

	case ScreenGames:
		t.menuItems = nil
		for _, info := range t.gameRegistry.List() {
			t.menuItems = append(t.menuItems, MenuItem{
				ID:          info.ID,
				Label:       info.Name,
				Description: info.Description,
			})
		}
		t.menuItems = append(t.menuItems, MenuItem{
			ID:    "back",
			Label: "← Back",
		})

	case ScreenAnimations:
		t.menuItems = nil
		for _, info := range t.animRegistry.List() {
			t.menuItems = append(t.menuItems, MenuItem{
				ID:          info.ID,
				Label:       info.Name,
				Description: info.Description,
			})
		}
		t.menuItems = append(t.menuItems, MenuItem{
			ID:    "back",
			Label: "← Back",
		})
	}
}

// Run starts the TUI
func (t *BubbleTeaTUI) Run() error {
	m := model{tui: t}
	t.program = tea.NewProgram(m, tea.WithAltScreen())
	_, err := t.program.Run()
	return err
}

// Stop gracefully stops the TUI
func (t *BubbleTeaTUI) Stop() {
	if t.program != nil {
		t.program.Quit()
	}
}

// SetScreen changes the current screen
func (t *BubbleTeaTUI) SetScreen(screen Screen) {
	t.screen = screen
	t.menuIndex = 0
	t.updateMenuItems()
}

// GetScreen returns the current screen
func (t *BubbleTeaTUI) GetScreen() Screen {
	return t.screen
}

// RegisterGameLauncher registers a game launcher callback
func (t *BubbleTeaTUI) RegisterGameLauncher(fn func(gameID string) error) {
	t.gameLauncher = fn
}

// RegisterAnimationLauncher registers an animation launcher callback
func (t *BubbleTeaTUI) RegisterAnimationLauncher(fn func(animID string) error) {
	t.animLauncher = fn
}

func (t *BubbleTeaTUI) handleSelection() tea.Cmd {
	if t.menuIndex >= len(t.menuItems) {
		return nil
	}

	selected := t.menuItems[t.menuIndex]

	switch t.screen {
	case ScreenHome:
		switch selected.ID {
		case "games":
			t.screen = ScreenGames
			t.menuIndex = 0
			t.updateMenuItems()
		case "animations":
			t.screen = ScreenAnimations
			t.menuIndex = 0
			t.updateMenuItems()
		case "store":
			t.screen = ScreenStore
			// Load store games if not loaded yet
			if len(t.storeGames) == 0 && !t.storeLoading && t.storeError == "" {
				t.storeLoading = true
				return loadStoreCmd(t.storeClient)
			}
		case "settings":
			t.screen = ScreenSettings
		case "help":
			t.screen = ScreenHelp
		case "quit":
			return tea.Quit
		}

	case ScreenGames:
		if selected.ID == "back" {
			t.screen = ScreenHome
			t.menuIndex = 0
			t.updateMenuItems()
		} else if t.gameLauncher != nil {
			// Launch game - this will take over the terminal
			t.program.ReleaseTerminal()
			t.gameLauncher(selected.ID)
			t.program.RestoreTerminal()
		}

	case ScreenAnimations:
		if selected.ID == "back" {
			t.screen = ScreenHome
			t.menuIndex = 0
			t.updateMenuItems()
		} else if t.animLauncher != nil {
			// Launch animation - this will take over the terminal
			t.program.ReleaseTerminal()
			t.animLauncher(selected.ID)
			t.program.RestoreTerminal()
		}
	}

	return nil
}

func (t *BubbleTeaTUI) renderHome() string {
	var b strings.Builder

	// Logo
	logo := t.titleStyle.Render("S L I P")
	b.WriteString(strings.Repeat("\n", 2))
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, logo))
	b.WriteString("\n")

	// Subtitle
	subtitle := t.descStyle.Render("Terminal Entertainment Platform")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, subtitle))
	b.WriteString("\n\n")

	// Menu items
	for i, item := range t.menuItems {
		var line string
		style := t.itemStyle
		prefix := fmt.Sprintf("%d. ", i+1)

		if i == t.menuIndex {
			style = t.selectedStyle
			prefix = fmt.Sprintf("▶ %d. ", i+1)
		}

		line = style.Render(prefix + item.Label)
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, line))
		b.WriteString("\n")

		if item.Description != "" {
			desc := t.descStyle.Render(item.Description)
			b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, desc))
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString(strings.Repeat("\n", max(0, t.height-b.Len()-5)))
	footer := t.footerStyle.Render("[↑↓] Navigate  [Enter] Select  [Q] Quit")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, footer))
	b.WriteString("\n")

	version := t.versionStyle.Render("v1.0.0")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Right, version))

	return b.String()
}

func (t *BubbleTeaTUI) renderGames() string {
	return t.renderMenuScreen("Games", "[↑↓] Navigate  [Enter] Play  [Esc] Back")
}

func (t *BubbleTeaTUI) renderAnimations() string {
	return t.renderMenuScreen("Animations", "[↑↓] Navigate  [Enter] Watch  [Esc] Back")
}

func (t *BubbleTeaTUI) renderMenuScreen(title, footer string) string {
	var b strings.Builder

	// Title
	titleText := t.titleStyle.Render(title)
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, titleText))
	b.WriteString("\n\n")

	// Menu items
	for i, item := range t.menuItems {
		var line string
		style := t.itemStyle
		prefix := fmt.Sprintf("%d. ", i+1)

		if i == t.menuIndex {
			style = t.selectedStyle
			prefix = fmt.Sprintf("▶ %d. ", i+1)
		}

		line = style.Render(prefix + item.Label)
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, line))
		b.WriteString("\n")

		if item.Description != "" && i == t.menuIndex {
			desc := t.descStyle.Render(item.Description)
			b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, desc))
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString(strings.Repeat("\n", max(0, t.height-b.Len()-3)))
	footerText := t.footerStyle.Render(footer)
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, footerText))

	return b.String()
}

func (t *BubbleTeaTUI) renderStore() string {
	var b strings.Builder

	// Title
	titleText := t.titleStyle.Render("Store")
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, titleText))
	b.WriteString("\n\n")

	// Loading state
	if t.storeLoading {
		loading := t.descStyle.Render("Loading store...")
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, loading))
		b.WriteString(strings.Repeat("\n", max(0, t.height-6)))
		footer := t.footerStyle.Render("[Esc] Back")
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, footer))
		return b.String()
	}

	// Error state
	if t.storeError != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		errorText := errorStyle.Render("Error loading store:")
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, errorText))
		b.WriteString("\n")
		errMsg := t.descStyle.Render(t.storeError)
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, errMsg))
		b.WriteString("\n\n")
		retry := t.descStyle.Render("Press Enter to retry")
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, retry))
		b.WriteString(strings.Repeat("\n", max(0, t.height-9)))
		footer := t.footerStyle.Render("[Enter] Retry  [Esc] Back")
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, footer))
		return b.String()
	}

	// Empty state
	if len(t.storeGames) == 0 {
		empty := t.descStyle.Render("No games available in store")
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, empty))
		b.WriteString(strings.Repeat("\n", max(0, t.height-6)))
		footer := t.footerStyle.Render("[Esc] Back")
		b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, footer))
		return b.String()
	}

	// Draw game list
	maxVisible := (t.height - 8) / 4 // Each item takes ~4 lines

	// Calculate scroll window
	scrollOffset := 0
	if t.menuIndex >= maxVisible {
		scrollOffset = t.menuIndex - maxVisible + 1
	}

	for i := 0; i < maxVisible && (scrollOffset+i) < len(t.storeGames); i++ {
		idx := scrollOffset + i
		game := t.storeGames[idx]

		var item strings.Builder

		// Game name
		name := game.Name
		if idx == t.menuIndex {
			item.WriteString(t.selectedStyle.Render("> " + name))
		} else {
			item.WriteString(t.itemStyle.Render("  " + name))
		}
		item.WriteString("\n")

		// Description
		desc := game.Description
		if len(desc) > t.width-8 {
			desc = desc[:t.width-11] + "..."
		}
		item.WriteString(t.descStyle.Render("    " + desc))
		item.WriteString("\n")

		// Author and version
		info := fmt.Sprintf("    by %s • v%s", game.Author, game.Version)
		item.WriteString(t.versionStyle.Render(info))
		item.WriteString("\n\n")

		b.WriteString(item.String())
	}

	// Scroll indicator
	if len(t.storeGames) > maxVisible {
		indicator := fmt.Sprintf("(%d/%d)", t.menuIndex+1, len(t.storeGames))
		b.WriteString(t.footerStyle.Render(indicator))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString(strings.Repeat("\n", max(0, t.height-strings.Count(b.String(), "\n")-4)))
	footer := t.footerStyle.Render("[↑↓/j/k] Navigate  [Enter] Install  [Esc] Back")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, footer))

	return b.String()
}

func (t *BubbleTeaTUI) renderHelp() string {
	var b strings.Builder

	// Title
	titleText := t.titleStyle.Render("Help")
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, titleText))
	b.WriteString("\n\n")

	// Help text
	help := []string{
		t.highlightStyle.Render("Welcome to Slip!"),
		"",
		"Slip is a terminal entertainment platform designed to",
		"keep you entertained while waiting for long-running",
		"commands like AI assistants, builds, or deploys.",
		"",
		t.highlightStyle.Render("Navigation:"),
		"  [↑/↓]     Move selection up/down",
		"  [Enter]   Select item",
		"  [Esc]     Go back",
		"  [Q]       Quit",
		"",
		t.highlightStyle.Render("In Games:"),
		"  [P]       Pause",
		"  [R]       Restart",
		"  [Q]       Return to menu",
		"",
		t.highlightStyle.Render("Tips:"),
		"  - Use 'slip summon' in tmux to run alongside your work",
		"  - Use 'slip play snake' to jump right into a game",
	}

	box := t.boxStyle.Render(strings.Join(help, "\n"))
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, box))

	// Footer
	b.WriteString(strings.Repeat("\n", max(0, t.height-b.Len()-3)))
	footer := t.footerStyle.Render("[Esc] or [Enter] to return")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, footer))

	return b.String()
}

func (t *BubbleTeaTUI) renderSettings() string {
	var b strings.Builder

	// Title
	titleText := t.titleStyle.Render("Settings")
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, titleText))
	b.WriteString("\n\n")

	// Load current config
	cfg, _ := state.LoadConfig()

	// Settings display
	settings := []string{
		t.highlightStyle.Render("Current Settings:"),
		"",
		fmt.Sprintf("  Theme:        %s", cfg.Theme),
		fmt.Sprintf("  TUI Mode:     %s", cfg.TUIMode),
		fmt.Sprintf("  Dock:         %s", cfg.DefaultDock),
		fmt.Sprintf("  Size:         %s", cfg.DefaultSize),
		fmt.Sprintf("  FPS:          %d", cfg.DefaultFPS),
		fmt.Sprintf("  Player:       %s", cfg.PlayerName),
		fmt.Sprintf("  Show FPS:     %v", cfg.ShowFPSCounter),
		"",
		"Use 'slip config set <key> <value>' to change settings",
	}

	box := t.boxStyle.Render(strings.Join(settings, "\n"))
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, box))

	// Footer
	b.WriteString(strings.Repeat("\n", max(0, t.height-b.Len()-3)))
	footer := t.footerStyle.Render("[Esc] or [Enter] to return")
	b.WriteString(lipgloss.PlaceHorizontal(t.width, lipgloss.Center, footer))

	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
