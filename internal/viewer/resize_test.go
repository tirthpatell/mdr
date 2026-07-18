package viewer

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModel_ResizeClampsOffset(t *testing.T) {
	lines := make([]string, 20)
	for i := range lines {
		lines[i] = "line"
	}
	m := NewModel(strings.Join(lines, "\n"))
	m.height = 5

	// Scroll to the bottom, then grow the window: the offset must be
	// clamped so the view doesn't stay scrolled past the content.
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	model := updated.(Model)

	updated, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 30})
	model = updated.(Model)

	if model.offset != model.maxOffset() {
		t.Fatalf("expected offset clamped to maxOffset %d after resize, got %d", model.maxOffset(), model.offset)
	}
}
