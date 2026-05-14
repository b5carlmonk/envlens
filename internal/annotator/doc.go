// Package annotator classifies environment variable changes by attaching
// semantic tags and human-readable notes to each change entry.
//
// Each change is evaluated against a set of configurable rules:
//
//   - Sensitive: keys matching patterns like PASSWORD, TOKEN, SECRET, etc.
//   - Deprecated: keys containing user-supplied deprecated substrings.
//   - Internal: keys whose prefix matches a list of internal prefixes.
//   - Required: keys that are explicitly listed as required.
//   - Unknown: keys that do not match any rule.
//
// Usage:
//
//	opts := annotator.Options{
//		DeprecatedKeys:   []string{"LEGACY", "OLD_"},
//		RequiredKeys:     []string{"DATABASE_URL", "APP_ENV"},
//		InternalPrefixes: []string{"_INTERNAL"},
//	}
//	result := annotator.Apply(changes, opts)
//	for _, a := range result.Annotations {
//		fmt.Printf("%s [%s]: %s\n", a.Key, a.Tag, a.Note)
//	}
package annotator
