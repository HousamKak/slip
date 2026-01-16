package animations

import (
	"slip/internal/engine"
	"slip/internal/plugins"
)

// AnimationFactory is a function that creates a new animation instance
type AnimationFactory func() engine.Animation

// Registry holds all registered animations
type Registry struct {
	animations map[string]AnimationFactory
	order      []string
	loader     *plugins.Loader
}

// NewRegistry creates a new animation registry with built-in animations
func NewRegistry() *Registry {
	r := &Registry{
		animations: make(map[string]AnimationFactory),
		loader:     plugins.NewLoader(),
	}

	// Register built-in animations
	r.Register("clipster", func() engine.Animation { return NewClipster() })
	r.Register("starfield", func() engine.Animation { return NewStarfield() })
	r.Register("matrix", func() engine.Animation { return NewMatrix() })
	r.Register("dvd", func() engine.Animation { return NewDVD() })
	r.Register("fire", func() engine.Animation { return NewFire() })

	// Load plugins
	r.loadPlugins()

	return r
}

// loadPlugins scans and loads plugin animations
func (r *Registry) loadPlugins() {
	if err := r.loader.Scan(); err != nil {
		// Silently fail - plugins are optional
		return
	}

	for _, manifest := range r.loader.List() {
		// Only load animations, not games
		if manifest.Category != "animation" {
			continue
		}

		// Create a closure to capture the manifest ID
		id := manifest.ID
		r.Register(id, func() engine.Animation {
			game, err := r.loader.Load(id)
			if err != nil || game == nil {
				return nil
			}
			// If plugin implements Animation interface directly, use it
			// Otherwise wrap the game
			// For now, plugins are games, so we can't use them as animations directly
			// This will be enhanced when we add proper animation plugin support
			return nil
		})
	}
}

// Register adds an animation to the registry
func (r *Registry) Register(id string, factory AnimationFactory) {
	if _, exists := r.animations[id]; !exists {
		r.order = append(r.order, id)
	}
	r.animations[id] = factory
}

// Get returns an animation factory by ID
func (r *Registry) Get(id string) (AnimationFactory, bool) {
	factory, ok := r.animations[id]
	return factory, ok
}

// Create creates a new instance of an animation
func (r *Registry) Create(id string) (engine.Animation, bool) {
	factory, ok := r.animations[id]
	if !ok {
		return nil, false
	}
	return factory(), true
}

// CreateAsGame creates an animation wrapped as a Game
func (r *Registry) CreateAsGame(id string) (engine.Game, bool) {
	anim, ok := r.Create(id)
	if !ok {
		return nil, false
	}
	return engine.WrapAnimation(anim), true
}

// List returns info about all registered animations
func (r *Registry) List() []engine.GameInfo {
	var infos []engine.GameInfo
	for _, id := range r.order {
		if factory, ok := r.animations[id]; ok {
			anim := factory()
			infos = append(infos, anim.Info())
		}
	}
	return infos
}

// IDs returns all registered animation IDs in order
func (r *Registry) IDs() []string {
	return r.order
}

// Count returns the number of registered animations
func (r *Registry) Count() int {
	return len(r.animations)
}

// RegisterPlugin registers a plugin animation
func (r *Registry) RegisterPlugin(manifest *plugins.Manifest, loader *plugins.Loader) {
	id := manifest.ID
	r.Register(id, func() engine.Animation {
		game, err := loader.Load(id)
		if err != nil || game == nil {
			return nil
		}
		// Plugins loaded as games can't be used as animations directly yet
		// This will be enhanced when we add proper animation plugin support
		return nil
	})
}
