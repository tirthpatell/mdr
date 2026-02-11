package viewer

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := NewModel("# Hello\n\nWorld")
	if m.content == "" {
		t.Fatal("expected content to be set")
	}
}

func TestModel_ScrollDown(t *testing.T) {
	m := NewModel("# Hello\n\nLine 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10")
	m.height = 5

	msg := tea.KeyMsg{Type: tea.KeyDown}
	updated, _ := m.Update(msg)
	model := updated.(Model)
	if model.offset <= 0 {
		t.Fatal("expected offset to increase after scroll down")
	}
}

func TestModel_ScrollUp(t *testing.T) {
	m := NewModel("# Hello\n\nContent")
	m.height = 5
	m.offset = 3

	msg := tea.KeyMsg{Type: tea.KeyUp}
	updated, _ := m.Update(msg)
	model := updated.(Model)
	if model.offset >= 3 {
		t.Fatal("expected offset to decrease after scroll up")
	}
}

func TestModel_Quit(t *testing.T) {
	m := NewModel("# Hello")

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(msg)
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}
