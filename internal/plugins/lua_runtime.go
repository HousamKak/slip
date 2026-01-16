package plugins

import (
	"fmt"
	"math/rand"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
	"slip/internal/engine"
)

// drawCmd represents a drawing command
type drawCmd struct {
	typ   string // "set", "text", "box", "clear"
	x, y  int
	w, h  int
	char  rune
	text  string
	style engine.Style
}

// LuaGame wraps a Lua script as a Game
type LuaGame struct {
	state    *lua.LState
	manifest *Manifest
	dir      string

	// Game context
	ctx *engine.GameContext

	// Drawing
	drawCmds     []drawCmd
	currentStyle engine.Style
}

// NewLuaGame creates a new Lua game from a manifest
func NewLuaGame(manifest *Manifest, dir string) (*LuaGame, error) {
	L := lua.NewState()

	game := &LuaGame{
		state:    L,
		manifest: manifest,
		dir:      dir,
		currentStyle: engine.Style{
			FG: engine.ColorWhite,
			BG: engine.ColorBlack,
		},
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

	// slip.random(min, max)
	L.SetField(slip, "random", L.NewFunction(func(L *lua.LState) int {
		min := int(L.CheckNumber(1))
		max := int(L.CheckNumber(2))
		L.Push(lua.LNumber(min + rand.Intn(max-min+1)))
		return 1
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
	return g.callLuaMethod("init", ctx)
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
			screen.DrawBox(cmd.x, cmd.y, cmd.w, cmd.h, engine.BoxStyleSingle, cmd.style)
		}
	}

	return nil
}

// HandleInput handles input
func (g *LuaGame) HandleInput(input engine.Input) error {
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
