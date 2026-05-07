// Package masker provides utilities for redacting sensitive environment
// variable values before display or export, ensuring secrets are not
// accidentally leaked in reports or logs.
package masker

import "strings"

// sensitivePatterns holds substrings that indicate a key may contain a secret.
var sensitivePatterns = []string{
	"secret", "password", "passwd", "token", "apikey", "api_key",
	"auth", "credential", "private", "key", "cert", "pwd",
}

// Options controls masking behaviour.
type Options struct {
	// MaskChar is the character used to build the masked string. Defaults to "*".
	MaskChar string
	// VisibleSuffix is the number of trailing characters to leave visible.
	// Set to 0 to hide the value entirely.
	VisibleSuffix int
}

// DefaultOptions returns sensible masking defaults.
func DefaultOptions() Options {
	return Options{
		MaskChar:      "*",
		VisibleSuffix: 4,
	}
}

// IsSensitive reports whether the given key name looks like it holds a secret.
func IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// Mask returns a redacted version of value according to opts.
// If the value is shorter than or equal to VisibleSuffix, the entire
// value is replaced with mask characters.
func Mask(value string, opts Options) string {
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}
	const maskLen = 6
	if opts.VisibleSuffix <= 0 || len(value) <= opts.VisibleSuffix {
		return strings.Repeat(opts.MaskChar, maskLen)
	}
	visible := value[len(value)-opts.VisibleSuffix:]
	return strings.Repeat(opts.MaskChar, maskLen) + visible
}

// MaskIfSensitive masks value only when the key is considered sensitive.
// Otherwise it returns the original value unchanged.
func MaskIfSensitive(key, value string, opts Options) string {
	if IsSensitive(key) {
		return Mask(value, opts)
	}
	return value
}
