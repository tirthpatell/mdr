package linter

import (
	"strings"
	"testing"
)

func TestHeadingLevelSkip(t *testing.T) {
	input := []byte("# Title\n\n### Skipped H2\n")
	issues := checkHeadingHierarchy(input)
	if len(issues) == 0 {
		t.Fatal("expected heading skip issue")
	}
	if issues[0].Rule != "heading-hierarchy" {
		t.Fatalf("expected rule 'heading-hierarchy', got %q", issues[0].Rule)
	}
}

func TestHeadingLevelSkip_NoIssue(t *testing.T) {
	input := []byte("# Title\n\n## Subtitle\n\n### Sub-sub\n")
	issues := checkHeadingHierarchy(input)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestDuplicateHeadings(t *testing.T) {
	input := []byte("# Title\n\n# Title\n")
	issues := checkDuplicateHeadings(input)
	if len(issues) == 0 {
		t.Fatal("expected duplicate heading issue")
	}
}

func TestEmptyLinks(t *testing.T) {
	input := []byte("[Empty link]()\n")
	issues := checkEmptyLinks(input)
	if len(issues) == 0 {
		t.Fatal("expected empty link issue")
	}
}

func TestEmptyLinks_NoIssue(t *testing.T) {
	input := []byte("[Valid link](https://example.com)\n")
	issues := checkEmptyLinks(input)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestDuplicateHeadings_WithFormattedText(t *testing.T) {
	input := []byte("# Title with **bold** text\n\n# Title with **bold** text\n")
	issues := checkDuplicateHeadings(input)
	if len(issues) == 0 {
		t.Fatal("expected duplicate heading issue for headings with bold text")
	}
}

func TestDuplicateHeadings_DifferentText(t *testing.T) {
	input := []byte("# Hello\n\n# World\n")
	issues := checkDuplicateHeadings(input)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestTrailingWhitespace(t *testing.T) {
	input := []byte("# Title  \n\nClean line\nTrailing space \n")
	issues := checkTrailingWhitespace(input)
	if len(issues) != 2 {
		t.Fatalf("expected 2 trailing whitespace issues, got %d", len(issues))
	}
	if issues[0].Rule != "trailing-whitespace" {
		t.Fatalf("expected rule 'trailing-whitespace', got %q", issues[0].Rule)
	}
}

func TestTrailingWhitespace_NoIssue(t *testing.T) {
	input := []byte("# Title\n\nClean content\n")
	issues := checkTrailingWhitespace(input)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestEmptySections(t *testing.T) {
	input := []byte("# Title\n\n## Empty Section\n\n## Next Section\n\nSome content\n")
	issues := checkEmptySections(input)
	if len(issues) == 0 {
		t.Fatal("expected empty section issue")
	}
	if issues[0].Rule != "no-empty-sections" {
		t.Fatalf("expected rule 'no-empty-sections', got %q", issues[0].Rule)
	}
}

func TestEmptySections_NoIssue(t *testing.T) {
	input := []byte("# Title\n\nContent here.\n\n## Section\n\nMore content.\n")
	issues := checkEmptySections(input)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestEmptySections_EndOfDocument(t *testing.T) {
	input := []byte("# Title\n\nContent\n\n## Trailing\n")
	issues := checkEmptySections(input)
	if len(issues) == 0 {
		t.Fatal("expected empty section issue for heading at end of document")
	}
}

func TestIssue_String(t *testing.T) {
	issue := Issue{
		Rule:     "test-rule",
		Message:  "test message",
		Line:     42,
		Severity: SeverityError,
	}
	got := issue.String()
	expected := "line 42: [error] test message (test-rule)"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestIssue_String_Warning(t *testing.T) {
	issue := Issue{
		Rule:     "trailing-whitespace",
		Message:  "line has trailing whitespace",
		Line:     1,
		Severity: SeverityWarning,
	}
	got := issue.String()
	if !strings.Contains(got, "[warning]") {
		t.Fatalf("expected [warning] in string, got %q", got)
	}
}

func TestCheckHeadingHierarchy_MultipleSkips(t *testing.T) {
	input := []byte("# Title\n\n### Skip1\n\n##### Skip2\n")
	issues := checkHeadingHierarchy(input)
	if len(issues) != 2 {
		t.Fatalf("expected 2 heading skip issues, got %d", len(issues))
	}
}

func TestCheckEmptyLinks_MultipleLinks(t *testing.T) {
	input := []byte("[a]() and [b](https://example.com) and [c]()\n")
	issues := checkEmptyLinks(input)
	if len(issues) != 2 {
		t.Fatalf("expected 2 empty link issues, got %d", len(issues))
	}
}

func TestCheckTrailingWhitespace_Tabs(t *testing.T) {
	input := []byte("hello\t\nworld\n")
	issues := checkTrailingWhitespace(input)
	if len(issues) != 1 {
		t.Fatalf("expected 1 trailing whitespace issue, got %d", len(issues))
	}
}

func TestCheckEmptySections_SingleHeadingWithContent(t *testing.T) {
	input := []byte("# Title\n\nSome content here.\n")
	issues := checkEmptySections(input)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestCheckDuplicateHeadings_DifferentLevels(t *testing.T) {
	// Same text but different levels should not be flagged
	input := []byte("# Title\n\n## Title\n")
	issues := checkDuplicateHeadings(input)
	if len(issues) != 0 {
		t.Fatalf("expected no issues for different heading levels, got %d", len(issues))
	}
}

func TestLineNumber(t *testing.T) {
	input := []byte("# Title\n\nParagraph\n")
	doc := parseAST(input)
	line := lineNumber(input, doc.FirstChild())
	if line < 1 {
		t.Fatalf("expected positive line number, got %d", line)
	}
}
