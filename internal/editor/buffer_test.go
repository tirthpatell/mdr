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
