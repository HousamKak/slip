# Slip - Terminal Entertainment Platform

<div align="center">

**Keep yourself entertained while waiting for long-running commands**

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

[Quick Start](#quick-start) • [Features](#features) • [Installation](#installation) • [Usage](#usage) • [Games](#games) • [Plugin Development](PLUGIN_GUIDE.md)

</div>

---

## What is Slip?

Slip is a terminal-based entertainment platform designed to keep you engaged while waiting for:
- 🤖 AI assistant responses
- 🏗️ Long build processes
- 🚀 Deployment pipelines
- 📦 Package installations
- Any command that makes you wait

Run Slip in a **tmux pane** alongside your work terminal, or launch it directly when you need a break.

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/slip.git
cd slip

# Build
go build ./cmd/slip

# Run Slip
./slip run
```

### Basic Usage

```bash
# Start Slip interactively (in current terminal)
slip run

# Play a specific game directly
slip play snake

# Run a specific animation
slip animate starfield

# List available content
slip games
slip animations

# View high scores
slip scores
slip scores snake
```

### Tmux Integration

Slip works great with tmux, running in a side pane while you work:

```bash
# Summon Slip in a tmux pane (right side by default)
slip summon

# Summon on a specific side
slip summon --dock bottom
slip summon --dock left

# Hide the Slip pane
slip hide

# Move Slip to different position
slip dock bottom
slip dock right
```

## Features

### 🎮 Four Built-in Games

| Game | Description | Controls |
|------|-------------|----------|
| **Snake** | Classic snake - eat food, grow longer | ←↑↓→ or WASD |
| **Pong** | Single player vs AI paddle game | ↑↓ or WS |
| **Breakout** | Break bricks with paddle and ball | ←→ or AD |
| **Blockfall** | Tetris-style block stacking | ←↑↓→, Space |

### ✨ Five Built-in Animations

- **Clipster** - Bouncing Clippy assistant
- **Starfield** - Flying through stars
- **Matrix** - Matrix-style falling characters
- **DVD Bounce** - Classic DVD logo bounce
- **Fire** - Animated fire effect

**Animation Controls:** Press `Q`, `Esc`, or `Ctrl+C` to exit

### 🎨 Five Color Themes

Choose from: `default`, `gruvbox`, `nord`, `dracula`, `monokai`

```bash
slip config set theme gruvbox
```

### 🔌 Plugin System

Create your own games using **Lua**! See [PLUGIN_GUIDE.md](PLUGIN_GUIDE.md) for details.

```bash
# Install a plugin from the store
slip store search pong
slip store install lua-pong

# List installed plugins
slip store installed
```

### 🎯 Two TUI Modes

Switch between two different terminal UI implementations:

- **Custom** - Hand-rolled ANSI rendering, minimal dependencies
- **Bubble Tea** - Modern Elm-architecture framework

```bash
slip run --tui bubbletea
slip run --tui custom

# Set default
slip config set tui_mode bubbletea
```

## Installation

### Download Pre-built Binaries (Recommended)

Download the latest release for your platform from the [GitHub Releases](https://github.com/yourusername/slip/releases) page.

**Linux:**
```bash
# Download and extract (replace VERSION with actual version, e.g., 1.0.0)
wget https://github.com/yourusername/slip/releases/download/vVERSION/slip-linux-amd64.tar.gz
tar -xzf slip-linux-amd64.tar.gz

# Make executable and move to PATH
chmod +x slip-linux-amd64
sudo mv slip-linux-amd64 /usr/local/bin/slip

# Verify installation
slip --help
```

**macOS:**
```bash
# Intel Mac
wget https://github.com/yourusername/slip/releases/download/vVERSION/slip-darwin-amd64.tar.gz
tar -xzf slip-darwin-amd64.tar.gz

# Apple Silicon Mac (M1/M2/M3)
wget https://github.com/yourusername/slip/releases/download/vVERSION/slip-darwin-arm64.tar.gz
tar -xzf slip-darwin-arm64.tar.gz

# Make executable and move to PATH
chmod +x slip-darwin-*
sudo mv slip-darwin-* /usr/local/bin/slip

# Verify installation
slip --help
```

**Windows:**
```powershell
# Download slip-windows-amd64.zip from releases page
# Extract the .exe file
# Add to PATH or run directly

.\slip-windows-amd64.exe --help
```

**Verify Checksums:**
```bash
# Download the .sha256 file and verify
sha256sum -c slip-linux-amd64.sha256
# or on macOS:
shasum -a 256 -c slip-darwin-arm64.sha256
```

### From Source

**Prerequisites:**
- Go 1.22 or higher
- Git
- Terminal with ANSI color support
- (Optional) tmux for pane integration

```bash
# Clone the repository
git clone https://github.com/yourusername/slip.git
cd slip

# Install dependencies
go mod download

# Build
go build ./cmd/slip

# (Optional) Install to PATH
go install ./cmd/slip

# Verify installation
slip --help
```

### Build from Source (All Platforms)

```bash
# Clone repository
git clone https://github.com/yourusername/slip.git
cd slip

# Build for all platforms using build scripts
./scripts/build.sh v1.0.0      # Unix/Linux/macOS
# or
.\scripts\build.ps1 v1.0.0     # Windows PowerShell

# Binaries will be in build/ directory
```

### Platform Support

| Platform | Status | Notes |
|----------|--------|-------|
| Linux | ✅ Fully Supported | Native terminal, tmux integration |
| macOS | ✅ Fully Supported | Native terminal, tmux integration |
| Windows | ⚠️ Partial | Works in Windows Terminal, WSL recommended |

## Usage

### Running Games

#### Interactive Menu

```bash
# Launch Slip's interactive menu
slip run

# Navigate with arrow keys or j/k
# Press Enter to select
# Press Esc or q to go back/quit
```

#### Direct Launch

```bash
# Play a specific game
slip play snake
slip play pong
slip play breakout
slip play blockfall

# Run an animation
slip animate matrix
slip animate starfield
slip animate fire
```

### Game Controls

**Universal Controls:**
- `P` - Pause/Unpause
- `R` - Restart game
- `Q` or `Esc` - Quit to menu

**Game-Specific:**

**Snake:**
- `←↑↓→` or `WASD` - Move direction
- Eat food (red dots) to grow
- Don't hit walls or yourself!

**Pong:**
- `↑↓` or `WS` - Move paddle
- Hit ball with paddle
- Don't let ball pass you!

**Breakout:**
- `←→` or `AD` - Move paddle
- Break all bricks to win
- Catch power-ups for bonuses

**Blockfall:**
- `←→` - Move block left/right
- `↓` - Soft drop
- `↑` - Rotate block
- `Space` - Hard drop
- Complete lines to score

### Configuration

```bash
# View all config
slip config

# Get specific value
slip config get theme

# Set values
slip config set theme nord
slip config set tui_mode bubbletea
slip config set player_name "YourName"
slip config set default_fps 60
```

**Available Config Options:**

| Option | Default | Description |
|--------|---------|-------------|
| `default_dock` | `right` | Tmux pane position: right, left, bottom, top |
| `default_size` | `40x20` | Terminal size: WIDTHxHEIGHT |
| `default_fps` | `30` | Frames per second |
| `theme` | `default` | Color theme |
| `tui_mode` | `custom` | UI framework: custom or bubbletea |
| `show_fps_counter` | `false` | Show FPS in games |
| `player_name` | `Player` | Your name for high scores |
| `registry_url` | (default) | Plugin registry URL |
| `auto_update` | `true` | Auto-check for plugin updates |

### High Scores

```bash
# View all high scores
slip scores

# View scores for specific game
slip scores snake
slip scores pong

# Scores are automatically saved when you beat them
```

### Store (Plugin Management)

```bash
# Browse available plugins
slip store list

# Search for plugins
slip store search tetris

# View plugin details
slip store info lua-pong

# Install a plugin
slip store install lua-pong

# View installed plugins
slip store installed

# Uninstall a plugin
slip store uninstall lua-pong
```

**Note:** After installing a plugin, restart Slip to see it in the menu.

### Tmux Integration

Slip can run in a tmux pane alongside your work:

```bash
# Start Slip in a tmux pane
slip summon

# Summon on different sides
slip summon --dock bottom
slip summon --dock left
slip summon --dock right

# Custom size (WIDTHxHEIGHT)
slip summon --size 60x30

# Hide (but keep running)
slip hide

# Move to different position
slip dock left
slip dock bottom

# Check status
slip status
```

**Workflow Example:**

```bash
# Terminal 1: Your work
$ npm run build

# Terminal 2: Summon Slip while build runs
$ slip summon --dock right

# Now you have:
# ┌─────────────────┬─────────────┐
# │ npm run build   │   SLIP      │
# │ [building...]   │  [ Snake ]  │
# │                 │  [  Play  ] │
# └─────────────────┴─────────────┘

# When build finishes, hide Slip
$ slip hide
```

## Games

### Snake 🐍

Classic snake game with smooth controls and increasing difficulty.

- **Objective:** Eat food to grow, avoid walls and yourself
- **Scoring:** 10 points per food
- **Speed:** Increases as you eat more
- **Controls:** Arrow keys or WASD

### Pong 🏓

Single-player pong against an AI opponent.

- **Objective:** Don't let the ball pass your paddle
- **Scoring:** 1 point per successful block
- **AI:** Gets harder as you score more
- **Controls:** Up/Down arrows or W/S

### Breakout 🧱

Classic brick breaker with power-ups.

- **Objective:** Break all bricks
- **Power-ups:**
  - Wide paddle
  - Fast ball
  - Slow ball
- **Lives:** 3 starting lives
- **Controls:** Left/Right arrows or A/D

### Blockfall 📦

Tetris-style falling blocks.

- **Objective:** Complete lines to score
- **Blocks:** 7 different tetromino shapes
- **Scoring:**
  - 1 line: 100 points
  - 2 lines: 300 points
  - 3 lines: 500 points
  - 4 lines: 800 points
- **Controls:** Arrows, Space for hard drop

## Themes

Slip includes 5 carefully crafted color themes:

```bash
# Switch themes
slip config set theme gruvbox
slip config set theme nord
slip config set theme dracula
slip config set theme monokai
slip config set theme default
```

Themes affect:
- Menu colors
- Game UI borders
- Text highlights
- Game elements

## Plugin Development

Want to create your own games? Slip supports **Lua plugins**!

See the complete [Plugin Development Guide](PLUGIN_GUIDE.md) for:
- Getting started with Lua games
- API reference
- Example code
- Best practices

### Quick Plugin Example

```lua
-- game.lua
local game = {}

function game.init(ctx)
  game.x = ctx.width / 2
  game.y = ctx.height / 2
end

function game.update(dt)
  -- Update game logic
end

function game.render()
  slip.set(game.x, game.y, '●')
end

function game.input(key)
  if key == "up" then game.y = game.y - 1 end
  if key == "down" then game.y = game.y + 1 end
end

return game
```

Install in `~/.cache/slip/games/mygame/` with a `manifest.json` and restart Slip!

## Troubleshooting

### Game not responding

- Check if terminal is in focus
- Try restarting Slip
- Verify terminal supports ANSI escape codes

### Tmux pane issues

```bash
# Kill any stuck Slip processes
pkill -f slip

# Check tmux sessions
tmux ls

# Kill specific tmux session if needed
tmux kill-session -t slip
```

### Plugins not showing

```bash
# Check plugins are installed
ls ~/.cache/slip/games/

# Verify manifest.json exists
cat ~/.cache/slip/games/mygame/manifest.json

# Restart Slip (plugins load at startup)
```

### Performance issues

```bash
# Lower FPS
slip config set default_fps 20

# Try custom TUI (lighter weight)
slip config set tui_mode custom
```

## Development

### Project Structure

```
slip/
├── cmd/slip/           # Entry point
├── internal/
│   ├── app/           # Main application
│   ├── tui/           # UI frameworks (Custom & Bubble Tea)
│   ├── engine/        # Game engine (30 FPS, rendering, input)
│   ├── games/         # Built-in games
│   ├── animations/    # Built-in animations
│   ├── plugins/       # Plugin system (Lua runtime)
│   ├── store/         # Plugin store client
│   ├── state/         # Config & score persistence
│   └── ui/            # UI components
└── examples/          # Example plugins
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/games/snake -v
go test ./internal/engine -v
go test ./internal/state -v

# Run with coverage
go test ./... -cover
```

### Building

```bash
# Development build
go build ./cmd/slip

# Production build with optimizations
go build -ldflags="-s -w" -o slip ./cmd/slip

# Build all platforms with version info using build scripts
./scripts/build.sh v1.0.0      # Creates binaries for all platforms
.\scripts\build.ps1 v1.0.0     # Windows PowerShell version

# Manual cross-compile for specific platform
GOOS=linux GOARCH=amd64 go build -o slip-linux ./cmd/slip
GOOS=darwin GOARCH=amd64 go build -o slip-macos ./cmd/slip
GOOS=windows GOARCH=amd64 go build -o slip.exe ./cmd/slip
```

### Release Process

Slip uses GitHub Actions for automated releases:

1. **Create and push a version tag:**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically:**
   - Runs all tests on Linux, macOS, and Windows
   - Builds binaries for all platforms (Linux, macOS Intel/ARM, Windows)
   - Generates SHA256 checksums
   - Creates release archives
   - Publishes GitHub release with binaries attached

3. **Users download from:** [GitHub Releases](https://github.com/yourusername/slip/releases)

See `.github/workflows/release.yml` for the complete workflow.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Credits

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - TUI styling
- [gopher-lua](https://github.com/yuin/gopher-lua) - Lua runtime

Inspired by classic terminal games and the need for productive procrastination.

## Support

- 📖 [Documentation](PLUGIN_GUIDE.md)
- 🐛 [Issue Tracker](https://github.com/yourusername/slip/issues)
- 💬 [Discussions](https://github.com/yourusername/slip/discussions)

---

<div align="center">

**Made with ❤️ for terminal enthusiasts**

[⬆ Back to Top](#slip---terminal-entertainment-platform)

</div>
