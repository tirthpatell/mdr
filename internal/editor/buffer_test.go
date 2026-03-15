package editor

import (
	"testing"
)

func TestNewBuffer(t *testing.T) {
	b := NewBuffer("hello\nworld")
	if b.LineCount() != 2 {
		t.Fatalf("expected 2 lines, got %d", b.LineCount())
	}
}

func TestNewBuffer_Empty(t *testing.T) {
	b := NewBuffer("")
	if b.LineCount() != 1 {
		t.Fatalf("expected 1 line for empty buffer, got %d", b.LineCount())
	}
}

func TestBuffer_GetLine(t *testing.T) {
	b := NewBuffer("line one\nline two")
	if got := b.GetLine(0); got != "line one" {
		t.Fatalf("expected 'line one', got %q", got)
	}
	if got := b.GetLine(1); got != "line two" {
		t.Fatalf("expected 'line two', got %q", got)
	}
}

func TestBuffer_InsertChar(t *testing.T) {
	b := NewBuffer("hello")
	b.InsertChar(0, 5, '!')
	if got := b.GetLine(0); got != "hello!" {
		t.Fatalf("expected 'hello!', got %q", got)
	}
}

func TestBuffer_InsertChar_Middle(t *testing.T) {
	b := NewBuffer("hllo")
	b.InsertChar(0, 1, 'e')
	if got := b.GetLine(0); got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestBuffer_DeleteChar(t *testing.T) {
	b := NewBuffer("hello")
	b.DeleteChar(0, 4) // delete 'o'
	if got := b.GetLine(0); got != "hell" {
		t.Fatalf("expected 'hell', got %q", got)
	}
}

func TestBuffer_DeleteChar_Backspace(t *testing.T) {
	b := NewBuffer("hello")
	b.Backspace(0, 5) // delete 'o' via backspace at pos 5
	if got := b.GetLine(0); got != "hell" {
		t.Fatalf("expected 'hell', got %q", got)
	}
}

func TestBuffer_InsertNewline(t *testing.T) {
	b := NewBuffer("hello world")
	b.InsertNewline(0, 5) // split at position 5
	if b.LineCount() != 2 {
		t.Fatalf("expected 2 lines, got %d", b.LineCount())
	}
	if got := b.GetLine(0); got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
	if got := b.GetLine(1); got != " world" {
		t.Fatalf("expected ' world', got %q", got)
	}
}

func TestBuffer_JoinLines(t *testing.T) {
	b := NewBuffer("hello\nworld")
	col := b.JoinLines(1) // join line 1 with line 0
	if b.LineCount() != 1 {
		t.Fatalf("expected 1 line, got %d", b.LineCount())
	}
	if got := b.GetLine(0); got != "helloworld" {
		t.Fatalf("expected 'helloworld', got %q", got)
	}
	if col != 5 {
		t.Fatalf("expected cursor col 5, got %d", col)
	}
}

func TestBuffer_String(t *testing.T) {
	b := NewBuffer("hello\nworld")
	b.InsertChar(0, 5, '!')
	expected := "hello!\nworld"
	if got := b.String(); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestBuffer_Modified(t *testing.T) {
	b := NewBuffer("hello")
	if b.Modified() {
		t.Fatal("new buffer should not be modified")
	}
	b.InsertChar(0, 5, '!')
	if !b.Modified() {
		t.Fatal("buffer should be modified after insert")
	}
}

func TestBuffer_GetLine_OutOfBounds(t *testing.T) {
	b := NewBuffer("hello")
	if got := b.GetLine(-1); got != "" {
		t.Fatalf("expected empty string for negative row, got %q", got)
	}
	if got := b.GetLine(5); got != "" {
		t.Fatalf("expected empty string for row beyond end, got %q", got)
	}
}

func TestBuffer_LineLen_OutOfBounds(t *testing.T) {
	b := NewBuffer("hello")
	if got := b.LineLen(-1); got != 0 {
		t.Fatalf("expected 0 for negative row, got %d", got)
	}
	if got := b.LineLen(5); got != 0 {
		t.Fatalf("expected 0 for row beyond end, got %d", got)
	}
}

func TestBuffer_InsertChar_OutOfBounds(t *testing.T) {
	b := NewBuffer("hello")
	// Invalid row — should be no-op
	b.InsertChar(-1, 0, 'x')
	b.InsertChar(5, 0, 'x')
	if got := b.GetLine(0); got != "hello" {
		t.Fatalf("expected unchanged line, got %q", got)
	}
	// Col beyond line length — should clamp
	b.InsertChar(0, 100, '!')
	if got := b.GetLine(0); got != "hello!" {
		t.Fatalf("expected 'hello!' with clamped col, got %q", got)
	}
}

func TestBuffer_DeleteChar_OutOfBounds(t *testing.T) {
	b := NewBuffer("hello")
	// Invalid row
	b.DeleteChar(-1, 0)
	b.DeleteChar(5, 0)
	// Invalid col
	b.DeleteChar(0, -1)
	b.DeleteChar(0, 10)
	if got := b.GetLine(0); got != "hello" {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestBuffer_Backspace_AtStart(t *testing.T) {
	b := NewBuffer("hello")
	b.Backspace(0, 0) // at start of line — should be no-op
	if got := b.GetLine(0); got != "hello" {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestBuffer_InsertNewline_OutOfBounds(t *testing.T) {
	b := NewBuffer("hello")
	b.InsertNewline(-1, 0)
	b.InsertNewline(5, 0)
	if b.LineCount() != 1 {
		t.Fatalf("expected 1 line after out-of-bounds newline inserts, got %d", b.LineCount())
	}
}

func TestBuffer_InsertNewline_ColClamp(t *testing.T) {
	b := NewBuffer("hello")
	b.InsertNewline(0, 100) // col beyond length, should clamp
	if b.LineCount() != 2 {
		t.Fatalf("expected 2 lines, got %d", b.LineCount())
	}
	if got := b.GetLine(0); got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
	if got := b.GetLine(1); got != "" {
		t.Fatalf("expected empty second line, got %q", got)
	}
}

func TestBuffer_JoinLines_OutOfBounds(t *testing.T) {
	b := NewBuffer("hello\nworld")
	// row 0 — can't join with row above
	col := b.JoinLines(0)
	if col != 0 {
		t.Fatalf("expected 0 for out-of-bounds join, got %d", col)
	}
	// row beyond end
	col = b.JoinLines(5)
	if col != 0 {
		t.Fatalf("expected 0 for out-of-bounds join, got %d", col)
	}
	if b.LineCount() != 2 {
		t.Fatalf("expected unchanged line count 2, got %d", b.LineCount())
	}
}

func TestBuffer_LineLen_Unicode(t *testing.T) {
	b := NewBuffer("héllo 世界")
	if got := b.LineLen(0); got != 8 {
		t.Fatalf("expected 8 runes, got %d", got)
	}
}

func TestBuffer_ResetModified(t *testing.T) {
	b := NewBuffer("hello")
	b.InsertChar(0, 0, 'x')
	if !b.Modified() {
		t.Fatal("expected modified after insert")
	}
	b.ResetModified()
	if b.Modified() {
		t.Fatal("expected not modified after reset")
	}
}
