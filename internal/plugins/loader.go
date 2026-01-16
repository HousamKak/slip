package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"slip/internal/engine"
	"slip/internal/state"
)

// Manifest describes a plugin/game package
type Manifest struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Version     string            `json:"version"`
	Type        string            `json:"type"` // "lua", "wasm", "native"
	Entry       string            `json:"entry"`
	MinWidth    int               `json:"min_width"`
	MinHeight   int               `json:"min_height"`
	Category    string            `json:"category"` // "game", "animation"
	Assets      []string          `json:"assets"`
	Platforms   map[string]string `json:"platforms"` // For native plugins
}

// Loader loads and manages plugins
type Loader struct {
	gamesDir  string
	manifests map[string]*Manifest
}

// NewLoader creates a new plugin loader
func NewLoader() *Loader {
	return &Loader{
		gamesDir:  state.GamesDir(),
		manifests: make(map[string]*Manifest),
	}
}

// Scan scans the games directory for installed plugins
func (l *Loader) Scan() error {
	if err := os.MkdirAll(l.gamesDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(l.gamesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(l.gamesDir, entry.Name(), "manifest.json")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		manifest, err := l.loadManifest(manifestPath)
		if err != nil {
			continue
		}

		l.manifests[manifest.ID] = manifest
	}

	return nil
}

func (l *Loader) loadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// List returns all installed plugins
func (l *Loader) List() []*Manifest {
	var list []*Manifest
	for _, m := range l.manifests {
		list = append(list, m)
	}
	return list
}

// Get returns a plugin manifest by ID
func (l *Loader) Get(id string) (*Manifest, bool) {
	m, ok := l.manifests[id]
	return m, ok
}

// Load loads a plugin as a Game
func (l *Loader) Load(id string) (engine.Game, error) {
	manifest, ok := l.manifests[id]
	if !ok {
		return nil, fmt.Errorf("plugin not found: %s", id)
	}

	pluginDir := filepath.Join(l.gamesDir, id)

	switch manifest.Type {
	case "lua":
		return l.loadLuaPlugin(manifest, pluginDir)
	case "wasm":
		return nil, fmt.Errorf("wasm plugins not yet supported")
	case "native":
		return nil, fmt.Errorf("native plugins not yet supported")
	default:
		return nil, fmt.Errorf("unknown plugin type: %s", manifest.Type)
	}
}

func (l *Loader) loadLuaPlugin(manifest *Manifest, dir string) (engine.Game, error) {
	return NewLuaGame(manifest, dir)
}

// Install installs a plugin from a directory or archive
func (l *Loader) Install(source string) error {
	// Future: download from URL, extract archive, validate manifest
	return fmt.Errorf("plugin installation not yet implemented")
}

// Uninstall removes a plugin
func (l *Loader) Uninstall(id string) error {
	manifest, ok := l.manifests[id]
	if !ok {
		return fmt.Errorf("plugin not found: %s", id)
	}

	pluginDir := filepath.Join(l.gamesDir, manifest.ID)
	if err := os.RemoveAll(pluginDir); err != nil {
		return err
	}

	delete(l.manifests, id)
	return nil
}

// GamesDir returns the games directory path
func (l *Loader) GamesDir() string {
	return l.gamesDir
}
