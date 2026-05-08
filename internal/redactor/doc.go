// Package redactor provides value-masking for environment variable changes
// before they are rendered or exported.
//
// It wraps the masker package and applies redaction rules across a slice of
// differ.Change values, returning new copies with sensitive values replaced
// by masked equivalents. Original Change values are never mutated.
//
// Usage:
//
//	opts := redactor.DefaultOptions()
//	opts.CustomSensitiveKeys = []string{"stripe", "twilio"}
//	safeChanges := redactor.Apply(changes, opts)
//
// The RedactAll option can be used to mask every value unconditionally,
// which is useful when sharing diffs in public or untrusted contexts.
package redactor
