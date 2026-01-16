package pong

import (
	"bytes"
	"testing"

	"slip/internal/engine"
)

func TestPongNew(t *testing.T) {
	game := New()
	if game == nil {
		t.Fatal("New() returned nil")
	}
	if game.winScore != 5 {
		t.Errorf("Expected winScore 5, got %d", game.winScore)
	}
	if game.rng == nil {
		t.Error("Expected rng to be initialized")
	}
}

func TestPongInfo(t *testing.T) {
	game := New()
	info := game.Info()

	if info.ID != "pong" {
		t.Errorf("Expected ID 'pong', got '%s'", info.ID)
	}
	if info.Name != "Pong" {
		t.Errorf("Expected Name 'Pong', got '%s'", info.Name)
	}
	if info.MinWidth != 30 {
		t.Errorf("Expected MinWidth 30, got %d", info.MinWidth)
	}
	if info.MinHeight != 15 {
		t.Errorf("Expected MinHeight 15, got %d", info.MinHeight)
	}
}

func TestPongInit(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{
		Width:  40,
		Height: 20,
	}

	err := game.Init(ctx)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Ball should be in center
	expectedBallX := float64(game.areaX + game.areaW/2)
	expectedBallY := float64(game.areaY + game.areaH/2)
	if game.ballX != expectedBallX {
		t.Errorf("Expected ballX %f, got %f", expectedBallX, game.ballX)
	}
	if game.ballY != expectedBallY {
		t.Errorf("Expected ballY %f, got %f", expectedBallY, game.ballY)
	}

	// Scores should be 0
	if game.playerScore != 0 {
		t.Errorf("Expected playerScore 0, got %d", game.playerScore)
	}
	if game.cpuScore != 0 {
		t.Errorf("Expected cpuScore 0, got %d", game.cpuScore)
	}

	// Game should not be over
	if game.gameOver {
		t.Error("Game should not be over at start")
	}

	// Countdown should be active
	if game.countdown != 3 {
		t.Errorf("Expected countdown 3, got %d", game.countdown)
	}
}

func TestPongCountdown(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	initialCountdown := game.countdown
	ballX := game.ballX

	// Update for less than 1 second - countdown should not change
	game.Update(0.5)
	if game.countdown != initialCountdown {
		t.Error("Countdown should not change before 1 second")
	}

	// Ball should not move during countdown
	if game.ballX != ballX {
		t.Error("Ball should not move during countdown")
	}

	// Complete countdown
	game.Update(0.5) // Now at 1 second
	if game.countdown != initialCountdown-1 {
		t.Errorf("Expected countdown %d, got %d", initialCountdown-1, game.countdown)
	}
}

func TestPongBallMovement(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Skip countdown
	game.countdown = 0
	game.ballVX = 10.0
	game.ballVY = 5.0

	initialX := game.ballX
	initialY := game.ballY

	game.Update(0.1)

	// Ball should have moved
	expectedX := initialX + 10.0*0.1
	expectedY := initialY + 5.0*0.1

	if game.ballX != expectedX {
		t.Errorf("Expected ballX %f, got %f", expectedX, game.ballX)
	}
	if game.ballY != expectedY {
		t.Errorf("Expected ballY %f, got %f", expectedY, game.ballY)
	}
}

func TestPongBallWallBounce(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Skip countdown
	game.countdown = 0

	// Place ball near top wall moving up
	game.ballY = float64(game.areaY) + 0.5
	game.ballVY = -10.0

	game.Update(0.1)

	// Ball should bounce (VY should become positive)
	if game.ballVY <= 0 {
		t.Error("Ball should bounce off top wall (VY should be positive)")
	}
}

func TestPongPlayerPaddleMovement(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	initialY := game.playerY

	// Move up
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyUp})
	if game.playerY >= initialY {
		t.Error("Paddle should move up (Y decrease)")
	}

	// Reset and move down
	game.playerY = initialY
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyDown})
	if game.playerY <= initialY {
		t.Error("Paddle should move down (Y increase)")
	}
}

func TestPongWASDControls(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	initialY := game.playerY

	// Move up with W
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'w'})
	if game.playerY >= initialY {
		t.Error("Paddle should move up with 'w'")
	}

	// Reset and move down with S
	game.playerY = initialY
	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 's'})
	if game.playerY <= initialY {
		t.Error("Paddle should move down with 's'")
	}
}

func TestPongPlayerScores(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Skip countdown
	game.countdown = 0

	// Move ball past CPU paddle (right side)
	game.ballX = float64(game.areaX + game.areaW + 1)
	game.ballVX = 10.0

	game.Update(0.01)

	if game.playerScore != 1 {
		t.Errorf("Expected playerScore 1, got %d", game.playerScore)
	}
}

func TestPongCPUScores(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Skip countdown
	game.countdown = 0

	// Move ball past player paddle (left side)
	game.ballX = float64(game.areaX) - 1
	game.ballVX = -10.0

	game.Update(0.01)

	if game.cpuScore != 1 {
		t.Errorf("Expected cpuScore 1, got %d", game.cpuScore)
	}
}

func TestPongPlayerWins(t *testing.T) {
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

	// Set player to win condition
	game.playerScore = game.winScore - 1
	game.countdown = 0
	game.ballX = float64(game.areaX + game.areaW + 1)
	game.ballVX = 10.0

	game.Update(0.01)

	if !game.gameOver {
		t.Error("Game should be over when player wins")
	}
	if game.winner != "Player" {
		t.Errorf("Expected winner 'Player', got '%s'", game.winner)
	}
	if reportedScore != game.playerScore*100 {
		t.Errorf("Expected score %d, got %d", game.playerScore*100, reportedScore)
	}
}

func TestPongCPUWins(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Set CPU to win condition
	game.cpuScore = game.winScore - 1
	game.countdown = 0
	game.ballX = float64(game.areaX) - 1
	game.ballVX = -10.0

	game.Update(0.01)

	if !game.gameOver {
		t.Error("Game should be over when CPU wins")
	}
	if game.winner != "CPU" {
		t.Errorf("Expected winner 'CPU', got '%s'", game.winner)
	}
}

func TestPongPause(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'p'})

	if !game.paused {
		t.Error("Game should be paused")
	}

	// Ball should not move while paused
	game.countdown = 0
	ballX := game.ballX
	game.ballVX = 10.0
	game.Update(0.1)

	if game.ballX != ballX {
		t.Error("Ball should not move while paused")
	}
}

func TestPongRestart(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	game.playerScore = 3
	game.cpuScore = 2
	game.gameOver = true

	game.HandleInput(engine.Input{Type: engine.InputKeyPress, Key: engine.KeyRune, Rune: 'r'})

	if game.playerScore != 0 {
		t.Error("Player score should be reset")
	}
	if game.cpuScore != 0 {
		t.Error("CPU score should be reset")
	}
	if game.gameOver {
		t.Error("Game over should be reset")
	}
}

func TestPongPaddleClamp(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	// Try to move paddle above play area
	game.playerY = float64(game.areaY) - 10
	game.countdown = 0
	game.Update(0.01)

	if game.playerY < float64(game.areaY) {
		t.Error("Paddle Y should be clamped to play area minimum")
	}
}

func TestPongRender(t *testing.T) {
	game := New()
	ctx := &engine.GameContext{Width: 40, Height: 20}
	game.Init(ctx)

	var buf bytes.Buffer
	screen := engine.NewScreen(40, 20, &buf)

	err := game.Render(screen)
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}
}
