package annotator

import (
	"strings"

	"github.com/yourusername/envlens/internal/differ"
)

// Tag represents a label attached to an environment variable change.
type Tag string

const (
	TagSensitive  Tag = "sensitive"
	TagDeprecated Tag = "deprecated"
	TagRequired   Tag = "required"
	TagInternal   Tag = "internal"
	TagUnknown    Tag = "unknown"
)

// Annotation holds metadata about a single change.
type Annotation struct {
	Key   string
	Tag   Tag
	Note  string
}

// Options configures annotation behaviour.
type Options struct {
	// DeprecatedKeys is a list of key substrings considered deprecated.
	DeprecatedKeys []string
	// RequiredKeys is a list of exact keys that must be present.
	RequiredKeys []string
	// InternalPrefixes marks keys with these prefixes as internal.
	InternalPrefixes []string
}

// Result contains all annotations produced for a set of changes.
type Result struct {
	Annotations []Annotation
}

// Apply examines each change and attaches an appropriate tag and note.
func Apply(changes []differ.Change, opts Options) Result {
	var annotations []Annotation

	for _, c := range changes {
		annotations = append(annotations, classify(c, opts))
	}

	return Result{Annotations: annotations}
}

func classify(c differ.Change, opts Options) Annotation {
	key := c.Key

	if isSensitive(key) {
		return Annotation{Key: key, Tag: TagSensitive, Note: "key appears to contain sensitive data"}
	}

	for _, dep := range opts.DeprecatedKeys {
		if containsIgnoreCase(key, dep) {
			return Annotation{Key: key, Tag: TagDeprecated, Note: "key matches a deprecated pattern: " + dep}
		}
	}

	for _, prefix := range opts.InternalPrefixes {
		if strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(prefix)) {
			return Annotation{Key: key, Tag: TagInternal, Note: "key uses an internal prefix: " + prefix}
		}
	}

	for _, req := range opts.RequiredKeys {
		if strings.EqualFold(key, req) {
			return Annotation{Key: key, Tag: TagRequired, Note: "key is marked as required"}
		}
	}

	return Annotation{Key: key, Tag: TagUnknown, Note: "no matching annotation rule"}
}

func isSensitive(key string) bool {
	sensitivePatterns := []string{"SECRET", "PASSWORD", "TOKEN", "PRIVATE", "API_KEY", "CREDENTIAL"}
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substr))
}
