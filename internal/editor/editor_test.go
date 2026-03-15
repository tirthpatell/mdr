package editor

import (
	"os"
	"path/filepath"
	"strings"
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

func TestEditorModel_SaveResetsModified(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.md")
	os.WriteFile(path, []byte("hello"), 0644)

	m, err := NewModelFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m.width = 80
	m.height = 24

	// Type a character to set modified
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model := updated.(Model)
	if !model.buffer.Modified() {
		t.Fatal("buffer should be modified after typing")
	}

	// Save
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	model = updated.(Model)
	if model.buffer.Modified() {
		t.Fatal("buffer should not be modified after save")
	}
	if model.saveMsg != "Saved!" {
		t.Fatalf("expected save message 'Saved!', got %q", model.saveMsg)
	}
}

func TestEditorModel_SaveError_DisplayedInView(t *testing.T) {
	m := NewModel("hello", "/nonexistent/path/file.md")
	m.width = 80
	m.height = 24
	m.editWidth = 40

	// Trigger save
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	model := updated.(Model)
	if model.err == nil {
		t.Fatal("expected save error for invalid path")
	}
	view := model.View()
	if !strings.Contains(view, "Error:") {
		t.Fatal("expected error message in view output")
	}
}

func TestEditorModel_QuitConfirmation_Modified(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Type to set modified
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model := updated.(Model)

	// Press Ctrl+C — should NOT quit, should show confirmation
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	model = updated.(Model)
	if cmd != nil {
		t.Fatal("should not quit immediately with unsaved changes")
	}
	if !model.confirmQuit {
		t.Fatal("expected confirmQuit to be true")
	}

	// View should show the confirmation prompt
	view := model.View()
	if !strings.Contains(view, "Unsaved changes") {
		t.Fatal("expected confirmation prompt in view")
	}

	// Press 'n' to cancel
	updated, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	model = updated.(Model)
	if cmd != nil {
		t.Fatal("pressing 'n' should not quit")
	}
	if model.confirmQuit {
		t.Fatal("confirmQuit should be false after 'n'")
	}
}

func TestEditorModel_QuitConfirmation_Accept(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Type to set modified
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model := updated.(Model)

	// Press Ctrl+C
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	model = updated.(Model)

	// Press 'y' to confirm quit
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if cmd == nil {
		t.Fatal("pressing 'y' should quit")
	}
}

func TestEditorModel_QuitUnmodified(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Press Ctrl+C on unmodified buffer — should quit immediately
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("should quit immediately when no unsaved changes")
	}
}

func TestEditorModel_Init(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	if cmd := m.Init(); cmd != nil {
		t.Fatal("Init should return nil")
	}
}

func TestEditorModel_WindowSizeMsg(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model := updated.(Model)
	if model.width != 120 || model.height != 40 {
		t.Fatalf("expected 120x40, got %dx%d", model.width, model.height)
	}
	if model.editWidth != 60 {
		t.Fatalf("expected editWidth 60, got %d", model.editWidth)
	}
}

func TestEditorModel_LeftWrap(t *testing.T) {
	m := NewModel("hello\nworld", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorRow = 1
	m.cursorCol = 0

	// Press Left at beginning of line 1 — should wrap to end of line 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model := updated.(Model)
	if model.cursorRow != 0 {
		t.Fatalf("expected row 0, got %d", model.cursorRow)
	}
	if model.cursorCol != 5 {
		t.Fatalf("expected col 5, got %d", model.cursorCol)
	}
}

func TestEditorModel_RightWrap(t *testing.T) {
	m := NewModel("hello\nworld", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorCol = 5 // end of "hello"

	// Press Right at end of line 0 — should wrap to start of line 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	model := updated.(Model)
	if model.cursorRow != 1 {
		t.Fatalf("expected row 1, got %d", model.cursorRow)
	}
	if model.cursorCol != 0 {
		t.Fatalf("expected col 0, got %d", model.cursorCol)
	}
}

func TestEditorModel_HomeEnd(t *testing.T) {
	m := NewModel("hello world", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorCol = 5

	// End key
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnd})
	model := updated.(Model)
	if model.cursorCol != 11 {
		t.Fatalf("expected col 11 after End, got %d", model.cursorCol)
	}

	// Home key
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyHome})
	model = updated.(Model)
	if model.cursorCol != 0 {
		t.Fatalf("expected col 0 after Home, got %d", model.cursorCol)
	}
}

func TestEditorModel_Tab(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := updated.(Model)
	if got := model.buffer.GetLine(0); got != "\thello" {
		t.Fatalf("expected tab at start, got %q", got)
	}
	if model.cursorCol != 1 {
		t.Fatalf("expected col 1, got %d", model.cursorCol)
	}
}

func TestEditorModel_Delete(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Delete at start of line
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDelete})
	model := updated.(Model)
	if got := model.buffer.GetLine(0); got != "ello" {
		t.Fatalf("expected 'ello', got %q", got)
	}
}

func TestEditorModel_DeleteJoinsLines(t *testing.T) {
	m := NewModel("hello\nworld", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorCol = 5 // end of "hello"

	// Delete at end of line should join with next
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDelete})
	model := updated.(Model)
	if model.buffer.LineCount() != 1 {
		t.Fatalf("expected 1 line, got %d", model.buffer.LineCount())
	}
	if got := model.buffer.GetLine(0); got != "helloworld" {
		t.Fatalf("expected 'helloworld', got %q", got)
	}
}

func TestEditorModel_BackspaceJoinsLines(t *testing.T) {
	m := NewModel("hello\nworld", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorRow = 1
	m.cursorCol = 0

	// Backspace at start of line 1 joins with line 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model := updated.(Model)
	if model.buffer.LineCount() != 1 {
		t.Fatalf("expected 1 line, got %d", model.buffer.LineCount())
	}
	if model.cursorRow != 0 {
		t.Fatalf("expected cursor on row 0, got %d", model.cursorRow)
	}
	if model.cursorCol != 5 {
		t.Fatalf("expected cursor at col 5, got %d", model.cursorCol)
	}
}

func TestEditorModel_CtrlH_ToggleHelp(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.editWidth = 40

	// Toggle help on
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlH})
	model := updated.(Model)
	if !model.showHelp {
		t.Fatal("expected showHelp true")
	}
	view := model.View()
	if !strings.Contains(view, "Editor Help") {
		t.Fatal("expected help text in view")
	}

	// Toggle help off
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlH})
	model = updated.(Model)
	if model.showHelp {
		t.Fatal("expected showHelp false")
	}
}

func TestEditorModel_ConfirmQuit_Escape(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Type to set modified
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model := updated.(Model)

	// Ctrl+C to get confirmation
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	model = updated.(Model)

	// Escape to cancel
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updated.(Model)
	if cmd != nil {
		t.Fatal("escape should not quit")
	}
	if model.confirmQuit {
		t.Fatal("confirmQuit should be false after escape")
	}
}

func TestEditorModel_ConfirmQuit_UnknownKey(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Type to set modified
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model := updated.(Model)

	// Ctrl+C to get confirmation
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	model = updated.(Model)

	// Press an unrelated key — should stay in confirmation mode
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	model = updated.(Model)
	if cmd != nil {
		t.Fatal("unknown key should not quit")
	}
	if !model.confirmQuit {
		t.Fatal("confirmQuit should still be true")
	}
}

func TestEditorModel_DownClampCol(t *testing.T) {
	m := NewModel("hello world\nhi", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorCol = 10 // somewhere in "hello world"

	// Move down to shorter line — col should clamp
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	model := updated.(Model)
	if model.cursorCol != 2 {
		t.Fatalf("expected col clamped to 2, got %d", model.cursorCol)
	}
}

func TestEditorModel_UpClampCol(t *testing.T) {
	m := NewModel("hi\nhello world", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorRow = 1
	m.cursorCol = 10

	// Move up to shorter line — col should clamp
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	model := updated.(Model)
	if model.cursorCol != 2 {
		t.Fatalf("expected col clamped to 2, got %d", model.cursorCol)
	}
}

func TestEditorModel_View_Loading(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	// width=0, height=0
	view := m.View()
	if view != "Loading..." {
		t.Fatalf("expected 'Loading...', got %q", view)
	}
}

func TestEditorModel_View_ModifiedIndicator(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.editWidth = 40

	// Type to set modified
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model := updated.(Model)
	view := model.View()
	if !strings.Contains(view, "[+]") {
		t.Fatal("expected [+] modified indicator in view")
	}
}

func TestEditorModel_ClearTransientMessages(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.md")
	os.WriteFile(path, []byte("hello"), 0644)

	m, _ := NewModelFromFile(path)
	m.width = 80
	m.height = 24

	// Save to get save message
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	model := updated.(Model)
	if model.saveMsg != "Saved!" {
		t.Fatal("expected save message")
	}

	// Any non-save key should clear transient messages
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model = updated.(Model)
	if model.saveMsg != "" {
		t.Fatal("expected save message to be cleared")
	}
}

func TestEditorModel_NewModelFromFile_NotFound(t *testing.T) {
	_, err := NewModelFromFile("/nonexistent/path/file.md")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestEditorModel_ScrollToCursor(t *testing.T) {
	// Create content with many lines
	lines := make([]string, 50)
	for i := range lines {
		lines[i] = "line"
	}
	content := strings.Join(lines, "\n")
	m := NewModel(content, "/tmp/test.md")
	m.width = 80
	m.height = 10

	// Move cursor to bottom
	for i := 0; i < 20; i++ {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = updated.(Model)
	}
	if m.offsetRow == 0 {
		t.Fatal("expected offsetRow to increase with scrolling")
	}
}

func TestEditorModel_UpAtTop(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Up at row 0 — should be no-op
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	model := updated.(Model)
	if model.cursorRow != 0 {
		t.Fatalf("expected row 0, got %d", model.cursorRow)
	}
}

func TestEditorModel_DownAtBottom(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Down at last line — should be no-op
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	model := updated.(Model)
	if model.cursorRow != 0 {
		t.Fatalf("expected row 0, got %d", model.cursorRow)
	}
}

func TestEditorModel_LeftAtStart(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24

	// Left at row 0, col 0 — should be no-op
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	model := updated.(Model)
	if model.cursorRow != 0 || model.cursorCol != 0 {
		t.Fatal("expected no movement at start")
	}
}

func TestEditorModel_RightAtEnd(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorCol = 5

	// Right at end of only line — should be no-op
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	model := updated.(Model)
	if model.cursorRow != 0 || model.cursorCol != 5 {
		t.Fatal("expected no movement at end")
	}
}

func TestEditorModel_BackspaceInLine(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	m.width = 80
	m.height = 24
	m.cursorCol = 3

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model := updated.(Model)
	if got := model.buffer.GetLine(0); got != "helo" {
		t.Fatalf("expected 'helo', got %q", got)
	}
	if model.cursorCol != 2 {
		t.Fatalf("expected col 2, got %d", model.cursorCol)
	}
}

func TestEditorModel_TypeMultipleRunes(t *testing.T) {
	m := NewModel("", "/tmp/test.md")
	m.width = 80
	m.height = 24

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a', 'b', 'c'}}
	updated, _ := m.Update(msg)
	model := updated.(Model)
	if got := model.buffer.GetLine(0); got != "abc" {
		t.Fatalf("expected 'abc', got %q", got)
	}
	if model.cursorCol != 3 {
		t.Fatalf("expected col 3, got %d", model.cursorCol)
	}
}

func TestEditorModel_UnknownMsg(t *testing.T) {
	m := NewModel("hello", "/tmp/test.md")
	// Pass an unhandled message type
	updated, cmd := m.Update("unhandled string message")
	model := updated.(Model)
	if cmd != nil {
		t.Fatal("expected nil cmd for unknown msg")
	}
	if model.cursorRow != 0 || model.cursorCol != 0 {
		t.Fatal("expected no state change for unknown msg")
	}
}

func TestTruncate(t *testing.T) {
	if got := truncate("hello world", 5); got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
	if got := truncate("hi", 10); got != "hi" {
		t.Fatalf("expected 'hi', got %q", got)
	}
	if got := truncate("hello", 0); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
	if got := truncate("hello", -1); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
