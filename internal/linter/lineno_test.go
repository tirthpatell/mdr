package linter

import "testing"

func TestEmptyLinkNoTextLineNumber(t *testing.T) {
	// A link with no text and no destination has no Text descendant, so
	// the line must come from the nearest block ancestor, not default to 0.
	input := []byte("Some text\n\nA line with []() empty link\n")
	issues := checkEmptyLinks(parseAST(input), input)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Line != 3 {
		t.Fatalf("expected line 3, got %d", issues[0].Line)
	}
}

func TestEmptyLinkWithTextLineNumber(t *testing.T) {
	input := []byte("# Title\n\nSee [broken]() here\n")
	issues := checkEmptyLinks(parseAST(input), input)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Line != 3 {
		t.Fatalf("expected line 3, got %d", issues[0].Line)
	}
}
