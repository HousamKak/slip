# Slip: Terminal Entertainment Platform
## Comprehensive Vision and Implementation Document

---

## Executive Summary

**Slip** is a terminal-based entertainment platform that transforms idle terminal time into an engaging experience. While developers wait for AI assistants (Claude Code, Codex, Aider), builds, or long-running processes, Slip provides a dedicated terminal space for animations and games.

The project evolved from **Clipster** (an ASCII paperclip animation POC) into a full-fledged "mini operating system" for terminal entertainment, featuring:

- **Persistent companion pane** alongside your work terminal
- **Animated ASCII art** (mascots, visualizations, idle animations)
- **Full-fledged terminal games** (Snake, Pong, Tetris-style, etc.)
- **Navigation menu/home screen** for switching between content
- **Remote game installation** from a curated store/registry

---

## Part 1: Foundation - Current State Analysis

### 1.1 Clipster POC Capabilities

The existing POC demonstrates:

| Component | Status | Description |
|-----------|--------|-------------|
| tmux Integration | Working | Split panes, pane tagging, summon/hide/dock |
| Windows Terminal | Working | Fallback via `wt.exe` for non-WSL Windows |
| ASCII Renderer | Working | Fixed-size framebuffer, ANSI cursor control |
| Image-to-ASCII | Working | PNG/JPG/GIF conversion with density mapping |
| Speech Bubbles | Working | Overlay text with auto-wrapping |
| State Management | Working | JSON-based heartbeat and messaging |

### 1.2 Architecture Overview (Current)

```
clipster (binary)
├── cmd/clipster/main.go        # Entry point
├── internal/
│   ├── cli/cli.go              # Command dispatch
│   ├── render/
│   │   ├── render.go           # Frame rendering
│   │   ├── bubble.go           # Speech bubbles
│   │   ├── art.go              # Frame loading
│   │   ├── image.go            # Image-to-ASCII
│   │   └── clipster.go         # Built-in frames
│   ├── tmux/tmux.go            # tmux pane management
│   ├── term/term.go            # ANSI sequences
│   ├── platform/platform.go    # OS detection
│   └── state/state.go          # Persistent state
```

### 1.3 Key Design Principles (Preserved)

1. **Pane isolation** - Never interfere with the work terminal
2. **Single binary** - No runtime dependencies (Go cross-compilation)
3. **tmux as window manager** - Proven, reliable, cross-platform (via WSL)
4. **Stateless by default** - Environment variables for configuration

---

## Part 2: The New Vision - Slip

### 2.1 Conceptual Model

```
┌─────────────────────────────────────┬──────────────────────┐
│                                     │                      │
│  PANE A (Work)                      │  PANE B (Slip)       │
│                                     │                      │
│  $ claude code                      │  ╔══════════════════╗│
│  > Analyzing codebase...            │  ║  SLIP v1.0       ║│
│  > Found 15 files...                │  ╠══════════════════╣│
│  > Proposing changes...             │  ║  > Snake         ║│
│                                     │  ║    Pong          ║│
│                                     │  ║    Animations    ║│
│                                     │  ║    Settings      ║│
│                                     │  ╚══════════════════╝│
│                                     │                      │
└─────────────────────────────────────┴──────────────────────┘
```

### 2.2 User Experience Flow

```
1. User starts tmux session
2. User runs: slip summon
3. Slip pane appears with home screen
4. User selects "Snake" from menu
5. Snake game starts in Slip pane
6. User runs "claude code" in work pane
7. User plays Snake while waiting for Claude
8. Press ESC to return to Slip home menu
9. Switch to "Animations" → watch starfield
10. slip hide when done
```

### 2.3 Core Features

#### A. Home Screen / Navigation Shell

```
╔════════════════════════════════════╗
║           S L I P                  ║
║     Terminal Entertainment         ║
╠════════════════════════════════════╣
║                                    ║
║   [G] Games                        ║
║   [A] Animations                   ║
║   [S] Store                        ║
║   [C] Config                       ║
║   [Q] Quit                         ║
║                                    ║
╚════════════════════════════════════╝
```

#### B. Games Subsystem

- Built-in games: Snake, Pong, Breakout, Tetris-clone
- Game loop with fixed timestep
- Input handling (arrow keys, WASD)
- High score persistence
- Pause/resume functionality

#### C. Animations Subsystem

- Clipster (existing paperclip buddy)
- Starfield (parallax star effect)
- Matrix rain
- Bouncing DVD logo
- Custom ASCII art files
- Image/GIF to ASCII

#### D. Store / Remote Installation

- Fetch game manifests from registry server
- Download game bundles (WASM or native plugins)
- Local game catalog management
- Version checking and updates

---

## Part 3: Technical Architecture

### 3.1 Proposed Directory Structure

```
slip/
├── cmd/slip/main.go
├── internal/
│   ├── app/                    # Application shell
│   │   ├── app.go              # Main TUI application
│   │   ├── home.go             # Home screen
│   │   ├── menu.go             # Menu navigation
│   │   └── router.go           # Screen routing
│   │
│   ├── engine/                 # Game engine core
│   │   ├── engine.go           # Main game loop
│   │   ├── input.go            # Input handling
│   │   ├── render.go           # Frame rendering
│   │   ├── collision.go        # Collision detection
│   │   ├── entity.go           # Entity system
│   │   └── audio.go            # Terminal bell/beep
│   │
│   ├── games/                  # Built-in games
│   │   ├── snake/
│   │   │   ├── snake.go
│   │   │   └── snake_test.go
│   │   ├── pong/
│   │   │   └── pong.go
│   │   ├── breakout/
│   │   │   └── breakout.go
│   │   └── registry.go         # Game registration
│   │
│   ├── animations/             # Animation content
│   │   ├── clipster.go
│   │   ├── starfield.go
│   │   ├── matrix.go
│   │   ├── dvd.go
│   │   └── registry.go
│   │
│   ├── store/                  # Remote content
│   │   ├── client.go           # HTTP client
│   │   ├── manifest.go         # Game manifests
│   │   ├── download.go         # Download manager
│   │   └── catalog.go          # Local catalog
│   │
│   ├── plugins/                # Plugin system
│   │   ├── loader.go           # Plugin loading
│   │   ├── interface.go        # Game interface
│   │   └── sandbox.go          # Security sandbox
│   │
│   ├── tmux/                   # tmux integration (existing)
│   ├── term/                   # Terminal helpers (existing)
│   ├── platform/               # OS detection (existing)
│   ├── state/                  # Persistent state
│   │   ├── state.go
│   │   ├── scores.go           # High scores
│   │   └── config.go           # User config
│   │
│   └── ui/                     # UI components
│       ├── box.go              # Box drawing
│       ├── menu.go             # Menu widget
│       ├── text.go             # Text rendering
│       └── styles.go           # Color/style definitions
│
├── games/                      # External game data
│   └── .gitkeep
│
├── go.mod
├── go.sum
└── README.md
```

### 3.2 Game Engine Design

#### 3.2.1 Game Interface

All games (built-in and plugins) implement a common interface:

```go
package engine

type Game interface {
    // Lifecycle
    Init(ctx GameContext) error
    Start() error
    Stop() error

    // Game loop (called at fixed timestep)
    Update(dt float64) error
    Render(screen *Screen) error

    // Input handling
    HandleInput(input Input) error

    // Metadata
    Info() GameInfo
}

type GameInfo struct {
    ID          string
    Name        string
    Description string
    Author      string
    Version     string
    MinWidth    int
    MinHeight   int
}

type GameContext struct {
    Width      int
    Height     int
    FPS        int
    StateDir   string
    OnExit     func()
}
```

#### 3.2.2 Game Loop Architecture

```go
package engine

type Engine struct {
    screen      *Screen
    currentGame Game
    running     bool
    targetFPS   int
}

func (e *Engine) Run() error {
    ticker := time.NewTicker(time.Second / time.Duration(e.targetFPS))
    defer ticker.Stop()

    lastTime := time.Now()

    for e.running {
        select {
        case input := <-e.inputChan:
            if input.IsQuit() {
                return nil
            }
            e.currentGame.HandleInput(input)

        case <-ticker.C:
            now := time.Now()
            dt := now.Sub(lastTime).Seconds()
            lastTime = now

            // Update game state
            e.currentGame.Update(dt)

            // Render frame
            e.screen.Clear()
            e.currentGame.Render(e.screen)
            e.screen.Flush()
        }
    }
    return nil
}
```

#### 3.2.3 Screen Buffer

```go
package engine

type Screen struct {
    width    int
    height   int
    buffer   [][]Cell
    styles   [][]Style
    dirty    bool
}

type Cell struct {
    Rune  rune
    Style Style
}

type Style struct {
    FG        Color
    BG        Color
    Bold      bool
    Underline bool
}

func (s *Screen) Set(x, y int, r rune, style Style) {
    if x >= 0 && x < s.width && y >= 0 && y < s.height {
        s.buffer[y][x] = Cell{Rune: r, Style: style}
        s.dirty = true
    }
}

func (s *Screen) DrawText(x, y int, text string, style Style) {
    for i, r := range text {
        s.Set(x+i, y, r, style)
    }
}

func (s *Screen) DrawBox(x, y, w, h int, style BoxStyle) {
    // Draw box characters: ┌─┐│└┘ or ╔═╗║╚╝
}

func (s *Screen) Flush() {
    if !s.dirty {
        return
    }
    // Build single output string and write to stdout
    // Use ANSI sequences for positioning and styling
}
```

### 3.3 Input System

```go
package engine

type Input struct {
    Type InputType
    Key  Key
    Rune rune
}

type InputType int

const (
    InputKeyPress InputType = iota
    InputKeyRelease
    InputResize
)

type Key int

const (
    KeyUp Key = iota
    KeyDown
    KeyLeft
    KeyRight
    KeyEnter
    KeyEscape
    KeySpace
    KeyTab
    KeyBackspace
    KeyRune  // For regular characters
)

func StartInputListener(ch chan<- Input) {
    // Read from stdin in raw mode
    // Parse ANSI escape sequences for arrow keys
    // Send Input events to channel
}
```

### 3.4 Plugin System for Remote Games

#### 3.4.1 Plugin Manifest Format

```json
{
    "id": "snake-pro",
    "name": "Snake Pro",
    "version": "1.0.0",
    "author": "GameDev",
    "description": "Enhanced snake with power-ups",
    "min_slip_version": "1.0.0",
    "type": "native",
    "entry": "snake-pro",
    "assets": [
        "levels.json",
        "sprites.txt"
    ],
    "platforms": {
        "linux/amd64": "snake-pro-linux-amd64",
        "darwin/amd64": "snake-pro-darwin-amd64",
        "darwin/arm64": "snake-pro-darwin-arm64",
        "windows/amd64": "snake-pro-windows-amd64.exe"
    }
}
```

#### 3.4.2 Plugin Types

**Option A: Native Binary Plugins**
- Downloaded as platform-specific executables
- Communicate via stdin/stdout protocol (JSON-RPC)
- Maximum performance, most work to distribute

**Option B: WASM Plugins (Recommended for v2)**
- Single cross-platform binary
- Run in sandboxed WASM runtime (wasmtime/wazero)
- Safer, easier distribution
- Slight performance overhead

**Option C: Lua Scripts (Simplest)**
- Game logic in Lua files
- Interpreted by embedded Lua runtime (gopher-lua)
- Very easy to create/distribute
- Limited performance for complex games

#### 3.4.3 Plugin Loading

```go
package plugins

type PluginLoader struct {
    gamesDir string
}

func (l *PluginLoader) LoadGame(id string) (engine.Game, error) {
    manifest, err := l.loadManifest(id)
    if err != nil {
        return nil, err
    }

    switch manifest.Type {
    case "native":
        return l.loadNativePlugin(manifest)
    case "wasm":
        return l.loadWasmPlugin(manifest)
    case "lua":
        return l.loadLuaPlugin(manifest)
    default:
        return nil, fmt.Errorf("unknown plugin type: %s", manifest.Type)
    }
}
```

### 3.5 Store / Registry System

#### 3.5.1 Registry API

```
GET  /api/v1/games              # List all games
GET  /api/v1/games/{id}         # Get game details
GET  /api/v1/games/{id}/download?platform=linux/amd64
POST /api/v1/games/{id}/report  # Report issues
```

#### 3.5.2 Local Catalog

```go
package store

type Catalog struct {
    path   string
    games  map[string]InstalledGame
}

type InstalledGame struct {
    ID          string
    Name        string
    Version     string
    InstalledAt time.Time
    Path        string
    Manifest    Manifest
}

func (c *Catalog) Install(id string) error {
    // 1. Fetch manifest from registry
    // 2. Download binary for current platform
    // 3. Verify checksum
    // 4. Extract to games directory
    // 5. Update catalog
}

func (c *Catalog) Uninstall(id string) error {
    // Remove from filesystem and catalog
}

func (c *Catalog) CheckUpdates() ([]Update, error) {
    // Compare installed versions with registry
}
```

### 3.6 State Management

#### 3.6.1 High Scores

```go
package state

type ScoreBoard struct {
    path   string
    scores map[string][]Score
}

type Score struct {
    Value     int
    Player    string
    Timestamp time.Time
    GameID    string
}

func (s *ScoreBoard) AddScore(gameID string, score Score) error {
    // Add to sorted list, keep top 10
}

func (s *ScoreBoard) GetTopScores(gameID string, limit int) []Score {
    // Return top N scores for game
}
```

#### 3.6.2 User Configuration

```go
package state

type Config struct {
    // Display
    DefaultDock    string  `json:"default_dock"`     // "right", "bottom", "left"
    DefaultSize    string  `json:"default_size"`     // "40x20"
    Theme          string  `json:"theme"`            // "default", "retro", "minimal"

    // Games
    DefaultFPS     int     `json:"default_fps"`
    ShowFPSCounter bool    `json:"show_fps_counter"`

    // Store
    RegistryURL    string  `json:"registry_url"`
    AutoUpdate     bool    `json:"auto_update"`

    // Player
    PlayerName     string  `json:"player_name"`
}
```

---

## Part 4: Built-in Games Specification

### 4.1 Snake

```
┌──────────────────────────────────────┐
│ SNAKE                    Score: 42   │
├──────────────────────────────────────┤
│                                      │
│                                      │
│        ████                          │
│           █                          │
│           █                          │
│           █                          │
│                              *       │
│                                      │
│                                      │
├──────────────────────────────────────┤
│ [←↑↓→] Move   [P] Pause   [Q] Quit   │
└──────────────────────────────────────┘
```

**Features:**
- Classic snake gameplay
- Growing snake on food collection
- Wall collision = game over
- Self collision = game over
- Speed increases with score
- High score tracking

### 4.2 Pong

```
┌──────────────────────────────────────┐
│ PONG              3 : 2              │
├──────────────────────────────────────┤
│                                      │
│  █                              █    │
│  █                              █    │
│  █                              █    │
│  █            o                 █    │
│  █                              █    │
│  █                              █    │
│  █                              █    │
│                                      │
├──────────────────────────────────────┤
│ [W/S] Move   [P] Pause   [Q] Quit    │
└──────────────────────────────────────┘
```

**Features:**
- Single player vs AI
- Paddle physics
- Ball speed increases over rally
- First to 5 points wins

### 4.3 Breakout

```
┌──────────────────────────────────────┐
│ BREAKOUT          Score: 150  Lives: 3│
├──────────────────────────────────────┤
│  ████████████████████████████████    │
│  ████████████████████████████████    │
│  ████████████████████████████████    │
│  ████████        ████████████████    │
│                                      │
│                                      │
│                  o                   │
│                                      │
│                                      │
│              ═══════                 │
├──────────────────────────────────────┤
│ [←→] Move   [SPACE] Launch   [Q] Quit│
└──────────────────────────────────────┘
```

**Features:**
- Brick-breaking gameplay
- Multiple brick types (1-hit, 2-hit, 3-hit)
- Power-ups (wider paddle, multi-ball, slow ball)
- Level progression

### 4.4 Tetris-style (Blockfall)

```
┌──────────────────────────────────────┐
│ BLOCKFALL          Score: 2400       │
├──────────────────────────────────────┤
│    │          │    Next:             │
│    │    ██    │    ▀▀                │
│    │    ██    │    ▀▀                │
│    │          │                      │
│    │          │    Level: 4          │
│    │          │    Lines: 12         │
│    │          │                      │
│    │   █      │                      │
│    │  ███   ██│                      │
│    │█████████ │                      │
│    └──────────┘                      │
├──────────────────────────────────────┤
│ [←→] Move [↑] Rotate [↓] Drop [Q]Quit│
└──────────────────────────────────────┘
```

**Features:**
- 7 standard tetromino shapes
- Rotation system
- Line clearing with scoring
- Level progression (increased speed)
- Next piece preview

---

## Part 5: Built-in Animations Specification

### 5.1 Clipster (Existing)

- Animated paperclip buddy
- Multiple states: idle, wave, think
- Speech bubble overlay
- Blink animation

### 5.2 Starfield

```
             ·
    ·              *
         ·                   ·
    *          ·
                      ·
  ·     ·                    *
              ·       ·
       *                     ·
```

**Features:**
- Parallax effect (multiple star layers)
- Stars move from right to left
- Different star sizes: · * ✦
- Configurable speed

### 5.3 Matrix Rain

```
    1   ケ       0   ト
    ム   0       ロ   0   ハ
    0   ス   1   0       ト   1
    キ       ム   ネ   0       0
    0   1   0       フ   1   0
        0   キ   0       0
    0           0   0
```

**Features:**
- Falling character columns
- Green color scheme
- Random characters (katakana, numbers)
- Variable column speeds

### 5.4 DVD Bounce

```


         ╔══════╗
         ║ DVD  ║
         ╚══════╝


```

**Features:**
- Logo bounces off screen edges
- Color changes on corner hit
- Celebrate when hitting exact corner

### 5.5 Fire Effect

```
    ░░▒▒▓▓██▓▓▒▒░░
  ░░▒▒▓▓████████▓▓▒▒░░
░░▒▒▓▓██████████████▓▓▒▒░░
▒▒▓▓████████████████████▓▓▒▒
▓▓██████████████████████████▓▓
████████████████████████████████
```

**Features:**
- Classic demoscene fire
- Palette cycling
- Procedural flame generation

---

## Part 6: CLI Command Design

### 6.1 Primary Commands

```bash
# Summon Slip in a tmux pane
slip summon [--dock right|bottom|left] [--size 40x20]

# Run Slip directly (takes over current terminal)
slip run [--width 40] [--height 20]

# Hide Slip pane
slip hide

# Reposition Slip pane
slip dock [--dock right|bottom|left] [--size 40x20]

# Check status
slip status
```

### 6.2 Game Commands

```bash
# Launch a specific game directly
slip play snake
slip play pong
slip play breakout

# List available games
slip games list

# Show high scores
slip scores [game]
slip scores snake
```

### 6.3 Animation Commands

```bash
# Run a specific animation
slip animate clipster
slip animate starfield
slip animate matrix

# List available animations
slip animations list
```

### 6.4 Store Commands

```bash
# Browse available games
slip store browse

# Search for games
slip store search "puzzle"

# Install a game
slip store install snake-pro

# Uninstall a game
slip store uninstall snake-pro

# Update installed games
slip store update

# Show installed games
slip store installed
```

### 6.5 Configuration Commands

```bash
# Show current config
slip config show

# Set config value
slip config set player_name "Player1"
slip config set theme "retro"
slip config set default_dock "bottom"

# Reset to defaults
slip config reset
```

---

## Part 7: Implementation Phases

### Phase 0: Foundation Refactoring (1-2 weeks)

**Goal:** Refactor Clipster into Slip architecture

**Tasks:**
1. Rename project from `clipster` to `slip`
2. Introduce Bubble Tea for proper TUI management
3. Create basic home screen with menu navigation
4. Extract existing animation into `animations/clipster.go`
5. Test tmux integration continues to work

**Deliverables:**
- `slip summon` opens home screen
- Can navigate menu with arrow keys
- Can select "Clipster" animation
- Press ESC returns to home

### Phase 1: Game Engine Core (2-3 weeks)

**Goal:** Build the foundational game engine

**Tasks:**
1. Implement `engine.Screen` with double-buffering
2. Implement `engine.Input` system with raw terminal mode
3. Implement `engine.Engine` game loop
4. Create `engine.Game` interface
5. Build Snake as first game implementation

**Deliverables:**
- Working Snake game
- Smooth rendering at 30 FPS
- Input handling for arrow keys
- Score display
- Game over detection

### Phase 2: Expand Content (2-3 weeks)

**Goal:** Add more games and animations

**Tasks:**
1. Implement Pong (single player vs AI)
2. Implement Breakout
3. Implement starfield animation
4. Implement matrix rain animation
5. Add high score system

**Deliverables:**
- 4 playable games
- 4 animations (clipster, starfield, matrix, dvd)
- Persistent high scores

### Phase 3: Polish and UX (1-2 weeks)

**Goal:** Make it feel like a product

**Tasks:**
1. Add visual themes (default, retro, minimal)
2. Implement help screens
3. Add pause functionality to all games
4. Handle terminal resize gracefully
5. Add configuration system
6. Improve menu aesthetics with box drawing

**Deliverables:**
- Polished UI
- Theme support
- Comprehensive help
- Graceful error handling

### Phase 4: Plugin System (2-3 weeks)

**Goal:** Enable third-party games

**Tasks:**
1. Define plugin manifest format
2. Implement Lua scripting runtime (gopher-lua)
3. Create Lua game API
4. Build plugin loader
5. Create example Lua game

**Deliverables:**
- Lua-based games work
- Documentation for game creators
- Example game source code

### Phase 5: Store / Registry (2-3 weeks)

**Goal:** Remote game installation

**Tasks:**
1. Build simple registry server (can be static JSON on GitHub)
2. Implement store client in Slip
3. Add browse/search/install commands
4. Implement update checking
5. Add download progress display

**Deliverables:**
- `slip store browse` shows available games
- `slip store install <id>` downloads and installs
- Installed games appear in games menu

### Phase 6: Distribution and Launch (1-2 weeks)

**Goal:** Ship it

**Tasks:**
1. Build binaries for all platforms
2. Create GitHub releases
3. Set up Gumroad/LemonSqueezy page
4. Write README and documentation
5. Create demo GIFs/videos
6. Announce on social media

---

## Part 8: Registry Server Design

### 8.1 Minimal Registry (v1)

For v1, the registry can be a static JSON file hosted on GitHub:

**games.json**
```json
{
    "games": [
        {
            "id": "snake-pro",
            "name": "Snake Pro",
            "description": "Enhanced snake with power-ups",
            "author": "GameDev",
            "version": "1.0.0",
            "downloads": {
                "linux/amd64": "https://releases.example.com/snake-pro/1.0.0/snake-pro-linux-amd64.tar.gz",
                "darwin/arm64": "https://releases.example.com/snake-pro/1.0.0/snake-pro-darwin-arm64.tar.gz",
                "windows/amd64": "https://releases.example.com/snake-pro/1.0.0/snake-pro-windows-amd64.zip"
            },
            "checksum": {
                "linux/amd64": "sha256:abc123...",
                "darwin/arm64": "sha256:def456...",
                "windows/amd64": "sha256:ghi789..."
            },
            "type": "lua",
            "min_slip_version": "1.0.0"
        }
    ],
    "featured": ["snake-pro"],
    "updated_at": "2025-01-15T00:00:00Z"
}
```

### 8.2 Full Registry Server (v2+)

For scale, implement a proper API:

```
POST   /api/v1/games                # Submit new game
GET    /api/v1/games                # List games
GET    /api/v1/games/:id            # Get game details
GET    /api/v1/games/:id/download   # Download game
DELETE /api/v1/games/:id            # Remove game (admin)
POST   /api/v1/games/:id/review     # Submit review
GET    /api/v1/games/:id/reviews    # Get reviews
```

---

## Part 9: Monetization Strategy

### 9.1 Core Product

**Slip (Free)**
- All built-in games
- All built-in animations
- Plugin system
- Open source

### 9.2 Premium Content

**Game Packs ($0.99-$2.99)**
- "Retro Pack" - 5 classic arcade games
- "Puzzle Pack" - 5 brain teasers
- "Animation Pack" - 10 visualizations

### 9.3 Creator Program

**For Game Developers:**
- Submit games to the store
- 70/30 revenue split
- Games priced $0.49-$4.99

### 9.4 Sponsorship

**For Companies:**
- Custom branded animations
- Company logo idle screen
- Team high score boards

---

## Part 10: Technical Considerations

### 10.1 Performance

- Target: 30 FPS minimum on all platforms
- Use string builder for frame assembly
- Minimize allocations in game loop
- Profile with `go tool pprof`

### 10.2 Terminal Compatibility

- Test on: iTerm2, Terminal.app, GNOME Terminal, Windows Terminal, Alacritty, Kitty
- Use only standard ANSI sequences
- Graceful degradation for limited terminals

### 10.3 Input Handling

- Use raw terminal mode for immediate key response
- Handle ANSI escape sequences for arrow keys
- Debounce rapid repeated input
- Restore terminal state on exit (always!)

### 10.4 Cross-Platform

- Go makes this easy
- Test on: Linux x64, macOS Intel, macOS ARM, Windows x64
- WSL testing for Windows + tmux

### 10.5 Security (Plugins)

- Lua sandbox: disable io, os, loadfile, dofile
- WASM: natural sandbox
- Native: sign binaries, verify checksums

---

## Part 11: Dependencies

### 11.1 Go Modules

```go
// go.mod
require (
    github.com/charmbracelet/bubbletea v0.25.0  // TUI framework
    github.com/charmbracelet/lipgloss v0.9.0    // Styling
    github.com/charmbracelet/bubbles v0.17.0    // UI components
    github.com/yuin/gopher-lua v1.1.0           // Lua scripting
    golang.org/x/term v0.15.0                   // Terminal utilities
)
```

### 11.2 Optional (v2+)

```go
require (
    github.com/tetratelabs/wazero v1.5.0        // WASM runtime
)
```

---

## Part 12: Success Metrics

### 12.1 Technical

- Game loop maintains 30 FPS
- Input latency < 50ms
- Memory usage < 50MB
- Binary size < 20MB

### 12.2 User Experience

- Time to first game < 10 seconds
- Zero crashes in normal use
- Works on 95% of terminals

### 12.3 Adoption

- 1000 downloads in first month
- 10 community-contributed games in first year
- Positive sentiment on social media

---

## Appendix A: Example Snake Implementation

```go
package snake

import (
    "slip/internal/engine"
    "math/rand"
)

type Direction int

const (
    Up Direction = iota
    Down
    Left
    Right
)

type Point struct {
    X, Y int
}

type SnakeGame struct {
    ctx       engine.GameContext
    snake     []Point
    direction Direction
    food      Point
    score     int
    gameOver  bool
    moveTimer float64
    moveDelay float64
}

func New() *SnakeGame {
    return &SnakeGame{
        moveDelay: 0.15, // seconds between moves
    }
}

func (g *SnakeGame) Info() engine.GameInfo {
    return engine.GameInfo{
        ID:          "snake",
        Name:        "Snake",
        Description: "Classic snake game",
        Author:      "Slip",
        Version:     "1.0.0",
        MinWidth:    20,
        MinHeight:   10,
    }
}

func (g *SnakeGame) Init(ctx engine.GameContext) error {
    g.ctx = ctx
    g.reset()
    return nil
}

func (g *SnakeGame) reset() {
    centerX := g.ctx.Width / 2
    centerY := g.ctx.Height / 2
    g.snake = []Point{
        {centerX, centerY},
        {centerX - 1, centerY},
        {centerX - 2, centerY},
    }
    g.direction = Right
    g.score = 0
    g.gameOver = false
    g.spawnFood()
}

func (g *SnakeGame) spawnFood() {
    for {
        g.food = Point{
            X: rand.Intn(g.ctx.Width - 2) + 1,
            Y: rand.Intn(g.ctx.Height - 4) + 2,
        }
        // Make sure food isn't on snake
        collision := false
        for _, p := range g.snake {
            if p == g.food {
                collision = true
                break
            }
        }
        if !collision {
            break
        }
    }
}

func (g *SnakeGame) Start() error { return nil }
func (g *SnakeGame) Stop() error  { return nil }

func (g *SnakeGame) Update(dt float64) error {
    if g.gameOver {
        return nil
    }

    g.moveTimer += dt
    if g.moveTimer < g.moveDelay {
        return nil
    }
    g.moveTimer = 0

    // Calculate new head position
    head := g.snake[0]
    newHead := head
    switch g.direction {
    case Up:
        newHead.Y--
    case Down:
        newHead.Y++
    case Left:
        newHead.X--
    case Right:
        newHead.X++
    }

    // Check wall collision
    if newHead.X <= 0 || newHead.X >= g.ctx.Width-1 ||
       newHead.Y <= 1 || newHead.Y >= g.ctx.Height-2 {
        g.gameOver = true
        return nil
    }

    // Check self collision
    for _, p := range g.snake {
        if p == newHead {
            g.gameOver = true
            return nil
        }
    }

    // Move snake
    g.snake = append([]Point{newHead}, g.snake...)

    // Check food collision
    if newHead == g.food {
        g.score += 10
        g.spawnFood()
        // Speed up slightly
        if g.moveDelay > 0.05 {
            g.moveDelay -= 0.005
        }
    } else {
        // Remove tail if no food eaten
        g.snake = g.snake[:len(g.snake)-1]
    }

    return nil
}

func (g *SnakeGame) Render(screen *engine.Screen) error {
    style := engine.Style{FG: engine.White}

    // Draw border
    screen.DrawBox(0, 1, g.ctx.Width, g.ctx.Height-2, engine.SingleBox)

    // Draw header
    screen.DrawText(1, 0, "SNAKE", style)
    screen.DrawText(g.ctx.Width-15, 0,
        fmt.Sprintf("Score: %d", g.score), style)

    // Draw snake
    for i, p := range g.snake {
        char := '█'
        if i == 0 {
            char = '█' // Head could be different
        }
        screen.Set(p.X, p.Y, char, engine.Style{FG: engine.Green})
    }

    // Draw food
    screen.Set(g.food.X, g.food.Y, '*', engine.Style{FG: engine.Red})

    // Draw footer
    screen.DrawText(1, g.ctx.Height-1,
        "[←↑↓→] Move   [P] Pause   [Q] Quit",
        engine.Style{FG: engine.Gray})

    // Draw game over
    if g.gameOver {
        msg := "GAME OVER - Press R to restart"
        x := (g.ctx.Width - len(msg)) / 2
        y := g.ctx.Height / 2
        screen.DrawText(x, y, msg, engine.Style{FG: engine.Red, Bold: true})
    }

    return nil
}

func (g *SnakeGame) HandleInput(input engine.Input) error {
    if g.gameOver {
        if input.Key == engine.KeyRune && input.Rune == 'r' {
            g.reset()
        }
        return nil
    }

    switch input.Key {
    case engine.KeyUp:
        if g.direction != Down {
            g.direction = Up
        }
    case engine.KeyDown:
        if g.direction != Up {
            g.direction = Down
        }
    case engine.KeyLeft:
        if g.direction != Right {
            g.direction = Left
        }
    case engine.KeyRight:
        if g.direction != Left {
            g.direction = Right
        }
    }
    return nil
}
```

---

## Appendix B: Lua Game API

```lua
-- Example Lua game: Simple Pong

local game = {}

function game.info()
    return {
        id = "lua-pong",
        name = "Lua Pong",
        description = "Simple pong in Lua",
        author = "Community",
        version = "1.0.0",
        min_width = 30,
        min_height = 15
    }
end

local ball_x, ball_y
local ball_vx, ball_vy
local paddle_y
local score

function game.init(ctx)
    ball_x = ctx.width / 2
    ball_y = ctx.height / 2
    ball_vx = 1
    ball_vy = 0.5
    paddle_y = ctx.height / 2
    score = 0
end

function game.update(dt)
    -- Move ball
    ball_x = ball_x + ball_vx
    ball_y = ball_y + ball_vy

    -- Bounce off top/bottom
    if ball_y <= 1 or ball_y >= slip.height() - 2 then
        ball_vy = -ball_vy
    end

    -- Check paddle collision
    if ball_x <= 3 and ball_y >= paddle_y - 2 and ball_y <= paddle_y + 2 then
        ball_vx = -ball_vx
        score = score + 1
    end

    -- Ball out of bounds
    if ball_x < 0 then
        game.init({width = slip.width(), height = slip.height()})
    end

    -- Bounce off right wall
    if ball_x >= slip.width() - 1 then
        ball_vx = -ball_vx
    end
end

function game.render()
    slip.clear()

    -- Draw border
    slip.box(0, 0, slip.width(), slip.height())

    -- Draw paddle
    for i = -2, 2 do
        slip.set(2, paddle_y + i, "|")
    end

    -- Draw ball
    slip.set(math.floor(ball_x), math.floor(ball_y), "o")

    -- Draw score
    slip.text(slip.width() / 2 - 3, 0, "Score: " .. score)
end

function game.input(key)
    if key == "up" and paddle_y > 3 then
        paddle_y = paddle_y - 1
    elseif key == "down" and paddle_y < slip.height() - 4 then
        paddle_y = paddle_y + 1
    end
end

return game
```

---

## Appendix C: Comparison with Alternatives

| Feature | Slip | cool-retro-term | ASCII Movie Player |
|---------|------|-----------------|-------------------|
| Games | Yes | No | No |
| Animations | Yes | Visual only | Yes |
| tmux integration | Yes | No | No |
| Plugin system | Yes | No | No |
| Cross-platform | Yes | Limited | Yes |
| Single binary | Yes | No | Python |
| Works alongside AI CLIs | Yes | No | No |

---

## Conclusion

Slip represents a unique opportunity to create a beloved developer tool that makes terminal waiting time enjoyable. By building on the solid foundation of Clipster and expanding into a full entertainment platform, we can create something that developers will genuinely want to use and share.

The phased approach ensures we deliver value incrementally while building toward the full vision. The plugin system enables community growth, and the store creates potential for sustainable revenue.

**Next Steps:**
1. Review this document
2. Approve high-level architecture
3. Begin Phase 0 implementation
4. Set up project tracking

---

*Document Version: 1.0*
*Last Updated: January 2025*
