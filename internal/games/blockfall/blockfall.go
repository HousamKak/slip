package blockfall

import (
	"fmt"
	"math/rand"
	"time"

	"slip/internal/engine"
)

// Tetromino shapes (each rotation state)
var tetrominoes = map[rune][][][]int{
	'I': {
		{{1, 1, 1, 1}},
		{{1}, {1}, {1}, {1}},
	},
	'O': {
		{{1, 1}, {1, 1}},
	},
	'T': {
		{{0, 1, 0}, {1, 1, 1}},
		{{1, 0}, {1, 1}, {1, 0}},
		{{1, 1, 1}, {0, 1, 0}},
		{{0, 1}, {1, 1}, {0, 1}},
	},
	'S': {
		{{0, 1, 1}, {1, 1, 0}},
		{{1, 0}, {1, 1}, {0, 1}},
	},
	'Z': {
		{{1, 1, 0}, {0, 1, 1}},
		{{0, 1}, {1, 1}, {1, 0}},
	},
	'J': {
		{{1, 0, 0}, {1, 1, 1}},
		{{1, 1}, {1, 0}, {1, 0}},
		{{1, 1, 1}, {0, 0, 1}},
		{{0, 1}, {0, 1}, {1, 1}},
	},
	'L': {
		{{0, 0, 1}, {1, 1, 1}},
		{{1, 0}, {1, 0}, {1, 1}},
		{{1, 1, 1}, {1, 0, 0}},
		{{1, 1}, {0, 1}, {0, 1}},
	},
}

var tetrominoColors = map[rune]engine.Color{
	'I': engine.ColorCyan,
	'O': engine.ColorYellow,
	'T': engine.ColorMagenta,
	'S': engine.ColorGreen,
	'Z': engine.ColorRed,
	'J': engine.ColorBlue,
	'L': engine.ColorBrightRed, // Orange-ish
}

// Piece represents a falling tetromino
type Piece struct {
	Type     rune
	Rotation int
	X, Y     int
}

func (p *Piece) Shape() [][]int {
	shapes := tetrominoes[p.Type]
	return shapes[p.Rotation%len(shapes)]
}

func (p *Piece) Color() engine.Color {
	return tetrominoColors[p.Type]
}

// Game implements the Blockfall (Tetris-style) game
type Game struct {
	ctx *engine.GameContext

	// Board
	board     [][]rune // 0 = empty, else tetromino type
	boardW    int
	boardH    int

	// Current and next piece
	current *Piece
	next    *Piece

	// Game state
	score     int
	lines     int
	level     int
	highScore int
	gameOver  bool
	paused    bool

	// Timing
	fallTimer float64
	fallDelay float64

	// Play area position
	boardX, boardY int

	rng *rand.Rand
}

// New creates a new Blockfall game
func New() *Game {
	return &Game{
		boardW: 10,
		boardH: 20,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Info returns game metadata
func (g *Game) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "blockfall",
		Name:        "Blockfall",
		Description: "Stack falling blocks - clear lines to score!",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    30,
		MinHeight:   24,
		Category:    "game",
	}
}

// Init initializes the game
func (g *Game) Init(ctx *engine.GameContext) error {
	g.ctx = ctx

	// Position the board
	g.boardX = (ctx.Width - g.boardW*2 - 10) / 2
	if g.boardX < 2 {
		g.boardX = 2
	}
	g.boardY = 2

	// Load high score
	if ctx.GetHighScore != nil {
		g.highScore = ctx.GetHighScore()
	}

	g.reset()
	return nil
}

func (g *Game) reset() {
	// Initialize board
	g.board = make([][]rune, g.boardH)
	for i := range g.board {
		g.board[i] = make([]rune, g.boardW)
	}

	g.score = 0
	g.lines = 0
	g.level = 1
	g.gameOver = false
	g.paused = false
	g.fallDelay = 0.5

	g.next = g.randomPiece()
	g.spawnPiece()
}

func (g *Game) randomPiece() *Piece {
	types := []rune{'I', 'O', 'T', 'S', 'Z', 'J', 'L'}
	return &Piece{
		Type:     types[g.rng.Intn(len(types))],
		Rotation: 0,
		X:        g.boardW/2 - 1,
		Y:        0,
	}
}

func (g *Game) spawnPiece() {
	g.current = g.next
	g.current.X = g.boardW/2 - 1
	g.current.Y = 0
	g.next = g.randomPiece()

	// Check if spawn position is blocked
	if !g.canPlace(g.current, g.current.X, g.current.Y) {
		g.gameOver = true
		g.endGame()
	}
}

func (g *Game) endGame() {
	if g.ctx.OnScore != nil && g.score > 0 {
		if g.score > g.highScore {
			g.highScore = g.score
		}
		g.ctx.OnScore(g.score)
	}
}

func (g *Game) canPlace(piece *Piece, x, y int) bool {
	shape := piece.Shape()
	for dy, row := range shape {
		for dx, cell := range row {
			if cell == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx < 0 || nx >= g.boardW || ny >= g.boardH {
				return false
			}
			if ny >= 0 && g.board[ny][nx] != 0 {
				return false
			}
		}
	}
	return true
}

func (g *Game) placePiece() {
	shape := g.current.Shape()
	for dy, row := range shape {
		for dx, cell := range row {
			if cell == 0 {
				continue
			}
			ny := g.current.Y + dy
			nx := g.current.X + dx
			if ny >= 0 && ny < g.boardH && nx >= 0 && nx < g.boardW {
				g.board[ny][nx] = g.current.Type
			}
		}
	}
}

func (g *Game) clearLines() int {
	cleared := 0
	for y := g.boardH - 1; y >= 0; y-- {
		full := true
		for x := 0; x < g.boardW; x++ {
			if g.board[y][x] == 0 {
				full = false
				break
			}
		}
		if full {
			// Move everything above down
			for yy := y; yy > 0; yy-- {
				copy(g.board[yy], g.board[yy-1])
			}
			// Clear top row
			for x := 0; x < g.boardW; x++ {
				g.board[0][x] = 0
			}
			cleared++
			y++ // Check this row again
		}
	}
	return cleared
}

func (g *Game) move(dx int) {
	newX := g.current.X + dx
	if g.canPlace(g.current, newX, g.current.Y) {
		g.current.X = newX
	}
}

func (g *Game) rotate() {
	oldRotation := g.current.Rotation
	g.current.Rotation++

	// Wall kick - try original position, then left/right
	if g.canPlace(g.current, g.current.X, g.current.Y) {
		return
	}
	if g.canPlace(g.current, g.current.X-1, g.current.Y) {
		g.current.X--
		return
	}
	if g.canPlace(g.current, g.current.X+1, g.current.Y) {
		g.current.X++
		return
	}

	// Can't rotate
	g.current.Rotation = oldRotation
}

func (g *Game) drop() {
	for g.canPlace(g.current, g.current.X, g.current.Y+1) {
		g.current.Y++
	}
	g.landPiece()
}

func (g *Game) landPiece() {
	g.placePiece()
	cleared := g.clearLines()
	if cleared > 0 {
		g.lines += cleared
		// Scoring: more lines at once = more points
		points := []int{0, 100, 300, 500, 800}
		g.score += points[cleared] * g.level

		// Level up every 10 lines
		newLevel := g.lines/10 + 1
		if newLevel > g.level {
			g.level = newLevel
			g.fallDelay *= 0.85 // Speed up
			if g.fallDelay < 0.1 {
				g.fallDelay = 0.1
			}
		}
	}
	g.spawnPiece()
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

	g.fallTimer += dt
	if g.fallTimer >= g.fallDelay {
		g.fallTimer = 0
		if g.canPlace(g.current, g.current.X, g.current.Y+1) {
			g.current.Y++
		} else {
			g.landPiece()
		}
	}

	return nil
}

// Render draws the game
func (g *Game) Render(screen *engine.Screen) error {
	// Draw title
	title := "BLOCKFALL"
	screen.DrawText(2, 0, title, engine.Style{FG: engine.ColorBrightCyan, Bold: true})

	// Draw board border
	boardDisplayW := g.boardW*2 + 2
	boardDisplayH := g.boardH + 2
	screen.DrawBox(g.boardX-1, g.boardY-1, boardDisplayW, boardDisplayH,
		engine.BoxStyleSingle, engine.Style{FG: engine.ColorWhite})

	// Draw board contents
	for y := 0; y < g.boardH; y++ {
		for x := 0; x < g.boardW; x++ {
			cell := g.board[y][x]
			screenX := g.boardX + x*2
			screenY := g.boardY + y
			if cell != 0 {
				color := tetrominoColors[cell]
				screen.Set(screenX, screenY, '█', engine.Style{FG: color})
				screen.Set(screenX+1, screenY, '█', engine.Style{FG: color})
			}
		}
	}

	// Draw current piece
	if g.current != nil && !g.gameOver {
		shape := g.current.Shape()
		color := g.current.Color()
		for dy, row := range shape {
			for dx, cell := range row {
				if cell == 0 {
					continue
				}
				screenX := g.boardX + (g.current.X+dx)*2
				screenY := g.boardY + g.current.Y + dy
				if screenY >= g.boardY {
					screen.Set(screenX, screenY, '█', engine.Style{FG: color})
					screen.Set(screenX+1, screenY, '█', engine.Style{FG: color})
				}
			}
		}

		// Draw ghost piece (where it will land)
		ghostY := g.current.Y
		for g.canPlace(g.current, g.current.X, ghostY+1) {
			ghostY++
		}
		if ghostY > g.current.Y {
			for dy, row := range shape {
				for dx, cell := range row {
					if cell == 0 {
						continue
					}
					screenX := g.boardX + (g.current.X+dx)*2
					screenY := g.boardY + ghostY + dy
					if screenY >= g.boardY && screenY != g.boardY+g.current.Y+dy {
						screen.Set(screenX, screenY, '░', engine.Style{FG: engine.ColorBrightBlack})
						screen.Set(screenX+1, screenY, '░', engine.Style{FG: engine.ColorBrightBlack})
					}
				}
			}
		}
	}

	// Draw side panel
	panelX := g.boardX + g.boardW*2 + 3

	// Next piece
	screen.DrawText(panelX, g.boardY, "Next:", engine.Style{FG: engine.ColorWhite})
	if g.next != nil {
		shape := g.next.Shape()
		color := g.next.Color()
		for dy, row := range shape {
			for dx, cell := range row {
				if cell != 0 {
					screen.Set(panelX+dx*2, g.boardY+2+dy, '█', engine.Style{FG: color})
					screen.Set(panelX+dx*2+1, g.boardY+2+dy, '█', engine.Style{FG: color})
				}
			}
		}
	}

	// Score
	screen.DrawText(panelX, g.boardY+7, "Score:", engine.Style{FG: engine.ColorWhite})
	screen.DrawText(panelX, g.boardY+8, fmt.Sprintf("%d", g.score),
		engine.Style{FG: engine.ColorBrightYellow})

	// Level
	screen.DrawText(panelX, g.boardY+10, "Level:", engine.Style{FG: engine.ColorWhite})
	screen.DrawText(panelX, g.boardY+11, fmt.Sprintf("%d", g.level),
		engine.Style{FG: engine.ColorBrightGreen})

	// Lines
	screen.DrawText(panelX, g.boardY+13, "Lines:", engine.Style{FG: engine.ColorWhite})
	screen.DrawText(panelX, g.boardY+14, fmt.Sprintf("%d", g.lines),
		engine.Style{FG: engine.ColorBrightCyan})

	// Draw footer/controls
	controls := "[←→] Move [↑] Rotate [↓] Soft Drop [SPACE] Hard Drop [P] Pause [Q] Quit"
	if len(controls) > g.ctx.Width-4 {
		controls = "←→:Move ↑:Rot ↓:Drop SPACE:Drop P:Pause Q:Quit"
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
	boxH := 9
	boxX := (g.ctx.Width - boxW) / 2
	boxY := (g.ctx.Height - boxH) / 2

	screen.DrawFilledBox(boxX, boxY, boxW, boxH, engine.BoxStyleDouble,
		engine.Style{FG: engine.ColorRed},
		' ', engine.Style{BG: engine.ColorBlack})

	screen.DrawTextCentered(boxY+1, "GAME OVER", engine.Style{FG: engine.ColorBrightRed, Bold: true})
	screen.DrawTextCentered(boxY+3, fmt.Sprintf("Score: %d", g.score), engine.Style{FG: engine.ColorWhite})
	screen.DrawTextCentered(boxY+4, fmt.Sprintf("Level: %d", g.level), engine.Style{FG: engine.ColorWhite})
	screen.DrawTextCentered(boxY+5, fmt.Sprintf("Lines: %d", g.lines), engine.Style{FG: engine.ColorWhite})
	screen.DrawTextCentered(boxY+7, "[R] Restart  [Q] Quit", engine.Style{FG: engine.ColorBrightBlack})
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

	switch input.Key {
	case engine.KeyLeft:
		g.move(-1)
	case engine.KeyRight:
		g.move(1)
	case engine.KeyUp:
		g.rotate()
	case engine.KeyDown:
		// Soft drop
		if g.canPlace(g.current, g.current.X, g.current.Y+1) {
			g.current.Y++
			g.score++ // Bonus for soft drop
		}
	case engine.KeySpace:
		g.drop()
	case engine.KeyRune:
		switch input.Rune {
		case 'a', 'A':
			g.move(-1)
		case 'd', 'D':
			g.move(1)
		case 'w', 'W':
			g.rotate()
		case 's', 'S':
			if g.canPlace(g.current, g.current.X, g.current.Y+1) {
				g.current.Y++
				g.score++
			}
		}
	}

	return nil
}
