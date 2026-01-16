# Slip Implementation Plan

## Overview

This document provides a step-by-step implementation plan for completing the Slip terminal entertainment platform. Each task is designed to be self-contained and can be picked up by any developer.

**Current Status:** ~70% complete
**Target:** 100% feature complete with tests and documentation

---

## Phase 4A: TUI Framework Abstraction (Bubble Tea Toggle)

### Goal
Create an abstraction layer that allows switching between the custom TUI and Bubble Tea, configurable via `slip config set tui_mode bubbletea|custom`.

### Task 4A.1: Define TUI Interface

**File:** `internal/tui/interface.go`

```go
package tui

import "slip/internal/engine"

// TUI is the interface that both custom and Bubble Tea implementations must satisfy
type TUI interface {
    // Run starts the TUI application
    Run() error

    // Stop gracefully stops the TUI
    Stop()

    // SetScreen changes the current screen
    SetScreen(screen Screen)

    // GetScreen returns the current screen
    GetScreen() Screen

    // RegisterGame registers a game launcher callback
    RegisterGameLauncher(fn func(gameID string) error)

    // RegisterAnimation registers an animation launcher callback
    RegisterAnimationLauncher(fn func(animID string) error)
}

// Screen represents different screens in the application
type Screen int

const (
    ScreenHome Screen = iota
    ScreenGames
    ScreenAnimations
    ScreenStore
    ScreenSettings
    ScreenHelp
)

// Config holds TUI configuration
type Config struct {
    Width  int
    Height int
    FPS    int
    Theme  string
}

// MenuItem represents a menu item
type MenuItem struct {
    ID          string
    Label       string
    Description string
}

// GameInfo represents game metadata
type GameInfo struct {
    ID          string
    Name        string
    Description string
}
```

**Acceptance Criteria:**
- [ ] Interface defined with all necessary methods
- [ ] Screen constants defined
- [ ] Config struct defined
- [ ] No implementation yet, just the contract

---

### Task 4A.2: Create TUI Factory

**File:** `internal/tui/factory.go`

```go
package tui

import "slip/internal/state"

// NewTUI creates a TUI based on configuration
func NewTUI(cfg Config) (TUI, error) {
    // Load user preference
    userCfg, _ := state.LoadConfig()
    mode := userCfg.TUIMode
    if mode == "" {
        mode = "custom" // Default to custom for backwards compatibility
    }

    switch mode {
    case "bubbletea":
        return NewBubbleTeaTUI(cfg)
    case "custom":
        return NewCustomTUI(cfg)
    default:
        return NewCustomTUI(cfg)
    }
}
```

**Acceptance Criteria:**
- [ ] Factory function reads `tui_mode` from config
- [ ] Returns appropriate TUI implementation
- [ ] Defaults to "custom" if not set

---

### Task 4A.3: Wrap Existing Custom TUI

**File:** `internal/tui/custom.go`

Refactor the existing `internal/app/app.go` logic into a struct that implements the `TUI` interface.

**Steps:**
1. Create `CustomTUI` struct
2. Move rendering logic from `app.go`
3. Implement all `TUI` interface methods
4. Keep existing functionality identical

**Key Code Structure:**
```go
package tui

type CustomTUI struct {
    width       int
    height      int
    fps         int
    screen      Screen
    running     bool

    // Menus
    homeMenu       *ui.Menu
    gamesMenu      *ui.Menu
    animationsMenu *ui.Menu

    // Callbacks
    gameLauncher func(string) error
    animLauncher func(string) error

    // Theme
    theme ui.Theme
}

func NewCustomTUI(cfg Config) (*CustomTUI, error) {
    // Initialize with existing logic from app.go
}

func (t *CustomTUI) Run() error {
    // Existing run loop from app.go
}

// ... implement remaining interface methods
```

**Acceptance Criteria:**
- [ ] All existing functionality preserved
- [ ] Implements `TUI` interface
- [ ] `slip run` works exactly as before
- [ ] No regressions in menu navigation

---

### Task 4A.4: Implement Bubble Tea TUI

**File:** `internal/tui/bubbletea.go`

**Dependencies to add to go.mod:**
```
github.com/charmbracelet/bubbletea v0.25.0
github.com/charmbracelet/lipgloss v0.9.0
github.com/charmbracelet/bubbles v0.17.0
```

**Implementation:**
```go
package tui

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type BubbleTeaTUI struct {
    width       int
    height      int
    screen      Screen

    // Bubble Tea specific
    program     *tea.Program

    // Menu state
    menuIndex   int
    menuItems   []MenuItem

    // Callbacks
    gameLauncher func(string) error
    animLauncher func(string) error

    // Styles
    titleStyle    lipgloss.Style
    itemStyle     lipgloss.Style
    selectedStyle lipgloss.Style
}

// Model for Bubble Tea
type model struct {
    tui *BubbleTeaTUI
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "up", "k":
            if m.tui.menuIndex > 0 {
                m.tui.menuIndex--
            }
        case "down", "j":
            if m.tui.menuIndex < len(m.tui.menuItems)-1 {
                m.tui.menuIndex++
            }
        case "enter":
            return m, m.tui.handleSelection()
        case "esc":
            if m.tui.screen != ScreenHome {
                m.tui.screen = ScreenHome
                m.tui.menuIndex = 0
            }
        }
    case tea.WindowSizeMsg:
        m.tui.width = msg.Width
        m.tui.height = msg.Height
    }
    return m, nil
}

func (m model) View() string {
    switch m.tui.screen {
    case ScreenHome:
        return m.tui.renderHome()
    case ScreenGames:
        return m.tui.renderGames()
    case ScreenAnimations:
        return m.tui.renderAnimations()
    case ScreenStore:
        return m.tui.renderStore()
    case ScreenSettings:
        return m.tui.renderSettings()
    case ScreenHelp:
        return m.tui.renderHelp()
    default:
        return m.tui.renderHome()
    }
}

func NewBubbleTeaTUI(cfg Config) (*BubbleTeaTUI, error) {
    t := &BubbleTeaTUI{
        width:  cfg.Width,
        height: cfg.Height,
    }

    // Initialize styles
    t.titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("51")). // Cyan
        MarginBottom(1)

    t.itemStyle = lipgloss.NewStyle().
        PaddingLeft(2)

    t.selectedStyle = lipgloss.NewStyle().
        PaddingLeft(2).
        Foreground(lipgloss.Color("229")). // Yellow
        Bold(true)

    t.initMenus()
    return t, nil
}

func (t *BubbleTeaTUI) Run() error {
    m := model{tui: t}
    t.program = tea.NewProgram(m, tea.WithAltScreen())
    _, err := t.program.Run()
    return err
}

func (t *BubbleTeaTUI) Stop() {
    if t.program != nil {
        t.program.Quit()
    }
}

// Implement remaining interface methods...
```

**Acceptance Criteria:**
- [ ] Bubble Tea dependencies added
- [ ] Home screen renders with menu
- [ ] Navigation with arrow keys works
- [ ] Can switch to Games/Animations screens
- [ ] Escape returns to previous screen
- [ ] Q quits the application
- [ ] Visual styling matches or improves on custom TUI

---

### Task 4A.5: Add TUI Mode to Config

**File:** `internal/state/state.go`

Add `TUIMode` field to Config struct:

```go
type Config struct {
    DefaultDock    string `json:"default_dock"`
    DefaultSize    string `json:"default_size"`
    DefaultFPS     int    `json:"default_fps"`
    Theme          string `json:"theme"`
    ShowFPSCounter bool   `json:"show_fps_counter"`
    RegistryURL    string `json:"registry_url"`
    AutoUpdate     bool   `json:"auto_update"`
    PlayerName     string `json:"player_name"`
    TUIMode        string `json:"tui_mode"` // NEW: "custom" or "bubbletea"
}
```

**File:** `internal/cli/cli.go`

Update config commands to handle `tui_mode`:

```go
// In configCmd, add validation for tui_mode
case "set":
    if args[1] == "tui_mode" {
        if args[2] != "custom" && args[2] != "bubbletea" {
            fmt.Fprintln(os.Stderr, "tui_mode must be 'custom' or 'bubbletea'")
            return 2
        }
    }
    // ... existing set logic
```

**Acceptance Criteria:**
- [ ] Config struct has TUIMode field
- [ ] `slip config set tui_mode bubbletea` works
- [ ] `slip config set tui_mode custom` works
- [ ] Invalid values are rejected
- [ ] `slip config` shows current tui_mode

---

### Task 4A.6: Update App to Use TUI Factory

**File:** `internal/app/app.go`

Refactor to use the TUI abstraction:

```go
package app

import (
    "slip/internal/tui"
    "slip/internal/engine"
    "slip/internal/games"
    "slip/internal/animations"
)

type App struct {
    tui          tui.TUI
    gameRegistry *games.Registry
    animRegistry *animations.Registry
    engine       *engine.Engine
    width        int
    height       int
    fps          int
}

func New(width, height, fps int) *App {
    app := &App{
        width:        width,
        height:       height,
        fps:          fps,
        gameRegistry: games.NewRegistry(),
        animRegistry: animations.NewRegistry(),
    }

    // Create TUI using factory
    cfg := tui.Config{
        Width:  width,
        Height: height,
        FPS:    fps,
    }

    t, err := tui.NewTUI(cfg)
    if err != nil {
        // Fallback to custom
        t, _ = tui.NewCustomTUI(cfg)
    }

    app.tui = t

    // Register callbacks
    t.RegisterGameLauncher(app.launchGame)
    t.RegisterAnimationLauncher(app.launchAnimation)

    return app
}

func (a *App) Run() error {
    return a.tui.Run()
}

func (a *App) launchGame(id string) error {
    game, ok := a.gameRegistry.Create(id)
    if !ok {
        return fmt.Errorf("game not found: %s", id)
    }

    eng := engine.NewEngine(a.width, a.height, a.fps)
    eng.SetGame(game)
    return eng.Run()
}

func (a *App) launchAnimation(id string) error {
    anim, ok := a.animRegistry.CreateAsGame(id)
    if !ok {
        return fmt.Errorf("animation not found: %s", id)
    }

    eng := engine.NewEngine(a.width, a.height, a.fps)
    eng.SetGame(anim)
    return eng.Run()
}
```

**Acceptance Criteria:**
- [ ] App uses TUI factory
- [ ] `slip config set tui_mode bubbletea && slip run` uses Bubble Tea
- [ ] `slip config set tui_mode custom && slip run` uses custom TUI
- [ ] Game launching works from both TUIs
- [ ] Animation launching works from both TUIs

---

### Task 4A.7: Add CLI Flag Override

**File:** `internal/cli/cli.go`

Add `--tui` flag to `run` command:

```go
func runCmd(args []string) int {
    fs := flag.NewFlagSet("run", flag.ContinueOnError)
    width := fs.Int("width", envInt("SLIP_WIDTH", defaultWidth), "Screen width")
    height := fs.Int("height", envInt("SLIP_HEIGHT", defaultHeight), "Screen height")
    size := fs.String("size", envString("SLIP_SIZE", ""), "Screen size (WxH)")
    fps := fs.Int("fps", envInt("SLIP_FPS", defaultFPS), "Frames per second")
    tuiMode := fs.String("tui", "", "TUI mode: custom or bubbletea (overrides config)")

    // ... existing parsing

    // If --tui flag provided, temporarily set it
    if *tuiMode != "" {
        // Validate
        if *tuiMode != "custom" && *tuiMode != "bubbletea" {
            fmt.Fprintln(os.Stderr, "--tui must be 'custom' or 'bubbletea'")
            return 2
        }
        // Set for this session
        os.Setenv("SLIP_TUI_MODE", *tuiMode)
    }

    // ... rest of function
}
```

Update factory to check env var:

```go
func NewTUI(cfg Config) (TUI, error) {
    // Check env var first (for CLI override)
    mode := os.Getenv("SLIP_TUI_MODE")
    if mode == "" {
        userCfg, _ := state.LoadConfig()
        mode = userCfg.TUIMode
    }
    if mode == "" {
        mode = "custom"
    }
    // ...
}
```

**Acceptance Criteria:**
- [ ] `slip run --tui bubbletea` works
- [ ] `slip run --tui custom` works
- [ ] CLI flag overrides config setting
- [ ] Environment variable `SLIP_TUI_MODE` also works

---

## Phase 4B: Plugin System - Lua Runtime

### Goal
Integrate gopher-lua to enable running Lua-based games/animations.

### Task 4B.1: Add Lua Dependency

**File:** `go.mod`

Add:
```
github.com/yuin/gopher-lua v1.1.1
```

Run:
```bash
go get github.com/yuin/gopher-lua
```

**Acceptance Criteria:**
- [ ] Dependency added to go.mod
- [ ] `go build` succeeds
- [ ] No import errors

---

### Task 4B.2: Create Lua Game Runtime

**File:** `internal/plugins/lua_runtime.go`

```go
package plugins

import (
    "fmt"
    "path/filepath"

    lua "github.com/yuin/gopher-lua"
    "slip/internal/engine"
)

// LuaGame wraps a Lua script as a Game
type LuaGame struct {
    state    *lua.LState
    manifest *Manifest
    dir      string

    // Cached Lua functions
    fnUpdate func(float64)
    fnRender func(*engine.Screen)
    fnInput  func(engine.Input)

    // Game context
    ctx *engine.GameContext
}

// NewLuaGame creates a new Lua game from a manifest
func NewLuaGame(manifest *Manifest, dir string) (*LuaGame, error) {
    L := lua.NewState()

    game := &LuaGame{
        state:    L,
        manifest: manifest,
        dir:      dir,
    }

    // Register the slip API
    game.registerAPI()

    // Load the entry script
    entryPath := filepath.Join(dir, manifest.Entry)
    if err := L.DoFile(entryPath); err != nil {
        L.Close()
        return nil, fmt.Errorf("failed to load lua script: %w", err)
    }

    return game, nil
}

func (g *LuaGame) registerAPI() {
    L := g.state

    // Create slip table
    slip := L.NewTable()
    L.SetGlobal("slip", slip)

    // slip.width() -> returns screen width
    L.SetField(slip, "width", L.NewFunction(func(L *lua.LState) int {
        if g.ctx != nil {
            L.Push(lua.LNumber(g.ctx.Width))
        } else {
            L.Push(lua.LNumber(40))
        }
        return 1
    }))

    // slip.height() -> returns screen height
    L.SetField(slip, "height", L.NewFunction(func(L *lua.LState) int {
        if g.ctx != nil {
            L.Push(lua.LNumber(g.ctx.Height))
        } else {
            L.Push(lua.LNumber(20))
        }
        return 1
    }))

    // slip.set(x, y, char) -> sets a character on screen
    // Note: This is stored and applied during Render
    L.SetField(slip, "set", L.NewFunction(func(L *lua.LState) int {
        // Implementation stores draw commands
        return 0
    }))

    // slip.text(x, y, str) -> draws text
    L.SetField(slip, "text", L.NewFunction(func(L *lua.LState) int {
        return 0
    }))

    // slip.clear() -> clears pending draw commands
    L.SetField(slip, "clear", L.NewFunction(func(L *lua.LState) int {
        return 0
    }))

    // slip.box(x, y, w, h) -> draws a box
    L.SetField(slip, "box", L.NewFunction(func(L *lua.LState) int {
        return 0
    }))

    // slip.color(fg, bg) -> sets current color
    L.SetField(slip, "color", L.NewFunction(func(L *lua.LState) int {
        return 0
    }))
}

// Info returns game metadata
func (g *LuaGame) Info() engine.GameInfo {
    return engine.GameInfo{
        ID:          g.manifest.ID,
        Name:        g.manifest.Name,
        Description: g.manifest.Description,
        Author:      g.manifest.Author,
        Version:     g.manifest.Version,
        MinWidth:    g.manifest.MinWidth,
        MinHeight:   g.manifest.MinHeight,
    }
}

// Init initializes the game
func (g *LuaGame) Init(ctx *engine.GameContext) error {
    g.ctx = ctx

    // Call game.init(ctx) in Lua
    if err := g.callLuaMethod("init", ctx); err != nil {
        return err
    }

    return nil
}

// Start starts the game
func (g *LuaGame) Start() error {
    return g.callLuaMethod("start")
}

// Stop stops the game
func (g *LuaGame) Stop() error {
    err := g.callLuaMethod("stop")
    g.state.Close()
    return err
}

// Update updates game state
func (g *LuaGame) Update(dt float64) error {
    return g.callLuaMethod("update", dt)
}

// Render renders the game
func (g *LuaGame) Render(screen *engine.Screen) error {
    // Set screen reference for API calls
    g.currentScreen = screen
    defer func() { g.currentScreen = nil }()

    return g.callLuaMethod("render")
}

// HandleInput handles input
func (g *LuaGame) HandleInput(input engine.Input) error {
    // Convert input to Lua-friendly format
    keyStr := inputToString(input)
    return g.callLuaMethod("input", keyStr)
}

func (g *LuaGame) callLuaMethod(method string, args ...interface{}) error {
    L := g.state

    // Get game table
    game := L.GetGlobal("game")
    if game == lua.LNil {
        return nil // No game table, skip
    }

    // Get method
    fn := L.GetField(game, method)
    if fn == lua.LNil {
        return nil // Method not defined, skip
    }

    // Push function
    L.Push(fn)

    // Push arguments
    for _, arg := range args {
        switch v := arg.(type) {
        case float64:
            L.Push(lua.LNumber(v))
        case int:
            L.Push(lua.LNumber(v))
        case string:
            L.Push(lua.LString(v))
        case *engine.GameContext:
            // Create context table
            ctx := L.NewTable()
            L.SetField(ctx, "width", lua.LNumber(v.Width))
            L.SetField(ctx, "height", lua.LNumber(v.Height))
            L.Push(ctx)
        default:
            L.Push(lua.LNil)
        }
    }

    // Call
    if err := L.PCall(len(args), 0, nil); err != nil {
        return err
    }

    return nil
}

func inputToString(input engine.Input) string {
    switch input.Key {
    case engine.KeyUp:
        return "up"
    case engine.KeyDown:
        return "down"
    case engine.KeyLeft:
        return "left"
    case engine.KeyRight:
        return "right"
    case engine.KeyEnter:
        return "enter"
    case engine.KeyEscape:
        return "escape"
    case engine.KeySpace:
        return "space"
    default:
        if input.Rune != 0 {
            return string(input.Rune)
        }
        return ""
    }
}
```

**Acceptance Criteria:**
- [ ] LuaGame struct implements engine.Game interface
- [ ] Lua state properly initialized and closed
- [ ] slip.* API functions registered
- [ ] Lua game methods (init, update, render, input) called correctly

---

### Task 4B.3: Implement Lua Drawing API

**File:** `internal/plugins/lua_runtime.go` (extend)

Add drawing command buffering:

```go
type drawCmd struct {
    typ   string // "set", "text", "box", "clear"
    x, y  int
    w, h  int
    char  rune
    text  string
    style engine.Style
}

type LuaGame struct {
    // ... existing fields
    drawCmds      []drawCmd
    currentScreen *engine.Screen
    currentStyle  engine.Style
}

func (g *LuaGame) registerAPI() {
    L := g.state
    slip := L.NewTable()
    L.SetGlobal("slip", slip)

    // slip.clear()
    L.SetField(slip, "clear", L.NewFunction(func(L *lua.LState) int {
        g.drawCmds = nil
        return 0
    }))

    // slip.set(x, y, char)
    L.SetField(slip, "set", L.NewFunction(func(L *lua.LState) int {
        x := int(L.CheckNumber(1))
        y := int(L.CheckNumber(2))
        char := L.CheckString(3)
        r := []rune(char)
        if len(r) > 0 {
            g.drawCmds = append(g.drawCmds, drawCmd{
                typ:   "set",
                x:     x,
                y:     y,
                char:  r[0],
                style: g.currentStyle,
            })
        }
        return 0
    }))

    // slip.text(x, y, str)
    L.SetField(slip, "text", L.NewFunction(func(L *lua.LState) int {
        x := int(L.CheckNumber(1))
        y := int(L.CheckNumber(2))
        text := L.CheckString(3)
        g.drawCmds = append(g.drawCmds, drawCmd{
            typ:   "text",
            x:     x,
            y:     y,
            text:  text,
            style: g.currentStyle,
        })
        return 0
    }))

    // slip.box(x, y, w, h)
    L.SetField(slip, "box", L.NewFunction(func(L *lua.LState) int {
        x := int(L.CheckNumber(1))
        y := int(L.CheckNumber(2))
        w := int(L.CheckNumber(3))
        h := int(L.CheckNumber(4))
        g.drawCmds = append(g.drawCmds, drawCmd{
            typ:   "box",
            x:     x,
            y:     y,
            w:     w,
            h:     h,
            style: g.currentStyle,
        })
        return 0
    }))

    // slip.color(fg) or slip.color(fg, bg)
    L.SetField(slip, "color", L.NewFunction(func(L *lua.LState) int {
        fg := L.CheckString(1)
        bg := L.OptString(2, "")
        g.currentStyle.FG = colorFromString(fg)
        if bg != "" {
            g.currentStyle.BG = colorFromString(bg)
        }
        return 0
    }))

    // slip.width()
    L.SetField(slip, "width", L.NewFunction(func(L *lua.LState) int {
        if g.ctx != nil {
            L.Push(lua.LNumber(g.ctx.Width))
        } else {
            L.Push(lua.LNumber(40))
        }
        return 1
    }))

    // slip.height()
    L.SetField(slip, "height", L.NewFunction(func(L *lua.LState) int {
        if g.ctx != nil {
            L.Push(lua.LNumber(g.ctx.Height))
        } else {
            L.Push(lua.LNumber(20))
        }
        return 1
    }))

    // slip.random(min, max) - helper for games
    L.SetField(slip, "random", L.NewFunction(func(L *lua.LState) int {
        min := int(L.CheckNumber(1))
        max := int(L.CheckNumber(2))
        L.Push(lua.LNumber(min + rand.Intn(max-min+1)))
        return 1
    }))
}

func (g *LuaGame) Render(screen *engine.Screen) error {
    g.currentScreen = screen
    g.drawCmds = nil // Clear previous frame

    // Call Lua render
    if err := g.callLuaMethod("render"); err != nil {
        return err
    }

    // Apply draw commands to screen
    for _, cmd := range g.drawCmds {
        switch cmd.typ {
        case "set":
            screen.Set(cmd.x, cmd.y, cmd.char, cmd.style)
        case "text":
            screen.DrawText(cmd.x, cmd.y, cmd.text, cmd.style)
        case "box":
            screen.DrawBox(cmd.x, cmd.y, cmd.w, cmd.h, engine.SingleBox, cmd.style)
        }
    }

    g.currentScreen = nil
    return nil
}

func colorFromString(s string) engine.Color {
    switch s {
    case "black":
        return engine.ColorBlack
    case "red":
        return engine.ColorRed
    case "green":
        return engine.ColorGreen
    case "yellow":
        return engine.ColorYellow
    case "blue":
        return engine.ColorBlue
    case "magenta":
        return engine.ColorMagenta
    case "cyan":
        return engine.ColorCyan
    case "white":
        return engine.ColorWhite
    default:
        return engine.ColorWhite
    }
}
```

**Acceptance Criteria:**
- [ ] slip.set() draws characters to screen
- [ ] slip.text() draws strings to screen
- [ ] slip.box() draws boxes
- [ ] slip.color() changes drawing color
- [ ] slip.clear() clears draw buffer
- [ ] Drawing commands properly buffered and applied

---

### Task 4B.4: Update Plugin Loader to Use Lua Runtime

**File:** `internal/plugins/loader.go`

Update `loadLuaPlugin`:

```go
func (l *Loader) loadLuaPlugin(manifest *Manifest, dir string) (engine.Game, error) {
    return NewLuaGame(manifest, dir)
}
```

**Acceptance Criteria:**
- [ ] Lua plugins load correctly
- [ ] Error handling for missing/invalid Lua files
- [ ] Plugin appears in game list after installation

---

### Task 4B.5: Create Example Lua Game

**File:** `examples/lua-pong/manifest.json`

```json
{
    "id": "lua-pong",
    "name": "Lua Pong",
    "description": "Simple pong game written in Lua",
    "author": "Slip Examples",
    "version": "1.0.0",
    "type": "lua",
    "entry": "game.lua",
    "category": "game",
    "min_width": 30,
    "min_height": 15
}
```

**File:** `examples/lua-pong/game.lua`

```lua
-- Lua Pong Example Game

local game = {}

-- Game state
local ball = { x = 0, y = 0, vx = 1, vy = 0.5 }
local paddle = { y = 0 }
local score = 0
local width, height = 40, 20

function game.init(ctx)
    width = ctx.width
    height = ctx.height
    ball.x = width / 2
    ball.y = height / 2
    paddle.y = height / 2
    score = 0
end

function game.start()
    -- Nothing to do
end

function game.stop()
    -- Nothing to do
end

function game.update(dt)
    -- Move ball
    ball.x = ball.x + ball.vx
    ball.y = ball.y + ball.vy

    -- Bounce off top/bottom
    if ball.y <= 1 or ball.y >= height - 2 then
        ball.vy = -ball.vy
    end

    -- Bounce off right wall
    if ball.x >= width - 2 then
        ball.vx = -ball.vx
    end

    -- Check paddle collision
    if ball.x <= 3 and ball.y >= paddle.y - 2 and ball.y <= paddle.y + 2 then
        ball.vx = -ball.vx
        score = score + 1
    end

    -- Ball out of bounds (left)
    if ball.x < 1 then
        -- Reset
        ball.x = width / 2
        ball.y = height / 2
        ball.vx = 1
        score = 0
    end
end

function game.render()
    slip.clear()

    -- Draw border
    slip.box(0, 0, slip.width(), slip.height())

    -- Draw score
    slip.color("cyan")
    slip.text(slip.width() / 2 - 5, 0, "Score: " .. score)

    -- Draw paddle
    slip.color("green")
    for i = -2, 2 do
        slip.set(2, math.floor(paddle.y) + i, "|")
    end

    -- Draw ball
    slip.color("yellow")
    slip.set(math.floor(ball.x), math.floor(ball.y), "o")
end

function game.input(key)
    if key == "up" and paddle.y > 3 then
        paddle.y = paddle.y - 1
    elseif key == "down" and paddle.y < height - 4 then
        paddle.y = paddle.y + 1
    end
end

return game
```

**Acceptance Criteria:**
- [ ] Example game loads and runs
- [ ] Ball bounces correctly
- [ ] Paddle moves with up/down
- [ ] Score increments
- [ ] Can be installed to plugins directory

---

### Task 4B.6: Add Lua Plugin Test

**File:** `internal/plugins/lua_runtime_test.go`

```go
package plugins

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLuaGameLoad(t *testing.T) {
    // Create temp directory with test game
    dir := t.TempDir()

    // Write manifest
    manifest := `{
        "id": "test-game",
        "name": "Test Game",
        "description": "Test",
        "author": "Test",
        "version": "1.0.0",
        "type": "lua",
        "entry": "game.lua",
        "category": "game"
    }`
    os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(manifest), 0644)

    // Write simple game
    game := `
    local game = {}
    function game.init(ctx) end
    function game.start() end
    function game.stop() end
    function game.update(dt) end
    function game.render() slip.clear() end
    function game.input(key) end
    return game
    `
    os.WriteFile(filepath.Join(dir, "game.lua"), []byte(game), 0644)

    // Load manifest
    m := &Manifest{
        ID:    "test-game",
        Name:  "Test Game",
        Type:  "lua",
        Entry: "game.lua",
    }

    // Create game
    g, err := NewLuaGame(m, dir)
    if err != nil {
        t.Fatalf("Failed to create Lua game: %v", err)
    }
    defer g.Stop()

    // Test Info
    info := g.Info()
    if info.ID != "test-game" {
        t.Errorf("Expected ID 'test-game', got '%s'", info.ID)
    }
}

func TestLuaDrawAPI(t *testing.T) {
    // Test that drawing commands work
    // ...
}
```

**Acceptance Criteria:**
- [ ] Test loads a simple Lua game
- [ ] Test verifies game info
- [ ] Test verifies draw API works
- [ ] All tests pass

---

## Phase 5: Store System Completion

### Task 5.1: Implement Checksum Verification

**File:** `internal/store/client.go`

```go
import (
    "crypto/sha256"
    "encoding/hex"
    "io"
    "os"
)

func (c *Client) verifyChecksum(filePath, expectedChecksum string) error {
    if expectedChecksum == "" {
        return nil // No checksum to verify
    }

    // Remove "sha256:" prefix if present
    expected := expectedChecksum
    if len(expected) > 7 && expected[:7] == "sha256:" {
        expected = expected[7:]
    }

    // Calculate file checksum
    f, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return err
    }

    actual := hex.EncodeToString(h.Sum(nil))

    if actual != expected {
        return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
    }

    return nil
}

// Update Download to verify
func (c *Client) Download(game *GameEntry, destPath string) error {
    // ... existing download code ...

    // After writing file, verify checksum
    plat := platform.GetPlatformString()
    if expectedChecksum, ok := game.Checksum[plat]; ok && expectedChecksum != "" {
        if err := c.verifyChecksum(tmpFile, expectedChecksum); err != nil {
            os.Remove(tmpFile) // Remove bad file
            return fmt.Errorf("checksum verification failed: %w", err)
        }
    }

    // ... rest of function
}
```

**Acceptance Criteria:**
- [ ] SHA256 checksum verification implemented
- [ ] Bad downloads are rejected and cleaned up
- [ ] Verification is skipped if no checksum provided
- [ ] Works with "sha256:" prefix and without

---

### Task 5.2: Add Download Progress Display

**File:** `internal/store/client.go`

```go
// ProgressWriter wraps an io.Writer to track progress
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

func (c *Client) DownloadWithProgress(game *GameEntry, destPath string, progress func(pct int)) error {
    url, err := c.GetDownloadURL(game)
    if err != nil {
        return err
    }

    // ... setup code ...

    resp, err := c.httpClient.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

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
    out, err := os.Create(tmpFile)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, io.TeeReader(resp.Body, pw))
    if err != nil {
        return err
    }

    // ... rest of function
}
```

**File:** `internal/cli/cli.go`

Update install command:

```go
func storeInstallCmd(client *store.Client, loader *plugins.Loader, id string) int {
    // ... existing checks ...

    fmt.Printf("Installing '%s' (%s)...\n", game.Name, game.Version)

    // Download with progress
    lastPct := -1
    err = client.DownloadWithProgress(game, "", func(pct int) {
        if pct != lastPct && pct % 10 == 0 {
            fmt.Printf("\rDownloading: %d%%", pct)
            lastPct = pct
        }
    })
    fmt.Println() // New line after progress

    if err != nil {
        fmt.Fprintf(os.Stderr, "Installation failed: %s\n", err)
        return 1
    }

    fmt.Printf("Successfully installed '%s'\n", game.Name)
    return 0
}
```

**Acceptance Criteria:**
- [ ] Progress percentage shown during download
- [ ] Updates at reasonable intervals (not every byte)
- [ ] Final 100% shown
- [ ] Works with unknown content length (no progress)

---

### Task 5.3: Implement In-App Store Browser

**File:** `internal/tui/store_screen.go` (new file)

Create a browsable store screen:

```go
package tui

import (
    "slip/internal/store"
)

type StoreScreen struct {
    client      *store.Client
    games       []store.GameEntry
    selectedIdx int
    loading     bool
    error       string

    // For custom TUI
    width, height int
}

func NewStoreScreen(client *store.Client) *StoreScreen {
    return &StoreScreen{
        client: client,
    }
}

func (s *StoreScreen) Load() {
    s.loading = true
    s.error = ""

    games, err := s.client.ListGames()
    if err != nil {
        s.error = err.Error()
        s.loading = false
        return
    }

    s.games = games
    s.loading = false
}

func (s *StoreScreen) HandleInput(input engine.Input) bool {
    switch input.Key {
    case engine.KeyUp:
        if s.selectedIdx > 0 {
            s.selectedIdx--
        }
        return true
    case engine.KeyDown:
        if s.selectedIdx < len(s.games)-1 {
            s.selectedIdx++
        }
        return true
    case engine.KeyEnter:
        if len(s.games) > 0 {
            s.installSelected()
        }
        return true
    }
    return false
}

func (s *StoreScreen) installSelected() {
    if s.selectedIdx >= len(s.games) {
        return
    }

    game := s.games[s.selectedIdx]
    // Trigger installation...
}

func (s *StoreScreen) Render(screen *engine.Screen) {
    if s.loading {
        screen.DrawTextCentered(s.height/2, "Loading store...", engine.Style{})
        return
    }

    if s.error != "" {
        screen.DrawTextCentered(s.height/2, "Error: "+s.error, engine.Style{FG: engine.ColorRed})
        screen.DrawTextCentered(s.height/2+2, "Press Enter to retry", engine.Style{})
        return
    }

    // Draw title
    screen.DrawTextCentered(1, "Store", engine.Style{Bold: true})

    // Draw games list
    startY := 4
    for i, game := range s.games {
        if startY+i*2 >= s.height-3 {
            break
        }

        style := engine.Style{FG: engine.ColorWhite}
        prefix := "  "
        if i == s.selectedIdx {
            style = engine.Style{FG: engine.ColorYellow, Bold: true}
            prefix = "> "
        }

        screen.DrawText(2, startY+i*2, prefix+game.Name, style)
        screen.DrawText(4, startY+i*2+1, game.Description, engine.Style{FG: engine.ColorBrightBlack})
    }

    // Draw footer
    footer := "[↑↓] Navigate  [Enter] Install  [Esc] Back"
    screen.DrawTextCentered(s.height-2, footer, engine.Style{FG: engine.ColorBrightBlack})
}
```

Update `app.go` to use the new store screen instead of placeholder.

**Acceptance Criteria:**
- [ ] Store screen shows list of available games
- [ ] Can navigate with up/down
- [ ] Selected item highlighted
- [ ] Enter triggers installation
- [ ] Loading state shown while fetching
- [ ] Error state shown on fetch failure
- [ ] Replaces "Coming Soon" placeholder

---

### Task 5.4: Register Installed Plugins in Game Menu

**File:** `internal/games/registry.go`

Add method to register plugins:

```go
func (r *Registry) RegisterPlugin(manifest *plugins.Manifest, loader *plugins.Loader) {
    r.factories[manifest.ID] = func() engine.Game {
        game, err := loader.Load(manifest.ID)
        if err != nil {
            return nil
        }
        return game
    }
    r.order = append(r.order, manifest.ID)
    r.info[manifest.ID] = engine.GameInfo{
        ID:          manifest.ID,
        Name:        manifest.Name,
        Description: manifest.Description,
        Author:      manifest.Author,
        Version:     manifest.Version,
    }
}
```

**File:** `internal/app/app.go`

Scan and register plugins on startup:

```go
func New(width, height, fps int) *App {
    // ... existing code ...

    // Load plugins
    loader := plugins.NewLoader()
    if err := loader.Scan(); err == nil {
        for _, manifest := range loader.List() {
            if manifest.Category == "game" {
                app.gameRegistry.RegisterPlugin(manifest, loader)
            } else if manifest.Category == "animation" {
                app.animRegistry.RegisterPlugin(manifest, loader)
            }
        }
    }

    app.initMenus() // Rebuild menus with plugins
    return app
}
```

**Acceptance Criteria:**
- [ ] Installed plugins appear in game/animation menus
- [ ] Plugins can be launched from menu
- [ ] Built-in games still work
- [ ] Menu updates after new plugin installed (on restart)

---

## Phase 6: Testing

### Task 6.1: Game Logic Tests

**File:** `internal/games/snake/snake_test.go`

```go
package snake

import (
    "testing"
    "slip/internal/engine"
)

func TestSnakeInit(t *testing.T) {
    game := New()
    ctx := &engine.GameContext{
        Width:  40,
        Height: 20,
    }

    err := game.Init(ctx)
    if err != nil {
        t.Fatalf("Init failed: %v", err)
    }

    // Snake should start in center with length 3
    if len(game.snake) != 3 {
        t.Errorf("Expected snake length 3, got %d", len(game.snake))
    }
}

func TestSnakeMovement(t *testing.T) {
    game := New()
    ctx := &engine.GameContext{Width: 40, Height: 20}
    game.Init(ctx)

    startHead := game.snake[0]

    // Move right (default direction)
    game.moveTimer = game.moveDelay // Force update
    game.Update(0.01)

    newHead := game.snake[0]
    if newHead.X != startHead.X+1 {
        t.Errorf("Snake should move right")
    }
}

func TestSnakeCollision(t *testing.T) {
    game := New()
    ctx := &engine.GameContext{Width: 40, Height: 20}
    game.Init(ctx)

    // Move snake to wall
    game.snake[0] = engine.Point{X: 0, Y: 5}
    game.direction = Left
    game.moveTimer = game.moveDelay
    game.Update(0.01)

    if !game.gameOver {
        t.Error("Expected game over on wall collision")
    }
}

func TestSnakeFoodCollision(t *testing.T) {
    game := New()
    ctx := &engine.GameContext{Width: 40, Height: 20}
    game.Init(ctx)

    // Place food in front of snake
    head := game.snake[0]
    game.food = engine.Point{X: head.X + 1, Y: head.Y}

    initialLen := len(game.snake)
    initialScore := game.score

    game.moveTimer = game.moveDelay
    game.Update(0.01)

    if len(game.snake) != initialLen+1 {
        t.Error("Snake should grow after eating")
    }
    if game.score <= initialScore {
        t.Error("Score should increase after eating")
    }
}
```

Create similar tests for:
- `internal/games/pong/pong_test.go`
- `internal/games/breakout/breakout_test.go`
- `internal/games/blockfall/blockfall_test.go`

**Acceptance Criteria:**
- [ ] Snake tests pass
- [ ] Pong tests pass
- [ ] Breakout tests pass
- [ ] Blockfall tests pass
- [ ] All game tests cover init, movement, collision

---

### Task 6.2: Input Parsing Tests

**File:** `internal/engine/input_test.go`

```go
package engine

import (
    "testing"
)

func TestParseArrowKeys(t *testing.T) {
    tests := []struct {
        input    []byte
        expected Key
    }{
        {[]byte{0x1b, '[', 'A'}, KeyUp},
        {[]byte{0x1b, '[', 'B'}, KeyDown},
        {[]byte{0x1b, '[', 'C'}, KeyRight},
        {[]byte{0x1b, '[', 'D'}, KeyLeft},
    }

    for _, tt := range tests {
        result := parseEscapeSequence(tt.input)
        if result.Key != tt.expected {
            t.Errorf("Input %v: expected %v, got %v", tt.input, tt.expected, result.Key)
        }
    }
}

func TestParseRegularKeys(t *testing.T) {
    tests := []struct {
        input    byte
        expected Key
        rune     rune
    }{
        {13, KeyEnter, 0},
        {27, KeyEscape, 0},
        {32, KeySpace, 0},
        {'a', KeyRune, 'a'},
        {'q', KeyRune, 'q'},
    }

    for _, tt := range tests {
        result := parseKeypress(tt.input)
        if result.Key != tt.expected {
            t.Errorf("Input %d: expected key %v, got %v", tt.input, tt.expected, result.Key)
        }
        if tt.rune != 0 && result.Rune != tt.rune {
            t.Errorf("Input %d: expected rune %c, got %c", tt.input, tt.rune, result.Rune)
        }
    }
}
```

**Acceptance Criteria:**
- [ ] Arrow key escape sequences parsed correctly
- [ ] Enter, Escape, Space parsed correctly
- [ ] Regular characters parsed as KeyRune
- [ ] Tests pass on Windows and Unix

---

### Task 6.3: Score Persistence Tests

**File:** `internal/state/state_test.go`

```go
package state

import (
    "os"
    "path/filepath"
    "testing"
    "time"
)

func TestScoreBoardSaveLoad(t *testing.T) {
    // Use temp directory
    tmpDir := t.TempDir()
    oldCacheDir := os.Getenv("SLIP_CACHE_DIR")
    os.Setenv("SLIP_CACHE_DIR", tmpDir)
    defer os.Setenv("SLIP_CACHE_DIR", oldCacheDir)

    // Create scoreboard and add scores
    sb := NewScoreBoard()
    sb.AddScore("snake", Score{
        Value:     100,
        Player:    "Test",
        Timestamp: time.Now(),
        GameID:    "snake",
    })

    err := sb.Save()
    if err != nil {
        t.Fatalf("Failed to save: %v", err)
    }

    // Load in new instance
    sb2, err := LoadScores()
    if err != nil {
        t.Fatalf("Failed to load: %v", err)
    }

    scores := sb2.GetTopScores("snake", 10)
    if len(scores) != 1 {
        t.Errorf("Expected 1 score, got %d", len(scores))
    }
    if scores[0].Value != 100 {
        t.Errorf("Expected score 100, got %d", scores[0].Value)
    }
}

func TestScoreBoardTop10(t *testing.T) {
    sb := NewScoreBoard()

    // Add 15 scores
    for i := 0; i < 15; i++ {
        sb.AddScore("snake", Score{
            Value:  i * 10,
            Player: "Test",
            GameID: "snake",
        })
    }

    // Should only keep top 10
    scores := sb.GetTopScores("snake", 10)
    if len(scores) != 10 {
        t.Errorf("Expected 10 scores, got %d", len(scores))
    }

    // Highest should be first
    if scores[0].Value != 140 {
        t.Errorf("Expected highest score 140, got %d", scores[0].Value)
    }
}

func TestConfigSaveLoad(t *testing.T) {
    tmpDir := t.TempDir()
    oldConfigDir := os.Getenv("SLIP_CONFIG_DIR")
    os.Setenv("SLIP_CONFIG_DIR", tmpDir)
    defer os.Setenv("SLIP_CONFIG_DIR", oldConfigDir)

    cfg := Config{
        Theme:      "gruvbox",
        PlayerName: "TestPlayer",
        TUIMode:    "bubbletea",
    }

    err := SaveConfig(cfg)
    if err != nil {
        t.Fatalf("Failed to save: %v", err)
    }

    loaded, err := LoadConfig()
    if err != nil {
        t.Fatalf("Failed to load: %v", err)
    }

    if loaded.Theme != "gruvbox" {
        t.Errorf("Expected theme 'gruvbox', got '%s'", loaded.Theme)
    }
    if loaded.TUIMode != "bubbletea" {
        t.Errorf("Expected tui_mode 'bubbletea', got '%s'", loaded.TUIMode)
    }
}
```

**Acceptance Criteria:**
- [ ] Scores save and load correctly
- [ ] Top 10 limit enforced
- [ ] Scores sorted by value descending
- [ ] Config saves and loads correctly
- [ ] Tests use temp directories

---

### Task 6.4: Store Client Tests

**File:** `internal/store/client_test.go`

```go
package store

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestFetchRegistry(t *testing.T) {
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{
            "games": [
                {
                    "id": "test-game",
                    "name": "Test Game",
                    "description": "A test game",
                    "author": "Test",
                    "version": "1.0.0"
                }
            ],
            "featured": ["test-game"],
            "updated_at": "2025-01-01T00:00:00Z"
        }`))
    }))
    defer server.Close()

    client := &Client{
        registryURL: server.URL,
        httpClient:  http.DefaultClient,
        cacheTTL:    0, // Disable cache for testing
    }

    registry, err := client.FetchRegistry()
    if err != nil {
        t.Fatalf("Failed to fetch: %v", err)
    }

    if len(registry.Games) != 1 {
        t.Errorf("Expected 1 game, got %d", len(registry.Games))
    }
    if registry.Games[0].ID != "test-game" {
        t.Errorf("Expected ID 'test-game', got '%s'", registry.Games[0].ID)
    }
}

func TestSearchGames(t *testing.T) {
    client := &Client{
        cache: &Registry{
            Games: []GameEntry{
                {ID: "snake", Name: "Snake", Description: "Classic snake"},
                {ID: "pong", Name: "Pong", Description: "Classic pong"},
                {ID: "snake-pro", Name: "Snake Pro", Description: "Enhanced snake"},
            },
        },
    }

    results, _ := client.SearchGames("snake")
    if len(results) != 2 {
        t.Errorf("Expected 2 results, got %d", len(results))
    }
}

func TestChecksumVerification(t *testing.T) {
    // Create temp file
    tmpDir := t.TempDir()
    tmpFile := filepath.Join(tmpDir, "test.txt")
    os.WriteFile(tmpFile, []byte("hello world"), 0644)

    client := &Client{}

    // Correct checksum for "hello world"
    correctSum := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

    err := client.verifyChecksum(tmpFile, correctSum)
    if err != nil {
        t.Errorf("Valid checksum should pass: %v", err)
    }

    err = client.verifyChecksum(tmpFile, "invalid")
    if err == nil {
        t.Error("Invalid checksum should fail")
    }
}
```

**Acceptance Criteria:**
- [ ] Registry fetch tested with mock server
- [ ] Search functionality tested
- [ ] Checksum verification tested
- [ ] All tests pass

---

## Phase 7: Documentation

### Task 7.1: Create User README

**File:** `README.md`

```markdown
# Slip

Terminal Entertainment Platform - Games and animations for your terminal.

## Features

- **4 Built-in Games**: Snake, Pong, Breakout, Blockfall
- **5 Animations**: Clipster, Starfield, Matrix, DVD Bounce, Fire
- **tmux Integration**: Run alongside your work in a split pane
- **Plugin System**: Install community games from the store
- **Themes**: 5 color themes (default, gruvbox, nord, dracula, monokai)
- **High Scores**: Persistent score tracking

## Installation

### From Source

```bash
git clone https://github.com/yourusername/slip
cd slip
go build -o slip ./cmd/slip
```

### Binary Releases

Download from [Releases](https://github.com/yourusername/slip/releases).

## Quick Start

```bash
# Run in current terminal
slip

# Or summon in a tmux pane (if in tmux)
slip summon

# Play a specific game
slip play snake

# Watch an animation
slip animate starfield
```

## Usage

### Commands

| Command | Description |
|---------|-------------|
| `slip` | Start interactive menu |
| `slip run` | Run in current terminal |
| `slip summon` | Open in tmux pane |
| `slip hide` | Hide tmux pane |
| `slip play <game>` | Play a specific game |
| `slip animate <anim>` | Run an animation |
| `slip games` | List available games |
| `slip animations` | List animations |
| `slip scores` | View high scores |
| `slip store list` | Browse available plugins |
| `slip store install <id>` | Install a plugin |
| `slip config` | View/set configuration |

### Configuration

```bash
# Set theme
slip config set theme gruvbox

# Set player name
slip config set player_name "YourName"

# Switch TUI mode
slip config set tui_mode bubbletea  # or 'custom'
```

### Available Themes

- `default` - Cyan/bright colors
- `gruvbox` - Warm yellows/reds
- `nord` - Cool blues/cyans
- `dracula` - Magenta/purple
- `monokai` - Bright green

## Creating Plugins

See [PLUGIN_GUIDE.md](PLUGIN_GUIDE.md) for creating Lua games.

## License

MIT
```

**Acceptance Criteria:**
- [ ] README covers installation
- [ ] README covers basic usage
- [ ] README lists all commands
- [ ] README explains configuration
- [ ] README has examples

---

### Task 7.2: Create Plugin Development Guide

**File:** `PLUGIN_GUIDE.md`

```markdown
# Slip Plugin Development Guide

## Overview

Slip supports Lua-based plugins for games and animations.

## Getting Started

### Directory Structure

```
my-game/
├── manifest.json    # Required: Plugin metadata
└── game.lua         # Required: Game logic
```

### Manifest Format

```json
{
    "id": "my-game",
    "name": "My Game",
    "description": "A fun game",
    "author": "Your Name",
    "version": "1.0.0",
    "type": "lua",
    "entry": "game.lua",
    "category": "game",
    "min_width": 30,
    "min_height": 15
}
```

## Lua API

### Game Structure

```lua
local game = {}

function game.init(ctx)
    -- ctx.width: screen width
    -- ctx.height: screen height
end

function game.start()
    -- Called when game starts
end

function game.stop()
    -- Called when game stops
end

function game.update(dt)
    -- dt: delta time in seconds
    -- Update game logic here
end

function game.render()
    -- Draw to screen here
end

function game.input(key)
    -- key: "up", "down", "left", "right", "enter", "escape", "space", or character
end

return game
```

### Drawing Functions

```lua
slip.clear()              -- Clear draw buffer
slip.set(x, y, char)      -- Draw character at position
slip.text(x, y, str)      -- Draw string at position
slip.box(x, y, w, h)      -- Draw box outline
slip.color(fg)            -- Set foreground color
slip.color(fg, bg)        -- Set foreground and background color
slip.width()              -- Get screen width
slip.height()             -- Get screen height
slip.random(min, max)     -- Random integer between min and max
```

### Available Colors

- `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`

## Example: Simple Game

```lua
local game = {}

local player = { x = 10, y = 10 }
local score = 0

function game.init(ctx)
    player.x = ctx.width / 2
    player.y = ctx.height / 2
end

function game.update(dt)
    -- Game logic
end

function game.render()
    slip.clear()
    slip.box(0, 0, slip.width(), slip.height())
    slip.color("cyan")
    slip.text(2, 0, "Score: " .. score)
    slip.color("green")
    slip.set(player.x, player.y, "@")
end

function game.input(key)
    if key == "up" then player.y = player.y - 1
    elseif key == "down" then player.y = player.y + 1
    elseif key == "left" then player.x = player.x - 1
    elseif key == "right" then player.x = player.x + 1
    end
end

return game
```

## Installing Your Plugin

1. Copy your plugin folder to `~/.cache/slip/games/`
2. Restart Slip
3. Your game appears in the Games menu

## Testing

```bash
# Install locally
cp -r my-game ~/.cache/slip/games/

# Run Slip
slip

# Or run directly
slip play my-game
```
```

**Acceptance Criteria:**
- [ ] Guide explains manifest format
- [ ] Guide documents all Lua API functions
- [ ] Guide has complete example
- [ ] Guide explains installation

---

## Phase 8: Distribution

### Task 8.1: Create Build Script

**File:** `scripts/build.sh`

```bash
#!/bin/bash

VERSION=${1:-"dev"}
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

mkdir -p dist

for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%/*}"
    ARCH="${PLATFORM#*/}"
    OUTPUT="dist/slip-${VERSION}-${OS}-${ARCH}"

    if [ "$OS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "Building for ${OS}/${ARCH}..."
    GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w -X main.Version=${VERSION}" -o "$OUTPUT" ./cmd/slip

    if [ $? -eq 0 ]; then
        echo "  Created: $OUTPUT"
    else
        echo "  Failed: $OUTPUT"
    fi
done

echo "Done!"
```

**File:** `scripts/build.ps1` (Windows)

```powershell
param(
    [string]$Version = "dev"
)

$platforms = @(
    @{OS="linux"; Arch="amd64"},
    @{OS="linux"; Arch="arm64"},
    @{OS="darwin"; Arch="amd64"},
    @{OS="darwin"; Arch="arm64"},
    @{OS="windows"; Arch="amd64"}
)

New-Item -ItemType Directory -Force -Path "dist" | Out-Null

foreach ($p in $platforms) {
    $output = "dist/slip-$Version-$($p.OS)-$($p.Arch)"
    if ($p.OS -eq "windows") {
        $output += ".exe"
    }

    Write-Host "Building for $($p.OS)/$($p.Arch)..."
    $env:GOOS = $p.OS
    $env:GOARCH = $p.Arch
    go build -ldflags="-s -w -X main.Version=$Version" -o $output ./cmd/slip

    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Created: $output"
    } else {
        Write-Host "  Failed: $output"
    }
}

Write-Host "Done!"
```

**Acceptance Criteria:**
- [ ] Build script creates binaries for all platforms
- [ ] Version embedded in binary
- [ ] Binaries are stripped (-s -w flags)
- [ ] Works on both Unix and Windows

---

### Task 8.2: Create GitHub Actions Workflow

**File:** `.github/workflows/release.yml`

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Build binaries
        run: |
          chmod +x scripts/build.sh
          ./scripts/build.sh ${{ github.ref_name }}

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*
          generate_release_notes: true
```

**Acceptance Criteria:**
- [ ] Workflow triggers on version tags
- [ ] Builds all platform binaries
- [ ] Creates GitHub release
- [ ] Attaches binaries to release

---

## Summary Checklist

### Phase 4A: TUI Abstraction
- [ ] Task 4A.1: Define TUI Interface
- [ ] Task 4A.2: Create TUI Factory
- [ ] Task 4A.3: Wrap Existing Custom TUI
- [ ] Task 4A.4: Implement Bubble Tea TUI
- [ ] Task 4A.5: Add TUI Mode to Config
- [ ] Task 4A.6: Update App to Use TUI Factory
- [ ] Task 4A.7: Add CLI Flag Override

### Phase 4B: Lua Plugin System
- [ ] Task 4B.1: Add Lua Dependency
- [ ] Task 4B.2: Create Lua Game Runtime
- [ ] Task 4B.3: Implement Lua Drawing API
- [ ] Task 4B.4: Update Plugin Loader
- [ ] Task 4B.5: Create Example Lua Game
- [ ] Task 4B.6: Add Lua Plugin Test

### Phase 5: Store Completion
- [ ] Task 5.1: Checksum Verification
- [ ] Task 5.2: Download Progress Display
- [ ] Task 5.3: In-App Store Browser
- [ ] Task 5.4: Register Plugins in Menu

### Phase 6: Testing
- [ ] Task 6.1: Game Logic Tests
- [ ] Task 6.2: Input Parsing Tests
- [ ] Task 6.3: Score Persistence Tests
- [ ] Task 6.4: Store Client Tests

### Phase 7: Documentation
- [ ] Task 7.1: User README
- [ ] Task 7.2: Plugin Development Guide

### Phase 8: Distribution
- [ ] Task 8.1: Build Script
- [ ] Task 8.2: GitHub Actions Workflow

---

## Task Execution Order (Recommended)

1. **Start with Phase 4A** (TUI Abstraction) - This enables comparing Bubble Tea vs Custom
2. **Then Phase 4B** (Lua Plugins) - Core feature for extensibility
3. **Then Phase 5** (Store) - Depends on plugin system
4. **Phase 6** (Tests) - Can be done in parallel with features
5. **Phase 7** (Docs) - After features stabilize
6. **Phase 8** (Distribution) - Final step

Each task is designed to be completed independently within 1-4 hours by a developer familiar with Go.
