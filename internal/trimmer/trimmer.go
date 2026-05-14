// Package trimmer detects and removes leading/trailing whitespace
// from environment variable values across a parsed env map.
package trimmer

import (
	"strings"

	"github.com/yourusername/envlens/internal/differ"
)

// Result holds the outcome of a trim operation.
type Result struct {
	// Trimmed maps key names to their original (untrimmed) values.
	Trimmed map[string]string
	// Clean is the sanitized env map with all values trimmed.
	Clean map[string]string
}

// Apply scans the provided env map for values that contain leading or
// trailing whitespace, trims them, and returns a Result describing what
// was changed. The original map is never mutated.
func Apply(env map[string]string) Result {
	clean := make(map[string]string, len(env))
	trimmed := make(map[string]string)

	for k, v := range env {
		t := strings.TrimSpace(v)
		clean[k] = t
		if t != v {
			trimmed[k] = v
		}
	}

	return Result{
		Trimmed: trimmed,
		Clean:   clean,
	}
}

// ToChanges converts a Result into a slice of differ.Change values so the
// trimmed entries can be fed into downstream pipeline stages (reporter,
// exporter, etc.) as modified keys.
func ToChanges(r Result) []differ.Change {
	changes := make([]differ.Change, 0, len(r.Trimmed))
	for k, orig := range r.Trimmed {
		changes = append(changes, differ.Change{
			Key:    k,
			Type:   differ.Modified,
			OldVal: orig,
			NewVal: r.Clean[k],
		})
	}
	return changes
}
