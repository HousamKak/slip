package engine

import (
	"fmt"
	"io"
	"os"
	"time"

	"slip/internal/state"
	"slip/internal/term"
)

// Engine is the main game engine
type Engine struct {
	screen      *Screen
	input       *InputReader
	currentGame Game
	running     bool
	paused      bool
	targetFPS   int
	width       int
	height      int
	writer      io.Writer

	// Callbacks
	onExit func()

	// Stats
	frameCount int
	lastFPS    float64
	showFPS    bool

	// Config
	config     state.Config
	scoreBoard *state.ScoreBoard
	playerName string
}

// NewEngine creates a new game engine
func NewEngine(width, height, fps int) *Engine {
	if width < 10 {
		width = 40
	}
	if height < 5 {
		height = 20
	}
	if fps < 1 {
		fps = 30
	}

	// Load config
	cfg, _ := state.LoadConfig()

	// Load scoreboard
	sb, _ := state.LoadScores()
	if sb == nil {
		sb = state.NewScoreBoard()
	}

	playerName := cfg.PlayerName
	if playerName == "" {
		playerName = "Player"
	}

	e := &Engine{
		width:      width,
		height:     height,
		targetFPS:  fps,
		writer:     os.Stdout,
		config:     cfg,
		scoreBoard: sb,
		playerName: playerName,
		showFPS:    cfg.ShowFPSCounter,
	}
	e.screen = NewScreen(width, height, e.writer)
	e.input = NewInputReader()

	return e
}

// SetGame sets the current game
func (e *Engine) SetGame(game Game) error {
	if e.currentGame != nil {
		e.currentGame.Stop()
	}
	e.currentGame = game

	gameInfo := game.Info()

	ctx := &GameContext{
		Width:      e.width,
		Height:     e.height,
		FPS:        e.targetFPS,
		StateDir:   state.GamesDir(),
		PlayerName: e.playerName,
		OnExit:     e.requestExit,
		OnScore: func(score int) {
			// Save score to scoreboard
			if e.scoreBoard != nil && score > 0 {
				_ = e.scoreBoard.AddScore(gameInfo.ID, state.Score{
					Value:  score,
					Player: e.playerName,
				})
			}
		},
		GetConfig: func(key string) string {
			val, _ := state.GetConfig(key)
			return val
		},
		GetHighScore: func() int {
			if e.scoreBoard != nil {
				return e.scoreBoard.GetHighScore(gameInfo.ID)
			}
			return 0
		},
		IsHighScore: func(score int) bool {
			if e.scoreBoard != nil {
				return e.scoreBoard.IsHighScore(gameInfo.ID, score)
			}
			return true
		},
	}

	if err := game.Init(ctx); err != nil {
		return err
	}

	return game.Start()
}

// SetOnExit sets the callback for when the engine exits
func (e *Engine) SetOnExit(fn func()) {
	e.onExit = fn
}

// SetShowFPS enables FPS display
func (e *Engine) SetShowFPS(show bool) {
	e.showFPS = show
}

// Width returns the screen width
func (e *Engine) Width() int {
	return e.width
}

// Height returns the screen height
func (e *Engine) Height() int {
	return e.height
}

// Run starts the game loop
func (e *Engine) Run() error {
	// Set up terminal
	term.EnableAlternateScreen(e.writer)
	term.HideCursor(e.writer)
	term.ClearScreenFull(e.writer)

	defer func() {
		term.ShowCursor(e.writer)
		term.DisableAlternateScreen(e.writer)
		term.ResetStyle(e.writer)
	}()

	// Start input reader
	if err := e.input.Start(); err != nil {
		return fmt.Errorf("failed to start input reader: %w", err)
	}
	defer e.input.Stop()

	e.running = true
	ticker := time.NewTicker(time.Second / time.Duration(e.targetFPS))
	defer ticker.Stop()

	fpsTimer := time.NewTicker(time.Second)
	defer fpsTimer.Stop()

	resizeTimer := time.NewTicker(500 * time.Millisecond)
	defer resizeTimer.Stop()

	lastTime := time.Now()
	frameThisSecond := 0

	for e.running {
		// Process all pending inputs first for responsiveness
		// This prevents input lag when inputs arrive faster than frame rate
		for {
			select {
			case input := <-e.input.Channel():
				if input.IsQuit() {
					e.running = false
					break
				}
				if input.Key == KeyEscape || input.Key == KeyRune && (input.Rune == 'q' || input.Rune == 'Q') {
					e.running = false
					break
				}
				if e.currentGame != nil && !e.paused {
					e.currentGame.HandleInput(input)
				}
			default:
				// No more inputs to process
				goto doneProcessingInputs
			}
		}
	doneProcessingInputs:

		// Now wait for next frame or other events
		select {
		case input := <-e.input.Channel():
			// Handle any new input that arrived while we were waiting
			if input.IsQuit() {
				e.running = false
				continue
			}
			if input.Key == KeyEscape || input.Key == KeyRune && (input.Rune == 'q' || input.Rune == 'Q') {
				e.running = false
				continue
			}
			if e.currentGame != nil && !e.paused {
				e.currentGame.HandleInput(input)
			}

		case <-ticker.C:
			now := time.Now()
			dt := now.Sub(lastTime).Seconds()
			lastTime = now

			if e.currentGame != nil && !e.paused {
				if err := e.currentGame.Update(dt); err != nil {
					// Game requested exit
					e.running = false
					continue
				}
			}

			// Render
			e.screen.Clear()
			if e.currentGame != nil {
				e.currentGame.Render(e.screen)
			}

			// Draw FPS if enabled
			if e.showFPS {
				fpsText := fmt.Sprintf("FPS: %.1f", e.lastFPS)
				e.screen.DrawText(e.width-len(fpsText)-1, 0, fpsText,
					Style{FG: ColorBrightBlack})
			}

			e.screen.Flush()
			frameThisSecond++
			e.frameCount++

		case <-fpsTimer.C:
			e.lastFPS = float64(frameThisSecond)
			frameThisSecond = 0

		case <-resizeTimer.C:
			// Check for terminal resize
			w, h, err := term.Size()
			if err == nil && (w != e.width || h != e.height) {
				e.handleResize(w, h)
			}
		}
	}

	if e.currentGame != nil {
		e.currentGame.Stop()
	}

	if e.onExit != nil {
		e.onExit()
	}

	return nil
}

// Stop stops the engine
func (e *Engine) Stop() {
	e.running = false
}

// handleResize handles terminal resize events
func (e *Engine) handleResize(width, height int) {
	e.width = width
	e.height = height
	e.screen.Resize(width, height)
	term.ClearScreenFull(e.writer)
}

// Pause pauses the engine
func (e *Engine) Pause() {
	e.paused = true
}

// Resume resumes the engine
func (e *Engine) Resume() {
	e.paused = false
}

// IsPaused returns whether the engine is paused
func (e *Engine) IsPaused() bool {
	return e.paused
}

// TogglePause toggles pause state
func (e *Engine) TogglePause() {
	e.paused = !e.paused
}

func (e *Engine) requestExit() {
	e.running = false
}

// Screen returns the screen for direct access
func (e *Engine) Screen() *Screen {
	return e.screen
}

// RunOnce runs a single frame (useful for testing or external loops)
func (e *Engine) RunOnce(dt float64) error {
	// Process any pending input
	select {
	case input := <-e.input.Channel():
		if input.IsQuit() {
			return fmt.Errorf("quit requested")
		}
		if e.currentGame != nil {
			e.currentGame.HandleInput(input)
		}
	default:
	}

	// Update
	if e.currentGame != nil && !e.paused {
		if err := e.currentGame.Update(dt); err != nil {
			return err
		}
	}

	// Render
	e.screen.Clear()
	if e.currentGame != nil {
		e.currentGame.Render(e.screen)
	}
	e.screen.Flush()

	return nil
}
