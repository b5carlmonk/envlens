// Package ignorer provides utilities for excluding specific environment
// variable keys from diff and audit operations.
//
// Keys can be excluded explicitly by name or by prefix pattern. Ignore
// rules can be defined programmatically via Options or loaded from a
// plain-text ignore file (similar in spirit to .gitignore).
//
// Ignore file format:
//
//	# lines starting with # are comments
//	EXACT_KEY          # ignores a key by exact name
//	CI_*               # ignores all keys with prefix CI_
//
// Example usage:
//
//	opts, err := ignorer.FromFile(".envignore")
//	filtered := ignorer.Apply(parsed, opts)
package ignorer
