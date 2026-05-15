package tracer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenderText returns a human-readable text representation of the trace.
func RenderText(t *Trace) string {
	if t.Len() == 0 {
		return "No trace entries recorded.\n"
	}

	var sb strings.Builder
	sb.WriteString("=== Trace Log ===\n")
	for i, e := range t.Entries() {
		sb.WriteString(fmt.Sprintf("[%d] %s | %s -> %s | %s\n",
			i+1,
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Source,
			e.Target,
			e.Label,
		))
		for _, c := range e.Changes {
			sb.WriteString(fmt.Sprintf("    [%s] %s\n", c.Type, c.Key))
		}
	}
	return sb.String()
}

// RenderJSON returns a JSON-encoded representation of all trace entries.
func RenderJSON(t *Trace) (string, error) {
	type jsonEntry struct {
		Timestamp string `json:"timestamp"`
		Label     string `json:"label"`
		Source    string `json:"source"`
		Target    string `json:"target"`
		ChangeCount int  `json:"change_count"`
	}

	var entries []jsonEntry
	for _, e := range t.Entries() {
		entries = append(entries, jsonEntry{
			Timestamp:   e.Timestamp.Format("2006-01-02T15:04:05Z"),
			Label:       e.Label,
			Source:      e.Source,
			Target:      e.Target,
			ChangeCount: len(e.Changes),
		})
	}

	if entries == nil {
		entries = []jsonEntry{}
	}

	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return "", fmt.Errorf("tracer: failed to render JSON: %w", err)
	}
	return string(b), nil
}
