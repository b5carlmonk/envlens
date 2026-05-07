package linter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenderText returns a human-readable summary of lint results.
func RenderText(r Result) string {
	if !r.HasIssues() {
		return "✔ No lint issues found.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("⚠ %d lint issue(s) found:\n", len(r.Issues)))
	for _, iss := range r.Issues {
		symbol := "[warn]"
		if iss.Level == LevelError {
			symbol = "[err] "
		}
		sb.WriteString(fmt.Sprintf("  %s %-30s %s\n", symbol, iss.Key, iss.Message))
	}
	return sb.String()
}

// RenderJSON returns a JSON-encoded representation of lint results.
func RenderJSON(r Result) (string, error) {
	type jsonIssue struct {
		Key     string `json:"key"`
		Level   string `json:"level"`
		Message string `json:"message"`
	}
	type jsonResult struct {
		Total  int          `json:"total"`
		Issues []jsonIssue  `json:"issues"`
	}

	out := jsonResult{Total: len(r.Issues)}
	for _, iss := range r.Issues {
		out.Issues = append(out.Issues, jsonIssue{
			Key:     iss.Key,
			Level:   string(iss.Level),
			Message: iss.Message,
		})
	}
	if out.Issues == nil {
		out.Issues = []jsonIssue{}
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
