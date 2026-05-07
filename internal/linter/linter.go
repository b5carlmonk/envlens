package linter

import (
	"fmt"
	"strings"

	"github.com/yourorg/envlens/internal/differ"
)

// IssueLevel represents the severity of a lint issue.
type IssueLevel string

const (
	LevelWarning IssueLevel = "warning"
	LevelError   IssueLevel = "error"
)

// Issue represents a single lint finding for an environment variable change.
type Issue struct {
	Key     string
	Level   IssueLevel
	Message string
}

// Result holds all issues found during linting.
type Result struct {
	Issues []Issue
}

// HasIssues returns true if any issues were found.
func (r *Result) HasIssues() bool {
	return len(r.Issues) > 0
}

// Rule is a function that inspects a change and returns zero or more issues.
type Rule func(c differ.Change) []Issue

// Lint runs all built-in rules against the provided changes.
func Lint(changes []differ.Change) Result {
	rules := []Rule{
		ruleNoEmptyValue,
		ruleKeyMustBeUpperSnake,
		ruleNoWhitespaceInValue,
	}

	var issues []Issue
	for _, c := range changes {
		for _, rule := range rules {
			issues = append(issues, rule(c)...)
		}
	}
	return Result{Issues: issues}
}

func ruleNoEmptyValue(c differ.Change) []Issue {
	if c.Type == differ.Added || c.Type == differ.Modified {
		if strings.TrimSpace(c.NewValue) == "" {
			return []Issue{{
				Key:     c.Key,
				Level:   LevelWarning,
				Message: "value is empty",
			}}
		}
	}
	return nil
}

func ruleKeyMustBeUpperSnake(c differ.Change) []Issue {
	for _, ch := range c.Key {
		if ch >= 'a' && ch <= 'z' {
			return []Issue{{
				Key:     c.Key,
				Level:   LevelWarning,
				Message: fmt.Sprintf("key %q is not UPPER_SNAKE_CASE", c.Key),
			}}
		}
	}
	return nil
}

func ruleNoWhitespaceInValue(c differ.Change) []Issue {
	val := c.NewValue
	if c.Type == differ.Removed {
		val = c.OldValue
	}
	if strings.ContainsAny(val, " \t") {
		return []Issue{{
			Key:     c.Key,
			Level:   LevelWarning,
			Message: "value contains unquoted whitespace",
		}}
	}
	return nil
}
