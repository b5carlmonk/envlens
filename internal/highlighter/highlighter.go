package highlighter

import (
	"fmt"
	"strings"

	"github.com/envlens/internal/differ"
)

// Style represents a terminal color/style code.
type Style string

const (
	StyleReset  Style = "\033[0m"
	StyleRed    Style = "\033[31m"
	StyleGreen  Style = "\033[32m"
	StyleYellow Style = "\033[33m"
	StyleCyan   Style = "\033[36m"
	StyleBold   Style = "\033[1m"
)

// Options controls highlighter behaviour.
type Options struct {
	// NoColor disables ANSI escape codes (useful for piped output).
	NoColor bool
	// ShowOldValue includes the previous value for modified keys.
	ShowOldValue bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		NoColor:      false,
		ShowOldValue: true,
	}
}

// Line is a single highlighted output line.
type Line struct {
	Key    string
	Text   string
	Change differ.ChangeType
}

// Apply returns a slice of highlighted Lines for the given changes.
func Apply(changes []differ.Change, opts Options) []Line {
	lines := make([]Line, 0, len(changes))
	for _, c := range changes {
		lines = append(lines, formatLine(c, opts))
	}
	return lines
}

func formatLine(c differ.Change, opts Options) Line {
	var prefix, body string
	var style Style

	switch c.Type {
	case differ.Added:
		prefix = "+"
		style = StyleGreen
		body = fmt.Sprintf("%s=%s", c.Key, c.NewValue)
	case differ.Removed:
		prefix = "-"
		style = StyleRed
		body = fmt.Sprintf("%s=%s", c.Key, c.OldValue)
	case differ.Modified:
		prefix = "~"
		style = StyleYellow
		if opts.ShowOldValue {
			body = fmt.Sprintf("%s=%s  (was: %s)", c.Key, c.NewValue, c.OldValue)
		} else {
			body = fmt.Sprintf("%s=%s", c.Key, c.NewValue)
		}
	default:
		prefix = " "
		style = StyleReset
		body = fmt.Sprintf("%s=%s", c.Key, c.NewValue)
	}

	raw := fmt.Sprintf("%s %s", prefix, body)
	var text string
	if opts.NoColor {
		text = raw
	} else {
		text = colorize(style, raw)
	}

	return Line{Key: c.Key, Text: text, Change: c.Type}
}

func colorize(s Style, text string) string {
	return fmt.Sprintf("%s%s%s", s, text, StyleReset)
}

// Render returns all highlighted lines joined by newlines.
func Render(lines []Line) string {
	parts := make([]string, len(lines))
	for i, l := range lines {
		parts[i] = l.Text
	}
	return strings.Join(parts, "\n")
}
