package editor

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestEditorModel_CursorVisibleOnLongLine(t *testing.T) {
	long := strings.Repeat("abcdefghij", 10) // 100 chars, wider than the edit pane
	m := NewModel(long, "/tmp/test.md")

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 60, Height: 10})
	model := updated.(Model)

	// Move the cursor to the end of the line
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnd})
	model = updated.(Model)

	view := model.View()
	if !strings.Contains(view, "█") {
		t.Fatal("expected cursor block to be visible when cursor is past the pane width")
	}
	// The characters immediately before the cursor should be visible too
	if !strings.Contains(view, "hij█") {
		t.Fatalf("expected end of line to be visible next to the cursor, view:\n%s", view)
	}
}
