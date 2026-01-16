package ui

import (
	"strings"

	"slip/internal/engine"
)

// Dialog is a modal dialog box
type Dialog struct {
	Title    string
	Message  string
	Buttons  []string
	Selected int
	Width    int
	X, Y     int
	Style    DialogStyle
	OnSelect func(buttonIndex int)
}

// DialogStyle defines the visual style of a dialog
type DialogStyle struct {
	Border  engine.Style
	Title   engine.Style
	Message engine.Style
	Button  engine.Style
	ButtonSelected engine.Style
	BoxStyle engine.BoxStyle
}

// DefaultDialogStyle returns the default dialog style
func DefaultDialogStyle() DialogStyle {
	return DialogStyle{
		Border:  engine.Style{FG: engine.ColorYellow},
		Title:   engine.Style{FG: engine.ColorBrightWhite, Bold: true},
		Message: engine.Style{FG: engine.ColorWhite},
		Button:  engine.Style{FG: engine.ColorWhite},
		ButtonSelected: engine.Style{FG: engine.ColorBlack, BG: engine.ColorYellow, Bold: true},
		BoxStyle: engine.BoxStyleDouble,
	}
}

// NewDialog creates a new dialog
func NewDialog(title, message string, buttons []string) *Dialog {
	// Calculate width based on content
	width := len(title) + 4

	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if len(line)+4 > width {
			width = len(line) + 4
		}
	}

	// Account for buttons
	buttonsWidth := 0
	for _, b := range buttons {
		buttonsWidth += len(b) + 4 // [ button ] + spacing
	}
	if buttonsWidth+4 > width {
		width = buttonsWidth + 4
	}

	if width < 20 {
		width = 20
	}

	return &Dialog{
		Title:    title,
		Message:  message,
		Buttons:  buttons,
		Selected: 0,
		Width:    width,
		Style:    DefaultDialogStyle(),
	}
}

// Height returns the dialog height
func (d *Dialog) Height() int {
	lines := strings.Split(d.Message, "\n")
	return len(lines) + 6 // border top + title + padding + message + padding + buttons + border bottom
}

// MoveLeft moves button selection left
func (d *Dialog) MoveLeft() {
	d.Selected--
	if d.Selected < 0 {
		d.Selected = len(d.Buttons) - 1
	}
}

// MoveRight moves button selection right
func (d *Dialog) MoveRight() {
	d.Selected++
	if d.Selected >= len(d.Buttons) {
		d.Selected = 0
	}
}

// Select activates the current button
func (d *Dialog) Select() {
	if d.OnSelect != nil {
		d.OnSelect(d.Selected)
	}
}

// HandleInput processes input for the dialog
func (d *Dialog) HandleInput(input engine.Input) bool {
	switch input.Key {
	case engine.KeyLeft:
		d.MoveLeft()
		return true
	case engine.KeyRight:
		d.MoveRight()
		return true
	case engine.KeyTab:
		d.MoveRight()
		return true
	case engine.KeyEnter, engine.KeySpace:
		d.Select()
		return true
	}
	return false
}

// Render draws the dialog to the screen
func (d *Dialog) Render(screen *engine.Screen) {
	x, y := d.X, d.Y
	h := d.Height()

	// Draw border
	screen.DrawBox(x, y, d.Width, h, d.Style.BoxStyle, d.Style.Border)

	// Draw title
	titleX := x + (d.Width-len(d.Title))/2
	screen.DrawText(titleX, y, d.Title, d.Style.Title)

	// Draw message
	lines := strings.Split(d.Message, "\n")
	for i, line := range lines {
		msgX := x + (d.Width-len(line))/2
		screen.DrawText(msgX, y+2+i, line, d.Style.Message)
	}

	// Draw buttons
	buttonY := y + h - 2
	totalButtonsWidth := 0
	for _, b := range d.Buttons {
		totalButtonsWidth += len(b) + 4
	}
	buttonX := x + (d.Width-totalButtonsWidth)/2

	for i, button := range d.Buttons {
		style := d.Style.Button
		if i == d.Selected {
			style = d.Style.ButtonSelected
		}

		label := "[ " + button + " ]"
		screen.DrawText(buttonX, buttonY, label, style)
		buttonX += len(label) + 1
	}
}

// CenterOn centers the dialog on the given screen dimensions
func (d *Dialog) CenterOn(screenWidth, screenHeight int) {
	d.X = (screenWidth - d.Width) / 2
	d.Y = (screenHeight - d.Height()) / 2
}

// ConfirmDialog creates a Yes/No confirmation dialog
func ConfirmDialog(title, message string, onConfirm func(confirmed bool)) *Dialog {
	d := NewDialog(title, message, []string{"Yes", "No"})
	d.OnSelect = func(idx int) {
		onConfirm(idx == 0)
	}
	return d
}

// AlertDialog creates an OK-only alert dialog
func AlertDialog(title, message string, onClose func()) *Dialog {
	d := NewDialog(title, message, []string{"OK"})
	d.OnSelect = func(idx int) {
		if onClose != nil {
			onClose()
		}
	}
	return d
}
