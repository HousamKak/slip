package pong

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"slip/internal/engine"
)

// Game implements the Pong game
type Game struct {
	ctx *engine.GameContext

	// Ball
	ballX, ballY   float64
	ballVX, ballVY float64

	// Paddles
	playerY float64
	cpuY    float64
	paddleH int

	// Scores
	playerScore int
	cpuScore    int
	winScore    int
	highScore   int

	// Game state
	gameOver   bool
	paused     bool
	winner     string
	countdown  int
	countdownT float64

	// Play area
	areaX, areaY int
	areaW, areaH int

	rng *rand.Rand
}

// New creates a new Pong game
func New() *Game {
	return &Game{
		winScore: 5,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Info returns game metadata
func (g *Game) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "pong",
		Name:        "Pong",
		Description: "Classic Pong - beat the CPU to 5 points!",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    30,
		MinHeight:   15,
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

	g.paddleH = g.areaH / 4
	if g.paddleH < 3 {
		g.paddleH = 3
	}

	// Load high score
	if ctx.GetHighScore != nil {
		g.highScore = ctx.GetHighScore()
	}

	g.reset()
	return nil
}

func (g *Game) reset() {
	g.playerScore = 0
	g.cpuScore = 0
	g.gameOver = false
	g.winner = ""
	g.resetBall()
	g.resetPaddles()
}

func (g *Game) resetBall() {
	g.ballX = float64(g.areaX + g.areaW/2)
	g.ballY = float64(g.areaY + g.areaH/2)

	// Random direction
	angle := (g.rng.Float64()*0.5 - 0.25) * math.Pi // -45 to +45 degrees
	speed := 15.0
	if g.rng.Intn(2) == 0 {
		g.ballVX = -speed * math.Cos(angle)
	} else {
		g.ballVX = speed * math.Cos(angle)
	}
	g.ballVY = speed * math.Sin(angle)

	// Start countdown
	g.countdown = 3
	g.countdownT = 0
}

func (g *Game) resetPaddles() {
	g.playerY = float64(g.areaY + g.areaH/2 - g.paddleH/2)
	g.cpuY = float64(g.areaY + g.areaH/2 - g.paddleH/2)
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

	// Handle countdown
	if g.countdown > 0 {
		g.countdownT += dt
		if g.countdownT >= 1.0 {
			g.countdown--
			g.countdownT = 0
		}
		return nil
	}

	// Move ball
	g.ballX += g.ballVX * dt
	g.ballY += g.ballVY * dt

	// Ball collision with top/bottom walls
	if g.ballY <= float64(g.areaY) {
		g.ballY = float64(g.areaY)
		g.ballVY = -g.ballVY
	}
	if g.ballY >= float64(g.areaY+g.areaH-1) {
		g.ballY = float64(g.areaY + g.areaH - 1)
		g.ballVY = -g.ballVY
	}

	// Ball collision with player paddle (left side)
	paddleX := float64(g.areaX + 2)
	if g.ballX <= paddleX+1 && g.ballX >= paddleX &&
		g.ballY >= g.playerY && g.ballY < g.playerY+float64(g.paddleH) {
		g.ballX = paddleX + 1
		g.ballVX = math.Abs(g.ballVX) * 1.05 // Speed up slightly
		// Add spin based on where it hit the paddle
		hitPos := (g.ballY - g.playerY) / float64(g.paddleH)
		g.ballVY += (hitPos - 0.5) * 10
	}

	// Ball collision with CPU paddle (right side)
	cpuPaddleX := float64(g.areaX + g.areaW - 3)
	if g.ballX >= cpuPaddleX-1 && g.ballX <= cpuPaddleX &&
		g.ballY >= g.cpuY && g.ballY < g.cpuY+float64(g.paddleH) {
		g.ballX = cpuPaddleX - 1
		g.ballVX = -math.Abs(g.ballVX) * 1.05
		hitPos := (g.ballY - g.cpuY) / float64(g.paddleH)
		g.ballVY += (hitPos - 0.5) * 10
	}

	// Ball out of bounds
	if g.ballX < float64(g.areaX) {
		// CPU scores
		g.cpuScore++
		if g.cpuScore >= g.winScore {
			g.gameOver = true
			g.winner = "CPU"
			g.endGame()
		} else {
			g.resetBall()
		}
	}
	if g.ballX > float64(g.areaX+g.areaW) {
		// Player scores
		g.playerScore++
		if g.playerScore >= g.winScore {
			g.gameOver = true
			g.winner = "Player"
			g.endGame()
		} else {
			g.resetBall()
		}
	}

	// CPU AI - follow the ball with some delay
	cpuSpeed := 20.0 * dt
	cpuCenter := g.cpuY + float64(g.paddleH)/2
	if g.ballVX > 0 { // Only move when ball is coming towards CPU
		if cpuCenter < g.ballY-1 {
			g.cpuY += cpuSpeed
		} else if cpuCenter > g.ballY+1 {
			g.cpuY -= cpuSpeed
		}
	}

	// Clamp paddle positions
	g.playerY = clamp(g.playerY, float64(g.areaY), float64(g.areaY+g.areaH-g.paddleH))
	g.cpuY = clamp(g.cpuY, float64(g.areaY), float64(g.areaY+g.areaH-g.paddleH))

	return nil
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func (g *Game) endGame() {
	// Only save score if player won
	if g.winner == "Player" && g.ctx.OnScore != nil {
		score := g.playerScore * 100
		if score > g.highScore {
			g.highScore = score
		}
		g.ctx.OnScore(score)
	}
}

// Render draws the game
func (g *Game) Render(screen *engine.Screen) error {
	// Draw border
	screen.DrawBox(0, 1, g.ctx.Width, g.ctx.Height-2, engine.BoxStyleSingle,
		engine.Style{FG: engine.ColorCyan})

	// Draw header
	title := "PONG"
	screen.DrawText(2, 0, title, engine.Style{FG: engine.ColorBrightCyan, Bold: true})

	scoreText := fmt.Sprintf("%d : %d", g.playerScore, g.cpuScore)
	screen.DrawTextCentered(0, scoreText, engine.Style{FG: engine.ColorBrightWhite, Bold: true})

	targetText := fmt.Sprintf("First to %d", g.winScore)
	screen.DrawText(g.ctx.Width-len(targetText)-2, 0, targetText,
		engine.Style{FG: engine.ColorBrightBlack})

	// Draw center line
	for y := g.areaY; y < g.areaY+g.areaH; y++ {
		if y%2 == 0 {
			screen.Set(g.areaX+g.areaW/2, y, '│', engine.Style{FG: engine.ColorBrightBlack})
		}
	}

	// Draw player paddle
	for i := 0; i < g.paddleH; i++ {
		screen.Set(g.areaX+2, int(g.playerY)+i, '█', engine.Style{FG: engine.ColorGreen})
	}

	// Draw CPU paddle
	for i := 0; i < g.paddleH; i++ {
		screen.Set(g.areaX+g.areaW-3, int(g.cpuY)+i, '█', engine.Style{FG: engine.ColorRed})
	}

	// Draw ball
	if g.countdown == 0 && !g.gameOver {
		screen.Set(int(g.ballX), int(g.ballY), '●', engine.Style{FG: engine.ColorBrightWhite, Bold: true})
	}

	// Draw countdown
	if g.countdown > 0 {
		countStr := fmt.Sprintf("%d", g.countdown)
		screen.DrawTextCentered(g.ctx.Height/2, countStr,
			engine.Style{FG: engine.ColorBrightYellow, Bold: true})
	}

	// Draw footer/controls
	controls := "[W/S] or [↑/↓] Move  [P] Pause  [R] Restart  [Q] Quit"
	if len(controls) > g.ctx.Width-4 {
		controls = "W/S:Move P:Pause R:Restart Q:Quit"
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

	var borderColor engine.Color
	if g.winner == "Player" {
		borderColor = engine.ColorGreen
	} else {
		borderColor = engine.ColorRed
	}

	screen.DrawFilledBox(boxX, boxY, boxW, boxH, engine.BoxStyleDouble,
		engine.Style{FG: borderColor},
		' ', engine.Style{BG: engine.ColorBlack})

	screen.DrawTextCentered(boxY+1, g.winner+" WINS!", engine.Style{FG: borderColor, Bold: true})
	screen.DrawTextCentered(boxY+3, fmt.Sprintf("%d - %d", g.playerScore, g.cpuScore),
		engine.Style{FG: engine.ColorWhite})
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

	if g.gameOver || g.paused {
		return nil
	}

	// Paddle movement
	speed := 2.0
	switch input.Key {
	case engine.KeyUp:
		g.playerY -= speed
	case engine.KeyDown:
		g.playerY += speed
	case engine.KeyRune:
		switch input.Rune {
		case 'w', 'W':
			g.playerY -= speed
		case 's', 'S':
			g.playerY += speed
		}
	}

	// Clamp paddle position
	g.playerY = clamp(g.playerY, float64(g.areaY), float64(g.areaY+g.areaH-g.paddleH))

	return nil
}
