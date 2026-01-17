package ui

import (
	"fmt"
	"strings"
)

// QuitCommand - quit the application
type QuitCommand struct {
	BaseCommand
}

func NewQuitCommand() *QuitCommand {
	return &QuitCommand{
		BaseCommand: BaseCommand{
			name:        "quit",
			description: "Exit Slip",
			usage:       ":quit or :q",
			aliases:     []string{"q", "exit"},
		},
	}
}

func (c *QuitCommand) Execute(ctx *CommandContext, args []string) error {
	// Signal quit - implementation in TUI will handle this
	if ctx.SetScreen != nil {
		ctx.SetScreen("quit")
	}
	return nil
}

// ThemeCommand - change theme
type ThemeCommand struct {
	BaseCommand
	themes []string
}

func NewThemeCommand(themes []string) *ThemeCommand {
	return &ThemeCommand{
		BaseCommand: BaseCommand{
			name:        "theme",
			description: "Change color theme",
			usage:       ":theme <name>",
			aliases:     []string{"t"},
		},
		themes: themes,
	}
}

func (c *ThemeCommand) Complete(args []string) []string {
	if len(args) == 0 {
		return c.themes
	}
	var matches []string
	for _, theme := range c.themes {
		if strings.HasPrefix(theme, args[0]) {
			matches = append(matches, theme)
		}
	}
	return matches
}

func (c *ThemeCommand) Execute(ctx *CommandContext, args []string) error {
	if len(args) == 0 {
		// List themes
		msg := "Available themes: " + strings.Join(c.themes, ", ")
		ctx.Output(msg, nil)
		return nil
	}

	themeName := args[0]
	if ctx.SetTheme != nil {
		return ctx.SetTheme(themeName)
	}
	return fmt.Errorf("theme change not available")
}

// PlayCommand - launch a game
type PlayCommand struct {
	BaseCommand
	games []string
}

func NewPlayCommand(games []string) *PlayCommand {
	return &PlayCommand{
		BaseCommand: BaseCommand{
			name:        "play",
			description: "Launch a game",
			usage:       ":play <game>",
			aliases:     []string{"p"},
		},
		games: games,
	}
}

func (c *PlayCommand) Complete(args []string) []string {
	if len(args) == 0 {
		return c.games
	}
	var matches []string
	for _, game := range c.games {
		if strings.HasPrefix(game, args[0]) {
			matches = append(matches, game)
		}
	}
	return matches
}

func (c *PlayCommand) Execute(ctx *CommandContext, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: :play <game>")
	}

	gameID := args[0]
	if ctx.LaunchGame != nil {
		return ctx.LaunchGame(gameID)
	}
	return fmt.Errorf("game launch not available")
}

// AnimateCommand - launch an animation
type AnimateCommand struct {
	BaseCommand
	animations []string
}

func NewAnimateCommand(animations []string) *AnimateCommand {
	return &AnimateCommand{
		BaseCommand: BaseCommand{
			name:        "animate",
			description: "Launch an animation",
			usage:       ":animate <animation>",
			aliases:     []string{"anim", "a"},
		},
		animations: animations,
	}
}

func (c *AnimateCommand) Complete(args []string) []string {
	if len(args) == 0 {
		return c.animations
	}
	var matches []string
	for _, anim := range c.animations {
		if strings.HasPrefix(anim, args[0]) {
			matches = append(matches, anim)
		}
	}
	return matches
}

func (c *AnimateCommand) Execute(ctx *CommandContext, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: :animate <animation>")
	}

	animID := args[0]
	if ctx.LaunchAnim != nil {
		return ctx.LaunchAnim(animID)
	}
	return fmt.Errorf("animation launch not available")
}

// NavigationCommand - navigate to screens
type NavigationCommand struct {
	BaseCommand
	screen string
}

func NewNavigationCommand(name, screen, desc string) *NavigationCommand {
	return &NavigationCommand{
		BaseCommand: BaseCommand{
			name:        name,
			description: desc,
			usage:       ":" + name,
			aliases:     []string{},
		},
		screen: screen,
	}
}

func (c *NavigationCommand) Execute(ctx *CommandContext, args []string) error {
	if ctx.SetScreen != nil {
		ctx.SetScreen(c.screen)
	}
	return nil
}

// HelpCommand - show help
type HelpCommand struct {
	BaseCommand
	registry *CommandRegistry
}

func NewHelpCommand(registry *CommandRegistry) *HelpCommand {
	return &HelpCommand{
		BaseCommand: BaseCommand{
			name:        "help",
			description: "Show command help",
			usage:       ":help [command]",
			aliases:     []string{"h", "?"},
		},
		registry: registry,
	}
}

func (c *HelpCommand) Execute(ctx *CommandContext, args []string) error {
	if len(args) == 0 {
		// List all commands
		var cmds []string
		for _, cmd := range c.registry.List() {
			cmds = append(cmds, fmt.Sprintf(":%s - %s", cmd.Name(), cmd.Description()))
		}
		msg := "Available commands:\n" + strings.Join(cmds, "\n")
		ctx.Output(msg, nil)
		return nil
	}

	// Show help for specific command
	cmdName := args[0]
	if cmd, ok := c.registry.Get(cmdName); ok {
		msg := fmt.Sprintf("Command: %s\nDescription: %s\nUsage: %s",
			cmd.Name(), cmd.Description(), cmd.Usage())
		if len(cmd.Aliases()) > 0 {
			msg += fmt.Sprintf("\nAliases: %s", strings.Join(cmd.Aliases(), ", "))
		}
		ctx.Output(msg, nil)
		return nil
	}

	return fmt.Errorf("unknown command: %s", cmdName)
}

// ClearCommand - clear screen
type ClearCommand struct {
	BaseCommand
}

func NewClearCommand() *ClearCommand {
	return &ClearCommand{
		BaseCommand: BaseCommand{
			name:        "clear",
			description: "Clear screen",
			usage:       ":clear",
			aliases:     []string{"cls"},
		},
	}
}

func (c *ClearCommand) Execute(ctx *CommandContext, args []string) error {
	// The TUI will handle actually clearing the screen
	return nil
}
