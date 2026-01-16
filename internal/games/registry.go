package games

import (
	"slip/internal/engine"
	"slip/internal/games/blockfall"
	"slip/internal/games/breakout"
	"slip/internal/games/pong"
	"slip/internal/games/snake"
	"slip/internal/plugins"
)

// GameFactory is a function that creates a new game instance
type GameFactory func() engine.Game

// Registry holds all registered games
type Registry struct {
	games  map[string]GameFactory
	order  []string // Maintains insertion order
	loader *plugins.Loader
}

// NewRegistry creates a new game registry with built-in games
func NewRegistry() *Registry {
	r := &Registry{
		games:  make(map[string]GameFactory),
		loader: plugins.NewLoader(),
	}

	// Register built-in games
	r.Register("snake", func() engine.Game { return snake.New() })
	r.Register("pong", func() engine.Game { return pong.New() })
	r.Register("breakout", func() engine.Game { return breakout.New() })
	r.Register("blockfall", func() engine.Game { return blockfall.New() })

	// Load plugins
	r.loadPlugins()

	return r
}

// loadPlugins scans and loads plugin games
func (r *Registry) loadPlugins() {
	if err := r.loader.Scan(); err != nil {
		// Silently fail - plugins are optional
		return
	}

	for _, manifest := range r.loader.List() {
		// Only load games, not animations
		if manifest.Category != "game" {
			continue
		}

		// Create a closure to capture the manifest ID
		id := manifest.ID
		r.Register(id, func() engine.Game {
			game, err := r.loader.Load(id)
			if err != nil {
				// Return a dummy game that shows error
				return nil
			}
			return game
		})
	}
}

// Register adds a game to the registry
func (r *Registry) Register(id string, factory GameFactory) {
	if _, exists := r.games[id]; !exists {
		r.order = append(r.order, id)
	}
	r.games[id] = factory
}

// Get returns a game factory by ID
func (r *Registry) Get(id string) (GameFactory, bool) {
	factory, ok := r.games[id]
	return factory, ok
}

// Create creates a new instance of a game
func (r *Registry) Create(id string) (engine.Game, bool) {
	factory, ok := r.games[id]
	if !ok {
		return nil, false
	}
	return factory(), true
}

// List returns info about all registered games
func (r *Registry) List() []engine.GameInfo {
	var infos []engine.GameInfo
	for _, id := range r.order {
		if factory, ok := r.games[id]; ok {
			game := factory()
			infos = append(infos, game.Info())
		}
	}
	return infos
}

// IDs returns all registered game IDs in order
func (r *Registry) IDs() []string {
	return r.order
}

// Count returns the number of registered games
func (r *Registry) Count() int {
	return len(r.games)
}

// RegisterPlugin registers a plugin game
func (r *Registry) RegisterPlugin(manifest *plugins.Manifest, loader *plugins.Loader) {
	id := manifest.ID
	r.Register(id, func() engine.Game {
		game, err := loader.Load(id)
		if err != nil {
			return nil
		}
		return game
	})
}
