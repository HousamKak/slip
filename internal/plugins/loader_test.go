package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader()
	if loader == nil {
		t.Fatal("NewLoader() returned nil")
	}
	if loader.manifests == nil {
		t.Error("manifests map should be initialized")
	}
	if loader.gamesDir == "" {
		t.Error("gamesDir should not be empty")
	}
}

func TestLoadManifest(t *testing.T) {
	loader := NewLoader()

	// Create temporary manifest file
	tmpDir := t.TempDir()
	manifestPath := filepath.Join(tmpDir, "manifest.json")

	manifest := Manifest{
		ID:          "test-game",
		Name:        "Test Game",
		Description: "A test game",
		Author:      "Test Author",
		Version:     "1.0.0",
		Type:        "lua",
		Entry:       "game.lua",
		Category:    "game",
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal manifest: %v", err)
	}

	err = os.WriteFile(manifestPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	// Test loading manifest
	loaded, err := loader.loadManifest(manifestPath)
	if err != nil {
		t.Fatalf("loadManifest() failed: %v", err)
	}

	if loaded.ID != manifest.ID {
		t.Errorf("Expected ID %s, got %s", manifest.ID, loaded.ID)
	}
	if loaded.Name != manifest.Name {
		t.Errorf("Expected Name %s, got %s", manifest.Name, loaded.Name)
	}
	if loaded.Type != manifest.Type {
		t.Errorf("Expected Type %s, got %s", manifest.Type, loaded.Type)
	}
}

func TestLoadManifestError(t *testing.T) {
	loader := NewLoader()

	// Test with non-existent file
	_, err := loader.loadManifest("/nonexistent/manifest.json")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Test with invalid JSON
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "invalid.json")
	err = os.WriteFile(invalidPath, []byte("not valid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid file: %v", err)
	}

	_, err = loader.loadManifest(invalidPath)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestScan(t *testing.T) {
	loader := NewLoader()

	// Create temporary games directory
	tmpDir := t.TempDir()
	loader.gamesDir = tmpDir

	// Create test plugin directories
	games := []struct {
		id   string
		name string
	}{
		{"game1", "Game 1"},
		{"game2", "Game 2"},
	}

	for _, g := range games {
		gameDir := filepath.Join(tmpDir, g.id)
		err := os.MkdirAll(gameDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create game directory: %v", err)
		}

		manifest := Manifest{
			ID:   g.id,
			Name: g.name,
			Type: "lua",
		}

		data, _ := json.Marshal(manifest)
		manifestPath := filepath.Join(gameDir, "manifest.json")
		err = os.WriteFile(manifestPath, data, 0644)
		if err != nil {
			t.Fatalf("Failed to write manifest: %v", err)
		}
	}

	// Test scanning
	err := loader.Scan()
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if len(loader.manifests) != 2 {
		t.Errorf("Expected 2 manifests, got %d", len(loader.manifests))
	}

	// Verify manifests were loaded
	for _, g := range games {
		if _, ok := loader.manifests[g.id]; !ok {
			t.Errorf("Manifest for %s not loaded", g.id)
		}
	}
}

func TestScanIgnoresNonDirectories(t *testing.T) {
	loader := NewLoader()

	tmpDir := t.TempDir()
	loader.gamesDir = tmpDir

	// Create a file (not a directory)
	filePath := filepath.Join(tmpDir, "not-a-game.txt")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Scan should not error and should not load the file
	err = loader.Scan()
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if len(loader.manifests) != 0 {
		t.Errorf("Expected 0 manifests, got %d", len(loader.manifests))
	}
}

func TestScanIgnoresDirectoriesWithoutManifest(t *testing.T) {
	loader := NewLoader()

	tmpDir := t.TempDir()
	loader.gamesDir = tmpDir

	// Create directory without manifest
	gameDir := filepath.Join(tmpDir, "incomplete-game")
	err := os.MkdirAll(gameDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Scan should not error
	err = loader.Scan()
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}

	if len(loader.manifests) != 0 {
		t.Errorf("Expected 0 manifests, got %d", len(loader.manifests))
	}
}

func TestList(t *testing.T) {
	loader := NewLoader()

	// Add some manifests
	loader.manifests["game1"] = &Manifest{ID: "game1", Name: "Game 1"}
	loader.manifests["game2"] = &Manifest{ID: "game2", Name: "Game 2"}
	loader.manifests["game3"] = &Manifest{ID: "game3", Name: "Game 3"}

	list := loader.List()
	if len(list) != 3 {
		t.Errorf("Expected 3 manifests, got %d", len(list))
	}
}

func TestListEmpty(t *testing.T) {
	loader := NewLoader()

	list := loader.List()
	// List() may return nil or empty slice for empty manifests
	if list != nil && len(list) != 0 {
		t.Errorf("Expected 0 manifests, got %d", len(list))
	}
}

func TestGet(t *testing.T) {
	loader := NewLoader()

	// Add a manifest
	expected := &Manifest{ID: "test-game", Name: "Test Game"}
	loader.manifests["test-game"] = expected

	// Test getting existing manifest
	manifest, ok := loader.Get("test-game")
	if !ok {
		t.Error("Expected to find manifest")
	}
	if manifest.ID != expected.ID {
		t.Errorf("Expected ID %s, got %s", expected.ID, manifest.ID)
	}

	// Test getting non-existent manifest
	_, ok = loader.Get("nonexistent")
	if ok {
		t.Error("Expected not to find non-existent manifest")
	}
}

func TestLoadPluginNotFound(t *testing.T) {
	loader := NewLoader()

	_, err := loader.Load("nonexistent")
	if err == nil {
		t.Error("Expected error when loading non-existent plugin")
	}
}

func TestLoadUnsupportedType(t *testing.T) {
	loader := NewLoader()

	// Add manifest with unsupported type
	loader.manifests["test"] = &Manifest{
		ID:   "test",
		Type: "wasm",
	}

	_, err := loader.Load("test")
	if err == nil {
		t.Error("Expected error for unsupported plugin type")
	}

	// Test unknown type
	loader.manifests["test2"] = &Manifest{
		ID:   "test2",
		Type: "unknown",
	}

	_, err = loader.Load("test2")
	if err == nil {
		t.Error("Expected error for unknown plugin type")
	}
}

func TestUninstall(t *testing.T) {
	loader := NewLoader()

	tmpDir := t.TempDir()
	loader.gamesDir = tmpDir

	// Create test plugin
	gameID := "test-game"
	gameDir := filepath.Join(tmpDir, gameID)
	err := os.MkdirAll(gameDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create game directory: %v", err)
	}

	// Create a file in the plugin directory
	testFile := filepath.Join(gameDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add to manifests
	loader.manifests[gameID] = &Manifest{ID: gameID, Name: "Test Game"}

	// Test uninstall
	err = loader.Uninstall(gameID)
	if err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}

	// Verify directory was removed
	if _, err := os.Stat(gameDir); !os.IsNotExist(err) {
		t.Error("Plugin directory should be removed")
	}

	// Verify manifest was removed
	if _, ok := loader.manifests[gameID]; ok {
		t.Error("Manifest should be removed from map")
	}
}

func TestUninstallNotFound(t *testing.T) {
	loader := NewLoader()

	err := loader.Uninstall("nonexistent")
	if err == nil {
		t.Error("Expected error when uninstalling non-existent plugin")
	}
}

func TestInstall(t *testing.T) {
	loader := NewLoader()

	// Install is not yet implemented
	err := loader.Install("test-source")
	if err == nil {
		t.Error("Expected error for unimplemented Install()")
	}
}

func TestManifestFields(t *testing.T) {
	manifest := Manifest{
		ID:          "test",
		Name:        "Test",
		Description: "Test Description",
		Author:      "Test Author",
		Version:     "1.0.0",
		Type:        "lua",
		Entry:       "main.lua",
		MinWidth:    40,
		MinHeight:   20,
		Category:    "game",
		Assets:      []string{"asset1.png", "asset2.png"},
		Platforms: map[string]string{
			"linux":   "linux-binary",
			"windows": "windows.exe",
		},
	}

	// Verify all fields are accessible
	if manifest.ID != "test" {
		t.Error("ID field not set correctly")
	}
	if manifest.MinWidth != 40 {
		t.Error("MinWidth field not set correctly")
	}
	if manifest.MinHeight != 20 {
		t.Error("MinHeight field not set correctly")
	}
	if len(manifest.Assets) != 2 {
		t.Error("Assets field not set correctly")
	}
	if len(manifest.Platforms) != 2 {
		t.Error("Platforms field not set correctly")
	}
}
