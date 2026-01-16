package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"slip/internal/platform"
	"slip/internal/plugins"
	"slip/internal/state"
)

// GameEntry represents a game in the store
type GameEntry struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Version     string            `json:"version"`
	Type        string            `json:"type"`
	Category    string            `json:"category"`
	Downloads   map[string]string `json:"downloads"` // platform -> URL
	Checksum    map[string]string `json:"checksum"`  // platform -> sha256
	Featured    bool              `json:"featured"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// Registry represents the game registry response
type Registry struct {
	Games     []GameEntry `json:"games"`
	Featured  []string    `json:"featured"`
	UpdatedAt string      `json:"updated_at"`
}

// Client is the store client
type Client struct {
	registryURL string
	httpClient  *http.Client
	cache       *Registry
	cacheTime   time.Time
	cacheTTL    time.Duration
}

// NewClient creates a new store client
func NewClient() *Client {
	cfg, _ := state.LoadConfig()
	registryURL := cfg.RegistryURL
	if registryURL == "" {
		registryURL = "https://raw.githubusercontent.com/slip-games/registry/main/games.json"
	}

	return &Client{
		registryURL: registryURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cacheTTL: 5 * time.Minute,
	}
}

// FetchRegistry fetches the game registry
func (c *Client) FetchRegistry() (*Registry, error) {
	// Check cache
	if c.cache != nil && time.Since(c.cacheTime) < c.cacheTTL {
		return c.cache, nil
	}

	resp, err := c.httpClient.Get(c.registryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var registry Registry
	if err := json.Unmarshal(body, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	c.cache = &registry
	c.cacheTime = time.Now()

	return &registry, nil
}

// ListGames returns all games from the registry
func (c *Client) ListGames() ([]GameEntry, error) {
	registry, err := c.FetchRegistry()
	if err != nil {
		return nil, err
	}
	return registry.Games, nil
}

// SearchGames searches for games by name or description
func (c *Client) SearchGames(query string) ([]GameEntry, error) {
	registry, err := c.FetchRegistry()
	if err != nil {
		return nil, err
	}

	var results []GameEntry
	query = toLowerCase(query)
	for _, game := range registry.Games {
		if contains(toLowerCase(game.Name), query) ||
			contains(toLowerCase(game.Description), query) ||
			contains(toLowerCase(game.ID), query) {
			results = append(results, game)
		}
	}
	return results, nil
}

// GetGame returns a specific game from the registry
func (c *Client) GetGame(id string) (*GameEntry, error) {
	registry, err := c.FetchRegistry()
	if err != nil {
		return nil, err
	}

	for _, game := range registry.Games {
		if game.ID == id {
			return &game, nil
		}
	}
	return nil, fmt.Errorf("game not found: %s", id)
}

// GetDownloadURL returns the download URL for the current platform
func (c *Client) GetDownloadURL(game *GameEntry) (string, error) {
	plat := platform.GetPlatformString()
	url, ok := game.Downloads[plat]
	if !ok {
		return "", fmt.Errorf("game not available for platform %s", plat)
	}
	return url, nil
}

// verifyChecksum verifies the SHA256 checksum of a file
func (c *Client) verifyChecksum(filePath, expectedChecksum string) error {
	if expectedChecksum == "" {
		return nil // No checksum to verify
	}

	// Remove "sha256:" prefix if present
	expected := expectedChecksum
	if strings.HasPrefix(strings.ToLower(expected), "sha256:") {
		expected = expected[7:]
	}

	// Calculate file checksum
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for checksum: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("failed to read file for checksum: %w", err)
	}

	actual := hex.EncodeToString(h.Sum(nil))

	if actual != expected {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	}

	return nil
}

// ProgressWriter wraps an io.Writer to track download progress
type ProgressWriter struct {
	Total      int64
	Written    int64
	OnProgress func(written, total int64)
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Written += int64(n)
	if pw.OnProgress != nil {
		pw.OnProgress(pw.Written, pw.Total)
	}
	return n, nil
}

// Download downloads a game to the specified path
func (c *Client) Download(game *GameEntry, destPath string) error {
	url, err := c.GetDownloadURL(game)
	if err != nil {
		return err
	}

	// Determine destination
	if destPath == "" {
		destPath = filepath.Join(state.GamesDir(), game.ID)
	}

	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Download file
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	// Determine file type from URL or content-type
	fileName := filepath.Base(url)
	if fileName == "" || fileName == "." {
		fileName = game.ID + ".tar.gz"
	}

	tmpFile := filepath.Join(destPath, fileName)
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy data
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Close file before verification
	out.Close()

	// Verify checksum if available
	plat := platform.GetPlatformString()
	if expectedChecksum, ok := game.Checksum[plat]; ok && expectedChecksum != "" {
		if err := c.verifyChecksum(tmpFile, expectedChecksum); err != nil {
			os.Remove(tmpFile) // Remove corrupted file
			return fmt.Errorf("checksum verification failed: %w", err)
		}
	}

	// Create manifest
	manifest := plugins.Manifest{
		ID:          game.ID,
		Name:        game.Name,
		Description: game.Description,
		Author:      game.Author,
		Version:     game.Version,
		Type:        game.Type,
		Category:    game.Category,
	}

	manifestPath := filepath.Join(destPath, "manifest.json")
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// DownloadWithProgress downloads a game with progress callback
func (c *Client) DownloadWithProgress(game *GameEntry, destPath string, progress func(pct int)) error {
	url, err := c.GetDownloadURL(game)
	if err != nil {
		return err
	}

	// Determine destination
	if destPath == "" {
		destPath = filepath.Join(state.GamesDir(), game.ID)
	}

	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Download file
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	// Determine file type from URL or content-type
	fileName := filepath.Base(url)
	if fileName == "" || fileName == "." {
		fileName = game.ID + ".tar.gz"
	}

	tmpFile := filepath.Join(destPath, fileName)
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Get content length
	total := resp.ContentLength

	// Create progress writer
	pw := &ProgressWriter{
		Total: total,
		OnProgress: func(written, total int64) {
			if total > 0 && progress != nil {
				pct := int(float64(written) / float64(total) * 100)
				progress(pct)
			}
		},
	}

	// Copy with progress
	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Close file before verification
	out.Close()

	// Verify checksum if available
	plat := platform.GetPlatformString()
	if expectedChecksum, ok := game.Checksum[plat]; ok && expectedChecksum != "" {
		if err := c.verifyChecksum(tmpFile, expectedChecksum); err != nil {
			os.Remove(tmpFile) // Remove corrupted file
			return fmt.Errorf("checksum verification failed: %w", err)
		}
	}

	// Create manifest
	manifest := plugins.Manifest{
		ID:          game.ID,
		Name:        game.Name,
		Description: game.Description,
		Author:      game.Author,
		Version:     game.Version,
		Type:        game.Type,
		Category:    game.Category,
	}

	manifestPath := filepath.Join(destPath, "manifest.json")
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// Helper functions

func toLowerCase(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
