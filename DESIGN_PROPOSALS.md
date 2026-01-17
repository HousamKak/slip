# Slip Design Proposals

This document outlines proposed improvements for Slip's Custom TUI aesthetics and a new integrated terminal command system.

---

## Table of Contents

1. [Custom UI Visual Improvements](#1-custom-ui-visual-improvements)
2. [Integrated Terminal Command System](#2-integrated-terminal-command-system)
3. [Windows Support & Behavior](#3-windows-support--behavior)

---

## 1. Custom UI Visual Improvements

The current Custom TUI is functional but visually plain. Here are proposals to make it more appealing.

### 1.1 Current Problems

| Issue | Description |
|-------|-------------|
| **Flat appearance** | No visual depth, borders feel basic |
| **Poor spacing** | Elements cramped together |
| **Inconsistent styling** | Some areas styled, others plain text |
| **No visual feedback** | Selections don't "pop" |
| **Missing decorations** | No ASCII art, logos, or visual interest |
| **Monochrome feel** | Not enough color contrast |

### 1.2 Proposed Box Drawing Improvements

**Current:**
```
┌─────────────────────┐
│ S L I P             │
├─────────────────────┤
│ > Games             │
│   Animations        │
│   Store             │
│   Settings          │
│   Quit              │
└─────────────────────┘
```

**Proposed - Double Border with Shadows:**
```
╔═══════════════════════════════════╗
║           ░▒▓ S L I P ▓▒░         ║
║         Terminal Entertainment     ║
╠═══════════════════════════════════╣
║                                   ║
║   ▸ Games          [4 available]  ║
║     Animations     [5 available]  ║
║     Store          [Browse]       ║
║     Settings       [Configure]    ║
║     Help           [?]            ║
║     Quit           [Esc]          ║
║                                   ║
╠═══════════════════════════════════╣
║  ↑↓ Navigate  Enter Select  Q Quit ║
╚═══════════════════════════════════╝
```

### 1.3 Color Scheme Enhancements

**Proposed Color Usage:**

```
┌─────────────────────────────────────────┐
│ Element          │ Color                │
├─────────────────────────────────────────┤
│ Title            │ Bold + Primary       │
│ Selected item    │ Inverse/Highlight BG │
│ Unselected items │ Normal foreground    │
│ Hints/shortcuts  │ Dim/Gray             │
│ Borders          │ Accent color         │
│ Descriptions     │ Cyan/Secondary       │
│ Errors           │ Red + Bold           │
│ Success          │ Green                │
│ Loading          │ Yellow/Animated      │
└─────────────────────────────────────────┘
```

### 1.4 ASCII Art Logo

Add a stylized header logo:

```
    _____ _      _____ _____
   / ____| |    |_   _|  __ \
  | (___ | |      | | | |__) |
   \___ \| |      | | |  ___/
   ____) | |____ _| |_| |
  |_____/|______|_____|_|

   Terminal Entertainment Platform
```

Or a more compact version:

```
  ╭━━━╮╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱
  ┃╭━╮┃╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱
  ┃╰━━┳╮╭┳━━╮╱╱╱╱╱╱╱╱╱
  ╰━━╮┃┃┃┃╭╮┃╱╱╱╱╱╱╱╱╱
  ┃╰━╯┃╰╯┃╰╯┃╱╱╱╱╱╱╱╱╱
  ╰━━━┻━━┫╭━╯╱╱╱╱╱╱╱╱╱
  ╱╱╱╱╱╱╱┃┃╱╱╱╱╱╱╱╱╱╱╱
  ╱╱╱╱╱╱╱╰╯╱╱╱╱╱╱╱╱╱╱╱
```

### 1.5 Selection Highlighting Styles

**Option A - Full Row Highlight:**
```
   Games
 ▸ Animations    ◀━━━━ Selected (inverse colors)
   Store
```

**Option B - Bracket Selection:**
```
   Games
 [ Animations ]  ◀━━━━ Selected
   Store
```

**Option C - Arrow + Underline:**
```
   Games
 ▸ Animations
   ‾‾‾‾‾‾‾‾‾‾
   Store
```

**Option D - Box Selection (Recommended):**
```
   Games

 ┌─────────────┐
 │ Animations  │  ◀━━ Selected item in box
 └─────────────┘

   Store
```

### 1.6 Menu Item Layout Improvements

**Current Layout:**
```
> Snake
  Pong
  Breakout
```

**Proposed Layout with Details:**
```
╭─────────────────────────────────────────╮
│  🐍 Snake                               │
│     Classic snake game - eat & grow     │
│     High Score: 1,250    ★★★☆☆         │
╰─────────────────────────────────────────╯
   🏓 Pong
      AI paddle opponent
      High Score: 15

   🧱 Breakout
      Break all the bricks
      High Score: 5,400
```

### 1.7 Animated Elements

Add subtle animations for visual appeal:

1. **Loading Spinner:**
   ```
   ⠋ Loading...
   ⠙ Loading...
   ⠹ Loading...
   ⠸ Loading...
   ⠼ Loading...
   ⠴ Loading...
   ⠦ Loading...
   ⠧ Loading...
   ⠇ Loading...
   ⠏ Loading...
   ```

2. **Pulsing Selection:**
   - Alternate between `▸` and `▹` every 500ms

3. **Scrolling Title:**
   - Marquee effect for long titles

4. **Fade Transitions:**
   - Brief dim when switching screens

### 1.8 Footer Status Bar

Add a persistent footer:

```
╔═══════════════════════════════════════════════════════╗
║                      SLIP v1.0                        ║
╠═══════════════════════════════════════════════════════╣
║                                                       ║
║                    (menu content)                     ║
║                                                       ║
╠═══════════════════════════════════════════════════════╣
║ ↑↓ Navigate │ Enter Select │ : Command │ Q Quit      ║
╚═══════════════════════════════════════════════════════╝
```

### 1.9 Implementation Priorities

| Priority | Feature | Effort |
|----------|---------|--------|
| P0 | Better borders (double-line) | Low |
| P0 | Improved color contrast | Low |
| P1 | Selection highlighting (inverse) | Low |
| P1 | Footer with controls | Medium |
| P1 | Proper spacing/padding | Low |
| P2 | ASCII art logo | Low |
| P2 | Loading spinner | Medium |
| P3 | Animated elements | High |
| P3 | Emoji support (with fallback) | Medium |

---

## 2. Integrated Terminal Command System

### 2.1 Concept

Add an in-app command terminal that can be invoked with `:` (like Vim) to execute Slip commands without leaving the TUI.

```
╔═══════════════════════════════════════╗
║            S L I P                    ║
╠═══════════════════════════════════════╣
║                                       ║
║   ▸ Games                             ║
║     Animations                        ║
║     Store                             ║
║                                       ║
╠═══════════════════════════════════════╣
║ : theme gruvbox█                      ║  ◀━━ Command input
╚═══════════════════════════════════════╝
```

### 2.2 Activation

- Press `:` to enter command mode (Vim-style)
- Press `/` for search mode
- Press `Escape` to cancel
- Press `Enter` to execute

### 2.3 Proposed Commands

#### Navigation Commands
```
:games              - Go to games menu
:animations         - Go to animations menu
:store              - Go to store
:settings           - Go to settings
:home               - Go to home screen
:back               - Go back one screen
```

#### Game/Animation Commands
```
:play snake         - Launch snake directly
:play pong          - Launch pong directly
:animate matrix     - Launch matrix animation
:animate starfield  - Launch starfield
```

#### Configuration Commands
```
:theme <name>       - Change theme (gruvbox, nord, dracula, etc.)
:theme              - List available themes
:fps <number>       - Set FPS (15-60)
:size <WxH>         - Set terminal size (e.g., :size 80x24)
:tui <mode>         - Switch TUI mode (custom, bubbletea)
:player <name>      - Set player name for scores
```

#### Store Commands
```
:install <id>       - Install a game/animation from store
:uninstall <id>     - Remove installed plugin
:update             - Check for updates
:search <query>     - Search store
:refresh            - Refresh store cache
```

#### Information Commands
```
:scores             - Show high scores
:scores <game>      - Show scores for specific game
:info <game>        - Show game/animation info
:version            - Show Slip version
:help               - Show command help
:help <command>     - Show help for specific command
```

#### System Commands
```
:quit / :q          - Quit Slip
:restart            - Restart Slip
:clear              - Clear screen
:debug              - Toggle debug mode
:fps-counter        - Toggle FPS display
```

### 2.4 Command Autocomplete

Tab completion for commands:

```
:th█
   ├─ theme
   └─ [Tab to complete]

:theme █
   ├─ default
   ├─ gruvbox
   ├─ nord
   ├─ dracula
   └─ monokai
```

### 2.5 Command History

- `↑` / `↓` to navigate command history
- History persisted across sessions
- `:history` to view all past commands

### 2.6 Search Mode

Press `/` to search within current context:

```
/sna█

Results:
  ▸ Snake (game)

Press Enter to select, Esc to cancel
```

### 2.7 Visual Design for Command Mode

**Inactive (normal mode):**
```
╠═══════════════════════════════════════╣
║ Press : for commands, / to search     ║
╚═══════════════════════════════════════╝
```

**Active (command mode):**
```
╠═══════════════════════════════════════╣
║ :█                                    ║
╚═══════════════════════════════════════╝
```

**With input:**
```
╠═══════════════════════════════════════╣
║ :theme gruvbox█                       ║
╚═══════════════════════════════════════╝
```

**With autocomplete:**
```
╠═══════════════════════════════════════╣
║ :theme gru█                           ║
║ ┌─────────────────────┐               ║
║ │ gruvbox (suggested) │               ║
║ └─────────────────────┘               ║
╚═══════════════════════════════════════╝
```

### 2.8 Implementation Architecture

```
┌─────────────────────────────────────────┐
│              TUI Layer                  │
├─────────────────────────────────────────┤
│         Command Parser                  │
│  ┌─────────────────────────────────┐    │
│  │ Input: ":theme gruvbox"         │    │
│  │ Parsed: {cmd: "theme",          │    │
│  │          args: ["gruvbox"]}     │    │
│  └─────────────────────────────────┘    │
├─────────────────────────────────────────┤
│         Command Registry                │
│  ┌─────────────────────────────────┐    │
│  │ "theme" -> ThemeCommand         │    │
│  │ "play"  -> PlayCommand          │    │
│  │ "quit"  -> QuitCommand          │    │
│  │ ...                             │    │
│  └─────────────────────────────────┘    │
├─────────────────────────────────────────┤
│         Command Executor                │
│  ┌─────────────────────────────────┐    │
│  │ Execute command                 │    │
│  │ Return result/error             │    │
│  │ Update TUI state                │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

### 2.9 Command Interface

```go
// Command represents an executable command
type Command interface {
    Name() string
    Description() string
    Usage() string
    Aliases() []string
    Complete(args []string) []string  // For autocomplete
    Execute(args []string) error
}

// CommandRegistry manages available commands
type CommandRegistry struct {
    commands map[string]Command
}

// Example command implementation
type ThemeCommand struct {
    tui *CustomTUI
}

func (c *ThemeCommand) Name() string { return "theme" }
func (c *ThemeCommand) Description() string { return "Change color theme" }
func (c *ThemeCommand) Usage() string { return ":theme <name>" }
func (c *ThemeCommand) Aliases() []string { return []string{"t"} }
func (c *ThemeCommand) Complete(args []string) []string {
    return []string{"default", "gruvbox", "nord", "dracula", "monokai"}
}
func (c *ThemeCommand) Execute(args []string) error {
    if len(args) == 0 {
        // List themes
        return nil
    }
    return c.tui.SetTheme(args[0])
}
```

---

## 3. Windows Support & Behavior

### 3.1 Current Windows Behavior

When running `slip summon` on Windows (without WSL):

**What Happens:**
```
┌─────────────────────────────────────────────┐
│  PowerShell Window                          │
├─────────────────────────────────────────────┤
│  ┌───────────────────────────────────────┐  │
│  │  Your normal terminal (top)           │  │
│  │  > commands here                      │  │
│  │                                       │  │
│  ├───────────────────────────────────────┤  │
│  │  Slip running (bottom)                │  │
│  │  ╔═══════════════════════════════╗    │  │
│  │  ║         S L I P               ║    │  │
│  │  ╚═══════════════════════════════╝    │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

**This is actually a cool feature!** It creates a split view where:
- Top: Your normal PowerShell/terminal for commands
- Bottom: Slip running for entertainment

### 3.2 Why Dock Doesn't Work on Windows

**The Issue:**

On Linux/macOS with tmux:
- `slip dock` creates a tmux split pane
- tmux manages the window splitting natively
- Works seamlessly in any terminal

On Windows without WSL:
- No tmux available
- Windows Terminal has its own pane system (wt.exe)
- PowerShell/cmd.exe don't support pane splitting

**Current Behavior:**
- Opens a new Windows Terminal tab/window with Slip
- Cannot dock to the side of existing terminal
- Works but not as integrated as tmux

### 3.3 Windows Terminal Integration

Windows Terminal (`wt.exe`) DOES support panes! Current implementation in `platform/platform.go`:

```go
func LaunchWindowsTerminal(command []string, dock string) error {
    // Current: Opens new tab
    args := []string{"new-tab", "--title", "Slip"}

    // With pane support:
    switch dock {
    case "bottom":
        args = []string{"split-pane", "-V"}  // Vertical split
    case "right":
        args = []string{"split-pane", "-H"}  // Horizontal split
    case "left":
        args = []string{"split-pane", "-H", "-d", "."}
    }
}
```

**Requirement:** Must be running INSIDE Windows Terminal already for splits to work.

### 3.4 Windows Support Matrix

| Feature | Windows Terminal | PowerShell/CMD | WSL |
|---------|-----------------|----------------|-----|
| Run Slip | ✅ | ✅ | ✅ |
| Full colors | ✅ | ⚠️ Limited | ✅ |
| Unicode/box drawing | ✅ | ⚠️ Depends | ✅ |
| Dock right | ✅ (if inside WT) | ❌ | ✅ (tmux) |
| Dock bottom | ✅ (if inside WT) | ❌ | ✅ (tmux) |
| Summon | ✅ New window | ✅ New window | ✅ tmux |
| Input responsiveness | ✅ | ✅ | ✅ |

### 3.5 Improving Windows Experience

#### Option 1: Detect Windows Terminal

```go
func InWindowsTerminal() bool {
    return os.Getenv("WT_SESSION") != ""
}

func Dock(position string) error {
    if runtime.GOOS == "windows" {
        if InWindowsTerminal() {
            // Use wt.exe split-pane
            return LaunchWindowsTerminalPane(position)
        }
        // Fallback: inform user
        return fmt.Errorf("docking requires Windows Terminal; use 'slip summon' instead")
    }
    // Linux/macOS: use tmux
    return TmuxDock(position)
}
```

#### Option 2: Windows Terminal JSON Profile

Users can add Slip to Windows Terminal settings:

```json
{
    "profiles": {
        "list": [
            {
                "guid": "{your-guid-here}",
                "name": "Slip",
                "commandline": "slip run",
                "icon": "🎮",
                "startingDirectory": "%USERPROFILE%",
                "hidden": false
            }
        ]
    },
    "actions": [
        {
            "command": {
                "action": "splitPane",
                "split": "horizontal",
                "profile": "Slip"
            },
            "keys": "ctrl+shift+s"
        }
    ]
}
```

#### Option 3: ConEmu/Cmder Support

For users with ConEmu or Cmder:

```go
func HasConEmu() bool {
    return os.Getenv("ConEmuPID") != ""
}

func ConEmuSplit(command string) error {
    // ConEmu supports -new_console:s[H|V] for splits
    args := fmt.Sprintf("-new_console:sH %s", command)
    return exec.Command("cmd", "/c", args).Start()
}
```

### 3.6 Recommended Windows Workflow

**For best experience on Windows:**

1. **Install Windows Terminal** (from Microsoft Store)
2. **Run Slip inside Windows Terminal:**
   ```powershell
   # In Windows Terminal
   slip run
   ```

3. **Use keyboard shortcut for pane:**
   - `Alt+Shift+Plus` - Split pane right
   - `Alt+Shift+Minus` - Split pane down
   - Run `slip run` in new pane

4. **Or use slip summon (opens new window):**
   ```powershell
   slip summon
   ```

### 3.7 Proposed Windows Improvements

| Improvement | Description | Priority |
|-------------|-------------|----------|
| WT pane detection | Detect if in Windows Terminal, enable proper docking | P0 |
| Better fallback message | Explain why dock doesn't work outside WT | P0 |
| WT profile generator | Command to generate Windows Terminal profile JSON | P1 |
| ConEmu support | Add ConEmu split support | P2 |
| Installer with WT integration | Create installer that adds Slip to WT | P3 |

### 3.8 Error Messages Improvement

**Current (confusing):**
```
Error: docking failed
```

**Proposed (helpful):**
```
╭─────────────────────────────────────────────────╮
│  Docking not available                          │
├─────────────────────────────────────────────────┤
│  Docking requires:                              │
│  • Windows Terminal (wt.exe) on Windows         │
│  • tmux on Linux/macOS                          │
│                                                 │
│  Alternatives:                                  │
│  • Run 'slip summon' to open in new window      │
│  • Use Windows Terminal's built-in pane split:  │
│    Alt+Shift+Plus (right) or Alt+Shift+- (down) │
│  • Install tmux via WSL for full dock support   │
╰─────────────────────────────────────────────────╯
```

---

## 4. Implementation Roadmap

### Phase 1: Visual Polish (1-2 days)
- [ ] Double-line borders
- [ ] Better color contrast
- [ ] Inverse selection highlighting
- [ ] Footer with controls
- [ ] Proper padding/margins

### Phase 2: Command System Core (2-3 days)
- [ ] Command input mode (`:`)
- [ ] Basic command parser
- [ ] Core commands (theme, quit, play)
- [ ] Command execution

### Phase 3: Command System Advanced (2-3 days)
- [ ] Tab autocomplete
- [ ] Command history
- [ ] Search mode (`/`)
- [ ] All commands implemented

### Phase 4: Windows Improvements (1-2 days)
- [ ] Windows Terminal pane detection
- [ ] Proper dock support in WT
- [ ] Better error messages
- [ ] WT profile generator command

---

## 5. Quick Wins (Can Do Now)

These can be implemented immediately with minimal effort:

1. **Change border style** - Just change `┌` to `╔` etc.
2. **Add padding** - Add empty lines around content
3. **Inverse selection** - Use `\033[7m` ANSI code
4. **Footer bar** - Add static help text at bottom
5. **Better WT detection** - Check `WT_SESSION` env var

---

## Appendix A: Box Drawing Characters Reference

```
Single:  ┌ ┐ └ ┘ ─ │ ├ ┤ ┬ ┴ ┼
Double:  ╔ ╗ ╚ ╝ ═ ║ ╠ ╣ ╦ ╩ ╬
Rounded: ╭ ╮ ╰ ╯ ─ │
Heavy:   ┏ ┓ ┗ ┛ ━ ┃ ┣ ┫ ┳ ┻ ╋
Mixed:   ╒ ╕ ╘ ╛ ╞ ╡ ╤ ╧ ╪
```

## Appendix B: ANSI Color Codes

```
Reset:      \033[0m
Bold:       \033[1m
Dim:        \033[2m
Italic:     \033[3m
Underline:  \033[4m
Blink:      \033[5m
Inverse:    \033[7m
Hidden:     \033[8m
Strike:     \033[9m

Foreground: \033[30-37m (standard), \033[90-97m (bright)
Background: \033[40-47m (standard), \033[100-107m (bright)
256-color:  \033[38;5;Nm (fg), \033[48;5;Nm (bg)
RGB:        \033[38;2;R;G;Bm (fg), \033[48;2;R;G;Bm (bg)
```

## Appendix C: Spinner Characters

```
Dots:     ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
Line:     - \ | /
Bounce:   ⠁ ⠂ ⠄ ⠂
Arrows:   ← ↖ ↑ ↗ → ↘ ↓ ↙
Circle:   ◐ ◓ ◑ ◒
Square:   ◰ ◳ ◲ ◱
Block:    ▖ ▘ ▝ ▗
```
