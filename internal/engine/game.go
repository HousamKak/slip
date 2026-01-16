package engine

// Game is the interface that all games must implement
type Game interface {
	// Info returns metadata about the game
	Info() GameInfo

	// Init initializes the game with the given context
	Init(ctx *GameContext) error

	// Start is called when the game starts
	Start() error

	// Stop is called when the game stops
	Stop() error

	// Update is called every frame with delta time in seconds
	Update(dt float64) error

	// Render draws the game to the screen
	Render(screen *Screen) error

	// HandleInput processes input events
	HandleInput(input Input) error
}

// GameInfo contains metadata about a game
type GameInfo struct {
	ID          string
	Name        string
	Description string
	Author      string
	Version     string
	MinWidth    int
	MinHeight   int
	Category    string // "game", "animation", "utility"
}

// GameContext provides context to games
type GameContext struct {
	Width        int
	Height       int
	FPS          int
	StateDir     string
	PlayerName   string
	OnExit       func()
	OnScore      func(score int)
	GetConfig    func(key string) string
	GetHighScore func() int
	IsHighScore  func(score int) bool
}

// Animation is a simplified interface for animations (non-interactive)
type Animation interface {
	// Info returns metadata about the animation
	Info() GameInfo

	// Init initializes the animation
	Init(ctx *GameContext) error

	// Update updates the animation state
	Update(dt float64) error

	// Render draws the animation
	Render(screen *Screen) error
}

// AnimationWrapper wraps an Animation to implement Game interface
type AnimationWrapper struct {
	anim Animation
}

// WrapAnimation wraps an Animation as a Game
func WrapAnimation(anim Animation) Game {
	return &AnimationWrapper{anim: anim}
}

func (w *AnimationWrapper) Info() GameInfo {
	return w.anim.Info()
}

func (w *AnimationWrapper) Init(ctx *GameContext) error {
	return w.anim.Init(ctx)
}

func (w *AnimationWrapper) Start() error {
	return nil
}

func (w *AnimationWrapper) Stop() error {
	return nil
}

func (w *AnimationWrapper) Update(dt float64) error {
	return w.anim.Update(dt)
}

func (w *AnimationWrapper) Render(screen *Screen) error {
	return w.anim.Render(screen)
}

func (w *AnimationWrapper) HandleInput(input Input) error {
	// Animations don't handle input (except for quit which is handled by engine)
	return nil
}
