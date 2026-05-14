// Package duplicator detects duplicate values across environment variable keys.
// It identifies cases where multiple keys share the same value, which may indicate
// redundant configuration or copy-paste errors in .env files.
package duplicator

import (
	"fmt"
	"sort"
	"strings"
)

// Match represents a group of keys that share the same value.
type Match struct {
	Value string
	Keys  []string
}

// Result holds the outcome of a duplicate value scan.
type Result struct {
	Source  string
	Matches []Match
}

// Detect scans the provided env map and returns groups of keys that share
// identical values. Keys are sorted for deterministic output.
func Detect(source string, env map[string]string) Result {
	index := make(map[string][]string)
	for k, v := range env {
		if v == "" {
			continue
		}
		index[v] = append(index[v], k)
	}

	var matches []Match
	for val, keys := range index {
		if len(keys) < 2 {
			continue
		}
		sort.Strings(keys)
		matches = append(matches, Match{Value: val, Keys: keys})
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Keys[0] < matches[j].Keys[0]
	})

	return Result{Source: source, Matches: matches}
}

// RenderText returns a human-readable summary of duplicate value findings.
func RenderText(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Duplicate Value Report — %s\n", r.Source)
	fmt.Fprintf(&sb, "Duplicates found: %d\n", len(r.Matches))
	if len(r.Matches) == 0 {
		sb.WriteString("  No duplicate values detected.\n")
		return sb.String()
	}
	sb.WriteString("\n")
	for _, m := range r.Matches {
		fmt.Fprintf(&sb, "  Value: %q\n", m.Value)
		fmt.Fprintf(&sb, "  Keys:  %s\n\n", strings.Join(m.Keys, ", "))
	}
	return sb.String()
}
