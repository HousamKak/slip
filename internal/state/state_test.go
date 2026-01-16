package state

import (
	"testing"
	"time"
)

func TestScoreBoardAddScore(t *testing.T) {
	sb := NewScoreBoard()

	score := Score{
		Value:     100,
		Player:    "Alice",
		Timestamp: time.Now().Unix(),
	}

	err := sb.AddScore("snake", score)
	if err != nil {
		t.Fatalf("AddScore failed: %v", err)
	}

	scores := sb.GetTopScores("snake", 10)
	if len(scores) != 1 {
		t.Fatalf("Expected 1 score, got %d", len(scores))
	}

	if scores[0].Value != 100 {
		t.Errorf("Expected score 100, got %d", scores[0].Value)
	}
	if scores[0].Player != "Alice" {
		t.Errorf("Expected player 'Alice', got '%s'", scores[0].Player)
	}
}

func TestScoreBoardSorting(t *testing.T) {
	sb := NewScoreBoard()

	scores := []Score{
		{Value: 50, Player: "Bob", Timestamp: time.Now().Unix()},
		{Value: 200, Player: "Alice", Timestamp: time.Now().Unix()},
		{Value: 150, Player: "Charlie", Timestamp: time.Now().Unix()},
		{Value: 100, Player: "Dave", Timestamp: time.Now().Unix()},
	}

	for _, score := range scores {
		sb.AddScore("snake", score)
	}

	top := sb.GetTopScores("snake", 10)

	// Should be sorted descending by value
	if len(top) != 4 {
		t.Fatalf("Expected 4 scores, got %d", len(top))
	}

	if top[0].Value != 200 {
		t.Errorf("Top score should be 200, got %d", top[0].Value)
	}
	if top[1].Value != 150 {
		t.Errorf("Second score should be 150, got %d", top[1].Value)
	}
	if top[2].Value != 100 {
		t.Errorf("Third score should be 100, got %d", top[2].Value)
	}
	if top[3].Value != 50 {
		t.Errorf("Fourth score should be 50, got %d", top[3].Value)
	}
}

func TestScoreBoardHighScore(t *testing.T) {
	sb := NewScoreBoard()

	scores := []Score{
		{Value: 100, Player: "Alice", Timestamp: time.Now().Unix()},
		{Value: 200, Player: "Bob", Timestamp: time.Now().Unix()},
		{Value: 150, Player: "Charlie", Timestamp: time.Now().Unix()},
	}

	for _, score := range scores {
		sb.AddScore("breakout", score)
	}

	highScore := sb.GetHighScore("breakout")
	if highScore != 200 {
		t.Errorf("Expected high score 200, got %d", highScore)
	}

	// Non-existent game should return 0
	nonexistent := sb.GetHighScore("nonexistent")
	if nonexistent != 0 {
		t.Errorf("Expected high score 0 for non-existent game, got %d", nonexistent)
	}
}

func TestScoreBoardIsHighScore(t *testing.T) {
	sb := NewScoreBoard()

	sb.AddScore("pong", Score{Value: 100, Player: "Alice", Timestamp: time.Now().Unix()})
	sb.AddScore("pong", Score{Value: 150, Player: "Bob", Timestamp: time.Now().Unix()})

	// 200 should be a high score
	if !sb.IsHighScore("pong", 200) {
		t.Error("200 should be a high score")
	}

	// 120 should not be a high score (less than current high of 150)
	if sb.IsHighScore("pong", 120) {
		t.Error("120 should not be a high score")
	}

	// For new game, any score > 0 should be a high score
	if !sb.IsHighScore("newgame", 10) {
		t.Error("First score should always be a high score")
	}
}

func TestScoreBoardMultipleGames(t *testing.T) {
	sb := NewScoreBoard()

	sb.AddScore("snake", Score{Value: 100, Player: "Alice", Timestamp: time.Now().Unix()})
	sb.AddScore("pong", Score{Value: 200, Player: "Bob", Timestamp: time.Now().Unix()})
	sb.AddScore("snake", Score{Value: 150, Player: "Charlie", Timestamp: time.Now().Unix()})

	snakeScores := sb.GetTopScores("snake", 10)
	pongScores := sb.GetTopScores("pong", 10)

	if len(snakeScores) != 2 {
		t.Errorf("Expected 2 snake scores, got %d", len(snakeScores))
	}
	if len(pongScores) != 1 {
		t.Errorf("Expected 1 pong score, got %d", len(pongScores))
	}

	if snakeScores[0].Value != 150 {
		t.Errorf("Top snake score should be 150, got %d", snakeScores[0].Value)
	}
	if pongScores[0].Value != 200 {
		t.Errorf("Top pong score should be 200, got %d", pongScores[0].Value)
	}
}

func TestStateLoadSave(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/state.json"

	state := State{
		PID:        12345,
		StartedAt:  time.Now().Unix(),
		CurrentApp: "snake",
		Message:    "Test message",
	}

	// Save
	err := Save(tmpFile, state)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loaded, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.PID != state.PID {
		t.Errorf("PID mismatch: %d != %d", loaded.PID, state.PID)
	}
	if loaded.CurrentApp != state.CurrentApp {
		t.Errorf("CurrentApp mismatch: %s != %s", loaded.CurrentApp, state.CurrentApp)
	}
	if loaded.Message != state.Message {
		t.Errorf("Message mismatch: %s != %s", loaded.Message, state.Message)
	}
}

func TestConfigPaths(t *testing.T) {
	// These functions should not panic or return empty
	defaultPath := DefaultPath()
	if defaultPath == "" {
		t.Error("DefaultPath() returned empty string")
	}

	configPath := ConfigPath()
	if configPath == "" {
		t.Error("ConfigPath() returned empty string")
	}

	scoresPath := ScoresPath()
	if scoresPath == "" {
		t.Error("ScoresPath() returned empty string")
	}

	gamesDir := GamesDir()
	if gamesDir == "" {
		t.Error("GamesDir() returned empty string")
	}
}

func TestNewScoreBoard(t *testing.T) {
	sb := NewScoreBoard()
	if sb == nil {
		t.Fatal("NewScoreBoard() returned nil")
	}
	if sb.Scores == nil {
		t.Error("Scores map should be initialized")
	}
}
