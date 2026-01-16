package breakout

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"slip/internal/engine"
)

// Brick represents a breakable brick
type Brick struct {
	X, Y   int
	Width  int
	Hits   int // Hits remaining
	Color  engine.Color
	Points int
}

// Game implements the Breakout game
type Game struct {
	ctx *engine.GameContext

	// Ball
	ballX, ballY   float64
	ballVX, ballVY float64
	ballLaunched   bool

	// Paddle
	paddleX int
	paddleW int

	// Bricks
	bricks []Brick

	// Game state
	score     int
	lives     int
	level     int
	highScore int
	gameOver  bool
	paused    bool
	victory   bool

	// Play area
	areaX, areaY int
	areaW, areaH int

	rng *rand.Rand
}

// New creates a new Breakout game
func New() *Game {
	return &Game{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Info returns game metadata
func (g *Game) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "breakout",
		Name:        "Breakout",
		Description: "Break all the bricks! Don't let the ball fall!",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    30,
		MinHeight:   18,
		Category:    "game",
	}
}

// Init initializes the game
func (g *Game) Init(ctx *engine.GameContext) error {
	g.ctx = ctx

	// Define play area
	g.areaX = 1
	g.areaY = 2
	g.areaW = ctx.Width - 2
	g.areaH = ctx.Height - 4

	g.paddleW = 7

	// Load high score
	if ctx.GetHighScore != nil {
		g.highScore = ctx.GetHighScore()
	}

	g.reset()
	return nil
}

func (g *Game) reset() {
	g.score = 0
	g.lives = 3
	g.level = 1
	g.gameOver = false
	g.victory = false
	g.createLevel()
	g.resetBall()
}

func (g *Game) createLevel() {
	g.bricks = nil

	// Calculate brick dimensions
	brickW := 4
	brickH := 1
	cols := (g.areaW - 4) / (brickW + 1)
	rows := 4 + g.level // More rows at higher levels

	if rows > 8 {
		rows = 8
	}

	startX := g.areaX + (g.areaW-cols*(brickW+1))/2
	startY := g.areaY + 2

	colors := []engine.Color{
		engine.ColorRed,
		engine.ColorYellow,
		engine.ColorGreen,
		engine.ColorCyan,
		engine.ColorMagenta,
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			hits := 1
			if row < 2 && g.level > 1 {
				hits = 2 // Top rows are tougher
			}

			g.bricks = append(g.bricks, Brick{
				X:      startX + col*(brickW+1),
				Y:      startY + row*brickH,
				Width:  brickW,
				Hits:   hits,
				Color:  colors[row%len(colors)],
				Points: (rows - row) * 10,
			})
		}
	}
}

func (g *Game) resetBall() {
	g.paddleX = g.areaX + g.areaW/2 - g.paddleW/2
	g.ballX = float64(g.paddleX + g.paddleW/2)
	g.ballY = float64(g.areaY + g.areaH - 3)
	g.ballVX = 0
	g.ballVY = 0
	g.ballLaunched = false
}

func (g *Game) launchBall() {
	if g.ballLaunched {
		return
	}
	g.ballLaunched = true
	speed := 20.0
	angle := math.Pi/4 + g.rng.Float64()*math.Pi/4 // 45-90 degrees
	if g.rng.Intn(2) == 0 {
		angle = math.Pi - angle
	}
	g.ballVX = speed * math.Cos(angle)
	g.ballVY = -speed * math.Sin(angle)
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
	if g.gameOver || g.paused || g.victory {
		return nil
	}

	if !g.ballLaunched {
		// Ball follows paddle
		g.ballX = float64(g.paddleX + g.paddleW/2)
		return nil
	}

	// Move ball
	g.ballX += g.ballVX * dt
	g.ballY += g.ballVY * dt

	// Wall collisions
	if g.ballX <= float64(g.areaX) {
		g.ballX = float64(g.areaX)
		g.ballVX = math.Abs(g.ballVX)
	}
	if g.ballX >= float64(g.areaX+g.areaW-1) {
		g.ballX = float64(g.areaX + g.areaW - 1)
		g.ballVX = -math.Abs(g.ballVX)
	}
	if g.ballY <= float64(g.areaY) {
		g.ballY = float64(g.areaY)
		g.ballVY = math.Abs(g.ballVY)
	}

	// Ball fell below paddle
	if g.ballY >= float64(g.areaY+g.areaH) {
		g.lives--
		if g.lives <= 0 {
			g.gameOver = true
			g.endGame()
		} else {
			g.resetBall()
		}
		return nil
	}

	// Paddle collision
	paddleTop := float64(g.areaY + g.areaH - 2)
	if g.ballY >= paddleTop && g.ballY <= paddleTop+1 &&
		g.ballX >= float64(g.paddleX) && g.ballX <= float64(g.paddleX+g.paddleW) {
		g.ballY = paddleTop - 1
		// Angle based on where ball hit paddle
		hitPos := (g.ballX - float64(g.paddleX)) / float64(g.paddleW)
		angle := (hitPos - 0.5) * math.Pi / 3 // -60 to +60 degrees
		speed := math.Sqrt(g.ballVX*g.ballVX + g.ballVY*g.ballVY)
		g.ballVX = speed * math.Sin(angle)
		g.ballVY = -speed * math.Cos(angle)
	}

	// Brick collisions
	ballBounds := engine.Rect{
		X: int(g.ballX), Y: int(g.ballY),
		Width: 1, Height: 1,
	}

	for i := len(g.bricks) - 1; i >= 0; i-- {
		brick := &g.bricks[i]
		brickBounds := engine.Rect{
			X: brick.X, Y: brick.Y,
			Width: brick.Width, Height: 1,
		}

		if ballBounds.Intersects(brickBounds) {
			brick.Hits--
			if brick.Hits <= 0 {
				g.score += brick.Points
				g.bricks = append(g.bricks[:i], g.bricks[i+1:]...)
			}

			// Bounce ball
			// Determine which side was hit
			ballCenterX := g.ballX
			brickCenterX := float64(brick.X) + float64(brick.Width)/2

			if math.Abs(ballCenterX-brickCenterX) > float64(brick.Width)/2*0.8 {
				g.ballVX = -g.ballVX
			} else {
				g.ballVY = -g.ballVY
			}
			break
		}
	}

	// Check victory
	if len(g.bricks) == 0 {
		g.level++
		if g.level > 5 {
			g.victory = true
			g.endGame()
		} else {
			g.createLevel()
			g.resetBall()
		}
	}

	return nil
}

func (g *Game) endGame() {
	if g.ctx.OnScore != nil && g.score > 0 {
		if g.score > g.highScore {
			g.highScore = g.score
		}
		g.ctx.OnScore(g.score)
	}
}

// Render draws the game
func (g *Game) Render(screen *engine.Screen) error {
	// Draw border
	screen.DrawBox(0, 1, g.ctx.Width, g.ctx.Height-2, engine.BoxStyleSingle,
		engine.Style{FG: engine.ColorCyan})

	// Draw header
	title := "BREAKOUT"
	screen.DrawText(2, 0, title, engine.Style{FG: engine.ColorBrightCyan, Bold: true})

	scoreText := fmt.Sprintf("Score: %d", g.score)
	screen.DrawText(g.ctx.Width/2-len(scoreText)/2, 0, scoreText,
		engine.Style{FG: engine.ColorBrightWhite})

	livesText := fmt.Sprintf("Lives: %d  Lvl: %d", g.lives, g.level)
	screen.DrawText(g.ctx.Width-len(livesText)-2, 0, livesText,
		engine.Style{FG: engine.ColorBrightYellow})

	// Draw bricks
	for _, brick := range g.bricks {
		style := engine.Style{FG: brick.Color}
		if brick.Hits > 1 {
			style.Bold = true
		}
		for i := 0; i < brick.Width; i++ {
			char := '█'
			if brick.Hits > 1 {
				char = '▓'
			}
			screen.Set(brick.X+i, brick.Y, char, style)
		}
	}

	// Draw paddle
	for i := 0; i < g.paddleW; i++ {
		screen.Set(g.paddleX+i, g.areaY+g.areaH-2, '═', engine.Style{FG: engine.ColorWhite, Bold: true})
	}

	// Draw ball
	if !g.gameOver && !g.victory {
		screen.Set(int(g.ballX), int(g.ballY), '●', engine.Style{FG: engine.ColorBrightWhite})
	}

	// Draw footer/controls
	var controls string
	if !g.ballLaunched {
		controls = "[←→] Move  [SPACE] Launch  [P] Pause  [Q] Quit"
	} else {
		controls = "[←→] Move  [P] Pause  [R] Restart  [Q] Quit"
	}
	if len(controls) > g.ctx.Width-4 {
		controls = "←→:Move SPACE:Launch P:Pause Q:Quit"
	}
	screen.DrawText(2, g.ctx.Height-1, controls, engine.Style{FG: engine.ColorBrightBlack})

	// Draw game over overlay
	if g.gameOver {
		g.drawGameOver(screen)
	}

	// Draw victory overlay
	if g.victory {
		g.drawVictory(screen)
	}

	// Draw pause overlay
	if g.paused && !g.gameOver && !g.victory {
		g.drawPaused(screen)
	}

	return nil
}

func (g *Game) drawGameOver(screen *engine.Screen) {
	boxW := 24
	boxH := 7
	boxX := (g.ctx.Width - boxW) / 2
	boxY := (g.ctx.Height - boxH) / 2

	screen.DrawFilledBox(boxX, boxY, boxW, boxH, engine.BoxStyleDouble,
		engine.Style{FG: engine.ColorRed},
		' ', engine.Style{BG: engine.ColorBlack})

	screen.DrawTextCentered(boxY+1, "GAME OVER", engine.Style{FG: engine.ColorBrightRed, Bold: true})
	screen.DrawTextCentered(boxY+3, fmt.Sprintf("Score: %d", g.score), engine.Style{FG: engine.ColorWhite})
	screen.DrawTextCentered(boxY+5, "[R] Restart  [Q] Quit", engine.Style{FG: engine.ColorBrightBlack})
}

func (g *Game) drawVictory(screen *engine.Screen) {
	boxW := 24
	boxH := 7
	boxX := (g.ctx.Width - boxW) / 2
	boxY := (g.ctx.Height - boxH) / 2

	screen.DrawFilledBox(boxX, boxY, boxW, boxH, engine.BoxStyleDouble,
		engine.Style{FG: engine.ColorGreen},
		' ', engine.Style{BG: engine.ColorBlack})

	screen.DrawTextCentered(boxY+1, "YOU WIN!", engine.Style{FG: engine.ColorBrightGreen, Bold: true})
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

	if g.gameOver || g.paused || g.victory {
		return nil
	}

	// Handle launch
	if input.Key == engine.KeySpace {
		g.launchBall()
		return nil
	}

	// Paddle movement
	speed := 3
	switch input.Key {
	case engine.KeyLeft:
		g.paddleX -= speed
	case engine.KeyRight:
		g.paddleX += speed
	case engine.KeyRune:
		switch input.Rune {
		case 'a', 'A':
			g.paddleX -= speed
		case 'd', 'D':
			g.paddleX += speed
		}
	}

	// Clamp paddle position
	if g.paddleX < g.areaX {
		g.paddleX = g.areaX
	}
	if g.paddleX > g.areaX+g.areaW-g.paddleW {
		g.paddleX = g.areaX + g.areaW - g.paddleW
	}

	return nil
}
