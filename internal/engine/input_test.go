package engine

import (
	"testing"
)

func TestInputIsQuit(t *testing.T) {
	tests := []struct {
		name     string
		input    Input
		expected bool
	}{
		{"Ctrl+C", Input{Key: KeyCtrlC}, true},
		{"Ctrl+D", Input{Key: KeyCtrlD}, true},
		{"Ctrl+Q", Input{Key: KeyCtrlQ}, true},
		{"Quit type", Input{Type: InputQuit}, true},
		{"Escape", Input{Key: KeyEscape}, false},
		{"Enter", Input{Key: KeyEnter}, false},
		{"Regular rune", Input{Key: KeyRune, Rune: 'a'}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsQuit()
			if result != tt.expected {
				t.Errorf("IsQuit() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestInputIsEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    Input
		expected bool
	}{
		{"Escape key", Input{Key: KeyEscape}, true},
		{"Enter key", Input{Key: KeyEnter}, false},
		{"Ctrl+C", Input{Key: KeyCtrlC}, false},
		{"Regular rune", Input{Key: KeyRune, Rune: 'a'}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsEscape()
			if result != tt.expected {
				t.Errorf("IsEscape() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestInputIsEnter(t *testing.T) {
	tests := []struct {
		name     string
		input    Input
		expected bool
	}{
		{"Enter key", Input{Key: KeyEnter}, true},
		{"Escape key", Input{Key: KeyEscape}, false},
		{"Space", Input{Key: KeySpace}, false},
		{"Regular rune", Input{Key: KeyRune, Rune: 'a'}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsEnter()
			if result != tt.expected {
				t.Errorf("IsEnter() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestInputIsArrow(t *testing.T) {
	tests := []struct {
		name     string
		input    Input
		expected bool
	}{
		{"Up arrow", Input{Key: KeyUp}, true},
		{"Down arrow", Input{Key: KeyDown}, true},
		{"Left arrow", Input{Key: KeyLeft}, true},
		{"Right arrow", Input{Key: KeyRight}, true},
		{"Enter", Input{Key: KeyEnter}, false},
		{"Space", Input{Key: KeySpace}, false},
		{"Regular rune", Input{Key: KeyRune, Rune: 'w'}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsArrow()
			if result != tt.expected {
				t.Errorf("IsArrow() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestInputDirection(t *testing.T) {
	tests := []struct {
		name     string
		input    Input
		expected Direction
	}{
		{"Up arrow", Input{Key: KeyUp}, DirUp},
		{"Down arrow", Input{Key: KeyDown}, DirDown},
		{"Left arrow", Input{Key: KeyLeft}, DirLeft},
		{"Right arrow", Input{Key: KeyRight}, DirRight},
		{"No arrow", Input{Key: KeyEnter}, DirNone},
		{"Regular rune", Input{Key: KeyRune, Rune: 'a'}, DirNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Direction()
			if result != tt.expected {
				t.Errorf("Direction() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestParseInput(t *testing.T) {
	reader := NewInputReader()

	tests := []struct {
		name     string
		buf      []byte
		expected []Input
	}{
		{
			name: "Up arrow",
			buf:  []byte{0x1b, '[', 'A'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyUp},
			},
		},
		{
			name: "Down arrow",
			buf:  []byte{0x1b, '[', 'B'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyDown},
			},
		},
		{
			name: "Left arrow",
			buf:  []byte{0x1b, '[', 'D'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyLeft},
			},
		},
		{
			name: "Right arrow",
			buf:  []byte{0x1b, '[', 'C'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyRight},
			},
		},
		{
			name: "Enter",
			buf:  []byte{0x0d},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyEnter},
			},
		},
		{
			name: "Escape",
			buf:  []byte{0x1b},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyEscape},
			},
		},
		{
			name: "Space",
			buf:  []byte{' '},
			expected: []Input{
				{Type: InputKeyPress, Key: KeySpace},
			},
		},
		{
			name: "Ctrl+C",
			buf:  []byte{0x03},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyCtrlC},
			},
		},
		{
			name: "Ctrl+D",
			buf:  []byte{0x04},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyCtrlD},
			},
		},
		{
			name: "Tab",
			buf:  []byte{0x09},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyTab},
			},
		},
		{
			name: "Backspace",
			buf:  []byte{0x7f},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyBackspace},
			},
		},
		{
			name: "Regular character 'a'",
			buf:  []byte{'a'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyRune, Rune: 'a'},
			},
		},
		{
			name: "Regular character 'Z'",
			buf:  []byte{'Z'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyRune, Rune: 'Z'},
			},
		},
		{
			name: "Multiple characters",
			buf:  []byte{'h', 'e', 'l', 'l', 'o'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyRune, Rune: 'h'},
				{Type: InputKeyPress, Key: KeyRune, Rune: 'e'},
				{Type: InputKeyPress, Key: KeyRune, Rune: 'l'},
				{Type: InputKeyPress, Key: KeyRune, Rune: 'l'},
				{Type: InputKeyPress, Key: KeyRune, Rune: 'o'},
			},
		},
		{
			name: "Home key",
			buf:  []byte{0x1b, '[', 'H'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyHome},
			},
		},
		{
			name: "End key",
			buf:  []byte{0x1b, '[', 'F'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyEnd},
			},
		},
		{
			name: "Delete key",
			buf:  []byte{0x1b, '[', '3', '~'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyDelete},
			},
		},
		{
			name: "Page Up",
			buf:  []byte{0x1b, '[', '5', '~'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyPageUp},
			},
		},
		{
			name: "Page Down",
			buf:  []byte{0x1b, '[', '6', '~'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyPageDown},
			},
		},
		{
			name: "F1 key",
			buf:  []byte{0x1b, 'O', 'P'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyF1},
			},
		},
		{
			name: "F2 key",
			buf:  []byte{0x1b, 'O', 'Q'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyF2},
			},
		},
		{
			name: "F3 key",
			buf:  []byte{0x1b, 'O', 'R'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyF3},
			},
		},
		{
			name: "F4 key",
			buf:  []byte{0x1b, 'O', 'S'},
			expected: []Input{
				{Type: InputKeyPress, Key: KeyF4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reader.parseInput(tt.buf)

			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d inputs, got %d", len(tt.expected), len(result))
			}

			for i := range result {
				if result[i].Type != tt.expected[i].Type {
					t.Errorf("Input %d: Type = %v, expected %v", i, result[i].Type, tt.expected[i].Type)
				}
				if result[i].Key != tt.expected[i].Key {
					t.Errorf("Input %d: Key = %v, expected %v", i, result[i].Key, tt.expected[i].Key)
				}
				if result[i].Rune != tt.expected[i].Rune {
					t.Errorf("Input %d: Rune = %c, expected %c", i, result[i].Rune, tt.expected[i].Rune)
				}
			}
		})
	}
}

func TestNewInputReader(t *testing.T) {
	reader := NewInputReader()
	if reader == nil {
		t.Fatal("NewInputReader() returned nil")
	}
	if reader.ch == nil {
		t.Error("Input channel should be initialized")
	}
	if reader.done == nil {
		t.Error("Done channel should be initialized")
	}
}

func TestInputReaderChannel(t *testing.T) {
	reader := NewInputReader()
	ch := reader.Channel()
	if ch == nil {
		t.Error("Channel() returned nil")
	}
}
