package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// State represents the current Slip runtime state
type State struct {
	PID          int    `json:"pid"`
	StartedAt    int64  `json:"started_at"`
	LastSeen     int64  `json:"last_seen"`
	CurrentApp   string `json:"current_app"`
	Message      string `json:"message"`
	MessageUntil int64  `json:"message_until"`
}

// Config represents user configuration
type Config struct {
	DefaultDock    string `json:"default_dock"`
	DefaultSize    string `json:"default_size"`
	DefaultFPS     int    `json:"default_fps"`
	Theme          string `json:"theme"`
	ShowFPSCounter bool   `json:"show_fps_counter"`
	RegistryURL    string `json:"registry_url"`
	AutoUpdate     bool   `json:"auto_update"`
	PlayerName     string `json:"player_name"`
	TUIMode        string `json:"tui_mode"` // "custom" or "bubbletea"
}

// Score represents a game score
type Score struct {
	Value     int    `json:"value"`
	Player    string `json:"player"`
	Timestamp int64  `json:"timestamp"`
	GameID    string `json:"game_id"`
}

// ScoreBoard holds all high scores
type ScoreBoard struct {
	Scores map[string][]Score `json:"scores"`
	mu     sync.RWMutex
}

var (
	defaultConfig = Config{
		DefaultDock:    "right",
		DefaultSize:    "40x20",
		DefaultFPS:     30,
		Theme:          "default",
		ShowFPSCounter: false,
		RegistryURL:    "https://raw.githubusercontent.com/slip-games/registry/main/games.json",
		AutoUpdate:     true,
		PlayerName:     "Player",
	}
)

// DefaultPath returns the default state file path
func DefaultPath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "slip", "state.json")
}

// ConfigPath returns the config file path
func ConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir, _ = os.UserHomeDir()
	}
	return filepath.Join(dir, "slip", "config.json")
}

// ScoresPath returns the scores file path
func ScoresPath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "slip", "scores.json")
}

// GamesDir returns the directory for installed games
func GamesDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "slip", "games")
}

// Load loads the state from disk
func Load(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return State{}, err
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, err
	}
	return st, nil
}

// Save saves the state to disk
func Save(path string, st State) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(dir, fmt.Sprintf(".state_%d.tmp", time.Now().UnixNano()))
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return os.WriteFile(path, data, 0o644)
	}
	return nil
}

// UpdateHeartbeat updates the heartbeat in state
func UpdateHeartbeat(path string, pid int, startedAt time.Time, currentApp string) error {
	st, err := Load(path)
	if err != nil {
		st = State{}
	}
	st.PID = pid
	st.StartedAt = startedAt.Unix()
	st.LastSeen = time.Now().Unix()
	st.CurrentApp = currentApp
	return Save(path, st)
}

// WriteMessage writes a message to state
func WriteMessage(path string, msg string, until time.Time) error {
	st, err := Load(path)
	if err != nil {
		st = State{}
	}
	st.Message = msg
	st.MessageUntil = until.Unix()
	return Save(path, st)
}

// LoadConfig loads configuration from disk
func LoadConfig() (Config, error) {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig, nil
		}
		return Config{}, err
	}
	cfg := defaultConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// SaveConfig saves configuration to disk
func SaveConfig(cfg Config) error {
	path := ConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// GetConfig returns the current config value for a key
func GetConfig(key string) (string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return "", err
	}
	switch key {
	case "default_dock":
		return cfg.DefaultDock, nil
	case "default_size":
		return cfg.DefaultSize, nil
	case "default_fps":
		return fmt.Sprintf("%d", cfg.DefaultFPS), nil
	case "theme":
		return cfg.Theme, nil
	case "show_fps_counter":
		return fmt.Sprintf("%v", cfg.ShowFPSCounter), nil
	case "registry_url":
		return cfg.RegistryURL, nil
	case "auto_update":
		return fmt.Sprintf("%v", cfg.AutoUpdate), nil
	case "player_name":
		return cfg.PlayerName, nil
	case "tui_mode":
		return cfg.TUIMode, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// SetConfig sets a config value
func SetConfig(key, value string) error {
	cfg, err := LoadConfig()
	if err != nil {
		cfg = defaultConfig
	}
	switch key {
	case "default_dock":
		cfg.DefaultDock = value
	case "default_size":
		cfg.DefaultSize = value
	case "default_fps":
		var fps int
		fmt.Sscanf(value, "%d", &fps)
		if fps > 0 {
			cfg.DefaultFPS = fps
		}
	case "theme":
		cfg.Theme = value
	case "show_fps_counter":
		cfg.ShowFPSCounter = value == "true" || value == "1" || value == "yes"
	case "registry_url":
		cfg.RegistryURL = value
	case "auto_update":
		cfg.AutoUpdate = value == "true" || value == "1" || value == "yes"
	case "player_name":
		cfg.PlayerName = value
	case "tui_mode":
		if value != "custom" && value != "bubbletea" {
			return fmt.Errorf("tui_mode must be 'custom' or 'bubbletea'")
		}
		cfg.TUIMode = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return SaveConfig(cfg)
}

// NewScoreBoard creates a new scoreboard
func NewScoreBoard() *ScoreBoard {
	return &ScoreBoard{
		Scores: make(map[string][]Score),
	}
}

// LoadScores loads scores from disk
func LoadScores() (*ScoreBoard, error) {
	path := ScoresPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewScoreBoard(), nil
		}
		return nil, err
	}
	sb := NewScoreBoard()
	if err := json.Unmarshal(data, sb); err != nil {
		return nil, err
	}
	return sb, nil
}

// Save saves scores to disk
func (sb *ScoreBoard) Save() error {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	path := ScoresPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(sb, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// AddScore adds a score for a game
func (sb *ScoreBoard) AddScore(gameID string, score Score) error {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	score.GameID = gameID
	if score.Timestamp == 0 {
		score.Timestamp = time.Now().Unix()
	}

	scores := sb.Scores[gameID]
	scores = append(scores, score)

	// Sort descending by value
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Value > scores[j].Value
	})

	// Keep only top 10
	if len(scores) > 10 {
		scores = scores[:10]
	}

	sb.Scores[gameID] = scores
	return sb.Save()
}

// GetTopScores returns top scores for a game
func (sb *ScoreBoard) GetTopScores(gameID string, limit int) []Score {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	scores := sb.Scores[gameID]
	if limit <= 0 || limit > len(scores) {
		limit = len(scores)
	}
	if limit == 0 {
		return nil
	}
	result := make([]Score, limit)
	copy(result, scores[:limit])
	return result
}

// GetHighScore returns the high score for a game
func (sb *ScoreBoard) GetHighScore(gameID string) int {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	scores := sb.Scores[gameID]
	if len(scores) == 0 {
		return 0
	}
	return scores[0].Value
}

// IsHighScore checks if a score is a new high score
func (sb *ScoreBoard) IsHighScore(gameID string, value int) bool {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	scores := sb.Scores[gameID]
	if len(scores) == 0 {
		return true
	}
	if len(scores) < 10 {
		return true
	}
	return value > scores[len(scores)-1].Value
}
