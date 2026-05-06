// Package filter provides utilities for filtering diff results
// based on key patterns, prefixes, or change types.
package filter

import (
	"strings"

	"github.com/user/envlens/internal/differ"
)

// Options holds the configuration for filtering diff results.
type Options struct {
	// Prefix filters changes to keys that start with this string.
	Prefix string
	// KeyContains filters changes to keys that contain this substring.
	KeyContains string
	// Types restricts results to specific change types (e.g., "added", "removed", "modified").
	Types []string
}

// Apply filters a slice of differ.Change according to the given Options.
// It returns a new slice containing only the changes that match all specified criteria.
func Apply(changes []differ.Change, opts Options) []differ.Change {
	var result []differ.Change

	typeSet := make(map[string]bool, len(opts.Types))
	for _, t := range opts.Types {
		typeSet[strings.ToLower(t)] = true
	}

	for _, c := range changes {
		if opts.Prefix != "" && !strings.HasPrefix(c.Key, opts.Prefix) {
			continue
		}
		if opts.KeyContains != "" && !strings.Contains(c.Key, opts.KeyContains) {
			continue
		}
		if len(typeSet) > 0 {
			kind := strings.ToLower(string(c.Type))
			if !typeSet[kind] {
				continue
			}
		}
		result = append(result, c)
	}

	return result
}
