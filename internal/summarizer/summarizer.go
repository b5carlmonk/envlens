// Package summarizer provides functionality to generate human-readable
// summaries of environment variable diffs, including counts by change type
// and a brief narrative description of the overall diff.
package summarizer

import (
	"fmt"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// Summary holds aggregated statistics and a narrative about a diff result.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Unchanged int
	Total    int
	Narrative string
}

// Summarize computes a Summary from a slice of differ.Change values.
func Summarize(changes []differ.Change) Summary {
	s := Summary{}
	for _, c := range changes {
		s.Total++
		switch c.Type {
		case differ.Added:
			s.Added++
		case differ.Removed:
			s.Removed++
		case differ.Modified:
			s.Modified++
		case differ.Unchanged:
			s.Unchanged++
		}
	}
	s.Narrative = buildNarrative(s)
	return s
}

// buildNarrative constructs a plain-English description of the diff.
func buildNarrative(s Summary) string {
	if s.Total == 0 {
		return "No environment variables were compared."
	}

	parts := []string{}
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s.Added))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s.Removed))
	}
	if s.Modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", s.Modified))
	}
	if s.Unchanged > 0 {
		parts = append(parts, fmt.Sprintf("%d unchanged", s.Unchanged))
	}

	if len(parts) == 0 {
		return "No changes detected."
	}

	return fmt.Sprintf("Compared %d variable(s): %s.", s.Total, strings.Join(parts, ", "))
}
