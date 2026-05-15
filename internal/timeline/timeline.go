// Package timeline provides ordered tracking of environment diff events
// across multiple snapshots, enabling historical analysis of how environment
// variables have evolved over time.
package timeline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envlens/internal/differ"
)

// Entry represents a single point-in-time diff recorded in the timeline.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Label     string           `json:"label"`
	Source    string           `json:"source"`
	Target    string           `json:"target"`
	Changes   []differ.Change  `json:"changes"`
}

// Timeline holds an ordered sequence of diff entries.
type Timeline struct {
	Entries []Entry `json:"entries"`
}

// New returns an empty Timeline.
func New() *Timeline {
	return &Timeline{Entries: []Entry{}}
}

// Add appends a new entry to the timeline.
func (t *Timeline) Add(label, source, target string, changes []differ.Change) {
	t.Entries = append(t.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Source:    source,
		Target:    target,
		Changes:   changes,
	})
}

// Len returns the number of entries in the timeline.
func (t *Timeline) Len() int {
	return len(t.Entries)
}

// Save writes the timeline to a JSON file at the given path.
func (t *Timeline) Save(path string) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("timeline: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("timeline: write failed: %w", err)
	}
	return nil
}

// Load reads a timeline from a JSON file at the given path.
func Load(path string) (*Timeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("timeline: read failed: %w", err)
	}
	var tl Timeline
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, fmt.Errorf("timeline: unmarshal failed: %w", err)
	}
	return &tl, nil
}

// KeyHistory returns all entries that contain a change for the given key.
func (t *Timeline) KeyHistory(key string) []Entry {
	var result []Entry
	for _, e := range t.Entries {
		for _, c := range e.Changes {
			if c.Key == key {
				result = append(result, e)
				break
			}
		}
	}
	return result
}
