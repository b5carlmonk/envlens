package stats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// RenderText returns a human-readable text representation of a stats Result.
func RenderText(r Result) string {
	var b strings.Builder
	b.WriteString("=== Environment Diff Statistics ===\n")
	b.WriteString(fmt.Sprintf("  Total changes : %d\n", r.Total))
	b.WriteString(fmt.Sprintf("  Added         : %d\n", r.Added))
	b.WriteString(fmt.Sprintf("  Removed       : %d\n", r.Removed))
	b.WriteString(fmt.Sprintf("  Modified      : %d\n", r.Modified))

	if len(r.TopPrefixes) > 0 {
		b.WriteString("\nTop key prefixes:\n")
		for _, p := range r.TopPrefixes {
			b.WriteString(fmt.Sprintf("  %-20s %d change(s)\n", p.Prefix, p.Count))
		}
	}
	return b.String()
}

// RenderJSON returns a JSON-encoded representation of a stats Result.
func RenderJSON(r Result) (string, error) {
	type jsonResult struct {
		Total       int           `json:"total"`
		Added       int           `json:"added"`
		Removed     int           `json:"removed"`
		Modified    int           `json:"modified"`
		TopPrefixes []PrefixCount `json:"top_prefixes"`
	}

	jr := jsonResult{
		Total:       r.Total,
		Added:       r.Added,
		Removed:     r.Removed,
		Modified:    r.Modified,
		TopPrefixes: r.TopPrefixes,
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jr); err != nil {
		return "", fmt.Errorf("stats: json encode: %w", err)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}
