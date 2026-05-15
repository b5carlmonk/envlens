// Package normalizer provides utilities for normalizing environment variable
// maps before comparison or diffing. It can trim whitespace, normalize key
// casing, and remove empty entries based on configurable options.
package normalizer

import (
	"strings"
)

// Options controls which normalization steps are applied.
type Options struct {
	// TrimValues removes leading and trailing whitespace from all values.
	TrimValues bool

	// UppercaseKeys converts all keys to UPPER_SNAKE_CASE.
	UppercaseKeys bool

	// RemoveEmpty drops entries where the value is empty after trimming.
	RemoveEmpty bool
}

// DefaultOptions returns a sensible default set of normalization options.
func DefaultOptions() Options {
	return Options{
		TrimValues:    true,
		UppercaseKeys: false,
		RemoveEmpty:   false,
	}
}

// Result holds the output of a normalization pass.
type Result struct {
	// Env is the normalized environment map.
	Env map[string]string

	// Trimmed lists keys whose values were changed by whitespace trimming.
	Trimmed []string

	// Removed lists keys that were dropped because their value was empty.
	Removed []string

	// Renamed maps original key names to their normalized (uppercased) forms.
	Renamed map[string]string
}

// Apply normalizes the given environment map according to opts and returns
// a Result describing every transformation that was performed.
func Apply(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	result := Result{
		Renamed: make(map[string]string),
	}

	for k, v := range env {
		newKey := k
		if opts.UppercaseKeys {
			upper := strings.ToUpper(k)
			if upper != k {
				result.Renamed[k] = upper
				newKey = upper
			}
		}

		newVal := v
		if opts.TrimValues {
			trimmed := strings.TrimSpace(v)
			if trimmed != v {
				result.Trimmed = append(result.Trimmed, newKey)
			}
			newVal = trimmed
		}

		if opts.RemoveEmpty && newVal == "" {
			result.Removed = append(result.Removed, newKey)
			continue
		}

		out[newKey] = newVal
	}

	result.Env = out
	return result
}
