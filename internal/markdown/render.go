package markdown

import (
	"io"
	"os"

	"github.com/charmbracelet/glamour"
)

// Render takes a markdown string and returns ANSI-styled terminal output.
func Render(input string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return "", err
	}
	return renderer.Render(input)
}

// RenderFile reads a file and renders its markdown content.
func RenderFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return Render(string(data))
}

// RenderFromReader reads from an io.Reader and renders the markdown content.
func RenderFromReader(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return Render(string(data))
}
