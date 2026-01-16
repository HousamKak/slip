package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.httpClient == nil {
		t.Error("httpClient should be initialized")
	}
	if client.registryURL == "" {
		t.Error("registryURL should not be empty")
	}
	if client.cacheTTL == 0 {
		t.Error("cacheTTL should be set")
	}
}

func TestFetchRegistry(t *testing.T) {
	// Create mock registry
	mockRegistry := Registry{
		Games: []GameEntry{
			{
				ID:          "test-game",
				Name:        "Test Game",
				Description: "A test game",
				Author:      "Test Author",
				Version:     "1.0.0",
				Type:        "lua",
				Category:    "game",
			},
		},
		Featured:  []string{"test-game"},
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockRegistry)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &Client{
		registryURL: server.URL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		cacheTTL: 5 * time.Minute,
	}

	registry, err := client.FetchRegistry()
	if err != nil {
		t.Fatalf("FetchRegistry() failed: %v", err)
	}

	if len(registry.Games) != 1 {
		t.Errorf("Expected 1 game, got %d", len(registry.Games))
	}

	if registry.Games[0].ID != "test-game" {
		t.Errorf("Expected game ID 'test-game', got '%s'", registry.Games[0].ID)
	}
}

func TestFetchRegistryCache(t *testing.T) {
	callCount := 0

	// Create test server that counts calls
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		registry := Registry{
			Games: []GameEntry{
				{ID: "test-game", Name: "Test Game"},
			},
		}
		json.NewEncoder(w).Encode(registry)
	}))
	defer server.Close()

	client := &Client{
		registryURL: server.URL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		cacheTTL: 1 * time.Minute,
	}

	// First call - should hit server
	_, err := client.FetchRegistry()
	if err != nil {
		t.Fatalf("First FetchRegistry() failed: %v", err)
	}

	// Second call - should use cache
	_, err = client.FetchRegistry()
	if err != nil {
		t.Fatalf("Second FetchRegistry() failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 server call (cached), got %d", callCount)
	}
}

func TestFetchRegistryError(t *testing.T) {
	// Create server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &Client{
		registryURL: server.URL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		cacheTTL: 5 * time.Minute,
	}

	_, err := client.FetchRegistry()
	if err == nil {
		t.Error("Expected error when server returns 500")
	}
}

func TestListGames(t *testing.T) {
	mockRegistry := Registry{
		Games: []GameEntry{
			{ID: "game1", Name: "Game 1"},
			{ID: "game2", Name: "Game 2"},
			{ID: "game3", Name: "Game 3"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockRegistry)
	}))
	defer server.Close()

	client := &Client{
		registryURL: server.URL,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
		cacheTTL:    5 * time.Minute,
	}

	games, err := client.ListGames()
	if err != nil {
		t.Fatalf("ListGames() failed: %v", err)
	}

	if len(games) != 3 {
		t.Errorf("Expected 3 games, got %d", len(games))
	}
}

func TestSearchGames(t *testing.T) {
	mockRegistry := Registry{
		Games: []GameEntry{
			{ID: "snake", Name: "Snake Game", Description: "Classic snake"},
			{ID: "pong", Name: "Pong", Description: "Table tennis game"},
			{ID: "breakout", Name: "Breakout", Description: "Brick breaker"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockRegistry)
	}))
	defer server.Close()

	client := &Client{
		registryURL: server.URL,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
		cacheTTL:    5 * time.Minute,
	}

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{"Search by name", "snake", 1},
		{"Search by description", "game", 2}, // "snake game" and "table tennis game"
		{"Search by ID", "pong", 1},
		{"Case insensitive", "SNAKE", 1},
		{"No matches", "tetris", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := client.SearchGames(tt.query)
			if err != nil {
				t.Fatalf("SearchGames() failed: %v", err)
			}
			if len(results) != tt.expected {
				t.Errorf("Expected %d results, got %d", tt.expected, len(results))
			}
		})
	}
}

func TestGetGame(t *testing.T) {
	mockRegistry := Registry{
		Games: []GameEntry{
			{ID: "game1", Name: "Game 1"},
			{ID: "game2", Name: "Game 2"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockRegistry)
	}))
	defer server.Close()

	client := &Client{
		registryURL: server.URL,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
		cacheTTL:    5 * time.Minute,
	}

	// Test existing game
	game, err := client.GetGame("game1")
	if err != nil {
		t.Fatalf("GetGame() failed: %v", err)
	}
	if game.ID != "game1" {
		t.Errorf("Expected game ID 'game1', got '%s'", game.ID)
	}

	// Test non-existing game
	_, err = client.GetGame("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent game")
	}
}

func TestGetDownloadURL(t *testing.T) {
	game := &GameEntry{
		ID:   "test-game",
		Name: "Test Game",
		Downloads: map[string]string{
			"linux/amd64":   "https://example.com/linux.tar.gz",
			"darwin/amd64":  "https://example.com/macos.tar.gz",
			"darwin/arm64":  "https://example.com/macos-arm.tar.gz",
			"windows/amd64": "https://example.com/windows.zip",
		},
	}

	client := NewClient()

	// Test with existing platform
	url, err := client.GetDownloadURL(game)
	if err != nil {
		t.Fatalf("GetDownloadURL() failed: %v", err)
	}
	if url == "" {
		t.Error("Expected non-empty download URL")
	}

	// Test with game missing current platform
	gameNoPlatform := &GameEntry{
		ID:        "test-game",
		Downloads: map[string]string{},
	}
	_, err = client.GetDownloadURL(gameNoPlatform)
	if err == nil {
		t.Error("Expected error when platform not available")
	}
}

func TestVerifyChecksum(t *testing.T) {
	client := NewClient()

	// Create temporary file with known content
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello, World!")

	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate expected checksum
	h := sha256.New()
	h.Write(content)
	expectedChecksum := hex.EncodeToString(h.Sum(nil))

	tests := []struct {
		name        string
		checksum    string
		shouldError bool
	}{
		{"Correct checksum", expectedChecksum, false},
		{"Checksum with sha256: prefix", "sha256:" + expectedChecksum, false},
		{"Incorrect checksum", "0000000000000000000000000000000000000000000000000000000000000000", true},
		{"Empty checksum (skip verification)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.verifyChecksum(testFile, tt.checksum)
			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestProgressWriter(t *testing.T) {
	progressCalls := 0
	var lastWritten, lastTotal int64

	pw := &ProgressWriter{
		Total: 100,
		OnProgress: func(written, total int64) {
			progressCalls++
			lastWritten = written
			lastTotal = total
		},
	}

	// Write some data
	data := []byte("Hello, World!")
	n, err := pw.Write(data)
	if err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}

	if pw.Written != int64(len(data)) {
		t.Errorf("Expected Written = %d, got %d", len(data), pw.Written)
	}

	if progressCalls != 1 {
		t.Errorf("Expected 1 progress callback, got %d", progressCalls)
	}

	if lastWritten != int64(len(data)) {
		t.Errorf("Expected progress callback with written=%d, got %d", len(data), lastWritten)
	}

	if lastTotal != 100 {
		t.Errorf("Expected progress callback with total=100, got %d", lastTotal)
	}
}

func TestProgressWriterWithoutCallback(t *testing.T) {
	pw := &ProgressWriter{
		Total: 100,
	}

	// Should not panic without callback
	data := []byte("test")
	_, err := pw.Write(data)
	if err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	if pw.Written != int64(len(data)) {
		t.Errorf("Expected Written = %d, got %d", len(data), pw.Written)
	}
}

func TestToLowerCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HELLO", "hello"},
		{"Hello", "hello"},
		{"hello", "hello"},
		{"", ""},
		{"123", "123"},
	}

	for _, tt := range tests {
		result := toLowerCase(tt.input)
		if result != tt.expected {
			t.Errorf("toLowerCase(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		haystack string
		needle   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "WORLD", false}, // Case sensitive
		{"hello world", "foo", false},
		{"", "", true},
		{"hello", "", true},
		{"", "hello", false},
	}

	for _, tt := range tests {
		result := contains(tt.haystack, tt.needle)
		if result != tt.expected {
			t.Errorf("contains(%q, %q) = %v, expected %v", tt.haystack, tt.needle, result, tt.expected)
		}
	}
}
