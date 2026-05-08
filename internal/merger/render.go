package merger

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// RenderText returns a human-readable summary of a merge Result.
func RenderText(r Result) string {
	var sb strings.Builder

	keys := make([]string, 0, len(r.Merged))
	for k := range r.Merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sb.WriteString(fmt.Sprintf("Merged keys: %d\n", len(r.Merged)))
	sb.WriteString(fmt.Sprintf("Conflicts resolved: %d\n", len(r.Conflicts)))

	if len(r.Conflicts) > 0 {
		sb.WriteString("\nConflicts:\n")
		for _, c := range r.Conflicts {
			sb.WriteString(fmt.Sprintf("  ~ %s  source=%q  target=%q  resolved=%q\n",
				c.Key, c.SourceValue, c.TargetValue, c.Resolved))
		}
	}

	sb.WriteString("\nFinal env:\n")
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("  %s=%s\n", k, r.Merged[k]))
	}

	return sb.String()
}

// RenderJSON returns a JSON-encoded representation of a merge Result.
func RenderJSON(r Result) (string, error) {
	type jsonResult struct {
		Merged    map[string]string `json:"merged"`
		Conflicts []Conflict        `json:"conflicts"`
	}
	out := jsonResult{Merged: r.Merged, Conflicts: r.Conflicts}
	if out.Conflicts == nil {
		out.Conflicts = []Conflict{}
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("merger: json render: %w", err)
	}
	return string(b), nil
}
