package linter

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

type Issue struct {
	Rule     string
	Message  string
	Line     int
	Severity Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("line %d: [%s] %s (%s)", i.Line, i.Severity, i.Message, i.Rule)
}

func parseAST(source []byte) ast.Node {
	md := goldmark.New()
	reader := text.NewReader(source)
	return md.Parser().Parse(reader)
}

func checkHeadingHierarchy(doc ast.Node, source []byte) []Issue {
	var issues []Issue
	prevLevel := 0

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if heading, ok := n.(*ast.Heading); ok {
			if prevLevel > 0 && heading.Level > prevLevel+1 {
				line := lineNumber(source, n)
				issues = append(issues, Issue{
					Rule:     "heading-hierarchy",
					Message:  fmt.Sprintf("heading level skipped from H%d to H%d", prevLevel, heading.Level),
					Line:     line,
					Severity: SeverityWarning,
				})
			}
			prevLevel = heading.Level
		}
		return ast.WalkContinue, nil
	})
	return issues
}

func checkDuplicateHeadings(doc ast.Node, source []byte) []Issue {
	var issues []Issue
	seen := make(map[string]int) // text -> first line

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if heading, ok := n.(*ast.Heading); ok {
			text := extractText(heading, source)
			key := fmt.Sprintf("h%d:%s", heading.Level, text)
			line := lineNumber(source, n)
			if firstLine, exists := seen[key]; exists {
				issues = append(issues, Issue{
					Rule:     "duplicate-heading",
					Message:  fmt.Sprintf("duplicate heading %q (first at line %d)", text, firstLine),
					Line:     line,
					Severity: SeverityWarning,
				})
			} else {
				seen[key] = line
			}
		}
		return ast.WalkContinue, nil
	})
	return issues
}

func checkEmptyLinks(doc ast.Node, source []byte) []Issue {
	var issues []Issue

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if link, ok := n.(*ast.Link); ok {
			if len(link.Destination) == 0 {
				line := lineNumber(source, n)
				issues = append(issues, Issue{
					Rule:     "empty-link",
					Message:  "link has empty destination",
					Line:     line,
					Severity: SeverityError,
				})
			}
		}
		return ast.WalkContinue, nil
	})
	return issues
}

func extractText(n ast.Node, source []byte) string {
	var buf bytes.Buffer
	ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := child.(*ast.Text); ok {
			buf.Write(t.Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return buf.String()
}

func checkTrailingWhitespace(source []byte) []Issue {
	var issues []Issue
	lines := bytes.Split(source, []byte("\n"))
	for i, line := range lines {
		trimmed := bytes.TrimRight(line, " \t")
		if len(trimmed) < len(line) {
			issues = append(issues, Issue{
				Rule:     "trailing-whitespace",
				Message:  "line has trailing whitespace",
				Line:     i + 1,
				Severity: SeverityWarning,
			})
		}
	}
	return issues
}

func checkEmptySections(doc ast.Node, source []byte) []Issue {
	var issues []Issue
	var prevHeading *ast.Heading
	var prevLine int

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if heading, ok := n.(*ast.Heading); ok {
			if prevHeading != nil && prevHeading.NextSibling() == heading {
				issues = append(issues, Issue{
					Rule:     "no-empty-sections",
					Message:  fmt.Sprintf("empty section under %q", extractText(prevHeading, source)),
					Line:     prevLine,
					Severity: SeverityWarning,
				})
			}
			prevHeading = heading
			prevLine = lineNumber(source, n)
		}
		return ast.WalkContinue, nil
	})

	// Flag a heading at the end of the document with no content after it
	if prevHeading != nil && prevHeading.NextSibling() == nil {
		issues = append(issues, Issue{
			Rule:     "no-empty-sections",
			Message:  fmt.Sprintf("empty section under %q", extractText(prevHeading, source)),
			Line:     prevLine,
			Severity: SeverityWarning,
		})
	}

	return issues
}

func lineNumber(source []byte, n ast.Node) int {
	// Only block nodes support Lines(); inline nodes panic
	if n.Type() == ast.TypeBlock {
		lines := n.Lines()
		if lines.Len() > 0 {
			seg := lines.At(0)
			return bytes.Count(source[:seg.Start], []byte("\n")) + 1
		}
	}
	// fallback: first text descendant at any depth
	if start, ok := firstTextStart(n); ok {
		return bytes.Count(source[:start], []byte("\n")) + 1
	}
	// fallback for inline nodes with no text (e.g. `[]()`): nearest block ancestor
	for parent := n.Parent(); parent != nil; parent = parent.Parent() {
		if parent.Type() == ast.TypeBlock {
			lines := parent.Lines()
			if lines.Len() > 0 {
				seg := lines.At(0)
				return bytes.Count(source[:seg.Start], []byte("\n")) + 1
			}
		}
	}
	return 0
}

// firstTextStart returns the source offset of the first Text descendant of n.
func firstTextStart(n ast.Node) (int, bool) {
	if t, ok := n.(*ast.Text); ok {
		return t.Segment.Start, true
	}
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if start, ok := firstTextStart(child); ok {
			return start, true
		}
	}
	return 0, false
}
