package tui

// TUI is the interface that both custom and Bubble Tea implementations must satisfy
type TUI interface {
	// Run starts the TUI application
	Run() error

	// Stop gracefully stops the TUI
	Stop()

	// SetScreen changes the current screen
	SetScreen(screen Screen)

	// GetScreen returns the current screen
	GetScreen() Screen

	// RegisterGameLauncher registers a game launcher callback
	RegisterGameLauncher(fn func(gameID string) error)

	// RegisterAnimationLauncher registers an animation launcher callback
	RegisterAnimationLauncher(fn func(animID string) error)
}

// Screen represents different screens in the application
type Screen int

const (
	ScreenHome Screen = iota
	ScreenGames
	ScreenAnimations
	ScreenStore
	ScreenSettings
	ScreenHelp
)

// Config holds TUI configuration
type Config struct {
	Width  int
	Height int
	FPS    int
	Theme  string
}

// MenuItem represents a menu item
type MenuItem struct {
	ID          string
	Label       string
	Description string
}

// GameInfo represents game metadata
type GameInfo struct {
	ID          string
	Name        string
	Description string
}
