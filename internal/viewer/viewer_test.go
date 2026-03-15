package viewer

import (
	"strings"
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

func TestModel_ScrollToBottom(t *testing.T) {
	// Create content with 20 lines
	lines := make([]string, 20)
	for i := range lines {
		lines[i] = "line"
	}
	content := strings.Join(lines, "\n")
	m := NewModel(content)
	m.height = 10

	// Press G to go to bottom
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	model := updated.(Model)

	// The maxOffset should account for the status bar
	if model.offset != model.maxOffset() {
		t.Fatalf("expected offset %d to equal maxOffset %d", model.offset, model.maxOffset())
	}

	// Verify the view includes the last line of content
	view := model.View()
	viewLines := strings.Split(view, "\n")
	// Last visible content line (before the status line) should be "line"
	if len(viewLines) < 2 {
		t.Fatal("expected at least 2 view lines")
	}
	lastContentLine := viewLines[len(viewLines)-2]
	if lastContentLine != "line" {
		t.Fatalf("expected last content line to be 'line', got %q", lastContentLine)
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

func TestModel_SearchMode(t *testing.T) {
	m := NewModel("Line one\nLine two\nLine three")
	m.height = 10
	m.width = 80

	// Enter search mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model := updated.(Model)
	if !model.searching {
		t.Fatal("expected searching to be true")
	}

	// Type query
	for _, r := range "two" {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = updated.(Model)
	}
	if model.searchInput != "two" {
		t.Fatalf("expected search input 'two', got %q", model.searchInput)
	}

	// Press enter to execute search
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)
	if model.searching {
		t.Fatal("searching should be false after enter")
	}
	if model.searchQuery != "two" {
		t.Fatalf("expected search query 'two', got %q", model.searchQuery)
	}
	if len(model.matchLines) != 1 {
		t.Fatalf("expected 1 match, got %d", len(model.matchLines))
	}
}

func TestModel_SearchNavigate(t *testing.T) {
	// Create content with "match" appearing on lines 0, 10, 20
	lines := make([]string, 30)
	for i := range lines {
		if i%10 == 0 {
			lines[i] = "match line"
		} else {
			lines[i] = "other line"
		}
	}
	content := strings.Join(lines, "\n")
	m := NewModel(content)
	m.height = 10
	m.width = 80

	// Enter search mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model := updated.(Model)

	// Type query and execute
	for _, r := range "match" {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = updated.(Model)
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)

	if len(model.matchLines) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(model.matchLines))
	}
	if model.matchIndex != 0 {
		t.Fatalf("expected match index 0, got %d", model.matchIndex)
	}

	// Press 'n' to go to next match
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	model = updated.(Model)
	if model.matchIndex != 1 {
		t.Fatalf("expected match index 1 after 'n', got %d", model.matchIndex)
	}

	// Press 'N' to go to previous match
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
	model = updated.(Model)
	if model.matchIndex != 0 {
		t.Fatalf("expected match index 0 after 'N', got %d", model.matchIndex)
	}

	// Press 'N' again — should wrap to last match
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
	model = updated.(Model)
	if model.matchIndex != 2 {
		t.Fatalf("expected match index 2 after wrap, got %d", model.matchIndex)
	}
}

func TestModel_SearchEscapeClear(t *testing.T) {
	m := NewModel("Line one\nLine two\nLine three")
	m.height = 10
	m.width = 80

	// Set up a search result
	m.searchQuery = "two"
	m.matchLines = []int{1}
	m.matchIndex = 0

	// Press Escape to clear search
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model := updated.(Model)
	if model.searchQuery != "" {
		t.Fatal("expected search query to be cleared")
	}
	if model.matchLines != nil {
		t.Fatal("expected match lines to be nil")
	}
}

func TestModel_SearchCancelInput(t *testing.T) {
	m := NewModel("Line one\nLine two")
	m.height = 10
	m.width = 80

	// Enter search mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model := updated.(Model)

	// Type something
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updated.(Model)

	// Cancel with Escape
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updated.(Model)
	if model.searching {
		t.Fatal("expected searching to be false after escape")
	}
	if model.searchInput != "" {
		t.Fatal("expected search input to be cleared")
	}
}

func TestModel_SearchViewStatus(t *testing.T) {
	m := NewModel("Line one\nLine two\nLine three")
	m.height = 10
	m.width = 80

	// During search, View should show search prompt
	m.searching = true
	m.searchInput = "test"
	view := m.View()
	if !strings.Contains(view, "/test") {
		t.Fatal("expected search prompt in view during search")
	}

	// With active results, View should show match count
	m.searching = false
	m.searchQuery = "Line"
	m.matchLines = []int{0, 1, 2}
	m.matchIndex = 1
	view = m.View()
	if !strings.Contains(view, "[2/3]") {
		t.Fatal("expected match count in view")
	}
}

func TestModel_Init(t *testing.T) {
	m := NewModel("hello")
	if cmd := m.Init(); cmd != nil {
		t.Fatal("Init should return nil")
	}
}

func TestModel_WindowSizeMsg(t *testing.T) {
	m := NewModel("hello")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model := updated.(Model)
	if model.width != 120 || model.height != 40 {
		t.Fatalf("expected 120x40, got %dx%d", model.width, model.height)
	}
}

func TestModel_CtrlC(t *testing.T) {
	m := NewModel("hello")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected quit command for Ctrl+C")
	}
}

func TestModel_PageDown(t *testing.T) {
	lines := make([]string, 30)
	for i := range lines {
		lines[i] = "line"
	}
	m := NewModel(strings.Join(lines, "\n"))
	m.height = 10

	// Press 'd' for half-page down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	model := updated.(Model)
	if model.offset == 0 {
		t.Fatal("expected offset to increase after page down")
	}
}

func TestModel_PageUp(t *testing.T) {
	lines := make([]string, 30)
	for i := range lines {
		lines[i] = "line"
	}
	m := NewModel(strings.Join(lines, "\n"))
	m.height = 10
	m.offset = 15

	// Press 'u' for half-page up
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
	model := updated.(Model)
	if model.offset >= 15 {
		t.Fatal("expected offset to decrease after page up")
	}
}

func TestModel_PageUpAtTop(t *testing.T) {
	m := NewModel("line 1\nline 2\nline 3")
	m.height = 10
	m.offset = 0

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
	model := updated.(Model)
	if model.offset != 0 {
		t.Fatal("expected offset to stay at 0")
	}
}

func TestModel_PageDownClamp(t *testing.T) {
	lines := make([]string, 10)
	for i := range lines {
		lines[i] = "line"
	}
	m := NewModel(strings.Join(lines, "\n"))
	m.height = 10
	m.offset = 5

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	model := updated.(Model)
	if model.offset > model.maxOffset() {
		t.Fatalf("offset %d should not exceed maxOffset %d", model.offset, model.maxOffset())
	}
}

func TestModel_GoToTop(t *testing.T) {
	lines := make([]string, 20)
	for i := range lines {
		lines[i] = "line"
	}
	m := NewModel(strings.Join(lines, "\n"))
	m.height = 10
	m.offset = 10

	// Press 'g' to go to top
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	model := updated.(Model)
	if model.offset != 0 {
		t.Fatalf("expected offset 0 after g, got %d", model.offset)
	}
}

func TestModel_ScrollDownAtBottom(t *testing.T) {
	m := NewModel("one\ntwo")
	m.height = 10

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	model := updated.(Model)
	if model.offset != 0 {
		t.Fatal("expected offset 0 when content fits")
	}
}

func TestModel_ScrollUpAtTop(t *testing.T) {
	m := NewModel("one\ntwo")
	m.height = 10
	m.offset = 0

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	model := updated.(Model)
	if model.offset != 0 {
		t.Fatal("expected offset 0 when already at top")
	}
}

func TestModel_SearchBackspace(t *testing.T) {
	m := NewModel("hello")
	m.height = 10
	m.searching = true
	m.searchInput = "tes"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model := updated.(Model)
	if model.searchInput != "te" {
		t.Fatalf("expected 'te', got %q", model.searchInput)
	}
}

func TestModel_SearchBackspaceUnicode(t *testing.T) {
	m := NewModel("hello")
	m.height = 10
	m.searching = true
	m.searchInput = "hé"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model := updated.(Model)
	if model.searchInput != "h" {
		t.Fatalf("expected 'h' after unicode backspace, got %q", model.searchInput)
	}
}

func TestModel_SearchBackspaceEmpty(t *testing.T) {
	m := NewModel("hello")
	m.height = 10
	m.searching = true
	m.searchInput = ""

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model := updated.(Model)
	if model.searchInput != "" {
		t.Fatalf("expected empty, got %q", model.searchInput)
	}
}

func TestModel_SearchNoMatch(t *testing.T) {
	m := NewModel("hello\nworld")
	m.height = 10
	m.width = 80
	m.searching = true
	m.searchInput = "zzz"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model := updated.(Model)
	if len(model.matchLines) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(model.matchLines))
	}

	// View should show [no matches]
	view := model.View()
	if !strings.Contains(view, "[no matches]") {
		t.Fatal("expected [no matches] in view")
	}
}

func TestModel_SearchUnhandledKey(t *testing.T) {
	m := NewModel("hello")
	m.height = 10
	m.searching = true

	// Send an unhandled key type while searching
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	model := updated.(Model)
	if !model.searching {
		t.Fatal("expected still in search mode")
	}
}

func TestModel_SearchNextWrap(t *testing.T) {
	m := NewModel("match\nother\nmatch")
	m.height = 10
	m.width = 80
	m.searchQuery = "match"
	m.matchLines = []int{0, 2}
	m.matchIndex = 1

	// Press 'n' to wrap to first match
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	model := updated.(Model)
	if model.matchIndex != 0 {
		t.Fatalf("expected match index 0 after wrap, got %d", model.matchIndex)
	}
}

func TestModel_NNoMatches(t *testing.T) {
	m := NewModel("hello")
	m.height = 10

	// Press 'n' with no matches — should be no-op
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	model := updated.(Model)
	if model.matchIndex != -1 {
		t.Fatal("expected matchIndex to remain -1")
	}
}

func TestModel_ViewSmallHeight(t *testing.T) {
	m := NewModel("hello\nworld")
	m.height = 1
	view := m.View()
	if view != m.content {
		t.Fatal("expected raw content for height <= 1")
	}
}

func TestModel_ViewNormal(t *testing.T) {
	m := NewModel("line 1\nline 2\nline 3")
	m.height = 5
	m.width = 80
	view := m.View()
	if !strings.Contains(view, "line 1") {
		t.Fatal("expected line 1 in view")
	}
	if !strings.Contains(view, "scroll") {
		t.Fatal("expected help status in view")
	}
}

func TestModel_UnknownMsg(t *testing.T) {
	m := NewModel("hello")
	updated, cmd := m.Update("unhandled")
	model := updated.(Model)
	if cmd != nil {
		t.Fatal("expected nil cmd for unknown msg")
	}
	if model.offset != 0 {
		t.Fatal("expected no state change")
	}
}

func TestModel_VimKeys(t *testing.T) {
	lines := make([]string, 30)
	for i := range lines {
		lines[i] = "line"
	}
	m := NewModel(strings.Join(lines, "\n"))
	m.height = 10

	// 'j' for scroll down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model := updated.(Model)
	if model.offset != 1 {
		t.Fatalf("expected offset 1 after j, got %d", model.offset)
	}

	// 'k' for scroll up
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(Model)
	if model.offset != 0 {
		t.Fatalf("expected offset 0 after k, got %d", model.offset)
	}
}

func TestModel_PgDownPgUp(t *testing.T) {
	lines := make([]string, 30)
	for i := range lines {
		lines[i] = "line"
	}
	m := NewModel(strings.Join(lines, "\n"))
	m.height = 10

	// PgDown key
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	model := updated.(Model)
	if model.offset == 0 {
		t.Fatal("expected offset to increase after PgDown")
	}

	// PgUp key
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	model = updated.(Model)
	if model.offset != 0 {
		t.Fatalf("expected offset 0 after PgUp, got %d", model.offset)
	}
}
