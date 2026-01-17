package app

import (
	"fmt"

	"slip/internal/animations"
	"slip/internal/engine"
	"slip/internal/games"
	"slip/internal/plugins"
	"slip/internal/tui"
)

// App is the main Slip application
type App struct {
	tui          tui.TUI
	gameRegistry *games.Registry
	animRegistry *animations.Registry
	width        int
	height       int
	fps          int
}

// New creates a new Slip application
func New(width, height, fps int) *App {
	app := &App{
		width:        width,
		height:       height,
		fps:          fps,
		gameRegistry: games.NewRegistry(),
		animRegistry: animations.NewRegistry(),
	}

	// Load plugins and register them
	loader := plugins.NewLoader()
	if err := loader.Scan(); err == nil {
		for _, manifest := range loader.List() {
			if manifest.Category == "game" {
				app.gameRegistry.RegisterPlugin(manifest, loader)
			} else if manifest.Category == "animation" {
				app.animRegistry.RegisterPlugin(manifest, loader)
			}
		}
	}

	// Create TUI using factory
	cfg := tui.Config{
		Width:  width,
		Height: height,
		FPS:    fps,
	}

	t, err := tui.NewTUI(cfg)
	if err != nil {
		// Fallback to custom
		t, _ = tui.NewCustomTUI(cfg)
	}

	app.tui = t

	// Register callbacks
	t.RegisterGameLauncher(app.launchGame)
	t.RegisterAnimationLauncher(app.launchAnimation)

	return app
}

// Run starts the application
func (a *App) Run() error {
	return a.tui.Run()
}

func (a *App) launchGame(id string) error {
	game, ok := a.gameRegistry.Create(id)
	if !ok {
		return fmt.Errorf("game not found: %s", id)
	}

	eng := engine.NewEngine(a.width, a.height, a.fps)
	if provider, ok := a.tui.(interface{ InputReader() *engine.InputReader }); ok {
		if reader := provider.InputReader(); reader != nil {
			eng.UseInputReader(reader)
		}
	}
	if err := eng.SetGame(game); err != nil {
		return err
	}
	return eng.Run()
}

func (a *App) launchAnimation(id string) error {
	anim, ok := a.animRegistry.CreateAsGame(id)
	if !ok {
		return fmt.Errorf("animation not found: %s", id)
	}

	eng := engine.NewEngine(a.width, a.height, a.fps)
	if provider, ok := a.tui.(interface{ InputReader() *engine.InputReader }); ok {
		if reader := provider.InputReader(); reader != nil {
			eng.UseInputReader(reader)
		}
	}
	if err := eng.SetGame(anim); err != nil {
		return err
	}
	return eng.Run()
}

// RunGame runs a specific game directly
func (a *App) RunGame(id string) error {
	return a.launchGame(id)
}

// RunAnimation runs a specific animation directly
func (a *App) RunAnimation(id string) error {
	return a.launchAnimation(id)
}

// ListGames returns list of available games
func (a *App) ListGames() []engine.GameInfo {
	return a.gameRegistry.List()
}

// ListAnimations returns list of available animations
func (a *App) ListAnimations() []engine.GameInfo {
	return a.animRegistry.List()
}
