package redactor

import (
	"strings"

	"github.com/yourusername/envlens/internal/differ"
	"github.com/yourusername/envlens/internal/masker"
)

// Options controls redaction behaviour.
type Options struct {
	// RedactAll masks every value regardless of key sensitivity.
	RedactAll bool
	// CustomSensitiveKeys are additional key substrings treated as sensitive.
	CustomSensitiveKeys []string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		RedactAll:           false,
		CustomSensitiveKeys: []string{},
	}
}

// Apply returns a new slice of Change values with sensitive values masked.
// Original Change values are never mutated.
func Apply(changes []differ.Change, opts Options) []differ.Change {
	result := make([]differ.Change, 0, len(changes))
	for _, c := range changes {
		result = append(result, redact(c, opts))
	}
	return result
}

func redact(c differ.Change, opts Options) differ.Change {
	sensitive := opts.RedactAll || masker.IsSensitive(c.Key) || isCustomSensitive(c.Key, opts.CustomSensitiveKeys)
	if !sensitive {
		return c
	}
	copied := c
	if copied.OldValue != "" {
		copied.OldValue = masker.Mask(copied.OldValue)
	}
	if copied.NewValue != "" {
		copied.NewValue = masker.Mask(copied.NewValue)
	}
	return copied
}

func isCustomSensitive(key string, customs []string) bool {
	lower := strings.ToLower(key)
	for _, sub := range customs {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
