package blockfall

import (
	"bytes"
	"testing"

	"slip/internal/engine"
)

func TestBlockfallNew(t *testing.T) {
	game := New()
	if game == nil {
		t.Fatal("New() returned nil")
	}
	if game.boardW != 10 {
		t.Errorf("Expected boardW 10, got %d", game.boardW)
	}
	if game.boardH != 20 {
		t.Errorf("Expected boardH 20, got %d", game.boardH)
	}
	if game.rng == nil {
		t.Error("Expected rng to be initialized")
	}
}

func TestBlockfallInfo(t *testing.T) {
	game := New()
	info := game.Info()

	if info.ID != "blockfall" {
		t.Errorf("Expected ID 'blockfall', got '%s'", info.ID)
	}
	if info.Name != "Blockfall" {
		t.Errorf("Expected Name 'Blockfall', got '%s'", info.Name)
	}
	if info.MinWidth != 30 {
		t.Errorf("Expected MinWidth 30, got %d", info.MinWidth)
	}
	if info.MinHeight != 24 {
		t.Errorf("Expected MinHeight 24, got %d", info.MinHeight)
	}
}

func TestBlockfallInit(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{
		Width:  40,
		Height: 26,
	}

	err := game.Init(ctx)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Score should be 0
	if game.score != 0 {
		t.Errorf("Expected score 0, got %d", game.score)
	}

	// Lines should be 0
	if game.lines != 0 {
		t.Errorf("Expected lines 0, got %d", game.lines)
	}

	// Level should be 1
	if game.level != 1 {
		t.Errorf("Expected level 1, got %d", game.level)
	}

	// Should have a current piece
	if game.current == nil {
		t.Error("Expected current piece to be set")
	}

	// Should have a next piece
	if game.next == nil {
		t.Error("Expected next piece to be set")
	}

	// Board should be initialized
	if len(game.board) != game.boardH {
		t.Errorf("Expected board height %d, got %d", game.boardH, len(game.board))
	}
	if len(game.board[0]) != game.boardW {
		t.Errorf("Expected board width %d, got %d", game.boardW, len(game.board[0]))
	}
}

func TestBlockfallPieceShape(t *testing.T) {
	piece := &Piece{Type: 'I', Rotation: 0}
	shape := piece.Shape()

	if len(shape) == 0 {
		t.Error("Expected shape to have rows")
	}

	// I piece horizontal should be 1 row, 4 columns
	if len(shape) != 1 || len(shape[0]) != 4 {
		t.Errorf("I piece rotation 0 should be 1x4, got %dx%d", len(shape), len(shape[0]))
	}

	// Rotate
	piece.Rotation = 1
	shape = piece.Shape()

	// I piece vertical should be 4 rows, 1 column
	if len(shape) != 4 || len(shape[0]) != 1 {
		t.Errorf("I piece rotation 1 should be 4x1, got %dx%d", len(shape), len(shape[0]))
	}
}

func TestBlockfallPieceColor(t *testing.T) {
	tests := []struct {
		pieceType rune
		expected  engine.Color
	}{
		{'I', engine.ColorCyan},
		{'O', engine.ColorYellow},
		{'T', engine.ColorMagenta},
		{'S', engine.ColorGreen},
		{'Z', engine.ColorRed},
		{'J', engine.ColorBlue},
		{'L', engine.ColorBrightRed},
	}

	for _, tt := range tests {
		piece := &Piece{Type: tt.pieceType}
		if piece.Color() != tt.expected {
			t.Errorf("Piece %c: expected color %v, got %v", tt.pieceType, tt.expected, piece.Color())
		}
	}
}

func TestBlockfallMove(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	initialX := game.current.X

	// Move left
	game.move(-1)
	if game.current.X != initialX-1 {
		t.Errorf("Expected X %d after moving left, got %d", initialX-1, game.current.X)
	}

	// Move right
	game.move(1)
	if game.current.X != initialX {
		t.Errorf("Expected X %d after moving right, got %d", initialX, game.current.X)
	}
}

func TestBlockfallMoveBlocked(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	// Move to left edge
	game.current.X = 0

	// Try to move left (should be blocked)
	game.move(-1)

	if game.current.X < 0 {
		t.Error("Piece should not move past left edge")
	}
}

func TestBlockfallRotate(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	// Use T piece for rotation test
	game.current = &Piece{Type: 'T', Rotation: 0, X: 5, Y: 5}

	initialRotation := game.current.Rotation
	game.rotate()

	if game.current.Rotation == initialRotation {
		t.Error("Piece should rotate")
	}
}

func TestBlockfallCanPlace(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	piece := &Piece{Type: 'O', Rotation: 0, X: 5, Y: 5}

	// Should be able to place in empty board
	if !game.canPlace(piece, 5, 5) {
		t.Error("Should be able to place piece in empty area")
	}

	// Fill a cell where piece would go
	game.board[5][5] = 'X'

	// Should not be able to place now
	if game.canPlace(piece, 5, 5) {
		t.Error("Should not be able to place piece on occupied cell")
	}
}

func TestBlockfallCanPlaceOutOfBounds(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	piece := &Piece{Type: 'I', Rotation: 0} // Horizontal I is 4 wide

	// Try to place too far right
	if game.canPlace(piece, game.boardW-1, 5) {
		t.Error("Should not be able to place piece out of bounds")
	}

	// Try to place too far left
	if game.canPlace(piece, -1, 5) {
		t.Error("Should not be able to place piece out of bounds (left)")
	}
}

func TestBlockfallDrop(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.current = &Piece{Type: 'O', Rotation: 0, X: 5, Y: 0}

	// Track initial board state
	initialEmptyCount := 0
	for y := 0; y < game.boardH; y++ {
		for x := 0; x < game.boardW; x++ {
			if game.board[y][x] == 0 {
				initialEmptyCount++
			}
		}
	}

	game.drop()

	// After drop, the board should have fewer empty cells (piece was placed)
	finalEmptyCount := 0
	for y := 0; y < game.boardH; y++ {
		for x := 0; x < game.boardW; x++ {
			if game.board[y][x] == 0 {
				finalEmptyCount++
			}
		}
	}

	if finalEmptyCount >= initialEmptyCount {
		t.Error("After drop, piece should be placed on board")
	}
}

func TestBlockfallClearLines(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	// Fill bottom row
	for x := 0; x < game.boardW; x++ {
		game.board[game.boardH-1][x] = 'X'
	}

	cleared := game.clearLines()

	if cleared != 1 {
		t.Errorf("Expected 1 line cleared, got %d", cleared)
	}

	// Bottom row should now be empty
	for x := 0; x < game.boardW; x++ {
		if game.board[game.boardH-1][x] != 0 {
			t.Error("Bottom row should be cleared")
			break
		}
	}
}

func TestBlockfallClearMultipleLines(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	// Fill bottom 4 rows
	for y := game.boardH - 4; y < game.boardH; y++ {
		for x := 0; x < game.boardW; x++ {
			game.board[y][x] = 'X'
		}
	}

	cleared := game.clearLines()

	if cleared != 4 {
		t.Errorf("Expected 4 lines cleared, got %d", cleared)
	}
}

func TestBlockfallScoring(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.level = 1
	initialScore := game.score

	// Fill bottom row minus one cell
	for x := 0; x < game.boardW-1; x++ {
		game.board[game.boardH-1][x] = 'X'
	}

	// Set current piece to fill the gap and complete the line
	game.current = &Piece{Type: 'O', Rotation: 0, X: game.boardW - 2, Y: game.boardH - 2}

	// Land the piece (this triggers placePiece and clearLines with scoring)
	game.landPiece()

	// Score should have increased
	if game.score <= initialScore {
		t.Errorf("Expected score to increase, got %d (was %d)", game.score, initialScore)
	}
}

func TestBlockfallLevelUp(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.lines = 9
	initialLevel := game.level
	initialDelay := game.fallDelay

	// Clear a line to trigger level up
	for x := 0; x < game.boardW; x++ {
		game.board[game.boardH-1][x] = 'X'
	}

	game.landPiece()

	if game.level != initialLevel+1 {
		t.Errorf("Expected level %d, got %d", initialLevel+1, game.level)
	}

	if game.fallDelay >= initialDelay {
		t.Error("Fall delay should decrease on level up")
	}
}

func TestBlockfallGameOver(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	// Fill top rows to trigger game over on spawn
	for x := 0; x < game.boardW; x++ {
		game.board[0][x] = 'X'
		game.board[1][x] = 'X'
	}

	game.spawnPiece()

	if !game.gameOver {
		t.Error("Game should be over when spawn is blocked")
	}
}

func TestBlockfallInputMovement(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	initialX := game.current.X

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyLeft})
	if game.current.X >= initialX {
		t.Error("Left arrow should move piece left")
	}

	game.current.X = initialX
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRight})
	if game.current.X <= initialX {
		t.Error("Right arrow should move piece right")
	}
}

func TestBlockfallInputRotate(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.current = &Piece{Type: 'T', Rotation: 0, X: 5, Y: 5}
	initialRotation := game.current.Rotation

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyUp})

	if game.current.Rotation == initialRotation {
		t.Error("Up arrow should rotate piece")
	}
}

func TestBlockfallInputSoftDrop(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.current = &Piece{Type: 'O', Rotation: 0, X: 5, Y: 5}
	initialY := game.current.Y
	initialScore := game.score

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyDown})

	if game.current.Y <= initialY {
		t.Error("Down arrow should soft drop piece")
	}
	if game.score <= initialScore {
		t.Error("Soft drop should give bonus points")
	}
}

func TestBlockfallInputHardDrop(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.current = &Piece{Type: 'O', Rotation: 0, X: 5, Y: 0}

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeySpace})

	// After hard drop, a new piece should be spawned
	// The current piece's Y should be reset to top
	if game.current.Y > 2 {
		t.Error("After hard drop, new piece should spawn at top")
	}
}

func TestBlockfallPause(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'p'})

	if !game.paused {
		t.Error("Game should be paused")
	}

	// Piece should not fall while paused
	initialY := game.current.Y
	game.fallTimer = game.fallDelay
	game.Update(0.1)

	if game.current.Y != initialY {
		t.Error("Piece should not fall while paused")
	}
}

func TestBlockfallRestart(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.score = 1000
	game.lines = 20
	game.level = 5
	game.gameOver = true

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'r'})

	if game.score != 0 {
		t.Error("Score should be reset")
	}
	if game.lines != 0 {
		t.Error("Lines should be reset")
	}
	if game.level != 1 {
		t.Error("Level should be reset")
	}
	if game.gameOver {
		t.Error("Game over should be reset")
	}
}

func TestBlockfallWASDControls(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	initialX := game.current.X

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'a'})
	if game.current.X >= initialX {
		t.Error("'a' should move piece left")
	}

	game.current.X = initialX
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'd'})
	if game.current.X <= initialX {
		t.Error("'d' should move piece right")
	}

	game.current = &Piece{Type: 'T', Rotation: 0, X: 5, Y: 5}
	initialRotation := game.current.Rotation
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'w'})
	if game.current.Rotation == initialRotation {
		t.Error("'w' should rotate piece")
	}
}

func TestBlockfallFallTimer(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	game.current = &Piece{Type: 'O', Rotation: 0, X: 5, Y: 5}
	initialY := game.current.Y

	// Update without enough time should not move piece
	game.Update(game.fallDelay / 2)
	if game.current.Y != initialY {
		t.Error("Piece should not fall before fallDelay")
	}

	// Update with enough time should move piece
	game.Update(game.fallDelay)
	if game.current.Y <= initialY {
		t.Error("Piece should fall after fallDelay")
	}
}

func TestBlockfallScoreCallback(t *testing.T) {
	game := New()
	var reportedScore int
	ctx := &engine.GameContext{
		Width:  40,
		Height: 26,
		OnScore: func(score int) {
			reportedScore = score
		},
	}
	game.Init(ctx)

	game.score = 5000

	// Trigger game over
	for x := 0; x < game.boardW; x++ {
		game.board[0][x] = 'X'
		game.board[1][x] = 'X'
	}
	game.spawnPiece()

	if reportedScore != 5000 {
		t.Errorf("Expected reported score 5000, got %d", reportedScore)
	}
}

func TestBlockfallRender(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 26}
	game.Init(ctx)

	var buf bytes.Buffer
	screen := engine.NewScreen(40, 26, &buf)

	err := game.Render(screen)
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}
}
