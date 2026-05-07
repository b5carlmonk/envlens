// Package validator provides utilities for validating environment variable
// keys and values against common rules and conventions.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// validKeyPattern matches conventional env var key names: uppercase letters,
// digits, and underscores, starting with a letter or underscore.
var validKeyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// Issue represents a single validation problem found in a change set.
type Issue struct {
	Key     string
	Message string
	Severity string // "warn" or "error"
}

// Result holds all issues found during validation.
type Result struct {
	Issues []Issue
}

// HasIssues returns true if any issues were found.
func (r Result) HasIssues() bool {
	return len(r.Issues) > 0
}

// Summary returns a short human-readable summary of the result.
func (r Result) Summary() string {
	if !r.HasIssues() {
		return "No validation issues found."
	}
	return fmt.Sprintf("%d validation issue(s) found.", len(r.Issues))
}

// Validate inspects a slice of differ.Change entries and returns a Result
// containing any detected issues such as non-conventional key names or
// suspiciously empty values.
func Validate(changes []differ.Change) Result {
	var issues []Issue

	for _, c := range changes {
		// Check key naming convention (only warn, not block).
		if !validKeyPattern.MatchString(c.Key) {
			issues = append(issues, Issue{
				Key:      c.Key,
				Message:  "key does not follow UPPER_SNAKE_CASE convention",
				Severity: "warn",
			})
		}

		// Warn when an added or modified value is empty.
		if c.Type == differ.Added || c.Type == differ.Modified {
			if strings.TrimSpace(c.NewValue) == "" {
				issues = append(issues, Issue{
					Key:      c.Key,
					Message:  "value is empty after add/modify",
					Severity: "warn",
				})
			}
		}

		// Error when a key name is blank.
		if strings.TrimSpace(c.Key) == "" {
			issues = append(issues, Issue{
				Key:      c.Key,
				Message:  "key name is blank",
				Severity: "error",
			})
		}
	}

	return Result{Issues: issues}
}
