// Package truncator provides utilities for truncating long environment variable
// values in diffs and reports, making output more readable in terminals and logs.
package truncator

import (
	"github.com/user/envlens/internal/differ"
)

// Options controls truncation behaviour.
type Options struct {
	// MaxLength is the maximum number of characters to display for a value.
	// Values longer than this are truncated and suffixed with an ellipsis.
	// Defaults to 64 if zero.
	MaxLength int

	// Ellipsis is the string appended to truncated values. Defaults to "...".
	Ellipsis string
}

// DefaultOptions returns sensible defaults for truncation.
func DefaultOptions() Options {
	return Options{
		MaxLength: 64,
		Ellipsis:  "...",
	}
}

// Apply truncates the OldValue and NewValue fields of each change according to
// the provided options. The original slice is not mutated; a new slice is
// returned.
func Apply(changes []differ.Change, opts Options) []differ.Change {
	if opts.MaxLength <= 0 {
		opts.MaxLength = 64
	}
	if opts.Ellipsis == "" {
		opts.Ellipsis = "..."
	}

	result := make([]differ.Change, len(changes))
	for i, c := range changes {
		c.OldValue = truncate(c.OldValue, opts)
		c.NewValue = truncate(c.NewValue, opts)
		result[i] = c
	}
	return result
}

// truncate shortens s to opts.MaxLength runes, appending opts.Ellipsis when
// the string was actually shortened.
func truncate(s string, opts Options) string {
	runes := []rune(s)
	if len(runes) <= opts.MaxLength {
		return s
	}
	return string(runes[:opts.MaxLength]) + opts.Ellipsis
}
