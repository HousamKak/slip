package ui

import (
	"fmt"
	"strings"
)

// Command represents an executable command
type Command interface {
	Name() string
	Description() string
	Usage() string
	Aliases() []string
	Complete(args []string) []string
	Execute(ctx *CommandContext, args []string) error
}

// CommandContext provides context for command execution
type CommandContext struct {
	TUI         interface{} // Reference to TUI for state changes
	Output      func(string, error)
	SetScreen   func(string)
	SetTheme    func(string) error
	LaunchGame  func(string) error
	LaunchAnim  func(string) error
	GetConfig   func(string) string
	SetConfig   func(string, string) error
}

// CommandRegistry manages available commands
type CommandRegistry struct {
	commands map[string]Command
	aliases  map[string]string
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
		aliases:  make(map[string]string),
	}
}

// Register registers a command
func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
	for _, alias := range cmd.Aliases() {
		r.aliases[alias] = cmd.Name()
	}
}

// Get retrieves a command by name or alias
func (r *CommandRegistry) Get(name string) (Command, bool) {
	// Try direct lookup
	if cmd, ok := r.commands[name]; ok {
		return cmd, true
	}
	// Try alias lookup
	if cmdName, ok := r.aliases[name]; ok {
		return r.commands[cmdName], true
	}
	return nil, false
}

// List returns all registered commands
func (r *CommandRegistry) List() []Command {
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// Complete provides autocomplete suggestions for a partial command
func (r *CommandRegistry) Complete(partial string) []string {
	var matches []string
	for name := range r.commands {
		if strings.HasPrefix(name, partial) {
			matches = append(matches, name)
		}
	}
	for alias := range r.aliases {
		if strings.HasPrefix(alias, partial) {
			matches = append(matches, alias)
		}
	}
	return matches
}

// Parse parses a command line into command name and arguments
func ParseCommandLine(line string) (string, []string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", nil
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", nil
	}

	return parts[0], parts[1:]
}

// CommandHistory manages command history
type CommandHistory struct {
	history []string
	index   int
	maxSize int
}

// NewCommandHistory creates a new command history
func NewCommandHistory(maxSize int) *CommandHistory {
	return &CommandHistory{
		history: make([]string, 0, maxSize),
		index:   -1,
		maxSize: maxSize,
	}
}

// Add adds a command to history
func (h *CommandHistory) Add(cmd string) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}

	// Don't add duplicates of the last command
	if len(h.history) > 0 && h.history[len(h.history)-1] == cmd {
		h.index = -1
		return
	}

	h.history = append(h.history, cmd)
	if len(h.history) > h.maxSize {
		h.history = h.history[1:]
	}
	h.index = -1
}

// Previous returns the previous command in history
func (h *CommandHistory) Previous() string {
	if len(h.history) == 0 {
		return ""
	}

	if h.index == -1 {
		h.index = len(h.history) - 1
	} else if h.index > 0 {
		h.index--
	}

	return h.history[h.index]
}

// Next returns the next command in history
func (h *CommandHistory) Next() string {
	if len(h.history) == 0 || h.index == -1 {
		return ""
	}

	if h.index < len(h.history)-1 {
		h.index++
		return h.history[h.index]
	}

	h.index = -1
	return ""
}

// List returns all commands in history
func (h *CommandHistory) List() []string {
	return append([]string{}, h.history...)
}

// BaseCommand provides default implementations for Command interface
type BaseCommand struct {
	name        string
	description string
	usage       string
	aliases     []string
}

func (c *BaseCommand) Name() string        { return c.name }
func (c *BaseCommand) Description() string { return c.description }
func (c *BaseCommand) Usage() string       { return c.usage }
func (c *BaseCommand) Aliases() []string   { return c.aliases }
func (c *BaseCommand) Complete(args []string) []string { return nil }
func (c *BaseCommand) Execute(ctx *CommandContext, args []string) error {
	return fmt.Errorf("command not implemented")
}
