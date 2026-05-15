package tracer

import (
	"fmt"
	"time"

	"github.com/envlens/internal/differ"
)

// Entry represents a single traced change event with metadata.
type Entry struct {
	Timestamp time.Time
	Label     string
	Source    string
	Target    string
	Changes   []differ.Change
}

// Trace holds an ordered sequence of traced diff entries.
type Trace struct {
	entries []Entry
}

// New returns an empty Trace.
func New() *Trace {
	return &Trace{}
}

// Add appends a new entry to the trace with the current timestamp.
func (t *Trace) Add(label, source, target string, changes []differ.Change) {
	t.entries = append(t.entries, Entry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Source:    source,
		Target:    target,
		Changes:   changes,
	})
}

// Entries returns all recorded trace entries.
func (t *Trace) Entries() []Entry {
	return t.entries
}

// Len returns the number of entries in the trace.
func (t *Trace) Len() int {
	return len(t.entries)
}

// FilterByLabel returns all entries whose label matches the given string.
func (t *Trace) FilterByLabel(label string) []Entry {
	var result []Entry
	for _, e := range t.entries {
		if e.Label == label {
			result = append(result, e)
		}
	}
	return result
}

// Summary returns a brief human-readable overview of the trace.
func (t *Trace) Summary() string {
	if len(t.entries) == 0 {
		return "trace: no entries recorded"
	}
	total := 0
	for _, e := range t.entries {
		total += len(e.Changes)
	}
	return fmt.Sprintf("trace: %d entries, %d total changes", len(t.entries), total)
}
