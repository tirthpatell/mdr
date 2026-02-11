package editor

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewEditorModel(t *testing.T) {
	m := NewModel("# Hello\n\nWorld", "/tmp/test.md")
	if m.filePath != "/tmp/test.md" {
		t.Fatalf("expected path '/tmp/test.md', got %q", m.filePath)
	}
	if m.buffer.LineCount() != 3 {
		t.Fatalf("expected 3 lines, got %d", m.buffer.LineCount())
	}
}

func TestEditorModel_CursorMovement(t *testing.T) {
	m := NewModel("hello\nworld", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Move right
	msg := tea.KeyMsg{Type: tea.KeyRight}
	updated, _ := m.Update(msg)
	model := updated.(Model)
	if model.cursorCol != 1 {
		t.Fatalf("expected col 1, got %d", model.cursorCol)
	}

	// Move down
	msg = tea.KeyMsg{Type: tea.KeyDown}
	updated, _ = model.Update(msg)
	model = updated.(Model)
	if model.cursorRow != 1 {
		t.Fatalf("expected row 1, got %d", model.cursorRow)
	}
}

func TestEditorModel_TypeCharacter(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}}
	updated, _ := m.Update(msg)
	model := updated.(Model)
	if got := model.buffer.GetLine(0); got != "!hello" {
		t.Fatalf("expected '!hello', got %q", got)
	}
}

func TestEditorModel_EnterKey(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorCol = 3

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ := m.Update(msg)
	model := updated.(Model)
	if model.buffer.LineCount() != 2 {
		t.Fatalf("expected 2 lines, got %d", model.buffer.LineCount())
	}
	if model.cursorRow != 1 {
		t.Fatalf("expected cursor on row 1, got %d", model.cursorRow)
	}
	if model.cursorCol != 0 {
		t.Fatalf("expected cursor at col 0, got %d", model.cursorCol)
	}
}
