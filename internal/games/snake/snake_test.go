package snake

import (
	"bytes"
	"testing"

	"slip/internal/engine"
)

func TestSnakeNew(t *testing.T) {
	game := New()
	if game == nil {
		t.Fatal("New() returned nil")
	}
	if game.moveDelay != 0.12 {
		t.Errorf("Expected moveDelay 0.12, got %f", game.moveDelay)
	}
	if game.rng == nil {
		t.Error("Expected rng to be initialized")
	}
}

func TestSnakeInfo(t *testing.T) {
	game := New()
	info := game.Info()

	if info.ID != "snake" {
		t.Errorf("Expected ID 'snake', got '%s'", info.ID)
	}
	if info.Name != "Snake" {
		t.Errorf("Expected Name 'Snake', got '%s'", info.Name)
	}
	if info.MinWidth != 20 {
		t.Errorf("Expected MinWidth 20, got %d", info.MinWidth)
	}
	if info.MinHeight != 12 {
		t.Errorf("Expected MinHeight 12, got %d", info.MinHeight)
	}
	if info.Category != "game" {
		t.Errorf("Expected Category 'game', got '%s'", info.Category)
	}
}

func TestSnakeInit(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{
		Width:  40,
		Height: 20,
	}

	err := game.Init(ctx)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Snake should start with length 3
	if len(game.snake) != 3 {
		t.Errorf("Expected snake length 3, got %d", len(game.snake))
	}

	// Direction should be right
	if game.direction != engine.DirRight {
		t.Errorf("Expected direction DirRight, got %v", game.direction)
	}

	// Game should not be over
	if game.gameOver {
		t.Error("Game should not be over at start")
	}

	// Score should be 0
	if game.score != 0 {
		t.Errorf("Expected score 0, got %d", game.score)
	}
}

func TestSnakeInitWithHighScore(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{
		Width:  40,
		Height: 20,
		GetHighScore: func() int {
			return 500
		},
	}

	err := game.Init(ctx)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if game.highScore != 500 {
		t.Errorf("Expected highScore 500, got %d", game.highScore)
	}
}

func TestSnakeMovement(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	startHead := game.snake[0]

	// Force update by setting moveTimer
	game.moveTimer = game.moveDelay
	game.Update(0.01)

	newHead := game.snake[0]

	// Snake should move right (default direction)
	if newHead.X != startHead.X+1 {
		t.Errorf("Expected head X to increase by 1, got %d -> %d", startHead.X, newHead.X)
	}
	if newHead.Y != startHead.Y {
		t.Errorf("Expected head Y to stay same, got %d -> %d", startHead.Y, newHead.Y)
	}
}

func TestSnakeDirectionChange(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Change direction to up
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyUp})

	// Direction should be queued
	if game.nextDir != engine.DirUp {
		t.Errorf("Expected nextDir DirUp, got %v", game.nextDir)
	}
}

func TestSnakePrevent180Turn(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Default direction is right, trying to go left should be prevented
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyLeft})

	// nextDir should still be right (opposite direction blocked)
	if game.nextDir != engine.DirRight {
		t.Errorf("180-degree turn should be prevented, nextDir is %v", game.nextDir)
	}
}

func TestSnakeWallCollision(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Move snake head to left edge
	game.snake[0] = engine.Point{X: game.areaX, Y: game.areaY + 5}
	game.direction = engine.DirLeft
	game.nextDir = engine.DirLeft

	// Trigger update
	game.moveTimer = game.moveDelay
	game.Update(0.01)

	if !game.gameOver {
		t.Error("Expected game over on wall collision")
	}
}

func TestSnakeSelfCollision(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Create a snake that will collide with itself
	// Snake going right, but we put the head in a position where it will hit the body
	game.snake = []engine.Point{
		{X: 10, Y: 10}, // head
		{X: 11, Y: 10}, // body
		{X: 12, Y: 10}, // body
		{X: 12, Y: 11}, // body
		{X: 11, Y: 11}, // body
		{X: 10, Y: 11}, // body
		{X: 9, Y: 11},  // body
		{X: 9, Y: 10},  // body - this is where head will go if we move left
	}
	game.direction = engine.DirLeft
	game.nextDir = engine.DirLeft

	// Trigger update
	game.moveTimer = game.moveDelay
	game.Update(0.01)

	if !game.gameOver {
		t.Error("Expected game over on self collision")
	}
}

func TestSnakeFoodCollision(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Place food directly in front of snake head
	head := game.snake[0]
	game.food = engine.Point{X: head.X + 1, Y: head.Y}

	initialLen := len(game.snake)
	initialScore := game.score

	// Trigger update
	game.moveTimer = game.moveDelay
	game.Update(0.01)

	// Snake should grow
	if len(game.snake) != initialLen+1 {
		t.Errorf("Snake should grow after eating, expected %d, got %d", initialLen+1, len(game.snake))
	}

	// Score should increase
	if game.score != initialScore+10 {
		t.Errorf("Score should increase by 10, expected %d, got %d", initialScore+10, game.score)
	}
}

func TestSnakePause(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Pause the game
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'p'})

	if !game.paused {
		t.Error("Game should be paused after pressing 'p'")
	}

	// Store head position
	headBefore := game.snake[0]

	// Try to update while paused
	game.moveTimer = game.moveDelay
	game.Update(0.01)

	// Snake should not move
	if !game.snake[0].Equals(headBefore) {
		t.Error("Snake should not move while paused")
	}

	// Unpause
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'P'})

	if game.paused {
		t.Error("Game should be unpaused after pressing 'P' again")
	}
}

func TestSnakeRestart(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Change game state
	game.score = 100
	game.gameOver = true

	// Restart
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'r'})

	if game.score != 0 {
		t.Errorf("Score should be reset to 0, got %d", game.score)
	}
	if game.gameOver {
		t.Error("Game over should be reset")
	}
	if len(game.snake) != 3 {
		t.Errorf("Snake should be reset to length 3, got %d", len(game.snake))
	}
}

func TestSnakeWASDControls(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	tests := []struct {
		rune     rune
		expected engine.Direction
	}{
		{'w', engine.DirUp},
		{'W', engine.DirUp},
		{'s', engine.DirDown},
		{'S', engine.DirDown},
		{'a', engine.DirLeft},
		{'d', engine.DirRight},
	}

	for _, tt := range tests {
		game.direction = engine.DirNone // Reset to allow any direction
		game.nextDir = engine.DirNone

		game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: tt.rune})

		if game.nextDir != tt.expected {
			t.Errorf("Key '%c': expected direction %v, got %v", tt.rune, tt.expected, game.nextDir)
		}
	}
}

func TestSnakeSpeedIncrease(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	initialDelay := game.moveDelay

	// Place food in front and eat it
	head := game.snake[0]
	game.food = engine.Point{X: head.X + 1, Y: head.Y}

	game.moveTimer = game.moveDelay
	game.Update(0.01)

	if game.moveDelay >= initialDelay {
		t.Error("Move delay should decrease after eating food (game should speed up)")
	}
}

func TestSnakeScoreCallback(t *testing.T) {
	game := New()
	var reportedScore int
	ctx := &engine.GameContext{
		Width:  40,
		Height: 20,
		OnScore: func(score int) {
			reportedScore = score
		},
	}
	game.Init(ctx)

	// Eat some food to get a score
	head := game.snake[0]
	game.food = engine.Point{X: head.X + 1, Y: head.Y}
	game.moveTimer = game.moveDelay
	game.Update(0.01)

	// Then trigger game over
	game.snake[0] = engine.Point{X: game.areaX, Y: game.areaY + 5}
	game.direction = engine.DirLeft
	game.nextDir = engine.DirLeft
	game.moveTimer = game.moveDelay
	game.Update(0.01)

	if reportedScore != 10 {
		t.Errorf("Expected reported score 10, got %d", reportedScore)
	}
}

func TestSnakeRender(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Create a mock screen
	var buf bytes.Buffer
	screen := engine.NewScreen(40, 20, &buf)

	// Should not panic
	err := game.Render(screen)
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}
}
