package templater

import (
	"fmt"
	"sort"
	"strings"
)

// Template represents a .env template with required and optional keys.
type Template struct {
	Required []string
	Optional []string
}

// ValidationResult holds the outcome of validating an env map against a template.
type ValidationResult struct {
	MissingRequired []string
	UnexpectedKeys  []string
	Present         []string
}

// HasIssues returns true if there are any missing required or unexpected keys.
func (r ValidationResult) HasIssues() bool {
	return len(r.MissingRequired) > 0 || len(r.UnexpectedKeys) > 0
}

// FromMap builds a Template from a map where values indicate "required" or "optional".
func FromMap(m map[string]string) (*Template, error) {
	t := &Template{}
	for k, v := range m {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "required":
			t.Required = append(t.Required, k)
		case "optional":
			t.Optional = append(t.Optional, k)
		default:
			return nil, fmt.Errorf("unknown directive %q for key %q: use 'required' or 'optional'", v, k)
		}
	}
	sort.Strings(t.Required)
	sort.Strings(t.Optional)
	return t, nil
}

// Validate checks an env map against the template and returns a ValidationResult.
func Validate(t *Template, env map[string]string) ValidationResult {
	result := ValidationResult{}

	known := make(map[string]bool)
	for _, k := range t.Required {
		known[k] = true
		if _, ok := env[k]; ok {
			result.Present = append(result.Present, k)
		} else {
			result.MissingRequired = append(result.MissingRequired, k)
		}
	}
	for _, k := range t.Optional {
		known[k] = true
		if _, ok := env[k]; ok {
			result.Present = append(result.Present, k)
		}
	}

	for k := range env {
		if !known[k] {
			result.UnexpectedKeys = append(result.UnexpectedKeys, k)
		}
	}

	sort.Strings(result.MissingRequired)
	sort.Strings(result.UnexpectedKeys)
	sort.Strings(result.Present)
	return result
}
