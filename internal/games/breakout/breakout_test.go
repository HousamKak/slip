package breakout

import (
	"bytes"
	"testing"

	"slip/internal/engine"
)

func TestBreakoutNew(t *testing.T) {
	game := New()
	if game == nil {
		t.Fatal("New() returned nil")
	}
	if game.rng == nil {
		t.Error("Expected rng to be initialized")
	}
}

func TestBreakoutInfo(t *testing.T) {
	game := New()
	info := game.Info()

	if info.ID != "breakout" {
		t.Errorf("Expected ID 'breakout', got '%s'", info.ID)
	}
	if info.Name != "Breakout" {
		t.Errorf("Expected Name 'Breakout', got '%s'", info.Name)
	}
	if info.MinWidth != 30 {
		t.Errorf("Expected MinWidth 30, got %d", info.MinWidth)
	}
	if info.MinHeight != 18 {
		t.Errorf("Expected MinHeight 18, got %d", info.MinHeight)
	}
}

func TestBreakoutInit(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{
		Width:  40,
		Height: 24,
	}

	err := game.Init(ctx)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Lives should be 3
	if game.lives != 3 {
		t.Errorf("Expected lives 3, got %d", game.lives)
	}

	// Score should be 0
	if game.score != 0 {
		t.Errorf("Expected score 0, got %d", game.score)
	}

	// Level should be 1
	if game.level != 1 {
		t.Errorf("Expected level 1, got %d", game.level)
	}

	// Should have bricks
	if len(game.bricks) == 0 {
		t.Error("Expected bricks to be created")
	}

	// Ball should not be launched
	if game.ballLaunched {
		t.Error("Ball should not be launched at start")
	}
}

func TestBreakoutBallFollowsPaddleBeforeLaunch(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	// Move paddle
	game.paddleX = 15

	game.Update(0.1)

	// Ball should follow paddle
	expectedBallX := float64(game.paddleX + game.paddleW/2)
	if game.ballX != expectedBallX {
		t.Errorf("Ball X should follow paddle center, expected %f, got %f", expectedBallX, game.ballX)
	}
}

func TestBreakoutLaunchBall(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeySpace})

	if !game.ballLaunched {
		t.Error("Ball should be launched after space press")
	}
	if game.ballVX == 0 && game.ballVY == 0 {
		t.Error("Ball should have velocity after launch")
	}
}

func TestBreakoutPaddleMovement(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	initialX := game.paddleX

	// Move left
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyLeft})
	if game.paddleX >= initialX {
		t.Error("Paddle should move left")
	}

	// Reset and move right
	game.paddleX = initialX
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRight})
	if game.paddleX <= initialX {
		t.Error("Paddle should move right")
	}
}

func TestBreakoutPaddleClamp(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	// Try to move paddle past left edge
	game.paddleX = game.areaX - 10
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyLeft})

	if game.paddleX < game.areaX {
		t.Error("Paddle should be clamped to left edge")
	}

	// Try to move paddle past right edge
	game.paddleX = game.areaX + game.areaW
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRight})

	if game.paddleX > game.areaX+game.areaW-game.paddleW {
		t.Error("Paddle should be clamped to right edge")
	}
}

func TestBreakoutBallWallBounce(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.ballLaunched = true

	// Test left wall bounce
	game.ballX = float64(game.areaX) - 0.5
	game.ballVX = -10.0
	game.Update(0.01)

	if game.ballVX < 0 {
		t.Error("Ball should bounce off left wall (VX positive)")
	}

	// Test right wall bounce
	game.ballX = float64(game.areaX+game.areaW) + 0.5
	game.ballVX = 10.0
	game.Update(0.01)

	if game.ballVX > 0 {
		t.Error("Ball should bounce off right wall (VX negative)")
	}

	// Test top wall bounce
	game.ballY = float64(game.areaY) - 0.5
	game.ballVY = -10.0
	game.Update(0.01)

	if game.ballVY < 0 {
		t.Error("Ball should bounce off top wall (VY positive)")
	}
}

func TestBreakoutBallFallsLoseLife(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.ballLaunched = true
	initialLives := game.lives

	// Ball falls below paddle
	game.ballY = float64(game.areaY + game.areaH + 1)
	game.ballVY = 10.0

	game.Update(0.01)

	if game.lives != initialLives-1 {
		t.Errorf("Expected lives %d, got %d", initialLives-1, game.lives)
	}
	if game.ballLaunched {
		t.Error("Ball should be reset (not launched)")
	}
}

func TestBreakoutGameOverNoLives(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.ballLaunched = true
	game.lives = 1

	// Ball falls
	game.ballY = float64(game.areaY + game.areaH + 1)
	game.Update(0.01)

	if !game.gameOver {
		t.Error("Game should be over when no lives left")
	}
}

func TestBreakoutBrickCollision(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.ballLaunched = true
	initialBricks := len(game.bricks)
	initialScore := game.score

	if initialBricks == 0 {
		t.Fatal("No bricks to test collision with")
	}

	// Position ball at first brick (save initial hits before update)
	brick := game.bricks[0]
	initialHits := brick.Hits
	game.ballX = float64(brick.X + brick.Width/2)
	game.ballY = float64(brick.Y)
	game.ballVY = 1.0 // Moving down into brick

	game.Update(0.01)

	// Either brick is destroyed or damaged
	brickDestroyed := len(game.bricks) < initialBricks
	brickDamaged := len(game.bricks) == initialBricks && game.bricks[0].Hits < initialHits

	if !brickDestroyed && !brickDamaged {
		t.Error("Brick should be damaged or destroyed on collision")
	}

	// Score should increase if brick destroyed
	if brickDestroyed && game.score <= initialScore {
		t.Error("Score should increase when brick is destroyed")
	}
}

func TestBreakoutLevelComplete(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.ballLaunched = true
	initialLevel := game.level

	// Remove all bricks
	game.bricks = nil

	game.Update(0.01)

	if game.level != initialLevel+1 {
		t.Errorf("Expected level %d, got %d", initialLevel+1, game.level)
	}
}

func TestBreakoutVictory(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.ballLaunched = true
	game.level = 5 // Max level

	// Remove all bricks
	game.bricks = nil

	game.Update(0.01)

	if !game.victory {
		t.Error("Should have victory after completing all levels")
	}
}

func TestBreakoutPause(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'p'})

	if !game.paused {
		t.Error("Game should be paused")
	}

	// Ball should not move while paused
	game.ballLaunched = true
	ballX := game.ballX
	game.ballVX = 10.0
	game.Update(0.1)

	if game.ballX != ballX {
		t.Error("Ball should not move while paused")
	}
}

func TestBreakoutRestart(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	game.score = 500
	game.lives = 1
	game.level = 3
	game.gameOver = true

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'r'})

	if game.score != 0 {
		t.Error("Score should be reset")
	}
	if game.lives != 3 {
		t.Error("Lives should be reset to 3")
	}
	if game.level != 1 {
		t.Error("Level should be reset to 1")
	}
	if game.gameOver {
		t.Error("Game over should be reset")
	}
}

func TestBreakoutADControls(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	initialX := game.paddleX

	// Move left with A
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'a'})
	if game.paddleX >= initialX {
		t.Error("Paddle should move left with 'a'")
	}

	// Reset and move right with D
	game.paddleX = initialX
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'd'})
	if game.paddleX <= initialX {
		t.Error("Paddle should move right with 'd'")
	}
}

func TestBreakoutScoreCallback(t *testing.T) {
	game := New()
	var reportedScore int
	ctx := &engine.GameContext{
		Width:  40,
		Height: 24,
		OnScore: func(score int) {
			reportedScore = score
		},
	}
	game.Init(ctx)

	game.score = 1000
	game.lives = 1
	game.ballLaunched = true

	// Ball falls - game over
	game.ballY = float64(game.areaY + game.areaH + 1)
	game.Update(0.01)

	if reportedScore != 1000 {
		t.Errorf("Expected reported score 1000, got %d", reportedScore)
	}
}

func TestBreakoutMultiHitBrick(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	// Create a 2-hit brick
	game.bricks = []Brick{
		{X: 10, Y: 5, Width: 4, Hits: 2, Points: 50},
	}

	game.ballLaunched = true
	game.ballX = 12.0
	game.ballY = 5.0
	game.ballVY = 1.0 // Moving down into brick

	game.Update(0.01)

	// Brick should still exist but with 1 hit (or possibly 0 depending on collision logic)
	if len(game.bricks) != 1 {
		t.Error("2-hit brick should not be destroyed on first hit")
	}
	// Brick should have been hit at least once
	if game.bricks[0].Hits >= 2 {
		t.Errorf("Expected brick to be damaged (hits < 2), got %d", game.bricks[0].Hits)
	}
}

func TestBreakoutRender(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 24}
	game.Init(ctx)

	var buf bytes.Buffer
	screen := engine.NewScreen(40, 24, &buf)

	err := game.Render(screen)
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}
}
