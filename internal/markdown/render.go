package markdown

import (
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/glamour"
)

// The renderer is cached because glamour.NewTermRenderer is expensive
// (style detection and setup) and the editor renders a preview on every
// keystroke. TermRenderer is not safe for concurrent use, so guard it.
var (
	rendererMu   sync.Mutex
	renderer     *glamour.TermRenderer
	rendererErr  error
	rendererOnce sync.Once
)

// Render takes a markdown string and returns ANSI-styled terminal output.
func Render(input string) (string, error) {
	rendererOnce.Do(func() {
		renderer, rendererErr = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(0),
		)
	})
	if rendererErr != nil {
		return "", rendererErr
	}
	rendererMu.Lock()
	defer rendererMu.Unlock()
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
