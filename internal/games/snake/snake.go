package snake

import (
	"fmt"
	"math/rand"
	"time"

	"slip/internal/engine"
)

// Game implements the classic Snake game
type Game struct {
	ctx       *engine.GameContext
	snake     []engine.Point
	direction engine.Direction
	nextDir   engine.Direction
	food      engine.Point
	score     int
	highScore int
	gameOver  bool
	paused    bool
	moveTimer float64
	moveDelay float64
	rng       *rand.Rand

	// Play area bounds (inside the border)
	areaX, areaY int
	areaW, areaH int
}

// New creates a new Snake game
func New() *Game {
	return &Game{
		moveDelay: 0.12,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Info returns game metadata
func (g *Game) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "snake",
		Name:        "Snake",
		Description: "Classic snake game - eat food, grow longer, don't hit walls or yourself!",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    20,
		MinHeight:   12,
		Category:    "game",
	}
}

// Init initializes the game
func (g *Game) Init(ctx *engine.GameContext) error {
	g.ctx = ctx

	// Define play area (with margin for border and UI)
	g.areaX = 1
	g.areaY = 2
	g.areaW = ctx.Width - 2
	g.areaH = ctx.Height - 4

	// Load high score from context
	if ctx.GetHighScore != nil {
		g.highScore = ctx.GetHighScore()
	}

	g.reset()
	return nil
}

func (g *Game) reset() {
	// Start snake in center
	centerX := g.areaX + g.areaW/2
	centerY := g.areaY + g.areaH/2

	g.snake = []engine.Point{
		{X: centerX, Y: centerY},
		{X: centerX - 1, Y: centerY},
		{X: centerX - 2, Y: centerY},
	}
	g.direction = engine.DirRight
	g.nextDir = engine.DirRight
	g.score = 0
	g.gameOver = false
	g.paused = false
	g.moveTimer = 0
	g.moveDelay = 0.12 // Reset speed

	g.spawnFood()
}

func (g *Game) spawnFood() {
	for {
		g.food = engine.Point{
			X: g.areaX + g.rng.Intn(g.areaW),
			Y: g.areaY + g.rng.Intn(g.areaH),
		}

		// Make sure food isn't on the snake
		collision := false
		for _, p := range g.snake {
			if p.Equals(g.food) {
				collision = true
				break
			}
		}
		if !collision {
			break
		}
	}
}

func (g *Game) endGame() {
	g.gameOver = true
	if g.score > g.highScore {
		g.highScore = g.score
	}
	// Report score to engine for persistence
	if g.ctx.OnScore != nil && g.score > 0 {
		g.ctx.OnScore(g.score)
	}
}

// Start is called when the game starts
func (g *Game) Start() error {
	return nil
}

// Stop is called when the game stops
func (g *Game) Stop() error {
	return nil
}

// Update updates the game state
func (g *Game) Update(dt float64) error {
	if g.gameOver || g.paused {
		return nil
	}

	g.moveTimer += dt
	if g.moveTimer < g.moveDelay {
		return nil
	}
	g.moveTimer = 0

	// Apply direction change
	g.direction = g.nextDir

	// Calculate new head position
	head := g.snake[0]
	delta := g.direction.Delta()
	newHead := head.Add(delta)

	// Check wall collision
	if newHead.X < g.areaX || newHead.X >= g.areaX+g.areaW ||
		newHead.Y < g.areaY || newHead.Y >= g.areaY+g.areaH {
		g.endGame()
		return nil
	}

	// Check self collision
	for _, p := range g.snake {
		if p.Equals(newHead) {
			g.endGame()
			return nil
		}
	}

	// Move snake
	g.snake = append([]engine.Point{newHead}, g.snake...)

	// Check food collision
	if newHead.Equals(g.food) {
		g.score += 10
		g.spawnFood()
		// Speed up slightly
		if g.moveDelay > 0.04 {
			g.moveDelay -= 0.002
		}
	} else {
		// Remove tail if no food eaten
		g.snake = g.snake[:len(g.snake)-1]
	}

	return nil
}

// Render draws the game
func (g *Game) Render(screen *engine.Screen) error {
	// Draw border
	screen.DrawBox(0, 1, g.ctx.Width, g.ctx.Height-2, engine.BoxStyleSingle,
		engine.Style{FG: engine.ColorCyan})

	// Draw header
	title := "SNAKE"
	screen.DrawText(2, 0, title, engine.Style{FG: engine.ColorBrightCyan, Bold: true})

	scoreText := fmt.Sprintf("Score: %d", g.score)
	screen.DrawText(g.ctx.Width-len(scoreText)-2, 0, scoreText,
		engine.Style{FG: engine.ColorBrightWhite})

	if g.highScore > 0 {
		hiText := fmt.Sprintf("Hi: %d", g.highScore)
		screen.DrawText(g.ctx.Width/2-len(hiText)/2, 0, hiText,
			engine.Style{FG: engine.ColorYellow})
	}

	// Draw snake
	for i, p := range g.snake {
		style := engine.Style{FG: engine.ColorGreen}
		char := '█'
		if i == 0 {
			// Head
			style = engine.Style{FG: engine.ColorBrightGreen, Bold: true}
			switch g.direction {
			case engine.DirUp:
				char = '▲'
			case engine.DirDown:
				char = '▼'
			case engine.DirLeft:
				char = '◀'
			case engine.DirRight:
				char = '▶'
			}
		}
		screen.Set(p.X, p.Y, char, style)
	}

	// Draw food
	screen.Set(g.food.X, g.food.Y, '●', engine.Style{FG: engine.ColorRed})

	// Draw footer/controls
	controls := "[←↑↓→] Move  [P] Pause  [R] Restart  [Q] Quit"
	if len(controls) > g.ctx.Width-4 {
		controls = "←↑↓→:Move P:Pause R:Restart Q:Quit"
	}
	screen.DrawText(2, g.ctx.Height-1, controls, engine.Style{FG: engine.ColorBrightBlack})

	// Draw game over overlay
	if g.gameOver {
		g.drawGameOver(screen)
	}

	// Draw pause overlay
	if g.paused && !g.gameOver {
		g.drawPaused(screen)
	}

	return nil
}

func (g *Game) drawGameOver(screen *engine.Screen) {
	boxW := 24
	boxH := 7
	boxX := (g.ctx.Width - boxW) / 2
	boxY := (g.ctx.Height - boxH) / 2

	// Draw box
	screen.DrawFilledBox(boxX, boxY, boxW, boxH, engine.BoxStyleDouble,
		engine.Style{FG: engine.ColorRed},
		' ', engine.Style{BG: engine.ColorBlack})

	// Draw text
	screen.DrawTextCentered(boxY+1, "GAME OVER", engine.Style{FG: engine.ColorBrightRed, Bold: true})
	screen.DrawTextCentered(boxY+3, fmt.Sprintf("Score: %d", g.score), engine.Style{FG: engine.ColorWhite})
	screen.DrawTextCentered(boxY+5, "[R] Restart  [Q] Quit", engine.Style{FG: engine.ColorBrightBlack})
}

func (g *Game) drawPaused(screen *engine.Screen) {
	boxW := 16
	boxH := 5
	boxX := (g.ctx.Width - boxW) / 2
	boxY := (g.ctx.Height - boxH) / 2

	screen.DrawFilledBox(boxX, boxY, boxW, boxH, engine.BoxStyleRounded,
		engine.Style{FG: engine.ColorYellow},
		' ', engine.Style{BG: engine.ColorBlack})

	screen.DrawTextCentered(boxY+1, "PAUSED", engine.Style{FG: engine.ColorBrightYellow, Bold: true})
	screen.DrawTextCentered(boxY+3, "[P] Resume", engine.Style{FG: engine.ColorBrightBlack})
}

// HandleInput processes input
func (g *Game) HandleInput(input engine.Input) error {
	// Handle restart
	if input.Key == engine.KeyRune && (input.Rune == 'r' || input.Rune == 'R') {
		g.reset()
		return nil
	}

	// Handle pause
	if input.Key == engine.KeyRune && (input.Rune == 'p' || input.Rune == 'P') {
		g.paused = !g.paused
		return nil
	}

	// Handle quit
	if input.Key == engine.KeyRune && (input.Rune == 'q' || input.Rune == 'Q') {
		if g.ctx.OnExit != nil {
			g.ctx.OnExit()
		}
		return nil
	}

	// If game over or paused, don't process movement
	if g.gameOver || g.paused {
		return nil
	}

	// Handle direction changes
	var newDir engine.Direction
	switch input.Key {
	case engine.KeyUp:
		newDir = engine.DirUp
	case engine.KeyDown:
		newDir = engine.DirDown
	case engine.KeyLeft:
		newDir = engine.DirLeft
	case engine.KeyRight:
		newDir = engine.DirRight
	case engine.KeyRune:
		switch input.Rune {
		case 'w', 'W':
			newDir = engine.DirUp
		case 's', 'S':
			newDir = engine.DirDown
		case 'a', 'A':
			newDir = engine.DirLeft
		case 'd', 'D':
			newDir = engine.DirRight
		default:
			return nil
		}
	default:
		return nil
	}

	// Prevent 180-degree turns
	if newDir != engine.DirNone && newDir != g.direction.Opposite() {
		g.nextDir = newDir
	}

	return nil
}
