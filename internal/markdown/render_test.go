package markdown

import (
	"os"
	"testing"
	"testing/iotest"
)

func TestRender_SimpleMarkdown(t *testing.T) {
	input := "# Hello\n\nThis is **bold** text.\n"
	result, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty rendered output")
	}
}

func TestRender_EmptyInput(t *testing.T) {
	result, err := Render("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty output even for empty input")
	}
}

func TestRenderFile_SimpleFixture(t *testing.T) {
	result, err := RenderFile("../../testdata/simple.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty rendered output")
	}
}

func TestRenderFile_NotFound(t *testing.T) {
	_, err := RenderFile("nonexistent.md")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestRenderFromReader(t *testing.T) {
	f, err := os.Open("../../testdata/simple.md")
	if err != nil {
		t.Fatalf("could not open fixture: %v", err)
	}
	defer f.Close()

	result, err := RenderFromReader(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty rendered output")
	}
}

func TestRenderFromReader_Error(t *testing.T) {
	_, err := RenderFromReader(iotest.ErrReader(os.ErrClosed))
	if err == nil {
		t.Fatal("expected error from bad reader")
	}
}

func TestRenderFile_ComplexFixture(t *testing.T) {
	result, err := RenderFile("../../testdata/complex.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty output for complex.md")
	}
}

func TestRender_WithCodeBlock(t *testing.T) {
	input := "# Code\n\n```go\nfmt.Println(\"hello\")\n```\n"
	result, err := Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty output with code block")
	}
}
