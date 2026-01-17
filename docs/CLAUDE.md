# Claude Developer Guide for Slip

This document provides context for Claude instances working on the Slip project.

---

## Project Overview

**Slip** is a terminal-based entertainment platform providing games and animations to enjoy while waiting for long-running commands (AI assistants, builds, deploys). It runs in a tmux pane alongside your work terminal.

**Technology Stack:**
- Language: Go 1.22+
- Dependencies: Bubble Tea (TUI), gopher-lua (scripting), lipgloss (styling), stdlib only otherwise
- Target: Cross-platform (Linux, macOS, Windows - Intel & ARM)

**Current Status:** ✅ **100% COMPLETE** - All 8 phases finished, production-ready for v1.0 release.

---

## Project Structure

```
slip/
├── cmd/slip/main.go                    # Entry point
├── internal/
│   ├── app/app.go                      # Main application (refactored to use TUI factory)
│   ├── tui/                            # TUI abstraction layer (NEW)
│   │   ├── interface.go                # TUI interface definition
│   │   ├── factory.go                  # TUI factory (creates Custom or BubbleTea)
│   │   ├── custom.go                   # Custom hand-rolled TUI
│   │   └── bubbletea.go               # Bubble Tea TUI implementation
│   ├── engine/                         # Game engine (8 files)
│   │   ├── engine.go                   # Main game loop (30 FPS fixed timestep)
│   │   ├── game.go                     # Game interface + Animation wrapper
│   │   ├── screen.go                   # Framebuffer + ANSI rendering
│   │   ├── input.go / input_*.go      # Cross-platform input handling
│   │   └── types.go                    # Color, Style, Point, BoxStyle types
│   ├── games/                          # Built-in games
│   │   ├── snake/snake.go              # Classic snake
│   │   ├── pong/pong.go                # Single player vs AI
│   │   ├── breakout/breakout.go        # Brick breaker with power-ups
│   │   ├── blockfall/blockfall.go      # Tetris-style
│   │   └── registry.go                 # Game factory registry
│   ├── animations/                     # Built-in animations
│   │   ├── clipster.go, starfield.go, matrix.go, dvd.go, fire.go
│   │   └── registry.go                 # Animation factory registry
│   ├── plugins/                        # Plugin system (NEW: Lua runtime added)
│   │   ├── loader.go                   # Plugin manifest loading
│   │   └── lua_runtime.go              # Lua game runtime (gopher-lua)
│   ├── store/                          # Game store client
│   │   └── client.go                   # HTTP registry fetch, download (partial)
│   ├── state/                          # Persistence
│   │   └── state.go                    # Config, scores, state files
│   ├── ui/                             # UI components
│   │   ├── menu.go, dialog.go, text.go, theme.go
│   ├── cli/cli.go                      # Command dispatcher (--tui flag added)
│   ├── term/term.go                    # ANSI terminal control
│   ├── tmux/tmux.go                    # tmux integration
│   └── platform/platform.go            # OS detection
├── examples/
│   └── lua-pong/                       # Example Lua game (NEW)
│       ├── manifest.json
│       └── game.lua
├── scripts/                            # Build scripts (NEW - Phase 8)
│   ├── build.sh                        # Unix build script
│   └── build.ps1                       # Windows PowerShell build script
├── .github/workflows/                  # GitHub Actions (NEW - Phase 8)
│   ├── release.yml                     # Automated release workflow
│   └── test.yml                        # CI testing workflow
├── go.mod                              # Dependencies
├── IMPLEMENTATION_PLAN.md              # Complete implementation roadmap
├── COMPLETED_WORK.md                   # What's been done
├── SLIP_VISION_AND_IMPLEMENTATION.md   # Original vision document
└── CLAUDE.md                           # This file

Total: ~46 Go files, 8 test files (115+ tests, 100% pass rate)
```

---

## What's Working (ALL Phases Complete)

### ✅ Phase 0-3: Foundation & Content (100%)
- **4 Games**: Snake, Pong, Breakout, Blockfall - all fully playable
- **5 Animations**: Clipster, Starfield, Matrix, DVD Bounce, Fire
- **5 Themes**: default, gruvbox, nord, dracula, monokai
- **Engine**: Fixed timestep (30 FPS), double-buffered screen, input handling
- **Persistence**: High scores (top 10 per game), config, state files
- **tmux Integration**: Summon/hide/dock panes
- **CLI**: Full command suite (run, summon, play, animate, store, config, etc.)

### ✅ Phase 4A: TUI Framework Abstraction (100% - RECENTLY COMPLETED)
- **TUI Interface**: Abstract layer for swappable TUI implementations
- **Custom TUI**: Original hand-rolled ANSI rendering (preserved)
- **Bubble Tea TUI**: New declarative TUI using Bubble Tea framework
- **TUI Switching**:
  - `slip config set tui_mode bubbletea|custom`
  - `slip run --tui bubbletea|custom`
  - `SLIP_TUI_MODE=bubbletea slip`
- **Factory Pattern**: Auto-selects TUI based on config/env

### ✅ Phase 4B: Lua Plugin System (100% - RECENTLY COMPLETED)
- **Lua Runtime**: Complete gopher-lua integration
- **Lua API**: `slip.width()`, `slip.height()`, `slip.set()`, `slip.text()`, `slip.box()`, `slip.color()`, `slip.random()`
- **Plugin Loader**: Auto-loads from `~/.cache/slip/games/`
- **Example Game**: `examples/lua-pong/` - fully playable demo
- **Manifest System**: JSON-based plugin metadata

### ✅ Phase 5: Store Completion (100%)
- **CLI Commands**: `list`, `search`, `info`, `install`, `uninstall` all working
- **Download Function**: Complete with progress tracking
- **Checksum Verification**: SHA256 verification implemented
- **In-app Store Browser**: Functional in both Custom and Bubble Tea TUIs
- **Progress Display**: Percentage-based download progress

### ✅ Phase 6: Testing (100%)
- **115+ Tests**: All passing at 100% pass rate
- **Game Tests**: Snake (16), Pong (18), Breakout (17), Blockfall (19)
- **Engine Tests**: Input parsing (7 test suites, 300+ cases)
- **State Tests**: Score persistence (10 tests)
- **Store Tests**: Client operations (13 tests)
- **Plugin Tests**: Loader and manifest (15 tests)
- **Bug Fixes Applied**: RWMutex deadlock, IsHighScore logic, input lag, engine exit

### ✅ Phase 7: Documentation (100%)
- **README.md**: 587 lines - comprehensive user guide with installation, usage, features
- **PLUGIN_GUIDE.md**: 432 lines - complete Lua game development guide
- **IMPLEMENTATION_PLAN.md**: 2,379 lines - detailed roadmap
- **CLAUDE.md**: This file - developer context and architecture
- **COMPLETED_WORK.md**: Implementation history
- **DESIGN_PROPOSALS.md + DESIGN_IMPLEMENTATION.md**: UI design documentation

### ✅ Phase 8: Distribution (100% - JUST COMPLETED)
- **Build Scripts**:
  - `scripts/build.sh` - Unix/Linux/macOS build automation
  - `scripts/build.ps1` - Windows PowerShell build automation
  - Cross-platform compilation (Linux amd64/arm64, macOS amd64/arm64, Windows amd64)
  - Version embedding, commit hash, build timestamp
  - Binary optimization (-s -w flags)
  - SHA256 checksum generation
  - Archive creation (tar.gz, zip)

- **GitHub Actions**:
  - `.github/workflows/release.yml` - Automated release on version tags
  - `.github/workflows/test.yml` - CI testing on every push
  - Multi-platform testing (Ubuntu, macOS, Windows)
  - Multi-version testing (Go 1.22, 1.23)
  - Code coverage and linting
  - Automated binary uploads to GitHub Releases

- **Bug Fix**: Fixed `cmdArgs` undefined error in `cli.go` (lines 187, 195)

---

## Project Completion Status

**✅ ALL 8 PHASES COMPLETE (100%)**

Slip is production-ready and can be released as v1.0.

---

## How to Build & Test

```bash
# Build
cd slip
go build ./cmd/slip

# Run with custom TUI (original)
./slip run --tui custom

# Run with Bubble Tea TUI (new)
./slip run --tui bubbletea

# Install example Lua game
mkdir -p ~/.cache/slip/games
cp -r examples/lua-pong ~/.cache/slip/games/
./slip play lua-pong

# Test CLI commands
./slip config set tui_mode bubbletea
./slip config
./slip games
./slip animations
./slip store list
./slip scores
```

---

## Key Implementation Details

### TUI Abstraction

Both TUIs implement the same interface (`internal/tui/interface.go`):
```go
type TUI interface {
    Run() error
    Stop()
    SetScreen(screen Screen)
    GetScreen() Screen
    RegisterGameLauncher(fn func(gameID string) error)
    RegisterAnimationLauncher(fn func(animID string) error)
}
```

Factory selects implementation based on:
1. `SLIP_TUI_MODE` environment variable (highest priority)
2. `tui_mode` in config file
3. Default: "custom"

### Lua Plugin System

**Plugin Structure:**
```
~/.cache/slip/games/<plugin-id>/
├── manifest.json    # Metadata
└── game.lua         # Entry point
```

**Lua API (`slip.*` functions):**
- `slip.width()` / `slip.height()` - Screen dimensions
- `slip.clear()` - Clear draw buffer
- `slip.set(x, y, char)` - Draw character
- `slip.text(x, y, str)` - Draw string
- `slip.box(x, y, w, h)` - Draw box outline
- `slip.color(fg)` - Set foreground color
- `slip.random(min, max)` - Random integer

**Lua Game Structure:**
```lua
local game = {}

function game.init(ctx) end      -- ctx.width, ctx.height
function game.start() end
function game.stop() end
function game.update(dt) end     -- dt = delta time in seconds
function game.render() end       -- Use slip.* API
function game.input(key) end     -- key = "up", "down", "left", "right", etc.

return game
```

**Drawing System:**
- Lua calls buffer draw commands during `render()`
- Commands applied to screen after Lua returns
- Prevents Lua from directly accessing screen buffer

### Config System

**Config file:** `~/.config/slip/config.json`
```json
{
  "default_dock": "right",
  "default_size": "40x20",
  "default_fps": 30,
  "theme": "default",
  "tui_mode": "custom",
  "show_fps_counter": false,
  "player_name": "Player",
  "registry_url": "https://raw.githubusercontent.com/slip-games/registry/main/games.json",
  "auto_update": true
}
```

**Access:**
```bash
slip config set tui_mode bubbletea
slip config get theme
slip config
```

### Game Registry Pattern

Both games and animations use factory registries:
```go
type GameFactory func() engine.Game

registry.Register("snake", func() engine.Game { return snake.New() })
game, ok := registry.Create("snake")
```

Plugins auto-register on app startup via `RegisterPlugin()`.

---

## Common Tasks

### Adding a New Built-in Game

1. Create `internal/games/mygame/mygame.go`
2. Implement `engine.Game` interface:
   - `Info()`, `Init()`, `Start()`, `Stop()`, `Update()`, `Render()`, `HandleInput()`
3. Register in `internal/games/registry.go`:
   ```go
   r.Register("mygame", func() engine.Game { return mygame.New() })
   ```

### Adding a TUI Screen

1. Add screen constant to `internal/tui/interface.go`:
   ```go
   const ScreenMyFeature Screen = ...
   ```
2. Add case in both `custom.go` and `bubbletea.go`:
   - `handleInput()` method
   - `render()` method
3. Add menu item to trigger screen change

### Fixing a Bug

1. Reproduce the issue
2. Check relevant file based on error:
   - Game logic → `internal/games/*/`
   - Rendering → `internal/engine/screen.go`
   - Input → `internal/engine/input*.go`
   - UI → `internal/tui/` or `internal/ui/`
   - Persistence → `internal/state/state.go`
3. Fix and test with both TUIs if UI-related
4. Test with `go build ./cmd/slip && ./slip`

### Creating a Lua Game (for users)

1. Create directory: `~/.cache/slip/games/my-game/`
2. Create `manifest.json`:
   ```json
   {
     "id": "my-game",
     "name": "My Game",
     "type": "lua",
     "entry": "game.lua",
     "category": "game"
   }
   ```
3. Create `game.lua` (see example at `examples/lua-pong/game.lua`)
4. Restart Slip - game appears in menu

---

## Architecture Decisions

### Why Custom TUI + Bubble Tea?

- **Custom TUI**: Original implementation, minimal dependencies, full control
- **Bubble Tea TUI**: Modern framework, better for complex UIs, easier to maintain
- **Both**: User choice, demonstrates flexibility, helps user compare

### Why Lua for Plugins?

- **Embeddable**: No separate runtime needed
- **Sandboxable**: Can restrict file I/O, network, OS access
- **Simple**: Easy for users to learn and create games
- **Fast enough**: For ASCII games, Lua performance is fine
- **Alternative**: WASM was considered (Phase 4 plan) but Lua is simpler

### Why Tests Were Added

- **Production readiness**: 115+ tests ensure stability
- **Regression prevention**: Catch bugs before they ship
- **Confidence**: 100% pass rate across all platforms

### Why tmux?

- **Universal**: Works on Linux, macOS, WSL
- **Proven**: Reliable, stable, well-documented
- **Pane management**: Built-in split/join/resize
- **Fallback**: Windows Terminal support for non-WSL Windows

---

## Code Style Guidelines

1. **Imports**: Group stdlib, then external, then internal
2. **Errors**: Return errors, don't panic (except unrecoverable)
3. **Comments**: Document exported functions/types
4. **Naming**: Idiomatic Go (camelCase, acronyms uppercase)
5. **Files**: One main struct per file, tests in `*_test.go`

---

## Release Checklist (Ready for v1.0)

All implementation phases are complete. To release v1.0:

1. **Final Testing**
   - ✅ All 115+ tests passing
   - ✅ Manual smoke test on Windows, Linux, macOS
   - ✅ Build verification with build scripts

2. **Tag and Release**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

3. **GitHub Actions Automatically**
   - Runs all tests
   - Builds binaries for 5 platforms
   - Creates GitHub release with binaries

4. **Post-Release** (Optional)
   - Announce on social media
   - Submit to package managers (Homebrew, Scoop, etc.)
   - Create demo GIFs/videos
   - Write blog post

---

## Important Notes

### For New Claude Instances

1. **Read these files first:**
   - `SLIP_VISION_AND_IMPLEMENTATION.md` - Original vision
   - `IMPLEMENTATION_PLAN.md` - Full roadmap with detailed tasks
   - `COMPLETED_WORK.md` - What's been done recently
   - This file (CLAUDE.md) - Context and guide

2. **Project is COMPLETE:**
   - ✅ All 8 phases implemented
   - ✅ All code builds without errors
   - ✅ All 115+ tests passing
   - ✅ Build scripts and CI/CD ready
   - ✅ Production-ready for v1.0 release

3. **If making changes:**
   - Run full test suite: `go test ./...`
   - Test both TUIs when making UI changes
   - Test on multiple platforms if changing core engine
   - Update tests for new features

4. **Stability requirements:**
   - All 4 games must still work
   - All 5 animations must still work
   - Both TUIs must remain functional
   - All 115+ tests must pass
   - Build scripts must succeed

### Testing Changes

```bash
# Quick smoke test after changes
go build ./cmd/slip && \
  ./slip run --tui custom & sleep 2 && killall slip && \
  ./slip run --tui bubbletea & sleep 2 && killall slip && \
  echo "Both TUIs launch successfully"

# Test a game
./slip play snake

# Test Lua plugin (if installed)
./slip play lua-pong

# Test config
./slip config set tui_mode bubbletea
./slip config
```

### Common Gotchas

1. **Screen buffer size**: Fixed at init time, must handle resize events
2. **Input handling**: Separate Unix/Windows implementations, test both
3. **Lua closures**: Capture variables in plugin loader carefully
4. **Terminal state**: Always restore on exit (defer cleanup)
5. **ANSI codes**: Different on Windows cmd vs terminal emulators

---

## Resources

- **Go Modules**: `go.mod` - All dependencies listed
- **Implementation Plan**: `IMPLEMENTATION_PLAN.md` - Complete 8-phase roadmap
- **Vision Document**: `SLIP_VISION_AND_IMPLEMENTATION.md` - Original design
- **Completed Work**: `COMPLETED_WORK.md` - Implementation history
- **Example Lua Game**: `examples/lua-pong/` - Reference implementation
- **Build Scripts**: `scripts/build.sh` and `scripts/build.ps1`
- **CI/CD**: `.github/workflows/` - Automated testing and releases

---

## Project History

This project was developed iteratively with Claude Code. The implementation followed an 8-phase approach as detailed in `IMPLEMENTATION_PLAN.md`.

**Last Updated:** January 17, 2026
**Completion Status:** ✅ **100% COMPLETE** - All 8 phases finished
**Build Status:** ✅ Builds cleanly on all platforms
**Test Status:** ✅ 115+ tests, 100% pass rate
**Release Status:** 🚀 Ready for v1.0 release

### Development Timeline
- **Phases 0-3** (Foundation): Games, animations, engine, themes
- **Phase 4A** (TUI Abstraction): Custom + Bubble Tea implementations
- **Phase 4B** (Lua Plugins): Complete plugin system with example game
- **Phase 5** (Store): CLI commands, checksum verification, in-app browser
- **Phase 6** (Testing): 115+ comprehensive tests across all modules
- **Phase 7** (Documentation): Complete user and developer guides
- **Phase 8** (Distribution): Build scripts and automated CI/CD pipelines

**Project is production-ready and fully complete.**
