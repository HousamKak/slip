package ui

import (
	"slip/internal/engine"
)

// MenuItem represents a menu item
type MenuItem struct {
	ID          string
	Label       string
	Description string
	Disabled    bool
	Action      func()
}

// Menu is a selectable menu widget
type Menu struct {
	Title       string
	Items       []MenuItem
	Selected    int
	X, Y        int
	Width       int
	Style       MenuStyle
	ShowBorder  bool
	ShowNumbers bool
}

// MenuStyle defines the visual style of a menu
type MenuStyle struct {
	Border         engine.Style
	Title          engine.Style
	Item           engine.Style
	ItemSelected   engine.Style
	ItemDisabled   engine.Style
	Description    engine.Style
	Highlight      rune
	BoxStyle       engine.BoxStyle
}

// DefaultMenuStyle returns the default menu style
func DefaultMenuStyle() MenuStyle {
	return MenuStyle{
		Border:       engine.Style{FG: engine.ColorCyan},
		Title:        engine.Style{FG: engine.ColorBrightWhite, Bold: true},
		Item:         engine.Style{FG: engine.ColorWhite},
		ItemSelected: engine.Style{FG: engine.ColorBlack, BG: engine.ColorCyan, Bold: true, Reverse: false},
		ItemDisabled: engine.Style{FG: engine.ColorBrightBlack},
		Description:  engine.Style{FG: engine.ColorBrightBlack},
		Highlight:    '▸',
		BoxStyle:     engine.BoxStyleDouble,
	}
}

// NewMenu creates a new menu
func NewMenu(title string, items []MenuItem) *Menu {
	// Calculate width: label + "▸ " (2) + "N. " (3) + padding (4) + border (2)
	width := len(title) + 6
	for _, item := range items {
		itemWidth := len(item.Label) + 11 // 2 (indicator) + 3 (number) + 4 (padding) + 2 (border)
		if itemWidth > width {
			width = itemWidth
		}
	}
	// Ensure minimum width for descriptions
	if width < 25 {
		width = 25
	}
	return &Menu{
		Title:      title,
		Items:      items,
		Selected:   0,
		Width:      width,
		Style:      DefaultMenuStyle(),
		ShowBorder: true,
	}
}

// Height returns the height of the menu
func (m *Menu) Height() int {
	h := len(m.Items)
	if m.ShowBorder {
		h += 4 // Title + top border + bottom border + padding
	}
	return h
}

// MoveUp moves selection up
func (m *Menu) MoveUp() {
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(m.Items) - 1
	}
	// Skip disabled items
	for m.Items[m.Selected].Disabled && m.Selected > 0 {
		m.Selected--
	}
}

// MoveDown moves selection down
func (m *Menu) MoveDown() {
	m.Selected++
	if m.Selected >= len(m.Items) {
		m.Selected = 0
	}
	// Skip disabled items
	for m.Items[m.Selected].Disabled && m.Selected < len(m.Items)-1 {
		m.Selected++
	}
}

// Select activates the current selection
func (m *Menu) Select() {
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		item := m.Items[m.Selected]
		if !item.Disabled && item.Action != nil {
			item.Action()
		}
	}
}

// GetSelected returns the currently selected item
func (m *Menu) GetSelected() *MenuItem {
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		return &m.Items[m.Selected]
	}
	return nil
}

// HandleInput processes input for the menu
func (m *Menu) HandleInput(input engine.Input) bool {
	switch input.Key {
	case engine.KeyUp:
		m.MoveUp()
		return true
	case engine.KeyDown:
		m.MoveDown()
		return true
	case engine.KeyEnter, engine.KeySpace:
		m.Select()
		return true
	case engine.KeyRune:
		// Number selection
		if input.Rune >= '1' && input.Rune <= '9' {
			idx := int(input.Rune - '1')
			if idx < len(m.Items) && !m.Items[idx].Disabled {
				m.Selected = idx
				m.Select()
				return true
			}
		}
		// Letter selection (first letter of item)
		for i, item := range m.Items {
			if !item.Disabled && len(item.Label) > 0 {
				if rune(item.Label[0]) == input.Rune ||
					rune(item.Label[0])+32 == input.Rune ||
					rune(item.Label[0])-32 == input.Rune {
					m.Selected = i
					m.Select()
					return true
				}
			}
		}
	}
	return false
}

// Render draws the menu to the screen
func (m *Menu) Render(screen *engine.Screen) {
	x, y := m.X, m.Y

	if m.ShowBorder {
		// Draw border
		screen.DrawBox(x, y, m.Width, m.Height(), m.Style.BoxStyle, m.Style.Border)

		// Draw title
		titleX := x + (m.Width-len(m.Title))/2
		screen.DrawText(titleX, y, m.Title, m.Style.Title)
		y += 2
	}

	// Draw items
	for i, item := range m.Items {
		style := m.Style.Item
		prefix := "  "

		if item.Disabled {
			style = m.Style.ItemDisabled
		} else if i == m.Selected {
			style = m.Style.ItemSelected
			prefix = string(m.Style.Highlight) + " "
		}

		// Draw item background for selected
		if i == m.Selected && !item.Disabled {
			for j := 0; j < m.Width-2; j++ {
				screen.Set(x+1+j, y+i, ' ', style)
			}
		}

		label := prefix + item.Label
		if m.ShowNumbers && i < 9 {
			label = prefix + string(rune('1'+i)) + ". " + item.Label
		}

		// Truncate if too long
		maxLen := m.Width - 3
		if len(label) > maxLen {
			label = label[:maxLen-2] + ".."
		}

		screen.DrawText(x+1, y+i, label, style)
	}

	// Draw description of selected item
	if m.Selected >= 0 && m.Selected < len(m.Items) {
		desc := m.Items[m.Selected].Description
		if desc != "" && m.ShowBorder {
			descY := m.Y + m.Height() - 1
			descX := x + 2
			maxLen := m.Width - 4
			if len(desc) > maxLen {
				desc = desc[:maxLen-2] + ".."
			}
			screen.DrawText(descX, descY, desc, m.Style.Description)
		}
	}
}

// CenterOn centers the menu on the given screen dimensions
func (m *Menu) CenterOn(screenWidth, screenHeight int) {
	m.X = (screenWidth - m.Width) / 2
	m.Y = (screenHeight - m.Height()) / 2
}
