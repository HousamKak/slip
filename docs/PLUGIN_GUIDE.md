# Slip Plugin Development Guide

Create your own games and animations for Slip using Lua!

## Table of Contents

- [Getting Started](#getting-started)
- [Quick Example](#quick-example)
- [Plugin Structure](#plugin-structure)
- [Lua API Reference](#lua-api-reference)
- [Game Lifecycle](#game-lifecycle)
- [Drawing Functions](#drawing-functions)
- [Input Handling](#input-handling)
- [Best Practices](#best-practices)
- [Example Games](#example-games)
- [Debugging](#debugging)
- [Publishing](#publishing)

## Getting Started

### Prerequisites

- Slip installed and working
- Basic Lua knowledge (very simple language!)
- Text editor

### Your First Plugin

1. **Create plugin directory:**

```bash
mkdir -p ~/.cache/slip/games/my-first-game
cd ~/.cache/slip/games/my-first-game
```

2. **Create `manifest.json`:**

```json
{
  "id": "my-first-game",
  "name": "My First Game",
  "description": "A simple bouncing ball game",
  "author": "Your Name",
  "version": "1.0.0",
  "type": "lua",
  "entry": "game.lua",
  "category": "game"
}
```

3. **Create `game.lua`:**

```lua
local game = {}

-- Initialize game
function game.init(ctx)
  game.x = ctx.width / 2
  game.y = ctx.height / 2
  game.dx = 1
  game.dy = 1
end

-- Start game (called once)
function game.start()
  -- Optional: setup code
end

-- Update game logic (called every frame)
function game.update(dt)
  game.x = game.x + game.dx
  game.y = game.y + game.dy

  -- Bounce off edges
  if game.x <= 0 or game.x >= slip.width() - 1 then
    game.dx = -game.dx
  end
  if game.y <= 0 or game.y >= slip.height() - 1 then
    game.dy = -game.dy
  end
end

-- Draw the game (called every frame)
function game.render()
  slip.clear()
  slip.set(math.floor(game.x), math.floor(game.y), '●')
end

-- Handle input (called when key pressed)
function game.input(key)
  if key == "q" then
    -- Handle quit
  end
end

-- Stop game (called on exit)
function game.stop()
  -- Cleanup if needed
end

return game
```

4. **Test it:**

```bash
# Restart Slip to load the plugin
slip run

# Your game will appear in the menu!
```

## Quick Example

Here's a complete mini-game (bouncing ball you can control):

```lua
local game = {}

function game.init(ctx)
  game.width = ctx.width
  game.height = ctx.height
  game.x = ctx.width / 2
  game.y = ctx.height / 2
  game.vx = 0
  game.vy = 0
  game.score = 0
  game.food = {
    x = slip.random(2, ctx.width - 2),
    y = slip.random(2, ctx.height - 2)
  }
end

function game.update(dt)
  -- Apply velocity
  game.x = game.x + game.vx * dt * 20
  game.y = game.y + game.vy * dt * 20

  -- Friction
  game.vx = game.vx * 0.95
  game.vy = game.vy * 0.95

  -- Boundary check
  if game.x < 1 then game.x = 1 end
  if game.x > game.width - 2 then game.x = game.width - 2 end
  if game.y < 1 then game.y = 1 end
  if game.y > game.height - 2 then game.y = game.height - 2 end

  -- Check food collision
  local px = math.floor(game.x)
  local py = math.floor(game.y)
  if px == game.food.x and py == game.food.y then
    game.score = game.score + 1
    game.food.x = slip.random(2, game.width - 2)
    game.food.y = slip.random(2, game.height - 2)
  end
end

function game.render()
  slip.clear()

  -- Draw border
  slip.box(0, 0, game.width, game.height)

  -- Draw player
  slip.color("green")
  slip.set(math.floor(game.x), math.floor(game.y), '●')

  -- Draw food
  slip.color("red")
  slip.set(game.food.x, game.food.y, '◆')

  -- Draw score
  slip.color("white")
  slip.text(2, 0, "Score: " .. game.score)
end

function game.input(key)
  local speed = 2
  if key == "up" or key == "w" then
    game.vy = -speed
  elseif key == "down" or key == "s" then
    game.vy = speed
  elseif key == "left" or key == "a" then
    game.vx = -speed
  elseif key == "right" or key == "d" then
    game.vx = speed
  end
end

return game
```

## Plugin Structure

### Directory Layout

```
~/.cache/slip/games/your-game/
├── manifest.json          # Required: Plugin metadata
├── game.lua              # Required: Entry point
├── utils.lua             # Optional: Helper functions
└── assets/               # Optional: Data files
    └── levels.lua
```

### Manifest Schema

```json
{
  "id": "unique-game-id",           // Required: lowercase, hyphens only
  "name": "Display Name",           // Required: shown in menu
  "description": "Short description", // Required: what the game does
  "author": "Your Name",            // Required: creator name
  "version": "1.0.0",               // Required: semantic version
  "type": "lua",                    // Required: must be "lua"
  "entry": "game.lua",              // Required: entry point file
  "category": "game"                // Required: "game" or "animation"
}
```

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier (lowercase, hyphens) |
| `name` | string | Display name (shown in menu) |
| `description` | string | Brief description |
| `author` | string | Your name or username |
| `version` | string | Semantic version (1.0.0) |
| `type` | string | Must be "lua" |
| `entry` | string | Entry point filename |
| `category` | string | "game" or "animation" |

## Lua API Reference

### Screen Information

```lua
-- Get screen dimensions
width = slip.width()    -- Returns terminal width
height = slip.height()  -- Returns terminal height
```

### Drawing Functions

```lua
-- Clear the screen (call at start of render)
slip.clear()

-- Set a single character at position
slip.set(x, y, char)
slip.set(10, 5, '█')
slip.set(10, 5, 'A')

-- Draw text string
slip.text(x, y, string)
slip.text(5, 10, "Hello World!")
slip.text(0, 0, "Score: " .. score)

-- Draw a box outline
slip.box(x, y, width, height)
slip.box(0, 0, slip.width(), slip.height())  -- Full border
slip.box(5, 5, 10, 8)                        -- Small box

-- Set drawing color (affects subsequent draws)
slip.color(colorname)

-- Available colors:
-- "black", "red", "green", "yellow", "blue",
-- "magenta", "cyan", "white", "default"
-- Also: "bright_*" versions (bright_red, bright_green, etc.)

slip.color("red")
slip.text(10, 5, "ERROR")
slip.color("green")
slip.text(10, 6, "OK")
```

### Utility Functions

```lua
-- Generate random integer (inclusive)
num = slip.random(min, max)

x = slip.random(0, 10)      -- Random 0-10
y = slip.random(5, 15)      -- Random 5-15

-- Example: Random position
food_x = slip.random(1, slip.width() - 2)
food_y = slip.random(1, slip.height() - 2)
```

### Character Set

Use any ASCII or Unicode character:

```lua
-- Blocks
slip.set(x, y, '█')  -- Full block
slip.set(x, y, '▓')  -- Dark shade
slip.set(x, y, '▒')  -- Medium shade
slip.set(x, y, '░')  -- Light shade

-- Shapes
slip.set(x, y, '●')  -- Circle
slip.set(x, y, '■')  -- Square
slip.set(x, y, '◆')  -- Diamond
slip.set(x, y, '▲')  -- Triangle

-- Lines
slip.set(x, y, '─')  -- Horizontal
slip.set(x, y, '│')  -- Vertical
slip.set(x, y, '┌')  -- Corners
slip.set(x, y, '┐')
slip.set(x, y, '└')
slip.set(x, y, '┘')

-- Regular ASCII
slip.set(x, y, '*')
slip.set(x, y, '#')
slip.set(x, y, '@')
```

## Game Lifecycle

### Required Functions

Your game **must** implement these functions:

```lua
local game = {}

-- 1. Initialize (called once on load)
function game.init(ctx)
  -- ctx contains:
  --   ctx.width  - screen width
  --   ctx.height - screen height

  -- Initialize game state here
  game.player_x = ctx.width / 2
  game.player_y = ctx.height / 2
  game.score = 0
end

-- 2. Start (called when game starts)
function game.start()
  -- Optional: reset game state
  game.running = true
end

-- 3. Update (called every frame)
function game.update(dt)
  -- dt is delta time in seconds (typically 0.016 for 60 FPS)
  -- Update game logic here

  game.player_x = game.player_x + game.velocity * dt
end

-- 4. Render (called every frame after update)
function game.render()
  -- Clear screen first
  slip.clear()

  -- Draw everything
  slip.set(game.player_x, game.player_y, '●')
end

-- 5. Input (called when key is pressed)
function game.input(key)
  -- key is a string: "up", "down", "left", "right",
  -- "enter", "space", or single character "a", "b", etc.

  if key == "up" then
    game.player_y = game.player_y - 1
  end
end

-- 6. Stop (called when game exits)
function game.stop()
  -- Optional: cleanup
  game.running = false
end

-- Must return the game table
return game
```

### Execution Order

```
1. init(ctx)     - Once when plugin loads
2. start()       - When user starts playing
   Loop:
3.   update(dt)  - Game logic (30-60 times/sec)
4.   render()    - Drawing (30-60 times/sec)
5.   input(key)  - When key pressed
   End loop when user quits
6. stop()        - Cleanup
```

## Input Handling

### Available Keys

```lua
function game.input(key)
  -- Arrow keys
  if key == "up" then end
  if key == "down" then end
  if key == "left" then end
  if key == "right" then end

  -- Special keys
  if key == "enter" then end
  if key == "space" then end
  if key == "escape" then end

  -- Letter keys (lowercase)
  if key == "w" then end  -- WASD
  if key == "a" then end
  if key == "s" then end
  if key == "d" then end

  -- Number keys
  if key == "1" then end
  if key == "2" then end

  -- Other characters
  if key == "p" then  -- Pause
    game.paused = not game.paused
  end
end
```

### Input Patterns

**Continuous Movement (velocity-based):**

```lua
function game.input(key)
  local speed = 5
  if key == "up" then game.vy = -speed end
  if key == "down" then game.vy = speed end
  if key == "left" then game.vx = -speed end
  if key == "right" then game.vx = speed end
end

function game.update(dt)
  game.x = game.x + game.vx * dt
  game.y = game.y + game.vy * dt

  -- Apply friction
  game.vx = game.vx * 0.9
  game.vy = game.vy * 0.9
end
```

**Direct Movement (grid-based):**

```lua
function game.input(key)
  if key == "up" and game.can_move_up() then
    game.y = game.y - 1
  end
  if key == "down" and game.can_move_down() then
    game.y = game.y + 1
  end
end
```

## Best Practices

### Performance

```lua
-- ✅ Good: Local variables are faster
function game.update(dt)
  local w = slip.width()
  local h = slip.height()
  -- Use w and h multiple times
end

-- ❌ Bad: Calling function repeatedly
function game.update(dt)
  if x > slip.width() then end
  if y > slip.height() then end
end
```

### Collision Detection

```lua
-- Point collision
function game.check_collision(x1, y1, x2, y2)
  return math.floor(x1) == math.floor(x2) and
         math.floor(y1) == math.floor(y2)
end

-- Rectangle collision (AABB)
function game.box_collision(x1, y1, w1, h1, x2, y2, w2, h2)
  return x1 < x2 + w2 and
         x1 + w1 > x2 and
         y1 < y2 + h2 and
         y1 + h1 > y2
end

-- Boundary check
function game.in_bounds(x, y)
  return x >= 0 and x < slip.width() and
         y >= 0 and y < slip.height()
end
```

### Game State Management

```lua
local game = {}

-- Use state machine
game.state = "menu"  -- "menu", "playing", "paused", "gameover"

function game.update(dt)
  if game.state == "menu" then
    -- Menu logic
  elseif game.state == "playing" then
    -- Game logic
  elseif game.state == "gameover" then
    -- Show game over
  end
end

function game.input(key)
  if game.state == "menu" and key == "enter" then
    game.state = "playing"
    game.reset()
  end
end
```

### Code Organization

```lua
-- Separate concerns
local game = {}
local player = {}
local enemies = {}

function player.new(x, y)
  return {
    x = x,
    y = y,
    vx = 0,
    vy = 0,
    health = 100
  }
end

function player.update(p, dt)
  p.x = p.x + p.vx * dt
  p.y = p.y + p.vy * dt
end

function player.render(p)
  slip.color("green")
  slip.set(math.floor(p.x), math.floor(p.y), '●')
end

function game.init(ctx)
  game.player = player.new(ctx.width / 2, ctx.height / 2)
end

function game.update(dt)
  player.update(game.player, dt)
end

function game.render()
  slip.clear()
  player.render(game.player)
end

return game
```

## Example Games

### Example 1: Simple Catch Game

```lua
local game = {}

function game.init(ctx)
  game.width = ctx.width
  game.height = ctx.height
  game.player_x = ctx.width / 2
  game.items = {}
  game.score = 0
  game.spawn_timer = 0
end

function game.update(dt)
  -- Spawn falling items
  game.spawn_timer = game.spawn_timer + dt
  if game.spawn_timer > 1 then
    table.insert(game.items, {
      x = slip.random(2, game.width - 2),
      y = 1
    })
    game.spawn_timer = 0
  end

  -- Update items
  for i = #game.items, 1, -1 do
    local item = game.items[i]
    item.y = item.y + 5 * dt

    -- Check catch
    if math.floor(item.y) == game.height - 2 and
       math.floor(item.x) == math.floor(game.player_x) then
      game.score = game.score + 1
      table.remove(game.items, i)
    -- Check miss
    elseif item.y > game.height then
      table.remove(game.items, i)
    end
  end
end

function game.render()
  slip.clear()
  slip.box(0, 0, game.width, game.height)

  -- Player
  slip.color("green")
  slip.set(math.floor(game.player_x), game.height - 2, '═')

  -- Items
  slip.color("yellow")
  for _, item in ipairs(game.items) do
    slip.set(math.floor(item.x), math.floor(item.y), '◆')
  end

  -- Score
  slip.color("white")
  slip.text(2, 0, "Score: " .. game.score)
end

function game.input(key)
  if key == "left" and game.player_x > 1 then
    game.player_x = game.player_x - 1
  elseif key == "right" and game.player_x < game.width - 2 then
    game.player_x = game.player_x + 1
  end
end

return game
```

### Example 2: Maze Explorer

See `examples/lua-pong/game.lua` for a complete working example!

## Debugging

### Print Debugging

```lua
-- Print to title area
function game.render()
  slip.clear()
  slip.text(0, 0, "Debug: x=" .. game.x .. " y=" .. game.y)
  -- Rest of rendering
end
```

### Common Issues

**Game not appearing in menu:**
- Check manifest.json syntax (valid JSON?)
- Restart Slip after installing
- Verify files in `~/.cache/slip/games/your-game/`

**Nothing renders:**
- Did you call `slip.clear()` at start of `render()`?
- Check x,y coordinates are in bounds
- Verify `return game` at end of file

**Lua errors:**
- Check Slip output for error messages
- Verify all required functions exist
- Check for typos in function names

**Performance issues:**
- Reduce number of draw calls
- Cache `slip.width()` and `slip.height()`
- Avoid creating tables in `update()`/`render()`

## Publishing

### Preparing for Release

1. **Test thoroughly:**
   - Play through multiple times
   - Test edge cases
   - Verify on different terminal sizes

2. **Polish manifest.json:**
   ```json
   {
     "id": "awesome-game",
     "name": "Awesome Game",
     "description": "Detailed description with gameplay info",
     "author": "Your Name",
     "version": "1.0.0",
     "type": "lua",
     "entry": "game.lua",
     "category": "game"
   }
   ```

3. **Add README (optional):**
   Create `README.md` in your plugin directory explaining controls and gameplay.

4. **Share:**
   - Zip your plugin directory
   - Share on GitHub, forums, etc.
   - Submit to Slip plugin registry (coming soon!)

### Installation Instructions for Users

```bash
# Download plugin
cd ~/.cache/slip/games/
# Extract/copy your plugin directory here

# Restart Slip
slip run
```

## Advanced Topics

### Multiple Files

```lua
-- game.lua
local utils = require("utils")

local game = {}

function game.init(ctx)
  game.helpers = utils.create_helpers(ctx)
end

return game
```

```lua
-- utils.lua
local utils = {}

function utils.create_helpers(ctx)
  return {
    width = ctx.width,
    height = ctx.height
  }
end

return utils
```

### Animations vs Games

Set `"category": "animation"` in manifest for non-interactive content:

```lua
-- Animation example
local anim = {}

function anim.init(ctx)
  anim.particles = {}
  for i = 1, 50 do
    table.insert(anim.particles, {
      x = slip.random(0, ctx.width),
      y = slip.random(0, ctx.height),
      vx = slip.random(-2, 2),
      vy = slip.random(-2, 2)
    })
  end
end

function anim.update(dt)
  for _, p in ipairs(anim.particles) do
    p.x = (p.x + p.vx * dt) % slip.width()
    p.y = (p.y + p.vy * dt) % slip.height()
  end
end

function anim.render()
  slip.clear()
  slip.color("cyan")
  for _, p in ipairs(anim.particles) do
    slip.set(math.floor(p.x), math.floor(p.y), '*')
  end
end

function anim.input(key) end  -- Animations ignore input

return anim
```

## Resources

- [Lua 5.1 Reference Manual](https://www.lua.org/manual/5.1/)
- [Example Plugin](examples/lua-pong/) - Complete working game
- [Slip Documentation](README.md) - Main documentation

## Getting Help

- Check existing plugins for examples
- Read error messages carefully
- Start simple and build up
- Ask in Discussions or Issues

---

**Happy coding! Can't wait to see what you create! 🎮**
