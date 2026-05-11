package patcher

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenderText returns a human-readable summary of the patch result.
func RenderText(r Result) string {
	var sb strings.Builder

	mode := "applied"
	if r.DryRun {
		mode = "dry-run"
	}
	fmt.Fprintf(&sb, "Patch result [%s]\n", mode)
	fmt.Fprintf(&sb, "  Applied  : %d\n", len(r.Applied))
	fmt.Fprintf(&sb, "  Skipped  : %d\n", len(r.Skipped))
	fmt.Fprintf(&sb, "  Conflicts: %d\n", len(r.Conflicts))

	if len(r.Applied) > 0 {
		sb.WriteString("\nApplied keys:\n")
		for _, k := range r.Applied {
			fmt.Fprintf(&sb, "  + %s\n", k)
		}
	}
	if len(r.Skipped) > 0 {
		sb.WriteString("\nSkipped keys:\n")
		for _, k := range r.Skipped {
			fmt.Fprintf(&sb, "  ~ %s\n", k)
		}
	}
	if len(r.Conflicts) > 0 {
		sb.WriteString("\nConflict keys (overwritten):\n")
		for _, k := range r.Conflicts {
			fmt.Fprintf(&sb, "  ! %s\n", k)
		}
	}

	return sb.String()
}

// RenderJSON returns the patch result as a JSON string.
func RenderJSON(r Result) (string, error) {
	type payload struct {
		DryRun    bool     `json:"dry_run"`
		Applied   []string `json:"applied"`
		Skipped   []string `json:"skipped"`
		Conflicts []string `json:"conflicts"`
	}
	p := payload{
		DryRun:    r.DryRun,
		Applied:   r.Applied,
		Skipped:   r.Skipped,
		Conflicts: r.Conflicts,
	}
	if p.Applied == nil {
		p.Applied = []string{}
	}
	if p.Skipped == nil {
		p.Skipped = []string{}
	}
	if p.Conflicts == nil {
		p.Conflicts = []string{}
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
