package linter

import (
	"strings"
	"testing"
)

func TestLint_AllRules(t *testing.T) {
	input := []byte("# Title\n\n### Skip\n\n# Title\n\n[empty]()\n")
	issues := Lint(input)
	if len(issues) < 3 {
		t.Fatalf("expected at least 3 issues, got %d", len(issues))
	}

	rules := make(map[string]bool)
	for _, issue := range issues {
		rules[issue.Rule] = true
	}
	for _, expected := range []string{"heading-hierarchy", "duplicate-heading", "empty-link"} {
		if !rules[expected] {
			t.Errorf("expected rule %q in issues", expected)
		}
	}
}

func TestLintFile(t *testing.T) {
	issues, err := LintFile("../../testdata/bad-headings.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issues in bad-headings.md")
	}
}

func TestLintFile_NotFound(t *testing.T) {
	_, err := LintFile("nonexistent.md")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLint_CleanFile(t *testing.T) {
	input := []byte("# Title\n\nSome intro text.\n\n## Subtitle\n\n[Link](https://example.com)\n")
	issues := Lint(input)
	if len(issues) != 0 {
		t.Fatalf("expected no issues for clean file, got %d: %v", len(issues), issues)
	}
}

func TestLintReader(t *testing.T) {
	r := strings.NewReader("# Title\n\n### Skip\n")
	issues, err := LintReader(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issues from reader input")
	}
}

func TestLint_Sorted(t *testing.T) {
	// Issues from multiple rules should be sorted by line number
	input := []byte("# Title  \n\n### Skipped\n\n# Title\n\n[empty]()\n")
	issues := Lint(input)
	for i := 1; i < len(issues); i++ {
		if issues[i].Line < issues[i-1].Line {
			t.Fatalf("issues not sorted: line %d before line %d", issues[i-1].Line, issues[i].Line)
		}
	}
}

func TestLintFile_CleanFile(t *testing.T) {
	issues, err := LintFile("../../testdata/simple.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected no issues for simple.md, got %d", len(issues))
	}
}

func TestLintReader_Clean(t *testing.T) {
	r := strings.NewReader("# Title\n\nContent.\n")
	issues, err := LintReader(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected no issues for clean content, got %d", len(issues))
	}
}
