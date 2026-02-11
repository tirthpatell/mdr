# mdr

[![Go Report Card](https://goreportcard.com/badge/github.com/tirthpatell/mdr)](https://goreportcard.com/report/github.com/tirthpatell/mdr)
[![Go Reference](https://pkg.go.dev/badge/github.com/tirthpatell/mdr.svg)](https://pkg.go.dev/github.com/tirthpatell/mdr)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/tirthpatell/mdr)](https://github.com/tirthpatell/mdr/releases)

Markdown renderer, editor, and linter for the terminal.

<p align="center">
  <img src="demo.gif" alt="mdr demo" width="600">
</p>

## Install

### Homebrew

```bash
brew install tirthpatell/tap/mdr
```

### Go

```bash
go install github.com/tirthpatell/mdr@latest
```

### GitHub Releases

Download binaries from the [releases page](https://github.com/tirthpatell/mdr/releases).

## Usage

### View

Render a markdown file in a scrollable TUI:

```bash
mdr view README.md
```

Pipe from stdin:

```bash
cat README.md | mdr view
```

Print rendered output without TUI:

```bash
mdr view --raw README.md
```

### Edit

Open the TUI editor with live preview:

```bash
mdr edit README.md
```

**Editor keys:** Arrow keys to move, type to edit, Ctrl+S to save, Ctrl+C to quit.

### Lint

Check markdown files for structural issues:

```bash
mdr lint README.md docs/*.md
```

**Rules:**
- `heading-hierarchy` - Heading levels should not skip (e.g., H1 to H3)
- `duplicate-heading` - Headings at the same level should not repeat
- `empty-link` - Links must have a destination

## License

MIT
