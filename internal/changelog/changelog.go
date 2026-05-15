package changelog

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/envlens/internal/differ"
)

// Entry represents a single changelog entry for a deployment diff.
type Entry struct {
	Timestamp time.Time
	Source    string
	Target    string
	Changes   []differ.Change
}

// Changelog holds an ordered list of diff entries over time.
type Changelog struct {
	Entries []Entry
}

// New creates a new Changelog with no entries.
func New() *Changelog {
	return &Changelog{}
}

// Add appends a new entry to the changelog.
func (c *Changelog) Add(source, target string, changes []differ.Change) {
	c.Entries = append(c.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Source:    source,
		Target:    target,
		Changes:   changes,
	})
}

// Len returns the number of entries.
func (c *Changelog) Len() int {
	return len(c.Entries)
}

// RenderText returns a human-readable changelog summary.
func RenderText(c *Changelog) string {
	if len(c.Entries) == 0 {
		return "No changelog entries recorded.\n"
	}

	var sb strings.Builder
	sb.WriteString("=== Changelog ===\n")

	for i, e := range c.Entries {
		sb.WriteString(fmt.Sprintf("\n[%d] %s\n", i+1, e.Timestamp.Format(time.RFC3339)))
		sb.WriteString(fmt.Sprintf("    Source : %s\n", e.Source))
		sb.WriteString(fmt.Sprintf("    Target : %s\n", e.Target))
		sb.WriteString(fmt.Sprintf("    Changes: %d\n", len(e.Changes)))

		for _, ch := range e.Changes {
			sb.WriteString(fmt.Sprintf("      [%s] %s\n", strings.ToUpper(string(ch.Type)), ch.Key))
		}
	}

	return sb.String()
}
