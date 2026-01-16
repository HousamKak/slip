package ui

import "slip/internal/engine"

// Theme represents a color theme
type Theme struct {
	Name        string
	Primary     engine.Color
	Secondary   engine.Color
	Accent      engine.Color
	Background  engine.Color
	Text        engine.Color
	TextDim     engine.Color
	Success     engine.Color
	Warning     engine.Color
	Error       engine.Color
	Border      engine.Color
}

// Themes holds all available themes
var Themes = map[string]Theme{
	"default": {
		Name:       "Default",
		Primary:    engine.ColorBrightCyan,
		Secondary:  engine.ColorCyan,
		Accent:     engine.ColorBrightYellow,
		Background: engine.ColorBlack,
		Text:       engine.ColorWhite,
		TextDim:    engine.ColorBrightBlack,
		Success:    engine.ColorGreen,
		Warning:    engine.ColorYellow,
		Error:      engine.ColorRed,
		Border:     engine.ColorCyan,
	},
	"gruvbox": {
		Name:       "Gruvbox",
		Primary:    engine.ColorBrightYellow,
		Secondary:  engine.ColorYellow,
		Accent:     engine.ColorBrightRed,
		Background: engine.ColorBlack,
		Text:       engine.ColorWhite,
		TextDim:    engine.ColorBrightBlack,
		Success:    engine.ColorGreen,
		Warning:    engine.ColorYellow,
		Error:      engine.ColorRed,
		Border:     engine.ColorYellow,
	},
	"nord": {
		Name:       "Nord",
		Primary:    engine.ColorBrightCyan,
		Secondary:  engine.ColorCyan,
		Accent:     engine.ColorBrightBlue,
		Background: engine.ColorBlack,
		Text:       engine.ColorWhite,
		TextDim:    engine.ColorBrightBlack,
		Success:    engine.ColorGreen,
		Warning:    engine.ColorYellow,
		Error:      engine.ColorRed,
		Border:     engine.ColorBrightBlue,
	},
	"dracula": {
		Name:       "Dracula",
		Primary:    engine.ColorBrightMagenta,
		Secondary:  engine.ColorMagenta,
		Accent:     engine.ColorBrightCyan,
		Background: engine.ColorBlack,
		Text:       engine.ColorWhite,
		TextDim:    engine.ColorBrightBlack,
		Success:    engine.ColorGreen,
		Warning:    engine.ColorYellow,
		Error:      engine.ColorRed,
		Border:     engine.ColorMagenta,
	},
	"monokai": {
		Name:       "Monokai",
		Primary:    engine.ColorBrightGreen,
		Secondary:  engine.ColorGreen,
		Accent:     engine.ColorBrightYellow,
		Background: engine.ColorBlack,
		Text:       engine.ColorWhite,
		TextDim:    engine.ColorBrightBlack,
		Success:    engine.ColorGreen,
		Warning:    engine.ColorYellow,
		Error:      engine.ColorRed,
		Border:     engine.ColorGreen,
	},
}

// GetTheme returns a theme by name, or default if not found
func GetTheme(name string) Theme {
	if theme, ok := Themes[name]; ok {
		return theme
	}
	return Themes["default"]
}

// ListThemes returns names of all available themes
func ListThemes() []string {
	names := make([]string, 0, len(Themes))
	for name := range Themes {
		names = append(names, name)
	}
	return names
}
